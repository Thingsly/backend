package main

import (
	"github.com/Thingsly/backend/data"
)

func main() {
	broker := "103.124.93.210"
	port := 1883
 
	// // Normal MQTT
	// MQTT_username := "fdcf8b60-f795-25e8-b62"
	// MQTT_password := "2e03a3a"
	// MQTT_clientId := "mqtt_f4b16b00-529"

	// Smart Agriculture
	MQTT_username_smart_agriculture := "5e77568b-fd7e-ee0a-8f0"
	MQTT_password_smart_agriculture := "d4773c8"
	MQTT_clientId_smart_agriculture := "mqtt_d21e6836-41a"

	// // Smart Home
	// MQTT_username_smart_home := "417b098a-f1c6-468a-a68"
	// MQTT_password_smart_home := "548c409"
	// MQTT_clientId_smart_home := "mqtt_a0efdd55-e2e"
	// MQTT_deviceId_smart_home := "a0efdd55-e2ef-d7ff-0b79-601340e90b5f"

	// data.StartMQTTClient(broker, port, MQTT_username, MQTT_password, MQTT_clientId)
	// data.StartAbnormalMQTTClient(broker, port, MQTT_username, MQTT_password, MQTT_clientId)
	data.StartSmartAgricultureMQTTClient(broker, port, MQTT_username_smart_agriculture, MQTT_password_smart_agriculture, MQTT_clientId_smart_agriculture)
	// data.StartSmartHomeMQTTClient(broker, port, MQTT_username_smart_home, MQTT_password_smart_home, MQTT_clientId_smart_home, MQTT_deviceId_smart_home)
}
