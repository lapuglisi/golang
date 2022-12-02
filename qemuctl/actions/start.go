package qemuctl_actions

import (
	"flag"
	"fmt"
	"os"
	"os/exec"

	helpers "luizpuglisi.com/qemuctl/helpers"
	runtime "luizpuglisi.com/qemuctl/runtime"
)

func init() {
}

const (
	QemuDefaultSystemBin string = "qemu-system-x86_64"
)

type StartAction struct {
	machineName string
	configFile  string
	qemuBinary  string
}

func (action *StartAction) Run(arguments []string) (err error) {
	var flagSet *flag.FlagSet = flag.NewFlagSet("qemuctl start", flag.ExitOnError)

	flagSet.StringVar(&action.machineName, "name", "", "machine name")
	flagSet.StringVar(&action.configFile, "config", "", "YAML configuration file")
	flagSet.StringVar(&action.qemuBinary, "qemu", "qemu-system-x86_64", "qemu binary to use")

	err = flagSet.Parse(arguments)
	if err != nil {
		return err
	}

	/* Do flags validation */
	if len(action.configFile) == 0 {
		flagSet.Usage()
		return fmt.Errorf("--config is mandatory")
	}

	if len(action.machineName) == 0 {
		flagSet.Usage()
		return fmt.Errorf("--name is mandatory")
	}

	/* Do proper handling */
	err = action.handleStart()
	if err != nil {
		return err
	}

	return nil
}

func (action *StartAction) launchQemu(qemuBinary string, configData *helpers.ConfigurationData) (err error) {

	var qemuArgs []string = nil
	var qemuPath string

	var procAttrs *os.ProcAttr = nil

	qemuPath, err = exec.LookPath(qemuBinary)
	if err != nil {
		return err
	}

	qemuArgs, err = configData.GetQemuArgs(qemuPath, action.machineName)
	if err != nil {
		return err
	}

	// TODO: use the log feature
	/*
		fmt.Println("[INFO] Executing QEMU with:")
		fmt.Printf("qemu_path .......... %s\n", qemuPath)
		fmt.Printf("qemu_args .......... %s\n", strings.Join(qemuArgs, " "))
		return nil
	*/

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

func (action *StartAction) handleStart() (err error) {
	var configData *helpers.ConfigurationData = nil
	var qemuBinary string = action.qemuBinary

	err = nil

	configHandle := helpers.NewConfigHandler(action.configFile)

	configData, err = configHandle.ParseConfigFile()
	if err != nil {
		runtime.UpdateMachineStatus(action.machineName, "error")
		return err
	}

	// Now that we have configData, launch qemu
	err = action.launchQemu(qemuBinary, configData)
	if err != nil {
		runtime.UpdateMachineStatus(action.machineName, "error")
		return err
	}

	runtime.UpdateMachineStatus(action.machineName, "running")
	return nil
}
