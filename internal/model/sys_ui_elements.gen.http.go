package model

// CreateUiElementsReq represents the request structure for creating UI elements, including properties like parent ID, element type, permissions, and description.
type CreateUiElementsReq struct {
	ID           string  `json:"id"`                                        // Primary key ID
	ParentID     string  `json:"parent_id" validate:"required,max=36"`      // Parent element ID
	ElementCode  string  `json:"element_code" validate:"required,max=100"`  // Element identifier
	ElementType  int     `json:"element_type" validate:"omitempty,max=10"`  // Element type: 1 - Menu, 2 - Directory, 3 - Button, 4 - Route
	Orders       int     `json:"orders" validate:"omitempty,max=10000"`     // Sorting order
	Param1       *string `json:"param1" validate:"omitempty,max=255"`
	Param2       *string `json:"param2" validate:"omitempty,max=255"`
	Param3       *string `json:"param3" validate:"omitempty,max=255"`
	Authority    string  `json:"authority" validate:"required"`           // Permissions (multiple choices): TENANT_ADMIN - Tenant Admin, SYS_ADMIN - System Admin
	Description  *string `json:"description" validate:"omitempty,max=36"` // Description
	Remark       *string `json:"remark" validate:"omitempty,max=255"`     // Remark
	Multilingual *string `json:"multilingual" validate:"omitempty,max=255"` // Multilingual support
	RoutePath    *string `json:"route_path" validate:"omitempty,max=255"`   // Route path
}

// UpdateUiElementsReq represents the request structure for updating UI elements, including properties like parent ID, element code, and route path.
type UpdateUiElementsReq struct {
	Id           string  `json:"id" validate:"required,max=36"`            // Element ID
	ParentID     *string `json:"parent_id" form:"parent_id" validate:"required,max=36"`         // Parent element ID
	ElementCode  *string `json:"element_code" form:"element_code" validate:"required,max=100"`  // Element identifier
	ElementType  *int16  `json:"element_type" form:"element_type" validate:"omitempty,max=10"`  // Element type: 1 - Menu, 2 - Directory, 3 - Button, 4 - Route
	Orders       *int16  `json:"orders" form:"orders" validate:"omitempty,max=10000"`           // Sorting order
	Param1       *string `json:"param1" form:"param1" validate:"omitempty,max=255"`
	Param2       *string `json:"param2" form:"param2" validate:"omitempty,max=255"`
	Param3       *string `json:"param3" form:"param3" validate:"omitempty,max=255"`
	Authority    *string `json:"authority" form:"authority" validate:"required"`                // Permissions (multiple choices): TENANT_ADMIN - Tenant Admin, SYS_ADMIN - System Admin
	Description  *string `json:"description" form:"description" validate:"omitempty,max=36"`    // Description
	Multilingual *string `json:"multilingual" form:"multilingual" validate:"omitempty,max=255"` // Multilingual support
	RoutePath    *string `json:"route_path" form:"route_path" validate:"omitempty,max=255"`     // Route path
	Remark       *string `json:"remark" form:"remark" validate:"omitempty,max=255"`             // Remark
}

// UiElementsListReq represents the request structure for retrieving a list of UI elements, with filters like parent ID, element code, and authority.
type UiElementsListReq struct {
	Id           string  `json:"id" form:"id" validate:"required,max=36"`            // Element ID (required)
	ParentID     string  `json:"parent_id" form:"parent_id" validate:"required,max=36"` // Parent element ID (required)
	ElementCode  string  `json:"element_code" form:"element_code" validate:"required,max=100"` // Element identifier (required)
	ElementType  *int16  `json:"element_type" form:"element_type" validate:"omitempty,max=10"` // Element type: 1 - Menu, 2 - Directory, 3 - Button, 4 - Route
	Orders       *int16  `json:"orders" form:"orders" validate:"omitempty,max=10000"`       // Sorting order
	Param1       *string `json:"param1" form:"param1" validate:"omitempty,max=255"`        // Parameter 1
	Param2       *string `json:"param2" form:"param2" validate:"omitempty,max=255"`        // Parameter 2
	Param3       *string `json:"param3" form:"param3" validate:"omitempty,max=255"`        // Parameter 3
	Authority    string  `json:"authority" form:"authority" validate:"required"`           // Permissions (multiple options): TENANT_ADMIN - Tenant Admin, SYS_ADMIN - System Admin
	Description  *string `json:"description" form:"description" validate:"omitempty,max=36"` // Description
	Multilingual *string `json:"multilingual" form:"multilingual" validate:"omitempty,max=255"` // Multilingual support
	RoutePath    *string `json:"route_path" form:"route_path" validate:"omitempty,max=255"`  // Route path
}

