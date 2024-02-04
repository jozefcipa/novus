package config

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/jozefcipa/novus/internal/logger"
	"gopkg.in/yaml.v3"
)

const configFileName = "novus.yml"

type NovusConfig struct {
	Routes []struct {
		Url      string `yaml:"url"`
		Upstream string `yaml:"upstream"`
	}
}

func loadOrCreateFile() []byte {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current working directory: %v\n", err)
	}

	file, err := os.ReadFile(cwd + "/" + configFileName)
	if err != nil {
		logger.Debugf("Configuration file not found. Creating one now.")

		// if we didn't find config let's create one
		err = createDefaultConfig()

		// if we weren't able to create a default config, throw an error
		if err != nil {
			logger.Errorf("%v\n", err)
			// we failed to create a default config file
			logger.Errorf("Configuration file not found. Make sure \"%s\" exists.\n", configFileName)
			os.Exit(1)
		}

		// exit now, to let the users update the config
		logger.Messagef("No configuration file found.\n")
		logger.Successf("Creating a new one...\n")
		logger.Infof("Open \"%s\" and define your routes.\n", configFileName)
		os.Exit(0)
	}

	logger.Checkf("Configuration file found.")

	return file
}

func createDefaultConfig() error {
	srcFile, err := os.Open("./assets/config.default.yml")
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(fmt.Sprintf("./%s", configFileName))
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}

func (config *NovusConfig) validate() {
	logger.Debugf("Validating configuration file")
	// TODO: make sure the loaded file is in the expected format
	// https://github.com/go-playground/validator
	// Example: https://github.com/go-playground/validator/blob/master/_examples/simple/main.go
}

func Load() NovusConfig {
	file := loadOrCreateFile()

	config := NovusConfig{}

	err := yaml.Unmarshal(file, &config)
	if err != nil {
		logger.Errorf("Failed to parse the config file: %v", err)
		os.Exit(1)
	}

	config.validate()

	return config
}
