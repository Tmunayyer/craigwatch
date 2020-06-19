package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupDBTestCase(t *testing.T) (connection, func(t *testing.T), error) {
	t.Helper()

	c, err := newDBClient()
	if err != nil {
		return nil, func(t *testing.T) {}, err
	}

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

	args := craigslistQuery{
		Email: "TESTING@gmail.com",
		URL:   "www.TESTING.com",
	}

	record, err := c.saveURL(args)
	assert.NoError(t, err)

	assert.Equal(t, args.Email, record.Email)
	assert.Equal(t, args.URL, record.URL)
	assert.False(t, record.Confirmed)
	assert.Less(t, 0, record.ID)
}
