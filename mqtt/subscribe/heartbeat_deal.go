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
	if device.DeviceConfigID == nil {
		return
	}

	deviceConfig, err := dal.GetDeviceConfigByID(*device.DeviceConfigID)
	if err != nil {
		return
	}

	if deviceConfig.OtherConfig == nil {
		return
	}

	type OtherConfig struct {
		OnlineTimeout int `json:"online_timeout"`
		Heartbeat     int `json:"heartbeat"`
	}

	var otherConfig OtherConfig
	err = json.Unmarshal([]byte(*deviceConfig.OtherConfig), &otherConfig)
	if err != nil {
		return
	}

	if otherConfig.Heartbeat > 0 {
		heartbeatKey := fmt.Sprintf("device:%s:heartbeat", device.ID)
		hasHeartbeat, err := global.STATUS_REDIS.Get(context.Background(), heartbeatKey).Result()
		if err == nil && hasHeartbeat == "1" {
			if device.IsOnline != int16(1) {
				DeviceOnline([]byte("1"), "devices/status/"+device.ID)
			}
		}

		err = global.STATUS_REDIS.Set(context.Background(),
			heartbeatKey,
			1,
			time.Duration(otherConfig.Heartbeat)*time.Second,
		).Err()
		if err != nil {
			logrus.Error(err)
			return
		}

		return
	}

	if otherConfig.OnlineTimeout > 0 {
		timeoutKey := fmt.Sprintf("device:%s:timeout", device.ID)
		hasTimeout, err := global.STATUS_REDIS.Get(context.Background(), timeoutKey).Result()
		if err == nil && hasTimeout == "1" {
			if device.IsOnline != int16(1) {
				DeviceOnline([]byte("1"), "devices/status/"+device.ID)
			}
		}

		err = global.STATUS_REDIS.Set(context.Background(),
			timeoutKey,
			1,
			time.Duration(otherConfig.OnlineTimeout)*time.Second,
		).Err()
		if err != nil {
			logrus.Error(err)
			return
		}
	}
}
