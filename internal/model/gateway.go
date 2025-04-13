package model

const (
	MQTT_RESPONSE_RESULT_SUCESS = 0 // Success
	MQTT_RESPONSE_RESULT_FAIL   = 1 // Failure
)

// Event/Command
type EventInfo struct {
	Method string                 `json:"method"`
	Params map[string]interface{} `json:"params"`
}

type MqttResponse struct {
	Result  int    `json:"result"`
	Errcode string `json:"errcode"`
	Message string `json:"message"`
	Ts      int64  `json:"ts"`
	Method  string `json:"method"`
}

type GatewayCommandPulish struct {
	GatewayData   *EventInfo            `json:"gateway_data"`
	SubDeviceData *map[string]EventInfo `json:"sub_device_data"`
}

type GatewayPublish struct {
	GatewayData   *map[string]interface{}            `json:"gateway_data"`
	SubDeviceData *map[string]map[string]interface{} `json:"sub_device_data"`
}

type GatewayAttributeGet struct {
	GatewayData   *[]string            `json:"gateway_data"`
	SubDeviceData *map[string][]string `json:"sub_device_data"`
}

type GatewayResponse struct {
	GatewayData   *MqttResponse           `json:"gateway_data"`
	SubDeviceData map[string]MqttResponse `json:"sub_device_data"`
}
