package dnsmasq

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/jozefcipa/novus/internal/brew"
	"github.com/jozefcipa/novus/internal/fs"
	"github.com/jozefcipa/novus/internal/logger"
)

var dnsmasqConfFile string

func init() {
	dnsmasqConfFile = filepath.Join(brew.BrewPath, "/etc/dnsmasq.conf")
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

func Configure() bool {
	// open DNSMasq configuration file
	logger.Debugf("DNSMasq: Reading configuration file [%s]", dnsmasqConfFile)
	confFile := fs.ReadFileOrExit(dnsmasqConfFile)

	// Enable reading DNSMasq configurations from /etc/dnsmasq.d/* directory
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

func CreateTLDConfig(tld string) (bool, string) {
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
