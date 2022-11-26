package main

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"http-server/internal"
)

func main() {
	db, err := gorm.Open(sqlite.Open("requests.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Read all records
	var requestInfos []internal.RequestInfo
	result := db.Find(&requestInfos)
	fmt.Printf("Found: %d request\n", result.RowsAffected)
	for i, product := range requestInfos {
		fmt.Printf("[%d] %v\n", i, product)
	}
}
