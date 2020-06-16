package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	initializeAPI()

	fmt.Println("listening on: localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
