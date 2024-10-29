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

var removeCmd = &cobra.Command{
	Use:   "remove [app-name]",
	Short: "Remove routing configuration for [app-name]",
	Long:  "Remove all domains registered in the configuration for the given app",
	Run: func(cmd *cobra.Command, args []string) {
		appName, appState := tui.ParseAppFromArgs(args, "remove")

		if !tui.Confirm(fmt.Sprintf("Do you want to remove \"%s\" configuration?", appName)) {
			os.Exit(0)
		}

		// Delete all routes
		domain_cleanup_manager.RemoveDomains(appState.Routes, appName, novus.GetState())

		// Remove NGINX configuration
		nginx.RemoveConfiguration(appName)

		// Remove app from Novus state
		novus.RemoveAppState(appName)

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
	rootCmd.AddCommand(removeCmd)
}
