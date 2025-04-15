package apps

import (
	"github.com/Thingsly/backend/internal/api"
	"github.com/gin-gonic/gin"
)

type MessagePush struct {
}

func (*MessagePush) Init(Router *gin.RouterGroup) {
	url := Router.Group("message_push")
	{

		url.POST("", api.Controllers.MessagePushApi.CreateMessagePush)

		url.POST("/logout", api.Controllers.MessagePushApi.MessagePushMangeLogout)

		url.GET("/config", api.Controllers.MessagePushApi.GetMessagePushConfig)

		url.POST("/config", api.Controllers.MessagePushApi.SetMessagePushConfig)
	}
}
