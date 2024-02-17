package fs

import (
	"os"
	"os/exec"
	"os/user"

	"github.com/jozefcipa/novus/internal/logger"
)

var UserHomeDir string

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		logger.Errorf("Failed to get user home directory\n%v\n", err)
		os.Exit(1)
	}

	UserHomeDir = homeDir
}

func ReadFileOrExit(path string) string {
	file, err := os.ReadFile(path)
	if err != nil {
		logger.Errorf("Failed to read a file\n%v\n", err)
		os.Exit(1)
	}

	return string(file)
}

func ReadFile(path string) (string, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(file), nil
}

func WriteFileOrExit(path string, data string) {
	err := os.WriteFile(path, []byte(data), 0644)
	if err != nil {
		logger.Errorf("Failed to write to a file\n%v\n", err)
		os.Exit(1)
	}
}

func WriteFileWithSudoOrExit(path string, data string) {
	if _, err := exec.Command("sudo", "touch", path).Output(); err != nil {
		logger.Errorf("Failed to create file %s\n%v\n", path, err)
		os.Exit(1)
	}

	usr, _ := user.Current()
	if _, err := exec.Command("sudo", "chown", usr.Username, path).Output(); err != nil {
		logger.Errorf("Failed to call `chown` on file %s\n%v\n", path, err)
		os.Exit(1)
	}

	err := os.WriteFile(path, []byte(data), 0644)

	if err != nil {
		logger.Errorf("Failed to write to a file\n%v\n", err)
		os.Exit(1)
	}
}

func MakeDirOrExit(path string) {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		logger.Errorf("Failed to create directory %s\n%v\n", path, err)
		os.Exit(1)
	}
}

func MakeDirWithSudoOrExit(path string) {
	if _, err := exec.Command("sudo", "mkdir", "-p", path).Output(); err != nil {
		logger.Errorf("Failed to create directory $s\n%v\n", path, err)
		os.Exit(1)
	}
}

func GetCurrentDir() string {
	cwd, err := os.Getwd()
	if err != nil {
		logger.Errorf("Failed to get current working directory\n%v\n", err)
		os.Exit(1)
	}

	return cwd
}

func FileExists(path string) bool {
	fInfo, _ := os.Stat(path)

	return fInfo != nil
}
