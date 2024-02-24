package ssl_manager

import (
	"path/filepath"

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

	logger.Checkf("SSL certificates created.")

	return domainCerts, hasNewCerts
}

func createCert(domain string) (shared.Certificate, bool) {
	appState := novus.GetState()

	// check if the certificate already exists
	storedCert, exists := appState.SSLCertificates[domain]
	if exists {
		// TODO: check certificate expiration

		logger.Debugf("SSL certificate already exists [%s]", storedCert.CertFilePath)
		return storedCert, false
	}

	// create a directory for the domain certificate (~/.novus/certs/{domain})
	domainCertDir := filepath.Join(certsDir, domain)
	fs.MakeDirOrExit(domainCertDir)

	// generate certificate
	logger.Debugf("Creating SSL certificate [%s]", domain)
	cert := mkcert.GenerateSSLCert(domain, domainCertDir)

	// save cert in state
	appState.SSLCertificates[domain] = cert

	logger.Successf("SSL certificate generated [%s]\n", domain)

	return cert, true
}
