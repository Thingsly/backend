package app

import (
	tptodb "github.com/Thingsly/backend/third_party/grpc/tptodb_client"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// GRPCService Implement the gRPC client service
type GRPCService struct {
	initialized bool
}

// NewGRPCService Create a gRPC service instance
func NewGRPCService() *GRPCService {
	return &GRPCService{
		initialized: false,
	}
}

// Name Return the service name
func (s *GRPCService) Name() string {
	return "gRPC client service"
}

// Start Start the gRPC service
func (s *GRPCService) Start() error {

	// Check if gRPC is enabled
	if !viper.GetBool("grpc.enabled") {
		logrus.Info("gRPC client service is disabled, skipping initialization")
		return nil
	}

	logrus.Info("Initializing gRPC client...")

	// Initialize the gRPC client
	tptodb.GrpcTptodbInit()

	s.initialized = true
	logrus.Info("gRPC client initialized")
	return nil
}

// Stop Stop the gRPC service
func (s *GRPCService) Stop() error {
	if !s.initialized {
		return nil
	}

	logrus.Info("Stopping gRPC client...")
	// If there is a close method, call it here

	logrus.Info("gRPC client stopped")
	return nil
}

// WithGRPCService Add the gRPC service to the application
func WithGRPCService() Option {
	return func(app *Application) error {
		service := NewGRPCService()
		app.RegisterService(service)
		return nil
	}
}
