package main

import (
	"github.com/Thingsly/backend/data"
)

func main() {
	broker := "localhost"
	port := 1883
	username := "b8dadde2-01ca-8758-71c"
	password := "24e43c5"
	clientId := "mqtt_9a63acb4-faa"

	data.StartMQTTClient(broker, port, username, password, clientId)
	// data.StartAbnormalMQTTClient(broker, port, username, password, clientId)
}
