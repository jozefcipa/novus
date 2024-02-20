package config

import (
	"os"

	"github.com/jozefcipa/novus/internal/fs"
	"github.com/jozefcipa/novus/internal/logger"
	"github.com/jozefcipa/novus/internal/shared"
	"gopkg.in/yaml.v3"
)

const configFileName = "novus.yml"

var AppName = "default"

type NovusConfig struct {
	Routes []shared.Route
}

func loadOrCreateFile() []byte {
	cwd := fs.GetCurrentDir()

	file, err := os.ReadFile(cwd + "/" + configFileName)
	if err == nil {
		logger.Checkf("Configuration file found.")
		return file
	}

	logger.Messagef("No configuration file found.\n")

	// if we didn't find config file, let's create the default one
	err = fs.Copy("./assets/novus.example.yml", "./"+configFileName)

	// if we weren't able to create a default config, throw an error
	if err != nil {
		logger.Errorf("%v\n", err)
		// we failed to create a default config file
		logger.Errorf("Configuration file not found. Make sure \"%s\" exists.\n", configFileName)
		os.Exit(1)
	}

	// exit now to let the users update the config
	logger.Successf("Creating a new config file.\n")
	logger.Infof("Open \"%s\" and define your routes.\n", configFileName)
	os.Exit(0)

	return []byte{}
}

func (config *NovusConfig) validate() {
	logger.Debugf("Validating configuration file")
	// TODO: make sure the loaded file is in the expected format
	// https://github.com/go-playground/validator
	// Example: https://github.com/go-playground/validator/blob/master/_examples/simple/main.go
}

func Load() NovusConfig {
	configFile := loadOrCreateFile()

	config := NovusConfig{}

	err := yaml.Unmarshal(configFile, &config)
	if err != nil {
		logger.Errorf("Failed to parse the config file: %v", err)
		os.Exit(1)
	}

	config.validate()

	return config
}
