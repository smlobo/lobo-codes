package main

import (
	"flag"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"subdomain-http-server/internal"
)

func main() {
	// Input arguments
	requestsDbPtr := flag.String("requestsdb", "amelia/requests.db", "db to read requests")
	allRecordsPtr := flag.Bool("all", true, "print all records")
	countryCountPtr := flag.Bool("countries", false, "select country, count group by country")
	flag.Parse()

	db, err := gorm.Open(sqlite.Open(*requestsDbPtr), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Read all records
	if *allRecordsPtr {
		var requestInfos []internal.RequestInfo
		result := db.Find(&requestInfos)
		fmt.Printf("Found: %d request\n", result.RowsAffected)
		for i, product := range requestInfos {
			fmt.Printf("[%d] %v\n", i, product)
		}
	}

	// Read all countries
	if *countryCountPtr {
		var uniqueCountries []internal.Country
		result := db.Table("request_infos").
			Select("country_short, country_long, count(country_short) as count").
			Group("country_short").
			Order("count").
			Find(&uniqueCountries)
		fmt.Printf("Found: %d request\n", result.RowsAffected)
		for i, country := range uniqueCountries {
			fmt.Printf("[%d] %v\n", i, country)
		}
	}
}
