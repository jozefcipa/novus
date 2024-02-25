package cmd

import (
	"github.com/jozefcipa/novus/internal/config"
	"github.com/jozefcipa/novus/internal/dnsmasq"
	"github.com/jozefcipa/novus/internal/logger"
	"github.com/jozefcipa/novus/internal/nginx"
	"github.com/jozefcipa/novus/internal/novus"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		novus.LoadState() // Load application state
		state := novus.GetState()

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
			// All good, show the information
			logger.Checkf("Routing [%s=%s]", config.AppName, state.Directory)
			for _, route := range state.Routes {
				logger.Infof("  - %s -> ", route.Upstream)
				logger.Successf("https://%s\n", route.Domain)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
