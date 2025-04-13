// pkg/errcode/language.go
package errcode

import (
	"sort"
	"strconv"
	"strings"
)

type Language struct {
	Tag    string
	Weight float64
}

// "fr-FR,fr;q=0.9,en-US;q=0.8,en;q=0.7"
func ParseAcceptLanguage(header string) []Language {
	if header == "" {
		return nil
	}

	parts := strings.Split(header, ",")
	langs := make([]Language, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		langParts := strings.Split(part, ";")
		lang := strings.TrimSpace(langParts[0])

		weight := 1.0

		if len(langParts) > 1 {
			qPart := strings.TrimSpace(langParts[1])
			if strings.HasPrefix(qPart, "q=") {
				w := strings.TrimPrefix(qPart, "q=")
				if parsedWeight, err := strconv.ParseFloat(w, 64); err == nil {
					weight = parsedWeight
				}
			}
		}

		langs = append(langs, Language{Tag: lang, Weight: weight})
	}

	sort.Slice(langs, func(i, j int) bool {
		return langs[i].Weight > langs[j].Weight
	})

	return langs
}

func NormalizeLanguage(lang string) string {
	if idx := strings.Index(lang, ";"); idx != -1 {
		lang = lang[:idx]
	}

	return strings.ReplaceAll(lang, "-", "_")
}
