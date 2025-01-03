package paths

import (
	"path/filepath"

	"github.com/jozefcipa/novus/internal/logger"
)

// Used to store SSL certificate files (~/.novus/certs)
var SSLCertificatesDir string

func resolveSSLCertDirs() {
	SSLCertificatesDir = filepath.Join(NovusStateDir, "certs")

	logger.Debugf("SSL paths resolved.\n\tSSLCertificatesDir = %s", SSLCertificatesDir)
}
