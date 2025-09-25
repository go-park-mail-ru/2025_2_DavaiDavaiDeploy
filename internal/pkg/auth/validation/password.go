package validation

import (
	"errors"
	"kinopoisk/internal/pkg/auth/constants"
	"strings"
)

func ValidatePassword(password string) error {
	if len(password) < 6 || len(password) > 20 {
		return errors.New("invalid password length")
	}

	validChars := constants.ValidChars
	for _, char := range password {
		if !strings.ContainsRune(validChars, char) {
			return errors.New("password contains invalid characters")
		}
	}
	return nil

}
