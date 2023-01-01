package internal

import (
	"gorm.io/gorm"
	"html/template"
	"net/http"
)

func ReturnHandler(directory string, db *gorm.DB) http.HandlerFunc {
	pathMap := map[string]*template.Template{}

	// Initialize with the index.html file
	pathMap["index.html"] = template.Must(template.ParseFiles(directory + "/index.html"))
	// Also, not found
	pathMap["notfound.html"] = template.Must(template.ParseFiles("static/notfound.html"))

	return func(writer http.ResponseWriter, request *http.Request) {
		//log.Printf("Got request from: %s, %s, %s", request.RemoteAddr, request.Header.Get("User-Agent"),
		//	request.RequestURI)

		// Log request to the sqlite3 db
		requestInfo(request, db)

		url := request.URL
		key := ""
		if url.Path == "" || url.Path == "/" {
			key = "index.html"
		} else {
			key = "notfound.html"
		}
		err := pathMap[key].Execute(writer, nil)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
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

func TemplateHandler(tmpl *template.Template, db *gorm.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		visitorPageData := VisitorsPage{}

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
		_ = tmpl.Execute(w, visitorPageData)
	}
}
