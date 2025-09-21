package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"regexp"
)

func HashPassword(password string) string {
	h := sha256.New()
	h.Write([]byte(password))
	return hex.EncodeToString(h.Sum(nil))
}

func ValidatePassword(password string) bool {
	if len(password) < 6 || len(password) > 50 {
		return false
	}

	matched, err := regexp.MatchString(`^[a-zA-Z0-9!@#$%^&*()_+]+$`, password)
	if err != nil {
		return false
	}
	return matched
}
