// pkg/metrics/metrics.go
package metrics

import (
	"os"
	"runtime"
	"sort"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/process"
	"github.com/sirupsen/logrus"
)

type Metrics struct {
	MemoryUsage     prometheus.Gauge
	MemoryAllocated prometheus.Gauge
	MemoryObjects   prometheus.Gauge
	CPUUsage        prometheus.Gauge
	DiskUsage       prometheus.Gauge
	GoroutinesTotal prometheus.Gauge
	GCPauseTotal    prometheus.Gauge
	GCRuns          prometheus.Counter

	APIRequestTotal *prometheus.CounterVec
	APIErrorTotal   *prometheus.CounterVec
	APILatency      *prometheus.HistogramVec

	BusinessErrorTotal *prometheus.CounterVec
	CriticalErrorTotal prometheus.Counter

	SlowRequestTotal   *prometheus.CounterVec
	LargeResponseTotal *prometheus.CounterVec
	// History data storage
	historyStorage HistoryStorage
}

// HistoryStorage defines the history data storage interface
type HistoryStorage interface {
	SaveMetrics(timestamp time.Time, cpuUsage, memoryUsage, diskUsage float64) error
	GetHistoryData(metric string, duration time.Duration) ([]MetricDataPoint, error)
	GetCurrentData() (*SystemMetrics, error)
}

// MetricDataPoint represents a metric data point
type MetricDataPoint struct {
	Timestamp time.Time `json:"timestamp"`
	Value     float64   `json:"value"`
}

// SystemMetrics current value of system metrics
type SystemMetrics struct {
	CPUUsage    float64   `json:"cpu_usage"`    // CPU usage percentage
	MemoryUsage float64   `json:"memory_usage"` // Memory usage percentage
	DiskUsage   float64   `json:"disk_usage"`   // Disk usage percentage
	Timestamp   time.Time `json:"timestamp"`    // Timestamp
}

// MetricsTimePoint represents all metric data at a specific time
type MetricsTimePoint struct {
	Timestamp   time.Time `json:"timestamp"` // Timestamp
	CPUUsage    float64   `json:"cpu"`       // CPU usage percentage
	MemoryUsage float64   `json:"memory"`    // Memory usage percentage
	DiskUsage   float64   `json:"disk"`      // Disk usage percentage
}

func NewMetrics(namespace string) *Metrics {
	return &Metrics{

		MemoryUsage: promauto.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "memory_usage_bytes",
				Help:      "Current memory usage in bytes",
			},
		),

		MemoryAllocated: promauto.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "memory_allocated_bytes",
				Help:      "Total memory allocated in bytes",
			},
		),

		MemoryObjects: promauto.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "memory_objects",
				Help:      "Number of allocated objects",
			},
		),

		CPUUsage: promauto.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "cpu_usage_percent",
				Help:      "CPU usage percentage",
			},
		),

		DiskUsage: promauto.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "disk_usage_percent",
				Help:      "Disk usage percentage",
			},
		),

		GoroutinesTotal: promauto.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "goroutines_total",
				Help:      "Total number of goroutines",
			},
		),

		GCPauseTotal: promauto.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "gc_pause_total_seconds",
				Help:      "Total GC pause time in seconds",
			},
		),

		GCRuns: promauto.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "gc_runs_total",
				Help:      "Total number of completed GC cycles",
			},
		),

		APIRequestTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "api_requests_total",
				Help:      "Total number of API requests by path and method",
			},
			[]string{"path", "method"},
		),

		APIErrorTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "api_errors_total",
				Help:      "Total number of API errors by error type",
			},
			[]string{"type"},
		),

		APILatency: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "api_latency_seconds",
				Help:      "API latency in seconds",
				Buckets:   []float64{0.1, 0.3, 0.5, 1.0, 2.0, 5.0},
			},
			[]string{"path"},
		),

		BusinessErrorTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "business_errors_total",
				Help:      "Total number of business errors by module",
			},
			[]string{"module", "code"},
		),

		CriticalErrorTotal: promauto.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "critical_errors_total",
				Help:      "Total number of critical errors that need immediate attention",
			},
		),

		SlowRequestTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "slow_requests_total",
				Help:      "Total number of slow requests (>1s) by path",
			},
			[]string{"path"},
		),

		LargeResponseTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "large_responses_total",
				Help:      "Total number of large responses (>1MB) by path",
			},
			[]string{"path"},
		),
	}
}

// SetHistoryStorage set the history data storage implementation
func (m *Metrics) SetHistoryStorage(storage HistoryStorage) {
	m.historyStorage = storage
}

