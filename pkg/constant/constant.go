package constant

const (
	SYS_ADMIN    string = "SYS_ADMIN"    
	TENANT_USER  string = "TENANT_USER"  
	TENANT_ADMIN string = "TENANT_ADMIN" 
)

const EMPTY string = ""

const (
	DIRECT_CONNECTION  int = iota + 1 
	GATEWAY_DEVICE                    
	GATEWAY_SON_DEVICE                
)

const (
	Manual int = iota + 1 
	Auto                  
)

const (
	StatusOK              int = iota + 1 
	StatusFailed                         
	ResponseStatusOk                     
	ResponseSStatusFailed                
)

// DeviceModelSource
type DeviceModelSource string

const (
	TelemetrySource DeviceModelSource = "telemetry"
	AttributeSource DeviceModelSource = "attributes"
	EventSource     DeviceModelSource = "event"
	CommandSource   DeviceModelSource = "command"
)


type FormType string

const (
	CONFIG_FORM       FormType = "CFG"  
	VOUCHER_FORM      FormType = "VCR"  
	VOUCHER_TYPE_FORM FormType = "VCRT" 
)


const (
	DEVICE_TYPE_1 string = "1" 
	DEVICE_TYPE_2 string = "2" 
	DEVICE_TYPE_3 string = "3" 
)

// sys_function.enable_flag
const (
	EnableFlag  string = "enable"
	DisableFlag string = "disable"
)
