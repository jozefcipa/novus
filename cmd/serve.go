package cmd

import (
	"os"
	"slices"
	"strings"

	"github.com/jozefcipa/novus/internal/config"
	"github.com/jozefcipa/novus/internal/config_manager"
	"github.com/jozefcipa/novus/internal/diff_manager"
	"github.com/jozefcipa/novus/internal/dns_manager"
	"github.com/jozefcipa/novus/internal/dnsmasq"
	"github.com/jozefcipa/novus/internal/domain_cleanup_manager"
	"github.com/jozefcipa/novus/internal/homebrew"
	"github.com/jozefcipa/novus/internal/logger"
	"github.com/jozefcipa/novus/internal/mkcert"
	"github.com/jozefcipa/novus/internal/nginx"
	"github.com/jozefcipa/novus/internal/novus"
	"github.com/jozefcipa/novus/internal/paths"
	"github.com/jozefcipa/novus/internal/ports"
	"github.com/jozefcipa/novus/internal/sharedtypes"
	"github.com/jozefcipa/novus/internal/ssl_manager"
	"github.com/jozefcipa/novus/internal/tui"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve [domain?] [upstream?]",
	Short: "Configure URLs and start routing",
	Long:  `Install Nginx, DNSMasq and mkcert and automatically expose HTTPs URLs for the endpoints defined in the config.`,
	Run: func(cmd *cobra.Command, args []string) {
		// If the binaries are missing, exit here, user needs to run `novus init` first
		if err := homebrew.CheckIfRequiredBinariesInstalled(); err != nil {
			logger.Hintf("Run \"novus init\" first to initialize Novus.")
			os.Exit(1)
		}

		var conf config.NovusConfig
		novusState := novus.GetState()

		// If inline domain is provided, prioritise that
		if len(args) > 0 {
			var upstream string
			if len(args) == 2 {
				upstream = args[1]
			} else {
				upstream = tui.AskUser("Enter an upstream address (e.g. http://localhost:3000): ")
			}

			// Ensure novus state for the global app exists
			if _, appStateExists := novus.GetAppState(novus.GlobalAppName); !appStateExists {
				novus.InitializeAppState(novus.GlobalAppName, paths.NovusStateDir)
			}

			// Load Novus config for the global app
			conf = config_manager.LoadConfigurationFromState(novus.GlobalAppName, *novusState)

			// Append the new route to the config
			conf.Routes = append(conf.Routes, sharedtypes.Route{Domain: args[0], Upstream: upstream})

			// Validate input
			if errors := config_manager.ValidateConfig(conf, config_manager.ValidationErrorsGlobalAppInput); len(errors) > 0 {
				logger.Errorf("Invalid configuration:\n   %s", strings.Join(errors, "\n   "))
				os.Exit(1)
			}
		} else {
			// Otherwise, load configuration file
			var exists bool
			conf, exists = config_manager.LoadConfigurationFromFile(*novusState)
			if !exists {
				logger.Warnf("Novus is not initialized in this directory (%s file does not exist).", config.ConfigFileName)
				logger.Hintf("Run \"novus init\" to create a configuration file.")
				os.Exit(1)
			}
		}

		config_manager.ValidateConfigDomainsUniqueness(conf, *novusState)
		appName := config.AppName()

		// Load application state
		appState, appStateExists := novus.GetAppState(appName)
		if !appStateExists {
			appState = novus.InitializeAppState(appName, paths.CurrentDir)
		}

		// Compare state and current config to detect changes
		addedRoutes, deletedRoutes := diff_manager.DetectConfigDiff(conf, *appState)

		// Remove domains that are no longer in config
		if len(deletedRoutes) > 0 {
			domain_cleanup_manager.RemoveDomains(deletedRoutes, appName, novusState)
		}

		if len(addedRoutes) > 0 {
			if len(addedRoutes) == 1 {
				logger.Successf("Found a new domain [%s]", addedRoutes[0].Domain)
			} else {
				logger.Successf("Found %d new domains:", len(addedRoutes))
				for _, newRoute := range addedRoutes {
					logger.Infof("   - %s", newRoute.Domain)
				}
			}
		}

		// Check if ports are available
		portsUsage := ports.CheckIfAvailable(slices.Concat(nginx.Ports, []string{dnsmasq.Port})...)
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
		nginxLoader := logger.Loadingf("Checking Nginx status")
		isNginxRunning := nginx.IsRunning()
		if nginxConfigUpdated || hasNewCerts || !isNginxRunning {
			nginxLoader.Done()
			nginx.Restart()
		} else {
			nginxLoader.Checkf("Nginx running")
		}

		// DNSMasq
		dnsmasqLoader := logger.Loadingf("Checking DNSMasq status")
		isDNSMasqRunning := dnsmasq.IsRunning()
		if dnsUpdated || !isDNSMasqRunning {
			dnsmasqLoader.Done()
			dnsmasq.Restart()
		} else {
			dnsmasqLoader.Checkf("DNSMasq running")
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
