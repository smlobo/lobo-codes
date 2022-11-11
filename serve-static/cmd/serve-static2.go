package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"serve-static/internal"
)

func main() {
	// Input arguments required : port & static directory to serve
	if len(os.Args) != 4 {
		fmt.Printf("Usage: %s <port> <static-directory> <geo-db>\n", os.Args[0])
		os.Exit(1)
	}

	port := os.Args[1]
	directory := os.Args[2]
	if err := internal.InitGeoDB(os.Args[3]); err != nil {
		fmt.Printf("Failed opening ip2location db file %s : %s\n", os.Args[3], err)
		os.Exit(2)
	}

	log.Printf("Listening on port %s ...\n", port)
	http.ListenAndServe(port, internal.ReturnHandler(directory))
}
