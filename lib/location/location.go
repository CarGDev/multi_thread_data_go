// Package location provides location-related types and functionality
package location

import (
	"fmt"
	"math"
	"sync"
)

const earthRadiusKm = 6371.0

type Location struct {
	LocationID int
	Latitude   float64
	Longitude  float64
	Address    string
}

var (
	locations   []*Location
	locationsMu sync.Mutex
)

func NewLocation(
	latitude float64,
	longitude float64,
	address string,
) *Location {
	locationsMu.Lock()
	defer locationsMu.Unlock()

	// For now create a new location every time, regardless if the address exist
	id := len(locations) + 1
	location := &Location{
		LocationID: id,
		Latitude:   latitude,
		Longitude:  longitude,
		Address:    address,
	}

	locations = append(locations, location)

	return location
}

// DistanceTo returns the great-circle distance in kilometers (haversine).
func (l *Location) DistanceTo(other *Location) float64 {
	toRad := func(deg float64) float64 { return deg * math.Pi / 180 }

	dLat := toRad(other.Latitude - l.Latitude)
	dLon := toRad(other.Longitude - l.Longitude)
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(toRad(l.Latitude))*math.Cos(toRad(other.Latitude))*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	return earthRadiusKm * 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}

// Equals reports whether both locations share the same coordinates.
func (l *Location) Equals(other *Location) bool {
	return l.Latitude == other.Latitude && l.Longitude == other.Longitude
}

func (l *Location) String() string {
	return fmt.Sprintf("%s (%.5f, %.5f)", l.Address, l.Latitude, l.Longitude)
}

func GetLocation(id int) *Location {
	locationsMu.Lock()
	defer locationsMu.Unlock()

	if id < 1 || id > len(locations) {
		return nil
	}
	return locations[id-1]
}
