package api

import (
	"github.com/Thingsly/backend/pkg/global"
	"github.com/Thingsly/backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

type SystemApi struct{}

// HandleSystime
// @Summary Get system time
// @Description Get system time
// @Tags system
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/systime [get]
func (*SystemApi) HandleSystime(c *gin.Context) {
	c.Set("data", map[string]interface{}{"systime": utils.GetSecondTimestamp()})
}

// Health check
// @Summary Health check
// @Description Health check
// @Tags system
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health [get]
func (*SystemApi) HealthCheck(c *gin.Context) {
	c.Set("data", nil)
}

// HandleSysVersion
// @Summary Get system version
// @Description Get system version
// @Tags system
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/sys_version [get]
func (*SystemApi) HandleSysVersion(c *gin.Context) {
	// c.Set("data", map[string]interface{}{"version": global.SYSTEM_VERSION})
	c.Set("data", map[string]interface{}{"version": global.BE_VERSION})
}
