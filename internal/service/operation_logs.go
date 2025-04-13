package service

import (
	dal "github.com/HustIoTPlatform/backend/internal/dal"
	model "github.com/HustIoTPlatform/backend/internal/model"
	"github.com/HustIoTPlatform/backend/pkg/errcode"
	utils "github.com/HustIoTPlatform/backend/pkg/utils"

	"github.com/sirupsen/logrus"
)

type OperationLogs struct{}

func (*OperationLogs) CreateOperationLogs(operationLog *model.OperationLog) error {
	err := dal.CreateOperationLogs(operationLog)

	if err != nil {
		logrus.Error(err)
		return errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}

	return err
}


func (*OperationLogs) GetListByPage(Params *model.GetOperationLogListByPageReq, userClaims *utils.UserClaims) (map[string]interface{}, error) {

	total, list, err := dal.GetListByPage(Params, userClaims)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}
	OperationLogsListRsp := make(map[string]interface{})
	OperationLogsListRsp["total"] = total
	OperationLogsListRsp["list"] = list

	return OperationLogsListRsp, err
}
