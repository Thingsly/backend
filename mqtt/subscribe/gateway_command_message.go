package subscribe

import (
	"encoding/json"
	"strings"

	"github.com/Thingsly/backend/internal/model"
	config "github.com/Thingsly/backend/mqtt"

	"github.com/sirupsen/logrus"
)

// @description Gateway device command response
// param payload []byte
// param topic string
// @return error
func GatewayDeviceCommandResponse(payload []byte, topic string) {
	var messageId string
	topicList := strings.Split(topic, "/")
	if len(topicList) >= 4 {
		messageId = topicList[3]
	}
	if messageId == "" {
		return
	}
	logrus.Debug("payload:", string(payload))

	attributePayload, err := verifyPayload(payload)
	if err != nil {
		return
	}
	logrus.Debug("payload:", string(attributePayload.Values))
	result := model.GatewayResponse{}
	if err := json.Unmarshal(attributePayload.Values, &result); err != nil {
		return
	}
	if ch, ok := config.GatewayResponseFuncMap[messageId]; ok {
		ch <- result
	}
}
