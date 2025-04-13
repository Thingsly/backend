package apps

import (
	"github.com/HustIoTPlatform/backend/internal/api"

	"github.com/gin-gonic/gin"
)

type CommandData struct{}

func (*CommandData) InitCommandData(Router *gin.RouterGroup) {
	commandDataApi := Router.Group("command/datas")
	{

		commandDataApi.GET("set/logs", api.Controllers.CommandSetLogApi.ServeSetLogsDataListByPage)

		commandDataApi.POST("pub", api.Controllers.CommandSetLogApi.CommandPutMessage)

		commandDataApi.GET(":id", api.Controllers.CommandSetLogApi.HandleCommandList)
	}
}
