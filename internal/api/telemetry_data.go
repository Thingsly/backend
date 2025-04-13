package api

import (
	"encoding/json"
	"strconv"
	"sync"
	"time"

	ws_subscribe "github.com/HustIoTPlatform/backend/mqtt/ws_subscribe"
	"github.com/HustIoTPlatform/backend/pkg/constant"
	"github.com/HustIoTPlatform/backend/pkg/errcode"
	"github.com/HustIoTPlatform/backend/pkg/utils"

	model "github.com/HustIoTPlatform/backend/internal/model"
	service "github.com/HustIoTPlatform/backend/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type TelemetryDataApi struct{}

// GetCurrentData
// @Router   /api/v1/telemetry/datas/current/{id} [get]
func (*TelemetryDataApi) HandleCurrentData(c *gin.Context) {
	deviceId := c.Param("id")
	date, err := service.GroupApp.TelemetryData.GetCurrentTelemetrData(deviceId)
	if err != nil {
		c.Error(err)
		return
	}

	c.Set("data", date)
}

// Query telemetry current value based on device ID and key
// @Router /api/v1/telemetry/datas/current/keys [get]
func (*TelemetryDataApi) HandleCurrentDataKeys(c *gin.Context) {
	var req model.GetTelemetryCurrentDataKeysReq
	if !BindAndValidate(c, &req) {
		return
	}

	date, err := service.GroupApp.TelemetryData.GetCurrentTelemetrDataKeys(&req)
	if err != nil {
		c.Error(err)
		return
	}

	c.Set("data", date)
}

// ServeHistoryData
// @Router   /api/v1/telemetry/datas/history [get]
func (*TelemetryDataApi) ServeHistoryData(c *gin.Context) {
	var req model.GetTelemetryHistoryDataReq
	if !BindAndValidate(c, &req) {
		return
	}
	date, err := service.GroupApp.TelemetryData.GetTelemetrHistoryData(&req)
	if err != nil {
		c.Error(err)
		return
	}

	c.Set("data", date)
}

// DeleteData
// @Router   /api/v1/telemetry/datas [delete]
func (*TelemetryDataApi) DeleteData(c *gin.Context) {
	var req model.DeleteTelemetryDataReq
	if !BindAndValidate(c, &req) {
		return
	}
	err := service.GroupApp.TelemetryData.DeleteTelemetrData(&req)
	if err != nil {
		c.Error(err)
		return
	}

	c.Set("data", nil)
}

// GetCurrentData Retrieve the latest telemetry data based on device ID
// @Router   /api/v1/telemetry/datas/current/detail/{id} [get]
func (*TelemetryDataApi) ServeCurrentDetailData(c *gin.Context) {
	deviceId := c.Param("id")
	date, err := service.GroupApp.TelemetryData.GetCurrentTelemetrDetailData(deviceId)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", date)
}

// ServeHistoryData Device historical value query (pagination)
// @Router   /api/v1/telemetry/datas/history/pagination [get]
func (*TelemetryDataApi) ServeHistoryDataByPage(c *gin.Context) {
	var req model.GetTelemetryHistoryDataByPageReq
	if !BindAndValidate(c, &req) {
		return
	}

	// Time range limited to within one month
	// if req.EndTime.Sub(req.StartTime) > time.Hour*24*30 {
	// 	ErrorHandler(c, http.StatusBadRequest, fmt.Errorf("time range should be within 30 days"))
	// 	return
	// }

	date, err := service.GroupApp.TelemetryData.GetTelemetrHistoryDataByPageV2(&req)
	if err != nil {
		c.Error(err)
		return
	}

	c.Set("data", date)
}

// ServeHistoryData Device historical value query (pagination)
// @Router   /api/v1/telemetry/datas/history/page [get]
func (*TelemetryDataApi) ServeHistoryDataByPageV2(c *gin.Context) {
	var req model.GetTelemetryHistoryDataByPageReq
	if !BindAndValidate(c, &req) {
		return
	}

	// Time range limited to within one month
	// if req.EndTime.Sub(req.StartTime) > time.Hour*24*30 {
	// 	ErrorHandler(c, http.StatusBadRequest, fmt.Errorf("time range should be within 30 days"))
	// 	return
	// }

	date, err := service.GroupApp.TelemetryData.GetTelemetrHistoryDataByPageV2(&req)
	if err != nil {
		c.Error(err)
		return
	}

	c.Set("data", date)
}

// ServeSetLogsDataListByPage Telemetry data dispatch record query (with pagination)
// @Router   /api/v1/telemetry/datas/set/logs [get]
func (*TelemetryDataApi) ServeSetLogsDataListByPage(c *gin.Context) {
	var req model.GetTelemetrySetLogsListByPageReq
	if !BindAndValidate(c, &req) {
		return
	}

	date, err := service.GroupApp.TelemetryData.GetTelemetrSetLogsDataListByPage(&req)
	if err != nil {
		c.Error(err)
		return
	}

	c.Set("data", date)
}

// Retrieve the echo data of telemetry sent by the simulated device
// /api/v1/telemetry/datas/simulation [get]
func (*TelemetryDataApi) ServeEchoData(c *gin.Context) {
	var req model.ServeEchoDataReq
	if !BindAndValidate(c, &req) {
		return
	}

	date, err := service.GroupApp.TelemetryData.ServeEchoData(&req)
	if err != nil {
		c.Error(err)
		return
	}

	c.Set("data", date)
}

// Simulate device sending telemetry data
// /api/v1/telemetry/datas/simulation [post]
func (*TelemetryDataApi) SimulationTelemetryData(c *gin.Context) {
	var req model.SimulationTelemetryDataReq
	if !BindAndValidate(c, &req) {
		return
	}
	_, err := service.GroupApp.TelemetryData.TelemetryPub(req.Command)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", nil)
}

// ServeCurrentDataByWS Handle real-time telemetry data from devices via WebSocket
// @Router   /api/v1/telemetry/datas/current/ws [get]
func (*TelemetryDataApi) ServeCurrentDataByWS(c *gin.Context) {
	// Upgrade HTTP connection to a WebSocket connection
	conn, err := Wsupgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.Error(errcode.WithData(errcode.CodeSystemError, "WebSocket upgrade failed"))
		return
	}
	defer conn.Close()

	clientIP := conn.RemoteAddr().String()
	logrus.Info("Received a new WebSocket connection:", clientIP)

	// Read the first message sent by the client
	msgType, msg, err := conn.ReadMessage()
	if err != nil {
		logrus.Error("Failed to read the initial message:", err)
		conn.WriteMessage(websocket.TextMessage, []byte("Failed to read message"))
		return
	}

	// Parse the JSON message
	var msgMap map[string]string
	if err := json.Unmarshal(msg, &msgMap); err != nil {
		logrus.Error("Invalid JSON format:", err)
		conn.WriteMessage(msgType, []byte("Invalid message format"))
		return
	}

	// Validate required fields
	deviceID, ok := msgMap["device_id"]
	if !ok || deviceID == "" {
		conn.WriteMessage(msgType, []byte("device_id is required"))
		return
	}

	token, ok := msgMap["token"]
	if !ok || token == "" {
		conn.WriteMessage(msgType, []byte("token is required"))
		return
	}

	logrus.Infof("WebSocket connection established - Device ID: %s", deviceID)

	// Get current telemetry data
	data, err := service.GroupApp.TelemetryData.GetCurrentTelemetrDataForWs(deviceID)
	if err != nil {
		logrus.Error("Failed to retrieve telemetry data:", err)
		conn.WriteMessage(msgType, []byte("Failed to get telemetry data"))
		return
	}

	// If data is available, send it to the client
	if data != nil {
		dataByte, err := json.Marshal(data)
		if err != nil {
			logrus.Error("Failed to serialize data:", err)
			conn.WriteMessage(msgType, []byte("Failed to process telemetry data"))
			return
		}
		if err := conn.WriteMessage(msgType, dataByte); err != nil {
			logrus.Error("Failed to send data:", err)
			return
		}
	}

	// Subscribe to real-time updates
	var mu sync.Mutex
	var mqttClient ws_subscribe.WsMqttClient
	if err := mqttClient.SubscribeDeviceTelemetry(deviceID, conn, msgType, &mu); err != nil {
		logrus.Error("Failed to subscribe to telemetry data:", err)
		conn.WriteMessage(msgType, []byte("Failed to subscribe to telemetry updates"))
		return
	}
	defer mqttClient.Close()

	// Handle heartbeat message
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			// Log error message
			logrus.Error("WebSocket read error:", err)

			// Try to send error message to client
			closeMsg := []byte("connection closed due to error")
			// Use WriteControl to send a close message with a 1-second timeout
			deadline := time.Now().Add(time.Second)
			conn.WriteControl(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseInternalServerErr, string(closeMsg)),
				deadline)

			// It is now safe to exit
			return
		}

		// Handle heartbeat message
		if string(msg) == "ping" {
			mu.Lock()
			if err := conn.WriteMessage(msgType, []byte("pong")); err != nil {
				logrus.Error("Failed to send pong message:", err)

				// Try to send error message
				deadline := time.Now().Add(time.Second)
				conn.WriteControl(websocket.CloseMessage,
					websocket.FormatCloseMessage(websocket.CloseInternalServerErr, "failed to send pong"),
					deadline)

				mu.Unlock()
				return
			}
			mu.Unlock()
		}
	}
}

