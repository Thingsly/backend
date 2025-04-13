package model

type UpdateDataPolicyReq struct {
	Id            string  `json:"id" validate:"required,max=36"`       // ID
	RetentionDays int32   `json:"retention_days" validate:"required"`  // Data retention period (days)
	Enabled       string  `json:"enabled" validate:"required"`         // Enabled status: 1-Enabled, 2-Disabled
	Remark        *string `json:"remark" validate:"required,max=2000"` // Remark
}

type GetDataPolicyListByPageReq struct {
	PageReq
}
