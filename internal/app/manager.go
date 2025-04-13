package app

import (
	"github.com/HustIoTPlatform/backend/mqtt/device"
	"github.com/HustIoTPlatform/backend/pkg/global"

	"github.com/sirupsen/logrus"
)

// Manager Service Manager
type Manager struct {
	deviceListener *device.DeviceListener
}

// NewManager Creates a new service manager
func NewManager() *Manager {
	return &Manager{
		deviceListener: device.NewDeviceListener(global.STATUS_REDIS),
	}
}

// Start Starts all services
func (m *Manager) Start() error {
	// Start the device status listener
	if err := m.deviceListener.Start(); err != nil {
		return err
	}

	logrus.Info("All services started successfully")
	return nil
}

// Stop Stops all services
func (m *Manager) Stop() {
	logrus.Info("Stopping all services...")

	if err := m.deviceListener.Stop(); err != nil {
		logrus.WithError(err).Error("Failed to stop device listener")
	}

	logrus.Info("All services have been stopped")
}
