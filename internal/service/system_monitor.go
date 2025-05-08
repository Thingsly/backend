package service

import (
	"time"

	"github.com/Thingsly/backend/pkg/metrics"
)

// SystemMonitor
type SystemMonitor struct{}

// Global metrics manager
var metricsManager *metrics.Metrics

// SetMetricsManager Set metrics manager
func SetMetricsManager(m *metrics.Metrics) {
	metricsManager = m
}

// GetCurrentMetrics Get current system metrics
func (s *SystemMonitor) GetCurrentMetrics() (*metrics.SystemMetrics, error) {
	if metricsManager == nil {
		return nil, nil
	}
	return metricsManager.GetCurrentMetrics()
}

// GetHistoryData Get historical data
func (s *SystemMonitor) GetHistoryData(metricType string, duration time.Duration) ([]metrics.MetricDataPoint, error) {
	if metricsManager == nil {
		return nil, nil
	}
	return metricsManager.GetHistoryData(metricType, duration)
}

// GetCombinedHistoryData Get combined historical data
func (s *SystemMonitor) GetCombinedHistoryData(duration time.Duration) ([]metrics.MetricsTimePoint, error) {
	if metricsManager == nil {
		return nil, nil
	}
	return metricsManager.GetCombinedHistoryData(duration)
}
