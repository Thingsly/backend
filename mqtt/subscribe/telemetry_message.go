package subscribe

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	initialize "github.com/Thingsly/backend/initialize"
	dal "github.com/Thingsly/backend/internal/dal"
	model "github.com/Thingsly/backend/internal/model"
	service "github.com/Thingsly/backend/internal/service"
	config "github.com/Thingsly/backend/mqtt"
	"github.com/Thingsly/backend/mqtt/publish"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Queue processing, data insertion into the database
func MessagesChanHandler(messages <-chan map[string]interface{}) {
	logrus.Println("Batch write coroutine started")
	var telemetryList []*model.TelemetryData

	batchSize := config.MqttConfig.Telemetry.BatchSize
	logrus.Println("Maximum number of records per batch:", batchSize)
	for {
		for i := 0; i < batchSize; i++ {
			// Retrieve message
			// logrus.Debug("Pipeline message count:", len(messages))
			message, ok := <-messages
			if !ok {
				break
			}

			// If a different database is configured, telemetry data will not be written to the original database
			dbType := viper.GetString("grpc.tptodb_type")
			if dbType == "TSDB" || dbType == "KINGBASE" || dbType == "POLARDB" {
				continue
			}

			// Check if telemetry data exists in the message
			if tskv, ok := message["telemetryData"].(model.TelemetryData); ok {
				telemetryList = append(telemetryList, &tskv)
			} else {
				logrus.Error("Pipeline message format error")
			}

			// If no messages in the pipeline, check for data insertion
			if len(messages) > 0 {
				continue
			}
			break
		}

		// If telemetryList has data, insert into the database
		if len(telemetryList) > 0 {
			logrus.Info("Batch insert telemetry data records:", len(telemetryList))
			err := dal.CreateTelemetrDataBatch(telemetryList)
			if err != nil {
				logrus.Error(err)
			}

			// Update the current value table
			err = dal.UpdateTelemetrDataBatch(telemetryList)
			if err != nil {
				logrus.Error(err)
			}

			// Clear telemetryList
			telemetryList = []*model.TelemetryData{}
		}
	}
}

// Process messages
func TelemetryMessages(payload []byte, topic string) {
	// If a different database is configured, do not write telemetry data to the original database
	dbType := viper.GetString("grpc.tptodb_type")
	if dbType == "TSDB" || dbType == "KINGBASE" || dbType == "POLARDB" {
		logrus.Infof("do not insert db for dbType:%v", dbType)
		return
	}

	logrus.Debugln(string(payload))
	// Validate message validity
	telemetryPayload, err := verifyPayload(payload)
	if err != nil {
		logrus.Error(err.Error(), topic)
		return
	}
	device, err := initialize.GetDeviceCacheById(telemetryPayload.DeviceId)
	if err != nil {
		logrus.Error(err.Error())
		return
	}
	TelemetryMessagesHandle(device, telemetryPayload.Values, topic)
}

// Attempt to parse the value as a JSON string
func tryParseAsJSON(value interface{}) (string, bool) {
	// Try converting the value to a string
	str := fmt.Sprint(value)

	// Check if it looks like JSON (simple check)
	trimmed := strings.TrimSpace(str)
	if (strings.HasPrefix(trimmed, "{") && strings.HasSuffix(trimmed, "}")) ||
		(strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]")) {

		// Try parsing as JSON
		var js interface{}
		if err := json.Unmarshal([]byte(str), &js); err == nil {
			// If it's valid JSON, reformat it to standard JSON
			if jsonBytes, err := json.Marshal(js); err == nil {
				return string(jsonBytes), true
			}
		}
	}

	return str, false
}

// Handles telemetry messages
func TelemetryMessagesHandle(device *model.Device, telemetryBody []byte, topic string) {
	// TODO: Script processing
	if device.DeviceConfigID != nil && *device.DeviceConfigID != "" {
		// Execute script to modify telemetry data based on device configuration
		newtelemetryBody, err := service.GroupApp.DataScript.Exec(device, "A", telemetryBody, topic)
		if err != nil {
			logrus.Error(err.Error())
			return
		}
		if newtelemetryBody != nil {
			telemetryBody = newtelemetryBody
		}
	}

	// Forward telemetry message to another service
	err := publish.ForwardTelemetryMessage(device.ID, telemetryBody)
	if err != nil {
		logrus.Error("telemetry forward error:", err.Error())
	}

	// Handle heartbeat in a separate goroutine
	go HeartbeatDeal(device)

	// Convert byte slice to map for further processing
	reqMap := make(map[string]interface{})
	err = json.Unmarshal(telemetryBody, &reqMap)
	if err != nil {
		logrus.Error(err.Error())
		return
	}

	// Get the current timestamp in milliseconds
	ts := time.Now().UTC()
	milliseconds := ts.UnixNano() / int64(time.Millisecond)
	logrus.Debug(device, ts)

	var (
		triggerParam  []string
		triggerValues = make(map[string]interface{})
	)

	// Loop through the message map and process each key-value pair
	for k, v := range reqMap {
		logrus.Debug(k, "(", v, ")")
		d := model.TelemetryData{
			DeviceID: device.ID,
			Key:      k,
			T:        milliseconds,
			TenantID: &device.TenantID,
		}

		// Set the correct value field based on the type of the value
		switch value := v.(type) {
		case string:
			d.StringV = &value
		case bool:
			d.BoolV = &value
		case float64:
			d.NumberV = &value
		case int:
			// Handle integer types
			f := float64(value)
			d.NumberV = &f
		case int64:
			// Handle long integer types
			f := float64(value)
			d.NumberV = &f
		case []interface{}, map[string]interface{}:
			// Handle JSON objects or arrays
			if jsonBytes, err := json.Marshal(value); err == nil {
				s := string(jsonBytes)
				d.StringV = &s
			} else {
				s := fmt.Sprint(value)
				d.StringV = &s
			}
		default:
			// Try to detect if the value is a JSON string
			if jsonStr, ok := tryParseAsJSON(value); ok {
				d.StringV = &jsonStr
			} else {
				s := fmt.Sprint(value)
				d.StringV = &s
			}
		}

		// Prepare trigger parameters for automation
		triggerParam = append(triggerParam, k)
		triggerValues[k] = v

		// Send telemetry data to the channel for batch processing
		TelemetryMessagesChan <- map[string]interface{}{
			"telemetryData": d,
		}
	}

	// Execute automation in a separate goroutine
	go func() {
		err = service.GroupApp.Execute(device, service.AutomateFromExt{
			TriggerParamType: model.TRIGGER_PARAM_TYPE_TEL,
			TriggerParam:     triggerParam,
			TriggerValues:    triggerValues,
		})
		if err != nil {
			logrus.Error("Automation execution failed, err: %w", err)
		}
	}()
}
