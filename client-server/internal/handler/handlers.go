package handler

import (
	"log"
	"net/http"
)

//func ReturnHandler(countryPath string) func(http.ResponseWriter, *http.Request) {
func ReturnHandler(countryPath string) http.HandlerFunc {
	var responseString string
	var count int
	if countryPath == "asia" {
		responseString = "India";
	} else if countryPath == "america" {
		responseString = "United States"
	}

	return func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != "GET" {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Printf("[%d] Received %s request from: %s", count, countryPath, request.UserAgent())
		count++
		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte(responseString))
	}
}