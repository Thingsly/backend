package utils

import "strings"

func FormatLangCode(acceptLanguage string) string {
	// If empty, return the default value "en_US"
	if acceptLanguage == "" {
		return "en_US"
	}

	// Split accept-language and take the first one
	langs := strings.Split(acceptLanguage, ",")
	primaryLang := strings.TrimSpace(langs[0])

	// Handle possible weight values like vi-VN;q=0.9
	primaryLang = strings.Split(primaryLang, ";")[0]

	// Replace "-" with "_"
	primaryLang = strings.Replace(primaryLang, "-", "_", 1)

	// Handle special cases
	switch primaryLang {
	case "vi":
		return "vi_VN"
	case "en":
		return "en_US"
	}

	// If it is already in the correct format, return it directly
	if len(primaryLang) == 5 && primaryLang[2] == '_' {
		return primaryLang
	}

	// In other cases, return the default value
	return "vi_VN"
}
