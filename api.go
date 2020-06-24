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

func (s *apiService) handleMonitor(w http.ResponseWriter, req *http.Request) {
	type requestBody struct {
		Name string
		URL  string
	}

	d := json.NewDecoder(req.Body)
	body := requestBody{}
	err := d.Decode(&body)
	if err != nil {
		apiErrorHandler(w, http.StatusBadRequest, "handleMonitorURL", "could not decode user payload", err)
		return
	}

	// validate the url before we put it in the DB
	_, err = s.cl.GetListings(req.Context(), body.URL)
	if err != nil {
		apiErrorHandler(w, http.StatusBadRequest, "handleMonitorURL", "url provided is not a compatible with craigslist", err)
		return
	}

	record, err := s.db.saveSearch(clSearch{
		Name: body.Name,
		URL:  body.URL,
	})
	if err != nil {
		apiErrorHandler(w, http.StatusInternalServerError, "handleMonitorURL", "could not save the information", err)
		return
	}

	// spin up the process to monitor the record
	go s.ps.poll(context.TODO(), record)

	data, err := json.Marshal(record)
	if err != nil {
		apiErrorHandler(w, http.StatusInternalServerError, "handleMonitorURL", "problems formatting the data", err)
		return
	}

	w.Write(data)
}

func (s *apiService) handleListing(w http.ResponseWriter, req *http.Request) {
	type resObject struct {
		HasNewListings bool
		Listings       []craigslist.Listing
	}

	queryValues := req.URL.Query()
	ID, err := strconv.Atoi(queryValues["ID"][0])
	if err != nil {
		apiErrorHandler(w, http.StatusBadRequest, "handleNewListings", "invalid id provided", err)
		return
	}

	listings, err := s.ps.flush(ID)
	if err != nil {
		apiErrorHandler(w, http.StatusBadRequest, "handleNewListings", "invalid id provided", err)
		return
	}

	// no new listings yet
	if len(listings) < 1 {
		data, err := json.Marshal(resObject{
			HasNewListings: false,
		})
		if err != nil {
			apiErrorHandler(w, http.StatusInternalServerError, "handleNewListings", "problems formatting the data", err)
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
		apiErrorHandler(w, http.StatusInternalServerError, "handleNewListings", "problems formatting the data", err)
		return
	}

	w.Write(data)
}
