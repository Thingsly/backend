package service

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Thingsly/backend/initialize"
	config "github.com/Thingsly/backend/mqtt"
	"github.com/Thingsly/backend/mqtt/publish"
	simulationpublish "github.com/Thingsly/backend/mqtt/simulation_publish"
	"github.com/Thingsly/backend/pkg/constant"
	"github.com/Thingsly/backend/pkg/errcode"
	"github.com/Thingsly/backend/pkg/utils"

	"github.com/go-basic/uuid"
	"github.com/mintance/go-uniqid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/xuri/excelize/v2"

	dal "github.com/Thingsly/backend/internal/dal"
	model "github.com/Thingsly/backend/internal/model"
)

type TelemetryData struct{}

func (*TelemetryData) GetCurrentTelemetrData(device_id string) (interface{}, error) {
	// d, err := dal.GetCurrentTelemetrData(device_id)
	d, err := dal.GetCurrentTelemetryDataEvolution(device_id)
	if err != nil {
		return nil, err
	}

	deviceInfo, err := dal.GetDeviceByID(device_id)
	if err != nil {
		return nil, err
	}
	telemetryModelMap := make(map[string]*model.DeviceModelTelemetry)
	telemetryModelUintMap := make(map[string]interface{})
	telemetryModelRWMap := make(map[string]interface{})

	if deviceInfo.DeviceConfigID != nil {

		deviceConfig, err := dal.GetDeviceConfigByID(*deviceInfo.DeviceConfigID)
		if err != nil {
			return nil, err
		}

		if deviceConfig.DeviceTemplateID != nil {

			telemetryModel, err := dal.GetDeviceModelTelemetryDataList(*deviceConfig.DeviceTemplateID)
			if err != nil {
				return nil, err
			}
			if len(telemetryModel) > 0 {

				for _, v := range telemetryModel {
					telemetryModelMap[v.DataIdentifier] = v
					telemetryModelUintMap[v.DataIdentifier] = v.Unit
					telemetryModelRWMap[v.DataIdentifier] = v.ReadWriteFlag
				}
			}
		}
	}

	data := make([]map[string]interface{}, 0)
	if len(d) > 0 {
		for _, v := range d {
			tmp := make(map[string]interface{})
			tmp["device_id"] = v.DeviceID
			tmp["key"] = v.Key
			tmp["ts"] = v.T
			tmp["tenant_id"] = v.TenantID
			if v.BoolV != nil {
				tmp["value"] = v.BoolV
			}
			if v.NumberV != nil {
				tmp["value"] = v.NumberV
			}
			if v.StringV != nil {
				tmp["value"] = v.StringV
			}

			if len(telemetryModelMap) > 0 {
				telemetryModel, ok := telemetryModelMap[v.Key]
				if ok {
					tmp["label"] = telemetryModel.DataName
					tmp["unit"] = telemetryModelUintMap[v.Key]
					tmp["read_write_flag"] = telemetryModelRWMap[v.Key]
					tmp["data_type"] = telemetryModel.DataType
					if telemetryModel.DataType != nil && *telemetryModel.DataType == "Enum" {
						var enumItems []model.EnumItem
						json.Unmarshal([]byte(*telemetryModel.AdditionalInfo), &enumItems)
						tmp["enum"] = enumItems
					}
				}
			}
			data = append(data, tmp)
		}
	}

	return data, err
}

