// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameRGroupDevice = "r_group_device"

// RGroupDevice mapped from table <r_group_device>
type RGroupDevice struct {
	GroupID  string `gorm:"column:group_id;not null" json:"group_id"`
	DeviceID string `gorm:"column:device_id;not null" json:"device_id"`
	TenantID string `gorm:"column:tenant_id;not null" json:"tenant_id"`
}

// TableName RGroupDevice's table name
func (*RGroupDevice) TableName() string {
	return TableNameRGroupDevice
}
