package cmd

import (
	"fmt"
	"os"

	"github.com/jozefcipa/novus/internal/config"
	"github.com/jozefcipa/novus/internal/config_manager"
	"github.com/jozefcipa/novus/internal/homebrew"
	"github.com/jozefcipa/novus/internal/logger"
	"github.com/jozefcipa/novus/internal/novus"
	"github.com/jozefcipa/novus/internal/stringutils"
	"github.com/jozefcipa/novus/internal/tui"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Novus configuration",
	Long:  fmt.Sprintf("Initialize Novus configuration by creating the %s file and installs all required binaries if not installed yet.", config.ConfigFileName),
	Run: func(cmd *cobra.Command, args []string) {
		// Install nginx, dnsmasq and mkcert if not installed
		if err := homebrew.InstallBinaries(); err != nil {
			logger.Errorf(err.Error())

			if _, ok := err.(*homebrew.HomebrewMissingError); ok {
				logger.Hintf("You can install it from %shttps://brew.sh/%s", logger.UNDERLINE, logger.RESET)
			}
			os.Exit(1)
		}

		// Check if novus.yml config exists
		exists := config_manager.ConfigFileExists()
		if !exists {
			// If config doesn't exist, create a new one
			input := tui.AskUser("Enter a new app name: ")
			appName := stringutils.ToKebabCase(input)

			err := config_manager.CreateNewConfiguration(appName, *novus.GetState())
			if err != nil {
				logger.Errorf(err.Error())
				os.Exit(1)
			}
			logger.Successf("Novus has been initialized.")
			logger.Hintf("Open %s to add your route definitions.", config.ConfigFileName)
		} else {
			logger.Checkf("Novus is already initialized (%s file exists).", config.ConfigFileName)
			logger.Hintf("Run \"novus serve\" to start routing.")
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
