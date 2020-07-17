package main

import (
	"io/ioutil"
	"os"
	"strings"
)

func setEvnironmentVariables() {
	envFilePath := "./.env"
	if os.Getenv("MODE") == "production" {
		// there is nothing to do here in production
		return
	}

	data, err := ioutil.ReadFile(envFilePath)
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
