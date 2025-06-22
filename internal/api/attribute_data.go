package api

import (
	"strconv"

	"github.com/Thingsly/backend/internal/model"
	"github.com/Thingsly/backend/internal/service"
	"github.com/Thingsly/backend/pkg/constant"
	"github.com/Thingsly/backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

type AttributeDataApi struct{}

// GetDataList
// @Router   /api/v1/attribute/datas/{id} [get]
func (*AttributeDataApi) HandleDataList(c *gin.Context) {
	id := c.Param("id")
	data, err := service.GroupApp.AttributeData.GetAttributeDataList(id)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", data)
}

// Query device attributes by key
// /api/v1/attribute/datas/key [get]
func (*AttributeDataApi) HandleAttributeDataByKey(c *gin.Context) {
	var req model.GetDataListByKeyReq
	if !BindAndValidate(c, &req) {
		return
	}
	data, err := service.GroupApp.AttributeData.GetAttributeDataByKey(req)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", data)
}

// DeleteData - Delete Attribute Data
// @Router   /api/v1/attribute/datas/{id} [delete]
func (*AttributeDataApi) DeleteData(c *gin.Context) {
	id := c.Param("id")
	err := service.GroupApp.AttributeData.DeleteAttributeData(id)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", nil)
}

// GetAttributeSetLogsDataListByPage - Attribute Distribution Record Query (Pagination)
// @Router   /api/v1/attribute/datas/set/logs [get]
func (*AttributeDataApi) HandleAttributeSetLogsDataListByPage(c *gin.Context) {
	var req model.GetAttributeSetLogsListByPageReq
	if !BindAndValidate(c, &req) {
		return
	}
	data, err := service.GroupApp.AttributeData.GetAttributeSetLogsDataListByPage(req)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", data)
}

// /api/v1/attribute/datas/pub [post]
// @Summary Send attribute data
// @Description Send attribute data to the server
// @Tags Attribute Data
// @Accept json
// @Produce json
// @Param attribute_data body model.AttributePutMessage true "Attribute data to send"
// @Success 200 {object} model.AttributePutMessage "Attribute data sent successfully"
func (*AttributeDataApi) AttributePutMessage(c *gin.Context) {
	var req model.AttributePutMessage
	if !BindAndValidate(c, &req) {
		return
	}

	userClaims := c.MustGet("claims").(*utils.UserClaims)
	err := service.GroupApp.AttributeData.AttributePutMessage(c, userClaims.ID, &req, strconv.Itoa(constant.Manual))
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", nil)
}

// Send request to get attributes
// /api/v1/attribute/datas/get
// @Summary Get attribute data
// @Description Get attribute data from the server
// @Tags Attribute Data
// @Accept json
// @Produce json
// @Param attribute_data body model.AttributeGetMessageReq true "Attribute data to get"
// @Success 200 {object} model.AttributeGetMessageReq "Attribute data retrieved successfully"
func (*AttributeDataApi) AttributeGetMessage(c *gin.Context) {
	var req model.AttributeGetMessageReq
	if !BindAndValidate(c, &req) {
		return
	}
	userClaims := c.MustGet("claims").(*utils.UserClaims)
	err := service.GroupApp.AttributeData.AttributeGetMessage(userClaims, &req)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", nil)
}
