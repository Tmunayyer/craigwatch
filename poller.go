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
	polledAsOf      time.Time
	pollingInterval int // seconds, default 60
	emptyPollCount  int // count the time a poll comes back as nothing new to use as input to interval adjustments
}

type pollingClient struct {
	cl      craigslist.API
	db      connection
	mu      sync.Mutex // guard records
	records map[int]*pollingRecord
}

func newPollingService(cl craigslist.API, db connection) pollingService {
	pc := pollingClient{
		cl:      cl,
		db:      db,
		records: make(map[int]*pollingRecord),
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
	record, has := pc.records[search.ID]
	if !has {
		// this is a new record, set up accordingly
		record = &pollingRecord{}
		pc.records[search.ID] = record

		record.polledAsOf = time.Unix(int64(search.UnixCutoffDate), 0)
	}

	fmt.Println("polling:", search.URL) // intentional
	result, err := pc.cl.GetNewListings(ctx, search.URL, record.polledAsOf)
	if err != nil {
		fmt.Println("err getting listings from fn poll():", err)
	}

	// new cutoff date
	newCutoff := record.polledAsOf
	// polling interaval
	if len(result.Listings) > 0 {
		listingsToSave, maxUnixDate := pc.processNewListings(result)

		fmt.Println("saving ", len(result.Listings), " new listings...") // intentional
		pc.db.saveListingMulti(search.ID, listingsToSave)

		newCutoff = pc.calculateCutoff(maxUnixDate)

		// calculate polling interval
		activity, err := pc.db.getSearchActivity(search.ID)
		if err != nil {
			fmt.Println("err getting search activity in fn poll():", err)
		}

		record.pollingInterval = activity.InSeconds
		record.emptyPollCount = 1
	} else {
		if record.emptyPollCount < 10 {
			record.emptyPollCount++
		}
	}
	record.polledAsOf = newCutoff

	pc.mu.Unlock()
	interval := time.Duration(record.pollingInterval*record.emptyPollCount) * time.Second
	time.AfterFunc(interval, func() { pc.poll(ctx, search) })
}

func (pc *pollingClient) processNewListings(data *craigslist.Result) ([]clListing, int64) {
	listingsToSave := []clListing{}
	var maxUnixDate int64

	for _, l := range data.Listings {
		var price int = 0
		if len(l.Price) > 0 {
			num, err := strconv.Atoi(l.Price[1:])
			if err != nil {
				fmt.Println("err converting from fn poll():", err)
			}
			price = num
		}

		unixDate := newUnixDate(l.Date)
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

	return listingsToSave, maxUnixDate
}

func (pc *pollingClient) calculateCutoff(maxUnixDate int64) time.Time {
	newCutoff := time.Unix(maxUnixDate, 0).UTC()
	layout := "2006-01-02 15:04"
	newCutoff, err := time.Parse(layout, newCutoff.String()[:16])
	if err != nil {
		fmt.Println("err parsing cutoff time", err)
	}

	// there is a bug from GetNewListings that is returning
	// a date equal to the currentCutoff, until its fixes, this
	// should be a decent hack. Issue is opened on github
	newCutoff = newCutoff.Add(1 * time.Second)
	return newCutoff
}
