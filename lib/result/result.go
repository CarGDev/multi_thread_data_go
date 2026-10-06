// Package result provides the outcome of a processed task and a store for them.
package result

import (
	"ridesharing/lib/location"
	"time"
)

// RideInfo is the detail of a ride attached to a Result. Locations are nil
// when they do not apply (for example, no driver was found).
type RideInfo struct {
	RiderID          int
	RiderName        string
	Pickup           *location.Location
	Destination      *location.Location
	DriverID         int
	DriverName       string
	DriverStart      *location.Location // where the driver was when assigned
	DriverEnd        *location.Location // where the driver ended the ride
	DriverToPickupKm float64
	TripDistanceKm   float64
	Fare             float64
}

type Result struct {
	taskID      int
	workerID    int
	success     bool
	message     string
	ride        RideInfo
	startedAt   time.Time
	completedAt time.Time
}

func NewResult(taskID int, success bool, message string) *Result {
	return &Result{
		taskID:      taskID,
		success:     success,
		message:     message,
		completedAt: time.Now(),
	}
}

func (r *Result) GetTaskID() int {
	return r.taskID
}

func (r *Result) IsSuccess() bool {
	return r.success
}

func (r *Result) GetMessage() string {
	return r.message
}

func (r *Result) GetCompletedAt() time.Time {
	return r.completedAt
}

func (r *Result) SetWorker(workerID int) {
	r.workerID = workerID
}

func (r *Result) GetWorkerID() int {
	return r.workerID
}

func (r *Result) SetStartedAt(t time.Time) {
	r.startedAt = t
}

func (r *Result) GetStartedAt() time.Time {
	return r.startedAt
}

// GetDuration returns how long the task took, or 0 if the start is unknown.
func (r *Result) GetDuration() time.Duration {
	if r.startedAt.IsZero() {
		return 0
	}
	return r.completedAt.Sub(r.startedAt)
}

func (r *Result) SetRide(info RideInfo) {
	r.ride = info
}

func (r *Result) GetRide() RideInfo {
	return r.ride
}
