package utils

import (
	"regexp"
)

func ValidateIndianPhoneNumber(phoneNumber string) bool {
	match := regexp.MustCompile(`^((\+?91)|0)?[0-9]{10}$`)
	return match.MatchString(phoneNumber)
}
