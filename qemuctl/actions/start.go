package qemuctl_actions

import (
	"fmt"
	"log"

	helpers "luizpuglisi.com/qemuctl/helpers"
	qemuctl_qemu "luizpuglisi.com/qemuctl/qemu"
	runtime "luizpuglisi.com/qemuctl/runtime"
)

func init() {
}

type StartAction struct {
	machineName string
	configFile  string
	qemuBinary  string
}

func (action *StartAction) Run(arguments []string) (err error) {
	/* Check for machine name */
	if len(arguments) < 1 {
		return fmt.Errorf("machine name is mandatory")
	}
	action.machineName = arguments[0]

	fmt.Printf("[start] starting machine '%s'...", action.machineName)

	/* Do proper handling */
	err = action.handleStart()
	if err != nil {
		fmt.Printf(" \033[31;1merror\033[0m: %s\n", err.Error())
		return err
	}

	fmt.Printf(" \033[32;1mok!\033[0m\n")
	return nil
}

func (action *StartAction) handleStart() (err error) {
	var machine *runtime.Machine

	log.Printf("[start] starting machine '%s'", action.machineName)
	machine = runtime.NewMachine(action.machineName)

	if !machine.Exists() {
		return fmt.Errorf("machine '%s' dos not exist", action.machineName)
	}

	if machine.IsStarted() {
		return fmt.Errorf("[start] machine '%s' is already started", action.machineName)
	}

	/* in this release, starting a machine means creating it again */
	log.Printf("[start] relaunching machine '%s' (%s)", machine.Name, machine.ConfigFile)

	log.Printf("[start] parsing config file '%s'", machine.ConfigFile)
	configHandle := helpers.NewConfigHandler(machine.ConfigFile)
	configData, err := configHandle.ParseConfigFile()
	if err != nil {
		return err
	}

	log.Printf("[start] creating qemuMonitor instance")
	qemuMonitor := qemuctl_qemu.NewQemuMonitor(machine)

	log.Printf("[start] launching qemu command")
	qemu := qemuctl_qemu.NewQemuCommand(configData, qemuMonitor)

	qemuPid, err := qemu.Launch()

	if err == nil {
		machine.QemuPid = qemuPid
		machine.SSHLocalPort = configData.SSH.LocalPort
		machine.UpdateStatus(runtime.MachineStatusStarted)
	} else {
		machine.UpdateStatus(runtime.MachineStatusDegraded)
	}

	return err
}
