package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Thingsly/backend/initialize"
	"github.com/Thingsly/backend/internal/dal"
	"github.com/Thingsly/backend/internal/model"
	"github.com/Thingsly/backend/pkg/common"

	"github.com/go-basic/uuid"
	pkgerrors "github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Automate struct {
	device  *model.Device
	formExt AutomateFromExt
	mu      sync.Mutex
}

var conditionAfterDecoration = []ConditionAfterFunc{
	ConditionAfterAlarm,
}

var actionAfterDecoration = []ActionAfterFunc{
	ActionAfterAlarm,
}

type ConditionAfterFunc = func(ok bool, conditions initialize.DTConditions, deviceId string, contents []string) error
type ActionAfterFunc = func(actions []model.ActionInfo, err error) error

type AutomateFromExt struct {
	TriggerParamType string
	TriggerParam     []string
	TriggerValues    map[string]interface{}
}

func (a *Automate) conditionAfterDecorationRun(ok bool, conditions initialize.DTConditions, deviceId string, contents []string) {
	defer a.ErrorRecover()
	for _, fc := range conditionAfterDecoration {
		err := fc(ok, conditions, deviceId, contents)
		if err != nil {
			logrus.Error(err)
		}
	}
}

func (a *Automate) actionAfterDecorationRun(actions []model.ActionInfo, err error) {
	defer a.ErrorRecover()
	for _, fc := range actionAfterDecoration {
		err := fc(actions, err)
		if err != nil {
			logrus.Error(err)
		}
	}
}

func (*Automate) ErrorRecover() func() {
	return func() {
		if r := recover(); r != nil {
			// Get the current call stack
			stack := string(debug.Stack())
			// Log the stack trace information
			logrus.Error("Automation execution exception:\n", r, "\nStack trace:\n", stack)
		}
	}
}

// Execute
// @description Executes telemetry setting reporting automation (reads cache information, queries the database if no information is in the cache, and saves it to the cache)
// @params deviceInfo *model.Device
// @return error
func (a *Automate) Execute(deviceInfo *model.Device, fromExt AutomateFromExt) error {
	defer a.ErrorRecover()
	a.device = deviceInfo
	a.formExt = fromExt
	//

	// Single-category device t
	if deviceInfo.DeviceConfigID != nil {
		deviceConfigId := *deviceInfo.DeviceConfigID
		err := a.telExecute(deviceInfo.ID, deviceConfigId, fromExt)
		if err != nil {
			logrus.Error("Automation execution failed", err)
		}
	}
	return a.telExecute(deviceInfo.ID, "", fromExt)
}

func (a *Automate) telExecute(deviceId, deviceConfigId string, fromExt AutomateFromExt) error {
	info, resultInt, err := initialize.NewAutomateCache().GetCacheByDeviceId(deviceId, deviceConfigId)
	logrus.Debugf("Automation execution started: info:%#v, resultInt:%d", info, resultInt)
	if err != nil {
		return pkgerrors.Wrap(err, "Failed to query cache information")
	}
	// No automation task for the current device
	if resultInt == initialize.AUTOMATE_CACHE_RESULT_NOT_TASK {
		return nil
	}
	// Cache data not found, query the database and set the cache
	if resultInt == initialize.AUTOMATE_CACHE_RESULT_NOT_FOUND {
		info, resultInt, err = a.QueryAutomateInfoAndSetCache(deviceId, deviceConfigId)
		if err != nil {
			return pkgerrors.Wrap(err, "Failed to query and set cache")
		}
		// No automation task for the current device
		if resultInt == initialize.AUTOMATE_CACHE_RESULT_NOT_TASK {
			return nil
		}
	}
	logrus.Debugf("Automation execution started 2: info:%#v, resultInt:%d", info, resultInt)
	// Filter automation trigger conditions
	info = a.AutomateFilter(info, fromExt)
	logrus.Debugf("Automation execution started 3: info:%#v, resultInt:%v", info, fromExt)
	// Execute automation
	return a.ExecuteRun(info)
}

