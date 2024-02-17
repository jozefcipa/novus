package ssl_manager

import (
	"path/filepath"

	"github.com/jozefcipa/novus/internal/config"
	"github.com/jozefcipa/novus/internal/fs"
	"github.com/jozefcipa/novus/internal/logger"
	"github.com/jozefcipa/novus/internal/mkcert"
	"github.com/jozefcipa/novus/internal/novus"
)

var certsDir string

type DomainCertificates map[string]mkcert.Certificate

func init() {
	// create a directory for the SSL certificates (~/.novus/certs)
	certsDir = filepath.Join(novus.NovusStateDir, "/certs")
	fs.MakeDirOrExit(certsDir)
}

func EnsureSSLCertificates(conf config.NovusConfig) DomainCertificates {
	domainCerts := make(DomainCertificates, len(conf.Routes))

	for _, route := range conf.Routes {
		cert := createCert(route.Domain)
		domainCerts[route.Domain] = cert
	}

	logger.Checkf("SSL certificates created.")

	return domainCerts
}

func createCert(domain string) mkcert.Certificate {
	domainCertDir := filepath.Join(certsDir, domain)

	// create a directory for the domain certificate (~/.novus/certs/{domain})
	fs.MakeDirOrExit(domainCertDir)

	// if certExists := fs.FileExists(domainCertDir); certExists {
	// 	logger.Debugf("SSL certificate already exists [%s]", domain)
	// 	return mkcert.Certificate{} // TODO: we have to return the existing certificate here
	// }

	logger.Debugf("Creating SSL certificate [%s]", domain)
	cert := mkcert.GenerateSSLCert(domain, domainCertDir)

	logger.Successf("SSL certificate generated [%s]\n", domain)

	return cert
}
