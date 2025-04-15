package service

import (
	"github.com/Thingsly/backend/initialize"
	"github.com/Thingsly/backend/internal/dal"
	model "github.com/Thingsly/backend/internal/model"
	"github.com/Thingsly/backend/pkg/errcode"
	utils "github.com/Thingsly/backend/pkg/utils"

	"github.com/go-basic/uuid"
	"github.com/sirupsen/logrus"
)

type SceneAutomation struct{}

func (s *SceneAutomation) CreateSceneAutomation(req *model.CreateSceneAutomationReq, u *utils.UserClaims) (string, error) {

	var sceneAutomationID string

	// Start transaction
	logrus.Info("Starting transaction")
	tx, err := dal.StartTransaction()
	if err != nil {
		return sceneAutomationID, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}

	// Insert scene_automation record
	var sceneAutomation = model.SceneAutomation{}
	sceneAutomation.ID = uuid.New()

	sceneAutomationID = sceneAutomation.ID

	sceneAutomation.Name = req.Name
	sceneAutomation.Description = &req.Description
	//sceneAutomation.Enabled = req.Enabled
	sceneAutomation.Enabled = "N"
	sceneAutomation.TenantID = u.TenantID
	sceneAutomation.Creator = u.ID
	sceneAutomation.Updator = u.ID
	sceneAutomation.CreatedAt = utils.GetUTCTime()
	sceneAutomation.UpdatedAt = &sceneAutomation.CreatedAt
	sceneAutomation.Remark = &req.Remark
	// Create the scene automation
	logrus.Info("Creating scene automation record")
	err = dal.CreateSceneAutomation(&sceneAutomation, tx)
	if err != nil {
		dal.Rollback(tx)
		return "", errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}

	// Handle trigger condition groups
	for _, v := range req.TriggerConditionGroups {
		groupID := uuid.New()
		var (
			oneCondition      bool
			multipleCondition bool
		)
		for _, v2 := range v {
			switch v2.TriggerConditionsType {
			case "10", "11", "22":
				if v2.TriggerConditionsType == "10" {
					oneCondition = true
				}
				if v2.TriggerConditionsType == "11" {
					multipleCondition = true
				}
				// Insert device trigger condition
				var dtc = model.DeviceTriggerCondition{}
				dtc.ID = uuid.New()
				dtc.SceneAutomationID = sceneAutomationID
				dtc.Enabled = req.Enabled
				dtc.GroupID = groupID
				dtc.TriggerConditionType = v2.TriggerConditionsType
				dtc.TriggerSource = v2.TriggerSource
				dtc.TriggerParamType = v2.TriggerParamType
				dtc.TriggerParam = v2.TriggerParam
				dtc.TriggerOperator = v2.TriggerOperator
				if v2.TriggerValue != nil {
					dtc.TriggerValue = *v2.TriggerValue
				}
				dtc.Enabled = req.Enabled
				dtc.TenantID = u.TenantID
				// Create device trigger condition
				logrus.Info("Creating device trigger condition")
				err = dal.CreateDeviceTriggerCondition(dtc, tx)
				if err != nil {
					dal.Rollback(tx)
					return "", errcode.WithData(errcode.CodeDBError, map[string]interface{}{
						"sql_error": err.Error(),
					})
				}

			case "20":
				// Insert one-time task
				var ott = model.OneTimeTask{}
				ott.ID = uuid.New()
				ott.SceneAutomationID = sceneAutomationID
				if v2.ExecutionTime != nil {
					ott.ExecutionTime = *v2.ExecutionTime
				}
				if v2.ExpirationTime != nil {
					ott.ExpirationTime = int64(*v2.ExpirationTime)
				}
				ott.ExecutingState = "NEX"
				ott.Enabled = req.Enabled
				// Create one-time task
				logrus.Info("Creating one-time task")
				err = dal.CreateOneTimeTask(ott, tx)
				if err != nil {
					dal.Rollback(tx)
					return "", errcode.WithData(errcode.CodeDBError, map[string]interface{}{
						"sql_error": err.Error(),
					})
				}
			case "21":
				// Insert periodic task
				var pt = model.PeriodicTask{}
				pt.ID = uuid.New()
				pt.SceneAutomationID = sceneAutomationID
				if v2.TaskType != nil {
					pt.TaskType = *v2.TaskType
				}
				if v2.Params != nil {
					pt.Param = *v2.Params
				}
				if v2.ExecutionTime != nil {
					pt.ExecutionTime = *v2.ExecutionTime
				}
				if v2.ExpirationTime != nil {
					pt.ExpirationTime = int64(*v2.ExpirationTime)
				}
				pt.Enabled = "Y"
				// Create periodic task
				logrus.Info("Creating periodic task")
				err = dal.CreatePeriodicTask(pt, tx)
				if err != nil {
					dal.Rollback(tx)
					return "", errcode.WithData(errcode.CodeDBError, map[string]interface{}{
						"sql_error": err.Error(),
					})
				}
			default:
				dal.Rollback(tx)
				return "", errcode.WithData(errcode.CodeDBError, map[string]interface{}{
					"sql_error": err.Error(),
				})
			}

		}
		if oneCondition && multipleCondition {
			dal.Rollback(tx)
			return "", errcode.New(200060) // Use a standard error code
		}
	}

	// Handle actions
	for _, v := range req.Actions {
		// Insert action info
		var actionInfo = model.ActionInfo{}
		actionInfo.ID = uuid.New()
		actionInfo.SceneAutomationID = sceneAutomationID
		actionInfo.ActionTarget = &v.ActionTarget
		actionInfo.ActionType = v.ActionType
		actionInfo.ActionParamType = &v.ActionParamType
		actionInfo.ActionParam = &v.ActionParam
		actionInfo.ActionValue = &v.ActionValue
		// Create action info
		logrus.Info("Creating action info")
		err = dal.CreateActionInfo(actionInfo, tx)
		if err != nil {
			dal.Rollback(tx)
			return "", errcode.WithData(errcode.CodeDBError, map[string]interface{}{
				"sql_error": err.Error(),
			})
		}

	}

	// Commit transaction
	dal.Commit(tx)
	// Save automation cache information
	go func() {
		if req.Enabled == "Y" {
			err := s.AutomateCacheSet(sceneAutomationID)
			if err != nil {
				logrus.Error("Failed to save automation cache information for the new scene automation, err:", err)
			}
		}
	}()
	return sceneAutomationID, nil
}

