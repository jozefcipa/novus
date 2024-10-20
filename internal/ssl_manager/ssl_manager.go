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

func ResolveDirs() {
	// ~/.novus/certs
	certsDir = filepath.Join(novus.NovusStateDir, "certs")
}

func EnsureSSLCertificates(conf config.NovusConfig, novusState *novus.NovusState, appName string) (shared.DomainCertificates, bool) {
	// Create a directory for the SSL certificates
	certsDir = filepath.Join(novus.NovusStateDir, "certs")
	fs.MakeDirOrExit(certsDir)

	// First make sure we have certs for internal routes
	internalAppState := novusState.Apps[novus.NovusInternalAppName]
	internalConfig := config.NovusConfig{
		AppName: novus.NovusInternalAppName,
		Routes:  internalAppState.Routes,
	}
	internalDomainCerts, hasNewInternalCerts := createCertsForConfig(internalConfig, internalAppState)

	// Now, create certs for the specific app
	domainCerts, hasNewCerts := createCertsForConfig(conf, novusState.Apps[appName])

	if hasNewCerts || hasNewInternalCerts {
		logger.Checkf("SSL certificates updated")
	} else {
		logger.Checkf("SSL certificates are up to date")
	}

	return shared.MergeMaps[shared.Certificate, shared.DomainCertificates](domainCerts, internalDomainCerts), hasNewCerts || hasNewInternalCerts
}

func createCertsForConfig(conf config.NovusConfig, appState *novus.AppState) (shared.DomainCertificates, bool) {
	domainCerts := make(shared.DomainCertificates, len(conf.Routes))
	hasNewCerts := false

	for _, route := range conf.Routes {
		cert, isNew := createCert(route.Domain, appState)
		if isNew {
			hasNewCerts = true
		}
		domainCerts[route.Domain] = cert
	}

	return domainCerts, hasNewCerts
}

func getCertificateDirectory(domain string) string {
	return filepath.Join(certsDir, domain)
}

func createCert(domain string, appState *novus.AppState) (shared.Certificate, bool) {
	timeNow := time.Now()

	// Check if the certificate already exists
	storedCert, exists := appState.SSLCertificates[domain]
	if exists {
		// Check certificate expiration
		// If the certificate expires in less than a month, we will renew it
		if timeNow.After(storedCert.ExpiresAt.AddDate(0, -1, 0)) {
			logger.Debugf("SSL certificate for domain [%s] expires in <1 month [%s]", domain, storedCert.CertFilePath)
		} else {
			logger.Debugf("SSL certificate for domain %s already exists [%s]", domain, storedCert.CertFilePath)
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

	logger.Debugf("SSL certificate generated [%s]", domain)

	return cert, true
}

func DeleteCert(domain string, appState *novus.AppState) {
	logger.Debugf("Deleting SSL certificate [%s]", domain)

	// Remove directory with SSL certificate for the given domain
	domainCertDir := getCertificateDirectory(domain)
	fs.DeleteDir(domainCertDir)

	// Remove cert from state
	delete(appState.SSLCertificates, domain)
}
