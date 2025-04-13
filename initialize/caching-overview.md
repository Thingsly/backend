# System Cache Overview

> Note: All cache retrieval methods are written in `initialize/redis_init.go`.

## Common Cache Methods

- Serialize a map or struct object to a JSON string and store it in Redis:
  - `initialize/redis_init.go/SetRedisForJsondata`
- Retrieve JSON from Redis and deserialize it into the specified object:
  - `initialize/redis_init.go/GetRedisForJsondata`

## Cache Data Overview

### Device Information

- Key: `deviceId`

#### Device ID Cache

- Key: `voucher`
- This cache has been removed from VerneMQ because it caused many difficult-to-handle issues.

### Data Script

**Explanation**: Data scripts and device configurations have a many-to-one relationship (however, a telemetry script can only have one configuration). The data script includes `device_config_id`.
**Cache Reason**: Telemetry data reporting needs to check `{deviceId}_telemetry_script_flag`. If not empty, it indicates that the script should be executed before inserting data into the database.
**Business Logic**: Data received from devices must be processed immediately upon reception to convert it to platform standards; data sent to devices must be converted to device standards before being sent.

- Key: `{scriptId}_script`
- Key: `{deviceId}_telemetry_script_flag`
  - Do not maintain this cache for additions or updates; simply delete it.
  - When changing device configurations, delete this flag.
  - Directly delete the flag when enabling or disabling scripts.
  - When no available telemetry script exists, the value is `""`; otherwise, it is `script_id`.
  - If `{deviceId}_telemetry_script_flag` does not exist, it indicates cache loss or deletion, requiring cache restoration.

### Device Configuration

- Key: `{deviceConfigId}_config`

### Automation

#### 1. Enabling and Stopping Scene Linkage Cache (For Single Device Telemetry Only)

##### Cache 1

When a device reports data, check if `auto_{device_id}` exists (method: `cache.GetCacheByDeviceId(deviceInfo.ID)`). If it does not exist, consider the cache lost or unavailable, and call the **restore cache method** (method: `cache.SetCacheByDeviceId(deviceId, groups, actionInfos)`).
Restore Cache Method Logic: `key: auto_{device_id}`

- Enable/disable operations must update the automation flag (**Enable Scene Linkage Method**).
- New additions are disabled.
- After editing, disable the cache and directly call the **Disable Method**.
- Delete calls should disable the cache.
- Key: `auto_{device_id}`
  - Contains the device's "AND" condition group IDs and the corresponding scene linkage IDs for the groups.
  - A scene linkage with a single device condition cannot have a single device condition.
  - If group A meets the condition, groups B and C do not need to be checked.
  - When disabling automation strategies, iterate through the cache to delete the `scene_automation_id` associated structures.

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

##### Cạche 2

- Key: `auto_group_{group_id}`

  - Restore Cache Method

  - Enable/disable automation strategies should also update this cache.

  - Stores the full data from device_trigger_condition where the condition type is 10 (single device) or 22 (time range) corresponding to group_id.

```json
[{data}, {data}]
```

##### Cache 3

- Key: auto_action_{scene_automation_id}

- Restore Cache Method

- Enable/disable automation strategies should also update this cache.

```json
{
  "group_id": ["", ""],
  "actions": [{data}, {data}]
}

```

##### Cache 4 - Alarm Related Cache in Scene Linkage

Key: auto_alarm_{group_id}

Value:

```json
{
  "scene_automation_id": "xxx",
  "alarm_config_id_list": ["xxx", "xxx"],
  "alarm_device_id_list": ["xxx"] // Saved only when triggered by device configuration
}
```

Note: A scene linkage may contain multiple alarms, so this stores a list of alarm configuration IDs.

When the condition involves a single device or multiple devices, the alarm history is updated by condition group (group_id) and the corresponding device list. For example:

Condition group 1 triggers alarm A, generating an alarm history entry with the group ID + device ID list.

Condition group 1 triggers alarm A again, no new history is generated.

Condition group 2 triggers alarm A, generating a new alarm history with the group ID + device ID list.

Condition group 1 no longer triggers alarm, generating an alarm recovery history with the group ID + device ID list.

When the condition is device configuration, the alarm history is updated by condition group and the individual device in that configuration.

Usage Instructions:

When the scene linkage is operating, the update method: Upon device reporting, first use cache1 to find the corresponding group_id, then use cache2 to check if the condition group is met.

If met, trigger an alarm and check if the group_id exists in the cache:

If exists and the condition is not device configuration, it means the alarm is already triggered, no update is needed.

If exists and the condition is device configuration, check alarm_device_id_list. If the device exists, do nothing. If not, add it and add alarm status records in alarm_history.

If not exists and the condition is not device configuration, insert into the cache and add alarm records.

If not exists and the condition is device configuration, insert into the cache and add alarm records for only the triggering device.

If not met, check if the group_id exists in the cache:

If exists and the condition is not device configuration, add alarm history for normal status and clear the cache.

If exists and the condition is device configuration, remove the device from alarm_device_id_list and add recovery alarm records.

If not exists, no action is needed.

When disabling scene linkage:

Delete the corresponding group_id cache.

When deleting scene linkage or alarm configuration:

Delete corresponding cache and alarm history records.

Restore Method: Restore from alarm_history table by alarm_config_id and alarm_device_id_list.

#### 2. Enabling and Stopping Scene Linkage Cache (For a Single Device Category)

Conditions include only "AND" with a single device condition or a single device + time range condition.

When a single device condition triggers an action for a single device, the device ID of the trigger condition must be passed into the action logic. This means the device being controlled in the action is the triggering device itself.

Cache 1
When a device reports, check if auto_devconfig_{device_config_id} exists. If it doesn't, consider the cache lost or unavailable and call the restore cache method. Restore Cache Method Logic:

Enable/disable operations must update the automation flag (Enable Scene Linkage Method).

New additions are disabled.

After editing, disable the cache and directly call the Disable Method.

Delete calls should disable the cache.

Key: auto_devconfig_{device_config_id}

Contains this device configuration's "AND" condition group IDs and the corresponding scene linkage IDs for the groups.

A scene linkage with a single device condition cannot have a single device condition.

If one group (e.g., A) satisfies the condition, groups B and C are not checked.

When disabling automation strategies, iterate through the cache to delete the scene_automation_id associated structures.

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

Cache 2
Key: auto_devconfig_group_{group_id}

Restore Cache Method

Enable/disable automation strategies should also update this cache.

Stores the full data from device_trigger_condition where condition type is 11 (single device) and 22 (time range) corresponding to group_id.

```json
[{data}, {data}]
```

Cache 3
Key: auto_devconfig_action_{scene_automation_id}

Restore Cache Method

Enable/disable automation strategies should also update this cache.

Stores group_id from device_trigger_condition and actions from action_info corresponding to scene_automation_id.

```json
{
  "group_id": ["group_id_1", "group_id_2", "group_id_3"],
  "actions": [{data}, {data}]
}
```

#### 3. Attribute Reporting and Event Reporting Automation Workflow

This is divided into two categories: Single Device and Single Device Category.

Single Device
Cache 1: Key auto_attr_event_{device_id}, other logic follows the same pattern as single device telemetry cache.

Cache 2, Cache 3: Reuse single device telemetry cache.

Single Device Category
Cache 1: Key auto_devconfig_attr_event_{device_id}, other logic follows the same pattern as single device category telemetry cache.

Cache 2, Cache 3: Reuse single device category telemetry cache.
