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
	poll(ctx context.Context, url string)
	flush() ([]craigslist.Listing, error)
}

type pollingClient struct {
	cl       craigslist.Client
	db       connection
	listings []craigslist.Listing
}

func newPollingService(cl craigslist.Client, db connection) pollingService {
	pc := pollingClient{
		cl: cl,
		db: db,
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
		go pc.poll(ctx, s.URL)
	}

	return nil
}

func (pc *pollingClient) shutdown() error {
	return nil
}

func (pc *pollingClient) flush() ([]craigslist.Listing, error) {
	newRecords := pc.listings
	pc.listings = []craigslist.Listing{}

	return newRecords, nil
}

func (pc *pollingClient) poll(ctx context.Context, url string) {
	records, err := pc.cl.GetListings(ctx, url)
	if err != nil {
		panic(err)
	}

	pc.listings = append(pc.listings, records...)

	time.AfterFunc(time.Duration(15*time.Second), func() { pc.poll(ctx, url) })
}
