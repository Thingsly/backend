package service

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	dal "github.com/Thingsly/backend/internal/dal"
	model "github.com/Thingsly/backend/internal/model"
	query "github.com/Thingsly/backend/internal/query"
	"github.com/Thingsly/backend/mqtt/publish"
	"github.com/Thingsly/backend/pkg/common"
	global "github.com/Thingsly/backend/pkg/global"
	utils "github.com/Thingsly/backend/pkg/utils"

	"github.com/go-basic/uuid"
	"github.com/sirupsen/logrus"
)

type OTA struct{}

func (*OTA) CreateOTAUpgradePackage(req *model.CreateOTAUpgradePackageReq, tenantID string) error {
	var ota = model.OtaUpgradePackage{}
	ota.ID = uuid.New()
	ota.Name = req.Name
	ota.Version = req.Version
	ota.TargetVersion = req.TargetVersion

	ota.DeviceConfigID = req.DeviceConfigID
	ota.Module = req.Module
	ota.PackageType = *req.PackageType
	ota.SignatureType = req.SignatureType

	fileurl := *req.PackageUrl
	filepath := strings.Replace(fileurl, "/api/v1/ota/download", "", 1)
	signature, err := utils.FileSign(filepath, *req.SignatureType)
	if err != nil {
		return err
	}
	ota.Signature = &signature

	ota.AdditionalInfo = req.AdditionalInfo
	defaultAdditionalInfo := "{}"
	if req.AdditionalInfo == nil || *req.AdditionalInfo == "" {
		ota.AdditionalInfo = &defaultAdditionalInfo
	}
	ota.Description = req.Description
	ota.PackageURL = req.PackageUrl
	ota.TenantID = &tenantID

	t := time.Now().UTC()
	ota.CreatedAt = t
	ota.UpdatedAt = &t
	ota.Remark = req.Remark
	err = dal.CreateOtaUpgradePackage(&ota)
	return err
}

func (*OTA) UpdateOTAUpgradePackage(req *model.UpdateOTAUpgradePackageReq) error {

	oldota, err := dal.GetOtaUpgradePackageByID(req.Id)
	if err != nil {
		return err
	}

	var ota = model.OtaUpgradePackage{}
	ota.ID = req.Id

	ota.Name = req.Name
	// ota.Version = req.Version
	// ota.TargetVersion = req.TargetVersion
	// ota.DeviceConfigsID = req.DeviceConfigsID
	// ota.Module = req.Module
	// ota.PackageType = *req.PackageType
	// ota.SignatureType = req.SignatureType
	ota.AdditionalInfo = req.AdditionalInfo
	ota.Description = req.Description
	ota.PackageURL = req.PackageUrl
	if req.PackageUrl != oldota.PackageURL {
		fileurl := *req.PackageUrl
		filepath := strings.Replace(fileurl, "/api/v1/ota/download", "", 1)
		signature, err := utils.FileSign(filepath, *req.SignatureType)
		if err != nil {
			return err
		}
		ota.Signature = &signature
	}

	t := time.Now().UTC()
	ota.UpdatedAt = &t
	ota.Remark = req.Remark
	info, err := dal.UpdateOtaUpgradePackage(&ota)
	if err != nil {
		return err
	}
	if info.RowsAffected == 0 {
		return fmt.Errorf("no data updated")
	}
	return nil
}

func (*OTA) DeleteOTAUpgradePackage(packageId string) error {
	err := dal.DeleteOtaUpgradePackage(packageId)
	return err
}

func (*OTA) GetOTAUpgradePackageListByPage(req *model.GetOTAUpgradePackageLisyByPageReq, userClaims *utils.UserClaims) (map[string]interface{}, error) {
	total, list, err := dal.GetOtaUpgradePackageListByPage(req, userClaims.TenantID)
	if err != nil {
		return nil, err
	}
	packageListRspMap := make(map[string]interface{})
	packageListRspMap["total"] = total
	packageListRspMap["list"] = list
	return packageListRspMap, nil

}

func (o *OTA) CreateOTAUpgradeTask(req *model.CreateOTAUpgradeTaskReq) error {
	tasks, err := dal.CreateOTAUpgradeTaskWithDetail(req)
	if err == nil {
		go func() {
			for _, t := range tasks {
				o.PushOTAUpgradePackage(t)
			}
		}()
	}
	return err
}

func (*OTA) DeleteOTAUpgradeTask(id string) error {
	err := dal.DeleteOTAUpgradeTask(id)
	return err
}

func (*OTA) GetOTAUpgradeTaskListByPage(req *model.GetOTAUpgradeTaskListByPageReq) (map[string]interface{}, error) {
	total, list, err := dal.GetOtaUpgradeTaskListByPage(req)
	if err != nil {
		return nil, err
	}
	dataMap := make(map[string]interface{})
	dataMap["total"] = total
	dataMap["list"] = list
	return dataMap, nil
}

func (*OTA) GetOTAUpgradeTaskDetailListByPage(req *model.GetOTAUpgradeTaskDetailReq) (map[string]interface{}, error) {
	total, list, statistics, err := dal.GetOtaUpgradeTaskDetailListByPage(req)
	if err != nil {
		return nil, err
	}
	dataMap := make(map[string]interface{})
	dataMap["total"] = total
	dataMap["statistics"] = statistics
	dataMap["list"] = list
	return dataMap, nil
}

