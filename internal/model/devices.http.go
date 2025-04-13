package model

import "time"

type CreateDeviceReq struct {
	Name           *string `json:"name" validate:"omitempty,max=255"`            // Device name
	Voucher        *string `json:"voucher" validate:"omitempty,max=500"`         // Authentication token
	DeviceNumber   *string `json:"device_number" validate:"omitempty,max=36"`    // Device number
	ProductID      *string `json:"product_id" validate:"omitempty,max=36"`       // Product ID
	ParentID       *string `json:"parent_id" validate:"omitempty,max=36"`        // Parent device ID
	Protocol       *string `json:"protocol" validate:"omitempty,max=36"`         // Protocol
	Label          *string `json:"label" validate:"omitempty,max=255"`           // Tag/Label
	Location       *string `json:"location" validate:"omitempty,max=36"`         // Location
	SubDeviceAddr  *string `json:"sub_device_addr" validate:"omitempty,max=36"`  // Sub-device address
	CurrentVersion *string `json:"current_version" validate:"omitempty,max=36"`  // Current version
	AdditionalInfo *string `json:"additional_info" validate:"omitempty"`         // Additional information
	ProtocolConfig *string `json:"protocol_config" validate:"omitempty"`         // Protocol configuration
	Remark1        *string `json:"remark1" validate:"omitempty,max=255"`         // Remark 1
	Remark2        *string `json:"remark2" validate:"omitempty,max=255"`         // Remark 2
	Remark3        *string `json:"remark3" validate:"omitempty,max=255"`         // Remark 3
	DeviceConfigId *string `json:"device_config_id" validate:"omitempty,max=36"` // Device configuration ID
	AccessWay      *string `json:"access_way" validate:"omitempty,max=36"`       // Access method
	Description    *string `json:"description" validate:"omitempty,max=500"`     // Description
}

type BatchCreateDeviceReq struct {
	ServiceAccessId string `json:"service_access_id" validate:"required,max=36"` // Service access point ID
	DeviceList      []struct {
		DeviceName     string  `json:"device_name" validate:"required,max=255"`     // Device name
		DeviceNumber   string  `json:"device_number" validate:"required,max=36"`    // Device number
		Description    *string `json:"description" validate:"omitempty,max=500"`    // Description
		DeviceConfigId string  `json:"device_config_id" validate:"required,max=36"` // Device configuration ID
	} `json:"device_list" validate:"required"`
}

type UpdateDeviceReq struct {
	Id             string  `json:"id" validate:"required,max=36"`                // Device ID
	Name           *string `json:"name" validate:"omitempty,max=255"`            // Device name
	Voucher        *string `json:"voucher" validate:"omitempty,max=500"`         // Authentication token
	DeviceNumber   *string `json:"device_number" validate:"omitempty,max=36"`    // Device number
	ProductID      *string `json:"product_id" validate:"omitempty,max=36"`       // Product ID
	ParentID       *string `json:"parent_id" validate:"omitempty,max=36"`        // Parent device ID
	Label          *string `json:"label" validate:"omitempty,max=255"`           // Tag/Label
	Location       *string `json:"location" validate:"omitempty,max=100"`        // Location
	SubDeviceAddr  *string `json:"sub_device_addr" validate:"omitempty,max=36"`  // Sub-device address
	CurrentVersion *string `json:"current_version" validate:"omitempty,max=36"`  // Current version
	AdditionalInfo *string `json:"additional_info" validate:"omitempty"`         // Additional information
	ProtocolConfig *string `json:"protocol_config" validate:"omitempty"`         // Protocol configuration
	Remark1        *string `json:"remark1" validate:"omitempty,max=255"`         // Remark 1
	Remark2        *string `json:"remark2" validate:"omitempty,max=255"`         // Remark 2
	Remark3        *string `json:"remark3" validate:"omitempty,max=255"`         // Remark 3
	DeviceConfigId *string `json:"device_config_id" validate:"omitempty,max=36"` // Device configuration ID
	AccessWay      *string `json:"access_way" validate:"omitempty,max=36"`       // Access method
	Description    *string `json:"description" validate:"omitempty,max=500"`     // Description
	IsOnline       *int16  `json:"is_online" validate:"omitempty"`               // Online status (1: online, 0: offline)
}

