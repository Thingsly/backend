# Biểu đồ cơ chế phát hiện trạng thái thiết bị IoT

## Tổng quan hệ thống

```mermaid
graph TB
    subgraph "IoT Devices"
        D1[Device 1]
        D2[Device 2]
        D3[Device N]
    end
    
    subgraph "Backend System"
        subgraph "Redis-based Detection"
            R1[Redis Server]
            R2[DeviceListener]
            R3[Expiration Events]
        end
        
        subgraph "MQTT-based Detection"
            M1[MQTT Broker]
            M2[StatusManager]
            M3[Status Messages]
        end
        
        subgraph "Database & Cache"
            DB[(Database)]
            CACHE[(Redis Cache)]
        end
    end
    
    subgraph "Frontend"
        F1[WebSocket Clients]
        F2[Status Dashboard]
    end
    
    D1 --> M1
    D2 --> M1
    D3 --> M1
    
    D1 --> R1
    D2 --> R1
    D3 --> R1
    
    M1 --> M2
    M2 --> M3
    M3 --> DB
    
    R1 --> R2
    R2 --> R3
    R3 --> DB
    
    DB --> CACHE
    CACHE --> F1
    F1 --> F2
```

## Chi tiết cơ chế Redis-based Detection

```mermaid
sequenceDiagram
    participant Device as IoT Device
    participant Backend as Backend System
    participant Redis as Redis Server
    participant Listener as DeviceListener
    participant DB as Database
    participant WS as WebSocket Clients
    
    Note over Device,WS: Redis-based Automatic Detection
    
    Device->>Backend: Gửi Telemetry Data
    Backend->>Backend: HeartbeatDeal() triggered
    Backend->>Redis: Set device:{deviceID}:heartbeat với TTL
    Note right of Redis: TTL = online_timeout (60s)
    
    loop Monitor Redis Expiration
        Redis->>Listener: Expiration Event
        Listener->>Listener: handleExpiredKey()
        Listener->>DB: DeviceOnline(deviceID, "0")
        Listener->>WS: Broadcast status change
        Note right of WS: Device offline detected
    end
```

## Chi tiết cơ chế MQTT-based Detection

```mermaid
sequenceDiagram
    participant Device as IoT Device
    participant Broker as MQTT Broker
    participant Manager as StatusManager
    participant DB as Database
    participant WS as WebSocket Clients
    
    Note over Device,WS: MQTT-based Manual Detection
    
    Device->>Broker: Publish status message
    Note right of Device: Topic: devices/status/{deviceID}<br/>Payload: "1" (online) / "0" (offline)
    
    Broker->>Manager: Forward message
    Manager->>Manager: MessageHandler()
    Manager->>DB: DeviceOnline(deviceID, status)
    Manager->>WS: Broadcast status change
    Note right of WS: Immediate status update
```

## So sánh hai cơ chế

```mermaid
graph LR
    subgraph "Redis-based Detection"
        R1[Automatic Timeout]
        R2[Expiration Events]
        R3[Reliable Detection]
        R4[After Timeout Period]
    end
    
    subgraph "MQTT-based Detection"
        M1[Manual Updates]
        M2[Status Messages]
        M3[Immediate Response]
        M4[Real-time Updates]
    end
    
    subgraph "Integration"
        I1[Dual Monitoring]
        I2[Redundancy]
        I3[Fallback System]
    end
    
    R1 --> I1
    M1 --> I1
    R2 --> I2
    M2 --> I2
    R3 --> I3
    M3 --> I3
```

## Workflow tổng hợp

```mermaid
flowchart TD
    A[Device Activity] --> B{Type of Activity}
    
    B -->|Telemetry Data| C[HeartbeatDeal]
    B -->|Status Message| D[StatusManager]
    
    C --> E[Set Redis Key with TTL]
    E --> F[DeviceListener Monitor]
    F --> G{Key Expired?}
    G -->|Yes| H[Device Offline]
    G -->|No| F
    
    D --> I[Process Status Message]
    I --> J[Update Database]
    J --> K[Broadcast Status]
    
    H --> L[Update Database]
    L --> M[Broadcast Status]
    
    K --> N[WebSocket Clients]
    M --> N
    
    N --> O[Frontend Dashboard]
    
    style A fill:#e1f5fe
    style H fill:#ffebee
    style K fill:#e8f5e8
    style M fill:#ffebee
    style O fill:#f3e5f5
```

## Cấu hình và tham số

```mermaid
graph TD
    subgraph "Redis Configuration"
        RC1[Database: 10]
        RC2[TTL: 60 seconds]
        RC3[Key Pattern: device:{deviceID}:heartbeat]
    end
    
    subgraph "MQTT Configuration"
        MC1[Broker: localhost:1883]
        MC2[Topic: devices/status/{deviceID}]
        MC3[QoS: 1]
    end
    
    subgraph "Device Configuration"
        DC1[online_timeout: 60]
        DC2[heartbeat: 30]
        DC3[status_report: manual]
    end
    
    subgraph "Performance Metrics"
        PM1[Memory Usage]
        PM2[CPU Usage]
        PM3[Network Latency]
        PM4[Detection Accuracy]
    end
    
    RC1 --> PM1
    MC1 --> PM3
    DC1 --> PM4
```

## Xử lý lỗi và fallback

```mermaid
flowchart TD
    A[System Startup] --> B[Initialize Both Systems]
    
    B --> C[Redis Connection]
    B --> D[MQTT Connection]
    
    C --> E{Redis Available?}
    D --> F{MQTT Available?}
    
    E -->|Yes| G[Enable Redis Detection]
    E -->|No| H[Disable Redis Detection]
    
    F -->|Yes| I[Enable MQTT Detection]
    F -->|No| J[Disable MQTT Detection]
    
    G --> K[Dual Monitoring Active]
    I --> K
    
    H --> L[Single System Mode]
    J --> L
    
    K --> M[High Reliability]
    L --> N[Reduced Reliability]
    
    style M fill:#e8f5e8
    style N fill:#fff3e0
```

## Kết luận

Hệ thống sử dụng **dual monitoring approach** với:

### **Redis-based Detection (Primary)**
- **Ưu điểm**: Tự động, đáng tin cậy, không phụ thuộc vào device
- **Nhược điểm**: Chỉ detect offline sau timeout
- **Use case**: Automatic offline detection

### **MQTT-based Detection (Secondary)**
- **Ưu điểm**: Real-time, immediate response
- **Nhược điểm**: Phụ thuộc vào device implementation
- **Use case**: Manual status updates, immediate changes

### **Integration Benefits**
- **Redundancy**: Cả hai system hoạt động song song
- **Reliability**: Fallback khi một system fails
- **Flexibility**: Hỗ trợ cả automatic và manual detection
- **Scalability**: Có thể scale từng system độc lập 