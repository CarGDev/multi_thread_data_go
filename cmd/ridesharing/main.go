package main

import (
	"fmt"
	"os"

	"ridesharing/lib/simulation"
)

const workerCount = 4

func main() {
	if err := simulation.LoadDrivers(); err != nil {
		fmt.Println("loading drivers:", err)
		os.Exit(1)
	}

	if err := simulation.LoadRides(); err != nil {
		fmt.Println("loading rides:", err)
		os.Exit(1)
	}
}
