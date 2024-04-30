package cmd

import (
	"github.com/jozefcipa/novus/internal/dnsmasq"
	"github.com/jozefcipa/novus/internal/logger"
	"github.com/jozefcipa/novus/internal/nginx"
	"github.com/jozefcipa/novus/internal/novus"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show status of services and registered routes",
	Long: `Show whether Nginx and DNSMasq services are running,
and print a list of all URLs that are registered by Novus.`,
	Run: func(cmd *cobra.Command, args []string) {
		novusState := novus.GetState()

		nginxChan := make(chan bool)
		dnsMasqChan := make(chan bool)
		go func() {
			nginxChan <- nginx.IsRunning()
		}()
		go func() {
			dnsMasqChan <- dnsmasq.IsRunning()
		}()

		isNginxRunning := <-nginxChan
		isDNSMasqRunning := <-dnsMasqChan

		if isNginxRunning {
			logger.Successf("✅ Nginx running.\n")
			logger.Debugf("Nginx configuration loaded from %s", nginx.NginxServersDir)
		} else {
			logger.Errorf("❌ Nginx not running.\n")
		}

		if isDNSMasqRunning {
			logger.Successf("✅ DNSMasq running.\n")
		} else {
			logger.Errorf("❌ DNSMasq not running.\n")
		}

		if !isNginxRunning || !isDNSMasqRunning {
			logger.Errorf("Please run `novus serve` to initialize the services.\n")
		} else {
			// All good, show the routing info
			for appName, appState := range novusState.Apps {
				logger.Checkf("Routing %s [%s]", appName, appState.Directory)
				for _, route := range appState.Routes {
					logger.Infof("  - %s -> ", route.Upstream)
					logger.Successf("https://%s\n", route.Domain)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
