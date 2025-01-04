package sudo

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/jozefcipa/novus/internal/fs"
	"github.com/jozefcipa/novus/internal/logger"
	"github.com/jozefcipa/novus/internal/paths"
)

type SudoCommand string

// These are the commands that are available in the `sudo-helper` bash script
const (
	CheckPorts SudoCommand = "check-ports"
	MakeDir    SudoCommand = "mkdir"
	RemoveFile SudoCommand = "rm"
	Chown      SudoCommand = "chown"
	Touch      SudoCommand = "touch"
)

var hasSudoersFile bool

func init() {
	hasSudoersFile = false
}

func sudo(command SudoCommand, args []string) error {
	cmdString := append([]string{paths.SudoHelperPath, string(command)}, args...)
	logger.Debugf("Running \"sudo %s\"", strings.Join(cmdString, " "))

	if out, err := exec.Command("sudo", cmdString...).Output(); err != nil {
		return errors.New(string(out))
	}

	return nil
}

func ensureSudoHelper() {
	if hasSudoersFile {
		return
	}

	if !fs.FileExists(paths.SudoHelperPath) {
		logger.Debugf("Sudo helper doesn't exist, creating one now.")
		createSudoHelper(paths.SudoAllowedPaths)
	}

	hasSudoersFile = true
}

func createSudoHelper(allowedPaths []string) {
	// Read sudo helper template content
	sudoHelperContent := fs.ReadFileOrExit(filepath.Join(paths.AssetsDir, "sudo-helper.template.sh"))

	// Replace variables
	sudoHelperContent = strings.ReplaceAll(
		sudoHelperContent,
		"--ALLOWED-PATHS--",
		// Define all directories that can be modified by passwordless `sudo`
		strings.Join(allowedPaths, " "),
	)

	// Create sudo helper file
	logger.Infof("Creating sudo helper")
	logger.Debugf("Creating sudo-helper file [%s]", paths.SudoHelperPath)
	if _, err := exec.Command("sudo", "touch", paths.SudoHelperPath).Output(); err != nil {
		logger.Errorf("Failed to create file %s\n   Reason: %v", paths.SudoHelperPath, err)
		os.Exit(1)
	}

	// We need to change the file owner to the current user in order to be able to write to the file
	user, _ := user.Current()
	if _, err := exec.Command("sudo", "chown", user.Username, paths.SudoHelperPath).Output(); err != nil {
		logger.Errorf("Failed to call `sudo chown %s %s`\n   Reason: %v", user.Username, paths.SudoHelperPath, err)
		os.Exit(1)
	}

	if err := os.WriteFile(paths.SudoHelperPath, []byte(sudoHelperContent), 0644); err != nil {
		logger.Errorf("Failed to write to a file %s\n   Reason: %v", paths.SudoHelperPath, err)
		os.Exit(1)
	}

	// Make it executable
	logger.Debugf("Making sudo-helper executable (0744)")
	err := os.Chmod(paths.SudoHelperPath, 0744)
	if err != nil {
		logger.Errorf("Failed to make sudo-helper executable: %v", err)
		os.Exit(1)
	}

	// Change the sudo-helper ownership to root to avoid direct file modification by users,
	// as this file now has passwordless sudo permissions
	logger.Debugf("Changing ownership of sudo-helper to root")
	if _, err := exec.Command("sudo", "chown", "root", paths.SudoHelperPath).Output(); err != nil {
		logger.Errorf("Failed to call `sudo chown root %s`\n   Reason: %v", paths.SudoHelperPath, err)
		os.Exit(1)
	}
	logger.Debugf("Sudo helper has been created.")
}

func RegisterSudoersFile() {
	user, _ := user.Current()
	sudoersContent := fmt.Sprintf("%s ALL=(ALL) NOPASSWD: %s", user.Username, paths.SudoHelperPath)

	logger.Debugf("Writing to %s file", paths.SudoersFilePath)
	WriteFileOrExit(paths.SudoersFilePath, sudoersContent)

	logger.Debugf("Changing ownership of %s to root", paths.SudoersFilePath)
	ChownOrExit("root", paths.SudoersFilePath)
}

func CheckPortsOrExit(ports []string) string {
	ensureSudoHelper()

	commandString := []string{paths.SudoHelperPath, string(CheckPorts), strings.Join(ports, ",")}
	logger.Debugf("Running \"sudo %s\"", strings.Join(commandString, " "))

	cmd := exec.Command("sudo", commandString...)
	out, err := cmd.CombinedOutput()
	result := string(out)

	// `lsof` command returns exit code 1 even if there is no script error but no information was found
	// therefore, let's only return error if the exit code is 1 and there is some actual output
	// https://stackoverflow.com/a/29843137/4480179
	if err != nil && result != "" {
		logger.Errorf("Failed to run \"%s\": %v\n%s", commandString, err, result)
		os.Exit(1)
	}

	return result
}

func MakeDirOrExit(filePath string) {
	ensureSudoHelper()

	if err := sudo(MakeDir, []string{filePath}); err != nil {
		logger.Errorf("Failed to create directory %s\n  Reason: %v", filePath, err)
		os.Exit(1)
	}
}

func DeleteFile(filePath string) error {
	ensureSudoHelper()

	if err := sudo(RemoveFile, []string{filePath}); err != nil {
		return errors.New(fmt.Sprintf("Failed to delete file %s\n   Reason: %v", filePath, err))
	}

	return nil
}

func ChownOrExit(userName string, filePath string) {
	ensureSudoHelper()

	if err := sudo(Chown, []string{userName, filePath}); err != nil {
		logger.Errorf("Failed to call `chown` on file %s\n  Reason: %v", filePath, err)
		os.Exit(1)
	}
}

func WriteFileOrExit(filePath string, data string) {
	ensureSudoHelper()

	if err := sudo(Touch, []string{filePath}); err != nil {
		logger.Errorf("Failed to create file %s\n  Reason: %v", filePath, err)
		os.Exit(1)
	}

	// We need to change the file owner to the current user in order to be able to write to the file
	user, _ := user.Current()
	ChownOrExit(user.Username, filePath)
	err := os.WriteFile(filePath, []byte(data), 0644)

	if err != nil {
		logger.Errorf("Failed to write to a file %s\n  Reason: %v", filePath, err)
		os.Exit(1)
	}
}
