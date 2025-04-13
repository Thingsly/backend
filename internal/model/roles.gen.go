// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameRole = "roles"

// Role mapped from table <roles>
type Role struct {
	ID          string     `gorm:"column:id;primaryKey" json:"id"`        
	Name        string     `gorm:"column:name;not null" json:"name"`      
	Description *string    `gorm:"column:description" json:"description"` 
	CreatedAt   *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   *time.Time `gorm:"column:updated_at" json:"updated_at"`
	TenantID    *string    `gorm:"column:tenant_id" json:"tenant_id"`  
}

// TableName Role's table name
func (*Role) TableName() string {
	return TableNameRole
}
