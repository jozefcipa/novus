package config

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/jozefcipa/novus/internal/fs"
	"github.com/jozefcipa/novus/internal/logger"
	"github.com/jozefcipa/novus/internal/shared"
	"gopkg.in/yaml.v3"
)

const configFileName = "novus.yml"

var AppName = ""

type NovusConfig struct {
	AppName string         `yaml:"appName" validate:"required"`
	Routes  []shared.Route `yaml:"routes" validate:"required,unique_routes,dive"`
}

func (config *NovusConfig) validate() {
	logger.Debugf("Validating configuration file")

	validate := validator.New(validator.WithRequiredStructEnabled())

	// Register custom `unique_routes` rule
	shared.RegisterUniqueRoutesValidator(validate)

	if err := validate.Struct(config); err != nil {
		logger.Errorf("Configuration file contains errors.\n\n%s", err.(validator.ValidationErrors))
		os.Exit(1)
	}

	if err := validateAppName(config.AppName); err != nil {
		logger.Errorf("Configuration file contains errors.\n\n%s", err.Error())
		os.Exit(1)
	}
}

func validateAppName(appName string) error {
	isValid, _ := regexp.MatchString("^[A-Za-z0-9-_]+$", appName)
	if !isValid {
		return fmt.Errorf("Invalid app name. Only alphanumeric characters are allowed.")
	}

	return nil
}

func CreateDefaultConfigFile(appName string) error {
	// Read the config file template
	configTemplate := fs.ReadFileOrExit(filepath.Join(fs.AssetsDir, "novus.template.yml"))

	// Set app name in the config
	configTemplate = strings.Replace(configTemplate, "--APP_NAME--", appName, 1)

	if err := validateAppName(appName); err != nil {
		return err
	}

	// TODO: check in state file if appName is already being used
	isDuplicateName := false
	if isDuplicateName {
		return fmt.Errorf("You already have a configuration with the name \"%s\"", appName)
	}

	// Create a new config file
	fs.WriteFileOrExit(filepath.Join(fs.CurrentDir, configFileName), configTemplate)
	return nil
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

	AppName = config.AppName

	return config, true
}
