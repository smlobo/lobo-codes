package main

import (
	"html/template"
	"log"
	"net/http"
)

func main() {
	tmpl := template.Must(template.ParseFiles("static/index.html"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Log the request
		log.Printf("Got request from: %s, %s\n", r.RemoteAddr, r.Header.Get("User-Agent"))
		tmpl.Execute(w, nil)
	})
	log.Print("Listening on :5000...")
	http.ListenAndServe(":5000", nil)
}
