package validation

import (
	"errors"
	"strings"
)

func ValidatePassword(password string) error {
	if len(password) < 6 || len(password) > 20 {
		return errors.New("validation error")
	}

	validChars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789.!?"
	for _, char := range password {
		if !strings.ContainsRune(validChars, char) {
			return errors.New("validation error")
		}
	}
	return nil
}
