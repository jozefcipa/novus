package cmd

import (
	"os"
	"slices"

	"github.com/jozefcipa/novus/internal/brew"
	"github.com/jozefcipa/novus/internal/config"
	"github.com/jozefcipa/novus/internal/config_manager"
	"github.com/jozefcipa/novus/internal/diff_manager"
	"github.com/jozefcipa/novus/internal/dns_manager"
	"github.com/jozefcipa/novus/internal/dnsmasq"
	"github.com/jozefcipa/novus/internal/domain_cleanup_manager"
	"github.com/jozefcipa/novus/internal/logger"
	"github.com/jozefcipa/novus/internal/mkcert"
	"github.com/jozefcipa/novus/internal/net"
	"github.com/jozefcipa/novus/internal/nginx"
	"github.com/jozefcipa/novus/internal/novus"
	"github.com/jozefcipa/novus/internal/ssl_manager"
	"github.com/jozefcipa/novus/internal/tui"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Configure URLs and start routing",
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

		// Load Novus state & validate config file
		novusState := novus.GetState()
		config_manager.ValidateConfig(conf, *novusState)
		appName := config.AppName()

		// Load application state
		appState, appStateExists := novus.GetAppState(appName)
		if !appStateExists {
			appState = novus.InitializeAppState(appName)
		}

		// Compare state and current config to detect changes
		addedRoutes, deletedRoutes := diff_manager.DetectConfigDiff(conf, *appState)

		// Remove domains that are no longer in config
		if len(deletedRoutes) > 0 {
			domain_cleanup_manager.RemoveDomains(deletedRoutes, appName, novusState)
		}

		if len(addedRoutes) > 0 {
			for _, newRoute := range addedRoutes {
				logger.Successf("Found new domain [%s]", newRoute.Domain)
			}
		}

		// Check if ports are available
		portsUsage := net.CheckPortsUsage(slices.Concat(nginx.Ports, []string{dnsmasq.Port})...)
		nginx.EnsurePortsAvailable(portsUsage)
		dnsmasq.EnsurePortAvailable(portsUsage)

		// Configure SSL
		mkcert.Configure(conf)
		domainCerts, hasNewCerts := ssl_manager.EnsureSSLCertificates(conf, novusState, appName)

		// Configure Nginx
		nginxConfigUpdated := nginx.Configure(conf, domainCerts, appState)

		// Configure DNS
		dnsUpdated := dns_manager.Configure(conf, novusState)

		// Restart services
		// Nginx
		isNginxRunning := nginx.IsRunning()
		if nginxConfigUpdated || hasNewCerts || !isNginxRunning {
			nginx.Restart()
		} else {
			logger.Checkf("Nginx running 🚀")
		}

		// DNSMasq
		isDNSMasqRunning := dnsmasq.IsRunning()
		if dnsUpdated || !isDNSMasqRunning {
			dnsmasq.Restart()
		} else {
			logger.Checkf("DNSMasq running 🚀")
		}

		// If app has been paused, make sure to set it to ACTIVE
		appState.Status = novus.APP_ACTIVE

		// Everything's set, start routing
		tui.PrintRoutingTable(*novusState)

		// Save application state
		novus.SaveState()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
