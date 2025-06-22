package initialize

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Thingsly/backend/internal/dal"
	model "github.com/Thingsly/backend/internal/model"
	global "github.com/Thingsly/backend/pkg/global"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

func RedisInit() (*redis.Client, error) {
	conf, err := loadConfig()
	if err != nil {
		return nil, fmt.Errorf("Failed to load Redis configuration: %v", err)
	}

	statusConf, err := loadStatusConfig()
	if err != nil {
		return nil, fmt.Errorf("Failed to load Redis configuration: %v", err)
	}

	client := connectRedis(conf)
	statusClient := connectRedis(statusConf)

	if checkRedisClient(client) != nil {
		return nil, fmt.Errorf("Failed to connect to Redis: %v", err)
	}
	if checkRedisClient(statusClient) != nil {
		return nil, fmt.Errorf("Failed to connect to Redis: %v", err)
	}
	global.REDIS = client
	global.STATUS_REDIS = statusClient
	// Start SSE
	go global.InitSSEManager()
	return client, nil
}

func connectRedis(conf *RedisConfig) *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     conf.Addr,
		Password: conf.Password,
		DB:       conf.DB,
	})
	// If nil is returned, create this DB

	return redisClient
}

func checkRedisClient(redisClient *redis.Client) error {
	// Check if the connection to the redis server is successful by calling client.Ping()
	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		return err
	} else {
		log.Println("Connected to Redis successfully...")
		return nil
	}
}

func loadConfig() (*RedisConfig, error) {
	redisConfig := &RedisConfig{
		Addr:     viper.GetString("db.redis.addr"),
		Password: viper.GetString("db.redis.password"),
		DB:       viper.GetInt("db.redis.db"),
	}

	if redisConfig.Addr == "" {
		redisConfig.Addr = "localhost:6379"
	}
	return redisConfig, nil
}

func loadStatusConfig() (*RedisConfig, error) {
	db := viper.GetInt("db.redis.db1")
	if db == 0 {
		db = 10 // Default to the 11th DB 
	}
	redisConfig := &RedisConfig{
		Addr:     viper.GetString("db.redis.addr"),
		Password: viper.GetString("db.redis.password"),
		DB:       db,
	}

	if redisConfig.Addr == "" {
		redisConfig.Addr = "localhost:6379"
	}
	return redisConfig, nil
}

// setRedis Serialize map or struct object to JSON string and store in Redis
func SetRedisForJsondata(key string, value interface{}, expiration time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return global.REDIS.Set(context.Background(), key, jsonData, expiration).Err()
}

// getRedis Get JSON from Redis and deserialize to a specified object
func GetRedisForJsondata(key string, dest interface{}) error {
	val, err := global.REDIS.Get(context.Background(), key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

// Get device information from redis by device id
// First get device information from redis, if not, get device information from database, and then store the device information in redis
func GetDeviceCacheById(deviceId string) (*model.Device, error) {
	var device model.Device
	err := GetRedisForJsondata(deviceId, &device)
	if err == nil {
		return &device, nil
	}
	// Get device information from database
	deviceFromDB, err := dal.GetDeviceCacheById(deviceId)
	if err != nil {
		return nil, err
	}
	// Store device information in redis
	err = SetRedisForJsondata(deviceId, deviceFromDB, 0)
	if err != nil {
		return nil, err
	}
	return deviceFromDB, nil
}

// Get script from redis by device and script type
func GetScriptByDeviceAndScriptType(device *model.Device, script_type string) (*model.DataScript, error) {
	var script *model.DataScript
	script = &model.DataScript{}
	if device.DeviceConfigID == nil {
		return nil, fmt.Errorf("Device config id is empty")
	}
	key := *device.DeviceConfigID + "_" + script_type + "_script"
	err := GetRedisForJsondata(key, script)
	if err != nil {
		logrus.Debug("Get redis_cache key:"+key+" failed with err:", err.Error())
		script, err = dal.GetDataScriptByDeviceConfigIdAndScriptType(device.DeviceConfigID, script_type)
		if err != nil {
			return nil, err
		}
		if script == nil {
			return nil, nil
		}
		err = SetRedisForJsondata(key, script, 0)
		if err != nil {
			logrus.Debug("Set redis_cache key:"+key+" failed with err:", err.Error())
			return nil, err
		}
		logrus.Debug("Set redis_cache key:"+key+" successed with ", script)
	}
	return script, nil
}

// Clear device information cache
func DelDeviceCache(deviceId string) error {
	err := global.REDIS.Del(context.Background(), deviceId).Err()
	if err != nil {
		logrus.Warn("del redis_cache key(deviceId):", deviceId, " failed with err:", err.Error())
	}
	return err
}

// Clear device configuration information cache
func DelDeviceConfigCache(deviceConfigId string) error {
	err := global.REDIS.Del(context.Background(), deviceConfigId+"_config").Err()
	if err != nil {
		logrus.Warn("del redis_cache key(deviceConfigId):", deviceConfigId+"_config", " failed with err:", err.Error())
	}
	return err
}

// Clear device corresponding script cache
func DelDeviceDataScriptCache(deviceConfigID string) error {
	scriptType := []string{"A", "B", "C", "D", "E", "F"}
	var key []string
	for _, scriptType := range scriptType {
		key = append(key, deviceConfigID+"_"+scriptType+"_script")
	}

	err := global.REDIS.Del(context.Background(), key...).Err()
	if err != nil {
		logrus.Warn("del redis_cache key:", key, " failed with err:", err.Error())
	}
	return err
}
