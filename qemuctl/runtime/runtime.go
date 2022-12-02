package qemuctl_runtime

import (
	"fmt"
	"log"
	"os"
)

func init() {

}

func GetUserDataDir() string {
	return fmt.Sprintf("%s/.qemuctl", os.ExpandEnv("$HOME"))
}

func GetMachineDirectory(machineName string) string {
	return fmt.Sprintf("%s/machines/%s", GetUserDataDir(), machineName)
}

func SetupRuntimeData() (err error) {
	var qemuctlDir string = GetUserDataDir()

	log.Printf("checking for directory '%s'\n", qemuctlDir)

	/* Create directory {userHome}/.qemuctl if it does not exits */
	_, err = os.Stat(qemuctlDir)
	if os.IsNotExist(err) {
		/* Create qemuctl directory */
		log.Printf("creating directory '%s'\n", qemuctlDir)

		err = os.Mkdir(qemuctlDir, os.ModeDir|os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}

func MachineExists(machineName string) bool {
	var err error
	var machineDir string = GetMachineDirectory(machineName)

	/* Check if machine exists */
	_, err = os.Stat(machineDir)
	return !os.IsNotExist(err)
}

func UpdateMachineStatus(machineName string, status string) error {
	var statusFile string = fmt.Sprintf("%s/status", GetMachineDirectory(machineName))

	return os.WriteFile(statusFile, []byte(status), os.ModePerm)
}