func (a *Automate) AutomateFilter(info initialize.AutomateExecteParams, fromExt AutomateFromExt) initialize.AutomateExecteParams {
	var sceneInfo []initialize.AutomateExecteSceneInfo
	for _, scene := range info.AutomateExecteSceeInfos {
		var isExists bool
		for _, cond := range scene.GroupsCondition {
			if cond.TriggerParamType == nil || cond.TriggerParam == nil {
				continue
			}
			condTriggerParamType := strings.ToUpper(*cond.TriggerParamType)
			switch fromExt.TriggerParamType {
			case model.TRIGGER_PARAM_TYPE_TEL:
				if condTriggerParamType == model.TRIGGER_PARAM_TYPE_TEL || condTriggerParamType == model.TRIGGER_PARAM_TYPE_TELEMETRY {
					if a.containString(fromExt.TriggerParam, *cond.TriggerParam) {
						isExists = true
					}
				}
			case model.TRIGGER_PARAM_TYPE_STATUS:
				if condTriggerParamType == model.TRIGGER_PARAM_TYPE_STATUS {
					isExists = true
				}
			case model.TRIGGER_PARAM_TYPE_EVT:
				if (condTriggerParamType == model.TRIGGER_PARAM_TYPE_EVT || condTriggerParamType == model.TRIGGER_PARAM_TYPE_EVENT) && a.containString(fromExt.TriggerParam, *cond.TriggerParam) {
					isExists = true
				}
			case model.TRIGGER_PARAM_TYPE_ATTR:
				if condTriggerParamType == model.TRIGGER_PARAM_TYPE_ATTR && a.containString(fromExt.TriggerParam, *cond.TriggerParam) {
					isExists = true
				}
			}
		}
		if isExists {
			sceneInfo = append(sceneInfo, scene)
		}
	}
	info.AutomateExecteSceeInfos = sceneInfo
	return info
}

func (*Automate) containString(slice []string, str string) bool {
	for _, v := range slice {
		logrus.Info(v, str)
		if v == str {
			return true
		}
	}
	return false
}

// Rate Limiting Implementation: 1 request per second per scene automation
func (*Automate) LimiterAllow(id string) bool {
	return initialize.NewAutomateLimiter().GetLimiter(fmt.Sprintf("SceneAutomationId:%s", id)).Allow()
}

// ExecuteRun
// @description Executes automation scene linkage actions
// @params info initialize.AutomateExecteParams
// @return error
func (a *Automate) ExecuteRun(info initialize.AutomateExecteParams) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	for _, v := range info.AutomateExecteSceeInfos {
		// Scene frequency limit (based on scene ID)
		if !a.LimiterAllow(v.SceneAutomationId) {
			continue
		}
		logrus.Debugf("Checking if automation is disabled 1: info:%#v,", v.SceneAutomationId)
		// Check if automation is disabled
		if a.CheckSceneAutomationHasClose(v.SceneAutomationId) {
			continue
		}
		logrus.Debugf("Checking if automation is disabled 2: info:%#v,", info)
		// Condition check
		if !a.AutomateConditionCheck(v.GroupsCondition, info.DeviceId) {
			continue
		}
		// Execute scene linkage actions
		err := a.SceneAutomateExecute(v.SceneAutomationId, []string{info.DeviceId}, v.Actions)
		// Post-action decoration
		a.actionAfterDecorationRun(v.Actions, err)
	}

	return nil
}

// CheckSceneAutomationHasClose
// @description Check if the automation has been disabled
func (*Automate) CheckSceneAutomationHasClose(sceneAutomationId string) bool {
	ok := dal.CheckSceneAutomationHasClose(sceneAutomationId)
	// Remove cache
	if ok {
		_ = initialize.NewAutomateCache().DeleteCacheBySceneAutomationId(sceneAutomationId)
	}
	return ok
}

// SceneAutomateExecute
// @description Executes scene linkage actions
// @params info initialize.AutomateExecteParams
// @return error
func (a *Automate) SceneAutomateExecute(sceneAutomationId string, deviceIds []string, actions []model.ActionInfo) error {
	tenantID := dal.GetSceneAutomationTenantID(context.Background(), sceneAutomationId)

	// Execute actions
	details, err := a.AutomateActionExecute(sceneAutomationId, deviceIds, actions, tenantID)

	_ = a.sceneExecuteLogSave(sceneAutomationId, details, err)

	return err
}

