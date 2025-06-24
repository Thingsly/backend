package dal

import (
	"context"
	"errors"
	"fmt"
	"time"

	model "github.com/Thingsly/backend/internal/model"
	"github.com/Thingsly/backend/internal/query"
	utils "github.com/Thingsly/backend/pkg/utils"

	"github.com/go-basic/uuid"
	"github.com/sirupsen/logrus"
)

func CreateSceneInfo(req model.CreateSceneReq, claims *utils.UserClaims) (string, error) {

	tx, err := StartTransaction()
	if err != nil {
		return "", err
	}

	sceneInfo := model.SceneInfo{}

	t := time.Now().UTC()
	sceneInfo.ID = uuid.New()

	sceneInfo.Name = req.Name
	sceneInfo.Description = &req.Description
	sceneInfo.TenantID = claims.TenantID
	sceneInfo.Creator = claims.ID
	sceneInfo.Updator = &claims.ID
	sceneInfo.CreatedAt = t
	sceneInfo.UpdatedAt = &t

	err = tx.SceneInfo.Create(&sceneInfo)
	if err != nil {
		return "", err
	}

	for _, v := range req.Actions {
		sceneAction := model.SceneActionInfo{}
		sceneAction.ID = uuid.New()
		sceneAction.SceneID = sceneInfo.ID
		sceneAction.ActionTarget = v.ActionTarget
		sceneAction.ActionType = v.ActionType
		sceneAction.ActionParamType = v.ActionParamType
		sceneAction.ActionParam = v.ActionParam
		sceneAction.ActionValue = v.ActionValue
		sceneAction.CreatedAt = t
		sceneAction.UpdatedAt = &t
		sceneAction.TenantID = claims.TenantID
		sceneAction.Remark = v.Remark
		err = tx.SceneActionInfo.Create(&sceneAction)
		if err != nil {
			Rollback(tx)
			return "", err
		}
	}

	err = Commit(tx)
	if err != nil {
		return "", err
	}

	return sceneInfo.ID, nil

}

func UpdateSceneInfo(req model.UpdateSceneReq, claims *utils.UserClaims) (string, error) {

	tx, err := StartTransaction()
	if err != nil {
		return "", err
	}

	sceneInfo := model.SceneInfo{}

	t := time.Now().UTC()
	//sceneInfo.ID = req.ID
	sceneInfo.Name = req.Name
	sceneInfo.Description = &req.Description
	sceneInfo.Updator = &claims.ID
	sceneInfo.UpdatedAt = &t
	//err = tx.SceneInfo.Save(&sceneInfo)
	result, err := tx.SceneInfo.Where(tx.SceneInfo.ID.Eq(req.ID)).Updates(sceneInfo)
	if err != nil {
		Rollback(tx)
		return "", err
	}
	if result.RowsAffected == 0 {
		Rollback(tx)
		return "", errors.New("edit failed")
	}

	_, err = tx.SceneActionInfo.Where(query.SceneActionInfo.SceneID.Eq(req.ID)).Delete()
	if err != nil {
		Rollback(tx)
		return "", err
	}

	for _, v := range req.Actions {
		sceneAction := model.SceneActionInfo{}
		sceneAction.ID = uuid.New()
		sceneAction.SceneID = req.ID
		sceneAction.ActionTarget = v.ActionTarget
		sceneAction.ActionType = v.ActionType
		sceneAction.ActionParamType = v.ActionParamType
		sceneAction.ActionParam = v.ActionParam
		sceneAction.ActionValue = v.ActionValue
		sceneAction.CreatedAt = t
		sceneAction.UpdatedAt = &t
		sceneAction.TenantID = claims.TenantID
		sceneAction.Remark = v.Remark
		err = tx.SceneActionInfo.Create(&sceneAction)
		if err != nil {
			Rollback(tx)
			return "", err
		}
	}

	err = Commit(tx)
	if err != nil {
		return "", err
	}

	return sceneInfo.ID, nil

}

func DeleteSceneInfo(scene_id string) error {
	_, err := query.SceneInfo.Where(query.SceneInfo.ID.Eq(scene_id)).Delete()
	if err != nil {
		logrus.Error(err)
	}
	return err
}

func GetSceneInfo(scene_id string) (*model.SceneInfo, error) {
	sceneInfo, err := query.SceneInfo.Where(query.SceneInfo.ID.Eq(scene_id)).First()
	if err != nil {
		logrus.Error(err)
	}
	return sceneInfo, err
}

func GetSceneInfoByPage(req *model.GetSceneListByPageReq, tenant_id string) (int64, []*model.SceneInfo, error) {
	q := query.SceneInfo
	var count int64
	queryBuilder := q.WithContext(context.Background())
	if req.Name != nil && *req.Name != "" {
		queryBuilder = queryBuilder.Where(q.Name.Like(fmt.Sprintf("%%%s%%", *req.Name)))
	}

	queryBuilder = queryBuilder.Where(q.TenantID.Eq(tenant_id))

	count, err := queryBuilder.Count()
	if err != nil {
		logrus.Error(err)
		return count, nil, err
	}

	if req.Page != 0 && req.PageSize != 0 {
		queryBuilder = queryBuilder.Limit(req.PageSize)
		queryBuilder = queryBuilder.Offset((req.Page - 1) * req.PageSize)
	}

	queryBuilder = queryBuilder.Order(q.CreatedAt.Desc())

	sceneList, err := queryBuilder.Find()
	if err != nil {
		return count, sceneList, err
	}
	return count, sceneList, nil
}

func GetSceneActionsInfo(scene_id string) ([]*model.SceneActionInfo, error) {
	sceneActionInfo, err := query.SceneActionInfo.Where(query.SceneActionInfo.SceneID.Eq(scene_id)).Find()
	if err != nil {
		logrus.Error(err)
	}
	return sceneActionInfo, err
}