func (*TelemetryData) GetCurrentTelemetrDataKeys(req *model.GetTelemetryCurrentDataKeysReq) (interface{}, error) {
	// d, err := dal.GetCurrentTelemetrData(device_id)

	d, err := dal.GetCurrentTelemetryDataEvolutionByKeys(req.DeviceID, req.Keys)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}

	deviceInfo, err := dal.GetDeviceByID(req.DeviceID)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}
	telemetryModelMap := make(map[string]*model.DeviceModelTelemetry)
	telemetryModelUintMap := make(map[string]interface{})

	if deviceInfo.DeviceConfigID != nil {

		deviceConfig, err := dal.GetDeviceConfigByID(*deviceInfo.DeviceConfigID)
		if err != nil {
			return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
				"sql_error": err.Error(),
			})
		}

		if deviceConfig.DeviceTemplateID != nil {

			telemetryModel, err := dal.GetDeviceModelTelemetryDataList(*deviceConfig.DeviceTemplateID)
			if err != nil {
				return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
					"sql_error": err.Error(),
				})
			}
			if len(telemetryModel) > 0 {

				for _, v := range telemetryModel {
					telemetryModelMap[v.DataIdentifier] = v
					telemetryModelUintMap[v.DataIdentifier] = v.Unit
				}
			}
		}
	}

	data := make([]map[string]interface{}, 0)
	if len(d) > 0 {
		for _, v := range d {
			tmp := make(map[string]interface{})

			tmp["device_id"] = v.DeviceID
			tmp["key"] = v.Key
			tmp["ts"] = v.T
			tmp["tenant_id"] = v.TenantID
			if v.BoolV != nil {
				tmp["value"] = v.BoolV
			}
			if v.NumberV != nil {
				tmp["value"] = v.NumberV
			}
			if v.StringV != nil {
				tmp["value"] = v.StringV
			}

			if len(telemetryModelMap) > 0 {
				telemetryModel, ok := telemetryModelMap[v.Key]
				if ok {
					tmp["label"] = telemetryModel.DataName
					tmp["unit"] = telemetryModelUintMap[v.Key]
					tmp["data_type"] = telemetryModel.DataType
					if telemetryModel.DataType != nil && *telemetryModel.DataType == "Enum" {
						var enumItems []model.EnumItem
						json.Unmarshal([]byte(*telemetryModel.AdditionalInfo), &enumItems)
						tmp["enum"] = enumItems
					}
				}
			}
			data = append(data, tmp)
		}
	}

	return data, err
}

func (*TelemetryData) GetCurrentTelemetrDataForWs(device_id string) (interface{}, error) {
	// d, err := dal.GetCurrentTelemetrData(device_id)

	d, err := dal.GetCurrentTelemetryDataEvolution(device_id)
	if err != nil {
		return nil, err
	}

	data := make(map[string]interface{})
	if len(d) > 0 {
		for _, v := range d {
			if v.BoolV != nil {
				data[v.Key] = v.BoolV
			}
			if v.NumberV != nil {
				data[v.Key] = v.NumberV
			}
			if v.StringV != nil {
				data[v.Key] = v.StringV
			}
		}
	}
	return data, err
}

func (*TelemetryData) GetCurrentTelemetrDataKeysForWs(device_id string, keys []string) (interface{}, error) {
	// d, err := dal.GetCurrentTelemetrData(device_id)

	d, err := dal.GetCurrentTelemetryDataEvolutionByKeys(device_id, keys)
	if err != nil {
		return nil, err
	}

	data := make(map[string]interface{})
	if len(d) > 0 {
		for _, v := range d {
			if v.BoolV != nil {
				data[v.Key] = v.BoolV
			}
			if v.NumberV != nil {
				data[v.Key] = v.NumberV
			}
			if v.StringV != nil {
				data[v.Key] = v.StringV
			}
		}
	}
	return data, err
}

func (*TelemetryData) GetTelemetrHistoryData(req *model.GetTelemetryHistoryDataReq) (interface{}, error) {

	sT := req.StartTime * 1000
	eT := req.EndTime * 1000

	d, err := dal.GetHistoryTelemetrData(req.DeviceID, req.Key, sT, eT)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}

	data := make([]map[string]interface{}, 0)
	if len(d) > 0 {
		for _, v := range d {
			tmp := make(map[string]interface{})

			tmp["device_id"] = v.DeviceID
			tmp["key"] = v.Key
			tmp["ts"] = v.T
			tmp["tenant_id"] = v.TenantID
			if v.BoolV != nil {
				tmp["value"] = v.BoolV
			}
			if v.NumberV != nil {
				tmp["value"] = v.NumberV
			}
			if v.StringV != nil {
				tmp["value"] = v.StringV
			}
			data = append(data, tmp)
		}
	}

	return data, nil
}

func (*TelemetryData) DeleteTelemetrData(req *model.DeleteTelemetryDataReq) error {
	err := dal.DeleteTelemetrData(req.DeviceID, req.Key)
	if err != nil {
		return err
	}

	err = dal.DeleteCurrentTelemetryData(req.DeviceID, req.Key)
	return err
}

