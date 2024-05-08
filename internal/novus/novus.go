package novus

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/go-playground/validator/v10"
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

func initStateFile() {
	// Create a directory ~/.novus
	// where we can store generated SSL certificates and application state
	NovusStateDir = filepath.Join(fs.UserHomeDir, ".novus")
	fs.MakeDirOrExit(NovusStateDir)
	novusStateFilePath = filepath.Join(NovusStateDir, "novus.json")
}

func loadState() {
	initStateFile()

	file, err := fs.ReadFile(novusStateFilePath)
	logger.Debugf("Loading state file [%s]", novusStateFilePath)
	// If there's an error, probably we didn't find the state, so initialize a new one
	if err != nil {
		logger.Debugf("State file not found. Creating a new one...")
		state = NovusState{
			DnsFiles: map[string]*DnsFiles{},
			Apps:     map[string]*AppState{},
		}
		return
	}

	if err := json.Unmarshal([]byte(file), &state); err != nil {
		logger.Errorf("Corrupted state file.\n%v", err)
		os.Exit(1)
	}

	state.validate()
}

func GetState() *NovusState {
	// If state is empty, load the state file first
	if len(state.Apps) == 0 {
		loadState()
	}

	return &state
}

func GetAppState(appName string) (appState *AppState, isNewState bool) {
	appState, exists := GetState().Apps[appName]
	if !exists {
		// Init empty state
		state.Apps[appName] = &AppState{
			Directory:       fs.CurrentDir,
			SSLCertificates: shared.DomainCertificates{},
			Routes:          []shared.Route{},
		}
		appState = state.Apps[appName]
		return appState, true
	}

	return appState, false
}

func SaveState() {
	// Validate config before saving it
	state.validate()

	// Encode JSON
	jsonState, err := json.MarshalIndent(state, "", "    ")
	if err != nil {
		logger.Errorf("Failed to save state file\n%v", err)
		os.Exit(1)
	}

	// Save file
	logger.Debugf("Saving novus state [%s]", novusStateFilePath)
	fs.WriteFileOrExit(novusStateFilePath, string(jsonState))
}
