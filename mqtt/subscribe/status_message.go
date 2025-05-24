package subscribe

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	initialize "github.com/Thingsly/backend/initialize"
	dal "github.com/Thingsly/backend/internal/dal"
	"github.com/Thingsly/backend/internal/model"
	service "github.com/Thingsly/backend/internal/service"
	"github.com/Thingsly/backend/pkg/global"

	"github.com/sirupsen/logrus"
)

func validateStatus(payload []byte) (int16, error) {
	str := string(payload)
	switch str {
	case "0":
		return 0, nil
	case "1":
		return 1, nil
	default:
		return 0, fmt.Errorf("The status value can only be 0 or 1, current value: %s", str)
	}
}

func DeviceOnline(payload []byte, topic string) {
	logrus.Debugf("Received status message for topic %s: %s", topic, string(payload))

	status, err := validateStatus(payload)
	if err != nil {
		logrus.Errorf("Invalid status for topic %s: %v", topic, err)
		return
	}

	deviceId := strings.Split(topic, "/")[2]
	logrus.Debugf("Processing status update for device %s: %d", deviceId, status)

	err = dal.UpdateDeviceStatus(deviceId, status)
	if err != nil {
		logrus.Errorf("Failed to update device status for device %s: %v", deviceId, err)
		return
	}

	if status == int16(1) {
		logrus.Debugf("Device %s is online, sending expected data", deviceId)
		time.Sleep(3 * time.Second)
		err := service.GroupApp.ExpectedData.Send(context.Background(), deviceId)
		if err != nil {
			logrus.Errorf("Failed to send expected data for device %s: %v", deviceId, err)
		}
	}

	initialize.DelDeviceCache(deviceId)

	var device *model.Device
	device, err = dal.GetDeviceCacheById(deviceId)
	if err != nil {
		logrus.Errorf("Failed to get device cache for device %s: %v", deviceId, err)
		return
	}

	go toUserClient(device, status)

	go func() {
		var loginStatus string
		if status == 1 {
			loginStatus = "ON-LINE"
		} else {
			loginStatus = "OFF-LINE"
		}
		logrus.Debugf("Executing automation for device %s with status %s", deviceId, loginStatus)
		err := service.GroupApp.Execute(device, service.AutomateFromExt{
			TriggerParamType: model.TRIGGER_PARAM_TYPE_STATUS,
			TriggerParam:     []string{},
			TriggerValues: map[string]interface{}{
				"login": loginStatus,
			},
		})
		if err != nil {
			logrus.Errorf("Automation execution failed for device %s: %v", deviceId, err)
		}
	}()

	err = initialize.SetRedisForJsondata(deviceId, device, 0)
	if err != nil {
		logrus.Errorf("Failed to set Redis JSON data for device %s: %v", deviceId, err)
		return
	}

	logrus.Debugf("Successfully processed status update for device %s", deviceId)
}

func toUserClient(device *model.Device, status int16) {
	var deviceName string
	sseEvent := global.SSEEvent{
		Type:     "device_online",
		TenantID: device.TenantID,
	}

	if device.Name != nil {
		deviceName = *device.Name
	} else {
		deviceName = device.DeviceNumber
	}

	logrus.Debugf("Sending status update to user client for device %s (%s): %d", device.DeviceNumber, deviceName, status)

	if status == int16(1) {
		jsonBytes, err := json.Marshal(map[string]interface{}{
			"device_id":   device.DeviceNumber,
			"device_name": deviceName,
			"is_online":   true,
			"timestamp":   time.Now().Unix(),
		})
		if err != nil {
			logrus.Errorf("Failed to marshal online status for device %s: %v", device.DeviceNumber, err)
			return
		}
		sseEvent.Message = string(jsonBytes)
	} else {
		jsonBytes, err := json.Marshal(map[string]interface{}{
			"device_id":   device.DeviceNumber,
			"device_name": deviceName,
			"is_online":   false,
			"timestamp":   time.Now().Unix(),
		})
		if err != nil {
			logrus.Errorf("Failed to marshal offline status for device %s: %v", device.DeviceNumber, err)
			return
		}
		sseEvent.Message = string(jsonBytes)
	}

	global.TPSSEManager.BroadcastEventToTenant(device.TenantID, sseEvent)
	logrus.Debugf("Successfully sent status update to user client for device %s", device.DeviceNumber)
}
