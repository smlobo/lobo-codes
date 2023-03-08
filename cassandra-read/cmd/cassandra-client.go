package main

import (
	"cassandra-client/internal"
	"net/http"
)

func main() {
	server := http.NewServeMux()
	server.HandleFunc("/", internal.ReturnHandler())
	_ = http.ListenAndServe(":8080", server)
}
