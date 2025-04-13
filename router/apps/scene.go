package apps

import (
	"github.com/HustIoTPlatform/backend/internal/api"

	"github.com/gin-gonic/gin"
)

type Scene struct{}

func (*Scene) Init(Router *gin.RouterGroup) {
	url := Router.Group("scene")
	{

		url.POST("", api.Controllers.SceneApi.CreateScene)

		url.DELETE(":id", api.Controllers.SceneApi.DeleteScene)

		url.GET("", api.Controllers.SceneApi.HandleSceneByPage)

		url.GET("/detail/:id", api.Controllers.SceneApi.HandleScene)

		url.PUT("", api.Controllers.SceneApi.UpdateScene)

		url.POST("active/:id", api.Controllers.SceneApi.ActiveScene)

		url.GET("log", api.Controllers.SceneApi.HandleSceneLog)

	}
}
