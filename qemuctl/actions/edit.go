package qemuctl_actions

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	runtime "luizpuglisi.com/qemuctl/runtime"
)

type EditAction struct {
	machineName string
}

func (action *EditAction) Run(arguments []string) (err error) {
	var machine *runtime.Machine

	/* Check for machine name */
	if len(arguments) < 1 {
		return fmt.Errorf("machine name is mandatory")
	}
	action.machineName = arguments[0]

	fmt.Printf("[edit] editing machine '%s'...", action.machineName)

	machine = runtime.NewMachine(action.machineName)

	if !machine.Exists() {
		return fmt.Errorf("machine '%s' does not exist", action.machineName)
	}

	if machine.IsStarted() {
		return fmt.Errorf("cannot edit a running machine ('%s' is started)", action.machineName)
	}

	log.Printf("[edit] looking for a valid EDITOR")
	editorBin := os.ExpandEnv("$EDITOR")
	if len(editorBin) == 0 {
		editorBin = "vim"
	}

	editorPath, err := exec.LookPath(editorBin)
	if err != nil {
		return err
	}

	log.Printf("[edit] using editor '%s'", editorPath)
	log.Printf("[edit] launching '%s %s'", editorPath, machine.ConfigFile)

	err = nil
	procAttrs := &os.ProcAttr{
		Dir: machine.RuntimeDirectory,
		Env: os.Environ(),
		Files: []*os.File{
			os.Stdin,
			os.Stdout,
			os.Stderr,
		},
		Sys: nil,
	}

	procHandle, err := os.StartProcess(editorPath, []string{editorPath, machine.ConfigFile}, procAttrs)
	if err != nil {
		return err
	}

	procHandle.Wait()

	log.Printf("[edit] action executed successfully")
	return nil
}
