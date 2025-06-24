package dal

import (
	"errors"
	"time"

	"github.com/Thingsly/backend/internal/model"
	query "github.com/Thingsly/backend/internal/query"

	"gorm.io/gorm"
)

// Get current telemetry data from telemetry_current_datas to replace telemetry_datas
func GetCurrentTelemetryDataEvolution(deviceId string) ([]*model.TelemetryCurrentData, error) {
	// dbType := viper.GetString("grpc.tptodb_type")
	// if dbType == "TSDB" || dbType == "KINGBASE" || dbType == "POLARDB" {
	// 	var telemetry []*model.TelemetryCurrentData
	// 	request := &pb.GetDeviceAttributesCurrentsRequest{
	// 		DeviceId: deviceId,
	// 	}

	// 	r, err := tptodb.TptodbClient.GetDeviceAttributesCurrents(context.Background(), request)
	// 	if err != nil {
	// 		logrus.Printf("GetDeviceAttributesCurrents err:%+v", err)
	// 		return nil, err
	// 	}
	// 	logrus.Printf("data: %+v", r.Data)
	// 	err = json.Unmarshal([]byte(r.Data), &telemetry)
	// 	if err != nil {
	// 		logrus.Printf("Unmarshal err:%v", err)
	// 		return nil, err
	// 	}
	// 	return telemetry, nil
	// }

	data, err := query.TelemetryCurrentData.Where(query.TelemetryCurrentData.DeviceID.Eq(deviceId)).Order(query.TelemetryCurrentData.T.Desc()).Find()
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Retrieve current telemetry data from telemetry_current_datas to replace telemetry_datas
func GetCurrentTelemetryDataEvolutionByKeys(deviceId string, keys []string) ([]*model.TelemetryCurrentData, error) {
	// dbType := viper.GetString("grpc.tptodb_type")
	// if dbType == "TSDB" || dbType == "KINGBASE" || dbType == "POLARDB" {
	// 	data := make([]*model.TelemetryCurrentData, 0)
	// 	fields := make([]map[string]interface{}, 0)
	// 	request := &pb.GetDeviceAttributesCurrentsRequest{
	// 		DeviceId:  deviceId,
	// 		Attribute: keys,
	// 	}
	// 	r, err := tptodb.TptodbClient.GetDeviceAttributesCurrents(context.Background(), request)
	// 	if err != nil {
	// 		logrus.Printf("err: %+v", err)
	// 		return nil, err
	// 	}
	// 	err = json.Unmarshal([]byte(r.Data), &fields)
	// 	if err != nil {
	// 		logrus.Printf("err: %+v", err)
	// 		return nil, err
	// 	}
	// 	logrus.Printf("fields: %+v", fields)

	// 	err = json.Unmarshal([]byte(r.Data), &data)
	// 	if err != nil {
	// 		logrus.Printf("Unmarshal err:%v", err)
	// 		return nil, err
	// 	}

	// 	return data, nil
	// }

	data, err := query.TelemetryCurrentData.Where(query.TelemetryCurrentData.DeviceID.Eq(deviceId), query.TelemetryCurrentData.Key.In(keys...)).Order(query.TelemetryCurrentData.T.Desc()).Find()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func GetCurrentTelemetryDataOneKeys(deviceId string, keys string) (interface{}, error) {
	data, err := query.TelemetryCurrentData.Where(query.TelemetryCurrentData.DeviceID.Eq(deviceId), query.TelemetryCurrentData.Key.Eq(keys)).Order(query.TelemetryCurrentData.T.Desc()).First()
	var result interface{}
	if err != nil {
		return result, err
		//} else if err == gorm.ErrRecordNotFound {
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		return result, nil
	}
	if data.BoolV != nil {
		//result = fmt.Sprintf("%t", *data.BoolV)
		result = *data.BoolV
	}
	if data.NumberV != nil {
		//result = fmt.Sprintf("%f", *data.NumberV)
		result = *data.NumberV
	}
	if data.StringV != nil {
		result = *data.StringV
	}
	return result, nil
}

func DeleteCurrentTelemetryData(deviceId string, key string) error {
	_, err := query.TelemetryCurrentData.Where(query.TelemetryCurrentData.DeviceID.Eq(deviceId), query.TelemetryCurrentData.Key.Eq(key)).Delete()
	return err
}

func DeleteCurrentTelemetryDataByDeviceId(deviceId string, tx *query.QueryTx) error {
	_, err := tx.TelemetryCurrentData.Where(query.TelemetryCurrentData.DeviceID.Eq(deviceId)).Delete()
	return err
}

type NewDeviceData struct {
	DeviceID  string    `json:"device_id"`
	Timestamp time.Time `json:"timestamp"`
}

func GetTenantTelemetryData(tenantId string) ([]NewDeviceData, error) {
	subQuery := query.TelemetryCurrentData.Select(
		query.TelemetryCurrentData.DeviceID.As("device_id"),
		query.TelemetryCurrentData.T.Max().As("max_t"),
	).Where(
		query.TelemetryCurrentData.TenantID.Eq(tenantId),
	).Group(query.TelemetryCurrentData.DeviceID).Order(query.TelemetryCurrentData.T.Max().Desc()).Limit(3)

	type DeviceData struct {
		DeviceID string    `json:"device_id"`
		MaxT     time.Time `json:"max_t"`
	}

	var devices []DeviceData
	err := subQuery.Scan(&devices)
	if err != nil {
		return nil, err
	}

	result := make([]NewDeviceData, 0, len(devices))
	for _, device := range devices {
		deviceInfo := NewDeviceData{
			DeviceID:  device.DeviceID,
			Timestamp: device.MaxT,
		}
		result = append(result, deviceInfo)
	}

	return result, nil
}
