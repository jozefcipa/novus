package cmd

import (
	"os"

	"github.com/jozefcipa/novus/internal/fs"
	"github.com/jozefcipa/novus/internal/logger"
	"github.com/jozefcipa/novus/internal/paths"
	"github.com/jozefcipa/novus/internal/sudo"
	"github.com/jozefcipa/novus/internal/tui"
	"github.com/spf13/cobra"
)

var revokeTrustFlag bool

var trustCmd = &cobra.Command{
	Use:   "trust",
	Short: "Enable passwordless use of Novus",
	Long:  "Create a sudoers file so Novus can be run without prompting for `sudo` password",
	Run: func(cmd *cobra.Command, args []string) {
		sudoersExists := fs.FileExists(paths.SudoersFilePath)

		// Revoke trust
		if revokeTrustFlag {
			if sudoersExists {
				if yes := tui.Confirm("Do you really want to revoke the sudo trust?"); !yes {
					os.Exit(0)
				}

				err := sudo.DeleteFile(paths.SudoersFilePath)
				if err != nil {
					logger.Errorf(err.Error())
					os.Exit(1)
				}

				logger.Infof("🚫 Novus trust revoked")
			} else {
				logger.Hintf("Novus trust is already revoked.")
			}

			os.Exit(0)
		}

		// If the sudoers exists, just exit
		if sudoersExists {
			logger.Checkf("Novus is already trusted and can be used without sudo password")
			os.Exit(0)
		}

		// Otherwise create a sudoers file
		logger.Infof("Creating sudoers file")
		sudo.RegisterSudoersFile()
		logger.Checkf("Novus is now trusted and can be used without sudo password.")
	},
}

func init() {
	trustCmd.Flags().BoolVar(&revokeTrustFlag, "revoke", false, "Remove the sudoers files and start prompting for sudo password again")
	rootCmd.AddCommand(trustCmd)
}
