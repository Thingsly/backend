```
type GatewayResponse struct {
    GatewayData   *MqttResponse           `json:"gateway_data"`
    SubDeviceData map[string]MqttResponse `json:"sub_device_data"`
}

type MqttResponse struct {
    Result  int    `json:"result"`   // 0: success, 1: failure
    Errcode string `json:"errcode"`  // Error code (optional)
    Message string `json:"message"`  // Response message
    Ts      int64  `json:"ts"`       // Timestamp
    Method  string `json:"method"`   // Method name (optional)
}
```

| Operation      | Publish Topic                                     | Response Topic                              |
| -------------- | ------------------------------------------------- | ------------------------------------------- |
| Command        | gateway/command/{deviceNumber}/{messageId}        | gateway/command/response/{messageId}        |
| Set Attributes | gateway/attributes/set/{deviceNumber}/{messageId} | gateway/attributes/set/response/{messageId} |
| Get Attributes | gateway/attributes/get/{deviceNumber}             | gateway/attributes/{messageId}              |


# Command Response Flow

```
// Backend gửi command
GatewayPublishCommandMessage(ctx, deviceInfo, messageId, command, callbackFunc)
```

Topic: gateway/command/{deviceNumber}/{messageId}

```
{
  "gateway_data": {
    "method": "restart",
    "params": {
      "delay": 5
    }
  },
  "sub_device_data": {
    "sub_device_001": {
      "method": "setMode",
      "params": {
        "mode": "auto"
      }
    }
  }
}
```

Gateway Process Command
Gateway nhận command và thực thi, sau đó gửi response:
Topic: gateway/command/response/{messageId}

```
{
  "gateway_data": {
    "result": 0,
    "errcode": "",
    "message": "Command executed successfully",
    "ts": 1703123456,
    "method": "restart"
  },
  "sub_device_data": {
    "sub_device_001": {
      "result": 0,
      "errcode": "",
      "message": "Mode set successfully",
      "ts": 1703123456,
      "method": "setMode"
    }
  }
}
```

Backend Receive Response

```
func GatewayDeviceCommandResponse(payload []byte, topic string) {
    // Extract messageId từ topic
    messageId := topicList[3]
    
    // Parse response
    result := model.GatewayResponse{}
    json.Unmarshal(attributePayload.Values, &result)
    
    // Send to waiting channel
    if ch, ok := config.GatewayResponseFuncMap[messageId]; ok {
        ch <- result
    }
}
```

