package apps

import (
	"github.com/HustIoTPlatform/backend/internal/api"

	"github.com/gin-gonic/gin"
)

type UpLoad struct{}

func (*UpLoad) Init(Router *gin.RouterGroup) {
	uploadapi := Router.Group("file")
	{

		uploadapi.POST("up", api.Controllers.UpLoadApi.UpFile)
	}
}
