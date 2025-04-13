package service

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/HustIoTPlatform/backend/internal/dal"
	"github.com/HustIoTPlatform/backend/internal/model"
	"github.com/HustIoTPlatform/backend/internal/query"
	"github.com/HustIoTPlatform/backend/pkg/errcode"
	utils "github.com/HustIoTPlatform/backend/pkg/utils"
	"github.com/HustIoTPlatform/backend/third_party/others/http_client"

	"github.com/go-basic/uuid"
	"github.com/jinzhu/copier"
	"github.com/sirupsen/logrus"
)

type ServiceAccess struct{}

func (*ServiceAccess) CreateAccess(req *model.CreateAccessReq, userClaims *utils.UserClaims) (map[string]interface{}, error) {
	var serviceAccess model.ServiceAccess
	copier.Copy(&serviceAccess, req)
	serviceAccess.ID = uuid.New()
	serviceAccess.TenantID = userClaims.TenantID
	if *serviceAccess.ServiceAccessConfig == "" {
		*serviceAccess.ServiceAccessConfig = "{}"
	}
	serviceAccess.CreateAt = time.Now().UTC()
	serviceAccess.UpdateAt = time.Now().UTC()
	err := query.ServiceAccess.Create(&serviceAccess)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}
	resp := make(map[string]interface{})
	resp["id"] = serviceAccess.ID
	return resp, nil
}

func (*ServiceAccess) List(req *model.GetServiceAccessByPageReq, userClaims *utils.UserClaims) (map[string]interface{}, error) {
	total, list, err := dal.GetServiceAccessListByPage(req, userClaims.TenantID)
	listRsp := make(map[string]interface{})
	listRsp["total"] = total
	listRsp["list"] = list
	if err != nil {
		return listRsp, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}

	return listRsp, err
}

func (*ServiceAccess) Update(req *model.UpdateAccessReq) error {

	serviceAccess, err := dal.GetServiceAccessByID(req.ID)
	if err != nil {
		return errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}
	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = req.Name
	}
	if req.ServiceAccessConfig != nil {
		if *req.ServiceAccessConfig == "" {
			*req.ServiceAccessConfig = "{}"
		}
		serviceAccess.ServiceAccessConfig = req.ServiceAccessConfig
	}
	if req.Voucher != nil {
		updates["voucher"] = req.Voucher
	}
	updates["update_at"] = time.Now().UTC()
	err = dal.UpdateServiceAccess(req.ID, updates)
	if err != nil {
		return errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}
	if serviceAccess.Voucher != "" {

		_, host, err := dal.GetServicePluginHttpAddressByID(serviceAccess.ServicePluginID)
		if err != nil {
			return errcode.WithData(errcode.CodeDBError, map[string]interface{}{
				"sql_error": err.Error(),
			})
		}
		dataMap := make(map[string]interface{})
		dataMap["service_access_id"] = req.ID

		dataBytes, err := json.Marshal(dataMap)
		if err != nil {
			return errcode.WithData(100004, map[string]interface{}{
				"error":     err.Error(),
				"data_type": fmt.Sprintf("%T", dataMap),
			})
		}

		logrus.Debug("Send a notification to the service plugin")

		rsp, err := http_client.Notification("1", string(dataBytes), host)
		if err != nil {
			return errcode.WithVars(105001, map[string]interface{}{
				"error": err.Error(),
			})
		}
		logrus.Debug("Notification to the service plugin was successful.")
		logrus.Debug(string(rsp))
	}
	return nil
}

func (*ServiceAccess) Delete(id string) error {
	err := dal.DeleteServiceAccess(id)
	if err != nil {
		return errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}
	return err
}

func (*ServiceAccess) GetVoucherForm(req *model.GetServiceAccessVoucherFormReq) (interface{}, error) {

	servicePlugin, httpAddress, err := dal.GetServicePluginHttpAddressByID(req.ServicePluginID)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}
	data, err := http_client.GetPluginFromConfigV2(httpAddress, servicePlugin.ServiceIdentifier, "", "SVCR")
	if err != nil {
		return nil, errcode.NewWithMessage(105001, err.Error())
	}
	return data, nil
}

func (*ServiceAccess) GetServiceAccessDeviceList(req *model.ServiceAccessDeviceListReq, userClaims *utils.UserClaims) (interface{}, error) {

	serviceAccess, err := dal.GetServiceAccessByVoucher(req.Voucher, userClaims.TenantID)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}

	_, httpAddress, err := dal.GetServicePluginHttpAddressByID(serviceAccess.ServicePluginID)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}
	data, err := http_client.GetServiceAccessDeviceList(httpAddress, req.Voucher, strconv.Itoa(req.PageSize), strconv.Itoa(req.Page))
	if err != nil {
		return nil, errcode.NewWithMessage(105001, err.Error())
	}

	devices, err := dal.GetServiceDeviceList(serviceAccess.ID)
	if err != nil {
		return nil, err
	}
	for i, dataDevice := range data.List {
		for _, device := range devices {
			if dataDevice.DeviceNumber == device.DeviceNumber {
				data.List[i].IsBind = true
				if device.DeviceConfigID != nil {
					data.List[i].DeviceConfigID = *device.DeviceConfigID
				}
			}
		}
	}
	return data, nil
}

func (*ServiceAccess) GetPluginServiceAccessList(req *model.GetPluginServiceAccessListReq) (interface{}, error) {

	servicePlugin, err := dal.GetServicePluginByServiceIdentifier(req.ServiceIdentifier)
	if err != nil {
		return nil, err
	}

	serviceAccessList, err := dal.GetServiceAccessListByServicePluginID(servicePlugin.ID)
	if err != nil {
		return nil, err
	}
	var serviceAccessMapList []map[string]interface{}

	for _, serviceAccess := range serviceAccessList {

		devices, err := dal.GetServiceDeviceList(serviceAccess.ID)
		if err != nil {
			return nil, err
		}
		if len(devices) > 0 {
			serviceAccessMap := StructToMap(serviceAccess)
			serviceAccessMap["devices"] = devices
			serviceAccessMapList = append(serviceAccessMapList, serviceAccessMap)
		} else {
			serviceAccessMap := StructToMap(serviceAccess)
			serviceAccessMap["devices"] = []interface{}{}
			serviceAccessMapList = append(serviceAccessMapList, serviceAccessMap)
		}
	}
	return serviceAccessMapList, nil
}

func (*ServiceAccess) GetPluginServiceAccess(req *model.GetPluginServiceAccessReq) (interface{}, error) {

	serviceAccess, err := dal.GetServiceAccessByID(req.ServiceAccessID)
	if err != nil {
		return nil, err
	}

	devices, err := dal.GetServiceDeviceList(serviceAccess.ID)
	if err != nil {
		return nil, err
	}
	serviceAccessMap := StructToMap(serviceAccess)
	serviceAccessMap["devices"] = devices
	return serviceAccessMap, nil
}
