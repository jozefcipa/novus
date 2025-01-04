package paths

import (
	"path/filepath"
	"strings"

	"github.com/jozefcipa/novus/internal/logger"
)

// Sudoers file that contains configuration for Novus
const SudoersFilePath = "/etc/sudoers.d/novus"

// A shell script that directly calls `sudo` commands
var SudoHelperPath string

// This is an array of paths that the above-mentioned shell script can modify
// As the script has root access (passwordless sudo), this config is a security measure
// to limit the scope as much as possible
var SudoAllowedPaths []string

func resolveSudoDirs() {
	logger.Debugf("Resolving sudo paths")

	SudoAllowedPaths = []string{
		DNSResolverDir,
		SudoersFilePath,
	}

	SudoHelperPath = filepath.Join(NovusBinaryDir, "sudo-helper")

	logger.Debugf(
		"Sudo paths resolved.\n"+
			"\tSudoHelperPath = %s\n"+
			"\tSudoAllowedPaths = [%s]",
		SudoHelperPath,
		strings.Join(SudoAllowedPaths, ", "),
	)
}