func (*TelemetryData) GetCurrentTelemetrDetailData(device_id string) (interface{}, error) {
	data, err := dal.GetCurrentTelemetrDetailData(device_id)
	if err != nil {
		return nil, err
	}

	dataMap := make(map[string]interface{})

	dataMap["device_id"] = data.DeviceID
	dataMap["key"] = data.Key
	dataMap["ts"] = data.T
	dataMap["tenant_id"] = data.TenantID

	if data.BoolV != nil {
		dataMap["value"] = data.BoolV
	}

	if data.NumberV != nil {
		dataMap["value"] = data.NumberV
	}

	if data.StringV != nil {
		dataMap["value"] = data.StringV
	}

	return dataMap, err
}

func (*TelemetryData) GetTelemetrHistoryDataByPage(req *model.GetTelemetryHistoryDataByPageReq) (interface{}, error) {
	if *req.ExportExcel {
		var addr string
		f := excelize.NewFile()
		f.SetCellValue("Sheet1", "A1", "Time")
		f.SetCellValue("Sheet1", "B1", "Value")

		batchSize := 100000
		offset := 0
		rowNumber := 2

		for {
			datas, err := dal.GetHistoryTelemetrDataByExport(req, offset, batchSize)
			if err != nil {
				return addr, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
					"sql_error": err.Error(),
				})
			}
			if len(datas) == 0 {
				break
			}
			for _, data := range datas {
				t := time.Unix(0, data.T*int64(time.Millisecond))
				f.SetCellValue("Sheet1", fmt.Sprintf("A%d", rowNumber), t.Format("2006-01-02 15:04:05"))

				cellRef := fmt.Sprintf("B%d", rowNumber)

				if data.StringV != nil && *data.StringV != "" {

					f.SetCellValue("Sheet1", cellRef, *data.StringV)
				} else if data.NumberV != nil {

					f.SetCellValue("Sheet1", cellRef, *data.NumberV)
				} else if data.BoolV != nil {

					f.SetCellValue("Sheet1", cellRef, *data.BoolV)
				} else {

					f.SetCellValue("Sheet1", cellRef, "")
				}
				rowNumber++
			}
			offset += batchSize
		}

		uploadDir := "./files/excel/"
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			return "", errcode.WithVars(errcode.CodeFilePathGenError, map[string]interface{}{
				"error": err.Error(),
			})
		}

		uniqidStr := uniqid.New(uniqid.Params{
			Prefix:      "excel",
			MoreEntropy: true,
		})

		addr = "files/excel/DataList" + uniqidStr + ".xlsx"

		if err := f.SaveAs(addr); err != nil {
			return "", errcode.WithVars(errcode.CodeFileSaveError, map[string]interface{}{
				"error": err.Error(),
			})
		}

		return addr, nil
	}

	_, data, err := dal.GetHistoryTelemetrDataByPage(req)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}

	var easyData []map[string]interface{}
	for _, v := range data {
		d := make(map[string]interface{})
		d["ts"] = v.T
		d["key"] = v.Key
		if v.StringV != nil {
			d["value"] = v.StringV
		} else if v.NumberV != nil {
			d["value"] = v.NumberV
		} else if v.BoolV != nil {
			d["value"] = v.BoolV
		} else {
			d["value"] = ""
		}

		easyData = append(easyData, d)
	}
	return easyData, nil
}

func (*TelemetryData) GetTelemetrHistoryDataByPageV2(req *model.GetTelemetryHistoryDataByPageReq) (interface{}, error) {
	if req.ExportExcel != nil && *req.ExportExcel {
		var addr string
		f := excelize.NewFile()
		f.SetCellValue("Sheet1", "A1", "Time")
		f.SetCellValue("Sheet1", "B1", "Value")

		batchSize := 100000
		offset := 0
		rowNumber := 2

		for {
			datas, err := dal.GetHistoryTelemetrDataByExport(req, offset, batchSize)
			if err != nil {
				return addr, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
					"sql_error": err.Error(),
				})
			}
			if len(datas) == 0 {
				break
			}
			for _, data := range datas {
				t := time.Unix(0, data.T*int64(time.Millisecond))
				f.SetCellValue("Sheet1", fmt.Sprintf("A%d", rowNumber), t.Format("2006-01-02 15:04:05"))
				f.SetCellValue("Sheet1", fmt.Sprintf("B%d", rowNumber), *data.NumberV)
				rowNumber++
			}
			offset += batchSize
		}

		uploadDir := "./files/excel/"
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			return nil, errcode.WithVars(errcode.CodeFilePathGenError, map[string]interface{}{
				"error": err.Error(),
			})
		}

		uniqidStr := uniqid.New(uniqid.Params{
			Prefix:      "excel",
			MoreEntropy: true,
		})

		fileName := "DataList" + uniqidStr + ".xlsx"
		filePath := "files/excel/" + fileName

		if err := f.SaveAs(filePath); err != nil {
			return nil, errcode.WithVars(errcode.CodeFileSaveError, map[string]interface{}{
				"error": err.Error(),
			})
		}

		result := map[string]interface{}{
			"filePath":   filePath,
			"fileName":   fileName,
			"fileType":   "excel",
			"createTime": time.Now().Format("2006-01-02T15:04:05-0700"),
		}

		return result, nil
	}

	total, data, err := dal.GetHistoryTelemetrDataByPage(req)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}

	var easyData []map[string]interface{}
	for _, v := range data {
		d := make(map[string]interface{})
		d["ts"] = v.T
		d["key"] = v.Key
		if v.StringV != nil {
			d["value"] = v.StringV
		} else if v.NumberV != nil {
			d["value"] = v.NumberV
		} else if v.BoolV != nil {
			d["value"] = v.BoolV
		} else {
			d["value"] = ""
		}

		easyData = append(easyData, d)
	}
	dataRsp := make(map[string]interface{})
	dataRsp["total"] = total
	dataRsp["list"] = easyData
	return dataRsp, nil
}

