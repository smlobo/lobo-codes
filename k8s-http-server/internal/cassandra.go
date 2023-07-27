package internal

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
)

var cassandraCluster *gocql.ClusterConfig
var cassandraSession *gocql.Session

func InitCassandra(cassandraServer string) *gocql.Session {
	var err error

	cassandraCluster = gocql.NewCluster(cassandraServer)
	cassandraCluster.Keyspace = "lobo_codes"
	cassandraCluster.Consistency = gocql.Quorum

	cassandraSession, err = cassandraCluster.CreateSession()
	if err != nil {
		log.Fatalf("failed to create session with cassandra database: %s; %s",
			cassandraCluster.Hosts[0], err.Error())
	}

	return cassandraSession
}

func cassandraLogRequest(info *RequestInfo, tableName string, request *http.Request) {
	//info.RemoteAddress = strings.Trim(strings.Split(remoteAddress, ":")[0], "[]")
	//log.Printf("Remote: %s -> %s", remoteAddress, info.RemoteAddress)

	// TODO: IP (v6?) address not found, log and skip
	if info.RemoteAddress == "" {
		log.Printf("WARNING: IP addr not found: %s; XFF: %s", info.OrigRemoteAddress, info.RemoteAddress)
		return
	}

	// inserting into the DB
	_, span := otel.Tracer("k8s-http-server").Start(request.Context(), "db-insert")
	defer span.End()

	// Find a pre-existing IP Address & UserAgent entry
	var foundId string
	queryString := fmt.Sprintf("SELECT id,count,country_short FROM %s WHERE remote_address=? AND "+
		"user_agent=?", tableName)
	err := cassandraSession.Query(queryString, info.RemoteAddress, info.UserAgent).Scan(&foundId, &info.Count,
		&info.CountryShort)

	if err != nil && err != gocql.ErrNotFound {
		log.Printf("WARNING: Error during lookup of IP: %s; %s", info.RemoteAddress, err.Error())
		return
	}

	// Update the visitor
	info.Count += 1
	info.UpdatedAt = time.Now()

	if err == gocql.ErrNotFound {
		// New visitor - get geo info
		geoInfo(info)
		info.CreatedAt = info.UpdatedAt

		info.Id, err = uuid.NewRandom()
		if err != nil {
			log.Printf("WARNING: Error during UUID generation for IP: %s; %s", info.RemoteAddress, err.Error())
			return
		}

		insertString := fmt.Sprintf("INSERT INTO %s "+
			"(id,created_at,updated_at,remote_address,user_agent,count,country_short,country_long,region,city,latitude,longitude,zipcode,timezone,elevation) "+
			"VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)", tableName)
		err = cassandraSession.Query(insertString, info.Id.String(), info.CreatedAt, info.UpdatedAt, info.RemoteAddress,
			info.UserAgent, info.Count, info.CountryShort, info.CountryLong, info.Region, info.City, info.Latitude,
			info.Longitude, info.Zipcode, info.Timezone, info.Elevation).Exec()
		if err != nil {
			log.Printf("WARNING: Error inserting %s; %s", info, err.Error())
			return
		}

		log.Printf("INFO: Inserted in %s : %v\n", tableName, info)
	} else {
		// Existing visitor
		updateString := fmt.Sprintf("UPDATE %s SET count=?,updated_at=? where remote_address=? AND user_agent=?",
			tableName)

		err = cassandraSession.Query(updateString, info.Count, info.UpdatedAt, info.RemoteAddress, info.UserAgent).Exec()

		if err != nil {
			log.Printf("WARNING: Error updating %s; %s", info, err.Error())
			return
		}

		log.Printf("INFO: Updated in %s : %v\n", tableName, info)
	}
}

func cassandraGetCountriesCities(tableName string, request *http.Request) (countryCount map[string]int,
	cityCount map[string]City) {

	// Getting info from DB
	_, span := otel.Tracer("k8s-http-server").Start(request.Context(), "db-query")

	// Read country name & count
	// Also, the city & region to count
	queryString := fmt.Sprintf("SELECT country_short, city, region FROM %s ALLOW FILTERING", tableName)
	scanner := cassandraSession.Query(queryString).Iter().Scanner()

	// Processing the data
	span.End()
	_, span = otel.Tracer("k8s-http-server").Start(request.Context(), "db-process")
	defer span.End()

	countryCount = make(map[string]int)
	cityCount = make(map[string]City)

	for scanner.Next() {
		var countryShort, city, region string
		err := scanner.Scan(&countryShort, &city, &region)
		if err != nil {
			continue
		}

		if count, ok := countryCount[countryShort]; !ok {
			countryCount[countryShort] = 1
		} else {
			countryCount[countryShort] = count + 1
		}

		if count, ok := cityCount[city]; !ok {
			cityCount[city] = City{
				City:         city,
				Region:       region,
				CountryShort: countryShort,
				Count:        1,
			}
		} else {
			count.Count += 1
			cityCount[city] = count
		}
	}
	return
}
