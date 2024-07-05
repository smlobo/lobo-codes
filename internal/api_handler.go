package internal

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

const FredApi = "https://api.stlouisfed.org/fred/series"
const CensusApi = "https://api.census.gov/data/timeseries/idb/5year"

const FredApiKey = "e12c0787f51ff6db24ac8029710fa175"
const CensusApiKey = "fee00b0d52c7cd8f170a09ce4785218a70523396"

type seriesData struct {
	Date  string
	Value string
}
type graphData struct {
	Title string
	Units string
	Data  []seriesData
}

func fredData(fredSeries string, writer http.ResponseWriter, request *http.Request) (graphData, bool) {
	log.Printf("Graph request for fredSeries: %s", fredSeries)

	var returnData graphData
	returnData.Data = make([]seriesData, 0, 100)

	// Concurrently get data
	statusChan := make(chan bool, 1)
	titleChan := make(chan string, 1)
	unitsChan := make(chan string, 1)

	// Get the title & units
	go func() {
		response, err := otelhttp.Get(request.Context(), FredApi+"?file_type=json&series_id="+fredSeries+
			"&api_key="+FredApiKey)
		//response, err := http.Get(FredApi + "?file_type=json&series_id=" + fredSeries + "&api_key=" + FredApiKey)
		if err != nil || response.StatusCode > http.StatusOK {
			writer.WriteHeader(http.StatusBadRequest)
			statusChan <- false
			return
		}
		body, err := io.ReadAll(response.Body)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			statusChan <- false
			return
		}

		// Unmarshal
		var metaData map[string]interface{}
		err = json.Unmarshal(body, &metaData)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			statusChan <- false
			return
		}
		series0 := metaData["seriess"].([]interface{})[0]
		seriesMap := series0.(map[string]interface{})
		titleChan <- seriesMap["title"].(string)
		unitsChan <- seriesMap["units"].(string)
		statusChan <- true
	}()

	// Get the data
	//response, err := otelhttp.Get(request.Context(), FredApi+"/observations?file_type=json&series_id="+
	//	fredSeries+"&api_key="+FredApiKey)
	response, err := http.Get(FredApi + "/observations?file_type=json&series_id=" + fredSeries + "&api_key=" + FredApiKey)
	if err != nil || response.StatusCode > http.StatusOK {
		writer.WriteHeader(http.StatusBadRequest)
		return returnData, false
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return returnData, false
	}

	// Unmarshal
	var metaData map[string]interface{}
	err = json.Unmarshal(body, &metaData)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return returnData, false
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
	if !<-statusChan {
		return returnData, false
	}
	returnData.Title = <-titleChan
	returnData.Units = <-unitsChan
	log.Printf("Got: %d : %s", response.StatusCode, returnData.Title)

	return returnData, true
}

func censusData(censusSeries string, writer http.ResponseWriter, request *http.Request) (graphData, bool) {
	log.Printf("Graph request for Population of: %s", censusSeries)

	var returnData graphData
	returnData.Data = make([]seriesData, 0, 100)

	// Get the data
	response, err := otelhttp.Get(request.Context(), CensusApi+"?get=POP,NAME,GENC&YR=&key="+CensusApiKey)
	if err != nil || response.StatusCode > http.StatusOK {
		writer.WriteHeader(http.StatusBadRequest)
		return returnData, false
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return returnData, false
	}

	// Unmarshal
	var populationArray [][]string
	err = json.Unmarshal(body, &populationArray)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return returnData, false
	}
	for index, populationEntry := range populationArray {
		// Ignore title index
		if index == 0 {
			continue
		}
		// Name encoding does not match
		if populationEntry[2] != censusSeries {
			continue
		}
		returnData.Title = populationEntry[1]
		returnData.Data = append(returnData.Data, seriesData{
			Date:  populationEntry[3],
			Value: populationEntry[0],
		})
	}

	log.Printf("Got: %d : %s", response.StatusCode, returnData.Title)

	return returnData, true
}

func GraphApiHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		writer.Header().Set("Content-Type", "application/json")

		fredSeries := chi.URLParam(request, "fredSeries")
		censusSeries := chi.URLParam(request, "censusSeries")

		var returnData graphData
		var status bool
		if fredSeries != "" {
			returnData, status = fredData(fredSeries, writer, request)
		} else if censusSeries != "" {
			returnData, status = censusData(censusSeries, writer, request)
		} else {
			writer.WriteHeader(http.StatusNoContent)
			return
		}

		if !status {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		encodedReturnData, _ := json.Marshal(returnData)
		writer.Write(encodedReturnData)
	}
}

func ApiNotFoundHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
}
