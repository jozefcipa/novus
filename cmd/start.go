package cmd

import (
	"github.com/jozefcipa/novus/internal/brew"
	"github.com/jozefcipa/novus/internal/config"
	"github.com/jozefcipa/novus/internal/nginx"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Make sure we have the necessary binaries available
		brew.InstallBinaries()

		// Load configuration file
		conf := config.Load()

		// Start services
		// nginx.Start()
		// TODO dnsmasq ...

		nginx.Configure(conf)

	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
