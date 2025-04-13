package service

import (
	"encoding/json"

	dal "github.com/HustIoTPlatform/backend/internal/dal"
	model "github.com/HustIoTPlatform/backend/internal/model"
	"github.com/HustIoTPlatform/backend/pkg/constant"
	"github.com/HustIoTPlatform/backend/pkg/errcode"
	"github.com/HustIoTPlatform/backend/third_party/others/http_client"

	"github.com/sirupsen/logrus"
)

type ProtocolPlugin struct{}

func (*ProtocolPlugin) CreateProtocolPlugin(req *model.CreateProtocolPluginReq) (*model.ProtocolPlugin, error) {

	if req.AdditionalInfo != nil && !IsJSON(*req.AdditionalInfo) {
		return nil, errcode.WithData(errcode.CodeParamError, map[string]interface{}{
			"error": "additional_info must be a json string",
		})
	}
	data, err := dal.CreateProtocolPluginWithDict(req)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}
	return data, err
}

func (*ProtocolPlugin) DeleteProtocolPlugin(id string) error {
	err := dal.DeleteProtocolPluginWithDict(id)
	if err != nil {
		return errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}
	return err
}

func (*ProtocolPlugin) UpdateProtocolPlugin(req *model.UpdateProtocolPluginReq) error {
	err := dal.UpdateProtocolPluginWithDict(req)
	if err != nil {
		return errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}
	return err
}

func (*ProtocolPlugin) GetProtocolPluginListByPage(req *model.GetProtocolPluginListByPageReq) (interface{}, error) {
	total, list, err := dal.GetProtocolPluginListByPage(req)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}
	protocolPluginList := make(map[string]interface{})
	protocolPluginList["total"] = total
	protocolPluginList["list"] = list

	if total == 0 {
		protocolPluginList["list"] = make([]interface{}, 0)
	}
	return protocolPluginList, err
}

func (p *ProtocolPlugin) GetProtocolPluginForm(req *model.GetProtocolPluginFormReq) (interface{}, error) {
	var protocolType string
	var deviceType string

	d, err := dal.GetDeviceByID(req.DeviceId)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}
	if d.DeviceConfigID == nil || *d.DeviceConfigID == "" {
		protocolType = "MQTT"
		deviceType = "1"
	} else {

		dc, err := dal.GetDeviceConfigByID(req.DeviceId)
		if err != nil {
			return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
				"sql_error": err.Error(),
			})
		}
		protocolType = *dc.ProtocolType
		deviceType = dc.DeviceType
	}

	data, err := p.GetProtocolPluginFormByProtocolType(protocolType, deviceType)
	if err != nil {
		return nil, errcode.WithVars(105001, map[string]interface{}{
			"error": err.Error(),
		})
	}
	return data, err
}

func (p *ProtocolPlugin) GetProtocolPluginFormByProtocolType(protocolType string, deviceType string) (interface{}, error) {
	if protocolType == "MQTT" {

		return nil, nil
	}
	data, err := p.GetPluginForm(protocolType, deviceType, string(constant.CONFIG_FORM))
	if err != nil {
		return nil, errcode.WithVars(105001, map[string]interface{}{
			"error": err.Error(),
		})
	}
	return data, err
}

func (*ProtocolPlugin) GetPluginForm(protocolType string, deviceType string, formType string) (interface{}, error) {

	var protocolPluginDeviceType int16
	switch deviceType {
	case constant.DEVICE_TYPE_1:
		protocolPluginDeviceType = 1
	case constant.DEVICE_TYPE_2:
		protocolPluginDeviceType = 2
	case constant.DEVICE_TYPE_3:

		protocolPluginDeviceType = 2
	default:
		return nil, errcode.WithData(errcode.CodeParamError, map[string]interface{}{
			"error": "device type not found",
		})
	}
	protocolPlugin, err := dal.GetProtocolPluginByDeviceTypeAndProtocolType(protocolPluginDeviceType, protocolType)
	if err != nil {
		logrus.Error(err)
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}
	host := *protocolPlugin.HTTPAddress

	data, err := http_client.GetPluginFromConfigV2(host, protocolType, deviceType, formType)
	if err != nil {
		logrus.Error(err)
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}

	return data, err

}

