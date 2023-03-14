package internal

import (
	"fmt"
	"github.com/gocql/gocql"
	"log"
	"net/http"
	"time"
)

var CassandraServer string

func requestInfo(request *http.Request, tableName string) {
	var info RequestInfo

	info.UserAgent = request.Header.Get("User-Agent")
	remoteAddress := request.RemoteAddr
	info.RemoteAddress = request.Header.Get("X-Forwarded-For")

	// Once the information is extracted from the request, the remainder of processing can be
	// done concurrently
	go func() {
		//info.RemoteAddress = strings.Trim(strings.Split(remoteAddress, ":")[0], "[]")
		//log.Printf("Remote: %s -> %s", remoteAddress, info.RemoteAddress)

		// TODO: IP (v6?) address not found, log and skip
		if info.RemoteAddress == "" {
			log.Printf("WARNING: IP addr not found: %s; XFF: %s", remoteAddress, info.RemoteAddress)
			return
		}

		// Cassandra session
		cluster := gocql.NewCluster(CassandraServer)
		cluster.Keyspace = "lobo_codes"
		cluster.Consistency = gocql.Quorum
		session, err := cluster.CreateSession()
		if err != nil {
			log.Printf("WARNING: failed to create session with cassandra database: %s; %s", CassandraServer, err.Error())
			return
		}
		defer session.Close()

		// Find a pre-existing IP Address & UserAgent entry
		queryString := fmt.Sprintf("SELECT id,created_at,updated_at,deleted_at,remote_address,user_agent,count,country_short,country_long,region,city,latitude,longitude,zipcode,timezone,elevation "+
			"FROM %s WHERE remote_address = ? AND user_agent = ? LIMIT 1 ALLOW FILTERING ", tableName)
		err = session.Query(queryString, info.RemoteAddress, info.UserAgent).
			Scan(&info.ID, &info.CreatedAt, &info.UpdatedAt, &info.UpdatedAt, &info.RemoteAddress, &info.UserAgent,
				&info.Count, &info.CountryShort, &info.CountryLong, &info.Region, &info.City, &info.Latitude,
				&info.Longitude, &info.Zipcode, &info.Timezone, &info.Elevation)

		if err != nil && err != gocql.ErrNotFound {
			log.Printf("WARNING: Error during lookup of IP: %s; %s", info.RemoteAddress, err.Error())
			return
		}

		// Update the visitor
		info.Count += 1
		info.UpdatedAt = time.Now()

		// New visitor - get geo info
		if err == gocql.ErrNotFound {
			geoInfo(&info)
			info.CreatedAt = info.UpdatedAt
		}

		insertString := fmt.Sprintf("INSERT INTO %s "+
			"(id,created_at,updated_at,deleted_at,remote_address,user_agent,count,country_short,country_long,region,city,latitude,longitude,zipcode,timezone,elevation) "+
			"VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)", tableName)
		err = session.Query(insertString, info.ID, info.CreatedAt, info.UpdatedAt, time.Time{}, info.RemoteAddress, info.UserAgent,
			info.Count, info.CountryShort, info.CountryLong, info.Region, info.City, info.Latitude, info.Longitude, info.Zipcode,
			info.Timezone, info.Elevation).Exec()
		if err != nil {
			log.Printf("WARNING: Error inserting %+v; %s", info, err.Error())
		}
		log.Printf("INFO: Inserted %v\n", info)
	}()
}