// ServeDeviceStatusByWS Retrieves the device online status through WebSocket
// @Summary      Get device online status
// @Description  Retrieve real-time device online status via WebSocket connection
// @Tags         Device
// @Accept       json
// @Produce      json
// @Router       /api/v1/device/online/status/ws [get]
func (*TelemetryDataApi) ServeDeviceStatusByWS(c *gin.Context) {
	// Upgrade to WebSocket connection
	conn, err := Wsupgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.Error(errcode.WithData(errcode.CodeSystemError, "WebSocket upgrade failed"))
		return
	}
	defer conn.Close()

	clientIP := conn.RemoteAddr().String()
	logrus.Info("New WebSocket connection received:", clientIP)

	// Read initial message
	msgType, msg, err := conn.ReadMessage()
	if err != nil {
		logrus.Error("Failed to read initial message:", err)
		conn.WriteMessage(websocket.TextMessage, []byte("Failed to read message"))
		return
	}

	// Parse JSON message
	var msgMap map[string]string
	if err := json.Unmarshal(msg, &msgMap); err != nil {
		logrus.Error("Invalid JSON format:", err)
		conn.WriteMessage(msgType, []byte("Invalid message format"))
		return
	}

	// Validate required fields
	deviceID, ok := msgMap["device_id"]
	if !ok || deviceID == "" {
		conn.WriteMessage(msgType, []byte("device_id is required"))
		return
	}

	token, ok := msgMap["token"]
	if !ok || token == "" {
		conn.WriteMessage(msgType, []byte("token is required"))
		return
	}

	logrus.Infof("WebSocket connection established - Device ID: %s", deviceID)
	// TODO: Validate token

	// Subscribe to device online status
	var mu sync.Mutex
	logrus.Info("User SubscribeOnlineOffline")
	var mqttClient ws_subscribe.WsMqttClient
	if err := mqttClient.SubscribeOnlineOffline(deviceID, conn, msgType, &mu); err != nil {
		logrus.Error("Failed to subscribe to device status:", err)
		conn.WriteMessage(msgType, []byte("Failed to subscribe to device status"))
		return
	}
	defer mqttClient.Close()

	// Handle heartbeats
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			// Log the error
			logrus.Error("WebSocket read error:", err)

			// Try to send an error message to the client
			closeMsg := []byte("connection closed due to error")
			// Use WriteControl to send the close message, set a 1-second timeout
			deadline := time.Now().Add(time.Second)
			conn.WriteControl(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseInternalServerErr, string(closeMsg)),
				deadline)

			// Now it's safe to exit
			return
		}

		// Handle heartbeat messages
		if string(msg) == "ping" {
			mu.Lock()
			if err := conn.WriteMessage(msgType, []byte("pong")); err != nil {
				logrus.Error("Failed to send pong message:", err)

				// Try to send an error message
				deadline := time.Now().Add(time.Second)
				conn.WriteControl(websocket.CloseMessage,
					websocket.FormatCloseMessage(websocket.CloseInternalServerErr, "failed to send pong"),
					deadline)

				mu.Unlock()
				return
			}
			mu.Unlock()
		}
	}
}

