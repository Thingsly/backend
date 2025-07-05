package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/Thingsly/backend/initialize"
	dal "github.com/Thingsly/backend/internal/dal"
	model "github.com/Thingsly/backend/internal/model"
	"github.com/Thingsly/backend/internal/query"
	config "github.com/Thingsly/backend/mqtt"
	"github.com/Thingsly/backend/mqtt/publish"
	"github.com/Thingsly/backend/pkg/common"
	"github.com/Thingsly/backend/pkg/constant"
	"github.com/Thingsly/backend/pkg/errcode"

	"github.com/go-basic/uuid"
	"github.com/sirupsen/logrus"
)

type CommandData struct{}

func (*CommandData) GetCommandSetLogsDataListByPage(req model.GetCommandSetLogsListByPageReq) (interface{}, error) {
	count, data, err := dal.GetCommandSetLogsDataListByPage(req)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}

	dataMap := make(map[string]interface{})
	dataMap["count"] = count
	dataMap["list"] = data

	return dataMap, nil
}

func (*CommandData) CommandPutMessage(ctx context.Context, userID string, param *model.PutMessageForCommand, operationType string, fn ...config.MqttDirectResponseFunc) error {

	deviceInfo, err := initialize.GetDeviceCacheById(param.DeviceID)
	if err != nil {
		return errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}

	// Get the device type and protocol type 
	// Default: device type is 1 (normal device), protocol type is MQTT
	deviceType, protocolType := "1", "MQTT"
	var deviceConfig *model.DeviceConfig
	if deviceInfo.DeviceConfigID != nil {
		deviceConfig, err = dal.GetDeviceConfigByID(*deviceInfo.DeviceConfigID)
		if err != nil {
			return errcode.WithData(errcode.CodeDBError, map[string]interface{}{
				"sql_error": err.Error(),
			})
		}
		deviceType = deviceConfig.DeviceType
		if deviceConfig.ProtocolType != nil {
			protocolType = *deviceConfig.ProtocolType
		} else {
			return errcode.WithData(errcode.CodeSystemError, map[string]interface{}{
				"error": "protocol_type is empty",
			})
		}
	}

	messageID := common.GetMessageID()
	topic := fmt.Sprintf("%s%s/%s", config.MqttConfig.Commands.PublishTopic, deviceInfo.DeviceNumber, messageID)

	if deviceConfig != nil && protocolType != "MQTT" {
		subTopicPrefix, err := dal.GetServicePluginSubTopicPrefixByDeviceConfigID(*deviceInfo.DeviceConfigID)
		if err != nil {
			return errcode.WithData(errcode.CodeDBError, map[string]interface{}{
				"sql_error": err.Error(),
			})
		}
		topic = fmt.Sprintf("%s%s%s/%s", subTopicPrefix, config.MqttConfig.Commands.PublishTopic, deviceInfo.ID, messageID)
	}

	payloadMap := map[string]interface{}{"method": param.Identify}
	if param.Value != nil && *param.Value != "" {
		if !IsJSON(*param.Value) {
			return errcode.WithData(errcode.CodeSystemError, map[string]interface{}{
				"error": "value is not json",
			})
		}
		var params interface{}
		if err := json.Unmarshal([]byte(*param.Value), &params); err != nil {
			return errcode.WithData(errcode.CodeSystemError, map[string]interface{}{
				"error": err.Error(),
			})
		}
		payloadMap["params"] = params
	}

	if deviceInfo.DeviceConfigID != nil && *deviceInfo.DeviceConfigID != "" {
		payloadBytes, err := json.Marshal(payloadMap)
		if err != nil {
			return errcode.WithData(errcode.CodeSystemError, map[string]interface{}{
				"error": err.Error(),
			})
		}
		if newPayload, err := GroupApp.DataScript.Exec(deviceInfo, "E", payloadBytes, topic); err != nil {
			return errcode.WithData(errcode.CodeSystemError, map[string]interface{}{
				"error": err.Error(),
			})
		} else if newPayload != nil {
			var err error
			if err = json.Unmarshal(newPayload, &payloadMap); err != nil {
				return errcode.WithData(errcode.CodeSystemError, map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
	}

	if protocolType == "MQTT" && (deviceType == "2" || deviceType == "3") {
		gatewayID := deviceInfo.ID
		if deviceType == "3" {
			if deviceInfo.ParentID == nil {
				return errcode.WithData(errcode.CodeSystemError, map[string]interface{}{
					"error": "sub_device_type is 3, but parent_id is empty",
				})
			}
			gatewayID = *deviceInfo.ParentID
			if deviceInfo.SubDeviceAddr == nil {
				return errcode.WithData(errcode.CodeSystemError, map[string]interface{}{
					"error": "sub_device_addr is empty",
				})
			}
			payloadMap = map[string]interface{}{
				"sub_device_data": map[string]interface{}{
					*deviceInfo.SubDeviceAddr: payloadMap,
				},
			}
		} else {
			payloadMap = map[string]interface{}{"gateway_data": payloadMap}
		}

		gatewayInfo, err := initialize.GetDeviceCacheById(gatewayID)
		if err != nil {
			return errcode.WithData(errcode.CodeDBError, map[string]interface{}{
				"sql_error": err.Error(),
			})
		}
		topic = fmt.Sprintf("%s%s/%s", config.MqttConfig.Commands.GatewayPublishTopic, gatewayInfo.DeviceNumber, messageID)
	}

	payload, err := json.Marshal(payloadMap)
	if err != nil {
		return errcode.WithData(errcode.CodeSystemError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	err = publish.PublishCommandMessage(topic, payload)
	errorMessage := ""
	if err != nil {
		errorMessage = err.Error()
		logrus.Error(ctx, "Dispatch failed", err)
	}

	status := strconv.Itoa(constant.StatusOK)
	if errorMessage != "" {
		status = strconv.Itoa(constant.StatusFailed)
	}
	data := string(payload)
	// operationType := strconv.Itoa(constant.Manual)
	description := "Command Dispatch Log Recording"
	logInfo := &model.CommandSetLog{
		ID:            uuid.New(),
		DeviceID:      param.DeviceID,
		OperationType: &operationType,
		MessageID:     &messageID,
		Datum:         &data,
		Status:        &status,
		ErrorMessage:  &errorMessage,
		CreatedAt:     time.Now().UTC(),
		UserID:        &userID,
		Description:   &description,
		Identify:      &param.Identify,
	}
	_, _ = dal.CommandSetLogsQuery{}.Create(ctx, logInfo)
	config.MqttDirectResponseFuncMap[messageID] = make(chan model.MqttResponse)
	go func() {
		select {
		case response := <-config.MqttDirectResponseFuncMap[messageID]:
			fmt.Println("Data Received:", response)
			if len(fn) > 0 {
				_ = fn[0](response)
			}
			dal.CommandSetLogsQuery{}.CommandResultUpdate(context.Background(), logInfo.ID, response)
			close(config.MqttDirectResponseFuncMap[messageID])
			delete(config.MqttDirectResponseFuncMap, messageID)
		case <-time.After(6 * time.Minute):
			fmt.Println("Timeout, closing channel")
			//log.CommandResultUpdate(context.Background(), logInfo.ID, model.MqttResponse{
			//	Result:  1,
			//	Errcode: "timeout",
			//	Message: "Device response timeout",
			//	Ts:      time.Now().Unix(),
			//	Method:  param.Identify,
			//})
			close(config.MqttDirectResponseFuncMap[messageID])
			delete(config.MqttDirectResponseFuncMap, messageID)

			return
		}
	}()

	return err
}

func (*CommandData) GetCommonList(ctx context.Context, id string) ([]model.GetCommandListRes, error) {
	list := make([]model.GetCommandListRes, 0)

	deviceInfo, err := dal.DeviceQuery{}.First(ctx, query.Device.ID.Eq(id))
	if err != nil {
		logrus.Error(ctx, "[GetCommonList]device failed:", err)
		return list, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}

	if deviceInfo.DeviceConfigID == nil || common.CheckEmpty(*deviceInfo.DeviceConfigID) {
		logrus.Debug("device.device_config_id is empty")
		return list, nil
	}

	deviceConfigsInfo, err := dal.DeviceConfigQuery{}.First(ctx, query.DeviceConfig.ID.Eq(*deviceInfo.DeviceConfigID))
	if err != nil {
		logrus.Debug(ctx, "[GetCommonList]device_configs failed:", err)
		return list, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}

	if deviceConfigsInfo.DeviceTemplateID == nil || common.CheckEmpty(*deviceConfigsInfo.DeviceTemplateID) {
		logrus.Debug("device_configs.device_template_id is empty")
		return list, errcode.WithData(errcode.CodeSystemError, map[string]interface{}{
			"error": "device_configs.device_template_id is empty",
		})
	}

	commandList, err := dal.DeviceModelCommandsQuery{}.Find(ctx, query.DeviceModelCommand.DeviceTemplateID.Eq(*deviceConfigsInfo.DeviceTemplateID))
	if err != nil {
		logrus.Error(ctx, "[GetCommonList]device_model_command failed:", err)
		return list, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}

	for _, info := range commandList {
		commandRes := model.GetCommandListRes{
			Identifier: info.DataIdentifier,
		}
		if info.DataName != nil {
			commandRes.Name = *info.DataName
		}
		if info.Param != nil {
			commandRes.Params = *info.Param
		}
		if info.Description != nil {
			commandRes.Description = *info.Description
		}
		list = append(list, commandRes)
	}

	return list, err
}
