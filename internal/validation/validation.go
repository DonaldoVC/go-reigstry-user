// Package validation contains input validation and normalization helpers.
package validation

import (
	"net/mail"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	emailShape = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
	phoneShape = regexp.MustCompile(`^\+?[1-9][0-9]{9,14}$`)
)

// NormalizeEmail makes uniqueness checks case-insensitive.
func NormalizeEmail(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

// NormalizePhone removes common presentation characters before validation.
func NormalizePhone(value string) string {
	replacer := strings.NewReplacer(" ", "", "-", "", "(", "", ")", "")
	return replacer.Replace(strings.TrimSpace(value))
}

func ValidEmail(value string) bool {
	if len(value) > 254 || !emailShape.MatchString(value) {
		return false
	}
	parsed, err := mail.ParseAddress(value)
	return err == nil && parsed.Address == value
}

// ValidPhone accepts 10 to 15 digits, optionally prefixed by '+'.
func ValidPhone(value string) bool {
	return phoneShape.MatchString(value)
}

// ValidPassword enforces length and character-class requirements.
func ValidPassword(value string) bool {
	length := utf8.RuneCountInString(value)
	if length < 6 || length > 12 {
		return false
	}
	var upper, lower, number bool
	for _, character := range value {
		upper = upper || unicode.IsUpper(character)
		lower = lower || unicode.IsLower(character)
		number = number || unicode.IsDigit(character)
	}
	return upper && lower && number && strings.ContainsAny(value, "@$&")
}
