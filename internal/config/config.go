package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/jozefcipa/novus/internal/fs"
	"github.com/jozefcipa/novus/internal/logger"
	"github.com/jozefcipa/novus/internal/shared"
	"gopkg.in/yaml.v3"
)

const configFileName = "novus.yml"

var AppName = "default"

type NovusConfig struct {
	Routes []shared.Route `yaml:"routes" validate:"required,unique_routes,dive"`
}

func (config *NovusConfig) validate() {
	logger.Debugf("Validating configuration file")

	validate := validator.New(validator.WithRequiredStructEnabled())

	// Register custom `unique_routes` rule
	shared.RegisterUniqueRoutesValidator(validate)

	err := validate.Struct(config)
	if err != nil {
		logger.Errorf("Configuration file contains errors.\n\n%s\n", err.(validator.ValidationErrors))
		os.Exit(1)
	}
}

func CreateDefaultConfigFile(appName string) (isDuplicateName bool) {
	// Read the config file template
	configTemplate := fs.ReadFileOrExit(filepath.Join(fs.AssetsDir, "novus.template.yml"))

	// Set app name in the config
	configTemplate = strings.Replace(configTemplate, "--APP_NAME--", appName, 1)

	// TODO: check in state file if appName is already being used
	// TODO: validate app name format
	isDuplicateName = false

	// Create a new config file
	fs.WriteFileOrExit(filepath.Join(fs.CurrentDir, configFileName), configTemplate)
	return
}

func Load() (NovusConfig, bool) {
	configPath := filepath.Join(fs.CurrentDir, configFileName)

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

	config.validate()

	return config, true
}
