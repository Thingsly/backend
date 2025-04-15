package ws_publish

import (
	"encoding/json"
	"strconv"
	"sync"
	"time"

	config "github.com/Thingsly/backend/mqtt"
	"github.com/Thingsly/backend/mqtt/subscribe"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-basic/uuid"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type WsMqttClient struct {
	Client mqtt.Client
}

func (w *WsMqttClient) CreateMqttClient() error {

	opts := mqtt.NewClientOptions()
	opts.AddBroker(config.MqttConfig.Broker)
	opts.SetUsername(config.MqttConfig.User)
	opts.SetPassword(config.MqttConfig.Pass)
	opts.SetClientID("ws_mqtt_" + uuid.New()[0:8])

	opts.SetCleanSession(true)

	opts.SetResumeSubs(false)

	opts.SetAutoReconnect(true)
	opts.SetConnectRetryInterval(5 * time.Second)
	opts.SetMaxReconnectInterval(20 * time.Second)

	opts.SetOrderMatters(false)
	opts.SetOnConnectHandler(func(_ mqtt.Client) {
		logrus.Println("ws mqtt connect success")
	})

	w.Client = mqtt.NewClient(opts)
	if token := w.Client.Connect(); token.Wait() && token.Error() != nil {
		logrus.Error("Ws MQTT Broker connection failed:", token.Error())
		return token.Error()
	}
	return nil
}

func (w *WsMqttClient) SubscribeDeviceTelemetry(deviceId string, conn *websocket.Conn, msgType int, mu *sync.Mutex) error {
	err := w.CreateMqttClient()
	if err != nil {
		return err
	}

	deviceTelemetryHandler := func(_ mqtt.Client, d mqtt.Message) {

		var valuesMap map[string]interface{}
		if err := json.Unmarshal(d.Payload(), &valuesMap); err != nil {
			logrus.Error(err)
			mu.Lock()
			conn.WriteMessage(msgType, []byte(err.Error()))
			mu.Unlock()
			return
		}

		valuesMap["systime"] = time.Now().UTC()

		data, err := json.Marshal(valuesMap)
		if err != nil {
			logrus.Error(err)
			mu.Lock()
			conn.WriteMessage(msgType, []byte(err.Error()))
			mu.Unlock()
			return
		}
		mu.Lock()
		err = conn.WriteMessage(msgType, data)
		mu.Unlock()
		if err != nil {
			logrus.Error(err)
			conn.WriteMessage(msgType, []byte(err.Error()))
			return
		}
	}
	telemetryTopic := config.MqttConfig.Telemetry.SubscribeTopic + "/" + deviceId
	telemetryQos := byte(0)
	if token := w.Client.Subscribe(telemetryTopic, telemetryQos, deviceTelemetryHandler); token.Wait() && token.Error() != nil {
		logrus.Error(token.Error())
		return token.Error()
	}
	return nil
}

func (w *WsMqttClient) SubscribeDeviceTelemetryByKeys(deviceId string, conn *websocket.Conn, msgType int, mu *sync.Mutex, keys []string) error {
	err := w.CreateMqttClient()
	if err != nil {
		return err
	}

	deviceTelemetryHandler := func(_ mqtt.Client, d mqtt.Message) {

		var valuesMap map[string]interface{}
		var rspMap = make(map[string]interface{})
		if err := json.Unmarshal(d.Payload(), &valuesMap); err != nil {
			logrus.Error(err)
			mu.Lock()
			conn.WriteMessage(msgType, []byte(err.Error()))
			mu.Unlock()
			return
		}

		for _, key := range keys {
			if value, ok := valuesMap[key]; ok {
				rspMap[key] = value
			}
		}

		rspMap["systime"] = time.Now().UTC()

		data, err := json.Marshal(rspMap)
		if err != nil {
			logrus.Error(err)
			mu.Lock()
			conn.WriteMessage(msgType, []byte(err.Error()))
			mu.Unlock()
			return
		}
		mu.Lock()
		err = conn.WriteMessage(msgType, data)
		mu.Unlock()
		if err != nil {
			logrus.Error(err)
			conn.WriteMessage(msgType, []byte(err.Error()))
			return
		}
	}
	telemetryTopic := config.MqttConfig.Telemetry.SubscribeTopic + "/" + deviceId
	telemetryQos := byte(0)
	if token := w.Client.Subscribe(telemetryTopic, telemetryQos, deviceTelemetryHandler); token.Wait() && token.Error() != nil {
		logrus.Error(token.Error())
		return token.Error()
	}
	return nil
}

func (w *WsMqttClient) SubscribeOnlineOffline(deviceId string, conn *websocket.Conn, msgType int, mu *sync.Mutex) error {
	err := w.CreateMqttClient()
	if err != nil {
		return err
	}

	onlineOfflineHandler := func(_ mqtt.Client, d mqtt.Message) {

		payloadInt, err := strconv.Atoi(string(d.Payload()))
		if err != nil {
			logrus.Error(err.Error())
			return
		}

		payloadMap := make(map[string]interface{})
		payloadMap["is_online"] = payloadInt

		data, err := json.Marshal(payloadMap)
		if err != nil {
			logrus.Error(err)
			mu.Lock()
			conn.WriteMessage(msgType, []byte(err.Error()))
			mu.Unlock()
			return
		}
		mu.Lock()
		err = conn.WriteMessage(msgType, data)
		mu.Unlock()
		if err != nil {
			logrus.Error(err)
			conn.WriteMessage(msgType, []byte(err.Error()))
			return
		}
	}
	onlineOfflineTopic := "devices/status/" + deviceId
	onlineOfflineTopic = subscribe.GenTopic(onlineOfflineTopic)
	onlineOfflineQos := byte(0)
	if token := w.Client.Subscribe(onlineOfflineTopic, onlineOfflineQos, onlineOfflineHandler); token.Wait() && token.Error() != nil {
		logrus.Error(token.Error())
		return token.Error()
	}
	return nil
}

func (w *WsMqttClient) Close() {
	w.Client.Disconnect(250)
}
