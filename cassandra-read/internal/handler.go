package internal

import (
	"fmt"
	"github.com/gocql/gocql"
	"log"
	"net/http"
	"runtime"
)

func ReturnHandler() http.HandlerFunc {
	var count int

	//cluster := gocql.NewCluster("10.152.183.196")
	cluster := gocql.NewCluster("cassandra-internal")
	cluster.Keyspace = "cycling"
	cluster.Consistency = gocql.Quorum

	return func(writer http.ResponseWriter, request *http.Request) {
		session, _ := cluster.CreateSession()
		defer session.Close()

		if request.Method != "GET" {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		log.Printf("[%d] Received request from: %s", count, request.UserAgent())
		count++
		writer.Header().Add("Server", runtime.Version())
		writer.WriteHeader(http.StatusOK)

		var id gocql.UUID
		var name string

		responseString := ""

		iter := session.Query(`SELECT id, lastname FROM cyclist_alt_stats`).Iter()
		for iter.Scan(&id, &name) {
			//fmt.Println("Cyclist stats:", id, name)
			responseString = fmt.Sprintf("%sCyclist Stats: %s, %s\n", responseString, id, name)
		}
		if err := iter.Close(); err != nil {
			log.Fatal(err)
		}

		_, _ = writer.Write([]byte(responseString))
	}
}
