package data

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type SensorData struct {
	Humidity    float64 `json:"humidity"`
	Temperature float64 `json:"temperature"`
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connection lost: %v\n", err)
}

func generateRandomData() SensorData {
	// Generate random humidity around 30 with small variations
	humidity := 25 + (rand.Float64()-0.5)*5 // ±2.5 variation around 25
	// Generate random temperature around 60 with small variations
	temperature := 60 + (rand.Float64()-0.5)*3 // ±1.5 variation around 60

	return SensorData{
		Humidity:    humidity,
		Temperature: temperature,
	}
}

func StartMQTTClient(broker string, port int, username string, password string, clientId string) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID(clientId)
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// Start sending data every 5 seconds
	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				data := generateRandomData()
				jsonData, err := json.Marshal(data)
				if err != nil {
					fmt.Printf("Error marshaling data: %v\n", err)
					continue
				}

				token := client.Publish("devices/telemetry", 0, false, jsonData)
				token.Wait()
				fmt.Printf("Published: %s\n", jsonData)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	// Keep the program running
	select {}
}
