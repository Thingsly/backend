package api

import (
	model "github.com/Thingsly/backend/internal/model"
	service "github.com/Thingsly/backend/internal/service"
	utils "github.com/Thingsly/backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

type EventDataApi struct{}

// GetEventDatasListByPage
// @Router   /api/v1/event/datas [get]
func (*EventDataApi) HandleEventDatasListByPage(c *gin.Context) {
	var req model.GetEventDatasListByPageReq
	if !BindAndValidate(c, &req) {
		return
	}
	var userClaims = c.MustGet("claims").(*utils.UserClaims)
	data, err := service.GroupApp.EventData.GetEventDatasListByPage(&req, userClaims)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", data)
}
