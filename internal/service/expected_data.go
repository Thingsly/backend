package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/HustIoTPlatform/backend/internal/dal"
	model "github.com/HustIoTPlatform/backend/internal/model"
	"github.com/HustIoTPlatform/backend/pkg/errcode"
	utils "github.com/HustIoTPlatform/backend/pkg/utils"

	"github.com/go-basic/uuid"
	"github.com/sirupsen/logrus"
)

type ExpectedData struct{}

func mergeIdentifyAndPayload(identify string, paramsStr *string) (string, error) {

	mergedData := map[string]interface{}{
		"method": identify,
	}

	if paramsStr != nil {
		var params any
		err := json.Unmarshal([]byte(*paramsStr), &params)
		if err != nil {
			return "", fmt.Errorf("error parsing payload JSON: %v", err)
		}
		mergedData["params"] = params
	}

	mergedJSON, err := json.Marshal(mergedData)
	if err != nil {
		return "", fmt.Errorf("error marshaling merged data to JSON: %v", err)
	}

	return string(mergedJSON), nil
}

func (e *ExpectedData) Create(ctx context.Context, req *model.CreateExpectedDataReq, userClaims *utils.UserClaims) (*model.ExpectedData, error) {
	if req.SendType == "command" {
		if req.Identify == nil {
			return nil, errcode.WithData(errcode.CodeParamError, map[string]interface{}{
				"identify": "identify is required",
			})
		}

		payload, err := mergeIdentifyAndPayload(*req.Identify, req.Payload)
		if err != nil {
			return nil, errcode.WithData(errcode.CodeParamError, map[string]interface{}{
				"payload": err.Error(),
			})
		}
		req.Payload = &payload
	} else if req.Payload == nil {
		return nil, errcode.WithData(errcode.CodeParamError, map[string]interface{}{
			"payload": "payload is required",
		})
	}

	ed := &model.ExpectedData{
		ID:         uuid.New(),
		DeviceID:   req.DeviceID,
		SendType:   req.SendType,
		Payload:    *req.Payload,
		CreatedAt:  time.Now(),
		Status:     "pending",
		ExpiryTime: req.Expiry,
		Label:      req.Label,
		TenantID:   userClaims.TenantID,
	}
	err := dal.ExpectedDataDal{}.Create(ctx, ed)
	if err != nil {
		logrus.Error(err)
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}

	expectedData, err := dal.ExpectedDataDal{}.GetByID(ctx, ed.ID)
	if err != nil {
		logrus.Error(err)
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}

	deviceStatus, err := GroupApp.Device.GetDeviceOnlineStatus(req.DeviceID)
	if err != nil {
		logrus.Error(err)
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}
	if deviceStatus["is_online"] == 1 {

		err := e.Send(ctx, req.DeviceID)
		if err != nil {
			logrus.Error(err)
			return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
				"sql_error": err.Error(),
			})
		}
	}

	return expectedData, nil

}

func (*ExpectedData) Delete(ctx context.Context, id string) error {
	return dal.ExpectedDataDal{}.Delete(ctx, id)
}

func (*ExpectedData) PageList(ctx context.Context, req *model.GetExpectedDataPageReq, userClaims *utils.UserClaims) (map[string]interface{}, error) {
	total, list, err := dal.ExpectedDataDal{}.PageList(ctx, req, userClaims.TenantID)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}
	return map[string]interface{}{
		"total": total,
		"list":  list,
	}, nil
}

