package api

import (
	"context"

	"github.com/Thingsly/backend/internal/model"
	service "github.com/Thingsly/backend/internal/service"
	"github.com/Thingsly/backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

type ExpectedDataApi struct{}

// @Summary Query expected data list
// @Description Query expected data list
// @Tags expected
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param page query int true "Page number"
// @Param page_size query int true "Page size"
// @Success 200 {object} model.GetExpectedDataPageRes
// @Router /api/v1/expected/data/list [get]
func (*ExpectedDataApi) HandleExpectedDataList(c *gin.Context) {
	var req model.GetExpectedDataPageReq
	if !BindAndValidate(c, &req) {
		return
	}
	var userClaims = c.MustGet("claims").(*utils.UserClaims)
	resp, err := service.GroupApp.ExpectedData.PageList(context.Background(), &req, userClaims)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", resp)
}

// Add expected data
// @Summary Add expected data
// @Description Add expected data
// @Tags expected
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param expected_data body model.CreateExpectedDataReq true "Expected data"
// @Success 200 {object} model.CreateExpectedDataRes
// @Router /api/v1/expected/data [post]
func (*ExpectedDataApi) CreateExpectedData(c *gin.Context) {
	var req model.CreateExpectedDataReq
	if !BindAndValidate(c, &req) {
		return
	}
	var userClaims = c.MustGet("claims").(*utils.UserClaims)
	resp, err := service.GroupApp.ExpectedData.Create(c, &req, userClaims)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", resp)
}

// Delete expected data
// @Summary Delete expected data
// @Description Delete expected data
// @Tags expected
// @Accept json
// @Produce json
// @Param x-token header string true "Authentication token"
// @Param id path string true "Expected data ID"
// @Success 200 {object} model.DeleteExpectedDataRes
// @Router /api/v1/expected/data/{id} [delete]
func (*ExpectedDataApi) DeleteExpectedData(c *gin.Context) {
	id := c.Param("id")
	err := service.GroupApp.ExpectedData.Delete(c, id)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", map[string]interface{}{})
}
