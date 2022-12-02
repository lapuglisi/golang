package qemuctl_actions

import (
	"flag"
	"fmt"

	monitor "luizpuglisi.com/qemuctl/qemu"
	runtime "luizpuglisi.com/qemuctl/runtime"
)

func init() {

}

type StopAction struct {
	machineName string
}

func (action *StopAction) Run(arguments []string) (err error) {
	var flagSet *flag.FlagSet = flag.NewFlagSet("qemuctl stop", flag.ExitOnError)
	var machine *runtime.Machine

	flagSet.StringVar(&action.machineName, "name", "", "machine name")

	err = flagSet.Parse(arguments)
	if err != nil {
		return err
	}

	if len(action.machineName) == 0 {
		flagSet.Usage()
		return fmt.Errorf("--name is mandatory")
	}

	machine = runtime.NewMachine(action.machineName)

	err = monitor.ShutdownMachine(action.machineName)
	if err != nil {
		return err
	}

	// Now, update machine status
	machine.UpdateStatus(runtime.MachineStatusStopped)

	return nil
}
