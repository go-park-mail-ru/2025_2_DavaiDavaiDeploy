package validation

import (
	"errors"
	"strings"
)

func ValidateLogin(login string) error {
	if len(login) < 6 || len(login) > 20 {
		return errors.New("validation error")
	}

	validChars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	for _, char := range login {
		if !strings.ContainsRune(validChars, char) {
			return errors.New("validation error")
		}
	}
	return nil
}
