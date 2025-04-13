package model

// CreateOTAUpgradePackageReq represents a request to create an OTA upgrade package.
type CreateOTAUpgradePackageReq struct {
	Name           string  `json:"name" validate:"required,max=200"`                    // Upgrade Package Name
	Version        string  `json:"version" validate:"required,max=36"`                  // Version Number
	TargetVersion  *string `json:"target_version" validate:"omitempty,max=36"`          // Target Version Number
	DeviceConfigID string  `json:"device_config_id" validate:"required,max=36"`         // Device Configuration ID
	Module         *string `json:"module" validate:"omitempty,max=36"`                  // Module Name
	PackageType    *int16  `json:"package_type" validate:"required,oneof=1 2"`          // Package Type: 1-Differential, 2-Full Package
	SignatureType  *string `json:"signature_type" validate:"omitempty,oneof=MD5 SHA256"`// Signature Algorithm: MD5, SHA256
	AdditionalInfo *string `json:"additional_info" validate:"omitempty" example:"{}"`   // Additional Information (JSON format)
	Description    *string `json:"description" validate:"omitempty,max=500"`            // Description
	PackageUrl     *string `json:"package_url" validate:"omitempty,max=500"`            // Package URL
	Remark         *string `json:"remark" validate:"omitempty,max=255"`                 // Remarks
}

// UpdateOTAUpgradePackageReq represents a request to update an OTA upgrade package.
type UpdateOTAUpgradePackageReq struct {
	Id             string  `json:"id" validate:"required,max=36"`                       // Upgrade Package ID
	Name           string  `json:"name" validate:"omitempty,max=200"`                   // Upgrade Package Name
	Version        string  `json:"version" validate:"omitempty,max=36"`                 // Version Number
	TargetVersion  *string `json:"target_version" validate:"omitempty,max=36"`          // Target Version Number
	DeviceConfigID string  `json:"device_config_id" validate:"omitempty,max=36"`        // Device Configuration ID
	Module         *string `json:"module" validate:"omitempty,max=36"`                  // Module Name
	PackageType    *int16  `json:"package_type" validate:"omitempty,oneof=1 2"`         // Package Type
	SignatureType  *string `json:"signature_type" validate:"omitempty,oneof=MD5 SHA256"`// Signature Algorithm
	AdditionalInfo *string `json:"additional_info" validate:"omitempty"`                // Additional Information (JSON format)
	Description    *string `json:"description" validate:"omitempty,max=500"`            // Description
	PackageUrl     *string `json:"package_url" validate:"omitempty,max=500"`            // Package URL
	Remark         *string `json:"remark" validate:"omitempty,max=255"`                 // Remarks
}

// GetOTAUpgradePackageLisyByPageReq represents a request to get a paginated list of OTA upgrade packages.
type GetOTAUpgradePackageLisyByPageReq struct {
	PageReq
	DeviceConfigID string `json:"device_configs_id" form:"device_config_id" validate:"omitempty,max=36" example:"uuid"` // Device Configuration ID
	Name           string `json:"name" form:"name" validate:"omitempty,max=200"`                                          // Upgrade Package Name
}

// GetOTAUpgradeTaskListByPageRsp represents a response with a paginated list of OTA upgrade tasks.
type GetOTAUpgradeTaskListByPageRsp struct {
	OtaUpgradePackage
	DeviceConfigName string `json:"device_config_name" validate:"omitempty,max=200"` // Device Configuration Name
}

