package apps

import (
	"github.com/Thingsly/backend/internal/api"

	"github.com/gin-gonic/gin"
)

type EventData struct{}

func (*EventData) InitEventData(Router *gin.RouterGroup) {
	evnetDataApi := Router.Group("event/datas")
	{
		evnetDataApi.GET("", api.Controllers.EventDataApi.HandleEventDatasListByPage)
	}
}