// ActiveSceneExecute
// @description Activates the scene
// @params info initialize.AutomateExecteParams
// @return error
func (a *Automate) ActiveSceneExecute(scene_id, tenantID string) error {

	actions, err := dal.GetActionInfoListBySceneId([]string{scene_id})
	if err != nil {
		return nil
	}
	var (
		deviceIds      []string
		deviceConfigId []string
	)
	for _, v := range actions {
		if v.ActionType == model.AUTOMATE_ACTION_TYPE_MULTIPLE && v.ActionTarget != nil {
			deviceConfigId = append(deviceConfigId, *v.ActionTarget)
		}
	}
	if len(deviceConfigId) > 0 {
		deviceIds, err = dal.GetDeviceIdsByDeviceConfigId(deviceConfigId)
		if err != nil {
			return err
		}
	}
	details, err := a.AutomateActionExecute(scene_id, deviceIds, actions, tenantID)
	var exeResult string
	if err == nil {
		exeResult = "S"
	} else {
		exeResult = "F"
	}
	logrus.Debug(details)
	return dal.SceneLogInsert(&model.SceneLog{
		ID:              uuid.New(),
		SceneID:         scene_id,
		ExecutedAt:      time.Now().UTC(),
		Detail:          details,
		ExecutionResult: exeResult,
		TenantID:        tenantID,
	})
}

// @description sceneExecuteLogSave Executes automation scene linkage actions and logs the result
// @params info initialize.AutomateExecteParams
// @return error
func (*Automate) sceneExecuteLogSave(scene_id, details string, err error) error {
	var exeResult string
	if err == nil {
		exeResult = "S"
	} else {
		exeResult = "F"
	}
	logrus.Debug(details)
	return dal.SceneAutomationLogInsert(&model.SceneAutomationLog{
		SceneAutomationID: scene_id,
		ExecutedAt:        time.Now().UTC(),
		Detail:            details,
		ExecutionResult:   exeResult,
		TenantID:          dal.GetSceneAutomationTenantID(context.Background(), scene_id),
	})
}

// AutomateConditionCheck
// @description  Automation condition check. Returns true if any one group of conditions is satisfied.
// @params conditions []initialize.DTConditions
// @return bool - true indicates the actions can be executed
func (a *Automate) AutomateConditionCheck(conditions initialize.DTConditions, deviceId string) bool {
	logrus.Debug("Starting condition check...")
	// Group conditions by group ID
	conditionsByGroupId := make(map[string]initialize.DTConditions)
	for _, v := range conditions {
		conditionsByGroupId[v.GroupID] = append(conditionsByGroupId[v.GroupID], v)
	}
	var result bool
	for _, val := range conditionsByGroupId {
		ok, contents := a.AutomateConditionCheckWithGroup(val, deviceId)
		if ok {
			result = true
		}
		// Post-check group condition hook
		a.conditionAfterDecorationRun(ok, val, deviceId, contents)
	}
	return result
}

// AutomateConditionCheckWithGroup
// @description  Validates a group of conditions; returns false if any single condition in the group fails
// @params conditions initialize.DTConditions
// @return bool
func (a *Automate) AutomateConditionCheckWithGroup(conditions initialize.DTConditions, deviceId string) (bool, []string) {
	var (
		result   []string
		resultOk bool = true
	)
	for _, val := range conditions {
		ok, content := a.AutomateConditionCheckWithGroupOne(val, deviceId)
		result = append(result, content)
		if !ok {
			resultOk = false
			break
		}
	}
	return resultOk, result
}

