package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)


func GenerateAPIKey() (string, error) {
	
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("Failed to generate API key: %v", err)
	}

	
	return fmt.Sprintf("sk_%s", hex.EncodeToString(bytes)), nil
}
