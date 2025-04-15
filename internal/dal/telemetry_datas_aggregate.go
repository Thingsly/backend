package dal

import (
	"context"
	"fmt"

	global "github.com/Thingsly/backend/pkg/global"
)

type TelemetryDatasAggregate struct {
	AggregateWindow   int64  `json:"aggregate_window"`
	AggregateFunction string `json:"aggregate_function"`
	STime             int64  `json:"s_time"`
	ETime             int64  `json:"e_time"`
	Count             int64  `json:"count"`
	DeviceID          string `json:"device_id"`
	Key               string `json:"key"`
}

func GetTelemetryDatasAggregate(_ context.Context, telemetryDatasAggregate TelemetryDatasAggregate) ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	var queryString string

	switch telemetryDatasAggregate.AggregateFunction {
	case "avg", "max", "min", "sum":
		queryString = GetQueryString1(telemetryDatasAggregate.AggregateFunction)
	case "diff":
		queryString = GetQueryString2(telemetryDatasAggregate.AggregateFunction)

	default:
		return nil, fmt.Errorf("Unsupported aggregation function: %s", telemetryDatasAggregate.AggregateFunction)
	}

	resultData := global.DB.Raw(queryString, telemetryDatasAggregate.AggregateWindow, telemetryDatasAggregate.STime, telemetryDatasAggregate.ETime, telemetryDatasAggregate.Key, telemetryDatasAggregate.DeviceID, telemetryDatasAggregate.AggregateWindow).Scan(&data)
	if resultData.Error != nil {
		return nil, resultData.Error
	}

	return data, nil

}

// Get queryString, supports average, maximum, minimum, and sum
func GetQueryString1(aggregateFunction string) string {
	queryString := fmt.Sprintf(
		`WITH TimeIntervals AS (
				SELECT 
					ts - (ts %% ?) AS x, 
					CAST(%s(number_v) AS NUMERIC(16,4)) AS y 
				FROM 
					telemetry_datas 
				WHERE 
					ts BETWEEN ? AND ? AND key = ? AND device_id = ? 
				GROUP BY 
					x
			)
			SELECT 
				x, 
				x + ? AS x2, 
				y 
			FROM 
				TimeIntervals 
			WHERE 
				y IS NOT NULL 
			ORDER BY 
				x asc;`,
		aggregateFunction,
	)
	return queryString
}

// Get queryString, supports delta (difference) calculation
func GetQueryString2(_ string) string {

	queryString := fmt.Sprintf(
		`WITH TimeIntervals AS (
				SELECT 
					ts - (ts %% ?) AS x, 
					MAX(number_v) - MIN(number_v) AS y 
				FROM 
					telemetry_datas 
				WHERE 
					ts BETWEEN ? AND ? AND key = ? AND device_id = ? 
				GROUP BY 
					x
			)
			SELECT 
				x, 
				x + ? AS x2, 
				y 
			FROM 
				TimeIntervals 
			WHERE 
				y IS NOT NULL 
			ORDER BY 
				x ASC;`,
	)

	return queryString
}
