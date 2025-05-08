package main

import (
	"github.com/Thingsly/backend/data"
)

func main() {
	broker := "localhost"
	port := 1883
	username := "d9eef969-db86-fc81-499"
	password := "7517c59"
	clientId := "mqtt_b4dc6235-dfd"

	// data.StartMQTTClient(broker, port, username, password, clientId)
	data.StartAbnormalMQTTClient(broker, port, username, password, clientId)
}
