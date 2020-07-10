package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

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
		Date:         "2020-01-02 16:04",
		Title:        "Bananas, organic, not eaten",
		Link:         "www.craigslist.com/bananapost",
		Price:        "$100",
		Hood:         "newyork",
	},
	{
		DataPID:      "2",
		DataRepostOf: "",
		Date:         "2020-01-02 13:04",
		Title:        "The best crackers ever",
		Link:         "www.craigslist.com/crackers",
		Price:        "$20",
		Hood:         "newyork",
	},
}

var fakeSearch = clSearch{
	ID:        99,
	Name:      "bladerunner",
	URL:       "www.bladerunner.com",
	Confirmed: false,
	Interval:  0,
	CreatedOn: time.Time{},
}

type mockCraigslistClient struct {
	location string
}

func (m *mockCraigslistClient) FormatURL(term string, o craigslist.Options) string {
	return ""
}

func (m *mockCraigslistClient) GetListings(ctx context.Context, url string) (*craigslist.Result, error) {
	if url == badCraigslistURL {
		return &craigslist.Result{}, fmt.Errorf("invalid url: %v", url)
	}

	fakeResult := craigslist.Result{
		Listings: fakeListings,
	}

	return &fakeResult, nil
}

func (m *mockCraigslistClient) GetNewListings(ctx context.Context, url string, date time.Time) (*craigslist.Result, error) {
	if url == badCraigslistURL {
		return &craigslist.Result{}, fmt.Errorf("invalid url: %v", url)
	}

	fakeResult := craigslist.Result{
		Listings: fakeListings,
	}

	return &fakeResult, nil
}

type mockDBClient struct {
	saveSearchCallCount    int
	saveSearchURLs         map[string]clSearch
	saveListingsCallCount  int
	saveListingsCalledWith []clListing
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
	if m.saveSearchURLs == nil {
		m.saveSearchURLs = make(map[string]clSearch)
	}

	_, has := m.saveSearchURLs[data.URL]
	if has {
		return data, errors.New(duplicateURLErrorMessage)
	}

	newRecord := clSearch{ID: 1, Name: data.Name, URL: data.URL, Confirmed: false}
	m.saveSearchURLs[data.URL] = newRecord
	return newRecord, nil
}
func (m *mockDBClient) getSearch(searchID int) (clSearch, error) {
	if searchID == 99 {
		return clSearch{
			ID:   1,
			Name: "Test seach 1",
			URL:  "www.testing.com",
		}, nil
	}
	return clSearch{}, nil
}
func (m *mockDBClient) getSearchMulti() ([]clSearch, error) {
	return []clSearch{
		{ID: 1, Name: "Test seach 1", URL: "www.testing.com", Confirmed: false},
		{ID: 1, Name: "Test search 2", URL: "www.bigpotatotest.com", Confirmed: false},
	}, nil
}
func (m *mockDBClient) deleteSearch(id int) error {
	return nil
}

func (m *mockDBClient) saveListingMulti(monitorID int, listings []clListing) error {
	m.saveListingsCallCount++
	m.saveListingsCalledWith = listings
	return nil
}
func (m *mockDBClient) deleteListingMulti(monitorID int) error {
	return nil
}
func (m *mockDBClient) getListingMulti(id int) ([]clListing, error) {
	return []clListing{}, nil
}
func (m *mockDBClient) getListingMultiAfter(id int, date int64) ([]clListing, error) {
	output := []clListing{}
	for _, l := range fakeListings {
		p, err := strconv.Atoi(l.Price[1:])
		if err != nil {
			fmt.Println("err converting from fn poll():", err)
		}
		listing := clListing{
			ID:           123456,
			SearchID:     fakeSearch.ID,
			DataPID:      l.DataPID,
			DataRepostOf: l.DataRepostOf,
			UnixDate:     newUnixDate(l.Date),
			Title:        l.Title,
			Link:         l.Link,
			Price:        p,
			Hood:         l.Hood,
		}

		output = append(output, listing)
	}
	return output, nil
}

