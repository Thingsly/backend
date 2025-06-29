// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameDeviceTriggerCondition = "device_trigger_condition"

// DeviceTriggerCondition mapped from table <device_trigger_condition>
type DeviceTriggerCondition struct {
	ID                   string  `gorm:"column:id;primaryKey" json:"id"` // Id
	SceneAutomationID    string  `gorm:"column:scene_automation_id;not null" json:"scene_automation_id"` // Scene automation ID (foreign key - cascade delete)
	Enabled              string  `gorm:"column:enabled;not null" json:"enabled"` // Is enabled
	GroupID              string  `gorm:"column:group_id;not null" json:"group_id"` // UUID
	TriggerConditionType string  `gorm:"column:trigger_condition_type;not null" json:"trigger_condition_type"` // Condition type: 10 = single device, 11 = device type, 22 = time range
	TriggerSource        *string `gorm:"column:trigger_source" json:"trigger_source"` // Trigger source: for type 10 = device ID, for type 11 = device config ID
	TriggerParamType     *string `gorm:"column:trigger_param_type" json:"trigger_param_type"` // Parameter type: telemetry (TEL), attribute (ATTR), event (EVT), status (STATUS)
	TriggerParam         *string `gorm:"column:trigger_param" json:"trigger_param"` // Trigger parameter (e.g., temperature)
	TriggerOperator      *string `gorm:"column:trigger_operator" json:"trigger_operator"` // Operator: =, !=, >, <, >=, <=, between, in
	TriggerValue         string  `gorm:"column:trigger_value;not null" json:"trigger_value"` // Trigger value. For types 10/11 with 'between': format is min-max (e.g., 2-6). For 'in': comma-separated values. For type 22: example format 137|HH:mm:ss+00:00|HH:mm:ss+00:00
	Remark               *string `gorm:"column:remark" json:"remark"` // Remark
	TenantID             string  `gorm:"column:tenant_id;not null" json:"tenant_id"` // Tenant ID
}

// TableName DeviceTriggerCondition's table name
func (*DeviceTriggerCondition) TableName() string {
	return TableNameDeviceTriggerCondition
}
