package internal

import (
	"fmt"
	"gorm.io/gorm"
)

type RequestInfo struct {
	gorm.Model
	RemoteAddress string
	UserAgent     string
	GeoLocation
}

type GeoLocation struct {
	CountryShort string
	CountryLong  string
	Region       string
	City         string
	Latitude     float32
	Longitude    float32
	Zipcode      string
	Timezone     string
	Elevation    float32
}

func (ri RequestInfo) String() string {
	return fmt.Sprintf("<%s> %s : %s [%s]", ri.CreatedAt, ri.RemoteAddress, ri.UserAgent, ri.GeoLocation)
}

func (gl GeoLocation) String() string {
	return fmt.Sprintf("%s, %s, %s, %s", gl.CountryShort, gl.Region, gl.City, gl.Zipcode)
}
