package cmd

import (
	"github.com/jozefcipa/novus/internal/config"
	"github.com/jozefcipa/novus/internal/dnsmasq"
	"github.com/jozefcipa/novus/internal/nginx"

	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop Nginx and DNSMasq services",
	Long: `Running this command will stop the HTTP and DNS servers,
so Novus will no longer serve application requests to the URLs
defined in the ` + config.ConfigFileName + ` configuration file.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		nginx.Stop()
		dnsmasq.Stop()
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
