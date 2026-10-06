package task

import (
	"fmt"
	"math/rand"
	"time"

	"ridesharing/lib/driver"
	"ridesharing/lib/location"
	"ridesharing/lib/result"
	"ridesharing/lib/rider"
)

const (
	baseFare     = 2.5
	farePerKm    = 1.2
	noDriverText = "no driver available"
	minWorkMs    = 50
	maxWorkMs    = 200
)

type RideRequest struct {
	Base
	rider       *rider.Rider
	pickup      *location.Location
	destination *location.Location
}

func NewRideRequest(
	taskID int,
	r *rider.Rider,
	pickup *location.Location,
	destination *location.Location,
) *RideRequest {
	rr := &RideRequest{
		rider:       r,
		pickup:      pickup,
		destination: destination,
	}
	rr.Init(taskID)
	return rr
}

func (rr *RideRequest) Process() *result.Result {
	rr.SetStatus(Processing)

	d := rr.FindDriver()
	if d == nil {
		rr.SetStatus(Failed)

		res := result.NewResult(rr.taskID, false, noDriverText)
		res.SetRide(result.RideInfo{
			RiderID:     rr.rider.GetID(),
			RiderName:   rr.rider.GetName(),
			Pickup:      rr.pickup,
			Destination: rr.destination,
		})
		return res
	}
	// The driver ends the ride at the rider's destination.
	defer d.CompleteRideAt(rr.destination)

	// Where the driver was when assigned, before the ride moves them.
	driverStart := d.GetLocation()

	// Simulate the ride; the driver stays busy meanwhile.
	time.Sleep(time.Duration(minWorkMs+rand.Intn(maxWorkMs-minWorkMs)) * time.Millisecond)

	fare := rr.CalculateFare()
	rr.SetStatus(Completed)

	res := result.NewResult(rr.taskID, true, fmt.Sprintf(
		"rider %s assigned to driver %s, fare %.2f",
		rr.rider.GetName(), d.GetName(), fare,
	))
	res.SetRide(result.RideInfo{
		RiderID:          rr.rider.GetID(),
		RiderName:        rr.rider.GetName(),
		Pickup:           rr.pickup,
		Destination:      rr.destination,
		DriverID:         d.GetID(),
		DriverName:       d.GetName(),
		DriverStart:      driverStart,
		DriverEnd:        rr.destination,
		DriverToPickupKm: driverStart.DistanceTo(rr.pickup),
		TripDistanceKm:   rr.TripDistanceKm(),
		Fare:             fare,
	})
	return res
}

// FindDriver reserves the available driver nearest to the pickup; the caller
// must free it with CompleteRide or CompleteRideAt.
func (rr *RideRequest) FindDriver() *driver.Driver {
	return driver.AcquireNearest(rr.pickup)
}

// TripDistanceKm is the straight-line distance from pickup to destination.
func (rr *RideRequest) TripDistanceKm() float64 {
	return rr.pickup.DistanceTo(rr.destination)
}

func (rr *RideRequest) CalculateFare() float64 {
	return baseFare + farePerKm*rr.TripDistanceKm()
}

func (rr *RideRequest) GetRider() *rider.Rider {
	return rr.rider
}

func (rr *RideRequest) GetPickup() *location.Location {
	return rr.pickup
}

func (rr *RideRequest) GetDestination() *location.Location {
	return rr.destination
}
