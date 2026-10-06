// Package simulation loads the sample drivers, riders and ride requests from embedded CSV files.
package simulation

import (
	"embed"
	"encoding/csv"
	"fmt"
	"strconv"

	"ridesharing/lib/driver"
	"ridesharing/lib/location"
	"ridesharing/lib/rider"
	"ridesharing/lib/task"
)

//go:embed Drivers.csv Riders.csv
var fileContent embed.FS

// readRows returns the data rows of an embedded CSV, with the header removed.
func readRows(name string, columns int) ([][]string, error) {
	file, err := fileContent.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	rows, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, err
	}

	if len(rows) < 2 {
		return nil, fmt.Errorf("%s there is not data rows", name)
	}

	for i, row := range rows[1:] {
		if len(row) < columns {
			return nil, fmt.Errorf("%s line %d: expected %d columns, got %d", name, i+2, columns, len(row))
		}
	}
	return rows[1:], nil
}

func parseFloat(name string, line int, field, value string) (float64, error) {
	f, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, fmt.Errorf("%s line %d: bad %s %q", name, line, field, value)
	}
	return f, nil
}

// LoadDrivers registers every driver found in Drivers.csv.
func LoadDrivers() error {
	const name = "Drivers.csv"
	rows, err := readRows(name, 5)
	if err != nil {
		return err
	}

	for i, row := range rows {
		line := i + 2
		id, err := strconv.Atoi(row[0])
		if err != nil {
			return fmt.Errorf("%s line %d: bad id %q", name, line, row[0])
		}
		lat, err := parseFloat(name, line, "latitude", row[2])
		if err != nil {
			return err
		}
		lon, err := parseFloat(name, line, "longitude", row[3])
		if err != nil {
			return err
		}

		driver.NewDriver(id, row[1], lat, lon, row[4])
	}

	fmt.Printf("loaded %d drivers\n", len(rows))
	return nil
}

// LoadRides registers every rider found in Riders.csv and returns one ride
// request per rider: the pickup is the rider's own location and the
// destination comes from the destination_* columns.
func LoadRides() ([]task.Task, error) {
	const name = "Riders.csv"

	rows, err := readRows(name, 8)
	if err != nil {
		return nil, err
	}

	tasks := make([]task.Task, 0, len(rows))
	for i, row := range rows {
		line := i + 2

		id, err := strconv.Atoi(row[0])
		if err != nil {
			return nil, fmt.Errorf("%s line %d: bad id %q", name, line, row[0])
		}
		lat, err := parseFloat(name, line, "latitude", row[2])
		if err != nil {
			return nil, err
		}
		lon, err := parseFloat(name, line, "longitude", row[3])
		if err != nil {
			return nil, err
		}
		destLat, err := parseFloat(name, line, "destination_lat", row[5])
		if err != nil {
			return nil, err
		}
		destLon, err := parseFloat(name, line, "destination_lon", row[6])
		if err != nil {
			return nil, err
		}

		r := rider.NewRider(id, row[1], lat, lon, row[4])
		destination := location.NewLocation(destLat, destLon, row[7])

		tasks = append(tasks, task.NewRideRequest(id, r, r.GetLocation(), destination))
	}

	fmt.Printf("loaded %d riders / ride requests\n", len(tasks))
	return tasks, nil
}
