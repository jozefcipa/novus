package cmd

import (
	"github.com/jozefcipa/novus/internal/brew"
	"github.com/jozefcipa/novus/internal/config"
	"github.com/jozefcipa/novus/internal/dnsmasq"
	"github.com/jozefcipa/novus/internal/logger"
	"github.com/jozefcipa/novus/internal/mkcert"
	"github.com/jozefcipa/novus/internal/nginx"
	"github.com/jozefcipa/novus/internal/novus"
	"github.com/jozefcipa/novus/internal/ssl_manager"

	"github.com/spf13/cobra"
)

var shouldCreateConfigFile bool

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		brew.InstallBinaries()                      // Make sure we have the necessary binaries available
		novus.LoadState()                           // Load application state
		conf := config.Load(shouldCreateConfigFile) // Load configuration file

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
			logger.Checkf("Nginx running")
		}

		// DNSMasq
		isDNSMasqRunning := dnsmasq.IsRunning()
		if dnsMasqConfigUpdated || !isDNSMasqRunning {
			dnsmasq.Restart()
			logger.Checkf("DNSMasq restarted")
		} else {
			logger.Checkf("DNSMasq running")
		}

		// Everything's set, start routing
		logger.Checkf("Routing started 🚀")
		for _, route := range conf.Routes {
			logger.Infof("  - %s -> ", route.Upstream)
			logger.Successf("https://%s\n", route.Domain)
		}

		// Save application state
		novus.SaveState()
	},
}

func init() {
	serveCmd.PersistentFlags().BoolVar(&shouldCreateConfigFile, "create-config", false, "create a configuration file")
	rootCmd.AddCommand(serveCmd)
}
