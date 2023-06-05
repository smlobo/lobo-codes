package internal

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"io"
	"log"
	"net/http"
)

const FRED_API = "https://api.stlouisfed.org/fred/series"

const API_KEY = "e12c0787f51ff6db24ac8029710fa175"

func GraphApiHandler(writer http.ResponseWriter, request *http.Request) {

	writer.Header().Set("Content-Type", "application/json")

	fredSeries := chi.URLParam(request, "fredSeries")
	if fredSeries == "" {
		writer.WriteHeader(http.StatusNoContent)
		return
	}

	log.Printf("Graph request for fredSeries: %s", fredSeries)

	type seriesData struct {
		Date  string
		Value string
	}
	type graphData struct {
		Title string
		Units string
		Data  []seriesData
	}
	var returnData graphData

	// Get the title
	response, err := http.Get(FRED_API + "?file_type=json&series_id=" + fredSeries + "&api_key=" + API_KEY)
	if err != nil || response.StatusCode > http.StatusOK {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// Unmarshal
	var metaData map[string]interface{}
	err = json.Unmarshal(body, &metaData)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	series0 := metaData["seriess"].([]interface{})[0]
	seriesMap := series0.(map[string]interface{})
	returnData.Title = seriesMap["title"].(string)
	returnData.Units = seriesMap["units"].(string)
	returnData.Data = make([]seriesData, 0, 100)

	log.Printf("Got: %d : %s", response.StatusCode, returnData.Title)

	// Get the data
	response, err = http.Get(FRED_API + "/observations?file_type=json&series_id=" + fredSeries + "&api_key=" + API_KEY)
	if err != nil || response.StatusCode > http.StatusOK {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	body, err = io.ReadAll(response.Body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// Unmarshal
	err = json.Unmarshal(body, &metaData)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	observations := metaData["observations"].([]interface{})
	for _, observation := range observations {
		observationMap := observation.(map[string]interface{})
		returnData.Data = append(returnData.Data, seriesData{
			Date:  observationMap["date"].(string),
			Value: observationMap["value"].(string),
		})
	}

	encodedReturnData, _ := json.Marshal(returnData)
	writer.Write(encodedReturnData)
}

func ApiNotFoundHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
}
