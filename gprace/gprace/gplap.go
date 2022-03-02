package gprace

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// RaceLap struct
type RaceLap struct {
	DriverID     int
	DriverName   string
	Number       int
	FinishedAt   time.Time
	TimeElapsed  time.Duration
	AverageSpeed float64
}

// ToString function
func (lap *RaceLap) ToString() (s string) {
	s = fmt.Sprintf("{%s, %d, %s, %s, %.3f}", lap.DriverName, lap.Number,
		lap.FinishedAt.Format("15:04:05.000"), lap.TimeElapsed.String(),
		lap.AverageSpeed)

	return s
}

func (lap *RaceLap) reset() {
	lap.Number = 0
	lap.AverageSpeed = 0
}

func (lap *RaceLap) setDriver(id string, name string) {
	lap.DriverID, _ = strconv.Atoi(id)
	lap.DriverName = name
}

func (lap *RaceLap) setNumber(n string) {
	lap.Number, _ = strconv.Atoi(n)
}

func (lap *RaceLap) setFinishedAt(t string) error {
	var err error
	lap.FinishedAt, err = time.Parse("15:04:05.000", t)

	return err
}

func (lap *RaceLap) setDuration(d string) error {
	var err error
	d = strings.Replace(d, ":", "m", 1)
	d = strings.Replace(d, ".", "s", 1)
	d = d + "ms"
	lap.TimeElapsed, err = time.ParseDuration(d)

	return err
}

func (lap *RaceLap) setAverageSpeed(speed string) error {
	var err error
	speed = strings.Replace(speed, ",", ".", -1)
	lap.AverageSpeed, err = strconv.ParseFloat(speed, 64)

	return err
}

func startLapListener(lin chan string, lout chan RaceLap) {

	for entry := range lin {

		lap, err := lapFromString(entry)

		if err == nil {
			lout <- lap
		}
	}

}

// ParseLapEntry is a function
func lapFromString(entry string) (lap RaceLap, err error) {

	// Set up a reader to scan for entries
	br := strings.NewReader(entry)

	// Set up our scanner to read each field in lap's info
	scan := bufio.NewScanner(br)
	scan.Split(bufio.ScanWords)

	// First entry in line is Lap Timestamp
	scan.Scan()
	err = lap.setFinishedAt(scan.Text())

	// If some error occured while trying to parse the lap timestamp
	// (i.e. malformatted line), we return
	if err != nil {
		lap.reset()
		return lap, err
	}

	// Second entry is "Driver ID - Driver Name"
	scan.Scan()
	driverID := scan.Text()

	scan.Scan()
	scan.Scan()

	driverName := scan.Text()

	lap.setDriver(driverID, driverName)

	// Third entry is lapNumber
	scan.Scan()
	lap.setNumber(scan.Text())

	// Fourth entry is lapDuration
	scan.Scan()
	lap.setDuration(scan.Text())

	// Last entry is averageSpeed
	scan.Scan()
	lap.setAverageSpeed(scan.Text())

	return lap, nil
}
