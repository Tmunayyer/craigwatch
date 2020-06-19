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
}

func newAPIService(cl craigslist.Client, db connection) *apiService {
	api := apiService{
		cl: cl,
		db: db,
	}

	return &api
}

func (s *apiService) handleMonitorURL(w http.ResponseWriter, req *http.Request) {
	type requestBody struct {
		Email string
		URL   string
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

	_, err = s.db.saveURL(url{
		email: body.Email,
		url:   body.URL,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("could not save the information: %v", err), http.StatusInternalServerError)
	}
	// TODO: refactgor this to return the db record instead of listings, this wont actually
	// hit the craigslist api in the end when im done

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
