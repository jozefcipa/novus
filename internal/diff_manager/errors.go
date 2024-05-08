package diff_manager

import "fmt"

type DuplicateDomainError struct {
	DuplicateDomain       string
	OriginalAppWithDomain string
}

func (e *DuplicateDomainError) Error() string {
	return fmt.Sprintf("Domain %s is already defined by app \"%s\"", e.DuplicateDomain, e.OriginalAppWithDomain)
}