// Device status modification (request parameters: 1 - Cancel upgrade, 2 - Retry upgrade)
// 1 - Pending push, 2 - Pushed, 3 - Upgrading, change to "Cancelled"
// 5 - Upgrade failed, change to "Pending push"
// 4 - Upgrade successful, 6 - Already cancelled, do not modify
func (o *OTA) UpdateOTAUpgradeTaskStatus(req *model.UpdateOTAUpgradeTaskStatusReq) error {
	// Retrieve the task details by task ID
	taskDetail, err := query.OtaUpgradeTaskDetail.Where(query.OtaUpgradeTaskDetail.ID.Eq(req.Id)).First()
	if err != nil {
		return err
	}

	// If the status is 4 (Upgrade successful) or 6 (Cancelled), do not modify
	if taskDetail.Status == 4 || taskDetail.Status == 6 {
		return fmt.Errorf("the task status cannot be modified")
	}

	// A task that has failed cannot be cancelled
	if req.Action == 6 && taskDetail.Status == 5 {
		return fmt.Errorf("the task status cannot be modified")
	}

	// A task in status 1 (Pending push), 2 (Pushed), or 3 (Upgrading) cannot be retried
	if req.Action == 1 && taskDetail.Status <= 3 {
		return fmt.Errorf("the task is upgrading")
	}

	// Current time for updating the task
	t := time.Now().UTC()

	// Action 6: Cancel upgrade
	if req.Action == 6 {
		taskDetail.Status = 6
		taskDetail.UpdatedAt = &t
		desc := "Manually cancelled upgrade"
		taskDetail.StatusDescription = &desc
		_, err := query.OtaUpgradeTaskDetail.Updates(taskDetail)
		return err
	}

	// Action 1: Retry upgrade
	if req.Action == 1 {
		desc := "Manually restarting upgrade"
		startStep := int16(0)
		taskDetail.Status = 1
		taskDetail.UpdatedAt = &t
		taskDetail.StatusDescription = &desc
		taskDetail.Step = &startStep

		_, err := query.OtaUpgradeTaskDetail.Updates(taskDetail)
		if err != nil {
			return err
		}
		// After retrying, push the upgrade package
		err = o.PushOTAUpgradePackage(taskDetail)
		return err
	}

	// Return error if no valid action was performed
	return err
}

func (*OTA) PushOTAUpgradePackage(taskDetail *model.OtaUpgradeTaskDetail) error {
	// Check if the device is online
	device := &model.Device{}
	device, err := query.Device.Where(query.Device.ID.Eq(taskDetail.DeviceID)).First()
	if err != nil {
		return err
	}
	if device.IsOnline != 1 {
		// Update device upgrade task status if the device is offline
		taskDetail.Status = 5
		desc := "Device offline"
		taskDetail.StatusDescription = &desc
		t := time.Now().UTC()
		taskDetail.UpdatedAt = &t
		_, err := query.OtaUpgradeTaskDetail.Updates(taskDetail)
		if err != nil {
			return err
		}
		return fmt.Errorf("the device is offline")
	}

	// Check if the device has other ongoing upgrade tasks
	count, err := query.OtaUpgradeTaskDetail.Where(query.OtaUpgradeTaskDetail.DeviceID.Eq(taskDetail.DeviceID), query.OtaUpgradeTaskDetail.Status.Lt(4)).Count()
	if err != nil {
		return err
	}
	if count > 0 {
		// Update device upgrade task status if there's an ongoing upgrade
		taskDetail.Status = 5
		desc := "Previous upgrade not completed"
		taskDetail.StatusDescription = &desc
		t := time.Now().UTC()
		taskDetail.UpdatedAt = &t
		_, err := query.OtaUpgradeTaskDetail.Updates(taskDetail)
		if err != nil {
			return err
		}
		return fmt.Errorf("the device is upgrading")
	}

	// Push the upgrade package
	taskQuery, err := query.OtaUpgradeTask.Select(query.OtaUpgradeTask.ID).Where(query.OtaUpgradeTask.ID.Eq(taskDetail.OtaUpgradeTaskID)).First()
	if err != nil {
		return err
	}
	otataskid := taskQuery.ID
	otapackage, err := query.OtaUpgradePackage.Where(query.OtaUpgradePackage.ID.Eq(otataskid)).First()
	if err != nil {
		return err
	}

	// Prepare the OTA message
	var otamsg = make(map[string]interface{})
	// Get a random 9-digit number and convert it to string
	randNum, err := common.GetRandomNineDigits()
	if err != nil {
		return err
	}
	otamsg["id"] = randNum
	otamsg["code"] = "200"

	// Prepare the OTA message parameters
	var otamsgparams = make(map[string]interface{})
	otamsgparams["version"] = otapackage.Version
	otamsgparams["size"] = "0"
	otamsgparams["url"] = global.OtaAddress + strings.TrimPrefix(*otapackage.PackageURL, ".")
	otamsgparams["signMethod"] = otapackage.SignatureType
	otamsgparams["sign"] = ""
	otamsgparams["module"] = otapackage.Module

	// Parse additional configuration data into a map
	var m map[string]interface{}
	err = json.Unmarshal([]byte(*otapackage.AdditionalInfo), &m)
	if err != nil {
		logrus.Error(err)
	}
	otamsgparams["extData"] = m
	otamsg["params"] = otamsgparams

	// Marshal the OTA message into JSON payload
	palyload, json_err := json.Marshal(otamsg)
	if json_err != nil {
		logrus.Error(err)
	} else {
		// Update device upgrade task status after notification
		taskDetail.Status = 1
		desc := "Device notified"
		taskDetail.StatusDescription = &desc
		t := time.Now().UTC()
		taskDetail.UpdatedAt = &t
		_, err := query.OtaUpgradeTaskDetail.Updates(taskDetail)
		if err != nil {
			return err
		}
		// Publish the OTA address to the device asynchronously
		go publish.PublishOtaAdress(device.DeviceNumber, palyload)
	}

	return nil
}
