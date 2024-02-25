package dnsmasq

import (
	"fmt"
	"strings"

	"github.com/jozefcipa/novus/internal/brew"
	"github.com/jozefcipa/novus/internal/config"
	"github.com/jozefcipa/novus/internal/fs"
	"github.com/jozefcipa/novus/internal/logger"
)

var dnsmasqConfFile string
var dnsResolverDir string

func init() {
	dnsmasqConfFile = brew.BrewPath + "/etc/dnsmasq.conf"
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

func Configure(config config.NovusConfig) bool {
	// update main dnsmasq config
	updated := updateDNSMasqConfig()

	// create the DNS resolver directory
	// https://www.manpagez.com/man/5/resolver/
	fs.MakeDirWithSudoOrExit(dnsResolverDir)

	// create configs for each domain
	domains := getDomains(config)
	for _, domain := range domains {
		// create DNSMasq domain config
		configCreated := createDNSMasqDomainConfig(domain)

		if !updated && configCreated {
			updated = true
		}

		// register the system's DNS domain resolver
		resolverCreated := registerDomainResolver(domain)
		if !updated && resolverCreated {
			updated = true
		}
	}

	if updated {
		logger.Checkf("DNSMasq configuration updated")
		return true
	} else {
		logger.Checkf("DNSMasq config is up to date")
		return false
	}
}

func updateDNSMasqConfig() bool {
	// open file
	logger.Debugf("Reading DNSMasq configuration file %s", dnsmasqConfFile)
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
		logger.Debugf("Updating DNSMasq configuration file")
		fs.WriteFileOrExit(dnsmasqConfFile, updatedConf)

		return true
	} else {
		logger.Debugf("DNSMasq configuration is up to date")

		return false
	}
}

func createDNSMasqDomainConfig(domain string) bool {
	configPath := fmt.Sprintf("%s/etc/dnsmasq.d/%s.conf", brew.BrewPath, domain)

	// first check if the file already exists
	if confExists := fs.FileExists(configPath); confExists {
		logger.Debugf("DNSMasq [*.%s]: Domain config already exists [%s]", domain, configPath)

		return false
	}

	logger.Debugf("[*.%s] Creating DNSMasq domain config", domain)

	// prepare the configuration
	configContent := fmt.Sprintf("address=/%s/127.0.0.1", domain)

	// create a configuration file
	fs.WriteFileOrExit(configPath, configContent)
	logger.Debugf("DNSMasq [*.%s]: Domain config saved [%s]", domain, configPath)

	return true
}

func registerDomainResolver(domain string) bool {
	configPath := fmt.Sprintf("%s/%s", dnsResolverDir, domain)

	// first check if the file already exists
	if fExists := fs.FileExists(configPath); fExists {
		logger.Debugf("DNSMasq [*.%s]: Domain resolver already exists [%s]", domain, configPath)
		return false
	}

	logger.Debugf("DNSMasq [*.%s]: Creating domain resolver", domain)

	// create a configuration file
	configContent := "nameserver 127.0.0.1\n"
	fs.WriteFileWithSudoOrExit(configPath, configContent)
	logger.Debugf("DNSMasq [*.%s]: Domain resolver saved [%s]", domain, configPath)

	return true
}

func getDomains(conf config.NovusConfig) []string {
	var domains = make(map[string]bool)

	for _, route := range conf.Routes {
		urlParts := strings.Split(route.Domain, ".")
		domain := urlParts[len(urlParts)-1]

		if _, ok := domains[domain]; !ok {
			domains[domain] = true
		}
	}

	// get domain keys
	keys := make([]string, 0, len(domains))
	for k := range domains {
		keys = append(keys, k)
	}
	return keys
}
