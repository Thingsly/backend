package data

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func generateAbnormalData() SensorData {
	// Generate high humidity between 80 and 95
	humidity := 80 + rand.Float64()*15
	// Generate high temperature between 80 and 90
	temperature := 80 + rand.Float64()*10

	return SensorData{
		Humidity:    humidity,
		Temperature: temperature,
	}
}

func StartAbnormalMQTTClient(broker string, port int, username string, password string, clientId string) {
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

	// Start sending abnormal data every 5 seconds
	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				data := generateAbnormalData()
				jsonData, err := json.Marshal(data)
				if err != nil {
					fmt.Printf("Error marshaling data: %v\n", err)
					continue
				}

				token := client.Publish("devices/telemetry", 0, false, jsonData)
				token.Wait()
				fmt.Printf("Published abnormal data: %s\n", jsonData)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	// Keep the program running
	select {}
}
