package subscribe

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Thingsly/backend/internal/dal"
	"github.com/Thingsly/backend/internal/model"
	"github.com/Thingsly/backend/pkg/global"

	"github.com/sirupsen/logrus"
)

func HeartbeatDeal(device *model.Device) {
	logrus.Debugf("Processing heartbeat for device %s", device.ID)

	if device.DeviceConfigID == nil {
		logrus.Errorf("Device %s has no config ID", device.ID)
		return
	}

	deviceConfig, err := dal.GetDeviceConfigByID(*device.DeviceConfigID)
	if err != nil {
		logrus.Errorf("Failed to get device config for device %s: %v", device.ID, err)
		return
	}

	if deviceConfig.OtherConfig == nil {
		logrus.Errorf("Device %s has no other config", device.ID)
		return
	}

	type OtherConfig struct {
		OnlineTimeout int `json:"online_timeout"`
		Heartbeat     int `json:"heartbeat"`
	}

	var otherConfig OtherConfig
	err = json.Unmarshal([]byte(*deviceConfig.OtherConfig), &otherConfig)
	if err != nil {
		logrus.Errorf("Failed to unmarshal other config for device %s: %v", device.ID, err)
		return
	}

	if otherConfig.Heartbeat > 0 {
		heartbeatKey := fmt.Sprintf("device:%s:heartbeat", device.ID)
		hasHeartbeat, err := global.STATUS_REDIS.Get(context.Background(), heartbeatKey).Result()
		if err != nil && err.Error() != "redis: nil" {
			logrus.Errorf("Failed to get heartbeat key for device %s: %v", device.ID, err)
			return
		}

		if hasHeartbeat == "1" {
			if device.IsOnline != int16(1) {
				logrus.Debugf("Updating device %s to online status", device.ID)
				DeviceOnline([]byte("1"), "devices/status/"+device.ID)
			}
		}

		err = global.STATUS_REDIS.Set(context.Background(),
			heartbeatKey,
			1,
			time.Duration(otherConfig.Heartbeat)*time.Second,
		).Err()
		if err != nil {
			logrus.Errorf("Failed to set heartbeat key for device %s: %v", device.ID, err)
			return
		}

		logrus.Debugf("Successfully processed heartbeat for device %s", device.ID)
		return
	}

	if otherConfig.OnlineTimeout > 0 {
		timeoutKey := fmt.Sprintf("device:%s:timeout", device.ID)
		hasTimeout, err := global.STATUS_REDIS.Get(context.Background(), timeoutKey).Result()
		if err != nil && err.Error() != "redis: nil" {
			logrus.Errorf("Failed to get timeout key for device %s: %v", device.ID, err)
			return
		}

		if hasTimeout == "1" {
			if device.IsOnline != int16(1) {
				logrus.Debugf("Updating device %s to online status", device.ID)
				DeviceOnline([]byte("1"), "devices/status/"+device.ID)
			}
		}

		err = global.STATUS_REDIS.Set(context.Background(),
			timeoutKey,
			1,
			time.Duration(otherConfig.OnlineTimeout)*time.Second,
		).Err()
		if err != nil {
			logrus.Errorf("Failed to set timeout key for device %s: %v", device.ID, err)
			return
		}

		logrus.Debugf("Successfully processed timeout for device %s", device.ID)
	}
}
