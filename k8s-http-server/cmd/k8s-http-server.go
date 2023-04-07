package main

import (
	"flag"
	"fmt"
	"github.com/gocql/gocql"
	"html/template"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/go-chi/chi/v5"
	"k8s-http-server/hostrouter"
	"k8s-http-server/internal"
)

func main() {
	var err error

	// Input arguments
	portPtr := flag.String("port", ":80", "http port (:80)")
	sslPortPtr := flag.String("ssl-port", ":443", "https port (:443)")
	sslKeyDirPtr := flag.String("ssl-key-dir", "./ssl-certificates",
		"directory with SSL key and cert files")
	logPtr := flag.String("log", "stdout", "log file (stdout)")
	cassandraPtr := flag.String("cassandra-db", "cassandra-internal", "cassandra server/service")
	testOsReleasePtr := flag.String("os-release", "/etc/os-release", "os-release file")

	flag.Parse()

	// Ip2location db file init
	geoDbFile := "IP2LOCATION-LITE-DB11.IPV6.BIN"
	if err := internal.InitGeoDB(geoDbFile); err != nil {
		fmt.Printf("Failed opening ip2location db file %s : %s\n", geoDbFile, err)
		os.Exit(1)
	}

	// Parse os-release
	internal.InitOsRelease(*testOsReleasePtr)

	// Get kubernetes version
	internal.InitKubernetesInfo()

	// Subdomains served
	internal.HandlerInfoMap = map[string]internal.HandlerInfo{
		"amelia":  {},
		"ryan":    {},
		"sheldon": {},
		"domain":  {},
	}

	// Init cassandra db
	internal.CassandraCluster = gocql.NewCluster(*cassandraPtr)
	internal.CassandraCluster.Keyspace = "lobo_codes"
	internal.CassandraCluster.Consistency = gocql.Quorum
	session, err := internal.CassandraCluster.CreateSession()
	if err != nil {
		log.Fatalf("failed to create session with cassandra database: %s; %s",
			internal.CassandraCluster.Hosts[0], err.Error())
	}
	session.Close()

	// Init cassandra table for each subdomain
	// Also, save and parse html templates
	for subdomain, handlerInfo := range internal.HandlerInfoMap {
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
	internal.NotFoundTemplate = template.Must(template.ParseFiles("common/notfound.html"))

	// Setup chi with hostrouter
	router := chi.NewRouter()
	hostRouter := hostrouter.New()

	hostRouter.Map("amelia.lobo.codes", ameliaRouter())
	hostRouter.Map("ryan.lobo.codes", ryanRouter())
	hostRouter.Map("sheldon.lobo.codes", sheldonRouter())
	hostRouter.Map("lobo.codes", domainRouter())

	// Testing locally
	hostRouter.Map("amelia.lobo.codes"+*portPtr, ameliaRouter())
	hostRouter.Map("ryan.lobo.codes"+*portPtr, ryanRouter())
	hostRouter.Map("sheldon.lobo.codes"+*portPtr, sheldonRouter())
	hostRouter.Map("lobo.codes"+*portPtr, domainRouter())

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
	r.Get("/visitors", internal.RyanHandler)
	r.Get("/visitors.html", internal.RyanHandler)
	r.NotFound(internal.NotFoundHandler)

	// Other static content
	ryanFileServer := http.FileServer(http.Dir("./ryan"))
	r.Handle("/static/*", ryanFileServer)

	return r
}

func sheldonRouter() chi.Router {
	r := chi.NewRouter()
	r.Get("/", internal.SheldonHandler)
	r.Get("/visitors", internal.SheldonHandler)
	r.Get("/visitors.html", internal.SheldonHandler)
	r.NotFound(internal.NotFoundHandler)

	// Other static content
	sheldonFileServer := http.FileServer(http.Dir("./sheldon"))
	r.Handle("/static/*", sheldonFileServer)

	return r
}

func domainRouter() chi.Router {
	r := chi.NewRouter()
	r.Get("/", internal.DomainHandler)
	r.Get("/visitors", internal.DomainHandler)
	r.Get("/visitors.html", internal.DomainHandler)
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
	notFoundFileServer := http.FileServer(http.Dir("./common"))
	r.Handle("/static/*", notFoundFileServer)

	return r
}
