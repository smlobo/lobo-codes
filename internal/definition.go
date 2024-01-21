package internal

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Country struct {
	CountryShort string
	Count        int
}

type City struct {
	City         string
	Region       string
	CountryShort string
	Count        int
}

type VisitorsPage struct {
	UniqueCountries int
	Countries       []Country
	Cities          []City
	PoweredBy       PoweredBy
}

type IndexPage struct {
	PoweredBy PoweredBy
}

type PoweredBy struct {
	GoVersion         string
	KubernetesVersion string
	OsVersion         string
	RqliteVersion     string
	PodName           string
	NodeName          string
}

type RequestInfo struct {
	Id                uuid.UUID
	CreatedAt         time.Time
	UpdatedAt         time.Time
	RemoteAddress     string
	OrigRemoteAddress string
	UserAgent         string
	Count             int64
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
	return fmt.Sprintf("[%s] <%s + %s> %s {%d} %s [%s]", ri.Id.String(), ri.CreatedAt, ri.UpdatedAt,
		ri.RemoteAddress, ri.Count, ri.UserAgent, ri.GeoLocation)
}

func (gl GeoLocation) String() string {
	return fmt.Sprintf("%s, %s, %s, %s", gl.CountryShort, gl.Region, gl.City, gl.Zipcode)
}