// AutomateCacheSet saves automation cache information
func (*SceneAutomation) AutomateCacheSet(sceneAutomationID string) error {
	logrus.Info("Starting to save automation cache information")

	// Fetch device trigger conditions for the scene automation
	groupInfoPtrs, err := dal.GetDeviceTriggerCondition(sceneAutomationID)
	if err != nil {
		return errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}

	// Fetch action information for the scene automation
	actionInfoPtrs, err := dal.GetActionInfo(sceneAutomationID)
	if err != nil {
		return errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}

	// Filter out disabled groupInfos and collect enabled ones
	var groupInfos []model.DeviceTriggerCondition
	for _, groupInfo := range groupInfoPtrs {
		if groupInfo != nil && groupInfo.Enabled == "Y" {
			groupInfos = append(groupInfos, *groupInfo)
		}
	}

	// Collect actionInfos
	var actionInfos []model.ActionInfo
	for _, actionInfo := range actionInfoPtrs {
		if actionInfo != nil {
			actionInfos = append(actionInfos, *actionInfo)
		}
	}

	// Set cache by scene automation ID
	err = initialize.NewAutomateCache().SetCacheBySceneAutomationId(sceneAutomationID, groupInfos, actionInfos)
	if err != nil {
		return errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}

	return nil
}

func (*SceneAutomation) DeleteSceneAutomation(scene_automation_id string) error {
	err := dal.DeleteSceneAutomation(scene_automation_id, nil)
	if err != nil {
		return errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}
	return nil
}