func (*TelemetryData) ServeEchoData(req *model.ServeEchoDataReq) (interface{}, error) {

	deviceInfo, err := dal.GetDeviceByID(req.DeviceId)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}
	voucher := deviceInfo.Voucher

	if !IsJSON(voucher) {
		return nil, errcode.NewWithMessage(errcode.CodeParamError, "voucher is not json")
	}
	var voucherMap map[string]interface{}
	err = json.Unmarshal([]byte(voucher), &voucherMap)
	if err != nil {
		return nil, err
	}

	var username, password, host, post, payload, clientID string
	if _, ok := voucherMap["username"]; !ok {
		return nil, errcode.NewWithMessage(errcode.CodeParamError, "username is not exist")
	}
	username = voucherMap["username"].(string)

	if _, ok := voucherMap["password"]; !ok {
		password = ""
	} else {
		password = voucherMap["password"].(string)
	}

	accessAddress := viper.GetString("mqtt.access_address")
	if accessAddress == "" {
		return nil, errcode.NewWithMessage(errcode.CodeParamError, "mqtt access address is not exist")
	}
	accessAddressList := strings.Split(accessAddress, ":")
	host = accessAddressList[0]
	post = accessAddressList[1]
	topic := config.MqttConfig.Telemetry.SubscribeTopic
	clientID = "mqtt_" + uuid.New()[0:12]
	payload = `{\"test_data1\":25.5,\"test_data2\":60}`

	command := utils.BuildMosquittoPubCommand(host, post, username, password, topic, payload, clientID)
	return command, nil
}

func (*TelemetryData) TelemetryPub(mosquittoCommand string) (interface{}, error) {

	params, err := utils.ParseMosquittoPubCommand(mosquittoCommand)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeParamError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	var voucher string
	if params.Password == "" {
		voucher = fmt.Sprintf("{\"username\":\"%s\"}", params.Username)
	} else {
		voucher = fmt.Sprintf("{\"username\":\"%s\",\"password\":\"%s\"}", params.Username, params.Password)
	}

	deviceInfo, err := dal.GetDeviceByVoucher(voucher)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}
	var isOnline int
	if deviceInfo.IsOnline == int16(1) {
		isOnline = 1
	}

	logrus.Debug("params:", params)
	err = simulationpublish.PublishMessage(params.Host, params.Port, params.Topic, params.Payload, params.Username, params.Password, params.ClientId)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeSystemError, map[string]interface{}{
			"error": err.Error(),
		})
	}
	go func() {
		time.Sleep(3 * time.Second)

		if isOnline == 1 {
			dal.UpdateDeviceOnlineStatus(deviceInfo.ID, int16(isOnline))

			err = publish.PublishOnlineMessage(deviceInfo.ID, []byte("1"))
			if err != nil {
				logrus.Error("publish online message failed:", err)
			}
		}
	}()
	return nil, nil
}