// AutomateConditionCheckWithGroupOne
// @description  Validates a single condition
// @params cond model.DeviceTriggerCondition
// @return bool
func (a *Automate) AutomateConditionCheckWithGroupOne(cond model.DeviceTriggerCondition, deviceId string) (bool, string) {
	logrus.Debug("Condition type:", cond.TriggerConditionType)
	switch cond.TriggerConditionType {
	case model.DEVICE_TRIGGER_CONDITION_TYPE_TIME:
		return a.automateConditionCheckWithTime(cond), ""
	case model.DEVICE_TRIGGER_CONDITION_TYPE_ONE, model.DEVICE_TRIGGER_CONDITION_TYPE_MULTIPLE:
		return a.automateConditionCheckWithDevice(cond, deviceId)
	default:
		return true, ""
	}
}

// automateConditionCheckWithTime
// @description  Validates time-range conditions
// @params cond model.DeviceTriggerCondition
// @return bool
func (*Automate) automateConditionCheckWithTime(cond model.DeviceTriggerCondition) bool {
	logrus.Debug("Starting time range comparison... Condition:", cond.TriggerValue)
	nowTime := time.Now().UTC()
	if cond.TriggerValue == "" {
		return false
	}
	valParts := strings.Split(cond.TriggerValue, "|")
	if len(valParts) < 3 {
		return false
	}

	// Determine if today is included in the specified weekdays
	isDayMatched := false
	weekDay := common.GetWeekDay(nowTime)
	for _, char := range valParts[0] {
		num, _ := strconv.Atoi(string(char))
		if weekDay == num {
			isDayMatched = true
			break
		}
	}
	if !isDayMatched {
		return false
	}

	// Compare current time with time range
	nowTimeOfDay, _ := time.Parse("15:04:05-07:00", nowTime.Format("15:04:05-07:00"))
	startTime, err := time.Parse("15:04:05-07:00", valParts[1])
	if err != nil {
		logrus.Error("Invalid time format in condition string:", cond.TriggerValue)
		return false
	}
	if startTime.After(nowTimeOfDay) {
		return false
	}

	endTime, err := time.Parse("15:04:05-07:00", valParts[2])
	if err != nil {
		logrus.Error("Invalid time format in condition string:", cond.TriggerValue)
		return false
	}
	if endTime.Before(nowTimeOfDay) {
		return false
	}
	logrus.Debug("Time range comparison passed.")
	return true
}

func (a *Automate) getActualValue(deviceId string, key string, triggerParamType string) (interface{}, error) {
	for k, v := range a.formExt.TriggerValues {
		if key == k {
			return v, nil
		}
	}
	switch triggerParamType {
	case model.TRIGGER_PARAM_TYPE_TEL:
		return dal.GetCurrentTelemetryDataOneKeys(deviceId, key)
	case model.TRIGGER_PARAM_TYPE_ATTR:
		return dal.GetAttributeOneKeys(deviceId, key)
	case model.TRIGGER_PARAM_TYPE_EVT:
		return dal.GetDeviceEventOneKeys(deviceId, key)
	case model.TRIGGER_PARAM_TYPE_STATUS:
		return dal.GetDeviceCurrentStatus(deviceId)
	}

	return nil, nil
}

