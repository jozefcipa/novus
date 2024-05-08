package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/jozefcipa/novus/internal/fs"
	"github.com/jozefcipa/novus/internal/logger"
	"github.com/jozefcipa/novus/internal/shared"
	"gopkg.in/yaml.v3"
)

const ConfigFileName = "novus.yml"

var appName = ""

type NovusConfig struct {
	AppName string         `yaml:"appName" validate:"required"`
	Routes  []shared.Route `yaml:"routes" validate:"required,unique_routes,dive"`
}

func SetAppName(name string) {
	logger.Debugf("Setting app [app=%s]", name)
	appName = name
}

func AppName() string {
	if appName != "" {
		return appName
	}

	// This should not happen normally,
	// but let's throw an error if the program tries to access config.AppName() when not set
	logger.Errorf("[Internal error]: No app set, make sure to call `config.SetAppName()`")
	os.Exit(1)
	return ""
}

func WriteDefaultFile(appName string) {
	// Read the config file template
	configTemplate := fs.ReadFileOrExit(filepath.Join(fs.AssetsDir, "novus.template.yml"))

	// Set app name in the config
	configTemplate = strings.Replace(configTemplate, "--APP_NAME--", appName, 1)

	// Create a new config file
	fs.WriteFileOrExit(filepath.Join(fs.CurrentDir, ConfigFileName), configTemplate)
}

func LoadFile() (NovusConfig, bool) {
	configPath := filepath.Join(fs.CurrentDir, ConfigFileName)

	logger.Debugf("Loading configuration file [%s]", configPath)
	configFile, err := fs.ReadFile(configPath)
	if err != nil {
		return NovusConfig{}, false
	}

	config := NovusConfig{}
	err = yaml.Unmarshal([]byte(configFile), &config)
	if err != nil {
		logger.Errorf("Failed to parse the config file: %v", err)
		os.Exit(1)
	}

	return config, true
}
