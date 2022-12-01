package qemuctl_helpers

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
	if cd.EnableKVM {
		qemuArgs = append(qemuArgs, "-enable-kvm")
	}

	// -- Machine spec (type and accel)
	{
		machineSpec = fmt.Sprintf("type=%s", cd.MachineType)
		if len(cd.AccelType) > 0 {
			machineSpec = fmt.Sprintf("%s,accel=%s", machineSpec, cd.AccelType)
		}

		qemuArgs = cd.appendQemuArg(qemuArgs, "-machine", machineSpec)
	}

	// -- Machine Name
	if len(cd.MachineName) > 0 {
		qemuArgs = cd.appendQemuArg(qemuArgs, "-name", cd.MachineName)
	}

	// -- Memory
	qemuArgs = cd.appendQemuArg(qemuArgs, "-m", cd.Memory)

	// -- cpus
	qemuArgs = cd.appendQemuArg(qemuArgs, "-smp", fmt.Sprintf("%d", cd.CPUs))

	// -- CDROM
	if len(cd.ISOCDrom) > 0 {
		qemuArgs = cd.appendQemuArg(qemuArgs, "-cdrom", cd.ISOCDrom)
	}

	// -- VGA
	qemuArgs = cd.appendQemuArg(qemuArgs, "-vga", cd.VGAType)

	// -- Display
	qemuArgs = cd.appendQemuArg(qemuArgs, "-display", cd.DisplaySpec)

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
	if len(cd.BiosFile) > 0 {
		qemuArgs = cd.appendQemuArg(qemuArgs, "-bios", cd.BiosFile)
	}

	// -- Boot menu & Boot order (exclusive)
	if cd.EnableBootMenu {
		qemuArgs = cd.appendQemuArg(qemuArgs, "-boot", "menu=on")
	} else if len(cd.BootOrder) > 0 {
		qemuArgs = cd.appendQemuArg(qemuArgs, "-boot", "order="+cd.BootOrder)
	}

	// -- Background?
	if cd.RunAsDaemon {
		qemuArgs = append(qemuArgs, "-daemonize")
	}

	// -- Network spec
	{
		netSpec = "user,model=virtio-net-pci"
		if len(cd.NetIPSubnet) > 0 {
			netSpec = fmt.Sprintf("%s,net=%s", netSpec, cd.NetIPSubnet)
		}

		if len(cd.NetID) > 0 {
			netSpec = fmt.Sprintf("%s,id=%s", netSpec, cd.NetID)
		}

		if cd.SSHLocalPort > 0 {
			netSpec = fmt.Sprintf("%s,hostfwd=tcp::%d-:22", netSpec, cd.SSHLocalPort)
		}

		qemuArgs = cd.appendQemuArg(qemuArgs, "-nic", netSpec)
	}

	// -- Finally, add hard disk info
	qemuArgs = append(qemuArgs, cd.HardDiskFile)

	return qemuArgs, nil
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
