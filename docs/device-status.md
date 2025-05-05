# Device Online/Offline Status Management System

## Overview

This document describes how the system manages and detects device online/offline status using a combination of MQTT, Redis, and PostgreSQL.

## Architecture

### Components

- **MQTT Broker**: Handles device communication and telemetry data
- **Redis**: Manages heartbeat expiration and status notifications
- **PostgreSQL**: Stores device status and configuration
- **WebSocket**: Real-time status updates to clients

### Data Structure

#### Redis Keys

```redis
device:{deviceId}:heartbeat  // Expiration key for heartbeat device
device:{deviceId}:timeout    // Expiration key for timed-out device
```

#### PostgreSQL Table (devices)

```sql
CREATE TABLE devices (
    id VARCHAR(36) PRIMARY KEY,
    is_online INT,           // Online status (1) / offline (0)
    device_config_id VARCHAR(36), // Device configuration ID
    // ... other fields
);
```

#### Device Configuration

```json
{
    "heartbeat": 60,        // Heartbeat interval in seconds
    "online_timeout": 5     // Timeout period in minutes
}
```

## Workflow

### 1. Device Connection and Telemetry

When a device connects to the MQTT broker, it sends telemetry data through the topic `devices/telemetry/{deviceId}`. The system processes this as follows:

```go
func TelemetryMessagesHandle(device *model.Device, telemetryBody []byte, topic string) {
    // Process telemetry data
    // ...
    
    // Handle heartbeat in a separate goroutine
    go HeartbeatDeal(device)
    
    // Process telemetry data
    // ...
}
```

### 2. Heartbeat Processing

When telemetry data is received, the system calls `HeartbeatDeal()` to handle the heartbeat:

```go
func HeartbeatDeal(device *model.Device) {
    // 1. Get device configuration
    deviceConfig, err := dal.GetDeviceConfigByID(*device.DeviceConfigID)
    
    // 2. Parse configuration
    type OtherConfig struct {
        OnlineTimeout int `json:"online_timeout"`
        Heartbeat     int `json:"heartbeat"`
    }
    
    var otherConfig OtherConfig
    json.Unmarshal([]byte(*deviceConfig.OtherConfig), &otherConfig)

    // 3. Handle heartbeat
    if otherConfig.Heartbeat > 0 {
        // Update to online if currently offline
        if device.IsOnline != int16(1) {
            DeviceOnline([]byte("1"), "devices/status/"+device.ID)
        }

        // Set Redis key with TTL = heartbeat time
        global.STATUS_REDIS.Set(
            fmt.Sprintf("device:%s:heartbeat", device.ID),
            1,
            time.Duration(otherConfig.Heartbeat)*time.Second,
        )
    }
}
```

### 3. Online Status Update

When a device is updated to online status:

```go
func DeviceOnline(payload []byte, topic string) {
    // 1. Validate status
    status, err := validateStatus(payload) // 0 or 1
    
    // 2. Update database
    deviceId := strings.Split(topic, "/")[2]
    err = dal.UpdateDeviceStatus(deviceId, status)
    
    // 3. Send WebSocket notification
    go toUserClient(device, status)
    
    // 4. Execute automation if configured
    go func() {
        loginStatus := "ON-LINE"
        if status == 0 {
            loginStatus = "OFF-LINE" 
        }
        service.GroupApp.Execute(device, service.AutomateFromExt{
            TriggerParamType: model.TRIGGER_PARAM_TYPE_STATUS,
            TriggerValues: map[string]interface{}{
                "login": loginStatus,
            },
        })
    }()
}
```

### 4. Offline Detection

When the heartbeat Redis key expires (device hasn't sent telemetry within heartbeat interval):

```go
func (l *DeviceListener) handleExpiredKey(msg *redis.Message) {
    // Update device to offline when key expires
    deviceID := strings.Split(msg.Payload, ":")[1]
    subscribe.DeviceOnline([]byte("0"), "devices/status/"+deviceID)
}
```

## Example Scenario

Consider a device with the following configuration:

```json
{
    "heartbeat": 60,        // Send heartbeat every 60 seconds
    "online_timeout": 5     // Timeout after 5 minutes
}
```

### Device Online Flow

1. Device sends telemetry data every 60 seconds
2. System receives telemetry → calls HeartbeatDeal()
3. Updates online status in database
4. Sets Redis key with 60s TTL

### Device Offline Flow

1. If no telemetry received within 60 seconds
2. Redis key expires
3. System updates device status to offline
4. Sends WebSocket notification
5. Executes configured automation

## Key Features

### 1. Automatic Detection

- Uses Redis expiration for automatic offline detection
- No need for explicit offline notifications from devices

### 2. Performance

- Handles heartbeat in separate goroutines
- Uses Redis to reduce database load
- Efficient concurrent processing

### 3. Reliability

- Retry mechanism for MQTT/Redis connection losses
- Concurrent handling of multiple devices
- Robust error handling

### 4. Scalability

- Supports multiple device types with different configurations
- Easy to add new notification types and handlers
- Flexible automation system

## Best Practices

1. **Configuration**
   - Set appropriate heartbeat intervals based on device type
   - Configure timeout periods according to business requirements
   - Monitor system performance and adjust as needed

2. **Monitoring**
   - Monitor Redis key expiration events
   - Track device connection status
   - Log important status changes

3. **Error Handling**
   - Implement proper error handling for all operations
   - Log errors for debugging
   - Implement retry mechanisms for failed operations

4. **Security**
   - Validate all incoming messages
   - Implement proper authentication
   - Secure WebSocket connections

## Conclusion

The device online/offline status management system provides a robust and efficient way to track device connectivity. By leveraging MQTT for communication, Redis for expiration handling, and PostgreSQL for status storage, the system ensures reliable device status monitoring with minimal overhead.
