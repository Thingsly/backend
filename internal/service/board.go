package service

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	dal "github.com/Thingsly/backend/internal/dal"
	model "github.com/Thingsly/backend/internal/model"
	query "github.com/Thingsly/backend/internal/query"
	common "github.com/Thingsly/backend/pkg/common"
	"github.com/Thingsly/backend/pkg/errcode"
	global "github.com/Thingsly/backend/pkg/global"
	utils "github.com/Thingsly/backend/pkg/utils"

	"github.com/go-basic/uuid"
	"github.com/sirupsen/logrus"
)

type Board struct{}

func (*Board) CreateBoard(ctx context.Context, CreateBoardReq *model.CreateBoardReq) (*model.Board, error) {
	var (
		board = model.Board{}
		db    = dal.BoardQuery{}
	)

	board.ID = uuid.New()
	board.Name = CreateBoardReq.Name
	if CreateBoardReq.Config != nil && !IsJSON(*CreateBoardReq.Config) {
		return nil, errcode.NewWithMessage(errcode.CodeParamError, "config is not a valid JSON")
	}
	board.Config = CreateBoardReq.Config
	board.MenuFlag = &CreateBoardReq.MenuFlag
	board.Description = CreateBoardReq.Description
	board.Remark = CreateBoardReq.Remark
	board.UpdatedAt = time.Now().UTC()
	board.CreatedAt = time.Now().UTC()
	board.TenantID = CreateBoardReq.TenantID
	board.HomeFlag = CreateBoardReq.HomeFlag

	if CreateBoardReq.HomeFlag == "Y" {
		err := db.UpdateHomeFlagN(ctx, CreateBoardReq.TenantID)
		if err != nil {
			logrus.Error(err)
			return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
				"sql_error": err.Error(),
			})
		}
	}
	boardInfo, err := db.Create(ctx, &board)
	if err != nil {
		logrus.Error(err)
		err = errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}

	return boardInfo, err
}

func (*Board) UpdateBoard(ctx context.Context, UpdateBoardReq *model.UpdateBoardReq) (*model.Board, error) {
	db := dal.BoardQuery{}
	board := model.Board{}
	board.ID = UpdateBoardReq.Id
	board.Name = UpdateBoardReq.Name

	if UpdateBoardReq.Config != nil && !IsJSON(*UpdateBoardReq.Config) {
		return nil, errcode.WithVars(100002, map[string]interface{}{
			"field": "config",
			"error": "config is not a valid JSON",
		})
	}
	board.Config = UpdateBoardReq.Config
	board.HomeFlag = UpdateBoardReq.HomeFlag
	board.MenuFlag = &UpdateBoardReq.MenuFlag
	board.Description = UpdateBoardReq.Description
	board.Remark = UpdateBoardReq.Remark
	board.UpdatedAt = time.Now().UTC()
	if UpdateBoardReq.Id != "" {
		if board.HomeFlag == "Y" {
			_, err := db.First(ctx, query.Board.TenantID.Eq(UpdateBoardReq.TenantID), query.Board.HomeFlag.Eq("Y"), query.Board.ID.Neq(UpdateBoardReq.Id))
			if err != nil {
				logrus.Error(err)
			} else {

				err := db.UpdateHomeFlagN(ctx, UpdateBoardReq.TenantID)
				if err != nil {
					logrus.Error(err)
					return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
						"sql_error": err.Error(),
					})
				}
			}
		}
		err := dal.UpdateBoard(&board)
		if err != nil {
			logrus.Error(err)
			return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
				"sql_error": err.Error(),
			})
		}

	} else {

		if board.Name == "" {
			return nil, fmt.Errorf("name is required")
		}
		if board.HomeFlag == "" {
			board.HomeFlag = "N"
		}
		board.ID = uuid.New()
		board.TenantID = UpdateBoardReq.TenantID

		if board.HomeFlag == "Y" {
			_, err := db.First(ctx, query.Board.TenantID.Eq(UpdateBoardReq.TenantID), query.Board.HomeFlag.Eq("Y"))
			if err != nil {
				logrus.Error(err)
			} else {
				return nil, errcode.New(203004)
			}
		}
		boardInfo, err := db.Create(ctx, &board)
		if err != nil {
			logrus.Error(err)
			err = errcode.WithData(errcode.CodeDBError, map[string]interface{}{
				"sql_error": err.Error(),
			})
		}
		return boardInfo, err
	}
	return &board, nil
}

func (*Board) DeleteBoard(id string) error {
	err := dal.DeleteBoard(id)
	if err != nil {
		return errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}
	return err
}

func (*Board) GetBoardListByPage(Params *model.GetBoardListByPageReq, U *utils.UserClaims) (map[string]interface{}, error) {
	total, list, err := dal.GetBoardListByPage(Params, U.TenantID)
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

func (*Board) GetBoard(id string, U *utils.UserClaims) (interface{}, error) {
	board, err := dal.GetBoard(id, U.TenantID)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}

	return board, err
}

func (*Board) GetBoardListByTenantId(tenantid string) (interface{}, error) {
	_, data, err := dal.GetBoardListByTenantId(tenantid)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}
	return data, err
}

