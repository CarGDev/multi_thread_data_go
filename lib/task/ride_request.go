package task

import (
	"fmt"

	"ridesharing/lib/driver"
	"ridesharing/lib/location"
	"ridesharing/lib/result"
	"ridesharing/lib/rider"
)

const (
	baseFare     = 2.5
	farePerKm    = 1.2
	noDriverText = "no driver available"
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
	return &RideRequest{
		Base:        NewBase(taskID),
		rider:       r,
		pickup:      pickup,
		destination: destination,
	}
}

func (rr *RideRequest) Process() *result.Result {
	rr.status = Processing

	d := rr.FindDriver()
	if d == nil {
		rr.status = Failed
		return result.NewResult(rr.taskID, false, noDriverText)
	}
	defer d.CompleteRide()

	fare := rr.CalculateFare()
	rr.status = Completed

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
