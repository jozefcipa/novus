package cmd

import (
	"errors"
	"os"

	"github.com/jozefcipa/novus/internal/brew"
	"github.com/jozefcipa/novus/internal/config"
	"github.com/jozefcipa/novus/internal/logger"
	"github.com/jozefcipa/novus/internal/shared"
	"github.com/jozefcipa/novus/internal/tui"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Novus configuration",
	Long:  "Initialize Novus configuration by creating the " + config.AppName + " file and installs all required binaries if not installed yet.",
	Run: func(cmd *cobra.Command, args []string) {
		// Install nginx, dnsmasq and mkcert if not installed
		if err := brew.InstallBinaries(); err != nil {
			if errors.Is(err, &brew.BrewMissingError{}) {
				logger.Errorf(err.Error())
				logger.Hintf("You can install it from \033[4mhttps://brew.sh/\033[0m")
			}
			os.Exit(1)
		}

		// Check if novus.yml config exists
		_, exists := config.Load()
		if !exists {
			// If config doesn't exist, create a new one
			input := tui.AskUser("Enter a new app name: ")
			appName := shared.ToKebabCase(input)

			err := config.CreateDefaultConfigFile(appName)
			if err != nil {
				logger.Errorf("%s", err.Error())
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
