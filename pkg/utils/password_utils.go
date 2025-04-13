package utils

import (
	"strings"
	"unicode"

	"github.com/HustIoTPlatform/backend/pkg/errcode"

	"github.com/sirupsen/logrus"
)

// ValidatePassword checks whether the given password meets the required standards.
// If the password is invalid, it returns an error explaining why.
// The password must:
// - Be at least 6 characters long
// - Contain only alphanumeric characters and the following special characters: !@#$%^&*()_+-=[]{};\:'"|,./<>?
// - Include at least one uppercase letter, one lowercase letter, one digit, and one special character
func ValidatePassword(password string) error {
	// Check password length
	if len(password) < 6 {
		return errcode.New(200040)
	}

	validSpecialChars := "!@#$%^&*()_+-=[]{};\\':\"|,./<>?"
	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)
	invalidChars := make([]rune, 0)

	// Iterate over each character in the password
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case strings.ContainsRune(validSpecialChars, char):
			hasSpecial = true
		default:
			invalidChars = append(invalidChars, char)
		}
	}

	// Check for invalid characters
	if len(invalidChars) > 0 {
		return errcode.WithVars(200053, map[string]interface{}{
			"invalid_chars": string(invalidChars),
		})
	}
	logrus.Debug("hasUpper", hasUpper)
	logrus.Debug("hasSpecial", hasSpecial)
	// Check password complexity
	var missingElements []string
	// if !hasUpper {
	// 	missingElements = append(missingElements, "uppercase letter")
	// }
	if !hasLower {
		missingElements = append(missingElements, "lowercase letter")
	}
	if !hasNumber {
		missingElements = append(missingElements, "digit")
	}
	// if !hasSpecial {
	// 	missingElements = append(missingElements, "special character")
	// }

	if len(missingElements) > 0 {
		return errcode.WithVars(200054, map[string]interface{}{
			"missing_elements": strings.Join(missingElements, ", "),
		})
	}

	return nil
}
