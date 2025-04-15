package apps

import (
	"github.com/Thingsly/backend/internal/api"

	"github.com/gin-gonic/gin"
)

type SysFunction struct{}

func (*SysFunction) Init(Router *gin.RouterGroup) {
	url := Router.Group("sys_function")
	{

		url.PUT(":id", api.Controllers.SysFunctionApi.UpdateSysFcuntion)

	}
}
