package apps

import (
	"github.com/HustIoTPlatform/backend/internal/api"

	"github.com/gin-gonic/gin"
)

type ServicePlugin struct {
}

func (*ServicePlugin) Init(Router *gin.RouterGroup) {
	url := Router.Group("service")
	{
		url.POST("", api.Controllers.ServicePluginApi.Create)

		url.GET("list", api.Controllers.ServicePluginApi.HandleList)

		url.GET("/detail/:id", api.Controllers.ServicePluginApi.Handle)

		url.PUT("", api.Controllers.ServicePluginApi.Update)

		url.DELETE(":id", api.Controllers.ServicePluginApi.Delete)

		url.GET("/plugin/select", api.Controllers.ServicePluginApi.HandleServiceSelect)

		url.GET("/plugin/info", api.Controllers.ServicePluginApi.HandleServicePluginByServiceIdentifier)

		access := url.Group("access")
		access.POST("", api.Controllers.ServiceAccessApi.Create)

		access.GET("/list", api.Controllers.ServiceAccessApi.HandleList)

		access.PUT("", api.Controllers.ServiceAccessApi.Update)

		access.DELETE(":id", api.Controllers.ServiceAccessApi.Delete)

		access.GET("/voucher/form", api.Controllers.ServiceAccessApi.HandleVoucherForm)

		access.GET("/device/list", api.Controllers.ServiceAccessApi.HandleDeviceList)
	}
}
