package main

import (
	"client-server/internal"
	"client-server/internal/handler"
	"net/http"
)

func main() {
	countryServer := http.NewServeMux()
	countryServer.HandleFunc(internal.AsiaPath, handler.ReturnHandler("asia"))
	countryServer.HandleFunc(internal.AmericaPath, handler.ReturnHandler("america"))
	_ = http.ListenAndServe(internal.Endpoint, countryServer)
}
