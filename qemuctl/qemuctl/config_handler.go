package qemuctl

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// ConfigurationData holds the power of the serominers
type ConfigurationData struct {
	EnableKVM      bool
	MachineName    string
	MachineType    string
	VNCConfig      string
	RunAsDaemon    bool
	Memory         string
	CPUs           int64
	NetID          string
	NetIPSubnet    string
	SSHLocalPort   int64
	HardDiskFile   string
	ISOCDrom       string
	VGAType        string
	DisplaySpec    string
	BiosFile       string
	EnableBootMenu bool
	BootOrder      string
	AccelType      string
}

// ConfigurationHandler is one hell of a seroclockers
type ConfigurationHandler struct {
	filePath string
}

func init() {
}

/* ConfigurationData implementation */
func NewConfigData() (configData *ConfigurationData) {
	return &ConfigurationData{
		EnableKVM:    true,
		MachineType:  "q35",
		AccelType:    "hvf",
		VGAType:      "virtio",
		DisplaySpec:  "default",
		SSHLocalPort: -1,
		CPUs:         1,
	}
}

/* ConfigurationHandler implementation */
func NewConfigHandler(configFile string) (configHandler *ConfigurationHandler) {
	return &ConfigurationHandler{
		filePath: configFile,
	}
}

func (ch *ConfigurationHandler) parseConfigEntry(configEntry string, configData *ConfigurationData) (err error) {
	var equalPos int = -1
	var configKey string
	var configValue string

	var commentPattern regexp.Regexp = *regexp.MustCompile(`^\s*#`)
	var entryBytes []byte = []byte(configEntry)

	if commentPattern.Match(entryBytes) {
		fmt.Printf("[info] parseConfigEntry: skipping comment '%s'\n", configEntry)
		return nil
	}

	equalPos = strings.IndexByte(configEntry, '=')
	if equalPos == -1 {
		err = fmt.Errorf("invalid config entry: %s", configEntry)
		return err
	}

	configKey = strings.TrimSpace(configEntry[0 : equalPos-1])
	configValue = strings.TrimSpace(configEntry[equalPos+1:])

	/*
		fmt.Printf("[parseConfigEntry] configEntry is .... '%s'\n", configEntry)
		fmt.Printf("[parseConfigEntry] configKey is ...... '%s'\n", configKey)
		fmt.Printf("[parseConfigEntry] configValue is .... '%s'\n", configValue)
	*/

	err = nil
	switch configKey {
	case "name":
		{
			configData.MachineName = configValue
			break
		}
	case "memory":
		{
			configData.Memory = configValue
			break
		}
	case "cpus":
		{
			configData.CPUs, err = strconv.ParseInt(configValue, 10, 0)
			break

		}
	case "disk":
		{
			configData.HardDiskFile = configValue
			break
		}
	case "cdrom":
		{
			configData.ISOCDrom = configValue
			break
		}
	case "vga":
		{
			configData.VGAType = configValue
			break
		}
	case "vnclisten":
		{
			configData.VNCConfig = configValue
			break
		}
	case "display":
		{
			configData.DisplaySpec = configValue
			break
		}
	case "bios":
		{
			configData.BiosFile = configValue
			break
		}
	case "boot_menu":
		{
			configData.EnableBootMenu, err = strconv.ParseBool(configValue)
			break
		}
	case "boot_order":
		{
			configData.BootOrder = configValue
			break
		}
	case "background":
		{
			configData.RunAsDaemon, err = strconv.ParseBool(configValue)
			break
		}
	case "machine_type":
		{
			configData.MachineType = configValue
			break
		}
	case "accel":
		{
			configData.AccelType = configValue
			break
		}
	case "ip_subnet":
		{
			configData.NetIPSubnet = configValue
			break
		}
	case "net_id":
		{
			configData.NetID = configValue
			break
		}
	case "ssh_local_port":
		{
			configData.SSHLocalPort, err = strconv.ParseInt(configValue, 10, 0)
			break
		}

	default:
		{
			err = fmt.Errorf("unknown config entry '%s'", configEntry)
		}
	}

	return err
}

func (ch *ConfigurationHandler) ParseConfigFile() (configData *ConfigurationData, err error) {
	var osErr error

	// Open file
	fileHandle, osErr := os.OpenFile(ch.filePath, os.O_RDONLY, 0644)
	if osErr != nil {
		err = fmt.Errorf("could not open file '%s': %s", ch.filePath, osErr.Error())
		return nil, err
	}
	defer fileHandle.Close()

	// Read lines
	bufScanner := bufio.NewReader(fileHandle)

	configData = NewConfigData()
	osErr = nil
	for osErr == nil {
		lineBytes, _ /* isPrefix */, osErr := bufScanner.ReadLine()
		if osErr != nil {
			break
		}

		// Parse entry here
		configEntry := string(lineBytes)
		if err = ch.parseConfigEntry(configEntry, configData); err != nil {
			fmt.Printf("[warning] parseConfigEntry: %s\n", err.Error())
		}
	}

	return configData, nil
}
