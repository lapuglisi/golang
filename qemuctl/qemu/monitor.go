package qemuctl_monitor

import (
	"fmt"
	"net"

	runtime "luizpuglisi.com/qemuctl/runtime"
)

func init() {

}

const (
	QemuMonitorShutdownCommand string = "system_powerdown"
)

func ShutdownMachine(machineName string) error {
	var monitorSocket string = fmt.Sprintf("%s/qemu-monitor.sock",
		runtime.GetMachineDirectory(machineName))

	if !runtime.MachineExists(machineName) {
		return fmt.Errorf("machine %s dos not exist", machineName)
	}

	cn, err := net.Dial("unix", monitorSocket)
	if err != nil {
		return err
	}

	_, err = cn.Write([]byte(QemuMonitorShutdownCommand))
	if err != nil {
		return err
	}

	return nil
}
