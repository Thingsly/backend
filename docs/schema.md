# Schema Thingsly

Dưới đây là mô tả chi tiết về các bảng trong cơ sở dữ liệu Thingsly, bao gồm cấu trúc bảng, các trường dữ liệu, ràng buộc và mối quan hệ giữa các bảng.

## 📘 Table: `action_info`

Bảng `action_info` định nghĩa các hành động sẽ được thực hiện trong một **tự động hóa cảnh (scene automation)**. Mỗi hành động tương ứng với một rule trong hệ thống logic điều khiển thiết bị IoT.

---

### 🧩 Columns

| Column Name           | Type           | Nullable | Default | Description                                                               |
| --------------------- | -------------- | -------- | ------- | ------------------------------------------------------------------------- |
| `id`                  | `varchar(36)`  | No       | —       | Khóa chính (UUID) của hành động                                           |
| `scene_automation_id` | `varchar(36)`  | No       | —       | Khóa ngoại đến bảng `scene_automations.id`                                |
| `action_target`       | `varchar(255)` | Yes      | —       | Đối tượng tác động (ví dụ: device ID, group ID, hoặc cảnh khác)           |
| `action_type`         | `varchar(10)`  | No       | —       | Loại hành động (ví dụ: "write", "delay", "notification", "trigger_scene") |
| `action_param_type`   | `varchar(20)`  | Yes      | —       | Kiểu tham số hành động (dạng dữ liệu hoặc danh mục)                       |
| `action_param`        | `varchar(50)`  | Yes      | —       | Tham số hành động cụ thể (ví dụ: tên thuộc tính cần điều khiển)           |
| `action_value`        | `text`         | Yes      | —       | Giá trị truyền vào hành động (JSON, số, chuỗi, v.v.)                      |
| `remark`              | `varchar(255)` | Yes      | —       | Ghi chú tùy ý                                                             |

---

### 🔐 Constraints

- `PRIMARY KEY (id)`
- `FOREIGN KEY (scene_automation_id)` → `scene_automations(id)` (ON DELETE CASCADE)

---

### 🔁 Relationship

Mỗi bản ghi trong `action_info` **thuộc về một cảnh tự động hóa** (một rule logic), được định nghĩa ở bảng:

- 🔗 `scene_automations.id`

---

### 🧠 Ghi chú

- `action_type` là trường then chốt quyết định hành vi của rule:
  - `"write"` → ghi giá trị vào thiết bị
  - `"delay"` → tạm dừng thực hiện rule
  - `"notification"` → gửi cảnh báo
  - `"trigger_scene"` → gọi cảnh khác

- Cột `action_value` cho phép mở rộng linh hoạt, có thể chứa JSON nếu cần truyền nhiều giá trị cùng lúc.

---

📌 *Ví dụ minh họa*: Một rule trong bảng `scene_automations` có thể gồm nhiều hành động trong `action_info`, ví dụ:

```json
[
  {
    "action_type": "write",
    "action_target": "device-abc",
    "action_param": "relay",
    "action_value": "1"
  },
  {
    "action_type": "delay",
    "action_value": "5s"
  },
  {
    "action_type": "notification",
    "action_value": "Thiết bị đã bật relay."
  }
]
```

## 🚨 Table: `alarm_config`

Bảng `alarm_config` định nghĩa các **cấu hình cảnh báo** trong hệ thống, giúp xác định mức độ cảnh báo, nhóm nhận thông báo và các thông tin liên quan đến cảnh báo của thiết bị IoT.

---

### 🧩 Columns

| Column Name             | Type                       | Nullable | Default | Description                                          |
| ----------------------- | -------------------------- | -------- | ------- | ---------------------------------------------------- |
| `id`                    | `varchar(36)`              | No       | —       | UUID định danh duy nhất cho cấu hình cảnh báo        |
| `name`                  | `varchar(255)`             | No       | —       | Tên cảnh báo                                         |
| `description`           | `varchar(255)`             | Yes      | —       | Mô tả thêm về cấu hình cảnh báo                      |
| `alarm_level`           | `varchar(10)`              | No       | —       | Mức độ cảnh báo (ví dụ: info, warning, critical)     |
| `notification_group_id` | `varchar(36)`              | No       | —       | Khóa ngoại đến nhóm nhận thông báo                   |
| `created_at`            | `timestamp with time zone` | No       | —       | Thời điểm tạo cấu hình                               |
| `updated_at`            | `timestamp with time zone` | No       | —       | Thời điểm cập nhật cuối cùng                         |
| `tenant_id`             | `varchar(36)`              | No       | —       | Định danh tổ chức chủ quản                           |
| `remark`                | `varchar(255)`             | Yes      | —       | Ghi chú tùy ý                                        |
| `enabled`               | `varchar(10)`              | No       | —       | Trạng thái bật/tắt cảnh báo (ví dụ: "true", "false") |

---

### 🔐 Constraints

- `PRIMARY KEY (id)`

---

### 🔁 Relationships

- 🔗 **Được tham chiếu bởi**:
  - `alarm_info.alarm_config_id` → `alarm_config.id` (ON DELETE CASCADE)

