package initialize

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	global "github.com/HustIoTPlatform/backend/pkg/global"

	pkgerrors "github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

var (
	alarmCache *AlarmCache
	alarmMu    sync.Mutex
)

type AlarmCache struct {
	client   *redis.Client
	expireIn time.Duration
}

// Cache 1: alarm_group_id – Stores the triggered alarm's config_id and device_id, using the scene group ID as the key.
// Cache 2: alarm_device_id – Stores the group ID, using the device ID as the key.
// Cache 3: alarm_config_id – Stores the group ID, using the alarm config ID as the key.
// Cache 4: scene_automation_id – Stores the group ID, using the scene automation ID as the key.

func NewAlarmCache() *AlarmCache {
	alarmMu.Lock()
	defer alarmMu.Unlock()
	if alarmCache == nil {
		alarmCache = &AlarmCache{
			client:   global.REDIS,
			expireIn: time.Hour * 24 * 6,
		}
	}
	return alarmCache
}

//	{
//	    "scene_automation_id":"xxx",
//	    "alarm_config_id_list": ["xxx","xxx"],
//	    "alarm_device_id_list":["xxx"]
//	}
type AlarmCacheGroup struct {
	SceneAutomationId  string   `json:"scene_automation_id"`
	AlarmConfigIdList  []string `json:"alarm_config_id_list"`
	AlaramDeviceIdList []string `json:"alaram_device_id_list"`
	Contents           []string `json:"contents"`
}

func (a *AlarmCacheGroup) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, a)
}

type SliceString []string

func (a *SliceString) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, a)
}

func (*AlarmCache) getCacheKeyByGroupId(group_id string) string {
	return fmt.Sprintf("alarm_cache_group_v5_%s", group_id)
}

func (*AlarmCache) getCacheKeyByDevice(device_id string) string {
	return fmt.Sprintf("alarm_cach_device_v5_%s", device_id)
}

func (*AlarmCache) getCacheKeyByAlarm(alarm_config_id string) string {
	return fmt.Sprintf("alarm_cach_alarm_v5_%s", alarm_config_id)
}

func (*AlarmCache) getCacheKeyByScene(scene_automation_id string) string {
	return fmt.Sprintf("alarm_cach_scene_v5_%s", scene_automation_id)
}

func (a *AlarmCache) set(key string, value interface{}) error {
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
	logrus.Debug(valueStr)
	return a.client.Set(context.Background(), key, valueStr, a.expireIn).Err()
}

// SetDevice
func (a *AlarmCache) SetDevice(group_id, scene_automation_id string, device_ids, contents []string) error {
	alarmMu.Lock()
	defer alarmMu.Unlock()
	var info AlarmCacheGroup
	cacheKey := a.getCacheKeyByGroupId(group_id)
	if count, err := a.client.Exists(context.Background(), cacheKey).Result(); err != nil {
		return pkgerrors.Wrap(err, "Failed to check if cache exists")
	} else if count > 0 {
		err := a.client.Get(context.Background(), cacheKey).Scan(&info)
		if err != nil {
			return pkgerrors.Wrap(err, "Failed to retrieve cache")
		}
		info.Contents = contents
	} else {
		info = AlarmCacheGroup{
			SceneAutomationId:  scene_automation_id,
			AlaramDeviceIdList: device_ids,
			Contents:           contents,
		}
	}
	logrus.Debugf("AlarmCacheGroupSet:%#v", info)
	err := a.set(cacheKey, info)
	if err != nil {
		return err
	}
	for _, device_id := range device_ids {
		cacheKey = a.getCacheKeyByDevice(device_id)
		err = a.groupCacheAdd(cacheKey, group_id)
		if err != nil {
			return err
		}
	}
	cacheKey = a.getCacheKeyByScene(scene_automation_id)
	logrus.Debug("SetDevice:", cacheKey, "==>", group_id)
	return a.groupCacheAdd(cacheKey, group_id)
}

func (a *AlarmCache) groupCacheAdd(cacheKey, groupId string) error {
	var groupIds SliceString
	err := a.client.Get(context.Background(), cacheKey).Scan(&groupIds)
	if err != nil && !errors.Is(err, redis.Nil) {
		return err
	}
	var isOk bool
	for _, g := range groupIds {
		if g == groupId {
			isOk = true
			break
		}
	}

	if isOk {
		return nil
	}
	groupIds = append(groupIds, groupId)
	err = a.set(cacheKey, groupIds)
	if err != nil {
		return err
	}
	return nil
}

func (a *AlarmCache) groupCacheDel(cachekey, group_id string) error {
	var groupIds SliceString
	err := a.client.Get(context.Background(), cachekey).Scan(&groupIds)
	if err != nil && err != redis.Nil {
		return err
	}
	for i, g := range groupIds {
		if g == group_id {
			groupIds = append(groupIds[:i], groupIds[i+1:]...)
		}
	}
	if len(groupIds) > 0 {
		err = a.set(cachekey, groupIds)
	} else {
		err = a.client.Del(context.Background(), cachekey).Err()
	}

	if err != nil {
		return err
	}
	return nil
}

func (a *AlarmCache) SetAlarm(group_id string, alarm_config_ids []string) error {
	alarmMu.Lock()
	defer alarmMu.Unlock()
	var info AlarmCacheGroup
	cachekey := a.getCacheKeyByGroupId(group_id)
	err := a.client.Get(context.Background(), cachekey).Scan(&info)
	if err != nil && err != redis.Nil {
		return err
	}
	info.AlarmConfigIdList = alarm_config_ids
	err = a.set(cachekey, info)
	if err != nil {
		return err
	}
	for _, alarm_id := range alarm_config_ids {
		cachekey = a.getCacheKeyByAlarm(alarm_id)
		err = a.groupCacheAdd(cachekey, group_id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *AlarmCache) GetByGroupId(group_id string) (AlarmCacheGroup, error) {
	var info AlarmCacheGroup
	cachekey := a.getCacheKeyByGroupId(group_id)
	err := a.client.Get(context.Background(), cachekey).Scan(&info)
	if err != nil && err != redis.Nil {
		return info, err
	}
	return info, nil
}

func (a *AlarmCache) GetBySceneAutomationId(scene_automation_id string) ([]string, error) {
	var groupIds SliceString
	cachekey := a.getCacheKeyByScene(scene_automation_id)
	err := a.client.Get(context.Background(), cachekey).Scan(&groupIds)
	if err != nil && err != redis.Nil {
		return groupIds, err
	}
	return groupIds, nil
}

func (a *AlarmCache) DeleteBygroupId(group_Id string) error {
	alarmMu.Lock()
	defer alarmMu.Unlock()
	info, err := a.GetByGroupId(group_Id)
	if err != nil {
		return err
	}
	for _, alarmId := range info.AlarmConfigIdList {
		cachekey := a.getCacheKeyByAlarm(alarmId)
		err = a.groupCacheDel(cachekey, group_Id)
		if err != nil {
			return err
		}
	}
	for _, deviceId := range info.AlaramDeviceIdList {
		cachekey := a.getCacheKeyByDevice(deviceId)
		err = a.groupCacheDel(cachekey, group_Id)
		if err != nil {
			return err
		}
	}
	cachekey := a.getCacheKeyByScene(info.SceneAutomationId)
	err = a.groupCacheDel(cachekey, group_Id)
	if err != nil {
		return err
	}

	cacheKey := a.getCacheKeyByGroupId(group_Id)

	return a.client.Del(context.Background(), cacheKey).Err()
}
