package service

import (
	"time"

	dal "github.com/Thingsly/backend/internal/dal"
	model "github.com/Thingsly/backend/internal/model"
	"github.com/Thingsly/backend/pkg/errcode"
	utils "github.com/Thingsly/backend/pkg/utils"

	"github.com/go-basic/uuid"
	"github.com/sirupsen/logrus"
)

type NotificationGroup struct{}

//	type CreateNotificationGroupReq struct {
//		Name               string    `json:"name" validate:"required"`
//		NotificationType   string    `json:"notification_type" validate:"required"`
//		Status             int       `json:"status" validate:"required"`
//		NotificationConfig *string    `json:"notification_config" validate:"omitempty"`
//		Description        string    `json:"description" validate:"required"`
//		TenantID           string    `json:"tenant_id" validate:"required"`
//		CreateTime         time.Time `json:"create_time" validate:"required"`
//		UpdateTime         time.Time `json:"update_time" validate:"required"`
//		Remark             string    `json:"remark" validate:"required"`
//	}
func (*NotificationGroup) CreateNotificationGroup(createNotificationgroupReq *model.CreateNotificationGroupReq, u *utils.UserClaims) (*model.NotificationGroup, error) {
	var notificationGroup model.NotificationGroup
	notificationGroup.ID = uuid.New()
	notificationGroup.Name = createNotificationgroupReq.Name
	notificationGroup.NotificationConfig = createNotificationgroupReq.NotificationConfig
	notificationGroup.NotificationType = createNotificationgroupReq.NotificationType
	notificationGroup.Status = createNotificationgroupReq.Status
	notificationGroup.Description = createNotificationgroupReq.Description
	notificationGroup.Remark = createNotificationgroupReq.Remark
	notificationGroup.UpdatedAt = time.Now().UTC()
	notificationGroup.CreatedAt = time.Now().UTC()
	notificationGroup.TenantID = u.TenantID
	err := dal.CreateNotificationGroup(&notificationGroup)

	if err != nil {
		logrus.Error(err)
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}

	return &notificationGroup, nil
}

func (*NotificationGroup) GetNotificationGroupById(id string) (notificationGroup *model.NotificationGroup, err error) {
	notificationGroup, err = dal.GetNotificationGroupById(id)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}
	return
}

func (*NotificationGroup) UpdateNotificationGroup(id string, updateNotificationgroupReq *model.UpdateNotificationGroupReq) (*model.NotificationGroup, error) {
	notificationGroup, err := dal.GetNotificationGroupById(id)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}
	utils.SerializeData(updateNotificationgroupReq, notificationGroup)

	notificationGroup.UpdatedAt = time.Now().UTC()
	err = dal.UpdateNotificationGroup(notificationGroup)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}
	return notificationGroup, nil
}

func (*NotificationGroup) DeleteNotificationGroup(id string) error {
	err := dal.DeleteNotificationGroup(id)
	if err != nil {
		return errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}
	return nil
}

func (*NotificationGroup) GetNotificationGroupListByPage(pageParam *model.GetNotificationGroupListByPageReq, u *utils.UserClaims) (map[string]interface{}, error) {
	total, list, err := dal.GetNotificationGroupListByPage(pageParam, u)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}
	notificationListRsp := make(map[string]interface{})
	notificationListRsp["total"] = total
	notificationListRsp["list"] = list

	return notificationListRsp, err
}

func (*NotificationGroup) GetNotificationGroupListByTenantId(tenantid string) (map[string]interface{}, error) {
	total, list, err := dal.GetNotificationGroupByTenantId(tenantid)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}
	notificationGroupListRsp := make(map[string]interface{})
	notificationGroupListRsp["total"] = total
	notificationGroupListRsp["list"] = list

	return notificationGroupListRsp, err
}

func (*NotificationGroup) GetNotificationByTenantId(tenantid string) (map[string]interface{}, error) {
	total, list, err := dal.GetBoardListByTenantId(tenantid)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}
	boardListRsp := make(map[string]interface{})
	boardListRsp["total"] = total
	boardListRsp["list"] = list

	return boardListRsp, err
}
