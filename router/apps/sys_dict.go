package apps

import (
	"github.com/Thingsly/backend/internal/api"

	"github.com/gin-gonic/gin"
)

type Dict struct {
}

func (*Dict) InitDict(Router *gin.RouterGroup) {
	dictapi := Router.Group("dict")
	{

		dictapi.POST("column", api.Controllers.DictApi.CreateDictColumn)

		dictapi.POST("language", api.Controllers.CreateDictLanguage)

		dictapi.GET("enum", api.Controllers.DictApi.HandleDict)

		dictapi.GET("", api.Controllers.DictApi.HandleDictLisyByPage)

		dictapi.GET("language/:id", api.Controllers.HandleDictLanguage)

		dictapi.DELETE("column/:id", api.Controllers.DictApi.DeleteDictColumn)

		dictapi.DELETE("language/:id", api.Controllers.DictApi.DeleteDictLanguage)

		dictapi.GET("protocol/service", api.Controllers.DictApi.HandleProtocolAndService)
	}
}
