package apps

import (
	"github.com/HustIoTPlatform/backend/internal/api"

	"github.com/gin-gonic/gin"
)

type NotificationHistoryGroup struct {
}

func (*NotificationHistoryGroup) InitNotificationHistory(Router *gin.RouterGroup) {
	url := Router.Group("notification_history")
	{

		url.GET("/list", api.Controllers.NotificationHistoryApi.HandleNotificationHistoryListByPage)

	}
}