func (*SceneAutomation) GetSceneAutomation(scene_automation_id string) (interface{}, error) {
	sceneAutomation, err := dal.GetSceneAutomation(scene_automation_id, nil)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}

	deviceTriggerCondition, err := dal.GetDeviceTriggerCondition(scene_automation_id)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}

	oneTimeTask, err := dal.GetOneTimeTask(scene_automation_id)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}

	periodicTask, err := dal.GetPeriodicTask(scene_automation_id)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}

	actionInfo, err := dal.GetActionInfo(scene_automation_id)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}

	res := make(map[string]interface{})
	res["id"] = sceneAutomation.ID
	res["name"] = sceneAutomation.Name
	res["description"] = sceneAutomation.Description
	res["enabled"] = sceneAutomation.Enabled
	res["tenant_id"] = sceneAutomation.TenantID
	res["creator"] = sceneAutomation.Creator
	res["updator"] = sceneAutomation.Updator

	triggerConditionGroups := make([][]map[string]interface{}, 0)

	if len(periodicTask) > 0 {
		tmp := make([][]map[string]interface{}, 0)
		for _, v := range periodicTask {
			mapList := make([]map[string]interface{}, 0)
			periodicTaskMap := make(map[string]interface{})
			periodicTaskMap["task_type"] = v.TaskType
			periodicTaskMap["expiration_time"] = v.ExpirationTime
			periodicTaskMap["params"] = v.Param
			periodicTaskMap["trigger_conditions_type"] = "21"
			mapList = append(mapList, periodicTaskMap)
			tmp = append(tmp, mapList)
		}
		triggerConditionGroups = append(triggerConditionGroups, tmp...)
	}

	if len(oneTimeTask) > 0 {
		tmp := make([][]map[string]interface{}, 0)
		for _, v := range oneTimeTask {
			mapList := make([]map[string]interface{}, 0)
			oneTimeTaskMap := make(map[string]interface{})
			oneTimeTaskMap["execution_time"] = v.ExecutionTime
			oneTimeTaskMap["expiration_time"] = v.ExpirationTime
			oneTimeTaskMap["trigger_conditions_type"] = "20"
			mapList = append(mapList, oneTimeTaskMap)
			tmp = append(tmp, mapList)
		}
		triggerConditionGroups = append(triggerConditionGroups, tmp...)
	}

	if len(deviceTriggerCondition) > 0 {
		tmp := make([][]map[string]interface{}, 0)

		rebuild := make(map[string][]*model.DeviceTriggerCondition)
		for _, v := range deviceTriggerCondition {
			rebuild[v.GroupID] = append(rebuild[v.GroupID], v)
		}

		for _, v := range rebuild {
			mapList := make([]map[string]interface{}, 0)
			for _, v2 := range v {
				deviceTriggerConditionMap := make(map[string]interface{})
				deviceTriggerConditionMap["id"] = v2.ID
				deviceTriggerConditionMap["group_id"] = v2.GroupID
				deviceTriggerConditionMap["trigger_conditions_type"] = v2.TriggerConditionType
				deviceTriggerConditionMap["trigger_source"] = v2.TriggerSource
				deviceTriggerConditionMap["trigger_param_type"] = v2.TriggerParamType
				deviceTriggerConditionMap["trigger_param"] = v2.TriggerParam
				deviceTriggerConditionMap["trigger_operator"] = v2.TriggerOperator
				deviceTriggerConditionMap["trigger_value"] = v2.TriggerValue
				mapList = append(mapList, deviceTriggerConditionMap)
			}
			tmp = append(tmp, mapList)
		}
		triggerConditionGroups = append(triggerConditionGroups, tmp...)

	}

	res["trigger_condition_groups"] = triggerConditionGroups

	if len(actionInfo) > 0 {
		actionInfoMap := make([]map[string]interface{}, 0)
		for _, v := range actionInfo {
			tmp := make(map[string]interface{})
			tmp["action_type"] = v.ActionType
			tmp["action_target"] = v.ActionTarget
			tmp["action_param_type"] = v.ActionParamType
			tmp["action_param"] = v.ActionParam
			tmp["action_value"] = v.ActionValue
			actionInfoMap = append(actionInfoMap, tmp)
		}
		res["actions"] = actionInfoMap
	}

	return res, err
}

