package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
)

const RqlitePort = "4001"

var rqliteURL string

func InitRqlite() {
	// Set the URL
	rqliteURL = "http://" + Config["RQLITE_SERVER"] + ":" + RqlitePort
	log.Printf("Using Rqlite URL: %s", rqliteURL)
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
	queryString := fmt.Sprintf("SELECT id,count FROM %s WHERE remote_address=%s AND user_agent=%s", tableName,
		info.RemoteAddress, info.UserAgent)
	rows, err := RqliteQuery(queryString)
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

		pattern := "INSERT INTO " + tableName + " (created_at,updated_at,remote_address,user_agent,count," +
			"country_short,country_long,region,city,latitude,longitude,zipcode,timezone,elevation) VALUES ('%s', '%s', " +
			"'%s', '%s', %d, '%s', '%s', '%s', '%s', %f, %f, '%s', '%s', %f)"

		insertQuery := fmt.Sprintf(pattern, info.CreatedAt.Format(time.RFC3339Nano),
			info.UpdatedAt.Format(time.RFC3339Nano), info.RemoteAddress, strings.Trim(info.UserAgent, "\""),
			info.Count, info.CountryShort, info.CountryLong, info.Region, info.City, info.Latitude, info.Longitude,
			info.Zipcode, info.Timezone, info.Elevation)

		err = RqliteExecute(insertQuery)
		if err != nil {
			log.Printf("WARN: error inserting new entry into rqlite %s: %s; %s", tableName, info, err.Error())
		}
		log.Printf("INFO: Inserted in %s : %s", tableName, info)
	} else {
		// Existing visitor
		rowMap, ok := rows[0].(map[string]interface{})
		if !ok {
			log.Printf("result not string-int map: %T", rows[0])
		}
		newCount := int(rowMap["count"].(float64)) + 1
		id := int(rowMap["id"].(float64))

		insertUpdate := fmt.Sprintf("UPDATE %s SET count = %d, updated_at = \"%s\" WHERE id = %d", tableName,
			newCount, info.UpdatedAt.Format(time.RFC3339Nano), id)

		err = RqliteExecute(insertUpdate)
		if err != nil {
			log.Printf("WARN: error updating entry into rqlite %s: %s, [count: %d, id: %d]; %s", tableName, info,
				newCount, id, err.Error())
		}
		log.Printf("INFO: Updated in %s: %s [count: %d, id: %d]", tableName, info, newCount, id)
	}
}

func rqliteGetCountriesCities(tableName string, request *http.Request) (countryCount map[string]int,
	cityCount map[string]City) {

	// Getting info from DB
	_, span := otel.Tracer("k8s-http-server").Start(request.Context(), "db-query")

	// Read country name & count
	// Also, the city & region to count
	queryString := fmt.Sprintf("SELECT country_short, city, region FROM %s", tableName)
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

func RqliteQuery(queryString string) ([]interface{}, error) {
	url := rqliteURL + "/db/query?associative"
	body := []byte(fmt.Sprintf("[\"%s\"]", queryString))

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("not OK status code: %d while querying rqlite", resp.StatusCode)
	}

	// Unmarshal response into Json
	var resultsJson map[string][]map[string]interface{}
	bytes, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(bytes, &resultsJson)
	if err != nil {
		return nil, err
	}

	rowsArray, ok := resultsJson["results"][0]["rows"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("rows array not found")
	}

	return rowsArray, nil
}

func RqliteExecute(queryString string) error {
	url := rqliteURL + "/db/execute"
	body := []byte(fmt.Sprintf("[\"%s\"]", queryString))

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("not OK status code: %d while querying rqlite", resp.StatusCode)
	}

	// Unmarshal response into Json
	return nil
}
