package model

import (
	"encoding/json"
	"time"
)

type CreateUserReq struct {
	AdditionalInfo *json.RawMessage `json:"additional_info" validate:"omitempty,max=10000"`
	Email          string           `json:"email"  validate:"required,email"`
	Password       string           `json:"password" validate:"required,min=6,max=255"`
	Name           *string          `json:"name" validate:"omitempty,min=2,max=50"`
	PhoneNumber    string           `json:"phone_number" validate:"required,max=50"`
	RoleIDs        []string         `json:"userRoles" validate:"omitempty"`
	Remark         *string          `json:"remark" validate:"omitempty,max=255"`
}

type LoginReq struct {
	Email    string `json:"email" validate:"required" example:"test@test.vn"`
	Password string `json:"password" validate:"required,min=6,max=512" example:"123456"`
	Salt     string `json:"salt" validate:"omitempty,max=512"`
}

type LoginRsp struct {
	Token     *string `gorm:"column:token" json:"token"`
	ExpiresIn int64   `json:"expires_in"`
}

type UserListReq struct {
	PageReq
	Email       *string `json:"email" form:"email" validate:"omitempty"`
	PhoneNumber *string `json:"phone_number" form:"phone_number" validate:"omitempty,max=50"`
	Name        *string `json:"name" form:"name" validate:"omitempty,max=50"`
	Status      *string `json:"status" form:"status" validate:"omitempty,oneof=N F"`
}

type UpdateUserReq struct {
	ID             string     `json:"id" validate:"required,uuid"`
	AdditionalInfo *string    `json:"additional_info" validate:"omitempty,max=10000"`
	Email          *string    `json:"email"  validate:"omitempty,email"`
	Name           *string    `json:"name" validate:"omitempty,min=2,max=50"`
	PhoneNumber    *string    `json:"phone_number" validate:"omitempty,max=50"`
	Remark         *string    `json:"remark" validate:"omitempty,max=255"`
	Status         *string    `json:"status" validate:"omitempty,oneof=N F"`
	Password       *string    `json:"password" validate:"omitempty,max=255"`
	UpdatedAt      *time.Time `json:"updated_at" validate:"omitempty"`
	RoleIDs        []string   `json:"userRoles" validate:"omitempty"`
}

type UpdateUserInfoReq struct {
	ID        string     `json:"id" validate:"required"`
	Name      *string    `json:"name" validate:"omitempty,min=2,max=50"`
	Remark    *string    `json:"remark" validate:"omitempty,max=255"`
	Password  *string    `json:"password" validate:"omitempty,min=6,max=255"`
	UpdatedAt *time.Time `json:"updated_at" validate:"omitempty"`
	Salt      string     `json:"salt"`
}

type TransformUserReq struct {
	BecomeUserID string `json:"become_user_id" validate:"required,uuid"`
}

type ResetPasswordReq struct {
	Email      string `json:"email" validate:"required,email"`
	VerifyCode string `json:"verify_code" validate:"required"`
	Password   string `json:"password" validate:"required,min=6,max=255"`
}

type EmailRegisterReq struct {
	Email           string  `json:"email" validate:"required,email"`
	VerifyCode      string  `json:"verify_code" validate:"required"`
	Password        string  `json:"password" validate:"required"`
	ConfirmPassword *string `json:"confirm_password" validate:"omitempty"`
	PhoneNumber     string  `json:"phone_number" validate:"required"`
	PhonePrefix     string  `json:"phone_prefix" validate:"required"`
	Salt            *string `json:"salt" validate:"omitempty"`
}
