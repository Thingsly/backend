package api

import (
	"time"
)

// // DeviceConfig mapped from table <device_configs>
// type DeviceConfig struct {
//     ID               string    `gorm:"column:id;primaryKey" json:"id"`                                           // Id
//     Name             string    `gorm:"column:name;not null" json:"name"`                                         // Name
//     DeviceTemplateID *string   `gorm:"column:device_template_id" json:"device_template_id"`        // Device Template ID
//     DeviceType       string    `gorm:"column:device_type;not null" json:"device_type"`                     // Device Type
//     ProtocolType     *string   `gorm:"column:protocol_type" json:"protocol_type"`                        // Protocol Type
//     VoucherType      *string   `gorm:"column:voucher_type" json:"voucher_type"`                          // Voucher Type
//     ProtocolConfig   *string   `gorm:"column:protocol_config" json:"protocol_config"`      // Protocol Form Configuration
//     DeviceConnType   *string   `gorm:"column:device_conn_type" json:"device_conn_type"` // Device Connection Type (default A) A - Device connects to platform, B - Platform connects to device
//     AdditionalInfo   *string   `gorm:"column:additional_info;default:{}" json:"additional_info"` // Additional Information
//     Description      *string   `gorm:"column:description" json:"description"`                              // Description
//     TenantID         string    `gorm:"column:tenant_id;not null" json:"tenant_id"`                           // Tenant ID
//     CreatedAt        time.Time `gorm:"column:created_at;not null" json:"created_at"`                     // Creation Time
//     UpdatedAt        time.Time `gorm:"column:updated_at;not null" json:"updated_at"`                       // Update Time
//     Remark           *string   `gorm:"column:remark" json:"remark"`                                              // Remark
// }

type DeviceConfigReadSchema struct {
	ID               string    `json:"id"`                 // Id
	Name             string    `json:"name"`               // Name
	DeviceTemplateID *string   `json:"device_template_id"` // Device Template ID
	DeviceType       string    `json:"device_type"`        // Device Type
	ProtocolType     *string   `json:"protocol_type"`      // Protocol Type
	VoucherType      *string   `json:"voucher_type"`       // Voucher Type
	ProtocolConfig   *string   `json:"protocol_config"`    // Protocol Form Configuration
	DeviceConnType   *string   `json:"device_conn_type"`   // Device Connection Type (default A) A - Device connects to platform, B - Platform connects to device
	AdditionalInfo   *string   `json:"additional_info"`    // Additional Information
	Description      *string   `json:"description"`        // Description
	TenantID         string    `json:"tenant_id"`          // Tenant ID
	CreatedAt        time.Time `json:"created_at"`         // Creation Time
	UpdatedAt        time.Time `json:"updated_at"`         // Update Time
	Remark           *string   `json:"remark"`             // Remark
}

type GetDeviceConfigResponse struct {
	Code    int                    `json:"code" example:"200"`
	Message string                 `json:"message" example:"success"`
	Data    DeviceConfigReadSchema `json:"data"`
}

type GetDeviceConfigListResponse struct {
	Code    int                     `json:"code" example:"200"`
	Message string                  `json:"message" example:"success"`
	Data    GetDeviceConfigListData `json:"data"`
}

type GetDeviceConfigListData struct {
	Total int64                      `json:"total"`
	List  []DeviceTemplateReadSchema `json:"list"`
}
