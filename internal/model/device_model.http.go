package model

const (
	DEVICE_MODEL_TELEMETRY  = "DEVICE_MODEL_TELEMETRY"
	DEVICE_MODEL_ATTRIBUTES = "DEVICE_MODEL_ATTRIBUTES"
	DEVICE_MODEL_EVENTS     = "DEVICE_MODEL_EVENTS"
	DEVICE_MODEL_COMMANDS   = "DEVICE_MODEL_COMMANDS"
)

// Device Model Creation (telemetry and attributes)
type CreateDeviceModelReq struct {
	DeviceTemplateId string  `json:"device_template_id" validate:"required,max=36"`          // Device template ID
	DataName         *string `json:"data_name" validate:"omitempty,max=255"`                 // Data name
	DataIdentifier   string  `json:"data_identifier" validate:"required,max=255"`            // Data identifier
	ReadWriteFlag    *string `json:"read_write_flag" validate:"omitempty,max=10,oneof=R RW"` // Read/Write flag R-Read RW-ReadWrite
	DataType         *string `json:"data_type" validate:"omitempty,max=50"`                  // Data type: String, Number, Boolean
	Unit             *string `json:"unit" validate:"omitempty,max=50"`                       // Unit
	Description      *string `json:"description" validate:"omitempty,max=500"`               // Description
	AdditionalInfo   *string `json:"additional_info" validate:"omitempty"`                   // Additional info
	Remark           *string `json:"remark" validate:"omitempty,max=255"`                    // Remark
}

// Device Model Creation (events and commands)
type CreateDeviceModelV2Req struct {
	DeviceTemplateId string  `json:"device_template_id" validate:"required,max=36"` // Device template ID
	DataName         *string `json:"data_name" validate:"omitempty,max=255"`        // Data name
	DataIdentifier   string  `json:"data_identifier" validate:"required,max=255"`   // Data identifier
	Params           *string `json:"params" validate:"omitempty"`                   // Parameters
	Description      *string `json:"description" validate:"omitempty,max=500"`      // Description
	AdditionalInfo   *string `json:"additional_info" validate:"omitempty"`          // Additional info
	Remark           *string `json:"remark" validate:"omitempty,max=255"`           // Remark
}

// Device Model Update (telemetry and attributes)
type UpdateDeviceModelReq struct {
	ID             string  `json:"id" validate:"required,max=36"`                          // ID
	DataName       *string `json:"data_name" validate:"omitempty,max=255"`                 // Data name
	DataIdentifier string  `json:"data_identifier" validate:"required,max=255"`            // Data identifier
	ReadWriteFlag  *string `json:"read_write_flag" validate:"omitempty,max=10,oneof=R RW"` // Read/Write flag R-Read RW-ReadWrite
	DataType       *string `json:"data_type" validate:"omitempty,max=50"`                  // Data type: String, Number, Boolean
	Unit           *string `json:"unit" validate:"omitempty,max=50"`                       // Unit
	Description    *string `json:"description" validate:"omitempty,max=500"`               // Description
	AdditionalInfo *string `json:"additional_info" validate:"omitempty"`                   // Additional info
	Remark         *string `json:"remark" validate:"omitempty,max=255"`                    // Remark
}

// Device Model Update (events and commands)
type UpdateDeviceModelV2Req struct {
	ID             string  `json:"id" validate:"required,max=36"`               // ID
	DataName       *string `json:"data_name" validate:"omitempty,max=255"`      // Data name
	DataIdentifier string  `json:"data_identifier" validate:"required,max=255"` // Data identifier
	Params         *string `json:"params" validate:"omitempty"`                 // Parameters
	Description    *string `json:"description" validate:"omitempty,max=500"`    // Description
	AdditionalInfo *string `json:"additional_info" validate:"omitempty"`        // Additional info
	Remark         *string `json:"remark" validate:"omitempty,max=255"`         // Remark
}

// Get Device Model List by Page Request
type GetDeviceModelListByPageReq struct {
	PageReq
	DeviceTemplateId string  `json:"device_template_id" form:"device_template_id"  validate:"required,max=36"` // Device template ID
	EnableStatus     *string `json:"enable_status"  form:"enable_status" validate:"omitempty,max=10"`          // Enable status
}


type GetModelSourceATRes struct {
	DataSourceTypeRes string     `json:"data_source_type"`
	Options           []*Options `json:"options"`
}

