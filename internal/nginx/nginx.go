package nginx

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jozefcipa/novus/internal/brew"
	"github.com/jozefcipa/novus/internal/config"
	"github.com/jozefcipa/novus/internal/logger"
)

var nginxServersDir string

func init() {
	// /opt/homebrew/etc/nginx/nginx.cosnf - main config
	// /opt/homebrew/etc/nginx/servers/* - directory of loaded configs
	nginxServersDir = brew.BrewPath + "/etc/nginx/servers"
}

func Restart() {
	brew.RestartBrewService("nginx")
	logger.Checkf("Nginx restarted.")
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
		logger.Debugf("New nginx config: \n\n%s", newNginxConf)
		logger.Checkf("Routing configuration updated.")

		return true
	} else {
		logger.Checkf("Nginx config is up to date.")
		return false
	}
}

func readServerConfig(app string) string {
	path := fmt.Sprintf("%s/novus-%s.conf", nginxServersDir, app)
	logger.Debugf("Reading %s", path)

	file, err := os.ReadFile(path)
	if err != nil {
		logger.Debugf("Error ocurred %v", err)
		return ""
	}

	return string(file[:])
}

func writeServerConfig(app string, serverConfig string) {
	data := []byte(serverConfig)
	path := fmt.Sprintf("%s/novus-%s.conf", nginxServersDir, app)
	logger.Debugf("Writing %s", path)

	err := os.WriteFile(path, data, 0644)
	if err != nil {
		log.Fatalf("Failed to write into Nginx config: %v", err)
	}
}

func buildServerConfig(config config.NovusConfig) string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current working directory: %v\n", err)
	}

	configFile, err := os.ReadFile("./assets/nginx/config.template.conf")
	if err != nil {
		log.Fatalf("Failed to read config.template.conf file: %v", err)
	}
	serverSectionFile, err := os.ReadFile("./assets/nginx/server.template.conf")
	if err != nil {
		log.Fatalf("Failed to read server.template.conf file: %v", err)
	}

	configTemplate := string(configFile)
	serverConfigTemplate := string(serverSectionFile)

	// Iterate through all the routes and generate Nginx config
	serversSection := ""
	for _, route := range config.Routes {
		routeConfig := strings.Replace(serverConfigTemplate, "--SERVER_NAME--", route.Url, -1)
		routeConfig = strings.Replace(routeConfig, "--PORT--", "80", -1)
		routeConfig = strings.Replace(routeConfig, "--UPSTREAM_ADDR--", "http://"+route.Upstream, -1) // TODO: validate config to ensure http is either always presnet or never
		routeConfig = strings.Replace(routeConfig, "--ERRORS_DIR--", cwd+"/assets/nginx", -1)

		serversSection += routeConfig + "\n"
	}

	// Insert servers section into the main config
	serverConfig := strings.Replace(configTemplate, "--SERVERS_SECTION--", serversSection, 1)

	return serverConfig
}
