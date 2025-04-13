package model

type GetTelemetryHistoryDataReq struct {
	DeviceID  string `json:"device_id" form:"device_id" validate:"required,max=36"`
	Key       string `json:"key" form:"key" validate:"required,max=255"`
	StartTime int64  `json:"start_time" form:"start_time" validate:"required"`
	EndTime   int64  `json:"end_time" form:"end_time"  validate:"required"`
}

type DeleteTelemetryDataReq struct {
	DeviceID string `json:"device_id" form:"device_id" validate:"required,max=36"`
	Key      string `json:"key" form:"key" validate:"required,max=255"`
}

type GetTelemetryCurrentDataKeysReq struct {
	DeviceID string   `json:"device_id" form:"device_id" validate:"required,max=36"`
	Keys     []string `json:"key" form:"keys" validate:"required,max=255"`
}

type GetTelemetryHistoryDataByPageReq struct {
	DeviceID    string `json:"device_id" form:"device_id" validate:"required,max=36"`
	Key         string `json:"key" form:"key" validate:"required,max=255"`
	StartTime   int64  `json:"start_time" form:"start_time" validate:"required"`
	EndTime     int64  `json:"end_time" form:"end_time"  validate:"required"`
	ExportExcel *bool  `json:"export_excel" form:"export_excel" validate:"omitempty"`
	Page        *int   `json:"page" form:"page" validate:"omitempty"`
	PageSize    *int   `json:"page_size" form:"page_size" validate:"omitempty"`
}

// GetTelemetrySetLogsListByPageReq represents the request structure for fetching telemetry set logs with pagination and optional filters such as device ID, status, and operation type.
type GetTelemetrySetLogsListByPageReq struct {
	PageReq
	DeviceId      string  `json:"device_id" form:"device_id" validate:"required,max=36"`               // Device ID (required)
	Status        *string `json:"status" form:"status" validate:"omitempty,oneof=1 2"`                // Status: 1 - Sent successfully, 2 - Failed
	OperationType *string `json:"operation_type" form:"operation_type" validate:"omitempty,oneof=1 2"` // Operation type: 1 - Manual operation, 2 - Automatic trigger
}

// SimulationTelemetryDataReq represents the request structure for simulating telemetry data by providing a command.
type SimulationTelemetryDataReq struct {
	Command string `json:"command" form:"command" validate:"required,max=500"` // Command for mosquitto_pub (required)
}

// ServeEchoDataReq represents the request structure for serving echo data for a specific device.
type ServeEchoDataReq struct {
	DeviceId string `json:"device_id" form:"device_id" validate:"required,max=36"` // Device ID (required)
}

// GetTelemetryStatisticReq represents the request structure for retrieving telemetry statistics for a specific device, with optional start and end times, time range, aggregation window, and aggregation function.
type GetTelemetryStatisticReq struct {
	DeviceId          string `json:"device_id" form:"device_id" validate:"required,max=36"` // Device ID (required)
	Key               string `json:"key" form:"key" validate:"required"`                    // Statistic key (required)
	StartTime         int64  `json:"start_time" form:"start_time" validate:"omitempty"`      // Start time (optional)
	EndTime           int64  `json:"end_time" form:"end_time" validate:"omitempty"`          // End time (optional)
	TimeRange         string `json:"time_range" form:"time_range" validate:"required"`       // Time range (required)
	AggregateWindow   string `json:"aggregate_window" form:"aggregate_window" validate:"required"` // Aggregation interval (required)
	AggregateFunction string `json:"aggregate_function" form:"aggregate_function" validate:"omitempty,max=255"` // Aggregation method (optional)
	IsExport          bool   `json:"is_export" form:"is_export" validate:"omitempty"`        // Whether to export the data (optional)
}
