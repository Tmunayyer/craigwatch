package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	craigslist "github.com/tmunayyer/go-craigslist"
)

type apiService struct {
	cl craigslist.API
	db connection
	ps pollingService
}

func newAPIService(cl craigslist.API, db connection, ps pollingService) *apiService {
	api := apiService{
		cl: cl,
		db: db,
		ps: ps,
	}

	return &api
}

func apiErrorHandler(w http.ResponseWriter, status int, endpoint string, message string, err error) {
	fmt.Println("error from handler: " + endpoint)
	fmt.Println("-- outbound message: " + message)
	fmt.Println(fmt.Sprintf("-- the error: %v", err))

	http.Error(w, message, http.StatusBadRequest)
}

func (s *apiService) handleListing(w http.ResponseWriter, req *http.Request) {
	//=========
	//== GET ==
	//=========
	if req.Method == http.MethodGet {
		type resObject struct {
			HasNewListings bool
			Listings       []clListing
		}

		queryValues := req.URL.Query()
		qValID, has := queryValues["ID"]
		if !has {
			apiErrorHandler(w, http.StatusBadRequest, "handleListing", "missing query value: ID", nil)
			return
		}

		ID, err := strconv.Atoi(qValID[0])
		if err != nil {
			apiErrorHandler(w, http.StatusBadRequest, "handleListing", "invalid id provided", err)
			return
		}

		qValDatetime, has := queryValues["Datetime"]
		if !has {
			apiErrorHandler(w, http.StatusBadRequest, "handleListing", "missing query value: Datetime", nil)
			return
		}

		unixTimestamp, err := strconv.Atoi(qValDatetime[0])
		if err != nil {
			apiErrorHandler(w, http.StatusBadRequest, "handleListing", "invalid Datetime provided", err)
			return
		}

		listings, err := s.db.getListingMultiAfter(ID, int64(unixTimestamp))
		if err != nil {
			apiErrorHandler(w, http.StatusInternalServerError, "handleListing", "err retrieving listings from db", err)
			return
		}

		// no new listings yet
		if len(listings) < 1 {
			data, err := json.Marshal(resObject{
				HasNewListings: false,
			})
			if err != nil {
				apiErrorHandler(w, http.StatusInternalServerError, "handleListing", "problems formatting the data", err)
				return
			}

			w.Write(data)
			return
		}

		data, err := json.Marshal(resObject{
			HasNewListings: true,
			Listings:       listings,
		})
		if err != nil {
			apiErrorHandler(w, http.StatusInternalServerError, "handleListing", "problems formatting the data", err)
			return
		}

		w.Write(data)
		return
	}

	//=================
	//== UNSUPPORTED ==
	//=================
	apiErrorHandler(w, http.StatusNotImplemented, "handleListing", "method not supported: "+req.Method, nil)
}

func (s *apiService) handleSearch(w http.ResponseWriter, req *http.Request) {
	//=========
	//== GET ==
	//=========
	if req.Method == http.MethodGet {
		queryValues := req.URL.Query()
		qValID, has := queryValues["ID"]
		if has {
			searchID, err := strconv.Atoi(qValID[0])
			if err != nil {
				apiErrorHandler(w, http.StatusBadRequest, "handleSearch", "invalid ID sent", err)
				return
			}

			search, err := s.db.getSearch(searchID)
			if err != nil {
				apiErrorHandler(w, http.StatusInternalServerError, "handleSearch", "unable to retrieve data", err)
				return
			}

			data, err := json.Marshal(search)
			if err != nil {
				apiErrorHandler(w, http.StatusInternalServerError, "handleSearch", "unable to marshal data", err)
				return
			}

			w.Write(data)
			return
		}

		searches, err := s.db.getSearchMulti()
		if err != nil {
			apiErrorHandler(w, http.StatusInternalServerError, "handleSearch", "unable to retrieve data", err)
			return
		}

		data, err := json.Marshal(searches)
		if err != nil {
			apiErrorHandler(w, http.StatusInternalServerError, "handleSearch", "unable to marshal data", err)
			return
		}

		w.Write(data)
		return
	}

	//=========
	//= POST ==
	//=========
	if req.Method == http.MethodPost {
		type requestBody struct {
			Name string
			URL  string
		}

		d := json.NewDecoder(req.Body)
		body := requestBody{}
		err := d.Decode(&body)
		if err != nil {
			apiErrorHandler(w, http.StatusBadRequest, "handleMonitor", "could not decode user payload", err)
			return
		}

		// validate the url before we put it in the DB
		result, err := s.cl.GetListings(req.Context(), body.URL)
		if err != nil {
			apiErrorHandler(w, http.StatusBadRequest, "handleMonitor", "url provided is not a compatible with craigslist", err)
			return
		}

		record, err := s.db.saveSearch(clSearch{
			Name:     body.Name,
			URL:      body.URL,
			Timezone: result.Timezone,
		})
		if err != nil {
			if err.Error() == duplicateURLErrorMessage {
				apiErrorHandler(w, http.StatusBadRequest, "handleMonitor", "URLs must be unique", err)
			} else {
				apiErrorHandler(w, http.StatusInternalServerError, "handleMonitor", "could not save the information", err)
			}
			return
		}

		// spin up the process to monitor the record
		go s.ps.poll(context.TODO(), record)

		data, err := json.Marshal(record)
		if err != nil {
			apiErrorHandler(w, http.StatusInternalServerError, "handleMonitor", "problems formatting the data", err)
			return
		}

		w.Write(data)
		return
	}

	//=================
	//== UNSUPPORTED ==
	//=================
	apiErrorHandler(w, http.StatusNotImplemented, "handleSearch", "method not supported: "+req.Method, nil)
}

