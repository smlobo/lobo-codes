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
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
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
	}

	// Init rqlite
	internal.InitRqlite()

	// Init cassandra db
	//cassandraSession := internal.InitCassandra()
	//defer cassandraSession.Close()

	// Init cassandra table for each subdomain
	// Also, save and parse html templates
	for subdomain, handlerInfo := range internal.HandlerInfoMap {
		handlerInfo.PathMap = internal.GetPathMap(subdomain)

		internal.HandlerInfoMap[subdomain] = handlerInfo
	}
	// Other HTML templates
	internal.HandlerInfoMap["sheldon"].PathMap["graph"] = template.Must(template.ParseFiles("sheldon/graph.html"))
	internal.HandlerInfoMap["sheldon"].PathMap["resume"] = template.Must(template.ParseFiles("sheldon/resume.html"))
	internal.HandlerInfoMap["bliu"].PathMap["generic"] = template.Must(template.ParseFiles("bliu/generic.html"))
	//internal.HandlerInfoMap["bliu"].PathMap["elements"] = template.Must(template.ParseFiles("bliu/elements.html"))

	// Setup OpenTelemetry tracer
	tp := internal.InitTracerProvider("lobo-codes")
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	// Not found template
	internal.NotFoundTemplate = template.Must(template.ParseFiles("common/notfound.html"))

	// Footer template
	internal.FooterTemplate = template.Must(template.ParseFiles("common/footer.html"))

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

	hostRouter.Map("amelia.lobo.codes", ameliaRouter())
	hostRouter.Map("ryan.lobo.codes", ryanRouter())
	hostRouter.Map("bliu.lobo.codes", bliuRouter())
	hostRouter.Map("sheldon.lobo.codes", sheldonRouter())
	hostRouter.Map("lobo.codes", domainRouter())
	hostRouter.Map("test-vue.lobo.codes", testVueRouter())
	hostRouter.Map("api.lobo.codes", apiRouter())
	hostRouter.Map("wasm.lobo.codes", wasmRouter())
	hostRouter.Map("hikes.lobo.codes", hikesRouter())

	// Testing locally
	hostRouter.Map("amelia.lobo.codes"+port, ameliaRouter())
	hostRouter.Map("ryan.lobo.codes"+port, ryanRouter())
	hostRouter.Map("bliu.lobo.codes"+port, bliuRouter())
	hostRouter.Map("sheldon.lobo.codes"+port, sheldonRouter())
	hostRouter.Map("lobo.codes"+port, domainRouter())
	hostRouter.Map("test-vue.lobo.codes"+port, testVueRouter())
	hostRouter.Map("api.lobo.codes"+port, apiRouter())
	hostRouter.Map("wasm.lobo.codes"+port, wasmRouter())
	hostRouter.Map("hikes.lobo.codes"+port, hikesRouter())

	hostRouter.Map("*", notFoundRouter())
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

func ameliaRouter() chi.Router {
	r := chi.NewRouter()

	ameliaWrappedHandler := otelhttp.NewHandler(internal.AmeliaHandler(), "amelia-handler")
	r.Method("GET", "/", ameliaWrappedHandler)
	r.Method("GET", "/visitors", ameliaWrappedHandler)
	r.Method("GET", "/visitors.html", ameliaWrappedHandler)
	r.Get("/footer.html", internal.FooterHandlerFunc)
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
	r.Get("/footer.html", internal.FooterHandlerFunc)
	r.NotFound(internal.NotFoundHandlerFunc)

	// Other static content
	ryanFileServer := http.FileServer(http.Dir("./ryan"))
	r.Handle("/static/*", ryanFileServer)

	return r
}

func bliuRouter() chi.Router {
	r := chi.NewRouter()

	bliuWrappedHandler := otelhttp.NewHandler(internal.BliuHandler(), "bliu-handler")
	r.Method("GET", "/", bliuWrappedHandler)
	r.Method("GET", "/index.html", bliuWrappedHandler)
	r.Method("GET", "/generic.html", bliuWrappedHandler)
	r.Method("GET", "/visitors", bliuWrappedHandler)
	r.Method("GET", "/visitors.html", bliuWrappedHandler)
	//r.Method("GET", "/elements.html", bliuWrappedHandler)
	r.Get("/footer.html", internal.FooterHandlerFunc)
	r.NotFound(internal.NotFoundHandlerFunc)

	// Other static content
	bliuFileServer := http.FileServer(http.Dir("./bliu"))
	r.Handle("/images/*", bliuFileServer)
	r.Handle("/assets/*", bliuFileServer)
	r.Handle("/static/*", bliuFileServer)

	return r
}