// automateConditionCheckWithDevice
// @description  Validates a device condition by comparing actual device data (telemetry, attribute, event, or status)
//
//	against the expected trigger value using the specified operator.
//
// @params cond model.DeviceTriggerCondition - The condition to validate
// @params deviceId string - The ID of the device
// @return bool - True if condition is met, false otherwise
// @return string - Description of the evaluation result
func (a *Automate) automateConditionCheckWithDevice(cond model.DeviceTriggerCondition, deviceId string) (bool, string) {
	logrus.Debug("Starting device condition check...")

	// If trigger source is missing, return false
	if cond.TriggerSource == nil {
		return false, ""
	}

	// If the condition type is single-device, use the device ID from the condition
	if cond.TriggerConditionType == model.DEVICE_TRIGGER_CONDITION_TYPE_ONE {
		deviceId = *cond.TriggerSource
	}

	var (
		actualValue     interface{}
		trigger         string
		triggerValue    string
		triggerOperator string
		triggerKey      string
		result          string
		deviceName      string
	)

	if a.device.Name != nil {
		deviceName = *a.device.Name
	}

	if cond.TriggerOperator == nil {
		triggerOperator = "="
	} else {
		triggerOperator = *cond.TriggerOperator
	}

	logrus.Debug("Evaluating device condition type:", strings.ToUpper(*cond.TriggerParamType))

	switch strings.ToUpper(*cond.TriggerParamType) {
	case model.TRIGGER_PARAM_TYPE_TEL, model.TRIGGER_PARAM_TYPE_TELEMETRY: // Telemetry
		trigger = "Telemetry"
		actualValue, _ = a.getActualValue(deviceId, *cond.TriggerParam, model.TRIGGER_PARAM_TYPE_TEL)
		triggerValue = cond.TriggerValue
		triggerKey = *cond.TriggerParam
		logrus.Debugf("Telemetry comparison - operator:%s, param:%s, expected:%v, actual:%v", triggerOperator, triggerKey, triggerValue, actualValue)
		dataValue := a.getTriggerParamsValue(triggerKey, dal.GetIdentifierNameTelemetry())
		result = fmt.Sprintf("Device(%s) %s [%s]: %v %s %v", deviceName, trigger, dataValue, actualValue, triggerOperator, triggerValue)

	case model.TRIGGER_PARAM_TYPE_ATTR: // Attribute
		trigger = "Attribute"
		actualValue, _ = a.getActualValue(deviceId, *cond.TriggerParam, model.TRIGGER_PARAM_TYPE_ATTR)
		triggerValue = cond.TriggerValue
		triggerKey = *cond.TriggerParam
		dataValue := a.getTriggerParamsValue(triggerKey, dal.GetIdentifierNameAttribute())
		result = fmt.Sprintf("Device(%s) %s [%s]: %v %s %v", deviceName, trigger, dataValue, actualValue, triggerOperator, triggerValue)

	case model.TRIGGER_PARAM_TYPE_EVT, model.TRIGGER_PARAM_TYPE_EVENT: // Event
		trigger = "Event"
		actualValue, _ = a.getActualValue(deviceId, *cond.TriggerParam, model.TRIGGER_PARAM_TYPE_EVT)
		triggerValue = cond.TriggerValue
		triggerKey = *cond.TriggerParam
		logrus.Debugf("Event evaluation - actual:%#v, expected:%#v", actualValue, triggerValue)
		dataValue := a.getTriggerParamsValue(triggerKey, dal.GetIdentifierNameEvent())
		result = fmt.Sprintf("Device(%s) %s [%s]: %v %s %v", deviceName, trigger, dataValue, actualValue, triggerOperator, triggerValue)

	case model.TRIGGER_PARAM_TYPE_STATUS: // Online/Offline Status
		trigger = "Offline"
		actualValue, _ = a.getActualValue(deviceId, "login", model.TRIGGER_PARAM_TYPE_STATUS)
		triggerValue = *cond.TriggerParam
		if strings.ToUpper(actualValue.(string)) == "ON-LINE" {
			trigger = "Online"
		}
		result = fmt.Sprintf("Device(%s) is %s", deviceName, trigger)
		triggerOperator = "="
		// If trigger value is "ALL", any status is acceptable
		if strings.ToUpper(triggerValue) == "ALL" {
			return true, result
		}
	}

	logrus.Debug("Checking condition: operator=", triggerOperator, " expected=", triggerValue, " actual=", actualValue)
	ok := a.automateConditionCheckByOperator(triggerOperator, triggerValue, actualValue)
	logrus.Debugf("Comparison result: %t", ok)

	return ok, result
}

type DataIdentifierName func(device_template_id, identifier string) string

func (*Automate) getTriggerParamsValue(triggerKey string, fc DataIdentifierName) string {
	tempId, _ := dal.GetDeviceTemplateIdByDeviceId(triggerKey)
	if tempId == "" {
		return triggerKey
	}

	return fc(tempId, triggerKey)
}

