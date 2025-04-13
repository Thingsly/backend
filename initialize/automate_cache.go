package initialize

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/HustIoTPlatform/backend/initialize/automatecache"
	"github.com/HustIoTPlatform/backend/internal/model"
	global "github.com/HustIoTPlatform/backend/pkg/global"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

const (
	AUTOMATE_CACHE_RESULT_NOT_FOUND = iota + 1
	AUTOMATE_CACHE_RESULT_NOT_TASK
	AUTOMATE_CACHE_RESULT_OK
)

const AUTOMATE_CACHE_CONTENT_NOT_TASK = "NOTASK"

type AutomateCache struct {
	client    *redis.Client
	expiredIn time.Duration
	device    AutimateCacheKeyDevice
}

var (
	instance *AutomateCache
	mu       sync.Mutex
)

func NewAutomateCache() *AutomateCache {
	mu.Lock()
	defer mu.Unlock()
	if instance == nil {
		instance = &AutomateCache{
			client:    global.REDIS,
			expiredIn: time.Minute * 5,
			device:    automatecache.NewOneDeviceCache(),
		}
	}
	return instance
}

type AutimateCacheKeyDevice interface {
	GetAutomateCacheKeyPrefix() string
	GetDeviceTriggerConditionType() string
}

// @description getAutomateCacheKey
// @params cType string
// @params id sting
// @return string
func (c *AutomateCache) getAutomateCacheKey(cType string, id string) string {
	return fmt.Sprintf("automate:v3:%s:%s:%s", c.device.GetAutomateCacheKeyPrefix(), cType, id)
}

func (c *AutomateCache) getAutomateCacheKeyBase(deviceId string) string {
	return c.getAutomateCacheKey("_", deviceId)
}

func (c *AutomateCache) getAutomateCacheKeyGroup(groupId string) string {
	return c.getAutomateCacheKey("_group_", groupId)
}

func (c *AutomateCache) getAutomateCacheKeyAction(sceneAutomationId string) string {
	return c.getAutomateCacheKey("_action_", sceneAutomationId)
}

func (c *AutomateCache) GetDeviceType() (AutimateCacheKeyDevice, error) {
	if c.device == nil {
		return nil, errors.New("未设置设备类型 单一设备或单类设备")
	}
	return c.device, nil
}

func (c *AutomateCache) SetDeviceType(deviceType AutimateCacheKeyDevice) {
	c.device = deviceType
}

func (c *AutomateCache) set(key string, value interface{}, expiration time.Duration) error {
	var valueStr string
	if val, ok := value.(string); ok {
		valueStr = val
	} else {
		valBytes, err := json.Marshal(value)
		if err != nil {
			return nil
		}
		valueStr = string(valBytes)
	}
	return c.client.Set(context.Background(), key, valueStr, expiration).Err()
}

func (*AutomateCache) scan(stringCmd *redis.StringCmd, val interface{}) (int, error) {
	str, err := stringCmd.Result()
	if err == redis.Nil {
		return AUTOMATE_CACHE_RESULT_NOT_FOUND, nil
	} else if err != nil {
		return 0, err
	}
	if str == AUTOMATE_CACHE_CONTENT_NOT_TASK {
		return AUTOMATE_CACHE_RESULT_NOT_TASK, nil
	}
	return AUTOMATE_CACHE_RESULT_OK, stringCmd.Scan(val)
}

// DeleteCacheBySceneAutomationId
// @description
// @params sceneAutomationId string
// @return error
func (c *AutomateCache) DeleteCacheBySceneAutomationId(sceneAutomationId string) error {

	c.device = automatecache.NewMultipleDeviceCache()
	err := c.deleteCacheBySceneAutId(sceneAutomationId)
	if err != nil {
		return err
	}

	c.device = automatecache.NewOneDeviceCache()
	err = c.deleteCacheBySceneAutId(sceneAutomationId)
	if err != nil {
		return err
	}
	return nil
}

