package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	craigslist "github.com/tmunayyer/go-craigslist"
)

func init() {
	setEvnironmentVariables()
}

func main() {
	fs := http.FileServer(http.Dir("./dist"))
	http.Handle("/", fs)

	initializeAPI()

	port := ":3000"
	prodPort := os.Getenv("PORT")
	if prodPort != "" {
		port = ":" + prodPort
	}
	fmt.Println("listening on: localhost", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func initializeAPI() {
	cl := craigslist.NewClient("newyork")
	db := newDBClient()

	ps := newPollingService(cl, db)
	api := newAPIService(cl, db, ps)

	http.HandleFunc("/api/v1/search", api.handleSearch)
	http.HandleFunc("/api/v1/listing", api.handleListing)
	http.HandleFunc("/api/v1/metric", api.handleMetric)
	http.HandleFunc("/api/v1/activityChart", api.handleActivityChart)
	http.HandleFunc("/api/v1/priceDistribution", api.handlePriceDistribution)
}
