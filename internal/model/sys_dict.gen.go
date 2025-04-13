// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameSysDict = "sys_dict"

// SysDict mapped from table <sys_dict>
type SysDict struct {
	ID        string    `gorm:"column:id;primaryKey" json:"id"`               
	DictCode  string    `gorm:"column:dict_code;not null" json:"dict_code"`  
	DictValue string    `gorm:"column:dict_value;not null" json:"dict_value"`  
	CreatedAt time.Time `gorm:"column:created_at;not null" json:"created_at"` 
	Remark    *string   `gorm:"column:remark" json:"remark"`                   
}

// TableName SysDict's table name
func (*SysDict) TableName() string {
	return TableNameSysDict
}
