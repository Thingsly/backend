package api

import (
	model "github.com/Thingsly/backend/internal/model"
	service "github.com/Thingsly/backend/internal/service"

	"github.com/gin-gonic/gin"
)

type ProtocolPluginApi struct{}

// CreateProtocolPlugin
// @Summary Create protocol plugin
// @Description Create protocol plugin
// @Tags protocol_plugin
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param protocol_plugin body model.CreateProtocolPluginReq true "Protocol plugin"
// @Success 200 {object} model.CreateProtocolPluginRes
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
// @Summary Delete protocol plugin
// @Description Delete protocol plugin
// @Tags protocol_plugin
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param id path string true "Protocol plugin id"
// @Success 200 {object} model.DeleteProtocolPluginRes
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
// @Summary Update protocol plugin
// @Description Update protocol plugin
// @Tags protocol_plugin
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param protocol_plugin body model.UpdateProtocolPluginReq true "Protocol plugin"
// @Success 200 {object} model.UpdateProtocolPluginRes
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
// @Summary Get protocol plugin list by page
// @Description Get protocol plugin list by page
// @Tags protocol_plugin
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param page query int true "Page"
// @Param page_size query int true "Page size"
// @Success 200 {object} model.GetProtocolPluginListByPageRes
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
// @Summary Get protocol plugin form
// @Description Get protocol plugin form
// @Tags protocol_plugin
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param protocol_plugin body model.GetProtocolPluginFormReq true "Protocol plugin"
// @Success 200 {object} model.GetProtocolPluginFormRes
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
// @Summary Get protocol plugin form by protocol type
// @Description Get protocol plugin form by protocol type
// @Tags protocol_plugin
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param protocol_plugin body model.GetProtocolPluginFormByProtocolTypeReq true "Protocol plugin"
// @Success 200 {object} model.GetProtocolPluginFormByProtocolTypeRes
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

// HandleDeviceConfigForProtocolPlugin
// @Summary Protocol plugin retrieves device configuration
// @Description Protocol plugin retrieves device configuration
// @Tags protocol_plugin
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param protocol_plugin body model.GetDeviceConfigReq true "Protocol plugin"
// @Success 200 {object} model.GetDeviceConfigRes
// @Router   /api/v1/plugin/device/config [get]
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
