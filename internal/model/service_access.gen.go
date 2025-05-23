// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameServiceAccess = "service_access"

// ServiceAccess mapped from table <service_access>
// ServiceAccess represents the information about service access, including configuration and access credentials.
type ServiceAccess struct {
	ID                  string    `gorm:"column:id;primaryKey" json:"id"`                             // Access ID
	Name                string    `gorm:"column:name;not null" json:"name"`                             // Name
	ServicePluginID     string    `gorm:"column:service_plugin_id;not null" json:"service_plugin_id"` // Service Plugin ID
	Voucher             string    `gorm:"column:voucher;not null" json:"voucher"`                      // Access credentials
	Description         *string   `gorm:"column:description" json:"description"`                       // Description
	ServiceAccessConfig *string   `gorm:"column:service_access_config" json:"service_access_config"`  // Service configuration
	Remark              *string   `gorm:"column:remark" json:"remark"`                                 // Remarks
	CreateAt            time.Time `gorm:"column:create_at;not null" json:"create_at"`                  // Creation timestamp
	UpdateAt            time.Time `gorm:"column:update_at;not null" json:"update_at"`                  // Update timestamp
	TenantID            string    `gorm:"column:tenant_id;not null" json:"tenant_id"`                  // Tenant ID
}

// TableName ServiceAccess's table name
func (*ServiceAccess) TableName() string {
	return TableNameServiceAccess
}