func (s *apiService) handleMetric(w http.ResponseWriter, req *http.Request) {
	//=========
	//== GET ==
	//=========
	if req.Method == http.MethodGet {
		queryValues := req.URL.Query()
		qValID, has := queryValues["ID"]
		if !has {
			apiErrorHandler(w, http.StatusBadRequest, "handleSearch", "invalid ID sent", nil)
			return
		}

		searchID, err := strconv.Atoi(qValID[0])
		if err != nil {
			apiErrorHandler(w, http.StatusBadRequest, "handleSearch", "invalid ID sent", err)
			return
		}

		search, err := s.db.getSearch(searchID)
		if err != nil {
			apiErrorHandler(w, http.StatusInternalServerError, "handleSearch", "no search with that ID found", err)
			return
		}

		activity, err := s.db.getSearchActivity(search.ID)
		if err != nil {
			apiErrorHandler(w, http.StatusInternalServerError, "handleSearch", "unable to calculate activity", err)
			return
		}

		data, err := json.Marshal(activity)
		if err != nil {
			apiErrorHandler(w, http.StatusInternalServerError, "handleSearch", "unable to marshal data", err)
			return
		}

		w.Write(data)
		return
	}

	//=================
	//== UNSUPPORTED ==
	//=================
	apiErrorHandler(w, http.StatusNotImplemented, "handleSearch", "method not supported: "+req.Method, nil)
}

func (s *apiService) handleActivityChart(w http.ResponseWriter, req *http.Request) {
	//=========
	//== GET ==
	//=========
	if req.Method == http.MethodGet {
		queryValues := req.URL.Query()
		qValID, has := queryValues["ID"]
		if !has {
			apiErrorHandler(w, http.StatusBadRequest, "handleActivityChart", "invalid ID sent", nil)
			return
		}

		searchID, err := strconv.Atoi(qValID[0])
		if err != nil {
			apiErrorHandler(w, http.StatusBadRequest, "handleActivityChart", "invalid ID sent", err)
			return
		}

		search, err := s.db.getSearch(searchID)
		if err != nil {
			apiErrorHandler(w, http.StatusInternalServerError, "handleActivityChart", "no search with that ID found", err)
			return
		}

		activity, err := s.db.getSearchActivityByHour(search.ID)
		if err != nil {
			apiErrorHandler(w, http.StatusInternalServerError, "handleActivityChart", "unable to calculate activity", err)
			return
		}

		data, err := json.Marshal(activity)
		if err != nil {
			apiErrorHandler(w, http.StatusInternalServerError, "handleActivityChart", "unable to marshal data", err)
			return
		}

		w.Write(data)
		return
	}

	//=================
	//== UNSUPPORTED ==
	//=================
	apiErrorHandler(w, http.StatusNotImplemented, "handleActivityChart", "method not supported: "+req.Method, nil)
}
