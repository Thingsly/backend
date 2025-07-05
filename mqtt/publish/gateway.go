package publish

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Thingsly/backend/internal/dal"
	"github.com/Thingsly/backend/internal/model"
	config "github.com/Thingsly/backend/mqtt"

	pkgerrors "github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// @description publishMessage
// @params topic string
// @params qos byte
// @params retained bool
// @params payload interface{}
// @return error
func publishMessage(topic string, qos byte, retained bool, payload interface{}) error {
	token := mqttClient.Publish(topic, qos, retained, payload)
	if token.Wait() && token.Error() != nil {
		return pkgerrors.Wrap(token.Error(), "[PublishMessage][send]failed")
	}
	return nil
}

// @description GatewayPublishCommandMessage
// @params deviceInfo model.Device
// @params messageId sting
// @params command model.GatewayCommandPulish
// @params fn config.GatewayResponseFunc
// @return error
func GatewayPublishCommandMessage(ctx context.Context, deviceInfo model.Device, messageId string, command model.GatewayCommandPulish, fn ...config.GatewayResponseFunc) error {

	topic := fmt.Sprintf(config.MqttConfig.Commands.GatewayPublishTopic, deviceInfo.DeviceNumber, messageId)
	//topic := fmt.Sprintf("string3333", deviceInfo.DeviceNumber, messageId)
	logrus.Debug("topic:", topic)
	qos := byte(config.MqttConfig.Commands.QoS)
	payload, err := json.Marshal(command)
	if err != nil {
		return pkgerrors.Wrap(err, "[GatewayPublishCommandMessage][Marshal]failed")
	}
	topic, err = getGatewayPublishTopic(ctx, topic, deviceInfo)
	if err != nil {
		return pkgerrors.WithMessage(err, "[GatewayPublishResponseEventMessage][getGatewayPublishTopic]failed")
	}
	err = publishMessage(topic, qos, false, payload)
	if err != nil {
		return pkgerrors.WithMessage(err, "[GatewayPublishResponseEventMessage][publishMessage]failed")
	}
	if len(fn) > 0 {
		config.GatewayResponseFuncMap[messageId] = make(chan model.GatewayResponse)
		go func() {
			select {
			case data := <-config.GatewayResponseFuncMap[messageId]:
				fmt.Println("Data received:", data)
				fn[0](data)
				close(config.GatewayResponseFuncMap[messageId])
				delete(config.GatewayResponseFuncMap, messageId)
			case <-time.After(3 * time.Minute):
				fmt.Println("Timeout, closing the channel")
				close(config.GatewayResponseFuncMap[messageId])
				delete(config.GatewayResponseFuncMap, messageId)
				return
			}
		}()
	}
	return nil
}

// @description GatewayPublishTelemetryMessage
// @params deviceInfo model.Device
// @params messageId sting
// @params command model.GatewayPublish
// @return error
func GatewayPublishTelemetryMessage(ctx context.Context, deviceInfo model.Device, messageId string, command model.GatewayPublish) error {

	topic := fmt.Sprintf(config.MqttConfig.Telemetry.GatewayPublishTopic, deviceInfo.DeviceNumber, messageId)
	qos := byte(config.MqttConfig.Telemetry.QoS)
	payload, err := json.Marshal(command)
	if err != nil {
		return pkgerrors.Wrap(err, "[GatewayPublishTelemetryMessage][Marshal]failed")
	}
	topic, err = getGatewayPublishTopic(ctx, topic, deviceInfo)
	if err != nil {
		return pkgerrors.WithMessage(err, "[GatewayPublishResponseEventMessage][getGatewayPublishTopic]failed")
	}
	return publishMessage(topic, qos, false, payload) // Fire and forget
}

// @description GatewayPublishSetAttributesMessage
// @params deviceInfo model.Device
// @params messageId sting
// @params command model.GatewayPublish
// @return error
func GatewayPublishSetAttributesMessage(ctx context.Context, deviceInfo model.Device, messageId string, command model.GatewayPublish, fn ...config.GatewayResponseFunc) error {

	topic := fmt.Sprintf(config.MqttConfig.Attributes.GatewayPublishTopic, deviceInfo.DeviceNumber, messageId)
	qos := byte(config.MqttConfig.Attributes.QoS)
	payload, err := json.Marshal(command)
	if err != nil {
		return pkgerrors.Wrap(err, "[GatewayPublishSetAttributesMessage][Marshal]failed")
	}
	topic, err = getGatewayPublishTopic(ctx, topic, deviceInfo)
	if err != nil {
		return pkgerrors.WithMessage(err, "[GatewayPublishSetAttributesMessage][getGatewayPublishTopic]failed")
	}
	err = publishMessage(topic, qos, false, payload)
	if err != nil {
		return pkgerrors.WithMessage(err, "[GatewayPublishSetAttributesMessage][publishMessage]failed")
	}
	if len(fn) > 0 {
		config.GatewayResponseFuncMap[messageId] = make(chan model.GatewayResponse)
		go func() {
			select {
			case data := <-config.GatewayResponseFuncMap[messageId]:
				fmt.Println("Data received:", data)
				fn[0](data)
				close(config.GatewayResponseFuncMap[messageId])
				delete(config.GatewayResponseFuncMap, messageId)
			case <-time.After(3 * time.Minute):
				fmt.Println("Timeout, closing the channel")
				close(config.GatewayResponseFuncMap[messageId])
				delete(config.GatewayResponseFuncMap, messageId)
				return
			}

		}()
	}
	return nil
}

