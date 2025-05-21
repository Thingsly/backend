package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Thingsly/backend/internal/dal"
	model "github.com/Thingsly/backend/internal/model"
	"github.com/Thingsly/backend/pkg/errcode"

	"github.com/go-basic/uuid"
	"github.com/sirupsen/logrus"
)

type Alarm struct{}

// CreateAlarmConfig
func (*Alarm) CreateAlarmConfig(req *model.CreateAlarmConfigReq) (data *model.AlarmConfig, err error) {
	data = &model.AlarmConfig{}
	t := time.Now().UTC()
	data.ID = uuid.New()
	data.Name = req.Name
	data.Description = req.Description
	data.AlarmLevel = req.AlarmLevel
	data.NotificationGroupID = req.NotificationGroupID
	data.CreatedAt = t
	data.UpdatedAt = t
	data.TenantID = req.TenantID
	data.Remark = req.Remark
	data.Enabled = req.Enabled

	err = dal.CreateAlarmConfig(data)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}
	return
}

// DeleteAlarmConfig
func (*Alarm) DeleteAlarmConfig(id string) (err error) {
	err = dal.DeleteAlarmConfig(id)
	if err != nil {
		return errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}
	return
}

// UpdateAlarmConfig
func (*Alarm) UpdateAlarmConfig(req *model.UpdateAlarmConfigReq) (data *model.AlarmConfig, err error) {
	data = &model.AlarmConfig{}
	data.ID = req.ID
	if req.Name != nil {
		data.Name = *req.Name
	}
	if req.Description != nil {
		data.Description = req.Description
	}
	if req.AlarmLevel != nil {
		data.AlarmLevel = *req.AlarmLevel
	}
	if req.NotificationGroupID != nil {
		data.NotificationGroupID = *req.NotificationGroupID
	}
	data.UpdatedAt = time.Now().UTC()
	data.TenantID = *req.TenantID
	data.Remark = req.Remark
	if req.Enabled != nil {
		data.Enabled = *req.Enabled
	}

	err = dal.UpdateAlarmConfig(data)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}
	data, err = dal.GetAlarmByID(req.ID)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}
	return data, nil
}

// GetAlarmConfigListByPage
func (*Alarm) GetAlarmConfigListByPage(req *model.GetAlarmConfigListByPageReq) (data map[string]interface{}, err error) {

	total, list, err := dal.GetAlarmConfigListByPage(req)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}
	data = make(map[string]interface{})
	data["total"] = total
	data["list"] = list
	return
}

// UpdateAlarmInfo
func (*Alarm) UpdateAlarmInfo(req *model.UpdateAlarmInfoReq, userid string) (alarmInfo *model.AlarmInfo, err error) {

	alarmInfo, err = dal.GetAlarmInfoByID(req.Id)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}
	alarmInfo.Processor = &userid
	if req.ProcessingResult != nil && *req.ProcessingResult != "" {
		alarmInfo.ProcessingResult = *req.ProcessingResult
	}
	err = dal.UpdateAlarmInfo(alarmInfo)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}
	return
}

// UpdateAlarmInfoBatch
func (*Alarm) UpdateAlarmInfoBatch(req *model.UpdateAlarmInfoBatchReq, userid string) error {
	if len(req.Id) == 0 {
		return errcode.WithData(errcode.CodeParamError, map[string]interface{}{
			"id": "id is empty",
		})
	}
	err := dal.UpdateAlarmInfoBatch(req, userid)
	if err != nil {
		return errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}
	return err
}

// GetAlarmInfoListByPage
func (*Alarm) GetAlarmInfoListByPage(req *model.GetAlarmInfoListByPageReq) (data map[string]interface{}, err error) {

	total, list, err := dal.GetAlarmInfoListByPage(req)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}
	data = make(map[string]interface{})
	data["total"] = total
	data["list"] = list
	return
}

// GetAlarmHisttoryListByPage
func (*Alarm) GetAlarmHisttoryListByPage(req *model.GetAlarmHisttoryListByPage, tenantID string) (data map[string]interface{}, err error) {

	total, list, err := dal.GetAlarmHistoryListByPage(req, tenantID)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}
	data = make(map[string]interface{})
	data["total"] = total
	data["list"] = list
	return
}
func (*Alarm) AlarmHistoryDescUpdate(req *model.AlarmHistoryDescUpdateReq, tenantID string) (err error) {
	err = dal.AlarmHistoryDescUpdate(req, tenantID)
	if err != nil {
		return errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}
	return
}
func (*Alarm) GetDeviceAlarmStatus(req *model.GetDeviceAlarmStatusReq) bool {
	return dal.GetDeviceAlarmStatus(req)
}

