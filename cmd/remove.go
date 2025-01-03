package cmd

import (
	"fmt"
	"os"
	"slices"

	"github.com/jozefcipa/novus/internal/config_manager"
	"github.com/jozefcipa/novus/internal/dnsmasq"
	"github.com/jozefcipa/novus/internal/domain_cleanup_manager"
	"github.com/jozefcipa/novus/internal/logger"
	"github.com/jozefcipa/novus/internal/nginx"
	"github.com/jozefcipa/novus/internal/novus"
	"github.com/jozefcipa/novus/internal/sharedtypes"
	"github.com/jozefcipa/novus/internal/ssl_manager"
	"github.com/jozefcipa/novus/internal/tui"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove [app-name]",
	Short: "Remove routing configuration for [app-name]",
	Long:  "Remove all domains registered in the configuration for the given app",
	Run: func(cmd *cobra.Command, args []string) {
		_, appState := tui.ParseAppFromArgs(args, "remove")

		if appState == nil {
			// Deleting domain
			domain := args[0]
			novusState := novus.GetState()
			// Get global app configuration
			conf := config_manager.LoadConfigurationFromState(novus.GlobalAppName, *novusState)

			// Check if the domain exists
			idx := slices.IndexFunc(conf.Routes, func(route sharedtypes.Route) bool { return route.Domain == domain })
			if idx == -1 {
				logger.Errorf("Domain or app \"%s\" does not exist", domain)
				os.Exit(1)
			}

			// Confirm deleting
			if !tui.Confirm(fmt.Sprintf("Do you want to remove \"%s\" domain?", domain)) {
				os.Exit(0)
			}

			// Remove domain from the config
			conf.Routes = append(conf.Routes[:idx], conf.Routes[idx+1:]...)

			// Delete route
			domain_cleanup_manager.RemoveDomains([]sharedtypes.Route{{Domain: domain}}, novus.GlobalAppName, novusState)

			// Configure SSL
			domainCerts, _ := ssl_manager.EnsureSSLCertificates(conf, novusState, novus.GlobalAppName)

			// Update NGINX configuration
			appState, _ := novus.GetAppState(novus.GlobalAppName)
			nginx.Configure(conf, domainCerts, appState)

			logger.Checkf("Domain [%s] has been removed", domain)
		} else {
			// Deleting app
			appName := args[0]
			if !tui.Confirm(fmt.Sprintf("Do you want to remove \"%s\" configuration?", appName)) {
				os.Exit(0)
			}

			// Delete all routes
			domain_cleanup_manager.RemoveDomains(appState.Routes, appName, novus.GetState())

			// Remove NGINX configuration
			nginx.RemoveConfiguration(appName)

			// Remove app from Novus state
			novus.RemoveAppState(appName)

			logger.Checkf("App \"%s\" has been removed", appName)
		}

		// Restart services
		nginx.Restart()
		dnsmasq.Restart()

		tui.PrintRoutingTable(*novus.GetState())

		// Save state to file
		novus.SaveState()
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
