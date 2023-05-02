package internal

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"sort"
)

type HandlerInfo struct {
	PathMap map[string]*template.Template
}

var HandlerInfoMap map[string]HandlerInfo

func GetPathMap(directory string) map[string]*template.Template {
	// Todo: iterate over all html files in directory
	pathMap := map[string]*template.Template{}

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
		// Cassandra session
		session, err := CassandraCluster.CreateSession()
		if err != nil {
			log.Printf("WARNING: failed to create session with cassandra database: %s; %s",
				CassandraCluster.Hosts[0], err.Error())
			return
		}
		defer session.Close()

		tmpl := HandlerInfoMap[directory].PathMap["visitors"]

		// Read country name & count
		// Also, the city & region to count
		visitorPageData := VisitorsPage{}

		queryString := fmt.Sprintf("SELECT country_short, city, region FROM %s ALLOW FILTERING", directory)
		scanner := session.Query(queryString).Iter().Scanner()

		countryCountMap := make(map[string]int)
		cityCountMap := make(map[string]City)

		for scanner.Next() {
			var countryShort, city, region string
			err = scanner.Scan(&countryShort, &city, &region)
			if err != nil {
				continue
			}

			if count, ok := countryCountMap[countryShort]; !ok {
				countryCountMap[countryShort] = 1
			} else {
				countryCountMap[countryShort] = count + 1
			}

			if count, ok := cityCountMap[city]; !ok {
				cityCountMap[city] = City{
					City:         city,
					Region:       region,
					CountryShort: countryShort,
					Count:        1,
				}
			} else {
				count.Count += 1
				cityCountMap[city] = count
			}
		}

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
		visitorPageData.Cities = visitorPageData.Cities[:20]

		getpoweredBy(&visitorPageData.PoweredBy)

		_ = tmpl.Execute(writer, visitorPageData)
	}
}

func AmeliaHandler(writer http.ResponseWriter, request *http.Request) {
	directory := "amelia"
	handleIndexHtml(directory, writer, request)
	handleVisitorHtml(directory, writer, request)
}

func RyanHandler(writer http.ResponseWriter, request *http.Request) {
	directory := "ryan"
	handleIndexHtml(directory, writer, request)
	handleVisitorHtml(directory, writer, request)
}

func SheldonHandler(writer http.ResponseWriter, request *http.Request) {
	directory := "sheldon"
	handleIndexHtml(directory, writer, request)
	handleVisitorHtml(directory, writer, request)

	url := request.URL
	if url.Path == "/graph" {
		pageData := IndexPage{}
		getpoweredBy(&pageData.PoweredBy)
		_ = HandlerInfoMap[directory].PathMap["graph"].Execute(writer, pageData)
	}
}

func DomainHandler(writer http.ResponseWriter, request *http.Request) {
	directory := "domain"
	handleIndexHtml(directory, writer, request)
	handleVisitorHtml(directory, writer, request)
}

var NotFoundTemplate *template.Template

func NotFoundHandler(writer http.ResponseWriter, request *http.Request) {
	_ = NotFoundTemplate.Execute(writer, nil)
}

func getpoweredBy(poweredByPtr *PoweredBy) {
	poweredByPtr.GoVersion = os.Getenv("GOLANG_VERSION")
	poweredByPtr.PodName = os.Getenv("HOSTNAME")
	poweredByPtr.NodeName = os.Getenv("NODE_NAME")

	poweredByPtr.OsVersion = OsRelease.VersionId

	poweredByPtr.KubernetesVersion = Kubernetes.Version
}

type Country struct {
	CountryShort string
	Count        int
}

type City struct {
	City         string
	Region       string
	CountryShort string
	Count        int
}

type VisitorsPage struct {
	UniqueCountries int
	Countries       []Country
	Cities          []City
	PoweredBy       PoweredBy
}

type IndexPage struct {
	PoweredBy PoweredBy
}

type PoweredBy struct {
	GoVersion         string
	KubernetesVersion string
	OsVersion         string
	PodName           string
	NodeName          string
}
