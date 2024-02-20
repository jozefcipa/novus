package mkcert

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/jozefcipa/novus/internal/config"
	"github.com/jozefcipa/novus/internal/logger"
	"github.com/jozefcipa/novus/internal/shared"
)

func Configure(conf config.NovusConfig) {
	logger.Debugf("Initializing mkcert")

	err := exec.Command("mkcert", "-install").Run()
	if err != nil {
		logger.Errorf("Failed to run \"mkcert -install\": %v", err)
		os.Exit(1)
	}
}

func GenerateSSLCert(domain string, dirPath string) shared.Certificate {
	certFilePath := filepath.Join(dirPath, "cert.pem")
	keyFilePath := filepath.Join(dirPath, "key.pem")

	params := []string{
		"-cert-file",
		certFilePath,
		"-key-file",
		keyFilePath,
		domain,
	}

	err := exec.Command("mkcert", params...).Run()
	if err != nil {
		logger.Errorf("Failed to run \"mkcert %s\": %v", strings.Join(params, " "), err)
		os.Exit(1)
	}

	return shared.Certificate{
		CertFilePath: certFilePath,
		KeyFilePath:  keyFilePath,

		// From mkcert's code:
		// (https://github.com/FiloSottile/mkcert/blob/2a46726cebac0ff4e1f133d90b4e4c42f1edf44a/cert.go#L59C2-L62C43)
		//
		// Certificates last for 2 years and 3 months, which is always less than
		// 825 days, the limit that macOS/iOS apply to all certificates,
		// including custom roots. See https://support.apple.com/en-us/HT210176.
		ExpiresAt: time.Now().AddDate(2, 3, 0),
	}
}
