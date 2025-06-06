// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNamePeriodicTask = "periodic_tasks"

// PeriodicTask mapped from table <periodic_tasks>
type PeriodicTask struct {
	ID                string    `gorm:"column:id;primaryKey" json:"id"`                                  // ID
	SceneAutomationID string    `gorm:"column:scene_automation_id;not null" json:"scene_automation_id"` // Scene automation ID (foreign key, cascading delete)
	TaskType          string    `gorm:"column:task_type;not null" json:"task_type"`       // Task type (HOUR, DAY, WEEK, MONTH, CRON)
	Param             string    `gorm:"column:params;not null" json:"params"`                             // Parameters
	ExecutionTime     time.Time `gorm:"column:execution_time;not null" json:"execution_time"` // Execution time
	Enabled           string    `gorm:"column:enabled;not null" json:"enabled"`     // Is enabled (Y-enabled, N-disabled)
	Remark            *string   `gorm:"column:remark" json:"remark"`                                       // Remarks
	ExpirationTime    int64     `gorm:"column:expiration_time;not null" json:"expiration_time"` // Expiration time (default is 5 minutes after execution time), in minutes
}

// TableName PeriodicTask's table name
func (*PeriodicTask) TableName() string {
	return TableNamePeriodicTask
}