// ServeCurrentDataByKey Queries the current telemetry value by key
// @Router /api/v1/telemetry/datas/current/keys/ws [get]
func (*TelemetryDataApi) ServeCurrentDataByKey(c *gin.Context) {
	// Upgrade to WebSocket connection
	conn, err := Wsupgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.Error(errcode.WithData(errcode.CodeSystemError, "WebSocket upgrade failed"))
		return
	}
	defer conn.Close()

	clientIP := conn.RemoteAddr().String()
	logrus.Infof("New WebSocket connection received: %s", clientIP)

	// Read initial message
	msgType, msg, err := conn.ReadMessage()
	if err != nil {
		logrus.Error("Failed to read initial message:", err)
		conn.WriteMessage(websocket.TextMessage, []byte("Failed to read message"))
		return
	}

	// Parse JSON message
	var msgMap map[string]interface{}
	if err := json.Unmarshal(msg, &msgMap); err != nil {
		logrus.Error("Invalid JSON format:", err)
		conn.WriteMessage(msgType, []byte("Invalid message format"))
		return
	}

	// Validate and extract device ID
	deviceID, ok := msgMap["device_id"].(string)
	if !ok || deviceID == "" {
		conn.WriteMessage(msgType, []byte("device_id is required and must be string"))
		return
	}

	// Validate and extract keys
	keysInterface, ok := msgMap["keys"].([]interface{})
	if !ok {
		conn.WriteMessage(msgType, []byte("keys must be an array"))
		return
	}

	// Convert keys to string array
	var stringKeys []string
	for _, key := range keysInterface {
		strKey, ok := key.(string)
		if !ok || strKey == "" {
			conn.WriteMessage(msgType, []byte("keys must be non-empty strings"))
			return
		}
		stringKeys = append(stringKeys, strKey)
	}

	if len(stringKeys) == 0 {
		conn.WriteMessage(msgType, []byte("keys array cannot be empty"))
		return
	}

	// Validate token
	token, ok := msgMap["token"].(string)
	if !ok || token == "" {
		conn.WriteMessage(msgType, []byte("token is required"))
		return
	}
	// TODO: Validate token

	logrus.Infof("WebSocket connection established - Device ID: %s, Keys: %v", deviceID, stringKeys)

	// Retrieve telemetry data
	data, err := service.GroupApp.TelemetryData.GetCurrentTelemetrDataKeysForWs(deviceID, stringKeys)
	if err != nil {
		logrus.Error("Failed to get telemetry data:", err)
		conn.WriteMessage(msgType, []byte("Failed to get telemetry data"))
		return
	}

	// Send data to client
	if data != nil {
		dataByte, err := json.Marshal(data)
		if err != nil {
			logrus.Error("Failed to serialize data:", err)
			conn.WriteMessage(msgType, []byte("Failed to process telemetry data"))
			return
		}
		if err := conn.WriteMessage(msgType, dataByte); err != nil {
			logrus.Error("Failed to send data:", err)
			return
		}
	}

	// Subscribe to telemetry updates
	var mu sync.Mutex
	var mqttClient ws_subscribe.WsMqttClient
	if err := mqttClient.SubscribeDeviceTelemetryByKeys(deviceID, conn, msgType, &mu, stringKeys); err != nil {
		logrus.Error("Failed to subscribe to telemetry data:", err)
		conn.WriteMessage(msgType, []byte("Failed to subscribe to telemetry updates"))
		return
	}
	defer mqttClient.Close()

	// Handle heartbeats
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			// Log the error
			logrus.Error("WebSocket read error:", err)

			// Try to send an error message to the client
			closeMsg := []byte("connection closed due to error")
			// Use WriteControl to send the close message, set a 1-second timeout
			deadline := time.Now().Add(time.Second)
			conn.WriteControl(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseInternalServerErr, string(closeMsg)),
				deadline)

			// Now it's safe to exit
			return
		}

		// Handle heartbeat messages
		if string(msg) == "ping" {
			mu.Lock()
			if err := conn.WriteMessage(msgType, []byte("pong")); err != nil {
				logrus.Error("Failed to send pong message:", err)

				// Try to send an error message
				deadline := time.Now().Add(time.Second)
				conn.WriteControl(websocket.CloseMessage,
					websocket.FormatCloseMessage(websocket.CloseInternalServerErr, "failed to send pong"),
					deadline)

				mu.Unlock()
				return
			}
			mu.Unlock()
		}
	}
}

