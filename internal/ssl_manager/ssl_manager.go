package ssl_manager

import (
	"path/filepath"
	"time"

	"github.com/jozefcipa/novus/internal/config"
	"github.com/jozefcipa/novus/internal/fs"
	"github.com/jozefcipa/novus/internal/logger"
	"github.com/jozefcipa/novus/internal/mkcert"
	"github.com/jozefcipa/novus/internal/novus"
	"github.com/jozefcipa/novus/internal/shared"
)

var certsDir string

func EnsureSSLCertificates(conf config.NovusConfig) (shared.DomainCertificates, bool) {
	// Create a directory for the SSL certificates (~/.novus/certs)
	certsDir = filepath.Join(novus.NovusStateDir, "certs")
	fs.MakeDirOrExit(certsDir)

	domainCerts := make(shared.DomainCertificates, len(conf.Routes))
	hasNewCerts := false

	for _, route := range conf.Routes {
		cert, isNew := createCert(route.Domain)
		if isNew {
			hasNewCerts = true
		}
		domainCerts[route.Domain] = cert
	}

	if hasNewCerts {
		logger.Checkf("SSL certificates updated")
	} else {
		logger.Checkf("SSL certificates are up to date")
	}

	return domainCerts, hasNewCerts
}

func getCertificateDirectory(domain string) string {
	// (~/.novus/certs/{domain})
	return filepath.Join(certsDir, domain)
}

func createCert(domain string) (shared.Certificate, bool) {
	appState, _ := novus.GetAppState()
	timeNow := time.Now()

	// Check if the certificate already exists
	storedCert, exists := appState.SSLCertificates[domain]
	if exists {
		// Check certificate expiration
		// If the certificate expires in less than a month, we will renew it
		if timeNow.After(storedCert.ExpiresAt.AddDate(0, -1, 0)) {
			logger.Debugf("SSL certificate expires in <1 month [%s=%s]", domain, storedCert.CertFilePath)
		} else {
			logger.Debugf("SSL certificate already exists [%s=%s]", domain, storedCert.CertFilePath)
			return storedCert, false
		}
	}

	// Create a directory for the domain certificate
	domainCertDir := getCertificateDirectory(domain)
	fs.MakeDirOrExit(domainCertDir)

	// Generate certificate
	logger.Debugf("Creating SSL certificate [%s]", domain)
	cert := mkcert.GenerateSSLCert(domain, domainCertDir)

	// Save cert in state
	appState.SSLCertificates[domain] = cert

	logger.Successf("✅ SSL certificate generated [%s]\n", domain)

	return cert, true
}

func DeleteCert(domain string) {
	appState, _ := novus.GetAppState()
	logger.Debugf("Deleting SSL certificate [%s]", domain)

	// Remove directory with SSL certificate for the given domain
	domainCertDir := getCertificateDirectory(domain)
	fs.DeleteDir(domainCertDir)

	// Remove cert from state
	delete(appState.SSLCertificates, domain)
}
