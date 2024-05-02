package cmd

import (
	"os"

	"github.com/jozefcipa/novus/internal/logger"
	"github.com/jozefcipa/novus/internal/novus"
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
		// Check if sudoers file already exists
		if exists := novus.SudoersFileExists(); exists {
			logger.Checkf("Novus is already trusted.")
			os.Exit(0)
		}

		// Create file if it doesn't exist yet
		novus.CreateSudoersFile()
		logger.Successf("Novus is now registered and can be used without sudo password.")
	},
}

func init() {
	rootCmd.AddCommand(trustCmd)
}
