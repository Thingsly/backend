// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameSceneInfo = "scene_info"

// SceneInfo mapped from table <scene_info>
// SceneInfo represents a configurable scene that can be used for automation or grouping.
type SceneInfo struct {
	ID          string     `gorm:"column:id;primaryKey" json:"id"`                   // Scene ID
	Name        string     `gorm:"column:name;not null" json:"name"`                 // Scene name
	Description *string    `gorm:"column:description" json:"description"`            // Description
	TenantID    string     `gorm:"column:tenant_id;not null" json:"tenant_id"`       // Tenant ID
	Creator     string     `gorm:"column:creator;not null" json:"creator"`           // Creator ID
	Updator     *string    `gorm:"column:updator" json:"updator"`                    // Last updated by (user ID)
	CreatedAt   time.Time  `gorm:"column:created_at;not null" json:"created_at"`     // Creation timestamp
	UpdatedAt   *time.Time `gorm:"column:updated_at" json:"updated_at"`              // Last update timestamp
}

// TableName SceneInfo's table name
func (*SceneInfo) TableName() string {
	return TableNameSceneInfo
}
