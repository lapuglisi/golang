package qemuctl_helpers

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"

	"gopkg.in/yaml.v2"

	runtime "luizpuglisi.com/qemuctl/runtime"
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

func (cd *ConfigurationData) appendQemuArg(argsSlice []string, argKey string, argValue string) (newSlice []string) {
	return append(argsSlice, []string{argKey, argValue}...)
}

func (cd *ConfigurationData) GetQemuArgs(qemuPath string, machineName string) (qemuArgs []string, err error) {
	/* Config specific */
	var machineSpec string
	var netSpec string

	/* VNC Spec parser */
	var vncRegex regexp.Regexp = *regexp.MustCompile(`[0-9\.]+:\d+`)

	/* Pre-checks */
	if len(machineName) > 0 {
		qemuArgs = cd.appendQemuArg(qemuArgs, "-name", machineName)
	} else {
		qemuArgs = cd.appendQemuArg(qemuArgs, "-name", cd.Machine.MachineName)
	}

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

	/*
	 * Display specification
	 */
	if !cd.Display.EnableGraphics {
		qemuArgs = append(qemuArgs, "-nographic")
	} else {
		// -- VGA
		qemuArgs = cd.appendQemuArg(qemuArgs, "-vga", cd.Display.VGAType)

		// -- Display
		qemuArgs = cd.appendQemuArg(qemuArgs, "-display", cd.Display.DisplaySpec)
	}

	// VNC ?
	if len(cd.VNCConfig) > 0 {
		// Is it in the format "xxx.xxx.xxx.xxx:ddd" ?
		if vncRegex.Match([]byte(cd.VNCConfig)) {
			qemuArgs = cd.appendQemuArg(qemuArgs, "-vnc", cd.VNCConfig)
		} else {
			qemuArgs = cd.appendQemuArg(qemuArgs, "-vnc", fmt.Sprintf("127.0.0.1:%s", cd.VNCConfig))
		}
	}

	/**
	 * BIOS and Boot habling
	 */
	if len(cd.Boot.KernelPath) > 0 && len(cd.Boot.RamdiskPath) > 0 {
		// Do not use biosFile or boot related stuff. Boot directly to kernel
		qemuArgs = cd.appendQemuArg(qemuArgs, "-kernel", cd.Boot.KernelPath)
		qemuArgs = cd.appendQemuArg(qemuArgs, "-initrd", cd.Boot.RamdiskPath)
	} else {
		if len(cd.Boot.BiosFile) > 0 {
			qemuArgs = cd.appendQemuArg(qemuArgs, "-bios", cd.Boot.BiosFile)
		}

		// -- Boot menu & Boot order (exclusive)
		if cd.Boot.EnableBootMenu {
			qemuArgs = cd.appendQemuArg(qemuArgs, "-boot", "menu=on")
		} else if len(cd.Boot.BootOrder) > 0 {
			qemuArgs = cd.appendQemuArg(qemuArgs, "-boot", "order="+cd.Boot.BootOrder)
		}
	}

	// -- Background?
	if cd.RunAsDaemon {
		qemuArgs = append(qemuArgs, "-daemonize")
	}

	// -- Network spec
	{
		/* Configure user network device */
		netSpec = fmt.Sprintf("%s,netdev=%s", cd.Net.DeviceType, cd.Net.User.ID)
		qemuArgs = cd.appendQemuArg(qemuArgs, "-device", netSpec)

		/* Configure User NIC */
		netSpec = fmt.Sprintf("user,id=%s", cd.Net.User.ID)

		if len(cd.Net.User.IPSubnet) > 0 {
			netSpec = fmt.Sprintf("%s,net=%s", netSpec, cd.Net.User.IPSubnet)
		}

		if cd.SSH.LocalPort > 0 {
			netSpec = fmt.Sprintf("%s,hostfwd=tcp::%d-:22", netSpec, cd.SSH.LocalPort)
		}

		/* Port fowards come here */
		for _, _value := range cd.Net.User.PortForwards {
			netSpec = fmt.Sprintf("%s,hostfwd=tcp::%d-:%d", netSpec, _value.HostPort, _value.GuestPort)
		}

		qemuArgs = cd.appendQemuArg(qemuArgs, "-netdev", netSpec)

		/*
		 * Configure bridge, if any
		 */
		if len(cd.Net.Bridge.Interface) > 0 {
			//-- Device specification
			netSpec = fmt.Sprintf("%s,netdev=%s", cd.Net.DeviceType, cd.Net.Bridge.ID)
			if len(cd.Net.Bridge.MacAddress) > 0 {
				netSpec = fmt.Sprintf("%s,mac=", cd.Net.Bridge.MacAddress)
			}
			qemuArgs = cd.appendQemuArg(qemuArgs, "-device", netSpec)

			// Bridge definition
			netSpec = fmt.Sprintf("bridge,id=%s,br=%s", cd.Net.Bridge.ID, cd.Net.Bridge.Interface)
			if len(cd.Net.Bridge.Helper) > 0 {
				netSpec = fmt.Sprintf("%s,helper=%s", netSpec, cd.Net.Bridge.Helper)
			}
			qemuArgs = cd.appendQemuArg(qemuArgs, "-netdev", netSpec)
		}
	}

	/*
	 * Disk specification
	 */
	if len(cd.Disks.BlockDevice) > 0 { // TODO: Use stat to check whether it is a valid block device
		driveName := "xvda"
		// Appends drive/device specification
		qemuArgs = cd.appendQemuArg(qemuArgs, "-device", fmt.Sprintf("virtio-blk-pci,drive=%s", driveName))

		// Appends block device configuration
		qemuArgs = cd.appendQemuArg(qemuArgs,
			"-blockdev",
			fmt.Sprintf("node-name=%s,driver=raw,file.driver=host_device,file.filename=%s", driveName, cd.Disks.BlockDevice))
	} else {
		// -- Otherwise, we finally add hard disk info
		qemuArgs = append(qemuArgs, cd.Disks.HardDisk)
	}

	/* Add a monitor specfication to be able to operate on the machine */
	monitorSpec := fmt.Sprintf("unix:%s/%s/qemu-monitor.sock,server,nowait",
		runtime.GetUserDataDir(), cd.Machine.MachineName)
	qemuArgs = cd.appendQemuArg(qemuArgs, "-monitor", monitorSpec)

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
