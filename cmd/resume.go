package cmd

import (
	"os"

	"github.com/jozefcipa/novus/internal/config_manager"
	"github.com/jozefcipa/novus/internal/dns_manager"
	"github.com/jozefcipa/novus/internal/dnsmasq"
	"github.com/jozefcipa/novus/internal/logger"
	"github.com/jozefcipa/novus/internal/mkcert"
	"github.com/jozefcipa/novus/internal/nginx"
	"github.com/jozefcipa/novus/internal/novus"
	"github.com/jozefcipa/novus/internal/ssl_manager"
	"github.com/jozefcipa/novus/internal/tui"
	"github.com/spf13/cobra"
)

var resumeCmd = &cobra.Command{
	Use:   "resume [app-name]",
	Short: "Resume paused app in Novus",
	Long:  "Resume paused app in Novus so the routing will start again. Similar to running `novus serve` but the app has to be registered already",
	Run: func(cmd *cobra.Command, args []string) {
		appName, appState := tui.ParseAppFromArgs(args, "resume")

		if appState.Status == novus.APP_ACTIVE {
			logger.Checkf("\"%s\" is already active.", appName)
			tui.PrintRoutingTable(*novus.GetState())
			os.Exit(0)
		}

		novusState := novus.GetState()

		// Load config from state
		conf := config_manager.LoadConfigurationFromState(appName, *novusState)
		config_manager.ValidateConfig(conf, *novusState)

		// Configure SSL
		mkcert.Configure(conf)
		domainCerts, _ := ssl_manager.EnsureSSLCertificates(conf, novusState, appName)

		// Configure Nginx
		nginx.Configure(conf, domainCerts, appState)

		// Configure DNS
		dns_manager.Configure(conf, novusState)

		// Restart services
		nginx.Restart()
		logger.Checkf("Nginx restarted 🔄")
		dnsmasq.Restart()
		logger.Checkf("DNSMasq restarted 🔄")

		// If app has been paused, make sure to set it to ACTIVE
		appState.Status = novus.APP_ACTIVE

		// Everything's set, start routing
		tui.PrintRoutingTable(*novusState)

		// Save application state
		novus.SaveState()
	},
}

func init() {
	rootCmd.AddCommand(resumeCmd)
}
