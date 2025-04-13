package model

type GetCommandSetLogsListByPageReq struct {
	PageReq
	DeviceId      string  `json:"device_id" form:"device_id" validate:"required,max=36"`               // Device ID
	Identify      *string `json:"identify" form:"identify" validate:"omitempty,max=36"`                // Data identifier
	Status        *string `json:"status" form:"status" validate:"omitempty,oneof=1 2 3 4"`             // Status: 1 - Sent successfully, 2 - Failure, 3 - Returned successfully, 4 - Returned failure
	OperationType *string `json:"operation_type" form:"operation_type" validate:"omitempty,oneof=1 2"` // Operation type: 1 - Manual operation, 2 - Automatically triggered
	IdentifyName  *string `json:"identify_name" form:"identify_name" validate:"omitempty,max=100"`     // Data identifier name
}
