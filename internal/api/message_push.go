package api

import (
	"github.com/Thingsly/backend/internal/model"
	"github.com/Thingsly/backend/internal/service"
	"github.com/Thingsly/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type MessagePushApi struct {
}

// @Summary Create message push
// @Description Create message push
// @Tags message_push
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param message_push body model.CreateMessagePushReq true "Message push"
// @Success 200 {object} model.CreateMessagePushRes
// @Router /api/v1/message_push [post]
func (*MessagePushApi) CreateMessagePush(c *gin.Context) {
	var req model.CreateMessagePushReq
	if !BindAndValidate(c, &req) {
		return
	}
	var userClaims = c.MustGet("claims").(*utils.UserClaims)
	err := service.GroupApp.MessagePush.CreateMessagePush(&req, userClaims.ID)
	if err != nil {
		c.Error(err)
		return
	}

	c.Set("data", nil)
}

// @Summary Message push manage logout
// @Description Message push manage logout
// @Tags message_push
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param message_push body model.MessagePushMangeLogoutReq true "Message push"
// @Success 200 {object} model.MessagePushMangeLogoutRes
// @Router /api/v1/message_push/manage/logout [post]
func (*MessagePushApi) MessagePushMangeLogout(c *gin.Context) {
	var req model.MessagePushMangeLogoutReq
	if !BindAndValidate(c, &req) {
		return
	}
	var userClaims = c.MustGet("claims").(*utils.UserClaims)
	err := service.GroupApp.MessagePush.MessagePushMangeLogout(&req, userClaims.ID)
	if err != nil {
		c.Error(err)
		return
	}

	c.Set("data", nil)
}

// @Summary Get message push config
// @Description Get message push config
// @Tags message_push
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Success 200 {object} model.MessagePushConfigRes
// @Router /api/v1/message_push/config [get]
func (*MessagePushApi) GetMessagePushConfig(c *gin.Context) {
	res, err := service.GroupApp.MessagePush.GetMessagePushConfig()
	if err != nil {
		c.Error(err)
		return
	}

	c.Set("data", res)
}

// @Summary Set message push config
// @Description Set message push config
// @Tags message_push
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param message_push body model.MessagePushConfigReq true "Message push"
// @Success 200 {object} model.MessagePushConfigRes
// @Router /api/v1/message_push/config [post]
func (*MessagePushApi) SetMessagePushConfig(c *gin.Context) {
	var req model.MessagePushConfigReq
	if !BindAndValidate(c, &req) {
		return
	}
	err := service.GroupApp.MessagePush.SetMessagePushConfig(&req)
	if err != nil {
		c.Error(err)
		return
	}

	c.Set("data", nil)
}
