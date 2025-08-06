package main

import (
	"context"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"k8s-http-server/hostrouter"
	"k8s-http-server/internal"
)

func main() {
	// Input arguments
	sslKeyDirPtr := flag.String("ssl-key-dir", "./ssl-certificates",
		"directory with SSL key and cert files")
	logPtr := flag.String("log", "stdout", "log file (stdout)")
	testOsReleasePtr := flag.String("os-release", "/etc/os-release", "os-release file")
	testGoVersionPtr := flag.String("go-version", "/app/golang_version.txt", "go_version.txt file")

	flag.Parse()

	// Logging
	if *logPtr != "stdout" {
		// If the file doesn't exist, create it or append to the file
		file, err := os.OpenFile(*logPtr, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatal(err)
		}
		log.SetOutput(file)
	}

	// config (currently traces only; TODO all above flags)
	internal.SetupConfig()

	// Ip2location db file init
	geoDbFile := "IP2LOCATION-LITE-DB11.IPV6.BIN"
	if err := internal.InitGeoDB(geoDbFile); err != nil {
		log.Printf("Failed opening ip2location db file %s : %s", geoDbFile, err)
		os.Exit(1)
	}

	// Parse os-release
	internal.InitOsRelease(*testOsReleasePtr)

	// Get kubernetes version
	internal.InitKubernetesInfo()

	// Get Go version
	internal.InitGoVersion(*testGoVersionPtr)

	log.Printf("Platform info: OS%s; K8s%s; Go<%s>", internal.OsRelease, internal.Kubernetes,
		internal.GoVersion)

	// Subdomains served
	internal.HandlerInfoMap = map[string]internal.HandlerInfo{
		"amelia":   {},
		"ryan":     {},
		"bliu":     {},
		"sheldon":  {},
		"domain":   {},
		"test-vue": {},
		"wasm":     {},
		"hikes":    {},
	}

	// Init rqlite
	internal.InitRqlite()

	// Init cassandra db
	//cassandraSession := internal.InitCassandra()
	//defer cassandraSession.Close()

	// Locate and parse all html (with/without templates)
	for subdomain, handlerInfo := range internal.HandlerInfoMap {
		handlerInfo.PathMap = internal.GetPathMap(subdomain)
		internal.HandlerInfoMap[subdomain] = handlerInfo
	}

	// Common template
	internal.VisitorTemplate = template.Must(template.ParseFiles("common/visitors.html"))
	internal.NotFoundTemplate = template.Must(template.ParseFiles("common/notfound.html"))
	internal.FooterTemplate = template.Must(template.ParseFiles("common/footer.html"))
	internal.HikesNotFoundTemplate = template.Must(template.ParseFiles("hikes/public/404.html"))

	// Setup OpenTelemetry tracer
	tp := internal.InitTracerProvider("lobo-codes")
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

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

	port := ":" + internal.Config["HTTP_PORT"]
	sslPort := ":" + internal.Config["HTTPS_PORT"]

	hostRouter.Map("amelia.lobo.codes", internal.GenericRouter("amelia"))
	hostRouter.Map("ryan.lobo.codes", internal.GenericRouter("ryan"))
	hostRouter.Map("bliu.lobo.codes", internal.GenericRouter("bliu"))
	hostRouter.Map("sheldon.lobo.codes", internal.GenericRouter("sheldon"))
	hostRouter.Map("lobo.codes", internal.GenericRouter("domain"))
	hostRouter.Map("test-vue.lobo.codes", internal.TestVueRouter())
	hostRouter.Map("api.lobo.codes", internal.ApiRouter())
	hostRouter.Map("wasm.lobo.codes", internal.WasmRouter())
	hostRouter.Map("hikes.lobo.codes", internal.HikesRouter())

	// Testing locally
	hostRouter.Map("amelia.lobo.codes"+port, internal.GenericRouter("amelia"))
	hostRouter.Map("ryan.lobo.codes"+port, internal.GenericRouter("ryan"))
	hostRouter.Map("bliu.lobo.codes"+port, internal.GenericRouter("bliu"))
	hostRouter.Map("sheldon.lobo.codes"+port, internal.GenericRouter("sheldon"))
	hostRouter.Map("lobo.codes"+port, internal.GenericRouter("domain"))
	hostRouter.Map("test-vue.lobo.codes"+port, internal.TestVueRouter())
	hostRouter.Map("api.lobo.codes"+port, internal.ApiRouter())
	hostRouter.Map("wasm.lobo.codes"+port, internal.WasmRouter())
	hostRouter.Map("hikes.lobo.codes"+port, internal.HikesRouter())

	hostRouter.Map("*", internal.NotFoundRouter())
	router.Mount("/", hostRouter)

	// Wait for both http & https servers to finish
	var serversWaitGroup sync.WaitGroup
	serversWaitGroup.Add(2)

	go func() {
		log.Printf("Listening on port %s ...", port)
		err := http.ListenAndServe(port, router)
		if err != nil {
			log.Printf("Error serving on port %s : %s", port, err)
		}
		serversWaitGroup.Done()
	}()

	go func() {
		log.Printf("Listening on SSL port %s ...", sslPort)
		err := http.ListenAndServeTLS(sslPort, *sslKeyDirPtr+"/fullchain.pem",
			*sslKeyDirPtr+"/privkey.pem", router)
		if err != nil {
			log.Printf("Error serving with SSL on port %s : %s", sslPort, err)
		}
		serversWaitGroup.Done()
	}()

	serversWaitGroup.Wait()
}
