package main

import (
	"fmt"
	"os"

	qemuctl "luizpuglisi.com/qemuctl/handlers"
)

func main() {
	fmt.Println("Serominers nas alturas")

	var err error
	var execArgs []string = os.Args
	var argsLen int = len(execArgs)

	if argsLen < 2 {
		fmt.Printf("Too few arguments. Got %d, need more than 2\n", argsLen)
		os.Exit(1)
	}

	var action string = execArgs[1]

	switch action {
	case "start":
		{
			err = qemuctl.HandleStart(execArgs[2:])
		}

	default:
		{
			fmt.Printf("[error] Unknown action '%s'.", action)
		}
	}

	if err != nil {
		fmt.Printf("[\033[31merror\033[0m] %s\n", err.Error())
	}

	os.Exit(0)
}
