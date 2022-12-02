package qemuctl_qemu

import (
	"fmt"
	"net"

	runtime "luizpuglisi.com/qemuctl/runtime"
)

func init() {

}

const (
	QemuMonitorSocketFileName  string = "qemu-monitor.sock"
	QemuMonitorShutdownCommand string = "system_powerdown"
	QemuMonitorQuitCommand     string = "quit"
)

func ShutdownMachine(machineName string) error {
	var machine *runtime.Machine = runtime.NewMachine(machineName)
	var monitorSocket string = fmt.Sprintf("%s/%s",
		machine.RuntimeDirectory, QemuMonitorSocketFileName)

	if !machine.Exists() {
		return fmt.Errorf("machine '%s' does not exist", machineName)
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
