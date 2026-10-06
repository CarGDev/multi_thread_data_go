// Package driver provides driver-related types and functionality.
package driver

import (
	"math"
	"sync"

	"ridesharing/lib/location"
)

type Driver struct {
	driverID        int
	name            string
	currentLocation int
	status          DriverStatus
}

// drivers is the full registry; availableDrivers holds only the drivers that
// can take a ride right now, keyed by driver ID. Both are guarded by driversMu.
var (
	drivers          []*Driver
	availableDrivers = map[int]*Driver{}
	driversMu        sync.Mutex
)

func NewDriver(
	driverID int,
	name string,
	latitude float64,
	longitude float64,
	address string,
) *Driver {
	loc := location.NewLocation(latitude, longitude, address)
	d := &Driver{
		driverID:        driverID,
		name:            name,
		currentLocation: loc.LocationID,
		status:          Available,
	}

	driversMu.Lock()
	drivers = append(drivers, d)
	availableDrivers[d.driverID] = d
	driversMu.Unlock()

	return d
}

// GetAll returns a copy of the driver registry.
func GetAll() []*Driver {
	driversMu.Lock()
	defer driversMu.Unlock()

	return append([]*Driver(nil), drivers...)
}

// GetAllAvailable returns a copy of the drivers that can take a ride now.
func GetAllAvailable() []*Driver {
	driversMu.Lock()
	defer driversMu.Unlock()

	values := make([]*Driver, 0, len(availableDrivers))
	for _, d := range availableDrivers {
		values = append(values, d)
	}
	return values
}

func (d *Driver) UpdateLocation(
	latitude float64,
	longitude float64,
	address string,
) {
	loc := location.NewLocation(latitude, longitude, address)

	driversMu.Lock()
	defer driversMu.Unlock()
	d.currentLocation = loc.LocationID
}

// AcquireNearest atomically reserves the available driver closest to pickup,
// or returns nil when none is available. The scan, the status change and the
// removal from the available set happen under one lock, so two workers can
// never get the same driver.
func AcquireNearest(pickup *location.Location) *Driver {
	driversMu.Lock()
	defer driversMu.Unlock()

	var nearest *Driver
	best := math.MaxFloat64
	for _, d := range availableDrivers {
		loc := location.GetLocation(d.currentLocation)
		if loc == nil {
			continue
		}
		if dist := loc.DistanceTo(pickup); dist < best {
			best, nearest = dist, d
		}
	}
	if nearest == nil {
		return nil
	}

	nearest.status = Busy
	delete(availableDrivers, nearest.driverID)
	return nearest
}

func (d *Driver) AcceptRide() bool {
	driversMu.Lock()
	defer driversMu.Unlock()
	if d.status != Available {
		return false
	}
	d.status = Busy
	delete(availableDrivers, d.driverID)
	return true
}

// CompleteRide frees the driver where they currently are.
func (d *Driver) CompleteRide() {
	driversMu.Lock()
	defer driversMu.Unlock()
	if d.status != Busy {
		return
	}
	d.status = Available
	availableDrivers[d.driverID] = d
}

// CompleteRideAt frees the driver and leaves them at the ride's destination.
func (d *Driver) CompleteRideAt(destination *location.Location) {
	driversMu.Lock()
	defer driversMu.Unlock()
	if d.status != Busy {
		return
	}
	d.currentLocation = destination.LocationID
	d.status = Available
	availableDrivers[d.driverID] = d
}

// GoOffline makes an available driver unavailable; a busy driver is left alone.
func (d *Driver) GoOffline() bool {
	driversMu.Lock()
	defer driversMu.Unlock()
	if d.status != Available {
		return false
	}
	d.status = Offline
	delete(availableDrivers, d.driverID)
	return true
}

func (d *Driver) GoOnline() {
	driversMu.Lock()
	defer driversMu.Unlock()
	if d.status == Offline {
		d.status = Available
		availableDrivers[d.driverID] = d
	}
}

func (d *Driver) GetID() int {
	return d.driverID
}

func (d *Driver) GetName() string {
	return d.name
}

func (d *Driver) GetLocation() *location.Location {
	driversMu.Lock()
	id := d.currentLocation
	driversMu.Unlock()

	return location.GetLocation(id)
}

func (d *Driver) GetStatus() DriverStatus {
	driversMu.Lock()
	defer driversMu.Unlock()

	return d.status
}
