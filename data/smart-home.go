package data

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// ControlMessage represents the structure of control commands
type ControlMessage struct {
	Relay int64 `json:"relay"`
}

// SmartHomeData represents all sensor data for smart home
type SmartHomeData struct {
	// Environmental parameters
	IndoorTemperature float64 `json:"indoor_temperature"` // Celsius (°C)
	IndoorHumidity    float64 `json:"indoor_humidity"`    // percentage (%)
	AirQuality        float64 `json:"air_quality"`        // parts per million (ppm)
	CO2Level          float64 `json:"co2_level"`          // parts per million (ppm)
	TVOC              float64 `json:"tvoc"`               // parts per billion (ppb)
	PM25              float64 `json:"pm25"`               // micrograms per cubic meter (µg/m³)
	PM10              float64 `json:"pm10"`               // micrograms per cubic meter (µg/m³)
	Noise             float64 `json:"noise"`              // decibel (dB)
	LightLevel        float64 `json:"light_level"`        // lux

	// Energy parameters
	PowerConsumption float64 `json:"power_consumption"` // kilowatt-hour (kWh)
	Voltage          float64 `json:"voltage"`           // volt (V)
	Current          float64 `json:"current"`           // ampere (A)

	// Security parameters
	MotionDetected int64   `json:"motion_detected"` // 0: false, 1: true
	SmokeLevel     float64 `json:"smoke_level"`     // parts per million (ppm)
	GasLevel       float64 `json:"gas_level"`       // parts per million (ppm)

	// Device control
	RelayState int64 `json:"relay"`      // 0: OFF, 1: ON
	LEDStatus  int64 `json:"led_status"` // 0: OFF, 1: ON
}

// Global variable to store the current relay state
var currentRelayState int64 = 0

// Control message handler
func controlMessageHandler(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received control message on topic %s: %s\n", msg.Topic(), string(msg.Payload()))

	var controlMsg ControlMessage
	if err := json.Unmarshal(msg.Payload(), &controlMsg); err != nil {
		fmt.Printf("Error parsing control message: %v\n", err)
		return
	}

	// Update relay state based on control message
	if controlMsg.Relay == 0 || controlMsg.Relay == 1 {
		currentRelayState = controlMsg.Relay
		fmt.Printf("Relay state updated to: %d\n", currentRelayState)
	} else {
		fmt.Printf("Invalid relay state received: %d\n", controlMsg.Relay)
	}
}

func generateRandomSmartHomeData() SmartHomeData {
	// Environmental parameters with realistic ranges
	indoorTemp := roundToTwoDecimals(22 + (rand.Float64()-0.5)*6)      // 19-25°C (comfortable indoor temperature)
	indoorHumidity := roundToTwoDecimals(45 + (rand.Float64()-0.5)*20) // 35-55% (comfortable indoor humidity)
	airQuality := roundToTwoDecimals(50 + (rand.Float64()-0.5)*40)     // 30-70 ppm (indoor air quality)
	co2Level := roundToTwoDecimals(600 + (rand.Float64()-0.5)*400)     // 400-800 ppm (indoor CO2 level)
	tvoc := roundToTwoDecimals(200 + (rand.Float64()-0.5)*150)         // 125-275 ppb (TVOC level)
	pm25 := roundToTwoDecimals(10 + (rand.Float64()-0.5)*8)            // 6-14 µg/m³ (PM2.5 level)
	pm10 := roundToTwoDecimals(20 + (rand.Float64()-0.5)*15)           // 12.5-27.5 µg/m³ (PM10 level)
	noise := roundToTwoDecimals(35 + (rand.Float64()-0.5)*10)          // 30-40 dB (typical indoor noise)
	lightLevel := roundToTwoDecimals(300 + (rand.Float64()-0.5)*200)   // 200-400 lux (typical indoor lighting)

	// Energy parameters with realistic ranges
	powerConsumption := roundToTwoDecimals(0.5 + (rand.Float64()-0.5)*0.3) // 0.35-0.65 kWh
	voltage := roundToTwoDecimals(220 + (rand.Float64()-0.5)*10)           // 215-225V
	current := roundToTwoDecimals(2 + (rand.Float64()-0.5)*0.5)            // 1.75-2.25A

	// Security parameters with realistic ranges
	var motionDetected int64
	if rand.Float64() > 0.7 { // 30% chance of motion
		motionDetected = 1
	} else {
		motionDetected = 0
	}
	smokeLevel := roundToTwoDecimals(0.1 + (rand.Float64()-0.5)*0.1) // 0.05-0.15 ppm
	gasLevel := roundToTwoDecimals(0.2 + (rand.Float64()-0.5)*0.15)  // 0.125-0.275 ppm

	// Device control - relay state is only controlled manually
	relayState := currentRelayState
	ledStatus := relayState // LED status follows relay state

	return SmartHomeData{
		IndoorTemperature: indoorTemp,
		IndoorHumidity:    indoorHumidity,
		AirQuality:        airQuality,
		CO2Level:          co2Level,
		TVOC:              tvoc,
		PM25:              pm25,
		PM10:              pm10,
		Noise:             noise,
		LightLevel:        lightLevel,
		PowerConsumption:  powerConsumption,
		Voltage:           voltage,
		Current:           current,
		MotionDetected:    motionDetected,
		SmokeLevel:        smokeLevel,
		GasLevel:          gasLevel,
		RelayState:        relayState,
		LEDStatus:         ledStatus,
	}
}

func StartSmartHomeMQTTClient(broker string, port int, username string, password string, clientId string, deviceId string) {
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

	// Subscribe to control topic
	controlTopic := fmt.Sprintf("devices/telemetry/control/%s", deviceId)
	if token := client.Subscribe(controlTopic, 0, controlMessageHandler); token.Wait() && token.Error() != nil {
		fmt.Printf("Error subscribing to control topic: %v\n", token.Error())
	} else {
		fmt.Printf("Successfully subscribed to control topic: %s\n", controlTopic)
	}

	// Start sending data every 5 seconds
	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				data := generateRandomSmartHomeData()
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
