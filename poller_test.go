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
	getNewListingsFnCount := 0
	mockCL := mockCraigslistClient{
		location: "newyork",
		getNewListingsFn: func(ctx context.Context, url string) (*craigslist.Result, error) {
			getNewListingsFnCount++

			fakeResult := craigslist.Result{
				Listings: []craigslist.Listing{},
			}

			if getNewListingsFnCount > 4 {
				fakeResult = craigslist.Result{
					Listings: []craigslist.Listing{
						{
							DataPID:      "1",
							DataRepostOf: "",
							Date:         "2020-01-02 16:04",
							Title:        "Bananas, organic, not eaten",
							Link:         "www.craigslist.com/bananapost",
							Price:        "$100",
							Hood:         "newyork",
						},
						{
							DataPID:      "2",
							DataRepostOf: "",
							Date:         "2020-01-02 13:04",
							Title:        "The best crackers ever",
							Link:         "www.craigslist.com/crackers",
							Price:        "$20",
							Hood:         "newyork",
						},
					},
				}
			}

			return &fakeResult, nil
		},
	}
	mockDB := mockDBClient{}
	mockPoller := &pollingClient{
		cl:      &mockCL,
		db:      &mockDB,
		records: make(map[int]*pollingRecord),

		defaultPollingInterval:       5000, // make this a non-factor for this test, fake listings hsould have 90 min interval
		maxPollingIntervalMultiplier: 3,
	}

	s := clSearch{
		ID:             99,
		URL:            "www.testing.com",
		UnixCutoffDate: 0,
	}

	mockPoller.poll(context.Background(), s)

	r, has := mockPoller.records[s.ID]
	assert.True(t, has)

	// defaults should be set
	assert.Equal(t, 1, r.emptyPollCount)
	mockPoller.poll(context.Background(), s)
	assert.Equal(t, 2, r.emptyPollCount)
	mockPoller.poll(context.Background(), s)
	assert.Equal(t, 3, r.emptyPollCount)
	// should remain at the max interval
	mockPoller.poll(context.Background(), s)
	assert.Equal(t, 3, r.emptyPollCount)

	// set the interval in seconds being set
	mockPoller.poll(context.Background(), s)
	assert.Equal(t, 1, r.emptyPollCount)
	assert.Equal(t, 5400, r.pollingInterval)
}
