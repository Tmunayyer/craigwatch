package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	craigslist "github.com/tmunayyer/go-craigslist"
)

func TestPollingCutoffCalculation(t *testing.T) {
	// at this point, I know
	// - the underlying library returns listings
	// - the datastore properly saves and retrieves listings

	// the only thing to test here is the handling of the cutoff date

	// these are defined in api_test.go
	mockCL := mockCraigslistClient{
		location: "newyork",
		getNewListingsFn: func(ctx context.Context, url string) (*craigslist.Result, error) {
			if url == badCraigslistURL {
				return &craigslist.Result{}, fmt.Errorf("invalid url: %v", url)
			}

			fakeResult := craigslist.Result{
				Listings: fakeListings,
			}

			return &fakeResult, nil
		},
	}
	mockDB := mockDBClient{}
	mockPoller := &pollingClient{
		cl:      &mockCL,
		db:      &mockDB,
		records: make(map[int]*pollingRecord),
	}

	s := clSearch{
		ID:             99,
		URL:            "www.testing.com",
		UnixCutoffDate: 0,
	}

	mockPoller.poll(context.Background(), s)

	r, has := mockPoller.records[s.ID]
	assert.True(t, has)
	// - fake listings most recent listing date: 2020-01-02 16:04 and should be the polled as of
	//   after poll is run
	// - add 1 to account for adding 1 second in poll
	assert.Equal(t, newUnixDate(fakeListings[0].Date)+1, r.polledAsOf.Unix())
}

func TestPollingInterval(t *testing.T) {
	// these are defined in api_test.go
	mockCL := mockCraigslistClient{
		location: "newyork",
		getNewListingsFn: func(ctx context.Context, url string) (*craigslist.Result, error) {
			if url == badCraigslistURL {
				return &craigslist.Result{}, fmt.Errorf("invalid url: %v", url)
			}

			fakeResult := craigslist.Result{
				Listings: fakeListings,
			}

			return &fakeResult, nil
		},
	}
	mockDB := mockDBClient{}
	mockPoller := &pollingClient{
		cl:      &mockCL,
		db:      &mockDB,
		records: make(map[int]*pollingRecord),
	}

	s := clSearch{
		ID:             99,
		URL:            "www.testing.com",
		UnixCutoffDate: 0,
	}

	mockPoller.poll(context.Background(), s)

	_, has := mockPoller.records[s.ID]
	assert.True(t, has)

	// defaults should be set

}
