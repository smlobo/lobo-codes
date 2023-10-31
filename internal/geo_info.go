package internal

import (
	"github.com/ip2location/ip2location-go/v9"
	"log"
)

var Db *ip2location.DB

func InitGeoDB(ip2locationDb string) error {
	var err error
	Db, err = ip2location.OpenDB(ip2locationDb)
	return err
}

func geoInfo(info *RequestInfo) {
	results, err := Db.Get_all(info.RemoteAddress)

	if err != nil {
		log.Printf("Error getting geo location for IP %s : %s", info.RemoteAddress, err)
		return
	}

	info.CountryShort = results.Country_short
	info.CountryLong = results.Country_long
	info.Region = results.Region
	info.City = results.City
	info.Latitude = results.Latitude
	info.Longitude = results.Longitude
	info.Zipcode = results.Zipcode
	info.Timezone = results.Timezone
	info.Elevation = results.Elevation
}
