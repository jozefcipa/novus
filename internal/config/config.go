package config

import (
	"os"
	"path/filepath"

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

	// register custom `unique_routes` rule
	shared.RegisterUniqueRoutesValidator(validate)

	err := validate.Struct(config)
	if err != nil {
		logger.Errorf("Configuration file contains errors.\n\n%s\n", err.(validator.ValidationErrors))
		os.Exit(1)
	}
}

func createDefaultConfigFile() {
	// if we didn't find config file, let's create the default one
	err := fs.Copy(
		filepath.Join(fs.AssetsDir, "novus.example.yml"),
		filepath.Join(fs.CurrentDir, configFileName),
	)

	// if we weren't able to create a default config, throw an error
	if err != nil {
		// we failed to create a default config file
		logger.Errorf("Example configuration file not found: %v", err)
		os.Exit(1)
	}
}

func Load(shouldCreateIfNotExists bool) NovusConfig {
	configPath := filepath.Join(fs.CurrentDir, configFileName)
	configFile, err := fs.ReadFile(configPath)
	logger.Debugf("Loading configuration file [%s]", configPath)

	if err != nil {
		if shouldCreateIfNotExists {
			createDefaultConfigFile()

			logger.Successf("Config file created.\n")
			logger.Infof("Open \"%s\" and define your routes.\n", configFileName)
		} else {
			logger.Errorf("Configuration file not found.\n")
			logger.Messagef("Run \"novus serve --create-config\" to initialize configuration.\n")
		}

		// exit now so the user can either update the generated config or call the command again properly
		os.Exit(0)
	}

	config := NovusConfig{}

	err = yaml.Unmarshal([]byte(configFile), &config)
	if err != nil {
		logger.Errorf("Failed to parse the config file: %v", err)
		os.Exit(1)
	}

	config.validate()

	return config
}
