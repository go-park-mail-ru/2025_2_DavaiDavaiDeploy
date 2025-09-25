package validation

import (
	"errors"
	"kinopoisk/internal/pkg/auth/constants"
	"strings"
)

func ValidateLogin(login string) error {
	if len(login) < 6 || len(login) > 20 {
		return errors.New("invalid login length")
	}

	validChars := constants.ValidChars
	for _, char := range login {
		if !strings.ContainsRune(validChars, char) {
			return errors.New("login contains invalid characters")
		}
	}
	return nil
}