// DeleteCacheBySceneAutomationId
// @description
// @params sceneAutomationId string
// @return error
func (c *AutomateCache) deleteCacheBySceneAutId(sceneAutomationId string) error {

	var (
		action         AutomateActionInfo
		deleteCacheKes []string
		deviceIds      []string
	)
	actionCacheKey := c.getAutomateCacheKeyAction(sceneAutomationId)
	resultInt, err := c.scan(c.client.Get(context.Background(), actionCacheKey), &action)
	if err != nil {
		return err
	}

	//if resultInt == AUTOMATE_CACHE_RESULT_NOT_TASK || resultInt == AUTOMATE_CACHE_RESULT_NOT_FOUND {
	//	return nil
	//}
	if resultInt == AUTOMATE_CACHE_RESULT_NOT_FOUND {
		return nil
	}
	deleteCacheKes = append(deleteCacheKes, actionCacheKey)

	for _, groupId := range action.GroupIds {
		var (
			groupCacheKey = c.getAutomateCacheKeyGroup(groupId)
			groupInfos    DTConditions
		)
		deleteCacheKes = append(deleteCacheKes, groupCacheKey)
		err := c.client.Get(context.Background(), groupCacheKey).Scan(&groupInfos)
		if err != nil {
			continue
		}
		for _, v := range groupInfos {
			if v.TriggerConditionType == c.device.GetDeviceTriggerConditionType() {
				deviceIds = append(deviceIds, *v.TriggerSource)
			}
		}
	}

	err = c.automateDeviceCacheDeleteHandel(sceneAutomationId, deviceIds, &deleteCacheKes)
	if err != nil {
		return err
	}

	return c.client.Del(context.Background(), deleteCacheKes...).Err()
}

