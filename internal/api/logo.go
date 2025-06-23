package api

import (
	model "github.com/Thingsly/backend/internal/model"
	service "github.com/Thingsly/backend/internal/service"

	"github.com/gin-gonic/gin"
)

type LogoApi struct{}

// UpdateLogo
// @Summary Update logo
// @Description Update logo
// @Tags logo
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param logo body model.UpdateLogoReq true "Logo"
// @Success 200 {object} model.UpdateLogoRes
// @Router /api/v1/logo [put]
func (LogoApi) UpdateLogo(c *gin.Context) {
	var req model.UpdateLogoReq
	if !BindAndValidate(c, &req) {
		return
	}

	err := service.GroupApp.Logo.UpdateLogo(&req)
	if err != nil {
		c.Error(err)
		return
	}

	c.Set("data", nil)
}

// GetLogoListByPage
// @Summary Get logo list
// @Description Get logo list
// @Tags logo
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Success 200 {object} model.GetLogoListRes
// @Router /api/v1/logo [get]
func (LogoApi) HandleLogoList(c *gin.Context) {
	logoList, err := service.GroupApp.Logo.GetLogoList()
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", logoList)
}
