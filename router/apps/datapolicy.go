package apps

import (
	"github.com/HustIoTPlatform/backend/internal/api"

	"github.com/gin-gonic/gin"
)

type DataPolicy struct {
}

func (*DataPolicy) Init(Router *gin.RouterGroup) {
	url := Router.Group("datapolicy")
	{

		url.PUT("", api.Controllers.DataPolicyApi.UpdateDataPolicy)

		url.GET("", api.Controllers.DataPolicyApi.HandleDataPolicyListByPage)
	}
}
