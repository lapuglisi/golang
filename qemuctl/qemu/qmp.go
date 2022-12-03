package qemuctl_qemu

import (
	"encoding/json"
	"log"
	"net"
)

const (
	QmpCapabilitiesCommand    string = "qmp_capabilities"
	QmpQueryStatusCommand     string = "query-status"
	QmpSystemPowerdownCommand string = "system_powerdown"
	QmpDefaultBufferSize      int    = 1024
)

type QmpHeaderVersionQemu struct {
	Micro int `json:"micro"`
	Minor int `json:"minor"`
	Major int `json:"major"`
}

type QmpHeaderVersionData struct {
	Qemu QmpHeaderVersionQemu `json:"qemu"`
}

type QmpHeader struct {
	QMP struct {
		Version QmpHeaderVersionData `json:"version"`
		Package string               `json:"package"`
	} `json:"QMP"`
}

type QmpCommandArguments struct {
	Device string `json:"device"`
}

type QmpBasicCommand struct {
	Command string `json:"execute"`
}

type QmpCommandQueryStatus struct {
	Command string `json:"execute"`
}

type QmpQueryStatusResult struct {
	Return struct {
		Status     string `json:"status"`
		SingleStep bool   `json:"singlestep"`
		Running    bool   `json:"running"`
	} `json:"return"`
}

type QmpQueryStatusReturn struct {
	Return QmpQueryStatusResult `json:"result"`
}

type QmpBasicResult struct {
	Result interface{} `json:"result"`
}

type QmpEventResult struct {
	Result struct {
		Timestamp struct {
			Seconds      int64 `json:"seconds"`
			Microseconds int   `json:"microseconds"`
		} `json:"timestamp"`
		Event string      `json:"event"`
		Data  interface{} `json:"data"`
	} `json:"result"`
}

func (command *QmpBasicCommand) Execute(socket net.Conn) (result *QmpBasicResult, err error) {
	var jsonBytes []byte
	var buffer []byte = make([]byte, QmpDefaultBufferSize)
	var nBytes int

	result = &QmpBasicResult{}

	/* Marshal struct into bytes */
	log.Printf("[QmpCommand] execute: creating json bytes")
	jsonBytes, err = json.Marshal(command)
	if err != nil {
		return nil, err
	}

	log.Printf("[QmpCommand] execute: json bytes is [%s]", string(jsonBytes))

	/* Send command data to socket */
	log.Printf("[QmpCommand] execute: sending command to socket")
	nBytes, err = socket.Write(jsonBytes)
	if err != nil {
		return nil, err
	}

	/* Read return JSON from socket */
	log.Printf("[QmpCommand] execute: reading result from socket")
	jsonBytes = make([]byte, 0)

	for nBytes, err = socket.Read(buffer); err == nil && nBytes > 0; {
		jsonBytes = append(jsonBytes, buffer[:nBytes]...)
		if nBytes < QmpDefaultBufferSize {
			break
		}
	}

	/* Unmarshal jsonBytes into result */
	err = json.Unmarshal(jsonBytes, &result)

	log.Printf("[QmpCommand] execute: result is [%v]", result)

	return result, err
}

func (command *QmpCommandQueryStatus) Execute(socket net.Conn) (result *QmpQueryStatusResult, err error) {
	var jsonBytes []byte
	var buffer []byte = make([]byte, QmpDefaultBufferSize)
	var nBytes int

	result = &QmpQueryStatusResult{}

	/* Marshal struct into bytes */
	log.Printf("[QmpCommandQuery] execute: creating json bytes")

	command.Command = QmpQueryStatusCommand

	jsonBytes, err = json.Marshal(command)
	if err != nil {
		return nil, err
	}

	log.Printf("[QmpCommandQuery] execute: json bytes is [%s]", string(jsonBytes))

	/* Send command data to socket */
	log.Printf("[QmpCommandQuery] execute: sending command to socket")
	_, err = socket.Write(jsonBytes)
	if err != nil {
		return nil, err
	}

	/* Read return JSON from socket */
	log.Printf("[QmpCommandQuery] execute: reading result from socket")
	jsonBytes = make([]byte, 0)

	for nBytes, err = socket.Read(buffer); err == nil && nBytes > 0; {
		jsonBytes = append(jsonBytes, buffer[:nBytes]...)
		if nBytes < QmpDefaultBufferSize {
			break
		}
	}

	/* Unmarshal jsonBytes into result */
	err = json.Unmarshal(jsonBytes, &result)

	log.Printf("[QmpCommandQuery] execute: result is [%v]", result)

	return result, err
}

func (event *QmpEventResult) ReadEvent(socket net.Conn) bool {
	var buffer []byte = make([]byte, QmpDefaultBufferSize)
	var jsonData []byte = make([]byte, 0)
	var err error
	var nBytes int

	log.Printf("[ReadEvent] reading data from socket")
	for nBytes, err = socket.Read(buffer); err == nil && nBytes > 0; {
		jsonData = append(jsonData, buffer[:nBytes]...)

		if nBytes < QmpDefaultBufferSize {
			break
		}
	}

	if err != nil || nBytes == 0 {
		return false
	}

	log.Printf("[ReadEvent] received from socket: [%s]", string(jsonData))

	err = json.Unmarshal(jsonData, &event)

	return (err == nil)
}
