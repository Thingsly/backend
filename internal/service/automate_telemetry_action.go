package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/Thingsly/backend/initialize"
	"github.com/Thingsly/backend/internal/dal"
	model "github.com/Thingsly/backend/internal/model"
	"github.com/Thingsly/backend/pkg/constant"

	"github.com/sirupsen/logrus"
)

const (
	AUTOMATE_ACTION_PARAM_TYPE_TEL          = "TEL"
	AUTOMATE_ACTION_PARAM_TYPE_TELEMETRY    = "telemetry"
	AUTOMATE_ACTION_PARAM_TYPE_C_TELEMETRY  = "c_telemetry"
	AUTOMATE_ACTION_PARAM_TYPE_ATTR         = "ATTR"
	AUTOMATE_ACTION_PARAM_TYPE_ATTRIBUTES   = "attributes"
	AUTOMATE_ACTION_PARAM_TYPE_C_ATTRIBUTES = "c_attribute"
	AUTOMATE_ACTION_PARAM_TYPE_CMD          = "CMD"
	AUTOMATE_ACTION_PARAM_TYPE_COMMAND      = "command"
	AUTOMATE_ACTION_PARAM_TYPE_C_COMMAND    = "c_command"
)

type AutomateTelemetryAction interface {
	AutomateActionRun(model.ActionInfo) (string, error)
}

func AutomateActionDeviceMqttSend(deviceId string, action model.ActionInfo, tenantID string) (string, error) {

	var executeMsg string
	// Get device cache information
	deviceInfo, err := initialize.GetDeviceCacheById(deviceId)
	if err != nil {
		executeMsg = fmt.Sprintf("Device ID: %s", deviceId)
	} else {
		executeMsg = fmt.Sprintf("Device Name: %s", *deviceInfo.Name)
	}

	if action.ActionParamType == nil {
		return executeMsg + " ActionParamType does not exist ", errors.New("ActionParamType does not exist")
	}
	if action.ActionValue == nil {
		return executeMsg + " Action target value does not exist ", errors.New("action target value does not exist")
	}
	// if action.ActionParam == nil {
	// 	return executeMsg + " Identifier does not exist", errors.New("Identifier does not exist")
	// }
	ctx := context.Background()

	var userId string
	userId, _ = dal.GetUserIdBYTenantID(tenantID)
	logrus.Debug("AutomateActionDeviceMqttSend:", tenantID, ", userId:", userId)
	operationType := strconv.Itoa(constant.Auto)
	//var valueMap = make(map[string]string)
	switch *action.ActionParamType {
	// Case 1: Telemetry
	case AUTOMATE_ACTION_PARAM_TYPE_TEL, AUTOMATE_ACTION_PARAM_TYPE_TELEMETRY, AUTOMATE_ACTION_PARAM_TYPE_C_TELEMETRY:
		msgReq := model.PutMessage{
			DeviceID: deviceId,
		}
		//valueMap = map[string]string{
		//	*action.ActionParam: *action.ActionValue,
		//}
		//valueStr, _ := json.Marshal(valueMap)
		//msgReq.Value = string(valueStr)
		msgReq.Value = *action.ActionValue
		logrus.Warning(msgReq)
		return executeMsg + fmt.Sprintf(" Telemetry command: %s", msgReq.Value), GroupApp.TelemetryData.TelemetryPutMessage(ctx, userId, &msgReq, operationType)

	// Case 2: Attribute
	case AUTOMATE_ACTION_PARAM_TYPE_ATTR, AUTOMATE_ACTION_PARAM_TYPE_ATTRIBUTES, AUTOMATE_ACTION_PARAM_TYPE_C_ATTRIBUTES:
		msgReq := model.AttributePutMessage{
			DeviceID: deviceId,
		}
		//valueMap = map[string]string{
		//	*action.ActionParam: *action.ActionValue,
		//}
		//valueStr, _ := json.Marshal(valueMap)
		//msgReq.Value = string(valueStr)
		msgReq.Value = *action.ActionValue
		return executeMsg + fmt.Sprintf(" Property setting: %s", msgReq.Value), GroupApp.AttributeData.AttributePutMessage(ctx, userId, &msgReq, operationType)

	// Case 3: Command
	case AUTOMATE_ACTION_PARAM_TYPE_CMD, AUTOMATE_ACTION_PARAM_TYPE_COMMAND, AUTOMATE_ACTION_PARAM_TYPE_C_COMMAND:
		type commandInfo struct {
			Method string      `json:"method"`
			Params interface{} `json:"params"`
		}
		var info commandInfo
		err := json.Unmarshal([]byte(*action.ActionValue), &info)
		if err != nil {
			return executeMsg + "Command dispatch data parsing failed", err
		}
		value, _ := json.Marshal(info.Params)
		valueStr := string(value)
		msgReq := model.PutMessageForCommand{
			DeviceID: deviceId,
			Value:    &valueStr,
			Identify: info.Method,
		}
		//msgReq := model.PutMessageForCommand{
		//	DeviceID: deviceId,
		//	Value:    action.ActionValue,
		//	Identify: *action.ActionParam,
		//}
		return executeMsg + fmt.Sprintf("Command dispatch: %s", *msgReq.Value), GroupApp.CommandData.CommandPutMessage(ctx, userId, &msgReq, operationType)
	default:

		return executeMsg + " unsupported type", errors.New("unsupported type")
	}
}

