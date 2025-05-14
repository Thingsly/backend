package main

import (
	"github.com/Thingsly/backend/data"
)

func main() {
	broker := "localhost"
	port := 1883
	username := "440cc282-4fc9-c7c5-53d"
	password := "ab57154"
	clientId := "mqtt_c32768b6-013"

	data.StartMQTTClient(broker, port, username, password, clientId)
	// data.StartAbnormalMQTTClient(broker, port, username, password, clientId)
}
