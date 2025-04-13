# 【Feature Level】Device Online Status Management Overview Design

## 1. Overall Solution

The offline detection of specific devices is implemented based on Redis expiration notifications, with PostgreSQL maintaining the status.

## 2. Technical Architecture

- **Redis**: Expiration notification mechanism
- **PostgreSQL**: Device status storage
- **MQTT**: Device data collection

## 3. Core Design

### 3.1 Data Structure

**Redis:**

- `device:{deviceId}:heartbeat`  // Expiration key for heartbeat device

- `device:{deviceId}:timeout`    // Expiration key for timed-out device

**PostgreSQL - devices table:**

- `is_online`: boolean      // Device online status
- `heartbeat_time`: int     // Heartbeat time, null indicates not set
- `timeout`: int            // Timeout time, null indicates not set

### 3.2 Key Workflow

- **MQTT Message Handling:**
  1. Get device configuration type
  2. Heartbeat device:
     - Only update to online if it is offline
     - Set heartbeat expiration key
  3. Timeout device:
     - Set timeout expiration key

- **Redis Expiration Handling:**
  - When corresponding key expires → update device to offline status

### 3.3 Key Technical Points

- Device type differentiation handling
- Redis expiration event reliability
- PostgreSQL concurrency control
