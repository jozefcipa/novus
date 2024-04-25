package fs

import (
	"io"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/jozefcipa/novus/internal/logger"
)

var UserHomeDir string
var CurrentDir string
var AssetsDir string

// Cannot use `init()` here because the order in which these init() functions are called across packages causes
// that `DebugEnabled` flag is not yet available here (`rootCmd.init()` is called after `fs.init()`)
func ResolveDirs() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		logger.Errorf("Failed to get user home directory\n%v\n", err)
		os.Exit(1)
	}
	UserHomeDir = homeDir

	currentDir, err := os.Getwd()
	if err != nil {
		logger.Errorf("Failed to get current working directory\n%v\n", err)
		os.Exit(1)
	}
	CurrentDir = currentDir

	executablePath, err := os.Executable()
	if err != nil {
		logger.Errorf("Failed to get novus binary directory\n%v\n", err)
		os.Exit(1)
	}
	novusBinaryDir := filepath.Dir(executablePath)
	// When running in development with `go run` it gives temporary directory,
	// therefore set the novus dir path to the current directory
	if strings.Contains(novusBinaryDir, "go-build") {
		novusBinaryDir = currentDir
		// In local develop environment, the ./assets are stored next to the output binary
		// - ./novus
		// - ./assets/...
		AssetsDir = filepath.Join(currentDir, "assets")
	} else {
		// If a non-develop binary is used, the ./assets directory is one level above in the filesystem
		// This is the Homebrew structure
		// - ./bin/novus
		// - ./assets/...
		AssetsDir = filepath.Join(novusBinaryDir, "..", "assets")
	}

	// Make sure assets directory is available
	if !FileExists(AssetsDir) {
		logger.Errorf("Assets directory not found: %s\n", AssetsDir)
		os.Exit(1)
	}

	logger.Debugf(
		"Filesystem initialized.\n"+
			"\tUser Home Directory = %s\n"+
			"\tAssets Directory = %s\n"+
			"\tNovus Binary Directory = %s\n"+
			"\tCurrent Directory = %s",
		UserHomeDir,
		AssetsDir,
		novusBinaryDir,
		CurrentDir,
	)
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

	// We need to change the file owner to the current user in order to be able to write to the file
	user, _ := user.Current()
	ChownOrExit(path, user.Username)
	err := os.WriteFile(path, []byte(data), 0644)

	if err != nil {
		logger.Errorf("Failed to write to a file\n%v\n", err)
		os.Exit(1)
	}
}

func DeleteFile(path string) error {
	if err := os.Remove(path); err != nil {
		logger.Errorf("Failed to delete file %s\n%v\n", path, err)
		return err
	}
	return nil
}

func DeleteFileWithSudo(path string) error {
	if _, err := exec.Command("sudo", "rm", path).Output(); err != nil {
		logger.Errorf("Failed to delete file %s\n%v\n", path, err)
		return err
	}

	return nil
}

func MakeDirOrExit(path string) {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		logger.Errorf("Failed to create directory %s\n%v\n", path, err)
		os.Exit(1)
	}
}

func DeleteDir(path string) error {
	if err := os.RemoveAll(path); err != nil {
		logger.Errorf("Failed to delete directory %s\n%v\n", path, err)
		return err
	}
	return nil
}

func MakeDirWithSudoOrExit(path string) {
	if _, err := exec.Command("sudo", "mkdir", "-p", path).Output(); err != nil {
		logger.Errorf("Failed to create directory $s\n%v\n", path, err)
		os.Exit(1)
	}
}

func FileExists(path string) bool {
	fInfo, _ := os.Stat(path)

	return fInfo != nil
}

func Copy(src string, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}

func ChownOrExit(path string, user string) {
	if _, err := exec.Command("sudo", "chown", user, path).Output(); err != nil {
		logger.Errorf("Failed to call `chown` on file %s\n%v\n", path, err)
		os.Exit(1)
	}
}
