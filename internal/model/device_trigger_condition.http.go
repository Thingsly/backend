package model

const (
	// Trigger condition types
	DEVICE_TRIGGER_CONDITION_TYPE_ONE      = "10" // Single device
	DEVICE_TRIGGER_CONDITION_TYPE_MULTIPLE = "11" // Devices of the same type
	DEVICE_TRIGGER_CONDITION_TYPE_TIME     = "22" // Time range

	// Trigger parameter types
	TRIGGER_PARAM_TYPE_TEL        = "TEL"        // Telemetry (abbreviation)
	TRIGGER_PARAM_TYPE_TELEMETRY  = "TELEMETRY"  // Telemetry (full name)
	TRIGGER_PARAM_TYPE_ATTR       = "ATTR"       // Attribute
	TRIGGER_PARAM_TYPE_ATTRIBUTES = "ATTRIBUTES" // Attributes
	TRIGGER_PARAM_TYPE_EVT        = "EVT"        // Event (abbreviation)
	TRIGGER_PARAM_TYPE_EVENT      = "EVENT"      // Event (full name)
	TRIGGER_PARAM_TYPE_STATUS     = "STATUS"     // Device status

	// Comparison operators
	CONDITION_TRIGGER_OPERATOR_EQ      = "="       // Equal to
	CONDITION_TRIGGER_OPERATOR_NEQ     = "!="      // Not equal to
	CONDITION_TRIGGER_OPERATOR_GT      = ">"       // Greater than
	CONDITION_TRIGGER_OPERATOR_LT      = "<"       // Less than
	CONDITION_TRIGGER_OPERATOR_GTE     = ">="      // Greater than or equal to
	CONDITION_TRIGGER_OPERATOR_LTE     = "<="      // Less than or equal to
	CONDITION_TRIGGER_OPERATOR_BETWEEN = "between" // Between (range)
	CONDITION_TRIGGER_OPERATOR_IN      = "in"      // Included in list

	// Automation action types
	AUTOMATE_ACTION_TYPE_ONE      = "10" // Single device
	AUTOMATE_ACTION_TYPE_MULTIPLE = "11" // Devices of the same type
	AUTOMATE_ACTION_TYPE_SCENE    = "20" // Activate scene
	AUTOMATE_ACTION_TYPE_ALARM    = "30" // Raise alarm
	AUTOMATE_ACTION_TYPE_SERVICE  = "40" // Invoke service
)
