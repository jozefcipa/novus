package novus

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/jozefcipa/novus/internal/fs"
	"github.com/jozefcipa/novus/internal/logger"
	"github.com/jozefcipa/novus/internal/shared"
)

var NovusStateDir string
var NovusStateFilePath string
var state NovusState
var stateLoaded bool

const NovusInternalDomain = "internal.novus"
const NovusIndexDomain = "index.novus"
const NovusInternalAppName = "_novus"

func init() {
	stateLoaded = false
}

func ResolveDirs() {
	// where we can store generated SSL certificates and application state
	NovusStateDir = filepath.Join(fs.UserHomeDir, ".novus")
	NovusStateFilePath = filepath.Join(NovusStateDir, "novus.json")
}

func initStateDir() {
	// Create a directory ~/.novus
	fs.MakeDirOrExit(NovusStateDir)
}

func loadState() {
	initStateDir()

	file, err := fs.ReadFile(NovusStateFilePath)
	logger.Debugf("Loading state file [%s]", NovusStateFilePath)

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

	// If the internal domain is not yet there, add it now
	if _, found := state.Apps[NovusInternalAppName]; !found {
		state.Apps[NovusInternalAppName] = &AppState{
			Directory:       NovusStateDir,
			Status:          APP_ACTIVE,
			SSLCertificates: shared.DomainCertificates{},
			Routes: []shared.Route{
				{Domain: NovusInternalDomain, Upstream: "http://example.org"},
				{Domain: NovusIndexDomain, Upstream: "http://example.org"},
			},
		}
	}

	state.validate()
	stateLoaded = true
}

func GetState() *NovusState {
	// If state hasn't been loaded yet, load it now
	if !stateLoaded {
		loadState()
	}

	return &state
}

func GetAppState(appName string) (*AppState, bool) {
	appState, exists := GetState().Apps[appName]
	return appState, exists
}

func InitializeAppState(appName string) *AppState {
	appState, exists := GetAppState(appName)

	if !exists {
		// Init empty state
		state.Apps[appName] = &AppState{
			Status:          APP_ACTIVE,
			Directory:       fs.CurrentDir,
			SSLCertificates: shared.DomainCertificates{},
			Routes:          []shared.Route{},
		}
		appState = state.Apps[appName]
	}

	return appState
}

func RemoveAppState(appName string) {
	logger.Debugf("Removing app configuration from state [%s]", appName)
	delete(state.Apps, appName)
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
	logger.Debugf("Saving novus state [%s]", NovusStateFilePath)
	fs.WriteFileOrExit(NovusStateFilePath, string(jsonState))
}
