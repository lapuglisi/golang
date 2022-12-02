package qemuctl_runtime

import (
	"fmt"
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
)

type Machine struct {
	Name             string
	Status           string
	RuntimeDirectory string
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

	return &Machine{
		Name:             machineName,
		Status:           machineStatus,
		RuntimeDirectory: runtimeDirectory,
		initialized:      true,
	}
}

func (m *Machine) Exists() bool {
	fileInfo, _ := os.Stat(m.RuntimeDirectory)
	return fileInfo.IsDir()
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

	switch status {
	case MachineStatusDegraded:
	case MachineStatusStarted:
	case MachineStatusStopped:
	case MachineStatusUnknown:
		{
			err = os.WriteFile(statusFile, []byte(status), os.ModePerm)
			break
		}

	default:
		{
			err = fmt.Errorf("invalid machine status '%s'", status)
		}
	}

	return err
}