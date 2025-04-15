package api

import (
	model "github.com/Thingsly/backend/internal/model"
	service "github.com/Thingsly/backend/internal/service"
	utils "github.com/Thingsly/backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

type DeviceConfigApi struct{}

// CreateDeviceConfig
// @Router   /api/v1/device_config [post]
func (*DeviceConfigApi) CreateDeviceConfig(c *gin.Context) {
	var req model.CreateDeviceConfigReq
	if !BindAndValidate(c, &req) {
		return
	}
	var userClaims = c.MustGet("claims").(*utils.UserClaims)
	data, err := service.GroupApp.DeviceConfig.CreateDeviceConfig(&req, userClaims)
	if err != nil {
		c.Error(err)
		return
	}

	c.Set("data", data)
}

// UpdateDeviceConfig
// @Router   /api/v1/device_config [put]
func (*DeviceConfigApi) UpdateDeviceConfig(c *gin.Context) {
	var req model.UpdateDeviceConfigReq
	if !BindAndValidate(c, &req) {
		return
	}

	data, err := service.GroupApp.DeviceConfig.UpdateDeviceConfig(req)
	if err != nil {
		c.Error(err)
		return
	}

	c.Set("data", data)
}

// DeleteDeviceConfig
// @Router   /api/v1/device_config/{id} [delete]
func (*DeviceConfigApi) DeleteDeviceConfig(c *gin.Context) {
	id := c.Param("id")
	err := service.GroupApp.DeviceConfig.DeleteDeviceConfig(id)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", nil)
}

// GetDeviceConfigById
// @Router   /api/v1/device_config/{id} [get]
func (*DeviceConfigApi) HandleDeviceConfigById(c *gin.Context) {
	id := c.Param("id")
	info, err := service.GroupApp.DeviceConfig.GetDeviceConfigByID(c, id)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", info)
}

// GetDeviceConfigListByPage
// @Router   /api/v1/device_config [get]
func (*DeviceConfigApi) HandleDeviceConfigListByPage(c *gin.Context) {
	var req model.GetDeviceConfigListByPageReq
	if !BindAndValidate(c, &req) {
		return
	}

	var userClaims = c.MustGet("claims").(*utils.UserClaims)

	deviceconfigList, err := service.GroupApp.DeviceConfig.GetDeviceConfigListByPage(&req, userClaims)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", deviceconfigList)
}

// @Router   /api/v1/device_config/menu [get]
func (*DeviceConfigApi) HandleDeviceConfigListMenu(c *gin.Context) {
	var req model.GetDeviceConfigListMenuReq
	if !BindAndValidate(c, &req) {
		return
	}

	var userClaims = c.MustGet("claims").(*utils.UserClaims)

	deviceconfigList, err := service.GroupApp.DeviceConfig.GetDeviceConfigListMenu(&req, userClaims)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", deviceconfigList)
}

// BatchUpdateDeviceConfig
// @Router   /api/v1/device_config/batch [put]
func (*DeviceConfigApi) BatchUpdateDeviceConfig(c *gin.Context) {
	var req model.BatchUpdateDeviceConfigReq
	if !BindAndValidate(c, &req) {
		return
	}

	err := service.GroupApp.DeviceConfig.BatchUpdateDeviceConfig(&req)
	if err != nil {
		c.Error(err)
		return
	}

	c.Set("data", nil)
}

// /api/v1/device_config/connect
func (*DeviceConfigApi) HandleDeviceConfigConnect(c *gin.Context) {
	var param model.DeviceIDReq
	if !BindAndValidate(c, &param) {
		return
	}
	data, err := service.GroupApp.DeviceConfig.GetDeviceConfigConnect(c, param.DeviceID)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", data)
}

// /api/v1/device_config/voucher_type
func (*DeviceConfigApi) HandleVoucherType(c *gin.Context) {
	var param model.GetVoucherTypeReq
	if !BindAndValidate(c, &param) {
		return
	}
	data, err := service.GroupApp.DeviceConfig.GetVoucherTypeForm(param.DeviceType, param.ProtocolType)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", data)
}

// Get the dropdown list for automation actions based on device configuration ID
// /api/v1/device_config/metrics/menu [get]
func (*DeviceConfigApi) HandleActionByDeviceConfigID(c *gin.Context) {
	var param model.GetActionByDeviceConfigIDReq
	if !BindAndValidate(c, &param) {
		return
	}
	list, err := service.GroupApp.DeviceConfig.GetActionByDeviceConfigID(param.DeviceConfigID)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", list)
}

// Retrieve the dropdown list for automation actions based on the device configuration ID
// /api/v1/device_config/metrics/condition/menu
func (*DeviceConfigApi) HandleConditionByDeviceConfigID(c *gin.Context) {
	var param model.GetActionByDeviceConfigIDReq
	if !BindAndValidate(c, &param) {
		return
	}
	list, err := service.GroupApp.DeviceConfig.GetConditionByDeviceConfigID(param.DeviceConfigID)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", list)
}
