package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
)

const RqlitePort = "4001"

var rqliteURL string

const RqliteDefaultVersion = "vX.Y.Z"

var RqliteVersion string

func InitRqlite() {

	RqliteVersion = RqliteDefaultVersion

	// Closure
	getRqliteVersion := func() bool {
		// Set the URL
		rqliteURL = "http://" + Config["RQLITE_SERVER"] + ":" + RqlitePort
		log.Printf("Using Rqlite URL: %s", rqliteURL)

		// Get the version
		resp, err := http.Get(rqliteURL + "/status")
		if err != nil {
			log.Printf("Error getting rqlite version: %s", err)
			return false
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Printf("Unexpected status code from rqlite GET to /status endpoint: %s", resp.Status)
			return false
		}

		// Parse json response
		var statusJson map[string]interface{}
		respBytes, _ := io.ReadAll(resp.Body)
		err = json.Unmarshal(respBytes, &statusJson)
		if err != nil {
			log.Printf("Error parsing rqlite GET /status: %s", err)
			return false
		}
		buildJson, ok := statusJson["build"].(map[string]interface{})
		if !ok {
			log.Printf("Rqlite GET /status has no build field that type asserts into map[string]interface{}")
		}
		RqliteVersion, ok = buildJson["version"].(string)
		if !ok {
			log.Printf("Rqlite GET /status build has no version field that type asserts into string")
			return false
		}

		return true
	}

	// Try 3 times
	// Update: The unmarshal problem was the entire structure was not map[string]map[string]interface{}
	// map["build"] fits this, but not all the others.
	// But, using a closure to retry is a good idea, so leaving this with a single try.
	for i := 0; i < 1; i++ {
		if getRqliteVersion() {
			log.Printf("  [%d] %s", i, RqliteVersion)
			break
		}
	}
	log.Printf("Got Rqlite version: %s", RqliteVersion)
}

func rqliteLogRequest(info *RequestInfo, tableName string, request *http.Request) {
	//info.RemoteAddress = strings.Trim(strings.Split(remoteAddress, ":")[0], "[]")
	//log.Printf("Remote: %s -> %s", remoteAddress, info.RemoteAddress)

	// TODO: IP (v6?) address not found, log and skip
	if info.RemoteAddress == "" {
		log.Printf("WARNING: IP addr not found: %s; XFF: %s", info.OrigRemoteAddress, info.RemoteAddress)
		return
	}

	// Find a pre-existing IP Address & UserAgent entry
	queryString := fmt.Sprintf("SELECT id FROM %s WHERE remote_address=? AND user_agent=?", rqliteTable(tableName))
	rows, err := RqliteQuery(queryString, info.RemoteAddress, info.UserAgent)
	if err != nil {
		log.Printf("WARNING: Error during lookup of IP: %s, user agent: %s; %s", info.RemoteAddress,
			info.UserAgent, err.Error())
		return
	}

	// Should be only 1 existing
	if len(rows) > 1 {
		log.Printf("Multiple %s / %s found: %d", info.RemoteAddress, info.UserAgent, len(rows))
		return
	}

	info.UpdatedAt = time.Now()

	if len(rows) == 0 {
		// New visitor - get geo info
		geoInfo(info)
		info.CreatedAt = info.UpdatedAt
		info.Count = 1

		insertQuery := fmt.Sprintf("INSERT INTO %s (created_at,updated_at,remote_address,user_agent,count,country_short,"+
			"country_long,region,city,latitude,longitude,zipcode,timezone,elevation) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)", rqliteTable(tableName))
		err = RqliteExecute(insertQuery, info.CreatedAt.Format(time.RFC3339Nano),
			info.UpdatedAt.Format(time.RFC3339Nano), info.RemoteAddress, info.UserAgent,
			info.Count, info.CountryShort, info.CountryLong, info.Region, info.City, info.Latitude, info.Longitude,
			info.Zipcode, info.Timezone, info.Elevation)
		if err != nil {
			log.Printf("WARN: error inserting new entry into rqlite %s: %s; %s", tableName, info, err.Error())
		}
		log.Printf("INFO: Inserted in %s : %s", tableName, info)
	} else {
		// Existing visitor
		rowMap, ok := rows[0].(map[string]interface{})
		if !ok {
			log.Printf("unexpected rqlite row type: %T", rows[0])
			return
		}
		idValue, ok := rowMap["id"].(float64)
		if !ok {
			log.Printf("unexpected rqlite id: %v", rowMap["id"])
			return
		}
		id := int(idValue)

		updateQuery := fmt.Sprintf("UPDATE %s SET count = count + 1, updated_at = ? WHERE id = ?", rqliteTable(tableName))
		err = RqliteExecute(updateQuery, info.UpdatedAt.Format(time.RFC3339Nano), id)
		if err != nil {
			log.Printf("WARN: error updating entry in rqlite %s: %s, [id: %d]; %s", tableName, info,
				id, err.Error())
			return
		}
		log.Printf("INFO: Updated in %s: %s [id: %d]", tableName, info, id)
	}
}

