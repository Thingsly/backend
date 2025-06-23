package api

import (
	model "github.com/Thingsly/backend/internal/model"
	service "github.com/Thingsly/backend/internal/service"
	utils "github.com/Thingsly/backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

type EventDataApi struct{}

// @Summary Get event data list by page
// @Description Get event data list by page
// @Tags event
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param page query int true "Page number"
// @Param page_size query int true "Page size"
// @Success 200 {object} model.GetEventDatasListByPageRes
// @Router /api/v1/event/datas [get]
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
