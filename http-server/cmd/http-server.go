package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"http-server/internal"
	"log"
	"net/http"
	"os"
)

func main() {
	// Input arguments required : port, static directory to serve, ip geodb mapping, request info db
	if len(os.Args) != 5 {
		fmt.Printf("Usage: %s <port> <static-directory> <geo-db> <requests-db>\n", os.Args[0])
		os.Exit(1)
	}

	// Port (as ":XXXX")
	port := os.Args[1]

	// Root of the static directory to serve
	directory := os.Args[2]

	// Ip2location db file init
	if err := internal.InitGeoDB(os.Args[3]); err != nil {
		fmt.Printf("Failed opening ip2location db file %s : %s\n", os.Args[3], err)
		os.Exit(1)
	}

	// Init gorm db
	requestsDb, err := gorm.Open(sqlite.Open(os.Args[4]), &gorm.Config{})
	if err != nil {
		fmt.Printf("Failed to connect to sqlite db %s : %s\n", os.Args[4], err)
		os.Exit(1)
	}
	// Migrate the schema
	//requestsDb.AutoMigrate(&internal.RequestInfo{})

	// Router
	router := chi.NewRouter()

	// Main pattern to be logged to db
	router.Get("/", internal.ReturnHandler(directory, requestsDb))

	// Other static content
	fileServer := http.FileServer(http.Dir(directory))
	//router.Handle("/*", http.StripPrefix("/", fileServer))
	router.Handle("/*", fileServer)

	log.Printf("Listening on port %s ...\n", port)
	http.ListenAndServe(port, router)
}
