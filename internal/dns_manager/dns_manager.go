package dns_manager

import (
	"path/filepath"

	"github.com/jozefcipa/novus/internal/config"
	"github.com/jozefcipa/novus/internal/dnsmasq"
	"github.com/jozefcipa/novus/internal/fs"
	"github.com/jozefcipa/novus/internal/logger"
	"github.com/jozefcipa/novus/internal/net"
	"github.com/jozefcipa/novus/internal/novus"
	"github.com/jozefcipa/novus/internal/shared"
)

// DNSMasq & DNS resolver setup
// https://gist.github.com/ogrrd/5831371

const dnsResolverDir = "/etc/resolver"

func GetTLDs(routes []shared.Route) []string {
	var tlds = make(map[string]bool)

	for _, route := range routes {
		tld := net.ExtractTLD(route.Domain)

		if _, ok := tlds[tld]; !ok {
			tlds[tld] = true
		}
	}

	return shared.MapKeys(tlds)
}

func Configure(config config.NovusConfig, novusState *novus.NovusState) bool {
	// Update main DNSMasq configuration
	updated := dnsmasq.Configure()

	// Create the DNS resolver directory if not exists
	// https://www.manpagez.com/man/5/resolver/
	fs.MakeDirWithSudoOrExit(dnsResolverDir)

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
		resolverCreated, resolverPath := registerTLDResolver(tld)
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
		logger.Checkf("DNS configuration is up to date")
		return false
	}
}

func registerTLDResolver(tld string) (bool, string) {
	configPath := filepath.Join(dnsResolverDir, tld)

	// First check if the file already exists
	if fExists := fs.FileExists(configPath); fExists {
		logger.Debugf("DNS resolver for TLD *.%s already exists [%s]", tld, configPath)
		return false, configPath
	}

	logger.Debugf("Creating DNS resolver [*.%s]", tld)

	// Create a configuration file
	configContent := "nameserver 127.0.0.1\n"
	fs.WriteFileWithSudoOrExit(configPath, configContent)
	logger.Debugf("DNS resolver for *.%s saved [%s]", tld, configPath)

	return true, configPath
}

func UnregisterTLD(tld string, novusState *novus.NovusState) {
	if novusState.DnsFiles[tld].DnsMasqConfig != "" {
		logger.Debugf("Deleting DNSMasq configuration for *.%s [%s]", tld, novusState.DnsFiles[tld].DnsMasqConfig)
		fs.DeleteFile(novusState.DnsFiles[tld].DnsMasqConfig)
	}

	if novusState.DnsFiles[tld].DnsResolver != "" {
		logger.Debugf("Deleting DNS resolver for *.%s [%s]", tld, novusState.DnsFiles[tld].DnsResolver)
		fs.DeleteFileWithSudo(novusState.DnsFiles[tld].DnsResolver)
	}

	// Remove from state
	delete(novusState.DnsFiles, tld)
}
