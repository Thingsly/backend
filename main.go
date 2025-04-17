package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"

	"time"

	"github.com/Thingsly/backend/initialize"
	"github.com/Thingsly/backend/initialize/croninit"
	"github.com/Thingsly/backend/internal/app"
	"github.com/Thingsly/backend/internal/query"
	"github.com/Thingsly/backend/mqtt"
	"github.com/Thingsly/backend/mqtt/device"
	"github.com/Thingsly/backend/mqtt/publish"
	"github.com/Thingsly/backend/mqtt/subscribe"
	grpc_tptodb "github.com/Thingsly/backend/third_party/grpc/tptodb_client"

	router "github.com/Thingsly/backend/router"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	initialize.ViperInit("./configs/conf.yml")
	//initialize.ViperInit("./configs/conf-localdev.yml")
	initialize.RsaDecryptInit("./configs/rsa_key/private_key.pem")
	initialize.LogInIt()
	db, err := initialize.PgInit()
	if err != nil {
		logrus.Fatal(err)
	}
	_, err = initialize.RedisInit()
	if err != nil {
		logrus.Fatal(err)
	}

	query.SetDefault(db)

	grpc_tptodb.GrpcTptodbInit()

	err = mqtt.MqttInit()
	if err != nil {
		logrus.Fatal(err)
	}
	go device.InitDeviceStatus()
	err = subscribe.SubscribeInit()
	if err != nil {
		logrus.Fatal(err)
	}
	publish.PublishInit()

	croninit.CronInit()
}

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
	manager := app.NewManager()
	if err := manager.Start(); err != nil {
		logrus.Fatalf("Failed to start services: %v", err)
	}
	defer manager.Stop()
	// gin.SetMode(gin.ReleaseMode)

	host, port := loadConfig()
	router := router.RouterInit()
	srv := initServer(host, port, router)

	go startServer(srv, host, port)

	gracefulShutdown(srv)

}

func loadConfig() (host, port string) {
	host = viper.GetString("service.http.host")
	if host == "" {
		host = "localhost"
		logrus.Println("Using default host:", host)
	}

	port = viper.GetString("service.http.port")
	if port == "" {
		port = "9999"
		logrus.Println("Using default port:", port)
	}

	return host, port
}

func initServer(host, port string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:         net.JoinHostPort(host, port),
		Handler:      handler,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}
}

func startServer(srv *http.Server, host, port string) {
	logrus.Println("Listening and serving HTTP on", host, ":", port)
	successInfo()
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logrus.Fatalf("listen: %s\n", err)
	}
}

func gracefulShutdown(srv *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logrus.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logrus.Fatal("Server Shutdown:", err)
	}
	logrus.Println("Server exiting")
}

func successInfo() {
	startTime := time.Now().Format("2006-01-02 15:04:05")

	fmt.Println("----------------------------------------")
	fmt.Println("        Thingsly Backend started successfully!")
	fmt.Println("----------------------------------------")
	fmt.Printf("Start time: %s\n", startTime)
	fmt.Println("Version: v1.0.0")
	fmt.Println("----------------------------------------")
	fmt.Println("Welcome to Thingsly IoT Platform!")
	fmt.Println("----------------------------------------")
}
