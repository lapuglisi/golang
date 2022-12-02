package qemuctl_qemu

import (
	"log"
	"os"
	"strings"
)

type QemuCommand struct {
	QemuPath    string
	QemuArgs    []string
	RunAsDaemon bool
}

func NewQemuCommand(qemuPath string, qemuArgs []string, runAsDaemon bool) (qemu *QemuCommand) {
	return &QemuCommand{
		QemuPath:    qemuPath,
		QemuArgs:    qemuArgs,
		RunAsDaemon: runAsDaemon,
	}
}

func (qemu *QemuCommand) Launch() (err error) {
	var procAttrs *os.ProcAttr

	// TODO: use the log feature
	log.Println("[QemuCommand::Launch] Executing QEMU with:")
	log.Printf("qemu_path .......... %s\n", qemu.QemuPath)
	log.Printf("qemu_args .......... %s\n", strings.Join(qemu.QemuArgs, " "))

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

	procHandle, err := os.StartProcess(qemu.QemuPath, qemu.QemuArgs, procAttrs)
	if err == nil {
		if qemu.RunAsDaemon {
			err = procHandle.Release()
		}
	}

	return err
}
