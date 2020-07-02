package main

import (
	"time"
)

func newDate(date string) time.Time {
	layout := "2006-01-02 15:04"
	formattedDate, err := time.Parse(layout, date)
	if err != nil {
		panic(err)
	}

	return formattedDate
}

func newUnixDate(date string) int64 {
	layout := "2006-01-02 15:04"
	formattedDate, err := time.Parse(layout, date)
	if err != nil {
		panic(err)
	}

	unixTime := formattedDate.Unix()

	return unixTime
}