func (*TelemetryData) GetTelemetrSetLogsDataListByPage(req *model.GetTelemetrySetLogsListByPageReq) (interface{}, error) {
	count, data, err := dal.GetTelemetrySetLogsListByPage(req)
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

/*
 1. Parameter Descriptions:
    - aggregate_window [Aggregation Interval]
        - no_aggregate: No aggregation
        - Options: "30s", "1m", "2m", "5m", "10m", "30m", "1h", "3h", "6h", "1d", "7d", "1mo"

    - time_range [Time Range]
        - Supported values: custom, last_5m, last_15m, last_30m, last_1h, last_3h
        - For "custom", the backend checks the start and end time range. If the range exceeds 3 hours, "no_aggregate" is not allowed.

    - aggregate_function [Aggregation Method]
        - avg: Average
        - max: Maximum

 2. Frontend Interaction Rules:
    - Page Initialization: Defaults to "Last 1 Hour" - "No Aggregation" - Aggregation method is hidden initially; shown when interval is selected (avg/max)
    - Last 5 Minutes – Show all aggregation intervals
    - Last 15 Minutes – Show all aggregation intervals
    - Last 30 Minutes – Show all aggregation intervals
    - Last 1 Hour – Show all aggregation intervals
    - Last 3 Hours – Default interval: "30s" (No aggregation is not selectable) – Default method: "avg"
    - Last 6 Hours – Default: "1m" (No aggregation and intervals ≤30s not selectable) – Default: "avg"
    - Last 12 Hours – Default: "2m" (No aggregation and intervals ≤1m not selectable) – Default: "avg"
    - Last 24 Hours – Default: "5m" (No aggregation and intervals ≤2m not selectable) – Default: "avg"
    - Last 3 Days – Default: "10m" (No aggregation and intervals ≤5m not selectable) – Default: "avg"
    - Last 7 Days – Default: "30m" (No aggregation and intervals ≤10m not selectable) – Default: "avg"
    - Last 15 Days – Default: "1h" (No aggregation and intervals ≤30m not selectable) – Default: "avg"
    - Last 30 Days – Default: "1h" (No aggregation and intervals ≤30m not selectable) – Default: "avg"
    - Last 60 Days – Default: "3h" (No aggregation and intervals ≤1h not selectable) – Default: "avg"
    - Last 90 Days – Default: "6h" (No aggregation and intervals ≤3h not selectable) – Default: "avg"
    - Last 6 Months – Default: "6h" (No aggregation and intervals ≤3h not selectable) – Default: "avg"
    - Last 1 Year – Default: "1mo" (No aggregation and intervals ≤7d not selectable) – Default: "avg"
    - Today – Default: "5m" (No aggregation and intervals ≤2m not selectable) – Default: "avg"
    - Yesterday – Default: "5m" (No aggregation and intervals ≤2m not selectable) – Default: "avg"
    - Day Before Yesterday – Default: "5m" (No aggregation and intervals ≤2m not selectable) – Default: "avg"
    - Same Day Last Week – Default: "5m" (No aggregation and intervals ≤2m not selectable) – Default: "avg"
    - This Week – Default: "30m" (No aggregation and intervals ≤10m not selectable) – Default: "avg"
    - Last Week – Default: "30m" (No aggregation and intervals ≤10m not selectable) – Default: "avg"
    - This Month – Default: "1h" (No aggregation and intervals ≤30m not selectable) – Default: "avg"
    - Last Month – Default: "1h" (No aggregation and intervals ≤30m not selectable) – Default: "avg"
    - This Year – Default: "1mo" (No aggregation and intervals ≤7d not selectable) – Default: "avg"
    - Last Year – Default: "1mo" (No aggregation and intervals ≤7d not selectable) – Default: "avg"

 Example Request Payload (frontend can directly use this format):

 // Custom Range, No Aggregation
 {
     "device_id": "4a5b326c-ba99-9ea2-34a9-1c484d69a1ab",
     "key": "temperature",
     "start_time": 1691048558615446,
     "end_time": 1691048693603021,
     "aggregate_window": "no_aggregate",
     "time_range": "custom"
 }

 // 30s Interval, Max Value
 {
     "device_id": "4a5b326c-ba99-9ea2-34a9-1c484d69a1ab",
     "key": "temperature",
     "start_time": 1691048558615446,
     "end_time": 1691048693603021,
     "aggregate_window": "30s",
     "aggregate_function": "max"
 }
*/

func (*TelemetryData) GetTelemetrServeStatisticData(req *model.GetTelemetryStatisticReq) (any, error) {

	if err := processTimeRange(req); err != nil {
		return nil, err
	}

	rspData, err := fetchTelemetryData(req)
	if err != nil {
		return nil, err
	}

	if !req.IsExport {
		if len(rspData) == 0 {
			return []map[string]interface{}{}, nil
		}
		return rspData, nil
	}

	data, err := exportToCSV(req, rspData)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}
	return data, nil
}

