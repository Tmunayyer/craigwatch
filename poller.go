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
	flush(ID int) ([]craigslist.Listing, error)
}

type pollingRecord struct {
	polledAsOf time.Time
	listings   []craigslist.Listing
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

	// err := pc.initiate(context.TODO())
	// if err != nil {
	// 	panic(err)
	// }

	return &pc
}

func (pc *pollingClient) initiate(ctx context.Context) error {
	searches, err := pc.db.getAllSearches()
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

func (pc *pollingClient) flush(ID int) ([]craigslist.Listing, error) {
	record, has := pc.records[ID]
	if !has {
		return []craigslist.Listing{}, fmt.Errorf("invalid ID provided")
	}

	newRecords := record.listings
	record.listings = []craigslist.Listing{}

	return newRecords, nil
}

func (pc *pollingClient) poll(ctx context.Context, search clSearch) {
	pc.mu.Lock()
	record, has := pc.records[search.ID]
	if !has {
		// this is a new record, set up accordingly
		record = &pollingRecord{}
		pc.records[search.ID] = record
	}

	// the first time, it should be time.Time nil value giving the entire first
	// page, the second time it should be updated
	currentCutoff := record.polledAsOf

	fmt.Println("polling:", search.URL)
	result, err := pc.cl.GetNewListings(ctx, search.URL, currentCutoff)
	if err != nil {
		fmt.Println("err getting listings from fn poll():", err)
	}

	// new cutoff date
	newCutoff := record.polledAsOf
	if len(result.Listings) > 0 {
		listingsToSave := []clListing{}
		for _, l := range result.Listings {
			p, err := strconv.Atoi(l.Price[1:])
			if err != nil {
				fmt.Println("err converting from fn poll():", err)
			}
			listingsToSave = append(listingsToSave, clListing{
				DataPID:      l.DataPID,
				DataRepostOf: l.DataRepostOf,
				UnixDate:     newUnixDate(l.Date),
				Title:        l.Title,
				Link:         l.Link,
				Price:        p,
				Hood:         l.Hood,
			})
		}
		pc.db.saveListings(search.ID, listingsToSave)

		layout := "2006-01-02 15:04"
		newCutoff, err = time.Parse(layout, result.Listings[0].Date)
		// there is a bug from GetNewListings that is returning
		// a date equal to the currentCutoff, until its fixes, this
		// should be a decent hack. Issue is opened on github
		newCutoff = newCutoff.Add(1 * time.Second)
		if err != nil {
			fmt.Println("err parsing cutoff time", err)
		}
	}

	record.polledAsOf = newCutoff
	record.listings = append(record.listings, result.Listings...)

	pc.mu.Unlock()
	time.AfterFunc(time.Duration(60*time.Second), func() { pc.poll(ctx, search) })
}
