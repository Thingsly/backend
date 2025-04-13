// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameSceneActionInfo = "scene_action_info"

// SceneActionInfo mapped from table <scene_action_info>
// SceneActionInfo represents the configuration of an action within an automation scene.
type SceneActionInfo struct {
	ID              string     `gorm:"column:id;primaryKey" json:"id"`                                                             // ID
	SceneID         string     `gorm:"column:scene_id;not null" json:"scene_id"`                                                   // Scene ID (with cascading delete)
	ActionTarget    string     `gorm:"column:action_target;not null" json:"action_target"`                                         // Action target ID (device ID, device config ID, scene ID, or alarm ID)
	ActionType      string     `gorm:"column:action_type;not null" json:"action_type"`                                             // Action type: 10 - Single Device, 11 - Device Category, 20 - Trigger Scene, 30 - Trigger Alarm, 40 - Service
	ActionParamType *string    `gorm:"column:action_param_type" json:"action_param_type"`                                          // Parameter type: TEL - Telemetry, ATTR - Attribute, CMD - Command
	ActionParam     *string    `gorm:"column:action_param" json:"action_param"`                                                    // Action parameter
	ActionValue     *string    `gorm:"column:action_value" json:"action_value"`                                                    // Target value
	CreatedAt       time.Time  `gorm:"column:created_at;not null" json:"created_at"`                                               // Creation time
	UpdatedAt       *time.Time `gorm:"column:updated_at" json:"updated_at"`                                                        // Last update time
	TenantID        string     `gorm:"column:tenant_id;not null" json:"tenant_id"`                                                 // Tenant ID
	Remark          *string    `gorm:"column:remark" json:"remark"`                                                                // Remarks
}

// TableName SceneActionInfo's table name
func (*SceneActionInfo) TableName() string {
	return TableNameSceneActionInfo
}
