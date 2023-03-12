package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/gocql/gocql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"k8s-http-server/internal"
)

func main() {
	// Input arguments
	requestsDbPtr := flag.String("sqliteDb", "", "sqlite db to read requests")
	cassandraDbPtr := flag.String("cassandraDb", "", "cassandra db to write to")
	cassandraTblPtr := flag.String("cassandraTbl", "", "cassandra table to write to")
	flag.Parse()

	db, err := gorm.Open(sqlite.Open(*requestsDbPtr), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to open sqlite database: %s", *requestsDbPtr)
	} else {
		log.Printf("Reading from sqlite database: %s\n", *requestsDbPtr)
	}

	cluster := gocql.NewCluster(*cassandraDbPtr)
	cluster.Keyspace = "lobo_codes"
	cluster.Consistency = gocql.Quorum
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("failed to create session with cassandra database: %s", *cassandraDbPtr)
	}
	defer session.Close()

	// Read all records
	var requestInfos []internal.RequestInfo
	result := db.Find(&requestInfos)
	fmt.Printf("Found: %d request\n", result.RowsAffected)
	//for i, product := range requestInfos {
	//	//fmt.Printf("[%d] %v\n", i, product)
	//}

	// Write to cassandra
	insertString := fmt.Sprintf("INSERT INTO %s "+
		"(id,created_at,updated_at,deleted_at,remote_address,user_agent,count,country_short,country_long,region,city,latitude,longitude,zipcode,timezone,elevation) "+
		"VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)", *cassandraTblPtr)
	zeroTime := time.Time{}
	for i, rI := range requestInfos {
		err := session.Query(insertString, rI.ID, rI.CreatedAt, rI.UpdatedAt, zeroTime, rI.RemoteAddress, rI.UserAgent,
			rI.Count, rI.CountryShort, rI.CountryLong, rI.Region, rI.City, rI.Latitude, rI.Longitude, rI.Zipcode,
			rI.Timezone, rI.Elevation).Exec()
		if err != nil {
			log.Fatalf("inserting [%d] %+v : %s", i, rI, err.Error())
		}
		log.Printf("[%d] Inserted %v\n", i, rI)
	}
}
