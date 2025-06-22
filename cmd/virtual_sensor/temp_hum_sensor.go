package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/Thingsly/backend/internal/model"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var (
	mqttClient        *mqtt.Client
	gatewayMqttClient *mqtt.Client
	switchStatus      int64 = 0
)

func TempHumSensor() {

	createClient()

	subscribeControlMessage()

	go publishTelemetryMessage("devices/telemetry")

	go publishAttributeMessage("devices/attributes/")

	go publishEventMessage("devices/event/")

	// createGatewayClient()

	// go publishGatewayTelemetryMessage("gateway/telemetry")

	// go publishGatewayAttributeMessage("gateway/attributes/")

	// go publishGatewayEventMessage("gateway/event/")

	select {}
}

func createClient() {

	opts := MqttConfig{
		Broker: "127.0.0.1:1883",
		// Broker: "103.124.93.210:1883",
		User:   "8cc60abf-40ab-b725-6d9",
		Pass:   "b7e693c",
	}
	mqttClient = CreateMqttClient(opts)
}

func createGatewayClient() {

	opts := MqttConfig{
		// Broker: "127.0.0.1:1883",
		Broker: "103.124.93.210:1883",
		User:   "2ff45516-be20-1d3e-afc",
		Pass:   "ed287f9",
	}
	gatewayMqttClient = CreateMqttClient(opts)
}

// Subscribe to control messages
func subscribeControlMessage() {
	topic := "devices/telemetry/control/5cd3b3a7-ff2b-02d4-63e6-12797e223c11"
	token := (*mqttClient).Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		var controlMsg map[string]interface{}
		err := json.Unmarshal(msg.Payload(), &controlMsg)
		if err != nil {
			log.Printf("Failed to parse control message: %v", err)
			return
		}

		// Check if the message contains switchStatus
		if status, ok := controlMsg["switchStatus"].(float64); ok {
			switchStatus = int64(status)
			log.Printf("Received switch control command: %d", switchStatus)

			// Immediately send a telemetry message for the switch status change
			message := make(map[string]interface{})
			message["switchStatus"] = switchStatus
			payload, err := json.Marshal(message)
			if err != nil {
				log.Printf("Failed to generate telemetry message: %v", err)
				return
			}
			token := (*mqttClient).Publish("devices/telemetry", 0, false, payload)
			token.Wait()
			log.Printf("Sent switch status change telemetry: %d", switchStatus)
		}
	})
	token.Wait()
	log.Printf("Subscribed to control topic: %s", topic)
}

// Publish telemetry message
func publishTelemetryMessage(topic string) {
	// Publish a message every 10 seconds
	for {
		message := make(map[string]interface{})
		// Generate a random temperature between -20 and 40 degrees, rounded to two decimal places
		t, err := generateRandomFloat()
		if err != nil {
			log.Println("generateRandomFloat failed:", err)
		}
		message["temperature"] = t
		// Round to two decimal places
		message["temperature"] = float64(int(message["temperature"].(float64)*100)) / 100
		// Generate a random humidity value between 0% and 100%
		h, err := generateRandomFloat()
		if err != nil {
			log.Println("generateRandomFloat failed:", err)
		}
		message["humidity"] = h

		// Add a boolean value for the device's online status
		message["isOnline"] = true

		// Use the global variable for the switch status
		message["switchStatus"] = switchStatus

		// Add a random device status string
		statuses := []string{"running", "idle", "error", "maintenance"}
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(statuses))))
		if err != nil {
			log.Println("generate random number failed:", err)
			message["deviceStatus"] = "running"
		} else {
			message["deviceStatus"] = statuses[n.Int64()]
		}

		// Convert to JSON format
		var payload []byte
		payload, err = json.Marshal(message)
		if err != nil {
			log.Println("json.Marshal failed:", err)
			return
		}
		token := (*mqttClient).Publish(topic, 0, false, payload)
		isSuccess := token.Wait()
		if !isSuccess {
			log.Println("Publish message failed", string(payload))
		} else {
			log.Println("Publish message successful:", string(payload))
		}
		// Wait for 30 seconds before publishing the next message
		<-time.After(30 * time.Second)
	}
}

func publishAttributeMessage(topic string) {
	// Publish a message every 30 seconds
	for {
		message := make(map[string]interface{})
		message["version"] = "1.0.0"
		message["status"] = "normal"
		message["mac"] = "00:11:22:33:44:55"
		message["longitude"] = 105.804817
		message["latitude"] = 21.028511
		// Convert to JSON format
		var payload []byte
		payload, err := json.Marshal(message)
		if err != nil {
			log.Println("json.Marshal failed:", err)
			return
		}
		messageId := GetMessageID()
		token := (*mqttClient).Publish(topic+messageId, 0, false, payload)
		isSuccess := token.Wait()
		if !isSuccess {
			log.Println("Publish message failed", string(payload))
		} else {
			log.Println("Publish message successful:", string(payload))
		}
		// Wait for 120 seconds before publishing the next message
		<-time.After(120 * time.Second)
	}
}

