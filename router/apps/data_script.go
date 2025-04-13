package apps

import (
	"github.com/HustIoTPlatform/backend/internal/api"

	"github.com/gin-gonic/gin"
)

type DataScript struct {
}

func (*DataScript) Init(Router *gin.RouterGroup) {
	url := Router.Group("data_script")
	{

		url.POST("", api.Controllers.DataScriptApi.CreateDataScript)

		url.DELETE(":id", api.Controllers.DataScriptApi.DeleteDataScript)

		url.PUT("", api.Controllers.DataScriptApi.UpdateDataScript)

		url.GET("", api.Controllers.DataScriptApi.HandleDataScriptListByPage)

		url.POST("quiz", api.Controllers.DataScriptApi.QuizDataScript)

		url.PUT("enable", api.Controllers.DataScriptApi.EnableDataScript)
	}
}
