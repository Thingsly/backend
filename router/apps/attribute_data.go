package apps

import (
	"github.com/HustIoTPlatform/backend/internal/api"

	"github.com/gin-gonic/gin"
)

type AttributeData struct{}

func (*AttributeData) InitAttributeData(Router *gin.RouterGroup) {
	attributedataapi := Router.Group("attribute/datas")
	{

		attributedataapi.GET(":id", api.Controllers.AttributeDataApi.HandleDataList)

		attributedataapi.GET("set/logs", api.Controllers.AttributeDataApi.HandleAttributeSetLogsDataListByPage)

		attributedataapi.DELETE(":id", api.Controllers.AttributeDataApi.DeleteData)

		attributedataapi.POST("pub", api.Controllers.AttributeDataApi.AttributePutMessage)

		attributedataapi.GET("get", api.Controllers.AttributeDataApi.AttributeGetMessage)

		attributedataapi.GET("key", api.Controllers.AttributeDataApi.HandleAttributeDataByKey)
	}
}
