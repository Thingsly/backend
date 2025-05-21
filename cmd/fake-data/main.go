package main

import (
	"github.com/Thingsly/backend/data"
)

func main() {
	broker := "103.124.93.210"
	port := 1883
	MQTT_username := "83caed50-f4e4-abbb-566"
	MQTT_password := "daf7dd7"
	MQTT_clientId := "mqtt_258f61e7-b3d"

	// data.StartMQTTClient(broker, port, MQTT_username, MQTT_password, MQTT_clientId)
	data.StartAbnormalMQTTClient(broker, port, MQTT_username, MQTT_password, MQTT_clientId)
}
