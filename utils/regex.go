package utils

import "regexp"

var (
	ToSnakeCaseRegex          = regexp.MustCompile(`([a-z0-9])([A-Z])`)
	NumericRegex              = regexp.MustCompile(`^[0-9]+$`)
	PasswordHasLowercaseRegex = regexp.MustCompile(`[a-z]`)
	PasswordHasSpecialRegex   = regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{}|;:,.<>?/~\x60]`)
	PasswordAllowedCharsRegex = regexp.MustCompile(`^[A-Z0-9!@#$%^&*()_+\-=\[\]{}|;:,.<>?/~\x60]+$`)
	WhitespaceRegex           = regexp.MustCompile(`\s+`)
)

func IsStringNumeric(s string) bool {
	return NumericRegex.MatchString(s)
}
