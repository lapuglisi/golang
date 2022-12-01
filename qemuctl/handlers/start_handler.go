package qemuctl_handlers

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	helpers "luizpuglisi.com/qemuctl/helpers"
)

func init() {
}

const (
	QemuDefaultSystemBin string = "qemu-system-x86_64"
)

func appendQemuArg(qemuArgs []string, argKey string, argValue string) (newArgs []string) {
	return append(qemuArgs, []string{argKey, argValue}...)
}

func launchQemu(qemuBinary string, configData *helpers.ConfigurationData) (err error) {

	var qemuArgs []string = nil
	var qemuPath string

	var procAttrs *os.ProcAttr = nil

	qemuPath, err = exec.LookPath(qemuBinary)
	if err != nil {
		return err
	}

	qemuArgs, err = configData.GetQemuArgs(qemuPath)
	if err != nil {
		return err
	}

	fmt.Println("[INFO] Executing QEMU with:")
	fmt.Printf("qemu_path .......... %s\n", qemuPath)
	fmt.Printf("qemu_args .......... %s\n", strings.Join(qemuArgs, " "))

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
	var configData *helpers.ConfigurationData = nil
	var qemuBinary string = QemuDefaultSystemBin

	err = nil
	if len(startArgs) < 1 {
		err = fmt.Errorf("[handle_start] Too few arguments")
		return err
	}

	configFile = startArgs[0]
	configHandle := helpers.NewConfigHandler(configFile)

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
