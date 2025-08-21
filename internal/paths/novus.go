package paths

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/jozefcipa/novus/internal/homebrew"
	"github.com/jozefcipa/novus/internal/logger"
)

// User home directory, in which we store the Novus state (~/)
var UserHomeDir string

// Holds current directory path from which the binary is executed
var CurrentDir string

// Path to the assets directory
var AssetsDir string

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

	// Novus executable
	executablePath, err := os.Executable()
	if err != nil {
		logger.Errorf("Failed to get novus binary directory\n   Reason: %v", err)
		os.Exit(1)
	}

	// Assets dir is relative to the executable path
	if strings.Contains(executablePath, "go-build") {
		// When running in development with `go run` it gives temporary directory,
		// therefore set the novus dir path to the current directory
		// .
		// ├── assets/
		AssetsDir = filepath.Join(currentDir, "assets")
	} else if strings.Contains(executablePath, homebrew.HomebrewPrefix) {
		// If running via Homebrew, the binary is in the Homebrew prefix directory
		// .
		// ├── {homebrew.HomebrewPrefix}/opt/
		// │   └── novus/
		// │       ├── bin/
		// │       │   └── novus
		// │       └── assets/
		AssetsDir = filepath.Join(homebrew.HomebrewPrefix, "/opt/novus/assets")
	} else {
		// Otherwise if built locally via `make build`, the binary is in the `bin` directory
		// .
		// ├── bin/
		// │   └── novus
		// ├── assets/
		novusExecutableDir := filepath.Dir(executablePath)
		AssetsDir = filepath.Join(novusExecutableDir, "../assets")
	}

	CurrentDir = currentDir
	NovusStateDir = filepath.Join(UserHomeDir, ".novus")
	NovusStateFilePath = filepath.Join(NovusStateDir, "novus.json")

	logger.Debugf(
		"Novus paths resolved.\n"+
			"\tUserHomeDir = %s\n"+
			"\tCurrentDir = %s\n"+
			"\tNovusStateDir = %s\n"+
			"\tNovusStateFilePath = %s\n"+
			"\tAssetsDir = %s",
		UserHomeDir,
		CurrentDir,
		NovusStateDir,
		NovusStateFilePath,
		AssetsDir,
	)
}
