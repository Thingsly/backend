package apps

import (
	"github.com/HustIoTPlatform/backend/internal/api"

	"github.com/gin-gonic/gin"
)

type Logo struct {
}

func (*Logo) Init(Router *gin.RouterGroup) {
	url := Router.Group("logo")
	{

		url.PUT("", api.Controllers.LogoApi.UpdateLogo)

		// url.GET("", api.Controllers.LogoApi.GetLogoList)
	}
}
