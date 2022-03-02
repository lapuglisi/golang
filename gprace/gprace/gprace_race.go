package gprace

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
)

// Race struct
type Race struct {
	inFile string
}

// Ready function
func Ready() *Race {
	return new(Race)
}

// Set function
func (race *Race) Set(raceFile string) {
	race.inFile = raceFile
}

// Go method
func (race *Race) Go() error {
	fmt.Printf("[Race] Start race at '%s'.\n", race.inFile)

	return race.startRace()
}

///
/// Privates methods of GPRace
///
func (race *Race) startRace() error {

	var err error

	// Read inFile data into a byte array
	fileData, err := ioutil.ReadFile(race.inFile)

	if err != nil {
		return fmt.Errorf("could not start race: %v", err)
	}

	// Create an instance of a Byte reader to handle
	// the lines in our inFile
	br := bufio.NewReader(bytes.NewReader(fileData))

	for err == nil {
		lineBytes, err := br.ReadBytes('\n')

		// TODO: Assign false to a variable and get the
		// corresponding value of our err variable
		if err != nil {
			return fmt.Errorf("error while reading race data (%v)", err)
		}

		ParseLapEntry(lineBytes)
	}

	return nil
}
