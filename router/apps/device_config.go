package apps

import (
	"github.com/HustIoTPlatform/backend/internal/api"

	"github.com/gin-gonic/gin"
)

type DeviceConfig struct {
}

func (*DeviceConfig) Init(Router *gin.RouterGroup) {
	url := Router.Group("device_config")
	{

		url.POST("", api.Controllers.DeviceConfigApi.CreateDeviceConfig)

		url.DELETE(":id", api.Controllers.DeviceConfigApi.DeleteDeviceConfig)

		url.PUT("", api.Controllers.DeviceConfigApi.UpdateDeviceConfig)

		url.GET("", api.Controllers.DeviceConfigApi.HandleDeviceConfigListByPage)

		url.GET("menu", api.Controllers.DeviceConfigApi.HandleDeviceConfigListMenu)

		url.GET("/:id", api.Controllers.DeviceConfigApi.HandleDeviceConfigById)

		url.PUT("batch", api.Controllers.DeviceConfigApi.BatchUpdateDeviceConfig)

		url.GET("connect", api.Controllers.DeviceConfigApi.HandleDeviceConfigConnect)

		url.GET("voucher_type", api.Controllers.DeviceConfigApi.HandleVoucherType)

		url.GET("metrics/menu", api.Controllers.DeviceConfigApi.HandleActionByDeviceConfigID)

		url.GET("metrics/condition/menu", api.Controllers.DeviceConfigApi.HandleConditionByDeviceConfigID)

	}
}