type ActiveDeviceReq struct {
	DeviceNumber string `json:"device_number" validate:"required,max=36"` // Device number
	Name         string `json:"name" validate:"max=255"`                  // Device name
}

type GetDeviceListByPageReq struct {
	PageReq
	ActivateFlag      *string `json:"activate_flag" form:"activate_flag" validate:"omitempty,max=36"`       // Activation status
	DeviceNumber      *string `json:"device_number" form:"device_number" validate:"omitempty,max=36"`       // Device number
	IsEnabled         *string `json:"is_enabled" form:"is_enabled" validate:"omitempty,max=36"`             // Enabled status
	ProductID         *string `json:"product_id" form:"product_id" validate:"omitempty,max=36"`             // Product ID
	ProtocolType      *string `json:"protocol_type" form:"protocol_type" validate:"omitempty,max=36"`       // Protocol type
	Label             *string `json:"label" form:"label" validate:"omitempty,max=255"`                      // Tag/Label
	Name              *string `json:"name" form:"name" validate:"omitempty,max=255"`                        // Device name
	CurrentVersion    *string `json:"current_version" form:"current_version" validate:"omitempty,max=36"`   // Current version
	GroupId           *string `json:"group_id" form:"group_id" validate:"omitempty,max=36"`                 // Device group ID
	DeviceConfigId    *string `json:"device_config_id" form:"device_config_id" validate:"omitempty,max=36"` // Device configuration ID
	IsOnline          *int    `json:"is_online" form:"is_online" validate:"omitempty,max=36"`               // Online status
	WarnStatus        *string `json:"warn_status" form:"warn_status" validate:"omitempty,max=36"`           // Alarm status (reserved)
	Search            *string `json:"search" form:"search" validate:"omitempty,max=36"`                     // Fuzzy match: name or number
	AccessWay         *string `json:"access_way" form:"access_way" validate:"omitempty,max=36"`             // Access method
	BatchNumber       *string `json:"batch_number" form:"batch_number" validate:"omitempty"`                // Batch number
	DeviceType        *string `json:"device_type" form:"device_type" validate:"omitempty,oneof=1 2 3"`          // Device type: 1-Gateway, 2-Subdevice, 3-Sub-subdevice
	ServiceIdentifier *string `json:"service_identifier" form:"service_identifier" validate:"omitempty,max=36"` // Service identifier
	ServiceAccessID   *string `json:"service_access_id" form:"service_access_id" validate:"omitempty,max=36"`   // Service access point ID
}

type GetDeviceListByPageRsp struct {
	ID               string     `json:"id"`                 // Device ID
	DeviceNumber     string     `json:"device_number"`      // Device number
	Name             string     `json:"name"`               // Device name
	DeviceConfigID   string     `json:"device_config_id"`   // Device configuration ID
	DeviceConfigName string     `json:"device_config_name"` // Device configuration name
	Ts               *time.Time `json:"ts"`                 // Last push timestamp
	ActivateFlag     string     `json:"activate_flag"`      // Activation status
	ActivateAt       *time.Time `json:"activate_at"`        // Activation time
	BatchNumber      string     `json:"batch_number"`       // Batch number
	CurrentVersion   string     `json:"current_version"`    // Current version
	CreatedAt        *time.Time `json:"created_at"`         // Creation timestamp
	IsOnline         int        `json:"is_online"`          // Online status
	Location         string     `json:"location"`           // Location
	AccessWay        string     `json:"access_way"`         // Access method
	ProtocolType     string     `json:"protocol_type"`      // Protocol type
	DeviceStatus     int        `json:"device_status"`      // Device status
	WarnStatus       string     `json:"warn_status"`        // Warning status (Y: warning, N: normal)
	DeviceType       string     `json:"device_type"`        // Device type: 1-Gateway, 2-Subdevice, 3-Sub-subdevice
}

type CreateDeviceGroupReq struct {
	ParentId    *string `json:"parent_id" validate:"omitempty,max=36"`    // Parent group ID
	Name        string  `json:"name" validate:"required,max=255"`         // Device group name
	Description *string `json:"description" validate:"omitempty,max=255"` // Description
	Remark      *string `json:"remark" validate:"omitempty,max=255"`      // Remark
}

