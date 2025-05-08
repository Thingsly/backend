package service

import (
	"github.com/Thingsly/backend/initialize"
	model "github.com/Thingsly/backend/internal/model"

	pkgerrors "github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// ActionAfterAlarm
// @description
// param actions []model.ActionInfo
// @return error
func ActionAfterAlarm(actions []model.ActionInfo, actionResultErr error) error {

	scene_automation_id := actions[0].SceneAutomationID
	alarmCache := initialize.NewAlarmCache()
	groupIds, err := alarmCache.GetBySceneAutomationId(scene_automation_id)
	if err != nil {
		return pkgerrors.Wrap(err, "Failed to get cache 1")
	}
	if len(groupIds) == 0 {
		return nil
	}
	logrus.Debug("ActionAfterAlarm:", groupIds)
	var alarm_config_ids []string

	for _, act := range actions {
		if act.ActionType == model.AUTOMATE_ACTION_TYPE_ALARM && act.ActionTarget != nil && *act.ActionTarget != "" {
			alarm_config_ids = append(alarm_config_ids, *act.ActionTarget)
		}
	}
	for _, group_id := range groupIds {

		if len(alarm_config_ids) == 0 {
			err = alarmCache.DeleteBygroupId(group_id)
		} else if actionResultErr == nil {
			err = alarmCache.SetAlarm(group_id, alarm_config_ids)
		}
		if err != nil {
			return pkgerrors.Wrap(err, "Cache deletion or setting failed.")
		}
	}

	return nil
}

// ConditionAfterAlarm
// @description
// param ok bool
// param conditions initialize.DTConditions
// param deviceId string
// @return error
func ConditionAfterAlarm(ok bool, conditions initialize.DTConditions, deviceId string, contents []string) error {
	var (
		device_ids          []string
		group_id            string
		scene_automation_id string
		alarmCache          = initialize.NewAlarmCache()
	)
	for _, cond := range conditions {
		group_id = cond.GroupID
		scene_automation_id = cond.SceneAutomationID
		if cond.TriggerConditionType == model.DEVICE_TRIGGER_CONDITION_TYPE_ONE {
			device_ids = append(device_ids, *cond.TriggerSource)
		}
		if cond.TriggerConditionType == model.DEVICE_TRIGGER_CONDITION_TYPE_MULTIPLE {
			device_ids = append(device_ids, deviceId)
		}
	}
	logrus.Debug("ConditionAfterAlarm:", group_id, device_ids, ok, contents)
	if len(device_ids) == 0 {
		return nil
	}

	//alarmCache.DeleteBygroupId(group_id)

	if ok {
		err := alarmCache.SetDevice(group_id, scene_automation_id, device_ids, contents)
		if err != nil {
			return pkgerrors.Wrap(err, "Cache setting failed.")
		}
		groupIds, _ := alarmCache.GetBySceneAutomationId(scene_automation_id)
		logrus.Debug("getGroupId:", groupIds)
	} else {

		err := AlarmRecovery(group_id, contents)
		if err != nil {
			return pkgerrors.WithMessage(err, "Failed to restore the alarm")
		}

		err = alarmCache.DeleteBygroupId(group_id)
		if err != nil {
			return pkgerrors.Wrap(err, "Failed to set the cache")
		}
		c, _ := alarmCache.GetByGroupId(group_id)
		logrus.Debug("Query after deletion: ", c)
	}
	return nil
}

// AlarmExecute
// @description
// param alarm_config_id string
// param scene_automation_id
// @return bool
func AlarmExecute(alarm_config_id, scene_automation_id string) (bool, string, string) {
	var (
		alarmName string
		resultOk  bool
		reason    string
	)

	alarmCache := initialize.NewAlarmCache()
	groupIds, err := alarmCache.GetBySceneAutomationId(scene_automation_id)
	logrus.Debugf("Cache 11: %#v, Scene ID: %#v", groupIds, scene_automation_id)
	if err != nil || len(groupIds) == 0 {
		reason = "Alarm cache does not exist"
		return resultOk, alarmName, reason
	}
	for _, group_id := range groupIds {
		cache, err := alarmCache.GetByGroupId(group_id)
		if err != nil {
			reason = "Alarm cache does not exist"
			return resultOk, alarmName, reason
		}
		logrus.Debugf("Query before alarm execution: %#v", cache)
		var isOk bool
		for _, acid := range cache.AlarmConfigIdList {
			if acid == alarm_config_id {
				isOk = true
				break
			}
		}
		if isOk {
			// Alarm ID already exists in cache, indicating the alarm has been triggered before
			reason = "Alarm already exists"
			logrus.Debugf("Alarm already exists in cache, skipping execution")
			continue
		}
		var content string
		content = "Scene automation triggers alarm"
		for _, strval := range cache.Contents {
			content += ";" + strval
		}
		resultOk, alarmName, reason = GroupApp.AlarmExecute(alarm_config_id, content, scene_automation_id, group_id, cache.AlaramDeviceIdList)
	}
	return resultOk, alarmName, reason
}

// AlarmRecovery
// @description
// param group_id
func AlarmRecovery(group_id string, contents []string) error {
	alarmCache := initialize.NewAlarmCache()
	cache, err := alarmCache.GetByGroupId(group_id)
	if err != nil {
		return err
	}
	logrus.Debug("AlarmRecovery:cache:", cache)
	for _, acid := range cache.AlarmConfigIdList {
		var content string
		content = "Scene automation restores alarm"
		for _, strval := range contents {
			content += ";" + strval
		}
		GroupApp.AlarmRecovery(acid, content, cache.SceneAutomationId, group_id, cache.AlaramDeviceIdList)

	}

	return nil
}
