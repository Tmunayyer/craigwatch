package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	craigslist "github.com/tmunayyer/go-craigslist"
)

const (
	badCraigslistURL = "www.badurl.com"
)

var fakeListings = []craigslist.Listing{
	{
		DataPID:      "1",
		DataRepostOf: "",
		Date:         "06/19/2020",
		Title:        "Bananas, organic, not eaten",
		Link:         "www.craigslist.com/bananapost",
		Price:        "$100",
		Hood:         "newyork",
	},
	{
		DataPID:      "2",
		DataRepostOf: "",
		Date:         "06/15/2020",
		Title:        "The best crackers ever",
		Link:         "www.craigslist.com/crackers",
		Price:        "$20",
		Hood:         "newyork",
	},
}

type mockCraigslistClient struct {
	location string
}

func (m *mockCraigslistClient) FormatURL(term string, o craigslist.Options) string {
	return ""
}

func (m *mockCraigslistClient) GetListings(ctx context.Context, url string) ([]craigslist.Listing, error) {
	if url == badCraigslistURL {
		return nil, fmt.Errorf("invalid url: %v", url)
	}

	return fakeListings, nil
}

type mockDBClient struct {
}

func (m *mockDBClient) connect() error {
	return nil
}
func (m *mockDBClient) shutdown() error {
	return nil
}
func (m *mockDBClient) testConnection() error {
	return nil
}
func (m *mockDBClient) applySchema() error {
	return nil
}
func (m *mockDBClient) saveURL(data url) (url, error) {
	return url{}, nil
}

func TestMonitorURL(t *testing.T) {
	mockCL := mockCraigslistClient{
		location: "newyork",
	}
	mockDB := mockDBClient{}

	api := newAPIService(&mockCL, &mockDB)

	t.Run("post - no body passed", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/", http.NoBody)
		assert.NoError(t, err)
		res := httptest.NewRecorder()

		api.handleMonitorURL(res, req)

		message, err := ioutil.ReadAll(res.Body)

		assert.Equal(t, http.StatusBadRequest, res.Code)
		assert.Equal(t, "could not decode user payload: EOF\n", string(message))
	})

	t.Run("post - invalid url", func(t *testing.T) {
		// make a body
		type body struct {
			URL string
		}

		b := body{URL: badCraigslistURL}
		data, err := json.Marshal(b)
		assert.NoError(t, err)
		reader := bytes.NewReader(data)

		req, err := http.NewRequest(http.MethodPost, "/", reader)
		assert.NoError(t, err)
		res := httptest.NewRecorder()

		api.handleMonitorURL(res, req)

		message, err := ioutil.ReadAll(res.Body)

		assert.Equal(t, http.StatusBadRequest, res.Code)
		assert.Equal(t, "url provided is not a compatible with craigslist\n", string(message))
	})

	t.Run("post - recieves data", func(t *testing.T) {
		// make a body
		type body struct {
			URL string
		}

		b := body{URL: "www.anything.com"}
		data, err := json.Marshal(b)
		assert.NoError(t, err)
		reader := bytes.NewReader(data)

		req, err := http.NewRequest(http.MethodPost, "/", reader)
		assert.NoError(t, err)
		res := httptest.NewRecorder()

		// call handelr
		api.handleMonitorURL(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		// just check for the pids
		responseData := struct {
			Searched string
			Results  []craigslist.Listing
		}{}

		// put body into a slice
		responseBytes, err := ioutil.ReadAll(res.Body)
		assert.NoError(t, err)
		err = json.Unmarshal(responseBytes, &responseData)
		assert.NoError(t, err)

		for i, listing := range responseData.Results {
			assert.Equal(t, fakeListings[i].DataPID, listing.DataPID)
		}

	})
}
