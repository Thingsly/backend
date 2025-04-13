package middleware

import (
	"time"

	"github.com/HustIoTPlatform/backend/pkg/metrics"

	"github.com/gin-gonic/gin"
)

// MetricsMiddleware Creates a monitoring middleware
func MetricsMiddleware(m *metrics.Metrics) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath() // Get the route path, not the actual URL

		// Record the request
		m.RecordAPIRequest(path, c.Request.Method)

		// Use defer to ensure metrics are recorded at the end of the request
		defer func() {
			// Record the response time
			duration := time.Since(start).Seconds()
			m.RecordAPILatency(path, duration)

			// Record the response size
			m.RecordResponseSize(path, float64(c.Writer.Size()))

			// Handle panic
			if err := recover(); err != nil {
				m.RecordAPIError("panic")
				m.RecordCriticalError()
				panic(err) // Re-throw the panic
			}
		}()

		c.Next()

		// Record errors
		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				if e.IsType(gin.ErrorTypePrivate) {
					m.RecordAPIError("system")
				} else {
					m.RecordAPIError("business")
				}
			}
		}
	}
}
