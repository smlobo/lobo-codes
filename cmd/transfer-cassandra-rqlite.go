package main

import (
	"flag"
	"fmt"
	"github.com/gocql/gocql"
	"gorm.io/gorm"
	"k8s-http-server/internal"
	"log"
)

type RequestInfo3 struct {
	gorm.Model
	RemoteAddress string
	UserAgent     string
	Count         int64
	internal.GeoLocation
}

func (ri RequestInfo3) String() string {
	return fmt.Sprintf("<%s / %s> %s {%d} %s [%s]", ri.CreatedAt.UTC(), ri.UpdatedAt, ri.RemoteAddress,
		ri.Count, ri.UserAgent, ri.GeoLocation)
}

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
	cassandraTblPtr := flag.String("cassandraTbl", "", "cassandra table to write to")
	flag.Parse()

	// config - for cassandra server
	internal.SetupConfig()

	// Init cassandra db
	cassandraSession := cassandraInit()
	defer cassandraSession.Close()

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
	lines := 0
	//for iter.Scan(&id, &created_at, &updated_at, &remote_address, &user_agent, &count, &country_short, &country_long,
	//	&region, &city, &latitude, &longitude, &zipcode, &timezone, &elevation) {
	for iter.Scan(&id, &ri.CreatedAt, &ri.UpdatedAt, &ri.RemoteAddress, &ri.UserAgent, &ri.Count, &ri.CountryShort,
		&ri.CountryLong, &ri.Region, &ri.City, &ri.Latitude, &ri.Longitude, &ri.Zipcode, &ri.Timezone, &ri.Elevation) {
		// Process each row
		//fmt.Printf("[%d], %s, %s, %s, %s, %s, %d, %.3f\n", lines, id, created_at, remote_address, user_agent,
		//	city, count, elevation)
		requestInfos = append(requestInfos, ri)
		lines++
	}

	if err := iter.Close(); err != nil {
		log.Fatal(err)
	}
	log.Printf("Read %d lines", lines)

	// Read all records
	for i, requestInfo := range requestInfos {
		fmt.Printf("[%d] %v\n", i, requestInfo)
	}
}
