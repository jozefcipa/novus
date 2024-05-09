package cmd

import (
	"os"

	"github.com/jozefcipa/novus/internal/brew"
	"github.com/jozefcipa/novus/internal/config"
	"github.com/jozefcipa/novus/internal/config_manager"
	"github.com/jozefcipa/novus/internal/diff_manager"
	"github.com/jozefcipa/novus/internal/dns_manager"
	"github.com/jozefcipa/novus/internal/dnsmasq"
	"github.com/jozefcipa/novus/internal/logger"
	"github.com/jozefcipa/novus/internal/mkcert"
	"github.com/jozefcipa/novus/internal/nginx"
	"github.com/jozefcipa/novus/internal/novus"
	"github.com/jozefcipa/novus/internal/shared"
	"github.com/jozefcipa/novus/internal/ssl_manager"
	"github.com/jozefcipa/novus/internal/tui"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Setup Novus to start serving URLs",
	Long:  `Install Nginx, DNSMasq and mkcert and automatically expose HTTPs URLs for the endpoints defined in the config.`,
	Run: func(cmd *cobra.Command, args []string) {
		// If the binaries are missing, exit here, user needs to run `novus init` first
		if err := brew.CheckIfRequiredBinariesInstalled(); err != nil {
			logger.Hintf("Run \"novus init\" first to initialize Novus.")
			os.Exit(1)
		}

		// Load configuration file
		conf, exists := config_manager.LoadConfiguration()
		if !exists {
			logger.Warnf("Novus is not initialized in this directory (" + config.ConfigFileName + " file does not exist).")
			logger.Hintf("Run \"novus init\" to create a configuration file.")
			os.Exit(1)
		}

		// Load application state
		appState, isNewState := novus.GetAppState(config.AppName())

		// Compare state and current config to detect changes
		addedRoutes, deletedRoutes := diff_manager.DetectConfigDiff(conf, *appState)

		// Remove domains that are no longer in config
		if len(deletedRoutes) > 0 {
			for _, deletedRoute := range deletedRoutes {
				logger.Errorf("Removing SSL certificate for domain [%s]", deletedRoute.Domain)
				ssl_manager.DeleteCert(deletedRoute.Domain)
			}

			// Remove DNS records for unused TLDs
			otherAppsRoutes := []shared.Route{}
			for appName, appState := range novus.GetState().Apps {
				// we want to find usage in other apps, not the current one
				if appName != config.AppName() {
					otherAppsRoutes = append(otherAppsRoutes, appState.Routes...)
				}
			}
			unusedTLDs := diff_manager.DetectUnusedTLDs(deletedRoutes, otherAppsRoutes)
			if len(unusedTLDs) > 0 {
				for _, tld := range unusedTLDs {
					logger.Errorf("Removing unused TLD domain [*.%s]", tld)
					dns_manager.UnregisterTLD(tld)
				}
			}
		}

		if !isNewState && len(addedRoutes) > 0 {
			for _, newRoute := range addedRoutes {
				logger.Successf("Found new domain [%s]", newRoute.Domain)
			}
		}

		// Configure SSL
		mkcert.Configure(conf)
		domainCerts, hasNewCerts := ssl_manager.EnsureSSLCertificates(conf)

		// Configure Nginx
		nginxConfigUpdated := nginx.Configure(conf, domainCerts)

		// Configure DNS
		dnsUpdated := dns_manager.Configure(conf)

		// Restart services
		// Nginx
		isNginxRunning := nginx.IsRunning()
		if nginxConfigUpdated || hasNewCerts || !isNginxRunning {
			nginx.Restart()
			logger.Checkf("Nginx restarted 🔄")
		} else {
			logger.Checkf("Nginx running 🚀")
		}

		// DNSMasq
		isDNSMasqRunning := dnsmasq.IsRunning()
		if dnsUpdated || !isDNSMasqRunning {
			dnsmasq.Restart()
			logger.Checkf("DNSMasq restarted 🔄")
		} else {
			logger.Checkf("DNSMasq running 🚀")
		}

		// Everything's set, start routing
		tui.PrintRoutingTable(novus.GetState().Apps)

		// Save application state
		novus.SaveState()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
