// Package rider provides rider-related types and functionality
package rider

import "ridesharing/lib/location"

type Rider struct {
	riderID         int
	name            string
	currentLocation int
}

var riders []Rider

func NewRider(
	riderID int,
	name string,
	latitude float64,
	longitude float64,
	address string,
) *Rider {
	loc := location.NewLocation(latitude, longitude, address)
	rider := Rider{
		riderID:         riderID,
		name:            name,
		currentLocation: loc.LocationID,
	}

	riders = append(riders, rider)
	return &rider

}

func (d *Rider) UpdateLocation(
	latitude float64,
	longitude float64,
	address string,
) {
	loc := location.NewLocation(latitude, longitude, address)

	d.currentLocation = loc.LocationID
}

func (d *Rider) GetID() int {
	return d.riderID
}

func (d *Rider) GetName() string {
	return d.name
}

func (d *Rider) GetLocation() *location.Location {
	return location.GetLocation(d.currentLocation)
}
