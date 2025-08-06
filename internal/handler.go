package internal

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
)

type HandlerInfo struct {
	PathMap map[string]*template.Template
}

var HandlerInfoMap map[string]HandlerInfo
var VisitorTemplate *template.Template
var NotFoundTemplate *template.Template
var FooterTemplate *template.Template
var HikesNotFoundTemplate *template.Template

func GetPathMap(directory string) map[string]*template.Template {
	log.Printf("Templating: %s", directory)

	pathMap := map[string]*template.Template{}

	// Domain/sub-domain specific html templates
	// Some are in a different path
	searchDir := directory
	switch directory {
	case "test-vue":
		searchDir = directory + "/dist"
	case "hikes":
		searchDir = directory + "/public"
	default:
	}

	// Add all html files to pathMap
	err := filepath.WalkDir(searchDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".html") {
			return nil
		}
		fileKey := path[len(searchDir)+len("/") : len(path)-len(".html")]
		log.Printf("  ~ %s -> %s", path, fileKey)
		pathMap[fileKey] = template.Must(template.ParseFiles(path))
		return nil
	})
	if err != nil {
		log.Fatalf("Error walking dir: %s, %v", searchDir, err)
	}

	return pathMap
}

func VisitorsHandlerGenerator(dbTable string) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		// Get Visitor page data for template processing
		visitorPageData := VisitorsPage{}

		countryCountMap, cityCountMap := rqliteGetCountriesCities(dbTable, request)

		// HTML Template processing
		_, span := otel.Tracer("k8s-http-server").Start(request.Context(), "template-process")
		defer span.End()

		visitorPageData.UniqueCountries = len(countryCountMap)

		visitorPageData.Countries = make([]Country, visitorPageData.UniqueCountries)
		index := 0
		for country, count := range countryCountMap {
			visitorPageData.Countries[index] = Country{
				CountryShort: country,
				Count:        count,
			}
			index++
		}

		visitorPageData.Cities = make([]City, len(cityCountMap))
		index = 0
		for _, city := range cityCountMap {
			visitorPageData.Cities[index] = city
			index++
		}
		sort.Slice(visitorPageData.Cities, func(i, j int) bool {
			return visitorPageData.Cities[i].Count > visitorPageData.Cities[j].Count
		})
		if len(cityCountMap) > 20 {
			visitorPageData.Cities = visitorPageData.Cities[:20]
		}

		getpoweredBy(&visitorPageData.PoweredBy)
		_ = VisitorTemplate.Execute(writer, visitorPageData)
	}
}

func GenericHandlerGenerator(subDomain string) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		pathKey := strings.TrimSuffix(strings.TrimLeft(request.URL.Path, "/"), ".html")
		if pathKey == "" || pathKey[len(pathKey)-1] == '/' {
			pathKey += "index"
		}
		log.Printf("handle any html: [%s] %s", subDomain, pathKey)

		// Only log requests to index.html
		if pathKey == "index" {
			// Log request to the rqlite db
			requestInfo(request, subDomain)
		}

		indexPageData := IndexPage{}
		getpoweredBy(&indexPageData.PoweredBy)
		_ = HandlerInfoMap[subDomain].PathMap[pathKey].Execute(writer, indexPageData)
	}
}

func CommonHandlerGenerator(template *template.Template) http.HandlerFunc {
	return func(writer http.ResponseWriter, _ *http.Request) {
		poweredByData := IndexPage{}
		getpoweredBy(&poweredByData.PoweredBy)
		_ = template.Execute(writer, poweredByData)
	}
}

func CommonNotFoundHandler(writer http.ResponseWriter, request *http.Request) {
	poweredByData := IndexPage{}
	getpoweredBy(&poweredByData.PoweredBy)
	_ = NotFoundTemplate.Execute(writer, poweredByData)

}

func CommonNotAllowedHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		CommonNotFoundHandler(writer, request)
		return
	}

	writer.WriteHeader(http.StatusMethodNotAllowed)

	responseString := fmt.Sprintf("Method Not Allowed: %s\n", request.Method)

	// Headers
	writer.Header().Add("Content-Length", strconv.Itoa(len(responseString)))
	writer.Header().Set("Content-Type", "text/plain")
	writer.Header().Add("Last-Modified", time.Now().UTC().Format(http.TimeFormat))

	_, _ = writer.Write([]byte(responseString))
}

func HikesNotFoundHandler(writer http.ResponseWriter, request *http.Request) {
	poweredByData := IndexPage{}
	getpoweredBy(&poweredByData.PoweredBy)
	_ = HikesNotFoundTemplate.Execute(writer, poweredByData)

}

func HikesNotAllowedHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		HikesNotFoundHandler(writer, request)
		return
	}

	writer.WriteHeader(http.StatusMethodNotAllowed)

	responseString := fmt.Sprintf("Method Not Allowed: %s\n", request.Method)

	// Headers
	writer.Header().Add("Content-Length", strconv.Itoa(len(responseString)))
	writer.Header().Set("Content-Type", "text/plain")
	writer.Header().Add("Last-Modified", time.Now().UTC().Format(http.TimeFormat))

	_, _ = writer.Write([]byte(responseString))
}

func HeadHandler(writer http.ResponseWriter, _ *http.Request) {
	writer.Header().Add("Content-Length", "0")
	writer.Header().Add("Content-Type", "text/plain")
	writer.Header().Add("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
}

func getpoweredBy(poweredByPtr *PoweredBy) {
	poweredByPtr.PodName = os.Getenv("HOSTNAME")
	poweredByPtr.NodeName = os.Getenv("NODE_NAME")

	poweredByPtr.OsVersion = OsRelease.VersionId

	poweredByPtr.KubernetesVersion = Kubernetes.Version

	poweredByPtr.GoVersion = string(GoVersion)

	poweredByPtr.RqliteVersion = RqliteVersion
}
