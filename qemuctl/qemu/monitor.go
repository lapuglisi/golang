package qemuctl_qemu

import (
	"encoding/json"
	"fmt"
	"log"
	"net"

	runtime "luizpuglisi.com/qemuctl/runtime"
)

func init() {

}

const (
	QemuMonitorSocketFileName string = "qemu-monitor.sock"
	QemuMonitorDefaultID      string = "qemu-mon-qmp"
)

type QemuMonitor struct {
	Machine *runtime.Machine
}

func NewQemuMonitor(machine *runtime.Machine) *QemuMonitor {
	return &QemuMonitor{
		Machine: machine,
	}
}

func (monitor *QemuMonitor) ReadQmpHeader(unix net.Conn) (qmpHeader *QmpHeader, err error) {
	const BufferSize = 1024

	var dataBytes []byte = make([]byte, 0)
	var buffer []byte = make([]byte, BufferSize)
	var nBytes int = 0

	qmpHeader = &QmpHeader{}

	dataBytes = make([]byte, 0)
	for nBytes, err = unix.Read(buffer); err == nil && nBytes > 0; {
		dataBytes = append(dataBytes, buffer[:nBytes]...)
		if nBytes < BufferSize {
			break
		}
	}

	if err != nil {
		return nil, err
	}

	/* Now that we have dataBytes, Unmarshal it to QmpHeader */
	err = json.Unmarshal(dataBytes, qmpHeader)
	if err != nil {
		return nil, fmt.Errorf("json error: %s", err.Error())
	}

	return qmpHeader, nil
}

func (monitor *QemuMonitor) GetUnixSocketPath() string {
	return fmt.Sprintf("%s/%s", monitor.Machine.RuntimeDirectory, QemuMonitorSocketFileName)
}

func (monitor *QemuMonitor) GetChardevSpec() string {
	return fmt.Sprintf("socket,id=%s,path=%s,server=on,wait=off",
		QemuMonitorDefaultID, monitor.GetUnixSocketPath())
}

func (monitor *QemuMonitor) GetMonitorSpec() string {
	return fmt.Sprintf("chardev:%s", QemuMonitorDefaultID)
}

func (monitor *QemuMonitor) GetControlSocket() (unix net.Conn, err error) {
	var qmpCommand QmpBasicCommand
	var socketData []byte

	log.Printf("[InitializeSocket] opening socket '%s'\n", monitor.GetUnixSocketPath())
	{
		unix, err = net.Dial("unix", monitor.GetUnixSocketPath())
		if err != nil {
			return nil, err
		}
	}

	log.Printf("[InitializeSocket] Reading QMP header")
	{
		_, err = monitor.ReadQmpHeader(unix)
		if err != nil {
			return nil, err
		}
	}

	log.Printf("[initialize] enabling QMP capabilities")
	{
		qmpCommand.Execute = QmpCapabilitiesCommand
		socketData, err = qmpCommand.GetJsonBytes()
		if err != nil {
			return nil, err
		}

		nBytes, err := unix.Write(socketData)
		if err != nil || nBytes == 0 {
			return nil, err
		}

		result := QmpBasicResult{}
		err = monitor.ReadQmpResult(unix, &result)
		if err != nil {
			return nil, err
		}
	}

	log.Printf("[InitializeSocket] socket initialized")
	return unix, nil
}

func (monitor *QemuMonitor) ReadQmpResult(unix net.Conn, out interface{}) (err error) {
	const BufferSize = 1024

	var dataBytes []byte = make([]byte, 0)
	var buffer []byte = make([]byte, BufferSize)
	var nBytes int = 0

	dataBytes = make([]byte, 0)
	for nBytes, err = unix.Read(buffer); err == nil && nBytes > 0; {
		dataBytes = append(dataBytes, buffer[:nBytes]...)
		if nBytes < BufferSize {
			break
		}
	}
	if err != nil {
		return err
	}

	/* Now that we have dataBytes, Unmarshal it to QmpHeader */
	err = json.Unmarshal(dataBytes, out)
	if err != nil {
		return fmt.Errorf("[ReadQmpResult] json error: %s", err.Error())
	}

	return nil
}

func (monitor *QemuMonitor) QueryStatus(result *QmpQueryStatusResult) (err error) {
	const BufferSize int = 512
	var unix net.Conn
	var buffer []byte = make([]byte, BufferSize)
	var socketData []byte

	/* Initialize socket */
	log.Printf("[QueryStatus] initializing socket\n")
	unix, err = monitor.GetControlSocket()
	if err != nil {
		return err
	}

	/* Create QueryStatus command and send it */
	log.Printf("[QueryStatus] create qeury-status command\n")
	qmpCommand := &QmpBasicCommand{
		Execute: QmpQueryStatusCommand,
	}
	socketData, err = qmpCommand.GetJsonBytes()
	if err != nil {
		return err
	}

	log.Printf("[QueryStatus] sending query-status command")
	_, err = unix.Write(socketData)
	if err != nil {
		log.Printf("[QueryStatus] error sending command: %s\n", err.Error())
		return err
	}

	/* Now parse command result into result */
	log.Printf("[QueryStatus] reading command response")
	socketData = make([]byte, 0)
	for nBytes, err := unix.Read(buffer); err == nil && nBytes > 0; {
		socketData = append(socketData, buffer[:nBytes]...)
		if nBytes < BufferSize {
			break
		}
	}
	if err != nil {
		return err
	}

	log.Printf("[QueryStatus] socketData is '%s'", string(socketData))

	/* Unmarshal socketData into result */
	return json.Unmarshal(socketData, &result)
}

func (monitor *QemuMonitor) SendShutdownCommand() (err error) {
	var unix net.Conn
	var shutdownCommand QmpBasicCommand
	var qmpResult QmpBasicResult

	var socketData []byte
	var nBytes int

	log.Printf("[SendShutdownCommand] initializing socket")
	unix, err = monitor.GetControlSocket()
	if err != nil {
		return err
	}

	log.Printf("[SendShutdownCommand] sending shutdown command")
	shutdownCommand = QmpBasicCommand{
		Execute: "system_powerdown",
	}

	socketData, err = shutdownCommand.GetJsonBytes()
	if err != nil {
		return err
	}

	log.Printf("[SendShutdownCommand] unix is %v\n", unix)
	nBytes, err = unix.Write(socketData)
	if err != nil || nBytes == 0 {
		log.Printf("[SendShutdownCommand] could not write to socket: %s\n", err.Error())
		return err
	} else {
		log.Printf("[SendShutdownCommand] wrote %d bytes to socket\n", nBytes)
	}

	/* Now read QMP results */
	err = monitor.ReadQmpResult(unix, &qmpResult)
	if err != nil {
		log.Printf("could not read QMP result: %s\n", err.Error())
	}

	log.Printf("QmpBasicResult: %v\n", qmpResult)

	return unix.Close()
}
