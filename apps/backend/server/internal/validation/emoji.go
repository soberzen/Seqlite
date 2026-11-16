package validation

import (
	"strings"
	"unicode/utf8"

	"github.com/go-playground/validator/v10"
)

var emoji validator.Func = func(fl validator.FieldLevel) bool {
	emoji, ok := fl.Field().Interface().(string)

	if !ok {
		return false
	}
	return utf8.RuneCountInString(emoji) <= 32 &&
		!strings.ContainsRune(emoji, 0)
}
