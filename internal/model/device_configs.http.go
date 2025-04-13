package model

import "time"

type CreateDeviceConfigReq struct {
	Name             string  `json:"name"  validate:"required,max=99"`                  // Device configuration name
	DeviceTemplateId *string `json:"device_template_id" validate:"omitempty,max=36"`    // Device template ID
	DeviceType       string  `json:"device_type" validate:"required,max=9,oneof=1 2 3"` // Device type
	ProtocolType     *string `json:"protocol_type" validate:"omitempty,max=36"`         // Protocol type (enumerated in dictionary)
	VoucherType      *string `json:"voucher_type" validate:"omitempty,max=36"`          // Voucher type (no specific enumeration, different protocols may have different types)
	ProtocolConfig   *string `json:"protocol_config" validate:"omitempty"`              // Protocol configuration
	DeviceConnType   *string `json:"device_conn_type" validate:"omitempty,oneof=A B"`   // Device connection type (default A) A - Device connects to platform, B - Platform connects to device
	AdditionalInfo   *string `json:"additional_info" validate:"omitempty"`              // Additional information
	Description      *string `json:"description" validate:"omitempty,max=255"`          // Description
	Remark           *string `json:"remark" validate:"omitempty,max=255"`               // Remark
}

type UpdateDeviceConfigReq struct {
	Id               string     `json:"id" validate:"required,max=36"`                   // Device configuration ID
	Name             *string    `json:"name"  validate:"omitempty,max=99"`               // Device configuration name
	DeviceTemplateId *string    `json:"device_template_id" validate:"omitempty,max=36"`  // Device template ID
	ProtocolType     *string    `json:"protocol_type" validate:"omitempty,max=100"`      // Protocol type (enumerated in dictionary)
	VoucherType      *string    `json:"voucher_type" validate:"omitempty,max=500"`       // Voucher type (no specific enumeration, different protocols may have different types)
	ProtocolConfig   *string    `json:"protocol_config" validate:"omitempty"`            // Protocol configuration
	DeviceConnType   *string    `json:"device_conn_type" validate:"omitempty,oneof=A B"` // Device connection type (default A) A - Device connects to platform, B - Platform connects to device
	AdditionalInfo   *string    `json:"additional_info" validate:"omitempty"`            // Additional information
	Description      *string    `json:"description" validate:"omitempty,max=255"`        // Description
	Remark           *string    `json:"remark" validate:"omitempty,max=255"`             // Remark
	UpdatedAt        *time.Time `json:"updated_at" validate:"omitempty"`                 // Update time
	OtherConfig      *string    `json:"other_config" validate:"omitempty"`               // Other configuration
}

type GetDeviceConfigListByPageReq struct {
	PageReq
	DeviceTemplateId *string `json:"device_template_id" form:"device_template_id" validate:"omitempty,max=36"` // Device template ID
	DeviceType       *string `json:"device_type" form:"device_type" validate:"omitempty,max=9,oneof=1 2 3"`    // Device type
	ProtocolType     *string `json:"protocol_type" form:"protocol_type" validate:"omitempty,max=36"`           // Protocol type
	Name             *string `json:"name" form:"name" validate:"omitempty,max=99"`                             // Device configuration name
}

type GetDeviceConfigListMenuReq struct {
	DeviceConfigName *string `json:"device_config_name" form:"device_config_name" validate:"omitempty,max=99"` // Device configuration name
	DeviceType       *string `json:"device_type" form:"device_type" validate:"omitempty,max=9,oneof=1 2 3"`    // Device type
	ProtocolType     *string `json:"protocol_type" form:"protocol_type" validate:"omitempty,max=50"`           // Protocol type
}

type BatchUpdateDeviceConfigReq struct {
	DeviceConfigID string   `json:"device_config_id" validate:"required,uuid"` // Device configuration ID
	DeviceIds      []string `json:"device_ids" validate:"omitempty,max=36"`    // Array of device IDs
}

type DeviceConfigRsp struct {
	*DeviceConfig
	DeviceCount int64 `json:"device_count"`
}

type DeviceConfigConnectRes struct {
	AccessToken string `json:"access_token"` // Access token for the connection
	Basic       string `json:"basic"`        // Basic authorization information
}

// DeviceConfigsRes
// Fixed return response for device configuration detail page
type DeviceConfigsRes struct {
	ID               string    `json:"id"`                 // Device configuration ID
	Name             string    `json:"name"`               // Device configuration name
	DeviceTemplateID string    `json:"device_template_id"` // Device template ID
	DeviceType       string    `json:"device_type"`        // Device type
	ProtocolType     string    `json:"protocol_type"`      // Protocol type
	VoucherType      string    `json:"voucher_type"`       // Voucher type
	ProtocolConfig   string    `json:"protocol_config"`    // Protocol configuration form
	DeviceConnType   string    `json:"device_conn_type"`   // Device connection type (default A) A - Device connects to platform, B - Platform connects to device
	AdditionalInfo   string    `json:"additional_info"`    // Additional information
	Description      string    `json:"description"`        // Description
	CreatedAt        time.Time `json:"created_at"`         // Creation time
	UpdatedAt        time.Time `json:"updated_at"`         // Last update time
	Remark           string    `json:"remark"`             // Remark
}

type DeviceOnline struct {
	DeviceConfigId *string `json:"device_config_id"` // Device configuration ID
	DeviceId       string  `json:"device_id"`        // Device ID
	//Online         int     `json:"online"`          // Online status
	//OtherConfig    *DeviceConfigOtherConfig `json:"other_config"`  // Additional device configuration
}

type DeviceConfigOtherConfig struct {
	OnlineTimeout int `json:"online_timeout"` // Online timeout in minutes
	Heartbeat     int `json:"heartbeat"`      // Heartbeat interval in seconds
}