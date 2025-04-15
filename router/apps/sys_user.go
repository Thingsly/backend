package apps

import (
	"github.com/Thingsly/backend/internal/api"

	"github.com/gin-gonic/gin"
)

type User struct {
}

func (*User) InitUser(Router *gin.RouterGroup) {
	userapi := Router.Group("user")
	{

		userapi.GET("detail", api.Controllers.UserApi.HandleUserDetail)
		userapi.PUT("update", api.Controllers.UserApi.UpdateUsers)
		userapi.GET("logout", api.Controllers.UserApi.Logout)
		userapi.GET("refresh", api.Controllers.UserApi.RefreshToken)

		userapi.GET("", api.Controllers.UserApi.HandleUserListByPage)
		userapi.POST("", api.Controllers.UserApi.CreateUser)
		userapi.PUT("", api.Controllers.UserApi.UpdateUser)
		userapi.DELETE(":id", api.Controllers.UserApi.DeleteUser)
		userapi.GET(":id", api.Controllers.UserApi.HandleUser)
		userapi.POST("transform", api.Controllers.UserApi.TransformUser)

		userapi.GET("/tenant/id", api.Controllers.UserApi.GetTenantID)

	}
}
