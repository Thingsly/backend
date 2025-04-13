package model

type GetEventDatasListByPageReq struct {
	PageReq
	DeviceId string  `json:"device_id" form:"device_id" validate:"required,max=36"` // Device ID
	Identify *string `json:"identify" form:"identify" validate:"omitempty,max=36"`  // Data identifier
}