func (*SceneAutomation) SwitchSceneAutomation(scene_automation_id, target string) error {

	tx, err := dal.StartTransaction()
	if err != nil {
		return errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}

	if target == "" {
		data, err := dal.GetSceneAutomation(scene_automation_id, tx)
		if err != nil {
			dal.Rollback(tx)
			return err
		}
		if data.Enabled == "Y" {
			target = "N"
		} else {
			target = "Y"
		}
	}

	err = dal.SwitchSceneAutomation(scene_automation_id, target, tx)
	if err != nil {
		dal.Rollback(tx)
		return errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}

	err = dal.SwitchDeviceTriggerCondition(scene_automation_id, target, tx)
	if err != nil {
		dal.Rollback(tx)
		return errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}

	err = dal.SwitchOneTimeTask(scene_automation_id, target, tx)
	if err != nil {
		dal.Rollback(tx)
		return errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}

	err = dal.SwitchPeriodicTask(scene_automation_id, target, tx)
	if err != nil {
		dal.Rollback(tx)
		return errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}

	dal.Commit(tx)

	go func() {
		// If the target is "Y", save automation cache information
		// Uncomment the following block if required:
		// if target == "Y" {
		//     err = s.AutomateCacheSet(scene_automation_id)
		//     if err != nil {
		//         logrus.Error("Failed to save automation cache information while editing scene automation, error: ", err)
		//     }
		// }

		// If the target is "N", delete the automation cache
		if target == "N" {
			err := initialize.NewAutomateCache().DeleteCacheBySceneAutomationId(scene_automation_id)
			if err != nil {
				logrus.Error("Failed to delete automation cache while editing, error: ", err)
			}
		}
	}()

	return nil
}

func (*SceneAutomation) GetSceneAutomationByPageReq(req *model.GetSceneAutomationByPageReq, u *utils.UserClaims) (interface{}, error) {
	total, sceneInfo, err := dal.GetSceneAutomationByPage(req, u.TenantID)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}
	sceneListMap := make(map[string]interface{})
	sceneListMap["total"] = total
	sceneListMap["list"] = sceneInfo
	return sceneListMap, nil
}

func (*SceneAutomation) GetSceneAutomationWithAlarmByPageReq(req *model.GetSceneAutomationsWithAlarmByPageReq, u *utils.UserClaims) (interface{}, error) {
	total, sceneInfo, err := dal.GetSceneAutomationWithAlarmByPageReq(req, u.TenantID)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}
	sceneListMap := make(map[string]interface{})
	sceneListMap["total"] = total
	sceneListMap["list"] = sceneInfo
	return sceneListMap, nil
}

