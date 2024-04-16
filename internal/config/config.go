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
	err := fs.Copy("./assets/novus.example.yml", "./"+configFileName)

	// if we weren't able to create a default config, throw an error
	if err != nil {
		logger.Errorf("%v\n", err)
		// we failed to create a default config file
		logger.Errorf("Example configuration file not found.\n%v", err)
		os.Exit(1)
	}
}

func Load(shouldCreateIfNotExists bool) NovusConfig {
	cwd := fs.GetCurrentDir()

	configFile, err := fs.ReadFile(filepath.Join(cwd, configFileName))
	if err != nil {
		logger.Errorf("No configuration file found.\n")

		if shouldCreateIfNotExists {
			createDefaultConfigFile()

			logger.Successf("Created a new config file.\n")
			logger.Infof("Open \"%s\" and define your routes.\n", configFileName)
		} else {
			logger.Messagef("Make sure %s exists in your directory or run `novus serve --create-config` to create a default configuration file.\n", configFileName)
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
