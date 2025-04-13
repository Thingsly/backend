# Logging Output Specifications

## Purpose

Logging output specifications are a key part of ensuring software quality, improving development efficiency, and maintaining performance.

1. **Rapid Problem Diagnosis**: Enhance fault diagnosis efficiency.
2. **Improve Code Maintainability**: Facilitate understanding and maintenance by team members.
3. **Ensure Consistency and Comparability**: A unified format facilitates analysis and monitoring.
4. **Meet Security and Compliance Requirements**: Record key operations and data access.
5. **Monitor Performance and Fault Recovery**: Assist with performance monitoring and system recovery.
6. **Facilitate Data Analysis**: Helps in generating business and technical insights.

## Explanation of Log Levels

Note: Use Go's `log` package to record key system startup information or logs of specific levels.
Note: Currently, we use the following four levels to output logs; others are not used.

### 1. Debug

- **Purpose**: Development and troubleshooting.
- **Content**: Detailed technical information, such as variable values, system status, execution paths.
- **Scenario**: Code debugging, issue tracking.

### 2. Information

- **Purpose**: Record normal operation status and important events.
- **Content**: Key operations and significant events.
- **Scenario**: System monitoring, auditing.

### 3. Warning

- **Purpose**: Identify situations that may affect performance or stability.
- **Content**: Potential issues and recommended actions.
- **Scenario**: Fault prevention, performance optimization.

### 4. Error

- **Purpose**: Record critical problems or system errors.
- **Content**: Error descriptions, impact scope, faulty components.
- **Scenario**: Error handling, system recovery.

### General Principles

- **Consistency**: Maintain a uniform format.
- **Conciseness**: Avoid redundancy.
- **Security**: Do not log sensitive information.
- **Performance Consideration**: Be mindful of performance impact.

## Code Example

Log content can also be in English.

```go
func main() {
    log.Println("Application started") 

    result, err := performOperation()
    if err != nil {
        logs.Error("Error during operation: %v", err) 
        return
    }

    if result < 0 {
        logs.Warning("Warning: Result is negative, result: %d", result) 
    } else {
        logs.Info("Operation completed successfully, result: %d", result) 
    }

    log.Println("Application finished") 
}

func performOperation() (int, error) {
    logs.Debug("Starting operation") 

    rand.Seed(time.Now().UnixNano())
    num := rand.Intn(20) - 10 // Random number between -10 and 9

    // Simulate error 10% of the time
    if rand.Float32() < 0.1 {
        logs.Debug("Error encountered during operation") 
        return 0, fmt.Errorf("Random error occurred") 
    }

    logs.Debug("Operation completed, result: %d", num) 
    return num, nil
}
```
