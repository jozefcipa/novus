package shared

import (
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/jozefcipa/novus/internal/logger"
	"github.com/jozefcipa/novus/internal/net"
)

// Make sure the config route doesn't use existing TLD domain
// We don't allow those as it would redirect the all URls with that TLD to our local DNS resolver
func nonExitentTLDValidator(fl validator.FieldLevel) bool {
	domain := fl.Field().Interface().(string)

	tld := strings.ToLower(net.ExtractTLD(domain))

	// check if tld exists
	if net.IsExistingTLD(tld) {
		return false
	}

	return true
}

func RegisterNonExistentTLDValidator(validate *validator.Validate) {
	err := validate.RegisterValidation("existing_tld", nonExitentTLDValidator)
	if err != nil {
		logger.Errorf("Failed to register custom validator rule %v", err)
		os.Exit(1)
	}
}