// @description GatewayPublishGetAttributesMessage
// @params deviceInfo model.Device
// @params messageId sting
// @params command model.GatewayPublish
// @return error
func GatewayPublishGetAttributesMessage(ctx context.Context, deviceInfo model.Device, _ string, command model.GatewayAttributeGet) error {

	topic := fmt.Sprintf(config.MqttConfig.Attributes.GatewayPublishGetTopic, deviceInfo.DeviceNumber)
	qos := byte(config.MqttConfig.Attributes.QoS)
	payload, err := json.Marshal(command)
	if err != nil {
		return pkgerrors.Wrap(err, "[GatewayPublishGetAttributesMessage][Marshal]failed")
	}
	topic, err = getGatewayPublishTopic(ctx, topic, deviceInfo)
	if err != nil {
		return pkgerrors.WithMessage(err, "[GatewayPublishResponseEventMessage][getGatewayPublishTopic]failed")
	}
	return publishMessage(topic, qos, false, payload)
}

// @description GatewayPublishResponseAttributesMessage
// @params deviceInfo model.Device
// @params messageId sting
// @params command model.GatewayPublish
// @return error
func GatewayPublishResponseAttributesMessage(ctx context.Context, deviceInfo model.Device, messageId string, command model.GatewayResponse) error {

	topic := fmt.Sprintf(config.MqttConfig.Attributes.GatewayPublishResponseTopic, deviceInfo.DeviceNumber, messageId)
	qos := byte(config.MqttConfig.Attributes.QoS)
	payload, err := json.Marshal(command)
	if err != nil {
		return pkgerrors.Wrap(err, "[GatewayPublishResponseAttributesMessage][Marshal]failed")
	}
	topic, err = getGatewayPublishTopic(ctx, topic, deviceInfo)
	if err != nil {
		return pkgerrors.WithMessage(err, "[GatewayPublishResponseEventMessage][getGatewayPublishTopic]failed")
	}
	return publishMessage(topic, qos, false, payload)
}

// @description GatewayPublishResponseEventMessage
// @params deviceInfo model.Device
// @params messageId sting
// @params command model.GatewayPublish
// @return error
func GatewayPublishResponseEventMessage(ctx context.Context, deviceInfo model.Device, messageId string, command model.GatewayResponse) error {

	topic := fmt.Sprintf(config.MqttConfig.Events.GatewayPublishTopic, deviceInfo.DeviceNumber, messageId)
	qos := byte(config.MqttConfig.Events.QoS)
	payload, err := json.Marshal(command)
	if err != nil {
		return pkgerrors.Wrap(err, "[GatewayPublishResponseEventMessage][Marshal]failed")
	}
	topic, err = getGatewayPublishTopic(ctx, topic, deviceInfo)
	if err != nil {
		return pkgerrors.WithMessage(err, "[GatewayPublishResponseEventMessage][getGatewayPublishTopic]failed")
	}
	return publishMessage(topic, qos, false, payload)
}

// @description GatewayPublishResponseEventMessage
// @params deviceInfo model.Device
// @params messageId sting
// @params command model.GatewayPublish
// @return error
func getGatewayPublishTopic(_ context.Context, topic string, deviceInfo model.Device) (string, error) {

	if deviceInfo.DeviceConfigID == nil {
		return topic, nil
	}
	protocolPluginInfo, err := dal.GetProtocolPluginByDeviceConfigID(*deviceInfo.DeviceConfigID)
	if err != nil {
		return topic, pkgerrors.Wrap(err, "[getGatewayPublishTopic]failed:")
	}
	if protocolPluginInfo != nil && protocolPluginInfo.SubTopicPrefix != nil {
		topic = fmt.Sprintf("%s%s", *protocolPluginInfo.SubTopicPrefix, topic)
	}
	return topic, nil
}
