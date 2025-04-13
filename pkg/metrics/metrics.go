// pkg/metrics/metrics.go
package metrics

import (
	"os"
	"runtime"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/shirou/gopsutil/process"
	"github.com/sirupsen/logrus"
)

type Metrics struct {
	MemoryUsage     prometheus.Gauge
	MemoryAllocated prometheus.Gauge
	MemoryObjects   prometheus.Gauge
	CPUUsage        prometheus.Gauge
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

			cpuPercent, err := process.Percent(0)
			if err == nil {
				m.CPUUsage.Set(cpuPercent)
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
