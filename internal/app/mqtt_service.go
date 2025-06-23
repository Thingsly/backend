package app

import (
	"github.com/Thingsly/backend/mqtt"
	"github.com/Thingsly/backend/mqtt/device"
	"github.com/Thingsly/backend/mqtt/publish"
	"github.com/Thingsly/backend/mqtt/subscribe"

	"github.com/sirupsen/logrus"
)

// MQTTService Implement the MQTT related services
type MQTTService struct {
	initialized bool
}

// NewMQTTService Create a new MQTT service instance
func NewMQTTService() *MQTTService {
	return &MQTTService{
		initialized: false,
	}
}

// Name Return the service name
func (s *MQTTService) Name() string {
	return "MQTT service"
}

// Start Start the MQTT service
func (s *MQTTService) Start() error {
	// Check if MQTT is enabled
	// if !viper.GetBool("mqtt.enabled") {
	// 	logrus.Info("MQTT service is disabled, skipping initialization")
	// 	return nil
	// }

	logrus.Info("Starting MQTT service...")

	// Initialize the MQTT client
	if err := mqtt.MqttInit(); err != nil {
		return err
	}

	// Initialize the device status
	go device.InitDeviceStatus()

	// Initialize the subscription
	if err := subscribe.SubscribeInit(); err != nil {
		return err
	}

	// Initialize the publication
	publish.PublishInit()

	s.initialized = true
	logrus.Info("MQTT service started")
	return nil
}

// Stop Stop the MQTT service
func (s *MQTTService) Stop() error {
	if !s.initialized {
		return nil
	}

	logrus.Info("Stopping MQTT service...")

	logrus.Info("MQTT service stopped")
	return nil
}

// WithMQTTService Add the MQTT service to the application
func WithMQTTService() Option {
	return func(app *Application) error {
		service := NewMQTTService()
		app.RegisterService(service)
		return nil
	}
}
