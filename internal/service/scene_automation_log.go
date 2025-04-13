package service

import (
	"github.com/HustIoTPlatform/backend/internal/dal"
	model "github.com/HustIoTPlatform/backend/internal/model"
	utils "github.com/HustIoTPlatform/backend/pkg/utils"
)

type SceneAutomationLog struct{}

func (*SceneAutomationLog) GetSceneAutomationLog(req *model.GetSceneAutomationLogReq, u *utils.UserClaims) (interface{}, error) {
	total, data, err := dal.GetSceneAutomationLog(req, u.TenantID)
	logList := make(map[string]interface{})
	logList["total"] = total
	logList["list"] = data

	return logList, err
}
