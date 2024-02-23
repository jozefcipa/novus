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
		domainCerts := ssl_manager.EnsureSSLCertificates(conf)

		// Configure Nginx
		/* nginxConfigUpdated := */
		nginx.Configure(conf, domainCerts)

		// Configure DNSMasq
		/* dnsMasqConfigUpdated := */
		dnsmasq.Configure(conf)

		// TODO: should start if not running
		// Reload services
		// if nginxConfigUpdated || certsUpdated {
		// nginx.Restart()
		// }
		// if dnsMasqConfigUpdated {
		// dnsmasq.Restart()
		// }

		// Everything's set, start routing
		logger.Messagef("🚀 starting routing...\n")
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
