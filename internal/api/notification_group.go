package api

import (
	model "github.com/HustIoTPlatform/backend/internal/model"
	service "github.com/HustIoTPlatform/backend/internal/service"
	utils "github.com/HustIoTPlatform/backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

type NotificationGroupApi struct{}

// CreateNotificationGroup
// @Router   /api/v1/notification_group [post]
func (*NotificationGroupApi) CreateNotificationGroup(c *gin.Context) {
	var req model.CreateNotificationGroupReq

	if !BindAndValidate(c, &req) {
		return
	}
	userClaims := c.MustGet("claims").(*utils.UserClaims)

	notificationGroup, err := service.GroupApp.NotificationGroup.CreateNotificationGroup(&req, userClaims)
	if err != nil {
		c.Error(err)
		return
	}

	notificationGroupOs, err := utils.SerializeData(*notificationGroup, ReadNotificationGroupOutSchema{})
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", notificationGroupOs)
}

// GetNotificationGroup
// @Router   /api/v1/notification_group/{id} [get]
func (*NotificationGroupApi) HandleNotificationGroupById(c *gin.Context) {
	id := c.Param("id")
	if ntfgroup, err := service.GroupApp.NotificationGroup.GetNotificationGroupById(id); err != nil {
		c.Error(err)
		return
	} else {
		notificationGroupOs, err := utils.SerializeData(*ntfgroup, ReadNotificationGroupOutSchema{})
		if err != nil {
			c.Error(err)
			return
		}
		c.Set("data", notificationGroupOs)
	}
}

// UpdateNotificationGroup
// @Router   /api/v1/notification_group/{id} [put]
func (*NotificationGroupApi) UpdateNotificationGroup(c *gin.Context) {
	id := c.Param("id")
	var req model.UpdateNotificationGroupReq
	if !BindAndValidate(c, &req) {
		return
	}

	if updated, err := service.GroupApp.NotificationGroup.UpdateNotificationGroup(id, &req); err != nil {
		c.Error(err)
		return
	} else {
		updateoutput, err := utils.SerializeData(updated, UpdateNotificationGroupOutSchema{})
		if err != nil {
			c.Error(err)
			return
		}
		c.Set("data", updateoutput)
	}
}

// DeleteNotificationGroup
// @Router   /api/v1/notification_group/{id} [delete]
func (*NotificationGroupApi) DeleteNotificationGroup(c *gin.Context) {
	id := c.Param("id")
	if err := service.GroupApp.NotificationGroup.DeleteNotificationGroup(id); err != nil {
		c.Error(err)
		return
	} else {
		c.Set("data", nil)
	}
}

// GetNotificationGroupListByPage
// @Router   /api/v1/notification_group/list [get]
func (*NotificationGroupApi) HandleNotificationGroupListByPage(c *gin.Context) {
	var req model.GetNotificationGroupListByPageReq
	if !BindAndValidate(c, &req) {
		return
	}

	userClaims := c.MustGet("claims").(*utils.UserClaims)
	notificationList, err := service.GroupApp.NotificationGroup.GetNotificationGroupListByPage(&req, userClaims)
	if err != nil {
		c.Error(err)
		return
	}
	ntfoutput, err := utils.SerializeData(notificationList, GetNotificationGroupListByPageOutSchema{})
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", ntfoutput)
}
