package main

import (
	"context"
	"fmt"
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
	record, has := pc.records[search.ID]
	if !has {
		// this is a new record, set up accordingly
		record = &pollingRecord{}
		pc.records[search.ID] = record
	}

	// current cutoff date
	currentCutoff := search.CreatedOn
	if record.polledAsOf.After(currentCutoff) {
		currentCutoff = record.polledAsOf
	}

	newRecords, err := pc.cl.GetNewListings(ctx, search.URL, currentCutoff)
	if err != nil {
		fmt.Println("err getting listings from fn poll():", err)
	}

	// new cutoff date
	newCutoff := record.polledAsOf
	if len(newRecords) > 0 {
		layout := "2006-01-02 15:04"
		newCutoff, err = time.Parse(layout, newRecords[0].Date)
		if err != nil {
			fmt.Println("err parsing cutoff time", err)
		}
	}

	record.polledAsOf = newCutoff
	record.listings = append(record.listings, newRecords...)

	time.AfterFunc(time.Duration(15*time.Second), func() { pc.poll(ctx, search) })
}
