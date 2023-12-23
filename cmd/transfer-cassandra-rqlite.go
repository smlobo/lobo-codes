package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/gocql/gocql"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	//"github.com/rqlite/gorqlite"
	"k8s-http-server/internal"
	"log"
)

const rqliteURL = "http://10.1.1.44:4001"

func cassandraInit() *gocql.Session {
	var err error

	cassandraServer := internal.Config["CASSANDRA_SERVER"]
	log.Printf("Connecting to Cassandra server: %s", cassandraServer)

	cassandraCluster := gocql.NewCluster(cassandraServer)
	cassandraCluster.Keyspace = "lobo_codes"
	cassandraCluster.Consistency = gocql.One

	cassandraSession, err := cassandraCluster.CreateSession()
	if err != nil {
		log.Fatalf("failed to create session with cassandra database: %s; %s",
			cassandraCluster.Hosts[0], err.Error())
	}

	log.Printf("Connected to Cassandra server: %s", cassandraCluster.Hosts)

	return cassandraSession
}

//func rqliteInit() *gorqlite.Connection {
//	rqliteServer := internal.Config["RQLITE_SERVER"]
//	log.Printf("Connecting to `rqlite` server: %s", rqliteServer)
//
//	conn, err := gorqlite.Open("http://" + rqliteServer + ":4001/")
//	if err != nil {
//		log.Fatalf("failed to create connection with rqlite db: %s; %s", rqliteServer, err.Error())
//	}
//
//	log.Printf("Connected to rqlite server with ID: %s", conn.ID)
//leader, err := conn.Leader()
//if err != nil {
//	log.Fatalf("could not get leader for rqlite: %s; %s", rqliteServer, err)
//}
//peers, err := conn.Peers()
//if err != nil {
//	log.Fatalf("could not get peers for rqlite: %s; %s", rqliteServer, err)
//}
//log.Printf("rqlite Leader: %s", leader)
//for i, peer := range peers {
//	log.Printf("rqlite Peer: [%d] %s", i, peer)
//}
//
//	return conn
//}

