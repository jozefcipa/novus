package sudoers

import (
	"fmt"
	"path/filepath"

	"github.com/jozefcipa/novus/internal/brew"
	"github.com/jozefcipa/novus/internal/fs"
	"github.com/jozefcipa/novus/internal/logger"
)

const sudoersFile = "/etc/sudoers.d/novus"

func Trust() {
	// Register Homebrew to sudoers file so it can be ran without sudo password
	brewBinPath := filepath.Join(brew.BrewPath, "bin/brew")
	novusBinPath := filepath.Join(fs.NovusBinaryDir, "novus")
	sudoPermissions := fmt.Sprintf(`Cmnd_Alias HOMEBREW = %s *
Cmnd_Alias NOVUS = %s *
%%admin ALL=(root) NOPASSWD: HOMEBREW, NOVUS
`,
		brewBinPath,
		novusBinPath,
	)

	// TODO: this only works when running it as `sudo novus ...` -> how to run it without sudo?

	logger.Debugf("Creating /etc/sudoers.d/novus file:\n\n%s", sudoPermissions)
	fs.WriteFileWithSudoOrExit(sudoersFile, sudoPermissions)

	// Sudoers file must be owned by root
	fs.ChownOrExit(sudoersFile, "root")
}

func IsTrusted() bool {
	logger.Debugf("Checking if %s exists", sudoersFile)
	return fs.FileExists(sudoersFile)
}
