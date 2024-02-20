package novus

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jozefcipa/novus/internal/fs"
	"github.com/jozefcipa/novus/internal/logger"
	"github.com/jozefcipa/novus/internal/shared"
)

var NovusStateDir string

var novusStateFilePath string
var state NovusState

type AppState struct {
	Directory       string                    `json:"directory"`
	SSLCertificates shared.DomainCertificates `json:"sslCertificates"`
	Routes          []shared.Route            `json:"routes"`
}

type NovusState map[string]*AppState

func (config *NovusState) validate() {
	logger.Debugf("Validating state file")
	// TODO: make sure the loaded file is in the expected format
	// https://github.com/go-playground/validator
	// Example: https://github.com/go-playground/validator/blob/master/_examples/simple/main.go
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

func GetState(appName string) *AppState {
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
	fmt.Print(string(jsonState))
	if err != nil {
		logger.Errorf("Failed to save state file.\n%v", err)
		os.Exit(1)
	}

	// save file
	logger.Debugf("Saving novus state.")
	fs.WriteFileOrExit(novusStateFilePath, string(jsonState))
}
