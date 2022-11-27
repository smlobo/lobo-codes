package internal

import (
	"gorm.io/gorm"
	"log"
	"net/http"
	"strings"
)

func requestInfo(request *http.Request, db *gorm.DB) {
	var info RequestInfo

	info.UserAgent = request.Header.Get("User-Agent")
	remoteAddress := request.RemoteAddr

	// Once the information is extracted from the request, the remainder of processing can be
	// done concurrently
	go func() {
		info.RemoteAddress = strings.Trim(strings.Split(remoteAddress, ":")[0], "[]")

		// TODO: IP (v6?) address not found, log and skip
		if info.RemoteAddress == "" {
			log.Printf("WARN: IP addr not parsed: %s\n", remoteAddress)
			return
		}

		// Find a pre-existing IP Address & UserAgent entry
		var prevInfo RequestInfo
		//dbTx := db.Where(&info).First(&prevInfo)
		dbTx := db.Where(&info).Limit(1).Find(&prevInfo)

		if dbTx.RowsAffected != 0 {
			prevInfo.Count += 1
			db.Save(&prevInfo)
			log.Printf("UPDATE: %s", prevInfo)
		} else {
			info.Count = 1
			geoInfo(&info)
			db.Create(&info)
			log.Printf("CREATE: %s", info)
		}
	}()
}
