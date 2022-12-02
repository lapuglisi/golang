package qemuctl_actions

import (
	"flag"
	"fmt"
	"log"

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

	flagSet.StringVar(&action.qemuBinary, "qemu", QemuDefaultSystemBin, "qemu binary")
	flagSet.StringVar(&action.configFile, "config", "", "YAML configuration file")

	err = flagSet.Parse(arguments)
	if err != nil {
		return err
	}

	/* Check for machine name */
	if len(flagSet.Args()) < 1 {
		return fmt.Errorf("machine name is mandatory")
	}
	action.machineName = flagSet.Arg(0)

	log.Printf("[start::Run] action.machineName is '%s'\n", action.machineName)

	/* Do flags validation */
	if len(action.configFile) == 0 {
		flagSet.Usage()
		return fmt.Errorf("--config is mandatory")
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
	var qemu *qemuctl_qemu.QemuCommand
	var machine *runtime.Machine

	err = nil

	machine = runtime.NewMachine(action.machineName)

	log.Printf("qemuctl: starting machine %s...\n", action.machineName)

	/* Check machine status */
	if machine.Exists() {
		if machine.IsStarted() {
			return fmt.Errorf("machine '%s' is already started", action.machineName)
		}
	} else {
		machine.CreateRuntime()
	}

	configHandle := helpers.NewConfigHandler(action.configFile)
	configData, err = configHandle.ParseConfigFile()
	if err != nil {
		machine.UpdateStatus(runtime.MachineStatusDegraded)
		return err
	}

	/* Get QemuCommand instance */
	qemuMonitor := qemuctl_qemu.NewQemuMonitor(machine)
	qemu = qemuctl_qemu.NewQemuCommand(qemuBinary, configData, qemuMonitor)
	err = qemu.Launch()
	if err != nil {
		machine.UpdateStatus(runtime.MachineStatusDegraded)
		return err
	}

	machine.UpdateStatus(runtime.MachineStatusStarted)
	return nil
}
