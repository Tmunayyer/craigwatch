package main

import (
	"io/ioutil"
	"os"
	"strings"
)

func setEvnironmentVariables() {
	data, err := ioutil.ReadFile("./.env")
	if err != nil {
		panic(err)
	}
	tuples := strings.Split(string(data), "\n")
	for _, tuple := range tuples {
		if tuple == "" {
			continue
		}
		keyval := strings.Split(tuple, "=")
		key := keyval[0]
		val := keyval[1]
		os.Setenv(key, val)
	}
}
