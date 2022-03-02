package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"wslex/wslexdef"
	"wslex/wslexreg"
)

func usage() {
	fmt.Println("")
	fmt.Println("usage:")
	fmt.Println("    wslex DISTRO_PATH")
	fmt.Println("")
}

func main() {
	// In this first version, we get the path to the distro
	// and save it in the registry
	if len(os.Args) < 2 {
		usage()
		return
	}

	var distroPath = os.Args[1]
	distroPath = strings.TrimRight(distroPath, "\\")
	var appManifest = distroPath + "\\" + wslexdef.AppxManifestXML

	// Check if distroPath contains a valid distro
	xmlData, err := ioutil.ReadFile(appManifest)
	if err != nil {
		fmt.Printf("[wslex] Could not open '%s': %v\n", appManifest, err)
		return
	}

	var wslData wslexdef.WslXMLRoot
	err = xml.Unmarshal(xmlData, &wslData)
	if err != nil {
		fmt.Printf("[wslex] Could not read xml: %v\n", err)
		return
	}

	// Now we should have the executable for the distro
	wslDistro, err := wslexreg.WslLoadDistro(distroPath, wslData.Applications.App.Id)
	if err != nil {
		fmt.Printf("[wslex] Could not load distro: %v\n", err)
		return
	}

	wslDistro.Executable = wslData.Applications.App.Executable

	// Dump distro
	fmt.Printf("[wslex] Loaded distro for '%s': \n", distroPath)
	fmt.Printf("[wslex] Distribution Name: '%s'. \n", wslDistro.Name)
	fmt.Printf("[wslex]       Default UID: %d \n", wslDistro.DefaultUID)
	fmt.Printf("[wslex]        Executable: %s \n", wslDistro.Executable)
	fmt.Printf("[wslex]        Is Default: %v\n", wslDistro.Default)

}
