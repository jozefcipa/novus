package paths

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/jozefcipa/novus/internal/logger"
)

// User home directory, in which we store the Novus state (~/)
var UserHomeDir string

// Holds current directory path from which the binary is executed
var CurrentDir string

// Path to the assets directory
var AssetsDir string

// Path to the Novus executable binary directory
var NovusBinaryDir string

// Main directory for storing all Novus application data, e.g. SSL certificates, configuration state, etc. (~/.novus)
var NovusStateDir string

// Configuration state file
var NovusStateFilePath string

func resolveNovusDirs() {
	// Home dir
	homeDir, err := os.UserHomeDir()
	if err != nil {
		logger.Errorf("Failed to get user home directory%s\n   Reason: %v", err)
		os.Exit(1)
	}
	UserHomeDir = homeDir

	// Current dir
	currentDir, err := os.Getwd()
	if err != nil {
		logger.Errorf("Failed to get current working directory%s\n   Reason: %v", err)
		os.Exit(1)
	}

	// Novus binary dir
	executablePath, err := os.Executable()
	if err != nil {
		logger.Errorf("Failed to get novus binary directory\n   Reason: %v", err)
		os.Exit(1)
	}
	// Homebrew creates symlinks for binaries, we need to get the original binary location
	executablePath, err = filepath.EvalSymlinks(executablePath)
	if err != nil {
		logger.Errorf("Failed to evaluate novus symlink\n   Reason: %v", err)
		os.Exit(1)
	}
	NovusBinaryDir = filepath.Dir(executablePath)

	// Assets dir
	// When running in development with `go run` it gives temporary directory,
	// therefore set the novus dir path to the current directory
	if strings.Contains(NovusBinaryDir, "go-build") {
		NovusBinaryDir = currentDir
		// In local develop environment, the ./assets are stored next to the output binary
		// - ./novus
		// - ./assets/...
		AssetsDir = filepath.Join(currentDir, "assets")
	} else {
		// If a non-develop binary is used, the ./assets directory is one level above in the filesystem
		// This is the Homebrew structure
		// - ./bin/novus
		// - ./assets/...
		AssetsDir = filepath.Join(NovusBinaryDir, "..", "assets")
	}

	CurrentDir = currentDir
	NovusStateDir = filepath.Join(UserHomeDir, ".novus")
	NovusStateFilePath = filepath.Join(NovusStateDir, "novus.json")

	logger.Debugf(
		"Novus paths resolved.\n"+
			"\tUserHomeDir = %s\n"+
			"\tUserHomeDir = %s\n"+
			"\tNovusBinaryDir = %s\n"+
			"\tNNovusStateDir = %s\n"+
			"\tNovusStateFilePath = %s\n"+
			"\tAssetsDir = %s",
		UserHomeDir,
		CurrentDir,
		NovusBinaryDir,
		NovusStateDir,
		NovusStateFilePath,
		AssetsDir,
	)
}
