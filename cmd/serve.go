package cmd

import (
	"os"
	"slices"

	"github.com/jozefcipa/novus/internal/brew"
	"github.com/jozefcipa/novus/internal/config"
	"github.com/jozefcipa/novus/internal/diff_manager"
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
			logger.Warnf("🙉 Novus is not initialized in this directory (no configuration found).")
			logger.Hintf("Run \"novus init\" to create a configuration file.")
			os.Exit(1)
		}
		appState, isNewState := novus.GetAppState(config.AppName) // Load application state

		// Handle config changes diff
		addedRoutes, deletedRoutes := diff_manager.DetectConfigDiff(conf, *appState)

		// Remove domains that are no longer in config
		if len(deletedRoutes) > 0 {
			for _, deletedRoute := range deletedRoutes {
				logger.Errorf("Removing SSL certificate for domain [%s]", deletedRoute.Domain)
				ssl_manager.DeleteCert(deletedRoute.Domain)
			}

			// Remove DNS records for unused TLDs
			unusedTLDs := diff_manager.DetectUnusedTLDs(conf, *appState)
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

		// TODO: refactor this into a function and move to some module
		// Check if there are duplicate domains across apps
		type AppDomain struct {
			App    string
			Domain string
		}
		allDomains := []AppDomain{}
		for appName, appConfig := range novus.GetState().Apps {
			for _, route := range appConfig.Routes {
				allDomains = append(allDomains, AppDomain{App: appName, Domain: route.Domain})
			}
		}
		for _, route := range addedRoutes {
			if idx := slices.IndexFunc(allDomains, func(appDomain AppDomain) bool { return appDomain.Domain == route.Domain }); idx != -1 {
				usedDomainAppName := allDomains[idx].App
				logger.Errorf("Domain %s is already defined by app \"%s\"", route.Domain, usedDomainAppName)
				logger.Hintf("Use a different domain name or temporarily stop \"%[1]s\" by running `novus stop %[1]s`", usedDomainAppName)
				os.Exit(1)
			}
			allDomains = append(allDomains, AppDomain{App: conf.AppName, Domain: route.Domain})
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
