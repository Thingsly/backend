package model

import "time"

// CreateOperationLogReq represents a request to create an operation log.
type CreateOperationLogReq struct {
	IP              string  `json:"ip" validate:"required,max=36"`                // Request IP
	Path            *string `json:"path" validate:"omitempty,max=2000"`           // Request URL
	UserID          string  `json:"user_id" validate:"required,max=36"`           // Operator User ID
	Name            *string `json:"name" validate:"omitempty,max=255"`            // API Name
	Latency         int64   `json:"latency" validate:"omitempty"`                 // Latency in milliseconds
	RequestMessage  *string `json:"request_message" validate:"omitempty,max=2000"` // Request content
	ResponseMessage *string `json:"response_message" validate:"omitempty,max=2000"` // Response content
	TenantID        string  `json:"tenant_id" validate:"required,max=36"`          // Tenant ID
	Remark          *string `json:"remark" validate:"omitempty,max=255"`           // Remarks
}

// GetOperationLogListByPageReq represents a request to get a paginated list of operation logs.
type GetOperationLogListByPageReq struct {
	PageReq
	IP        *string    `json:"ip" form:"ip" validate:"omitempty,max=36"`                            // Request IP
	StartTime *time.Time `json:"start_time,omitempty" form:"start_time" validate:"omitempty,max=50"`  // Start Date
	EndTime   *time.Time `json:"end_time,omitempty" form:"end_time" validate:"omitempty,max=50"`      // End Date
	UserName  *string    `json:"username" form:"username" validate:"omitempty,max=255"`               // Username
	Method    *string    `json:"method" form:"method" validate:"omitempty,max=255"`                   // Method
}

// GetOperationLogListByPageRsp represents the response of a paginated list of operation logs.
type GetOperationLogListByPageRsp struct {
	ID              string     `json:"id"`               // Primary key
	IP              string     `json:"ip"`               // Request IP
	Path            *string    `json:"path"`             // Request URL
	UserID          string     `json:"user_id"`          // Operator User ID
	Name            *string    `json:"name"`             // API Name
	Latency         int64      `json:"latency"`          // Latency in milliseconds
	RequestMessage  *string    `json:"request_message"`  // Request content
	ResponseMessage *string    `json:"response_message"` // Response content
	TenantID        string     `json:"tenant_id"`        // Tenant ID
	CreatedAt       *time.Time `json:"created_at"`       // Creation time
	Remark          *string    `json:"remark"`           // Remarks
	UserName        *string    `json:"username"`         // Username
	Email           *string    `json:"email"`            // Email
}

