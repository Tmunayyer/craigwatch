package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: do I really need to test this asynchronously or is it enough to
// make it synchronous and assume it will work?

func TestPoll(t *testing.T) {
	// these are defined in api_test.go
	mockCL := mockCraigslistClient{}
	mockDB := mockDBClient{}

	t.Run("should return listings on first flush call", func(t *testing.T) {
		mockPoller := newPollingService(&mockCL, &mockDB)

		mockPoller.poll(context.Background(), "www.testing.com")

		listings, err := mockPoller.flush()
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

		// put something into accumulator
		mockPoller.poll(context.Background(), "www.testing.com")

		listings, err := mockPoller.flush()
		assert.NoError(t, err)
		assert.Greater(t, len(listings), 0)

		listings, err = mockPoller.flush()
		assert.NoError(t, err)
		assert.Equal(t, 0, len(listings))
	})
}
