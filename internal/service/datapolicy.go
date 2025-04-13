package service

import (
	"time"

	dal "github.com/HustIoTPlatform/backend/internal/dal"
	model "github.com/HustIoTPlatform/backend/internal/model"
	"github.com/HustIoTPlatform/backend/pkg/errcode"
	"github.com/HustIoTPlatform/backend/pkg/utils"

	"github.com/sirupsen/logrus"
)

type DataPolicy struct{}

func (*DataPolicy) UpdateDataPolicy(UpdateDataPolicyReq *model.UpdateDataPolicyReq) error {
	var datapolicy = model.DataPolicy{}
	datapolicy.ID = UpdateDataPolicyReq.Id
	datapolicy.RetentionDay = UpdateDataPolicyReq.RetentionDays
	datapolicy.Enabled = UpdateDataPolicyReq.Enabled
	datapolicy.Remark = UpdateDataPolicyReq.Remark
	err := dal.UpdateDataPolicy(&datapolicy)
	if err != nil {
		logrus.Error(err)
		return errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}
	return err
}

func (*DataPolicy) GetDataPolicyListByPage(Params *model.GetDataPolicyListByPageReq) (map[string]interface{}, error) {

	total, list, err := dal.GetDataPolicyListByPage(Params)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}
	datapolicyListRsp := make(map[string]interface{})
	datapolicyListRsp["total"] = total
	datapolicyListRsp["list"] = list

	return datapolicyListRsp, err
}

func (*DataPolicy) CleanSystemDataByCron() error {
	data, err := dal.GetDataPolicy()
	if err != nil {
		return err
	}

	now := time.Now()
	for _, v := range data {

		if v.Enabled == "2" {
			continue
		}

		if utils.IsToday(*v.LastCleanupTime) {
			continue
		}

		if v.DataType == "1" {
			daysAgeInt64 := utils.MillisecondsTimestampDaysAgo(int(v.RetentionDay))
			daysAgeTime := utils.DaysAgo(int(v.RetentionDay))
			err := dal.DeleteTelemetrDataByTime(daysAgeInt64)
			if err != nil {
				return err
			}

			var datapolicy = model.DataPolicy{}
			datapolicy.ID = v.ID
			datapolicy.LastCleanupTime = &now
			datapolicy.LastCleanupDataTime = &daysAgeTime
			err = dal.UpdateDataPolicy(&datapolicy)
			if err != nil {
				return err
			}

		} else if v.DataType == "2" {

			daysAge := utils.DaysAgo(int(v.RetentionDay))
			err := dal.DeleteOperationLogsByTime(daysAge)
			if err != nil {
				return err
			}

			var datapolicy = model.DataPolicy{}
			datapolicy.ID = v.ID
			datapolicy.LastCleanupTime = &now
			datapolicy.LastCleanupDataTime = &daysAge
			err = dal.UpdateDataPolicy(&datapolicy)
			if err != nil {
				return err
			}
		}

	}
	return nil
}
