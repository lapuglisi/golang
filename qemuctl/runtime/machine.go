package qemuctl_runtime

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func init() {

}

// Machine constants
const (
	MachineBaseDirectoryName string = "machines"
	MachineStatusFileName    string = "status"
	MachineStatusStarted     string = "started"
	MachineStatusStopped     string = "stopped"
	MachineStatusDegraded    string = "degraded"
	MachineStatusUnknown     string = "unknown"
	MachineConfigFileName    string = "config.yaml"
)

type Machine struct {
	Name             string
	Status           string
	RuntimeDirectory string
	ConfigFile       string
	initialized      bool
}

func NewMachine(machineName string) *Machine {
	var runtimeDirectory string = fmt.Sprintf("%s/%s/%s",
		GetUserDataDir(), MachineBaseDirectoryName, machineName)
	var statusFile string = fmt.Sprintf("%s/%s", runtimeDirectory, MachineStatusFileName)

	var fileData []byte
	var machineStatus string

	fileData, err := os.ReadFile(statusFile)
	if err != nil {
		log.Printf("error: could not open status file: %s\n", err.Error())
		machineStatus = MachineStatusDegraded
	} else {
		machineStatus = string(fileData)
	}

	configFile := fmt.Sprintf("%s/%s", runtimeDirectory, MachineConfigFileName)

	return &Machine{
		Name:             machineName,
		Status:           machineStatus,
		RuntimeDirectory: runtimeDirectory,
		ConfigFile:       configFile,
		initialized:      true,
	}
}

func (m *Machine) Exists() bool {
	fileInfo, err := os.Stat(m.RuntimeDirectory)
	if os.IsNotExist(err) {
		return false
	}

	return fileInfo.IsDir()
}

func (m *Machine) Destroy() bool {
	log.Printf("qemuctl: destroying machine %s\n", m.Name)

	err := os.RemoveAll(m.RuntimeDirectory)

	return err == nil
}

func (m *Machine) IsStarted() bool {
	return (strings.Compare(MachineStatusStarted, m.Status) == 0)
}

func (m *Machine) IsStopped() bool {
	return (strings.Compare(MachineStatusStopped, m.Status) == 0)
}

func (m *Machine) IsDegraded() bool {
	return (strings.Compare(MachineStatusDegraded, m.Status) == 0)
}

func (m *Machine) IsUnknown() bool {
	return (strings.Compare(MachineStatusUnknown, m.Status) == 0)
}

func (m *Machine) UpdateStatus(status string) (err error) {
	var statusFile string = fmt.Sprintf("%s/%s", m.RuntimeDirectory, MachineStatusFileName)
	var fileHandle *os.File

	log.Printf("[UpdateStatus] opening file '%s'\n", statusFile)
	fileHandle, err = os.OpenFile(statusFile, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0755)
	if err != nil {
		return err
	}

	switch status {
	case MachineStatusDegraded, MachineStatusStarted, MachineStatusStopped, MachineStatusUnknown:
		{
			log.Printf("[UpdateStatus] writing '%s' to file '%s'.\n", status, statusFile)
			_, err = fileHandle.WriteString(status)
			break
		}
	default:
		{
			err = fmt.Errorf("invalid machine status '%s'", status)
		}
	}
	fileHandle.Close()

	return err
}

func (m *Machine) CreateRuntime() {
	os.Mkdir(m.RuntimeDirectory, 0744)
}

func (m *Machine) UpdateConfigFile(sourcePath string) (err error) {
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	targetFile, err := os.Create(m.ConfigFile)
	if err != nil {
		return err
	}
	defer targetFile.Close()

	_, err = io.Copy(targetFile, sourceFile)

	return err
}
