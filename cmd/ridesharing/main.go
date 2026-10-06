package main

import (
	"fmt"
	"os"
	"time"

	"ridesharing/lib/simulation"
	"ridesharing/lib/system"
)

const (
	workerCount = 8
	resultsFile = "results.txt"
)

func main() {
	if err := simulation.LoadDrivers(); err != nil {
		fmt.Println("loading drivers:", err)
		os.Exit(1)
	}

	tasks, err := simulation.LoadRides()
	if err != nil {
		fmt.Println("loading rides:", err)
		os.Exit(1)
	}

	sys := system.Initialize(workerCount)
	defer sys.Close()

	start := time.Now()
	if err := sys.Run(tasks, resultsFile); err != nil {
		fmt.Println("running:", err)
		sys.Close()
		os.Exit(1)
	}

	stats := sys.Stats()
	fmt.Printf("done in %s: %d total, %d successful, %d failed (results in %s)\n",
		time.Since(start).Round(time.Millisecond), stats.Total, stats.Success, stats.Failed, resultsFile)
}
