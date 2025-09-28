package validation

import (
	"kinopoisk/internal/pkg/auth/constants"
	"strings"
)

func ValidateLogin(login string) (string, bool) {
	if len(login) < 6 || len(login) > 20 {
		return "Invalid login length", false
	}

	validChars := constants.ValidChars
	for _, char := range login {
		if !strings.ContainsRune(validChars, char) {
			return "Login contains invalid characters", false
		}
	}
	return "Ok", true
}
