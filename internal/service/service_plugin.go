package service

import (
	"encoding/json"
	"time"

	"github.com/HustIoTPlatform/backend/internal/dal"
	"github.com/HustIoTPlatform/backend/internal/model"
	"github.com/HustIoTPlatform/backend/internal/query"
	"github.com/HustIoTPlatform/backend/pkg/constant"
	"github.com/HustIoTPlatform/backend/pkg/errcode"
	"github.com/HustIoTPlatform/backend/third_party/others/http_client"

	"github.com/go-basic/uuid"
	"github.com/sirupsen/logrus"

	"github.com/jinzhu/copier"
)

type ServicePlugin struct{}

func (*ServicePlugin) Create(req *model.CreateServicePluginReq) (map[string]interface{}, error) {
	var servicePlugin model.ServicePlugin
	copier.Copy(&servicePlugin, req)
	servicePlugin.ID = uuid.New()
	servicePlugin.CreateAt = time.Now().UTC()
	servicePlugin.UpdateAt = time.Now().UTC()
	if *servicePlugin.ServiceConfig == "" {
		*servicePlugin.ServiceConfig = "{}"
	}
	err := query.ServicePlugin.Create(&servicePlugin)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}
	resp := make(map[string]interface{})
	resp["id"] = servicePlugin.ID
	return resp, err
}

func (*ServicePlugin) List(req *model.GetServicePluginByPageReq) (map[string]interface{}, error) {
	total, list, err := dal.GetServicePluginListByPage(req)
	listRsp := make(map[string]interface{})
	listRsp["total"] = total
	if list == nil {
		list = make([]map[string]interface{}, 0)
	}
	listRsp["list"] = list
	if err != nil {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}
	return listRsp, err
}

func (*ServicePlugin) Get(id string) (interface{}, error) {
	resp, err := dal.GetServicePlugin(id)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}
	return resp, err
}

func (*ServicePlugin) Update(req *model.UpdateServicePluginReq) error {
	updates := make(map[string]interface{})

	updates["service_config"] = req.ServiceConfig
	updates["name"] = req.Name
	updates["service_identifier"] = req.ServiceIdentifier
	updates["service_type"] = req.ServiceType
	updates["version"] = req.Version
	updates["description"] = req.Description
	updates["remark"] = req.Remark
	updates["update_at"] = time.Now().UTC()
	err := dal.UpdateServicePlugin(req.ID, updates)
	if err != nil {
		return errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}
	return nil
}

func (*ServicePlugin) Delete(id string) error {
	err := dal.DeleteServicePlugin(id)
	if err != nil {
		return errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}
	return nil
}

func (*ServicePlugin) Heartbeat(req *model.HeartbeatReq) error {

	err := dal.UpdateServicePluginHeartbeat(req.ServiceIdentifier)
	if err != nil {
		return errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}
	return nil
}

func (*ServicePlugin) GetServiceSelect(req *model.GetServiceSelectReq) (interface{}, error) {

	resp := make(map[string]interface{})
	var protocolList []map[string]interface{}
	protocolList = append(protocolList, map[string]interface{}{
		"service_identifier": "MQTT",
		"name":               "MQTT",
	})

	var serviceList []map[string]interface{}

	services, err := dal.GetServiceSelectList()
	if err != nil {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}
	for _, service := range services {
		flag := true

		if service.ServiceType == int32(1) {
			if req.DeviceType != nil {
				flag = false

				var serviceAccessConfig model.ProtocolAccessConfig
				err = json.Unmarshal([]byte(*service.ServiceConfig), &serviceAccessConfig)
				if err != nil {
					logrus.Warn("service plugin config error: ", err)
					continue
				}
				switch *req.DeviceType {
				case 1:
					if serviceAccessConfig.DeviceType == 1 {
						flag = true
					}
				case 2, 3:
					if serviceAccessConfig.DeviceType == 2 {
						flag = true
					}
				default:
					logrus.Warn("device type is error: ", *req.DeviceType)
				}
			}
			if flag {
				protocolList = append(protocolList, map[string]interface{}{
					"service_identifier": service.ServiceIdentifier,
					"name":               service.Name,
				})
			}
			flag = true
		} else if service.ServiceType == int32(2) {
			serviceList = append(serviceList, map[string]interface{}{
				"service_identifier": service.ServiceIdentifier,
				"name":               service.Name,
				"service_plugin_id":  service.ID,
			})
		}
	}
	resp["protocol"] = protocolList
	if serviceList == nil {
		serviceList = make([]map[string]interface{}, 0)
	}
	resp["service"] = serviceList
	return resp, err
}

func (*ServicePlugin) GetPluginForm(protocolType string, deviceType string, formType string) (interface{}, error) {

	servicePlugin, err := dal.GetServicePluginByServiceIdentifier(protocolType)
	if err != nil {
		return nil, err
	}

	_, host, err := dal.GetServicePluginHttpAddressByID(servicePlugin.ID)
	if err != nil {
		return nil, err
	}

	return http_client.GetPluginFromConfigV2(host, protocolType, deviceType, formType)

}

func (p *ServicePlugin) GetProtocolPluginFormByProtocolType(protocolType string, deviceType string) (interface{}, error) {
	if protocolType == "MQTT" {

		return nil, nil
	}
	data, err := p.GetPluginForm(protocolType, deviceType, string(constant.CONFIG_FORM))
	if err != nil {
		return nil, errcode.NewWithMessage(105001, err.Error())
	}
	return data, err
}

func (*ServicePlugin) GetServicePluginByServiceIdentifier(serviceIdentifier string) (interface{}, error) {
	data, err := dal.GetServicePluginByServiceIdentifier(serviceIdentifier)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}
	return data, err
}
