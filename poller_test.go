package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPoll(t *testing.T) {
	// these are defined in api_test.go
	mockCL := mockCraigslistClient{}
	mockDB := mockDBClient{}

	t.Run("should return listings on first flush call", func(t *testing.T) {
		mockPoller := newPollingService(&mockCL, &mockDB)

		s := clSearch{
			ID:  99,
			URL: "www.testing.com",
		}

		mockPoller.poll(context.Background(), s)

		listings, err := mockPoller.flush(s.ID)
		assert.NoError(t, err)
		assert.Greater(t, len(listings), 0)
	})
}

func TestFlush(t *testing.T) {
	// defined in api_test.go
	mockCL := mockCraigslistClient{}
	mockDB := mockDBClient{}

	t.Run("should NOT return listings on subsequent calls", func(t *testing.T) {
		mockPoller := newPollingService(&mockCL, &mockDB)

		s := clSearch{
			ID:  99,
			URL: "www.testing.com",
		}

		// put something into accumulator
		mockPoller.poll(context.Background(), s)

		listings, err := mockPoller.flush(s.ID)
		assert.NoError(t, err)
		assert.Greater(t, len(listings), 0)

		listings, err = mockPoller.flush(s.ID)
		assert.NoError(t, err)
		assert.Equal(t, 0, len(listings))
	})
}
