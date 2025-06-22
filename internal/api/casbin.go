package api

import (
	model "github.com/Thingsly/backend/internal/model"
	service "github.com/Thingsly/backend/internal/service"
	"github.com/Thingsly/backend/pkg/errcode"

	"github.com/gin-gonic/gin"
)

type CasbinApi struct{}

var casbinService = service.GroupApp.Casbin

// AddFunctionToRole
// @Router   /api/v1/casbin/function [post]
// @Summary Add functions to role
// @Description Add functions to a specific role
// @Tags Casbin
// @Accept json
// @Produce json
// @Param functions_role_validate body model.FunctionsRoleValidate true "Functions and role details"
func (*CasbinApi) AddFunctionToRole(c *gin.Context) {
	var req model.FunctionsRoleValidate
	if !BindAndValidate(c, &req) {
		return
	}

	ok := casbinService.AddFunctionToRole(req.RoleID, req.FunctionsIDs)
	if !ok {
		c.Error(errcode.WithData(errcode.CodeParamError, map[string]interface{}{
			"role_id":      req.RoleID,
			"function_ids": req.FunctionsIDs,
			"error":        "AddFunctionToRole failed",
		}))
		return
	}

	c.Set("data", nil)
}

// GetFunctionFromRole
// @Router   /api/v1/casbin/function [get]
// @Summary Get functions from role
// @Description Get functions from a specific role
// @Tags Casbin
// @Accept json
// @Produce json
// @Param role_id path string true "Role ID"
// @Success 200 {object} model.RoleValidate "Functions retrieved successfully"
func (*CasbinApi) HandleFunctionFromRole(c *gin.Context) {
	var req model.RoleValidate
	if !BindAndValidate(c, &req) {
		return
	}

	roles, ok := casbinService.GetFunctionFromRole(req.RoleID)
	if !ok {
		c.Error(errcode.WithData(errcode.CodeParamError, map[string]interface{}{
			"role_id": req.RoleID,
			"error":   "GetFunctionFromRole failed",
		}))
		return
	}

	c.Set("data", roles)
}

// UpdateFunctionFromRole
// @Router   /api/v1/casbin/function [put]
// @Summary Update functions from role
// @Description Update functions for a specific role
// @Tags Casbin
// @Accept json
// @Produce json
// @Param functions_role_validate body model.FunctionsRoleValidate true "Functions and role details"
func (*CasbinApi) UpdateFunctionFromRole(c *gin.Context) {
	var req model.FunctionsRoleValidate
	if !BindAndValidate(c, &req) {
		return
	}

	if req.RoleID == "" && req.FunctionsIDs == nil {
		c.Error(errcode.WithData(errcode.CodeParamError, map[string]interface{}{
			"role_id":      req.RoleID,
			"function_ids": req.FunctionsIDs,
			"error":        "UpdateFunctionFromRole failed",
		}))
		return
	}

	f, _ := casbinService.GetFunctionFromRole(req.RoleID)
	if len(f) > 0 {
		// No record found, deletion will return false
		ok := casbinService.RemoveRoleAndFunction(req.RoleID)
		if !ok {
			c.Error(errcode.WithData(errcode.CodeParamError, map[string]interface{}{
				"role_id": req.RoleID,
				"error":   "RemoveRoleAndFunction failed",
			}))
			return
		}
	}
	ok := casbinService.AddFunctionToRole(req.RoleID, req.FunctionsIDs)
	if !ok {
		c.Error(errcode.WithData(errcode.CodeParamError, map[string]interface{}{
			"role_id":      req.RoleID,
			"function_ids": req.FunctionsIDs,
			"error":        "AddFunctionToRole failed",
		}))
	}
	c.Set("data", nil)
}