type UpdateDeviceGroupReq struct {
	Id          string  `json:"id" validate:"required,max=36"`            // Group ID
	ParentId    string  `json:"parent_id" validate:"required,max=36"`     // Parent group ID
	Name        string  `json:"name" validate:"required,max=255"`         // Group name
	Description *string `json:"description" validate:"omitempty,max=255"` // Description
	Remark      *string `json:"remark" validate:"omitempty,max=255"`      // Remark
}

type GetDeviceGroupsListByPageReq struct {
	PageReq
	ParentId *string `json:"parent_id" form:"parent_id" validate:"omitempty,max=36"` // Parent group ID
	Name     *string `json:"name" form:"name" validate:"omitempty,max=255"`          // Group name
}

type GetDeviceListByGroup struct {
	PageReq
	GroupId string `json:"group_id" form:"group_id" validate:"required,max=36"` // Device group ID
}

type GetDeviceListByGroupRsp struct {
	GroupId            string `json:"group_id"`
	Id                 string `json:"id"`
	DeviceNumber       string `json:"device_number"`
	Name               string `json:"name"`
	Device_config_name string `json:"device_config_name"`
}

type GetDeviceGroupListByDeviceIdReq struct {
	DeviceId string `json:"device_id" form:"device_id" validate:"required,max=36"` // Device ID
}

type CreateDeviceGroupRelationReq struct {
	GroupId      string   `json:"group_id" validate:"required,max=36"` // Device group ID
	DeviceIDList []string `json:"device_id_list" validate:"required"`  // List of device IDs
}

type DeleteDeviceGroupRelationReq struct {
	GroupId  string `json:"group_id" form:"group_id" validate:"required,max=36"`   // Device group ID
	DeviceId string `json:"device_id" form:"device_id" validate:"required,max=36"` // Device ID
}

type CreateDevicePreRegisterReq struct {
	ProductID      string  `json:"product_id" validate:"required,max=36"`             // Product ID
	BatchNumber    string  `json:"batch_number" validate:"required,max=36"`           // Batch number
	CurrentVersion *string `json:"current_version" validate:"omitempty,max=36"`       // Firmware version
	DeviceCount    *int    `json:"device_count" validate:"omitempty,min=1,max=10000"` // Number of devices (required if create_type is 1)
	CreateType     string  `json:"create_type" validate:"required,oneof=1 2"`         // Creation type: 1 - auto, 2 - file
	BatchFile      *string `json:"batch_file" validate:"omitempty,max=500"`           // Batch file path
}

type GetDevicePreRegisterListByPageReq struct {
	PageReq
	ProductID      string  `json:"product_id" form:"product_id" validate:"omitempty,max=36"`                       // Product ID
	BatchNumber    *string `json:"batch_number" form:"batch_number" validate:"omitempty"`                          // Batch number
	DeviceNumber   *string `json:"device_number" form:"device_number" validate:"omitempty"`                        // Device number
	IsEnabled      *string `json:"is_enabled" form:"is_enabled" validate:"omitempty"`                              // Enabled status
	ActivateFlag   *string `json:"activate_flag" form:"activate_flag" validate:"omitempty,oneof=active inactive"` // Activation status
	Name           *string `json:"name" form:"name" validate:"omitempty"`                                          // Device name
	DeviceConfigID *string `json:"device_config_id" form:"device_config_id" validate:"omitempty"`                  // Device configuration ID
}

type GetDevicePreRegisterListByPageRsp struct {
	ID             string     `json:"id"`              // Device ID
	Name           string     `json:"name"`            // Device name
	DeviceNumber   string     `json:"device_number"`   // Device number
	ActivateFlag   string     `json:"activate_flag"`   // Activation status
	ActivateAt     *time.Time `json:"activate_at"`     // Activation time
	BatchNumber    string     `json:"batch_number"`    // Batch number
	CurrentVersion string     `json:"current_version"` // Current firmware version
	CreatedAt      *time.Time `json:"created_at"`      // Creation time
}

type ExportPreRegisterReq struct {
	ProductID    string  `json:"product_id" form:"product_id" validate:"required,max=36"`                        // Product ID
	BatchNumber  *string `json:"batch_number" form:"batch_number" validate:"omitempty,max=36"`                   // Batch number
	ActivateFlag *string `json:"activate_flag" form:"activate_flag" validate:"omitempty,oneof=active inactive"` // Activation status
}

// Remove sub-device
type RemoveSonDeviceReq struct {
	SubDeviceId string `json:"sub_device_id" validate:"required,max=36"` // Sub-device ID
}

