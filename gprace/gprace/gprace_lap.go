package gprace

import (
	"bufio"
	"bytes"
	"fmt"
	"strconv"
	"time"
)

// RaceLap struct
type RaceLap struct {
	Driver       RaceDriver
	Number       int
	FinishedAt   time.Duration
	TimeElapsed  time.Duration
	AverageSpeed float64
}

func (lap *RaceLap) setDriver(id string, name string) {
	lap.Driver.ID, _ = strconv.Atoi(id)
	lap.Driver.Name = name
}

func (lap *RaceLap) setNumber(n string) {
	lap.Number, _ = strconv.Atoi(n)
}

func (lap *RaceLap) setFinishedAt(t string) {
	lap.FinishedAt, _ = time.ParseDuration(t)
}

func (lap *RaceLap) setDuration(d string) {
	lap.TimeElapsed, _ = time.ParseDuration(d)
}

func (lap *RaceLap) setAverageSpeed(speed string) {
	lap.AverageSpeed, _ = strconv.ParseFloat(speed, 32)
}

// ParseLapEntry is a function
func ParseLapEntry(entry []byte) (lap *RaceLap) {

	br := bytes.NewReader(entry)
	scan := bufio.NewScanner(br)

	scan.Split(bufio.ScanWords)

	gpLap := new(RaceLap)

	// First entry is Lap Time
	scan.Scan()
	gpLap.setFinishedAt(scan.Text())

	// Second entry is "Driver ID - Driver Name"
	scan.Scan()
	driverID := scan.Text()

	scan.Scan()
	scan.Scan()

	driverName := scan.Text()

	gpLap.setDriver(driverID, driverName)

	// Third entry is lapNumber
	scan.Scan()
	gpLap.setNumber(scan.Text())

	// Fourth entry is lapDuration
	scan.Scan()
	gpLap.setDuration(scan.Text())

	// Last entry is averageSpeed
	scan.Scan()
	gpLap.setAverageSpeed(scan.Text())

	fmt.Printf("[ParseLapEntry] Current LAP:\n")
	fmt.Printf("{\n")
	fmt.Printf("  Time  ......... %v\n", gpLap.FinishedAt)
	fmt.Printf("  Driver ID ..... %d\n", gpLap.Driver.ID)
	fmt.Printf("  Driver Name ... %s\n", gpLap.Driver.Name)
	fmt.Printf("  Lap Number .... %d\n", gpLap.Number)
	fmt.Printf("  Lap Duration .. %v\n", gpLap.TimeElapsed)
	fmt.Printf("  Avg Speed ..... %v\n", gpLap.AverageSpeed)
	fmt.Printf("}\n\n")

	// This entry must be ended by a new-line char.
	// If it is not, so line is malformatted
	fmt.Printf("[gprace] entry is '%s'\n", scan.Text())

	return nil
}