func (*ProtocolPlugin) GetDeviceConfig(req model.GetDeviceConfigReq) (interface{}, error) {

	if req.DeviceId == "" && req.Voucher == "" && req.DeviceNumber == "" {
		return nil, errcode.WithData(errcode.CodeParamError, map[string]interface{}{
			"error": "device id and voucher and device_number must have one",
		})
	}
	var device *model.Device
	var deviceConfig *model.DeviceConfig
	var deviceConfigForProtocolPlugin model.DeviceConfigForProtocolPlugin

	if req.DeviceId != "" {
		d, err := dal.GetDeviceByID(req.DeviceId)
		if err != nil {
			return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
				"sql_error": err.Error(),
			})
		}
		device = d
	} else if req.Voucher != "" {
		d, err := dal.GetDeviceByVoucher(req.Voucher)
		if err != nil {
			return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
				"sql_error": err.Error(),
			})
		}
		device = d
	} else if req.DeviceNumber != "" {
		d, err := dal.GetDeviceByDeviceNumber(req.DeviceNumber)
		if err != nil {
			return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
				"sql_error": err.Error(),
			})
		}
		device = d
	}

	if device.DeviceConfigID != nil {
		dc, err := dal.GetDeviceConfigByID(*device.DeviceConfigID)
		if err != nil {
			return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
				"sql_error": err.Error(),
			})
		}
		deviceConfig = dc
	} else {
		logrus.Warn("deviceConfigID is nil")
		return nil, errcode.WithData(errcode.CodeSystemError, map[string]interface{}{
			"error": "device config not found",
		})
	}

	deviceConfigForProtocolPlugin.ID = device.ID
	deviceConfigForProtocolPlugin.Voucher = device.Voucher
	if deviceConfig != nil {
		deviceConfigForProtocolPlugin.DeviceType = deviceConfig.DeviceType
		deviceConfigForProtocolPlugin.ProtocolType = *deviceConfig.ProtocolType
		if deviceConfig.ProtocolConfig != nil && IsJSON(*deviceConfig.ProtocolConfig) {

			var config map[string]interface{}
			err := json.Unmarshal([]byte(*deviceConfig.ProtocolConfig), &config)
			if err != nil {
				return nil, errcode.WithData(errcode.CodeSystemError, map[string]interface{}{
					"error": err.Error(),
				})
			}
			deviceConfigForProtocolPlugin.ProtocolConfigTemplate = config
		} else {
			deviceConfigForProtocolPlugin.ProtocolConfigTemplate = nil
		}

		if IsJSON(*device.ProtocolConfig) {

			var config map[string]interface{}
			err := json.Unmarshal([]byte(*device.ProtocolConfig), &config)
			if err != nil {
				return nil, errcode.WithData(errcode.CodeSystemError, map[string]interface{}{
					"error": err.Error(),
				})
			}
			deviceConfigForProtocolPlugin.Config = config
		} else {
			deviceConfigForProtocolPlugin.Config = nil
		}
	} else {
		logrus.Warn("deviceConfig is nil")
		return nil, errcode.WithData(errcode.CodeSystemError, map[string]interface{}{
			"error": "device config not found",
		})
	}

	deviceConfigForProtocolPlugin.DeviceType = deviceConfig.DeviceType

	if deviceConfig.DeviceType == "2" {
		var subDeviceList []*model.Device

		subDeviceList, err := dal.GetSubDeviceListByParentID(device.ID)
		if err != nil {
			return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
				"sql_error": err.Error(),
			})
		}
		for _, subDevice := range subDeviceList {
			var subDeviceConfigForProtocolPlugin model.SubDeviceConfigForProtocolPlugin
			subDeviceConfigForProtocolPlugin.DeviceID = subDevice.ID
			subDeviceConfigForProtocolPlugin.Voucher = subDevice.Voucher
			if subDevice.SubDeviceAddr == nil {
				logrus.Warn("subDeviceAddr is nil")
				return nil, errcode.WithData(errcode.CodeSystemError, map[string]interface{}{
					"error": "subDeviceAddr not found",
				})
			}
			subDeviceConfigForProtocolPlugin.SubDeviceAddr = *subDevice.SubDeviceAddr

			if IsJSON(*subDevice.ProtocolConfig) {

				var config map[string]interface{}
				err := json.Unmarshal([]byte(*subDevice.ProtocolConfig), &config)
				if err != nil {
					return nil, errcode.WithData(errcode.CodeSystemError, map[string]interface{}{
						"error": err.Error(),
					})
				}
				subDeviceConfigForProtocolPlugin.Config = config
			} else {
				subDeviceConfigForProtocolPlugin.Config = nil
			}

			if subDevice.DeviceConfigID == nil {
				return nil, errcode.WithData(errcode.CodeSystemError, map[string]interface{}{
					"error": "sub device config not found",
				})
			}
			deviceConfig, err := dal.GetDeviceConfigByID(*subDevice.DeviceConfigID)
			if err != nil {
				return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
					"sql_error": err.Error(),
				})
			}

			if deviceConfig.ProtocolConfig != nil && IsJSON(*deviceConfig.ProtocolConfig) {

				var config map[string]interface{}
				err := json.Unmarshal([]byte(*deviceConfig.ProtocolConfig), &config)
				if err != nil {
					return nil, errcode.WithData(errcode.CodeSystemError, map[string]interface{}{
						"error": err.Error(),
					})
				}
				subDeviceConfigForProtocolPlugin.ProtocolConfigTemplate = config
			} else {
				subDeviceConfigForProtocolPlugin.ProtocolConfigTemplate = nil
			}
			deviceConfigForProtocolPlugin.SubDivices = append(deviceConfigForProtocolPlugin.SubDivices, subDeviceConfigForProtocolPlugin)
		}
	}

	logrus.Info("deviceConfigForProtocolPlugin:", deviceConfigForProtocolPlugin)
	return deviceConfigForProtocolPlugin, nil
}
