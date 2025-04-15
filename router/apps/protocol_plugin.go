package apps

import (
	"github.com/Thingsly/backend/internal/api"

	"github.com/gin-gonic/gin"
)

type ProtocolPlugin struct {
}

func (*ProtocolPlugin) InitProtocolPlugin(Router *gin.RouterGroup) {
	protocolPluginApi := Router.Group("protocol_plugin")
	{

		protocolPluginApi.POST("", api.Controllers.ProtocolPluginApi.CreateProtocolPlugin)

		protocolPluginApi.DELETE(":id", api.Controllers.ProtocolPluginApi.DeleteProtocolPlugin)

		protocolPluginApi.PUT("", api.Controllers.ProtocolPluginApi.UpdateProtocolPlugin)

		protocolPluginApi.GET("", api.Controllers.ProtocolPluginApi.HandleProtocolPluginListByPage)

		protocolPluginApi.GET("device_config_form", api.Controllers.ProtocolPluginApi.HandleProtocolPluginForm)

		protocolPluginApi.GET("config_form", api.Controllers.ProtocolPluginApi.HandleProtocolPluginFormByProtocolType)
	}
}
