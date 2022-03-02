package main

import (
	"flag"
	"fmt"
	"gprace/gprace"
	"os"
)

func usage() {
	fmt.Println("")
	flag.Usage()
	fmt.Println("")
}

func main() {
	var inFile string
	///
	/// gprace must be executed as follows:
	/// gprace INFILE
	///
	/// If no arguments are provided, warn the user and exit
	///
	flag.StringVar(&inFile, "race", "", "race: input file with race information")

	flag.Parse()

	if inFile == "" {
		usage()
		return
	}

	//
	// Check whether 'inFile' is a valid file
	//
	_, err := os.Stat(inFile)

	if os.IsNotExist(err) {
		fmt.Printf("[gprace] File '%s' does not exist.\n", inFile)
		return
	}

	///
	/// Set our gpRace variable to start a new race
	///
	gpRace := gprace.OnYourMarks()

	gpRace.GetSet(inFile)

	err = gpRace.Go()
	if err != nil {
		fmt.Printf("[!] Race not started: %v\n", err)
	}

	fmt.Printf("\n")
	fmt.Printf("[Race summary]\n")
	fmt.Printf("  Race file ....... %s\n", inFile)
	fmt.Printf("  Winner .......... %s\n", gpRace.AndTheWinnerIs().ToString())
	fmt.Printf("  Best Race Lap ... %s\n", gpRace.GetBestLap().ToString())
	fmt.Printf("\n")

	gpRace.PrintPodium()

}
