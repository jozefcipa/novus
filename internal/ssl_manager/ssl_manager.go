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

func init() {
	// create a directory for the SSL certificates (~/.novus/certs)
	certsDir = filepath.Join(novus.NovusStateDir, "/certs")
	fs.MakeDirOrExit(certsDir)
}

func EnsureSSLCertificates(conf config.NovusConfig) (shared.DomainCertificates, bool) {
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

	// check if the certificate already exists
	storedCert, exists := appState.SSLCertificates[domain]
	if exists {
		// check certificate expiration
		// if the certificate expires in less than a month, we will renew it
		if timeNow.After(storedCert.ExpiresAt.AddDate(0, -1, 0)) {
			logger.Debugf("SSL certificate expires in <1 month [%s=%s]", domain, storedCert.CertFilePath)
		} else {
			logger.Debugf("SSL certificate already exists [%s=%s]", domain, storedCert.CertFilePath)
			return storedCert, false
		}
	}

	// create a directory for the domain certificate
	domainCertDir := getCertificateDirectory(domain)
	fs.MakeDirOrExit(domainCertDir)

	// generate certificate
	logger.Debugf("Creating SSL certificate [%s]", domain)
	cert := mkcert.GenerateSSLCert(domain, domainCertDir)

	// save cert in state
	appState.SSLCertificates[domain] = cert

	logger.Successf("SSL certificate generated [%s]\n", domain)

	return cert, true
}

func DeleteCert(domain string) {
	appState, _ := novus.GetAppState()
	logger.Debugf("Deleting SSL certificate [%s]", domain)

	// remove directory with SSL certificate for the given domain
	domainCertDir := getCertificateDirectory(domain)
	fs.DeleteDir(domainCertDir)

	// remove cert from state
	delete(appState.SSLCertificates, domain)
}
