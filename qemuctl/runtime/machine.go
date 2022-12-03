package qemuctl_runtime

import (
	"encoding/json"
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
	MachineDataFileName      string = "machine-data.json"
	MachineStatusStarted     string = "started"
	MachineStatusStopped     string = "stopped"
	MachineStatusDegraded    string = "degraded"
	MachineStatusUnknown     string = "unknown"
	MachineConfigFileName    string = "config.yaml"
)

type MachineData struct {
	QemuPid      int    `json:"qemuProcessPID"`
	State        string `json:"machineState"`
	SSHLocalPort int    `json:"sshLocalPort"`
}

type Machine struct {
	Name             string
	Status           string
	QemuPid          int
	SSHLocalPort     int
	RuntimeDirectory string
	ConfigFile       string
	initialized      bool
}

func NewMachine(machineName string) *Machine {
	var runtimeDirectory string = fmt.Sprintf("%s/%s/%s",
		GetUserDataDir(), MachineBaseDirectoryName, machineName)
	var dataFile string = fmt.Sprintf("%s/%s", runtimeDirectory, MachineDataFileName)

	var fileData []byte
	var machineData MachineData = MachineData{
		QemuPid:      0,
		State:        MachineStatusUnknown,
		SSHLocalPort: 0,
	}

	fileData, err := os.ReadFile(dataFile)
	if err != nil {
		log.Printf("error: could not open data file: %s\n", err.Error())
	} else {
		err = json.Unmarshal(fileData, &machineData)
		if err != nil {
			log.Printf("[machine] could not obtain machine data: %s", err.Error())
			return nil
		}
	}

	configFile := fmt.Sprintf("%s/%s", runtimeDirectory, MachineConfigFileName)

	return &Machine{
		Name:             machineName,
		Status:           machineData.State,
		QemuPid:          machineData.QemuPid,
		SSHLocalPort:     machineData.SSHLocalPort,
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
	var statusFile string = fmt.Sprintf("%s/%s", m.RuntimeDirectory, MachineDataFileName)
	var fileHandle *os.File
	var machineData MachineData

	log.Printf("[UpdateStatus] opening file '%s'\n", statusFile)
	fileHandle, err = os.OpenFile(statusFile, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0755)
	if err != nil {
		return err
	}

	/* populate new MachineData */
	machineData = MachineData{
		QemuPid:      m.QemuPid,
		SSHLocalPort: m.SSHLocalPort,
		State:        status,
	}

	switch status {
	case MachineStatusDegraded, MachineStatusStarted, MachineStatusStopped, MachineStatusUnknown:
		{
			log.Printf("[UpdateStatus] updating file '%s'.\n", statusFile)
			jsonBytes, err := json.Marshal(machineData)

			if err != nil {
				log.Printf("[UpdateStatus] error while generating new JSON: '%s'.\n", err.Error())
			} else {
				log.Printf("[UpdateStatus] writing [%s] to file '%s'.\n", string(jsonBytes), statusFile)
				_, err = fileHandle.Write(jsonBytes)

				if err != nil {
					log.Printf("[UpdateStatus] error while updating '%s': %s\n", statusFile, err.Error())
				}
			}
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
