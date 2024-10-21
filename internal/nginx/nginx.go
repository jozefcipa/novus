package nginx

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jozefcipa/novus/internal/brew"
	"github.com/jozefcipa/novus/internal/config"
	"github.com/jozefcipa/novus/internal/fs"
	"github.com/jozefcipa/novus/internal/logger"
	"github.com/jozefcipa/novus/internal/novus"
	"github.com/jozefcipa/novus/internal/shared"
)

var NginxServersDir string
var currentTime time.Time
var fileHeader string

func init() {
	// /opt/homebrew/etc/nginx/nginx.conf - main config
	// /opt/homebrew/etc/nginx/servers/* - directory of loaded configs
	NginxServersDir = filepath.Join(brew.BrewPath, "/etc/nginx/servers")

	currentTime = time.Now()
	fileHeader = `##################################################################
# DO NOT EDIT THIS FILE!!!
# ------------------------
# This configuration is auto-generated.
# Created at ` + currentTime.Format("2006-01-02 15:04") + ` by Novus
#################################################################

`
}

func Restart() {
	logger.Infof("🔄 Nginx restarting...")

	brew.RestartService("nginx")

	// Check if the restart was successful
	isNginxRunning := IsRunning()
	if !isNginxRunning {
		logger.Errorf("Failed to restart Nginx.")
		logger.Hintf("Try running one of the following commands for more info:\n- brew services info nginx --json\n- nginx -t")
		os.Exit(1)
	}
	logger.Checkf("Nginx restarted")
}

func Stop() {
	brew.StopService("nginx")
}

func IsRunning() bool {
	return brew.IsServiceRunning("nginx")
}

func Configure(novusConf config.NovusConfig, sslCerts shared.DomainCertificates, appState *novus.AppState) bool {
	// Create default server config if it doesn't exist
	nginxDefaultConf := readServerConfig(getDefaultConfigName())

	defaultConfig := fileHeader
	defaultConfig += fs.ReadFileOrExit(filepath.Join(fs.AssetsDir, "nginx/default-server.template.conf"))
	defaultConfig = strings.Replace(defaultConfig, "--NOVUS_HTML_DIR--", filepath.Join(fs.AssetsDir, "nginx/html"), -1)
	defaultConfig = strings.Replace(defaultConfig, "--NOVUS_ASSETS_DIR--", filepath.Join(fs.AssetsDir, "nginx"), -1)
	defaultConfig = strings.Replace(defaultConfig, "--NOVUS_STATE_FILE_PATH--", novus.NovusStateFilePath, -1)
	defaultConfig = strings.Replace(defaultConfig, "--NOVUS_INTERNAL_SERVER_NAME--", novus.NovusInternalDomain, -1)
	defaultConfig = strings.Replace(defaultConfig, "--NOVUS_INDEX_SERVER_NAME--", novus.NovusIndexDomain, -1)

	novusInternalDomainSSL, ok := sslCerts[novus.NovusInternalDomain]
	if !ok {
		logger.Errorf("Internal domain %s not found in SSL certs config\n", novus.NovusInternalDomain)
		os.Exit(1)
	}

	defaultConfig = strings.Replace(defaultConfig, "--NOVUS_INTERNAL_SSL_CERT_PATH--", novusInternalDomainSSL.CertFilePath, -1)
	defaultConfig = strings.Replace(defaultConfig, "--NOVUS_INTERNAL_SSL_KEY_PATH--", novusInternalDomainSSL.KeyFilePath, -1)

	novusIndexDomainSSL, ok := sslCerts[novus.NovusIndexDomain]
	if !ok {
		logger.Errorf("Internal domain %s not found in SSL certs config\n", novus.NovusIndexDomain)
		os.Exit(1)
	}
	defaultConfig = strings.Replace(defaultConfig, "--NOVUS_INDEX_SSL_CERT_PATH--", novusIndexDomainSSL.CertFilePath, -1)
	defaultConfig = strings.Replace(defaultConfig, "--NOVUS_INDEX_SSL_KEY_PATH--", novusIndexDomainSSL.KeyFilePath, -1)

	if nginxDefaultConf == "" || nginxDefaultConf != defaultConfig {
		logger.Debugf("Generated default server Nginx config: \n\n%s", defaultConfig)
		writeServerConfig(getDefaultConfigName(), defaultConfig)
	}

	// Create application server config if it doesn't exist
	nginxAppConf := readServerConfig(getAppConfigName(config.AppName()))
	newNginxAppConf := buildServerConfig(novusConf, sslCerts, appState)

	if nginxAppConf == "" || nginxAppConf != newNginxAppConf {
		logger.Debugf("Generated application server Nginx config: \n\n%s", newNginxAppConf)
		writeServerConfig(getAppConfigName(config.AppName()), newNginxAppConf)
		logger.Checkf("Nginx configuration updated")
		return true
	} else {
		logger.Checkf("Nginx configuration is up to date")
		return false
	}
}

func RemoveConfiguration(appName string) {
	configFilePath := filepath.Join(NginxServersDir, getAppConfigName(appName))

	logger.Debugf("Removing application server Nginx config for app %s [%s]", appName, configFilePath)

	fs.DeleteFile(configFilePath)
}

func readServerConfig(fileName string) string {
	path := filepath.Join(NginxServersDir, fileName)
	logger.Debugf("Reading Nginx config [%s]", path)

	// If file doesn't exist (an error is thrown) just return an empty string and we'll create a new config later
	file, _ := fs.ReadFile(path)

	return file
}

func writeServerConfig(fileName string, serverConfig string) {
	path := filepath.Join(NginxServersDir, fileName)
	logger.Debugf("Updating Nginx config [%s]", path)

	fs.WriteFileOrExit(path, serverConfig)
}

func buildServerConfig(appConfig config.NovusConfig, sslCerts shared.DomainCertificates, appState *novus.AppState) string {
	// Read template file
	serverConfigTemplate := fs.ReadFileOrExit(filepath.Join(fs.AssetsDir, "nginx/server.template.conf"))

	// Update routes in state
	appState.Routes = appConfig.Routes

	// Iterate through all the routes and generate Nginx config
	serverConfig := fileHeader
	for _, route := range appConfig.Routes {
		sslCert := sslCerts[route.Domain]

		// Create Nginx server block
		routeConfig := strings.ReplaceAll(serverConfigTemplate, "--SERVER_NAME--", route.Domain)
		routeConfig = strings.ReplaceAll(routeConfig, "--UPSTREAM_ADDR--", route.Upstream)
		routeConfig = strings.ReplaceAll(routeConfig, "--NOVUS_HTML_DIR--", filepath.Join(fs.AssetsDir, "nginx/html"))
		routeConfig = strings.ReplaceAll(routeConfig, "--SSL_CERT_PATH--", sslCert.CertFilePath)
		routeConfig = strings.ReplaceAll(routeConfig, "--SSL_KEY_PATH--", sslCert.KeyFilePath)

		serverConfig += routeConfig + "\n"
	}

	return serverConfig
}

func getDefaultConfigName() string {
	return "novus-default.conf"
}

func getAppConfigName(appName string) string {
	return fmt.Sprintf("novus-app-%s.conf", appName)
}
