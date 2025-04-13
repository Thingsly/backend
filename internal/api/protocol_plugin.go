package api

import (
	model "github.com/HustIoTPlatform/backend/internal/model"
	service "github.com/HustIoTPlatform/backend/internal/service"

	"github.com/gin-gonic/gin"
)

type ProtocolPluginApi struct{}

// CreateProtocolPlugin 
// @Router   /api/v1/protocol_plugin [post]
func (*ProtocolPluginApi) CreateProtocolPlugin(c *gin.Context) {
	var req model.CreateProtocolPluginReq
	if !BindAndValidate(c, &req) {
		return
	}
	data, err := service.GroupApp.ProtocolPlugin.CreateProtocolPlugin(&req)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", data)
}

// DeleteProtocolPlugin 
// @Router   /api/v1/protocol_plugin/{id} [delete]
func (*ProtocolPluginApi) DeleteProtocolPlugin(c *gin.Context) {
	id := c.Param("id")
	err := service.GroupApp.ProtocolPlugin.DeleteProtocolPlugin(id)

	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", id)
}

// UpdateProtocolPlugin
// @Router   /api/v1/protocol_plugin [put]
func (*ProtocolPluginApi) UpdateProtocolPlugin(c *gin.Context) {
	var req model.UpdateProtocolPluginReq
	if !BindAndValidate(c, &req) {
		return
	}
	err := service.GroupApp.ProtocolPlugin.UpdateProtocolPlugin(&req)

	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", nil)
}

// UpdateProtocolPlugin
// @Router   /api/v1/protocol_plugin [get]
func (*ProtocolPluginApi) HandleProtocolPluginListByPage(c *gin.Context) {
	var req model.GetProtocolPluginListByPageReq
	if !BindAndValidate(c, &req) {
		return
	}

	list, err := service.GroupApp.ProtocolPlugin.GetProtocolPluginListByPage(&req)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", list)
}

// GetProtocolPluginForm 
// @Router   /api/v1/protocol_plugin/device_config_form [get]
func (*ProtocolPluginApi) HandleProtocolPluginForm(c *gin.Context) {
	var req model.GetProtocolPluginFormReq
	if !BindAndValidate(c, &req) {
		return
	}

	data, err := service.GroupApp.ProtocolPlugin.GetProtocolPluginForm(&req)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", data)
}

// GetProtocolPluginForm
// @Router   /api/v1/protocol_plugin/config_form [get]
func (*ProtocolPluginApi) HandleProtocolPluginFormByProtocolType(c *gin.Context) {
	var req model.GetProtocolPluginFormByProtocolType
	if !BindAndValidate(c, &req) {
		return
	}

	data, err := service.GroupApp.ServicePlugin.GetProtocolPluginFormByProtocolType(req.ProtocolType, req.DeviceType)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", data)
}

// /api/v1/plugin/device/config
// Protocol plugin retrieves device configuration
func (*ProtocolPluginApi) HandleDeviceConfigForProtocolPlugin(c *gin.Context) {
	var req model.GetDeviceConfigReq
	if !BindAndValidate(c, &req) {
		return
	}

	data, err := service.GroupApp.ProtocolPlugin.GetDeviceConfig(req)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", data)
}