---

### 🧠 Ghi chú

- Trường `alarm_level` giúp phân loại độ nghiêm trọng của cảnh báo để quyết định cách xử lý.
- `notification_group_id` là mối liên kết tới các user/group sẽ nhận thông báo khi cảnh báo được kích hoạt.
- Các cảnh báo thực tế (thường là theo thời gian thực) sẽ được lưu trong bảng `alarm_info`.

---

📌 *Ví dụ minh họa*:
Một bản ghi cấu hình cảnh báo:

```json
{
  "id": "cfg-001",
  "name": "High Temperature Alarm",
  "alarm_level": "critical",
  "notification_group_id": "group-01",
  "enabled": "true"
}
```

## 🗂️ Table: `alarm_history`

Bảng `alarm_history` lưu trữ **lịch sử các lần kích hoạt cảnh báo**, bao gồm thông tin chi tiết về cảnh báo, thiết bị liên quan, trạng thái, và thời gian kích hoạt. Đây là bảng cực kỳ quan trọng để tra cứu, phân tích sự cố hoặc tạo báo cáo sau khi sự kiện xảy ra.

---

### 🧩 Columns

| Column Name           | Type                       | Nullable | Default | Description                                                            |
| --------------------- | -------------------------- | -------- | ------- | ---------------------------------------------------------------------- |
| `id`                  | `varchar(36)`              | No       | —       | UUID định danh duy nhất cho bản ghi lịch sử                            |
| `alarm_config_id`     | `varchar(36)`              | No       | —       | Liên kết tới bảng `alarm_config` để xác định cấu hình cảnh báo         |
| `group_id`            | `varchar(36)`              | No       | —       | Nhóm liên quan tới cảnh báo (thường là nhóm người dùng nhận thông báo) |
| `scene_automation_id` | `varchar(36)`              | No       | —       | Mã tự động hóa cảnh báo đã kích hoạt                                   |
| `name`                | `varchar(255)`             | No       | —       | Tên cảnh báo (được ghi lại tại thời điểm xảy ra)                       |
| `description`         | `varchar(255)`             | Yes      | —       | Mô tả thêm cho bản ghi cảnh báo                                        |
| `content`             | `text`                     | Yes      | —       | Nội dung chi tiết về cảnh báo                                          |
| `alarm_status`        | `varchar(3)`               | No       | —       | Trạng thái cảnh báo (`ON`, `OFF`, v.v.)                                |
| `tenant_id`           | `varchar(36)`              | No       | —       | Tổ chức sở hữu cảnh báo này                                            |
| `remark`              | `varchar(255)`             | Yes      | —       | Ghi chú tùy ý                                                          |
| `create_at`           | `timestamp with time zone` | No       | —       | Thời gian cảnh báo được kích hoạt                                      |
| `alarm_device_list`   | `jsonb`                    | No       | —       | Danh sách các thiết bị liên quan đến cảnh báo (dưới dạng JSON)         |

---

### 🔐 Constraints

- `PRIMARY KEY (id)`

---

### 🔁 Relationships

- 🔗 `alarm_config_id` → `alarm_config(id)` _(không có khai báo FK nhưng logic là quan hệ trực tiếp)_
- 🔗 `scene_automation_id` → `scene_automations(id)`
- 🔗 `group_id` → Nhóm người dùng (có thể liên kết với bảng `notification_group` hoặc tương đương)

---

### 🧠 Ghi chú

- Trường `alarm_device_list` là một mảng JSON chứa chi tiết thiết bị (có thể bao gồm `device_id`, `name`, `value`, trạng thái,...).
- Bản ghi trong bảng này không bị xóa mà giữ lại để **audit**, **truy vết lỗi** hoặc **thống kê cảnh báo**.
- `alarm_status` thường sẽ là `"ON"` khi bắt đầu cảnh báo, `"OFF"` khi kết thúc.

---

📌 *Ví dụ minh họa `alarm_device_list`*:
```json
[
  {
    "device_id": "dev-001",
    "device_name": "TempSensor-A1",
    "value": 78.6,
    "unit": "°C"
  },
  {
    "device_id": "dev-002",
    "device_name": "SmokeDetector-B2",
    "status": "triggered"
  }
]
```

## 🚨 Table: `alarm_info`

Bảng `alarm_info` lưu trữ **thông tin chi tiết về các cảnh báo đang hoạt động hoặc đã xảy ra**, gắn với cấu hình cảnh báo cụ thể và hỗ trợ quá trình xử lý, giám sát các tình huống bất thường.

---

### 🧩 Columns

