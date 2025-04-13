package service

import (
	dal "github.com/HustIoTPlatform/backend/internal/dal"
	model "github.com/HustIoTPlatform/backend/internal/model"
	"github.com/HustIoTPlatform/backend/pkg/errcode"
	utils "github.com/HustIoTPlatform/backend/pkg/utils"
)

type EventData struct{}

func (*EventData) GetEventDatasListByPage(req *model.GetEventDatasListByPageReq, _ *utils.UserClaims) (interface{}, error) {
	count, data, err := dal.GetEventDatasListByPage(req)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}

	dataMap := make(map[string]interface{})
	dataMap["count"] = count
	dataMap["list"] = data

	return dataMap, nil
}
