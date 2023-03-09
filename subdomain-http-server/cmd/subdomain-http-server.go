package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/go-chi/chi/v5"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"html/template"
	"log"
	"net/http"
	"os"
	"subdomain-http-server/hostrouter"
	"subdomain-http-server/internal"
	"sync"
)

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
		"domain":  {},
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

	// Not found template
	internal.NotFoundTemplate = template.Must(template.ParseFiles("notfound/notfound.html"))

	// Setup chi with hostrouter
	router := chi.NewRouter()
	hostRouter := hostrouter.New()
	hostRouter.Map("amelia.lobo.codes", ameliaRouter())
	hostRouter.Map("ryan.lobo.codes", ryanRouter())
	hostRouter.Map("sheldon.lobo.codes", sheldonRouter())
	hostRouter.Map("lobo.codes", domainRouter())
	hostRouter.Map("*", notFoundRouter())
	router.Mount("/", hostRouter)

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

func ameliaRouter() chi.Router {
	r := chi.NewRouter()
	r.Get("/", internal.AmeliaHandler)
	r.Get("/visitors", internal.AmeliaHandler)
	r.Get("/visitors.html", internal.AmeliaHandler)
	r.NotFound(internal.NotFoundHandler)

	// Other static content
	ameliaFileServer := http.FileServer(http.Dir("./amelia"))
	r.Handle("/static/*", ameliaFileServer)

	return r
}

func ryanRouter() chi.Router {
	r := chi.NewRouter()
	r.Get("/", internal.RyanHandler)
	r.NotFound(internal.NotFoundHandler)

	// Other static content
	ryanFileServer := http.FileServer(http.Dir("./ryan"))
	r.Handle("/static/*", ryanFileServer)

	return r
}

func sheldonRouter() chi.Router {
	r := chi.NewRouter()
	r.Get("/", internal.SheldonHandler)
	r.NotFound(internal.NotFoundHandler)

	// Other static content
	sheldonFileServer := http.FileServer(http.Dir("./sheldon"))
	r.Handle("/static/*", sheldonFileServer)

	return r
}

func domainRouter() chi.Router {
	r := chi.NewRouter()
	r.Get("/", internal.DomainHandler)
	r.NotFound(internal.NotFoundHandler)

	// Other static content
	domainFileServer := http.FileServer(http.Dir("./domain"))
	r.Handle("/static/*", domainFileServer)

	return r
}

func notFoundRouter() chi.Router {
	r := chi.NewRouter()
	r.NotFound(internal.NotFoundHandler)

	// Other static content
	notFoundFileServer := http.FileServer(http.Dir("./notfound"))
	r.Handle("/static/*", notFoundFileServer)

	return r
}
