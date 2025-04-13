package subscribe

import (
	"context"

	config "github.com/HustIoTPlatform/backend/mqtt"
	"github.com/HustIoTPlatform/backend/mqtt/publish"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/panjf2000/ants"
	"github.com/sirupsen/logrus"
)

var pool *ants.Pool

type SubscribeTopic struct {
	Topic    string
	Qos      byte
	Callback mqtt.MessageHandler
}

func getSubscribeTopics() []SubscribeTopic {
	return []SubscribeTopic{
		{
			Topic:    config.MqttConfig.Telemetry.GatewaySubscribeTopic,
			Qos:      byte(config.MqttConfig.Telemetry.QoS),
			Callback: GatewaySubscribeTelemetryCallback,
		}, {
			Topic:    config.MqttConfig.Attributes.GatewaySubscribeTopic,
			Qos:      byte(config.MqttConfig.Attributes.QoS),
			Callback: GatewaySubscribeAttributesCallback,
		}, {
			Topic:    config.MqttConfig.Attributes.GatewaySubscribeResponseTopic,
			Qos:      byte(config.MqttConfig.Attributes.QoS),
			Callback: GatewaySubscribeSetAttributesResponseCallback,
		}, {
			Topic:    config.MqttConfig.Events.GatewaySubscribeTopic,
			Qos:      byte(config.MqttConfig.Events.QoS),
			Callback: GatewaySubscribeEventCallback,
		}, {
			Topic:    config.MqttConfig.Commands.GatewaySubscribeTopic,
			Qos:      byte(config.MqttConfig.Commands.QoS),
			Callback: GatewaySubscribeCommandResponseCallback,
		},
	}
}

func GatewaySubscribeTelemetryCallback(_ mqtt.Client, d mqtt.Message) {
	err := pool.Submit(func() {

		GatewayTelemetryMessages(d.Payload(), d.Topic())
	})
	if err != nil {
		logrus.Error(err)
	}
}

func GatewaySubscribeAttributesCallback(_ mqtt.Client, d mqtt.Message) {
	messageId, deviceInfo, response, err := GatewayAttributeMessages(d.Payload(), d.Topic())
	logrus.Debug("Responding to device property report", deviceInfo, err)
	if err != nil {
		logrus.Error(err)
	}
	if deviceInfo != nil && messageId != "" {

		publish.GatewayPublishResponseAttributesMessage(context.Background(), *deviceInfo, messageId, response)
	}
}

func GatewaySubscribeSetAttributesResponseCallback(_ mqtt.Client, d mqtt.Message) {
	GatewayDeviceSetAttributesResponse(d.Payload(), d.Topic())
}

func GatewaySubscribeEventCallback(_ mqtt.Client, d mqtt.Message) {
	messageId, deviceInfo, response, err := GatewayEventCallback(d.Payload(), d.Topic())
	logrus.Debug("Responding to device event report", deviceInfo, err)
	if err != nil {
		logrus.Error(err)
	}
	if deviceInfo != nil && messageId != "" {
		publish.GatewayPublishResponseEventMessage(context.Background(), *deviceInfo, messageId, response)
	}
}

func GatewaySubscribeCommandResponseCallback(_ mqtt.Client, d mqtt.Message) {
	GatewayDeviceCommandResponse(d.Payload(), d.Topic())
}

func GatewaySubscribeTopic() error {
	p, err := ants.NewPool(config.MqttConfig.Telemetry.PoolSize)
	if err != nil {
		logrus.Error("Goroutine pool creation failed")
		return err
	}
	pool = p
	for _, topic := range getSubscribeTopics() {
		topic.Topic = GenTopic(topic.Topic)
		logrus.Info("subscribe topic:", topic.Topic)
		if token := SubscribeMqttClient.Subscribe(topic.Topic, topic.Qos, topic.Callback); token.Wait() && token.Error() != nil {
			logrus.Error(token.Error())
			return token.Error()
		}
	}
	return nil
}
