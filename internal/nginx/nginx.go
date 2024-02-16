package nginx

import (
	"fmt"
	"strings"

	"github.com/jozefcipa/novus/internal/brew"
	"github.com/jozefcipa/novus/internal/config"
	"github.com/jozefcipa/novus/internal/fs"
	"github.com/jozefcipa/novus/internal/logger"
)

var nginxServersDir string

func init() {
	// /opt/homebrew/etc/nginx/nginx.conf - main config
	// /opt/homebrew/etc/nginx/servers/* - directory of loaded configs
	nginxServersDir = brew.BrewPath + "/etc/nginx/servers"
}

func Restart() {
	brew.RestartBrewService("nginx")
	logger.Checkf("Nginx: Service restarted.")
}

func Stop() {
	brew.StopBrewService("nginx")
}

func Configure(config config.NovusConfig) bool {
	appName := "default"
	nginxConf := readServerConfig(appName)
	newNginxConf := buildServerConfig(config)

	if nginxConf == "" || nginxConf != newNginxConf {
		writeServerConfig(appName, newNginxConf)
		logger.Debugf("Nginx: Built new config: \n\n%s", newNginxConf)
		logger.Checkf("Nginx: Routing configuration updated.")

		return true
	} else {
		logger.Checkf("Nginx: Configuration is up to date.")
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

func buildServerConfig(config config.NovusConfig) string {
	cwd := fs.GetCurrentDir()

	// Read template files
	configTemplate := fs.ReadFileOrExit("./assets/nginx/config.template.conf")
	serverConfigTemplate := fs.ReadFileOrExit("./assets/nginx/server.template.conf")

	// Iterate through all the routes and generate Nginx config
	serversSection := ""
	for _, route := range config.Routes {
		routeConfig := strings.Replace(serverConfigTemplate, "--SERVER_NAME--", route.Url, -1)
		routeConfig = strings.Replace(routeConfig, "--UPSTREAM_ADDR--", "http://"+route.Upstream, -1) // TODO: validate config to ensure http is either always presnet or never
		routeConfig = strings.Replace(routeConfig, "--ERRORS_DIR--", cwd+"/assets/nginx", -1)

		serversSection += routeConfig + "\n"
	}

	// Insert servers section into the main config
	serverConfig := strings.Replace(configTemplate, "--ERRORS_DIR--", cwd+"/assets/nginx", -1)
	serverConfig = strings.Replace(serverConfig, "--SERVERS_SECTION--", serversSection, 1)

	return serverConfig
}
