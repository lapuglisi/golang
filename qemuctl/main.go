package main

import (
	"fmt"
	"os"

	actions "luizpuglisi.com/qemuctl/actions"
	runtime "luizpuglisi.com/qemuctl/runtime"
)

func usage() {
	fmt.Println()
	fmt.Println("usage:")
	fmt.Println("    qemuctl {start|stop|seila} OPTIONS")
}

func main() {
	var err error

	var execArgs []string = os.Args
	var action string

	/* Initialize qemuctl */
	err = runtime.SetupRuntimeData()
	if err != nil {
		fmt.Printf("[\033[31merror\033[0m] %s\n", err.Error())
	}

	if len(execArgs) < 2 {
		usage()
		os.Exit(1)
	}

	action = execArgs[1]
	execArgs = execArgs[2:]

	fmt.Println("")

	switch action {
	case "create":
		{
			action := actions.CreateAction{}
			err = action.Run(execArgs)
			break
		}
	case "destroy":
		{
			action := actions.DestroyAction{}
			err = action.Run(execArgs)
			break
		}
	case "start":
		{
			action := actions.StartAction{}
			err = action.Run(execArgs)
			break
		}

	case "stop":
		{
			action := actions.StopAction{}
			err = action.Run(execArgs)
			break
		}
	case "status":
		{
			action := actions.StatusAction{}
			err = action.Run(execArgs)
			break
		}
	case "edit":
		{
			action := actions.EditAction{}
			err = action.Run(execArgs)
		}
	case "list":
		{
			action := actions.ListAction{}
			err = action.Run(execArgs)
		}
	default:
		{
			fmt.Printf("[error] Unknown action '%s'\n", action)
		}
	}

	if err != nil {
		fmt.Printf("[\033[31merror\033[0m] %s\n", err.Error())
	}

	os.Exit(0)
}