func processTimeRange(req *model.GetTelemetryStatisticReq) error {
	if req.AggregateWindow == "no_aggregate" {

		if req.EndTime-req.StartTime > 24*time.Hour.Milliseconds() {
			return errcode.New(207001)
		}
	}
	if req.TimeRange == "custom" {
		if req.StartTime == 0 || req.EndTime == 0 || req.StartTime > req.EndTime {
			return errcode.New(207002)
		}
		return nil
	}

	timeRanges := map[string]time.Duration{
		"last_5m":  5 * time.Minute,
		"last_15m": 15 * time.Minute,
		"last_30m": 30 * time.Minute,
		"last_1h":  time.Hour,
		"last_3h":  3 * time.Hour,
		"last_6h":  6 * time.Hour,
		"last_12h": 12 * time.Hour,
		"last_24h": 24 * time.Hour,
		"last_3d":  72 * time.Hour,
		"last_7d":  7 * 24 * time.Hour,
		"last_15d": 15 * 24 * time.Hour,
		"last_30d": 30 * 24 * time.Hour,
		"last_60d": 60 * 24 * time.Hour,
		"last_90d": 90 * 24 * time.Hour,
		"last_6m":  180 * 24 * time.Hour,
		"last_1y":  365 * 24 * time.Hour,
	}

	duration, ok := timeRanges[req.TimeRange]
	if !ok {
		return errcode.WithVars(207003, map[string]interface{}{
			"time_range": req.TimeRange,
		})
	}

	now := time.Now()
	req.EndTime = now.UnixNano() / 1e6
	req.StartTime = now.Add(-duration).UnixNano() / 1e6
	return nil
}

func fetchTelemetryData(req *model.GetTelemetryStatisticReq) ([]map[string]interface{}, error) {
	if req.AggregateWindow == "no_aggregate" {
		data, err := dal.GetTelemetrStatisticData(req.DeviceId, req.Key, req.StartTime, req.EndTime)
		if err != nil {
			return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
				"sql_error": err.Error(),
			})
		}
		return data, nil
	}

	if err := validateAggregateWindow(req.StartTime, req.EndTime, req.AggregateWindow); err != nil {
		return nil, err
	}

	if req.AggregateFunction == "" {
		req.AggregateFunction = "avg"
	}

	return dal.GetTelemetrStatisticaAgregationData(
		req.DeviceId,
		req.Key,
		req.StartTime,
		req.EndTime,
		dal.StatisticAggregateWindowMillisecond[req.AggregateWindow],
		req.AggregateFunction,
	)
}

func exportToCSV(req *model.GetTelemetryStatisticReq, data []map[string]interface{}) (map[string]interface{}, error) {

	if len(data) == 0 {
		return nil, errcode.New(202100)
	}

	exportDir := "./files/excel/telemetry/"
	if err := os.MkdirAll(exportDir, os.ModePerm); err != nil {
		return nil, errcode.WithVars(202101, map[string]interface{}{
			"error": err.Error(),
		})
	}

	fileName := fmt.Sprintf("%s_%s_%d_%d.csv", req.DeviceId, req.Key, req.StartTime, req.EndTime)
	filePath := filepath.Join(exportDir, fileName)

	file, err := os.Create(filePath)
	if err != nil {
		return nil, errcode.WithVars(202102, map[string]interface{}{
			"error": err.Error(),
		})
	}

	defer func() {
		syncErr := file.Sync()
		closeErr := file.Close()
		if err == nil {
			err = syncErr
		}
		if err == nil {
			err = closeErr
		}
	}()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write([]string{"Timestamp", "Value"}); err != nil {
		return nil, errcode.WithVars(202103, map[string]interface{}{
			"error": err.Error(),
		})
	}

	for _, row := range data {
		timestamp, ok := row["x"].(int64)
		if !ok {
			return nil, errcode.New(202105)
		}

		value, ok := row["y"].(float64)
		if !ok {
			return nil, errcode.New(202106)
		}

		t := time.Unix(0, timestamp*int64(time.Millisecond))
		formattedTime := t.Format("2006-01-02 15:04:05.000")

		if err := writer.Write([]string{formattedTime, fmt.Sprintf("%.3f", value)}); err != nil {
			return nil, errcode.WithVars(202104, map[string]interface{}{
				"error": err.Error(),
			})
		}
	}

	logrus.Info("CSV file created: ", filePath)

	return map[string]interface{}{
		"file_name": fileName,
		"file_path": filePath,
	}, nil
}

