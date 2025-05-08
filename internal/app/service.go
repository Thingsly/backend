package app

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// Service Define the methods that all service components must implement
type Service interface {
	// Name Return the service name
	Name() string

	// Start Start the service, if it fails, return an error
	Start() error

	// Stop Stop the service and clean up resources
	Stop() error
}

// ServiceManager Manage the startup and shutdown of multiple services
type ServiceManager struct {
	services []Service
	wg       sync.WaitGroup
	mu       sync.Mutex
	started  bool
}

// NewServiceManager Create a new service manager
func NewServiceManager() *ServiceManager {
	return &ServiceManager{
		services: make([]Service, 0),
	}
}

// RegisterService Register a service to the manager
func (m *ServiceManager) RegisterService(service Service) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.services = append(m.services, service)
	logrus.Infof("Service %s registered", service.Name())
}

// StartAll Start all registered services
func (m *ServiceManager) StartAll() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.started {
		return fmt.Errorf("service already started")
	}

	for _, service := range m.services {
		logrus.Infof("Starting service: %s", service.Name())
		if err := service.Start(); err != nil {
			return fmt.Errorf("failed to start service %s: %v", service.Name(), err)
		}
		m.wg.Add(1)
	}

	m.started = true
	return nil
}

// StopAll Stop all services, in the reverse order of registration
func (m *ServiceManager) StopAll() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.started {
		return
	}

	// Create a context with a timeout to ensure the stop operation does not block indefinitely
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Iterate through the service list in reverse order to ensure the stop operation is performed in the correct order
	for i := len(m.services) - 1; i >= 0; i-- {
		service := m.services[i]
		logrus.Infof("Stopping service: %s", service.Name())

		// Create a channel to receive the stop completion signal
		done := make(chan error, 1)

		go func(s Service) {
			done <- s.Stop()
			m.wg.Done()
		}(service)

		// Wait for the service to stop or timeout
		select {
		case err := <-done:
			if err != nil {
				logrus.Errorf("failed to stop service %s: %v", service.Name(), err)
			} else {
				logrus.Infof("service %s stopped successfully", service.Name())
			}
		case <-ctx.Done():
			logrus.Warnf("stop service %s timeout", service.Name())
		}
	}

	m.started = false
}

// Wait Wait for all services to complete
func (m *ServiceManager) Wait() {
	m.wg.Wait()
}