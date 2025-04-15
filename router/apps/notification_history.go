package apps

import (
	"github.com/Thingsly/backend/internal/api"

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
