package wslexreg

import (
	"fmt"
	"strings"
	"wslex/wslexdef"

	"github.com/google/uuid"
	"golang.org/x/sys/windows/registry"
)

// wslLoadDistro func
func wslLoadDistro(distroPath string) (distro wslexdef.WslDistro, err error) {

	var wslDistro wslexdef.WslDistro

	baseKey, err := registry.OpenKey(registry.CURRENT_USER, wslexdef.WslexLxssKey,
		registry.QUERY_VALUE|registry.ENUMERATE_SUB_KEYS)
	if err != nil {
		return wslDistro, fmt.Errorf("Could not open registry key: %v", err)
	}

	// Retrieve value for default distribution
	defaultDistro, _, err := baseKey.GetStringValue(wslexdef.DefaultDistribution)
	if err != nil {
		return wslDistro, fmt.Errorf("Could not determine default distro: %v", err)
	}

	// Retrieves all distributions
	keys, err := baseKey.ReadSubKeyNames(0)
	if err != nil {
		return wslDistro, fmt.Errorf("Could not enumerate sub keys: %v", err)
	}

	err = fmt.Errorf("Distro not found")

	// Iterate through the distros
	var distroFound = false
	for _, keyName := range keys {

		distro.Default = (strings.Compare(defaultDistro, keyName) == 0)

		// Open sub key
		localKey, err := registry.OpenKey(baseKey, keyName, registry.QUERY_VALUE)
		if err != nil {
			return wslDistro, fmt.Errorf("Could not open sub keys: %v", err)
		}

		// Compare current key[BasePath] to distroPath
		basePath, _, err := localKey.GetStringValue(wslexdef.BasePath)
		if err != nil {
			return wslDistro, fmt.Errorf("[WslLoadDistro] GetBasePath: %v", err)
		}

		if strings.Compare(basePath, distroPath) == 0 {

			// Found distro entry
			wslDistro.BasePath = basePath
			wslDistro.Default = (strings.Compare(keyName, defaultDistro) == 0)
			wslDistro.DefaultUID, _, err = localKey.GetIntegerValue(wslexdef.DefaultUID)
			wslDistro.Name, _, err = localKey.GetStringValue(wslexdef.DistributionName)
			wslDistro.Flags, _, err = localKey.GetIntegerValue(wslexdef.Flags)
			wslDistro.State, _, err = localKey.GetIntegerValue(wslexdef.State)
			wslDistro.Version, _, err = localKey.GetIntegerValue(wslexdef.Version)

			distroFound = true
		}

		localKey.Close()

		if distroFound {
			return wslDistro, nil
		}
	}

	return wslDistro, err
}

func wslRegisterDistro(distroPath string, distroId string) (wslDistro wslexdef.WslDistro, err error) {

	baseKey, err := registry.OpenKey(registry.CURRENT_USER, wslexdef.WslexLxssKey,
		registry.SET_VALUE|registry.QUERY_VALUE|registry.ENUMERATE_SUB_KEYS)

	if err != nil {
		return wslDistro, err
	}
	defer baseKey.Close()

	// Add new entry for tthe current distro
	guid, err := uuid.NewUUID()
	if err != nil {
		return wslDistro, err
	}

	distroKey, _, err := registry.CreateKey(baseKey, guid.String(), registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		return wslDistro, err
	}
	defer distroKey.Close()

	// Set default values for distro
	wslDistro.BasePath = distroPath
	wslDistro.Name = distroId
	wslDistro.Default = false
	wslDistro.DefaultUID = 1000
	wslDistro.Flags = 7
	wslDistro.State = 1
	wslDistro.Version = 1

	err = nil

	err = distroKey.SetStringValue(wslexdef.BasePath, wslDistro.BasePath)
	err = distroKey.SetStringValue(wslexdef.DistributionName, wslDistro.Name)
	err = distroKey.SetQWordValue(wslexdef.DefaultUID, wslDistro.DefaultUID)
	err = distroKey.SetQWordValue(wslexdef.Flags, wslDistro.Flags)
	err = distroKey.SetQWordValue(wslexdef.State, wslDistro.State)
	err = distroKey.SetQWordValue(wslexdef.Version, wslDistro.Version)

	return wslDistro, err
}

// WslLoadDistro ...
func WslLoadDistro(distroPath string, distroId string) (wslDistro wslexdef.WslDistro, err error) {

	wslDistro, err = wslLoadDistro(distroPath)
	if err != nil {
		wslDistro, err = wslRegisterDistro(distroPath, distroId)
	}

	return wslDistro, err
}
