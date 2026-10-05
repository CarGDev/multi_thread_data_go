// Package driver provides driver-related types and functionality.
package driver

import (
	"sync"

	"ridesharing/lib/location"
)

type Driver struct {
	driverID        int
	name            string
	currentLocation int
	status          DriverStatus
}

var (
	drivers   []*Driver
	driversMu sync.Mutex
)

func NewDriver(
	driverID int,
	name string,
	latitude float64,
	longitude float64,
	address string,
) *Driver {
	loc := location.NewLocation(latitude, longitude, address)
	driver := Driver{
		driverID:        driverID,
		name:            name,
		currentLocation: loc.LocationID,
		status:          Available,
	}

	driversMu.Lock()
	drivers = append(drivers, &driver)
	driversMu.Unlock()

	return &driver
}

func (d *Driver) UpdateLocation(
	latitude float64,
	longitude float64,
	address string,
) {
	loc := location.NewLocation(latitude, longitude, address)

	d.currentLocation = loc.LocationID
}

// AcquireAvailable atomically reserves the first available driver, or returns nil.
func AcquireAvailable() *Driver {
	driversMu.Lock()
	snapshot := append([]*Driver(nil), drivers...)
	driversMu.Unlock()

	for _, d := range snapshot {
		if d.AcceptRide() {
			return d
		}
	}
	return nil
}

func (d *Driver) AcceptRide() bool {
	driversMu.Lock()
	defer driversMu.Unlock()
	if d.status != Available {
		return false
	}
	d.status = Busy
	return true

}

func (d *Driver) CompleteRide() {
	driversMu.Lock()
	defer driversMu.Unlock()
	if d.status != Busy {
		return
	}
	d.status = Available
}

func (d *Driver) GetID() int {
	return d.driverID
}

func (d *Driver) GetName() string {
	return d.name
}

func (d *Driver) GetLocation() *location.Location {
	return location.GetLocation(d.currentLocation)
}

func (d *Driver) GetStatus() DriverStatus {
	return d.status
}
