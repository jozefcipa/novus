package config_manager

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/jozefcipa/novus/internal/config"
	"github.com/jozefcipa/novus/internal/diff_manager"
	"github.com/jozefcipa/novus/internal/fs"
	"github.com/jozefcipa/novus/internal/logger"
	"github.com/jozefcipa/novus/internal/novus"
	"github.com/jozefcipa/novus/internal/shared"
)

var validationMessages = map[string]string{
	"required":      "Field '%s' is required",
	"url":           "Field '%s' is not a valid URL",
	"fqdn":          "Field '%s' is not a valid FQDN",
	"existing_tld":  "Field '%s' contains an existing TLD domain.",
	"unique_routes": "Field '%s' contains duplicate route definitions.",
}

func LoadConfiguration() (config.NovusConfig, bool) {
	conf, exists := config.LoadFile()
	return conf, exists
}

func LoadConfigurationFromState(appName string, novusState novus.NovusState) config.NovusConfig {
	return config.NovusConfig{
		AppName: appName,
		Routes:  novusState.Apps[appName].Routes,
	}
}

func getConfigFieldPath(structNamespace string) string {
	pathKeys := strings.Split(structNamespace, ".")[1:] // remove first item as it is the name of the config struct (NovusConfig)

	for i, key := range pathKeys {
		pathKeys[i] = shared.LowerFirst(key)
	}

	return strings.Join(pathKeys, ".")
}

func ValidateConfig(conf config.NovusConfig, novusState novus.NovusState) {
	// Validate configuration
	if err := validateConfigSyntax(conf); err != nil {

		errors := []string{}
		for _, err := range err.(validator.ValidationErrors) {
			validationRule := err.Tag()
			path := getConfigFieldPath(err.StructNamespace())

			// Default error message
			errorMessage := err.Error()

			// Check if we have custom error message defined
			if customErrorMesssage, ok := validationMessages[validationRule]; ok {
				errorMessage = fmt.Sprintf(customErrorMesssage, path)
			}

			errors = append(errors, errorMessage)
		}

		logger.Errorf("Configuration file contains errors:\n   %s", strings.Join(errors, "\n   "))

		os.Exit(1)
	}

	// Validate app name syntax and whether it is unique across apps
	if err := validateConfigAppName(conf.AppName, novusState); err != nil {
		logger.Errorf(err.Error())
		os.Exit(1)
	}
	config.SetAppName(conf.AppName)

	// Check if the config contains domains that are already registered in another app
	if err := checkForDuplicateDomains(conf, novusState); err != nil {
		logger.Errorf(err.Error())
		if err, ok := err.(*diff_manager.DuplicateDomainError); ok {
			logger.Hintf(
				"Use a different domain name or temporarily stop %[1]s by running \"novus pause %[1]s\"",
				err.OriginalAppWithDomain,
			)
		}
		os.Exit(1)
	}
}

func ConfigFileExists() bool {
	return config.ConfigFileExists()
}

func CreateNewConfiguration(appName string, novusState novus.NovusState) error {
	if err := validateConfigAppName(appName, novusState); err != nil {
		return err
	}

	config.WriteDefaultFile(appName)
	return nil
}

func validateConfigSyntax(conf config.NovusConfig) error {
	logger.Debugf("Validating configuration file syntax")

	validate := validator.New(validator.WithRequiredStructEnabled())

	// Register custom `unique_routes` rule
	shared.RegisterUniqueRoutesValidator(validate)
	// Register custom `existing_tld` rule
	shared.RegisterNonExistentTLDValidator(validate)

	return validate.Struct(conf)
}

func validateConfigAppName(appName string, novusState novus.NovusState) error {
	logger.Debugf("Validating configuration file app name [%s]", appName)

	isValid, _ := regexp.MatchString("^[A-Za-z0-9-_]+$", appName)
	if !isValid {
		return fmt.Errorf("Invalid app name. Only alphanumeric characters are allowed.")
	}

	if appName == novus.NovusInternalAppName {
		return fmt.Errorf("Reserved app name. This app is used internally by Novus.")
	}

	// Check in state file if appName is already being used
	for appNameFromConfig, appConfig := range novusState.GetActiveApps() {
		if appNameFromConfig == appName && appConfig.Directory != fs.CurrentDir {
			return fmt.Errorf("App \"%s\" is already defined in a different directory (%s)", appName, appConfig.Directory)
		}
	}

	return nil
}

func checkForDuplicateDomains(conf config.NovusConfig, novusState novus.NovusState) error {
	logger.Debugf("Checking for duplicate domains across apps")

	// pick all existing apps except the current one (based on the config)
	otherApps := map[string]novus.AppState{}
	for appName, appState := range novusState.GetActiveApps() {
		if appName != conf.AppName {
			otherApps[appName] = *appState
		}
	}
	return diff_manager.DetectDuplicateDomains(otherApps, conf.Routes)
}
