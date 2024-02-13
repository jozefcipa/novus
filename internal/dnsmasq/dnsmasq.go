package dnsmasq

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strings"

	"github.com/jozefcipa/novus/internal/brew"
	"github.com/jozefcipa/novus/internal/config"
	"github.com/jozefcipa/novus/internal/logger"
)

var dnsmasqConfFile string
var dnsResolverDir string

func init() {
	dnsmasqConfFile = brew.BrewPath + "/etc/dnsmasq.conf"
	dnsResolverDir = "/etc/resolver"
}

func Restart() {
	brew.RestartBrewServiceWithSudo("dnsmasq")
	logger.Checkf("DNSMasq restarted.")
}

func Stop() {
	brew.StopBrewService("dnsmasq")
}

func Configure(config config.NovusConfig) bool {
	// update main dnsmasq config
	updated := updateDNSMasqConfig()

	// create the DNS resolver directory
	// https://www.manpagez.com/man/5/resolver/
	if _, err := exec.Command("sudo", "mkdir", "-p", dnsResolverDir).Output(); err != nil {
		logger.Errorf("Failed to create %s directory: %v\n", dnsResolverDir, err)
		os.Exit(1)
	}

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
		logger.Checkf("DNSMasq configuration updated.")
		return true
	} else {
		logger.Checkf("DNSMasq config is up to date.")
		return false
	}
}

func updateDNSMasqConfig() bool {
	// open file
	confFile, err := os.ReadFile(dnsmasqConfFile)
	if err != nil {
		logger.Errorf("Failed to read DNSMasq config file: %v.\n", err)
		os.Exit(1)
	}

	// remove comment "#"
	updatedConf := []byte(strings.Replace(
		string(confFile),
		fmt.Sprintf("#conf-dir=%s/etc/dnsmasq.d/,*.conf", brew.BrewPath),
		fmt.Sprintf("conf-dir=%s/etc/dnsmasq.d/,*.conf", brew.BrewPath),
		1,
	))

	// if the config differs (there was an actual change), write the changes
	if !bytes.Equal(confFile, updatedConf) {
		err = os.WriteFile(dnsmasqConfFile, []byte(updatedConf), 0644)
		if err != nil {
			logger.Errorf("Failed to write DNSMasq config file: %v.\n", err)
			os.Exit(1)
		}

		logger.Debugf("DNSMasq config has been updated.")

		return true
	} else {
		logger.Debugf("DNSMasq config is up to date.")

		return false
	}
}

func createDNSMasqDomainConfig(domain string) bool {
	configPath := fmt.Sprintf("%s/etc/dnsmasq.d/%s.conf", brew.BrewPath, domain)

	// first check if the file already exists
	if out, _ := os.Stat(configPath); out != nil {
		logger.Debugf("DNSMasq [*.%s]: Domain config already exists [%s].", domain, configPath)

		return false
	}

	logger.Debugf("[*.%s] Creating DNSMasq domain config.", domain)

	// prepare the configuration
	configContent := []byte(
		fmt.Sprintf("address=/%s/127.0.0.1", domain),
	)

	// create a configuration file
	err := os.WriteFile(configPath, configContent, 0644)
	if err != nil {
		logger.Errorf("DNSMasq [*.%s]: Failed to create domain config: %v.\n", domain, err)
		os.Exit(1)
	} else {
		logger.Debugf("DNSMasq [*.%s]: Domain config saved [%s].", domain, configPath)
	}

	return true
}

func registerDomainResolver(domain string) bool {
	configPath := fmt.Sprintf("%s/%s", dnsResolverDir, domain)

	// first check if the file already exists
	if out, _ := os.Stat(configPath); out != nil {
		logger.Debugf("DNSMasq [*.%s]: Domain resolver already exists [%s].", domain, configPath)

		return false
	}

	logger.Debugf("DNSMasq [*.%s]: Creating domain resolver.", domain)

	// create a configuration file
	configContent := []byte("nameserver 127.0.0.1\n")

	if _, err := exec.Command("sudo", "touch", configPath).Output(); err != nil {
		logger.Errorf("Failed to create %s file: %v\n", configPath, err)
		os.Exit(1)
	}

	usr, _ := user.Current()
	if _, err := exec.Command("sudo", "chown", usr.Username, configPath).Output(); err != nil {
		logger.Errorf("Failed to create %s file: %v\n", configPath, err)
		os.Exit(1)
	}

	err := os.WriteFile(configPath, configContent, 0644)
	if err != nil {
		logger.Errorf("DNSMasq [*.%s]: Failed to create domain resolver: %v.\n", domain, err)
		os.Exit(1)
	} else {
		logger.Debugf("DNSMasq [*.%s]: Domain resolver saved [%s].", domain, configPath)
	}

	return true
}

func getDomains(conf config.NovusConfig) []string {
	var domains = make(map[string]bool)

	for _, route := range conf.Routes {
		urlParts := strings.Split(route.Url, ".")
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
