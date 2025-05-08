package app

import (
	"github.com/Thingsly/backend/initialize/croninit"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// CronService Implement the cron task service
type CronService struct {
	initialized bool
}

// NewCronService Create a cron task service instance
func NewCronService() *CronService {
	return &CronService{
		initialized: false,
	}
}

// Name Return the service name
func (s *CronService) Name() string {
	return "Cron task service"
}

// Start Start the cron task service
func (s *CronService) Start() error {
	// Check if cron task is enabled
	if !viper.GetBool("cron.enabled") {
		logrus.Info("Cron task service is disabled, skipping initialization")
		return nil
	}

	logrus.Info("Starting cron task service...")

	// Initialize cron tasks
	croninit.CronInit()

	s.initialized = true
	logrus.Info("Cron task service started")
	return nil
}

// Stop Stop the cron task service
func (s *CronService) Stop() error {
	if !s.initialized {
		return nil
	}

	logrus.Info("Stopping cron task service...")
	// If croninit provides a stop method, call it here

	logrus.Info("Cron task service stopped")
	return nil
}

// WithCronService Add the cron task service to the application
func WithCronService() Option {
	return func(app *Application) error {
		service := NewCronService()
		app.RegisterService(service)
		return nil
	}
}