func (*Alarm) GetConfigByDevice(req *model.GetDeviceAlarmStatusReq) ([]model.AlarmConfig, error) {
	data, err := dal.GetConfigByDevice(req)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}
	return data, nil
}

// AddAlarmInfo
func (*Alarm) AddAlarmInfo(alarmConfigID, content string) (bool, string) {

	alarmConfig, err := dal.GetAlarmByID(alarmConfigID)
	if err != nil {
		logrus.Error(err)
		return false, ""
	}

	if alarmConfig.Enabled != "Y" {
		return false, ""
	}

	if alarmConfig.NotificationGroupID != "" {
		alarmLevelText := ""
		headerColor := ""
		switch alarmConfig.AlarmLevel {
		case "H":
			alarmLevelText = "High"
			headerColor = "#dc3545" // Red
		case "M":
			alarmLevelText = "Middle"
			headerColor = "#ffc107" // Yellow
		case "L":
			alarmLevelText = "Low"
			headerColor = "#28a745" // Green
		}
		title := fmt.Sprintf("[%s Alert] %s %s", alarmLevelText, alarmConfig.Name, time.Now().Format("2006-01-02 15:04:05"))
		formattedContent := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { 
            font-family: 'Segoe UI', Arial, sans-serif; 
            line-height: 1.6; 
            color: #2c3e50;
            background-color: #f8f9fa;
            margin: 0;
            padding: 0;
        }
        .container { 
            max-width: 650px; 
            margin: 20px auto; 
            padding: 30px;
            background-color: #ffffff;
            border-radius: 10px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        .header { 
            background-color: #f8f9fa; 
            padding: 25px; 
            border-radius: 8px; 
            margin-bottom: 30px;
            border-left: 5px solid %s;
        }
        .alert-level { 
            font-weight: bold; 
            color: %s;
            font-size: 1.2em;
        }
        .section { 
            margin-bottom: 30px;
            background-color: #ffffff;
            padding: 20px;
            border-radius: 8px;
            box-shadow: 0 1px 3px rgba(0,0,0,0.05);
        }
        .section-title { 
            font-weight: bold; 
            color: #2c3e50; 
            border-bottom: 2px solid #e9ecef; 
            padding-bottom: 10px; 
            margin-bottom: 15px;
            font-size: 1.3em;
        }
        .info-row { 
            margin: 12px 0;
            font-size: 1.1em;
        }
        .label { 
            font-weight: 600; 
            color: #495057;
            display: inline-block;
            min-width: 120px;
        }
        .value { 
            color: #2c3e50;
            font-weight: 500;
        }
        .footer { 
            margin-top: 30px; 
            font-size: 0.9em; 
            color: #6c757d; 
            border-top: 1px solid #e9ecef; 
            padding-top: 20px;
            text-align: center;
        }
        h2 {
            font-size: 1.8em;
            margin: 0;
            color: %s;
            font-weight: 600;
        }
        .recommended-actions {
            background-color: #f8f9fa;
            padding: 15px;
            border-radius: 6px;
            margin-top: 10px;
        }
        .recommended-actions .value {
            color: #495057;
            line-height: 1.8;
        }
        .view-details {
            display: inline-block;
            margin-top: 20px;
            padding: 10px 20px;
            background-color: %s;
            color: white;
            text-decoration: none;
            border-radius: 5px;
            font-weight: 500;
        }
        .view-details:hover {
            opacity: 0.9;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h2>%s</h2>
            <p style="margin: 10px 0 0 0; font-size: 1.1em;">Alert Level: <span class="alert-level">%s</span></p>
        </div>

        <div class="section">
            <div class="section-title">Alert Details</div>
            <div class="info-row">
                <span class="label">Time:</span>
                <span class="value">%s</span>
            </div>
            %s
        </div>

        <div class="section">
            <div class="section-title">Alert Information</div>
            <div class="info-row">
                <span class="value">%s</span>
            </div>
        </div>

        <div class="section">
            <div class="section-title">Recommended Actions</div>
            <div class="recommended-actions">
                <div class="info-row">
                    <span class="value">• Please verify the alert condition<br>• Check the device status<br>• Review system logs if necessary</span>
                </div>
            </div>
        </div>

        <div style="text-align: center;">
            <a href="https://dangky.app/alarm/warning-message" class="view-details">View Full Alert Details</a>
        </div>

        <div class="footer">
            <p>This is an automated message from the Thingsly IoT Platform.</p>
            <p>Please do not reply to this email.</p>
        </div>
    </div>
</body>
</html>`, headerColor, headerColor, headerColor, headerColor, alarmConfig.Name, alarmLevelText, time.Now().Format("2006-01-02 15:04:05"),
			func() string {
				if alarmConfig.Description != nil && *alarmConfig.Description != "" {
					return fmt.Sprintf(`<div class="info-row">
                <span class="label">Description:</span>
                <span class="value">%s</span>
            </div>`, *alarmConfig.Description)
				}
				return ""
			}(), content)
		GroupApp.NotificationServicesConfig.ExecuteNotification(alarmConfig.NotificationGroupID, formattedContent, title)
	}

	id := uuid.New()
	t := time.Now().UTC()
	err = dal.CreateAlarmInfo(&model.AlarmInfo{
		ID:               id,
		Name:             alarmConfig.Name,
		AlarmConfigID:    alarmConfigID,
		AlarmLevel:       &alarmConfig.AlarmLevel,
		Content:          &content,
		AlarmTime:        t,
		Description:      alarmConfig.Description,
		ProcessingResult: "UND",
		TenantID:         alarmConfig.TenantID,
	})
	if err != nil {
		logrus.Error(err)
		return false, ""
	}
	return true, id
}

func (*Alarm) AlarmRecovery(alarmConfigID, content, scene_automation_id, group_id string, device_ids []string) (bool, string) {
	alarmConfig, err := dal.GetAlarmByID(alarmConfigID)
	if err != nil {
		logrus.Error(err)
		return false, ""
	}

	device_ids_str, _ := json.Marshal(device_ids)
	id := uuid.New()
	t := time.Now().UTC()
	err = dal.AlarmHistorySave(&model.AlarmHistory{
		ID:                id,
		Name:              alarmConfig.Name,
		AlarmConfigID:     alarmConfigID,
		Content:           &content,
		Description:       alarmConfig.Description,
		TenantID:          alarmConfig.TenantID,
		SceneAutomationID: scene_automation_id,
		GroupID:           group_id,
		AlarmDeviceList:   string(device_ids_str),
		AlarmStatus:       "N",
		CreateAt:          t,
	})
	if err != nil {
		logrus.Error(err)
		return false, ""
	}
	return true, id
}

func (*Alarm) AlarmExecute(alarmConfigID, content, scene_automation_id, group_id string, device_ids []string) (bool, string, string) {
	var alarmName string
	alarmConfig, err := dal.GetAlarmByID(alarmConfigID)
	if err != nil {
		logrus.Error(err)
		return false, alarmName, ""
	}

	if alarmConfig.Enabled != "Y" {
		return false, alarmName, "alarm is not enabled"
	}
	alarmName = alarmConfig.Name
	if alarmConfig.NotificationGroupID != "" {
		alarmLevelText := ""
		headerColor := ""
		switch alarmConfig.AlarmLevel {
		case "H":
			alarmLevelText = "High"
			headerColor = "#dc3545" // Red
		case "M":
			alarmLevelText = "Middle"
			headerColor = "#ffc107" // Yellow
		case "L":
			alarmLevelText = "Low"
			headerColor = "#28a745" // Green
		}
		title := fmt.Sprintf("[%s Alert] %s %s", alarmLevelText, alarmConfig.Name, time.Now().Format("2006-01-02 15:04:05"))
		formattedContent := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { 
            font-family: 'Segoe UI', Arial, sans-serif; 
            line-height: 1.6; 
            color: #2c3e50;
            background-color: #f8f9fa;
            margin: 0;
            padding: 0;
        }
        .container { 
            max-width: 650px; 
            margin: 20px auto; 
            padding: 30px;
            background-color: #ffffff;
            border-radius: 10px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        .header { 
            background-color: #f8f9fa; 
            padding: 25px; 
            border-radius: 8px; 
            margin-bottom: 30px;
            border-left: 5px solid %s;
        }
        .alert-level { 
            font-weight: bold; 
            color: %s;
            font-size: 1.2em;
        }
        .section { 
            margin-bottom: 30px;
            background-color: #ffffff;
            padding: 20px;
            border-radius: 8px;
            box-shadow: 0 1px 3px rgba(0,0,0,0.05);
        }
        .section-title { 
            font-weight: bold; 
            color: #2c3e50; 
            border-bottom: 2px solid #e9ecef; 
            padding-bottom: 10px; 
            margin-bottom: 15px;
            font-size: 1.3em;
        }
        .info-row { 
            margin: 12px 0;
            font-size: 1.1em;
        }
        .label { 
            font-weight: 600; 
            color: #495057;
            display: inline-block;
            min-width: 120px;
        }
        .value { 
            color: #2c3e50;
            font-weight: 500;
        }
        .footer { 
            margin-top: 30px; 
            font-size: 0.9em; 
            color: #6c757d; 
            border-top: 1px solid #e9ecef; 
            padding-top: 20px;
            text-align: center;
        }
        h2 {
            font-size: 1.8em;
            margin: 0;
            color: %s;
            font-weight: 600;
        }
        .recommended-actions {
            background-color: #f8f9fa;
            padding: 15px;
            border-radius: 6px;
            margin-top: 10px;
        }
        .recommended-actions .value {
            color: #495057;
            line-height: 1.8;
        }
        .view-details {
            display: inline-block;
            margin-top: 20px;
            padding: 10px 20px;
            background-color: %s;
            color: white;
            text-decoration: none;
            border-radius: 5px;
            font-weight: 500;
        }
        .view-details:hover {
            opacity: 0.9;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h2>%s</h2>
            <p style="margin: 10px 0 0 0; font-size: 1.1em;">Alert Level: <span class="alert-level">%s</span></p>
        </div>

        <div class="section">
            <div class="section-title">Alert Details</div>
            <div class="info-row">
                <span class="label">Time:</span>
                <span class="value">%s</span>
            </div>
            %s
        </div>

        <div class="section">
            <div class="section-title">Alert Information</div>
            <div class="info-row">
                <span class="value">%s</span>
            </div>
        </div>

        <div class="section">
            <div class="section-title">Recommended Actions</div>
            <div class="recommended-actions">
                <div class="info-row">
                    <span class="value">• Please verify the alert condition<br>• Check the device status<br>• Review system logs if necessary</span>
                </div>
            </div>
        </div>

        <div style="text-align: center;">
            <a href="https://dangky.app/alarm/warning-message" class="view-details">View Full Alert Details</a>
        </div>

        <div class="footer">
            <p>This is an automated message from the Thingsly IoT Platform.</p>
            <p>Please do not reply to this email.</p>
        </div>
    </div>
</body>
</html>`, headerColor, headerColor, headerColor, headerColor, alarmConfig.Name, alarmLevelText, time.Now().Format("2006-01-02 15:04:05"),
			func() string {
				if alarmConfig.Description != nil && *alarmConfig.Description != "" {
					return fmt.Sprintf(`<div class="info-row">
                <span class="label">Description:</span>
                <span class="value">%s</span>
            </div>`, *alarmConfig.Description)
				}
				return ""
			}(), content)
		GroupApp.NotificationServicesConfig.ExecuteNotification(alarmConfig.NotificationGroupID, formattedContent, title)
	}
	device_ids_str, _ := json.Marshal(device_ids)
	id := uuid.New()
	t := time.Now().UTC()
	err = dal.AlarmHistorySave(&model.AlarmHistory{
		ID:                id,
		Name:              alarmConfig.Name,
		AlarmConfigID:     alarmConfigID,
		Content:           &content,
		Description:       alarmConfig.Description,
		TenantID:          alarmConfig.TenantID,
		SceneAutomationID: scene_automation_id,
		GroupID:           group_id,
		AlarmDeviceList:   string(device_ids_str),
		AlarmStatus:       alarmConfig.AlarmLevel,
		CreateAt:          t,
	})
	if err != nil {
		logrus.Error(err)
		return false, alarmName, ""
	}
	for _, deviceId := range device_ids {
		deviceInfo, _ := dal.GetDeviceByID(deviceId)
		go GroupApp.AlarmMessagePushSend(alarmConfig.Name, id, deviceInfo)
	}
	return true, alarmName, ""
}

func (*Alarm) GetAlarmInfoHistoryByID(id string) (map[string]interface{}, error) {
	alarmInfo, err := dal.GetAlarmInfoHistoryByID(id)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}
	return alarmInfo, nil
}

func (a *Alarm) GetAlarmDeviceCountsByTenant(tenantID string) (*model.AlarmDeviceCountsResponse, error) {
	ctx := context.Background()
	db := &dal.LatestDeviceAlarmQuery{}

	totalCount, err := db.CountDevicesByTenantAndStatus(ctx, tenantID)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"operation": "count_alarm_devices",
			"error":     err.Error(),
		})
	}

	return &model.AlarmDeviceCountsResponse{
		AlarmDeviceTotal: int64(totalCount),
	}, nil
}
