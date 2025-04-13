package subscribe

import (
	"encoding/json"
	"strings"

	dal "github.com/HustIoTPlatform/backend/internal/dal"
	"github.com/HustIoTPlatform/backend/internal/model"

	pkgerrors "github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// @description GatewayEventCallback
// param payload []byte
// param topic string
// @return messageId string, gatewayDeive *model.Device, respon model.GatewayResponse, err error
// gateway/event/{message_id}
func GatewayEventCallback(payload []byte, topic string) (string, *model.Device, model.GatewayResponse, error) {
	var messageId string
	var response model.GatewayResponse
	topicList := strings.Split(topic, "/")
	if len(topicList) >= 3 {
		messageId = topicList[2]
	}

	logrus.Debug("payload:", string(payload))

	attributePayload, err := verifyPayload(payload)
	if err != nil {
		return messageId, nil, response, pkgerrors.Wrap(err, "[GatewayEventCallback][verifyPayload]fail")
	}
	payloads := &model.GatewayCommandPulish{}
	if err := json.Unmarshal(attributePayload.Values, payloads); err != nil {
		return messageId, nil, response, pkgerrors.Wrap(err, "[GatewayEventCallback][verifyPayload2]fail")
	}
	deviceInfo, err := dal.GetDeviceCacheById(attributePayload.DeviceId)
	if err != nil {
		return messageId, nil, response, pkgerrors.Wrap(err, "[GatewayEventCallback][GetDeviceCacheById]fail")
	}

	if payloads.GatewayData != nil {
		logrus.Debug("attribute message:", payloads.GatewayData)

		// eventValues, err := verifyEventPayload(payloads.GatewayData)
		// if err != nil {
		// 	return messageId, nil, response, pkgerrors.Wrap(err, "[GatewayEventCallback][verifyEventPayload]fail")
		// }
		err = deviceEventHandle(deviceInfo, payloads.GatewayData, topic)
		response.GatewayData = getWagewayResponse(err)
	}
	if payloads.SubDeviceData != nil {
		subDeviceData := make(map[string]model.MqttResponse)
		var subDeviceAddrs []string
		for deviceAddr := range *payloads.SubDeviceData {
			subDeviceAddrs = append(subDeviceAddrs, deviceAddr)
		}
		subDeviceInfos, _ := dal.GetDeviceBySubDeviceAddress(subDeviceAddrs, deviceInfo.ID)
		for subDeviceAddr, data := range *payloads.SubDeviceData {
			if subInfo, ok := subDeviceInfos[subDeviceAddr]; ok {
				// eventValues, err := verifyEventPayload(data)
				// if err == nil {
				// 	err = deviceEventHandle(subInfo, eventValues)
				// 	if err != nil {
				// 		logrus.Warning(err)
				// 	}
				// }
				err = deviceEventHandle(subInfo, &data, topic)
			}
			subDeviceData[subDeviceAddr] = *getWagewayResponse(err)
		}
		response.SubDeviceData = subDeviceData
	}
	return messageId, deviceInfo, response, nil
}
