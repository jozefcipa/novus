package nginx

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

var NginxServersDir string

func init() {
	// /opt/homebrew/etc/nginx/nginx.conf - main config
	// /opt/homebrew/etc/nginx/servers/* - directory of loaded configs
	NginxServersDir = filepath.Join(brew.BrewPath, "/etc/nginx/servers")
}

func Restart() {
	brew.RestartService("nginx")
}

func Stop() {
	brew.StopService("nginx")
}

func IsRunning() bool {
	return brew.IsServiceRunning("nginx")
}

func Configure(novusConf config.NovusConfig, sslCerts shared.DomainCertificates) bool {
	nginxConf := readServerConfig(config.AppName)
	newNginxConf := buildServerConfig(novusConf, sslCerts)

	if nginxConf == "" || nginxConf != newNginxConf {
		logger.Debugf("Generated new Nginx config: \n\n%s", newNginxConf)
		writeServerConfig(config.AppName, newNginxConf)
		logger.Checkf("Nginx configuration updated")
		return true
	} else {
		logger.Checkf("Nginx configuration is up to date")
		return false
	}
}

func readServerConfig(app string) string {
	path := filepath.Join(NginxServersDir, fmt.Sprintf("novus-%s.conf", app))
	logger.Debugf("Reading Nginx config [%s]", path)

	// If file doesn't exist (an error is thrown) just return an empty string and we'll create a new config later
	file, _ := fs.ReadFile(path)

	return file
}

func writeServerConfig(app string, serverConfig string) {
	path := filepath.Join(NginxServersDir, fmt.Sprintf("novus-%s.conf", app))
	logger.Debugf("Updating Nginx config [%s]", path)

	fs.WriteFileOrExit(path, serverConfig)
}

func buildServerConfig(novusConfig config.NovusConfig, sslCerts shared.DomainCertificates) string {
	// Read template files
	configTemplate := fs.ReadFileOrExit(filepath.Join(fs.AssetsDir, "nginx/config.template.conf"))
	serverConfigTemplate := fs.ReadFileOrExit(filepath.Join(fs.AssetsDir, "nginx/server.template.conf"))

	// update routes in state
	appState, _ := novus.GetAppState()
	appState.Routes = novusConfig.Routes

	// Iterate through all the routes and generate Nginx config
	serversSection := ""
	for _, route := range novusConfig.Routes {
		sslCert := sslCerts[route.Domain]

		// create Nginx server block
		routeConfig := strings.Replace(serverConfigTemplate, "--SERVER_NAME--", route.Domain, -1)
		routeConfig = strings.Replace(routeConfig, "--UPSTREAM_ADDR--", route.Upstream, -1)
		routeConfig = strings.Replace(routeConfig, "--ERRORS_DIR--", filepath.Join(fs.AssetsDir, "nginx"), -1)
		routeConfig = strings.Replace(routeConfig, "--SSL_CERT_PATH--", sslCert.CertFilePath, 1)
		routeConfig = strings.Replace(routeConfig, "--SSL_KEY_PATH--", sslCert.KeyFilePath, 1)

		serversSection += routeConfig + "\n"
	}

	// Insert servers section into the main config
	serverConfig := strings.Replace(configTemplate, "--ERRORS_DIR--", filepath.Join(fs.AssetsDir, "nginx"), -1)
	serverConfig = strings.Replace(serverConfig, "--SERVERS_SECTION--", serversSection, 1)

	return serverConfig
}