// automateConditionCheckByOperator
// @description  Applies the appropriate operator logic based on the data type of the actual value.
// @params operator string - The comparison operator (e.g., =, >, <, etc.)
// @params condValue string - The expected value from the condition (as string)
// @params actualValue interface{} - The value retrieved from the device
// @return bool - Whether the condition is met
func (a *Automate) automateConditionCheckByOperator(operator string, condValue string, actualValue interface{}) bool {
	switch value := actualValue.(type) {
	case string:
		return a.automateConditionCheckByOperatorWithString(operator, condValue, value)
	case float64:
		return a.automateConditionCheckByOperatorWithFloat(operator, condValue, value)
	case bool:
		// Convert bool to string and reuse string comparison logic
		return a.automateConditionCheckByOperatorWithString(operator, condValue, fmt.Sprintf("%t", value))
	}
	return false
}

func float64Equal(a, b float64) bool {
	const threshold = 1e-9
	return math.Abs(a-b) < threshold
}

// automateConditionCheckByOperatorWithFloat
// @description  Compares float values using a given operator.
// @params operator string - One of EQ, NEQ, GT, LT, GTE, LTE, BETWEEN, IN
// @params condValue string - Expected value or range (e.g., "10", "10-20", "10,15")
// @params actualValue float64 - Actual numeric value from device
// @return bool - Whether the comparison passes
func (*Automate) automateConditionCheckByOperatorWithFloat(operator string, condValue string, actualValue float64) bool {
	switch operator {
	case model.CONDITION_TRIGGER_OPERATOR_EQ:
		condFloat, err := strconv.ParseFloat(condValue, 64)
		if err != nil {
			logrus.Error(err)
			return false
		}
		return float64Equal(condFloat, actualValue)

	case model.CONDITION_TRIGGER_OPERATOR_NEQ:
		condFloat, err := strconv.ParseFloat(condValue, 64)
		if err != nil {
			logrus.Error(err)
			return false
		}
		return !float64Equal(condFloat, actualValue)

	case model.CONDITION_TRIGGER_OPERATOR_GT:
		condFloat, err := strconv.ParseFloat(condValue, 64)
		if err != nil {
			logrus.Error(err)
			return false
		}
		return actualValue > condFloat

	case model.CONDITION_TRIGGER_OPERATOR_LT:
		condFloat, err := strconv.ParseFloat(condValue, 64)
		if err != nil {
			logrus.Error(err)
			return false
		}
		return actualValue < condFloat

	case model.CONDITION_TRIGGER_OPERATOR_GTE:
		condFloat, err := strconv.ParseFloat(condValue, 64)
		if err != nil {
			logrus.Error(err)
			return false
		}
		return actualValue >= condFloat

	case model.CONDITION_TRIGGER_OPERATOR_LTE:
		condFloat, err := strconv.ParseFloat(condValue, 64)
		if err != nil {
			logrus.Error(err)
			return false
		}
		return actualValue <= condFloat

	case model.CONDITION_TRIGGER_OPERATOR_BETWEEN:
		// Format: "min-max"
		valParts := strings.Split(condValue, "-")
		if len(valParts) != 2 {
			return false
		}
		minVal, err1 := strconv.ParseFloat(valParts[0], 64)
		maxVal, err2 := strconv.ParseFloat(valParts[1], 64)
		if err1 != nil || err2 != nil {
			return false
		}
		return actualValue >= minVal && actualValue <= maxVal

	case model.CONDITION_TRIGGER_OPERATOR_IN:
		// Format: "v1,v2,v3"
		valParts := strings.Split(condValue, ",")
		for _, part := range valParts {
			val, err := strconv.ParseFloat(part, 64)
			if err != nil {
				continue
			}
			if float64Equal(val, actualValue) {
				return true
			}
		}
	}
	return false
}

