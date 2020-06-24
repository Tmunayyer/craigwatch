package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
