package validator

import (
	"slices"
	"strings"
	"unicode/utf8"
)

type Validator struct {
	FieldErrors map[string]string
}

// IsValid returns true if there are no errors, otherwise false.
func (v *Validator) IsValid() bool {
	return len(v.FieldErrors) == 0
}

// AddFieldError adds an error message for a specific field to the FieldErrors map.
func (v *Validator) AddFieldError(field, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[field]; !exists {
		v.FieldErrors[field] = message
	}
}

// CheckField checks a condition. If that fails, it adds an error message to the FieldErrors map.
func (v *Validator) CheckField(ok bool, field, message string) {
	if !ok {
		v.AddFieldError(field, message)
	}
}

// NotBlank checks if the value is not an empty string.
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// MaxChars checks if the value is not longer than max characters.
func MaxChars(value string, max int) bool {
	// NOTE: utf8.RuneCountInString counts the number of unicode code points
	return utf8.RuneCountInString(value) <= max
}

// PermittedValue checks if the value is in the list of permitted values.
func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}
