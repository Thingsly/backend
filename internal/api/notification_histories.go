package api

import (
	model "github.com/Thingsly/backend/internal/model"
	service "github.com/Thingsly/backend/internal/service"
	utils "github.com/Thingsly/backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

type NotificationHistoryApi struct{}

// GetNotificationHistoryListByPage
// @Summary Get notification history list by page
// @Description Get notification history list by page
// @Tags notification_history
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param page query int true "Page"
// @Param page_size query int true "Page size"
// @Success 200 {object} model.GetNotificationHistoryListByPageOutSchema
// @Router   /api/v1/notification_history/list [get]
func (*NotificationHistoryApi) HandleNotificationHistoryListByPage(c *gin.Context) {
	var req model.GetNotificationHistoryListByPageReq
	if !BindAndValidate(c, &req) {
		return
	}

	var userClaims = c.MustGet("claims").(*utils.UserClaims)
	req.TenantID = userClaims.TenantID
	notificationList, err := service.GroupApp.NotificationHisory.GetNotificationHistoryListByPage(&req)
	if err != nil {
		c.Error(err)
		return
	}
	ntfoutput, err := utils.SerializeData(notificationList, GetNotificationHistoryListByPageOutSchema{})
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", ntfoutput)
}
