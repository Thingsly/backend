# Scheduled Task Description

## Cache

- Redis Distributed Lock

## Related Tables

- One-Time Task Table: `one_time_tasks`
- Periodic Task Table: `periodic_tasks`

## Business Flow

Separate scheduled tasks for one-time and periodic tasks to scan the tables.

### One-Time Tasks

1. Acquire Lock
2. Retrieve the first 100 tasks (the quantity of 100 is configurable)
3. Select tasks where `execution_time` is greater than the current time, `enabled` is 'Y' (enabled), and `executing_state` is 'NEX' (not executed)
4. Directly update tasks to executed status
5. Release Lock
6. Check if the difference between the current time and `execution_time` exceeds `expiration_time`:
   - If it does, update the execution status to `EXP` (expired, not executed)
   - If not, call the **execute action method**, then update the status to `EXE` (executed)
   - After execution, update the `executing_state`

### Periodic Tasks

1. Acquire Lock
2. Retrieve tasks
3. Calculate the next execution time (using GPT for the algorithm) and update the `execution_time` field
4. Release Lock
5. Similar to the one-time task process, the **execute action method** is reused
