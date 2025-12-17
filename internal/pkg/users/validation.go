package users

import "strings"

const (
	ValidChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+-=[]{}|;:,.<>?`~"
)

func Validaton(login, password string) (string, bool) {
	if len(login) < 6 || len(login) > 15 {
		return "Invalid login length", false
	}

	if len(password) < 6 || len(password) > 15 {
		return "Invalid password length", false
	}

	for _, char := range login {
		if !strings.ContainsRune(ValidChars, char) {
			return "Login contains invalid characters", false
		}
	}

	for _, char := range password {
		if !strings.ContainsRune(ValidChars, char) {
			return "Password contains invalid characters", false
		}
	}
	return "Ok", true
}
