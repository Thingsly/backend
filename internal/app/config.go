package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// Config file path priority:
// 1. The path specified by the environment variable THINGSLY_CONFIG_PATH
// 2. The configs/conf.yml in the current directory
// 3. The .thingsly/conf.yml in the user's home directory
// 4. /etc/thingsly/conf.yml (Linux/Mac)
// 5. Built-in default configuration

// LoadConfig Load the configuration file based on priority
func LoadConfig() (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigType("yml")

	// 1. Check environment variable
	if configPath := os.Getenv("THINGSLY_CONFIG_PATH"); configPath != "" {
		v.SetConfigFile(configPath)
		if err := v.ReadInConfig(); err == nil {
			return v, nil
		}
	}

	// 2. Check current directory
	v.SetConfigFile("./configs/conf.yml")
	if err := v.ReadInConfig(); err == nil {
		return v, nil
	}

	// 3. Check user home directory
	home, err := os.UserHomeDir()
	if err == nil {
		userConfigPath := filepath.Join(home, ".thingsly", "conf.yml")
		v.SetConfigFile(userConfigPath)
		if err := v.ReadInConfig(); err == nil {
			return v, nil
		}
	}

	// 4. Check system directory (only for Linux/Mac)
	if os.Getenv("GOOS") != "windows" {
		v.SetConfigFile("/etc/thingsly/conf.yml")
		if err := v.ReadInConfig(); err == nil {
			return v, nil
		}
	}

	// 5. Use built-in default configuration
	return nil, fmt.Errorf("no available configuration file found")
}

// Load specific environment configuration
func LoadEnvironmentConfig(env string) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigType("yml")

	var configFile string
	switch env {
	case "dev":
		configFile = "./configs/conf-localdev.yml"
	case "test":
		configFile = "./configs/conf-test.yml"
	case "prod":
		configFile = "./configs/conf.yml"
	default:
		return nil, fmt.Errorf("unknown environment type: %s", env)
	}

	v.SetConfigFile(configFile)
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	// Set environment variable prefix and enable automatic reading of environment variables
	v.SetEnvPrefix("THINGSLY")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	return v, nil
}