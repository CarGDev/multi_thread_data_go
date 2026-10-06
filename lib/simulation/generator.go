// Package simulation will generate sample drivers, riders and ride requests.
package simulation

import (
	"embed"
	"encoding/csv"
	"fmt"
	"strconv"

	"ridesharing/lib/driver"
	"ridesharing/lib/rider"
)

//go:embed Drivers.csv Riders.csv
var fileContent embed.FS

type DriverInformation struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Address   string  `json:"address"`
}

type RiderInformation struct {
	ID                   int     `json:"id"`
	Name                 string  `json:"name"`
	Latitude             float64 `json:"latitude"`
	Longitude            float64 `json:"longitude"`
	Address              string  `json:"address"`
	DestinationLatitude  float64 `json:"dest_lat"`
	DestinationLongitude float64 `json:"dest_lon"`
	DestinationAddress   string  `json:"dest_address"`
}

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
		info := DriverInformation{ID: id, Name: row[1], Latitude: lat, Longitude: lon, Address: row[4]}

		driver.NewDriver(info.ID, info.Name, info.Latitude, info.Longitude, info.Address)

	}

	return nil
}

func LoadRides() error {
	const name = "Riders.csv"

	rows, err := readRows(name, 8)
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
		destLat, err := parseFloat(name, line, "destination_lat", row[5])
		if err != nil {
			return err
		}
		destLon, err := parseFloat(name, line, "destination_lon", row[6])
		if err != nil {
			return err
		}

		info := RiderInformation{
			ID:                   id,
			Name:                 row[1],
			Latitude:             lat,
			Longitude:            lon,
			Address:              row[4],
			DestinationLatitude:  destLat,
			DestinationLongitude: destLon,
			DestinationAddress:   row[7],
		}

		rider.NewRider(info.ID, info.Name, info.Latitude, info.Longitude, info.Address)
	}

	return nil
}
