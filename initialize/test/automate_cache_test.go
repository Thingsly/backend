package test

import (
	"fmt"
	"testing"

	"github.com/HustIoTPlatform/backend/initialize"
	"github.com/HustIoTPlatform/backend/internal/model"
)

func init() {
	initialize.ViperInit("../../configs/conf-localdev.yml")
	initialize.RedisInit()
	cache = initialize.NewAutomateCache()
}

var (
	sceneAutomateId = "sceneAutomateId_test"
	cache           *initialize.AutomateCache
)

func StringPoints(s string) *string {
	return &s
}

func getConditions(sceneAutomateId string) []model.DeviceTriggerCondition {
	var conditions []model.DeviceTriggerCondition

	condition := model.DeviceTriggerCondition{
		SceneAutomationID:    sceneAutomateId,
		GroupID:              "groupId",
		TriggerConditionType: "10",
		TriggerValue:         "30",
		TriggerSource:        StringPoints("condition_deviceIds01"),
		TriggerParamType:     StringPoints("TEL"),
		TriggerParam:         StringPoints("temperature"),
		TriggerOperator:      StringPoints(">"),
	}

	conditions = append(conditions, condition)

	condition1 := model.DeviceTriggerCondition{
		SceneAutomationID:    sceneAutomateId,
		GroupID:              "groupId",
		TriggerConditionType: "10",
		TriggerValue:         "30",
		TriggerSource:        StringPoints("condition_deviceIds02"),
		TriggerParamType:     StringPoints("TEL"),
		TriggerParam:         StringPoints("temperature"),
		TriggerOperator:      StringPoints(">"),
	}
	conditions = append(conditions, condition1)
	return conditions
}

func getActions(sceneAutomateId string) []model.ActionInfo {
	var actions []model.ActionInfo
	action1 := model.ActionInfo{
		SceneAutomationID: sceneAutomateId,
		ActionType:        "10",
		ActionTarget:      StringPoints("action_deviceIds01"),
		ActionParamType:   StringPoints("CMD"),
		ActionParam:       StringPoints("test_cmd"),
		ActionValue:       StringPoints("test_val"),
	}
	actions = append(actions, action1)
	action2 := model.ActionInfo{
		SceneAutomationID: sceneAutomateId,
		ActionType:        "11",
		ActionTarget:      StringPoints("action_deviceIds02"),
		ActionParamType:   StringPoints("CMD"),
		ActionParam:       StringPoints("test_cmd"),
		ActionValue:       StringPoints("test_val"),
	}
	actions = append(actions, action2)
	return actions
}

// conditions []model.DeviceTriggerCondition, actions []model.ActionInfo
func TestSetCacheBySceneAutomationId(t *testing.T) {
	fmt.Println("Testing cache creation...")
	// cache := initialize.NewAutomateCache()
	conditions := getConditions(sceneAutomateId)
	err := cache.SetCacheBySceneAutomationId(sceneAutomateId, conditions, getActions(sceneAutomateId))
	if err != nil {
		t.Error("Failed to save automation cache", err)
	}
}

func TestGetCacheByDeviceId(t *testing.T) {
	fmt.Println("Testing cache retrieval when data exists...")
	// cache := initialize.NewAutomateCache()
	res, resultCode, err := cache.GetCacheByDeviceId("condition_deviceIds01", "")
	if err != nil {
		t.Error("Failed to retrieve automation cache by device ID", err)
	}
	if resultCode != initialize.AUTOMATE_CACHE_RESULT_OK {
		t.Errorf("Unexpected query result. Status: %d, Result: %#v", resultCode, res)
	}
	fmt.Printf("Result: %#v", res)
}

func TestDeleteCacheBySceneAutomationId(t *testing.T) {
	fmt.Println("Testing deletion of scene automation cache...")
	// cache := initialize.NewAutomateCache()
	err := cache.DeleteCacheBySceneAutomationId(sceneAutomateId)
	if err != nil {
		t.Error("Failed to delete cache", err)
	}
}

func TestGetCacheByDeviceIdNotExists(t *testing.T) {
	fmt.Println("Testing cache retrieval for device with no data...")
	// cache := initialize.NewAutomateCache()
	res, resultCode, err := cache.GetCacheByDeviceId("condition_deviceIds0004", "")
	fmt.Printf("Query status: %d", resultCode)
	if err != nil {
		t.Error("Failed to retrieve automation cache by device ID", err)
	}
	if resultCode != initialize.AUTOMATE_CACHE_RESULT_NOT_FOUND {
		t.Errorf("Unexpected query result. Status: %d, Result: %#v", resultCode, res)
	}
	fmt.Printf("Result: %#v", res)
}

func TestSetCacheByDeviceIdWithNoTask(t *testing.T) {
	fmt.Println("Testing cache entry for device with no associated task...")
	err := cache.SetCacheByDeviceIdWithNoTask("condition_deviceIds00005", "")
	if err != nil {
		t.Error("Failed to cache device with no associated task", err)
	}
	_, resultCode, err := cache.GetCacheByDeviceId("condition_deviceIds00005", "")
	fmt.Printf("Query result: resultCode: %d", resultCode)
	if err != nil {
		t.Error("Failed to retrieve automation cache by device ID", err)
	}
	if resultCode != initialize.AUTOMATE_CACHE_RESULT_NOT_TASK {
		t.Errorf("Unexpected result code, expected NOT_TASK. Got: %d", resultCode)
	}
}

