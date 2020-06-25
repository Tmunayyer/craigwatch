package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testSearch = clSearch{
	Name: "test search 99",
	URL:  "www.TESTING.com",
}

var testListings = []clListing{
	{
		DataPID:      "123456",
		DataRepostOf: "",
		Date:         newDate("2020-06-01 12:00"),
		Title:        "testListingNumeroUno",
		Link:         "www.testing.com",
		Price:        106,
		Hood:         "dontbeamenacetosouthcentral",
	},
	{
		DataPID:      "654321",
		DataRepostOf: "",
		Date:         newDate("2020-05-01 12:00"),
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

func TestSaveURL(t *testing.T) {
	c, teardown, err := setupDBTestCase(t)
	assert.NoError(t, err)
	defer teardown(t)

	args := clSearch{
		Name: "test search 99",
		URL:  "www.TESTING.com",
	}

	record, err := c.saveSearch(args)
	assert.NoError(t, err)

	assert.Equal(t, args.Name, record.Name)
	assert.Equal(t, args.URL, record.URL)
	assert.False(t, record.Confirmed)
	assert.Less(t, 0, record.ID)

	// remove the record
	err = c.deleteSearch(record.ID)
	assert.NoError(t, err)
}

func TestGetAllURL(t *testing.T) {
	c, teardown, err := setupDBTestCase(t)
	assert.NoError(t, err)
	defer teardown(t)

	args := clSearch{
		Name: "test search 0100",
		URL:  "www.TESTING.com",
	}

	saved, err := c.saveSearch(args)
	assert.NoError(t, err)

	records, err := c.getAllSearches()
	assert.NoError(t, err)
	assert.Equal(t, args.URL, records[0].URL)

	// delete the saved records
	err = c.deleteSearch(saved.ID)
	assert.NoError(t, err)
}

func TestDeleteSearch(t *testing.T) {
	c, teardown, err := setupDBTestCase(t)
	assert.NoError(t, err)
	defer teardown(t)

	args := clSearch{
		Name: "testing 123",
		URL:  "www.TESTING.com",
	}

	saved, err := c.saveSearch(args)
	assert.NoError(t, err)

	// delete the record
	err = c.deleteSearch(saved.ID)
	assert.NoError(t, err)

	// make sure its gone
	records, err := c.getAllSearches()
	assert.NoError(t, err)
	for _, record := range records {
		if record.ID == saved.ID {
			t.Fail()
		}
	}
}

func TestSaveListing(t *testing.T) {
	c, teardown, err := setupDBTestCase(t)
	assert.NoError(t, err)
	defer teardown(t)

	search, err := c.saveSearch(testSearch)
	assert.NoError(t, err)

	err = c.saveListings(search.ID, testListings)
	assert.NoError(t, err)

	savedListings, err := c.getListings(search.ID)
	assert.NoError(t, err)

	// tests
	assert.Len(t, savedListings, 2)
	// should be ordered by date
	assert.Equal(t, savedListings[0].DataPID, testListings[0].DataPID)
	assert.Equal(t, savedListings[1].DataPID, testListings[1].DataPID)

	err = c.deleteListings(search.ID)
	assert.NoError(t, err)
	err = c.deleteListings(search.ID)
	assert.NoError(t, err)
}
