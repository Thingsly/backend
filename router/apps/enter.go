package apps

type apps struct {
	User
	Role
	Casbin
	Dict
	OTA
	UpLoad
	ProtocolPlugin
	Device
	UiElements
	Board
	EventData
	TelemetryData
	AttributeData
	CommandData
	OperationLog
	Logo
	DataPolicy
	DeviceConfig
	DataScript
	NotificationGroup
	NotificationHistoryGroup
	NotificationServicesConfig
	Alarm
	SceneAutomations
	Scene
	SysFunction
	ServicePlugin
	ExpectedData
	OpenAPIKey
	MessagePush
}

var Model = new(apps)
