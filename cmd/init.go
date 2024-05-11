package cmd

import (
	"os"

	"github.com/jozefcipa/novus/internal/brew"
	"github.com/jozefcipa/novus/internal/config"
	"github.com/jozefcipa/novus/internal/config_manager"
	"github.com/jozefcipa/novus/internal/logger"
	"github.com/jozefcipa/novus/internal/novus"
	"github.com/jozefcipa/novus/internal/shared"
	"github.com/jozefcipa/novus/internal/tui"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Novus configuration",
	Long:  "Initialize Novus configuration by creating the " + config.ConfigFileName + " file and installs all required binaries if not installed yet.",
	Run: func(cmd *cobra.Command, args []string) {
		// Install nginx, dnsmasq and mkcert if not installed
		if err := brew.InstallBinaries(); err != nil {
			logger.Errorf(err.Error())

			if _, ok := err.(*brew.BrewMissingError); ok {
				logger.Hintf("You can install it from \033[4mhttps://brew.sh/\033[0m")
			}
			os.Exit(1)
		}

		// Check if novus.yml config exists
		exists := config_manager.ConfigFileExists()
		if !exists {
			// If config doesn't exist, create a new one
			input := tui.AskUser("Enter a new app name: ")
			appName := shared.ToKebabCase(input)

			err := config_manager.CreateNewConfiguration(appName, *novus.GetState())
			if err != nil {
				logger.Errorf(err.Error())
				os.Exit(1)
			}
			logger.Successf("Novus has been initialized.")
			logger.Hintf("Open " + config.ConfigFileName + " to add your route definitions.")
		} else {
			logger.Checkf("Novus is already initialized.")
			logger.Hintf("Run \"novus serve\" to start the proxy.")
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
