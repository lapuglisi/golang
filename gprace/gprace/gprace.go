package gprace

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"time"
)

// Race struct
type Race struct {
	inFile string

	drivers raceDriversT

	Laps []RaceLap
}

// OnYourMarks function
func OnYourMarks() *Race {
	var race = new(Race)

	race.drivers = make(raceDriversT)

	return race
}

// GetSet function
func (race *Race) GetSet(raceFile string) {
	race.inFile = raceFile
}

// Go method
func (race *Race) Go() error {
	return race.startRace()
}

//
// GetBestLap returns the best race lap for the current race
// based on the race lap average speed
//
func (race *Race) GetBestLap() (lap *RaceLap) {
	var lapSpeed float64

	lapSpeed = math.MaxFloat64

	for _, raceLap := range race.Laps {
		if raceLap.AverageSpeed < lapSpeed {
			lapSpeed = raceLap.AverageSpeed
			lap = &raceLap
		}
	}

	return lap
}

// PrintPodium prints the result of the current race
func (race *Race) PrintPodium() {

	//
	// Initialize a drivers slice to be sorted and give us the actual podium
	// - Time complexity here is something around:
	// - O(n) to copy from 'race.drivers' to 'drivers'
	// - O(n * log(n)) to bubble-sort them
	//   Amortized O(n * log(n))
	//
	dlen := len(race.drivers)
	drivers := []*Driver{}

	// Copy from 'race.drivers' to 'drivers' (Time complexity: O(n))
	for _, dr := range race.drivers {
		drivers = append(drivers, dr)
	}

	// Bubble sort our drivers slice (amortized O(n * log(n)))
	for idx := 0; idx < dlen; idx++ {
		swapped := false
		for jdx := 0; jdx < dlen-1; jdx++ {
			if drivers[jdx].getRaceTime() > drivers[jdx+1].getRaceTime() {
				swapMe := drivers[jdx+1]
				drivers[jdx+1] = drivers[jdx]
				drivers[jdx] = swapMe

				swapped = true
			}
		}

		if !swapped {
			break
		}
	}

	fmt.Println()
	fmt.Println("##########################################################")
	fmt.Println("##### GP Race : Podium                               #####")
	fmt.Println("##########################################################")

	driver := drivers[0]
	fmt.Printf("[#1] %s\n", driver.Name)
	for idx := 1; idx < dlen; idx++ {
		currentDriver := drivers[idx]
		timeBehind := currentDriver.getRaceTime() - driver.getRaceTime()

		fmt.Printf("[#%d] %-15s (%s behind %s)\n",
			idx+1, currentDriver.Name, timeBehind.String(), driver.Name)
	}

	fmt.Println("##########################################################")
	fmt.Println()

}

// AndTheWinnerIs returns the winner for the current race
func (race *Race) AndTheWinnerIs() (driver *Driver) {
	var raceTime time.Duration

	raceTime = math.MaxInt64

	for _, d := range race.drivers {
		if d.getRaceTime() < raceTime {
			raceTime = d.getRaceTime()
			driver = d
		}
	}

	return driver
}

//
// raceOnTheGo: Goroutine to listen for a completed race lap
///
func (race *Race) raceOnTheGo(laps chan RaceLap) {
	for lap := range laps {
		// Update the list of race laps for the current race
		race.Laps = append(race.Laps, lap)

		// Check if the driver is already in our drivers list
		driver := updateRaceDrivers(race, lap.DriverID, lap.DriverName)

		// Append the current race lap to the list of race laps for
		// a specific driver
		driver.addRaceLap(lap)
	}
}

func (race *Race) startRace() error {

	var err error

	// Read inFile data into a byte array
	file, err := os.OpenFile(race.inFile, os.O_RDONLY, 0755)

	if err != nil {
		return fmt.Errorf("could not start race: %v", err)
	}

	// Set up a channel to exchange info between Lap and Race
	var lin chan (string)
	var lout chan RaceLap

	lin = make(chan string)
	lout = make(chan RaceLap)

	// Start our lap listener goroutine
	go startLapListener(lin, lout)

	// Start raceOnTheGo goroutine
	go race.raceOnTheGo(lout)

	// Create an instance of a scanner to handle
	// the lines in our inFile
	scanner := bufio.NewScanner(file)

	// Set up our scanner to read line by line
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {

		lin <- scanner.Text()

	}

	return err
}
