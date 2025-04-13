package model

import "time"


// device_id send_type payload expiry_time label
type CreateExpectedDataReq struct {
	DeviceID string     `json:"device_id" form:"device_id" validate:"required,max=36"`                                   // Device ID
	SendType string     `json:"send_type" form:"send_type" validate:"required,max=50,oneof=telemetry attribute command"` // Send type
	Payload  *string    `json:"payload" form:"payload" validate:"omitempty,max=9999"`                                    // Data content
	Expiry   *time.Time `json:"expiry" form:"expiry" validate:"omitempty"`                                               // Expiry time
	Label    *string    `json:"label" form:"label" validate:"omitempty,max=100"`                                         // Label
	Identify *string    `json:"identify" form:"identify" validate:"omitempty,max=100"`                                   // Identifier
}

type DeleteExpectedDataReq struct {
	ID string `json:"id" form:"id" validate:"required,max=36"` // Expected data ID
}

type GetExpectedDataPageReq struct {
	PageReq
	DeviceID string  `json:"device_id" form:"device_id" validate:"required,max=36"`  // Device ID
	SendType *string `json:"send_type" form:"send_type" validate:"omitempty,max=50"` // Send type
	Label    *string `json:"label" form:"label" validate:"omitempty,max=100"`        // Label
	Status   *string `json:"status" form:"status" validate:"omitempty,max=50"`       // Status
}

