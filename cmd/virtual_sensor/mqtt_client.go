package main

import (
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-basic/uuid"
)

type MqttConfig struct {
	Broker string
	User   string
	Pass   string
}

func CreateMqttClient(config MqttConfig) *mqtt.Client {

	opts := mqtt.NewClientOptions()
	opts.AddBroker(config.Broker)
	opts.SetUsername(config.User)
	if config.Pass != "" {
		opts.SetPassword(config.Pass)
	}
	opts.SetClientID(uuid.New())

	opts.SetCleanSession(true)
	opts.SetAutoReconnect(true)

	opts.SetResumeSubs(false)
	opts.SetAutoReconnect(true)
	opts.SetConnectRetryInterval(5 * time.Second)
	opts.SetMaxReconnectInterval(20 * time.Second)

	opts.SetOrderMatters(false)
	opts.SetOnConnectHandler(func(_ mqtt.Client) {
		log.Println("mqtt connect success")
	})

	opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		log.Println("mqtt connect  lost: ", err)

		for {
			token := client.Connect()
			if token.Wait() && token.Error() == nil {
				log.Println("Reconnected to MQTT broker")
				break
			}
			log.Printf("Reconnect failed: %v\n", token.Error())
			time.Sleep(5 * time.Second)
		}
	})
	mqttClient := mqtt.NewClient(opts)
	for {
		if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
			log.Println(token.Error())
			time.Sleep(15 * time.Second)
		} else {
			break
		}
	}
	return &mqttClient
}
