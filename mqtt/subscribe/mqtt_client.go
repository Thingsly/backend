package subscribe

import (
	"path"
	"time"

	"github.com/Thingsly/backend/initialize"
	config "github.com/Thingsly/backend/mqtt"

	"github.com/Thingsly/backend/mqtt/publish"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-basic/uuid"
	"github.com/panjf2000/ants"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var SubscribeMqttClient mqtt.Client
var TelemetryMessagesChan chan map[string]interface{}

func GenTopic(topic string) string {
	topic = path.Join("$share/mygroup", topic)
	return topic
}

func SubscribeInit() error {

	initialize.NewAutomateLimiter()

	subscribeMqttClient()

	telemetryMessagesChan()

	err := subscribe()
	return err
}

func subscribe() error {

	err := SubscribeAttribute()
	if err != nil {
		logrus.Error(err)
		return err
	}

	err = SubscribeSetAttribute()
	if err != nil {
		logrus.Error(err)
		return err
	}

	err = SubscribeEvent()
	if err != nil {
		logrus.Error(err)
		return err
	}

	err = SubscribeTelemetry()
	if err != nil {
		logrus.Error(err)
		return err
	}

	err = GatewaySubscribeTopic()
	if err != nil {
		logrus.Error(err)
		return err
	}

	err = SubscribeCommand()
	if err != nil {
		logrus.Error(err)
		return err
	}

	err = SubscribeOtaUpprogress()
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func subscribeMqttClient() {

	opts := mqtt.NewClientOptions()
	opts.AddBroker(config.MqttConfig.Broker)
	opts.SetUsername(config.MqttConfig.User)
	opts.SetPassword(config.MqttConfig.Pass)
	id := "thingsly-go-sub-" + uuid.New()[0:8]
	opts.SetClientID(id)
	logrus.Info("clientid: ", id)

	opts.SetCleanSession(true)

	opts.SetResumeSubs(true)

	opts.SetAutoReconnect(true)
	opts.SetConnectRetryInterval(5 * time.Second)
	opts.SetMaxReconnectInterval(200 * time.Second)

	opts.SetOrderMatters(false)
	opts.SetOnConnectHandler(func(_ mqtt.Client) {
		logrus.Println("mqtt connect success")
	})

	opts.SetConnectionLostHandler(func(_ mqtt.Client, err error) {
		logrus.Println("mqtt connect  lost: ", err)
		SubscribeMqttClient.Disconnect(250)
		for {
			if token := SubscribeMqttClient.Connect(); token.Wait() && token.Error() != nil {
				logrus.Error("MQTT Broker 1 connection failed:", token.Error())
				time.Sleep(5 * time.Second)
				continue
			}
			subscribe()
			break
		}
	})

	SubscribeMqttClient = mqtt.NewClient(opts)

	for {
		if token := SubscribeMqttClient.Connect(); token.Wait() && token.Error() != nil {
			logrus.Error("MQTT Broker 1 connection failed:", token.Error())
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}

}

func telemetryMessagesChan() {
	TelemetryMessagesChan = make(chan map[string]interface{}, config.MqttConfig.ChannelBufferSize)
	writeWorkers := config.MqttConfig.WriteWorkers
	for i := 0; i < writeWorkers; i++ {
		go MessagesChanHandler(TelemetryMessagesChan)
	}
}

func SubscribeTelemetry() error {

	dbType := viper.GetString("grpc.tptodb_type")
	if dbType == "TSDB" || dbType == "KINGBASE" || dbType == "POLARDB" {
		logrus.Infof("dbType:%v do not need subcribe topic: %v", dbType, config.MqttConfig.Telemetry.SubscribeTopic)
		return nil
	}

	p, err := ants.NewPool(config.MqttConfig.Telemetry.PoolSize)
	if err != nil {
		return err
	}
	deviceTelemetryMessageHandler := func(_ mqtt.Client, d mqtt.Message) {
		err = p.Submit(func() {

			TelemetryMessages(d.Payload(), d.Topic())
		})
		if err != nil {
			logrus.Error(err)
		}
	}

	topic := config.MqttConfig.Telemetry.SubscribeTopic
	topic = GenTopic(topic)
	logrus.Info("subscribe topic:", topic)

	qos := byte(config.MqttConfig.Telemetry.QoS)

	if token := SubscribeMqttClient.Subscribe(topic, qos, deviceTelemetryMessageHandler); token.Wait() && token.Error() != nil {
		logrus.Error(token.Error())
		return err
	}
	return nil
}

func SubscribeAttribute() error {

	deviceAttributeHandler := func(_ mqtt.Client, d mqtt.Message) {

		logrus.Debug("attribute message:", string(d.Payload()))
		deviceNumber, messageId, err := DeviceAttributeReport(d.Payload(), d.Topic())
		logrus.Debug("Responding to device property report", deviceNumber, err)
		if err != nil {
			logrus.Error(err)
		}
		if deviceNumber != "" && messageId != "" {

			publish.PublishAttributeResponseMessage(deviceNumber, messageId, err)
		}
	}
	topic := config.MqttConfig.Attributes.SubscribeTopic
	topic = GenTopic(topic)
	logrus.Info("subscribe topic:", topic)
	qos := byte(config.MqttConfig.Attributes.QoS)
	if token := SubscribeMqttClient.Subscribe(topic, qos, deviceAttributeHandler); token.Wait() && token.Error() != nil {
		logrus.Error(token.Error())
		return token.Error()
	}
	return nil
}

func SubscribeSetAttribute() error {

	deviceAttributeHandler := func(_ mqtt.Client, d mqtt.Message) {

		logrus.Debug("attribute message:", string(d.Payload()))
		DeviceSetAttributeResponse(d.Payload(), d.Topic())
	}
	topic := config.MqttConfig.Attributes.SubscribeResponseTopic
	topic = GenTopic(topic)
	logrus.Info("subscribe topic:", topic)
	qos := byte(config.MqttConfig.Attributes.QoS)
	if token := SubscribeMqttClient.Subscribe(topic, qos, deviceAttributeHandler); token.Wait() && token.Error() != nil {
		logrus.Error(token.Error())
		return token.Error()
	}
	return nil
}

func SubscribeCommand() error {

	deviceCommandHandler := func(_ mqtt.Client, d mqtt.Message) {

		messageID, err := DeviceCommand(d.Payload(), d.Topic())
		logrus.Debug("Device command response report", messageID, err)
		if err != nil || messageID == "" {
			logrus.Debug("Device command response report failed", messageID, err)
			logrus.Error(err)
		}
	}
	topic := config.MqttConfig.Commands.SubscribeTopic
	topic = GenTopic(topic)
	logrus.Info("subscribe topic:", topic)
	qos := byte(config.MqttConfig.Commands.QoS)
	if token := SubscribeMqttClient.Subscribe(topic, qos, deviceCommandHandler); token.Wait() && token.Error() != nil {
		logrus.Error(token.Error())
		return token.Error()
	}
	return nil
}

func SubscribeEvent() error {

	deviceEventHandler := func(_ mqtt.Client, d mqtt.Message) {

		logrus.Debug("event message:", string(d.Payload()))
		deviceNumber, messageId, method, err := DeviceEvent(d.Payload(), d.Topic())
		logrus.Debug("Responding to device property report", deviceNumber, err)
		if err != nil {
			logrus.Error(err)
		}
		if deviceNumber != "" && messageId != "" {

			publish.PublishEventResponseMessage(deviceNumber, messageId, method, err)
		}
	}
	topic := config.MqttConfig.Events.SubscribeTopic
	qos := byte(config.MqttConfig.Events.QoS)
	if token := SubscribeMqttClient.Subscribe(topic, qos, deviceEventHandler); token.Wait() && token.Error() != nil {
		logrus.Error(token.Error())
		return token.Error()
	}
	return nil
}

func SubscribeOtaUpprogress() error {

	otaUpgradeHandler := func(_ mqtt.Client, d mqtt.Message) {

		logrus.Debug("ota upgrade message:", string(d.Payload()))
		OtaUpgrade(d.Payload(), d.Topic())
	}
	topic := config.MqttConfig.OTA.SubscribeTopic
	qos := byte(config.MqttConfig.OTA.QoS)
	if token := SubscribeMqttClient.Subscribe(topic, qos, otaUpgradeHandler); token.Wait() && token.Error() != nil {
		logrus.Error(token.Error())
		return token.Error()
	}
	return nil
}
