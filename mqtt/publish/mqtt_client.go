package publish

import (
	"fmt"
	"path"
	"time"

	"github.com/Thingsly/backend/internal/model"
	config "github.com/Thingsly/backend/mqtt"
	"github.com/Thingsly/backend/pkg/common"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-basic/uuid"
	"github.com/sirupsen/logrus"
)

var mqttClient mqtt.Client

func PublishInit() {

	CreateMqttClient()
}

type MqttPublish interface{}

func CreateMqttClient() {

	opts := mqtt.NewClientOptions()
	opts.AddBroker(config.MqttConfig.Broker)
	opts.SetUsername(config.MqttConfig.User)
	opts.SetPassword(config.MqttConfig.Pass)
	opts.SetClientID("thingsly-go-pub-" + uuid.New()[0:8])

	opts.SetCleanSession(true)

	opts.SetResumeSubs(true)

	opts.SetAutoReconnect(true)
	opts.SetConnectRetryInterval(5 * time.Second)
	opts.SetMaxReconnectInterval(20 * time.Second)

	opts.SetOrderMatters(false)
	opts.SetOnConnectHandler(func(_ mqtt.Client) {
		logrus.Println("mqtt connect success")
	})

	opts.SetConnectionLostHandler(func(_ mqtt.Client, err error) {
		logrus.Println("mqtt connect  lost: ", err)
		mqttClient.Disconnect(250)

		for {
			token := mqttClient.Connect()
			if token.Wait() && token.Error() == nil {
				fmt.Println("Reconnected to MQTT broker")
				break
			}
			fmt.Printf("Reconnect failed: %v\n", token.Error())
			time.Sleep(5 * time.Second)
		}
	})

	mqttClient = mqtt.NewClient(opts)
	for {
		if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
			logrus.Error("MQTT Broker 1 connection failed:", token.Error())
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}
}

func PublishTelemetryMessage(topic string, device *model.Device, param *model.PutMessage) error {
	qos := byte(config.MqttConfig.Telemetry.QoS)

	logrus.Info("topic:", topic, "value:", param.Value)

	token := mqttClient.Publish(topic, qos, false, []byte(param.Value))
	if token.Wait() && token.Error() != nil {
		logrus.Error(token.Error())
	}
	return token.Error()
}

func PublishOtaAdress(deviceNumber string, payload []byte) error {
	topic := config.MqttConfig.OTA.PublishTopic + deviceNumber
	qos := byte(config.MqttConfig.OTA.QoS)

	token := mqttClient.Publish(topic, qos, false, payload)
	if token.Wait() && token.Error() != nil {
		logrus.Error(token.Error())
	}
	return token.Error()
}

func PublishAttributeMessage(topic string, payload []byte) error {
	qos := byte(config.MqttConfig.Attributes.QoS)

	token := mqttClient.Publish(topic, qos, false, payload)
	if token.Wait() && token.Error() != nil {
		logrus.Error(token.Error())
		return token.Error()
	}
	return nil
}

func PublishAttributeResponseMessage(deviceNumber string, messageId string, err error) error {
	qos := byte(config.MqttConfig.Attributes.QoS)
	topic := fmt.Sprintf("%s%s/%s", config.MqttConfig.Attributes.PublishResponseTopic, deviceNumber, messageId)

	payload := common.GetResponsePayload("", err)

	token := mqttClient.Publish(topic, qos, false, payload)
	if token.Wait() && token.Error() != nil {
		logrus.Error(token.Error())
	}
	return token.Error()
}

func PublishEventResponseMessage(deviceNumber string, messageId string, method string, err error) error {
	qos := byte(config.MqttConfig.Events.QoS)
	topic := fmt.Sprintf("%s%s/%s", config.MqttConfig.Events.PublishTopic, deviceNumber, messageId)

	payload := common.GetResponsePayload(method, err)

	token := mqttClient.Publish(topic, qos, false, payload)
	if token.Wait() && token.Error() != nil {
		logrus.Error(token.Error())
	}
	return token.Error()
}

func PublishGetAttributeMessage(deviceNumber string, payload []byte) error {
	topic := fmt.Sprintf("%s%s", config.MqttConfig.Attributes.PublishGetTopic, deviceNumber)
	qos := byte(config.MqttConfig.Attributes.QoS)

	token := mqttClient.Publish(topic, qos, false, payload)
	if token.Wait() && token.Error() != nil {
		logrus.Error(token.Error())
	}
	return token.Error()
}

func PublishEventMessage(payload []byte) error {
	topic := config.MqttConfig.Events.PublishTopic
	qos := byte(config.MqttConfig.Events.QoS)

	token := mqttClient.Publish(topic, qos, false, payload)
	if token.Wait() && token.Error() != nil {
		logrus.Error(token.Error())
	}
	return token.Error()
}

func PublishCommandMessage(topic string, payload []byte) error {
	qos := byte(config.MqttConfig.Commands.QoS)

	token := mqttClient.Publish(topic, qos, false, payload)
	if token.Wait() && token.Error() != nil {
		logrus.Error(token.Error())
	}

	logrus.Debug("Issued topic:", topic)
	logrus.Debug("Issued command:", string(payload))

	return token.Error()
}

func ForwardTelemetryMessage(deviceId string, payload []byte) error {
	telemetryTopic := config.MqttConfig.Telemetry.SubscribeTopic + "/" + deviceId
	qos := byte(config.MqttConfig.Telemetry.QoS)

	token := mqttClient.Publish(telemetryTopic, qos, false, payload)
	if token.Wait() && token.Error() != nil {
		logrus.Error(token.Error())
	}
	return token.Error()
}

func PublishOnlineMessage(deviceID string, payload []byte) error {
	topic := fmt.Sprintf("devices/status/%s", deviceID)
	topic = path.Join("$share/mygroup", topic)
	qos := byte(0)

	token := mqttClient.Publish(topic, qos, false, payload)
	if token.Wait() && token.Error() != nil {
		logrus.Error(token.Error())
	}
	return token.Error()
}
