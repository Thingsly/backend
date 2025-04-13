// internal/model/open_api_keys.http.go
package model

// OpenAPIKeyListReq represents a request to query the list of API keys.
type OpenAPIKeyListReq struct {
	PageReq        // Inherits the base pagination request
	Status  *int16 `json:"status" form:"status" validate:"omitempty,oneof=0 1"` // Status: 0 - Disabled, 1 - Enabled
}

// CreateOpenAPIKeyReq represents a request to create an API key.
type CreateOpenAPIKeyReq struct {
	TenantID string `json:"tenant_id" validate:"required,max=36"` // Tenant ID
	Name     string `json:"name" validate:"omitempty,max=200"`    // Name
}

// UpdateOpenAPIKeyReq represents a request to update an API key.
type UpdateOpenAPIKeyReq struct {
	ID     string  `json:"id" validate:"required,max=36"`         // Primary key ID
	Status *int16  `json:"status" validate:"omitempty,oneof=0 1"` // Status: 0 - Disabled, 1 - Enabled
	Name   *string `json:"name" validate:"omitempty,max=200"`     // Name
}