func sheldonRouter() chi.Router {
	r := chi.NewRouter()

	sheldonWrappedHandler := otelhttp.NewHandler(internal.SheldonHandler(), "sheldon-handler")
	r.Method("GET", "/", sheldonWrappedHandler)
	r.Method("GET", "/visitors", sheldonWrappedHandler)
	r.Method("GET", "/graph", sheldonWrappedHandler)
	r.Method("GET", "/resume", sheldonWrappedHandler)
	r.Get("/footer.html", internal.FooterHandlerFunc)
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
	r.Get("/footer.html", internal.FooterHandlerFunc)
	r.NotFound(internal.NotFoundHandlerFunc)

	// Other static content
	domainFileServer := http.FileServer(http.Dir("./domain"))
	r.Handle("/static/*", domainFileServer)

	return r
}

func testVueRouter() chi.Router {
	r := chi.NewRouter()

	testVueWrappedHandler := otelhttp.NewHandler(internal.TestVueHandler(), "test-vue-handler")
	r.Method("GET", "/", testVueWrappedHandler)
	r.Method("GET", "/visitors", testVueWrappedHandler)
	r.Method("GET", "/visitors.html", testVueWrappedHandler)
	r.Get("/footer.html", internal.FooterHandlerFunc)
	r.NotFound(internal.NotFoundHandlerFunc)

	// Other static content
	testVueFileServer := http.FileServer(http.Dir("./test-vue"))
	r.Handle("/static/*", testVueFileServer)

	// Vue distribution static content
	testVueDistFileServer := http.FileServer(http.Dir("./test-vue/dist"))
	r.Handle("/assets/*", testVueDistFileServer)
	r.Handle("/js/*", testVueDistFileServer)

	return r
}

func apiRouter() chi.Router {
	r := chi.NewRouter()

	graphApiWrappedHandler := otelhttp.NewHandler(internal.GraphApiHandler(), "api-handler")
	r.Method("GET", "/graph/fred/{fredSeries}", graphApiWrappedHandler)
	r.Method("GET", "/graph/census/{censusSeries}", graphApiWrappedHandler)

	r.NotFound(internal.ApiNotFoundHandler)

	return r
}

func wasmRouter() chi.Router {
	r := chi.NewRouter()

	wasmWrappedHandler := otelhttp.NewHandler(internal.WasmHandler(), "wasm-handler")
	r.Method("GET", "/", wasmWrappedHandler)
	r.Method("GET", "/visitors", wasmWrappedHandler)
	r.Method("GET", "/visitors.html", wasmWrappedHandler)
	r.Get("/footer.html", internal.FooterHandlerFunc)
	r.NotFound(internal.NotFoundHandlerFunc)

	// Other static content
	wasmFileServer := http.FileServer(http.Dir("./wasm"))
	r.Handle("/static/*", wasmFileServer)

	// wasm-rotating-cube
	wasmRotatingCubeFileServer := http.FileServer(http.Dir("./wasm/wasm-rotating-cube/dist"))
	r.Handle("/wasm-rotating-cube/*", wasmRotatingCubeFileServer)

	// h-tree
	wasmHTreeFileServer := http.FileServer(http.Dir("./wasm/h-tree/dist"))
	r.Handle("/h-tree/*", wasmHTreeFileServer)

	// fractal-circle
	wasmFractalCircleFileServer := http.FileServer(http.Dir("./wasm/fractal-circle/dist"))
	r.Handle("/fractal-circle/*", wasmFractalCircleFileServer)

	// julia-set
	wasmJuliaSetFileServer := http.FileServer(http.Dir("./wasm/julia-set/dist"))
	r.Handle("/julia-set/*", wasmJuliaSetFileServer)

	// collision-system
	wasmCollisionSystemFileServer := http.FileServer(http.Dir("./wasm/collision-system/dist"))
	r.Handle("/collision-system/*", wasmCollisionSystemFileServer)

	// shortest-path
	wasmShortestPathFileServer := http.FileServer(http.Dir("./wasm/shortest-path/dist"))
	r.Handle("/shortest-path/*", wasmShortestPathFileServer)

	return r
}

func hikesRouter() chi.Router {
	r := chi.NewRouter()

	r.NotFound(internal.NotFoundHandlerFunc)

	// All static content
	hikesFileServer := http.FileServer(http.Dir("./hikes"))
	r.Handle("/*", hikesFileServer)

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
