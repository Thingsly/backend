package model

import "time"

type GetAlarmInfoListByPageReq struct {
	PageReq
	StartTime        *time.Time `json:"start_time" form:"start_time" validate:"omitempty"`               
	EndTime          *time.Time `json:"end_time" form:"end_time" validate:"omitempty"`                   
	AlarmLevel       *string    `json:"alarm_level" form:"alarm_level" validate:"omitempty"`             
	ProcessingResult *string    `json:"processing_result" form:"processing_result" validate:"omitempty"` 
	TenantID         string     `json:"tenant_id" validate:"omitempty"`
}

type UpdateAlarmInfoReq struct {
	Id               string  `json:"id" validate:"required,max=36"`
	ProcessingResult *string `json:"processing_result" validate:"required"`
}

type UpdateAlarmInfoBatchReq struct {
	Id                     []string `json:"id" validate:"required"`
	ProcessingResult       *string  `json:"processing_result" validate:"required"`      
	ProcessingInstructions *string  `json:"processing_instructions" validate:"required"`
}

type GetAlarmHisttoryListByPage struct {
	PageReq
	StartTime   *time.Time `json:"start_time" form:"start_time" validate:"omitempty"`     
	EndTime     *time.Time `json:"end_time" form:"end_time" validate:"omitempty"`         
	AlarmStatus *string    `json:"alarm_status" form:"alarm_status" validate:"omitempty"` 
	DeviceId    *string    `json:"device_id" form:"device_id" validate:"omitempty"`      
}

type AlarmHistoryDescUpdateReq struct {
	AlarmHistoryId string `json:"id"  validate:"required"`         
	Description    string `json:"description" validate:"required"`
}
type GetDeviceAlarmStatusReq struct {
	DeviceId string `json:"device_id" form:"device_id" validate:"required"`
}
