package main

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"pacaur/aurjson"
	"strings"
)

const (
	AURMainUrl     string = "https://aur.archlinux.org"
	PackageName    string = "pacaur"
	TempDirPattern string = "pacaur-install-*"
)

func usage() {
	fmt.Println()
	fmt.Println("usage:")
	fmt.Printf("    %s {install | search | upgrade | list}\n", PackageName)
}

func main() {
	if len(os.Args) < 2 {
		usage()

		return
	}

	switch os.Args[1] {
	case "install":
		{
			doAurInstall(os.Args[2:])
			break
		}

	case "search":
		{
			if len(os.Args) < 3 {
				usage()
				break
			}

			doAurSearch(os.Args[2])
			break
		}

	case "upgrade":
		{
			doAurUpgrade()
			break
		}

	case "list":
		{
			doAurList()
			break
		}

	default:
		{
			usage()
		}
	}

}

func executeMakepkg(packageDir string, args ...string) (err error) {
	var cmdBin string
	var aurCmd *exec.Cmd = nil

	cmdBin, err = exec.LookPath("makepkg")
	if err != nil {
		return err
	}

	aurCmd = exec.Command(cmdBin, args...)
	aurCmd.Dir = packageDir
	aurCmd.Stdin = os.Stdin
	aurCmd.Stderr = os.Stderr
	aurCmd.Stdout = os.Stdout

	return aurCmd.Run()
}

func executePacman(args ...string) (err error) {
	var cmdBin string
	var aurCmd *exec.Cmd = nil
	var cmdArgs []string = make([]string, 1)

	cmdBin, err = exec.LookPath("sudo")
	if err != nil {
		return err
	}

	cmdArgs[0] = "pacman"
	cmdArgs = append(cmdArgs, args...)

	aurCmd = exec.Command(cmdBin, cmdArgs...)
	aurCmd.Stdin = os.Stdin
	aurCmd.Stderr = os.Stderr
	aurCmd.Stdout = os.Stdout

	return aurCmd.Run()
}

func doAurSearch(search string) {
	var result *aurjson.AURJson

	fmt.Printf("[\033[33mpacaur\033[0m] Searching for '\033[33m%s\033[0m'\n", search)

	result, err := aurjson.SearchAUR(search)
	if err != nil {
		fmt.Printf("\033[31mERROR\033[0m: Search failed: %s\n", err.Error())
		return
	}

	if result.ResultCount <= 0 {
		fmt.Printf("[\033[33m%s\033[0m] No results found for '%s'\n", PackageName, search)
		return
	}

	fmt.Println()
	for _, aurPack := range result.Results {
		fmt.Printf("\033[36maur\033[0m/%s (\033[33m%s\033[0m)\n", aurPack.Name, aurPack.Version)
		fmt.Printf("  %s\n", aurPack.Description)
		fmt.Println()
	}
}

func doAurInstall(aurPacks []string) {

	for _, aurPack := range aurPacks {
		fmt.Println()
		fmt.Printf("[\033[33m%s\033[0m] Installing package \033[1m%s\033[0m\n", PackageName, aurPack)

		result, err := aurjson.GetPackageInfo(aurPack)
		if err != nil {
			fmt.Printf("\033[31mERROR\033[0m: Could not get info for package '%s': %s\n", aurPack, err.Error())
			continue
		}

		if result.ResultCount != 1 {
			result_error := result.Error
			if len(result_error) == 0 {
				result_error = "404"
			}

			fmt.Printf("[\033[33m%s\033[0m] Package \033[1m%s\033[0m not found: %s\n",
				PackageName, aurPack, result_error)
			continue
		}

		packageInfo := result.Results[0]

		fmt.Printf("  Package URL: %s%s\n", AURMainUrl, packageInfo.URLPath)

		downloadAndInstall(packageInfo)
	}
}