func (*SceneAutomation) UpdateSceneAutomation(req *model.UpdateSceneAutomationReq, u *utils.UserClaims) (string, error) {

	var scene_automation_id string

	tx, err := dal.StartTransaction()
	if err != nil {
		return scene_automation_id, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}

	scene_automation_id = req.ID
	t := utils.GetUTCTime()

	var sceneAutomation = model.SceneAutomation{}
	sceneAutomation.ID = scene_automation_id
	sceneAutomation.Name = req.Name
	sceneAutomation.Description = &req.Description
	sceneAutomation.Enabled = req.Enabled
	sceneAutomation.TenantID = u.TenantID
	sceneAutomation.Updator = u.ID
	sceneAutomation.UpdatedAt = &t
	sceneAutomation.Remark = &req.Remark

	err = dal.SaveSceneAutomation(&sceneAutomation, tx)
	if err != nil {
		return "", errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}

	err = dal.DeleteDeviceTriggerCondition(scene_automation_id, tx)
	if err != nil {
		dal.Rollback(tx)
		return "", errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}

	err = dal.DeleteOneTimeTask(scene_automation_id, tx)
	if err != nil {
		dal.Rollback(tx)
		return "", errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}

	err = dal.DeletePeriodicTask(scene_automation_id, tx)
	if err != nil {
		dal.Rollback(tx)
		return "", errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}

	err = dal.DeleteActionInfo(scene_automation_id, tx)
	if err != nil {
		dal.Rollback(tx)
		return "", errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}

	for _, v := range req.TriggerConditionGroups {
		groupId := uuid.New()
		var (
			oneCondition      bool
			multipleCondition bool
		)
		for _, v2 := range v {
			if v2.TriggerConditionsType == "10" {
				oneCondition = true
			}
			if v2.TriggerConditionsType == "11" {
				multipleCondition = true
			}
			switch v2.TriggerConditionsType {
			case "10", "11", "22":

				var dtc = model.DeviceTriggerCondition{}
				dtc.ID = uuid.New()
				dtc.SceneAutomationID = scene_automation_id
				dtc.Enabled = req.Enabled
				dtc.GroupID = groupId
				dtc.TriggerConditionType = v2.TriggerConditionsType
				dtc.TriggerSource = v2.TriggerSource
				dtc.TriggerParamType = v2.TriggerParamType
				dtc.TriggerParam = v2.TriggerParam
				dtc.TriggerOperator = v2.TriggerOperator
				if v2.TriggerValue != nil {
					dtc.TriggerValue = *v2.TriggerValue
				}
				dtc.Enabled = req.Enabled
				dtc.TenantID = u.TenantID

				logrus.Info("Creating device trigger condition")
				err = dal.CreateDeviceTriggerCondition(dtc, tx)
				if err != nil {
					dal.Rollback(tx)
					return "", errcode.WithData(errcode.CodeDBError, map[string]interface{}{
						"sql_error": err.Error(),
					})
				}

			case "20":

				var ott = model.OneTimeTask{}
				ott.ID = uuid.New()
				ott.SceneAutomationID = scene_automation_id
				if v2.ExecutionTime != nil {
					ott.ExecutionTime = *v2.ExecutionTime
				}

				// if v2.ExecutionTime != nil {
				// 	orgTime := *v2.ExecutionTime
				// 	intOrgTime := orgTime.Unix()
				// 	ott.ExpirationTime = intOrgTime
				// }

				if v2.ExpirationTime != nil {
					ott.ExpirationTime = int64(*v2.ExpirationTime)
				}
				ott.ExecutingState = "NEX"
				ott.Enabled = req.Enabled
				logrus.Info("Creating one-time task information")
				err = dal.CreateOneTimeTask(ott, tx)
				if err != nil {
					dal.Rollback(tx)
					return "", errcode.WithData(errcode.CodeDBError, map[string]interface{}{
						"sql_error": err.Error(),
					})
				}
			case "21":
				var pt = model.PeriodicTask{}
				pt.ID = uuid.New()
				pt.SceneAutomationID = scene_automation_id
				if v2.TaskType != nil {
					pt.TaskType = *v2.TaskType
				}
				if v2.Params != nil {
					pt.Param = *v2.Params
				}
				if v2.ExpirationTime != nil {
					pt.ExpirationTime = int64(*v2.ExpirationTime)
				}
				// pt.ExecutionTime = *v2.ExecutionTime
				// pt.ExpirationTime = *v2.ExpirationTime
				pt.Enabled = "Y"
				logrus.Info("Creating periodic task information")
				err = dal.CreatePeriodicTask(pt, tx)
				if err != nil {
					dal.Rollback(tx)
					return "", errcode.WithData(errcode.CodeDBError, map[string]interface{}{
						"sql_error": err.Error(),
					})
				}
			default:
				dal.Rollback(tx)
				return "", errcode.WithData(errcode.CodeParamError, map[string]interface{}{
					"err": "not support trigger type",
				})
			}
		}
		if oneCondition && multipleCondition {
			dal.Rollback(tx)
			return "", errcode.New(200060)
		}
	}

	for _, v := range req.Actions {
		var actionInfo = model.ActionInfo{}
		actionInfo.ID = uuid.New()
		actionInfo.SceneAutomationID = scene_automation_id
		actionInfo.ActionTarget = &v.ActionTarget
		actionInfo.ActionType = v.ActionType
		actionInfo.ActionParamType = &v.ActionParamType
		actionInfo.ActionParam = &v.ActionParam
		actionInfo.ActionValue = &v.ActionValue
		err = dal.CreateActionInfo(actionInfo, tx)
		if err != nil {
			dal.Rollback(tx)
			return "", errcode.WithData(errcode.CodeDBError, map[string]interface{}{
				"sql_error": err.Error(),
			})
		}

	}

	dal.Commit(tx)

	// Commit transaction
	go func() {
		err1 := initialize.NewAutomateCache().DeleteCacheBySceneAutomationId(scene_automation_id)
		if err1 != nil {
			logrus.Error("Failed to delete automation cache while editing, error: ", err)
		}
		//data, index, _ := initialize.NewAutomateCache().GetCacheByDeviceId("2dffdf60-f937-8d60-b141-6faf9935a7ab", "")
		//logrus.Info(" ,data:", data, ";index:", index)
		//data, index, _ = initialize.NewAutomateCache().GetCacheByDeviceId("2dffdf60-f937-8d60-b141-6faf9935a7ab", "903aa8a2-e03b-9ab1-10b8-87d82e1a6216")
		//logrus.Info(" ,data:", data, ";index:", index)
		//if req.Enabled == "Y" {
		//	err1 = s.AutomateCacheSet(scene_automation_id)
		//	if err1 != nil {
		//		logrus.Error("Fail to set cache ，err: ", err)
		//	}
		//}
	}()
	//time.Sleep(time.Second * 5)
	return scene_automation_id, nil
}
