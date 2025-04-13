// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameServicePlugin = "service_plugins"

// ServicePlugin mapped from table <service_plugins>
// ServicePlugin represents the details of a service plugin, including its identifier, type, version, and configuration.
type ServicePlugin struct {
	ID                string     `gorm:"column:id;primaryKey" json:"id"`                                   // Service ID
	Name              string     `gorm:"column:name;not null" json:"name"`                                 // Service name
	ServiceIdentifier string     `gorm:"column:service_identifier;not null" json:"service_identifier"`    // Service identifier
	ServiceType       int32      `gorm:"column:service_type;not null" json:"service_type"`                 // Service type: 1 - Access protocol, 2 - Access service
	LastActiveTime    *time.Time `gorm:"column:last_active_time" json:"last_active_time"`                  // Last active time of the service
	Version           *string    `gorm:"column:version" json:"version"`                                    // Version number
	CreateAt          time.Time  `gorm:"column:create_at;not null" json:"create_at"`                       // Creation time
	UpdateAt          time.Time  `gorm:"column:update_at;not null" json:"update_at"`                       // Update time
	Description       *string    `gorm:"column:description" json:"description"`                            // Description
	ServiceConfig     *string    `gorm:"column:service_config" json:"service_config"`                      // Service configuration: Configuration for access protocol and access service
	Remark            *string    `gorm:"column:remark" json:"remark"`                                      // Remarks
}

// TableName ServicePlugin's table name
func (*ServicePlugin) TableName() string {
	return TableNameServicePlugin
}
