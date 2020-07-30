package main

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var testSearch = clSearch{
	Name:     "test search 99",
	URL:      "www.TESTING.com",
	Timezone: "testtimezone",
}

var testListings = []clListing{
	{
		DataPID:      "123456",
		DataRepostOf: "987654",
		UnixDate:     newUnixDate("2020-01-02 12:00:00", nil),
		Title:        "testListingNumeroUno",
		Link:         "www.testing.com",
		Price:        106,
		Hood:         "dontbeamenacetosouthcentral",
	},
	{
		DataPID:      "654321",
		DataRepostOf: "123498",
		UnixDate:     newUnixDate("2020-01-01 12:00:00", nil),
		Title:        "testListingNumeroDOS",
		Link:         "www.testing.com",
		Price:        999,
		Hood:         "gattaca",
	},
}

func setupDBTestCase(t *testing.T) (connection, func(t *testing.T), error) {
	t.Helper()

	c := newDBClient()

	teardown := func(t *testing.T) {
		t.Helper()
		c.shutdown()
	}

	return c, teardown, nil
}

func TestSaveSearch(t *testing.T) {
	c, teardown, err := setupDBTestCase(t)
	assert.NoError(t, err)
	defer teardown(t)

	args := testSearch

	record, err := c.saveSearch(args)
	assert.NoError(t, err)

	assert.Equal(t, args.Name, record.Name)
	assert.Equal(t, args.URL, record.URL)
	assert.Equal(t, args.Timezone, record.Timezone)
	assert.False(t, record.Confirmed)
	assert.Less(t, 0, record.ID)

	// remove the record
	err = c.deleteSearch(record.ID)
	assert.NoError(t, err)
}

func TestGetSearch(t *testing.T) {
	c, teardown, err := setupDBTestCase(t)
	assert.NoError(t, err)
	defer teardown(t)

	args := testSearch

	saved, err := c.saveSearch(args)
	assert.NoError(t, err)

	record, err := c.getSearch(saved.ID)
	assert.NoError(t, err)

	assert.Equal(t, args.Name, record.Name)
	assert.Equal(t, args.URL, record.URL)

	// delete the saved records
	err = c.deleteSearch(saved.ID)
	assert.NoError(t, err)
}

func TestGetSearchMulti(t *testing.T) {
	c, teardown, err := setupDBTestCase(t)
	assert.NoError(t, err)
	defer teardown(t)

	args := testSearch

	saved, err := c.saveSearch(args)
	assert.NoError(t, err)

	records, err := c.getSearchMulti()
	assert.NoError(t, err)

	doesExist := false
	for _, r := range records {
		if r.URL == args.URL {
			doesExist = true
			break
		}
	}
	assert.True(t, doesExist)

	// delete the saved records
	err = c.deleteSearch(saved.ID)
	assert.NoError(t, err)
}

func TestDeleteSearch(t *testing.T) {
	c, teardown, err := setupDBTestCase(t)
	assert.NoError(t, err)
	defer teardown(t)

	args := testSearch

	saved, err := c.saveSearch(args)
	assert.NoError(t, err)

	// delete the record
	err = c.deleteSearch(saved.ID)
	assert.NoError(t, err)

	// make sure its gone
	records, err := c.getSearchMulti()
	assert.NoError(t, err)
	for _, record := range records {
		if record.ID == saved.ID {
			t.Fail()
		}
	}
}

func TestSaveListingMulti(t *testing.T) {
	c, teardown, err := setupDBTestCase(t)
	assert.NoError(t, err)
	defer teardown(t)

	search, err := c.saveSearch(testSearch)
	assert.NoError(t, err)

	err = c.saveListingMulti(search.ID, testListings)
	assert.NoError(t, err)

	savedListings, err := c.getListingMulti(search.ID)
	assert.NoError(t, err)

	// tests
	assert.Len(t, savedListings, 2)
	// should be ordered by date
	assert.Equal(t, savedListings[0].DataPID, testListings[0].DataPID)
	assert.Equal(t, savedListings[1].DataPID, testListings[1].DataPID)

	err = c.deleteListingMulti(search.ID)
	assert.NoError(t, err)
	err = c.deleteSearch(search.ID)
	assert.NoError(t, err)
}

func TestGetListingMulti(t *testing.T) {
	c, teardown, err := setupDBTestCase(t)
	assert.NoError(t, err)
	defer teardown(t)

	t.Run("get all listings", func(t *testing.T) {
		search, err := c.saveSearch(testSearch)
		assert.NoError(t, err)

		err = c.saveListingMulti(search.ID, testListings)
		assert.NoError(t, err)

		savedListings, err := c.getListingMulti(search.ID)
		assert.NoError(t, err)

		// tests
		assert.Len(t, savedListings, 2)
		// should be ordered by date
		assert.Equal(t, savedListings[0].DataPID, testListings[0].DataPID)
		assert.Equal(t, savedListings[1].DataPID, testListings[1].DataPID)

		err = c.deleteListingMulti(search.ID)
		assert.NoError(t, err)
		err = c.deleteSearch(search.ID)
		assert.NoError(t, err)
	})

	t.Run("get all listings after specfic datetime", func(t *testing.T) {
		search, err := c.saveSearch(testSearch)
		assert.NoError(t, err)

		err = c.saveListingMulti(search.ID, testListings)
		assert.NoError(t, err)

		unixTime := newUnixDate("2020-01-01 12:00:00", nil)

		savedListings, err := c.getListingMultiAfter(search.ID, unixTime)
		assert.NoError(t, err)

		// tests
		assert.Len(t, savedListings, 1)
		// should be ordered by date
		assert.Equal(t, savedListings[0].DataPID, testListings[0].DataPID)

		err = c.deleteListingMulti(search.ID)
		assert.NoError(t, err)
		err = c.deleteSearch(search.ID)
		assert.NoError(t, err)
	})
}

