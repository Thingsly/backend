package apps

import (
	"github.com/Thingsly/backend/internal/api"

	"github.com/gin-gonic/gin"
)

type Alarm struct{}

func (*Alarm) Init(Router *gin.RouterGroup) {
	url := Router.Group("alarm")
	alarmconfig(url)
	alarminfo(url)
}

func alarmconfig(Router *gin.RouterGroup) {
	url := Router.Group("config")
	{

		url.POST("", api.Controllers.AlarmApi.CreateAlarmConfig)

		url.DELETE(":id", api.Controllers.AlarmApi.DeleteAlarmConfig)

		url.PUT("", api.Controllers.AlarmApi.UpdateAlarmConfig)

		url.GET("", api.Controllers.AlarmApi.ServeAlarmConfigListByPage)
	}
}

func alarminfo(Router *gin.RouterGroup) {
	url := Router.Group("info")
	{

		url.PUT("", api.Controllers.AlarmApi.UpdateAlarmInfo)

		url.PUT("batch", api.Controllers.AlarmApi.BatchUpdateAlarmInfo)

		url.GET("", api.Controllers.AlarmApi.HandleAlarmInfoListByPage)

		url.GET("history", api.Controllers.AlarmApi.HandleAlarmHisttoryListByPage)

		url.PUT("history", api.Controllers.AlarmApi.AlarmHistoryDescUpdate)

		url.GET("history/device", api.Controllers.AlarmApi.HandleDeviceAlarmStatus)

		url.GET("config/device", api.Controllers.AlarmApi.HandleConfigByDevice)

		url.GET("history/:id", api.Controllers.AlarmApi.HandleAlarmInfoHistory)
	}
}
