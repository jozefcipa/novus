package cmd

import (
	"os"
	"slices"

	"github.com/jozefcipa/novus/internal/dnsmasq"
	"github.com/jozefcipa/novus/internal/homebrew"
	"github.com/jozefcipa/novus/internal/logger"
	"github.com/jozefcipa/novus/internal/nginx"
	"github.com/jozefcipa/novus/internal/novus"
	"github.com/jozefcipa/novus/internal/ports"
	"github.com/jozefcipa/novus/internal/tui"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start Nginx and DNSMasq services",
	Long:  `Start Nginx, DNSMasq and start routing URLs.`,
	Run: func(cmd *cobra.Command, args []string) {
		// If the binaries are missing, exit here, user needs to run `novus init` first
		if err := homebrew.CheckIfRequiredBinariesInstalled(); err != nil {
			logger.Hintf("Run \"novus init\" first to initialize Novus.")
			os.Exit(1)
		}

		novusState := novus.GetState()

		// Check if ports are available
		portsUsage := ports.CheckIfAvailable(slices.Concat(nginx.Ports, []string{dnsmasq.Port})...)
		nginx.EnsurePortsAvailable(portsUsage)
		dnsmasq.EnsurePortAvailable(portsUsage)

		// Restart services
		// Nginx
		nginxLoader := logger.Loadingf("Checking Nginx status")
		isNginxRunning := nginx.IsRunning()
		if !isNginxRunning {
			nginxLoader.Done()
			nginx.Restart()
		} else {
			nginxLoader.Checkf("Nginx running")
		}

		// DNSMasq
		dnsmasqLoader := logger.Loadingf("Checking DNSMasq status")
		isDNSMasqRunning := dnsmasq.IsRunning()
		if !isDNSMasqRunning {
			dnsmasqLoader.Done()
			dnsmasq.Restart()
		} else {
			dnsmasqLoader.Checkf("DNSMasq running")
		}

		// Everything's set, start routing
		tui.PrintRoutingTable(*novusState)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
