package net

import (
	"path/filepath"
	"slices"
	"strings"

	"github.com/jozefcipa/novus/internal/fs"
	"github.com/jozefcipa/novus/internal/logger"
)

var existingTLDs []string

func init() {
	existingTLDs = []string{}
}

func ExtractTLD(domain string) string {
	urlParts := strings.Split(domain, ".")
	return urlParts[len(urlParts)-1]
}

func LoadExistingTLDsFile() []string {
	content, err := fs.ReadFile(filepath.Join(fs.AssetsDir, "iana-tlds-list.txt"))
	if err != nil {
		logger.Warnf("Failed to read TLDs file: %v", err)
		return []string{}
	}
	logger.Debugf("Loaded TLDs file.")

	for _, tld := range strings.Split(content, "\n") {
		// skip comments
		if strings.HasPrefix(tld, "#") {
			continue
		}

		existingTLDs = append(existingTLDs, strings.ToLower(tld))
	}

	return existingTLDs
}

func IsExistingTLD(tld string) bool {
	return slices.Contains[[]string](existingTLDs, tld)
}