type AggregateRule struct {
	Days         int
	MinInterval  string
	FriendlyDesc string
}

func validateAggregateWindow(startTime, endTime int64, aggregateWindow string) error {

	days := int((endTime - startTime) / (24 * 60 * 60 * 1000))

	// Rules defining the minimum allowed aggregation interval for different time spans (in days).
	rules := []AggregateRule{
		{365, "7d", "1 year"},
		{180, "1d", "6 months"},
		{90, "6h", "90 days"},
		{60, "3h", "60 days"},
		{30, "1h", "30 days"},
		{15, "30m", "15 days"},
		{7, "10m", "7 days"},
		{3, "5m", "3 days"},
		{1, "2m", "1 day"},
	}

	for _, rule := range rules {
		// If the selected time range exceeds this rule's day threshold,
		// and the aggregation window is shorter than allowed, return an error.
		if days > rule.Days && !isValidInterval(aggregateWindow, rule.MinInterval) {
			return errcode.WithVars(207004, map[string]interface{}{
				"time_range":         rule.FriendlyDesc,
				"min_interval":       rule.MinInterval,
				"current_time_range": fmt.Sprintf("%s to %s (%d days)", formatTime(startTime), formatTime(endTime), days),
				"aggregate_window":   aggregateWindow,
			})
		}
	}

	return nil
}

func isValidInterval(current, minInterval string) bool {

	weights := map[string]int{
		"30s": 1,
		"1m":  2,
		"2m":  3,
		"5m":  4,
		"10m": 5,
		"30m": 6,
		"1h":  7,
		"3h":  8,
		"6h":  9,
		"1d":  10,
		"7d":  11,
		"1mo": 12,
	}

	currentWeight, exists := weights[current]
	if !exists {
		return false
	}

	minWeight, exists := weights[minInterval]
	if !exists {
		return false
	}

	return currentWeight >= minWeight
}

func formatTime(timestamp int64) string {
	return time.Unix(timestamp/1000, 0).Format("2006-01-02 15:04:05")
}