| Column Name         | Type                       | Nullable | Default | Description                                                      |
| ------------------- | -------------------------- | -------- | ------- | ---------------------------------------------------------------- |
| `id`                | `varchar(36)`              | No       | —       | UUID định danh cảnh báo                                          |
| `alarm_config_id`   | `varchar(36)`              | No       | —       | Liên kết đến bảng `alarm_config` (cấu hình cảnh báo)             |
| `name`              | `varchar(255)`             | No       | —       | Tên của cảnh báo                                                 |
| `alarm_time`        | `timestamp with time zone` | No       | —       | Thời điểm cảnh báo được kích hoạt                                |
| `description`       | `varchar(255)`             | Yes      | —       | Mô tả ngắn gọn về cảnh báo                                       |
| `content`           | `text`                     | Yes      | —       | Nội dung chi tiết về cảnh báo                                    |
| `processor`         | `varchar(36)`              | Yes      | —       | ID của người hoặc hệ thống xử lý cảnh báo                        |
| `processing_result` | `varchar(10)`              | No       | —       | Kết quả xử lý cảnh báo (`solved`, `ignored`, `escalated`,...)    |
| `tenant_id`         | `varchar(36)`              | No       | —       | Thuộc về tenant (tổ chức hoặc khách hàng) nào                    |
| `remark`            | `varchar(255)`             | Yes      | —       | Ghi chú tùy ý                                                    |
| `alarm_level`       | `varchar(10)`              | Yes      | —       | Mức độ nghiêm trọng của cảnh báo (`low`, `medium`, `high`, v.v.) |

---

### 🔐 Constraints

- `PRIMARY KEY (id)`
- 🔗 `FOREIGN KEY (alarm_config_id)` → `alarm_config(id)` _(ON DELETE CASCADE)_

---

### 🔁 Relationships

- Mỗi cảnh báo gắn với một `alarm_config` (cấu hình cảnh báo định nghĩa ngưỡng và điều kiện).
- Có thể có liên kết logic tới người dùng hoặc thiết bị thông qua `processor`.

---

### 📌 Ghi chú sử dụng

- Bảng này thường lưu cảnh báo đang hoạt động hoặc đã ghi nhận, hỗ trợ theo dõi lịch sử xử lý.
- Trường `processing_result` là chìa khóa để lọc trạng thái cảnh báo: đã được xử lý hay chưa.
- Có thể kết hợp với `alarm_history` để lưu nhật ký chi tiết cho mỗi cảnh báo.

---

### 🧠 Ví dụ

```json
{
  "id": "alarm-001",
  "alarm_config_id": "config-123",
  "name": "High Temperature",
  "alarm_time": "2025-04-14T10:20:00+07:00",
  "description": "Temperature exceeded 85°C",
  "content": "Device A1 reported 87.3°C at 10:20 AM",
  "processor": "user-567",
  "processing_result": "solved",
  "tenant_id": "tenant-abc",
  "alarm_level": "high"
}
```

## 📦 Table: `attribute_datas`

Bảng `attribute_datas` lưu trữ **các thuộc tính thời gian thực của thiết bị** (real-time attributes). Đây là nơi ghi lại giá trị thuộc tính mới nhất của từng thiết bị, dùng cho việc hiển thị trạng thái thiết bị hoặc xử lý logic tự động.

---

### 🧩 Columns

| Column Name | Type                       | Nullable | Default | Description                                                         |
| ----------- | -------------------------- | -------- | ------- | ------------------------------------------------------------------- |
| `id`        | `varchar(36)`              | No       | —       | ID định danh của bản ghi thuộc tính                                 |
| `device_id` | `varchar(36)`              | No       | —       | Liên kết đến thiết bị trong bảng `devices`                          |
| `key`       | `varchar(255)`             | No       | —       | Tên thuộc tính (VD: `temperature`, `status`, `firmwareVersion`,...) |
| `ts`        | `timestamp with time zone` | No       | —       | Dấu thời gian cập nhật giá trị thuộc tính                           |
| `bool_v`    | `boolean`                  | Yes      | —       | Giá trị kiểu boolean (nếu applicable)                               |
| `number_v`  | `double precision`         | Yes      | —       | Giá trị kiểu số (VD: nhiệt độ, điện áp...)                          |
| `string_v`  | `text`                     | Yes      | —       | Giá trị kiểu chuỗi (VD: trạng thái, firmware name,...)              |
| `tenant_id` | `varchar(36)`              | Yes      | —       | Thuộc về tenant nào (tổ chức/khách hàng)                            |

---

### 🔐 Constraints

- `UNIQUE(device_id, key)` — Mỗi thiết bị chỉ có **1 giá trị cuối cùng** cho mỗi `key` (attribute).
- 🔗 `FOREIGN KEY(device_id)` → `devices(id)` _(ON DELETE RESTRICT)_

---

### 🔁 Relationships

- Liên kết chặt chẽ với bảng `devices`: mỗi bản ghi thuộc về 1 thiết bị cụ thể.
- Cập nhật liên tục thông qua telemetry hoặc API từ thiết bị.

---

### 🔎 Cách sử dụng

- Dùng để hiển thị **trạng thái hiện tại** của thiết bị trên giao diện.
- Là nguồn dữ liệu để kiểm tra điều kiện trong automation hoặc cảnh báo (`if temperature > 50°C`).
- Có thể được cập nhật từ **telemetry** hoặc khi người dùng **set attribute** qua API.

---

### 🧠 Ví dụ

