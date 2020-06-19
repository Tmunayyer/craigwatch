package main

import craigslist "github.com/tmunayyer/go-craigslist"

type pollingService interface {
	loadURL() ([]string, error)
	poll(url string)
}

type pollingClient struct {
	cl craigslist.Client
	db connection
}

func newPollingService(cl craigslist.Client, db connection) pollingService {
	pc := pollingClient{
		cl: cl,
		db: db,
	}

	return &pc
}

func (pc *pollingClient) loadURL() ([]string, error) {
	return []string{}, nil
}

func (pc *pollingClient) poll(url string) {}