type AutomateTelemetryActionOne struct {
	TenantID string
}

func (a *AutomateTelemetryActionOne) AutomateActionRun(action model.ActionInfo) (string, error) {

	if action.ActionTarget == nil {
		return "Single device execution, device ID does not exist", errors.New("Device ID does not exist")
	}
	return AutomateActionDeviceMqttSend(*action.ActionTarget, action, a.TenantID)
}

// Service 10: Multiple devices execution
type AutomateTelemetryActionMultiple struct {
	DeviceIds []string
	TenantID  string
}

func (a *AutomateTelemetryActionMultiple) AutomateActionRun(action model.ActionInfo) (string, error) {

	var (
		messages []string
		errs     error
	)
	for _, deviceId := range a.DeviceIds {
		msg, err := AutomateActionDeviceMqttSend(deviceId, action, a.TenantID)
		if err != nil && errs == nil {
			errs = err
		}
		messages = append(messages, msg)
	}

	return "Single-type setting: " + fmt.Sprintf("%s", messages), errs
}

// Service 20: Scene activation
type AutomateTelemetryActionScene struct {
	TenantID string
}

func (a *AutomateTelemetryActionScene) AutomateActionRun(action model.ActionInfo) (string, error) {

	if action.ActionTarget == nil {
		return "Scene activation", errors.New("Scene ID does not exist")
	}
	// Retrieve scene information
	sceneInfo, err := dal.GetSceneInfo(*action.ActionTarget)
	if err != nil {
		return "Scene activation", err
	}
	return fmt.Sprintf("Scene activation: %s", sceneInfo.Name), GroupApp.ActiveSceneExecute(*action.ActionTarget, a.TenantID)
}

// Service 30: Alarm service execution
type AutomateTelemetryActionAlarm struct{}

func (*AutomateTelemetryActionAlarm) AutomateActionRun(action model.ActionInfo) (string, error) {

	logrus.Debugf("Alarm service: %#v", *action.ActionTarget)
	// The alarm service has a decorator implementation, so no further handling is needed here
	if action.ActionTarget == nil || *action.ActionTarget == "" {
		return "Alarm service", errors.New("Alarm ID does not exist")
	}

	ok, alarmName, reason := AlarmExecute(*action.ActionTarget, action.SceneAutomationID)
	if ok {
		return fmt.Sprintf("Alarm service(%s)", alarmName), nil
	}
	// Retrieve alarm name; display it regardless of execution success
	alarmName = dal.GetAlarmNameWithCache(*action.ActionTarget)

	// Handle case where the alarm already exists – do not treat as an error
	if reason == "Alarm already exists" {
		logrus.Debugf("Alarm (%s) already exists, skipping re-trigger", alarmName)
		return fmt.Sprintf("Alarm service (%s) - already exists, skipping re-trigger", alarmName), nil
	}

	// Other failure cases
	errRsp := errors.New("Execution failed, " + reason)
	return fmt.Sprintf("Alarm service (%s)", alarmName), errRsp
}

// Service 40: Automation another service execution (to be implemented)
type AutomateTelemetryActionService struct{}

func (*AutomateTelemetryActionService) AutomateActionRun(_ model.ActionInfo) (string, error) {
	//todo To be implemented
	fmt.Println("Automation service action implementation")
	return "Service", nil
}
