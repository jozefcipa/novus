package cmd

import (
	"github.com/jozefcipa/novus/internal/brew"
	"github.com/jozefcipa/novus/internal/logger"

	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		brew.StopBrewService("nginx")
		logger.Messagef("🚫 Nginx stopped.\n")

		brew.StopBrewService("dnsmasq")
		logger.Messagef("🚫 DNSMasq stopped.\n")
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
