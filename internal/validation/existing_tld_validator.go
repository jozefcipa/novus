package validation

import (
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/jozefcipa/novus/internal/logger"
	"github.com/jozefcipa/novus/internal/tld"
)

// Make sure the config route doesn't use existing TLD domain
// We don't allow those as it would redirect the all URls with that TLD to our local DNS resolver
func nonExistentTLDValidator(fl validator.FieldLevel) bool {
	domain := fl.Field().Interface().(string)

	domainTld := strings.ToLower(tld.ExtractFromDomain(domain))

	return !tld.Exists(domainTld)
}

func RegisterNonExistentTLDValidator(validate *validator.Validate) {
	err := validate.RegisterValidation("existing_tld", nonExistentTLDValidator)
	if err != nil {
		logger.Errorf("Failed to register custom validator rule %v", err)
		os.Exit(1)
	}
}
