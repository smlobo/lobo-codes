package internal

import (
	"html/template"
	"log"
	"net/http"
)

type HandlerInfo struct {
	PathMap map[string]*template.Template
}

var HandlerInfoMap map[string]HandlerInfo

func GetPathMap(directory string) map[string]*template.Template {
	// Todo: iterate over all html files in directory
	pathMap := map[string]*template.Template{}
	pathMap["index"] = template.Must(template.ParseFiles(directory + "/index.html"))
	pathMap["visitors"] = template.Must(template.ParseFiles(directory + "/visitors.html"))
	return pathMap
}

func handleIndexHtml(directory string, writer http.ResponseWriter, request *http.Request) {
	url := request.URL
	key := ""
	if url.Path == "" || url.Path == "/" || url.Path == "/index.html" {
		// Log request to the sqlite3 db
		requestInfo(request, directory)

		key = "index"
		_ = HandlerInfoMap[directory].PathMap[key].Execute(writer, nil)
	}
}

func AmeliaHandler(writer http.ResponseWriter, request *http.Request) {
	//log.Printf("Amelia request: %s :: %s", request.RequestURI, request.URL.Path)

	handleIndexHtml("amelia", writer, request)

	url := request.URL
	if url.Path == "/visitors.html" || url.Path == "/visitors" {
		visitorPageData := VisitorsPage{}

		// Cassandra session
		session, err := CassandraCluster.CreateSession()
		if err != nil {
			log.Printf("WARNING: failed to create session with cassandra database: %s; %s",
				CassandraCluster.Hosts[0], err.Error())
			return
		}
		defer session.Close()

		tmpl := HandlerInfoMap["amelia"].PathMap["visitors"]

		//// Read the top 20 cities
		//_ = db.Table("request_infos").
		//	Select("city, region, country_short, count(city) as count").
		//	Group("city").
		//	Order("count desc").
		//	Limit(20).
		//	Find(&visitorPageData.Cities)

		// Read country name & count
		queryString := "SELECT country_short FROM amelia ALLOW FILTERING "
		scanner := session.Query(queryString).Iter().Scanner()
		countryCountMap := make(map[string]int)
		for scanner.Next() {
			var countryShort string
			err = scanner.Scan(&countryShort)
			if err != nil {
				continue
			}
			if count, ok := countryCountMap[countryShort]; !ok {
				countryCountMap[countryShort] = 1
			} else {
				countryCountMap[countryShort] = count + 1
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
		}

		_ = tmpl.Execute(writer, visitorPageData)
	}
}

func RyanHandler(writer http.ResponseWriter, request *http.Request) {
	handleIndexHtml("ryan", writer, request)
}

func SheldonHandler(writer http.ResponseWriter, request *http.Request) {
	handleIndexHtml("sheldon", writer, request)
}

func DomainHandler(writer http.ResponseWriter, request *http.Request) {
	handleIndexHtml("domain", writer, request)
}

var NotFoundTemplate *template.Template

func NotFoundHandler(writer http.ResponseWriter, request *http.Request) {
	_ = NotFoundTemplate.Execute(writer, nil)
}

type Country struct {
	CountryShort string
	//CountryLong  string
	Count int
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
}
