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
		return 0, fmt.Errorf("the status value can only be 0 or 1, current value: %s", str)
	}
}

func DeviceOnline(payload []byte, topic string) {
	status, err := validateStatus(payload)
	if err != nil {
		logrus.Error(err.Error())
		return
	}

	deviceId := strings.Split(topic, "/")[2]
	logrus.Debug(deviceId, " device status message:", status)

	err = dal.UpdateDeviceStatus(deviceId, status)
	if err != nil {
		logrus.Error(err.Error())
		return
	}
	if status == int16(1) {

		time.Sleep(3 * time.Second)
		err := service.GroupApp.ExpectedData.Send(context.Background(), deviceId)
		if err != nil {
			logrus.Error(err.Error())
		}

	}

	initialize.DelDeviceCache(deviceId)

	var device *model.Device
	device, err = dal.GetDeviceCacheById(deviceId)
	if err != nil {
		logrus.Error(err.Error())
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
		err := service.GroupApp.Execute(device, service.AutomateFromExt{
			TriggerParamType: model.TRIGGER_PARAM_TYPE_STATUS,
			TriggerParam:     []string{},
			TriggerValues: map[string]interface{}{
				"login": loginStatus,
			},
		})
		if err != nil {
			logrus.Error("Automation execution failed, err: %w", err)
		}
	}()

	err = initialize.SetRedisForJsondata(deviceId, device, 0)
	if err != nil {
		logrus.Error(err.Error())
		return
	}

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
	if status == int16(1) {
		jsonBytes, _ := json.Marshal(map[string]interface{}{
			"device_id":   device.DeviceNumber,
			"device_name": deviceName,
			"is_online":   true,
		})
		sseEvent.Message = string(jsonBytes)
	} else {
		jsonBytes, _ := json.Marshal(map[string]interface{}{
			"device_id":   device.DeviceNumber,
			"device_name": deviceName,
			"is_online":   false,
		})
		sseEvent.Message = string(jsonBytes)
	}
	global.TLSSEManager.BroadcastEventToTenant(device.TenantID, sseEvent)
}
