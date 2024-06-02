package internal

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	//"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

const FredApi = "https://api.stlouisfed.org/fred/series"

const ApiKey = "e12c0787f51ff6db24ac8029710fa175"

func GraphApiHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

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
		returnData.Data = make([]seriesData, 0, 100)

		// Concurrently get data
		titleChan := make(chan string)
		unitsChan := make(chan string)

		// Get the title & units
		go func() {
			//response, err := otelhttp.Get(request.Context(), FredApi+"?file_type=json&series_id="+fredSeries+
			//	"&api_key="+ApiKey)
			response, err := http.Get(FredApi + "?file_type=json&series_id=" + fredSeries + "&api_key=" + ApiKey)
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
			titleChan <- seriesMap["title"].(string)
			unitsChan <- seriesMap["units"].(string)
		}()

		// Get the data
		//response, err := otelhttp.Get(request.Context(), FredApi+"/observations?file_type=json&series_id="+
		//	fredSeries+"&api_key="+ApiKey)
		response, err := http.Get(FredApi + "/observations?file_type=json&series_id=" + fredSeries + "&api_key=" + ApiKey)
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
		observations := metaData["observations"].([]interface{})
		for _, observation := range observations {
			observationMap := observation.(map[string]interface{})
			returnData.Data = append(returnData.Data, seriesData{
				Date:  observationMap["date"].(string),
				Value: observationMap["value"].(string),
			})
		}

		// Block on the earlier request
		returnData.Title = <-titleChan
		returnData.Units = <-unitsChan
		log.Printf("Got: %d : %s", response.StatusCode, returnData.Title)

		encodedReturnData, _ := json.Marshal(returnData)
		writer.Write(encodedReturnData)
	}
}

func ApiNotFoundHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
}
