package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: do I really need to test this asynchronously or is it enough to
// make it synchronous and assume it will work?

func TestPoll(t *testing.T) {
	// these are defined in api_test.go
	mockCL := mockCraigslistClient{}
	mockDB := mockDBClient{}

	mockPoller := newPollingService(&mockCL, &mockDB)

	mockPoller.poll(context.Background(), "www.testing.com")

	listings, err := mockPoller.flush()
	assert.NoError(t, err)

	fmt.Printf("mockpolling listings: %+v", listings)

}
