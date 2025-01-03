package validation

import (
	"os"
	"slices"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/jozefcipa/novus/internal/logger"
	"github.com/jozefcipa/novus/internal/shared"
)

// Make sure the config doesn't contain duplicate domains
func uniqueRoutesValidator(fl validator.FieldLevel) bool {
	value := fl.Field().Interface().([]shared.Route)

	domains := []string{}

	for _, route := range value {
		domain := strings.ToLower(route.Domain)
		if slices.Contains(domains, domain) {
			return false
		}
		domains = append(domains, domain)
	}

	return true
}

func RegisterUniqueRoutesValidator(validate *validator.Validate) {
	err := validate.RegisterValidation("unique_routes", uniqueRoutesValidator)
	if err != nil {
		logger.Errorf("Failed to register custom validator rule %v", err)
		os.Exit(1)
	}
}
