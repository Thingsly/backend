// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameOtaUpgradeTaskDetail = "ota_upgrade_task_details"

// OtaUpgradeTaskDetail mapped from table <ota_upgrade_task_details>
type OtaUpgradeTaskDetail struct {
	ID                string     `gorm:"column:id;primaryKey" json:"id"`                                             // ID
	OtaUpgradeTaskID  string     `gorm:"column:ota_upgrade_task_id;not null" json:"ota_upgrade_task_id"` // OTA Upgrade Task ID (foreign key, cascading delete)
	DeviceID          string     `gorm:"column:device_id;not null" json:"device_id"`                       // Device ID (foreign key, prevent delete)
	Step              *int16     `gorm:"column:steps" json:"steps"`                                           // Upgrade progress (1-100)
	Status            int16      `gorm:"column:status;not null" json:"status"`      // Status (1-pending, 2-sent, 3-upgrading, 4-upgrade success, 5-upgrade fail, 6-cancelled)
	StatusDescription *string    `gorm:"column:status_description" json:"status_description"`                      // Status description
	UpdatedAt         *time.Time `gorm:"column:updated_at" json:"updated_at"` // Last updated time
	Remark            *string    `gorm:"column:remark" json:"remark"`                                                    // Remarks
}

// TableName OtaUpgradeTaskDetail's table name
func (*OtaUpgradeTaskDetail) TableName() string {
	return TableNameOtaUpgradeTaskDetail
}
