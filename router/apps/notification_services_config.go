package apps

import (
	"github.com/HustIoTPlatform/backend/internal/api"

	"github.com/gin-gonic/gin"
)

type NotificationServicesConfig struct{}

func (*NotificationServicesConfig) Init(Router *gin.RouterGroup) {
	url := Router.Group("notification/services/config")
	{

		url.POST("", api.Controllers.NotificationServicesConfigApi.SaveNotificationServicesConfig)

		url.GET(":type", api.Controllers.NotificationServicesConfigApi.HandleNotificationServicesConfig)

		url.POST("e-mail/test", api.Controllers.NotificationServicesConfigApi.SendTestEmail)
	}
}
