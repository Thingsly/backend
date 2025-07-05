```
IoT Devices → ChannelBufferSize → WriteWorkers → BatchSize → Database
     ↓              ↓                ↓            ↓
  Messages      Queue Size      Concurrent    Bulk Write
```

Flow hoạt động:

1. IoT devices gửi messages liên tục
2. ChannelBufferSize đảm bảo không bị mất messages
3. WriteWorkers xử lý song song các messages
4. PoolSize cung cấp objects để xử lý telemetry
5. BatchSize gom nhóm để tối ưu database writes



Device → Gửi Telemetry → Backend nhận → Trigger HeartbeatDeal() → Update Redis Key

```
1. Device gửi telemetry data
   ↓
2. Backend nhận telemetry
   ↓
3. Trigger HeartbeatDeal() trong goroutine riêng
   ↓
4. Update device status thành online (nếu chưa online)
   ↓
5. Set Redis key: device:{deviceID}:heartbeat với TTL
   ↓
6. DeviceListener monitor Redis expiration events
   ↓
7. Khi key expire → Device offline
   ↓
8. Publish MQTT status update: "0" (offline)
```

## Flow hoạt động

### 1. **Initialization**
```
InitDeviceStatus() → NewStatusManager() → connectWithRetry() → subscribe()
```

### 2. **Message Flow**
```
Device Status Update → MQTT Broker → StatusManager → DeviceOnline() → Database Update
```

### 3. **Topic Pattern**
```
devices/status/device001 → Payload: "1" (online) hoặc "0" (offline)
devices/status/device002 → Payload: "1" (online) hoặc "0" (offline)
```

## So sánh với DeviceListener (Redis-based)

| Aspect | StatusManager (MQTT) | DeviceListener (Redis) |
|--------|---------------------|------------------------|
| **Data Source** | MQTT messages | Redis expiration events |
| **Trigger** | Manual status updates | Automatic timeout detection |
| **Real-time** | Immediate | After timeout |
| **Reliability** | Depends on MQTT QoS | Redis persistence |
| **Scalability** | MQTT broker capacity | Redis capacity |

### **Chi tiết so sánh:**

#### **StatusManager (MQTT-based)**
- **Data Source**: Nhận messages trực tiếp từ MQTT broker
- **Trigger**: Device chủ động gửi status updates
- **Real-time**: Cập nhật ngay lập tức khi nhận message
- **Reliability**: Phụ thuộc vào MQTT QoS level và network stability
- **Scalability**: Giới hạn bởi MQTT broker capacity và network bandwidth

#### **DeviceListener (Redis-based)**
- **Data Source**: Monitor Redis expiration events
- **Trigger**: Tự động detect khi Redis key expire (timeout)
- **Real-time**: Cập nhật sau khi timeout period kết thúc
- **Reliability**: Redis persistence đảm bảo events không bị mất
- **Scalability**: Giới hạn bởi Redis capacity và memory usage

### **Workflow Comparison:**

#### **StatusManager Workflow:**
```
1. Device gửi status message qua MQTT
   ↓
2. MQTT Broker nhận và forward
   ↓
3. StatusManager subscribe và nhận message
   ↓
4. MessageHandler xử lý message
   ↓
5. DeviceOnline() update database
   ↓
6. Broadcast status change
```

#### **DeviceListener Workflow:**
```
1. Device gửi telemetry data
   ↓
2. HeartbeatDeal() set Redis key với TTL
   ↓
3. DeviceListener monitor Redis expiration events
   ↓
4. Khi key expire → nhận expiration event
   ↓
5. handleExpiredKey() xử lý event
   ↓
6. DeviceOnline() update database
   ↓
7. Broadcast status change
```

### **Integration Pattern:**

#### **Dual Monitoring System:**
```go
// Redis-based (automatic timeout)
DeviceListener → Redis expiration → Device offline

// MQTT-based (manual updates)  
StatusManager → MQTT messages → Device status update
```

#### **Redundancy Strategy:**
- **Primary**: Redis-based timeout detection (automatic)
- **Secondary**: MQTT-based manual status updates (immediate)
- **Fallback**: Cả hai method đều hoạt động song song

#### **Message Processing:**
```go
// Trong StatusManager messageHandler
subscribe.DeviceOnline(msg.Payload(), msg.Topic())
// ↓
// Update database và broadcast status change

// Trong DeviceListener handleExpiredKey
subscribe.DeviceOnline([]byte("0"), "devices/status/"+deviceID)
// ↓
// Update database và broadcast status change
```

### **Use Cases:**

#### **StatusManager Use Cases:**
1. **Device Manual Status Report**: Device chủ động report online/offline
2. **Gateway Status Updates**: Gateway update status cho sub-devices
3. **System Integration**: External systems update device status
4. **Immediate Status Changes**: Status updates không cần chờ timeout

#### **DeviceListener Use Cases:**
1. **Automatic Offline Detection**: Tự động detect device offline khi không có telemetry
2. **Heartbeat Monitoring**: Monitor device heartbeat patterns
3. **Timeout Management**: Quản lý timeout periods cho different devices
4. **Reliable Detection**: Đảm bảo device offline được detect ngay cả khi network có vấn đề

### **Configuration Examples:**

#### **StatusManager Config:**
```yaml
mqtt:
  broker: "localhost:1883"
  user: "root"
  pass: "root"
  # StatusManager tự động config từ MqttConfig
```

#### **DeviceListener Config:**
```yaml
db:
  redis:
    db1: 10  # Database cho device monitoring

# Device config trong database
{
  "online_timeout": 60,
  "heartbeat": 30
}
```

### **Performance Considerations:**

#### **StatusManager:**
- **Memory**: Low memory usage, chỉ lưu MQTT client
- **CPU**: Minimal CPU usage, chỉ xử lý incoming messages
- **Network**: Depends on MQTT message frequency
- **Latency**: Immediate processing

#### **DeviceListener:**
- **Memory**: Redis memory usage cho keys và events
- **CPU**: Minimal CPU usage, chỉ monitor Redis events
- **Network**: Depends on Redis connection
- **Latency**: Timeout-based latency



```
1. Backend → Gateway: Command/Attribute/Telemetry
2. Gateway → Sub-Devices: Forward messages
3. Sub-Devices → Gateway: Responses/Data
4. Gateway → Backend: Aggregated responses
```