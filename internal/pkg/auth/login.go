package auth

import (
	"regexp"
)

func ValidateLogin(login string) bool {
	if len(login) < 6 || len(login) > 20 {
		return false
	}

	matched, err := regexp.MatchString(`^[a-zA-Z0-9_]+$`, login)
	if err != nil {
		return false
	}
	return matched
}
