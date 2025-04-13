package model

type UpdateLogoReq struct {
	Id             string  `json:"id" validate:"required,max=36"`                // Id
	SystemName     string  `json:"system_name" validate:"omitempty,max=99"`      // System name
	LogoCache      *string `json:"logo_cache" validate:"omitempty,max=255"`      // Cached logo
	LogoBackground *string `json:"logo_background" validate:"omitempty,max=255"` // Site logo
	LogoLoading    *string `json:"logo_loading" validate:"omitempty,max=255"`    // Loading page logo
	HomeBackground *string `json:"home_background" validate:"omitempty,max=255"` // Home page background
	Remark         *string `json:"remark" validate:"omitempty,max=255"`
}

