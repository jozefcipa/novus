package cmd

import (
	"github.com/jozefcipa/novus/internal/dnsmasq"
	"github.com/jozefcipa/novus/internal/logger"
	"github.com/jozefcipa/novus/internal/nginx"
	"github.com/jozefcipa/novus/internal/novus"
	"github.com/jozefcipa/novus/internal/tui"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show status of services and registered routes",
	Long: `Show whether Nginx and DNSMasq services are running,
and print a list of all URLs that are registered by Novus.`,
	Run: func(cmd *cobra.Command, args []string) {
		isNginxRunning := nginx.IsRunning()
		isDNSMasqRunning := dnsmasq.IsRunning()

		if isNginxRunning {
			logger.Successf("Nginx running 🚀")
			logger.Debugf("Nginx configuration loaded from %s", nginx.NginxServersDir)
		} else {
			logger.Errorf("Nginx not running")
		}

		if isDNSMasqRunning {
			logger.Successf("DNSMasq running 🚀")
		} else {
			logger.Errorf("DNSMasq not running")
		}

		if !isNginxRunning || !isDNSMasqRunning {
			logger.Hintf("Run \"novus serve\" to initialize the services")
		} else {
			// All good, show the routing info
			novusState := novus.GetState()
			tui.PrintRoutingTable(*novusState)
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
