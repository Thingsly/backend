package api

import (
	"fmt"
	"net/http"

	"github.com/Thingsly/backend/pkg/errcode"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
)

type Controller struct {
	UserApi                       // User management
	DictApi                       // Dictionary management
	ProductApi                    // Product management
	OTAApi                        // OTA management
	UpLoadApi                     // File upload
	ProtocolPluginApi             // Protocol plugin
	DeviceApi                     // Device management
	DeviceModelApi                // Device model management
	UiElementsApi                 // UI elements control
	BoardApi                      // Dashboard
	TelemetryDataApi              // Telemetry data
	AttributeDataApi              // Attribute data
	EventDataApi                  // Event data
	CommandSetLogApi              // Command dispatch record
	OperationLogsApi              // System logs
	LogoApi                       // Logo management
	DataPolicyApi                 // Data cleanup policies
	DeviceConfigApi               // Device configuration
	DataScriptApi                 // Data processing scripts
	RoleApi                       // Role management
	CasbinApi                     // Permission management
	NotificationGroupApi          // Notification groups
	NotificationHistoryApi        // Notification history
	NotificationServicesConfigApi // Notification service configuration
	AlarmApi                      // Alarm management
	SceneAutomationsApi           // Scene automation
	SceneApi                      // Scene management
	SystemApi                     // System management
	SysFunctionApi                // System functionality settings
	VisPluginApi                  // Visualization plugins
	ServicePluginApi              // Service plugin management
	ServiceAccessApi              // Service access management
	ExpectedDataApi               // Expected data
	OpenAPIKeyApi                 // OpenAPI keys
	MessagePushApi                // Message push
	SystemMonitorApi              // System monitor
}

var (
	Controllers = new(Controller)
	Validate    *validator.Validate
)

func init() {
	Validate = validator.New()
}

// ValidateStruct validates the request structure
func ValidateStruct(i interface{}) error {
	err := Validate.Struct(i)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}

		var errors []string
		for _, err := range err.(validator.ValidationErrors) {
			var tError string
			switch err.Tag() {
			case "required":
				tError = fmt.Sprintf("Field '%s' is required", err.Field())
			case "email":
				tError = fmt.Sprintf("Field '%s' must be a valid email address", err.Field())
			case "gte":
				tError = fmt.Sprintf("The value of field '%s' must be at least %s", err.Field(), err.Param())
			case "lte":
				tError = fmt.Sprintf("The value of field '%s' must be at most %s", err.Field(), err.Param())
			default:
				tError = fmt.Sprintf("Field '%s' failed validation (%s)", err.Field(), validationErrorToText(err))
			}
			errors = append(errors, tError)
		}

		return fmt.Errorf("%s", errors[0])
	}
	return nil
}

// validationErrorToText converts validation errors to more descriptive text
func validationErrorToText(e validator.FieldError) string {
	switch e.Tag() {
	case "min":
		return fmt.Sprintf("At least %s characters", e.Param())
	case "max":
		return fmt.Sprintf("At most %s characters", e.Param())
	case "len":
		return fmt.Sprintf("Must be %s characters", e.Param())
	// Add more cases as needed
	default:
		return "Does not meet validation rules"
	}
}

// Define a unified HTTP response structure
type ApiResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorHandler - Unified error handling
// func ErrorHandler(c *gin.Context, code int, err error) {
// 	// if strings.Contains(err.Error(), "SQLSTATE 23503") {
// 	// 	// Handle foreign key constraint violation
// 	// 	err = fmt.Errorf("Operation cannot be completed: Please delete the associated data before trying again")
// 	// }
// 	// if strings.Contains(err.Error(), "SQLSTATE 23505") {
// 	// 	// Handle unique key constraint violation
// 	// 	err = fmt.Errorf("Operation cannot be completed: Identical data already exists")
// 	// }
// 	// fmt.Printf("%T\n", err)
// 	// // Check if the error is of type *pgconn.PgError
// 	// var pgErr *pgconn.PgError
// 	// if errors.As(err, &pgErr) {
// 	// 	logrus.Error("-----------------")
// 	// 	// Now pgErr contains the *pgconn.PgError part of err (if it exists)
// 	// 	if pgErr.SQLState() == "23503" {
// 	// 		// This is a foreign key constraint violation error
// 	// 		err = fmt.Errorf("Foreign key constraint violation: %w", err)
// 	// 	}
// 	// }
// 	c.JSON(http.StatusOK, ApiResponse{
// 		Code:    code,
// 		Message: err.Error(),
// 	})
// }

// // SuccessHandler - Unified success response
// func SuccessHandler(c *gin.Context, message string, data interface{}) {
// 	c.JSON(http.StatusOK, ApiResponse{
// 		Code:    http.StatusOK,
// 		Message: message,
// 		Data:    data,
// 	})
// }

// // SuccessOK - Unified success response with default success message
// func SuccessOK(c *gin.Context) {
// 	c.JSON(http.StatusOK, ApiResponse{
// 		Code:    http.StatusOK,
// 		Message: "Success",
// 	})
// }

func BindAndValidate(c *gin.Context, obj interface{}) bool {
	// Determine the request method
	switch c.Request.Method {
	case http.MethodGet:
		if err := c.ShouldBindQuery(obj); err != nil {
			c.Error(errcode.NewWithMessage(errcode.CodeParamError, err.Error()))
			return false
		}
	case http.MethodPost, http.MethodPut, http.MethodDelete:
		if err := c.ShouldBindJSON(obj); err != nil {
			c.Error(errcode.NewWithMessage(errcode.CodeParamError, err.Error()))
			return false
		}
	}

	if err := ValidateStruct(obj); err != nil {
		c.Error(errcode.NewWithMessage(errcode.CodeParamError, err.Error()))
		return false
	}

	return true
}

var Wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(_ *http.Request) bool {
		// No cross-origin checks
		return true
	},
}
