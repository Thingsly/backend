package apps

import (
	"github.com/Thingsly/backend/internal/api"

	"github.com/gin-gonic/gin"
)

type Casbin struct{}

func (*Casbin) Init(Router *gin.RouterGroup) {
	url := Router.Group("casbin")
	{

		url.POST("function", api.Controllers.CasbinApi.AddFunctionToRole)
		url.DELETE("function/:id", api.Controllers.CasbinApi.DeleteFunctionFromRole)
		url.PUT("function", api.Controllers.CasbinApi.UpdateFunctionFromRole)
		url.GET("function", api.Controllers.CasbinApi.HandleFunctionFromRole)

		url.POST("user", api.Controllers.CasbinApi.AddRoleToUser)
		url.DELETE("user/:id", api.Controllers.CasbinApi.DeleteRolesFromUser)
		url.PUT("user", api.Controllers.CasbinApi.UpdateRolesFromUser)
		url.GET("user", api.Controllers.CasbinApi.HandleRolesFromUser)
	}
}
