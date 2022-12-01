package qemuctl

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func init() {
}

const (
	QemuDefaultSystemBin string = "qemu-system-x86_64"
)

func appendQemuArg(qemuArgs []string, argKey string, argValue string) (newArgs []string) {
	return append(qemuArgs, []string{argKey, argValue}...)
}

func launchQemu(qemuBinary string, configData *ConfigurationData) (err error) {

	var qemuArgs []string = nil
	var qemuPath string
	var procAttrs *os.ProcAttr = nil

	/* Config specific */
	var machineSpec string
	var netSpec string

	/* VNC Spec parser */
	var vncRegex regexp.Regexp = *regexp.MustCompile(`[0-9\.]+:\d+`)

	// Get QEMU path
	qemuPath, err = exec.LookPath(qemuBinary)
	if err != nil {
		return err
	}

	/* Initialize qemuArgs */
	qemuArgs = append(qemuArgs, qemuPath)

	/* Do the config stuff */
	if configData.EnableKVM {
		qemuArgs = append(qemuArgs, "-enable-kvm")
	}

	// -- Machine spec (type and accel)
	{
		machineSpec = fmt.Sprintf("type=%s", configData.MachineType)
		if len(configData.AccelType) > 0 {
			machineSpec = fmt.Sprintf("%s,accel=%s", machineSpec, configData.AccelType)

		}

		qemuArgs = appendQemuArg(qemuArgs, "-machine", machineSpec)
	}

	// -- Machine Name
	if len(configData.MachineName) > 0 {
		qemuArgs = appendQemuArg(qemuArgs, "-name", configData.MachineName)
	}

	// -- Memory
	qemuArgs = appendQemuArg(qemuArgs, "-m", configData.Memory)

	// -- cpus
	qemuArgs = appendQemuArg(qemuArgs, "-smp", fmt.Sprintf("%d", configData.CPUs))

	// -- CDROM
	if len(configData.ISOCDrom) > 0 {
		qemuArgs = appendQemuArg(qemuArgs, "-cdrom", configData.ISOCDrom)
	}

	// -- VGA
	qemuArgs = appendQemuArg(qemuArgs, "-vga", configData.VGAType)

	// -- Display
	qemuArgs = appendQemuArg(qemuArgs, "-display", configData.DisplaySpec)

	// VNC ?
	if len(configData.VNCConfig) > 0 {
		// Is it in the format "xxx.xxx.xxx.xxx:ddd" ?
		if vncRegex.Match([]byte(configData.VNCConfig)) {
			qemuArgs = appendQemuArg(qemuArgs, "-vnc", configData.VNCConfig)
		} else {
			qemuArgs = appendQemuArg(qemuArgs, "-vnc", fmt.Sprintf("127.0.0.1:%s", configData.VNCConfig))
		}
	}

	// -- Bios file
	if len(configData.BiosFile) > 0 {
		qemuArgs = appendQemuArg(qemuArgs, "-bios", configData.BiosFile)
	}

	// -- Boot menu & Boot order (exclusive)
	if configData.EnableBootMenu {
		qemuArgs = appendQemuArg(qemuArgs, "-boot", "menu=on")
	} else if len(configData.BootOrder) > 0 {
		qemuArgs = appendQemuArg(qemuArgs, "-boot", "order="+configData.BootOrder)
	}

	// -- Background?
	if configData.RunAsDaemon {
		qemuArgs = append(qemuArgs, "-daemonize")
	}

	// -- Network spec
	{
		netSpec = "user,model=virtio-net-pci"
		if len(configData.NetIPSubnet) > 0 {
			netSpec = fmt.Sprintf("%s,net=%s", netSpec, configData.NetIPSubnet)
		}

		if len(configData.NetID) > 0 {
			netSpec = fmt.Sprintf("%s,id=%s", netSpec, configData.NetID)
		}

		if configData.SSHLocalPort > 0 {
			netSpec = fmt.Sprintf("%s,hostfwd=tcp::%d-:22", netSpec, configData.SSHLocalPort)
		}

		qemuArgs = appendQemuArg(qemuArgs, "-nic", netSpec)
	}

	// -- Finally, add hard disk info
	qemuArgs = append(qemuArgs, configData.HardDiskFile)

	fmt.Println("[INFO] Executing QEMU with:")
	fmt.Printf("qemu_path .......... %s\n", qemuPath)
	fmt.Printf("qemu_args .......... %s\n", strings.Join(qemuArgs, " "))

	err = nil
	procAttrs = &os.ProcAttr{
		Dir: os.ExpandEnv("$HOME"),
		Env: os.Environ(),
		Files: []*os.File{
			os.Stdin,
			os.Stdout,
			os.Stderr,
		},
		Sys: nil,
	}
	procHandle, err := os.StartProcess(qemuPath, qemuArgs, procAttrs)
	if err == nil {
		if configData.RunAsDaemon {
			err = procHandle.Release()
		}
	}

	return err
}

func HandleStart(startArgs []string) (err error) {
	var configFile string
	var configData *ConfigurationData = nil
	var qemuBinary string = QemuDefaultSystemBin

	err = nil
	if len(startArgs) < 1 {
		err = fmt.Errorf("[handle_start] Too few arguments")
		return err
	}

	configFile = startArgs[0]
	configHandle := NewConfigHandler(configFile)

	// set qemu binary, if specified
	if len(startArgs) > 1 {
		qemuBinary = startArgs[2]
	}

	configData, err = configHandle.ParseConfigFile()
	if err != nil {
		return err
	}

	// Now that we have configData, launch qemu
	err = launchQemu(qemuBinary, configData)
	if err != nil {
		return err
	}

	return err
}
