package diff_manager

import (
	"fmt"

	"github.com/jozefcipa/novus/internal/novus"
)

type DuplicateDomainError struct {
	DuplicateDomain       string
	OriginalAppWithDomain string
}

func (e *DuplicateDomainError) Error() string {
	if e.OriginalAppWithDomain == novus.GlobalAppName {
		return fmt.Sprintf("Domain %s is already defined in the global scope", e.DuplicateDomain)
	} else {
		return fmt.Sprintf("Domain %s is already defined by app \"%s\"", e.DuplicateDomain, e.OriginalAppWithDomain)
	}
}
