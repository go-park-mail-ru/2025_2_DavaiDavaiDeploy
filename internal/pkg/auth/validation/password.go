package validation

import (
	"kinopoisk/internal/pkg/auth/constants"
	"strings"
)

func ValidatePassword(password string) (string, bool) {
	if len(password) < 6 || len(password) > 20 {
		return "Invalid password length", false
	}

	validChars := constants.ValidChars
	for _, char := range password {
		if !strings.ContainsRune(validChars, char) {
			return "Password contains invalid characters", false
		}
	}
	return "Ok", true

}
