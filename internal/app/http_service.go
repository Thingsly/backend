package app

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	router "github.com/Thingsly/backend/router"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// HTTPService Implement the HTTP service
type HTTPService struct {
	server *http.Server
	config *HTTPConfig
}

// HTTPConfig Save the HTTP service configuration
type HTTPConfig struct {
	Host            string
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

// NewHTTPService Create a new HTTP service
func NewHTTPService() *HTTPService {
	return &HTTPService{
		config: &HTTPConfig{
			Host:            "localhost",
			Port:            "9999",
			ReadTimeout:     60 * time.Second,
			WriteTimeout:    60 * time.Second,
			ShutdownTimeout: 5 * time.Second,
		},
	}
}

// Name Return the service name
func (s *HTTPService) Name() string {
	return "HTTP service"
}

// SetConfig Set the HTTP service configuration
func (s *HTTPService) SetConfig(host, port string, readTimeout, writeTimeout, shutdownTimeout time.Duration) {
	s.config.Host = host
	s.config.Port = port
	s.config.ReadTimeout = readTimeout
	s.config.WriteTimeout = writeTimeout
	s.config.ShutdownTimeout = shutdownTimeout
}

// Start Start the HTTP service
func (s *HTTPService) Start() error {
	// Load the host and port from the configuration
	host := viper.GetString("service.http.host")
	if host == "" {
		host = s.config.Host
		logrus.Debugf("Using default host: %s", host)
	}

	port := viper.GetString("service.http.port")
	if port == "" {
		port = s.config.Port
		logrus.Debugf("Using default port: %s", port)
	}

	// Initialize the router
	handler := router.RouterInit()

	// Create the server
	s.server = &http.Server{
		Addr:         net.JoinHostPort(host, port),
		Handler:      handler,
		ReadTimeout:  s.config.ReadTimeout,
		WriteTimeout: s.config.WriteTimeout,
	}

	// Start the server asynchronously
	go func() {
		logrus.Infof("HTTP service is listening on %s:%s", host, port)

		// Print the success information
		successInfo()

		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Errorf("HTTP server error: %v", err)
		}
	}()

	return nil
}

// Stop Stop the HTTP service
func (s *HTTPService) Stop() error {
	if s.server == nil {
		return nil
	}

	logrus.Info("Stopping HTTP service...")
	ctx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("HTTP service graceful shutdown failed: %v", err)
	}

	logrus.Info("HTTP service stopped")
	return nil
}

// WithHTTPService Add the HTTP service to the application
func WithHTTPService() Option {
	return func(app *Application) error {
		service := NewHTTPService()
		app.RegisterService(service)
		return nil
	}
}

