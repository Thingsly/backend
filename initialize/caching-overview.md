# System Cache Description

> **Note:** All cache access methods are implemented in `initialize/redis_init.go`.

## Common Cache Methods

- Serialize a map or struct into a JSON string and store it in Redis  
  - `initialize/redis_init.go` → `SetRedisForJsondata`
- Retrieve JSON from Redis and deserialize it into the target object  
  - `initialize/redis_init.go` → `GetRedisForJsondata`

## Cache Data Overview

### Device Information

- `key`: `deviceId`

#### Device ID Cache

- `key`: `voucher`  
- This cache has been removed from VerneMQ due to issues caused by its complexity and unreliability.

### Data Scripts

- **Description:**  
  Data scripts and device configurations have a many-to-one relationship. However, only one telemetry script is allowed per configuration. Data scripts contain a `device_config_id`.

- **Purpose of Caching:**  
  When telemetry data is reported, the system checks whether `{deviceId}_telemetry_script_flag` exists and is non-empty. If so, it executes the script before storing the data.

- **Business Logic:**  
  - When receiving data from a device, it must be processed immediately to convert it to the platform standard.
  - When sending data to a device, it must be converted into the device’s native format beforehand.

- **Cache Keys:**
  - `{scriptId}_script`
  - `{deviceId}_telemetry_script_flag`  
    - This cache is always deleted upon create/update actions.
    - Must be deleted when the device configuration is changed.
    - Enable/disable script operations delete this flag.
    - Empty string means no available telemetry script; otherwise, the value is `script_id`.
    - If this key is missing, it indicates cache loss and must be repaired.

### Device Configuration

- `key`: `{deviceConfigId}_config`

### Automation

#### 1. Scene Automation for Single Device Telemetry

##### cache1

- `key`: `auto_{device_id}`  
- On telemetry report, the system checks if this cache exists via `cache.GetCacheByDeviceId(deviceInfo.ID)`. If not, it’s considered lost or missing and repaired using `cache.SetCacheByDeviceId(deviceId, groups, actionInfos)`.

- **Post-operation requirements:**
  - Start/stop operations must update automation flags.
  - Newly added automations are disabled by default.
  - Edited automations are disabled and must explicitly call the stop method.
  - Deletions must call the stop method.

- **Structure:**
  ```json
  [
    {
      "scene_automation_id": "xxxxx",
      "group_id": ["1", "2", "3"]
    },
    {
      "scene_automation_id": "xxxxx",
      "group_id": ["E", "F"]
    }
  ]
  ```

##### cache2

- `key`: `auto_group_{group_id}`  
- Stores rows from `device_trigger_condition` with condition types 10 (single device) and 22 (time range).

  ```json
  [{data}, {data}]
  ```

##### cache3

- `key`: `auto_action_{scene_automation_id}`  
- Stores actions and group IDs associated with the automation scene.

  ```json
  {
    "group_id": ["", ""],
    "actions": [{data}, {data}]
  }
  ```

##### cache4: Alarm-related Cache in Scene Automation

- `key`: `auto_alarm_{group_id}`

  ```json
  {
    "scene_automation_id": "xxx",
    "alarm_config_id_list": ["xxx", "xxx"],
    "alarm_device_id_list": ["xxx"]
  }
  ```

> ⚠️ A scene automation may include multiple alarms. This cache stores a list of alarm config IDs.

**Alarm history update logic** varies depending on the condition type (single device, cross-device, or device configuration). For each case, alarm history is updated or recovered based on group ID and involved device list.

**Usage Guidelines:**
- During scene automation operation:
  - On telemetry report:
    - Use cache1 to retrieve group_id.
    - Use cache2 to evaluate conditions.
    - If triggered:
      - Check if group_id exists in the alarm cache.
      - Based on presence and condition type, determine whether to insert or update alarm history.
    - If not triggered:
      - Remove or update alarm cache as appropriate.
- On automation deactivation:
  - Delete associated group_id cache entries.
- On automation or alarm config deletion:
  - Delete cache and related alarm history.
- **Cache Recovery Method:** Reconstruct from `alarm_history` using `alarm_config_id` and `alarm_device_id_list`.

#### 2. Scene Automation for Device Type (Device Config)

- Used when conditions involve only device type (group condition includes device type or device type + time range).

