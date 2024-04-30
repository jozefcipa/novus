package cmd

import (
	"os"

	"github.com/jozefcipa/novus/internal/brew"
	"github.com/jozefcipa/novus/internal/cli"
	"github.com/jozefcipa/novus/internal/config"
	"github.com/jozefcipa/novus/internal/logger"
	"github.com/jozefcipa/novus/internal/shared"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Novus configuration",
	Long:  "Initialize Novus configuration by creating the novus.yml file and installs all required binaries if not installed yet.",
	Run: func(cmd *cobra.Command, args []string) {
		// Install nginx, dnsmasq and mkcert if not installed
		brew.InstallBinaries()

		// Check if novus.yml config exists
		_, exists := config.Load()
		if !exists {
			// If config doesn't exist, create a new one
			input := cli.AskUser("Enter a new app name: ")
			appName := shared.ToKebabCase(input)

			err := config.CreateDefaultConfigFile(appName)
			if err != nil {
				logger.Errorf("%s\n", err.Error())
				os.Exit(1)
			}
			logger.Successf("✅ Novus has been initialized.\n")
			logger.Messagef("💡 Open \"novus.yml\" to add your route definitions.\n")
		} else {
			logger.Checkf(" Novus is already initialized.")
			logger.Messagef("💡 Run \"novus serve\" to start the proxy.\n")
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
