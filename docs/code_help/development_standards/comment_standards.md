# Commenting Conventions

Comment the code based on the situation, and reduce the need for comments through clear code structure and naming as much as possible.

The following three situations should be marked appropriately:

- For marking to-dos (TODO), defects to fix (FIXME), or other temporary solutions (HACK).

## Description

### 1. Overview

- **Clarity**: Comments should be clear, concise, and direct.
- **Accuracy**: Comments must accurately reflect the intent and behavior of the code.
- **Timely Updates**: Related comments should be updated when the code is modified.

### 2. Code Documentation Comments

- **Package Comments**: Every package should have a package comment, which is placed before the package declaration.
- **Function Comments**: Every public (exported) function should have a comment explaining the function's purpose and usage.
- **Type and Variable Comments**: Public structs, interfaces, type aliases, and global variables should have comments.
- **Godoc**: Comments should be compatible with the Godoc tool, using complete sentences and appropriate punctuation.

### 3. Special Comment Tags

- **TODO**: Marks work to be completed in the future or code that needs optimization.
- **FIXME**: Points to parts of the code that need fixing or refactoring.
- **HACK**: Identifies temporary solutions or code that needs improvement.
- **DEPRECATED**: Marks outdated code or functionality that is no longer recommended for use.

### 4. Behaviors to Avoid

- **Redundant Comments**: Do not comment obvious code.
- **Outdated Comments**: Ensure comments are synchronized with the code.
- **Commented-out Code**: Avoid keeping commented-out old code unless explicitly specified.

### 5. Other Considerations

- **Self-Documenting Code**: Reduce the need for comments by using clear code structure and naming conventions.
- **Internationalization**: Consider the language preferences of the team. Comments can be in English or the team's commonly used language.
- **Context-Specific Comments**: Additional comment rules might be needed based on the project or team needs, such as performance explanations, security considerations, etc.

## Using Godoc

- Run `go get golang.org/x/tools/cmd/godoc`
- In the terminal, run `godoc -http=:6060`
- View locally at `127.0.0.1:6060`

`Godoc` is the documentation generation tool for Go. It extracts comments from source code to generate documentation. To ensure your Go code generates clear and useful documentation, follow the commenting conventions.

## Comment Example

```go
package services

// Import necessary packages
// Example: "github.com/some/dependency"

// DeviceService provides service logic related to IoT devices.
// It encapsulates core functionality such as retrieving device status and updating configurations.
type DeviceService struct {
    // Declare any state or dependencies required for the service
    // Example: database connection
}

// NewDeviceService creates and returns a new instance of DeviceService.
// This function initializes necessary dependencies and state.
// TODO: Complete integration with the database.
func NewDeviceService() *DeviceService {
    // Return a new instance of DeviceService
    return &DeviceService{
        // Initialization logic, e.g., setting default values
    }
}

// GetDeviceStatus retrieves the status of a device based on its ID.
// This method may access the database or an external API to get the latest device information.
// FIXME: Optimize the performance of device status query.
func (s *DeviceService) GetDeviceStatus(deviceID string) string {
    // Logic to fetch the device status
    // Example return value
    return "Device status: Online"
}

// UpdateDeviceConfig updates the device configuration based on the device ID and configuration details.
// This includes validation of configuration parameters and applying changes.
// HACK: Currently uses a non-standard method for configuration updates; may change in the future.
func (s *DeviceService) UpdateDeviceConfig(deviceID string, config map[string]interface{}) {
    // Logic for updating the device configuration
    // Example: calling an API or updating a database record
}

// DeprecatedMethod is a method that has been deprecated.
// It contains old logic or algorithms.
// DEPRECATED: Please use NewMethod as a replacement.
func (s *DeviceService) DeprecatedMethod() {
    // Old handling logic
}

// NewMethod is a new method to replace DeprecatedMethod.
// It provides improved logic or algorithms.
func (s *DeviceService) NewMethod() {
    // New handling logic
}
```
