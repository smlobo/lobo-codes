package main

import (
	"context"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"k8s-http-server/hostrouter"
	"k8s-http-server/internal"
)

func main() {
	// Input arguments
	portPtr := flag.String("port", ":80", "http port (:80)")
	sslPortPtr := flag.String("ssl-port", ":443", "https port (:443)")
	sslKeyDirPtr := flag.String("ssl-key-dir", "./ssl-certificates",
		"directory with SSL key and cert files")
	logPtr := flag.String("log", "stdout", "log file (stdout)")
	cassandraPtr := flag.String("cassandra-db", "cassandra-internal", "cassandra server/service")
	testOsReleasePtr := flag.String("os-release", "/etc/os-release", "os-release file")
	testGoVersionPtr := flag.String("go-version", "/app/golang_version.txt", "go_version.txt file")

	flag.Parse()

	// Config (currently traces only; TODO all above flags)
	internal.SetupConfig()

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

	// Get Go version
	internal.InitGoVersion(*testGoVersionPtr)

	// Subdomains served
	internal.HandlerInfoMap = map[string]internal.HandlerInfo{
		"amelia":   {},
		"ryan":     {},
		"sheldon":  {},
		"domain":   {},
		"test-vue": {},
	}

	// Init cassandra db
	cassandraSession := internal.InitCassandra(*cassandraPtr)
	defer cassandraSession.Close()

	// Init cassandra table for each subdomain
	// Also, save and parse html templates
	for subdomain, handlerInfo := range internal.HandlerInfoMap {
		handlerInfo.PathMap = internal.GetPathMap(subdomain)

		internal.HandlerInfoMap[subdomain] = handlerInfo
	}
	// Other HTML templates
	internal.HandlerInfoMap["sheldon"].PathMap["graph"] = template.Must(template.ParseFiles("sheldon/graph.html"))
	internal.HandlerInfoMap["sheldon"].PathMap["resume"] = template.Must(template.ParseFiles("sheldon/resume.html"))

	// Logging
	if *logPtr != "stdout" {
		// If the file doesn't exist, create it or append to the file
		file, err := os.OpenFile(*logPtr, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatal(err)
		}
		log.SetOutput(file)
	}

	// Setup OpenTelemetry tracer
	tp := internal.InitTracerProvider("lobo-codes")
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	// Not found template
	internal.NotFoundTemplate = template.Must(template.ParseFiles("common/notfound.html"))

	// Setup chi with hostrouter
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	hostRouter := hostrouter.New()

	hostRouter.Map("amelia.lobo.codes", ameliaRouter())
	hostRouter.Map("ryan.lobo.codes", ryanRouter())
	hostRouter.Map("sheldon.lobo.codes", sheldonRouter())
	hostRouter.Map("lobo.codes", domainRouter())
	hostRouter.Map("test-vue.lobo.codes", testVueRouter())
	hostRouter.Map("api.lobo.codes", apiRouter())

	// Testing locally
	hostRouter.Map("amelia.lobo.codes"+*portPtr, ameliaRouter())
	hostRouter.Map("ryan.lobo.codes"+*portPtr, ryanRouter())
	hostRouter.Map("sheldon.lobo.codes"+*portPtr, sheldonRouter())
	hostRouter.Map("lobo.codes"+*portPtr, domainRouter())
	hostRouter.Map("test-vue.lobo.codes"+*portPtr, testVueRouter())
	hostRouter.Map("api.lobo.codes"+*portPtr, apiRouter())

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

	ameliaWrappedHandler := otelhttp.NewHandler(internal.AmeliaHandler(), "amelia-handler")
	r.Method("GET", "/", ameliaWrappedHandler)
	r.Method("GET", "/visitors", ameliaWrappedHandler)
	r.Method("GET", "/visitors.html", ameliaWrappedHandler)

	r.NotFound(internal.NotFoundHandlerFunc)

	// Other static content
	ameliaFileServer := http.FileServer(http.Dir("./amelia"))
	r.Handle("/static/*", ameliaFileServer)

	return r
}

func ryanRouter() chi.Router {
	r := chi.NewRouter()

	ryanWrappedHandler := otelhttp.NewHandler(internal.RyanHandler(), "ryan-handler")
	r.Method("GET", "/", ryanWrappedHandler)
	r.Method("GET", "/visitors", ryanWrappedHandler)
	r.Method("GET", "/visitors.html", ryanWrappedHandler)

	r.NotFound(internal.NotFoundHandlerFunc)

	// Other static content
	ryanFileServer := http.FileServer(http.Dir("./ryan"))
	r.Handle("/static/*", ryanFileServer)

	return r
}

func sheldonRouter() chi.Router {
	r := chi.NewRouter()

	sheldonWrappedHandler := otelhttp.NewHandler(internal.SheldonHandler(), "sheldon-handler")
	r.Method("GET", "/", sheldonWrappedHandler)
	r.Method("GET", "/visitors", sheldonWrappedHandler)
	r.Method("GET", "/graph", sheldonWrappedHandler)
	r.Method("GET", "/resume", sheldonWrappedHandler)

	r.NotFound(internal.NotFoundHandlerFunc)

	// Other static content
	sheldonFileServer := http.FileServer(http.Dir("./sheldon"))
	r.Handle("/static/*", sheldonFileServer)

	return r
}

func domainRouter() chi.Router {
	r := chi.NewRouter()

	domainWrappedHandler := otelhttp.NewHandler(internal.DomainHandler(), "domain-handler")
	r.Method("GET", "/", domainWrappedHandler)
	r.Method("GET", "/visitors", domainWrappedHandler)
	r.Method("GET", "/visitors.html", domainWrappedHandler)

	r.NotFound(internal.NotFoundHandlerFunc)

	// Other static content
	domainFileServer := http.FileServer(http.Dir("./domain"))
	r.Handle("/static/*", domainFileServer)

	return r
}

func testVueRouter() chi.Router {
	r := chi.NewRouter()

	testVueWrappedHandler := otelhttp.NewHandler(internal.TestVueHandler(), "domain-handler")
	r.Method("GET", "/", testVueWrappedHandler)

	r.NotFound(internal.NotFoundHandlerFunc)

	// Other static content
	testVueFileServer := http.FileServer(http.Dir("./test-vue/dist"))
	r.Handle("/assets/*", testVueFileServer)
	r.Handle("/js/*", testVueFileServer)

	return r
}

func apiRouter() chi.Router {
	r := chi.NewRouter()

	graphApiWrappedHandler := otelhttp.NewHandler(internal.GraphApiHandler(), "api-handler")
	r.Method("GET", "/graph/{fredSeries}", graphApiWrappedHandler)

	r.NotFound(internal.ApiNotFoundHandler)

	return r
}

func notFoundRouter() chi.Router {
	r := chi.NewRouter()
	r.NotFound(internal.NotFoundHandlerFunc)

	// Other static content
	notFoundFileServer := http.FileServer(http.Dir("./common"))
	r.Handle("/static/*", notFoundFileServer)

	return r
}
