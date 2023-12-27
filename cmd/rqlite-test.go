package main

import (
	"fmt"
	"k8s-http-server/internal"
	"log"
	"time"
)

func main() {
	// config - for rqlite server
	internal.SetupConfig()

	// Init rqlite
	internal.InitRqlite()

	// Get the first 5 entries
	for i := 1; i <= 5; i++ {
		queryString := fmt.Sprintf("SELECT id,count,country_short,city FROM sheldon WHERE id = %d", i)
		rows, err := internal.RqliteQuery(queryString)
		if err != nil {
			log.Printf("WARNING: Error during lookup of id 1; %s", err.Error())
			return
		}

		// Should be only 1 existing
		if len(rows) != 1 {
			log.Printf("Multiple found: 1")
			return
		}

		log.Printf("[%d]: %s", i, rows[0])
	}

	queryString := fmt.Sprintf("SELECT id,count FROM sheldon WHERE id = 40")
	rows, _ := internal.RqliteQuery(queryString)
	log.Printf("Test: %s", rows[0])
	rowMap, ok := rows[0].(map[string]interface{})
	if !ok {
		log.Printf("result not string-int map: %T", rows[0])
	}
	log.Printf("  count: %d", int(rowMap["count"].(float64)))
	log.Printf("  id: %d", int(rowMap["id"].(float64)))

	newCount := int(rowMap["count"].(float64)) + 1000
	id := int(rowMap["id"].(float64))

	insertUpdate := fmt.Sprintf("UPDATE sheldon SET count = %d, updated_at = '%s' WHERE id = %d", newCount,
		time.Now().Format(time.RFC3339Nano), id)
	//insertUpdate := fmt.Sprintf("UPDATE sheldon SET count = %d WHERE id = %d", newCount, id)
	log.Printf("Update string: %s", insertUpdate)

	err := internal.RqliteExecute(insertUpdate)
	if err != nil {
		log.Printf(" Update err: %s", err)
	}
	log.Printf("INFO: Updated [count: %d, id: %d]", newCount, id)
}
