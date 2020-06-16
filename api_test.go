package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	craigslist "github.com/tmunayyer/go-craigslist"
)

type mockCraigslistClient struct {
	location string
}

func (m *mockCraigslistClient) FormatURL(term string, o craigslist.Options) string {
	return ""
}

func (m *mockCraigslistClient) GetListings(ctx context.Context, url string) ([]craigslist.Listing, error) {
	listings := []craigslist.Listing{
		{
			DataPID: "1",
		},
		{
			DataPID: "2",
		},
		{
			DataPID: "3",
		},
	}

	return listings, nil
}

func TestMonitorURL(t *testing.T) {
	cl := mockCraigslistClient{
		location: "newyork",
	}

	mockHandlerEnv := handlerEnv{
		cl: &cl,
	}

	t.Run("post - should return data", func(t *testing.T) {
		// create a new request
		request, err := http.NewRequest(http.MethodPost, "/players/Pepper", nil)
		assert.NoError(t, err)
		// recorder to pass to spy on response
		response := httptest.NewRecorder()

		handleMonitorURL(response, request, &mockHandlerEnv)

		data, err := ioutil.ReadAll(response.Body)
		assert.NoError(t, err)

		expected := []string{"1", "2", "3"}
		var listings []struct {
			DataPID string
		}
		err = json.Unmarshal(data, &listings)
		assert.NoError(t, err)

		for i, listing := range listings {
			assert.Equal(t, expected[i], listing.DataPID)
		}
	})
}
