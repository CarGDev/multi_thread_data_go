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
	minWorkMs    = 100
	maxWorkMs    = 500
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

	// Simulate the computational work of matching and routing.
	time.Sleep(time.Duration(minWorkMs+rand.Intn(maxWorkMs-minWorkMs)) * time.Millisecond)

	d := rr.FindDriver()
	if d == nil {
		rr.SetStatus(Failed)
		return result.NewResult(rr.taskID, false, noDriverText)
	}
	defer d.CompleteRide()

	fare := rr.CalculateFare()
	rr.SetStatus(Completed)

	return result.NewResult(rr.taskID, true, fmt.Sprintf(
		"rider %s assigned to driver %s, fare %.2f",
		rr.rider.GetName(), d.GetName(), fare,
	))
}

// FindDriver reserves an available driver; the caller must call CompleteRide on it.
func (rr *RideRequest) FindDriver() *driver.Driver {
	return driver.AcquireAvailable()
}

func (rr *RideRequest) CalculateFare() float64 {
	return baseFare + farePerKm*rr.pickup.DistanceTo(rr.destination)
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
