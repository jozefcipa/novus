package nginx

import (
	"fmt"
	"strings"

	"github.com/jozefcipa/novus/internal/brew"
	"github.com/jozefcipa/novus/internal/config"
	"github.com/jozefcipa/novus/internal/fs"
	"github.com/jozefcipa/novus/internal/logger"
	"github.com/jozefcipa/novus/internal/novus"
	"github.com/jozefcipa/novus/internal/shared"
)

var nginxServersDir string

func init() {
	// /opt/homebrew/etc/nginx/nginx.conf - main config
	// /opt/homebrew/etc/nginx/servers/* - directory of loaded configs
	nginxServersDir = brew.BrewPath + "/etc/nginx/servers"
}

func Restart() {
	brew.RestartService("nginx")
	logger.Checkf("Nginx: Service restarted.")
}

func Stop() {
	brew.StopService("nginx")
}

func IsRunning() bool {
	return brew.IsServiceRunning("nginx")
}

func Configure(config config.NovusConfig, sslCerts shared.DomainCertificates) bool {
	appName := "default"
	nginxConf := readServerConfig(appName)
	newNginxConf := buildServerConfig(config, sslCerts)

	if nginxConf == "" || nginxConf != newNginxConf {
		writeServerConfig(appName, newNginxConf)
		logger.Debugf("Nginx: Built new config: \n\n%s", newNginxConf)
		logger.Checkf("Nginx: Routing configuration updated.")

		return true
	} else {
		logger.Checkf("Nginx configuration is up to date.")
		return false
	}
}

func readServerConfig(app string) string {
	path := fmt.Sprintf("%s/novus-%s.conf", nginxServersDir, app)
	logger.Debugf("Nginx: Reading config %s", path)

	// If file doesn't exist (an error is thrown) just return an empty string and we'll create a new config later
	file, _ := fs.ReadFile(path)

	return file
}

func writeServerConfig(app string, serverConfig string) {
	path := fmt.Sprintf("%s/novus-%s.conf", nginxServersDir, app)
	logger.Debugf("Nginx: Writing config [%s]", path)

	fs.WriteFileOrExit(path, serverConfig)
}

func buildServerConfig(novusConfig config.NovusConfig, sslCerts shared.DomainCertificates) string {
	cwd := fs.GetCurrentDir()

	// Read template files
	configTemplate := fs.ReadFileOrExit("./assets/nginx/config.template.conf")
	serverConfigTemplate := fs.ReadFileOrExit("./assets/nginx/server.template.conf")

	// update routes in state
	appState := novus.GetState()
	appState.Routes = novusConfig.Routes

	// Iterate through all the routes and generate Nginx config
	serversSection := ""
	for _, route := range novusConfig.Routes {
		sslCert := sslCerts[route.Domain]

		// create Nginx server block
		routeConfig := strings.Replace(serverConfigTemplate, "--SERVER_NAME--", route.Domain, -1)
		routeConfig = strings.Replace(routeConfig, "--UPSTREAM_ADDR--", route.Upstream, -1)
		routeConfig = strings.Replace(routeConfig, "--ERRORS_DIR--", cwd+"/assets/nginx", -1)
		routeConfig = strings.Replace(routeConfig, "--SSL_CERT_PATH--", sslCert.CertFilePath, 1)
		routeConfig = strings.Replace(routeConfig, "--SSL_KEY_PATH--", sslCert.KeyFilePath, 1)

		serversSection += routeConfig + "\n"
	}

	// Insert servers section into the main config
	serverConfig := strings.Replace(configTemplate, "--ERRORS_DIR--", cwd+"/assets/nginx", -1)
	serverConfig = strings.Replace(serverConfig, "--SERVERS_SECTION--", serversSection, 1)

	return serverConfig
}
