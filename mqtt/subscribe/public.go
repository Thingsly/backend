package subscribe

import (
	"encoding/json"
	"errors"

	"github.com/HustIoTPlatform/backend/internal/model"

	"github.com/sirupsen/logrus"
)

type publicPayload struct {
	DeviceId string `json:"device_id"`
	Values   []byte `json:"values"`
}

func verifyPayload(body []byte) (*publicPayload, error) {
	payload := &publicPayload{
		Values: make([]byte, 0),
	}
	if err := json.Unmarshal(body, payload); err != nil {
		logrus.Error("Failed to parse message:", err)
		return payload, err
	}
	if len(payload.DeviceId) == 0 {
		return payload, errors.New("DeviceId cannot be empty:" + payload.DeviceId)
	}
	if len(payload.Values) == 0 {
		return payload, errors.New("The values message content cannot be empty")
	}
	return payload, nil
}

func verifyEventPayload(values interface{}) (*model.EventInfo, error) {
	eventPayload := &model.EventInfo{}
	if err := json.Unmarshal(values.([]byte), eventPayload); err != nil {
		logrus.Error("Failed to parse message:", err)
		return eventPayload, err
	}
	if len(eventPayload.Method) == 0 {
		return eventPayload, errors.New("Method cannot be empty:" + eventPayload.Method)
	}
	if len(eventPayload.Params) == 0 {
		return eventPayload, errors.New("Params message content cannot be empty")
	}
	return eventPayload, nil
}

func verifyCommandResponsePayload(values interface{}) (*model.MqttResponse, error) {
	payload := &model.MqttResponse{}
	if err := json.Unmarshal(values.([]byte), payload); err != nil {
		logrus.Error("Failed to parse message:", err)
		return payload, err
	}
	if len(payload.Method) == 0 {
		return payload, errors.New("Method cannot be empty:" + payload.Method)
	}
	if len(payload.Message) == 0 {
		return payload, errors.New("Params message content cannot be empty")
	}
	return payload, nil
}

func verifyAttributeResponsePayload(values interface{}) (*model.MqttResponse, error) {
	payload := &model.MqttResponse{}
	if err := json.Unmarshal(values.([]byte), payload); err != nil {
		logrus.Error("Failed to parse message:", err)
		return payload, err
	}
	return payload, nil
}
