// router/apps/open_api_keys.go
package apps

import (
	"github.com/Thingsly/backend/internal/api"

	"github.com/gin-gonic/gin"
)

type OpenAPIKey struct{}

func (*OpenAPIKey) InitOpenAPIKey(Router *gin.RouterGroup) {
	openAPIRouter := Router.Group("open/keys")
	{

		openAPIRouter.POST("", api.Controllers.OpenAPIKeyApi.CreateOpenAPIKey)
		openAPIRouter.GET("", api.Controllers.OpenAPIKeyApi.GetOpenAPIKeyList)
		openAPIRouter.PUT("", api.Controllers.OpenAPIKeyApi.UpdateOpenAPIKey)
		openAPIRouter.DELETE(":id", api.Controllers.OpenAPIKeyApi.DeleteOpenAPIKey)
	}
}
