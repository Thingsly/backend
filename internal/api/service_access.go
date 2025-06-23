package api

import (
	"github.com/Thingsly/backend/internal/model"
	"github.com/Thingsly/backend/internal/service"
	"github.com/Thingsly/backend/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ServiceAccessApi struct{}

// Create service access
// @Summary Create service access
// @Description Create service access
// @Tags service_access
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param service_access body model.CreateAccessReq true "Service access"
// @Success 200 {object} model.CreateAccessRes
// @Router   /api/v1/service/access [post]
func (*ServiceAccessApi) Create(c *gin.Context) {
	var req model.CreateAccessReq
	if !BindAndValidate(c, &req) {
		return
	}
	var userClaims = c.MustGet("claims").(*utils.UserClaims)
	resp, err := service.GroupApp.ServiceAccess.CreateAccess(&req, userClaims)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", resp)
}

// Get service access by page
// @Summary Get service access by page
// @Description Get service access by page
// @Tags service_access
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param page query int true "Page"
// @Param page_size query int true "Page size"
// @Success 200 {object} model.GetServiceAccessByPageRes
// @Router   /api/v1/service/access/list [get]
func (*ServiceAccessApi) HandleList(c *gin.Context) {
	var req model.GetServiceAccessByPageReq
	if !BindAndValidate(c, &req) {
		return
	}
	var userClaims = c.MustGet("claims").(*utils.UserClaims)
	resp, err := service.GroupApp.ServiceAccess.List(&req, userClaims)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", resp)
}

// Update service access
// @Summary Update service access
// @Description Update service access
// @Tags service_access
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param service_access body model.UpdateAccessReq true "Service access"
// @Success 200 {object} model.UpdateAccessRes
// @Router   /api/v1/service/access [put]
func (*ServiceAccessApi) Update(c *gin.Context) {
	var req model.UpdateAccessReq
	if !BindAndValidate(c, &req) {
		return
	}
	err := service.GroupApp.ServiceAccess.Update(&req)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", nil)
}

// Delete service access
// @Summary Delete service access
// @Description Delete service access
// @Tags service_access
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param id path string true "Service access id"
// @Success 200 {object} model.DeleteAccessRes
// @Router   /api/v1/service/access/:id [delete]
func (*ServiceAccessApi) Delete(c *gin.Context) {
	id := c.Param("id")
	err := service.GroupApp.ServiceAccess.Delete(id)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", nil)
}

// Get service access voucher form
// @Summary Get service access voucher form
// @Description Get service access voucher form
// @Tags service_access
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Success 200 {object} model.GetServiceAccessVoucherFormRes
// @Router   /api/v1/service/access/voucher/form [get]
func (*ServiceAccessApi) HandleVoucherForm(c *gin.Context) {
	var req model.GetServiceAccessVoucherFormReq
	if !BindAndValidate(c, &req) {
		return
	}
	resp, err := service.GroupApp.ServiceAccess.GetVoucherForm(&req)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", resp)
}

// Get service access device list
// @Summary Get service access device list
// @Description Get service access device list
// @Tags service_access
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Success 200 {object} model.GetServiceAccessDeviceListRes
// @Router   /api/v1/service/access/device/list [get]
func (*ServiceAccessApi) HandleDeviceList(c *gin.Context) {
	var req model.ServiceAccessDeviceListReq
	if !BindAndValidate(c, &req) {
		return
	}
	var userClaims = c.MustGet("claims").(*utils.UserClaims)
	resp, err := service.GroupApp.ServiceAccess.GetServiceAccessDeviceList(&req, userClaims)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", resp)
}

// Get plugin service access list
// @Summary Get plugin service access list
// @Description Get plugin service access list
// @Tags service_access
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Success 200 {object} model.GetPluginServiceAccessListRes
// @Router   /api/v1/plugin/service/access/list [get]
func (*ServiceAccessApi) HandlePluginServiceAccessList(c *gin.Context) {
	logrus.Info("get plugin list")
	var req model.GetPluginServiceAccessListReq
	if !BindAndValidate(c, &req) {
		return
	}
	resp, err := service.GroupApp.ServiceAccess.GetPluginServiceAccessList(&req)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", resp)
}

// Get plugin service access
// @Summary Get plugin service access
// @Description Get plugin service access
// @Tags service_access
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Success 200 {object} model.GetPluginServiceAccessRes
// @Router   /api/v1/plugin/service/access [get]
func (*ServiceAccessApi) HandlePluginServiceAccess(c *gin.Context) {
	var req model.GetPluginServiceAccessReq
	if !BindAndValidate(c, &req) {
		return
	}
	resp, err := service.GroupApp.ServiceAccess.GetPluginServiceAccess(&req)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", resp)
}
