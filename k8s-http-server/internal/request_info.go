package internal

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gocql/gocql"
	"github.com/google/uuid"
)

var CassandraCluster *gocql.ClusterConfig

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
		session, err := CassandraCluster.CreateSession()
		if err != nil {
			log.Printf("WARNING: failed to create session with cassandra database: %s; %s",
				CassandraCluster.Hosts[0], err.Error())
			return
		}
		defer session.Close()

		// Find a pre-existing IP Address & UserAgent entry
		var foundId string
		queryString := fmt.Sprintf("SELECT id,count,country_short FROM %s WHERE remote_address=? AND "+
			"user_agent=?", tableName)
		err = session.Query(queryString, info.RemoteAddress, info.UserAgent).Scan(&foundId, &info.Count, &info.CountryShort)

		if err != nil && err != gocql.ErrNotFound {
			log.Printf("WARNING: Error during lookup of IP: %s; %s", info.RemoteAddress, err.Error())
			return
		}

		// Update the visitor
		info.Count += 1
		info.UpdatedAt = time.Now()

		if err == gocql.ErrNotFound {
			// New visitor - get geo info
			geoInfo(&info)
			info.CreatedAt = info.UpdatedAt

			if foundId != uuid.Nil.String() {
				log.Printf("WARNING: Request info not found, but UUID: %s for IP: %s; %s", foundId,
					info.RemoteAddress, err.Error())
				return
			}

			info.Id, err = uuid.NewRandom()
			if err != nil {
				log.Printf("WARNING: Error during UUID generation for IP: %s; %s", info.RemoteAddress, err.Error())
				return
			}

			insertString := fmt.Sprintf("INSERT INTO %s "+
				"(id,created_at,updated_at,remote_address,user_agent,count,country_short,country_long,region,city,latitude,longitude,zipcode,timezone,elevation) "+
				"VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)", tableName)
			err = session.Query(insertString, info.Id.String(), info.CreatedAt, info.UpdatedAt, info.RemoteAddress, info.UserAgent,
				info.Count, info.CountryShort, info.CountryLong, info.Region, info.City, info.Latitude, info.Longitude, info.Zipcode,
				info.Timezone, info.Elevation).Exec()
			if err != nil {
				log.Printf("WARNING: Error inserting %s; %s", info, err.Error())
				return
			}

			log.Printf("INFO: Inserted in %s : %v\n", tableName, info)
		} else {
			// Existing visitor
			updateString := fmt.Sprintf("UPDATE %s SET count=?,updated_at=? where remote_address=? AND user_agent=?",
				tableName)

			err = session.Query(updateString, info.Count, info.UpdatedAt, info.RemoteAddress, info.UserAgent).Exec()

			if err != nil {
				log.Printf("WARNING: Error updating %s; %s", info, err.Error())
				return
			}

			log.Printf("INFO: Updated in %s : %v\n", tableName, info)
		}
	}()
}
