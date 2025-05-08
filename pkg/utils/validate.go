package utils

import (
	"regexp"
	"strings"
)

type InputType string

const (
	Email InputType = "email"
	Phone InputType = "phone"
)

type ValidateResult struct {
	IsValid bool
	Type    InputType
	Message string
}

func ValidateInput(input string) ValidateResult {

	input = strings.TrimSpace(input)

	if input == "" {
		return ValidateResult{
			IsValid: false,
			Type:    "",
			Message: "input cannot be empty",
		}
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	phoneRegex := regexp.MustCompile(`^0[35789][0-9]{8}$`)

	if emailRegex.MatchString(input) {
		return ValidateResult{
			IsValid: true,
			Type:    Email,
			Message: "not a valid email",
		}
	}

	if phoneRegex.MatchString(input) {
		return ValidateResult{
			IsValid: true,
			Type:    Phone,
			Message: "not a valid phone number",
		}
	}

	if strings.Contains(input, "@") {
		return ValidateResult{
			IsValid: false,
			Type:    Email,
			Message: "email format is incorrect",
		}
	}

	numberRegex := regexp.MustCompile(`^\d+$`)
	if numberRegex.MatchString(input) {
		return ValidateResult{
			IsValid: false,
			Type:    Phone,
			Message: "phone number format is incorrect",
		}
	}

	return ValidateResult{
		IsValid: false,
		Type:    "",
		Message: "input format is incorrect",
	}
}

// ValidateEmail validate the email format
func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}