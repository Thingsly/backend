package api

import (
	model "github.com/Thingsly/backend/internal/model"
	service "github.com/Thingsly/backend/internal/service"
	"github.com/Thingsly/backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

type OperationLogsApi struct{}

// GetListByPage
// @Summary Get operation log list by page
// @Description Get operation log list by page
// @Tags operation_log
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param page query int true "Page"
// @Param page_size query int true "Page size"
// @Success 200 {object} model.GetOperationLogListByPageRes
// @Router   /api/v1/operation_logs [get]
func (*OperationLogsApi) HandleListByPage(c *gin.Context) {
	var req model.GetOperationLogListByPageReq
	if !BindAndValidate(c, &req) {
		return
	}
	var userClaims = c.MustGet("claims").(*utils.UserClaims)
	list, err := service.GroupApp.OperationLogs.GetListByPage(&req, userClaims)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", list)
}