func TestListingActivity(t *testing.T) {
	c, teardown, err := setupDBTestCase(t)
	assert.NoError(t, err)
	defer teardown(t)

	t.Run("should return metrics", func(t *testing.T) {
		search, err := c.saveSearch(testSearch)
		assert.NoError(t, err)

		err = c.saveListingMulti(search.ID, testListings)
		assert.NoError(t, err)

		activity, err := c.getSearchActivity(search.ID)
		assert.NoError(t, err)

		assert.Equal(t, float32(24), activity.InHours)

		err = c.deleteListingMulti(search.ID)
		assert.NoError(t, err)
		err = c.deleteSearch(search.ID)
		assert.NoError(t, err)
	})
}

func TestListingActivityByHour(t *testing.T) {
	c, teardown, err := setupDBTestCase(t)
	assert.NoError(t, err)
	defer teardown(t)

	t.Run("should return hourly data", func(t *testing.T) {
		search, err := c.saveSearch(testSearch)
		assert.NoError(t, err)

		// need listings within 24 hours
		recentListings := []clListing{}

		date := time.Now()
		for i := 0; i < 5; i++ {
			// 30 minutes ago
			recentListingDate := date.Add(-30 * time.Minute).UTC()
			recentListings = append(recentListings, clListing{
				DataPID:      "123456" + strconv.Itoa(i),
				DataRepostOf: "987654",
				UnixDate:     recentListingDate.Unix(),
				Title:        "testListingNumeroUno",
				Link:         "www.testing.com",
				Price:        106,
				Hood:         "dontbeamenacetosouthcentral",
			})
		}

		for i := 0; i < 5; i++ {
			// 30 minutes ago
			recentListingDate := date.Add(-90 * time.Minute).UTC()
			recentListings = append(recentListings, clListing{
				DataPID:      "654321" + strconv.Itoa(i),
				DataRepostOf: "987654",
				UnixDate:     recentListingDate.Unix(),
				Title:        "testListingNumeroUno",
				Link:         "www.testing.com",
				Price:        106,
				Hood:         "dontbeamenacetosouthcentral",
			})
		}

		err = c.saveListingMulti(search.ID, recentListings)
		assert.NoError(t, err)

		activity, err := c.getSearchActivityByHour(search.ID)
		assert.NoError(t, err)

		assert.Equal(t, activity[0].Count, 5)
		assert.Equal(t, activity[1].Count, 5)

		err = c.deleteListingMulti(search.ID)
		assert.NoError(t, err)
		err = c.deleteSearch(search.ID)
		assert.NoError(t, err)
	})
}

func TestPriceDistribution(t *testing.T) {
	c, teardown, err := setupDBTestCase(t)
	assert.NoError(t, err)
	defer teardown(t)

	t.Run("should return distribution data", func(t *testing.T) {
		search, err := c.saveSearch(testSearch)
		assert.NoError(t, err)

		// need listings within 24 hours
		recentListings := []clListing{}

		date := time.Now()
		for i := 0; i < 5; i++ {
			// 30 minutes ago
			recentListingDate := date.Add(-30 * time.Minute).UTC()
			recentListings = append(recentListings, clListing{
				DataPID:      "123456" + strconv.Itoa(i),
				DataRepostOf: "987654",
				UnixDate:     recentListingDate.Unix(),
				Title:        "testListingNumeroUno",
				Link:         "www.testing.com",
				Price:        57 + 3*i, // total = 315
				Hood:         "dontbeamenacetosouthcentral",
			})
		}

		for i := 0; i < 5; i++ {
			// 30 minutes ago
			recentListingDate := date.Add(-90 * time.Minute).UTC()
			recentListings = append(recentListings, clListing{
				DataPID:      "654321" + strconv.Itoa(i),
				DataRepostOf: "987654",
				UnixDate:     recentListingDate.Unix(),
				Title:        "testListingNumeroUno",
				Link:         "www.testing.com",
				Price:        120 + 9*i, // total = 690
				Hood:         "dontbeamenacetosouthcentral",
			})
		}

		err = c.saveListingMulti(search.ID, recentListings)
		assert.NoError(t, err)

		distribution, err := c.getPriceDistribution(search.ID)
		assert.NoError(t, err)

		// 315
		assert.Equal(t, distribution.AveragePrice, 1005/10+1)

		err = c.deleteListingMulti(search.ID)
		assert.NoError(t, err)
		err = c.deleteSearch(search.ID)
		assert.NoError(t, err)
	})
}
