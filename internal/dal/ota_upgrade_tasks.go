package dal

import (
	"context"
	"fmt"
	"time"

	model "github.com/HustIoTPlatform/backend/internal/model"
	query "github.com/HustIoTPlatform/backend/internal/query"
	global "github.com/HustIoTPlatform/backend/pkg/global"

	"github.com/go-basic/uuid"
	"github.com/sirupsen/logrus"
)

func CreateOTAUpgradeTaskWithDetail(req *model.CreateOTAUpgradeTaskReq) ([]*model.OtaUpgradeTaskDetail, error) {

	var task = model.OtaUpgradeTask{}
	var taskDetail = []*model.OtaUpgradeTaskDetail{}

	t := time.Now().UTC()
	taskId := uuid.New()

	task.ID = taskId
	task.Name = req.Name
	task.OtaUpgradePackageID = req.OTAUpgradePackageId
	task.Description = req.Description
	task.CreatedAt = t
	task.Remark = req.Remark

	for _, v := range req.DeviceIdList {
		detail := &model.OtaUpgradeTaskDetail{}
		detail.ID = uuid.New()
		detail.DeviceID = v
		detail.Status = 1
		detail.UpdatedAt = &t
		detail.OtaUpgradeTaskID = taskId
		taskDetail = append(taskDetail, detail)
	}

	tx := query.Use(global.DB).Begin()

	if tx.Error != nil {
		return nil, tx.Error
	}

	if err := tx.OtaUpgradeTask.Create(&task); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.OtaUpgradeTaskDetail.CreateInBatches(taskDetail, len(taskDetail)); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return taskDetail, nil

}

func DeleteOTAUpgradeTask(id string) error {
	_, err := query.OtaUpgradeTask.Where(query.OtaUpgradeTask.ID.Eq(id)).Delete()
	return err
}

func GetOtaUpgradeTaskListByPage(p *model.GetOTAUpgradeTaskListByPageReq) (int64, []map[string]interface{}, error) {

	whereClause := "WHERE t.ota_upgrade_package_id = ?"
	params := []interface{}{p.OTAUpgradePackageId}

	countSQL := `SELECT COUNT(*) FROM ota_upgrade_tasks t ` + whereClause

	var totalCount int64
	err := global.DB.Raw(countSQL, params...).Scan(&totalCount).Error
	if err != nil {
		return 0, nil, err
	}

	if totalCount == 0 || p.Page <= 0 || p.PageSize <= 0 {
		return 0, []map[string]interface{}{}, nil
	}

	dataSQL := `SELECT t.*, 
                       (SELECT COUNT(*) 
                        FROM ota_upgrade_task_details d 
                        WHERE d.ota_upgrade_task_id = t.id) AS device_count 
                FROM ota_upgrade_tasks t ` + whereClause +
		" ORDER BY t.created_at DESC LIMIT ? OFFSET ?"

	params = append(params, p.PageSize, (p.Page-1)*p.PageSize)

	var tasks []map[string]interface{}
	err = global.DB.Raw(dataSQL, params...).Scan(&tasks).Error
	if err != nil {
		return 0, nil, err
	}

	return totalCount, tasks, nil
}

func GetOtaUpgradeTaskDetailListByPage(p *model.GetOTAUpgradeTaskDetailReq) (int64, interface{}, interface{}, error) {

	var count int64
	type StatusCount struct {
		Status string `json:"status"`
		Count  int    `json:"count"`
	}
	detailDataMap := make([]map[string]interface{}, 0)

	statsResult := make([]StatusCount, 0)
	statsData := query.Device
	otaTaskDetail := query.OtaUpgradeTaskDetail

	queryBuilder := statsData.WithContext(context.Background())
	queryBuilder = queryBuilder.Join(otaTaskDetail, otaTaskDetail.DeviceID.EqCol(statsData.ID))
	if p.DeviceName != nil {
		queryBuilder = queryBuilder.Where(statsData.Name.Like(fmt.Sprintf("%%%s%%", *p.DeviceName)))
	}
	queryBuilder = queryBuilder.Where(otaTaskDetail.OtaUpgradeTaskID.Eq(p.OtaUpgradeTaskId))
	queryBuilder = queryBuilder.Select(otaTaskDetail.Status, otaTaskDetail.Status.Count()).Group(otaTaskDetail.Status)
	err := queryBuilder.Scan(&statsResult)
	if err != nil {
		logrus.Error(err)
		return count, nil, statsResult, err
	}

	detailData := query.Device
	detailDataBuilder := detailData.WithContext(context.Background())
	detailDataBuilder = detailDataBuilder.Join(otaTaskDetail, otaTaskDetail.DeviceID.EqCol(detailData.ID))

	if p.DeviceName != nil && *p.DeviceName != "" {
		detailDataBuilder = detailDataBuilder.Where(detailData.Name.Like(fmt.Sprintf("%%%s%%", *p.DeviceName)))
	}

	if p.TaskStatus != nil {
		detailDataBuilder = detailDataBuilder.Where(detailData.Name.Like(fmt.Sprintf("%%%s%%", *p.DeviceName)))
	}
	detailDataBuilder.Where(otaTaskDetail.OtaUpgradeTaskID.Eq(p.OtaUpgradeTaskId))

	if p.Page != 0 && p.PageSize != 0 {
		detailDataBuilder = detailDataBuilder.Limit(p.PageSize)
		detailDataBuilder = detailDataBuilder.Offset((p.Page - 1) * p.PageSize)
	}
	otaTask := query.OtaUpgradeTask
	otaPackage := query.OtaUpgradePackage

	detailDataBuilder = detailDataBuilder.Join(otaTask, otaTask.ID.EqCol(otaTaskDetail.OtaUpgradeTaskID))

	detailDataBuilder = detailDataBuilder.Join(otaPackage, otaPackage.ID.EqCol(otaTask.OtaUpgradePackageID))

	// Upgrade task details ID, device name, device serial number, device model, device version, upgrade package version, upgrade progress, last updated time, status, status details
	detailDataBuilder = detailDataBuilder.Select(
		otaTaskDetail.ID, otaTaskDetail.OtaUpgradeTaskID, detailData.DeviceNumber, detailData.Name, detailData.CurrentVersion,
		otaPackage.Version, otaTaskDetail.Step, otaTaskDetail.UpdatedAt,
		otaTaskDetail.Status, otaTaskDetail.StatusDescription,
	)

	err = detailDataBuilder.Scan(&detailDataMap)
	if err != nil {
		logrus.Error(err)
		return count, detailDataMap, statsResult, err
	}
	count, err = detailDataBuilder.Count()
	return count, detailDataMap, statsResult, err

}
