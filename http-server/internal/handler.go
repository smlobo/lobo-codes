package internal

import (
	"gorm.io/gorm"
	"html/template"
	"log"
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

		var info RequestInfo
		requestInfo(request, &info)
		//log.Printf("IP: %s", info.RemoteAddress)
		db.Create(&info)
		log.Printf("RequestInfo: %s", info)

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
