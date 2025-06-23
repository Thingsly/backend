package api

import (
	"net/http"

	model "github.com/Thingsly/backend/internal/model"
	service "github.com/Thingsly/backend/internal/service"
	"github.com/Thingsly/backend/pkg/errcode"
	utils "github.com/Thingsly/backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

type RoleApi struct{}

// CreateRole
// @Summary Create role
// @Description Create role
// @Tags role
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param role body model.CreateRoleReq true "Role"
// @Success 200 {object} model.CreateRoleRes
// @Router   /api/v1/role [post]
func (*RoleApi) CreateRole(c *gin.Context) {
	var req model.CreateRoleReq
	if !BindAndValidate(c, &req) {
		return
	}

	userClaims := c.MustGet("claims").(*utils.UserClaims)

	err := service.GroupApp.Role.CreateRole(&req, userClaims)
	if err != nil {
		c.Error(err)
		return
	}

	c.Set("data", nil)
}

// UpdateRole
// @Summary Update role
// @Description Update role
// @Tags role
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param role body model.UpdateRoleReq true "Role"
// @Success 200 {object} model.UpdateRoleRes
// @Router   /api/v1/role [put]
func (*RoleApi) UpdateRole(c *gin.Context) {
	var req model.UpdateRoleReq
	if !BindAndValidate(c, &req) {
		return
	}

	if req.Description == nil && req.Name == "" {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "Content to be modified cannot be empty"})
		return
	}

	data, err := service.GroupApp.Role.UpdateRole(&req)
	if err != nil {
		c.Error(err)
		return
	}

	c.Set("data", data)
}

// DeleteRole
// @Summary Delete role
// @Description Delete role
// @Tags role
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param id path string true "Role id"
// @Success 200 {object} model.DeleteRoleRes
// @Router   /api/v1/role/{id} [delete]
func (*RoleApi) DeleteRole(c *gin.Context) {
	id := c.Param("id")

	// Check if the role is in use
	if service.GroupApp.Casbin.HasRole(id) {
		c.Error(errcode.WithData(errcode.CodeParamError, map[string]interface{}{
			"role_id": id,
			"error":   "Role in use",
		}))
		return
	}

	err := service.GroupApp.Role.DeleteRole(id)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", nil)
}

// GetRoleListByPage
// @Summary Get role list by page
// @Description Get role list by page
// @Tags role
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param page query int true "Page"
// @Param page_size query int true "Page size"
// @Success 200 {object} model.GetRoleListByPageRes
// @Router   /api/v1/role [get]
func (*RoleApi) HandleRoleListByPage(c *gin.Context) {
	var req model.GetRoleListByPageReq
	if !BindAndValidate(c, &req) {
		return
	}

	var userClaims = c.MustGet("claims").(*utils.UserClaims)
	roleList, err := service.GroupApp.Role.GetRoleListByPage(&req, userClaims)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", roleList)
}
