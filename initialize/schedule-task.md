# Scheduled Task Specification

## Cache

- Distributed locking using **Redis**

## Related Tables

- `one_time_tasks`: Table for one-time scheduled tasks  
- `periodic_tasks`: Table for recurring scheduled tasks

## Business Logic

Separate scheduled jobs are responsible for scanning each of the two tables: one-time and recurring tasks.

### One-Time Tasks

1. **Acquire a distributed lock**
2. **Fetch up to 100 tasks**  
   - The limit (default: 100) is configurable
3. **Filter tasks** where:
   - `execution_time` > current time
   - `enabled` = `Y` (enabled)
   - `executing_state` = `NEX` (not executed)
4. **Mark tasks as executed**
5. **Release the lock**
6. **Process each task:**
   - If `current_time - execution_time` > `expiration_time`:
     - Set `executing_state` to `EXP` (expired, not executed)
   - Else:
     - Invoke the **action execution method**
     - Set `executing_state` to `EXE` (executed)

### Periodic Tasks

1. **Acquire a distributed lock**
2. **Fetch all active periodic tasks**
   Tính toán thời gian tiếp theo: Sử dụng `GetSceneExecuteTime()` để tính thời gian thực thi tiếp theo dựa trên:

- HOUR: Thực thi mỗi giờ (VD: "30" = phút thứ 30 mỗi giờ)
- DAY: Thực thi mỗi ngày (VD: "15:30:00+07:00")
- WEEK: Thực thi theo tuần (VD: "1|15:30:00+07:00" = thứ 2 lúc 15:30)
- MONTH: Thực thi theo tháng (VD: "15T15:30:00+07:00" = ngày 15 lúc 15:30)
- CRON: Sử dụng cron expression
  
1. **Calculate next execution time**  
   - Update `execution_time` field  
   - (Use a GPT-generated algorithm)
2. **Release the lock**
3. **Reuse one-time task execution logic**  
   - Invoke the **action execution method**
   - Update `executing_state` accordingly

```
func AcquireLock(lockKey string, expiration time.Duration) bool {
    ok, err := global.REDIS.SetNX(context.Background(), lockKey, true, expiration).Result()
    return ok
}
```

SETNX (SET if Not eXists)
