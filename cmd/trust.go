package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jozefcipa/novus/internal/brew"
	"github.com/jozefcipa/novus/internal/fs"
	"github.com/jozefcipa/novus/internal/logger"
	"github.com/spf13/cobra"
)

var trustCmd = &cobra.Command{
	Use:   "trust",
	Short: "Allow running `novus` commands without password",
	Long: `This command will create a record in "/etc/sudoers.d/novus",
which allows running Novus commands without asking for a password.

Novus needs sudo access for manipulating DNS records via DNSMasq.
`,
	Run: func(cmd *cobra.Command, args []string) {
		sudoersFile := "/etc/sudoers.d/novus"

		if fs.FileExists(sudoersFile) {
			logger.Checkf("Novus is already registered in %s.", sudoersFile)
			os.Exit(0)
		}

		// Register Homebrew to sudoers file so it can be ran without sudo password
		brewBinPath := filepath.Join(brew.BrewPath, "bin/brew")
		sudoPermissions := fmt.Sprintf("Cmnd_Alias HOMEBREW = %s *\n"+
			"%%admin ALL=(root) NOPASSWD: HOMEBREW\n",
			brewBinPath,
		)
		logger.Messagef("Creating /etc/sudoers.d file for Novus.\n")

		fs.WriteFileWithSudoOrExit(sudoersFile, sudoPermissions)
		fs.ChownOrExit(sudoersFile, "root") // sudoers file must be owned by root

		logger.Successf("Novus is now registered in %s and can be used without sudo password.\n", sudoersFile)
	},
}

func init() {
	rootCmd.AddCommand(trustCmd)
}
