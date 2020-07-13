package main

import (
	"fmt"
	"log"
	"net/http"

	craigslist "github.com/tmunayyer/go-craigslist"
)

func init() {
	setEvnironmentVariables()
}

func main() {
	fs := http.FileServer(http.Dir("./dist"))
	http.Handle("/", fs)

	initializeAPI()

	fmt.Println("listening on: localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func initializeAPI() {
	cl := craigslist.NewClient("newyork")
	db := newDBClient()

	ps := newPollingService(cl, db)
	api := newAPIService(cl, db, ps)

	http.HandleFunc("/api/v1/search", api.handleSearch)
	http.HandleFunc("/api/v1/listing", api.handleListing)
	http.HandleFunc("/api/v1/metric", api.handleMetric)
}
