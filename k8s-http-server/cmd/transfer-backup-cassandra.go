package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/google/uuid"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gocql/gocql"
	"gorm.io/gorm"

	"k8s-http-server/internal"
)

type RequestInfo2 struct {
	gorm.Model
	RemoteAddress string
	UserAgent     string
	Count         int64
	internal.GeoLocation
}

func (ri RequestInfo2) String() string {
	return fmt.Sprintf("<%s / %s> %s {%d} %s [%s]", ri.CreatedAt.UTC(), ri.UpdatedAt, ri.RemoteAddress,
		ri.Count, ri.UserAgent, ri.GeoLocation)
}

func main() {
	// Input arguments
	backupPtr := flag.String("backup", "", "backup to read from")
	cassandraDbPtr := flag.String("cassandraDb", "", "cassandra db to write to")
	cassandraTblPtr := flag.String("cassandraTbl", "", "cassandra table to write to")
	flag.Parse()

	cluster := gocql.NewCluster(*cassandraDbPtr)
	cluster.Keyspace = "lobo_codes"
	cluster.Consistency = gocql.Quorum
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("failed to create session with cassandra database: %s", *cassandraDbPtr)
	}
	defer session.Close()

	// Read the backup file
	log.Printf("Reading file: %s", *backupPtr)
	readFile, err := os.Open(*backupPtr)
	if err != nil {
		log.Fatalf("failed to read file: %s [%s]", *backupPtr, err)
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	count := 0
	requestInfos := make([]RequestInfo2, 0)
	for fileScanner.Scan() {
		//fmt.Println(fileScanner.Text())
		line := fileScanner.Text()
		words := strings.Split(line, "|")
		if len(words) != 15 {
			continue
		}
		if strings.TrimSpace(words[0]) == "remote_address" {
			continue
		}
		ri := RequestInfo2{}
		ri.RemoteAddress = strings.TrimSpace(words[0])
		ri.UserAgent = strings.TrimSpace(words[1])
		ri.City = strings.TrimSpace(words[2])
		if intVar, err := strconv.ParseInt(strings.TrimSpace(words[3]), 0, 64); err != nil {
			log.Fatalf("bad count %s for : %s", strings.TrimSpace(words[3]), line)
		} else {
			ri.Count = intVar
		}
		ri.CountryLong = strings.TrimSpace(words[4])
		ri.CountryShort = strings.TrimSpace(words[5])
		if theTime, err := time.Parse("2006-01-02 15:04:05.000000-0700", strings.TrimSpace(words[6])); err != nil {
			log.Fatalf("bad createdAt %s for : %s", strings.TrimSpace(words[6]), line)
		} else {
			ri.CreatedAt = theTime
		}
		if floatVar, err := strconv.ParseFloat(strings.TrimSpace(words[7]), 32); err != nil {
			log.Fatalf("bad elevation %s for : %s", strings.TrimSpace(words[3]), line)
		} else {
			ri.Elevation = float32(floatVar)
		}
		//ri.ID
		if floatVar, err := strconv.ParseFloat(strings.TrimSpace(words[9]), 32); err != nil {
			log.Fatalf("bad latitude %s for : %s", strings.TrimSpace(words[3]), line)
		} else {
			ri.Latitude = float32(floatVar)
		}
		if floatVar, err := strconv.ParseFloat(strings.TrimSpace(words[10]), 32); err != nil {
			log.Fatalf("bad longitude %s for : %s", strings.TrimSpace(words[3]), line)
		} else {
			ri.Longitude = float32(floatVar)
		}
		ri.Region = strings.TrimSpace(words[11])
		ri.Timezone = strings.TrimSpace(words[12])
		if theTime, err := time.Parse("2006-01-02 15:04:05.000000-0700", strings.TrimSpace(words[13])); err != nil {
			log.Fatalf("bad createdAt %s for : %s", strings.TrimSpace(words[6]), line)
		} else {
			ri.UpdatedAt = theTime
		}
		ri.Zipcode = strings.TrimSpace(words[14])

		requestInfos = append(requestInfos, ri)
		count++
	}
	_ = readFile.Close()
	log.Printf("Read %d lines", count)

	// Read all records
	for i, product := range requestInfos {
		fmt.Printf("[%d] %v\n", i, product)
	}

	// Write to cassandra
	insertString := fmt.Sprintf("INSERT INTO %s "+
		"(id,created_at,updated_at,remote_address,user_agent,count,country_short,country_long,region,city,latitude,longitude,zipcode,timezone,elevation) "+
		"VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)", *cassandraTblPtr)
	for i, rI := range requestInfos {
		id, err := uuid.NewRandom()
		if err != nil {
			log.Fatalf("uuid generation error: %s", err.Error())
		}

		err = session.Query(insertString, id.String(), rI.CreatedAt, rI.UpdatedAt, rI.RemoteAddress, rI.UserAgent,
			rI.Count, rI.CountryShort, rI.CountryLong, rI.Region, rI.City, rI.Latitude, rI.Longitude, rI.Zipcode,
			rI.Timezone, rI.Elevation).Exec()
		if err != nil {
			log.Fatalf("inserting [%d] %+v : %s", i, rI, err.Error())
		}
		log.Printf("[%d] Inserted %v\n", i, rI)
	}
}
