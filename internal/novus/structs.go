package novus

import (
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/jozefcipa/novus/internal/logger"
	"github.com/jozefcipa/novus/internal/shared"
)

type AppStatus string

const (
	APP_ACTIVE AppStatus = "active"
	APP_PAUSED AppStatus = "paused"
)

type AppState struct {
	Directory       string                    `json:"directory" validate:"required,dirpath"`
	Status          AppStatus                 `json:"appStatus" validate:"required"`
	SSLCertificates shared.DomainCertificates `json:"sslCertificates"`
	Routes          []shared.Route            `json:"routes" validate:"required,dive"`
}

type DnsFiles struct {
	DnsMasqConfig string `json:"dnsMasqConfig" validate:"required,dirpath"`
	DnsResolver   string `json:"dnsResolver" validate:"required,dirpath"`
}

type NovusState struct {
	// Track files that we create for DNS
	// As we write into shared directory, we can later on only delete files that we're sure have been created by us
	// e.g. /etc/resolver directory
	DnsFiles map[string]*DnsFiles `json:"dnsFiles" validate:"required"`
	// State for each of the apps
	Apps map[string]*AppState `json:"apps" validate:"required"`
}

func (state *NovusState) validate() {
	logger.Debugf("Validating state file")

	validate := validator.New(validator.WithRequiredStructEnabled())

	for _, appState := range state.Apps {
		err := validate.Struct(appState)
		if err != nil {
			logger.Errorf("Novus state file is corrupted.\n\n%s", err.(validator.ValidationErrors))
			os.Exit(1)
		}

		for _, sslCerts := range appState.SSLCertificates {
			err := validate.Struct(sslCerts)
			if err != nil {
				logger.Errorf("Novus state file is corrupted.\n\n%s", err.(validator.ValidationErrors))
				os.Exit(1)
			}
		}
	}
}

func (state *NovusState) GetAllApps() map[string]*AppState {
	return state.Apps
}

func (state *NovusState) GetActiveApps() map[string]*AppState {
	activeApps := map[string]*AppState{}

	for appName, appState := range state.Apps {
		if appState.Status == APP_ACTIVE {
			activeApps[appName] = appState
		}
	}

	return activeApps
}
