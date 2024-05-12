package cmd

import (
	"fmt"
	"os"

	"github.com/jozefcipa/novus/internal/dnsmasq"
	"github.com/jozefcipa/novus/internal/domain_cleanup_manager"
	"github.com/jozefcipa/novus/internal/logger"
	"github.com/jozefcipa/novus/internal/nginx"
	"github.com/jozefcipa/novus/internal/novus"
	"github.com/jozefcipa/novus/internal/tui"
	"github.com/spf13/cobra"
)

var pauseCmd = &cobra.Command{
	Use:   "pause [app-name]",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		appName, appState := tui.ParseAppFromArgs(args)

		if appState.Status == novus.APP_PAUSED {
			logger.Checkf("\"%s\" is already paused.", appName)
			tui.PrintRoutingTable(*novus.GetState())
			os.Exit(0)
		}

		if !tui.Confirm(fmt.Sprintf("Do you want to pause \"%s\"?", appName)) {
			os.Exit(0)
		}

		// Delete all routes
		domain_cleanup_manager.RemoveDomains(appState.Routes, appName, novus.GetState())

		// Remove NGINX configuration
		nginx.RemoveConfiguration(appName)

		// Mark app as paused so it won't be routed
		appState.Status = novus.APP_PAUSED

		// Restart services
		nginx.Restart()
		logger.Checkf("Nginx restarted 🔄")
		dnsmasq.Restart()
		logger.Checkf("DNSMasq restarted 🔄")

		tui.PrintRoutingTable(*novus.GetState())

		// Save state to file
		novus.SaveState()
	},
}

func init() {
	rootCmd.AddCommand(pauseCmd)
}
