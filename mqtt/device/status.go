package device

import (
	"time"

	config "github.com/HustIoTPlatform/backend/mqtt"
	"github.com/HustIoTPlatform/backend/mqtt/subscribe"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-basic/uuid"
	"github.com/sirupsen/logrus"
)

type StatusManager struct {
	mqttClient     mqtt.Client
	subscribeTopic string
	subscribeQos   byte
	messageHandler mqtt.MessageHandler

	retryInterval time.Duration
	maxRetries    int
}

type StatusConfig struct {
	Broker        string
	ClientID      string
	Username      string
	Password      string
	RetryInterval time.Duration
	MaxRetries    int
}

func InitDeviceStatus() error {
	uuid := uuid.New()

	config := StatusConfig{
		Broker:        config.MqttConfig.Broker,
		ClientID:      "device-status-" + uuid[0:10],
		Username:      config.MqttConfig.User,
		Password:      config.MqttConfig.Pass,
		RetryInterval: 5 * time.Second,
		MaxRetries:    0,
	}

	manager, err := NewStatusManager(config)
	if err != nil {
		logrus.WithError(err).Error("Failed to create status manager")
		return err
	}

	defer manager.Stop()

	if err := manager.Start(); err != nil {
		logrus.WithError(err).Error("Failed to start status monitoring")
		return err
	}

	logrus.Info("Device status monitoring started")

	select {}
}

func NewStatusManager(config StatusConfig) (*StatusManager, error) {
    messageHandler := func(_ mqtt.Client, msg mqtt.Message) {
        logrus.WithFields(logrus.Fields{
            "topic":   msg.Topic(),
            "payload": string(msg.Payload()),
        }).Debug("Received device status message")

        subscribe.DeviceOnline(msg.Payload(), msg.Topic())
    }

    manager := &StatusManager{
        subscribeTopic: "devices/status/+",
        subscribeQos:   byte(1),
        messageHandler: messageHandler,
        retryInterval:  config.RetryInterval,
        maxRetries:     config.MaxRetries,
    }

    opts := mqtt.NewClientOptions().
        AddBroker(config.Broker).
        SetClientID(config.ClientID).
        SetUsername(config.Username).
        SetPassword(config.Password).
        SetAutoReconnect(true).
        SetCleanSession(false)

    opts.SetOnConnectHandler(func(_ mqtt.Client) {
        logrus.Info("Connected to MQTT broker")
        if err := manager.subscribe(); err != nil {
            logrus.WithError(err).Error("Failed to resubscribe")
        }
    })

    opts.SetConnectionLostHandler(func(_ mqtt.Client, err error) {
        logrus.WithError(err).Warn("Lost connection to MQTT broker")
    })

    client := mqtt.NewClient(opts)
    manager.mqttClient = client

    if err := manager.connectWithRetry(); err != nil {
        return nil, err
    }

    return manager, nil
}

func (sm *StatusManager) connectWithRetry() error {
	retryCount := 0
	for {
		logrus.WithFields(logrus.Fields{
			"retry_count": retryCount,
			"max_retries": sm.maxRetries,
		}).Info("Attempting to connect to MQTT broker")

		token := sm.mqttClient.Connect()
		if token.WaitTimeout(10*time.Second) && token.Error() == nil {
			return nil
		}

		if sm.maxRetries > 0 && retryCount >= sm.maxRetries {
			return token.Error()
		}

		retryCount++
		logrus.WithFields(logrus.Fields{
			"retry_count": retryCount,
			"interval":    sm.retryInterval,
			"error":       token.Error(),
		}).Warn("Connection failed, retrying...")

		time.Sleep(sm.retryInterval)
	}
}

func (sm *StatusManager) subscribe() error {
    logrus.WithField("topic", sm.subscribeTopic).Info("Subscribing to device status topic")

    if token := sm.mqttClient.Subscribe(sm.subscribeTopic, sm.subscribeQos, sm.messageHandler); token.Wait() && token.Error() != nil {
        logrus.WithError(token.Error()).Error("Failed to subscribe to topic")
        return token.Error()
    }
    return nil
}

func (sm *StatusManager) Start() error {
	return sm.subscribe()
}

func (sm *StatusManager) Stop() {
	if sm.mqttClient != nil && sm.mqttClient.IsConnected() {
		logrus.Info("Disconnecting MQTT connection")
		sm.mqttClient.Disconnect(250)
	}
}
