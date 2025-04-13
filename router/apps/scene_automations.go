package apps

import (
	"github.com/HustIoTPlatform/backend/internal/api"

	"github.com/gin-gonic/gin"
)

type SceneAutomations struct{}

func (*SceneAutomations) Init(Router *gin.RouterGroup) {
	url := Router.Group("scene_automations")
	{

		url.POST("", api.Controllers.SceneAutomationsApi.CreateSceneAutomations)

		url.DELETE(":id", api.Controllers.SceneAutomationsApi.DeleteSceneAutomations)

		url.PUT("", api.Controllers.SceneAutomationsApi.UpdateSceneAutomations)

		url.POST("switch/:id", api.Controllers.SceneAutomationsApi.SwitchSceneAutomations)

		url.GET("list", api.Controllers.SceneAutomationsApi.HandleSceneAutomationsByPage)

		url.GET("detail/:id", api.Controllers.SceneAutomationsApi.HandleSceneAutomations)

		url.GET("log", api.Controllers.SceneAutomationsApi.HandleSceneAutomationsLog)

		url.GET("alarm", api.Controllers.SceneAutomationsApi.HandleSceneAutomationsWithAlarmByPage)

	}
}
