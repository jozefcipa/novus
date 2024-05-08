package dns_manager

import (
	"path/filepath"
	"strings"

	"github.com/jozefcipa/novus/internal/config"
	"github.com/jozefcipa/novus/internal/dnsmasq"
	"github.com/jozefcipa/novus/internal/fs"
	"github.com/jozefcipa/novus/internal/logger"
	"github.com/jozefcipa/novus/internal/novus"
	"github.com/jozefcipa/novus/internal/shared"
)

// TODO: describe the steps here, how it works
// what's /etc/resolver & what's DNSMasq

var dnsResolverDir string

func init() {
	dnsResolverDir = "/etc/resolver"
}

func GetTLDs(routes []shared.Route) []string {
	var tlds = make(map[string]bool)

	for _, route := range routes {
		urlParts := strings.Split(route.Domain, ".")
		tld := urlParts[len(urlParts)-1]

		if _, ok := tlds[tld]; !ok {
			tlds[tld] = true
		}
	}

	return shared.MapKeys(tlds)
}

func Configure(config config.NovusConfig) bool {
	state := novus.GetState()

	// Update main DNSMasq configuration
	updated := dnsmasq.Configure()

	// Create the DNS resolver directory if not exists
	// https://www.manpagez.com/man/5/resolver/
	fs.MakeDirWithSudoOrExit(dnsResolverDir)

	// Create configs for each TLD
	tlds := GetTLDs(config.Routes)
	for _, tld := range tlds {
		configCreated, configPath := dnsmasq.CreateTLDConfig(tld)

		// Initialize state struct if not exists
		if _, exists := state.DnsFiles[tld]; !exists {
			state.DnsFiles[tld] = &novus.DnsFiles{}
		}

		if configCreated {
			updated = true
			// Store config path in state
			state.DnsFiles[tld].DnsMasqConfig = configPath
		}

		// Register the system's DNS TLD resolver
		resolverCreated, resolverPath := registerTLDResolver(tld)
		if resolverCreated {
			updated = true
			// Store config path in state
			state.DnsFiles[tld].DnsResolver = resolverPath
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

func UnregisterTLD(tld string) {
	state := novus.GetState()

	if state.DnsFiles[tld].DnsMasqConfig != "" {
		logger.Debugf("Deleting DNSMasq configuration for *.%s [%s]", tld, state.DnsFiles[tld].DnsMasqConfig)
		fs.DeleteFile(state.DnsFiles[tld].DnsMasqConfig)
	}

	if state.DnsFiles[tld].DnsResolver != "" {
		logger.Debugf("Deleting DNS resolver for *.%s [%s]", tld, state.DnsFiles[tld].DnsResolver)
		fs.DeleteFileWithSudo(state.DnsFiles[tld].DnsResolver)
	}

	// Remove from state
	delete(state.DnsFiles, tld)
}
