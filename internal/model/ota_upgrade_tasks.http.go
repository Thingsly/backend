package model

// CreateOTAUpgradeTaskReq represents the request to create an OTA upgrade task.
type CreateOTAUpgradeTaskReq struct {
	Name                string   `json:"name" validate:"required,max=200"`                  // Task Name
	OTAUpgradePackageId string   `json:"ota_upgrade_package_id" validate:"required,max=36"` // OTA Upgrade Package ID
	Description         *string  `json:"description" validate:"omitempty,max=500"`          // Description
	Remark              *string  `json:"remark" validate:"omitempty,max=255"`               // Remarks
	DeviceIdList        []string `json:"device_id_list" validate:"required"`                // List of Device IDs
}

// GetOTAUpgradeTaskDetailReq represents the request to get the details of an OTA upgrade task.
type GetOTAUpgradeTaskDetailReq struct {
	PageReq
	DeviceName       *string `json:"device_name" form:"device_name" validate:"omitempty,max=200"`                // Device Name
	TaskStatus       *int16  `json:"task_status" form:"task_status" validate:"omitempty,max=10"`                 // Task Status: 1-Waiting to Push, 2-Pushed, 3-Upgrading, 4-Upgrade Success, 5-Upgrade Failure, 6-Canceled
	OtaUpgradeTaskId string  `json:"ota_upgrade_task_id" form:"ota_upgrade_task_id" validate:"required,max=36"` // OTA Upgrade Task ID
}

// GetOTAUpgradeTaskListByPageReq represents the request to get a list of OTA upgrade tasks with pagination.
type GetOTAUpgradeTaskListByPageReq struct {
	PageReq
	OTAUpgradePackageId string `json:"ota_upgrade_package_id" form:"ota_upgrade_package_id" validate:"required,max=36"` // OTA Upgrade Package ID
}

// UpdateOTAUpgradeTaskStatusReq represents the request to update the status of an OTA upgrade task.
type UpdateOTAUpgradeTaskStatusReq struct {
	Id     string `json:"id" validate:"required,max=36"`        // Task Detail ID
	Action int16  `json:"action" validate:"required,oneof=1 6"` // Action: 1-Retry Upgrade, 6-Cancel Upgrade
}