```json
{
  "id": "attr-001",
  "device_id": "dev-001",
  "key": "temperature",
  "ts": "2025-04-14T08:30:00+07:00",
  "number_v": 72.5
}
```

## 🧾 Table: `attribute_set_logs`

Bảng `attribute_set_logs` lưu lại **log của các thao tác set attribute** cho thiết bị, bao gồm dữ liệu gửi đi, phản hồi từ thiết bị, kết quả thực hiện và thông tin liên quan đến người dùng và thời điểm thực hiện.

---

### 🧩 Columns

| Column Name      | Type                       | Nullable | Description                                         |
| ---------------- | -------------------------- | -------- | --------------------------------------------------- |
| `id`             | `varchar(36)`              | No       | Mã định danh log                                    |
| `device_id`      | `varchar(36)`              | No       | Thiết bị được cập nhật attribute                    |
| `operation_type` | `varchar(255)`             | Yes      | Loại thao tác (set attribute / auto update /...)    |
| `message_id`     | `varchar(36)`              | Yes      | ID của message hoặc request liên quan               |
| `data`           | `text`                     | Yes      | Payload dữ liệu gửi tới thiết bị (JSON text format) |
| `rsp_data`       | `text`                     | Yes      | Phản hồi từ thiết bị (nếu có)                       |
| `status`         | `varchar(2)`               | Yes      | Trạng thái thao tác (`OK`, `NG`, `ER`...)           |
| `error_message`  | `varchar(500)`             | Yes      | Thông tin lỗi (nếu thao tác thất bại)               |
| `created_at`     | `timestamp with time zone` | No       | Thời điểm thao tác được ghi nhận                    |
| `user_id`        | `varchar(36)`              | Yes      | Người dùng thực hiện thao tác (nếu có)              |
| `description`    | `varchar(255)`             | Yes      | Ghi chú mô tả thao tác                              |

---

### 🔐 Constraints

- `PRIMARY KEY(id)`
- 🔗 `FOREIGN KEY(device_id)` → `devices(id)` _(ON DELETE CASCADE)_

---

### 🔎 Cách sử dụng

- Ghi lại lịch sử **đặt lại giá trị thuộc tính** của thiết bị.
- Phục vụ cho việc **kiểm tra lỗi**, **kiểm toán**, hoặc **debug** hệ thống.
- Dùng trong giao diện quản lý để xem các thay đổi attribute thủ công hoặc từ hệ thống tự động hóa.

---

### 🧠 Ví dụ

```json
{
  "id": "log-001",
  "device_id": "dev-001",
  "operation_type": "manual_set",
  "message_id": "msg-999",
  "data": "{\"power\": true}",
  "rsp_data": "{\"result\": \"ok\"}",
  "status": "OK",
  "created_at": "2025-04-14T08:45:00+07:00",
  "user_id": "user-001"
}
```

## 🧾 Table: `boards`

Bảng `boards` đại diện cho **giao diện bảng điều khiển (dashboard)** hoặc **bảng hiển thị dữ liệu** được cấu hình cho từng tenant. Mỗi bảng có thể chứa thông tin cấu hình widget, bố cục và metadata khác.

---

### 🧩 Columns

| Column Name   | Type                       | Nullable | Description                                                   |
| ------------- | -------------------------- | -------- | ------------------------------------------------------------- |
| `id`          | `varchar(36)`              | No       | Mã định danh của bảng                                         |
| `name`        | `varchar(255)`             | No       | Tên hiển thị của bảng                                         |
| `config`      | `json`                     | Yes      | Cấu hình JSON cho bảng (các widget, vị trí, kiểu hiển thị...) |
| `tenant_id`   | `varchar(36)`              | No       | Tenant sở hữu bảng này                                        |
| `created_at`  | `timestamp with time zone` | No       | Ngày tạo                                                      |
| `updated_at`  | `timestamp with time zone` | No       | Ngày cập nhật gần nhất                                        |
| `home_flag`   | `varchar(2)`               | No       | Có phải bảng chính (home) không (`"Y"`/`"N"`)                 |
| `description` | `varchar(500)`             | Yes      | Mô tả bảng                                                    |
| `remark`      | `varchar(255)`             | Yes      | Ghi chú bổ sung                                               |
| `menu_flag`   | `varchar(2)`               | Yes      | Có hiển thị trong menu không (`"Y"`/`"N"`)                    |

---

### 🔐 Constraints

- `PRIMARY KEY(id)`

---

### 🧠 Ghi chú

- Trường `config` lưu thông tin cấu hình chi tiết các **widget**, **filter**, **layout**, thường được trình bày dưới dạng JSON như:

```json
{
  "layout": "grid",
  "widgets": [
    {
      "type": "chart",
      "title": "Temperature",
      "dataSource": "device_001.attribute.temperature"
    }
  ]
}
```

## 🔐 Table: `casbin_rule`

