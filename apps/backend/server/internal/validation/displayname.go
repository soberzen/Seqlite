package validation

import (
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

var displayname validator.Func = func(fl validator.FieldLevel) bool {
	name := strings.TrimSpace(fl.Field().String())

	if name == "" {
		return false
	}

	if strings.ContainsRune(name, 0) {
		return false
	}
	for _, r := range name {
		if unicode.Is(unicode.Cf, r) || unicode.Is(unicode.Mn, r) ||
			unicode.Is(unicode.Me, r) {
			return false
		}
	}
	return true
}
