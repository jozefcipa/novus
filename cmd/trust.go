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
which will allow running Novus commands without asking for password.

Novus needs sudo access for manipulating DNS records via DNSMasq.
`,
	Run: func(cmd *cobra.Command, args []string) {
		sudoersFile := "/etc/sudoers.d/novus"

		if fs.FileExists(sudoersFile) {
			logger.Checkf("Novus is already registered in %s.", sudoersFile)
			os.Exit(0)
		}

		novusBinPath := filepath.Join(brew.BrewPath, "bin/novus")
		sudoPermissions := fmt.Sprintf("Cmnd_Alias NOVUS = %s *\n%%admin ALL=(root) NOPASSWD: NOVUS\n", novusBinPath)
		logger.Messagef("Creating /etc/sudoers.d file for Novus.\n")
		fs.WriteFileWithSudoOrExit(sudoersFile, sudoPermissions)
		logger.Successf("Novus is now registered in %s and can be used without sudo password.\n", sudoersFile)
	},
}

func init() {
	rootCmd.AddCommand(trustCmd)
}
