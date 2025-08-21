package paths

import (
	"strings"

	"github.com/jozefcipa/novus/internal/logger"
)

// Sudoers file that contains configuration for Novus
const SudoersFilePath = "/etc/sudoers.d/novus"

// A shell script that directly calls `sudo` commands
// This folder will be created and owned by root so it's not modifiable by the user
// and will also keep working between different Novus versions
const SudoHelperPath = "/usr/local/novus/sudo-helper"

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

	logger.Debugf(
		"Sudo paths resolved.\n"+
			"\tSudoHelperPath = %s\n"+
			"\tSudoAllowedPaths = [%s]",
		SudoHelperPath,
		strings.Join(SudoAllowedPaths, ", "),
	)
}
