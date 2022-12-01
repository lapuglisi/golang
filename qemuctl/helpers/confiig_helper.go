package qemuctl_helpers

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"gopkg.in/yaml.v2"
)

// ConfigurationData holds the power of the serominers
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
		NetID    string `yaml:"netId"`
		IPSubnet string `yaml:"ipSubnet"`
	} `yaml:"net"`
	SSH struct {
		LocalPort int `yaml:"localPort"`
	} `yaml:"ssh"`
	Disks struct {
		HardDisk string `yaml:"hardDisk"`
		ISOCDrom string `yaml:"cdrom"`
	} `yaml:"disks"`
	Display struct {
		VGAType     string `yaml:"vgaType"`
		DisplaySpec string `yaml:"displaySpec"`
	} `yaml:"display"`
	Boot struct {
		BiosFile       string `yaml:"biosFile"`
		EnableBootMenu bool   `yaml:"enableBootMenu"`
		BootOrder      string `yaml:"bootOrder"`
	} `yaml:"boot"`
}

// ConfigurationHandler is one hell of a seroclockers
type ConfigurationHandler struct {
	filePath string
}

func init() {
}

/* ConfigurationData implementation */
func NewConfigData() (configData *ConfigurationData) {
	return &ConfigurationData{}
}

func (cd *ConfigurationData) appendQemuArg(argsSlice []string, argKey string, argValue string) (newSlice []string) {
	return append(argsSlice, []string{argKey, argValue}...)
}

func (cd *ConfigurationData) GetQemuArgs(qemuPath string) (qemuArgs []string, err error) {
	/* Config specific */
	var machineSpec string
	var netSpec string

	/* VNC Spec parser */
	var vncRegex regexp.Regexp = *regexp.MustCompile(`[0-9\.]+:\d+`)

	/* Initialize qemuArgs */
	qemuArgs = append(qemuArgs, qemuPath)

	/* Do the config stuff */
	if cd.Machine.EnableKVM {
		qemuArgs = append(qemuArgs, "-enable-kvm")
	}

	// -- Machine spec (type and accel)
	{
		machineSpec = fmt.Sprintf("type=%s", cd.Machine.MachineType)
		if len(cd.Machine.AccelType) > 0 {
			machineSpec = fmt.Sprintf("%s,accel=%s", machineSpec, cd.Machine.AccelType)
		}

		qemuArgs = cd.appendQemuArg(qemuArgs, "-machine", machineSpec)
	}

	// -- Machine Name
	if len(cd.Machine.MachineName) > 0 {
		qemuArgs = cd.appendQemuArg(qemuArgs, "-name", cd.Machine.MachineName)
	}

	// -- Memory
	qemuArgs = cd.appendQemuArg(qemuArgs, "-m", cd.Memory)

	// -- cpus
	qemuArgs = cd.appendQemuArg(qemuArgs, "-smp", fmt.Sprintf("%d", cd.CPUs))

	// -- CDROM
	if len(cd.Disks.ISOCDrom) > 0 {
		qemuArgs = cd.appendQemuArg(qemuArgs, "-cdrom", cd.Disks.ISOCDrom)
	}

	// -- VGA
	qemuArgs = cd.appendQemuArg(qemuArgs, "-vga", cd.Display.VGAType)

	// -- Display
	qemuArgs = cd.appendQemuArg(qemuArgs, "-display", cd.Display.DisplaySpec)

	// VNC ?
	if len(cd.VNCConfig) > 0 {
		// Is it in the format "xxx.xxx.xxx.xxx:ddd" ?
		if vncRegex.Match([]byte(cd.VNCConfig)) {
			qemuArgs = cd.appendQemuArg(qemuArgs, "-vnc", cd.VNCConfig)
		} else {
			qemuArgs = cd.appendQemuArg(qemuArgs, "-vnc", fmt.Sprintf("127.0.0.1:%s", cd.VNCConfig))
		}
	}

	// -- Bios file
	if len(cd.Boot.BiosFile) > 0 {
		qemuArgs = cd.appendQemuArg(qemuArgs, "-bios", cd.Boot.BiosFile)
	}

	// -- Boot menu & Boot order (exclusive)
	if cd.Boot.EnableBootMenu {
		qemuArgs = cd.appendQemuArg(qemuArgs, "-boot", "menu=on")
	} else if len(cd.Boot.BootOrder) > 0 {
		qemuArgs = cd.appendQemuArg(qemuArgs, "-boot", "order="+cd.Boot.BootOrder)
	}

	// -- Background?
	if cd.RunAsDaemon {
		qemuArgs = append(qemuArgs, "-daemonize")
	}

	// -- Network spec
	{
		netSpec = "user,model=virtio-net-pci"
		if len(cd.Net.IPSubnet) > 0 {
			netSpec = fmt.Sprintf("%s,net=%s", netSpec, cd.Net.IPSubnet)
		}

		if len(cd.Net.NetID) > 0 {
			netSpec = fmt.Sprintf("%s,id=%s", netSpec, cd.Net.NetID)
		}

		if cd.SSH.LocalPort > 0 {
			netSpec = fmt.Sprintf("%s,hostfwd=tcp::%d-:22", netSpec, cd.SSH.LocalPort)
		}

		qemuArgs = cd.appendQemuArg(qemuArgs, "-nic", netSpec)
	}

	// -- Finally, add hard disk info
	qemuArgs = append(qemuArgs, cd.Disks.HardDisk)

	return qemuArgs, nil
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

	configBytes, err = ioutil.ReadAll(bufReader)
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
