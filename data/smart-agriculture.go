package data

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// SmartAgricultureData represents all sensor data for smart agriculture
type SmartAgricultureData struct {
	// Environmental parameters
	Noise           float64 `json:"noise"`            // decibel (dB)
	Humidity        float64 `json:"humidity"`         // percentage (%)
	AirPressure     float64 `json:"air_pressure"`     // hectopascal (hPa)
	LightValue      float64 `json:"light_value"`      // lux
	CarbonDioxide   float64 `json:"carbon_dioxide"`   // parts per million (ppm)
	UVIndex         float64 `json:"uv_index"`         // UV index (0-11+)
	Radiation       float64 `json:"radiation"`        // watts per square meter (W/m²)
	WindSpeed       float64 `json:"wind_speed"`       // meters per second (m/s)
	CurrentRainfall float64 `json:"current_rainfall"` // millimeters (mm)
	Temperature     float64 `json:"temperature"`      // Celsius (°C)

	// Soil parameters
	SoilMoisture     float64 `json:"soil_moisture"`     // percentage (%)
	SoilConductivity float64 `json:"soil_conductivity"` // microsiemens per centimeter (µS/cm)
	SoilPH           float64 `json:"soil_ph"`           // pH scale (0-14)
	SoilTemperature  float64 `json:"soil_temperature"`  // Celsius (°C)

	// Water parameters
	WaterLevel float64 `json:"water_level"` // centimeters (cm)
}

func roundToTwoDecimals(value float64) float64 {
	return math.Round(value*100) / 100
}

func generateRandomSmartAgricultureData() SmartAgricultureData {
	// Environmental parameters with realistic ranges
	noise := roundToTwoDecimals(30 + (rand.Float64()-0.5)*20)            // 20-40 dB (typical ambient noise)
	humidity := roundToTwoDecimals(40 + (rand.Float64()-0.5)*40)         // 20-60% (typical humidity range)
	airPressure := roundToTwoDecimals(1000 + (rand.Float64()-0.5)*20)    // 990-1010 hPa (normal atmospheric pressure)
	lightValue := roundToTwoDecimals(10000 + (rand.Float64()-0.5)*20000) // 0-20000 lux (daylight range)
	carbonDioxide := roundToTwoDecimals(400 + (rand.Float64()-0.5)*200)  // 300-500 ppm (normal outdoor CO2 levels)
	uvIndex := roundToTwoDecimals(rand.Float64() * 11)                   // 0-11+ (UV index scale)
	radiation := roundToTwoDecimals(100 + (rand.Float64()-0.5)*200)      // 0-300 W/m² (solar radiation)
	windSpeed := roundToTwoDecimals(rand.Float64() * 10)                 // 0-10 m/s (wind speed)
	currentRainfall := roundToTwoDecimals(rand.Float64() * 5)            // 0-5 mm (current rainfall)
	temperature := roundToTwoDecimals(20 + (rand.Float64()-0.5)*20)      // 10-30°C (typical temperature range)

	// Soil parameters with realistic ranges
	soilMoisture := roundToTwoDecimals(20 + (rand.Float64()-0.5)*40)       // 0-40% (soil moisture)
	soilConductivity := roundToTwoDecimals(100 + (rand.Float64()-0.5)*200) // 0-300 µS/cm (soil conductivity)
	soilPH := roundToTwoDecimals(5.5 + (rand.Float64()-0.5)*3)             // 4-7 pH (typical soil pH range)
	soilTemperature := roundToTwoDecimals(15 + (rand.Float64()-0.5)*20)    // 5-25°C (soil temperature)

	// Water parameters
	waterLevel := roundToTwoDecimals(20 + (rand.Float64()-0.5)*40) // 0-40 cm (water level)

	return SmartAgricultureData{
		Noise:            noise,
		Humidity:         humidity,
		AirPressure:      airPressure,
		LightValue:       lightValue,
		CarbonDioxide:    carbonDioxide,
		UVIndex:          uvIndex,
		Radiation:        radiation,
		WindSpeed:        windSpeed,
		CurrentRainfall:  currentRainfall,
		Temperature:      temperature,
		SoilMoisture:     soilMoisture,
		SoilConductivity: soilConductivity,
		SoilPH:           soilPH,
		SoilTemperature:  soilTemperature,
		WaterLevel:       waterLevel,
	}
}

func StartSmartAgricultureMQTTClient(broker string, port int, username string, password string, clientId string) {
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
				data := generateRandomSmartAgricultureData()
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