// Print the success information
func successInfo() {
	// Get the current time
	startTime := time.Now().Format("2006-01-02 15:04:05")

	// Print the success message
	fmt.Println("----------------------------------------")
	// fmt.Println("   ⣿⣿⣿⣿⣿⣿⣿⣿⣿⣻⣿⣿⣿⡿⢿⡿⠿⠿⣿⣿⣿⣿⣿⣿⡿⣿⣿⣿⡿⣿⣿")
	// fmt.Println("   ⣿⣿⣿⣿⠿⠿⢿⣿⣿⠟⣋⣭⣶⣶⣞⣿⣶⣶⣶⣬⣉⠻⣿⣿⡿⣋⣉⠻⣿⣿⣿")
	// fmt.Println("   ⣿⢻⣿⠃⣤⣤⣢⣍⣴⣿⢋⣵⣿⣿⣿⣿⣷⣶⣝⣿⣿⣧⣄⢉⣜⣥⣜⢷⢹⢇⢛")
	// fmt.Println("   ⡏⡦⡁⡸⢛⡴⢣⣾⢟⣿⣿⣿⢟⣾⣧⣙⢿⣿⣿⣿⣿⣿⣿⣿⢩⢳⣞⢿⡏⢷⣾")
	// fmt.Println("   ⣷⣵⡇⣗⡾⢁⣾⣟⣾⣿⡿⣻⣾⣿⣿⣿⡎⠛⡛⢿⣿⡟⣿⣿⡜⡜⢿⡌⠇⢾⣿")
	// fmt.Println("   ⣿⣿⠁⣾⠏⣾⣿⣿⣽⣑⣺⣥⣿⣿⣿⣿⣷⣶⣦⣖⢝⢿⣿⣿⣿⡀⠹⣿⣼⢸⣿")
	// fmt.Println("   ⣿⣿⢰⡏⢡⣿⣿⠐⣵⠿⠛⠛⣿⣿⣿⣿⣿⠍⠚⢙⠻⢦⣼⣿⣿⠁⣄⣿⣿⠘⣿")
	// fmt.Println("   ⣿⣿⢸⢹⢈⣿⣿⠘⣡⡞⠉⡀⢻⣿⣿⣿⣿⢃⠠⢈⢳⣌⣩⣿⣿⠰⠿⢼⣿⠀⣿")
	// fmt.Println("   ⣿⠿⣘⠯⠌⡟⣿⡟⣾⣇⢾⡵⣹⣟⣿⣿⣿⣮⣓⣫⣿⣟⢿⣿⢿⡾⡹⢆⣦⣤⢹")
	// fmt.Println("   ⣅⣛⠶⠽⣧⣋⠳⡓⢿⣿⣿⣿⣿⣿⢿⣿⣿⣿⣿⣿⣿⣫⣸⠏⡋⠷⣛⣫⡍⣶⣿")
	// fmt.Println("   ⣿⡿⢸⢳⣶⣶⠀⡇⣬⡛⠿⣿⣿⣿⣿⣿⣿⣿⠿⢟⣉⣕⡭⠀⢺⣸⣽⢻⡅⣿⣿")
	// fmt.Println("   ⣿⡇⣾⡾⣰⡯⠀⡗⣯⣿⣽⡶⠶⠂⢠⣾⣿⠐⠚⠻⢯⣿⠇⠎⡀⣳⣿⣼⡃⣿⣿")
	// fmt.Println("   ⣿⡇⣟⣇⡟⣧⠀⡗⣿⣿⡿⢡⢖⣀⠼⢟⣻⣤⣔⢦⢸⣿⢀⢆⢡⣿⣯⢹⡃⣿⣿")
	// fmt.Println("   ⣿⡇⡏⣿⡾⣸⣿⣇⠸⠟⣋⣼⣼⣿⢻⣿⣿⢿⣟⢾⣌⠫⠈⣶⣿⡿⣩⡿⢃⣿⣿")
	// fmt.Println("   ⣿⣷⡀⠻⡷⢪⢧⡙⠰⣾⣿⣿⣾⡽⣾⣿⡿⣺⣵⣾⣿⡇⡜⣽⠟⢷⣪⣴⣿⣿⣿")
	// fmt.Println("   ⣿⣿⣿⣾⣿⠏⣤⡁⣷⣽⣿⣿⣿⣿⣷⣶⣿⣿⣿⣿⣿⣱⠸⣱⣦⠙⣿⣿⣿⣾⣿")
	fmt.Println("   ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣼⣿⡄⠀⠀⠀⠀⠀⠀⠀⣠⣄⠀⠀⠀⠀⠀")
	fmt.Println("   ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣼⣿⣿⡇⠀⠀⠀⠀⠀⠀⢰⣿⣿⡄⠀⠀⠀⠀")
	fmt.Println("   ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢰⣿⣿⣿⣿⠀⠀⠀⠀⠀⠀⣿⣿⣿⡇⠀⠀⠀⠀")
	fmt.Println("   ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢘⣿⣿⣿⣯⣤⣤⣤⣀⣀⣸⣿⣿⣿⡇⠀⠀⠀⠀")
	fmt.Println("   ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣼⡿⠿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡏⠀⠀⠀⠀⠀")
	fmt.Println("   ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣰⣿⡏⢠⡦⠈⣿⣿⣿⣿⣿⣿⠟⠛⢻⣷⡄⠀⠀⠀⠀")
	fmt.Println("   ⠀⠀⠀⠀⠀⠀⡀⠀⠀⠀⠀⠀⠀⢠⣿⣿⣿⣦⣤⣾⣿⣿⣿⣿⣿⣿⠀⠿⢀⣿⣷⠄⠀⠀⠀")
	fmt.Println("   ⢠⣄⠀⠀⠀⣼⣿⡆⠀⠀⠀⠀⣰⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣶⣿⣿⣿⡇⠀⠀⠀")
	fmt.Println("   ⢸⣿⣷⣄⢀⣿⣿⣿⠀⠀⠀⢀⣿⣿⣿⠿⠋⠉⠁⠀⠀⠈⠉⠉⠻⢿⣿⣿⣿⣿⣿⣷⠀⠀⠀")
	fmt.Println("   ⠀⣿⣿⠿⣿⣿⡿⣛⢷⠀⠀⢸⣿⣿⠏⢀⣤⣄⠀⣠⣤⡄⠀⠀⠀⠀⢻⣿⣿⣿⣿⣿⣦⣄⠀")
	fmt.Println("   ⠀⣿⣇⣀⣽⣿⣷⣤⣾⣧⠀⠘⣿⠏⠀⠛⠋⠙⠀⠛⠙⠛⠀⠾⠿⣷⢸⣿⣿⣿⣿⣿⣿⣿⡇")
	fmt.Println("   ⢠⣿⣿⣿⣿⣿⣿⣿⣿⣿⡆⠀⠋⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⢿⡿⣿⣿⣿⣿⣿⡇")
	fmt.Println("   ⠘⣿⣿⣿⣿⣿⣿⣿⣿⣿⠇⠀⠐⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣿⣿⣿⣿⡿⠁")
	fmt.Println("   ⠀⢻⣿⣿⣿⣿⣿⣿⣿⡟⠀⠀⠀⠈⠢⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣼⣿⣿⡿⠋⠀⠀")
	fmt.Println("   ⠀⠀⠉⠛⠛⠛⠛⠛⠛⠁⠀⠀⠀⠀⠀⠘⠻⢲⠦⠤⠤⠀⠀⠀⠀⣤⢴⡿⠟⠁⠀⠀⠀⠀⠀")
	fmt.Println("")
	fmt.Println("  ______  __  __  ______   __  __  ____    ____    __       __    __")
	fmt.Println(" /\\__  _\\/\\ \\/\\ \\/\\__  _\\ /\\ \\/\\ \\/\\  _`\\ /\\  _`\\ /\\ \\     /\\ \\  /\\ \\")
	fmt.Println(" \\/_/\\ \\/\\ \\ \\_\\ \\/_/\\ \\/ \\ \\ `\\\\ \\ \\ \\L\\_\\ \\,\\L\\_\\ \\ \\    \\ `\\`\\\\/'/")
	fmt.Println("    \\ \\ \\ \\ \\  _  \\ \\ \\ \\  \\ \\ , ` \\ \\ \\L_L\\/_\\__ \\\\ \\ \\  __`\\ `\\ /'")
	fmt.Println("     \\ \\ \\ \\ \\ \\ \\ \\ \\_\\ \\__\\ \\ \\`\\ \\ \\ \\/, \\/\\ \\L\\ \\ \\ \\L\\ \\ `\\ \\ \\")
	fmt.Println("      \\ \\_\\ \\ \\_\\ \\_\\/\\_____\\\\ \\_\\ \\_\\ \\____/\\ `\\____\\ \\____/   \\ \\_\\")
	fmt.Println("       \\/_/  \\/_/\\/_/\\/_____/ \\/_/\\/_/\\/___/  \\/_____/\\/___/     \\/_/")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("        Thingsly Backend Started Successfully!")
	fmt.Println("----------------------------------------")
	fmt.Printf("Start Time: %s\n", startTime)
	fmt.Println("Version: v1.1.1")
	fmt.Println("----------------------------------------")
	fmt.Println("Welcome to Thingsly Backend!")
	fmt.Println("For help, visit: https://dangky.app")
	fmt.Println("----------------------------------------")
}
