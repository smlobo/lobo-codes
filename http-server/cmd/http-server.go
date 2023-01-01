package main

import (
	"flag"
	"fmt"
	"github.com/go-chi/chi/v5"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"html/template"
	"http-server/internal"
	"log"
	"net/http"
	"os"
	"sync"
)

func main() {
	// Input arguments
	portPtr := flag.String("port", ":80", "http port (:80)")
	sslPortPtr := flag.String("ssl-port", ":443", "https port (:443)")
	htmlDirPtr := flag.String("html-dir", "/var/www/amelia-html/portfolio-website",
		"directory with static html to serve")
	geoDbPtr := flag.String("geodb", "/home/pi/db/IP2LOCATION-LITE-DB11.IPV6.BIN",
		"ip2location db")
	requestsDbPtr := flag.String("requestsdb", "requests.db", "db to store incoming requests ip")
	sslKeyDirPtr := flag.String("ssl-key-dir", "/etc/letsencrypt/live/amelia.lobo.codes",
		"directory with SSL key and cert files")
	logPtr := flag.String("log", "stdout", "log file (stdout)")

	flag.Parse()

	// Ip2location db file init
	if err := internal.InitGeoDB(*geoDbPtr); err != nil {
		fmt.Printf("Failed opening ip2location db file %s : %s\n", *geoDbPtr, err)
		os.Exit(1)
	}

	// Init gorm db
	requestsDb, err := gorm.Open(sqlite.Open(*requestsDbPtr), &gorm.Config{})
	if err != nil {
		fmt.Printf("Failed to connect to sqlite db %s : %s\n", *requestsDbPtr, err)
		os.Exit(1)
	}
	// Migrate the schema
	//requestsDb.AutoMigrate(&internal.RequestInfo{})

	// Logging
	if *logPtr != "stdout" {
		// If the file doesn't exist, create it or append to the file
		file, err := os.OpenFile(*logPtr, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatal(err)
		}
		log.SetOutput(file)
	}

	// Router
	router := chi.NewRouter()

	// Main pattern to be logged to db
	router.Get("/", internal.ReturnHandler(*htmlDirPtr, requestsDb))

	// Template
	tmpl := template.Must(template.ParseFiles("template/visitors.html"))
	router.Get("/myvisitors", internal.TemplateHandler(tmpl, requestsDb))

	// Other static content
	fileServer := http.FileServer(http.Dir(*htmlDirPtr))
	//router.Handle("/*", http.StripPrefix("/", fileServer))
	router.Handle("/*", fileServer)

	// Wait for both http & https servers to finish
	var serversWaitGroup sync.WaitGroup
	serversWaitGroup.Add(2)

	go func() {
		fmt.Printf("Listening on port %s ...\n", *portPtr)
		err := http.ListenAndServe(*portPtr, router)
		if err != nil {
			fmt.Printf("Error serving on port %s : %s\n", *portPtr, err)
		}
		serversWaitGroup.Done()
	}()

	go func() {
		if *sslPortPtr != "" {
			fmt.Printf("Listening on SSL port %s ...\n", *sslPortPtr)
			err := http.ListenAndServeTLS(*sslPortPtr, *sslKeyDirPtr+"/fullchain.pem",
				*sslKeyDirPtr+"/privkey.pem", router)
			if err != nil {
				fmt.Printf("Error serving with SSL on port %s : %s\n", *sslPortPtr, err)
			}
		}
		serversWaitGroup.Done()
	}()

	serversWaitGroup.Wait()
}
