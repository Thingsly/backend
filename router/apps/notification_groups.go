package apps

import (
	"github.com/HustIoTPlatform/backend/internal/api"

	"github.com/gin-gonic/gin"
)

type NotificationGroup struct {
}

func (*NotificationGroup) InitNotificationGroup(Router *gin.RouterGroup) {
	url := Router.Group("notification_group")
	{

		url.POST("", api.Controllers.NotificationGroupApi.CreateNotificationGroup)

		url.DELETE("/:id", api.Controllers.NotificationGroupApi.DeleteNotificationGroup)

		url.PUT("/:id", api.Controllers.NotificationGroupApi.UpdateNotificationGroup)

		url.GET("/list", api.Controllers.NotificationGroupApi.HandleNotificationGroupListByPage)

		url.GET("/:id", api.Controllers.NotificationGroupApi.HandleNotificationGroupById)

	}
}