func downloadAndInstall(packageInfo aurjson.AURSearchResult) (err error) {

	var aurCmd *exec.Cmd = nil
	var cmdBin string
	var aurUrl string
	var aurFileName string
	var aurFileData []byte
	var resp *http.Response

	// -- Create temporary directory
	fmt.Printf(" [\033[36mpacaur\033[0m] Creating temporary directory\n")

	tempDir, err := os.MkdirTemp(os.TempDir(), TempDirPattern)
	if err != nil {
		fmt.Printf("\033[31mERROR\033[0m: Could not create temporary directory.\n")
		return err
	}

	// try to install package dependencies
	pacmanArgs := []string{"-Sy", "--noconfirm", "--needed"}
	pacmanArgs = append(pacmanArgs, packageInfo.Depends...)
	err = executePacman(pacmanArgs...)

	if err != nil {
		fmt.Printf("\033[31mERROR\033[0m: Could not install dependencies of '%s' (%s)\n", packageInfo.PackageBase, err.Error())
		goto cleanup
	}

	// -- Download package into tempDir
	aurUrl = fmt.Sprintf("%s%s", AURMainUrl, packageInfo.URLPath)
	fmt.Printf(" [\033[36mpacaur\033[0m] Downloading package from '%s'\n", aurUrl)

	resp, err = http.Get(aurUrl)
	if err != nil {
		fmt.Printf("\033[31mERROR\033[0m: Could not download from '%s'\n", aurUrl)
		goto cleanup
	}

	defer resp.Body.Close()

	// -- Reading data from HTTP request
	aurFileData, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("\033[31mERROR\033[0m: Could not read data from '%s'\n", aurUrl)
		goto cleanup
	}

	aurFileName = fmt.Sprintf("%s/%s.tgz", tempDir, packageInfo.Name)
	fmt.Printf(" [\033[36mpacaur\033[0m] Creating temporary file '%s'\n", aurFileName)

	err = os.WriteFile(aurFileName, aurFileData, fs.FileMode(os.O_CREATE)|fs.FileMode(os.O_RDWR))
	if err != nil {
		fmt.Printf("\033[31mERROR\033[0m: Could not write to file '%s'\n", aurFileName)
		goto cleanup
	}

	os.Chmod(aurFileName, fs.FileMode(0775))

	// -- By now, we should have our tgz downloaded
	// -- Unpack it and do the mambo jambo
	fmt.Printf(" [\033[36mpacaur\033[0m] Unpacking tarball '%s'\n", aurFileName)

	cmdBin, err = exec.LookPath("tar")
	if err != nil {
		fmt.Printf("\033[31mERROR\033[0m: Could not unpack tarball '%s': %s\n", aurFileName, err.Error())
		goto cleanup
	}

	aurCmd = exec.Command(cmdBin, "--overwrite", "-xvf", aurFileName)
	aurCmd.Dir = tempDir
	aurCmd.Stdin = os.Stdin
	aurCmd.Stderr = os.Stderr

	err = aurCmd.Run()
	if err != nil {
		fmt.Printf("\033[31mERROR\033[0m: Could not unpack tarball '%s': %s\n", aurFileName, err.Error())
		goto cleanup
	}

	// -- Execute makepkg on the file and show no mercy
	fmt.Printf(" [\033[32mpacaur\033[0m] Installing AUR package '\033[35m%s\033[0m' with 'makepkg'\n", packageInfo.Name)

	err = executeMakepkg(fmt.Sprintf("%s/%s", tempDir, packageInfo.Name), "-sif", "--noconfirm")
	if err != nil {
		fmt.Printf("\033[31mERROR\033[0m: Could not execute makepkg on '%s': %s\n", tempDir, err.Error())
		goto cleanup
	}

cleanup:
	// -- Time to remove temporary directory
	fmt.Printf(" [\033[36mpacaur\033[0m] Removing temporary directory '%s'\n", tempDir)
	err = os.RemoveAll(tempDir)
	if err != nil {
		fmt.Printf("\033[31mERROR\033[0m: Could not remove temporary directory '%s': %s\n", tempDir, err.Error())
	}

	return err
}

func doAurUpgrade() (err error) {

	var aurCmd *exec.Cmd = nil
	var aurCmdPath string

	fmt.Printf("[\033[32mpacaur\033[0m] Performing full AUR upgrade\n")

	// List all pacman -Qm packages
	aurCmdPath, err = exec.LookPath("pacman")
	if err != nil {
		fmt.Printf("\033[31mERROR\033[0m: Could not upgrade AUR packages: %s\n", err.Error())
		return err
	}

	aurCmd = exec.Command(aurCmdPath, "-Qm")
	aurCmdReader, err := aurCmd.StdoutPipe()
	if err != nil {
		fmt.Printf("\033[31mERROR\033[0m: Could not upgrade AUR packages: %s\n", err.Error())
		return err
	}

	err = aurCmd.Start()
	if err != nil {
		fmt.Printf("\033[31mERROR\033[0m: Could not upgrade AUR packages: %s\n", err.Error())
		return err
	}

	// For every other item on string, starting from 0
	dataScan := bufio.NewScanner(aurCmdReader)
	for dataScan.Scan() {
		aurCmdLinePackage := strings.Split(string(dataScan.Text()), " ")

		fmt.Println()
		fmt.Printf(" [\033[34mpacaur\033[0m] Upgrading AUR package '\033[1m%s\033[0m'\n", aurCmdLinePackage[0])

		// Get package info
		packageInfo, err := aurjson.GetPackageInfo(aurCmdLinePackage[0])
		if err != nil {
			fmt.Printf("\033[31mERROR\033[0m: Could not upgrade AUR package %s: %s\n", aurCmdLinePackage[0], err.Error())
			continue
		}

		if packageInfo.Type == "error" || packageInfo.ResultCount != 1 {
			fmt.Printf("\033[31mERROR\033[0m: Could not upgrade AUR package %s: %s\n", aurCmdLinePackage[0], packageInfo.Error)
			continue
		}

		if strings.Compare(packageInfo.Results[0].Version, aurCmdLinePackage[1]) == 0 {
			fmt.Printf(" [\033[35mpacaur\033[0m] AUR package '\033[1m%s\033[0m' is up-to-date\n", aurCmdLinePackage[0])
		} else {
			fmt.Printf(" [\033[35mpacaur\033[0m] '\033[1m%s\033[0m': \033[33m%s\033[0m -> \033[32m%s\033[0m\n",
				aurCmdLinePackage[0], aurCmdLinePackage[1], packageInfo.Results[0].Version)

			downloadAndInstall(packageInfo.Results[0])
		}

	}

	return nil
}

func doAurList() {

	var aurCmd *exec.Cmd = nil
	var aurCmdPath string
	var err error

	// List all pacman -Qm packages
	aurCmdPath, err = exec.LookPath("pacman")
	if err != nil {
		fmt.Printf("\033[31mERROR\033[0m: Could not list AUR packages: %s\n", err.Error())
		return
	}

	aurCmd = exec.Command(aurCmdPath, "-Qm", "--color", "always")
	aurCmd.Stdout = os.Stdout

	err = aurCmd.Run()
	if err != nil {
		fmt.Printf("\033[31mERROR\033[0m: Could not list AUR packages: %s\n", err.Error())
		return
	}
}
