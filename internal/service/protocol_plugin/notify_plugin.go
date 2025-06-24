package protocolplugin

import (
	"encoding/json"
	"fmt"

	"github.com/Thingsly/backend/internal/dal"
	"github.com/Thingsly/backend/third_party/others/http_client"

	"github.com/sirupsen/logrus"
)

func DeviceConfigUpdateAndDisconnect(deviceConfigID string, protocolType string, deviceType string) error {

	servicePlugin, err := dal.GetServicePluginByServiceIdentifier(protocolType)
	if err != nil {
		return err
	}

	_, host, err := dal.GetServicePluginHttpAddressByID(servicePlugin.ID)
	if err != nil {
		return err
	}

	switch deviceType {
	case "3":
		deviceIDs, err := dal.GetGatewayDevicesBySubDeviceConfigID(deviceConfigID)
		if err != nil {
			return err
		}
		for _, deviceID := range deviceIDs {
			DisconnectDevice(deviceID, host)
		}
	case "1", "2":
		devices, err := dal.GetDevicesByDeviceConfigID(deviceConfigID)
		if err != nil {
			return err
		}
		for _, device := range devices {
			DisconnectDevice(device.ID, host)
		}
		return nil
	}
	return nil
}

func DisconnectDevice(deviceID string, httpAddress string) error {
	type ReqData struct {
		DeviceID string `json:"device_id"`
	}
	reqData := ReqData{DeviceID: deviceID}
	reqDataBytes, err := json.Marshal(reqData)
	if err != nil {
		return err
	}
	rsp, err := http_client.DisconnectDevice(reqDataBytes, httpAddress)
	if err != nil {
		logrus.Warnf("update succeeded, but connect plugin failed: %s", err)
		return err
	}

	var rspData http_client.RspData
	err = json.NewDecoder(rsp.Body).Decode(&rspData)
	if err != nil {
		logrus.Warnf("update succeeded, but plugin rspdata decode failed: %s", err)
		return err
	}
	if rspData.Code != 200 {
		logrus.Warnf("update succeeded, but plugin rsp: %s", rspData.Message)
		return err
	}
	return nil
}

func DisconnectDeviceByDeviceID(deviceID string) error {

	device, err := dal.GetDeviceByID(deviceID)
	if err != nil {
		return err
	}
	if device.DeviceConfigID == nil {
		return nil
	}

	deviceConfig, err := dal.GetDeviceConfigByID(*device.DeviceConfigID)
	if err != nil {
		return err
	}
	if deviceConfig == nil {
		return nil
	}
	if deviceConfig.ProtocolType == nil {
		return fmt.Errorf("protocol type not found")
	}
	if *deviceConfig.ProtocolType == "MQTT" {
		return nil
	}

	servicePlugin, err := dal.GetServicePluginByServiceIdentifier(*deviceConfig.ProtocolType)
	if err != nil {
		return err
	}

	_, host, err := dal.GetServicePluginHttpAddressByID(servicePlugin.ID)
	if err != nil {
		return err
	}

	if deviceConfig.DeviceType == "3" {
		err = DisconnectDevice(*device.ParentID, host)
		if err != nil {
			return err
		}
	} else {
		err = DisconnectDevice(deviceID, host)
		if err != nil {
			return err
		}
	}
	return nil
}
