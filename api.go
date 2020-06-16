package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	craigslist "github.com/tmunayyer/go-craigslist"
)

type api interface {
	handleMonitorURL(w http.ResponseWriter, r *http.Request)
}

// env (environment) will house global references. Endpoints will be methods
// on the env to have access
type handlerEnv struct {
	cl craigslist.Client
}

// representation of an endpoint as a struct
type endpoint struct {
	path    string
	handler func(w http.ResponseWriter, r *http.Request, e *handlerEnv)
}

func initializeAPI() {
	cl := craigslist.NewClient("newyork")

	environment := handlerEnv{
		cl: cl,
	}

	for _, ep := range []endpoint{
		{
			path:    "/api/monitorurl",
			handler: handleMonitorURL,
		},
	} {
		http.HandleFunc(ep.path, func(w http.ResponseWriter, r *http.Request) {
			ep.handler(w, r, &environment)
		})
	}
}

func defaultResponse(w http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal("unsupported method")
	if err != nil {
		fmt.Println("error in defaultResposne:", err)
	}

	w.Write(data)
}

func internalError(e error, w http.ResponseWriter, r *http.Request) {
	if e != nil {
		fmt.Println(e)
		w.WriteHeader(500)
	}
}

func handleMonitorURL(w http.ResponseWriter, r *http.Request, e *handlerEnv) {
	w.WriteHeader(200)
	switch r.Method {
	case "POST":
		listings, err := e.cl.GetListings(r.Context(), "https://sfbay.craigslist.org/d/antiques/search/ata")
		internalError(err, w, r)

		data, err := json.Marshal(listings)
		internalError(err, w, r)

		w.Write(data)
	default:
		defaultResponse(w, r)
	}
}
