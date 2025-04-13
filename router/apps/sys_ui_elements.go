package apps

import (
	"github.com/HustIoTPlatform/backend/internal/api"

	"github.com/gin-gonic/gin"
)

type UiElements struct {
}

func (*UiElements) Init(Router *gin.RouterGroup) {
	url := Router.Group("ui_elements")
	{

		url.POST("", api.Controllers.UiElementsApi.CreateUiElements)

		url.DELETE(":id", api.Controllers.UiElementsApi.DeleteUiElements)

		url.PUT("", api.Controllers.UiElementsApi.UpdateUiElements)

		url.GET("", api.Controllers.UiElementsApi.ServeUiElementsListByPage)

		url.GET("menu", api.Controllers.UiElementsApi.ServeUiElementsListByAuthority)

		url.GET("select/form", api.Controllers.UiElementsApi.ServeUiElementsListByTenant)
	}
}
