package apps

import (
	"github.com/HustIoTPlatform/backend/internal/api"

	"github.com/gin-gonic/gin"
)

type OperationLog struct{}

func (*OperationLog) Init(Router *gin.RouterGroup) {
	url := Router.Group("operation_logs")
	{

		url.GET("", api.Controllers.OperationLogsApi.HandleListByPage)
	}
}
