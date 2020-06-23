package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	craigslist "github.com/tmunayyer/go-craigslist"
)

type apiService struct {
	cl craigslist.Client
	db connection
	ps pollingService
}

func newAPIService(cl craigslist.Client, db connection, ps pollingService) *apiService {
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

func (s *apiService) handleMonitorURL(w http.ResponseWriter, req *http.Request) {
	type requestBody struct {
		Email string
		URL   string
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
		Email: body.Email,
		URL:   body.URL,
	})
	if err != nil {
		apiErrorHandler(w, http.StatusInternalServerError, "handleMonitorURL", "could not save the information", err)
		return
	}

	data, err := json.Marshal(record)
	if err != nil {
		apiErrorHandler(w, http.StatusInternalServerError, "handleMonitorURL", "problems formatting the data", err)
		return
	}

	w.Write(data)
}

func (s *apiService) handleNewListings(w http.ResponseWriter, req *http.Request) {
	type resObject struct {
		hasNewListings bool
		listings       []craigslist.Listing
	}
	listings, err := s.ps.flush()
	if err != nil {
		apiErrorHandler(w, http.StatusInternalServerError, "handleNewListings", "problems retrieving polled data", err)
		return
	}

	// no new listings yet
	if len(listings) == 0 {
		data, err := json.Marshal(resObject{
			hasNewListings: false,
		})
		if err != nil {
			apiErrorHandler(w, http.StatusInternalServerError, "handleNewListings", "problems formatting the data", err)
			return
		}

		w.Write(data)
		return
	}

	data, err := json.Marshal(resObject{
		hasNewListings: true,
		listings:       listings,
	})

	w.Write(data)
}
