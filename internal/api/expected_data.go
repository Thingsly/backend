package api

import (
	"context"

	"github.com/HustIoTPlatform/backend/internal/model"
	service "github.com/HustIoTPlatform/backend/internal/service"
	"github.com/HustIoTPlatform/backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

type ExpectedDataApi struct{}

// Query expected data list
// /api/v1/expected/data/list
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
// /api/v1/expected/data POST
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
// /api/v1/expected/data/{id} DELETE
func (*ExpectedDataApi) DeleteExpectedData(c *gin.Context) {
	id := c.Param("id")
	err := service.GroupApp.ExpectedData.Delete(c, id)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", map[string]interface{}{})
}
