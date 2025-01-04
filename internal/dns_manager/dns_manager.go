package dns_manager

import (
	"fmt"
	"path/filepath"

	"github.com/jozefcipa/novus/internal/config"
	"github.com/jozefcipa/novus/internal/dnsmasq"
	"github.com/jozefcipa/novus/internal/fs"
	"github.com/jozefcipa/novus/internal/logger"
	"github.com/jozefcipa/novus/internal/maputils"
	"github.com/jozefcipa/novus/internal/novus"
	"github.com/jozefcipa/novus/internal/paths"
	"github.com/jozefcipa/novus/internal/sharedtypes"
	"github.com/jozefcipa/novus/internal/sudo"
	"github.com/jozefcipa/novus/internal/tld"
)

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
	// Update main DNSMasq configuration
	updated := dnsmasq.Configure()

	logger.Infof("Creating DNS resolvers")
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
	configPath := filepath.Join(paths.DNSResolverDir, tld)

	// First check if the file already exists
	if fExists := fs.FileExists(configPath); fExists {
		logger.Debugf("DNS resolver for TLD *.%s already exists [%s]", tld, configPath)
		return false, configPath
	}

	logger.Debugf("Creating DNS resolver [*.%s]", tld)

	// Create a configuration file
	configContent := fmt.Sprintf("nameserver 127.0.0.1\nport %s\n", dnsmasq.Port)
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
