package subscribe

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	initialize "github.com/Thingsly/backend/initialize"
	dal "github.com/Thingsly/backend/internal/dal"
	"github.com/Thingsly/backend/internal/model"
	service "github.com/Thingsly/backend/internal/service"
	config "github.com/Thingsly/backend/mqtt"

	"github.com/sirupsen/logrus"
)

func DeviceAttributeReport(payload []byte, topic string) (string, string, error) {

	var messageId string
	topicList := strings.Split(topic, "/")
	if len(topicList) < 3 {
		messageId = ""
	} else {
		messageId = topicList[2]
	}

	logrus.Debug("payload:", string(payload))

	attributePayload, err := verifyPayload(payload)
	if err != nil {
		logrus.Error(err.Error())
		return "", "", err
	}
	logrus.Debug("attribute message:", attributePayload)

	device, err := initialize.GetDeviceCacheById(attributePayload.DeviceId)
	if err != nil {
		logrus.Error(err.Error())
		return "", messageId, err
	}

	reqMap := make(map[string]interface{})
	err = json.Unmarshal(attributePayload.Values, &reqMap)
	if err != nil {
		logrus.Error(err.Error())
		return device.DeviceNumber, messageId, err
	}
	err = deviceAttributesHandle(device, reqMap, topic)
	if err != nil {
		logrus.Error(err.Error())
		return device.DeviceNumber, messageId, err
	}
	return device.DeviceNumber, messageId, err

}

func deviceAttributesHandle(device *model.Device, reqMap map[string]interface{}, topic string) error {

	if device.DeviceConfigID != nil && *device.DeviceConfigID != "" {
		scriptType := "C"
		attributesBody, _ := json.Marshal(reqMap)
		newAttributesBody, err := service.GroupApp.DataScript.Exec(device, scriptType, attributesBody, topic)
		if err != nil {
			logrus.Error("Error in attribute script processing: ", err.Error())
		}
		if newAttributesBody != nil {
			err = json.Unmarshal(newAttributesBody, &reqMap)
			if err != nil {
				logrus.Error("Error in attribute script processing: ", err.Error())
			}
		}
	}

	ts := time.Now().UTC()
	logrus.Debug(device, ts)
	var (
		triggerParam  []string
		triggerValues = make(map[string]interface{})
	)
	for k, v := range reqMap {
		logrus.Debug(k, "(", v, ")")

		d := model.AttributeData{
			DeviceID: device.ID,
			Key:      k,
			T:        ts,
			TenantID: &device.TenantID,
		}

		switch value := v.(type) {
		case string:
			d.StringV = &value
		case bool:
			d.BoolV = &value
		case float64:
			d.NumberV = &value
		case int:

			f := float64(value)
			d.NumberV = &f
		case int64:

			f := float64(value)
			d.NumberV = &f
		case []interface{}, map[string]interface{}:

			if jsonBytes, err := json.Marshal(value); err == nil {
				s := string(jsonBytes)
				d.StringV = &s
			} else {
				s := fmt.Sprint(value)
				d.StringV = &s
			}
		default:

			if jsonStr, ok := tryParseAsJSON(value); ok {
				d.StringV = &jsonStr
			} else {
				s := fmt.Sprint(value)
				d.StringV = &s
			}
		}
		triggerParam = append(triggerParam, k)
		triggerValues[k] = v
		logrus.Debug("attribute data:", d)
		_, err := dal.UpdateAttributeData(&d)
		if err != nil {
			logrus.Error(err.Error())
			return err
		}
	}

	go func() {
		err := service.GroupApp.Execute(device, service.AutomateFromExt{
			TriggerParam:     triggerParam,
			TriggerValues:    triggerValues,
			TriggerParamType: model.TRIGGER_PARAM_TYPE_ATTR,
		})
		if err != nil {
			logrus.Error("Automation execution failed, err: ", err)
		}
	}()
	return nil
}

func DeviceSetAttributeResponse(payload []byte, topic string) {
	logrus.Debug("command message:", string(payload))
	var messageId string
	topicList := strings.Split(topic, "/")
	if len(topicList) < 5 {
		messageId = ""
	} else {
		messageId = topicList[4]
	}

	attributePayload, err := verifyPayload(payload)
	if err != nil {
		return
	}
	logrus.Debug("command values message:", string(attributePayload.Values))

	commandResponsePayload, err := verifyAttributeResponsePayload(attributePayload.Values)
	if err != nil {
		logrus.Error(err.Error())
		return
	}
	logrus.Debug("command response message:", commandResponsePayload)

	if ch, ok := config.MqttDirectResponseFuncMap[messageId]; ok {
		ch <- *commandResponsePayload
	}
}
