package model

import "time"

type CreateBoardReq struct {
	Name        string  `json:"name" validate:"required,max=255"`         // Board Name (Required, max length 255)
	Config      *string `json:"config" validate:"omitempty"`              // Board Configuration (Optional)
	HomeFlag    string  `json:"home_flag"  validate:"required,max=2"`     // Homepage Flag, default 'N', 'Y' means homepage (Required, max length 2)
	MenuFlag    string  `json:"menu_flag"`                                // Menu Flag, default 'N', 'Y' means shown in menu
	Description *string `json:"description" validate:"omitempty,max=500"` // Description (Optional, max length 500)
	Remark      *string `json:"remark" validate:"omitempty,max=255"`      // Remark (Optional, max length 255)
	TenantID    string  `json:"tenant_id" validate:"omitempty,max=36"`    // Tenant ID (Optional, max length 36)
}

type UpdateBoardReq struct {
	Id          string  `json:"id" validate:"omitempty,max=36"`           // Board ID (Optional, max length 36)
	Name        string  `json:"name" validate:"omitempty,max=255"`        // Board Name (Optional, max length 255)
	Config      *string `json:"config" validate:"omitempty"`              // Board Configuration (Optional)
	HomeFlag    string  `json:"home_flag"  validate:"omitempty,max=2"`    // Homepage Flag (Optional, max length 2)
	MenuFlag    string  `json:"menu_flag"  validate:"omitempty,max=2"`    // Menu Flag (Optional, max length 2)
	Description *string `json:"description" validate:"omitempty,max=500"` // Description (Optional, max length 500)
	Remark      *string `json:"remark" validate:"omitempty,max=255"`      // Remark (Optional, max length 255)
	TenantID    string  `json:"tenant_id" validate:"omitempty,max=36"`    // Tenant ID (Optional, max length 36)
}

type GetBoardListByPageReq struct {
	PageReq
	Name     *string `json:"name" form:"name" validate:"omitempty,max=255"`
	HomeFlag *string `json:"home_flag" form:"home_flag"  validate:"omitempty,max=2"`
}

// DeviceTrendReq - Device Trend Request
type DeviceTrendReq struct {
	TenantID *string `form:"tenant_id" json:"tenant_id" validate:"omitempty,max=36"` // Tenant ID (Optional, max length 36)
}

// DeviceTrendPoint - Trend Data Point
type DeviceTrendPoint struct {
	Timestamp     time.Time `json:"timestamp"`      // Timestamp of the data point
	DeviceTotal   int64     `json:"device_total"`   // Total number of devices
	DeviceOnline  int64     `json:"device_online"`  // Number of online devices
	DeviceOffline int64     `json:"device_offline"` // Number of offline devices
}

// DeviceTrendRes - Device Trend Response
type DeviceTrendRes struct {
	Points []DeviceTrendPoint `json:"points"` // List of trend data points
}
