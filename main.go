package main

import (
	"fmt"
	"os"

	"github.com/Thingsly/backend/internal/app"

	"github.com/sirupsen/logrus"
)

// @title           Thingsly API
// @version         1.0
// @description     Thingsly API.
// @schemes         http
// @host      localhost:9999
// @BasePath
// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        x-token
func main() {
	// Initialize the Application structure
	application, err := app.NewApplication(
		// Basic configuration
		app.WithProductionConfig(),
		app.WithRsaDecrypt("./configs/rsa_key/private_key.pem"),
		app.WithLogger(),
		app.WithDatabase(),
		app.WithRedis(),

		// Services
		app.WithHTTPService(),
		app.WithGRPCService(),
		app.WithMQTTService(),
		app.WithCronService(),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Application initialization failed: %v\n", err)
		os.Exit(1)
	}

	// Start all services
	if err := application.Start(); err != nil {
		logrus.Fatalf("Failed to start services: %v", err)
	}

	// Wait for the services to run and handle exit
	application.Wait()

	// Automatically call the Shutdown method when the application is closed
	defer application.Shutdown()
}