package api

import (
	"github.com/HustIoTPlatform/backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

type SystemApi struct{}

// /api/v1/systime
func (*SystemApi) HandleSystime(c *gin.Context) {
	c.Set("data", map[string]interface{}{"systime": utils.GetSecondTimestamp()})
}

// Health check /health
func (*SystemApi) HealthCheck(c *gin.Context) {
	c.Set("data", nil)
}