func rqliteGetCountriesCities(tableName string, request *http.Request) (countryCount map[string]int,
	cityCount map[string]City) {

	// Getting info from DB
	_, span := otel.Tracer("k8s-http-server").Start(request.Context(), "db-query")

	// Read country name & count
	// Also, the city & region to count
	queryString := fmt.Sprintf("SELECT country_short, city, region FROM %s", rqliteTable(tableName))
	rows, err := RqliteQuery(queryString)
	if err != nil {
		log.Printf("WARNING: Error during country/city/region lookup for %s; %s", tableName, err.Error())
		return
	}

	// Processing the data
	span.End()
	_, span = otel.Tracer("k8s-http-server").Start(request.Context(), "db-process")
	defer span.End()

	countryCount = make(map[string]int)
	cityCount = make(map[string]City)

	// Iterate over result rows
	for _, row := range rows {
		rowMap, _ := row.(map[string]interface{})
		country := rowMap["country_short"].(string)
		city := rowMap["city"].(string)
		region := rowMap["region"].(string)

		if count, ok := countryCount[country]; !ok {
			countryCount[country] = 1
		} else {
			countryCount[country] = count + 1
		}

		if count, ok := cityCount[city]; !ok {
			cityCount[city] = City{
				City:         city,
				Region:       region,
				CountryShort: country,
				Count:        1,
			}
		} else {
			count.Count += 1
			cityCount[city] = count
		}
	}

	return
}

func rqliteTable(tableName string) string {
	return `"` + strings.ReplaceAll(tableName, `"`, `""`) + `"`
}

func rqliteRequestBody(queryString string, args ...interface{}) ([]byte, error) {
	if len(args) == 0 {
		return json.Marshal([]string{queryString})
	}
	statement := append([]interface{}{queryString}, args...)
	return json.Marshal([]interface{}{statement})
}

func RqliteQuery(queryString string, args ...interface{}) ([]interface{}, error) {
	url := rqliteURL + "/db/query?associative"
	body, err := rqliteRequestBody(queryString, args...)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("not OK status code: %d while querying rqlite", resp.StatusCode)
	}

	var result struct {
		Results []struct {
			Rows  []interface{} `json:"rows"`
			Error string        `json:"error"`
		} `json:"results"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	if len(result.Results) == 0 {
		return nil, fmt.Errorf("rqlite query returned no result")
	}
	if result.Results[0].Error != "" {
		return nil, fmt.Errorf("rqlite query: %s", result.Results[0].Error)
	}
	return result.Results[0].Rows, nil
}

func RqliteExecute(queryString string, args ...interface{}) error {
	url := rqliteURL + "/db/execute"
	body, err := rqliteRequestBody(queryString, args...)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("not OK status code: %d while querying rqlite", resp.StatusCode)
	}

	var result struct {
		Results []struct {
			Error string `json:"error"`
		} `json:"results"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}
	if len(result.Results) == 0 {
		return fmt.Errorf("rqlite execute returned no result")
	}
	if result.Results[0].Error != "" {
		return fmt.Errorf("rqlite execute: %s", result.Results[0].Error)
	}
	return nil
}
