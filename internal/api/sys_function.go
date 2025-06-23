package api

import (
	dal "github.com/Thingsly/backend/internal/dal"
	"github.com/Thingsly/backend/internal/service"
	"github.com/Thingsly/backend/pkg/errcode"
	"github.com/Thingsly/backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

type SysFunctionApi struct{}

// Get sys function
// @Summary Get sys function
// @Description Get sys function
// @Tags sys_function
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Success 200 {object} model.GetSysFunctionRes
// @Router   /api/v1/sys_function [get]
func (*SysFunctionApi) HandleSysFcuntion(c *gin.Context) {
	lang := c.GetHeader("Accept-Language")
	date, err := service.GroupApp.SysFunction.GetSysFuncion(lang)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", date)
}

// Update sys function
// @Summary Update sys function
// @Description Update sys function
// @Tags sys_function
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param function_id path string true "Function id"
// @Param function body model.UpdateSysFunctionReq true "Function"
// @Success 200 {object} model.UpdateSysFunctionRes
// @Router   /api/v1/sys_function/{function_id} [put]
func (*SysFunctionApi) UpdateSysFcuntion(c *gin.Context) {
	var userClaims = c.MustGet("claims").(*utils.UserClaims)
	if userClaims.Authority != dal.SYS_ADMIN {
		c.Error(errcode.WithData(errcode.CodeParamError, map[string]interface{}{
			"authority": "authority is not sys admin",
		}))
		return
	}
	id := c.Param("id")
	err := service.GroupApp.SysFunction.UpdateSysFuncion(id)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", nil)
}
