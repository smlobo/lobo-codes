package main

import (
	"flag"
	"fmt"
	"github.com/gocql/gocql"
	"strings"
	"time"

	//"github.com/rqlite/gorqlite"
	"k8s-http-server/internal"
	"log"
)

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

func main() {
	// Input arguments
	cassandraTblPtr := flag.String("cassandraTbl", "", "cassandra table to read from")
	flag.Parse()

	// config - for cassandra/rqlite server
	internal.SetupConfig()

	// Init rqlite
	internal.InitRqlite()

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

	// Create rqlite table (Table created by hand)
	//query = "create table " + *cassandraTblPtr + "(id INTEGER NOT NULL PRIMARY KEY, created_at datetime, " +
	//	"updated_at datetime, remote_address TEXT, user_agent TEXT, count INTEGER, country_short TEXT, " +
	//	"country_long TEXT, region TEXT, city TEXT, latitude REAL, longitude REAL, zipcode TEXT, timezone TEXT, " +
	//	"elevation REAL)"

	// Insert into the previously created table
	//statements := make([]string, len(requestInfos))
	pattern := "INSERT INTO " + *cassandraTblPtr + " (created_at,updated_at,remote_address,user_agent,count," +
		"country_short,country_long,region,city,latitude,longitude,zipcode,timezone,elevation) VALUES ('%s', '%s', " +
		"'%s', '%s', %d, '%s', '%s', '%s', '%s', %f, %f, '%s', '%s', %f)"
	for i, requestInfo := range requestInfos {
		// Weird Mozilliqa"<?=print... User-Agent
		if strings.HasPrefix(requestInfo.UserAgent, "Mozilliqa") {
			log.Printf("Skipping: [%d] %s", i, requestInfo)
			continue
		}
		//statements[i] = fmt.Sprintf(pattern, requestInfo.CreatedAt, requestInfo.UpdatedAt, requestInfo.RemoteAddress,
		//	requestInfo.UserAgent, requestInfo.Count, requestInfo.CountryShort, requestInfo.CountryLong,
		//	requestInfo.City, requestInfo.Region, requestInfo.Latitude, requestInfo.Longitude, requestInfo.Zipcode,
		//	requestInfo.Timezone, requestInfo.Elevation)

		insertQuery := fmt.Sprintf(pattern, requestInfo.CreatedAt.Format(time.RFC3339Nano),
			requestInfo.UpdatedAt.Format(time.RFC3339Nano), requestInfo.RemoteAddress,
			strings.Trim(requestInfo.UserAgent, "\""), requestInfo.Count, requestInfo.CountryShort,
			requestInfo.CountryLong, requestInfo.Region, requestInfo.City, requestInfo.Latitude, requestInfo.Longitude,
			requestInfo.Zipcode, requestInfo.Timezone, requestInfo.Elevation)

		err := internal.RqliteExecute(insertQuery)
		if err != nil {
			log.Fatalf("Error inserting [%d] %s; %s", i, insertQuery, err)
		} else {
			log.Printf("[%d] success: %s", i, insertQuery)
		}
	}
}
