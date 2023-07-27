package internal

import (
	"html/template"
	"net/http"
	"os"
	"sort"

	"go.opentelemetry.io/otel"
)

type HandlerInfo struct {
	PathMap map[string]*template.Template
}

var HandlerInfoMap map[string]HandlerInfo

func GetPathMap(directory string) map[string]*template.Template {
	// Todo: iterate over all html files in directory
	pathMap := map[string]*template.Template{}

	// Special case test-vue
	if directory == "test-vue" {
		pathMap["index"] = template.Must(template.ParseFiles(directory + "/dist/index.html"))
		return pathMap
	}

	// Domain/sub-domain specific html templates
	pathMap["index"] = template.Must(template.ParseFiles(directory + "/index.html"))

	// Common html templates
	pathMap["visitors"] = template.Must(template.ParseFiles("common/visitors.html"))

	return pathMap
}

func handleIndexHtml(directory string, writer http.ResponseWriter, request *http.Request) {
	url := request.URL
	if url.Path == "" || url.Path == "/" || url.Path == "/index.html" {
		// Log request to the Cassandra db
		requestInfo(request, directory)

		indexPageData := IndexPage{}
		getpoweredBy(&indexPageData.PoweredBy)
		_ = HandlerInfoMap[directory].PathMap["index"].Execute(writer, indexPageData)
	}
}

func handleVisitorHtml(directory string, writer http.ResponseWriter, request *http.Request) {
	url := request.URL
	if url.Path == "/visitors" {

		// Get Visitor page data for template processing
		visitorPageData := VisitorsPage{}

		countryCountMap, cityCountMap := cassandraGetCountriesCities(directory, request)

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

		tmpl := HandlerInfoMap[directory].PathMap["visitors"]
		_ = tmpl.Execute(writer, visitorPageData)
	}
}

func AmeliaHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		directory := "amelia"
		handleIndexHtml(directory, writer, request)
		handleVisitorHtml(directory, writer, request)
	}
}

func RyanHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		directory := "ryan"
		handleIndexHtml(directory, writer, request)
		handleVisitorHtml(directory, writer, request)
	}
}

func SheldonHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		directory := "sheldon"
		handleIndexHtml(directory, writer, request)
		handleVisitorHtml(directory, writer, request)

		url := request.URL
		if url.Path == "/graph" {
			pageData := IndexPage{}
			getpoweredBy(&pageData.PoweredBy)
			_ = HandlerInfoMap[directory].PathMap["graph"].Execute(writer, pageData)
		} else if url.Path == "/resume" {
			pageData := IndexPage{}
			getpoweredBy(&pageData.PoweredBy)
			_ = HandlerInfoMap[directory].PathMap["resume"].Execute(writer, pageData)
		}
	}
}

func DomainHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		directory := "domain"
		handleIndexHtml(directory, writer, request)
		handleVisitorHtml(directory, writer, request)
	}
}

func TestVueHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		directory := "test-vue"
		url := request.URL
		if url.Path == "" || url.Path == "/" {
			// Todo: Log request to the Cassandra db
			//requestInfo(request, directory)

			//indexPageData := IndexPage{}
			//getpoweredBy(&indexPageData.PoweredBy)
			_ = HandlerInfoMap[directory].PathMap["index"].Execute(writer, nil)
		}
	}
}

var NotFoundTemplate *template.Template

func NotFoundHandlerFunc(writer http.ResponseWriter, request *http.Request) {
	_ = NotFoundTemplate.Execute(writer, nil)
}

func getpoweredBy(poweredByPtr *PoweredBy) {
	poweredByPtr.PodName = os.Getenv("HOSTNAME")
	poweredByPtr.NodeName = os.Getenv("NODE_NAME")

	poweredByPtr.OsVersion = OsRelease.VersionId

	poweredByPtr.KubernetesVersion = Kubernetes.Version

	poweredByPtr.GoVersion = string(GoVersion)
}
