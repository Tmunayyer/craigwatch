package main

import (
	"time"
)

func newUnixDate(date string, tz *time.Location) int64 {
	layout := "2006-01-02 15:04"
	var formattedDate time.Time
	var err error
	if tz == nil {
		loc, _ := time.LoadLocation("America/New_York")
		formattedDate, err = time.ParseInLocation(layout, date[:16], loc)
	} else {
		formattedDate, err = time.ParseInLocation(layout, date[:16], tz)
	}

	if err != nil {
		panic(err)
	}

	unixTime := formattedDate.Unix()

	return unixTime
}
