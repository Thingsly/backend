package model

type FunctionsRoleValidate struct {
	RoleID       string   `json:"role_id"  valid:"Required; MaxSize(36)"`                   // Role
	FunctionsIDs []string `json:"functions_ids"  alias:"Function List" valid:"Required;MaxSize(36)"` // List of Functions
}

type RoleValidate struct {
	RoleID string `json:"role_id"  form:"role_id" valid:"Required; MaxSize(36)"` // Role
}

type RolesUserValidate struct {
	UserID   string   `json:"user_id"  valid:"Required; MaxSize(36)"`   // User
	RolesIDs []string `json:"roles_ids"   valid:"Required;MaxSize(36)"` // List of Roles
}

type UserValidate struct {
	UserID string `json:"user_id"  form:"user_id" valid:"Required; MaxSize(255)"` // User
}
