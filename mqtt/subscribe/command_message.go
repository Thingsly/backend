package subscribe

import (
	"strings"

	config "github.com/Thingsly/backend/mqtt"

	"github.com/sirupsen/logrus"
)

func DeviceCommand(payload []byte, topic string) (string, error) {

	logrus.Debug("command message:", string(payload))
	var messageId string
	topicList := strings.Split(topic, "/")
	if len(topicList) < 4 {
		messageId = ""
	} else {
		messageId = topicList[3]
	}

	attributePayload, err := verifyPayload(payload)
	if err != nil {
		return "", err
	}
	logrus.Debug("command values message:", string(attributePayload.Values))

	commandResponsePayload, err := verifyCommandResponsePayload(attributePayload.Values)
	if err != nil {
		logrus.Error(err.Error())
		return "", err
	}
	logrus.Debug("command response message:", commandResponsePayload)

	//log := dal.CommandSetLogsQuery{}

	//if m, err := log.FilterOneHourByMessageID(messageId); err != nil || m == nil {
	//	logrus.Error(err.Error())
	//	return "", err
	//}

	//logInfo := &model.CommandSetLog{
	//	MessageID: &messageId,
	//}
	//if commandResponsePayload.Result == 0 {
	//	execFail := "3"
	//	logInfo.Status = &execFail
	//} else {
	//	execSuccess := "4"
	//	logInfo.Status = &execSuccess
	//	logInfo.RspDatum = &commandResponsePayload.Errcode
	//	logInfo.ErrorMessage = &commandResponsePayload.Message
	//}
	//err = log.Update(nil, logInfo)
	if ch, ok := config.MqttDirectResponseFuncMap[messageId]; ok {
		ch <- *commandResponsePayload
	}
	return messageId, err
}
