package main

import (
	"errors"
	"flag"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"strings"
	"subdomain-http-server/internal"
	"sync"
)

type Mux struct {
	amelia, ryan, sheldon, main *http.ServeMux
}

func (mux *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Host == "lobo.codes" {
		mux.main.ServeHTTP(w, r)
		return
	}
	domainParts := strings.Split(r.Host, ".")
	if domainParts[0] == "amelia" {
		mux.amelia.ServeHTTP(w, r)
	} else if domainParts[0] == "ryan" {
		mux.ryan.ServeHTTP(w, r)
	} else if domainParts[0] == "sheldon" {
		mux.sheldon.ServeHTTP(w, r)
	} else {
		http.Error(w, "Not found", 404)
	}
}

func main() {
	var err error

	// Input arguments
	portPtr := flag.String("port", ":80", "http port (:80)")
	sslPortPtr := flag.String("ssl-port", ":443", "https port (:443)")
	sslKeyDirPtr := flag.String("ssl-key-dir", "/etc/letsencrypt/live/amelia.lobo.codes",
		"directory with SSL key and cert files")
	logPtr := flag.String("log", "stdout", "log file (stdout)")

	flag.Parse()

	// Ip2location db file init
	geoDbFile := "IP2LOCATION-LITE-DB11.IPV6.BIN"
	if err := internal.InitGeoDB(geoDbFile); err != nil {
		fmt.Printf("Failed opening ip2location db file %s : %s\n", geoDbFile, err)
		os.Exit(1)
	}

	// Subdomains served
	internal.HandlerInfoMap = map[string]internal.HandlerInfo{
		"amelia":  {},
		"ryan":    {},
		"sheldon": {},
	}

	// Init gorm db for each subdomain
	// Also, save and parse html templates
	for subdomain, handlerInfo := range internal.HandlerInfoMap {
		dbFile := fmt.Sprintf("%s/requests.db", subdomain)

		// Check if exists
		_, existsErr := os.Stat(dbFile)

		handlerInfo.RequestsDb, err = gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
		if err != nil {
			fmt.Printf("Failed to open sqlite3 db %s : %s\n", dbFile, err)
			os.Exit(1)
		}

		// If new db file, migrate the schema
		if errors.Is(existsErr, os.ErrNotExist) {
			err = handlerInfo.RequestsDb.AutoMigrate(&internal.RequestInfo{})
			if err != nil {
				fmt.Printf("Failed to migrate schema to new sqlite3 db %s : %s\n", dbFile, err)
				os.Exit(1)
			}
		}

		handlerInfo.PathMap = internal.GetPathMap(subdomain)

		internal.HandlerInfoMap[subdomain] = handlerInfo
	}

	// Logging
	if *logPtr != "stdout" {
		// If the file doesn't exist, create it or append to the file
		file, err := os.OpenFile(*logPtr, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatal(err)
		}
		log.SetOutput(file)
	}

	// Setup mux and handlers
	mux := &Mux{
		amelia:  http.NewServeMux(),
		ryan:    http.NewServeMux(),
		sheldon: http.NewServeMux(),
		main:    http.NewServeMux(),
	}

	mux.amelia.HandleFunc("/", internal.AmeliaHandler)
	mux.ryan.HandleFunc("/", internal.RyanHandler)
	mux.sheldon.HandleFunc("/", internal.SheldonHandler)
	mux.main.HandleFunc("/", internal.MainHandler)

	// Wait for both http & https servers to finish
	var serversWaitGroup sync.WaitGroup
	serversWaitGroup.Add(2)

	go func() {
		fmt.Printf("Listening on port %s ...\n", *portPtr)
		err := http.ListenAndServe(*portPtr, mux)
		if err != nil {
			fmt.Printf("Error serving on port %s : %s\n", *portPtr, err)
		}
		serversWaitGroup.Done()
	}()

	go func() {
		if *sslPortPtr != "" {
			fmt.Printf("Listening on SSL port %s ...\n", *sslPortPtr)
			err := http.ListenAndServeTLS(*sslPortPtr, *sslKeyDirPtr+"/fullchain.pem",
				*sslKeyDirPtr+"/privkey.pem", mux)
			if err != nil {
				fmt.Printf("Error serving with SSL on port %s : %s\n", *sslPortPtr, err)
			}
		}
		serversWaitGroup.Done()
	}()

	serversWaitGroup.Wait()
}
