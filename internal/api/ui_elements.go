package api

import (
	model "github.com/Thingsly/backend/internal/model"
	service "github.com/Thingsly/backend/internal/service"
	"github.com/Thingsly/backend/pkg/errcode"
	utils "github.com/Thingsly/backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

type UiElementsApi struct{}

// CreateUiElements
// @Summary Create UI elements
// @Description Create UI elements
// @Tags ui_elements
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Success 200 {object} map[string]interface{}
// @Router   /api/v1/ui_elements [post]
func (*UiElementsApi) CreateUiElements(c *gin.Context) {
	var req model.CreateUiElementsReq
	if !BindAndValidate(c, &req) {
		return
	}
	err := service.GroupApp.UiElements.CreateUiElements(&req)
	if err != nil {
		c.Error(err)
		return
	}

	c.Set("data", nil)
}

// UpdateUiElements
// @Summary Update UI elements
// @Description Update UI elements
// @Tags ui_elements
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Success 200 {object} map[string]interface{}
// @Router   /api/v1/ui_elements [put]
func (*UiElementsApi) UpdateUiElements(c *gin.Context) {
	var req model.UpdateUiElementsReq
	if !BindAndValidate(c, &req) {
		return
	}

	if req.ElementType == nil && req.Authority == nil {
		c.Error(errcode.WithData(errcode.CodeParamError, map[string]interface{}{
			"element_type": "element_type or authority is required",
		}))
		return
	}

	err := service.GroupApp.UiElements.UpdateUiElements(&req)
	if err != nil {
		c.Error(err)
		return
	}

	c.Set("data", nil)
}

// DeleteUiElements
// @Summary Delete UI elements
// @Description Delete UI elements
// @Tags ui_elements
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Success 200 {object} map[string]interface{}
// @Router   /api/v1/ui_elements/{id} [delete]
func (*UiElementsApi) DeleteUiElements(c *gin.Context) {
	id := c.Param("id")
	err := service.GroupApp.UiElements.DeleteUiElements(id)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", nil)
}

// ServeUiElementsListByPage
// @Summary Get UI elements list by page
// @Description Get UI elements list by page
// @Tags ui_elements
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Success 200 {object} map[string]interface{}
// @Router   /api/v1/ui_elements [get]
func (*UiElementsApi) ServeUiElementsListByPage(c *gin.Context) {
	var req model.ServeUiElementsListByPageReq
	if !BindAndValidate(c, &req) {
		return
	}

	UiElementsList, err := service.GroupApp.UiElements.ServeUiElementsListByPage(&req)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", UiElementsList)
}

// ServeUiElementsListByPage
// @Summary Get UI elements list by authority
// @Description Get UI elements list by authority
// @Tags ui_elements
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Success 200 {object} map[string]interface{}
// @Router   /api/v1/ui_elements/menu [get]
func (*UiElementsApi) ServeUiElementsListByAuthority(c *gin.Context) {
	var userClaims = c.MustGet("claims").(*utils.UserClaims)

	uiElementsList, err := service.GroupApp.UiElements.ServeUiElementsListByAuthority(userClaims)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", uiElementsList)
}

// ServeUiElementsListByTenant
// @Summary Get UI elements list by tenant
// @Description Get UI elements list by tenant
// @Tags ui_elements
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Success 200 {object} map[string]interface{}
// @Router   /api/v1/ui_elements/select/form [get]
func (*UiElementsApi) ServeUiElementsListByTenant(c *gin.Context) {
	uiElementsList, err := service.GroupApp.UiElements.GetTenantUiElementsList()
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", uiElementsList)
}