// automateConditionCheckByOperatorWithString
// @description  Evaluates a condition using a string comparison based on the specified operator.
// @params operator string - The operator (e.g., EQ, GT, IN)
// @params condValue string - The expected value or range (as a string)
// @params actualValue string - The actual value to compare
// @return bool - Whether the comparison result is true
func (*Automate) automateConditionCheckByOperatorWithString(operator string, condValue string, actualValue string) bool {
	logrus.Warningf("Compare: operator: %s, condValue: %s, actualValue: %s, result: %d",
		operator, condValue, actualValue, strings.Compare(actualValue, condValue))

	switch operator {
	case model.CONDITION_TRIGGER_OPERATOR_EQ:
		// Case-insensitive equality check
		return strings.EqualFold(strings.ToUpper(actualValue), strings.ToUpper(condValue))

	case model.CONDITION_TRIGGER_OPERATOR_NEQ:
		return strings.Compare(actualValue, condValue) != 0

	case model.CONDITION_TRIGGER_OPERATOR_GT:
		// Try numeric comparison, fallback to string comparison
		actualFloat, err1 := strconv.ParseFloat(actualValue, 64)
		condFloat, err2 := strconv.ParseFloat(condValue, 64)
		if err1 == nil && err2 == nil {
			return actualFloat > condFloat
		}
		return strings.Compare(actualValue, condValue) > 0

	case model.CONDITION_TRIGGER_OPERATOR_LT:
		actualFloat, err1 := strconv.ParseFloat(actualValue, 64)
		condFloat, err2 := strconv.ParseFloat(condValue, 64)
		if err1 == nil && err2 == nil {
			return actualFloat < condFloat
		}
		return strings.Compare(actualValue, condValue) < 0

	case model.CONDITION_TRIGGER_OPERATOR_GTE:
		actualFloat, err1 := strconv.ParseFloat(actualValue, 64)
		condFloat, err2 := strconv.ParseFloat(condValue, 64)
		if err1 == nil && err2 == nil {
			return actualFloat >= condFloat
		}
		return strings.Compare(actualValue, condValue) >= 0

	case model.CONDITION_TRIGGER_OPERATOR_LTE:
		actualFloat, err1 := strconv.ParseFloat(actualValue, 64)
		condFloat, err2 := strconv.ParseFloat(condValue, 64)
		if err1 == nil && err2 == nil {
			return actualFloat <= condFloat
		}
		return strings.Compare(actualValue, condValue) <= 0

	case model.CONDITION_TRIGGER_OPERATOR_BETWEEN:
		valParts := strings.Split(condValue, "-")
		if len(valParts) != 2 {
			return false
		}
		actualFloat, err := strconv.ParseFloat(actualValue, 64)
		val1Float, err1 := strconv.ParseFloat(valParts[0], 64)
		val2Float, err2 := strconv.ParseFloat(valParts[1], 64)

		if err == nil && err1 == nil && err2 == nil {
			return actualFloat >= val1Float && actualFloat <= val2Float
		}
		// Fallback: use string range check
		return actualValue >= valParts[0] && actualValue <= valParts[1]

	case model.CONDITION_TRIGGER_OPERATOR_IN:
		valParts := strings.Split(condValue, ",")
		for _, v := range valParts {
			if v == actualValue {
				return true
			}
		}
	}

	return false
}

// AutomateActionExecute
// @description  Executes automation actions for a list of devices.
// @params _ string - (Unused parameter)
// @params deviceIds []string - List of target device IDs
// @params actions []model.ActionInfo - List of automation actions to execute
// @params tenantID string - Tenant ID for context
// @return string, error - Summary result and first encountered error (if any)
func (*Automate) AutomateActionExecute(_ string, deviceIds []string, actions []model.ActionInfo, tenantID string) (string, error) {
	logrus.Debug("Starting automation actions execution:")
	var (
		result    string
		resultErr error
	)

	if len(actions) == 0 {
		return "No actions found to execute", errors.New("no actions provided")
	}

	for _, action := range actions {
		var actionService AutomateTelemetryAction
		logrus.Debug("ActionType:", action.ActionType)

		switch action.ActionType {
		case model.AUTOMATE_ACTION_TYPE_ONE:
			actionService = &AutomateTelemetryActionOne{TenantID: tenantID}
		case model.AUTOMATE_ACTION_TYPE_ALARM:
			actionService = &AutomateTelemetryActionAlarm{}
		case model.AUTOMATE_ACTION_TYPE_MULTIPLE:
			actionService = &AutomateTelemetryActionMultiple{DeviceIds: deviceIds, TenantID: tenantID}
		case model.AUTOMATE_ACTION_TYPE_SCENE:
			actionService = &AutomateTelemetryActionScene{TenantID: tenantID}
		case model.AUTOMATE_ACTION_TYPE_SERVICE:
			actionService = &AutomateTelemetryActionService{}
		}

		if actionService == nil {
			logrus.Error("Unsupported action type")
			return "Unsupported action type", errors.New("unsupported action type")
		}

		actionMessage, err := actionService.AutomateActionRun(action)
		if err != nil && resultErr == nil {
			resultErr = err
		}
		if err != nil {
			result += fmt.Sprintf("%s failed;", actionMessage)
		} else {
			result += fmt.Sprintf("%s succeeded;", actionMessage)
		}
	}

	logrus.Debug("Execution result:", result)
	return result, resultErr
}

