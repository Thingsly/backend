package api

import (
	"github.com/Thingsly/backend/internal/model"
	"github.com/Thingsly/backend/internal/service"

	"github.com/gin-gonic/gin"
)

type ServicePluginApi struct{}

// Create service plugin
// @Summary Create service plugin
// @Description Create service plugin
// @Tags service_plugin
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param service_plugin body model.CreateServicePluginReq true "Service plugin"
// @Success 200 {object} model.CreateServicePluginRes
// @Router   /api/v1/service [post]
func (*ServicePluginApi) Create(c *gin.Context) {
	var req model.CreateServicePluginReq
	if !BindAndValidate(c, &req) {
		return
	}
	//var userClaims = c.MustGet("claims").(*utils.UserClaims)
	resp, err := service.GroupApp.ServicePlugin.Create(&req)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", resp)
}

// Get service plugin by page
// @Summary Get service plugin by page
// @Description Get service plugin by page
// @Tags service_plugin
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param page query int true "Page"
// @Param page_size query int true "Page size"
// @Success 200 {object} model.GetServicePluginByPageRes
// @Router   /api/v1/service/list [get]
func (*ServicePluginApi) HandleList(c *gin.Context) {
	var req model.GetServicePluginByPageReq
	if !BindAndValidate(c, &req) {
		return
	}
	resp, err := service.GroupApp.ServicePlugin.List(&req)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", resp)
}

// Get service plugin by id
// @Summary Get service plugin by id
// @Description Get service plugin by id
// @Tags service_plugin
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param id path string true "Service plugin id"
// @Success 200 {object} model.GetServicePluginRes
// @Router   /api/v1/service/detail/{id} [get]
func (*ServicePluginApi) Handle(c *gin.Context) {
	id := c.Param("id")
	resp, err := service.GroupApp.ServicePlugin.Get(id)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", resp)
}

// Update service plugin
// @Summary Update service plugin
// @Description Update service plugin
// @Tags service_plugin
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param service_plugin body model.UpdateServicePluginReq true "Service plugin"
// @Success 200 {object} model.UpdateServicePluginRes
// @Router   /api/v1/service [put]
func (*ServicePluginApi) Update(c *gin.Context) {
	var req model.UpdateServicePluginReq
	if !BindAndValidate(c, &req) {
		return
	}
	err := service.GroupApp.ServicePlugin.Update(&req)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", map[string]interface{}{})
}

// Delete service plugin
// @Summary Delete service plugin
// @Description Delete service plugin
// @Tags service_plugin
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param id path string true "Service plugin id"
// @Success 200 {object} model.DeleteServicePluginRes
// @Router   /api/v1/service/{id} [delete]
func (*ServicePluginApi) Delete(c *gin.Context) {
	id := c.Param("id")
	err := service.GroupApp.ServicePlugin.Delete(id)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", map[string]interface{}{})
}

// Heartbeat
// @Summary Heartbeat
// @Description Heartbeat
// @Tags service_plugin
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param heartbeat body model.HeartbeatReq true "Heartbeat"
// @Success 200 {object} model.HeartbeatRes
// @Router   /api/v1/plugin/heartbeat [post]
func (*ServicePluginApi) Heartbeat(c *gin.Context) {
	var req model.HeartbeatReq
	if !BindAndValidate(c, &req) {
		return
	}
	err := service.GroupApp.ServicePlugin.Heartbeat(&req)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", map[string]interface{}{})
}

// Get service select
// @Summary Get service select
// @Description Get service select
// @Tags service_plugin
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Success 200 {object} model.GetServiceSelectRes
// @Router   /api/v1/service/plugin/select [get]
func (*ServicePluginApi) HandleServiceSelect(c *gin.Context) {
	var req model.GetServiceSelectReq
	if !BindAndValidate(c, &req) {
		return
	}
	resp, err := service.GroupApp.ServicePlugin.GetServiceSelect(&req)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", resp)
}

// Get service plugin by service identifier
// @Summary Get service plugin by service identifier
// @Description Get service plugin by service identifier
// @Tags service_plugin
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param service_identifier path string true "Service identifier"
// @Success 200 {object} model.GetServicePluginByServiceIdentifierRes
// @Router   /api/v1/service/plugin/info/{service_identifier} [get]
func (*ServicePluginApi) HandleServicePluginByServiceIdentifier(c *gin.Context) {
	var req model.GetServicePluginByServiceIdentifierReq
	if !BindAndValidate(c, &req) {
		return
	}
	data, err := service.GroupApp.ServicePlugin.GetServicePluginByServiceIdentifier(req.ServiceIdentifier)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", data)
}
