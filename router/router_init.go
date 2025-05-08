package router

import (
	"time"

	middleware "github.com/Thingsly/backend/internal/middleware"
	"github.com/Thingsly/backend/internal/middleware/response"
	"github.com/Thingsly/backend/pkg/global"
	"github.com/Thingsly/backend/pkg/metrics"
	"github.com/Thingsly/backend/router/apps"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"

	// gin-swagger middleware
	_ "github.com/Thingsly/backend/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	api "github.com/Thingsly/backend/internal/api"
	service "github.com/Thingsly/backend/internal/service"
)

// swagger embed files

func RouterInit() *gin.Engine {
	//gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	m := metrics.NewMetrics("Thingsly")
	// Create a memory storage implementation
	memStorage := metrics.NewMemoryStorage()
	// Set the storage implementation
	m.SetHistoryStorage(memStorage)

	m.StartMetricsCollection(15 * time.Second)

	router.Use(middleware.MetricsMiddleware(m))

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Set the metrics manager to the system monitoring service
	service.SetMetricsManager(m)

	router.StaticFile("/metrics-viewer", "./static/metrics-viewer.html")

	router.GET("/files/*filepath", func(c *gin.Context) {
		filepath := c.Param("filepath")
		c.File("./files" + filepath)
	})

	//router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Use(middleware.Cors())

	handler, err := response.NewHandler("configs/messages.yaml", "configs/messages_str.yaml")
	if err != nil {
		logrus.Fatalf("Failed to initialize response handler: %v", err)
	}

	router.Use(middleware.OperationLogs())

	global.ResponseHandler = handler

	router.Use(handler.Middleware())

	controllers := new(api.Controller)

	router.GET("/health", controllers.SystemApi.HealthCheck)

	api := router.Group("api")
	{

		v1 := api.Group("v1")
		{
			v1.POST("plugin/heartbeat", controllers.Heartbeat)
			v1.POST("plugin/device/config", controllers.HandleDeviceConfigForProtocolPlugin)
			v1.POST("plugin/service/access/list", controllers.HandlePluginServiceAccessList)
			v1.POST("plugin/service/access", controllers.HandlePluginServiceAccess)
			v1.POST("login", controllers.Login)
			v1.GET("verification/code", controllers.HandleVerificationCode)
			v1.POST("reset/password", controllers.ResetPassword)
			v1.GET("logo", controllers.HandleLogoList)

			v1.GET("telemetry/datas/current/ws", controllers.TelemetryDataApi.ServeCurrentDataByWS)

			v1.GET("device/online/status/ws", controllers.TelemetryDataApi.ServeDeviceStatusByWS)

			v1.GET("telemetry/datas/current/keys/ws", controllers.TelemetryDataApi.ServeCurrentDataByKey)
			v1.GET("ota/download/files/upgradePackage/:path/:file", controllers.OTAApi.DownloadOTAUpgradePackage)

			v1.GET("systime", controllers.SystemApi.HandleSystime)

			v1.GET("sys_function", controllers.SysFunctionApi.HandleSysFcuntion)

			v1.POST("/tenant/email/register", controllers.UserApi.EmailRegister)

			v1.POST("/device/gateway-register", controllers.DeviceApi.GatewayRegister)

			v1.POST("/device/gateway-sub-register", controllers.DeviceApi.GatewaySubRegister)

			// Get system version
			v1.GET("sys_version", controllers.SystemApi.HandleSysVersion)
		}

		v1.Use(middleware.JWTAuth())

		v1.Use(middleware.CasbinRBAC())

		SSERouter(v1)

		{
			apps.Model.User.InitUser(v1)

			apps.Model.Role.Init(v1)

			apps.Model.Casbin.Init(v1)

			apps.Model.Dict.InitDict(v1)

			apps.Model.OTA.InitOTA(v1)

			apps.Model.UpLoad.Init(v1)

			apps.Model.ProtocolPlugin.InitProtocolPlugin(v1)

			apps.Model.Device.InitDevice(v1)

			apps.Model.UiElements.Init(v1)

			apps.Model.Board.InitBoard(v1)

			apps.Model.EventData.InitEventData(v1)

			apps.Model.TelemetryData.InitTelemetryData(v1)

			apps.Model.AttributeData.InitAttributeData(v1)

			apps.Model.CommandData.InitCommandData(v1)

			apps.Model.OperationLog.Init(v1)

			apps.Model.Logo.Init(v1)

			apps.Model.DataPolicy.Init(v1)

			apps.Model.DeviceConfig.Init(v1)

			apps.Model.DataScript.Init(v1)

			apps.Model.NotificationGroup.InitNotificationGroup(v1)

			apps.Model.NotificationHistoryGroup.InitNotificationHistory(v1)

			apps.Model.NotificationServicesConfig.Init(v1)

			apps.Model.Alarm.Init(v1)

			apps.Model.Scene.Init(v1)

			apps.Model.SceneAutomations.Init(v1)

			apps.Model.SysFunction.Init(v1)

			apps.Model.ServicePlugin.Init(v1)

			apps.Model.ExpectedData.InitExpectedData(v1)

			apps.Model.OpenAPIKey.InitOpenAPIKey(v1)

			apps.Model.MessagePush.Init(v1)

			// Initialize system monitoring routes
			apps.Model.SystemMonitor.InitSystemMonitor(v1, m)
		}
	}

	return router
}
