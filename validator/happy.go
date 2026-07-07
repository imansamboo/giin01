package validator

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

func ValidateHappy(field validator.FieldLevel) bool {
	return strings.Contains(field.Field().String(), "happy")
}
