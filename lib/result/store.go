package result

import (
	"encoding/csv"
	"os"
	"strconv"
	"sync"

	"ridesharing/lib/errors"
	"ridesharing/lib/location"
)

const timeLayout = "2006-01-02T15:04:05.000Z07:00"

var csvHeader = []string{
	"task_id", "worker_id", "rider_id", "rider_name",
	"pickup_lat", "pickup_lon", "pickup_address",
	"dest_lat", "dest_lon", "dest_address",
	"driver_id", "driver_name",
	"driver_start_lat", "driver_start_lon", "driver_start_address",
	"driver_end_lat", "driver_end_lon", "driver_end_address",
	"driver_to_pickup_km", "trip_distance_km", "fare",
	"success", "message", "started_at", "completed_at", "duration_ms",
}

type ResultStore struct {
	results []*Result
	lock    sync.Mutex
}

func NewResultStore() *ResultStore {
	return &ResultStore{}
}

func (s *ResultStore) AddResult(r *Result) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.results = append(s.results, r)
}

// GetResults returns a copy of the stored results.
func (s *ResultStore) GetResults() []*Result {
	s.lock.Lock()
	defer s.lock.Unlock()

	return append([]*Result(nil), s.results...)
}

func (s *ResultStore) Count() int {
	s.lock.Lock()
	defer s.lock.Unlock()

	return len(s.results)
}

func (s *ResultStore) CountSuccess() int {
	s.lock.Lock()
	defer s.lock.Unlock()

	n := 0
	for _, r := range s.results {
		if r.success {
			n++
		}
	}
	return n
}

func (s *ResultStore) CountFailed() int {
	return s.Count() - s.CountSuccess()
}

// locationFields returns lat, lon and address, or three empty fields for nil.
func locationFields(loc *location.Location) []string {
	if loc == nil {
		return []string{"", "", ""}
	}
	return []string{
		strconv.FormatFloat(loc.Latitude, 'f', 6, 64),
		strconv.FormatFloat(loc.Longitude, 'f', 6, 64),
		loc.Address,
	}
}

func (r *Result) csvRecord() []string {
	ride := r.ride

	riderID, riderName := "", ""
	if ride.Pickup != nil {
		riderID, riderName = strconv.Itoa(ride.RiderID), ride.RiderName
	}

	driverID, driverName, toPickup, trip, fare := "", "", "", "", ""
	if ride.DriverStart != nil {
		driverID, driverName = strconv.Itoa(ride.DriverID), ride.DriverName
		toPickup = strconv.FormatFloat(ride.DriverToPickupKm, 'f', 3, 64)
		trip = strconv.FormatFloat(ride.TripDistanceKm, 'f', 3, 64)
		fare = strconv.FormatFloat(ride.Fare, 'f', 2, 64)
	}

	startedAt, duration := "", ""
	if !r.startedAt.IsZero() {
		startedAt = r.startedAt.Format(timeLayout)
		duration = strconv.FormatInt(r.GetDuration().Milliseconds(), 10)
	}

	record := []string{strconv.Itoa(r.taskID), strconv.Itoa(r.workerID), riderID, riderName}
	record = append(record, locationFields(ride.Pickup)...)
	record = append(record, locationFields(ride.Destination)...)
	record = append(record, driverID, driverName)
	record = append(record, locationFields(ride.DriverStart)...)
	record = append(record, locationFields(ride.DriverEnd)...)
	return append(record,
		toPickup, trip, fare,
		strconv.FormatBool(r.success), r.message,
		startedAt, r.completedAt.Format(timeLayout), duration,
	)
}

// WriteToCSV writes one row per result, in completion order.
func (s *ResultStore) WriteToCSV(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return errors.NewFileIOError(err.Error())
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	if err := writer.Write(csvHeader); err != nil {
		return errors.NewFileIOError(err.Error())
	}
	for _, r := range s.GetResults() {
		if err := writer.Write(r.csvRecord()); err != nil {
			return errors.NewFileIOError(err.Error())
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return errors.NewFileIOError(err.Error())
	}
	return nil
}
