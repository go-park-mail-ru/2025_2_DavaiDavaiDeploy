package validation

import (
	"strings"
)

const (
	ValidChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func ValidatePassword(password string) (string, bool) {
	if len(password) < 6 || len(password) > 20 {
		return "Invalid password length", false
	}

	for _, char := range password {
		if !strings.ContainsRune(ValidChars, char) {
			return "Password contains invalid characters", false
		}
	}
	return "Ok", true

}