func publishEventMessage(topic string) {
	// Publish a message every 60 seconds
	for {
		message := make(map[string]interface{})

		message["method"] = "alert"
		// params is a map type
		message["params"] = map[string]interface{}{
			"level":   "warning",
			"message": "temperature is too high",
		}
		// Convert to JSON format
		var payload []byte
		payload, err := json.Marshal(message)
		if err != nil {
			log.Println("json.Marshal failed:", err)
			return
		}
		messageId := GetMessageID()
		token := (*mqttClient).Publish(topic+messageId, 0, false, payload)
		isSuccess := token.Wait()
		if !isSuccess {
			log.Println("Publish message failed", string(payload))
		} else {
			log.Println("Publish message successful:", string(payload))
		}
		// Wait for 120 seconds before publishing the next message
		<-time.After(120 * time.Second)
	}
}

func getTelemetryMessageParams() *map[string]interface{} {
	message := make(map[string]interface{})
	t, err := generateRandomFloat()
	if err != nil {
		log.Println("generateRandomFloat failed:", err)
		return nil
	}
	// Generate a random temperature between -20 and 40 degrees, rounded to two decimal places
	message["temperature"] = t
	// Round to two decimal places
	message["temperature"] = float64(int(message["temperature"].(float64)*100)) / 100
	// Generate a random humidity value between 0% and 100%
	h, err := generateRandomFloat()
	if err != nil {
		log.Println("generateRandomFloat failed:", err)
		return nil
	}
	message["humidity"] = h

	return &message
}

func generateRandomFloat() (float64, error) {
	// Generate the integer part between [10.00, 99.99]
	integer, err := rand.Int(rand.Reader, big.NewInt(90))
	if err != nil {
		return 0, fmt.Errorf("failed to generate integer part: %v", err)
	}
	integer = integer.Add(integer, big.NewInt(10))

	// Generate the decimal part between [0, 99]
	decimal, err := rand.Int(rand.Reader, big.NewInt(100))
	if err != nil {
		return 0, fmt.Errorf("failed to generate decimal part: %v", err)
	}

	// Combine the integer and decimal parts
	result := float64(integer.Int64()) + float64(decimal.Int64())/100.0

	return result, nil
}

func getAttributeMessageParams() *map[string]interface{} {
	message := make(map[string]interface{})
	message["version"] = "1.0.0"
	message["status"] = "normal"
	message["mac"] = "00:11:22:33:44:55"
	message["longitude"] = 105.804817
	message["latitude"] = 21.028511

	return &message
}

func getEventMessageParams() *map[string]interface{} {
	message := make(map[string]interface{})

	message["method"] = "alert"
	message["params"] = map[string]interface{}{
		"level":   "warning",
		"message": "temperature is too high",
	}

	return &message
}

func publishGatewayTelemetryMessage(topic string) {

	for {
		subDevice := make(map[string]map[string]interface{})
		subDevice["fb7871c3"] = *getTelemetryMessageParams()
		payloads := &model.GatewayPublish{
			GatewayData:   getTelemetryMessageParams(),
			SubDeviceData: &subDevice,
		}

		var payload []byte
		payload, err := json.Marshal(payloads)
		if err != nil {
			log.Println("json.Marshal failed:", err)
			return
		}
		token := (*gatewayMqttClient).Publish(topic, 0, false, payload)
		token.Wait()
		log.Println("Publish message:", string(payload))

		<-time.After(50 * time.Second)
	}
}

func publishGatewayAttributeMessage(topic string) {

	for {
		subDevice := make(map[string]map[string]interface{})
		subDevice["fb7871c3"] = *getAttributeMessageParams()
		payloads := &model.GatewayPublish{
			GatewayData:   getAttributeMessageParams(),
			SubDeviceData: &subDevice,
		}

		var payload []byte
		payload, err := json.Marshal(payloads)
		if err != nil {
			log.Println("json.Marshal failed:", err)
			return
		}
		messageId := GetMessageID()
		token := (*gatewayMqttClient).Publish(topic+messageId, 0, false, payload)
		token.Wait()
		log.Println("Publish message:", string(payload))

		<-time.After(40 * time.Second)
	}
}

func publishGatewayEventMessage(topic string) {

	for {
		subDevice := make(map[string]map[string]interface{})
		subDevice["fb7871c3"] = *getEventMessageParams()
		payloads := &model.GatewayPublish{
			GatewayData:   getEventMessageParams(),
			SubDeviceData: &subDevice,
		}

		var payload []byte
		payload, err := json.Marshal(payloads)
		if err != nil {
			log.Println("json.Marshal failed:", err)
			return
		}
		messageId := GetMessageID()
		token := (*gatewayMqttClient).Publish(topic+messageId, 0, false, payload)
		token.Wait()
		log.Println("Publish message:", string(payload))

		<-time.After(30 * time.Second)
	}
}
