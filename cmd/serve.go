package cmd

import (
	"os"

	"github.com/jozefcipa/novus/internal/brew"
	"github.com/jozefcipa/novus/internal/config"
	"github.com/jozefcipa/novus/internal/diff_manager"
	"github.com/jozefcipa/novus/internal/dnsmasq"
	"github.com/jozefcipa/novus/internal/logger"
	"github.com/jozefcipa/novus/internal/mkcert"
	"github.com/jozefcipa/novus/internal/nginx"
	"github.com/jozefcipa/novus/internal/novus"
	"github.com/jozefcipa/novus/internal/ssl_manager"
	"github.com/jozefcipa/novus/internal/tui"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Setup Novus to start serving URLs",
	Long:  `Install Nginx, DNSMasq and mkcert and automatically expose HTTPs URLs for the endpoints defined in the config.`,
	Run: func(cmd *cobra.Command, args []string) {
		// if we don't have the binaries, exit here, the user needs to run `novus init` first.
		if !brew.CheckRequiredBinariesPresence() {
			os.Exit(1)
		}

		conf, exists := config.Load() // Load configuration file
		if !exists {
			logger.Warnf("🙉 Novus is not initialized in this directory (no configuration found).\n")
			logger.Messagef("💡 Run \"novus init\" to create a configuration file.\n")
			os.Exit(1)
		}
		appState, isNewState := novus.GetAppState() // Load application state

		// Handle config changes diff
		addedRoutes, deletedRoutes := diff_manager.DetectConfigDiff(conf, *appState)

		// Remove domains that are no longer in config
		if len(deletedRoutes) > 0 {
			for _, deletedRoute := range deletedRoutes {
				logger.Errorf("❌ Removing SSL certificate for domain [%s]\n", deletedRoute.Route.Domain)
				ssl_manager.DeleteCert(deletedRoute.Route.Domain)
			}

			// Remove DNS records for unused TLDs
			unusedTLDs := diff_manager.DetectUnusedTLDs(conf, *appState)
			if len(unusedTLDs) > 0 {
				for _, tld := range unusedTLDs {
					logger.Errorf("❌ Removing unused TLD domain [*.%s]\n", tld)
					dnsmasq.UnregisterTLD(tld)
				}
			}
		}

		if !isNewState && len(addedRoutes) > 0 {
			for _, newRoute := range addedRoutes {
				logger.Successf("Found new domain [%s]\n", newRoute.Route.Domain)
			}
		}

		// Configure SSL
		mkcert.Configure(conf)
		domainCerts, hasNewCerts := ssl_manager.EnsureSSLCertificates(conf)

		// Configure Nginx
		nginxConfigUpdated := nginx.Configure(conf, domainCerts)

		// Configure DNSMasq
		dnsMasqConfigUpdated := dnsmasq.Configure(conf)

		// Restart services
		// Nginx
		isNginxRunning := nginx.IsRunning()
		if nginxConfigUpdated || hasNewCerts || !isNginxRunning {
			nginx.Restart()
			logger.Checkf("Nginx restarted")
		} else {
			logger.Checkf("Nginx is up and running")
		}

		// DNSMasq
		isDNSMasqRunning := dnsmasq.IsRunning()
		if dnsMasqConfigUpdated || !isDNSMasqRunning {
			dnsmasq.Restart()
			logger.Checkf("DNSMasq restarted")
		} else {
			logger.Checkf("DNSMasq is up and running")
		}

		// Everything's set, start routing
		logger.Checkf("Routing has started 🚀")
		tui.PrintRoutingTable(novus.GetState().Apps)

		// Save application state
		novus.SaveState()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
