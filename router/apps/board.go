package apps

import (
	"github.com/Thingsly/backend/internal/api"

	"github.com/gin-gonic/gin"
)

type Board struct {
}

func (*Board) InitBoard(Router *gin.RouterGroup) {
	url := Router.Group("board")
	{

		url.POST("", api.Controllers.BoardApi.CreateBoard)

		url.DELETE(":id", api.Controllers.BoardApi.DeleteBoard)

		url.PUT("", api.Controllers.BoardApi.UpdateBoard)

		url.GET("", api.Controllers.BoardApi.HandleBoardListByPage)

		url.GET(":id", api.Controllers.BoardApi.HandleBoard)

		url.GET("home", api.Controllers.BoardApi.HandleBoardListByTenantId)

		url.GET("trend", api.Controllers.BoardApi.GetDeviceTrend)

	}

	devices(url)

	tenant(url)

	user(url)
}

func devices(Router *gin.RouterGroup) {
	url := Router.Group("device")

	url.GET("total", api.Controllers.BoardApi.HandleDeviceTotal)

	url.GET("", api.Controllers.BoardApi.HandleDevice)
}

func tenant(Router *gin.RouterGroup) {
	url := Router.Group("tenant")

	url.GET("", api.Controllers.BoardApi.HandleTenant)

	url.GET("user/info", api.Controllers.BoardApi.HandleTenantUserInfo)

	url.GET("device/info", api.Controllers.BoardApi.HandleTenantDeviceInfo)
}

func user(Router *gin.RouterGroup) {
	url := Router.Group("user")

	url.GET("info", api.Controllers.BoardApi.HandleUserInfo)

	url.POST("update", api.Controllers.BoardApi.UpdateUserInfo)

	url.POST("update/password", api.Controllers.BoardApi.UpdateUserInfoPassword)
}