func (*TelemetryData) TelemetryPutMessage(ctx context.Context, userID string, param *model.PutMessage, operationType string) error {
	var (
		log = dal.TelemetrySetLogsQuery{}

		errorMessage string
	)

	if !json.Valid([]byte(param.Value)) {
		errorMessage = "value must be json"
	}

	deviceInfo, err := initialize.GetDeviceCacheById(param.DeviceID)
	if err != nil {
		logrus.Error(ctx, "[TelemetryPutMessage][GetDeviceCacheById]failed:", err)
		return errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	var protocolType string
	var deviceConfig *model.DeviceConfig
	var deviceType string

	if deviceInfo.DeviceConfigID != nil {
		deviceConfig, err = dal.GetDeviceConfigByID(*deviceInfo.DeviceConfigID)
		if err != nil {
			logrus.Error(ctx, "[TelemetryPutMessage][GetDeviceConfigByID]failed:", err)
			return errcode.WithData(errcode.CodeDBError, map[string]interface{}{
				"error": err.Error(),
			})
		}
		deviceType = deviceConfig.DeviceType

		if deviceConfig.ProtocolType != nil {
			protocolType = *deviceConfig.ProtocolType
		} else {
			return errcode.WithData(errcode.CodeParamError, map[string]interface{}{
				"error": "protocolType is nil",
			})
		}
	} else {
		protocolType = "MQTT"
		deviceType = "1"

	}

	var topic string
	if protocolType == "MQTT" {

		// messageID := common.GetMessageID()
		topic, err = getTopicByDevice(deviceInfo, deviceType, param)
		if err != nil {
			logrus.Error(ctx, "failed to get topic", err)
			return errcode.WithData(errcode.CodeParamError, map[string]interface{}{
				"error": err.Error(),
			})
		}

		if deviceType == "3" || deviceType == "2" {

			var inputData map[string]interface{}
			if err := json.Unmarshal([]byte(param.Value), &inputData); err != nil {
				return errcode.WithData(errcode.CodeParamError, map[string]interface{}{
					"error": err.Error(),
				})
			}

			var outputData map[string]interface{}
			if deviceType == "3" {
				if deviceInfo.SubDeviceAddr == nil {
					return errcode.WithData(errcode.CodeParamError, map[string]interface{}{
						"error": "subDeviceAddr is nil",
					})
				}
				outputData = map[string]interface{}{
					"sub_device_data": map[string]interface{}{
						*deviceInfo.SubDeviceAddr: inputData,
					},
				}
			} else {
				outputData = map[string]interface{}{
					"gateway_data": inputData,
				}
			}

			output, err := json.Marshal(outputData)
			if err != nil {
				return errcode.WithData(errcode.CodeParamError, map[string]interface{}{
					"error": err.Error(),
				})
			}
			param.Value = string(output)
		}
	} else {

		subTopicPrefix, err := dal.GetServicePluginSubTopicPrefixByDeviceConfigID(*deviceInfo.DeviceConfigID)
		if err != nil {
			logrus.Error(ctx, "failed to get sub topic prefix", err)
			return errcode.WithData(errcode.CodeParamError, map[string]interface{}{
				"error": err.Error(),
			})
		}
		topic = fmt.Sprintf("%s%s%s", subTopicPrefix, config.MqttConfig.Telemetry.PublishTopic, deviceInfo.ID)

	}

	if deviceInfo.DeviceConfigID != nil && *deviceInfo.DeviceConfigID != "" {
		script, err := initialize.GetScriptByDeviceAndScriptType(deviceInfo, "B")
		if err != nil {
			logrus.Error(err.Error())
			return err
		}
		if script != nil && script.Content != nil && *script.Content != "" {
			msg, err := utils.ScriptDeal(*script.Content, []byte(param.Value), topic)
			if err != nil {
				logrus.Error(err.Error())
				return err
			}
			param.Value = msg
		}
	}
	err = publish.PublishTelemetryMessage(topic, deviceInfo, param)
	if err != nil {
		logrus.Error(ctx, "Failed to dispatch", err)
		errorMessage = err.Error()
	}
	// operationType := strconv.Itoa(constant.Manual)

	description := "Send telemetry log records"
	logInfo := &model.TelemetrySetLog{
		ID:            uuid.New(),
		DeviceID:      param.DeviceID,
		OperationType: &operationType,
		Datum:         &(param.Value),
		Status:        nil,
		ErrorMessage:  &errorMessage,
		CreatedAt:     time.Now().UTC(),
		Description:   &description,
		UserID:        &userID,
	}

	if userID == "" {
		logInfo.UserID = nil
	}
	if err != nil {
		logInfo.ErrorMessage = &errorMessage
		status := strconv.Itoa(constant.StatusFailed)
		logInfo.Status = &status
	} else {
		status := strconv.Itoa(constant.StatusOK)
		logInfo.Status = &status
	}
	_, err = log.Create(ctx, logInfo)
	if err != nil {
		logrus.Error(ctx, "failed to create telemetry set log", err)
		return errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"error": err.Error(),
		})
	}
	return err
}

func getTopicByDevice(deviceInfo *model.Device, deviceType string, param *model.PutMessage) (string, error) {
	switch deviceType {
	case "1":

		return fmt.Sprintf("%s%s", config.MqttConfig.Telemetry.PublishTopic, deviceInfo.DeviceNumber), nil
	case "2", "3":
		// Retrieve the gateway ID for device types "2" or "3".
		gatewayID := deviceInfo.ID

		// For device type "3", if ParentID is nil, return an error.
		if deviceType == "3" {
			if deviceInfo.ParentID == nil {
				return "", fmt.Errorf("parentID is null")
			}
			gatewayID = *deviceInfo.ParentID
		}

		// Fetch the gateway information using the gateway ID.
		gatewayInfo, err := initialize.GetDeviceCacheById(gatewayID)
		if err != nil {
			return "", fmt.Errorf("failed to retrieve gateway information: %v", err)
		}

		// Return the formatted MQTT topic for publishing telemetry data using the gateway's device number.
		return fmt.Sprintf(config.MqttConfig.Telemetry.GatewayPublishTopic, gatewayInfo.DeviceNumber), nil

	default:
		// Return an error for unknown device types.
		return "", fmt.Errorf("unknown device type")
	}
}

func (*TelemetryData) ServeMsgCountByTenantId(tenantId string) (int64, error) {
	cnt, err := dal.GetTelemetryDataCountByTenantId(tenantId)
	if err != nil {
		return 0, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}
	return cnt, err
}
