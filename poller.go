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
	polledAsOf time.Time
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

	fmt.Println("polling:", search.URL)
	result, err := pc.cl.GetNewListings(ctx, search.URL, record.polledAsOf)
	if err != nil {
		fmt.Println("err getting listings from fn poll():", err)
	}

	// new cutoff date
	newCutoff := record.polledAsOf
	if len(result.Listings) > 0 {
		listingsToSave := []clListing{}
		var maxUnixDate int64

		for _, l := range result.Listings {
			var price int = 0
			if len(l.Price) > 0 {
				price, err = strconv.Atoi(l.Price[1:])
				if err != nil {
					fmt.Println("err converting from fn poll():", err)
				}
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

		fmt.Println("saving ", len(result.Listings), " new listings...")
		pc.db.saveListingMulti(search.ID, listingsToSave)

		newCutoff = time.Unix(maxUnixDate, 0)
		layout := "2006-01-02 15:04"
		newCutoff, err = time.Parse(layout, newCutoff.String())
		// there is a bug from GetNewListings that is returning
		// a date equal to the currentCutoff, until its fixes, this
		// should be a decent hack. Issue is opened on github
		newCutoff = newCutoff.Add(1 * time.Second)
		if err != nil {
			fmt.Println("err parsing cutoff time", err)
		}
	}

	record.polledAsOf = newCutoff

	pc.mu.Unlock()
	time.AfterFunc(time.Duration(60*time.Second), func() { pc.poll(ctx, search) })
}
