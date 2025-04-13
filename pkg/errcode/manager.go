// pkg/errcode/manager.go
package errcode

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/patrickmn/go-cache"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Messages map[int]map[string]string `yaml:"messages"` // code -> language -> message
	Metadata struct {
		Version            string   `yaml:"version"`
		LastUpdated        string   `yaml:"last_updated"`
		SupportedLanguages []string `yaml:"supported_languages"`
	} `yaml:"metadata"`
}

type StringConfig struct {
	Messages map[string]map[string]string `yaml:"messages"` // key -> language -> message
	Metadata struct {
		Version            string   `yaml:"version"`
		LastUpdated        string   `yaml:"last_updated"`
		SupportedLanguages []string `yaml:"supported_languages"`
	} `yaml:"metadata"`
}

type ErrorManager struct {
	messages        map[int]map[string]string    // code -> language -> message
	messageStr      map[string]map[string]string // key -> language -> message
	cache           *cache.Cache
	defaultLanguage string
	configPath      string
	strConfigPath   string
}

func NewErrorManager(configPath string, strConfigPath string) *ErrorManager {
	manager := &ErrorManager{
		messages:        make(map[int]map[string]string),
		messageStr:      make(map[string]map[string]string),
		cache:           cache.New(10*time.Minute, 20*time.Minute),
		defaultLanguage: "en_US",
		configPath:      configPath,
		strConfigPath:   strConfigPath,
	}
	return manager
}

// LoadMessages loads error codes and string configurations
func (m *ErrorManager) LoadMessages() error {
	// Load error code configuration
	if err := m.loadErrorMessages(); err != nil {
		return fmt.Errorf("failed to load error code configuration: %w", err)
	}

	// Load string configuration
	if err := m.loadStringMessages(); err != nil {
		return fmt.Errorf("failed to load string configuration: %w", err)
	}

	return nil
}

// loadErrorMessages loads the error code configuration
func (m *ErrorManager) loadErrorMessages() error {
	data, err := os.ReadFile(filepath.Clean(m.configPath))
	if err != nil {
		return fmt.Errorf("failed to read error code configuration file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse error code configuration file: %w", err)
	}

	for code := range config.Messages {
		if !m.validateCode(code) {
			return fmt.Errorf("invalid error code format: %d", code)
		}
	}

	m.messages = config.Messages
	return nil
}

// loadStringMessages loads the string configuration
func (m *ErrorManager) loadStringMessages() error {
	data, err := os.ReadFile(filepath.Clean(m.strConfigPath))
	if err != nil {
		return fmt.Errorf("failed to read string configuration file: %w", err)
	}

	var config StringConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse string configuration file: %w", err)
	}

	m.messageStr = config.Messages
	return nil
}

// GetMessage retrieves the message for the specified error code
// Parameters:
//   - code: Error code
//   - lang: Language code, if empty, the default language is used
func (m *ErrorManager) GetMessage(code int, acceptLanguage string) string {
	// If no language is specified, use the default language
	if acceptLanguage == "" {
		return m.getMessageForLanguage(code, m.defaultLanguage)
	}

	// Parse the Accept-Language header
	languages := ParseAcceptLanguage(acceptLanguage)

	// Try each language by priority
	for _, lang := range languages {
		normalizedLang := NormalizeLanguage(lang.Tag)
		if msg := m.getMessageForLanguage(code, normalizedLang); msg != "" {
			return msg
		}
	}

	// If none found, use the default language
	return m.getMessageForLanguage(code, m.defaultLanguage)
}

// GetMessageStr retrieves the string message for the specified key
func (m *ErrorManager) GetMessageStr(key string, acceptLanguage string) string {
	// If no language is specified, use the default language
	if acceptLanguage == "" {
		return m.getStrMessageForLanguage(key, m.defaultLanguage)
	}

	languages := ParseAcceptLanguage(acceptLanguage)
	for _, lang := range languages {
		normalizedLang := NormalizeLanguage(lang.Tag)
		if msg := m.getStrMessageForLanguage(key, normalizedLang); msg != "" {
			return msg
		}
	}

	return m.getStrMessageForLanguage(key, m.defaultLanguage)
}

// getMessageForLanguage retrieves the message for the specified language
func (m *ErrorManager) getMessageForLanguage(code int, lang string) string {
	// Attempt to get from cache
	cacheKey := fmt.Sprintf("%d:%s", code, lang)
	if msg, found := m.cache.Get(cacheKey); found {
		return msg.(string)
	}

	// Retrieve the message from memory
	if messages, ok := m.messages[code]; ok {
		if msg, ok := messages[lang]; ok {
			m.cache.Set(cacheKey, msg, cache.DefaultExpiration)
			return msg
		}
	}

	// If the language is not the default and the message is not found, return an empty string
	// This allows continuing the search in other language options
	if lang != m.defaultLanguage {
		return ""
	}

	// Use the default error message for the default language
	defaultMsg := "Unknown Error"
	if lang != "en_US" {
		defaultMsg = "Unknown Error"
	}
	return defaultMsg
}

// getStrMessageForLanguage retrieves the string message for the specified language
func (m *ErrorManager) getStrMessageForLanguage(key string, lang string) string {
	// Check cache for string message
	cacheKey := fmt.Sprintf("str:%s:%s", key, lang)
	if msg, found := m.cache.Get(cacheKey); found {
		return msg.(string)
	}

	// Retrieve the string message from memory
	if messages, ok := m.messageStr[key]; ok {
		if msg, ok := messages[lang]; ok {
			m.cache.Set(cacheKey, msg, cache.DefaultExpiration)
			return msg
		}
	}

	// If the language is not the default and the string message is not found, return an empty string
	if lang != m.defaultLanguage {
		return ""
	}

	// If the string message is not found, return the key itself
	return key
}

// validateCode validates the error code format
func (m *ErrorManager) validateCode(code int) bool {
	// Special handling for success code
	if code == 200 {
		return true
	}

	// Check length and range
	if code < 100000 || code > 599999 {
		return false
	}

	// Check the allowed first digit (1, 2, 3, 4, 5)
	firstDigit := code / 100000
	switch firstDigit {
	case 1, 2, 3, 4, 5:
		return true
	default:
		return false
	}
}

// SetDefaultLanguage sets the default language
func (m *ErrorManager) SetDefaultLanguage(lang string) {
	m.defaultLanguage = lang
}

// ClearCache clears the cache
func (m *ErrorManager) ClearCache() {
	m.cache.Flush()
}

// Example usage:
/*
func main() {
    // Create error code manager
    manager := NewErrorManager("config/messages.yaml")

    // Load configuration
    if err := manager.LoadMessages(); err != nil {
        log.Fatalf("Failed to load error code configuration: %v", err)
    }

    // Retrieve error message
	msg := manager.GetMessage("100001", "zh_CN")
	fmt.Println(msg) // Output: Service temporarily unavailable

	// Use default language
	msg = manager.GetMessage("100001", "")
	fmt.Println(msg) // Output: Service temporarily unavailable

    // Use English
    msg = manager.GetMessage("100001", "en_US")
    fmt.Println(msg) // Output: Service Temporarily Unavailable
}
*/
