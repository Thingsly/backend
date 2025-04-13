package api

import "time"

type DeviceTemplateReadSchema struct {
	ID                string    `json:"id"`                  // Id
	Name              string    `json:"name"`                // Template Name
	Author            *string   `json:"author"`              // Author
	Version           *string   `json:"version"`             // Version Number
	Description       *string   `json:"description"`         // Description
	TenantID          string    `json:"tenant_id"`           // Tenant ID
	CreatedAt         time.Time `json:"created_at"`          // Creation Time
	UpdatedAt         time.Time `json:"updated_at"`          // Update Time
	Flag              *int16    `json:"flag" example:"1"`    // Flag, default 1
	Label             *string   `json:"label"`               // Label
	DeviceModelConfig *string   `json:"device_model_config"` // Device Model Configuration
	WebChartConfig    *string   `json:"web_chart_config"`    // Web Chart Configuration
	AppChartConfig    *string   `json:"app_chart_config"`    // App Chart Configuration
}

type GetDeviceTemplateListResponse struct {
	Code    int                       `json:"code" example:"200"`
	Message string                    `json:"message" example:"success"`
	Data    GetDeviceTemplateListData `json:"data"`
}

type GetDeviceTemplateListData struct {
	Total int64                      `json:"total"`
	List  []DeviceTemplateReadSchema `json:"list"`
}

type GetDeviceTemplateResponse struct {
	Code    int                      `json:"code" example:"200"`
	Message string                   `json:"message" example:"success"`
	Data    DeviceTemplateReadSchema `json:"data"`
}
