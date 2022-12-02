package qemuctl_actions

import (
	"flag"
	"fmt"
	"os/exec"

	helpers "luizpuglisi.com/qemuctl/helpers"
	qemuctl_qemu "luizpuglisi.com/qemuctl/qemu"
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

func (action *StartAction) handleStart() (err error) {
	var configData *helpers.ConfigurationData = nil
	var qemuBinary string = action.qemuBinary
	var qemuPath string
	var qemuArgs []string
	var qemu *qemuctl_qemu.QemuCommand
	var machine *runtime.Machine

	err = nil

	machine = runtime.NewMachine(action.machineName)

	/* Check machine status */
	if machine.Exists() {
		if machine.IsStarted() {
			return fmt.Errorf("machine '%s' is already started", action.machineName)
		}
	}

	configHandle := helpers.NewConfigHandler(action.configFile)
	configData, err = configHandle.ParseConfigFile()
	if err != nil {
		machine.UpdateStatus(runtime.MachineStatusDegraded)
		return err
	}

	// Get qemu real path
	qemuPath, err = exec.LookPath(qemuBinary)
	if err != nil {
		return err
	}

	// Now that we have configData, launch qemu
	qemuArgs, err = configData.GetQemuArgs(qemuPath, action.machineName)
	if err != nil {
		return err
	}

	/* Get QemuCommand instance */
	qemu = qemuctl_qemu.NewQemuCommand(qemuBinary, qemuArgs, configData.RunAsDaemon)
	err = qemu.Launch()
	if err != nil {
		machine.UpdateStatus(runtime.MachineStatusDegraded)
		return err
	}

	machine.UpdateStatus(runtime.MachineStatusStarted)
	return nil
}
