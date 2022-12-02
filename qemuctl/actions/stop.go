package qemuctl_actions

import (
	"flag"
	"fmt"

	monitor "luizpuglisi.com/qemuctl/qemu"
)

func init() {

}

type StopAction struct {
	machineName string
}

func (action *StopAction) Run(arguments []string) (err error) {
	var flagSet *flag.FlagSet = flag.NewFlagSet("qemuctl stop", flag.ExitOnError)

	flagSet.StringVar(&action.machineName, "name", "", "machine name")

	err = flagSet.Parse(arguments)
	if err != nil {
		return err
	}

	if len(action.machineName) == 0 {
		flagSet.Usage()
		return fmt.Errorf("--name is mandatory")
	}

	err = monitor.ShutdownMachine(action.machineName)
	if err != nil {
		return err
	}

	return nil
}
