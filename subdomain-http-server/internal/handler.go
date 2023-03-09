package internal

import (
	"gorm.io/gorm"
	"html/template"
	"net/http"
)

type HandlerInfo struct {
	PathMap    map[string]*template.Template
	RequestsDb *gorm.DB
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
		requestInfo(request, HandlerInfoMap[directory].RequestsDb)

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

		db := HandlerInfoMap["amelia"].RequestsDb
		tmpl := HandlerInfoMap["amelia"].PathMap["visitors"]
		// Read country name and count
		result := db.Table("request_infos").
			Select("country_short, country_long, count(country_short) as count").
			Group("country_short").
			Order("count desc").
			Find(&visitorPageData.Countries)

		// Read the top 20 cities
		_ = db.Table("request_infos").
			Select("city, region, country_short, count(city) as count").
			Group("city").
			Order("count desc").
			Limit(20).
			Find(&visitorPageData.Cities)

		visitorPageData.UniqueCountries = result.RowsAffected
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
	CountryLong  string
	Count        int
}

type City struct {
	City         string
	Region       string
	CountryShort string
	Count        int
}

type VisitorsPage struct {
	UniqueCountries int64
	Countries       []Country
	Cities          []City
}
