package dnsmasq

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jozefcipa/novus/internal/brew"
	"github.com/jozefcipa/novus/internal/fs"
	"github.com/jozefcipa/novus/internal/logger"
	"github.com/jozefcipa/novus/internal/net"
)

var dnsmasqConfFile string

// On some systems, port 53 might be already used by another DNS or some other service (e.g. PaloAltos GlobalProtect VPN),
// therefore let's use a different port to play it safe
const Port = "5053"

func init() {
	dnsmasqConfFile = filepath.Join(brew.BrewPath, "/etc/dnsmasq.conf")
}

func Restart() {
	logger.Infof("🔄 DNSMasq restarting...")

	brew.RestartService("dnsmasq")

	// Check if the restart was successful
	isDNSMasqRunning := IsRunning()
	if !isDNSMasqRunning {
		logger.Errorf("Failed to restart DNSMasq.")
		logger.Hintf("Try running \"brew services info dnsmasq --json\" for more info.")
		os.Exit(1)
	}
	logger.Checkf("DNSMasq restarted")
}

func Stop() {
	brew.StopService("dnsmasq")
}

func IsRunning() bool {
	return brew.IsServiceRunning("dnsmasq")
}

func EnsurePortAvailable(portsUsage net.PortUsage) {
	if portUsedBy, isUsed := portsUsage[Port]; isUsed && portUsedBy != "dnsmasq" {
		logger.Errorf("Cannot start DNSMasq: Port %s is already used by '%s'", Port, portUsedBy)
		os.Exit(1)
	}
}

func Configure() bool {
	// Open DNSMasq configuration file
	logger.Debugf("DNSMasq: Reading configuration file [%s]", dnsmasqConfFile)
	confFile := string(fs.ReadFileOrExit(dnsmasqConfFile))

	// Enable reading DNSMasq configurations from /etc/dnsmasq.d/* directory
	updatedConf := strings.Replace(
		confFile,
		fmt.Sprintf("#conf-dir=%s/etc/dnsmasq.d/,*.conf", brew.BrewPath),
		fmt.Sprintf("conf-dir=%s/etc/dnsmasq.d/,*.conf", brew.BrewPath),
		1,
	)
	// Enable alternative listening port
	updatedConf = strings.Replace(
		confFile,
		"#port=5353",
		fmt.Sprintf("port=%s", Port),
		1,
	)

	// If the config differs (there was an actual change), write the changes
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

	// First check if the file already exists
	if confExists := fs.FileExists(configPath); confExists {
		logger.Debugf("DNSMasq [*.%s]: Domain config already exists [%s]", tld, configPath)

		return false, configPath
	}

	logger.Debugf("DNSMasq [*.%s]: Creating domain config", tld)

	// Prepare the configuration
	configContent := fmt.Sprintf("address=/%s/127.0.0.1", tld)

	// Create a configuration file
	fs.WriteFileOrExit(configPath, configContent)
	logger.Debugf("DNSMasq [*.%s]: Domain config saved [%s]", tld, configPath)

	return true, configPath
}
