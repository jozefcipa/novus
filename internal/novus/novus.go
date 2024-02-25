package novus

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/go-playground/validator/v10"
	"github.com/jozefcipa/novus/internal/config"
	"github.com/jozefcipa/novus/internal/fs"
	"github.com/jozefcipa/novus/internal/logger"
	"github.com/jozefcipa/novus/internal/shared"
)

var NovusStateDir string

var novusStateFilePath string
var state NovusState

type AppState struct {
	Directory       string                    `json:"directory" validate:"required,dirpath"`
	SSLCertificates shared.DomainCertificates `json:"sslCertificates"`
	Routes          []shared.Route            `json:"routes" validate:"required,dive"`
}

type NovusState map[string]*AppState

func (config *NovusState) validate() {
	logger.Debugf("Validating state file")

	validate := validator.New(validator.WithRequiredStructEnabled())

	for _, appState := range *config {
		err := validate.Struct(appState)
		if err != nil {
			logger.Errorf("Novus state file is corrupted.\n\n%s\n", err.(validator.ValidationErrors))
			os.Exit(1)
		}

		for _, sslCerts := range appState.SSLCertificates {
			err := validate.Struct(sslCerts)
			if err != nil {
				logger.Errorf("Novus state file is corrupted.\n\n%s\n", err.(validator.ValidationErrors))
				os.Exit(1)
			}
		}
	}
}

func init() {
	// Create a directory ~/.novus
	// where we can store generated SSL certificates and application state
	NovusStateDir = filepath.Join(fs.UserHomeDir, ".novus")
	fs.MakeDirOrExit(NovusStateDir)
	novusStateFilePath = filepath.Join(NovusStateDir, "novus.json")
}

func initEmptyState() *AppState {
	return &AppState{
		Directory:       fs.GetCurrentDir(),
		SSLCertificates: shared.DomainCertificates{},
		Routes:          []shared.Route{},
	}
}

func LoadState() {
	file, err := fs.ReadFile(novusStateFilePath)
	// if there's an error, probably we didn't find the state, so initialize a new one
	if err != nil {
		logger.Debugf("State file not found. Creating a new one.")
		state = NovusState{
			"default": initEmptyState(),
		}
		return
	}

	if err := json.Unmarshal([]byte(file), &state); err != nil {
		logger.Errorf("Corrupted state file.\n%v\n", err)
		os.Exit(1)
	}

	state.validate()
}

func GetState() *AppState {
	appName := config.AppName

	// if state is empty, load the state file first
	if state == nil {
		LoadState()
	}

	appState, exists := state[appName]
	if !exists {
		state[appName] = initEmptyState()
		appState = state[appName]
	}

	return appState
}

func SaveState() {
	// validate config before saving it
	state.validate()

	// encode JSON
	jsonState, err := json.Marshal(state)
	if err != nil {
		logger.Errorf("Failed to save state file.\n%v", err)
		os.Exit(1)
	}

	// save file
	logger.Debugf("Saving novus state.")
	fs.WriteFileOrExit(novusStateFilePath, string(jsonState))
}
