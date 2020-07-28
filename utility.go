package main

import (
	"time"
)

func newUnixDate(date string, tz *time.Location) int64 {
	layout := "2006-01-02 15:04:05"
	var formattedDate time.Time
	var err error
	if tz == nil {
		loc, _ := time.LoadLocation("America/New_York")
		formattedDate, err = time.ParseInLocation(layout, date, loc)
	} else {
		formattedDate, err = time.ParseInLocation(layout, date, tz)
	}

	if err != nil {
		panic(err)
	}

	unixTime := formattedDate.Unix()

	return unixTime
}

func unixToLocal(unixTimestamp int64, tz *time.Location) time.Time {
	if tz == nil {
		tz, _ = time.LoadLocation("America/New_York")
	}

	date := time.Unix(unixTimestamp, 0)
	return date.In(tz)
}
