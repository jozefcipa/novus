package dns_manager

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jozefcipa/novus/internal/config"
	"github.com/jozefcipa/novus/internal/dnsmasq"
	"github.com/jozefcipa/novus/internal/fs"
	"github.com/jozefcipa/novus/internal/logger"
	"github.com/jozefcipa/novus/internal/maputils"
	"github.com/jozefcipa/novus/internal/novus"
	"github.com/jozefcipa/novus/internal/paths"
	"github.com/jozefcipa/novus/internal/ports"
	"github.com/jozefcipa/novus/internal/sharedtypes"
	"github.com/jozefcipa/novus/internal/sudo"
	"github.com/jozefcipa/novus/internal/tld"
	"github.com/jozefcipa/novus/internal/tui"
)

// A flag to indicate if DNS port was updated during the run
var dnsPortUpdated = false

// DNSMasq & DNS resolver setup
// https://gist.github.com/ogrrd/5831371

func GetTLDs(routes []sharedtypes.Route) []string {
	var tlds = make(map[string]bool)

	for _, route := range routes {
		tld := tld.ExtractFromDomain(route.Domain)

		if _, ok := tlds[tld]; !ok {
			tlds[tld] = true
		}
	}

	return maputils.MapKeys(tlds)
}

func Configure(config config.NovusConfig, novusState *novus.NovusState) bool {
	dnsPort := GetDNSPort(novusState)

	// Update main DNSMasq configuration
	updated := dnsmasq.Configure(dnsPort)

	// Create the DNS resolver directory if not exists
	// https://www.manpagez.com/man/5/resolver/
	sudo.MakeDirOrExit(paths.DNSResolverDir)

	// Create configs for each TLD
	tlds := GetTLDs(config.Routes)
	// Include internal domains
	tlds = append(tlds, GetTLDs(novusState.Apps[novus.NovusInternalAppName].Routes)...)

	for _, tld := range tlds {
		configCreated, configPath := dnsmasq.CreateTLDConfig(tld)

		// Initialize state struct if not exists
		if _, exists := novusState.DnsFiles[tld]; !exists {
			novusState.DnsFiles[tld] = &novus.DnsFiles{}
		}

		if configCreated {
			updated = true
			// Store config path in state
			novusState.DnsFiles[tld].DnsMasqConfig = configPath
		}

		// Register the system's DNS TLD resolver
		resolverCreated, resolverPath := registerTLDResolver(tld, dnsPort)
		if resolverCreated {
			updated = true
			// Store config path in state
			novusState.DnsFiles[tld].DnsResolver = resolverPath
		}
	}

	if updated {
		logger.Checkf("DNS configuration updated")
		return true
	} else {
		logger.Debugf("DNS configuration is up to date")
		return false
	}
}

func registerTLDResolver(tld string, dnsPort string) (bool, string) {
	configPath := filepath.Join(paths.DNSResolverDir, tld)

	// First check if the file already exists (but only if the port was not changed)
	if !dnsPortUpdated {
		if fExists := fs.FileExists(configPath); fExists {
			logger.Debugf("DNS resolver for TLD *.%s already exists [%s]", tld, configPath)
			return false, configPath
		}
	}

	logger.Debugf("Creating/updating DNS resolver [*.%s] (DNS port: %s)", tld, dnsPort)

	// Create a configuration file
	configContent := fmt.Sprintf("nameserver 127.0.0.1\nport %s\n", dnsPort)
	sudo.WriteFileOrExit(configPath, configContent)
	logger.Debugf("DNS resolver for *.%s saved [%s]", tld, configPath)

	return true, configPath
}

func UnregisterTLD(tld string, novusState *novus.NovusState) {
	if novusState.DnsFiles[tld].DnsMasqConfig != "" {
		logger.Debugf("Deleting DNSMasq configuration for *.%s [%s]", tld, novusState.DnsFiles[tld].DnsMasqConfig)
		fs.DeleteFile(novusState.DnsFiles[tld].DnsMasqConfig)
	}

	if novusState.DnsFiles[tld].DnsResolver != "" {
		logger.Infof("Deleting DNS resolver for *.%s", tld)
		logger.Debugf("*.%s resolver saved in %s", tld, novusState.DnsFiles[tld].DnsResolver)
		err := sudo.DeleteFile(novusState.DnsFiles[tld].DnsResolver)
		if err != nil {
			logger.Debugf(err.Error())
		}
	}

	// Remove from state
	delete(novusState.DnsFiles, tld)
}

func GetDNSPort(novusState *novus.NovusState) string {
	if novusState.DNSMasq.Port != "" && novusState.DNSMasq.Port != dnsmasq.DefaultPort {
		logger.Debugf("Using custom DNS port from state: %s", novusState.DNSMasq.Port)
		return novusState.DNSMasq.Port
	}

	return dnsmasq.DefaultPort
}

func EnsurePort(initialPortsUsage ports.PortUsage, novusState *novus.NovusState) {
	dnsPort := GetDNSPort(novusState)
	novusState.DNSMasq.Port = dnsPort

	if portUsedBy, isUsed := initialPortsUsage[dnsPort]; isUsed && portUsedBy != "dnsmasq" {
		logger.Errorf("Cannot start DNSMasq: Port %s is already used by '%s'", dnsPort, portUsedBy)

		// Ask user for an alternative port
		var alternativePort string
		attempts := 0
		for {
			if attempts >= 3 {
				logger.Errorf("Failed to set an alternative DNS port after %d attempts, exiting.", attempts)
				os.Exit(1)
			}

			alternativePort = tui.AskUser("Choose an alternative port for DNS: ")

			// Check if the port is a valid port number
			if !ports.IsValidPort(alternativePort) {
				logger.Errorf("'%s' is not a valid port number, please choose a port between 1 and 65535.", alternativePort)
				attempts += 1
				continue
			}

			// Check if the alternative port is available
			portUsage := ports.CheckPortsUsage(alternativePort)
			if portUsedBy, isUsed := portUsage[alternativePort]; isUsed {
				logger.Errorf("Port %s is already used by '%s', please choose another one.", alternativePort, portUsedBy)
				attempts += 1
				continue
			}

			break
		}

		// Update all TLD resolvers
		dnsPortUpdated = true
		for tld := range novusState.DnsFiles {
			registerTLDResolver(tld, alternativePort)
		}

		novusState.DNSMasq.Port = alternativePort
	}
}
