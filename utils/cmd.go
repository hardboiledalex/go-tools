package utils

import (
	"github.houston.softwaregrp.net/Hercules/gotools/utils/logging"
	"os"
	"os/exec"
	"strings"
)

//Run shell command and fail if error
func ExecShell(command string) {
	logging.Log(logging.DEBUG, "Running command ", command, " locally")
	cmd := exec.Command("sh", "-c", command)
	err := cmd.Run()
	if err != nil {
		logging.LogErrorf("Cannot run local shell command. Error: %s ", err.Error())
		os.Exit(1)
	}
}

//Run shell command in interactive mode
func ExecShellInteractive(command string) {
	cmd := exec.Command("sh", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		logging.LogErrorf("Cannot run shell command. Error: %s ", err.Error())
		os.Exit(1)
	}
}

func ExecShellInteractiveWithError(command string) error {
	cmd := exec.Command("sh", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err
}

//Run shell command and return output
func ExecShellWithOutput(command string) (string, error) {
	out, err := exec.Command("sh", "-c", command).CombinedOutput()
	return strings.TrimSuffix(string(out[:]), "\n"), err
}
