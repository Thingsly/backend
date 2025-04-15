package mqtt

import (
	"encoding/json"
	"fmt"

	"github.com/Thingsly/backend/internal/model"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type GatewayResponseFunc = func(model.GatewayResponse) error

type MqttDirectResponseFunc = func(response model.MqttResponse) error

var MqttConfig Config

var GatewayResponseFuncMap map[string]chan model.GatewayResponse

var MqttDirectResponseFuncMap map[string]chan model.MqttResponse

func MqttInit() error {

	err := loadConfig()
	if err != nil {
		return err
	}

	GatewayResponseFuncMap = make(map[string]chan model.GatewayResponse)
	MqttDirectResponseFuncMap = make(map[string]chan model.MqttResponse)
	return nil
}

func loadConfig() error {
	var configMap map[string]interface{}

	err := viper.Unmarshal(&configMap)
	if err != nil {
		return fmt.Errorf("unable to decode into struct, %s", err)
	}

	jsonStr, err := json.Marshal(configMap["mqtt"])
	if err != nil {
		return fmt.Errorf("unable to marshal config, %s", err)
	}

	err = json.Unmarshal(jsonStr, &MqttConfig)
	if err != nil {
		return fmt.Errorf("unable to unmarshal config, %s", err)
	}

	logrus.Debug("mqtt config:", MqttConfig)

	broker := viper.GetString("mqtt.broker")
	if broker == "" {
		broker = "localhost:1883"
		logrus.Println("Using default broker:", broker)
	}
	MqttConfig.Broker = broker

	user := viper.GetString("mqtt.user")
	if user == "" {
		user = "root"
		logrus.Println("Using default user:", user)
	}
	MqttConfig.User = user

	pass := viper.GetString("mqtt.pass")
	if pass == "" {
		pass = "root"
		logrus.Println("Using default pass:", pass)
	}
	MqttConfig.Pass = pass

	channelBufferSize := viper.GetInt("mqtt.channel_buffer_size")
	if channelBufferSize == 0 {
		channelBufferSize = 10000
		logrus.Println("Using default channel_buffer_size:", channelBufferSize)
	}
	MqttConfig.ChannelBufferSize = channelBufferSize

	writeWorkers := viper.GetInt("mqtt.write_workers")
	if writeWorkers == 0 {
		writeWorkers = 10
		logrus.Println("Using default write_workers:", writeWorkers)
	}
	MqttConfig.WriteWorkers = writeWorkers

	poolSize := viper.GetInt("mqtt.telemetry.pool_size")
	if poolSize == 0 {
		poolSize = 100
		logrus.Println("Using default pool_size:", poolSize)
	}
	MqttConfig.Telemetry.PoolSize = poolSize

	batchSize := viper.GetInt("mqtt.telemetry.batch_size")
	if batchSize == 0 {
		batchSize = 100
		logrus.Println("Using default batch_size:", batchSize)
	}
	return nil
}
