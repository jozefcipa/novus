package dnsmasq

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/jozefcipa/novus/internal/fs"
	"github.com/jozefcipa/novus/internal/homebrew"
	"github.com/jozefcipa/novus/internal/logger"
)

var dnsmasqConfFile string

// On some systems, port 53 might be already used by another DNS or some other service (e.g. PaloAltos GlobalProtect VPN),
// therefore we default to a different port
// However, if this port is also used, user will be prompted to provide an alternative port
const DefaultPort = "5053"

func init() {
	dnsmasqConfFile = filepath.Join(homebrew.HomebrewPrefix, "/etc/dnsmasq.conf")
}

func Restart() {
	dnsMasqLoader := logger.Loadingf("DNSMasq restarting")
	homebrew.RestartService("dnsmasq")

	// Check if the restart was successful
	isDNSMasqRunning := IsRunning()
	if !isDNSMasqRunning {
		dnsMasqLoader.Errorf("Failed to restart DNSMasq.")
		logger.Hintf("Try running \"brew services info dnsmasq --json\" for more info.")
		os.Exit(1)
	}

	dnsMasqLoader.Checkf("DNSMasq restarted")
}

func Stop() {
	nginxLoader := logger.Loadingf("Stopping DNSMasq")
	homebrew.StopService("dnsmasq")
	nginxLoader.Infof("🚫 DNSMasq stopped")
}

func IsRunning() bool {
	return homebrew.IsServiceRunning("dnsmasq")
}

func Configure(dnsPort string) bool {
	if dnsPort == "" {
		logger.Errorf("Called dnsmasq.Configure() with empty port")
		os.Exit(1)
	}

	// Open DNSMasq configuration file
	logger.Debugf("DNSMasq: Reading configuration file [%s]", dnsmasqConfFile)
	confFile := string(fs.ReadFileOrExit(dnsmasqConfFile))

	// Enable reading DNSMasq configurations from /etc/dnsmasq.d/* directory
	updatedConf := strings.Replace(
		confFile,
		fmt.Sprintf("#conf-dir=%s/etc/dnsmasq.d/,*.conf", homebrew.HomebrewPrefix),
		fmt.Sprintf("conf-dir=%s/etc/dnsmasq.d/,*.conf", homebrew.HomebrewPrefix),
		1,
	)

	// Enable alternative listening port
	// (matches both "#port=5353" and "port=1234" (any number))
	re := regexp.MustCompile(`#?port=\d+`)
	updatedConf = re.ReplaceAllString(confFile, fmt.Sprintf("port=%s", dnsPort))

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
	configPath := fmt.Sprintf(filepath.Join(homebrew.HomebrewPrefix, "/etc/dnsmasq.d/%s.conf"), tld)

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
