package model

type GetAttributeSetLogsListByPageReq struct {
	PageReq
	DeviceId      string  `json:"device_id" form:"device_id" validate:"required,max=36"`                // Device ID
	Status        *string `json:"status" form:"status" validate:"omitempty,oneof=1 2 3 4"`              // Status: 1 - Success, 2 - Failure, 3 - Response Success, 4 - Response Failure
	OperationType *string `json:"operation_type" form:"operation_type" validate:"omitempty,oneof=1 2"`   // Operation type: 1 - Manual operation, 2 - Automatic trigger
}

type AttributePutMessage struct {
	DeviceID string `json:"device_id" form:"device_id" validate:"required,max=36"`
	Value    string `json:"value" form:"value" validate:"required"`
}

type AttributeGetMessageReq struct {
	DeviceID string   `json:"device_id" form:"device_id" validate:"required,max=36"`
	Keys     []string `json:"keys" form:"keys" validate:"max=9999"`
}

type GetDataListByKeyReq struct {
	DeviceId string `json:"device_id" form:"device_id" validate:"required,max=36"`
	Key      string `json:"key" form:"key" validate:"required,max=255"`
}
