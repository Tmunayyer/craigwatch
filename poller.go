package main

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	craigslist "github.com/tmunayyer/go-craigslist"
)

type pollingService interface {
	initiate(context.Context) error
	shutdown() error
	poll(ctx context.Context, search clSearch)
}

type pollingRecord struct {
	polledAsOf      int64
	pollingInterval int // seconds, default 60
	emptyPollCount  int // count the time a poll comes back as nothing new to use as input to interval adjustments
	timezone        *time.Location
}

type pollingClient struct {
	cl      craigslist.API
	db      connection
	mu      sync.Mutex // guard records
	records map[int]*pollingRecord

	defaultPollingInterval       int
	maxPollingIntervalMultiplier int
}

func newPollingService(cl craigslist.API, db connection) pollingService {
	pc := pollingClient{
		cl:      cl,
		db:      db,
		records: make(map[int]*pollingRecord),

		defaultPollingInterval:       60, // seconds
		maxPollingIntervalMultiplier: 10,
	}

	err := pc.initiate(context.TODO())
	if err != nil {
		panic(err)
	}

	return &pc
}

func (pc *pollingClient) initiate(ctx context.Context) error {
	searches, err := pc.db.getSearchMulti()
	if err != nil {
		return err
	}

	fmt.Println("monitoring", len(searches), "searches...")
	for _, s := range searches {
		go pc.poll(ctx, s)
	}

	return nil
}

func (pc *pollingClient) shutdown() error {
	return nil
}

func (pc *pollingClient) poll(ctx context.Context, search clSearch) {
	pc.mu.Lock()
	// ===================
	// INITALIZATION
	// ===================
	record, has := pc.records[search.ID]
	if !has {
		// this is a new record, set up accordingly
		record = &pollingRecord{}
		pc.records[search.ID] = record

		loc, err := time.LoadLocation(search.Timezone)
		if err != nil {
			fmt.Printf("error loading location: %+v", err)
		}
		record.timezone = loc
		// 1. on initialization of a new poll, this is 0 so there is no need to apply a timezone to it,
		//    0 is going to tell the craigslist library to get all possible searches
		// 2. on initalization of an existing poll (found in DB) the unix cutoff will already
		//    have the timezone applied
		record.polledAsOf = search.UnixCutoffDate + 1
		record.pollingInterval = pc.defaultPollingInterval
		// set it to 0, if there are listings, its set to 1, if there are none 1 gets added
		record.emptyPollCount = 0
	}

	// ===================
	// GATHERING DATA
	// ===================
	fmt.Println("polling:", search.ID, search.URL) // intentional
	rawListings := []craigslist.Listing{}
	cutoff := unixToLocal(record.polledAsOf, record.timezone)
	result, err := pc.cl.GetNewListings(ctx, search.URL, cutoff)
	if err != nil {
		fmt.Println("err getting listings from fn poll():", err)
	}
	rawListings = append(rawListings, result.Listings...)

	for !result.Done {
		result, err = result.Next(ctx, cutoff)
		if err != nil {
			fmt.Println("err getting listings from fn poll():", err)
		}
		rawListings = append(rawListings, result.Listings...)
	}

	// ===================
	// PROCESSING DATA
	// ===================
	// new cutoff date
	newCutoff := record.polledAsOf
	// polling interaval
	if len(rawListings) > 0 {
		// process new listings also deduplicates data
		listingsToSave, maxUnixDate := pc.processNewListings(rawListings, record.timezone)

		fmt.Println("-- id", search.ID, "saving ", len(listingsToSave), " new listings...") // intentional
		err := pc.db.saveListingMulti(search.ID, listingsToSave)
		if err != nil {
			fmt.Println("err saving listings:", err)
		}
		// add 1 second to the max to avoid duplication, craigslist lib is inclusive
		newCutoff = maxUnixDate + 1

		// calculate polling interval
		activity, err := pc.db.getSearchActivity(search.ID)
		if err != nil {
			fmt.Println("err getting search activity in fn poll():", err)
		}

		// dont poll any faster than one a min
		if activity.InSeconds > pc.defaultPollingInterval {
			record.pollingInterval = activity.InSeconds
		} else {
			record.pollingInterval = pc.defaultPollingInterval
		}
		record.emptyPollCount = 1
	} else {
		if record.emptyPollCount < pc.maxPollingIntervalMultiplier {
			record.emptyPollCount++
		}
	}
	record.polledAsOf = newCutoff

	pc.mu.Unlock()
	interval := time.Duration(record.pollingInterval*record.emptyPollCount) * time.Second
	fmt.Println("-- id", search.ID, "polling again in", interval, "seconds") // intentional
	time.AfterFunc(interval, func() { pc.poll(ctx, search) })
}

func (pc *pollingClient) processNewListings(data []craigslist.Listing, tz *time.Location) ([]clListing, int64) {
	listingsToSave := []clListing{}
	var maxUnixDate int64

	for _, l := range data {
		var price int = 0
		if len(l.Price) > 0 {
			num, err := strconv.Atoi(l.Price[1:])
			if err != nil {
				fmt.Println("err converting from fn poll():", err)
			}
			price = num
		}

		unixDate := newUnixDate(l.Date, tz)
		if unixDate > maxUnixDate {
			maxUnixDate = unixDate
		}

		listingsToSave = append(listingsToSave, clListing{
			DataPID:      l.DataPID,
			DataRepostOf: l.DataRepostOf,
			UnixDate:     unixDate,
			Title:        l.Title,
			Link:         l.Link,
			Price:        price,
			Hood:         l.Hood,
		})
	}

	listingsToSave = pc.deduplicateData(listingsToSave)

	return listingsToSave, maxUnixDate
}

func (pc *pollingClient) deduplicateData(data []clListing) []clListing {
	// DataPID = clListing{}
	knownValues := make(map[string]clListing)

	for _, l := range data {
		record, has := knownValues[l.DataPID]
		if !has {
			knownValues[l.DataPID] = l
			continue
		}

		// use most recent listing
		if l.UnixDate > record.UnixDate {
			knownValues[l.DataPID] = l
		}
	}

	output := []clListing{}
	for _, v := range knownValues {
		output = append(output, v)
	}

	return output
}
