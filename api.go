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

	d := json.NewDecoder(req.Body)
	body := requestBody{}
	err := d.Decode(&body)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not decode user payload: %v", err), http.StatusBadRequest)
		return
	}

	// validate the url before we put it in the DB
	_, err = s.cl.GetListings(req.Context(), body.URL)
	if err != nil {
		http.Error(w, fmt.Sprint("url provided is not a compatible with craigslist"), http.StatusBadRequest)
		return
	}

	record, err := s.db.saveURL(craigslistQuery{
		Email: body.Email,
		URL:   body.URL,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("could not save the information: %v", err), http.StatusInternalServerError)
	}

	data, err := json.Marshal(record)
	if err != nil {
		http.Error(w, fmt.Sprintf("problems formatting the data: %v", err), http.StatusInternalServerError)
		return
	}

	w.Write(data)
}