// GetDeviceTotal
func (*Board) GetDeviceTotal(ctx context.Context, authority string, tenantID string) (int64, error) {
	var (
		total int64
		err   error
		db    = dal.DeviceQuery{}
	)
	if common.CheckUserIsAdmin(authority) {
		total, err = db.Count(ctx)
	} else {
		total, err = db.CountByTenantID(ctx, tenantID)
	}
	if err != nil {
		return 0, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}

	return total, err
}

// GetDevice
func (*Board) GetDevice(ctx context.Context, U *utils.UserClaims) (data *model.GetBoardDeviceRes, err error) {
	var (
		total, on int64
		device    = query.Device
		db        = dal.DeviceQuery{}
	)

	if !common.CheckUserIsAdmin(U.Authority) {
		total, err = db.CountByTenantID(ctx, U.TenantID)
	} else {
		total, err = db.Count(ctx)
	}
	if err != nil {
		logrus.Error(ctx, "[GetDevice]Device count failed:", err)
		err = errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
		return
	}
	if !common.CheckUserIsAdmin(U.Authority) {
		on, err = db.CountByWhere(ctx, device.ActivateFlag.Eq("active"), device.TenantID.Eq(U.TenantID))
	} else {
		on, err = db.CountByWhere(ctx, device.ActivateFlag.Eq("active"))
	}
	if err != nil {
		logrus.Error(ctx, "[GetDevice]Device count/on failed:", err)
		err = errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
		return
	}
	data = &model.GetBoardDeviceRes{
		DeviceTotal: total,
		DeviceOn:    on,
	}
	return
}

// GetDeviceByTenantID
func (*Board) GetDeviceByTenantID(ctx context.Context, tenantID string) (data *model.GetBoardDeviceRes, err error) {
	var (
		total, on int64
		device    = query.Device
		db        = dal.DeviceQuery{}
	)

	logrus.Debugf("Getting device counts for tenant %s", tenantID)

	// Count total active devices
	total, err = db.CountByWhere(ctx, device.TenantID.Eq(tenantID), device.ActivateFlag.Neq("inactive"))
	if err != nil {
		logrus.Errorf("Failed to count total devices for tenant %s: %v", tenantID, err)
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}

	// Count online devices
	on, err = db.CountByWhere(ctx, device.ActivateFlag.Eq("active"), device.TenantID.Eq(tenantID), device.IsOnline.Eq(1))
	if err != nil {
		logrus.Errorf("Failed to count online devices for tenant %s: %v", tenantID, err)
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}

	logrus.Debugf("Device counts for tenant %s - Total: %d, Online: %d", tenantID, total, on)

	data = &model.GetBoardDeviceRes{
		DeviceTotal: total,
		DeviceOn:    on,
	}
	return
}

// GetDeviceTrend
func (*Device) GetDeviceTrend(ctx context.Context, tenantID string) (*model.DeviceTrendRes, error) {

	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)

	dates := []string{
		yesterday.Format("2006-01-02"),
		now.Format("2006-01-02"),
	}

	var allPoints []model.DeviceTrendPoint

	for _, date := range dates {

		key := fmt.Sprintf("device_stats:%s:%s", tenantID, date)

		statsJsonList, err := global.REDIS.LRange(ctx, key, 0, -1).Result()
		if err != nil {
			logrus.Errorf("Failed to retrieve device statistics from Redis: %v", err)
			return nil, errcode.WithData(errcode.CodeCacheError, map[string]interface{}{
				"error": err.Error(),
				"key":   key,
			})
		}

		for _, statsJson := range statsJsonList {
			var statsData struct {
				DeviceTotal int64     `json:"device_total"`
				DeviceOn    int64     `json:"device_on"`
				Timestamp   time.Time `json:"timestamp"`
			}

			if err := json.Unmarshal([]byte(statsJson), &statsData); err != nil {
				logrus.Errorf("Failed to parse device statistics data: %v", err)
				continue
			}

			if statsData.Timestamp.Before(yesterday) {
				continue
			}

			point := model.DeviceTrendPoint{
				Timestamp:     statsData.Timestamp,
				DeviceTotal:   statsData.DeviceTotal,
				DeviceOnline:  statsData.DeviceOn,
				DeviceOffline: statsData.DeviceTotal - statsData.DeviceOn,
			}

			allPoints = append(allPoints, point)
		}
	}

	sort.Slice(allPoints, func(i, j int) bool {
		return allPoints[i].Timestamp.Before(allPoints[j].Timestamp)
	})

	if len(allPoints) == 0 {
		currentStats, err := GroupApp.Board.GetDeviceByTenantID(ctx, tenantID)
		if err != nil {
			return nil, err
		}

		point := model.DeviceTrendPoint{
			Timestamp:     now,
			DeviceTotal:   currentStats.DeviceTotal,
			DeviceOnline:  currentStats.DeviceOn,
			DeviceOffline: currentStats.DeviceTotal - currentStats.DeviceOn,
		}

		allPoints = append(allPoints, point)
	}

	return &model.DeviceTrendRes{
		Points: allPoints,
	}, nil
}
