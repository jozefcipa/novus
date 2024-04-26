package novus

import (
	"fmt"
	"path/filepath"

	"github.com/jozefcipa/novus/internal/brew"
	"github.com/jozefcipa/novus/internal/fs"
	"github.com/jozefcipa/novus/internal/logger"
)

const sudoersFile = "/etc/sudoers.d/novus"

func CreateSudoersFile() {
	// Register Homebrew to sudoers file so it can be ran without sudo password
	brewBinPath := filepath.Join(brew.BrewPath, "bin/brew")
	sudoPermissions := fmt.Sprintf("Cmnd_Alias HOMEBREW = %s *\n"+
		"%%admin ALL=(root) NOPASSWD: HOMEBREW\n",
		brewBinPath,
	)

	logger.Messagef("⏳ Creating /etc/sudoers.d file for Novus.\n")
	fs.WriteFileWithSudoOrExit(sudoersFile, sudoPermissions)

	// sudoers file must be owned by root
	fs.ChownOrExit(sudoersFile, "root")
}

func SudoersFileExists() bool {
	return fs.FileExists(sudoersFile)
}
