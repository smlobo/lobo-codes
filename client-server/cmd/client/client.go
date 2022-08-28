package main

import (
	"client-server/internal"
	"io"
	"log"
	"net/http"
	"time"
)

var requestCount int

func main() {
	// Keep hitting endpoints indefinitely
	for {
		makeRequest(internal.AsiaPath)
		makeRequest(internal.AmericaPath)
	}
}

func makeRequest(restPath string) {
	response, err := http.Get("http://" + internal.Endpoint + restPath)
	if err != nil {
		log.Printf("[%d] Failed", requestCount)
	}
	body, err := io.ReadAll(response.Body)
	log.Printf("[%d] Got: %d : %s", requestCount, response.StatusCode, string(body))
	requestCount++
	_ = response.Body.Close()
	time.Sleep(time.Second*2)
}