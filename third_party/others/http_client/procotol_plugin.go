package http_client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Thingsly/backend/pkg/errcode"

	"github.com/sirupsen/logrus"
)

type RspData struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type RspDeviceListData struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Data    ListData `json:"data"`
}
type ListData struct {
	Total int          `json:"total"`
	List  []DeviceData `json:"list"`
}
type DeviceData struct {
	DeviceName     string `json:"device_name"`
	DeviceNumber   string `json:"device_number"`
	Description    string `json:"description"`
	IsBind         bool   `json:"is_bind"`
	DeviceConfigID string `json:"device_config_id"`
}

// Get the plugin's form configuration
// CONFIG - Configuration form, VOUCHER - Voucher form, VOUCHER-TYPE - Voucher type form
// func GetPluginFromConfig(host string, protocol_type string, device_type string, form_type string, voucher_type string) ([]byte, error) {
//     return Get("http://" + host + "/api/v1/form/config?protocol_type=" + protocol_type + "&device_type=" + device_type + "&form_type=" + form_type + "&voucher_type=" + voucher_type)
// }

// /api/v2/form/config
// CFG - Configuration form, VCR - Voucher form, VCRT - Voucher type form, SVCRT - Service voucher form
func GetPluginFromConfigV2(host string, service_identifier string, device_type string, form_type string) (interface{}, error) {
	// Send GET request to retrieve form configuration
	b, err := Get("http://" + host + "/api/v1/form/config?protocol_type=" + service_identifier + "&device_type=" + device_type + "&form_type=" + form_type)
	if err != nil {
		logrus.Error(err)
		if err.Error() != "" && (strings.Contains(err.Error(), "connection refused")) {
			return nil, errcode.WithData(200068, err.Error())
		}
		return nil, errcode.WithData(200069, err.Error())
	}

	// Parse the response
	var rspdata RspData
	err = json.Unmarshal(b, &rspdata)
	if err != nil {
		logrus.Error(err)
		return nil, errcode.WithData(200070, err.Error())
	}

	// Check if the response code is 200
	if rspdata.Code != 200 {
		err = errcode.NewWithMessage(200070, rspdata.Message)
		logrus.Error(err)
	}

	return rspdata.Data, nil
}

// Disconnect a device to allow it to reconnect
func DisconnectDevice(reqdata []byte, host string) (*http.Response, error) {
	// Send a POST request to disconnect the device
	return PostJson("http://"+host+"/api/v1/device/disconnect", reqdata)
}

// Delete device or sub-device notification (device protocol change is also considered a deletion)
func DeleteDevice(reqdata []byte, host string) (*http.Response, error) {
	// Send a POST request to delete the device
	return PostJson("http://"+host+"/api/v1/device/delete", reqdata)
}

// Device or sub-device configuration change notification
func UpdateDeviceConfig(reqdata []byte, host string) (*http.Response, error) {
	// Send a POST request to update the device configuration
	return PostJson("http://"+host+"/api/v1/device/config/update", reqdata)
}

// Add device or sub-device notification (device protocol change is also considered an addition)
func AddDevice(reqdata []byte, host string) (*http.Response, error) {
	// Send a POST request to add the device
	return PostJson("http://"+host+"/api/v1/device/add", reqdata)
}

// messageType 1 - Service Configuration Modification
func Notification(messageType string, message string, host string) ([]byte, error) {
	// Define the request data structure
	type ReqData struct {
		MessageType string `json:"message_type"`
		Message     string `json:"message"`
	}

	// Prepare the request data
	reqData := ReqData{MessageType: messageType, Message: message}
	reqDataBytes, err := json.Marshal(reqData)
	if err != nil {
		return nil, err
	}

	// Send the POST request
	response, err := PostJson("http://"+host+"/api/v1/notify/event", reqDataBytes)
	if err != nil {
		logrus.Error(err)
		return nil, fmt.Errorf("post plugin notification failed: %s", err)
	}

	// Check if the response status is 200 OK
	if response.StatusCode != 200 {
		err = fmt.Errorf("protocol plugin response message: %s", response.Status)
		logrus.Error(err)
		return nil, err
	}

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		logrus.Error(err)
		return nil, fmt.Errorf("read plugin response body failed: %s", err)
	}
	logrus.Info(string(body))

	return body, nil
}

// /api/v1/service/access/device/list
// Third-party service device list query
func GetServiceAccessDeviceList(host string, voucher string, page_size string, page string) (*ListData, error) {
	// Send GET request to retrieve device list
	b, err := Get("http://" + host + "/api/v1/plugin/device/list?voucher=" + voucher + "&page_size=" + page_size + "&page=" + page)
	if err != nil {
		logrus.Error(err)
		logrus.Error("http://" + host + "/api/v1/plugin/device/list?voucher=" + voucher + "&page_size=" + page_size + "&page=" + page)
		return nil, fmt.Errorf("get plugin form failed: %s", err)
	}

	// Parse the response
	var rspdata RspDeviceListData
	err = json.Unmarshal(b, &rspdata)
	if err != nil {
		logrus.Error(err)
		return nil, fmt.Errorf("unmarshal response data failed: %s", err)
	}

	// Check if the response code is 200 OK
	if rspdata.Code != 200 {
		err = fmt.Errorf("protocol plugin response message: %s", rspdata.Message)
		logrus.Error(err)
		return nil, err
	}

	// If rspdata.Data.List is nil, return an empty array
	if rspdata.Data.List == nil {
		rspdata.Data.List = []DeviceData{}
	}
	return &rspdata.Data, nil
}