func (m *mockDBClient) getSearchActivity(searchID int) (searchActivity, error) {
	return searchActivity{}, nil
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
func (m *mockPollingService) poll(ctx context.Context, search clSearch) {
	// note: copy requires destination to have a predefined length
	listings := make([]craigslist.Listing, len(fakeListings))
	copy(listings, fakeListings)
	m.listings = listings
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

func TestHandleListing(t *testing.T) {
	api := setupTestAPI(t)

	type body struct {
		HasNewListings bool
		Listings       []clListing
	}

	t.Run("get - should return ALL new listings", func(t *testing.T) {
		// this test case basically rides off the fact that initializing the
		// mock clients will retrieve no new listings

		req, err := http.NewRequest(http.MethodGet, "/listing?ID=99&Datetime=0", nil)
		assert.NoError(t, err)
		res := httptest.NewRecorder()

		api.handleListing(res, req)

		resBody := body{}
		readBodyInto(t, res.Body, &resBody)

		assert.True(t, resBody.HasNewListings)
		assert.Len(t, resBody.Listings, 2)
	})
}

func TestHandleSearch(t *testing.T) {
	api := setupTestAPI(t)
	t.Run("get - gets a single search record", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/api/v1/search?ID=99", nil)
		assert.NoError(t, err)
		res := httptest.NewRecorder()

		api.handleSearch(res, req)

		resBody := clSearch{}
		readBodyInto(t, res.Body, &resBody)

		assert.Equal(t, resBody.Name, "Test seach 1")
	})

	t.Run("get - gets a list of searches", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/api/v1/search", nil)
		assert.NoError(t, err)
		res := httptest.NewRecorder()

		api.handleSearch(res, req)

		resBody := []clSearch{}
		readBodyInto(t, res.Body, &resBody)

		assert.Equal(t, 2, len(resBody))
	})

	t.Run("post - recieves data", func(t *testing.T) {
		// make a body
		type body struct {
			ID        int
			Name      string
			URL       string
			Confirmed bool
		}

		b := body{Name: "test monitor 1", URL: "www.anything.com"}
		data, err := json.Marshal(b)
		assert.NoError(t, err)
		reader := bytes.NewReader(data)

		req, err := http.NewRequest(http.MethodPost, "/", reader)
		assert.NoError(t, err)
		res := httptest.NewRecorder()

		// call handelr
		api.handleSearch(res, req)

		resBody := body{}
		readBodyInto(t, res.Body, &resBody)

		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, b.Name, resBody.Name)
		assert.Equal(t, b.URL, resBody.URL)
		assert.Equal(t, false, resBody.Confirmed)
	})

	t.Run("post - invalid url", func(t *testing.T) {
		type body struct {
			Name string
			URL  string
		}

		b := body{Name: "badurl", URL: badCraigslistURL}
		data, err := json.Marshal(b)
		assert.NoError(t, err)
		reader := bytes.NewReader(data)

		req, err := http.NewRequest(http.MethodPost, "/", reader)
		assert.NoError(t, err)
		res := httptest.NewRecorder()

		api.handleSearch(res, req)

		message, err := ioutil.ReadAll(res.Body)

		assert.Equal(t, http.StatusBadRequest, res.Code)
		assert.Equal(t, "url provided is not a compatible with craigslist\n", string(message))
	})

	t.Run("post - should not accept duplicate urls", func(t *testing.T) {
		type body struct {
			Name string
			URL  string
		}

		// ROUND 1
		b := body{Name: "goodurl", URL: "www.goodurl.com"}
		data, err := json.Marshal(b)
		assert.NoError(t, err)
		reader := bytes.NewReader(data)

		req, err := http.NewRequest(http.MethodPost, "/", reader)
		assert.NoError(t, err)
		res := httptest.NewRecorder()

		api.handleSearch(res, req)

		resBody := clSearch{}
		readBodyInto(t, res.Body, &resBody)

		assert.Greater(t, resBody.ID, 0)

		// ROUND 2
		b = body{Name: "goodurl", URL: "www.goodurl.com"}
		data, err = json.Marshal(b)
		assert.NoError(t, err)
		reader = bytes.NewReader(data)

		req, err = http.NewRequest(http.MethodPost, "/", reader)
		assert.NoError(t, err)
		res = httptest.NewRecorder()

		api.handleSearch(res, req)

		message, err := ioutil.ReadAll(res.Body)
		assert.NoError(t, err)

		assert.Equal(t, "URLs must be unique\n", string(message))

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