func TestSetCacheByDeviceId(t *testing.T) {
	fmt.Println("Testing caching of automation information by device ID...")
	deviceId := "condition_deviceIds_with_device"

	conditions := []model.DeviceTriggerCondition{
		{
			SceneAutomationID:    "sceneAutomateId_with_device_01",
			GroupID:              "groupId",
			TriggerConditionType: "10",
			TriggerValue:         "20-20",
			TriggerSource:        StringPoints(deviceId),
			TriggerParamType:     StringPoints("TEL"),
			TriggerParam:         StringPoints("temperature"),
			TriggerOperator:      StringPoints(">"),
		},
		{
			SceneAutomationID:    "sceneAutomateId_with_device_02",
			GroupID:              "groupId",
			TriggerConditionType: "10",
			TriggerValue:         "20-20",
			TriggerSource:        StringPoints(deviceId),
			TriggerParamType:     StringPoints("TEL"),
			TriggerParam:         StringPoints("temperature"),
			TriggerOperator:      StringPoints(">"),
		},
	}

	actions := []model.ActionInfo{
		{
			SceneAutomationID: "sceneAutomateId_with_device_01",
			ActionType:        "10",
			ActionTarget:      StringPoints("action_deviceIds01"),
			ActionParamType:   StringPoints("CMD"),
			ActionParam:       StringPoints("test_cmd"),
			ActionValue:       StringPoints("test_val"),
		},
		{
			SceneAutomationID: "sceneAutomateId_with_device_02",
			ActionType:        "10",
			ActionTarget:      StringPoints("action_deviceIds01"),
			ActionParamType:   StringPoints("CMD"),
			ActionParam:       StringPoints("test_cmd"),
			ActionValue:       StringPoints("test_val"),
		},
	}

	err := cache.SetCacheByDeviceId(deviceId, "", conditions, actions)
	if err != nil {
		t.Error("Failed to cache automation data by device ID", err)
	}
	res, resultCode, err := cache.GetCacheByDeviceId(deviceId, "")
	fmt.Printf("Query result: resultCode: %d; Cached data: %#v", resultCode, res)
	if err != nil {
		t.Error("Failed to retrieve automation cache by device ID", err)
	}
	if resultCode != initialize.AUTOMATE_CACHE_RESULT_OK {
		t.Errorf("Unexpected result code, expected OK. Got: %d", resultCode)
	}
}

func TestSetCacheByDeviceConfidId(t *testing.T) {
	fmt.Println("Testing automation cache creation using device ID and configuration ID...")

	deviceId := "condition_deviceIds_with_device"
	deviceConfigId := "condition_deviceIds_with_device_config_id"

	conditions := []model.DeviceTriggerCondition{
		{
			SceneAutomationID:    "sceneAutomateId_with_device_01",
			GroupID:              "groupId",
			TriggerConditionType: "11",
			TriggerValue:         "21",
			TriggerSource:        StringPoints(deviceConfigId),
			TriggerParamType:     StringPoints("TEL"),
			TriggerParam:         StringPoints("temperature"),
			TriggerOperator:      StringPoints(">"),
		},
		{
			SceneAutomationID:    "sceneAutomateId_with_device_01",
			GroupID:              "groupId",
			TriggerConditionType: "22",
			TriggerValue:         "137|06:30:00+00:00|16:30:00+00:00",
		},
		{
			SceneAutomationID:    "sceneAutomateId_with_device_02",
			GroupID:              "groupId02",
			TriggerConditionType: "11",
			TriggerValue:         "21",
			TriggerSource:        StringPoints(deviceConfigId),
			TriggerParamType:     StringPoints("TEL"),
			TriggerParam:         StringPoints("temperature"),
			TriggerOperator:      StringPoints(">"),
		},
	}

	actions := []model.ActionInfo{
		{
			SceneAutomationID: "sceneAutomateId_with_device_01",
			ActionType:        "10",
			ActionTarget:      StringPoints("action_deviceIds01"),
			ActionParamType:   StringPoints("CMD"),
			ActionParam:       StringPoints("test_cmd"),
			ActionValue:       StringPoints("test_val"),
		},
		{
			SceneAutomationID: "sceneAutomateId_with_device_02",
			ActionType:        "10",
			ActionTarget:      StringPoints("action_deviceIds01"),
			ActionParamType:   StringPoints("CMD"),
			ActionParam:       StringPoints("test_cmd"),
			ActionValue:       StringPoints("test_val"),
		},
	}

	err := cache.SetCacheByDeviceId(deviceId, deviceConfigId, conditions, actions)
	if err != nil {
		t.Error("Failed to save automation cache with device and config ID", err)
	}

	res, resultCode, err := cache.GetCacheByDeviceId(deviceId, deviceConfigId)
	fmt.Printf("Query result: resultCode: %d; Cached data: %#v\n", resultCode, res)

	if err != nil {
		t.Error("Failed to retrieve automation cache by device and config ID", err)
	}

	if resultCode != initialize.AUTOMATE_CACHE_RESULT_OK {
		t.Errorf("Unexpected result code, expected OK. Got: %d", resultCode)
	}
}
