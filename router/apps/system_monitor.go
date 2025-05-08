package apps

import (
	"github.com/Thingsly/backend/internal/api"
	"github.com/Thingsly/backend/pkg/metrics"

	"github.com/gin-gonic/gin"
)

// SystemMonitor system monitoring module
type SystemMonitor struct{}

// InitSystemMonitor initialize the system monitoring related routes
func (m *SystemMonitor) InitSystemMonitor(r *gin.RouterGroup, metricsManager *metrics.Metrics) {
	// Register routes
	r.GET("system/metrics/current", api.Controllers.SystemMonitorApi.GetCurrentSystemMetrics)
	r.GET("system/metrics/history", api.Controllers.SystemMonitorApi.GetHistorySystemMetrics)
}