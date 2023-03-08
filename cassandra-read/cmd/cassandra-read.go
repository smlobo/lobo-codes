package main

import (
	"fmt"
	"log"

	"github.com/gocql/gocql"
)

func main() {
	cluster := gocql.NewCluster("10.1.1.42")
	cluster.Keyspace = "cycling"
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()
	defer session.Close()

	//if err := session.Query(`INSERT INTO tweet (timeline, id, text) VALUES (?, ?, ?)`,
	//	"me", gocql.TimeUUID(), "hello world").Exec(); err != nil {
	//	log.Fatal(err)
	//}

	var id gocql.UUID
	var text string

	//if err := session.Query(`SELECT id, lastname FROM cyclist_alt_stats`).
	//	Consistency(gocql.One).Scan(&id, &text); err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println("Tweet:", id, text)

	iter := session.Query(`SELECT id, lastname FROM cyclist_alt_stats`).Iter()
	for iter.Scan(&id, &text) {
		fmt.Println("Cyclist stats:", id, text)
	}
	if err := iter.Close(); err != nil {
		log.Fatal(err)
	}
}
