package cmd

import (
	"github.com/jozefcipa/novus/internal/brew"
	"github.com/jozefcipa/novus/internal/config"
	"github.com/jozefcipa/novus/internal/dnsmasq"
	"github.com/jozefcipa/novus/internal/logger"
	"github.com/jozefcipa/novus/internal/mkcert"
	"github.com/jozefcipa/novus/internal/nginx"
	"github.com/jozefcipa/novus/internal/ssl_manager"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Make sure we have the necessary binaries available
		brew.InstallBinaries()

		// Load configuration file
		conf := config.Load()

		// Configure services
		shouldRestartNginx := nginx.Configure(conf)
		shouldRestartDNSMasq := dnsmasq.Configure(conf)
		mkcert.Configure(conf)
		ssl_manager.EnsureSSLCertificates(conf)

		// TODO: should start if not running
		logger.Debugf("should restart nginx: %t", shouldRestartNginx)
		logger.Debugf("should restart dnsmasq: %t", shouldRestartDNSMasq)

		// Reload services
		if shouldRestartNginx {
			nginx.Restart() // TODO: doesn't throw an error if fails to start, maybe we should call nginx -t before launching
		}
		if shouldRestartDNSMasq {
			dnsmasq.Restart()
		}

		// Everything's set, start routing
		logger.Messagef("🚀 Starting routing...\n")
		for _, route := range conf.Routes {
			logger.Infof("  - http://%s -> ", route.Upstream)
			logger.Successf("https://%s\n", route.Url)
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