Bảng `casbin_rule` lưu trữ các quy tắc phân quyền truy cập dựa trên mô hình [Casbin RBAC/ABAC](https://casbin.org/docs/en/supported-models). Đây là trung tâm của hệ thống kiểm soát truy cập.

---

### 📦 Columns

| Column Name | Type           | Nullable | Description                                                       |
| ----------- | -------------- | -------- | ----------------------------------------------------------------- |
| `id`        | `bigint`       | No       | Khóa chính tự tăng                                                |
| `ptype`     | `varchar(100)` | Yes      | Loại chính sách (ví dụ: `p` - policy, `g` - grouping)             |
| `v0`        | `varchar(100)` | Yes      | Chủ thể (subject) - thường là user/role                           |
| `v1`        | `varchar(100)` | Yes      | Đối tượng (object) - thường là resource, route, hoặc API endpoint |
| `v2`        | `varchar(100)` | Yes      | Hành động (action) - ví dụ: `read`, `write`                       |
| `v3`        | `varchar(100)` | Yes      | (Tùy chọn) - thường dùng trong ABAC để mở rộng logic              |
| `v4`        | `varchar(100)` | Yes      | (Tùy chọn)                                                        |
| `v5`        | `varchar(100)` | Yes      | (Tùy chọn)                                                        |

---

### 🔐 Constraints

- `PRIMARY KEY(id)`
- `UNIQUE(ptype, v0, v1, v2, v3, v4, v5)` — đảm bảo không có quy tắc trùng lặp.

---

### 🧠 Ghi chú

- Một policy thông thường có thể trông như:
  - `ptype = "p", v0 = "admin", v1 = "/devices", v2 = "GET"` — nghĩa là `admin` có quyền `GET` trên `/devices`.
- Grouping policy:
  - `ptype = "g", v0 = "user1", v1 = "admin"` — `user1` thuộc nhóm `admin`.
- Các cột `v3 ~ v5` được dùng khi bạn mở rộng mô hình, ví dụ thêm điều kiện theo `tenant_id`, `domain`, hoặc thuộc tính của đối tượng.

---

### 🔄 Tích hợp

Casbin thường được sử dụng trực tiếp trong code như:

```go
e.AddPolicy("admin", "/devices", "GET")
e.Enforce("admin", "/devices", "GET") // => true or false
```

## 📝 Table: `command_set_logs`

Bảng `command_set_logs` lưu trữ các bản ghi về lệnh được gửi đến các thiết bị, bao gồm thông tin về loại lệnh, trạng thái, và các dữ liệu phản hồi từ thiết bị.

---

### 📦 Columns

| Column Name      | Type                          | Nullable | Description                                                  |
| ---------------- | ----------------------------- | -------- | ------------------------------------------------------------ |
| `id`             | `character varying(36)`       | No       | Khóa chính tự tăng                                           |
| `device_id`      | `character varying(36)`       | No       | ID của thiết bị nhận lệnh (liên kết với bảng `devices`)      |
| `operation_type` | `character varying(255)`      | Yes      | Loại của lệnh (ví dụ: "SET", "RESET", "CONFIG")              |
| `message_id`     | `character varying(36)`       | Yes      | ID của thông điệp (nếu có)                                   |
| `data`           | `text`                        | Yes      | Dữ liệu gửi đến thiết bị (thông qua lệnh)                    |
| `rsp_data`       | `text`                        | Yes      | Dữ liệu phản hồi từ thiết bị                                 |
| `status`         | `character varying(2)`        | Yes      | Trạng thái lệnh (ví dụ: "OK", "FAIL")                        |
| `error_message`  | `character varying(500)`      | Yes      | Thông điệp lỗi (nếu có)                                      |
| `created_at`     | `timestamp(6) with time zone` | No       | Thời gian tạo bản ghi                                        |
| `user_id`        | `character varying(36)`       | Yes      | ID người dùng (nếu có liên kết với hành động này)            |
| `description`    | `character varying(255)`      | Yes      | Mô tả ngắn về lệnh hoặc trạng thái                           |
| `identify`       | `character varying(255)`      | Yes      | Dùng để nhận diện lệnh hoặc thiết bị (thêm chi tiết về lệnh) |

---

### 🔐 Constraints

- `PRIMARY KEY(id)` — Khóa chính của bảng.
- `FOREIGN KEY(device_id) REFERENCES devices(id) ON DELETE CASCADE` — Liên kết đến bảng `devices`. Nếu thiết bị bị xóa, các bản ghi lệnh liên quan cũng sẽ bị xóa.

---

### 🧠 Ghi chú

- Dữ liệu trong bảng này được sử dụng để theo dõi các lệnh được gửi và thực thi trên các thiết bị IoT.
- Các cột như `data`, `rsp_data`, `status`, và `error_message` giúp theo dõi quá trình giao tiếp và xử lý lệnh, giúp phát hiện lỗi và xử lý phản hồi hiệu quả.
- Cột `user_id` có thể liên kết đến bảng người dùng, nếu bạn muốn theo dõi ai đã gửi lệnh hoặc thao tác với thiết bị.
- Cột `identify` có thể được sử dụng để lưu trữ thông tin nhận diện thêm cho lệnh, có thể là mã nhận diện hoặc thông tin bổ sung từ thiết bị.

---

### 🔄 Tích hợp

Có thể sử dụng bảng này để kiểm tra và debug quá trình giao tiếp giữa hệ thống và các thiết bị, chẳng hạn như xác định xem một lệnh đã được thực thi thành công hay không, và nếu có lỗi, xác định lỗi đó là gì.

```go
db.Query("INSERT INTO command_set_logs (device_id, operation_type, data, rsp_data) VALUES (?, ?, ?, ?)", deviceID, operationType, requestData, responseData)
```

## 📝 Table: `data_policy`

Bảng `data_policy` lưu trữ các chính sách dữ liệu, bao gồm các thông tin liên quan đến loại dữ liệu, thời gian giữ dữ liệu, và thông tin về các lần dọn dẹp dữ liệu.

---

### 📦 Columns

| Column Name              | Type                          | Nullable | Description                                                                                  |
| ------------------------ | ----------------------------- | -------- | -------------------------------------------------------------------------------------------- |
| `id`                     | `character varying(36)`       | No       | Khóa chính của bảng.                                                                         |
| `data_type`              | `character varying(1)`        | No       | Loại dữ liệu (ví dụ: "A" cho dữ liệu loại A, "B" cho dữ liệu loại B).                        |
| `retention_days`         | `integer`                     | No       | Số ngày lưu trữ dữ liệu trước khi thực hiện dọn dẹp.                                         |
| `last_cleanup_time`      | `timestamp(6) with time zone` | Yes      | Thời gian thực hiện dọn dẹp dữ liệu gần nhất.                                                |
| `last_cleanup_data_time` | `timestamp(6) with time zone` | Yes      | Thời gian dọn dẹp dữ liệu cuối cùng, có thể khác với `last_cleanup_time`.                    |
| `enabled`                | `character varying(1)`        | No       | Trạng thái kích hoạt chính sách dữ liệu (ví dụ: "Y" cho kích hoạt, "N" cho không kích hoạt). |
| `remark`                 | `character varying(255)`      | Yes      | Ghi chú hoặc mô tả về chính sách dữ liệu.                                                    |

---

### 🔐 Constraints

- `PRIMARY KEY(id)` — Khóa chính của bảng.

---

### 🧠 Ghi chú

- Bảng này được sử dụng để quản lý các chính sách lưu trữ và dọn dẹp dữ liệu, ví dụ, các dữ liệu có thể được dọn dẹp sau một số ngày nhất định.
- Cột `data_type` có thể xác định loại dữ liệu mà chính sách này áp dụng, có thể là các loại dữ liệu khác nhau được phân loại.
- Cột `retention_days` quyết định thời gian lưu trữ dữ liệu trước khi dữ liệu được dọn dẹp hoặc xóa đi.
- `last_cleanup_time` và `last_cleanup_data_time` cho phép theo dõi các lần dọn dẹp dữ liệu đã thực hiện.

---

### 🔄 Tích hợp

Có thể sử dụng bảng này để lên lịch các tác vụ dọn dẹp dữ liệu tự động, giúp tiết kiệm không gian lưu trữ và đảm bảo rằng dữ liệu không cần thiết được loại bỏ kịp thời.

Ví dụ trong Go:

```go
db.Query("INSERT INTO data_policy (id, data_type, retention_days, enabled) VALUES (?, ?, ?, ?)", id, dataType, retentionDays, enabled)
```

## 📝 Table: `data_scripts`

Bảng `data_scripts` lưu trữ thông tin về các script dữ liệu liên quan đến các cấu hình thiết bị, bao gồm tên, nội dung script, trạng thái kích hoạt, và các thông tin mô tả khác.

---

### 📦 Columns

| Column Name         | Type                          | Nullable | Description                                             |
| ------------------- | ----------------------------- | -------- | ------------------------------------------------------- |
| `id`                | `character varying(36)`       | No       | Khóa chính của bảng, duy nhất cho mỗi bản ghi script.   |
| `name`              | `character varying(99)`       | No       | Tên của script dữ liệu.                                 |
| `device_config_id`  | `character varying(36)`       | No       | Tham chiếu đến cấu hình thiết bị mà script này áp dụng. |
| `enable_flag`       | `character varying(9)`        | No       | Cờ kích hoạt (ví dụ: "enabled", "disabled").            |
| `content`           | `text`                        | Yes      | Nội dung chi tiết của script.                           |
| `script_type`       | `character varying(9)`        | No       | Loại script (ví dụ: "lua", "python", ...).              |
| `last_analog_input` | `text`                        | Yes      | Thông tin về đầu vào analog cuối cùng được xử lý.       |
| `description`       | `character varying(255)`      | Yes      | Mô tả ngắn gọn về chức năng của script.                 |
| `created_at`        | `timestamp(6) with time zone` | Yes      | Thời gian tạo script.                                   |
| `updated_at`        | `timestamp(6) with time zone` | Yes      | Thời gian cập nhật script.                              |
| `remark`            | `character varying(255)`      | Yes      | Ghi chú hoặc chú thích về script.                       |

---

### 🔐 Constraints

- `PRIMARY KEY(id)` — Khóa chính của bảng.

### 🔗 Foreign Key

- `FOREIGN KEY (device_config_id) REFERENCES device_configs(id) ON DELETE CASCADE` — Liên kết đến bảng `device_configs`. Nếu một cấu hình thiết bị bị xóa, các script liên quan cũng sẽ bị xóa.

---

### 🧠 Ghi chú

- **Cột `enable_flag`**: Xác định xem script có đang được kích hoạt hay không. Trạng thái có thể là "enabled" hoặc "disabled".
- **Cột `content`**: Lưu trữ nội dung chi tiết của script. Cột này có thể chứa mã nguồn của script dưới dạng văn bản.
- **Cột `script_type`**: Xác định loại script, có thể là các loại như Lua, Python, v.v.
- **Cột `last_analog_input`**: Có thể lưu trữ dữ liệu về lần xử lý đầu vào analog cuối cùng trong script.

---

### 🔄 Tích hợp

Bảng `data_scripts` có thể được sử dụng để lưu trữ các script tùy chỉnh cho từng cấu hình thiết bị, giúp hệ thống có thể thực thi các tác vụ tự động liên quan đến thiết bị.

Ví dụ trong Go:

```go
db.Query("INSERT INTO data_scripts (id, name, device_config_id, enable_flag, script_type) VALUES (?, ?, ?, ?, ?)", id, name, deviceConfigID, enableFlag, scriptType)
```

## 📝 Table: `device_configs`

Bảng `device_configs` lưu trữ thông tin về cấu hình thiết bị, bao gồm các chi tiết về loại thiết bị, cấu hình giao thức, thông tin bổ sung và các thông số kỹ thuật khác.

---

### 📦 Columns

| Column Name          | Type                          | Nullable | Description                                                       |
| -------------------- | ----------------------------- | -------- | ----------------------------------------------------------------- |
| `id`                 | `character varying(36)`       | No       | Khóa chính của bảng, duy nhất cho mỗi cấu hình thiết bị.          |
| `name`               | `character varying(99)`       | No       | Tên của cấu hình thiết bị.                                        |
| `device_template_id` | `character varying(36)`       | Yes      | Tham chiếu đến mẫu thiết bị.                                      |
| `device_type`        | `character varying(9)`        | No       | Loại thiết bị (ví dụ: "sensor", "actuator", ...).                 |
| `protocol_type`      | `character varying(36)`       | Yes      | Loại giao thức sử dụng cho thiết bị (ví dụ: "MQTT", "HTTP", ...). |
| `voucher_type`       | `character varying(36)`       | Yes      | Loại mã kích hoạt (nếu có).                                       |
| `protocol_config`    | `json`                        | Yes      | Cấu hình giao thức thiết bị dưới dạng JSON.                       |
| `device_conn_type`   | `character varying(36)`       | Yes      | Loại kết nối của thiết bị (ví dụ: "WIFI", "Ethernet", ...).       |
| `additional_info`    | `json`                        | Yes      | Thông tin bổ sung về thiết bị dưới dạng JSON.                     |
| `description`        | `character varying(255)`      | Yes      | Mô tả chi tiết về cấu hình thiết bị.                              |
| `tenant_id`          | `character varying(36)`       | No       | ID của thuê bao sử dụng cấu hình thiết bị.                        |
| `created_at`         | `timestamp(6) with time zone` | No       | Thời gian tạo cấu hình thiết bị.                                  |
| `updated_at`         | `timestamp(6) with time zone` | No       | Thời gian cập nhật cấu hình thiết bị.                             |
| `remark`             | `character varying(255)`      | Yes      | Ghi chú về cấu hình thiết bị.                                     |
| `other_config`       | `json`                        | Yes      | Cấu hình khác dưới dạng JSON (tùy chỉnh).                         |

---

### 🔐 Constraints

- `PRIMARY KEY(id)` — Khóa chính của bảng.

### 🔗 Foreign Key

- `FOREIGN KEY (device_template_id) REFERENCES device_templates(id) ON DELETE RESTRICT` — Liên kết đến bảng `device_templates`. Nếu một mẫu thiết bị bị xóa, các cấu hình liên quan không bị ảnh hưởng nhưng không thể xóa mẫu thiết bị khi có cấu hình liên kết.
- `FOREIGN KEY (device_config_id) REFERENCES device_configs(id) ON DELETE CASCADE` — Các bảng như `data_scripts`, `devices`, `products` đều tham chiếu đến `device_configs`. Nếu cấu hình thiết bị bị xóa, các bản ghi trong các bảng này sẽ tự động bị xóa.

---

### 🧠 Ghi chú

- **Cột `device_type`**: Xác định loại thiết bị, có thể là các loại như "sensor" (cảm biến), "actuator" (thiết bị điều khiển), v.v.
- **Cột `protocol_config`**: Lưu trữ các cấu hình giao thức, giúp thiết bị có thể giao tiếp với các hệ thống khác.
- **Cột `additional_info` và `other_config`**: Các cột này có thể chứa các thông tin tùy chỉnh và cấu hình bổ sung cho thiết bị, với định dạng JSON linh hoạt.

---

### 🔄 Tích hợp

Bảng `device_configs` đóng vai trò quan trọng trong việc quản lý và cấu hình các thiết bị trong hệ thống. Nó liên kết với nhiều bảng khác như `devices`, `products` và `data_scripts`, giúp các thiết bị có thể được cấu hình và vận hành một cách hiệu quả.

Ví dụ trong Go:

```go
db.Query("INSERT INTO device_configs (id, name, device_type, protocol_type) VALUES (?, ?, ?, ?)", id, name, deviceType, protocolType)
```



## 📘 Table: `devices`

Bảng `devices` lưu thông tin các thiết bị IoT trong hệ thống Thingsly. Đây là một bảng trung tâm, được liên kết với nhiều bảng khác như telemetry, OTA upgrade, logs, và cấu hình.

### 🧩 Columns

| Column Name         | Type                       | Nullable | Default | Description                             |
| ------------------- | -------------------------- | -------- | ------- | --------------------------------------- |
| `id`                | `varchar(36)`              | No       | —       | Khóa chính (UUID của thiết bị)          |
| `name`              | `varchar(255)`             | Yes      | —       | Tên thiết bị                            |
| `voucher`           | `varchar(500)`             | No       | `''`    | Mã voucher duy nhất                     |
| `tenant_id`         | `varchar(36)`              | No       | `''`    | Thuộc tenant nào                        |
| `is_enabled`        | `varchar(36)`              | No       | `''`    | Trạng thái kích hoạt                    |
| `activate_flag`     | `varchar(36)`              | No       | `''`    | Đã kích hoạt chưa                       |
| `created_at`        | `timestamp with time zone` | Yes      | —       | Ngày tạo                                |
| `update_at`         | `timestamp with time zone` | Yes      | —       | Ngày cập nhật                           |
| `device_number`     | `varchar(36)`              | No       | `''`    | Mã số thiết bị (unique)                 |
| `product_id`        | `varchar(36)`              | Yes      | —       | FK → `products.id`                      |
| `parent_id`         | `varchar(36)`              | Yes      | —       | FK đến thiết bị cha (nếu là sub-device) |
| `protocol`          | `varchar(36)`              | Yes      | —       | Loại giao thức                          |
| `label`             | `varchar(255)`             | Yes      | —       | Nhãn/tên gọi khác                       |
| `location`          | `varchar(100)`             | Yes      | —       | Vị trí vật lý                           |
| `sub_device_addr`   | `varchar(36)`              | Yes      | —       | Địa chỉ thiết bị con                    |
| `current_version`   | `varchar(36)`              | Yes      | —       | Phiên bản hiện tại                      |
| `additional_info`   | `json`                     | Yes      | `'{}'`  | Thông tin mở rộng (JSON)                |
| `protocol_config`   | `json`                     | Yes      | `'{}'`  | Cấu hình giao thức (JSON)               |
| `remark1`           | `varchar(255)`             | Yes      | —       | Ghi chú mở rộng 1                       |
| `remark2`           | `varchar(255)`             | Yes      | —       | Ghi chú mở rộng 2                       |
| `remark3`           | `varchar(255)`             | Yes      | —       | Ghi chú mở rộng 3                       |
| `device_config_id`  | `varchar(36)`              | Yes      | —       | FK → `device_configs.id`                |
| `batch_number`      | `varchar(500)`             | Yes      | —       | Số lô sản xuất                          |
| `activate_at`       | `timestamp with time zone` | Yes      | —       | Ngày kích hoạt                          |
| `is_online`         | `smallint`                 | No       | `0`     | Trạng thái online: 0/1                  |
| `access_way`        | `varchar(10)`              | Yes      | —       | Phương thức truy cập                    |
| `description`       | `varchar(500)`             | Yes      | —       | Mô tả thiết bị                          |
| `service_access_id` | `varchar(36)`              | Yes      | —       | FK → `service_access.id`                |

---

## 🔗 Indexes & Constraints

- `PRIMARY KEY (id)`
- `UNIQUE (device_number)`
- `UNIQUE (voucher)`
- `btree index`: `lower(device_number)`, `lower(name)`

---

## 🔐 Foreign Keys

| Column              | References           | On Delete  |
| ------------------- | -------------------- | ---------- |
| `product_id`        | `products(id)`       | `RESTRICT` |
| `device_config_id`  | `device_configs(id)` | `RESTRICT` |
| `service_access_id` | `service_access(id)` | `RESTRICT` |

---

## 🔁 Referenced By (Liên kết ngược)

Các bảng khác có khóa ngoại tham chiếu đến `devices.id`:

- `attribute_datas(device_id)`
- `attribute_set_logs(device_id)`
- `command_set_logs(device_id)`
- `event_datas(device_id)`
- `expected_datas(device_id)`
- `r_group_device(device_id)`
- `ota_upgrade_task_details(device_id)`
- `telemetry_set_logs(device_id)`

