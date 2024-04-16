package dnsmasq

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/jozefcipa/novus/internal/brew"
	"github.com/jozefcipa/novus/internal/config"
	"github.com/jozefcipa/novus/internal/fs"
	"github.com/jozefcipa/novus/internal/logger"
	"github.com/jozefcipa/novus/internal/novus"
	"github.com/jozefcipa/novus/internal/shared"
)

var dnsmasqConfFile string
var dnsResolverDir string

func init() {
	dnsmasqConfFile = filepath.Join(brew.BrewPath, "/etc/dnsmasq.conf")
	dnsResolverDir = "/etc/resolver"
}

func Restart() {
	brew.RestartServiceWithSudo("dnsmasq")
}

func Stop() {
	brew.StopServiceWithSudo("dnsmasq")
}

func IsRunning() bool {
	return brew.IsSudoServiceRunning("dnsmasq")
}

func GetTLDs(routes []shared.Route) []string {
	var tlds = make(map[string]bool)

	for _, route := range routes {
		urlParts := strings.Split(route.Domain, ".")
		domain := urlParts[len(urlParts)-1]

		if _, ok := tlds[domain]; !ok {
			tlds[domain] = true
		}
	}

	// get domain keys
	keys := make([]string, 0, len(tlds))
	for k := range tlds {
		keys = append(keys, k)
	}
	return keys
}

func Configure(config config.NovusConfig) bool {
	state := novus.GetState()

	// update main dnsmasq config
	updated := updateDNSMasqConfig()

	// create the DNS resolver directory
	// https://www.manpagez.com/man/5/resolver/
	fs.MakeDirWithSudoOrExit(dnsResolverDir)

	// create configs for each TLD
	tlds := GetTLDs(config.Routes)
	for _, tld := range tlds {
		// create DNSMasq TLD config
		configCreated, configPath := createDNSMasqTLDConfig(tld)

		// initialize state struct if not exists
		if _, exists := state.DnsFiles[tld]; !exists {
			state.DnsFiles[tld] = &novus.DnsFiles{}
		}

		if configCreated {
			updated = true
			// store config path in state
			state.DnsFiles[tld].DnsMasqConfig = configPath
		}

		// register the system's DNS TLD resolver
		resolverCreated, resolverPath := registerTLDResolver(tld)
		if resolverCreated {
			updated = true
			// store config path in state
			state.DnsFiles[tld].DnsResolver = resolverPath
		}
	}

	if updated {
		logger.Checkf("DNSMasq: Configuration updated")
		return true
	} else {
		logger.Checkf("DNSMasq: Configuration is up to date")
		return false
	}
}

func UnregisterTLD(tld string) {
	state := novus.GetState()

	if state.DnsFiles[tld].DnsMasqConfig != "" {
		logger.Debugf("DNSMasq [*.%s]: Deleting DNSMasq config [%s]", tld, state.DnsFiles[tld].DnsMasqConfig)
		fs.DeleteFile(state.DnsFiles[tld].DnsMasqConfig)
	}

	if state.DnsFiles[tld].DnsResolver != "" {
		logger.Debugf("DNSMasq [*.%s]: Deleting DNS resolver [%s]", tld, state.DnsFiles[tld].DnsResolver)
		fs.DeleteFileWithSudo(state.DnsFiles[tld].DnsResolver)
	}

	// remove from state
	delete(state.DnsFiles, tld)
}

func updateDNSMasqConfig() bool {
	// open file
	logger.Debugf("DNSMasq: Reading configuration file [%s]", dnsmasqConfFile)
	confFile := fs.ReadFileOrExit(dnsmasqConfFile)

	// remove comment "#"
	updatedConf := strings.Replace(
		string(confFile),
		fmt.Sprintf("#conf-dir=%s/etc/dnsmasq.d/,*.conf", brew.BrewPath),
		fmt.Sprintf("conf-dir=%s/etc/dnsmasq.d/,*.conf", brew.BrewPath),
		1,
	)

	// if the config differs (there was an actual change), write the changes
	if confFile != updatedConf {
		logger.Debugf("DNSMasq: Updating configuration file [%s]", dnsmasqConfFile)
		fs.WriteFileOrExit(dnsmasqConfFile, updatedConf)

		return true
	} else {
		logger.Debugf("DNSMasq: Configuration file is up to date [%s]", dnsmasqConfFile)

		return false
	}
}

func createDNSMasqTLDConfig(tld string) (bool, string) {
	configPath := fmt.Sprintf(filepath.Join(brew.BrewPath, "/etc/dnsmasq.d/%s.conf"), tld)

	// first check if the file already exists
	if confExists := fs.FileExists(configPath); confExists {
		logger.Debugf("DNSMasq [*.%s]: Domain config already exists [%s]", tld, configPath)

		return false, configPath
	}

	logger.Debugf("DNSMasq [*.%s]: Creating domain config", tld)

	// prepare the configuration
	configContent := fmt.Sprintf("address=/%s/127.0.0.1", tld)

	// create a configuration file
	fs.WriteFileOrExit(configPath, configContent)
	logger.Debugf("DNSMasq [*.%s]: Domain config saved [%s]", tld, configPath)

	return true, configPath
}

func registerTLDResolver(tld string) (bool, string) {
	configPath := fmt.Sprintf("%s/%s", dnsResolverDir, tld)

	// first check if the file already exists
	if fExists := fs.FileExists(configPath); fExists {
		logger.Debugf("DNSMasq [*.%s]: Domain resolver already exists [%s]", tld, configPath)
		return false, configPath
	}

	logger.Debugf("DNSMasq [*.%s]: Creating domain resolver", tld)

	// create a configuration file
	configContent := "nameserver 127.0.0.1\n"
	fs.WriteFileWithSudoOrExit(configPath, configContent)
	logger.Debugf("DNSMasq [*.%s]: Domain resolver saved [%s]", tld, configPath)

	return true, configPath
}
