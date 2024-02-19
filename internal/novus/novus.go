package novus

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/jozefcipa/novus/internal/fs"
	"github.com/jozefcipa/novus/internal/logger"
)

var NovusStateDir string

var novusStateFilePath string
var state NovusState

type tldDomain struct {
	Tld           string
	CertExpiresAt string
}

type route struct {
	Url      string
	Upstream string
}

type AppState struct {
	Directory  string      `json:"directory"`
	TldDomains []tldDomain `json:"tldDomains"`
	Routes     []route     `json:"routes"`
}

type NovusState map[string]AppState

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

func LoadState() {
	file, err := fs.ReadFile(novusStateFilePath)
	// if there's an error, probably we didn't find the state, so initialize a new one
	if err != nil {
		logger.Debugf("State file not found. Creating a new one.")
		state = NovusState{
			"default": {
				Directory:  "",
				Routes:     []route{},
				TldDomains: []tldDomain{},
			},
		}
		return
	}

	if err := json.Unmarshal([]byte(file), &state); err != nil {
		logger.Errorf("Corrupted state file.\n%v\n", err)
		os.Exit(1)
	}

	state.validate()
}

func GetState() NovusState {
	// if state is empty, load the state file first
	if state == nil {
		LoadState()
	}

	return state
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
	fs.WriteFileOrExit(novusStateFilePath, string(jsonState))
}
