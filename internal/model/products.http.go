package model

// CreateProductReq represents the request to create a new product.
type CreateProductReq struct {
	Name           string  `json:"name" validate:"required,max=255"`             // Product name
	Description    *string `json:"description"  validate:"omitempty,max=255"`    // Product description
	ProductType    *string `json:"product_type" validate:"omitempty,max=36"`     // Product type
	ProductKey     *string `json:"product_key" validate:"omitempty,max=255"`     // Product key (if empty, the backend will auto-generate it)
	ProductModel   *string `json:"product_model" validate:"omitempty,max=100"`   // Product model
	ImageUrl       *string `json:"image_url" validate:"omitempty,max=500"`       // Product image URL
	AdditionalInfo *string `json:"additional_info" validate:"omitempty"`         // Additional information
	Remark         *string `json:"remark" validate:"omitempty,max=255"`          // Remarks
	DeviceConfigID *string `json:"device_config_id" validate:"omitempty,max=36"` // Device configuration ID
}

// UpdateProductReq represents the request to update an existing product.
type UpdateProductReq struct {
	Id           string  `json:"id" validate:"required,max=36"`              // Product ID
	Name         *string `json:"name" validate:"omitempty,max=255"`          // Product name
	Description  *string `json:"description"  validate:"omitempty,max=255"`  // Product description
	ProductModel *string `json:"product_model" validate:"omitempty,max=100"` // Product model
	ImageUrl     *string `json:"image_url" validate:"omitempty,max=500"`     // Product image URL
	ProductType  *string `json:"product_type" validate:"omitempty,max=36"`   // Product type
}

// GetProductListByPageReq represents the request to get a paginated list of products.
type GetProductListByPageReq struct {
	PageReq
	Name         *string `json:"name" form:"name" validate:"omitempty,max=255"`                    // Product name
	ProductModel *string `json:"product_model" form:"product_model"  validate:"omitempty,max=100"` // Product model
	ProductType  *string `json:"product_type" form:"product_type" validate:"omitempty,max=36"`     // Product type
}

type ProductList struct {
	Product
	DeviceConfigName *string `json:"device_config_name"`
}
