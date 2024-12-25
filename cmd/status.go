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
		nginxLoader := logger.Loadingf("Checking Nginx status")
		isNginxRunning := nginx.IsRunning()
		if isNginxRunning {
			nginxLoader.Checkf("Nginx running")
			logger.Debugf("Nginx configuration loaded from %s", nginx.NginxServersDir)
		} else {
			nginxLoader.Errorf("Nginx not running")
		}

		dnsmasqLoader := logger.Loadingf("Checking DNSMasq status")
		isDNSMasqRunning := dnsmasq.IsRunning()
		if isDNSMasqRunning {
			dnsmasqLoader.Checkf("DNSMasq running")
		} else {
			dnsmasqLoader.Errorf("DNSMasq not running")
		}

		if !isNginxRunning || !isDNSMasqRunning {
			logger.Hintf("Run \"novus start\" to start routing.")
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
