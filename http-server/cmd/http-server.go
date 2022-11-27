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
	// Input arguments required : port, ssl-port, static directory to serve, ip geodb mapping, request info db
	if len(os.Args) != 7 {
		fmt.Printf("Usage: %s <port> <ssl-port> <static-directory> <geo-db> <requests-db> <ssl-key-dir>\n",
			os.Args[0])
		os.Exit(1)
	}

	// Port (as ":XXXX")
	port := os.Args[1]
	sslPort := os.Args[2]

	// Root of the static directory to serve
	directory := os.Args[3]

	// Ip2location db file init
	if err := internal.InitGeoDB(os.Args[4]); err != nil {
		fmt.Printf("Failed opening ip2location db file %s : %s\n", os.Args[4], err)
		os.Exit(1)
	}

	// Init gorm db
	requestsDb, err := gorm.Open(sqlite.Open(os.Args[5]), &gorm.Config{})
	if err != nil {
		fmt.Printf("Failed to connect to sqlite db %s : %s\n", os.Args[5], err)
		os.Exit(1)
	}
	// Migrate the schema
	//requestsDb.AutoMigrate(&internal.RequestInfo{})

	// SSL key dir
	sslDir := os.Args[6]

	// Router
	router := chi.NewRouter()

	// Main pattern to be logged to db
	router.Get("/", internal.ReturnHandler(directory, requestsDb))

	// Other static content
	fileServer := http.FileServer(http.Dir(directory))
	//router.Handle("/*", http.StripPrefix("/", fileServer))
	router.Handle("/*", fileServer)

	log.Printf("Listening on port %s ...\n", port)
	go func() {
		err := http.ListenAndServe(port, router)
		if err != nil {
			fmt.Printf("Error serving on port %s : %s", port, err)
		}
	}()
	err = http.ListenAndServeTLS(sslPort, sslDir+"/fullchain.pem", sslDir+"/privkey.pem", router)
	if err != nil {
		fmt.Printf("Error serving with SSL on port %s : %s", sslPort, err)
	}
}
