package model

import "time"

// CreateRoleReq represents the request payload for creating a new role.
type CreateRoleReq struct {
	Name        string  `json:"name" validate:"required,max=255"`         // Role name
	Description *string `json:"description" validate:"omitempty,max=500"` // Role description
}

// UpdateRoleReq represents the request payload for updating an existing role.
type UpdateRoleReq struct {
	Id          string     `json:"id" validate:"required,max=36"`          // Role ID
	Name        string     `json:"name" validate:"required,max=255"`       // Role name
	Description *string    `json:"description" validate:"omitempty,max=500"` // Role description
	UpdatedAt   *time.Time `json:"updated_at" validate:"omitempty"`        // Update timestamp (not required from frontend)
	Authority   *string    `json:"authority" validate:"omitempty"`         // Permissions
}

type GetRoleListByPageReq struct {
	PageReq
	Name *string `json:"name" form:"name" validate:"omitempty,max=255"`
}
