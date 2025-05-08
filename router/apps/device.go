package apps

import (
	"github.com/Thingsly/backend/internal/api"

	"github.com/gin-gonic/gin"
)

type Device struct{}

func (*Device) InitDevice(Router *gin.RouterGroup) {

	deviceapi := Router.Group("device")
	{

		deviceapi.POST("", api.Controllers.DeviceApi.CreateDevice)

		deviceapi.DELETE(":id", api.Controllers.DeviceApi.DeleteDevice)

		deviceapi.PUT("", api.Controllers.DeviceApi.UpdateDevice)

		deviceapi.PUT("active", api.Controllers.DeviceApi.ActiveDevice)

		deviceapi.GET("/detail/:id", api.Controllers.DeviceApi.HandleDeviceByID)

		deviceapi.GET("", api.Controllers.DeviceApi.HandleDeviceListByPage)

		deviceapi.GET("check/:deviceNumber", api.Controllers.DeviceApi.CheckDeviceNumber)

		deviceapi.GET("tenant/list", api.Controllers.DeviceApi.HandleTenantDeviceList)

		deviceapi.GET("list", api.Controllers.DeviceApi.HandleDeviceList)

		deviceapi.POST("son/add", api.Controllers.DeviceApi.CreateSonDevice)

		deviceapi.GET("connect/form", api.Controllers.DeviceApi.DeviceConnectForm)

		deviceapi.GET("connect/info", api.Controllers.DeviceApi.DeviceConnect)

		deviceapi.POST("update/voucher", api.Controllers.DeviceApi.UpdateDeviceVoucher)

		deviceapi.GET("sub-list/:id", api.Controllers.DeviceApi.HandleSubList)

		deviceapi.PUT("sub-remove", api.Controllers.DeviceApi.RemoveSubDevice)

		deviceapi.GET("metrics/:id", api.Controllers.DeviceApi.HandleMetrics)

		deviceapi.GET("metrics/menu", api.Controllers.DeviceApi.HandleActionByDeviceID)

		deviceapi.GET("metrics/condition/menu", api.Controllers.DeviceApi.HandleConditionByDeviceID)

		deviceapi.GET("map/telemetry/:id", api.Controllers.DeviceApi.HandleMapTelemetry)

		deviceapi.PUT("update/config", api.Controllers.DeviceApi.UpdateDeviceConfig)

		deviceapi.GET("online/status/:id", api.Controllers.DeviceApi.HandleDeviceOnlineStatus)

		deviceapi.POST("service/access/batch", api.Controllers.DeviceApi.CreateDeviceBatch)

		deviceapi.GET("/metrics/chart", api.Controllers.DeviceApi.HandleDeviceMetricsChart)

		deviceapi.GET("/selector", api.Controllers.DeviceApi.HandleDeviceSelector)

		// The telemetry data of the three devices with the latest data from the tenant
		deviceapi.GET("/telemetry/latest", api.Controllers.DeviceApi.HandleTenantTelemetryData)
	}

	deviceTemplateapi := deviceapi.Group("template")
	{

		deviceTemplateapi.POST("", api.Controllers.DeviceApi.CreateDeviceTemplate)

		deviceTemplateapi.DELETE(":id", api.Controllers.DeviceApi.DeleteDeviceTemplate)

		deviceTemplateapi.PUT("", api.Controllers.DeviceApi.UpdateDeviceTemplate)

		deviceTemplateapi.GET("/detail/:id", api.Controllers.DeviceApi.HandleDeviceTemplateById)

		deviceTemplateapi.GET("", api.Controllers.DeviceApi.HandleDeviceTemplateListByPage)

		deviceTemplateapi.GET("/menu", api.Controllers.DeviceApi.HandleDeviceTemplateMenu)

		deviceTemplateapi.GET("/chart", api.Controllers.DeviceApi.HandleDeviceTemplateByDeviceId)

		deviceTemplateapi.GET("/chart/select", api.Controllers.DeviceApi.HandleDeviceTemplateChartSelect)
	}

	deviceGroupapi := deviceapi.Group("group")
	{

		deviceGroupapi.POST("", api.Controllers.DeviceApi.CreateDeviceGroup)

		deviceGroupapi.DELETE(":id", api.Controllers.DeviceApi.DeleteDeviceGroup)

		deviceGroupapi.PUT("", api.Controllers.DeviceApi.UpdateDeviceGroup)

		deviceGroupapi.GET("", api.Controllers.DeviceApi.HandleDeviceGroupByPage)

		deviceGroupapi.GET("tree", api.Controllers.DeviceApi.HandleDeviceGroupByTree)

		deviceGroupapi.GET("detail/:id", api.Controllers.DeviceApi.HandleDeviceGroupByDetail)
	}

	deviceGroupRapi := deviceGroupapi.Group("relation")
	{

		deviceGroupRapi.POST("", api.Controllers.DeviceApi.CreateDeviceGroupRelation)

		deviceGroupRapi.DELETE("", api.Controllers.DeviceApi.DeleteDeviceGroupRelation)

		deviceGroupRapi.GET("list", api.Controllers.DeviceApi.HandleDeviceGroupRelation)

		deviceGroupRapi.GET("", api.Controllers.DeviceApi.HandleDeviceGroupListByDeviceId)

	}

	deviceModelApi := deviceapi.Group("model")
	{

		deviceModelApi.GET("source/at/list", api.Controllers.DeviceModelApi.HandleModelSourceAT)
		deviceModelTelemetryApi := deviceModelApi.Group("telemetry")
		{
			deviceModelTelemetryApi.POST("", api.Controllers.DeviceModelApi.CreateDeviceModelTelemetry)
			deviceModelTelemetryApi.DELETE(":id", api.Controllers.DeviceModelApi.DeleteDeviceModelGeneral)
			deviceModelTelemetryApi.PUT("", api.Controllers.DeviceModelApi.UpdateDeviceModelGeneral)
			deviceModelTelemetryApi.GET("", api.Controllers.DeviceModelApi.HandleDeviceModelGeneral)
		}

		deviceModelAttributesApi := deviceModelApi.Group("attributes")
		{
			deviceModelAttributesApi.POST("", api.Controllers.DeviceModelApi.CreateDeviceModelAttributes)
			deviceModelAttributesApi.DELETE(":id", api.Controllers.DeviceModelApi.DeleteDeviceModelGeneral)
			deviceModelAttributesApi.PUT("", api.Controllers.DeviceModelApi.UpdateDeviceModelGeneral)
			deviceModelAttributesApi.GET("", api.Controllers.DeviceModelApi.HandleDeviceModelGeneral)
		}

		deviceModelEventsApi := deviceModelApi.Group("events")
		{
			deviceModelEventsApi.POST("", api.Controllers.DeviceModelApi.CreateDeviceModelEvents)
			deviceModelEventsApi.DELETE(":id", api.Controllers.DeviceModelApi.DeleteDeviceModelGeneral)
			deviceModelEventsApi.PUT("", api.Controllers.DeviceModelApi.UpdateDeviceModelGeneralV2)
			deviceModelEventsApi.GET("", api.Controllers.DeviceModelApi.HandleDeviceModelGeneral)
		}

		deviceModelCommandsApi := deviceModelApi.Group("commands")
		{
			deviceModelCommandsApi.POST("", api.Controllers.DeviceModelApi.CreateDeviceModelCommands)
			deviceModelCommandsApi.DELETE(":id", api.Controllers.DeviceModelApi.DeleteDeviceModelGeneral)
			deviceModelCommandsApi.PUT("", api.Controllers.DeviceModelApi.UpdateDeviceModelGeneralV2)
			deviceModelCommandsApi.GET("", api.Controllers.DeviceModelApi.HandleDeviceModelGeneral)

		}

		deviceModelCustomCommandsApi := deviceModelApi.Group("custom/commands")
		{
			deviceModelCustomCommandsApi.POST("", api.Controllers.DeviceModelApi.CreateDeviceModelCustomCommands)
			deviceModelCustomCommandsApi.DELETE(":id", api.Controllers.DeviceModelApi.DeleteDeviceModelCustomCommands)
			deviceModelCustomCommandsApi.PUT("", api.Controllers.DeviceModelApi.UpdateDeviceModelCustomCommands)
			deviceModelCustomCommandsApi.GET("", api.Controllers.DeviceModelApi.HandleDeviceModelCustomCommandsByPage)
			deviceModelCustomCommandsApi.GET(":deviceId", api.Controllers.DeviceModelApi.HandleDeviceModelCustomCommandsByDeviceId)
		}

		deviceModelCustomControlApi := deviceModelApi.Group("custom/control")
		{
			deviceModelCustomControlApi.POST("", api.Controllers.DeviceModelApi.CreateDeviceModelCustomControl)
			deviceModelCustomControlApi.DELETE(":id", api.Controllers.DeviceModelApi.DeleteDeviceModelCustomControl)
			deviceModelCustomControlApi.PUT("", api.Controllers.DeviceModelApi.UpdateDeviceModelCustomControl)
			deviceModelCustomControlApi.GET("", api.Controllers.DeviceModelApi.HandleDeviceModelCustomControl)
		}

	}
}