// Send processes and sends expected data for a specific device.
// It queries the expected data from the database, checks for expiration, 
// and sends the data based on its type (telemetry, attribute, command).
// It also updates the status of the data after sending.
func (*ExpectedData) Send(ctx context.Context, deviceID string) error {

	// Retrieve all expected data for the given device ID
	ed, err := dal.ExpectedDataDal{}.GetAllByDeviceID(ctx, deviceID)
	if err != nil {
		logrus.WithError(err).Error("Failed to query expected data")
		return err
	}
	logrus.WithField("deviceID", deviceID).Debug("Retrieved expected data", ed)

	// Loop through all expected data
	for _, v := range ed {
		// Check if the expected data has expired
		if v.ExpiryTime != nil && v.ExpiryTime.Before(time.Now()) {
			logrus.WithField("dataID", v.ID).Debug("Expected data has expired")
			// Update the status to 'expired' if the data is expired
			if err := updateStatus(ctx, v.ID, "expired", nil); err != nil {
				return err
			}
			continue
		}

		// Define the default status and message variables
		var (
			status  = "sent"
			message string
		)

		// Handle different data send types (telemetry, attribute, command)
		switch v.SendType {
		case "telemetry":
			message, err = sendTelemetry(ctx, deviceID, v.Payload)
		case "attribute":
			message, err = sendAttribute(ctx, deviceID, v.Payload)
		case "command":
			message, err = sendCommand(ctx, deviceID, v.Payload)
		default:
			// Handle unknown send type
			logrus.WithField("sendType", v.SendType).Error("Unknown send type")
			continue
		}

		// If an error occurred while sending data, update the status to 'expired'
		if err != nil {
			status = "expired"
			logrus.WithError(err).WithField("sendType", v.SendType).Error("Failed to send data")
		}

		// Update the status in the database after attempting to send data
		if err := updateStatus(ctx, v.ID, status, &message); err != nil {
			return err
		}
	}

	return nil
}

// sendTelemetry sends expected telemetry data to the device.
func sendTelemetry(ctx context.Context, deviceID, payload string) (string, error) {
	logrus.Debug("Sending expected telemetry data")

	// Create a PutMessage struct for telemetry data
	putMessage := &model.PutMessage{
		DeviceID: deviceID,
		Value:    payload,
	}

	// Send telemetry data using the TelemetryPutMessage function
	err := GroupApp.TelemetryData.TelemetryPutMessage(ctx, "", putMessage, "2")
	if err != nil {
		// If sending fails, return the error message
		return err.Error(), err
	}

	// Return success message if sending is successful
	return "send success", nil
}

// sendAttribute sends expected attribute data to the device.
func sendAttribute(ctx context.Context, deviceID, payload string) (string, error) {
	logrus.Debug("Sending expected attribute data")

	// Create a PutMessage struct for attribute data
	putMessage := &model.AttributePutMessage{
		DeviceID: deviceID,
		Value:    payload,
	}

	// Send attribute data using the AttributePutMessage function
	err := GroupApp.AttributeData.AttributePutMessage(ctx, "", putMessage, "2")
	if err != nil {
		// If sending fails, return the error message
		return err.Error(), err
	}

	// Return success message if sending is successful
	return "send success", nil
}

// sendCommand sends an expected command to the device.
func sendCommand(ctx context.Context, deviceID, payload string) (string, error) {
	logrus.Debug("Sending expected command data")

	// Parse the payload JSON into a map
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(payload), &data); err != nil {
		// If parsing fails, return an error message
		return fmt.Sprintf("Error parsing JSON payload: %s", err.Error()), err
	}

	// Retrieve the "method" field from the JSON
	method, ok := data["method"].(string)
	if !ok {
		// If "method" is missing, return an error
		return "identify is required", fmt.Errorf("identify is required")
	}

	// Optional: If "params" is provided in the payload, convert it to a string
	var paramsStr *string
	if params, exists := data["params"]; exists {
		paramsJSON, err := json.Marshal(params)
		if err != nil {
			// If marshalling params fails, return an error message
			return fmt.Sprintf("Error converting params to string: %s", err.Error()), err
		}
		p := string(paramsJSON)
		paramsStr = &p
	}

	// Create a PutMessageForCommand struct for sending the command
	putMessage := &model.PutMessageForCommand{
		DeviceID: deviceID,
		Identify: method,
		Value:    paramsStr,
	}

	// Send the command using the CommandPutMessage function
	err := GroupApp.CommandData.CommandPutMessage(ctx, "", putMessage, "2")
	if err != nil {
		// If sending fails, return the error message
		return err.Error(), err
	}

	// Return success message if sending is successful
	return "send success", nil
}

func updateStatus(ctx context.Context, id string, status string, message *string) error {
	var sendTime time.Time
	if status == "sent" {
		sendTime = time.Now()
	}

	err := dal.ExpectedDataDal{}.UpdateStatus(ctx, id, status, message, &sendTime)
	if err != nil {
		logrus.WithError(err).WithField("dataID", id).Error("Failed to update the expected data status")
	}
	return err
}
