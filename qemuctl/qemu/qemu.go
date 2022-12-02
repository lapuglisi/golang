package qemuctl_qemu

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"

	config "luizpuglisi.com/qemuctl/helpers"
	runtime "luizpuglisi.com/qemuctl/runtime"
)

type QemuCommand struct {
	QemuPath      string
	Configuration *config.ConfigurationData
	Monitor       *QemuMonitor
}

func NewQemuCommand(qemuBinary string, configData *config.ConfigurationData, qemuMonitor *QemuMonitor) (qemu *QemuCommand) {
	var qemuPath string

	qemuPath, err := exec.LookPath(qemuBinary)
	if err != nil {
		qemuPath = qemuBinary
	}

	return &QemuCommand{
		QemuPath:      qemuPath,
		Configuration: configData,
		Monitor:       qemuMonitor,
	}
}

func (qemu *QemuCommand) appendQemuArg(argSlice []string, argKey string, argValue string) []string {
	return append(argSlice, []string{argKey, argValue}...)
}

func (qemu *QemuCommand) getQemuArgs() (qemuArgs []string, err error) {
	/* Config specific */
	var machineSpec string
	var netSpec string

	var cd *config.ConfigurationData = qemu.Configuration

	var machine *runtime.Machine = runtime.NewMachine(cd.Machine.MachineName)
	var monitor *QemuMonitor = NewQemuMonitor(machine)

	/* VNC Spec parser */
	var vncRegex regexp.Regexp = *regexp.MustCompile(`[0-9\.]+:\d+`)

	qemuArgs = append(qemuArgs, qemu.QemuPath)

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

		qemuArgs = qemu.appendQemuArg(qemuArgs, "-machine", machineSpec)
	}

	// -- Machine Name
	if len(cd.Machine.MachineName) > 0 {
		qemuArgs = qemu.appendQemuArg(qemuArgs, "-name", cd.Machine.MachineName)
	}

	// -- Memory
	qemuArgs = qemu.appendQemuArg(qemuArgs, "-m", cd.Memory)

	// -- cpus
	qemuArgs = qemu.appendQemuArg(qemuArgs, "-smp", fmt.Sprintf("%d", cd.CPUs))

	// -- CDROM
	if len(cd.Disks.ISOCDrom) > 0 {
		qemuArgs = qemu.appendQemuArg(qemuArgs, "-cdrom", cd.Disks.ISOCDrom)
	}

	/*
	 * Display specification
	 */
	if !cd.Display.EnableGraphics {
		qemuArgs = append(qemuArgs, "-nographic")
	} else {
		// -- VGA
		qemuArgs = qemu.appendQemuArg(qemuArgs, "-vga", cd.Display.VGAType)

		// -- Display
		qemuArgs = qemu.appendQemuArg(qemuArgs, "-display", cd.Display.DisplaySpec)
	}

	// VNC ?
	if len(cd.VNCConfig) > 0 {
		// Is it in the format "xxx.xxx.xxx.xxx:ddd" ?
		if vncRegex.Match([]byte(cd.VNCConfig)) {
			qemuArgs = qemu.appendQemuArg(qemuArgs, "-vnc", cd.VNCConfig)
		} else {
			qemuArgs = qemu.appendQemuArg(qemuArgs, "-vnc", fmt.Sprintf("127.0.0.1:%s", cd.VNCConfig))
		}
	}

	/**
	 * BIOS and Boot habling
	 */
	if len(cd.Boot.KernelPath) > 0 && len(cd.Boot.RamdiskPath) > 0 {
		// Do not use biosFile or boot related stuff. Boot directly to kernel
		qemuArgs = qemu.appendQemuArg(qemuArgs, "-kernel", cd.Boot.KernelPath)
		qemuArgs = qemu.appendQemuArg(qemuArgs, "-initrd", cd.Boot.RamdiskPath)
	} else {
		if len(cd.Boot.BiosFile) > 0 {
			qemuArgs = qemu.appendQemuArg(qemuArgs, "-bios", cd.Boot.BiosFile)
		}

		// -- Boot menu & Boot order (exclusive)
		if cd.Boot.EnableBootMenu {
			qemuArgs = qemu.appendQemuArg(qemuArgs, "-boot", "menu=on")
		} else if len(cd.Boot.BootOrder) > 0 {
			qemuArgs = qemu.appendQemuArg(qemuArgs, "-boot", "order="+cd.Boot.BootOrder)
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
		qemuArgs = qemu.appendQemuArg(qemuArgs, "-device", netSpec)

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

		qemuArgs = qemu.appendQemuArg(qemuArgs, "-netdev", netSpec)

		/*
		 * Configure bridge, if any
		 */
		if len(cd.Net.Bridge.Interface) > 0 {
			//-- Device specification
			netSpec = fmt.Sprintf("%s,netdev=%s", cd.Net.DeviceType, cd.Net.Bridge.ID)
			if len(cd.Net.Bridge.MacAddress) > 0 {
				netSpec = fmt.Sprintf("%s,mac=", cd.Net.Bridge.MacAddress)
			}
			qemuArgs = qemu.appendQemuArg(qemuArgs, "-device", netSpec)

			// Bridge definition
			netSpec = fmt.Sprintf("bridge,id=%s,br=%s", cd.Net.Bridge.ID, cd.Net.Bridge.Interface)
			if len(cd.Net.Bridge.Helper) > 0 {
				netSpec = fmt.Sprintf("%s,helper=%s", netSpec, cd.Net.Bridge.Helper)
			}
			qemuArgs = qemu.appendQemuArg(qemuArgs, "-netdev", netSpec)
		}
	}

	/*
	 * Disk specification
	 */
	if len(cd.Disks.BlockDevice) > 0 { // TODO: Use stat to check whether it is a valid block device
		driveName := "xvda"
		// Appends drive/device specification
		qemuArgs = qemu.appendQemuArg(qemuArgs, "-device", fmt.Sprintf("virtio-blk-pci,drive=%s", driveName))

		// Appends block device configuration
		qemuArgs = qemu.appendQemuArg(qemuArgs,
			"-blockdev",
			fmt.Sprintf("node-name=%s,driver=raw,file.driver=host_device,file.filename=%s", driveName, cd.Disks.BlockDevice))
	} else {
		// -- Otherwise, we finally add hard disk info
		qemuArgs = append(qemuArgs, cd.Disks.HardDisk)
	}

	/* Add a monitor specfication to be able to operate on the machine */
	qemuArgs = qemu.appendQemuArg(qemuArgs, "-chardev", monitor.GetChardevSpec())
	qemuArgs = qemu.appendQemuArg(qemuArgs, "-qmp", monitor.GetMonitorSpec())

	return qemuArgs, nil
}

func (qemu *QemuCommand) Launch() (err error) {
	var procAttrs *os.ProcAttr

	var qemuArgs []string

	qemuArgs, err = qemu.getQemuArgs()
	if err != nil {
		return err
	}

	// TODO: use the log feature
	log.Println("[QemuCommand::Launch] Executing QEMU with:")
	log.Printf("qemu_path ....... %s\n", qemu.QemuPath)
	log.Printf("qemu_args ....... %s\n", strings.Join(qemuArgs, " "))

	/* Actual execution of QEMU */
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

	procHandle, err := os.StartProcess(qemu.QemuPath, qemuArgs, procAttrs)
	if err == nil {
		if qemu.Configuration.RunAsDaemon {
			err = procHandle.Release()
		}
	}

	return err
}