// automateDeviceCacheDeleteHandel
// @description
// @params sceneAutomationId string
// @params ids []string
// @params deleteCacheKes *[]string
// @return error
func (c *AutomateCache) automateDeviceCacheDeleteHandel(sceneAutomationId string, ids []string, deleteCacheKes *[]string) error {
	for _, deviceId := range ids {
		var (
			baseCacheKey = c.getAutomateCacheKeyBase(deviceId)
			//baseCacheKey        = getCachekey(deviceId)
			automateDeviceInfos AutomateDeviceInfos
		)
		err := c.client.Get(context.Background(), baseCacheKey).Scan(&automateDeviceInfos)
		if err != nil {
			continue
		}
		for index, val := range automateDeviceInfos {
			if val.SceneAutomationId == sceneAutomationId {
				copy(automateDeviceInfos[index:], automateDeviceInfos[index+1:])
				automateDeviceInfos = automateDeviceInfos[:len(automateDeviceInfos)-1]
				break
			}
		}
		if len(automateDeviceInfos) > 0 {
			*deleteCacheKes = append(*deleteCacheKes, baseCacheKey)
		} else {
			err := c.set(baseCacheKey, automateDeviceInfos, c.expiredIn)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// SetCacheBySceneAutomationId
// @description
// @params sceneAutomationId string
// @params contions []model.DeviceTriggerCondition
// @params actions []model.ActionInfo
// @return error
func (c *AutomateCache) SetCacheBySceneAutomationId(sceneAutomationId string, conditions []model.DeviceTriggerCondition, actions []model.ActionInfo) error {

	c.device = automatecache.NewMultipleDeviceCache()
	err := c.setCacheBySceneAutId(sceneAutomationId, conditions, actions)
	if err != nil {
		return err
	}

	c.device = automatecache.NewOneDeviceCache()
	err = c.setCacheBySceneAutId(sceneAutomationId, conditions, actions)
	if err != nil {
		return err
	}

	return nil
}

// @description SetCacheBySceneAutomationId
// @params sceneAutomationId string
// @params contions []model.DeviceTriggerCondition
// @params actions []model.ActionInfo
// @return error
func (c *AutomateCache) setCacheBySceneAutId(sceneAutomationId string, conditions []model.DeviceTriggerCondition, actions []model.ActionInfo) error {
	automateDeviceInfo := AutomateDeviceInfo{
		SceneAutomationId: sceneAutomationId,
	}
	actionInfos := AutomateActionInfo{
		Actions: actions,
	}
	var (
		groupInfosMap = make(map[string][]model.DeviceTriggerCondition)
		deviceIdsMap  = make(map[string]bool)
	)

	for _, v := range conditions {
		groupInfosMap[v.GroupID] = append(groupInfosMap[v.GroupID], v)
		if v.TriggerConditionType == c.device.GetDeviceTriggerConditionType() {
			if _, ok := deviceIdsMap[*v.TriggerSource]; !ok {
				deviceIdsMap[*v.TriggerSource] = true
			}
		}
	}

	if len(deviceIdsMap) == 0 {
		return nil
	}

	for groupId, val := range groupInfosMap {

		var ok bool
		for _, v := range val {
			if v.TriggerConditionType == c.device.GetDeviceTriggerConditionType() {
				ok = true
				break
			}
		}
		if !ok {
			continue
		}
		automateDeviceInfo.GroupIds = append(automateDeviceInfo.GroupIds, groupId)
		actionInfos.GroupIds = append(actionInfos.GroupIds, groupId)
		var groupCacheKey = c.getAutomateCacheKeyGroup(groupId)
		err := c.set(groupCacheKey, val, c.expiredIn)
		if err != nil {
			return err
		}
	}

	err := c.set(c.getAutomateCacheKeyAction(sceneAutomationId), actionInfos, c.expiredIn)
	if err != nil {
		return err
	}

	for deviceId := range deviceIdsMap {
		var automateDeviceInfos AutomateDeviceInfos
		var deviceCacheKey = c.getAutomateCacheKeyBase(deviceId)
		_, err = c.scan(c.client.Get(context.Background(), deviceCacheKey), &automateDeviceInfos)
		if err != nil {
			continue
		}
		automateDeviceInfos = append(automateDeviceInfos, automateDeviceInfo)
		err = c.set(deviceCacheKey, automateDeviceInfos, c.expiredIn)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetCacheByDeviceId
// @description
// @params deviceId string 1 - No data available 2 - No automation task 3 - Automation task found
// @return AutomateExecteParams error
func (c *AutomateCache) GetCacheByDeviceId(deviceId, deviceConfigId string) (AutomateExecteParams, int, error) {
	var deviceCacheKey string
	if deviceConfigId == "" {
		c.SetDeviceType(automatecache.NewOneDeviceCache())
		deviceCacheKey = c.getAutomateCacheKeyBase(deviceId)
	} else {
		c.SetDeviceType(automatecache.NewMultipleDeviceCache())
		deviceCacheKey = c.getAutomateCacheKeyBase(deviceConfigId)
	}
	return c.getCacheByDId(deviceId, deviceConfigId, deviceCacheKey)
}

// getCacheByDeviceId
// @description
// @params deviceId string 1 - No data available 2 - No automation task 3 - Automation task found
// @return AutomateExecteParams error
func (c *AutomateCache) getCacheByDId(deviceId, deviceConfigId, deviceCacheKey string) (AutomateExecteParams, int, error) {
	var (
		automateDeviceInfos   = make(AutomateDeviceInfos, 0)
		automateExecuteParams = AutomateExecteParams{
			DeviceId:       deviceId,
			DeviceConfigId: deviceConfigId,
		}
		resultInt int
	)
	stringCmd := c.client.Get(context.Background(), deviceCacheKey)
	resultInt, err := c.scan(stringCmd, &automateDeviceInfos)
	if err != nil || resultInt != AUTOMATE_CACHE_RESULT_OK {
		return automateExecuteParams, resultInt, err
	}

	for _, info := range automateDeviceInfos {
		automateExecuteSceneInfo := AutomateExecteSceneInfo{
			SceneAutomationId: info.SceneAutomationId,
		}
		for _, groupId := range info.GroupIds {
			var (
				groupCacheKey = c.getAutomateCacheKeyGroup(groupId)
				condition     DTConditions
			)
			err := c.client.Get(context.Background(), groupCacheKey).Scan(&condition)
			if err != nil {
				logrus.Warning("Redis did not find data (code: 1)", err)
				continue
			}
			automateExecuteSceneInfo.GroupsCondition = append(automateExecuteSceneInfo.GroupsCondition, condition...)
		}
		var (
			actionInfo     AutomateActionInfo
			actionCacheKey = c.getAutomateCacheKeyAction(info.SceneAutomationId)
		)
		err := c.client.Get(context.Background(), actionCacheKey).Scan(&actionInfo)
		if err != nil {
			logrus.Warning("Redis did not find data", err)
			continue
		}
		automateExecuteSceneInfo.Actions = actionInfo.Actions
		automateExecuteParams.AutomateExecteSceeInfos = append(automateExecuteParams.AutomateExecteSceeInfos, automateExecuteSceneInfo)
	}
	return automateExecuteParams, resultInt, nil
}

// SetCacheByDeviceId
// @description
// @params deviceId string
// @params deviceConfigId string
// @params conditions []model.DeviceTriggerCondition
// @params actions []model.ActionInfo
// @return AutomateExecteParams error
func (c *AutomateCache) SetCacheByDeviceId(deviceId, deviceConfigId string, conditions []model.DeviceTriggerCondition, actions []model.ActionInfo) error {
	if deviceConfigId == "" {
		c.device = automatecache.NewOneDeviceCache()
	} else {
		c.device = automatecache.NewMultipleDeviceCache()
	}

	return c.setCache(deviceId, deviceConfigId, conditions, actions)
}

// @description setCacheByDeviceId
// @params deviceId string
// @params deviceConfigId string
// @params conditions []model.DeviceTriggerCondition
// @params actions []model.ActionInfo
// @return AutomateExecteParams error
func (c *AutomateCache) setCache(deviceId, deviceConfigId string, conditions []model.DeviceTriggerCondition, actions []model.ActionInfo) error {
	var (
		groupInfosMap  = make(map[string][]model.DeviceTriggerCondition)
		deviceInfosMap = make(map[string]map[string]bool)
	)
	logrus.Debug("deviceConfigID:", deviceConfigId)
	for _, v := range conditions {
		groupInfosMap[v.GroupID] = append(groupInfosMap[v.GroupID], v)
		if deviceInfosMap[v.SceneAutomationID] == nil {
			deviceInfosMap[v.SceneAutomationID] = make(map[string]bool)
		}
		deviceInfosMap[v.SceneAutomationID][v.GroupID] = true
	}

	for groupId, val := range groupInfosMap {
		var groupCacheKey = c.getAutomateCacheKeyGroup(groupId)
		logrus.Info("groupCacheKey:", groupCacheKey)
		err := c.set(groupCacheKey, val, c.expiredIn)
		if err != nil {
			return err
		}
	}

	var (
		actionsMap = make(map[string][]model.ActionInfo)
	)
	for _, val := range actions {
		actionsMap[val.SceneAutomationID] = append(actionsMap[val.SceneAutomationID], val)
	}

	var automateDeviceInfos []AutomateDeviceInfo
	for sceneAutomationID, actions := range actionsMap {
		var groupIds []string
		if groupsMap, ok := deviceInfosMap[sceneAutomationID]; ok {
			for groupId := range groupsMap {
				groupIds = append(groupIds, groupId)
			}
			actionInfos := AutomateActionInfo{
				Actions:  actions,
				GroupIds: groupIds,
			}
			err := c.set(c.getAutomateCacheKeyAction(sceneAutomationID), actionInfos, c.expiredIn)
			if err != nil {
				return err
			}
			automateDeviceInfos = append(automateDeviceInfos, AutomateDeviceInfo{
				SceneAutomationId: sceneAutomationID,
				GroupIds:          groupIds,
			})
		}
	}
	var cacheKey string
	if deviceConfigId != "" {
		cacheKey = c.getAutomateCacheKeyBase(deviceConfigId)
	} else {
		cacheKey = c.getAutomateCacheKeyBase(deviceId)
	}

	logrus.Info("cacheKey:", cacheKey)

	return c.set(cacheKey, automateDeviceInfos, c.expiredIn)
}

// @description SetCacheByDeviceIdWithNoTask
// @params deviceId string
// @return error
func (c *AutomateCache) SetCacheByDeviceIdWithNoTask(deviceId, deviceConfigId string) error {
	if deviceConfigId == "" {
		c.device = automatecache.NewOneDeviceCache()
	} else {
		c.device = automatecache.NewMultipleDeviceCache()
	}
	cacheKey := c.getAutomateCacheKeyBase(deviceId)
	return c.set(cacheKey, AUTOMATE_CACHE_CONTENT_NOT_TASK, c.expiredIn)
}

type AutomateDeviceInfo struct {
	SceneAutomationId string   `json:"scene_automation_id"`
	GroupIds          []string `json:"group_id"`
}

type AutomateActionInfo struct {
	GroupIds []string           `json:"group_id"`
	Actions  []model.ActionInfo `json:"actions"`
}

type AutomateExecteSceneInfo struct {
	SceneAutomationId string             `json:"scene_automation_id"`
	GroupsCondition   DTConditions       `json:"groups_condition"`
	Actions           []model.ActionInfo `json:"actions"`
}

type AutomateExecteParams struct {
	DeviceId                string                    `json:"device_id"`
	DeviceConfigId          string                    `json:"device_config_id"`
	AutomateExecteSceeInfos []AutomateExecteSceneInfo `json:"automate_execte_scene_infos"`
}

type AutomateDeviceInfos []AutomateDeviceInfo

// UnmarshalBinary(data []byte) error
func (a *AutomateDeviceInfos) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, a)
}

type DTConditions []model.DeviceTriggerCondition

func (a *DTConditions) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, a)
}

//AutomateActionInfo

func (a *AutomateActionInfo) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, a)
}