- When triggered, the ID of the triggering device must be passed to the action logic to ensure correct control.

##### cache1

- `key`: `auto_devconfig_{device_config_id}`  
- Stores a list of condition group IDs and corresponding automation IDs.

  ```json
  [
    {
      "scene_automation_id": "xxxxx",
      "group_id": ["group_id_1", "group_id_2", "group_id_3"]
    },
    {
      "scene_automation_id": "xxxxx",
      "group_id": ["group_id_E", "group_id_F"]
    }
  ]
  ```

##### cache2

- `key`: `auto_devconfig_group_{group_id}`  
- Stores rows from `device_trigger_condition` with condition types 11 (device type) and 22 (time range).

  ```json
  [{data}, {data}]
  ```

##### cache3

- `key`: `auto_devconfig_action_{scene_automation_id}`  
- Stores:
  - `group_id`: condition group list
  - `actions`: action list from `action_info` related to the automation ID

  ```json
  {
    "group_id": ["group_id_1", "group_id_2", "group_id_3"],
    "actions": [{data}, {data}]
  }
  ```

#### 3. Attribute and Event Reporting Automation

These follow the same classification: single device and device type.

##### Single Device

- `cache1`: `auto_attr_event_{device_id}`  
  - All other logic reuses the single device telemetry automation caches.

##### Device Type

- `cache1`: `auto_devconfig_attr_event_{device_id}`  
  - All other logic reuses the device type telemetry automation caches.

## Example automate cache

>> KEYS "automate:v3:*"

1) "automate:v3:one:_:4b0d5805-21e8-7be9-4a74-8372e92eb630"
2) "automate:v3:one:_action_:de2137ca-fa12-02da-213e-a3a9180cd9a2"
3) "automate:v3:one:_group_:1a56940d-67cc-79b6-aa2b-dbc0b8444c6d"
>> GET "automate:v3:one:_:4b0d5805-21e8-7be9-4a74-8372e92eb630"

"[{\"scene_automation_id\":\"de2137ca-fa12-02da-213e-a3a9180cd9a2\",\"group_id\":[\"1a56940d-67cc-79b6-aa2b-dbc0b8444c6d\"]}]"
>> GET "automate:v3:one:_action_:de2137ca-fa12-02da-213e-a3a9180cd9a2"

"{\"group_id\":[\"1a56940d-67cc-79b6-aa2b-dbc0b8444c6d\"],\"actions\":[{\"id\":\"25c94f7d-e9a4-5929-a7ff-5b8eac0b4b70\",\"scene_automation_id\":\"de2137ca-fa12-02da-213e-a3a9180cd9a2\",\"action_target\":\"ae9afb32-99a1-09e8-0a43-fce70fe1b068\",\"action_type\":\"30\",\"action_param_type\":\"\",\"action_param\":\"\",\"action_value\":\"\",\"remark\":null},{\"id\":\"21003bf9-b5a1-cf09-f0d4-9da051db9633\",\"scene_automation_id\":\"de2137ca-fa12-02da-213e-a3a9180cd9a2\",\"action_target\":\"4b0d5805-21e8-7be9-4a74-8372e92eb630\",\"action_type\":\"10\",\"action_param_type\":\"attributes\",\"action_param\":\"status\",\"action_value\":\"{\\\"status\\\":\\\"ok\\\"}\",\"remark\":null}]}"
>> GET "automate:v3:one:_group_:1a56940d-67cc-79b6-aa2b-dbc0b8444c6d"

"[{\"id\":\"4cc2fd4e-6552-40a8-accd-95e215c678ca\",\"scene_automation_id\":\"de2137ca-fa12-02da-213e-a3a9180cd9a2\",\"enabled\":\"Y\",\"group_id\":\"1a56940d-67cc-79b6-aa2b-dbc0b8444c6d\",\"trigger_condition_type\":\"10\",\"trigger_source\":\"4b0d5805-21e8-7be9-4a74-8372e92eb630\",\"trigger_param_type\":\"telemetry\",\"trigger_param\":\"humidity\",\"trigger_operator\":\"\\u003e\",\"trigger_value\":\"20\",\"remark\":null,\"tenant_id\":\"d616bcbb\"}]"

