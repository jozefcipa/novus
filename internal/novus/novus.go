package novus

import (
	"path/filepath"

	"github.com/jozefcipa/novus/internal/fs"
)

var NovusStateDir string

func init() {
	// Create a directory ~/.novus
	// where we can store generated SSL certificates and application state
	NovusStateDir = filepath.Join(fs.UserHomeDir, ".novus")
	fs.MakeDirOrExit(NovusStateDir)
}