// Retrieve device dropdown list
type GetDeviceMenuReq struct {
	GroupId    string `json:"group_id" form:"group_id" validate:"omitempty,max=36"`        // Device group ID
	DeviceName string `json:"device_name" form:"device_name" validate:"omitempty,max=255"` // Device name
	BindConfig int    `json:"bind_config" form:"bind_config" validate:"omitempty"`         // Binding state: 0 - all, 1 - bound, 2 - unbound
}

type GetTenantDeviceListReq struct {
	ID               string `json:"id"`                 // Device ID
	Name             string `json:"name"`               // Device name
	DeviceConfigID   string `json:"device_config_id"`   // Device configuration ID
	DeviceConfigName string `json:"device_config_name"` // Device configuration name
}

type CreateSonDeviceRes struct {
	ID    string `json:"id" validate:"required,max=36"`       // Parent device ID
	SonID string `json:"son_id" validate:"required,max=3600"` // Sub-device ID(s), separated by commas
}

type DeviceConnectFormReq struct {
	DeviceID string `query:"device_id" form:"device_id" json:"device_id" validate:"required,max=36"`
}

type DeviceConnectFormRes struct {
	DataKey     string                       `json:"dataKey"`
	Label       string                       `json:"label"`
	Placeholder string                       `json:"placeholder"`
	Type        string                       `json:"type"`
	Validate    DeviceConnectFormValidateRes `json:"validate"`
}

type DeviceConnectFormValidateRes struct {
	Message  string `json:"message,omitempty"`
	Required bool   `json:"required"`
	Type     string `json:"type"`
}

type DeviceIDReq struct {
	DeviceID string `query:"device_id" form:"device_id" json:"device_id" validate:"required,max=36"`
}

type GetVoucherTypeReq struct {
	DeviceType   string `json:"device_type"  form:"device_type"  validate:"required,max=36,oneof=1 2 3"`
	ProtocolType string `json:"protocol_type"  form:"protocol_type"  validate:"required,max=255"`
}

type UpdateDeviceVoucherReq struct {
	DeviceID string `json:"device_id" validate:"required,max=36"`
	Voucher  any    `json:"voucher" validate:"required"`
}

type GetSubListResp struct {
	Name          string `json:"name"`
	Id            string `json:"id"`
	SubDeviceAddr string `json:"subDeviceAddr"`
}

type GetDeviceTemplateChartSelectReq struct {
	GroupID string `json:"group_id" form:"group_id" validate:"required,max=36"`
}

type GetActionByDeviceConfigIDReq struct {
	DeviceConfigID string `json:"device_config_id" form:"device_config_id" validate:"required,max=36"`
}

type GetActionByDeviceIDReq struct {
	DeviceID string `json:"device_id" form:"device_id" validate:"required,max=36"`
}

// Update device configuration
type ChangeDeviceConfigReq struct {
	DeviceID       string  `json:"device_id" validate:"required,max=36"` // Device ID
	DeviceConfigID *string `json:"device_config_id" validate:"max=36"`   // Device configuration ID
}

type GatewayRegisterReq struct {
	GatewayId string `json:"gateway_id"`
	TenantId  string `json:"tenant_id"`
	Model     string `json:"model"`
}

type GatewayRegisterRes struct {
	MqttUsername string `json:"mqtt_username"`
	MqttPassword string `json:"mqtt_password"`
	MqttClientId string `json:"mqtt_client_id"`
}

type DeviceVoucher struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type DeviceRegisterReq struct {
	Type      string          `json:"type"`
	DeviceId  string          `json:"device_id"`
	Registers []DeviceSubItem `json:"registers"`
}

type DeviceSubItem struct {
	SubAddr  string `json:"sub_addr"`
	Model    string `json:"model"`
	Protocol string `json:"protocol"`
}

type DeviceRegisterRes struct {
	Type         string                          `json:"type"`
	Status       string                          `json:"status"`
	Message      string                          `json:"message"`
	RegistersRes map[string]DeviceSubRegisterRes `json:"registersRes"`
}

type DeviceSubRegisterRes struct {
	Result    int    `json:"result"`
	Errorcode string `json:"errorcode"`
	Message   string `json:"message"`
	SubAddr   string `json:"sub_addr"`
}
