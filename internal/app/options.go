package app

import (
	"github.com/Thingsly/backend/initialize"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// WithConfig Use the already initialized Viper instance
func WithConfig(config *viper.Viper) Option {
	return func(app *Application) error {
		app.Config = config
		// Set the global viper instance (for compatibility with existing code)
		for _, key := range config.AllKeys() {
			viper.Set(key, config.Get(key))
		}
		return nil
	}
}

// WithEnvironment Load the configuration based on the environment name
func WithEnvironment(env string) Option {
	return func(app *Application) error {
		config, err := LoadEnvironmentConfig(env)
		if err != nil {
			return err
		}
		return WithConfig(config)(app)
	}
}

// WithProductionConfig Use the production environment configuration
func WithProductionConfig() Option {
	return WithEnvironment("prod")
}

// WithDevelopmentConfig Use the development environment configuration
func WithDevelopmentConfig() Option {
	return WithEnvironment("dev")
}

// WithTestConfig Use the test environment configuration
func WithTestConfig() Option {
	return WithEnvironment("test")
}

// WithRsaDecrypt Initialize the RSA decryption
func WithRsaDecrypt(keyPath string) Option {
	return func(app *Application) error {
		return initialize.RsaDecryptInit(keyPath)
	}
}

// WithLogger Configure the logging system
func WithLogger() Option {
	return func(app *Application) error {
		if err := initialize.LogInIt(); err != nil {
			return err
		}
		app.Logger = logrus.StandardLogger()
		return nil
	}
}

// WithDatabase Initialize the database connection
func WithDatabase() Option {
	return func(app *Application) error {
		db, err := initialize.PgInit()
		if err != nil {
			return err
		}
		app.DB = db
		return nil
	}
}

// WithRedis Initialize the Redis connection
func WithRedis() Option {
	return func(app *Application) error {
		client, err := initialize.RedisInit()
		if err != nil {
			return err
		}
		app.RedisClient = client
		return nil
	}
}