// DeleteFunctionFromRole
// @Router   /api/v1/casbin/function/{id} [delete]
// @Summary Delete functions from role
// @Description Delete functions from a specific role
// @Tags Casbin
// @Accept json
// @Produce json
// @Param id path string true "Role ID"
func (*CasbinApi) DeleteFunctionFromRole(c *gin.Context) {
	id := c.Param("id")
	ok := casbinService.RemoveRoleAndFunction(id)
	if !ok {
		c.Error(errcode.WithData(errcode.CodeParamError, map[string]interface{}{
			"role_id": id,
			"error":   "RemoveRoleAndFunction failed",
		}))
		return
	}
	c.Set("data", nil)
}

// AddRoleToUser
// @Router   /api/v1/casbin/user [post]
// @Summary Add roles to user
// @Description Add roles to a specific user
// @Tags Casbin
// @Accept json
// @Produce json
// @Param roles_user_validate body model.RolesUserValidate true "Roles and user details"
func (*CasbinApi) AddRoleToUser(c *gin.Context) {
	var req model.RolesUserValidate
	if !BindAndValidate(c, &req) {
		return
	}

	ok := casbinService.AddRolesToUser(req.UserID, req.RolesIDs)
	if !ok {
		c.Error(errcode.WithData(errcode.CodeParamError, map[string]interface{}{
			"user_id": req.UserID,
			"role_id": req.RolesIDs,
			"error":   "AddRolesToUser failed",
		}))
		return
	}

	c.Set("data", nil)

}

// GetRolesFromUser
// @Router   /api/v1/casbin/user [get]
// @Summary Get roles from user
// @Description Get roles from a specific user
// @Tags Casbin
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} model.UserValidate "Roles retrieved successfully"
func (*CasbinApi) HandleRolesFromUser(c *gin.Context) {
	var req model.UserValidate
	if !BindAndValidate(c, &req) {
		return
	}

	roles, ok := casbinService.GetRoleFromUser(req.UserID)
	if !ok {
		c.Error(errcode.WithData(errcode.CodeParamError, map[string]interface{}{
			"user_id": req.UserID,
			"error":   "GetRoleFromUser failed",
		}))
		return
	}

	c.Set("data", roles)

}

// UpdateRolesFromUser
// @Router   /api/v1/casbin/user [put]
// @Summary Update roles from user
// @Description Update roles for a specific user
// @Tags Casbin
// @Accept json
// @Produce json
// @Param roles_user_validate body model.RolesUserValidate true "Roles and user details"
func (*CasbinApi) UpdateRolesFromUser(c *gin.Context) {
	var req model.RolesUserValidate
	if !BindAndValidate(c, &req) {
		return
	}

	if req.UserID == "" && req.RolesIDs == nil {
		c.Error(errcode.WithData(errcode.CodeParamError, map[string]interface{}{
			"user_id": req.UserID,
			"role_id": req.RolesIDs,
			"error":   "UpdateRolesFromUser failed",
		}))
		return
	}

	casbinService.RemoveUserAndRole(req.UserID)
	ok := casbinService.AddRolesToUser(req.UserID, req.RolesIDs)
	if !ok {
		c.Error(errcode.WithData(errcode.CodeParamError, map[string]interface{}{
			"user_id": req.UserID,
			"role_id": req.RolesIDs,
			"error":   "AddRolesToUser failed",
		}))
	}
	c.Set("data", nil)
}

// DeleteRolesFromUser
// @Router   /api/v1/casbin/user/{id} [delete]
// @Summary Delete roles from user
// @Description Delete roles from a specific user
// @Tags Casbin
// @Accept json
// @Produce json
// @Param id path string true "User ID"
func (*CasbinApi) DeleteRolesFromUser(c *gin.Context) {
	id := c.Param("id")
	ok := casbinService.RemoveUserAndRole(id)
	if !ok {
		c.Error(errcode.WithData(errcode.CodeParamError, map[string]interface{}{
			"user_id": id,
			"error":   "RemoveUserAndRole failed",
		}))
		return
	}
	c.Set("data", nil)
}
