package cmd

import (
	"github.com/jozefcipa/novus/internal/dnsmasq"
	"github.com/jozefcipa/novus/internal/logger"
	"github.com/jozefcipa/novus/internal/nginx"

	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop Nginx and DNSMasq services",
	Long: `Running this command will stop the HTTP and DNS servers,
so Novus will no longer serve application requests to the URLs
defined in the novus.yml configuration file.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		nginx.Stop()
		logger.Messagef("🚫 Nginx stopped.\n")

		dnsmasq.Stop()
		logger.Messagef("🚫 DNSMasq stopped.\n")
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