// ServeStatisticData
// @Router   /api/v1/telemetry/datas/statistic [get]
func (*TelemetryDataApi) ServeStatisticData(c *gin.Context) {
	var req model.GetTelemetryStatisticReq
	if !BindAndValidate(c, &req) {
		return
	}

	date, err := service.GroupApp.TelemetryData.GetTelemetrServeStatisticData(&req)
	if err != nil {
		c.Error(err)
		return
	}

	c.Set("data", date)
}

// /api/v1/telemetry/datas/pub
func (*TelemetryDataApi) TelemetryPutMessage(c *gin.Context) {
	var req model.PutMessage
	if !BindAndValidate(c, &req) {
		return
	}

	userClaims := c.MustGet("claims").(*utils.UserClaims)
	err := service.GroupApp.TelemetryData.TelemetryPutMessage(c, userClaims.ID, &req, strconv.Itoa(constant.Manual))
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", nil)
}

// /api/v1/telemetry/datas/msg/count
func (*TelemetryDataApi) ServeMsgCountByTenant(c *gin.Context) {
	userClaims := c.MustGet("claims").(*utils.UserClaims)
	if userClaims.TenantID == "" {
		c.Error(errcode.New(201001))
		return
	}
	cnt, err := service.GroupApp.TelemetryData.ServeMsgCountByTenantId(userClaims.TenantID)
	if err != nil {
		c.Error(err)
		return
	}

	c.Set("data", map[string]interface{}{"msg": cnt})
}
