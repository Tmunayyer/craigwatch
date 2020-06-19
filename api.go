package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	craigslist "github.com/tmunayyer/go-craigslist"
)

func initializeAPI() {
	cl := craigslist.NewClient("newyork")
	api := newAPIService(cl)

	http.HandleFunc("/api/monitorurl", api.handleMonitorURL)
}

type apiService struct {
	cl craigslist.Client
}

func newAPIService(cl craigslist.Client) *apiService {
	api := apiService{
		cl: cl,
	}

	return &api
}

func (s *apiService) handleMonitorURL(w http.ResponseWriter, req *http.Request) {
	type requestBody struct {
		URL string
	}

	type responseObj struct {
		Searched string
		Results  []craigslist.Listing
	}

	d := json.NewDecoder(req.Body)
	body := requestBody{}
	err := d.Decode(&body)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not decode user payload: %v", err), http.StatusBadRequest)
		return
	}

	listings, err := s.cl.GetListings(req.Context(), body.URL)
	if err != nil {
		http.Error(w, fmt.Sprint("url provided is not a compatible with craigslist"), http.StatusBadRequest)
		return
	}

	resObj := responseObj{
		Searched: body.URL,
		Results:  listings,
	}

	data, err := json.Marshal(resObj)
	if err != nil {
		http.Error(w, fmt.Sprintf("problems formatting the data: %v", err), http.StatusInternalServerError)
		return
	}

	w.Write(data)
}