// UiElementsListRsp represents the response structure for the list of UI elements, with properties such as element code, type, and child elements.
type UiElementsListRsp struct {
	ID           string               `json:"id" form:"id" validate:"required,max=36"`               // Element ID (required)
	ParentID     string               `json:"parent_id" form:"parent_id" validate:"required,max=36"` // Parent element ID (required)
	ElementCode  string               `json:"element_code" form:"element_code" validate:"required,max=100"` // Element identifier (required)
	ElementType  *int16               `json:"element_type" form:"element_type" validate:"omitempty,max=10"` // Element type: 1 - Menu, 2 - Directory, 3 - Button, 4 - Route
	Orders       *int16               `json:"orders" form:"orders" validate:"omitempty,max=10000"`    // Sorting order
	Param1       *string              `json:"param1" form:"param1" validate:"omitempty,max=255"`      // Parameter 1
	Param2       *string              `json:"param2" form:"param2" validate:"omitempty,max=255"`      // Parameter 2
	Param3       *string              `json:"param3" form:"param3" validate:"omitempty,max=255"`      // Parameter 3
	Authority    string               `json:"authority" form:"authority" validate:"required"`         // Permissions (multiple options): TENANT_ADMIN - Tenant Admin, SYS_ADMIN - System Admin
	Description  *string              `json:"description" form:"description" validate:"omitempty,max=36"` // Description
	Remark       *string              `json:"remark" form:"remark" validate:"omitempty,max=255"`       // Remark
	Multilingual *string              `json:"multilingual" form:"multilingual" validate:"omitempty,max=255"` // Multilingual support
	RoutePath    *string              `json:"route_path" form:"route_path" validate:"omitempty,max=255"`  // Route path
	Children     []*UiElementsListRsp `json:"children" form:"children"`                               // Child elements
}

// UiElementsListRsp1 represents a simplified version of the UI element response structure, focusing on key details and child elements.
type UiElementsListRsp1 struct {
	ID          string                `json:"id" form:"id" validate:"required,max=36"`               // Element ID (required)
	ParentID    string                `json:"parent_id" form:"parent_id" validate:"required,max=36"` // Parent element ID (required)
	ElementCode string                `json:"element_code" form:"element_code" validate:"required,max=100"` // Element identifier (required)
	ElementType *int16                `json:"element_type" form:"element_type" validate:"omitempty,max=10"` // Element type: 1 - Menu, 2 - Directory, 3 - Button, 4 - Route
	Description *string               `json:"description" form:"description" validate:"omitempty,max=36"`  // Description
	Children    []*UiElementsListRsp1 `json:"children" form:"children"`                               // Child elements
}

type ServeUiElementsListByPageReq struct {
	PageReq
}

func (u *SysUIElement) ToRsp() *UiElementsListRsp {
	return &UiElementsListRsp{
		ID:           u.ID,
		ParentID:     u.ParentID,
		ElementCode:  u.ElementCode,
		ElementType:  &u.ElementType,
		Orders:       u.Order_,
		Param1:       u.Param1,
		Param2:       u.Param2,
		Param3:       u.Param3,
		Authority:    u.Authority,
		Description:  u.Description,
		Remark:       u.Remark,
		Multilingual: u.Multilingual,
		RoutePath:    u.RoutePath,
		Children:     []*UiElementsListRsp{},
	}
}
func (u *SysUIElement) ToRsp1() *UiElementsListRsp1 {
	return &UiElementsListRsp1{
		ID:          u.ID,
		ParentID:    u.ParentID,
		ElementCode: u.ElementCode,
		ElementType: &u.ElementType,
		Description: u.Description,
		Children:    []*UiElementsListRsp1{},
	}
}
