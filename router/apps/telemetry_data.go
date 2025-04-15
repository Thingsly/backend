package apps

import (
	"github.com/Thingsly/backend/internal/api"

	"github.com/gin-gonic/gin"
)

type TelemetryData struct{}

func (*TelemetryData) InitTelemetryData(Router *gin.RouterGroup) {
	telemetrydataapi := Router.Group("telemetry/datas")
	{

		telemetrydataapi.GET("current/:id", api.Controllers.TelemetryDataApi.HandleCurrentData)

		telemetrydataapi.GET("/current/keys", api.Controllers.TelemetryDataApi.HandleCurrentDataKeys)

		telemetrydataapi.GET("current/detail/:id", api.Controllers.TelemetryDataApi.ServeCurrentDetailData)

		telemetrydataapi.GET("history", api.Controllers.TelemetryDataApi.ServeHistoryData)

		telemetrydataapi.GET("history/pagination", api.Controllers.TelemetryDataApi.ServeHistoryDataByPage)

		telemetrydataapi.GET("history/page", api.Controllers.TelemetryDataApi.ServeHistoryDataByPage)

		telemetrydataapi.DELETE("", api.Controllers.TelemetryDataApi.DeleteData)

		telemetrydataapi.GET("statistic", api.Controllers.TelemetryDataApi.ServeStatisticData)

		telemetrydataapi.GET("set/logs", api.Controllers.TelemetryDataApi.ServeSetLogsDataListByPage)

		telemetrydataapi.POST("pub", api.Controllers.TelemetryDataApi.TelemetryPutMessage)

		telemetrydataapi.GET("simulation", api.Controllers.TelemetryDataApi.ServeEchoData)

		telemetrydataapi.POST("simulation", api.Controllers.TelemetryDataApi.SimulationTelemetryData)

		telemetrydataapi.GET("msg/count", api.Controllers.TelemetryDataApi.ServeMsgCountByTenant)

	}
}
