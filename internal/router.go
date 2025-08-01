package internal

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func GenericRouter(subDomain string) chi.Router {
	r := chi.NewRouter()

	// Iterate over all static html pages adding them to the pattern
	log.Printf("Router patterns for: %s", subDomain)
	for pathKey, _ := range HandlerInfoMap[subDomain].PathMap {
		log.Printf("  * %s", pathKey)
		r.Get("/"+pathKey, GenericHandlerGenerator(subDomain))
		r.Get("/"+pathKey+".html", GenericHandlerGenerator(subDomain))
	}
	// Default "/" path
	r.Get("/", GenericHandlerGenerator(subDomain))

	// Footer frame
	r.Get("/footer.html", CommonHandlerGenerator(FooterTemplate))

	// Visitor page
	r.Get("/visitors", VisitorsHandlerGenerator(subDomain))
	r.Get("/visitors.html", VisitorsHandlerGenerator(subDomain))

	// Pattern not found
	r.NotFound(CommonHandlerGenerator(NotFoundTemplate))
	r.MethodNotAllowed(CommonHandlerGenerator(NotFoundTemplate))

	// Other static content
	subDomainFileServer := http.FileServer(http.Dir("./" + subDomain))
	r.Handle("/images/*", subDomainFileServer)
	r.Handle("/assets/*", subDomainFileServer)
	r.Handle("/static/*", subDomainFileServer)

	// HEAD requests
	r.Head("/*", HeadHandler)

	return r
}

func TestVueRouter() chi.Router {
	r := chi.NewRouter()

	r.Get("/", GenericHandlerGenerator("test-vue"))
	r.Get("/footer.html", CommonHandlerGenerator(FooterTemplate))

	// Visitor page
	r.Get("/visitors", VisitorsHandlerGenerator("test-vue"))
	r.Get("/visitors.html", VisitorsHandlerGenerator("test-vue"))

	// Pattern not found
	r.NotFound(CommonHandlerGenerator(NotFoundTemplate))
	r.MethodNotAllowed(CommonHandlerGenerator(NotFoundTemplate))

	// Other static content
	testVueFileServer := http.FileServer(http.Dir("./test-vue"))
	r.Handle("/static/*", testVueFileServer)

	// Vue distribution static content
	testVueDistFileServer := http.FileServer(http.Dir("./test-vue/dist"))
	r.Handle("/*", testVueDistFileServer)

	// HEAD requests
	r.Head("/*", HeadHandler)

	return r
}

func ApiRouter() chi.Router {
	r := chi.NewRouter()

	graphApiWrappedHandler := otelhttp.NewHandler(GraphApiHandler(), "api-handler")
	r.Method("GET", "/graph/fred/{fredSeries}", graphApiWrappedHandler)
	r.Method("GET", "/graph/census/{censusSeries}", graphApiWrappedHandler)

	r.NotFound(ApiNotFoundHandler)

	// HEAD requests
	r.Head("/*", HeadHandler)

	return r
}

func WasmRouter() chi.Router {
	r := chi.NewRouter()

	r.Get("/", GenericHandlerGenerator("wasm"))
	r.Get("/footer.html", CommonHandlerGenerator(FooterTemplate))

	// Visitor page
	r.Get("/visitors", VisitorsHandlerGenerator("wasm"))
	r.Get("/visitors.html", VisitorsHandlerGenerator("wasm"))

	// Pattern not found
	r.NotFound(CommonHandlerGenerator(NotFoundTemplate))
	r.MethodNotAllowed(CommonHandlerGenerator(NotFoundTemplate))

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

	// HEAD requests
	r.Head("/*", HeadHandler)

	return r
}

func HikesRouter() chi.Router {
	r := chi.NewRouter()

	//r.NotFound(CommonHandlerGenerator(NotFoundTemplate))
	//r.MethodNotAllowed(CommonHandlerGenerator(NotFoundTemplate))

	// All static content
	hikesFileServer := http.FileServer(http.Dir("./hikes/public"))
	r.Handle("/*", hikesFileServer)

	return r
}

func NotFoundRouter() chi.Router {
	r := chi.NewRouter()

	r.NotFound(CommonHandlerGenerator(NotFoundTemplate))

	// Other static content
	notFoundFileServer := http.FileServer(http.Dir("./common"))
	r.Handle("/static/*", notFoundFileServer)

	return r
}
