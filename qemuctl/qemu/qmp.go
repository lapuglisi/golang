package qemuctl_qemu

import "encoding/json"

const (
	QmpCapabilitiesCommand    string = "qmp_capabilities"
	QmpQueryStatusCommand     string = "query-status"
	QmpSystemPowerdownCommand string = "system_powerdown"
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
	Execute string `json:"execute"`
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

func (command *QmpBasicCommand) GetJsonBytes() (jsonBytes []byte, err error) {
	jsonBytes, err = json.Marshal(command)
	if err != nil {
		jsonBytes = nil
	}

	return jsonBytes, err
}
