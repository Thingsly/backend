package model

import "time"

type CreateDataScriptReq struct {
	Name            string  `json:"name" validate:"required,max=99"`
	DeviceConfigId  string  `json:"device_config_id"  validate:"required,max=36"`
	Content         *string `json:"content" validate:"omitempty"`
	ScriptType      string  `json:"script_type" validate:"omitempty"`
	LastAnalogInput *string `json:"last_analog_input" validate:"omitempty"`
	Description     *string `json:"description" validate:"omitempty,max=255"`
	Remark          *string `json:"remark" validate:"omitempty,max=255"`
}

type UpdateDataScriptReq struct {
	Id              string     `json:"id" validate:"required,max=36"` // ID
	Name            string     `json:"name" validate:"required,max=99"` // Script name
	DeviceConfigId  string     `json:"device_config_id"  validate:"required,max=36"` // Device configuration ID
	Content         *string    `json:"content" validate:"omitempty"` // Script content
	ScriptType      string     `json:"script_type" validate:"required,oneof=A B C D E F H"` // Script type: A - Telemetry report preprocessing, B - Telemetry downlink preprocessing, C - Attribute report preprocessing, D - Attribute downlink preprocessing, E - Command downlink preprocessing, F - Event report preprocessing, H - Event downlink preprocessing
	LastAnalogInput *string    `json:"last_analog_input" validate:"omitempty"` // Last analog input
	Description     *string    `json:"description" validate:"omitempty,max=255"` // Description
	Remark          *string    `json:"remark" validate:"omitempty,max=255"` // Remark
	UpdatedAt       *time.Time `json:"updated_at" validate:"omitempty"` // Last updated time
}

type GetDataScriptListByPageReq struct {
	PageReq
	DeviceConfigId *string `json:"device_config_id" form:"device_config_id" validate:"required,max=36"`
	ScriptType     *string `json:"script_type" form:"script_type" validate:"omitempty"`
}

type QuizDataScriptReq struct {
	Content     string `json:"content" validate:"omitempty"`
	AnalogInput string `json:"last_analog_input" validate:"omitempty"`
	Topic       string `json:"topic" validate:"omitempty"`
}

type EnableDataScriptReq struct {
	Id         string `json:"id" validate:"required,max=36"`
	EnableFlag string `json:"enable_flag" validate:"required,oneof=Y N"`
}