### 1. Device Cache - automate:v3:one:_:4b0d5805-21e8-7be9-4a74-8372e92eb630

```json
[{
  "scene_automation_id": "de2137ca-fa12-02da-213e-a3a9180cd9a2",
  "group_id": ["1a56940d-67cc-79b6-aa2b-dbc0b8444c6d"]
}]
```

Ý nghĩa:

- Device ID: 4b0d5805-21e8-7be9-4a74-8372e92eb630 (device này có automation)
- Scene Automation ID: de2137ca-fa12-02da-213e-a3a9180cd9a2 (ID của scene automation)
- Group ID: 1a56940d-67cc-79b6-aa2b-dbc0b8444c6d (group điều kiện trigger)

Mục đích: Cache này cho biết device nào có automation nào, giúp tìm nhanh automation khi device gửi data.

### 2. Action Cache - automate:v3:one:_action_:de2137ca-fa12-02da-213e-a3a9180cd9a2

```json
{
  "group_id": ["1a56940d-67cc-79b6-aa2b-dbc0b8444c6d"],
  "actions": [
    {
      "id": "25c94f7d-e9a4-5929-a7ff-5b8eac0b4b70",
      "scene_automation_id": "de2137ca-fa12-02da-213e-a3a9180cd9a2",
      "action_target": "ae9afb32-99a1-09e8-0a43-fce70fe1b068",
      "action_type": "30",
      "action_param_type": "",
      "action_param": "",
      "action_value": "",
      "remark": null
    },
    {
      "id": "21003bf9-b5a1-cf09-f0d4-9da051db9633",
      "scene_automation_id": "de2137ca-fa12-02da-213e-a3a9180cd9a2",
      "action_target": "4b0d5805-21e8-7be9-4a74-8372e92eb630",
      "action_type": "10",
      "action_param_type": "attributes",
      "action_param": "status",
      "action_value": "{\"status\":\"ok\"}",
      "remark": null
    }
  ]
}
```

Ý nghĩa:

- Action 1: action_type: "30" - Trigger alarm
- Action 2: action_type: "10" - Cập nhật attribute của device
  - action_target: Device cần cập nhật
  - action_param: "status"
  - action_value: {"status":"ok"}
  
Mục đích: Lưu trữ các hành động sẽ thực hiện khi điều kiện trigger được thỏa mãn.

### 3. Group Cache - automate:v3:one:_group_:1a56940d-67cc-79b6-aa2b-dbc0b8444c6d

```json
[{
  "id": "4cc2fd4e-6552-40a8-accd-95e215c678ca",
  "scene_automation_id": "de2137ca-fa12-02da-213e-a3a9180cd9a2",
  "enabled": "Y",
  "group_id": "1a56940d-67cc-79b6-aa2b-dbc0b8444c6d",
  "trigger_condition_type": "10",
  "trigger_source": "4b0d5805-21e8-7be9-4a74-8372e92eb630",
  "trigger_param_type": "telemetry",
  "trigger_param": "humidity",
  "trigger_operator": ">",
  "trigger_value": "20",
  "remark": null,
  "tenant_id": "d616bcbb"
}]
```

Ý nghĩa:

- Điều kiện trigger: Khi humidity > 20
- Nguồn trigger: Device 4b0d5805-21e8-7be9-4a74-8372e92eb630
- Loại data: telemetry (dữ liệu cảm biến)
- Trạng thái: enabled: "Y" (đang hoạt động)

### Luồng hoạt động của automation này:

- 1. Device gửi data: Device 4b0d5805-21e8-7be9-4a74-8372e92eb630 gửi telemetry data
- 2. Kiểm tra cache: Hệ thống tìm trong device cache xem device này có automation không
- 3. Kiểm tra điều kiện: So sánh humidity với giá trị 20
- 4. Thực hiện actions: Nếu humidity > 20:
  - Thực hiện action 1 (Trigger alarm)
  - Cập nhật status của device thành {"status":"ok"}

## Example alarm cache

>> KEYS "*alarm*"

1) "alarm_cach_device_v5_4b0d5805-21e8-7be9-4a74-8372e92eb630"
2) "alarm_cache_group_v5_b2e0bff2-cdf8-2268-6741-874058b48c6b"
3) "alarm_cach_scene_v5_de2137ca-fa12-02da-213e-a3a9180cd9a2"
4) "alarm_cach_alarm_v5_ba722102-d7ef-7550-f5e7-b8628e37f2fd"

