package gprace

import (
	"fmt"
	"math"
	"time"
)

// Driver struct
type Driver struct {
	ID   int
	Name string
	Laps []RaceLap

	// private variables
	bestLap      *RaceLap
	bestLapSpeed float64

	averageRaceSpeed float64
}

// Typename for our list of drivers
type raceDriversT map[int]*Driver

// NewDriver is the constructor for the Driver 'class'
func NewDriver(id int, name string) (d *Driver) {
	d = new(Driver)

	d.ID = id
	d.Name = name

	d.averageRaceSpeed = 0.0
	d.bestLapSpeed = math.MaxFloat64

	return d
}

// ToString method
func (driver *Driver) ToString() (s string) {
	s = fmt.Sprintf("[ID: %d, Name: %s, BestLap: %s, AS: %.3f, Race Time: %s]",
		driver.ID, driver.Name, driver.bestLap.ToString(), driver.getAverageSpeed(),
		driver.getRaceTime().String())

	return s
}

//
// addRaceLap: adds a new race lap information for the current driver
//
func (driver *Driver) addRaceLap(lap RaceLap) {
	// Update the driver's race lap list
	driver.Laps = append(driver.Laps, lap)
	lenLaps := len(driver.Laps)

	// Calculate the best race lap for this driver, so far
	if lap.AverageSpeed < driver.bestLapSpeed {
		driver.bestLapSpeed = lap.AverageSpeed
		driver.bestLap = &(driver.Laps[lenLaps-1])
	}

	// Calculate the average speed of this driver for the current race
	driver.averageRaceSpeed += lap.AverageSpeed
}

//
// getAverageSpeed: calculates the average speed of a driver
// based on the race laps he/she completed
//
func (driver *Driver) getAverageSpeed() (as float64) {

	laps := float64(len(driver.Laps))
	return driver.averageRaceSpeed / laps
}

//
// getBestLap: calculates the best lap a driver completed
// based on the average speed of th race lap
//
func (driver *Driver) getBestLap() (bl int) {
	var lapSpeed float64
	var bestLap int

	lapSpeed = math.MaxFloat64
	bestLap = 0

	for _, lap := range driver.Laps {
		if lap.AverageSpeed < lapSpeed {
			lapSpeed = lap.AverageSpeed
			bestLap = lap.Number
		}
	}

	return bestLap
}

//
// getRaceTime: calculates the total time a driver
// spent in a race
//
func (driver *Driver) getRaceTime() time.Duration {
	var raceTime time.Duration
	for _, lap := range driver.Laps {
		raceTime += lap.TimeElapsed
	}

	return raceTime
}

///
/// updateRaceDrivers: gets or insert a driver from/onto our drivers list
///
func updateRaceDrivers(race *Race, id int, name string) (driver *Driver) {
	driver, ok := race.drivers[id]

	if !ok {
		driver = NewDriver(id, name)

		race.drivers[id] = driver
	}

	return driver
}

//
// getRaceDriver returns the race driver identified by id
//
func getRaceDriver(race *Race, id int) (driver *Driver, _ error) {
	driver, ok := race.drivers[id]

	if !ok {
		return nil, fmt.Errorf("driver with id %d not found", id)
	}

	return driver, nil
}
