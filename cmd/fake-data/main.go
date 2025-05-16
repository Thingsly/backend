package main

import (
	"github.com/Thingsly/backend/data"
)

func main() {
	broker := "localhost"
	port := 1883
	MQTT_username := "26fdbc41-9999-f981-26f"
	MQTT_password := "7e00453"
	MQTT_clientId := "mqtt_45d73c7f-b8a"

	// data.StartMQTTClient(broker, port, MQTT_username, MQTT_password, MQTT_clientId)
	data.StartAbnormalMQTTClient(broker, port, MQTT_username, MQTT_password, MQTT_clientId)
}