func main() {
	// Input arguments
	cassandraTblPtr := flag.String("cassandraTbl", "", "cassandra table to read from")
	flag.Parse()

	// config - for cassandra server
	internal.SetupConfig()

	// Init cassandra db
	cassandraSession := cassandraInit()
	//defer cassandraSession.Close()

	// Read the cassandra table
	log.Printf("Reading cassandra table: %s", *cassandraTblPtr)
	requestInfos := make([]internal.RequestInfo, 0, 500)

	// Query to read data from a table
	query := "SELECT id,created_at,updated_at,remote_address,user_agent,count,country_short,country_long,region,city," +
		"latitude,longitude,zipcode,timezone,elevation FROM " + *cassandraTblPtr
	iter := cassandraSession.Query(query).Iter()

	// Iterate over the results
	var id gocql.UUID
	//var created_at, updated_at time.Time
	//var remote_address, user_agent, country_short, country_long, region, city, zipcode, timezone string
	//var latitude, longitude, elevation float32
	//var count int64
	var ri internal.RequestInfo
	//for iter.Scan(&id, &created_at, &updated_at, &remote_address, &user_agent, &count, &country_short, &country_long,
	//	&region, &city, &latitude, &longitude, &zipcode, &timezone, &elevation) {
	for iter.Scan(&id, &ri.CreatedAt, &ri.UpdatedAt, &ri.RemoteAddress, &ri.UserAgent, &ri.Count, &ri.CountryShort,
		&ri.CountryLong, &ri.Region, &ri.City, &ri.Latitude, &ri.Longitude, &ri.Zipcode, &ri.Timezone, &ri.Elevation) {
		// Process each row
		//fmt.Printf("[%d], %s, %s, %s, %s, %s, %d, %.3f\n", lines, id, created_at, remote_address, user_agent,
		//	city, count, elevation)
		requestInfos = append(requestInfos, ri)
	}

	if err := iter.Close(); err != nil {
		log.Fatal(err)
	}
	cassandraSession.Close()

	log.Printf("Read %d lines", len(requestInfos))

	// Print all records
	//for i, requestInfo := range requestInfos {
	//	fmt.Printf("[%d] %v\n", i, requestInfo)
	//}

	// Create rqlite connection
	//rqliteConn := rqliteInit()
	//defer rqliteConn.Close()

	// Create rqlite table (Table created by hand)
	//query = "create table " + *cassandraTblPtr + "(id INTEGER NOT NULL PRIMARY KEY, created_at datetime, " +
	//	"updated_at datetime, remote_address TEXT, user_agent TEXT, count INTEGER, country_short TEXT, " +
	//	"country_long TEXT, region TEXT, city TEXT, latitude REAL, longitude REAL, zipcode TEXT, timezone TEXT, " +
	//	"elevation REAL)"

	// Create rqlite session
	//sessionID, err := createSession()
	//if err != nil {
	//	fmt.Println("Error creating session:", err)
	//	return
	//}

	// Insert into the previously created table
	//statements := make([]string, len(requestInfos))
	pattern := "INSERT INTO " + *cassandraTblPtr + " (created_at,updated_at,remote_address,user_agent,count," +
		"country_short,country_long,region,city,latitude,longitude,zipcode,timezone,elevation) VALUES ('%s', '%s', " +
		"'%s', '%s', %d, '%s', '%s', '%s', '%s', %f, %f, '%s', '%s', %f)"
	for i, requestInfo := range requestInfos {
		//statements[i] = fmt.Sprintf(pattern, requestInfo.CreatedAt, requestInfo.UpdatedAt, requestInfo.RemoteAddress,
		//	requestInfo.UserAgent, requestInfo.Count, requestInfo.CountryShort, requestInfo.CountryLong,
		//	requestInfo.City, requestInfo.Region, requestInfo.Latitude, requestInfo.Longitude, requestInfo.Zipcode,
		//	requestInfo.Timezone, requestInfo.Elevation)

		insertQuery := fmt.Sprintf(pattern, requestInfo.CreatedAt.Format(time.RFC3339Nano),
			requestInfo.UpdatedAt.Format(time.RFC3339Nano), requestInfo.RemoteAddress,
			strings.Trim(requestInfo.UserAgent, "\""), requestInfo.Count, requestInfo.CountryShort, requestInfo.CountryLong,
			requestInfo.City, requestInfo.Region, requestInfo.Latitude, requestInfo.Longitude, requestInfo.Zipcode,
			requestInfo.Timezone, requestInfo.Elevation)

		//err = executeQuery(sessionID, insertQuery)
		//if err != nil {
		//	fmt.Println("Error executing query [%d]:", i, err)
		//	return
		//}
		err := executeInsert(insertQuery)
		if err != nil {
			log.Fatalf("Error inserting [%d]; %s", i, err)
		} else {
			log.Printf("[%d] success: %s", i, insertQuery)
		}
	}
	//results, err := rqliteConn.Write(statements)
	//if err != nil {
	//	log.Fatalf("Error inserting into rqlite table: %s; %s", *cassandraTblPtr, err.Error())
	//}
	//for i, result := range results {
	//	log.Printf("[%d], %d rows affected", i, result.RowsAffected)
	//	if result.Err != nil {
	//		log.Fatalf("Error inserting into [%d] %s; %s", i, *cassandraTblPtr, err.Error())
	//	}
	//}
}

func createSession() (string, error) {
	url := rqliteURL + "/db/execute?pretty"
	body := []byte(`[{"stmt": "BEGIN;"}]`)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to create session, status code: %d", resp.StatusCode)
	}

	// Parse the response to get the session ID
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Extract the session ID from the response
	sessionID := string(respBody)
	return sessionID, nil
}

func executeQuery(sessionID, query string) error {
	url := rqliteURL + "/db/execute?pretty"
	body := []byte(fmt.Sprintf(`[{"stmt": " %s ", "transaction": "%s"}]`, query, sessionID))

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to execute query, status code: %d", resp.StatusCode)
	}

	return nil
}

func executeInsert(query string) error {
	url := rqliteURL + "/db/execute?pretty&timings"
	body := []byte(fmt.Sprintf("[\"%s\"]", query))

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to execute query %s, status code: %d", query, resp.StatusCode)
	}

	return nil
}