func (m *Metrics) StartMetricsCollection(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		var lastPauseNs uint64
		var lastNumGC uint32

		pid := os.Getpid()
		process, err := process.NewProcess(int32(pid))
		if err != nil {
			logrus.Warnf("Failed to get process: %v", err)
			return
		}

		for range ticker.C {

			var memStats runtime.MemStats
			runtime.ReadMemStats(&memStats)
			m.MemoryUsage.Set(float64(memStats.Alloc))
			m.MemoryAllocated.Set(float64(memStats.Sys))
			m.MemoryObjects.Set(float64(memStats.HeapObjects))

			// Process CPU usage - use a more reliable method
			cpuPercent := 0.0

			// Method 1: Use process.Percent to calculate CPU usage
			percent, err := process.Percent(time.Second)
			if err == nil && percent > 0 {
				cpuPercent = percent
			} else {
				// Method 2: Use the entire system CPU usage as a backup
				cpuStat, err := cpu.Percent(time.Second, false)
				if err == nil && len(cpuStat) > 0 {
					cpuPercent = cpuStat[0]
				}
			}

			m.CPUUsage.Set(cpuPercent)

			// Disk usage
			var diskUsagePercent float64
			// Use C drive on Windows, and root directory on other systems
			diskPath := "/"
			if runtime.GOOS == "windows" {
				diskPath = "C:\\"
			}

			diskStat, err := disk.Usage(diskPath)
			if err == nil {
				diskUsagePercent = diskStat.UsedPercent
				m.DiskUsage.Set(diskUsagePercent)
			}

			m.GoroutinesTotal.Set(float64(runtime.NumGoroutine()))

			if memStats.NumGC > lastNumGC {
				diff := memStats.NumGC - lastNumGC
				m.GCRuns.Add(float64(diff))
				lastNumGC = memStats.NumGC
			}

			if pauseNs := memStats.PauseTotalNs; pauseNs > lastPauseNs {
				pauseDiff := float64(pauseNs - lastPauseNs)
				m.GCPauseTotal.Add(pauseDiff / 1e9)
				lastPauseNs = pauseNs
			}

			// Store historical data
			if m.historyStorage != nil {
				// Get the current memory usage percentage
				var memoryUsagePercent float64
				if memStats.Sys > 0 {
					memoryUsagePercent = float64(memStats.Alloc) / float64(memStats.Sys) * 100
				}

				// Save the metric history
				err := m.historyStorage.SaveMetrics(
					time.Now(),
					cpuPercent,
					memoryUsagePercent,
					diskUsagePercent,
				)
				if err != nil {
					logrus.Warnf("Failed to save metrics history: %v", err)
				}
			}
		}
	}()
}

func (m *Metrics) RecordAPIRequest(path, method string) {
	m.APIRequestTotal.WithLabelValues(path, method).Inc()
}

func (m *Metrics) RecordAPIError(errorType string) {
	m.APIErrorTotal.WithLabelValues(errorType).Inc()
}

func (m *Metrics) RecordAPILatency(path string, duration float64) {
	m.APILatency.WithLabelValues(path).Observe(duration)

	if duration > 1.0 {
		m.SlowRequestTotal.WithLabelValues(path).Inc()
	}
}

func (m *Metrics) RecordBusinessError(module, code string) {
	m.BusinessErrorTotal.WithLabelValues(module, code).Inc()
}

func (m *Metrics) RecordCriticalError() {
	m.CriticalErrorTotal.Inc()
}

func (m *Metrics) RecordResponseSize(path string, sizeBytes float64) {
	if sizeBytes > 1024*1024 {
		m.LargeResponseTotal.WithLabelValues(path).Inc()
	}
}

// GetHistoryData get historical data
func (m *Metrics) GetHistoryData(metric string, duration time.Duration) ([]MetricDataPoint, error) {
	if m.historyStorage == nil {
		return nil, nil
	}
	return m.historyStorage.GetHistoryData(metric, duration)
}

// GetCurrentMetrics get the current system metrics
func (m *Metrics) GetCurrentMetrics() (*SystemMetrics, error) {
	if m.historyStorage == nil {
		return nil, nil
	}
	return m.historyStorage.GetCurrentData()
}

// GetCombinedHistoryData get combined historical data, each time point contains all metrics
func (m *Metrics) GetCombinedHistoryData(duration time.Duration) ([]MetricsTimePoint, error) {
	if m.historyStorage == nil {
		return nil, nil
	}

	// Get the historical data of each metric
	cpuData, err := m.historyStorage.GetHistoryData("cpu", duration)
	if err != nil {
		return nil, err
	}

	memoryData, err := m.historyStorage.GetHistoryData("memory", duration)
	if err != nil {
		return nil, err
	}

	diskData, err := m.historyStorage.GetHistoryData("disk", duration)
	if err != nil {
		return nil, err
	}

	// Build the time point mapping
	timeMap := make(map[time.Time]MetricsTimePoint)

	// Add CPU data
	for _, point := range cpuData {
		timeMap[point.Timestamp] = MetricsTimePoint{
			Timestamp: point.Timestamp,
			CPUUsage:  point.Value,
		}
	}

	// Add memory data
	for _, point := range memoryData {
		if tp, exists := timeMap[point.Timestamp]; exists {
			tp.MemoryUsage = point.Value
			timeMap[point.Timestamp] = tp
		} else {
			timeMap[point.Timestamp] = MetricsTimePoint{
				Timestamp:   point.Timestamp,
				MemoryUsage: point.Value,
			}
		}
	}

	// Add disk data
	for _, point := range diskData {
		if tp, exists := timeMap[point.Timestamp]; exists {
			tp.DiskUsage = point.Value
			timeMap[point.Timestamp] = tp
		} else {
			timeMap[point.Timestamp] = MetricsTimePoint{
				Timestamp: point.Timestamp,
				DiskUsage: point.Value,
			}
		}
	}

	// Convert the map to a slice and sort by time
	result := make([]MetricsTimePoint, 0, len(timeMap))
	for _, point := range timeMap {
		result = append(result, point)
	}

	// Sort by time
	sort.Slice(result, func(i, j int) bool {
		return result[i].Timestamp.Before(result[j].Timestamp)
	})

	return result, nil
}