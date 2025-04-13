package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/HustIoTPlatform/backend/internal/dal"
	"github.com/HustIoTPlatform/backend/internal/model"
	"github.com/go-basic/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type MessagePush struct {
}

func (receiver *MessagePush) CreateMessagePush(req *model.CreateMessagePushReq, userId string) error {
	exists, err := dal.GetMessagePushMangeExists(userId, req.PushId)
	if err != nil {
		return err
	}
	if exists {
		return dal.ActiveMessagePushMange(userId, req.PushId, req.DeviceType)
	}
	return dal.CreateMessagePushMange(&model.MessagePushManage{
		ID:         uuid.New(),
		UserID:     userId,
		PushID:     req.PushId,
		DeviceType: req.DeviceType,
		Status:     1,
		CreateTime: time.Now(),
	})
}

func (receiver *MessagePush) MessagePushMangeLogout(req *model.MessagePushMangeLogoutReq, userId string) error {
	exists, err := dal.GetMessagePushMangeExists(userId, req.PushId)
	if err != nil {
		return err
	}
	if exists {
		return dal.LogoutMessagePushMange(userId, req.PushId)
	}
	return errors.New("The current user's push ID does not exist")
}

func (receiver *MessagePush) GetMessagePushConfig() (*model.MessagePushConfigRes, error) {
	return dal.GetMessagePushConfig()
}

func (receiver *MessagePush) SetMessagePushConfig(req *model.MessagePushConfigReq) error {
	return dal.SetMessagePushConfig(req)
}

// MessagePushSend sends a push message to the configured URL
// @params message model.MessagePushSend
// @return res string, err error
func (receiver *MessagePush) MessagePushSend(message model.MessagePushSend) (res string, err error) {
	// Retrieve push configuration
	config, err := dal.GetMessagePushConfig()
	if err != nil {
		return
	}
	// If URL is empty, return
	if config.Url == "" {
		return
	}
	// Marshal message to JSON
	jsonData, _ := json.Marshal(message)
	logrus.Debug(fmt.Sprintf("Sending URL: %s, Request parameters: %s", config.Url, string(jsonData)))

	// Create new HTTP POST request with JSON body
	req, err := http.NewRequest("POST", config.Url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set request header to specify content type as JSON
	req.Header.Set("Content-Type", "application/json")

	// Send the request using HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	// Log the response
	logrus.Debug("Response:", string(body))
	return string(body), nil
}

// AlarmMessagePushSend sends an alarm push message to the user
// @params triggered string, alarmConfigId string, deviceInfo *model.Device
func (receiver *MessagePush) AlarmMessagePushSend(triggered, alarmConfigId string, deviceInfo *model.Device) {
	// Get user push IDs based on tenant ID
	pushManges, err := dal.GetUserMessagePushId(deviceInfo.TenantID)
	if err != nil {
		logrus.Error("Failed to query user push IDs:", err)
		return
	}

	// If no push IDs found, return
	if len(pushManges) == 0 {
		return
	}
	logrus.Debug(fmt.Sprintf("pushManges:%#v", len(pushManges)))

	// Prepare message with alarm details
	message := model.MessagePushSend{
		Title:   fmt.Sprintf("Alarm: %v", triggered),
		Content: deviceInfo.DeviceNumber,
		Payload: model.MessagePushSendPayload{
			AlarmConfigId: alarmConfigId,
			TenantId:      deviceInfo.TenantID,
		},
	}

	// Loop through all push IDs and send the message
	for _, v := range pushManges {
		// Skip if push ID is empty
		if v.PushID == "" {
			continue
		}
		message.CIds = v.PushID
		receiver.MessagePushSendAndLog(message, v, 1)
	}
}

// MessagePushSendAndLog sends a message and logs the result
// @params message model.MessagePushSend, mange model.MessagePushManage, messageType int64
func (receiver *MessagePush) MessagePushSendAndLog(message model.MessagePushSend, mange model.MessagePushManage, messageType int64) {
	// Send the message
	res, err := receiver.MessagePushSend(message)
	// Marshal message into JSON format for logging
	contents, _ := json.Marshal(message)

	// Create log entry for the message
	log := model.MessagePushLog{
		ID:          uuid.New(),
		UserID:      mange.UserID,
		MessageType: messageType,
		Content:     string(contents),
		CreateTime:  time.Now(),
	}

	// Check if there was an error sending the message
	if err != nil {
		log.ErrMessage = err.Error()
		log.Status = 2 // Status 2 indicates failure
	} else {
		// Parse the response into a map
		var result map[string]interface{}
		err = json.Unmarshal([]byte(res), &result)
		if err != nil {
			logrus.Error("Failed to parse the response to map:", err, "Response:", res)
			log.Status = 2 // Failure status
			log.ErrMessage = fmt.Sprintf("Failed to parse response to map: %v, Response: %v", err, res)
		} else if errCode, ok := result["errCode"]; ok {
			// Check for error codes in the response
			switch value := errCode.(type) {
			case float64:
				if value == 0 {
					log.Status = 1 // Status 1 indicates success
					log.ErrMessage = res
				} else {
					log.Status = 2 // Failure status
					log.ErrMessage = res
				}
			default:
				log.Status = 2 // Failure status
				log.ErrMessage = res
			}
		} else {
			log.Status = 2 // Failure status
			log.ErrMessage = res
		}
	}

	// Save the log entry to the database
	err = dal.MessagePushSendLogSave(&log)
	if err != nil {
		logrus.Error("Failed to record message push log:", err)
	}

	// Update the message push management record with the latest push time and error count
	updates := map[string]interface{}{
		"last_push_time": time.Now(),
	}
	if log.Status == 1 {
		updates["err_count"] = 0
	} else {
		updates["err_count"] = gorm.Expr("err_count + ?", 1)
	}
	err = dal.MessagePushMangeSendUpdate(mange.ID, updates)
	if err != nil {
		logrus.Error("Failed to update the last push time in message push management:", err)
	}
}

// MessagePushMangeClear clears the inactive or unsuccessful push users
func (receiver *MessagePush) MessagePushMangeClear() {
	// Retrieve users who have been marked as inactive for 7 days without being online
	err := dal.GetMessagePushMangeInactiveWithSeven()
	if err != nil {
		logrus.Error("Failed to retrieve users marked as inactive for 7 days and not online:", err)
		return
	}

	// Retrieve users with more than 3 consecutive push failures in 30 days and mark them as inactive
	err = dal.GetMessagePushMangeInactive()
	if err != nil {
		logrus.Error("Failed to retrieve users with more than 3 consecutive push failures and 30 days of inactivity:", err)
		return
	}
}
