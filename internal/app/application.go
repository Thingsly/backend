package app

import (
	"github.com/Thingsly/backend/internal/query"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type Application struct {
	Config         *viper.Viper
	Logger         *logrus.Logger
	DB             *gorm.DB
	RedisClient    *redis.Client
	ServiceManager *ServiceManager
}

func NewApplication(options ...Option) (*Application, error) {
	app := &Application{
		Logger:         logrus.New(),
		ServiceManager: NewServiceManager(),
	}

	for _, option := range options {
		if err := option(app); err != nil {
			return nil, err
		}
	}

	// Set query default DB
	if app.DB != nil {
		query.SetDefault(app.DB)
	}

	return app, nil
}

// RegisterService Register a service to the application
func (app *Application) RegisterService(service Service) {
	app.ServiceManager.RegisterService(service)
}

// Start Start all registered services
func (app *Application) Start() error {
	return app.ServiceManager.StartAll()
}

// Shutdown Gracefully close all resources
func (app *Application) Shutdown() {
	// Stop all services
	app.ServiceManager.StopAll()

	// Close Redis connection
	if app.RedisClient != nil {
		app.RedisClient.Close()
		app.Logger.Info("Redis connection closed")
	}

	// DB does not need to be explicitly closed, gorm.DB does not have a Close method

	app.Logger.Info("All resources have been successfully cleaned up")
}

// Wait Wait for all services to complete
func (app *Application) Wait() {
	app.ServiceManager.Wait()
}

// Option Define application initialization options
type Option func(*Application) error