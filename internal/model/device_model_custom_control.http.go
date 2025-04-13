package model

// CreateDeviceModelCustomControlReq represents the request body for creating a custom control for a device model
type CreateDeviceModelCustomControlReq struct {
	DeviceTemplateId string  `json:"device_template_id" validate:"required,max=36"` // Device template ID
	Name             string  `json:"name" validate:"required,max=36"`               // Name
	ControlType      string  `json:"control_type" validate:"required,max=50"`       // Control type
	Description      *string `json:"description" validate:"omitempty,max=500"`      // Description
	Content          *string `json:"content" validate:"omitempty"`                  // Instruction content
	EnableStatus     string  `json:"enable_status" validate:"required,max=10"`      // Enable status
	Remark           *string `json:"remark" validate:"omitempty,max=255"`           // Remarks
}

// UpdateDeviceModelCustomControlReq represents the request body for updating a custom control for a device model
type UpdateDeviceModelCustomControlReq struct {
	ID               string  `json:"id" validate:"required,max=36"`                  // ID
	DeviceTemplateId *string `json:"device_template_id" validate:"omitempty,max=36"` // Device template ID
	Name             *string `json:"name" validate:"omitempty,max=36"`               // Name
	ControlType      *string `json:"control_type" validate:"omitempty,max=50"`       // Control type
	Description      *string `json:"description" validate:"omitempty,max=500"`       // Description
	Content          *string `json:"content" validate:"omitempty"`                   // Instruction content
	EnableStatus     *string `json:"enable_status" validate:"omitempty,max=10"`      // Enable status
	Remark           *string `json:"remark" validate:"omitempty,max=255"`            // Remarks
}

