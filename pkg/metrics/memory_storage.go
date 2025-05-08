package metrics

import (
	"sync"
	"time"
)

const (
	// Default retention period of 24 hours
	DefaultRetentionPeriod = 24 * time.Hour
	// Default collection interval of 5 minutes
	DefaultCollectionInterval = 5 * time.Minute
)

// MemoryStorage implements a memory-based storage for historical metric data
type MemoryStorage struct {
	sync.RWMutex
	cpuData    []MetricDataPoint // CPU usage history data
	memoryData []MetricDataPoint // Memory usage history data
	diskData   []MetricDataPoint // Disk usage history data

	// Current latest data
	currentData SystemMetrics

	// Data retention period
	retentionPeriod time.Duration
	// Last cleanup time
	lastCleanup time.Time
}

// NewMemoryStorage creates a new memory storage instance
func NewMemoryStorage() *MemoryStorage {
	storage := &MemoryStorage{
		cpuData:         make([]MetricDataPoint, 0, 288), // 24 hours, 5 minutes per point = 24*12 = 288 points
		memoryData:      make([]MetricDataPoint, 0, 288),
		diskData:        make([]MetricDataPoint, 0, 288),
		retentionPeriod: DefaultRetentionPeriod,
		lastCleanup:     time.Now(),
	}

	// Start periodic cleanup of expired data
	go storage.periodicCleanup()

	return storage
}

// SaveMetrics save the current system metrics
func (s *MemoryStorage) SaveMetrics(timestamp time.Time, cpuUsage, memoryUsage, diskUsage float64) error {
	s.Lock()
	defer s.Unlock()

	// Update current data
	s.currentData = SystemMetrics{
		CPUUsage:    cpuUsage,
		MemoryUsage: memoryUsage,
		DiskUsage:   diskUsage,
		Timestamp:   timestamp,
	}

	// Save a data point every 5 minutes
	if len(s.cpuData) == 0 || time.Since(s.cpuData[len(s.cpuData)-1].Timestamp) >= DefaultCollectionInterval {
		point := MetricDataPoint{
			Timestamp: timestamp,
			Value:     cpuUsage,
		}
		s.cpuData = append(s.cpuData, point)

		point.Value = memoryUsage
		s.memoryData = append(s.memoryData, point)

		point.Value = diskUsage
		s.diskData = append(s.diskData, point)

		// Check if cleanup is needed
		if time.Since(s.lastCleanup) > time.Hour {
			s.cleanup()
			s.lastCleanup = time.Now()
		}
	}

	return nil
}

// GetHistoryData get historical data based on the metric type and time range
func (s *MemoryStorage) GetHistoryData(metric string, duration time.Duration) ([]MetricDataPoint, error) {
	s.RLock()
	defer s.RUnlock()

	var data []MetricDataPoint
	switch metric {
	case "cpu":
		data = s.cpuData
	case "memory":
		data = s.memoryData
	case "disk":
		data = s.diskData
	default:
		return nil, nil
	}

	// If no time range is specified or the time range is greater than the retention period, return all data
	if duration <= 0 || duration >= s.retentionPeriod {
		result := make([]MetricDataPoint, len(data))
		copy(result, data)
		return result, nil
	}

	// Filter data based on the time range
	cutoffTime := time.Now().Add(-duration)
	result := make([]MetricDataPoint, 0, len(data))

	for _, point := range data {
		if point.Timestamp.After(cutoffTime) {
			result = append(result, point)
		}
	}

	return result, nil
}

// GetCurrentData get the latest system metrics
func (s *MemoryStorage) GetCurrentData() (*SystemMetrics, error) {
	s.RLock()
	defer s.RUnlock()

	result := s.currentData
	return &result, nil
}

// Clean up expired data
func (s *MemoryStorage) cleanup() {
	cutoffTime := time.Now().Add(-s.retentionPeriod)

	s.cpuData = filterDataPoints(s.cpuData, cutoffTime)
	s.memoryData = filterDataPoints(s.memoryData, cutoffTime)
	s.diskData = filterDataPoints(s.diskData, cutoffTime)
}

// Periodically clean up expired data
func (s *MemoryStorage) periodicCleanup() {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		s.Lock()
		s.cleanup()
		s.lastCleanup = time.Now()
		s.Unlock()
	}
}

// Filter data points, keep data after the specified time
func filterDataPoints(data []MetricDataPoint, cutoffTime time.Time) []MetricDataPoint {
	if len(data) == 0 {
		return data
	}

	index := 0
	for i, point := range data {
		if point.Timestamp.After(cutoffTime) {
			index = i
			break
		}
	}

	if index == 0 && data[0].Timestamp.Before(cutoffTime) {
		return data[:0] // All expired
	}

	// Move data to the front
	if index > 0 {
		copy(data, data[index:])
		return data[:len(data)-index]
	}

	return data
}