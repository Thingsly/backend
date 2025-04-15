package apps

import (
	"github.com/Thingsly/backend/internal/api"

	"github.com/gin-gonic/gin"
)

type Role struct {
}

func (*Role) Init(Router *gin.RouterGroup) {
	url := Router.Group("role")
	{

		url.POST("", api.Controllers.RoleApi.CreateRole)

		url.DELETE(":id", api.Controllers.RoleApi.DeleteRole)

		url.PUT("", api.Controllers.RoleApi.UpdateRole)

		url.GET("", api.Controllers.RoleApi.HandleRoleListByPage)
	}
}