// QueryAutomateInfoAndSetCache
// @description  Queries automation configuration for a device and sets the result to cache.
// @params deviceId string - Device ID
// @params deviceConfigId string - Device configuration ID (for type: MULTIPLE)
// @return initialize.AutomateExecteParams, int, error - Automation parameters, cache result code, and error if any
func (*Automate) QueryAutomateInfoAndSetCache(deviceId, deviceConfigId string) (initialize.AutomateExecteParams, int, error) {
	automateExecuteParams := initialize.AutomateExecteParams{
		DeviceId:       deviceId,
		DeviceConfigId: deviceConfigId,
	}

	var (
		groups []model.DeviceTriggerCondition
		err    error
	)

	// If deviceConfigId is provided, use MULTIPLE-type trigger conditions
	if deviceConfigId != "" {
		groups, err = dal.GetDeviceTriggerConditionByDeviceId(deviceConfigId, model.DEVICE_TRIGGER_CONDITION_TYPE_MULTIPLE)
	} else {
		groups, err = dal.GetDeviceTriggerConditionByDeviceId(deviceId, model.DEVICE_TRIGGER_CONDITION_TYPE_ONE)
	}

	logrus.Debugf("DeviceConfigId: %s, Conditions: %v", deviceConfigId, groups)

	if err != nil {
		return automateExecuteParams, 0, pkgerrors.Wrap(err, "Failed to query automation conditions by device ID")
	}

	if len(groups) == 0 {
		err := initialize.NewAutomateCache().SetCacheByDeviceIdWithNoTask(deviceId, deviceConfigId)
		if err != nil {
			return automateExecuteParams, 0, pkgerrors.Wrap(err, "Failed to cache empty automation for device")
		}
		return automateExecuteParams, initialize.AUTOMATE_CACHE_RESULT_NOT_TASK, nil
	}

	sceneAutomateGroups := make(map[string]bool)
	var (
		sceneAutomateIds []string
		groupIds         []string
	)

	for _, groupInfo := range groups {
		if _, ok := sceneAutomateGroups[groupInfo.SceneAutomationID]; !ok {
			sceneAutomateIds = append(sceneAutomateIds, groupInfo.SceneAutomationID)
			sceneAutomateGroups[groupInfo.SceneAutomationID] = true
		}
		groupIds = append(groupIds, groupInfo.GroupID)
	}

	// Fetch all conditions by group ID
	groups, err = dal.GetDeviceTriggerConditionByGroupIds(groupIds)
	if err != nil {
		return automateExecuteParams, 0, pkgerrors.Wrap(err, "Failed to fetch conditions by group ID")
	}

	// Fetch automation actions by scene IDs
	actionInfos, err := dal.GetActionInfoListBySceneAutomationId(sceneAutomateIds)
	if err != nil {
		return automateExecuteParams, 0, pkgerrors.Wrap(err, "Failed to fetch automation actions")
	}

	logrus.Debugf("DeviceConfigId: %s, Conditions: %v, Actions: %v", deviceConfigId, groups, actionInfos)

	// Set cache
	err = initialize.NewAutomateCache().SetCacheByDeviceId(deviceId, deviceConfigId, groups, actionInfos)
	if err != nil {
		return automateExecuteParams, 0, pkgerrors.Wrap(err, "Failed to set automation cache")
	}

	return initialize.NewAutomateCache().GetCacheByDeviceId(deviceId, deviceConfigId)
}
