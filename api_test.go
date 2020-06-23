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
func (m *mockDBClient) saveSearch(data clSearch) (clSearch, error) {
	return clSearch{ID: 1, Email: data.Email, URL: data.URL, Confirmed: false}, nil
}
func (m *mockDBClient) getAllSearches() ([]clSearch, error) {
	return []clSearch{}, nil
}
func (m *mockDBClient) deleteSearch(id int) error {
	return nil
}

type mockPollingService struct {
	listings []craigslist.Listing
}

func (m *mockPollingService) initiate(context.Context) error {
	return nil
}
func (m *mockPollingService) shutdown() error {
	return nil
}
func (m *mockPollingService) poll(ctx context.Context, url string) {
	// note: copy requires destination to have a predefined length
	listings := make([]craigslist.Listing, len(fakeListings))
	copy(listings, fakeListings)
	m.listings = listings
}
func (m *mockPollingService) flush() ([]craigslist.Listing, error) {
	return m.listings, nil
}

func setupTestAPI(t *testing.T) *apiService {
	t.Helper()
	mockCL := mockCraigslistClient{
		location: "newyork",
	}
	mockDB := mockDBClient{}
	mockPS := mockPollingService{}

	return newAPIService(&mockCL, &mockDB, &mockPS)
}

func TestMonitorURL(t *testing.T) {
	api := setupTestAPI(t)

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
			ID        int
			Email     string
			URL       string
			Confirmed bool
		}

		b := body{Email: "testing@gmail.com", URL: "www.anything.com"}
		data, err := json.Marshal(b)
		assert.NoError(t, err)
		reader := bytes.NewReader(data)

		req, err := http.NewRequest(http.MethodPost, "/", reader)
		assert.NoError(t, err)
		res := httptest.NewRecorder()

		// call handelr
		api.handleMonitorURL(res, req)

		resBody := body{}
		readBodyInto(t, res.Body, &resBody)

		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, b.Email, resBody.Email)
		assert.Equal(t, b.URL, resBody.URL)
		assert.Equal(t, false, resBody.Confirmed)
	})
}

func TestHandleNewListings(t *testing.T) {
	api := setupTestAPI(t)

	type body struct {
		HasNewListings bool
		Listings       []craigslist.Listing
	}

	t.Run("get - should NOT return new listings", func(t *testing.T) {
		// this test case basically rides off the fact that initializing the
		// mock clients will retrieve no new listings

		req, err := http.NewRequest(http.MethodPost, "/", nil)
		assert.NoError(t, err)
		res := httptest.NewRecorder()

		api.handleNewListings(res, req)

		resBody := body{}
		readBodyInto(t, res.Body, &resBody)

		assert.False(t, resBody.HasNewListings)
	})

	t.Run("get - should return new listings", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/", nil)
		assert.NoError(t, err)
		res := httptest.NewRecorder()

		// before making the request, call poll to add some listings
		api.ps.poll(context.Background(), "anything.com")

		api.handleNewListings(res, req)

		resBody := body{}
		readBodyInto(t, res.Body, &resBody)

		assert.True(t, resBody.HasNewListings)
	})
}

// NOTE: before debugging here, make sure destination field are public
func readBodyInto(t *testing.T, b *bytes.Buffer, destination interface{}) {
	t.Helper()

	bodyBytes, err := ioutil.ReadAll(b)
	assert.NoError(t, err)
	err = json.Unmarshal(bodyBytes, destination)
	assert.NoError(t, err)
}
