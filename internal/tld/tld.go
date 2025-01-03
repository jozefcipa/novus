package tld

import (
	"path/filepath"
	"slices"
	"strings"

	"github.com/jozefcipa/novus/internal/fs"
	"github.com/jozefcipa/novus/internal/logger"
	"github.com/jozefcipa/novus/internal/paths"
)

var existingTLDs []string

func init() {
	existingTLDs = []string{}
}

func ExtractFromDomain(domain string) string {
	urlParts := strings.Split(domain, ".")
	return urlParts[len(urlParts)-1]
}

func LoadExistingTLDsFile() []string {
	tldsFilePath := filepath.Join(paths.AssetsDir, "iana-tlds-list.txt")

	content, err := fs.ReadFile(tldsFilePath)
	if err != nil {
		logger.Warnf("Failed to read TLDs file: %v", err)
		return []string{}
	}
	logger.Debugf("Loaded TLDs file = %s.", tldsFilePath)

	for _, tld := range strings.Split(content, "\n") {
		// skip comments
		if strings.HasPrefix(tld, "#") {
			continue
		}

		existingTLDs = append(existingTLDs, strings.ToLower(tld))
	}

	return existingTLDs
}

func Exists(tld string) bool {
	return slices.Contains[[]string](existingTLDs, tld)
}
