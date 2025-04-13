package api

import (
	"net/http"

	model "github.com/HustIoTPlatform/backend/internal/model"
	service "github.com/HustIoTPlatform/backend/internal/service"
	"github.com/HustIoTPlatform/backend/pkg/errcode"
	utils "github.com/HustIoTPlatform/backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

type RoleApi struct{}

// CreateRole 
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
