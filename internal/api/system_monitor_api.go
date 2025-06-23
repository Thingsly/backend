package api

import (
	"fmt"
	"time"

	"github.com/Thingsly/backend/internal/service"
	"github.com/Thingsly/backend/pkg/errcode"
	"github.com/Thingsly/backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

// SystemMonitorApi
type SystemMonitorApi struct{}

// GetCurrentSystemMetrics
// @Summary Get current system metrics
// @Description Get current system metrics
// @Tags system_monitor
// @Accept  json
// @Produce  json
// @Param x-token header string true "Authentication token"
// @Success 200 {object} response.Response
// @Router /api/v1/system/metrics/current [get]
func (api *SystemMonitorApi) GetCurrentSystemMetrics(c *gin.Context) {
	userClaims := c.MustGet("claims").(*utils.UserClaims)
	if userClaims.Authority != "SYS_ADMIN" {
		c.Error(errcode.New(errcode.CodeNoPermission))
		return
	}

	metrics, err := service.GroupApp.SystemMonitor.GetCurrentMetrics()
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", metrics)
}

// GetHistorySystemMetrics
// @Summary Get history system metrics
// @Description Get history system metrics
// @Tags system_monitor
// @Accept  json
// @Produce  json
// @Param x-token header string true "Authentication token"
// @Param hours query int false "query hours" default(24)
// @Success 200 {object} response.Response
// @Router /api/v1/system/metrics/history [get]
func (api *SystemMonitorApi) GetHistorySystemMetrics(c *gin.Context) {
	userClaims := c.MustGet("claims").(*utils.UserClaims)
	if userClaims.Authority != "SYS_ADMIN" {
		c.Error(errcode.New(errcode.CodeNoPermission))
		return
	}

	hours := 24
	if hoursStr := c.Query("hours"); hoursStr != "" {
		if _, err := fmt.Sscanf(hoursStr, "%d", &hours); err != nil {
			hours = 24
		}
	}

	if hours <= 0 {
		hours = 1
	} else if hours > 72 {
		hours = 72
	}

	duration := time.Duration(hours) * time.Hour
	data, err := service.GroupApp.SystemMonitor.GetCombinedHistoryData(duration)
	if err != nil {
		c.Error(err)
		return
	}

	c.Set("data", data)
}
