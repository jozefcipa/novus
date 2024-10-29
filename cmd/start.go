package cmd

import (
	"os"
	"slices"

	"github.com/jozefcipa/novus/internal/brew"
	"github.com/jozefcipa/novus/internal/dnsmasq"
	"github.com/jozefcipa/novus/internal/logger"
	"github.com/jozefcipa/novus/internal/net"
	"github.com/jozefcipa/novus/internal/nginx"
	"github.com/jozefcipa/novus/internal/novus"
	"github.com/jozefcipa/novus/internal/tui"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start Nginx and DNSMasq services",
	Long:  `Start Nginx, DNSMasq and start routing URLs.`,
	Run: func(cmd *cobra.Command, args []string) {
		// If the binaries are missing, exit here, user needs to run `novus init` first
		if err := brew.CheckIfRequiredBinariesInstalled(); err != nil {
			logger.Hintf("Run \"novus init\" first to initialize Novus.")
			os.Exit(1)
		}

		novusState := novus.GetState()

		// Check if ports are available
		portsUsage := net.CheckPortsUsage(slices.Concat(nginx.Ports, []string{dnsmasq.Port})...)
		nginx.EnsurePortsAvailable(portsUsage)
		dnsmasq.EnsurePortAvailable(portsUsage)

		// Restart services
		// Nginx
		isNginxRunning := nginx.IsRunning()
		if !isNginxRunning {
			nginx.Restart()
		} else {
			logger.Checkf("Nginx running 🚀")
		}

		// DNSMasq
		isDNSMasqRunning := dnsmasq.IsRunning()
		if !isDNSMasqRunning {
			dnsmasq.Restart()
		} else {
			logger.Checkf("DNSMasq running 🚀")
		}

		// Everything's set, start routing
		tui.PrintRoutingTable(*novusState)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
