package ssl_manager

import (
	"fmt"
	"path/filepath"

	"github.com/jozefcipa/novus/internal/config"
	"github.com/jozefcipa/novus/internal/fs"
	"github.com/jozefcipa/novus/internal/logger"
	"github.com/jozefcipa/novus/internal/mkcert"
	"github.com/jozefcipa/novus/internal/novus"
)

func EnsureSSLCertificates(conf config.NovusConfig) {
	// create a directory for the SSL certificates (~/.novus/certs)
	certsDir := filepath.Join(novus.NovusStateDir, "/certs")
	fs.MakeDirOrExit(certsDir)

	for _, route := range conf.Routes {
		domainCertDir := filepath.Join(certsDir, route.Url)

		// create a directory for the domain certificate (~/.novus/certs/{domain})
		fs.MakeDirOrExit(domainCertDir)

		if certExists := fs.FileExists(domainCertDir); certExists {
			logger.Debugf("SSL certificate already exists [%s]", route.Url)
			continue
		}

		logger.Debugf("Creating SSL certificate [%s]", route.Url)
		cert := mkcert.GenerateSSLCert(route.Url, domainCertDir)

		logger.Checkf("SSL certificate created [%s]", route.Url)

		fmt.Println(cert.CertFilePath) // TODO temp
	}

	logger.Checkf("SSL certificates created.")
}