### 1. Device Cache - alarm_cach_device_v5_4b0d5805-21e8-7be9-4a74-8372e92eb630

```json
["b2e0bff2-cdf8-2268-6741-874058b48c6b"]
```

Ý nghĩa:

- Device ID: 4b0d5805-21e8-7be9-4a74-8372e92eb630 (device này có alarm)
- Group ID: b2e0bff2-cdf8-2268-6741-874058b48c6b (group alarm liên quan)

Mục đích: Cho biết device nào có alarm nào, giúp tìm nhanh alarm khi device gửi data.

### 2. Group Cache - alarm_cache_group_v5_b2e0bff2-cdf8-2268-6741-874058b48c6b

```json
{
  "scene_automation_id": "de2137ca-fa12-02da-213e-a3a9180cd9a2",
  "alarm_config_id_list": ["ba722102-d7ef-7550-f5e7-b8628e37f2fd"],
  "alaram_device_id_list": ["4b0d5805-21e8-7be9-4a74-8372e92eb630"],
  "contents": ["Device(direct_sensor) Telemetry [humidity]: 37.01 > 20"]
}
```

Ý nghĩa:

- Scene Automation ID: de2137ca-fa12-02da-213e-a3a9180cd9a2 (scene automation liên quan)
- Alarm Config ID: ba722102-d7ef-7550-f5e7-b8628e37f2fd (cấu hình alarm)
- Device ID: 4b0d5805-21e8-7be9-4a74-8372e92eb630 (device bị alarm)
- Contents: Mô tả alarm - "Device(direct_sensor) Telemetry [humidity]: 37.01 > 20"

Mục đích: Lưu trữ thông tin chi tiết về alarm group, bao gồm điều kiện trigger và mô tả.

### 3. Scene Cache - alarm_cach_scene_v5_de2137ca-fa12-02da-213e-a3a9180cd9a2

```json
["b2e0bff2-cdf8-2268-6741-874058b48c6b"]
```

Ý nghĩa:

- Scene Automation ID: de2137ca-fa12-02da-213e-a3a9180cd9a2 (scene automation)
- Group ID: b2e0bff2-cdf8-2268-6741-874058b48c6b (group alarm thuộc scene này)

Mục đích: Cho biết scene automation nào có alarm nào.

### 4. Alarm Config Cache - alarm_cach_alarm_v5_ba722102-d7ef-7550-f5e7-b8628e37f2fd

```json
["b2e0bff2-cdf8-2268-6741-874058b48c6b"]
```

Ý nghĩa:

- Alarm Config ID: ba722102-d7ef-7550-f5e7-b8628e37f2fd (cấu hình alarm)
- Group ID: b2e0bff2-cdf8-2268-6741-874058b48c6b (group sử dụng cấu hình này)

Mục đích: Cho biết cấu hình alarm nào được sử dụng bởi group nào.

### Luồng hoạt động của alarm system:

Khi device gửi data:

- 1. Device: 4b0d5805-21e8-7be9-4a74-8372e92eb630 gửi telemetry data
- 2. Kiểm tra device cache: Tìm xem device này có alarm không
- 3. Tìm group: Từ device cache → group b2e0bff2-cdf8-2268-6741-874058b48c6b
- 4. Kiểm tra điều kiện: So sánh humidity: 37.01 với 20
- 5. Kích hoạt alarm: Vì 37.01 > 20 → tạo alarm
  
Khi có alarm:

- Lưu thông tin: Vào group cache với mô tả "Device(direct_sensor) Telemetry [humidity]: 37.01 > 20"
- Liên kết: Scene automation, alarm config, và device
- Xử lý: Có thể gửi notification, thực hiện action, v.v.

Cấu trúc cache key pattern:

- Device Cache: alarm_cach_device_v5_deviceId → [groupIds]
- Group Cache: alarm_cache_group_v5_groupId → {sceneAutomationId, alarmConfigIds, deviceIds, contents}
- Scene Cache: alarm_cach_scene_v5_sceneAutomationId → [groupIds]
- Alarm Config Cache: alarm_cach_alarm_v5_alarmConfigId → [groupIds]


