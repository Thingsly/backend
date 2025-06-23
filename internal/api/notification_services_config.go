package api

import (
	dal "github.com/Thingsly/backend/internal/dal"
	model "github.com/Thingsly/backend/internal/model"
	service "github.com/Thingsly/backend/internal/service"
	"github.com/Thingsly/backend/pkg/errcode"
	utils "github.com/Thingsly/backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

type NotificationServicesConfigApi struct{}

// SaveNotificationServicesConfig Create/Update notification service configuration (2-in-1 interface)
// @Summary Create/Update notification service configuration
// @Description Create/Update notification service configuration
// @Tags notification_services_config
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param notification_services_config body model.SaveNotificationServicesConfigReq true "Notification services config"
// @Success 200 {object} model.SaveNotificationServicesConfigRes
// @Router   /api/v1/notification/services/config [post]
func (*NotificationServicesConfigApi) SaveNotificationServicesConfig(c *gin.Context) {
	var req model.SaveNotificationServicesConfigReq
	if !BindAndValidate(c, &req) {
		return
	}
	userClaims := c.MustGet("claims").(*utils.UserClaims)

	// Validate SYS_ADMIN
	if userClaims.Authority != dal.SYS_ADMIN {
		c.Error(errcode.WithData(errcode.CodeSystemError, map[string]interface{}{
			"authority": "authority is not sys admin",
		}))
		return
	}

	// Validate notification type, currently supports email and SMS
	if req.NoticeType != model.NoticeType_Email && req.NoticeType != model.NoticeType_SME_CODE {
		c.Error(errcode.WithData(errcode.CodeSystemError, map[string]interface{}{
			"noticeType": "noticeType is not email or sme",
		}))
		return
	}

	// Toggle enumeration validation
	if req.Status != model.OPEN && req.Status != model.CLOSE {
		c.Error(errcode.WithData(errcode.CodeSystemError, map[string]interface{}{
			"status": "status is not open or close",
		}))
		return
	}

	data, err := service.GroupApp.NotificationServicesConfig.SaveNotificationServicesConfig(&req)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", data)
}

// GetNotificationServicesConfig Get notification service configuration
// @Summary Get notification service configuration
// @Description Get notification service configuration
// @Tags notification_services_config
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param type path string true "Notification type"
// @Success 200 {object} model.GetNotificationServicesConfigRes
// @Router   /api/v1/notification/services/config/{type} [get]
func (*NotificationServicesConfigApi) HandleNotificationServicesConfig(c *gin.Context) {
	noticeType := c.Param("type")
	userClaims := c.MustGet("claims").(*utils.UserClaims)
	// Validate SYS_ADMIN
	if userClaims.Authority != dal.SYS_ADMIN {
		c.Error(errcode.WithData(errcode.CodeSystemError, map[string]interface{}{
			"authority": "authority is not sys admin",
		}))
		return
	}
	data, err := service.GroupApp.NotificationServicesConfig.GetNotificationServicesConfig(noticeType)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", data)
}

// SendTestEmail Send test email
// @Summary Send test email
// @Description Send test email
// @Tags notification_services_config
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param send_test_email body model.SendTestEmailReq true "Send test email"
// @Router   /api/v1/notification/services/config/e-mail/test [post]
func (*NotificationServicesConfigApi) SendTestEmail(c *gin.Context) {
	var req model.SendTestEmailReq
	if !BindAndValidate(c, &req) {
		return
	}
	err := service.GroupApp.NotificationServicesConfig.SendTestEmail(&req)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", nil)
}