type Options struct {
	Key      string     `json:"key"`
	Label    *string    `json:"label"`
	DataType *string    `json:"data_type"`
	Enum     []EnumItem `json:"enum"`
}

type EnumItem struct {
	ValueType   string `json:"value_type"`
	Value       int    `json:"value"`
	Description string `json:"description"`
}

// Device Model Custom Command Creation Request
type CreateDeviceModelCustomCommandReq struct {
	DeviceTemplateId string  `json:"device_template_id" validate:"required,max=36"` // Device template ID
	ButtomName       string  `json:"buttom_name" validate:"required,max=36"`        // Button name
	DataIdentifier   string  `json:"data_identifier" validate:"required,max=255"`   // Data identifier
	Description      *string `json:"description" validate:"omitempty,max=500"`      // Description
	Instruct         *string `json:"instruct" validate:"omitempty"`                 // Instruction content
	EnableStatus     string  `json:"enable_status" validate:"required,max=10"`      // Enable status
	Remark           *string `json:"remark" validate:"omitempty,max=255"`           // Remark
}

// Device Model Custom Command Update Request
type UpdateDeviceModelCustomCommandReq struct {
	ID             string  `json:"id" validate:"required,max=36"`               // ID
	ButtomName     string  `json:"buttom_name" validate:"required,max=36"`      // Button name
	DataIdentifier string  `json:"data_identifier" validate:"required,max=255"` // Data identifier
	Description    *string `json:"description" validate:"omitempty,max=500"`    // Description
	Instruct       *string `json:"instruct" validate:"omitempty"`               // Instruction content
	EnableStatus   string  `json:"enable_status" validate:"required,max=10"`    // Enable status
	Remark         *string `json:"remark" validate:"omitempty,max=255"`         // Remark
}

// Device Metrics Chart Request
type GetDeviceMetricsChartReq struct {
	DeviceID          string  `json:"device_id" form:"device_id" validate:"required"`                                                                                                                                                  // Device ID
	DataType          string  `json:"data_type" form:"data_type" validate:"required,oneof=telemetry attribute command event"`                                                                                                          // Device data type
	DataMode          string  `json:"data_mode" form:"data_mode" validate:"required,oneof=latest history"`                                                                                                                             // Data mode
	Key               string  `json:"key" form:"key" validate:"required"`                                                                                                                                                              // Data identifier
	TimeRange         *string `json:"time_range" form:"time_range" validate:"omitempty,oneof=last_5m last_15m last_30m last_1h last_3h last_6h last_12h last_24h last_3d last_7d last_15d last_30d last_60d last_90d last_6m last_1y"` // Time range
	AggregateWindow   *string `json:"aggregate_window" form:"aggregate_window" validate:"omitempty,oneof=no_aggregate 30s 1m 2m 5m 10m 30m 1h 3h 6h 1d 7d 1mo"`                                                                        // Aggregate interval
	AggregateFunction *string `json:"aggregate_function" form:"aggregate_function" validate:"omitempty,oneof=avg max min sum diff"`                                                                                                    // Aggregate function
}

// Device Metrics Chart Data Response
type DeviceMetricsChartData struct {
	DeviceID          string       `json:"device_id"`          // Device ID
	DataType          string       `json:"data_type"`          // Device data type
	Key               string       `json:"key"`                // Data identifier
	AggregateWindow   *string      `json:"aggregate_window"`   // Aggregate interval
	AggregateFunction *string      `json:"aggregate_function"` // Aggregate function
	TimeRange         *string      `json:"time_range"`         // Time range
	Value             *interface{} `json:"value"`              // Latest value
	Timestamp         *int64       `json:"timestamp"`          // Latest value timestamp
	Points            *[]DataPoint `json:"points"`             // List of data points
}

// Data Point
type DataPoint struct {
	T int64   `json:"t"` // Timestamp
	V float64 `json:"v"` // Value
}

// Device Selector Request
type DeviceSelectorReq struct {
	PageReq
	// Whether the device template exists
	HasDeviceConfig *bool `json:"has_device_config" form:"has_device_config" validate:"omitempty"`
}

// Device Selector Response
type DeviceSelectorRes struct {
	Total int64                 `json:"total"`
	List  []*DeviceSelectorData `json:"list"`
}

type DeviceSelectorData struct {
	DeviceID   string `json:"device_id"`   // Device ID
	DeviceName string `json:"device_name"` // Device name
}
