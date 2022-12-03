package qemuctl_helpers

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v2"
)

// ConfigurationData holds the power of the serominers
type portForwards struct {
	GuestPort int `yaml:"guestPort"`
	HostPort  int `yaml:"hostPort"`
}

type ConfigurationData struct {
	Machine struct {
		EnableKVM   bool   `yaml:"enableKVM"`
		MachineName string `yaml:"name"`
		MachineType string `yaml:"type"`
		AccelType   string `yaml:"accel"`
	} `yaml:"machine"`
	VNCConfig   string `yaml:"vncListen"`
	RunAsDaemon bool   `yaml:"runAsDaemon"`
	Memory      string `yaml:"memory"`
	CPUs        int64  `yaml:"cpus"`
	Net         struct {
		DeviceType string `yaml:"deviceType"`
		User       struct {
			ID           string         `yaml:"id"`
			IPSubnet     string         `yaml:"ipSubnet"`
			PortForwards []portForwards `yaml:"portForwards"`
		} `yaml:"user"`
		Bridge struct {
			ID         string `yaml:"id"`
			Interface  string `yaml:"interface"`
			MacAddress string `yaml:"mac"`
			Helper     string `yaml:"helper"`
		}
	} `yaml:"net"`
	SSH struct {
		LocalPort int `yaml:"localPort"`
	} `yaml:"ssh"`
	Disks struct {
		BlockDevice string `yaml:"blockDevice"`
		HardDisk    string `yaml:"hardDisk"`
		ISOCDrom    string `yaml:"cdrom"`
	} `yaml:"disks"`
	Display struct {
		EnableGraphics bool   `yaml:"enableGraphics"`
		VGAType        string `yaml:"vgaType"`
		DisplaySpec    string `yaml:"displaySpec"`
	} `yaml:"display"`
	Boot struct {
		KernelPath     string `yaml:"kernelPath"`
		RamdiskPath    string `yaml:"ramdiskPath"`
		BiosFile       string `yaml:"biosFile"`
		EnableBootMenu bool   `yaml:"enableBootMenu"`
		BootOrder      string `yaml:"bootOrder"`
	} `yaml:"boot"`
	QemuBinary string `yaml:"qemuBinary"`
}

// ConfigurationHandler is one hell of a seroclockers
type ConfigurationHandler struct {
	filePath string
}

func init() {
}

/* ConfigurationData implementation */
func NewConfigData() (configData *ConfigurationData) {
	configData = &ConfigurationData{}

	configData.Machine.MachineType = "q35"
	configData.Machine.AccelType = "hvm"
	configData.Machine.EnableKVM = true

	configData.Net.DeviceType = "e1000"

	configData.Net.User.ID = "mynet0"

	configData.Net.Bridge.ID = "mybr0"

	configData.RunAsDaemon = false

	configData.Display.EnableGraphics = true
	configData.Display.VGAType = "none"
	configData.Display.DisplaySpec = "none"

	return configData
}

/* ConfigurationHandler implementation */
func NewConfigHandler(configFile string) (configHandler *ConfigurationHandler) {
	return &ConfigurationHandler{
		filePath: configFile,
	}
}

func (ch *ConfigurationHandler) ParseConfigFile() (configData *ConfigurationData, err error) {
	var configBytes []byte = nil
	var bufReader *bufio.Reader = nil

	// Open file
	fileHandle, osErr := os.OpenFile(ch.filePath, os.O_RDONLY, 0644)
	if osErr != nil {
		err = fmt.Errorf("could not open file '%s': %s", ch.filePath, osErr.Error())
		return nil, err
	}
	defer fileHandle.Close()

	// Read lines
	bufReader = bufio.NewReader(fileHandle)

	configData = NewConfigData()
	osErr = nil

	configBytes, err = io.ReadAll(bufReader)
	if err != nil {
		return nil, err
	}

	/* Now YAML the whole thing */
	err = yaml.Unmarshal(configBytes, &configData)
	if err != nil {
		return nil, err
	}

	return configData, nil
}
