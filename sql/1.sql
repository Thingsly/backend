-- public.alarm_config definition

-- Drop table
-- DROP TABLE public.alarm_config;

CREATE TABLE public.alarm_config (
	id varchar(36) NOT NULL, -- Alarm ID
	"name" varchar(255) NOT NULL, -- Alarm name
	description varchar(255) NULL, -- Alarm description
	alarm_level varchar(10) NOT NULL, -- Alarm level: H - High, M - Medium, L - Low
	notification_group_id varchar(36) NOT NULL, -- Notification group ID
	created_at timestamptz(6) NOT NULL, -- Creation time
	updated_at timestamptz(6) NOT NULL, -- Update time
	tenant_id varchar(36) NOT NULL, -- Tenant ID
	remark varchar(255) NULL, -- Remarks
	enabled varchar(10) NOT NULL, -- Whether enabled: Y - Enabled, N - Disabled
	CONSTRAINT alarm_config_pk PRIMARY KEY (id) -- Primary key constraint
);

-- Add comments to the table and columns
COMMENT ON TABLE public.alarm_config IS 'Alarm configuration';

-- Column comments
COMMENT ON COLUMN public.alarm_config."name" IS 'Alarm name';
COMMENT ON COLUMN public.alarm_config.description IS 'Alarm description';
COMMENT ON COLUMN public.alarm_config.alarm_level IS 'Alarm level: H - High, M - Medium, L - Low';
COMMENT ON COLUMN public.alarm_config.notification_group_id IS 'Notification group ID';
COMMENT ON COLUMN public.alarm_config.enabled IS 'Whether enabled: Y - Enabled, N - Disabled';


-- public.alarm_history definition

-- Drop table
-- DROP TABLE public.alarm_history;

CREATE TABLE public.alarm_history (
	id varchar(36) NOT NULL, -- Alarm ID
	alarm_config_id varchar(36) NOT NULL, -- Alarm configuration ID
	group_id varchar(36) NOT NULL, -- Group ID
	scene_automation_id varchar(36) NOT NULL, -- Scene automation ID
	"name" varchar(255) NOT NULL, -- Alarm name
	description varchar(255) NULL, -- Alarm description
	"content" text NULL, -- Content (reason for the alarm)
	alarm_status varchar(3) NOT NULL, -- Alarm status: L - Low, M - Medium, H - High, N - Normal
	tenant_id varchar(36) NOT NULL, -- Tenant ID
	remark varchar(255) NULL, -- Remarks
	create_at timestamptz(6) NOT NULL, -- Creation time
	alarm_device_list jsonb NOT NULL, -- Triggered device IDs
	CONSTRAINT alarm_history_pkey PRIMARY KEY (id) -- Primary key constraint
);

-- Column comments
COMMENT ON COLUMN public.alarm_history."name" IS 'Alarm name';
COMMENT ON COLUMN public.alarm_history.description IS 'Alarm description';
COMMENT ON COLUMN public.alarm_history."content" IS 'Content (reason for the alarm)';
COMMENT ON COLUMN public.alarm_history.alarm_status IS 'L - Low, M - Medium, H - High, N - Normal';
COMMENT ON COLUMN public.alarm_history.tenant_id IS 'Tenant ID';
COMMENT ON COLUMN public.alarm_history.create_at IS 'Creation time';
COMMENT ON COLUMN public.alarm_history.alarm_device_list IS 'Triggered device IDs';


-- public.boards definition

-- Drop table
-- DROP TABLE public.boards;

CREATE TABLE public.boards (
	id varchar(36) NOT NULL, -- ID
	"name" varchar(255) NOT NULL, -- Board name
	config json NULL DEFAULT '{}'::json, -- Board configuration
	tenant_id varchar(36) NOT NULL, -- Tenant ID (unique)
	created_at timestamptz(6) NOT NULL, -- Creation time
	updated_at timestamptz(6) NOT NULL, -- Update time
	home_flag varchar(2) NOT NULL, -- Home flag (default N, Y for home)
	description varchar(500) NULL, -- Description
	remark varchar(255) NULL, -- Remarks
	menu_flag varchar(2) NULL, -- Menu flag (default N, Y for menu)
	CONSTRAINT boards_pkey PRIMARY KEY (id) -- Primary key constraint
);

-- Column comments
COMMENT ON COLUMN public.boards.id IS 'ID';
COMMENT ON COLUMN public.boards."name" IS 'Board name';
COMMENT ON COLUMN public.boards.config IS 'Board configuration';
COMMENT ON COLUMN public.boards.tenant_id IS 'Tenant ID (unique)';
COMMENT ON COLUMN public.boards.home_flag IS 'Home flag (default N, Y for home)';
COMMENT ON COLUMN public.boards.description IS 'Description';
COMMENT ON COLUMN public.boards.remark IS 'Remarks';
COMMENT ON COLUMN public.boards.menu_flag IS 'Menu flag (default N, Y for menu)';

-- public.data_policy definition

-- Drop table

-- DROP TABLE public.data_policy;

-- public.data_policy definition

CREATE TABLE public.data_policy (
	id varchar(36) NOT NULL, -- ID
	data_type varchar(1) NOT NULL, -- Cleanup type: 1 - Device data, 2 - Operation logs
	retention_days int4 NOT NULL, -- Data retention time (days)
	last_cleanup_time timestamptz(6) NULL, -- Last cleanup time
	last_cleanup_data_time timestamptz(6) NULL, -- Last cleanup data time (actual data cleanup time)
	enabled varchar(1) NOT NULL, -- Whether enabled: 1 - Enabled, 2 - Disabled
	remark varchar(255) NULL, -- Remarks
	CONSTRAINT data_policy_pkey PRIMARY KEY (id) -- Primary key constraint
);

-- Column comments

COMMENT ON COLUMN public.data_policy.id IS 'ID';
COMMENT ON COLUMN public.data_policy.data_type IS 'Cleanup type: 1 - Device data, 2 - Operation logs';
COMMENT ON COLUMN public.data_policy.retention_days IS 'Data retention time (days)';
COMMENT ON COLUMN public.data_policy.last_cleanup_time IS 'Last cleanup time';
COMMENT ON COLUMN public.data_policy.last_cleanup_data_time IS 'Last cleanup data time (actual data cleanup time)';
COMMENT ON COLUMN public.data_policy.enabled IS 'Whether enabled: 1 - Enabled, 2 - Disabled';
COMMENT ON COLUMN public.data_policy.remark IS 'Remarks';


-- public.device_model_custom_commands definition

-- Drop table
-- DROP TABLE public.device_model_custom_commands;

CREATE TABLE public.device_model_custom_commands (
	id varchar(36) NOT NULL, -- ID
	device_template_id varchar(36) NOT NULL, -- Device template ID
	buttom_name varchar(255) NOT NULL, -- Button name
	data_identifier varchar(255) NOT NULL, -- Data identifier
	description varchar(500) NULL, -- Description
	instruct text NULL, -- Instruction content
	enable_status varchar(10) NOT NULL, -- Enable status: enable - Enabled, disable - Disabled
	remark varchar(255) NULL, -- Remarks
	tenant_id varchar NOT NULL, -- Tenant ID
	CONSTRAINT device_model_custom_commands_pk PRIMARY KEY (id) -- Primary key constraint
);

-- Column comments

COMMENT ON COLUMN public.device_model_custom_commands.id IS 'ID';
COMMENT ON COLUMN public.device_model_custom_commands.device_template_id IS 'Device template ID';
COMMENT ON COLUMN public.device_model_custom_commands.buttom_name IS 'Button name';
COMMENT ON COLUMN public.device_model_custom_commands.data_identifier IS 'Data identifier';
COMMENT ON COLUMN public.device_model_custom_commands.description IS 'Description';
COMMENT ON COLUMN public.device_model_custom_commands.instruct IS 'Instruction content';
COMMENT ON COLUMN public.device_model_custom_commands.enable_status IS 'Enable status: enable - Enabled, disable - Disabled';
COMMENT ON COLUMN public.device_model_custom_commands.remark IS 'Remarks';


-- public.device_templates definition

-- Drop table
-- DROP TABLE public.device_templates;

CREATE TABLE public.device_templates (
	id varchar(36) NOT NULL, -- ID
	"name" varchar(255) NOT NULL, -- Template name
	author varchar(36) NULL DEFAULT ''::character varying, -- Author
	"version" varchar(50) NULL DEFAULT ''::character varying, -- Version number
	description varchar(500) NULL DEFAULT ''::character varying, -- Description
	tenant_id varchar(36) NOT NULL DEFAULT ''::character varying, -- Tenant ID (unique)
	created_at timestamptz(6) NOT NULL, -- Creation time
	updated_at timestamptz(6) NOT NULL, -- Update time
	flag int2 NULL DEFAULT 1, -- Flag (default 1)
	"label" varchar(255) NULL, -- Label
	web_chart_config json NULL, -- Web chart configuration
	app_chart_config json NULL, -- App chart configuration
	remark varchar(255) NULL, -- Remarks
	"path" varchar(999) NULL, -- Image path
	CONSTRAINT device_templates_pkey PRIMARY KEY (id) -- Primary key constraint
);

-- Column comments

COMMENT ON COLUMN public.device_templates.id IS 'ID';
COMMENT ON COLUMN public.device_templates."name" IS 'Template name';
COMMENT ON COLUMN public.device_templates.author IS 'Author';
COMMENT ON COLUMN public.device_templates."version" IS 'Version number';
COMMENT ON COLUMN public.device_templates.description IS 'Description';
COMMENT ON COLUMN public.device_templates.flag IS 'Flag (default 1)';
COMMENT ON COLUMN public.device_templates."label" IS 'Label';
COMMENT ON COLUMN public.device_templates.web_chart_config IS 'Web chart configuration';
COMMENT ON COLUMN public.device_templates.app_chart_config IS 'App chart configuration';
COMMENT ON COLUMN public.device_templates.remark IS 'Remarks';
COMMENT ON COLUMN public.device_templates."path" IS 'Image path';


-- public.device_user_logs definition

-- Drop table

-- DROP TABLE public.device_user_logs;

-- public.device_user_logs definition

CREATE TABLE public.device_user_logs (
	id varchar(36) NOT NULL, -- ID
	device_nums int4 NOT NULL DEFAULT 0, -- Number of devices
	device_on int4 NOT NULL DEFAULT 0, -- Number of devices turned on
	created_at timestamptz(6) NOT NULL DEFAULT CURRENT_TIMESTAMP, -- Creation time
	tenant_id varchar(36) NOT NULL, -- Tenant ID
	CONSTRAINT device_user_logs_pkey PRIMARY KEY (id) -- Primary key constraint
);

-- Column comments

COMMENT ON COLUMN public.device_user_logs.tenant_id IS 'Tenant ID';


-- public."groups" definition

-- Drop table
-- DROP TABLE public."groups";

CREATE TABLE public."groups" (
	id varchar(36) NOT NULL, -- ID
	parent_id varchar(36) NULL DEFAULT 0, -- Default 0 is the parent group
	tier int4 NOT NULL DEFAULT 1, -- Level (starting from 1)
	"name" varchar(255) NOT NULL, -- Group name
	description varchar(255) NULL, -- Description
	created_at timestamptz(6) NOT NULL, -- Creation time
	updated_at timestamptz(6) NOT NULL, -- Update time
	remark varchar(255) NULL, -- Remarks
	tenant_id varchar(36) NOT NULL, -- Tenant ID
	CONSTRAINT groups_pkey PRIMARY KEY (id) -- Primary key constraint
);

-- Column comments

COMMENT ON COLUMN public."groups".parent_id IS 'Default 0 is the parent group';
COMMENT ON COLUMN public."groups".tier IS 'Level (starting from 1)';
COMMENT ON COLUMN public."groups"."name" IS 'Group name';
COMMENT ON COLUMN public."groups".description IS 'Description';
COMMENT ON COLUMN public."groups".created_at IS 'Creation time';
COMMENT ON COLUMN public."groups".updated_at IS 'Update time';
COMMENT ON COLUMN public."groups".tenant_id IS 'Tenant ID';


-- public.logo definition

-- Drop table
-- DROP TABLE public.logo;

CREATE TABLE public.logo (
	id varchar(36) NOT NULL, -- ID
	system_name varchar(99) NOT NULL, -- System name
	logo_cache varchar(255) NOT NULL, -- Site logo
	logo_background varchar(255) NOT NULL, -- Loading page logo
	logo_loading varchar(255) NOT NULL, -- Loading page logo
	home_background varchar(255) NOT NULL, -- Homepage background
	remark varchar(255) NULL, -- Remarks
	CONSTRAINT logo_pkey PRIMARY KEY (id) -- Primary key constraint
);

-- Column comments

COMMENT ON COLUMN public.logo.id IS 'ID';
COMMENT ON COLUMN public.logo.system_name IS 'System name';
COMMENT ON COLUMN public.logo.logo_cache IS 'Site logo';
COMMENT ON COLUMN public.logo.logo_background IS 'Loading page logo';
COMMENT ON COLUMN public.logo.logo_loading IS 'Loading page logo';
COMMENT ON COLUMN public.logo.home_background IS 'Homepage background';


-- public.notification_groups definition

-- Drop table

-- DROP TABLE public.notification_groups;

-- public.notification_groups definition

CREATE TABLE public.notification_groups (
	id varchar(36) NOT NULL, -- ID
	"name" varchar(99) NOT NULL, -- Name
	notification_type varchar(25) NOT NULL, -- Notification type: MEMBER - Member notification, EMAIL - Email notification, SME - SMS notification, VOICE - Voice notification, WEBHOOK - Webhook
	status varchar(10) NOT NULL, -- Notification status: ON - Enabled, OFF - Disabled
	notification_config jsonb NULL, -- Notification configuration
	description varchar(255) NULL, -- Description
	tenant_id varchar(36) NOT NULL, -- Tenant ID
	created_at timestamptz(6) NOT NULL, -- Creation time
	updated_at timestamptz(6) NOT NULL, -- Update time
	remark varchar(255) NULL, -- Remarks
	CONSTRAINT notification_groups_pkey PRIMARY KEY (id) -- Primary key constraint
);

-- Column comments

COMMENT ON COLUMN public.notification_groups."name" IS 'Name';
COMMENT ON COLUMN public.notification_groups.notification_type IS 'Notification type: MEMBER - Member notification, EMAIL - Email notification, SME - SMS notification, VOICE - Voice notification, WEBHOOK - Webhook';
COMMENT ON COLUMN public.notification_groups.status IS 'Notification status: ON - Enabled, OFF - Disabled';
COMMENT ON COLUMN public.notification_groups.notification_config IS 'Notification configuration';
COMMENT ON COLUMN public.notification_groups.description IS 'Description';
COMMENT ON COLUMN public.notification_groups.tenant_id IS 'Tenant ID';
COMMENT ON COLUMN public.notification_groups.created_at IS 'Creation time';
COMMENT ON COLUMN public.notification_groups.updated_at IS 'Update time';
COMMENT ON COLUMN public.notification_groups.remark IS 'Remarks';


-- public.notification_histories definition

CREATE TABLE public.notification_histories (
	id varchar(36) NOT NULL, -- ID
	send_time timestamptz(6) NOT NULL, -- Send time
	send_content text NULL, -- Send content
	send_target varchar(255) NOT NULL, -- Send target
	send_result varchar(25) NULL, -- Send result: SUCCESS - Success, FAILURE - Failure
	notification_type varchar(25) NOT NULL, -- Notification type: MEMBER - Member notification, EMAIL - Email notification, SME - SMS notification, VOICE - Voice notification, WEBHOOK - Webhook
	tenant_id varchar(36) NOT NULL, -- Tenant ID
	remark varchar(255) NULL, -- Remarks
	CONSTRAINT notification_histories_pkey PRIMARY KEY (id) -- Primary key constraint
);

-- Column comments

COMMENT ON COLUMN public.notification_histories.send_time IS 'Send time';
COMMENT ON COLUMN public.notification_histories.send_content IS 'Send content';
COMMENT ON COLUMN public.notification_histories.send_target IS 'Send target';
COMMENT ON COLUMN public.notification_histories.send_result IS 'Send result: SUCCESS - Success, FAILURE - Failure';
COMMENT ON COLUMN public.notification_histories.notification_type IS 'Notification type: MEMBER - Member notification, EMAIL - Email notification, SME - SMS notification, VOICE - Voice notification, WEBHOOK - Webhook';
COMMENT ON COLUMN public.notification_histories.tenant_id IS 'Tenant ID';
COMMENT ON COLUMN public.notification_histories.remark IS 'Remarks';


-- public.notification_services_config definition

CREATE TABLE public.notification_services_config (
	id varchar(36) NOT NULL, -- ID
	config json NULL, -- Notification configuration
	notice_type varchar(36) NOT NULL, -- Notification type: EMAIL - Email configuration, SME - SMS configuration
	status varchar(36) NOT NULL, -- Status: OPEN - Enabled, CLOSE - Disabled
	remark varchar(255) NULL, -- Remarks
	CONSTRAINT notification_services_config_pkey PRIMARY KEY (id) -- Primary key constraint
);

-- Column comments

COMMENT ON COLUMN public.notification_services_config.config IS 'Notification configuration';
COMMENT ON COLUMN public.notification_services_config.notice_type IS 'Notification type: EMAIL - Email configuration, SME - SMS configuration';
COMMENT ON COLUMN public.notification_services_config.status IS 'Status: OPEN - Enabled, CLOSE - Disabled';


-- public.operation_logs definition

CREATE TABLE public.operation_logs (
	id varchar(36) NOT NULL, -- ID
	ip varchar(36) NOT NULL, -- Request IP
	"path" varchar(2000) NULL, -- Request URL
	user_id varchar(36) NOT NULL, -- Operator user
	"name" varchar(255) NULL, -- Interface name
	created_at timestamptz(6) NOT NULL, -- Creation time
	latency int8 NULL, -- Latency (ms)
	request_message text NULL, -- Request content
	response_message text NULL, -- Response content
	tenant_id varchar(36) NOT NULL, -- Tenant ID
	remark varchar(255) NULL, -- Remarks
	CONSTRAINT operation_logs_pkey PRIMARY KEY (id) -- Primary key constraint
);

-- Column comments

COMMENT ON COLUMN public.operation_logs.id IS 'ID';
COMMENT ON COLUMN public.operation_logs.ip IS 'Request IP';
COMMENT ON COLUMN public.operation_logs."path" IS 'Request URL';
COMMENT ON COLUMN public.operation_logs.user_id IS 'Operator user';
COMMENT ON COLUMN public.operation_logs."name" IS 'Interface name';
COMMENT ON COLUMN public.operation_logs.created_at IS 'Creation time';
COMMENT ON COLUMN public.operation_logs.latency IS 'Latency (ms)';
COMMENT ON COLUMN public.operation_logs.request_message IS 'Request content';
COMMENT ON COLUMN public.operation_logs.response_message IS 'Response content';
COMMENT ON COLUMN public.operation_logs.tenant_id IS 'Tenant ID';

-- public.ota_upgrade_packages definition

-- Drop table

-- DROP TABLE public.ota_upgrade_packages;

-- public.ota_upgrade_packages definition

CREATE TABLE public.ota_upgrade_packages (
	id varchar(36) NOT NULL, -- ID
	"name" varchar(200) NOT NULL, -- Upgrade package name
	"version" varchar(36) NOT NULL, -- Upgrade package version number
	target_version varchar(36) NULL, -- Target version number for upgrade
	device_config_id varchar(36) NOT NULL, -- Device configuration ID
	"module" varchar(36) NULL, -- Module name
	package_type int2 NOT NULL, -- Upgrade package type: 1 - Differential, 2 - Full package
	signature_type varchar(36) NULL, -- Signature algorithm: MD5, SHA256
	additional_info json NULL DEFAULT '{}'::json, -- Additional information
	description varchar(500) NULL, -- Description
	package_url varchar(500) NULL, -- Package download URL
	created_at timestamptz(6) NOT NULL, -- Creation time
	updated_at timestamptz(6) NULL, -- Update time
	remark varchar(255) NULL, -- Remarks
	signature varchar(255) NULL, -- Upgrade package signature
	tenant_id varchar(36) NULL, -- Tenant ID
	CONSTRAINT ota_upgrade_packages_pkey PRIMARY KEY (id) -- Primary key constraint
);

-- Column comments

COMMENT ON COLUMN public.ota_upgrade_packages.id IS 'ID';
COMMENT ON COLUMN public.ota_upgrade_packages."name" IS 'Upgrade package name';
COMMENT ON COLUMN public.ota_upgrade_packages."version" IS 'Upgrade package version number';
COMMENT ON COLUMN public.ota_upgrade_packages.target_version IS 'Target version number for upgrade';
COMMENT ON COLUMN public.ota_upgrade_packages.device_config_id IS 'Device configuration ID';
COMMENT ON COLUMN public.ota_upgrade_packages."module" IS 'Module name';
COMMENT ON COLUMN public.ota_upgrade_packages.package_type IS 'Upgrade package type: 1 - Differential, 2 - Full package';
COMMENT ON COLUMN public.ota_upgrade_packages.signature_type IS 'Signature algorithm: MD5, SHA256';
COMMENT ON COLUMN public.ota_upgrade_packages.additional_info IS 'Additional information';
COMMENT ON COLUMN public.ota_upgrade_packages.description IS 'Description';
COMMENT ON COLUMN public.ota_upgrade_packages.package_url IS 'Package download URL';
COMMENT ON COLUMN public.ota_upgrade_packages.created_at IS 'Creation time';
COMMENT ON COLUMN public.ota_upgrade_packages.updated_at IS 'Update time';
COMMENT ON COLUMN public.ota_upgrade_packages.remark IS 'Remarks';
COMMENT ON COLUMN public.ota_upgrade_packages.signature IS 'Upgrade package signature';


-- public.ota_upgrade_tasks definition

CREATE TABLE public.ota_upgrade_tasks (
	id varchar(36) NOT NULL, -- ID
	"name" varchar(200) NOT NULL, -- Task name
	ota_upgrade_package_id varchar(36) NOT NULL, -- Upgrade package ID (foreign key, cascade delete)
	description varchar(500) NULL, -- Description
	created_at timestamptz(6) NOT NULL, -- Creation time
	remark varchar(255) NULL, -- Remarks
	CONSTRAINT ota_upgrade_tasks_pkey PRIMARY KEY (id) -- Primary key constraint
);

-- Column comments

COMMENT ON COLUMN public.ota_upgrade_tasks.id IS 'ID';
COMMENT ON COLUMN public.ota_upgrade_tasks."name" IS 'Task name';
COMMENT ON COLUMN public.ota_upgrade_tasks.ota_upgrade_package_id IS 'Upgrade package ID (foreign key, cascade delete)';
COMMENT ON COLUMN public.ota_upgrade_tasks.description IS 'Description';
COMMENT ON COLUMN public.ota_upgrade_tasks.created_at IS 'Creation time';
COMMENT ON COLUMN public.ota_upgrade_tasks.remark IS 'Remarks';


-- public.protocol_plugins definition

CREATE TABLE public.protocol_plugins (
	id varchar(36) NOT NULL, -- ID
	"name" varchar(36) NOT NULL, -- Plugin name
	device_type int2 NOT NULL DEFAULT 1, -- Device type (1 - Direct device, 2 - Gateway device, default is direct device)
	protocol_type varchar(50) NOT NULL, -- Protocol type
	access_address varchar(500) NULL, -- Access address
	http_address varchar(500) NULL, -- HTTP service address
	sub_topic_prefix varchar(500) NULL, -- Plugin subscription prefix
	description varchar(500) NULL, -- Description
	additional_info varchar(1000) NULL, -- Additional information
	created_at timestamptz(6) NOT NULL, -- Creation time
	update_at timestamptz(6) NOT NULL, -- Update time
	remark varchar(255) NULL, -- Remarks
	CONSTRAINT protocol_plugins_pkey PRIMARY KEY (id) -- Primary key constraint
);

-- Column comments

COMMENT ON COLUMN public.protocol_plugins.id IS 'ID';
COMMENT ON COLUMN public.protocol_plugins."name" IS 'Plugin name';
COMMENT ON COLUMN public.protocol_plugins.device_type IS 'Device type (1 - Direct device, 2 - Gateway device, default is direct device)';
COMMENT ON COLUMN public.protocol_plugins.protocol_type IS 'Protocol type';
COMMENT ON COLUMN public.protocol_plugins.access_address IS 'Access address';
COMMENT ON COLUMN public.protocol_plugins.http_address IS 'HTTP service address';
COMMENT ON COLUMN public.protocol_plugins.sub_topic_prefix IS 'Plugin subscription prefix';
COMMENT ON COLUMN public.protocol_plugins.description IS 'Description';
COMMENT ON COLUMN public.protocol_plugins.additional_info IS 'Additional information';
COMMENT ON COLUMN public.protocol_plugins.created_at IS 'Creation time';
COMMENT ON COLUMN public.protocol_plugins.update_at IS 'Update time';
COMMENT ON COLUMN public.protocol_plugins.remark IS 'Remarks';


-- public.roles definition

CREATE TABLE public.roles (
	id varchar(36) NOT NULL, -- ID
	"name" varchar(99) NOT NULL, -- Name
	description varchar(255) NULL, -- Description
	created_at timestamp NULL, -- Creation time
	updated_at timestamp NULL, -- Update time
	tenant_id varchar(36) NULL, -- Tenant ID
	CONSTRAINT roles_pkey PRIMARY KEY (id) -- Primary key constraint
);

-- Column comments

COMMENT ON COLUMN public.roles.id IS 'ID';
COMMENT ON COLUMN public.roles."name" IS 'Name';
COMMENT ON COLUMN public.roles.description IS 'Description';
COMMENT ON COLUMN public.roles.created_at IS 'Creation time';
COMMENT ON COLUMN public.roles.updated_at IS 'Update time';
COMMENT ON COLUMN public.roles.tenant_id IS 'Tenant ID';


-- public.scene_automations definition

-- Drop table

-- DROP TABLE public.scene_automations;

-- public.scene_automations definition

CREATE TABLE public.scene_automations (
	id varchar(36) NOT NULL, -- Automation ID
	"name" varchar(255) NOT NULL, -- Name
	description varchar(255) NULL, -- Description
	enabled varchar(10) NOT NULL, -- Enabled status: Y = Enabled, N = Disabled
	tenant_id varchar(36) NOT NULL, -- Tenant ID
	creator varchar(36) NOT NULL, -- Creator ID
	updator varchar(36) NOT NULL, -- Updator ID
	created_at timestamptz(6) NOT NULL, -- Creation time
	updated_at timestamptz(6) NULL, -- Update time
	remark varchar(255) NULL, -- Remarks
	CONSTRAINT scene_automations_pkey PRIMARY KEY (id) -- Primary key constraint
);

-- Column comments

COMMENT ON COLUMN public.scene_automations.id IS 'Automation ID';
COMMENT ON COLUMN public.scene_automations."name" IS 'Name';
COMMENT ON COLUMN public.scene_automations.description IS 'Description';
COMMENT ON COLUMN public.scene_automations.enabled IS 'Enabled status: Y = Enabled, N = Disabled';
COMMENT ON COLUMN public.scene_automations.tenant_id IS 'Tenant ID';
COMMENT ON COLUMN public.scene_automations.creator IS 'Creator ID';
COMMENT ON COLUMN public.scene_automations.updator IS 'Updator ID';
COMMENT ON COLUMN public.scene_automations.created_at IS 'Creation time';
COMMENT ON COLUMN public.scene_automations.updated_at IS 'Update time';


-- public.scene_info definition

CREATE TABLE public.scene_info (
	id varchar(36) NOT NULL, -- Scene ID
	"name" varchar(255) NOT NULL, -- Name
	description varchar(255) NULL, -- Description
	tenant_id varchar(36) NOT NULL, -- Tenant ID
	creator varchar(36) NOT NULL, -- Creator ID
	updator varchar(36) NULL, -- Updator ID
	created_at timestamptz(6) NOT NULL, -- Creation time
	updated_at timestamptz(6) NULL, -- Update time
	CONSTRAINT scene_info_pkey PRIMARY KEY (id) -- Primary key constraint
);

-- Column comments

COMMENT ON COLUMN public.scene_info."name" IS 'Name';
COMMENT ON COLUMN public.scene_info.description IS 'Description';
COMMENT ON COLUMN public.scene_info.tenant_id IS 'Tenant ID';
COMMENT ON COLUMN public.scene_info.creator IS 'Creator ID';
COMMENT ON COLUMN public.scene_info.updator IS 'Updator ID';
COMMENT ON COLUMN public.scene_info.created_at IS 'Creation time';
COMMENT ON COLUMN public.scene_info.updated_at IS 'Update time';


-- public.sys_dict definition

CREATE TABLE public.sys_dict (
	id varchar(36) NOT NULL, -- Primary key ID
	dict_code varchar(36) NOT NULL, -- Dictionary identifier
	dict_value varchar(255) NOT NULL, -- Dictionary value
	created_at timestamptz(6) NOT NULL, -- Creation time
	remark varchar(255) NULL, -- Remarks
	CONSTRAINT sys_dict_dict_code_dict_value_key UNIQUE (dict_code, dict_value), -- Unique constraint on dict_code and dict_value
	CONSTRAINT sys_dict_pkey PRIMARY KEY (id) -- Primary key constraint
);

-- Column comments

COMMENT ON COLUMN public.sys_dict.id IS 'Primary key ID';
COMMENT ON COLUMN public.sys_dict.dict_code IS 'Dictionary identifier';
COMMENT ON COLUMN public.sys_dict.dict_value IS 'Dictionary value';
COMMENT ON COLUMN public.sys_dict.created_at IS 'Creation time';
COMMENT ON COLUMN public.sys_dict.remark IS 'Remarks';

-- Constraint comments

COMMENT ON CONSTRAINT sys_dict_dict_code_dict_value_key ON public.sys_dict IS 'Unique constraint on dict_code and dict_value';


-- public.sys_function definition

-- Drop table

-- DROP TABLE public.sys_function;

-- public.sys_function definition

CREATE TABLE public.sys_function (
	id varchar(36) NOT NULL, -- ID
	"name" varchar(50) NOT NULL, -- Function name
	enable_flag varchar(20) NOT NULL, -- Enable flag: enable = Enabled, disable = Disabled
	description varchar(500) NULL, -- Description
	remark varchar(255) NULL, -- Remarks
	CONSTRAINT sys_function_pk PRIMARY KEY (id) -- Primary key constraint
);

-- Column comments

COMMENT ON COLUMN public.sys_function.id IS 'ID';
COMMENT ON COLUMN public.sys_function."name" IS 'Function name';
COMMENT ON COLUMN public.sys_function.enable_flag IS 'Enable flag: enable = Enabled, disable = Disabled';
COMMENT ON COLUMN public.sys_function.description IS 'Description';
COMMENT ON COLUMN public.sys_function.remark IS 'Remarks';


-- public.sys_ui_elements definition

CREATE TABLE public.sys_ui_elements (
	id varchar(36) NOT NULL, -- Primary key ID
	parent_id varchar(36) NOT NULL, -- Parent element ID
	element_code varchar(100) NOT NULL, -- Element identifier
	element_type int2 NOT NULL, -- Element type: 1 = Menu, 2 = Directory, 3 = Button, 4 = Route
	orders int2 NULL, -- Sorting order
	param1 varchar(255) NULL, -- Parameter 1
	param2 varchar(255) NULL, -- Parameter 2
	param3 varchar(255) NULL, -- Parameter 3
	authority json NOT NULL, -- Permissions (multiple choices) 1 = System Admin, 2 = Tenant, e.g., [1, 2]
	description varchar(255) NULL, -- Description
	created_at timestamptz(6) NOT NULL, -- Creation time
	remark varchar(255) NULL, -- Remarks
	multilingual varchar(100) NULL, -- Multilingual identifier
	route_path varchar(255) NULL, -- Route path
	CONSTRAINT sys_ui_elements_pkey PRIMARY KEY (id) -- Primary key constraint
);

COMMENT ON TABLE public.sys_ui_elements IS 'UI Elements Table';

-- Column comments

COMMENT ON COLUMN public.sys_ui_elements.id IS 'Primary key ID';
COMMENT ON COLUMN public.sys_ui_elements.parent_id IS 'Parent element ID';
COMMENT ON COLUMN public.sys_ui_elements.element_code IS 'Element identifier';
COMMENT ON COLUMN public.sys_ui_elements.element_type IS 'Element type: 1 = Menu, 2 = Directory, 3 = Button, 4 = Route';
COMMENT ON COLUMN public.sys_ui_elements.orders IS 'Sorting order';
COMMENT ON COLUMN public.sys_ui_elements.authority IS 'Permissions (multiple choices) 1 = System Admin, 2 = Tenant, e.g., [1, 2]';
COMMENT ON COLUMN public.sys_ui_elements.description IS 'Description';
COMMENT ON COLUMN public.sys_ui_elements.multilingual IS 'Multilingual identifier';


-- public.telemetry_current_datas definition

CREATE TABLE public.telemetry_current_datas (
	device_id varchar(36) NOT NULL, -- Device ID
	"key" varchar(255) NOT NULL, -- Data identifier
	ts timestamptz(6) NOT NULL, -- Report time
	bool_v bool NULL, -- Boolean value
	number_v float8 NULL, -- Numeric value
	string_v text NULL, -- String value
	tenant_id varchar(36) NULL, -- Tenant ID
	CONSTRAINT telemetry_current_datas_unique UNIQUE (device_id, key) -- Unique constraint
);
CREATE INDEX telemetry_datas_ts_idx_copy1 ON public.telemetry_current_datas USING btree (ts DESC); -- Index for report time

-- Column comments

COMMENT ON COLUMN public.telemetry_current_datas.device_id IS 'Device ID';
COMMENT ON COLUMN public.telemetry_current_datas."key" IS 'Data identifier';
COMMENT ON COLUMN public.telemetry_current_datas.ts IS 'Report time';


-- public.telemetry_datas definition

CREATE TABLE public.telemetry_datas (
	device_id varchar(36) NOT NULL, -- Device ID
	"key" varchar(255) NOT NULL, -- Data identifier
	ts int8 NOT NULL, -- Report time
	bool_v bool NULL, -- Boolean value
	number_v float8 NULL, -- Numeric value
	string_v text NULL, -- String value
	tenant_id varchar(36) NULL, -- Tenant ID
	CONSTRAINT telemetry_datas_device_id_key_ts_key UNIQUE (device_id, key, ts) -- Unique constraint
);
CREATE INDEX telemetry_datas_ts_idx ON public.telemetry_datas USING btree (ts DESC); -- Index for report time

-- Column comments

COMMENT ON COLUMN public.telemetry_datas.device_id IS 'Device ID';
COMMENT ON COLUMN public.telemetry_datas."key" IS 'Data identifier';
COMMENT ON COLUMN public.telemetry_datas.ts IS 'Report time';

-- Table Triggers

-- Create a hypertable for partitioning based on the timestamp column, with a 24-hour partition interval
SELECT create_hypertable('telemetry_datas', 'ts', chunk_time_interval => 86400000000);

-- public.users definition

-- Drop table

-- DROP TABLE public.users;

-- public.users definition

CREATE TABLE public.users (
	id varchar(36) NOT NULL, -- ID
	"name" varchar(255) NULL, -- User's name
	phone_number varchar(50) NOT NULL, -- Phone number
	email varchar(255) NOT NULL, -- Email address
	status varchar(2) NULL, -- User status: F = Frozen, N = Normal
	authority varchar(50) NULL, -- Authority type: TENANT_ADMIN = Tenant Admin, TENANT_USER = Tenant User, SYS_ADMIN = System Admin
	"password" varchar(255) NOT NULL, -- Password
	tenant_id varchar(36) NULL, -- Tenant ID
	remark varchar(255) NULL, -- Remarks
	additional_info json NULL DEFAULT '{}'::json, -- Additional information in JSON format
	created_at timestamptz(6) NULL, -- Creation timestamp
	updated_at timestamptz(6) NULL, -- Last updated timestamp
	CONSTRAINT users_pkey PRIMARY KEY (id), -- Primary key constraint
	CONSTRAINT users_un UNIQUE (email) -- Unique constraint for email
);
COMMENT ON TABLE public.users IS 'Users';

-- Column comments

COMMENT ON COLUMN public.users.status IS 'User status: F = Frozen, N = Normal';
COMMENT ON COLUMN public.users.authority IS 'Authority type: TENANT_ADMIN = Tenant Admin, TENANT_USER = Tenant User, SYS_ADMIN = System Admin';


-- public.vis_dashboard definition

-- Drop table

-- DROP TABLE public.vis_dashboard;

CREATE TABLE public.vis_dashboard (
	id varchar(36) NOT NULL, -- ID
	relation_id varchar(36) NULL, -- Related ID
	json_data json NULL DEFAULT '{}'::json, -- JSON data
	dashboard_name varchar(99) NULL, -- Dashboard name
	create_at timestamp NULL, -- Creation timestamp
	sort int4 NULL, -- Sorting order
	remark varchar(255) NULL, -- Remarks
	tenant_id varchar(36) NULL, -- Tenant ID
	share_id varchar(36) NULL, -- Share ID
	CONSTRAINT vis_dashboard_pk PRIMARY KEY (id) -- Primary key constraint
);
COMMENT ON TABLE public.vis_dashboard IS 'Visualization Plugin';

-- Column comments

COMMENT ON COLUMN public.vis_dashboard.sort IS 'Sorting order';
COMMENT ON COLUMN public.vis_dashboard.share_id IS 'Share ID';


-- public.vis_files definition

-- Drop table

-- DROP TABLE public.vis_files;

CREATE TABLE public.vis_files (
	id varchar(36) NOT NULL, -- ID
	vis_plugin_id varchar(36) NOT NULL, -- Visualization plugin ID
	file_name varchar(150) NULL, -- File name
	file_url varchar(150) NULL, -- File URL
	file_size varchar(20) NULL, -- File size
	create_at int8 NULL, -- Creation timestamp
	remark varchar(255) NULL, -- Remarks
	CONSTRAINT tp_vis_files_pkey PRIMARY KEY (id) -- Primary key constraint
);
COMMENT ON TABLE public.vis_files IS 'Visualization Files Table';

-- Column comments

COMMENT ON COLUMN public.vis_files.vis_plugin_id IS 'Visualization plugin ID';
COMMENT ON COLUMN public.vis_files.file_name IS 'File name';
COMMENT ON COLUMN public.vis_files.file_url IS 'File URL';
COMMENT ON COLUMN public.vis_files.file_size IS 'File size';


-- public.vis_plugin definition

-- Drop table

-- DROP TABLE public.vis_plugin;

CREATE TABLE public.vis_plugin (
	id varchar(36) NOT NULL, -- ID
	tenant_id varchar(36) NOT NULL, -- Tenant ID
	plugin_name varchar(150) NOT NULL, -- Visualization plugin name
	plugin_description varchar(150) NULL, -- Plugin description
	create_at int8 NULL, -- Creation timestamp
	remark varchar(255) NULL, -- Remarks
	CONSTRAINT tp_vis_plugin_pkey PRIMARY KEY (id) -- Primary key constraint
);
COMMENT ON TABLE public.vis_plugin IS 'Visualization Plugin Table';

-- Column comments

COMMENT ON COLUMN public.vis_plugin.tenant_id IS 'Tenant ID';
COMMENT ON COLUMN public.vis_plugin.plugin_name IS 'Visualization plugin name';
COMMENT ON COLUMN public.vis_plugin.plugin_description IS 'Plugin description';

-- public.action_info definition

-- Drop table

-- DROP TABLE public.action_info;

-- public.action_info definition

CREATE TABLE public.action_info (
	id varchar(36) NOT NULL, -- ID
	scene_automation_id varchar(36) NOT NULL, -- Scene automation ID (foreign key - cascading delete)
	action_target varchar(255) NULL, -- Action target ID: device ID, scene ID, alarm ID; if the condition is a single device, this is empty
	action_type varchar(10) NOT NULL, -- Action type: 
	-- 10: Single device 
	-- 11: Single device class 
	-- 20: Activate scene 
	-- 30: Trigger alarm 
	-- 40: Service
	action_param_type varchar(10) NULL, -- Telemetry TEL attribute ATTR command CMD
	action_param varchar(50) NULL, -- Action parameter, valid for action types 10, 11, identifier
	action_value text NULL, -- Target value
	remark varchar(255) NULL, -- Remarks
	CONSTRAINT action_info_pkey PRIMARY KEY (id), -- Primary key constraint
	CONSTRAINT action_info_scene_automations_fk FOREIGN KEY (scene_automation_id) REFERENCES public.scene_automations(id) ON DELETE CASCADE -- Foreign key constraint
);

-- Column comments

COMMENT ON COLUMN public.action_info.scene_automation_id IS 'Scene automation ID (foreign key - cascading delete)';
COMMENT ON COLUMN public.action_info.action_target IS 'Action target ID: device ID, scene ID, alarm ID; if the condition is a single device, this is empty';
COMMENT ON COLUMN public.action_info.action_type IS 'Action type: 10 = Single device, 11 = Single device class, 20 = Activate scene, 30 = Trigger alarm, 40 = Service';
COMMENT ON COLUMN public.action_info.action_param_type IS 'Telemetry TEL attribute ATTR command CMD';
COMMENT ON COLUMN public.action_info.action_param IS 'Action parameter, valid for action types 10, 11, identifier';
COMMENT ON COLUMN public.action_info.action_value IS 'Target value';


-- public.alarm_info definition

-- Drop table

-- DROP TABLE public.alarm_info;

CREATE TABLE public.alarm_info (
	id varchar(36) NOT NULL, -- ID
	alarm_config_id varchar(36) NOT NULL, -- Alarm configuration ID
	"name" varchar(255) NOT NULL, -- Alarm name
	alarm_time timestamptz(6) NOT NULL, -- Alarm time
	description varchar(255) NULL, -- Alarm description
	"content" text NULL, -- Content
	processor varchar(36) NULL, -- Processor ID
	processing_result varchar(10) NOT NULL, -- Processing result: DOP = Processed, UND = Not processed, IGN = Ignored
	tenant_id varchar(36) NOT NULL, -- Tenant ID
	remark varchar(255) NULL, -- Remarks
	alarm_level varchar(10) NULL, -- Alarm level: L = Low, M = Medium, H = High
	CONSTRAINT alarm_info_pk PRIMARY KEY (id), -- Primary key constraint
	CONSTRAINT alarm_info_fk FOREIGN KEY (alarm_config_id) REFERENCES public.alarm_config(id) ON DELETE CASCADE -- Foreign key constraint
);
COMMENT ON TABLE public.alarm_info IS 'Alarm Information';

-- Column comments

COMMENT ON COLUMN public.alarm_info.alarm_config_id IS 'Alarm configuration ID';
COMMENT ON COLUMN public.alarm_info."name" IS 'Alarm name';
COMMENT ON COLUMN public.alarm_info.alarm_time IS 'Alarm time';
COMMENT ON COLUMN public.alarm_info.description IS 'Alarm description';
COMMENT ON COLUMN public.alarm_info."content" IS 'Content';
COMMENT ON COLUMN public.alarm_info.processor IS 'Processor ID';
COMMENT ON COLUMN public.alarm_info.processing_result IS 'Processing result: DOP = Processed, UND = Not processed, IGN = Ignored';
COMMENT ON COLUMN public.alarm_info.tenant_id IS 'Tenant ID';
COMMENT ON COLUMN public.alarm_info.alarm_level IS 'Alarm level: L = Low, M = Medium, H = High';


-- public.device_configs definition

-- Drop table

-- DROP TABLE public.device_configs;

CREATE TABLE public.device_configs (
	id varchar(36) NOT NULL, -- ID
	"name" varchar(99) NOT NULL, -- Name
	device_template_id varchar(36) NULL, -- Device template ID
	device_type varchar(9) NOT NULL, -- Device type: 1 = Direct device, 2 = Gateway device, 3 = Gateway sub-device
	protocol_type varchar(36) NULL, -- Protocol type
	voucher_type varchar(36) NULL, -- Voucher type
	protocol_config json NULL, -- Protocol form configuration
	device_conn_type varchar(36) NULL, -- Device connection type (default A): A = Device connected to platform, B = Platform connected to device
	additional_info json NULL DEFAULT '{}'::json, -- Additional information
	description varchar(255) NULL, -- Description
	tenant_id varchar(36) NOT NULL, -- Tenant ID
	created_at timestamptz(6) NOT NULL, -- Creation time
	updated_at timestamptz(6) NOT NULL, -- Last update time
	remark varchar(255) NULL, -- Remarks
	other_config json NULL, -- Other configurations
	CONSTRAINT device_configs_pkey PRIMARY KEY (id), -- Primary key constraint
	CONSTRAINT device_configs_device_templates_fk FOREIGN KEY (device_template_id) REFERENCES public.device_templates(id) ON DELETE RESTRICT -- Foreign key constraint
);

-- Column comments

COMMENT ON COLUMN public.device_configs.id IS 'ID';
COMMENT ON COLUMN public.device_configs."name" IS 'Name';
COMMENT ON COLUMN public.device_configs.device_template_id IS 'Device template ID';
COMMENT ON COLUMN public.device_configs.device_type IS 'Device type: 1 = Direct device, 2 = Gateway device, 3 = Gateway sub-device';
COMMENT ON COLUMN public.device_configs.protocol_type IS 'Protocol type';
COMMENT ON COLUMN public.device_configs.voucher_type IS 'Voucher type';
COMMENT ON COLUMN public.device_configs.protocol_config IS 'Protocol form configuration';
COMMENT ON COLUMN public.device_configs.device_conn_type IS 'Device connection type (default A): A = Device connected to platform, B = Platform connected to device';
COMMENT ON COLUMN public.device_configs.additional_info IS 'Additional information';
COMMENT ON COLUMN public.device_configs.description IS 'Description';
COMMENT ON COLUMN public.device_configs.tenant_id IS 'Tenant ID';
COMMENT ON COLUMN public.device_configs.created_at IS 'Creation time';
COMMENT ON COLUMN public.device_configs.updated_at IS 'Last update time';
COMMENT ON COLUMN public.device_configs.remark IS 'Remarks';
COMMENT ON COLUMN public.device_configs.other_config IS 'Other configurations';

-- public.device_model_attributes definition

-- Drop table

-- DROP TABLE public.device_model_attributes;

-- public.device_model_attributes definition

CREATE TABLE public.device_model_attributes (
	id varchar(36) NOT NULL, -- ID
	device_template_id varchar(36) NOT NULL, -- Device template ID
	data_name varchar(255) NULL, -- Data name
	data_identifier varchar(255) NOT NULL, -- Data identifier
	read_write_flag varchar(10) NULL, -- Read-write flag: R = Read, W = Write, RW = Read/Write
	data_type varchar(50) NULL, -- Data type: String, Number, Boolean, Enum
	unit varchar(50) NULL, -- Unit
	description varchar(255) NULL, -- Description
	additional_info json NULL, -- Additional information
	created_at timestamptz(6) NOT NULL, -- Creation time
	updated_at timestamptz(6) NOT NULL, -- Update time
	remark varchar(255) NULL, -- Remarks
	tenant_id varchar(36) NOT NULL, -- Tenant ID
	CONSTRAINT device_model_attributes_unique UNIQUE (device_template_id, data_identifier), -- Unique constraint for device template ID and data identifier
	CONSTRAINT device_model_telemetry_copy1_pkey PRIMARY KEY (id), -- Primary key constraint
	CONSTRAINT device_model_attributes_device_templates_fk FOREIGN KEY (device_template_id) REFERENCES public.device_templates(id) ON DELETE CASCADE -- Foreign key constraint
);

-- Column comments

COMMENT ON COLUMN public.device_model_attributes.id IS 'ID';
COMMENT ON COLUMN public.device_model_attributes.device_template_id IS 'Device template ID';
COMMENT ON COLUMN public.device_model_attributes.data_name IS 'Data name';
COMMENT ON COLUMN public.device_model_attributes.data_identifier IS 'Data identifier';
COMMENT ON COLUMN public.device_model_attributes.read_write_flag IS 'Read-write flag: R = Read, W = Write, RW = Read/Write';
COMMENT ON COLUMN public.device_model_attributes.data_type IS 'Data type: String, Number, Boolean, Enum';
COMMENT ON COLUMN public.device_model_attributes.unit IS 'Unit';
COMMENT ON COLUMN public.device_model_attributes.description IS 'Description';
COMMENT ON COLUMN public.device_model_attributes.additional_info IS 'Additional information';
COMMENT ON COLUMN public.device_model_attributes.created_at IS 'Creation time';
COMMENT ON COLUMN public.device_model_attributes.updated_at IS 'Update time';
COMMENT ON COLUMN public.device_model_attributes.remark IS 'Remarks';


-- public.device_model_commands definition

-- Drop table

-- DROP TABLE public.device_model_commands;

CREATE TABLE public.device_model_commands (
	id varchar(36) NOT NULL, -- ID
	device_template_id varchar(36) NOT NULL, -- Device template ID
	data_name varchar(255) NULL, -- Data name
	data_identifier varchar(255) NOT NULL, -- Data identifier
	params json NULL, -- Parameters
	description varchar(255) NULL, -- Description
	additional_info json NULL, -- Additional information
	created_at timestamptz(6) NOT NULL, -- Creation time
	updated_at timestamptz(6) NOT NULL, -- Update time
	remark varchar(255) NULL, -- Remarks
	tenant_id varchar(36) NOT NULL, -- Tenant ID
	CONSTRAINT device_model_commands_unique UNIQUE (data_identifier, device_template_id), -- Unique constraint for data identifier and device template ID
	CONSTRAINT device_model_telemetry_copy1_pkey2 PRIMARY KEY (id), -- Primary key constraint
	CONSTRAINT device_model_commands_device_templates_fk FOREIGN KEY (device_template_id) REFERENCES public.device_templates(id) ON DELETE CASCADE -- Foreign key constraint
);

-- Column comments

COMMENT ON COLUMN public.device_model_commands.id IS 'ID';
COMMENT ON COLUMN public.device_model_commands.device_template_id IS 'Device template ID';
COMMENT ON COLUMN public.device_model_commands.data_name IS 'Data name';
COMMENT ON COLUMN public.device_model_commands.data_identifier IS 'Data identifier';
COMMENT ON COLUMN public.device_model_commands.params IS 'Parameters';
COMMENT ON COLUMN public.device_model_commands.description IS 'Description';
COMMENT ON COLUMN public.device_model_commands.additional_info IS 'Additional information';
COMMENT ON COLUMN public.device_model_commands.created_at IS 'Creation time';
COMMENT ON COLUMN public.device_model_commands.updated_at IS 'Update time';
COMMENT ON COLUMN public.device_model_commands.remark IS 'Remarks';


-- public.device_model_events definition

-- Drop table

-- DROP TABLE public.device_model_events;

CREATE TABLE public.device_model_events (
	id varchar(36) NOT NULL, -- ID
	device_template_id varchar(36) NOT NULL, -- Device template ID
	data_name varchar(255) NULL, -- Data name
	data_identifier varchar(255) NOT NULL, -- Data identifier
	params json NULL, -- Parameters
	description varchar(255) NULL, -- Description
	additional_info json NULL, -- Additional information
	created_at timestamptz(6) NOT NULL, -- Creation time
	updated_at timestamptz(6) NOT NULL, -- Update time
	remark varchar(255) NULL, -- Remarks
	tenant_id varchar(36) NOT NULL, -- Tenant ID
	CONSTRAINT device_model_events_unique UNIQUE (device_template_id, data_identifier), -- Unique constraint for device template ID and data identifier
	CONSTRAINT device_model_telemetry_copy1_pkey1 PRIMARY KEY (id), -- Primary key constraint
	CONSTRAINT device_model_events_device_templates_fk FOREIGN KEY (device_template_id) REFERENCES public.device_templates(id) ON DELETE CASCADE -- Foreign key constraint
);

-- Column comments

COMMENT ON COLUMN public.device_model_events.id IS 'ID';
COMMENT ON COLUMN public.device_model_events.device_template_id IS 'Device template ID';
COMMENT ON COLUMN public.device_model_events.data_name IS 'Data name';
COMMENT ON COLUMN public.device_model_events.data_identifier IS 'Data identifier';
COMMENT ON COLUMN public.device_model_events.params IS 'Parameters';
COMMENT ON COLUMN public.device_model_events.description IS 'Description';
COMMENT ON COLUMN public.device_model_events.additional_info IS 'Additional information';
COMMENT ON COLUMN public.device_model_events.created_at IS 'Creation time';
COMMENT ON COLUMN public.device_model_events.updated_at IS 'Update time';
COMMENT ON COLUMN public.device_model_events.remark IS 'Remarks';


-- public.device_model_telemetry definition

-- Drop table

-- DROP TABLE public.device_model_telemetry;

-- public.device_model_telemetry definition

CREATE TABLE public.device_model_telemetry (
	id varchar(36) NOT NULL, -- id
	device_template_id varchar(36) NOT NULL, -- Device Template ID
	data_name varchar(255) NULL, -- Data Name
	data_identifier varchar(255) NOT NULL, -- Data Identifier
	read_write_flag varchar(10) NULL, -- Read/Write Flag: R-Read W-Write RW-Read/Write
	data_type varchar(50) NULL, -- Data Type: String, Number, Boolean
	unit varchar(50) NULL, -- Unit
	description varchar(255) NULL, -- Description
	additional_info json NULL, -- Additional Information
	created_at timestamptz(6) NOT NULL, -- Creation Time
	updated_at timestamptz(6) NOT NULL, -- Update Time
	remark varchar(255) NULL, -- Remark
	tenant_id varchar(36) NOT NULL, -- Tenant ID
	CONSTRAINT device_model_telemetry_pkey PRIMARY KEY (id),
	CONSTRAINT device_model_telemetry_unique UNIQUE (device_template_id, data_identifier),
	CONSTRAINT device_model_telemetry_device_templates_fk FOREIGN KEY (device_template_id) REFERENCES public.device_templates(id) ON DELETE CASCADE
);

-- Column comments

COMMENT ON COLUMN public.device_model_telemetry.id IS 'id';
COMMENT ON COLUMN public.device_model_telemetry.device_template_id IS 'Device Template ID';
COMMENT ON COLUMN public.device_model_telemetry.data_name IS 'Data Name';
COMMENT ON COLUMN public.device_model_telemetry.data_identifier IS 'Data Identifier';
COMMENT ON COLUMN public.device_model_telemetry.read_write_flag IS 'Read/Write Flag: R-Read W-Write RW-Read/Write';
COMMENT ON COLUMN public.device_model_telemetry.data_type IS 'Data Type: String, Number, Boolean';
COMMENT ON COLUMN public.device_model_telemetry.unit IS 'Unit';
COMMENT ON COLUMN public.device_model_telemetry.description IS 'Description';
COMMENT ON COLUMN public.device_model_telemetry.additional_info IS 'Additional Information';
COMMENT ON COLUMN public.device_model_telemetry.created_at IS 'Creation Time';
COMMENT ON COLUMN public.device_model_telemetry.updated_at IS 'Update Time';
COMMENT ON COLUMN public.device_model_telemetry.remark IS 'Remark';

-- public.device_trigger_condition definition

CREATE TABLE public.device_trigger_condition (
	id varchar(36) NOT NULL, -- Id
	scene_automation_id varchar(36) NOT NULL, -- Scene Automation ID (Foreign Key - Delete Cascade)
	enabled varchar(10) NOT NULL, -- Is Enabled
	group_id varchar(36) NOT NULL, -- UUID
	trigger_condition_type varchar(10) NOT NULL, -- Condition Type: 10: Device Type - Single Device, 11: Device Type - Device Category, 2: Time Range
	trigger_source varchar(36) NULL, -- Trigger Source: Possible values depend on Condition Type
	trigger_param_type varchar(10) NULL, -- Telemetry (TEL), Attribute (ATTR), Event (EVT), Status (STATUS)
	trigger_param varchar(50) NULL, -- Trigger Parameter (e.g., temperature)
	trigger_operator varchar(10) NULL, -- Operator: =, !=, >, <, >=, <=, between, in
	trigger_value varchar(99) NOT NULL, -- Trigger Value (e.g., 2-6 for a range)
	remark varchar(255) NULL, -- Remark
	tenant_id varchar(36) NOT NULL, -- Tenant ID
	CONSTRAINT device_trigger_condition_pkey PRIMARY KEY (id),
	CONSTRAINT fk_scene_automation_id FOREIGN KEY (scene_automation_id) REFERENCES public.scene_automations(id) ON DELETE CASCADE
);

-- Column comments

COMMENT ON COLUMN public.device_trigger_condition.id IS 'Id';
COMMENT ON COLUMN public.device_trigger_condition.scene_automation_id IS 'Scene Automation ID (Foreign Key - Delete Cascade)';
COMMENT ON COLUMN public.device_trigger_condition.enabled IS 'Is Enabled';
COMMENT ON COLUMN public.device_trigger_condition.group_id IS 'UUID';
COMMENT ON COLUMN public.device_trigger_condition.trigger_condition_type IS 'Condition Type: 10: Device Type - Single Device, 11: Device Type - Device Category, 2: Time Range';
COMMENT ON COLUMN public.device_trigger_condition.trigger_source IS 'Trigger Source: Depends on Condition Type';
COMMENT ON COLUMN public.device_trigger_condition.trigger_param_type IS 'Telemetry (TEL), Attribute (ATTR), Event (EVT), Status (STATUS)';
COMMENT ON COLUMN public.device_trigger_condition.trigger_param IS 'Trigger Parameter (e.g., temperature)';
COMMENT ON COLUMN public.device_trigger_condition.trigger_operator IS 'Operator: =, !=, >, <, >=, <=, between, in';
COMMENT ON COLUMN public.device_trigger_condition.trigger_value IS 'Trigger Value (e.g., 2-6 for range)';
COMMENT ON COLUMN public.device_trigger_condition.tenant_id IS 'Tenant ID';

-- public.one_time_tasks definition

CREATE TABLE public.one_time_tasks (
	id varchar(36) NOT NULL, -- Id
	scene_automation_id varchar(36) NOT NULL, -- Scene Automation ID (Foreign Key - Delete Cascade)
	execution_time timestamptz(6) NOT NULL, -- Execution Time
	executing_state varchar(10) NOT NULL, -- Execution Status: NEX - Not Executed, EXE - Executed, EXP - Expired
	enabled varchar(10) NOT NULL, -- Is Enabled: Y - Enabled, N - Disabled
	remark varchar(255) NULL, -- Remark
	expiration_time int8 NOT NULL, -- Expiration Time (default greater than execution time by 5 minutes)
	CONSTRAINT one_time_tasks_pkey PRIMARY KEY (id),
	CONSTRAINT fk_scene_automation_id FOREIGN KEY (scene_automation_id) REFERENCES public.scene_automations(id) ON DELETE CASCADE
);

-- Column comments

COMMENT ON COLUMN public.one_time_tasks.scene_automation_id IS 'Scene Automation ID (Foreign Key - Delete Cascade)';
COMMENT ON COLUMN public.one_time_tasks.execution_time IS 'Execution Time';
COMMENT ON COLUMN public.one_time_tasks.executing_state IS 'Execution Status: NEX - Not Executed, EXE - Executed, EXP - Expired';
COMMENT ON COLUMN public.one_time_tasks.enabled IS 'Is Enabled: Y - Enabled, N - Disabled';
COMMENT ON COLUMN public.one_time_tasks.expiration_time IS 'Expiration Time (default greater than execution time by 5 minutes)';

-- public.periodic_tasks definition

CREATE TABLE public.periodic_tasks (
	id varchar(36) NOT NULL, -- Id
	scene_automation_id varchar(36) NOT NULL, -- Scene Automation ID (Foreign Key - Delete Cascade)
	task_type varchar(255) NOT NULL, -- Task Type: HOUR, DAY, WEEK, MONTH, CRON
	params varchar(50) NOT NULL, -- Parameters for task type
	execution_time timestamptz(6) NOT NULL, -- Execution Time
	enabled varchar(10) NOT NULL, -- Is Enabled: Y - Enabled, N - Disabled
	remark varchar(255) NULL, -- Remark
	expiration_time int8 NOT NULL, -- Expiration Time (default greater than execution time by 5 minutes)
	CONSTRAINT periodic_tasks_pkey PRIMARY KEY (id),
	CONSTRAINT scene_automation_id_fkey FOREIGN KEY (scene_automation_id) REFERENCES public.scene_automations(id) ON DELETE CASCADE
);

-- Column comments

COMMENT ON COLUMN public.periodic_tasks.scene_automation_id IS 'Scene Automation ID (Foreign Key - Delete Cascade)';
COMMENT ON COLUMN public.periodic_tasks.task_type IS 'Task Type: HOUR, DAY, WEEK, MONTH, CRON';
COMMENT ON COLUMN public.periodic_tasks.execution_time IS 'Execution Time';
COMMENT ON COLUMN public.periodic_tasks.enabled IS 'Is Enabled: Y - Enabled, N - Disabled';
COMMENT ON COLUMN public.periodic_tasks.expiration_time IS 'Expiration Time (default greater than execution time by 5 minutes)';

-- public.products definition

CREATE TABLE public.products (
	id varchar(36) NOT NULL, -- UUID
	"name" varchar(255) NOT NULL, -- Product Name
	description varchar(255) NULL, -- Description
	product_type varchar(36) NULL, -- Product Type
	product_key varchar(255) NULL, -- Product Key
	product_model varchar(100) NULL, -- Product Model (Code)
	image_url varchar(500) NULL, -- Image URL
	created_at timestamptz(6) NOT NULL, -- Creation Time
	remark varchar(500) NULL, -- Remark
	additional_info json NULL, -- Additional Information
	tenant_id varchar(36) NULL, -- Tenant ID
	device_config_id varchar(36) NULL, -- Device Config ID
	CONSTRAINT products_pkey PRIMARY KEY (id),
	CONSTRAINT products_device_configs_fk FOREIGN KEY (device_config_id) REFERENCES public.device_configs(id) ON DELETE RESTRICT
);

-- Column comments

COMMENT ON COLUMN public.products.id IS 'UUID';
COMMENT ON COLUMN public.products."name" IS 'Product Name';
COMMENT ON COLUMN public.products.description IS 'Description';
COMMENT ON COLUMN public.products.product_type IS 'Product Type';
COMMENT ON COLUMN public.products.product_key IS 'Product Key';
COMMENT ON COLUMN public.products.product_model IS 'Product Model (Code)';
COMMENT ON COLUMN public.products.image_url IS 'Image URL';
COMMENT ON COLUMN public.products.created_at IS 'Creation Time';
COMMENT ON COLUMN public.products.tenant_id IS 'Tenant ID';

-- public.scene_action_info definition

-- Drop table

-- DROP TABLE public.scene_action_info;

-- Table: public.scene_action_info

CREATE TABLE public.scene_action_info (
    id varchar(36) NOT NULL,
    scene_id varchar(36) NOT NULL, -- Scene ID (cascade delete)
    action_target varchar(36) NOT NULL, -- Action target ID (device ID, device config ID, scene ID, alarm ID)
    action_type varchar(10) NOT NULL, -- Action type 
    -- 10: Single device
    -- 11: Device group
    -- 20: Activate scene
    -- 30: Trigger alarm
    -- 40: Service
    action_param_type varchar(10) NULL, -- Parameter type: 
    -- TEL: Telemetry
    -- ATTR: Attribute
    -- CMD: Command
    action_param varchar(10) NULL, -- Action parameter
    action_value varchar(255) NULL, -- Target value
    created_at timestamptz(6) NOT NULL, -- Creation time
    updated_at timestamptz(6) NULL, -- Update time
    tenant_id varchar(36) NOT NULL,
    remark varchar(255) NULL,
    CONSTRAINT scene_action_info_pkey PRIMARY KEY (id),
    CONSTRAINT scene_action_info_scene_id_fkey FOREIGN KEY (scene_id) REFERENCES public.scene_info(id) ON DELETE CASCADE
);

-- Column comments

COMMENT ON COLUMN public.scene_action_info.scene_id IS 'Scene ID (cascade delete)';
COMMENT ON COLUMN public.scene_action_info.action_target IS 'Action target ID (device ID, device config ID, scene ID, alarm ID)';
COMMENT ON COLUMN public.scene_action_info.action_type IS 'Action type: 10 - Single device, 11 - Device group, 20 - Activate scene, 30 - Trigger alarm, 40 - Service';
COMMENT ON COLUMN public.scene_action_info.action_param_type IS 'Parameter type: TEL - Telemetry, ATTR - Attribute, CMD - Command';
COMMENT ON COLUMN public.scene_action_info.action_param IS 'Action parameter';
COMMENT ON COLUMN public.scene_action_info.action_value IS 'Target value';
COMMENT ON COLUMN public.scene_action_info.created_at IS 'Creation time';
COMMENT ON COLUMN public.scene_action_info.updated_at IS 'Update time';

-- Table: public.scene_automation_log

CREATE TABLE public.scene_automation_log (
    scene_automation_id varchar(36) NOT NULL, -- Scene automation ID (foreign key - cascade delete)
    executed_at timestamptz(6) NOT NULL, -- Execution time
    detail text NOT NULL, -- Execution details: Detailed execution process
    execution_result varchar(10) NOT NULL, -- Execution result: S - Success, F - Failure, success only if all tasks succeed
    tenant_id varchar(36) NOT NULL,
    remark varchar(255) NULL,
    CONSTRAINT scene_automation_log_scene_automation_id_fkey FOREIGN KEY (scene_automation_id) REFERENCES public.scene_automations(id) ON DELETE RESTRICT
);

-- Column comments

COMMENT ON COLUMN public.scene_automation_log.scene_automation_id IS 'Scene automation ID (foreign key - cascade delete)';
COMMENT ON COLUMN public.scene_automation_log.executed_at IS 'Execution time';
COMMENT ON COLUMN public.scene_automation_log.detail IS 'Execution details: Detailed execution process';
COMMENT ON COLUMN public.scene_automation_log.execution_result IS 'Execution result: S - Success, F - Failure, success only if all tasks succeed';

-- Table: public.scene_log

CREATE TABLE public.scene_log (
    scene_id varchar(36) NOT NULL, -- Scene ID (cascade delete)
    executed_at timestamptz(6) NOT NULL, -- Execution time
    detail text NOT NULL, -- Execution details: Detailed execution process
    execution_result varchar(10) NOT NULL, -- Execution status: S - Success, F - Failure, success only if all tasks succeed
    tenant_id varchar(36) NOT NULL,
    remark varchar(255) NULL,
    id varchar(36) NOT NULL,
    CONSTRAINT scene_log_pkey PRIMARY KEY (id),
    CONSTRAINT scene_log_scene_id_fkey FOREIGN KEY (scene_id) REFERENCES public.scene_info(id) ON DELETE CASCADE
);

-- Column comments

COMMENT ON COLUMN public.scene_log.scene_id IS 'Scene ID (cascade delete)';
COMMENT ON COLUMN public.scene_log.executed_at IS 'Execution time';
COMMENT ON COLUMN public.scene_log.detail IS 'Execution details: Detailed execution process';
COMMENT ON COLUMN public.scene_log.execution_result IS 'Execution status: S - Success, F - Failure, success only if all tasks succeed';

-- Table: public.sys_dict_language

CREATE TABLE public.sys_dict_language (
    id varchar(36) NOT NULL, -- Primary key ID
    dict_id varchar(36) NOT NULL, -- sys_dict.id
    language_code varchar(36) NOT NULL, -- Language code
    "translation" varchar(255) NOT NULL, -- Translation
    CONSTRAINT sys_dict_language_dict_id_language_code_key UNIQUE (dict_id, language_code),
    CONSTRAINT sys_dict_language_pkey PRIMARY KEY (id),
    CONSTRAINT sys_dict_language_dict_id_fkey FOREIGN KEY (dict_id) REFERENCES public.sys_dict(id) ON DELETE CASCADE
);

-- Column comments

COMMENT ON COLUMN public.sys_dict_language.id IS 'Primary key ID';
COMMENT ON COLUMN public.sys_dict_language.dict_id IS 'sys_dict.id';
COMMENT ON COLUMN public.sys_dict_language.language_code IS 'Language code';
COMMENT ON COLUMN public.sys_dict_language."translation" IS 'Translation';

-- Constraint comments

COMMENT ON CONSTRAINT sys_dict_language_dict_id_language_code_key ON public.sys_dict_language IS 'dict_id and language_code are unique';

-- Table: public.data_scripts

CREATE TABLE public.data_scripts (
    id varchar(36) NOT NULL, -- ID
    "name" varchar(99) NOT NULL, -- Name
    device_config_id varchar(36) NOT NULL, -- Device config ID (cascade delete)
    enable_flag varchar(9) NOT NULL, -- Enable flag: Y - Enabled, N - Disabled, default enabled
    "content" text NULL, -- Content
    script_type varchar(9) NOT NULL, -- Script type: A - Telemetry report pre-processing, B - Telemetry command pre-processing, C - Attribute report pre-processing, D - Attribute command pre-processing
    last_analog_input text NULL, -- Last analog input
    description varchar(255) NULL, -- Description
    created_at timestamptz(6) NULL, -- Creation time
    updated_at timestamptz(6) NULL, -- Update time
    remark varchar(255) NULL, -- Remarks
    CONSTRAINT data_scripts_pkey PRIMARY KEY (id),
    CONSTRAINT data_scripts_device_configs_fk FOREIGN KEY (device_config_id) REFERENCES public.device_configs(id) ON DELETE CASCADE
);

-- Column comments

COMMENT ON COLUMN public.data_scripts.id IS 'ID';
COMMENT ON COLUMN public.data_scripts."name" IS 'Name';
COMMENT ON COLUMN public.data_scripts.device_config_id IS 'Device config ID (cascade delete)';
COMMENT ON COLUMN public.data_scripts.enable_flag IS 'Enable flag: Y - Enabled, N - Disabled, default enabled';
COMMENT ON COLUMN public.data_scripts."content" IS 'Content';
COMMENT ON COLUMN public.data_scripts.script_type IS 'Script type: A - Telemetry report pre-processing, B - Telemetry command pre-processing, C - Attribute report pre-processing, D - Attribute command pre-processing';
COMMENT ON COLUMN public.data_scripts.last_analog_input IS 'Last analog input';
COMMENT ON COLUMN public.data_scripts.description IS 'Description';
COMMENT ON COLUMN public.data_scripts.created_at IS 'Creation time';
COMMENT ON COLUMN public.data_scripts.updated_at IS 'Update time';
COMMENT ON COLUMN public.data_scripts.remark IS 'Remarks';


-- public.devices definition

-- Drop table

-- DROP TABLE public.devices;

CREATE TABLE public.devices (
	id varchar(36) NOT NULL, -- Id
	"name" varchar(255) NULL, -- Device Name
	voucher varchar(500) NOT NULL DEFAULT ''::character varying, -- Voucher
	tenant_id varchar(36) NOT NULL DEFAULT ''::character varying, -- Tenant ID, foreign key, restrict deletion
	is_enabled varchar(36) NOT NULL DEFAULT ''::character varying, -- Enabled/Disabled: enabled - enabled, disabled - disabled, default disabled, enabled by default after activation
	activate_flag varchar(36) NOT NULL DEFAULT ''::character varying, -- Activation flag: inactive - not activated, active - activated
	created_at timestamptz(6) NULL, -- Creation Time
	update_at timestamptz(6) NULL, -- Update Time
	device_number varchar(36) NOT NULL DEFAULT ''::character varying, -- Device Number, if not provided defaults to the token value
	product_id varchar(36) NULL, -- Product ID, foreign key, restrict deletion
	parent_id varchar(36) NULL, -- Gateway ID of the sub-device
	protocol varchar(36) NULL, -- Communication Protocol
	"label" varchar(255) NULL, -- Label, single label, separated by commas
	"location" varchar(100) NULL, -- Location
	sub_device_addr varchar(36) NULL, -- Sub-device Address
	current_version varchar(36) NULL, -- Current Firmware Version
	additional_info json NULL DEFAULT '{}'::json, -- Additional Information, e.g., thresholds, images
	protocol_config json NULL DEFAULT '{}'::json, -- Protocol Form Configuration
	remark1 varchar(255) NULL,
	remark2 varchar(255) NULL,
	remark3 varchar(255) NULL,
	device_config_id varchar(36) NULL, -- Device Configuration ID (foreign key)
	batch_number varchar(500) NULL, -- Batch Number
	activate_at timestamptz(6) NULL, -- Activation Date
	is_online int2 NOT NULL DEFAULT 0, -- Online Status: 1 - Online, 0 - Offline
	access_way varchar(10) NULL, -- Access Method: A - via protocol, B - via service
	description varchar(500) NULL, -- Description
	CONSTRAINT device_pkey PRIMARY KEY (id),
	CONSTRAINT devices_unique UNIQUE (device_number),
	CONSTRAINT devices_unique_1 UNIQUE (voucher),
	CONSTRAINT fk_device_config_id FOREIGN KEY (device_config_id) REFERENCES public.device_configs(id) ON DELETE RESTRICT,
	CONSTRAINT fk_product_id FOREIGN KEY (product_id) REFERENCES public.products(id) ON DELETE RESTRICT
);

-- Column comments

COMMENT ON COLUMN public.devices.id IS 'Id';
COMMENT ON COLUMN public.devices."name" IS 'Device Name';
COMMENT ON COLUMN public.devices.voucher IS 'Voucher';
COMMENT ON COLUMN public.devices.tenant_id IS 'Tenant ID, foreign key, restrict deletion';
COMMENT ON COLUMN public.devices.is_enabled IS 'Enabled/Disabled: enabled - enabled, disabled - disabled, default disabled, enabled by default after activation';
COMMENT ON COLUMN public.devices.activate_flag IS 'Activation flag: inactive - not activated, active - activated';
COMMENT ON COLUMN public.devices.created_at IS 'Creation Time';
COMMENT ON COLUMN public.devices.update_at IS 'Update Time';
COMMENT ON COLUMN public.devices.device_number IS 'Device Number, if not provided defaults to the token value';
COMMENT ON COLUMN public.devices.product_id IS 'Product ID, foreign key, restrict deletion';
COMMENT ON COLUMN public.devices.parent_id IS 'Gateway ID of the sub-device';
COMMENT ON COLUMN public.devices.protocol IS 'Communication Protocol';
COMMENT ON COLUMN public.devices."label" IS 'Label, single label, separated by commas';
COMMENT ON COLUMN public.devices."location" IS 'Location';
COMMENT ON COLUMN public.devices.sub_device_addr IS 'Sub-device Address';
COMMENT ON COLUMN public.devices.current_version IS 'Current Firmware Version';
COMMENT ON COLUMN public.devices.additional_info IS 'Additional Information, e.g., thresholds, images';
COMMENT ON COLUMN public.devices.protocol_config IS 'Protocol Form Configuration';
COMMENT ON COLUMN public.devices.device_config_id IS 'Device Configuration ID (foreign key)';
COMMENT ON COLUMN public.devices.batch_number IS 'Batch Number';
COMMENT ON COLUMN public.devices.activate_at IS 'Activation Date';
COMMENT ON COLUMN public.devices.is_online IS 'Online Status: 1 - Online, 0 - Offline';
COMMENT ON COLUMN public.devices.access_way IS 'Access Method: A - via protocol, B - via service';
COMMENT ON COLUMN public.devices.description IS 'Description';

-- public.event_datas definition

-- Drop table

-- DROP TABLE public.event_datas;

CREATE TABLE public.event_datas (
	id varchar(36) NOT NULL,
	device_id varchar(36) NOT NULL, -- Device ID (foreign key - restrict deletion)
	identify varchar(255) NOT NULL, -- Data Identifier
	ts timestamptz(6) NOT NULL, -- Reporting Time
	"data" json NULL, -- Data
	tenant_id varchar(36) NULL,
	CONSTRAINT event_datas_pkey PRIMARY KEY (id),
	CONSTRAINT event_datas_device_id_fkey FOREIGN KEY (device_id) REFERENCES public.devices(id) ON DELETE RESTRICT
);

-- Column comments

COMMENT ON COLUMN public.event_datas.device_id IS 'Device ID (foreign key - restrict deletion)';
COMMENT ON COLUMN public.event_datas.identify IS 'Data Identifier';
COMMENT ON COLUMN public.event_datas.ts IS 'Reporting Time';
COMMENT ON COLUMN public.event_datas."data" IS 'Data';

-- public.ota_upgrade_task_details definition

-- Drop table

-- DROP TABLE public.ota_upgrade_task_details;

CREATE TABLE public.ota_upgrade_task_details (
	id varchar(36) NOT NULL, -- Id
	ota_upgrade_task_id varchar(200) NOT NULL, -- Upgrade Task ID (foreign key - cascade deletion)
	device_id varchar(200) NOT NULL, -- Device ID (foreign key - restrict deletion)
	steps int2 NULL, -- Upgrade Progress 1-100
	status int2 NOT NULL, -- Status: 1 - Pending, 2 - Pushed, 3 - Upgrading, 4 - Upgrade Success, 5 - Upgrade Failed, 6 - Canceled
	status_description varchar(500) NULL, -- Status Description
	updated_at timestamptz(6) NULL,
	remark varchar(255) NULL,
	CONSTRAINT ota_upgrade_task_details_pkey PRIMARY KEY (id),
	CONSTRAINT fk_ota_upgrade_tasks FOREIGN KEY (ota_upgrade_task_id) REFERENCES public.ota_upgrade_tasks(id) ON DELETE CASCADE,
	CONSTRAINT ota_upgrade_task_details_device_id_fkey FOREIGN KEY (device_id) REFERENCES public.devices(id) ON DELETE RESTRICT
);

-- Column comments

COMMENT ON COLUMN public.ota_upgrade_task_details.id IS 'Id';
COMMENT ON COLUMN public.ota_upgrade_task_details.ota_upgrade_task_id IS 'Upgrade Task ID (foreign key - cascade deletion)';
COMMENT ON COLUMN public.ota_upgrade_task_details.device_id IS 'Device ID (foreign key - restrict deletion)';
COMMENT ON COLUMN public.ota_upgrade_task_details.steps IS 'Upgrade Progress 1-100';
COMMENT ON COLUMN public.ota_upgrade_task_details.status IS 'Status: 1 - Pending, 2 - Pushed, 3 - Upgrading, 4 - Upgrade Success, 5 - Upgrade Failed, 6 - Canceled';
COMMENT ON COLUMN public.ota_upgrade_task_details.status_description IS 'Status Description';


-- public.r_group_device definition

-- Drop table

-- DROP TABLE public.r_group_device;

CREATE TABLE public.r_group_device (
	group_id varchar(36) NOT NULL,
	device_id varchar(36) NOT NULL,
	tenant_id varchar(36) NOT NULL,
	CONSTRAINT r_group_device_group_id_device_id_key UNIQUE (group_id, device_id),
	CONSTRAINT fk_group_device FOREIGN KEY (group_id) REFERENCES public."groups"(id) ON DELETE CASCADE,
	CONSTRAINT fk_group_device_2 FOREIGN KEY (device_id) REFERENCES public.devices(id) ON DELETE CASCADE
);


-- public.telemetry_set_logs definition

-- Drop table

-- DROP TABLE public.telemetry_set_logs;

CREATE TABLE public.telemetry_set_logs (
	id varchar(36) NOT NULL,
	device_id varchar(36) NOT NULL, -- Device ID (foreign key - cascade delete)
	operation_type varchar(255) NULL, -- Operation type: 1 - Manual operation, 2 - Auto trigger
	"data" json NULL, -- Sent content
	status varchar(2) NULL, -- Status: 1 - Success, 2 - Failure
	error_message varchar(500) NULL, -- Error message
	created_at timestamptz(6) NOT NULL, -- Creation time
	user_id varchar(36) NULL, -- Operation user
	description varchar(255) NULL, -- Description
	CONSTRAINT telemetry_set_logs_pkey PRIMARY KEY (id),
	CONSTRAINT telemetry_set_logs_device_id_fkey FOREIGN KEY (device_id) REFERENCES public.devices(id) ON DELETE RESTRICT
);

-- Column comments

COMMENT ON COLUMN public.telemetry_set_logs.device_id IS 'Device ID (foreign key - cascade delete)';
COMMENT ON COLUMN public.telemetry_set_logs.operation_type IS 'Operation type: 1 - Manual operation, 2 - Auto trigger';
COMMENT ON COLUMN public.telemetry_set_logs."data" IS 'Sent content';
COMMENT ON COLUMN public.telemetry_set_logs.status IS 'Status: 1 - Success, 2 - Failure';
COMMENT ON COLUMN public.telemetry_set_logs.error_message IS 'Error message';
COMMENT ON COLUMN public.telemetry_set_logs.created_at IS 'Creation time';
COMMENT ON COLUMN public.telemetry_set_logs.user_id IS 'Operation user';
COMMENT ON COLUMN public.telemetry_set_logs.description IS 'Description';


-- public.attribute_datas definition

-- Drop table

-- DROP TABLE public.attribute_datas;

CREATE TABLE public.attribute_datas (
	id varchar(36) NOT NULL,
	device_id varchar(36) NOT NULL, -- Device ID (foreign key - cascade delete)
	"key" varchar(255) NOT NULL, -- Data identifier
	ts timestamptz(6) NOT NULL, -- Report time
	bool_v bool NULL,
	number_v float8 NULL,
	string_v text NULL,
	tenant_id varchar(36) NULL,
	CONSTRAINT attribute_datas_device_id_key_key UNIQUE (device_id, key),
	CONSTRAINT attribute_datas_device_id_fkey FOREIGN KEY (device_id) REFERENCES public.devices(id) ON DELETE RESTRICT
);

-- Column comments

COMMENT ON COLUMN public.attribute_datas.device_id IS 'Device ID (foreign key - cascade delete)';
COMMENT ON COLUMN public.attribute_datas."key" IS 'Data identifier';
COMMENT ON COLUMN public.attribute_datas.ts IS 'Report time';


-- public.attribute_set_logs definition

-- Drop table

-- DROP TABLE public.attribute_set_logs;

CREATE TABLE public.attribute_set_logs (
	id varchar(36) NOT NULL,
	device_id varchar(36) NOT NULL, -- Device ID (foreign key - cascade delete)
	operation_type varchar(255) NULL, -- Operation type: 1 - Manual operation, 2 - Auto trigger
	message_id varchar(36) NULL, -- Message ID
	"data" text NULL, -- Sent content
	rsp_data text NULL, -- Response content
	status varchar(2) NULL, -- Status: 1 - Success, 2 - Failure
	error_message varchar(500) NULL, -- Error message
	created_at timestamptz(6) NOT NULL, -- Creation time
	user_id varchar(36) NULL, -- Operation user
	description varchar(255) NULL, -- Description
	CONSTRAINT attribute_set_logs_pkey PRIMARY KEY (id),
	CONSTRAINT attribute_set_logs_device_id_fkey FOREIGN KEY (device_id) REFERENCES public.devices(id) ON DELETE RESTRICT
);

-- Column comments

COMMENT ON COLUMN public.attribute_set_logs.device_id IS 'Device ID (foreign key - cascade delete)';
COMMENT ON COLUMN public.attribute_set_logs.operation_type IS 'Operation type: 1 - Manual operation, 2 - Auto trigger';
COMMENT ON COLUMN public.attribute_set_logs.message_id IS 'Message ID';
COMMENT ON COLUMN public.attribute_set_logs."data" IS 'Sent content';
COMMENT ON COLUMN public.attribute_set_logs.rsp_data IS 'Response content';
COMMENT ON COLUMN public.attribute_set_logs.status IS 'Status: 1 - Success, 2 - Failure';
COMMENT ON COLUMN public.attribute_set_logs.error_message IS 'Error message';
COMMENT ON COLUMN public.attribute_set_logs.created_at IS 'Creation time';
COMMENT ON COLUMN public.attribute_set_logs.user_id IS 'Operation user';
COMMENT ON COLUMN public.attribute_set_logs.description IS 'Description';


-- public.command_set_logs definition

-- Drop table

-- DROP TABLE public.command_set_logs;

CREATE TABLE public.command_set_logs (
	id varchar(36) NOT NULL,
	device_id varchar(36) NOT NULL, -- Device ID (foreign key - cascade delete)
	operation_type varchar(255) NULL, -- Operation type: 1 - Manual operation, 2 - Auto trigger
	message_id varchar(36) NULL, -- Message ID
	"data" text NULL, -- Sent content
	rsp_data text NULL, -- Response content
	status varchar(2) NULL, -- Status: 1 - Success, 2 - Failure
	error_message varchar(500) NULL, -- Error message
	created_at timestamptz(6) NOT NULL, -- Creation time
	user_id varchar(36) NULL, -- Operation user
	description varchar(255) NULL, -- Description
	identify varchar(255) NULL, -- Data identifier
	CONSTRAINT command_set_logs_pkey PRIMARY KEY (id),
	CONSTRAINT command_set_logs_device_id_fkey FOREIGN KEY (device_id) REFERENCES public.devices(id) ON DELETE RESTRICT
);
COMMENT ON TABLE public.command_set_logs IS 'Command execution log';

-- Column comments

COMMENT ON COLUMN public.command_set_logs.device_id IS 'Device ID (foreign key - cascade delete)';
COMMENT ON COLUMN public.command_set_logs.operation_type IS 'Operation type: 1 - Manual operation, 2 - Auto trigger';
COMMENT ON COLUMN public.command_set_logs.message_id IS 'Message ID';
COMMENT ON COLUMN public.command_set_logs."data" IS 'Sent content';
COMMENT ON COLUMN public.command_set_logs.rsp_data IS 'Response content';
COMMENT ON COLUMN public.command_set_logs.status IS 'Status: 1 - Success, 2 - Failure';
COMMENT ON COLUMN public.command_set_logs.error_message IS 'Error message';
COMMENT ON COLUMN public.command_set_logs.created_at IS 'Creation time';
COMMENT ON COLUMN public.command_set_logs.user_id IS 'Operation user';
COMMENT ON COLUMN public.command_set_logs.description IS 'Description';
COMMENT ON COLUMN public.command_set_logs.identify IS 'Data identifier';

-- Initialization SQL
INSERT INTO public.sys_dict (id, dict_code, dict_value, created_at, remark) 
VALUES('0013fb9e-e3be-95d4-9c96-f18d1f9ddfcd', 'GATEWAY_PROTOCOL', 'MQTT', '2024-01-18 15:39:38.469', NULL);
INSERT INTO public.sys_dict (id, dict_code, dict_value, created_at, remark) 
VALUES('7162fb9e-e3be-95d4-9c96-f18d1f9ddfcd', 'DIRECT_ATTACHED_PROTOCOL', 'MQTT', '2024-01-18 15:39:38.469', NULL);

INSERT INTO public.sys_dict_language (id, dict_id, language_code, "translation") 
VALUES('001c3960-3067-536d-5c97-7645351a687c', '7162fb9e-e3be-95d4-9c96-f18d1f9ddfcd', 'en', 'MQTT Protocol');
INSERT INTO public.sys_dict_language (id, dict_id, language_code, "translation") 
VALUES('002c3960-3067-536d-5c97-7645351a687b', '0013fb9e-e3be-95d4-9c96-f18d1f9ddfcd', 'en', 'MQTT Protocol (Gateway)');

INSERT INTO public.sys_function (id, "name", enable_flag, description, remark) 
VALUES('function_1', 'use_captcha', 'disable', 'Captcha Login', NULL);
INSERT INTO public.sys_function (id, "name", enable_flag, description, remark) 
VALUES('function_2', 'enable_reg', 'disable', 'Tenant Registration', NULL);

INSERT INTO public.users (id, "name", phone_number, email, status, authority, "password", tenant_id, remark, additional_info, created_at, updated_at) 
VALUES('00000000-4fe9-b409-67c3-000000000000', 'admin', '1231231321', 'demo-super@thingsly.vn', 'N', 'SYS_ADMIN', '$2a$10$XKA2iwqPTIJEXkn9aKlIQ.PmyAu8Ae1UhapHqj/ShkeCu0vInRGiG', 'aaaaaa', 'dolor', '{}'::json, NULL, '2024-03-06 14:52:52.390');
INSERT INTO public.users (id, "name", phone_number, email, status, authority, "password", tenant_id, remark, additional_info, created_at, updated_at) 
VALUES('11111111-4fe9-b409-67c3-111111111111', 'Tenant', '17366666666', 'demo-tenant@thingsly.vn', 'N', 'TENANT_ADMIN', '$2a$10$XKA2iwqPTIJEXkn9aKlIQ.PmyAu8Ae1UhapHqj/ShkeCu0vInRGiG', 'd616bcbb', '', '{}'::json, '2024-06-05 16:48:11.097', '2024-06-05 16:48:11.097');


INSERT INTO public.data_policy (id, data_type, retention_days, last_cleanup_time, last_cleanup_data_time, enabled, remark) VALUES('b', '2', 15, '2024-06-05 10:02:00.003', '2024-05-21 10:02:00.003', '1', '');
INSERT INTO public.data_policy (id, data_type, retention_days, last_cleanup_time, last_cleanup_data_time, enabled, remark) VALUES('a', '1', 15, '2024-06-05 10:02:00.003', '2024-05-21 10:02:00.101', '1', '');

INSERT INTO public.sys_ui_elements (id, parent_id, element_code, element_type, orders, param1, param2, param3, authority, description, created_at, remark, multilingual, route_path) 
VALUES('6e5e0963-46bf-bc27-d792-156e87a69f51', '6e5e0963-46bf-bc27-d792-156e87a69f51', 'alarm', 1, 115, '/alarm', 'simple-icons:antdesign', 'self', '["TENANT_ADMIN","SYS_ADMIN"]'::json, 'Alarm', '2024-03-07 21:46:40.055', '', 'route.alarm', 'layout.base');

INSERT INTO public.sys_ui_elements (id, parent_id, element_code, element_type, orders, param1, param2, param3, authority, description, created_at, remark, multilingual, route_path) 
VALUES('f9bd5f79-291e-26d2-1553-473c04b15ce4', 'e1ebd134-53df-3105-35f4-489fc674d173', 'management_setting', 3, 42, '/management/setting', 'uil:brightness-plus', 'self', '["SYS_ADMIN"]'::json, 'System Settings', '2024-02-18 17:52:08.236', '', 'route.management_setting', 'view.management_setting');

INSERT INTO public.sys_ui_elements (id, parent_id, element_code, element_type, orders, param1, param2, param3, authority, description, created_at, remark, multilingual, route_path) 
VALUES('51381989-1160-93cd-182e-d44a1c4ab89b', '676e8f33-875a-0473-e9ca-c82fd09fef57', 'automation_scene-manage', 3, 1142, '/automation/scene-manage', 'uil:brightness-plus', 'self', '["TENANT_ADMIN","SYS_ADMIN"]'::json, 'Scene Management', '2024-03-07 21:44:11.106', '', 'route.automation_scene-manage', 'view.automation_scene-manage');

INSERT INTO public.sys_ui_elements (id, parent_id, element_code, element_type, orders, param1, param2, param3, authority, description, created_at, remark, multilingual, route_path) 
VALUES('8f4e9058-e30d-2fb5-ac6d-784613234883', '6e5e0963-46bf-bc27-d792-156e87a69f51', 'alarm-information', 3, 1151, '/alarm/alarm-information', 'ph:alarm', 'basic', '["TENANT_ADMIN"]'::json, 'Alarm Information', '2024-03-07 21:47:22.817', '', 'default', '');

INSERT INTO public.sys_ui_elements (id, parent_id, element_code, element_type, orders, param1, param2, param3, authority, description, created_at, remark, multilingual, route_path) 
VALUES('b6d57a4a-d37a-9d9d-6e4e-be33b955ff04', '6e5e0963-46bf-bc27-d792-156e87a69f51', 'alarm_notification-group', 3, 1152, '/alarm/notification-group', 'simple-icons:apacheecharts', 'basic', '["TENANT_ADMIN"]'::json, 'Notification Group', '2024-03-07 21:48:15.416', '', 'route.alarm_notification-group', 'view.alarm_notification-group');

INSERT INTO public.sys_ui_elements (id, parent_id, element_code, element_type, orders, param1, param2, param3, authority, description, created_at, remark, multilingual, route_path) 
VALUES('faf7e607-00ae-3483-40a1-b74f9245b100', 'e1ebd134-53df-3105-35f4-489fc674d173', 'management_auth', 3, 43, '/management/auth', 'ic:baseline-security', 'self', '["SYS_ADMIN"]'::json, 'Menu Management', '2024-02-18 17:49:31.209', '', 'route.management_auth', 'view.management_auth');

INSERT INTO public.sys_ui_elements (id, parent_id, element_code, element_type, orders, param1, param2, param3, authority, description, created_at, remark, multilingual, route_path) 
VALUES('95e2a961-382b-f4a6-87b3-1898123c95bc', '0', 'visualization', 1, 113, '/visualization', 'icon-park-outline:data-server', 'self', '["TENANT_ADMIN","SYS_ADMIN"]'::json, 'Visualization', '2024-03-07 21:37:16.042', '', 'route.visualization', 'layout.base');

INSERT INTO public.sys_ui_elements (id, parent_id, element_code, element_type, orders, param1, param2, param3, authority, description, created_at, remark, multilingual, route_path) 
VALUES('676e8f33-875a-0473-e9ca-c82fd09fef57', '0', 'automation', 1, 114, '/automation', 'material-symbols:device-hub', 'self', '["TENANT_ADMIN"]'::json, 'Automation', '2024-03-07 21:41:17.921', '', 'route.automation', 'layout.base');

INSERT INTO public.sys_ui_elements (id, parent_id, element_code, element_type, orders, param1, param2, param3, authority, description, created_at, remark, multilingual, route_path) 
VALUES('e1ebd134-53df-3105-35f4-489fc674d173', '0', 'management', 1, 120, '/management', 'carbon:cloud-service-management', 'self', '["SYS_ADMIN","TENANT_ADMIN"]'::json, 'System Management', '2024-02-18 17:48:45.265', '', 'route.management', 'layout.base');

INSERT INTO public.sys_ui_elements (id, parent_id, element_code, element_type, orders, param1, param2, param3, authority, description, created_at, remark, multilingual, route_path) 
VALUES('96aa2fac-90b2-aca1-1ce0-51b5060f4081', '676e8f33-875a-0473-e9ca-c82fd09fef57', 'automation_linkage-edit', 3, 1143, '/automation/linkage-edit', '', '1', '["TENANT_ADMIN","SYS_ADMIN"]'::json, 'Scene Linkage Edit', '2024-03-15 01:36:03.938', '', 'route.automation_linkage-edit', 'view.automation_linkage-edit');


INSERT INTO public.sys_ui_elements (id, parent_id, element_code, element_type, orders, param1, param2, param3, authority, description, created_at, remark, multilingual, route_path) 
VALUES('e619f321-9823-b563-b24d-ecc16d7b23cc', '6e5e0963-46bf-bc27-d792-156e87a69f51', 'alarm_notification-record', 3, 1153, '/alarm/notification-record', 'mdi:monitor-dashboard', 'basic', '["TENANT_ADMIN"]'::json, 'Notification Record', '2024-03-07 21:48:56.415', '', 'route.alarm_notification-record', 'view.alarm_notification-record');

INSERT INTO public.sys_ui_elements (id, parent_id, element_code, element_type, orders, param1, param2, param3, authority, description, created_at, remark, multilingual, route_path) 
VALUES('01dab674-9556-cdd7-b800-78bcb366adb4', '676e8f33-875a-0473-e9ca-c82fd09fef57', 'automation_scene-linkage', 3, 1141, '/automation/scene-linkage', 'mdi:airplane-edit', 'self', '["TENANT_ADMIN","SYS_ADMIN"]'::json, 'Scene Linkage', '2024-03-07 21:43:33.920', '', 'route.automation_scene-linkage', 'view.automation_scene-linkage');

INSERT INTO public.sys_ui_elements (id, parent_id, element_code, element_type, orders, param1, param2, param3, authority, description, created_at, remark, multilingual, route_path) 
VALUES('c078182f-bf4b-b560-da97-02926fa98f78', '650bc444-7672-1123-1e41-7e37365b0186', 'alarm_notification-record', 3, 1, '/alarm/notification-record', 'icon-park-outline:editor', 'self', '["TENANT_ADMIN"]'::json, 'Notification Record', '2024-03-20 10:04:34.927', '', 'route.alarm_notification-record', 'view.alarm_notification-record');

INSERT INTO public.sys_ui_elements (id, parent_id, element_code, element_type, orders, param1, param2, param3, authority, description, created_at, remark, multilingual, route_path) 
VALUES('485c2a20-ebc5-2216-4871-26453470d290', '650bc444-7672-1123-1e41-7e37365b0186', 'alarm_warning-message', 3, 999, '/alarm/warning-message', 'mdi:airballoon', 'self', '["TENANT_ADMIN"]'::json, 'Warning Information', '2024-03-17 15:27:40.378', '', 'route.alarm_warning-message', 'view.alarm_warning-message');

INSERT INTO public.sys_ui_elements (id, parent_id, element_code, element_type, orders, param1, param2, param3, authority, description, created_at, remark, multilingual, route_path) 
VALUES('2f3ffd60-efec-aafb-a866-f1cb79f88390', 'e1ebd134-53df-3105-35f4-489fc674d173', 'system-management-user_system-log', 3, 1171, '/system-management-user/system-log', 'mdi:monitor-dashboard', 'basic', '["TENANT_ADMIN"]'::json, 'System Logs', '2024-03-07 22:23:08.576', '', 'route.system-management-user_system-log', 'view.system-management-user_system-log');

INSERT INTO public.sys_ui_elements (id, parent_id, element_code, element_type, orders, param1, param2, param3, authority, description, created_at, remark, multilingual, route_path) 
VALUES('e186a671-8e24-143a-5a2c-27a1f5f38bf3', '5373a6a2-1861-af35-eb4c-adfd5ca55ecd', 'device_config-edit', 3, 1128, '/device/config-edit', '', '1', '["TENANT_ADMIN"]'::json, 'Device Configuration Edit', '2024-03-11 21:49:34.952', '', 'route.device_config-edit', 'view.device_config-edit');

INSERT INTO public.sys_ui_elements (id, parent_id, element_code, element_type, orders, param1, param2, param3, authority, description, created_at, remark, multilingual, route_path) 
VALUES('a2c53126-029f-7138-4d7a-f45491f396da', '0', 'apply', 1, 3, '/apply', 'mdi:apps-box', '0', '["SYS_ADMIN"]'::json, 'Application Management', '2024-02-18 17:59:31.642', '', 'route.apply', 'layout.base');

INSERT INTO public.sys_ui_elements (id, parent_id, element_code, element_type, orders, param1, param2, param3, authority, description, created_at, remark, multilingual, route_path) 
VALUES('49857e46-2176-610e-98fc-892b4fde50f9', '5373a6a2-1861-af35-eb4c-adfd5ca55ecd', 'device_details', 3, 1124, '/device/details', 'mdi:monitor-dashboard', '1', '["TENANT_ADMIN"]'::json, 'Device Details', '2024-03-05 17:52:21.434', '', 'route.device_details', 'view.device_details');

INSERT INTO public.sys_ui_elements (id, parent_id, element_code, element_type, orders, param1, param2, param3, authority, description, created_at, remark, multilingual, route_path) 
VALUES('ed4a5cfa-03e7-ccc0-6cc8-bcadccd25541', '95e2a961-382b-f4a6-87b3-1898123c95bc', 'visualization_kanban-details', 3, 1132, '/visualization/kanban-details', 'ic:baseline-credit-card', '1', '["TENANT_ADMIN","SYS_ADMIN"]'::json, 'Kanban Details', '2024-03-12 10:14:50.152', '', 'route.visualization_panel-details', 'view.visualization_panel-details');

INSERT INTO public.sys_ui_elements (id, parent_id, element_code, element_type, orders, param1, param2, param3, authority, description, created_at, remark, multilingual, route_path) 
VALUES('502a0d6c-750e-92f6-a1a7-ffdd362dbbac', '95e2a961-382b-f4a6-87b3-1898123c95bc', 'visualization_panel-preview', 3, 1133, '/visualization/panel-preview', '', '1', '["TENANT_ADMIN","SYS_ADMIN"]'::json, 'Kanban Preview', '2024-03-12 10:16:29.336', '', 'route.visualization_panel-preview', 'view.visualization_panel-preview');

INSERT INTO public.sys_ui_elements (id, parent_id, element_code, element_type, orders, param1, param2, param3, authority, description, created_at, remark, multilingual, route_path) 
VALUES('75785418-a5af-d790-0783-e4ee4e42521e', '5373a6a2-1861-af35-eb4c-adfd5ca55ecd', 'device_grouping', 3, 1122, '/device/grouping', 'material-symbols:grid-on-outline-sharp', '0', '["TENANT_ADMIN"]'::json, 'Device Grouping', '2024-03-05 17:53:25.004', '', 'route.device_grouping', 'view.device_grouping');


INSERT INTO public.sys_ui_elements (id, parent_id, element_code, element_type, orders, param1, param2, param3, authority, description, created_at, remark, multilingual, route_path) 
VALUES ('8de46003-170c-a24d-6baf-84d1c7298aa3', '5373a6a2-1861-af35-eb4c-adfd5ca55ecd', 'device_grouping-details', 3, 1123, '/device/grouping-details', '', '1', '["TENANT_ADMIN"]'::json, 'Grouping Details', '2024-03-05 17:54:23.158', '', 'route.device_grouping-details', 'view.device_grouping-details');

INSERT INTO public.sys_ui_elements (id, parent_id, element_code, element_type, orders, param1, param2, param3, authority, description, created_at, remark, multilingual, route_path) 
VALUES ('5373a6a2-1861-af35-eb4c-adfd5ca55ecd', '0', 'device', 1, 112, '/device', 'icon-park-outline:workbench', '0', '["TENANT_ADMIN"]'::json, 'Device Access', '2024-03-05 17:51:19.298', '', 'route.device', 'layout.base');

INSERT INTO public.sys_ui_elements (id, parent_id, element_code, element_type, orders, param1, param2, param3, authority, description, created_at, remark, multilingual, route_path) 
VALUES ('7419e37e-c167-f12b-7ace-76e479144181', '5373a6a2-1861-af35-eb4c-adfd5ca55ecd', 'device_template', 3, 1127, '/device/template', 'simple-icons:apacheecharts', 'self', '["TENANT_ADMIN"]'::json, 'Function Template', '2024-03-05 18:01:29.826', 'Defines the physical model and displays charts', 'route.device_template', 'view.device_template');

INSERT INTO public.sys_ui_elements (id, parent_id, element_code, element_type, orders, param1, param2, param3, authority, description, created_at, remark, multilingual, route_path) 
VALUES ('774a716d-9861-bac9-857f-acaa25e7659f', '5373a6a2-1861-af35-eb4c-adfd5ca55ecd', 'device_config', 3, 1126, '/device/config', 'clarity:plugin-line', 'self', '["TENANT_ADMIN"]'::json, 'Configuration Template', '2024-03-05 22:06:53.842', 'Configuration for device protocols and other parameters', 'route.device_config', 'view.device_config');

INSERT INTO public.sys_ui_elements (id, parent_id, element_code, element_type, orders, param1, param2, param3, authority, description, created_at, remark, multilingual, route_path) 
VALUES ('36c4f5ce-3279-55f2-ede2-81b4a0bae24b', 'e1ebd134-53df-3105-35f4-489fc674d173', 'management_user', 3, 41, '/management/user', 'ic:round-manage-accounts', 'self', '["SYS_ADMIN"]'::json, 'Tenant Management', '2024-02-18 17:50:48.999', '', 'route.management_user', 'view.management_user');


INSERT INTO public.sys_ui_elements (id, parent_id, element_code, element_type, orders, param1, param2, param3, authority, description, created_at, remark, multilingual, route_path) 
VALUES ('c4dff952-3bf4-8102-6882-e9d3f3cffbda', '5373a6a2-1861-af35-eb4c-adfd5ca55ecd', 'device_manage', 3, 1121, '/device/manage', 'icon-park-outline:analysis', '0', '["TENANT_ADMIN"]'::json, 'Device Management', '2024-03-05 17:55:08.170', '', 'route.device_manage', 'view.device_manage');

INSERT INTO public.sys_ui_elements (id, parent_id, element_code, element_type, orders, param1, param2, param3, authority, description, created_at, remark, multilingual, route_path) 
VALUES ('fec91838-d30d-7d66-6715-0912f1b171d8', 'e1ebd134-53df-3105-35f4-489fc674d173', 'management_notification', 3, 44, '/management/notification', 'mdi:alert', 'self', '["SYS_ADMIN"]'::json, 'Notification Configuration', '2024-03-15 19:50:07.495', '', 'route.management_notification', 'view.management_notification');

INSERT INTO public.sys_ui_elements (id, parent_id, element_code, element_type, orders, param1, param2, param3, authority, description, created_at, remark, multilingual, route_path) 
VALUES ('82c46beb-9ec4-8a3d-c6e4-04ba426e525a', '650bc444-7672-1123-1e41-7e37365b0186', 'alarm_notification-group', 3, 1, '/alarm/notification-group', 'ic:round-supervisor-account', 'basic', '["TENANT_ADMIN"]'::json, 'Notification Group', '2024-03-20 10:03:19.955', '', 'route.alarm_notification-group', 'view.alarm_notification-group');

INSERT INTO public.sys_ui_elements (id, parent_id, element_code, element_type, orders, param1, param2, param3, authority, description, created_at, remark, multilingual, route_path) 
VALUES ('650bc444-7672-1123-1e41-7e37365b0186', '0', 'alarm', 1, 115, '/alarm', 'mdi:alert', 'self', '["TENANT_ADMIN"]'::json, 'Alarm', '2024-03-17 09:01:52.183', '', 'route.alarm', 'layout.base');

INSERT INTO public.sys_ui_elements (id, parent_id, element_code, element_type, orders, param1, param2, param3, authority, description, created_at, remark, multilingual, route_path) 
VALUES ('64f684f1-390c-b5f2-9994-36895025df8a', '676e8f33-875a-0473-e9ca-c82fd09fef57', 'automation_space-management', 3, 10, 'automation/space-management', 'ic:baseline-security', '1', '["TENANT_ADMIN","SYS_ADMIN"]'::json, 'Scene Management', '2024-03-22 13:25:38.820', '', 'default', 'view.automation space-management');

INSERT INTO public.sys_ui_elements (id, parent_id, element_code, element_type, orders, param1, param2, param3, authority, description, created_at, remark, multilingual, route_path) 
VALUES ('76bfc16e-ed22-bcc0-c688-d462666e8a8d', '0', 'personal-center', 3, 999, '/personal-center', 'carbon:user-role', '1', '["TENANT_ADMIN","SYS_ADMIN"]'::json, 'Personal Center', '2024-03-17 09:27:01.048', '', 'route.personal_center', 'layout.base$view.personal-center');

INSERT INTO public.sys_ui_elements (id, parent_id, element_code, element_type, orders, param1, param2, param3, authority, description, created_at, remark, multilingual, route_path) 
VALUES ('975c9550-5db9-7b4c-5dea-7a4c326a37ff', '676e8f33-875a-0473-e9ca-c82fd09fef57', 'automation_scene-edit', 3, 1, '/automation/scene-edit', 'mdi:apps-box', '1', '["TENANT_ADMIN"]'::json, 'Add Scene', '2024-04-04 10:50:43.219', '', 'route.automation_scene-edit', 'view.automation_scene-edit');


INSERT INTO public.sys_ui_elements (id, parent_id, element_code, element_type, orders, param1, param2, param3, authority, description, created_at, remark, multilingual, route_path) 
VALUES ('680cae76-6c50-90e6-c2f9-58d01389aa08', '9a11b3e4-9982-a0f0-996c-a9be6e738947', 'data-service_rule-engine', 3, 21, '/data-service/rule-engine', 'mdi:menu', '1', '["SYS_ADMIN"]'::json, 'Rule Engine', '2024-03-07 17:06:02.804', '', 'route.data-service_rule-engine', 'view.data-service_rule-engine');

INSERT INTO public.sys_ui_elements (id, parent_id, element_code, element_type, orders, param1, param2, param3, authority, description, created_at, remark, multilingual, route_path) 
VALUES ('a2654c98-3749-c88b-0472-b414049ca532', '95e2a961-382b-f4a6-87b3-1898123c95bc', 'route.visualization_kanban', 3, 1131, '/visualization/kanban', 'tabler:device-tv', 'self', '["TENANT_ADMIN","SYS_ADMIN"]'::json, 'Kanban', '2024-03-07 21:39:58.608', '', 'route.visualization_kanban', 'view.visualization_panel');

INSERT INTO public.sys_ui_elements (id, parent_id, element_code, element_type, orders, param1, param2, param3, authority, description, created_at, remark, multilingual, route_path) 
VALUES ('f960c45c-6d5b-e67a-c4ff-1f0e869c1625', '5373a6a2-1861-af35-eb4c-adfd5ca55ecd', 'device_service-details', 3, 1130, '/device/service-details', 'ph:align-bottom', '1', '["TENANT_ADMIN"]'::json, 'Service Details', '2024-07-01 23:16:56.668', '', 'route.device_service_details', '');

INSERT INTO public.sys_ui_elements (id, parent_id, element_code, element_type, orders, param1, param2, param3, authority, description, created_at, remark, multilingual, route_path) 
VALUES ('86cb08fa-8b08-3d99-4b3a-d6132ee93a0f', '5373a6a2-1861-af35-eb4c-adfd5ca55ecd', 'device_config-detail', 3, 1127, '/device/config-detail', 'icon-park-outline:data-server', '1', '["TENANT_ADMIN"]'::json, 'Device Configuration Details', '2024-03-10 11:13:25.253', '', 'route.device_config-detail', 'view.device_config-detail');

-- INSERT INTO public.sys_ui_elements (id, parent_id, element_code, element_type, orders, param1, param2, param3, authority, description, created_at, remark, multilingual, route_path) 
-- VALUES ('04c523e3-5d60-8554-fd7f-11303ac8b8fe', 'e1ebd134-53df-3105-35f4-489fc674d173', 'management_role', 3, 1, '/management/role', 'ic:round-supervisor-account', '0', '["TENANT_ADMIN"]'::json, 'Role Management', '2024-08-16 16:24:40.097', '', 'route.management_role', '');

INSERT INTO public.sys_ui_elements (id, parent_id, element_code, element_type, orders, param1, param2, param3, authority, description, created_at, remark, multilingual, route_path) 
VALUES ('075d9f19-5618-bb9b-6ccd-f382bfd3292b', '5373a6a2-1861-af35-eb4c-adfd5ca55ecd', 'device_service-access', 3, 1129, '/device/service-access', 'mdi:ab-testing', '0', '["TENANT_ADMIN"]'::json, 'Service Access Point Management', '2024-07-01 21:52:09.402', '', 'route.device_service_access', '');

-- INSERT INTO public.sys_ui_elements (id, parent_id, element_code, element_type, orders, param1, param2, param3, authority, description, created_at, remark, multilingual, route_path) 
-- VALUES ('8251d887-d63e-9529-f645-c5393a6afae0', 'e1ebd134-53df-3105-35f4-489fc674d173', 'management_ordinary-user', 3, 1788, '/management/ordinary-user', 'ph:android-logo', '0', '["TENANT_ADMIN"]'::json, 'User Management', '2024-08-27 14:48:08.191', '', '', '');

INSERT INTO public.sys_ui_elements (id, parent_id, element_code, element_type, orders, param1, param2, param3, authority, description, created_at, remark, multilingual, route_path) 
VALUES ('59612e2f-e297-acb7-fcf4-143bf6e66109', '5373a6a2-1861-af35-eb4c-adfd5ca55ecd', 'device_details-child', 3, 1124, '/device/details-child', '', '1', '["TENANT_ADMIN"]'::json, 'Child Device Details', '2024-05-10 20:33:34.869', '', 'route.device_details-child', 'view.device_details-child');

INSERT INTO public.sys_ui_elements (id, parent_id, element_code, element_type, orders, param1, param2, param3, authority, description, created_at, remark, multilingual, route_path) 
VALUES ('29a684f9-c2bb-1a6f-6045-314944bef580', 'a2c53126-029f-7138-4d7a-f45491f396da', 'plug_in', 3, 32, '/apply/plugin', 'mdi:emoticon', '0', '["SYS_ADMIN"]'::json, 'Plugin Management', '2024-06-29 01:04:51.301', '', 'route.apply_in', '');

INSERT INTO public.sys_ui_elements (id, parent_id, element_code, element_type, orders, param1, param2, param3, authority, description, created_at, remark, multilingual, route_path) 
VALUES ('a190f7a5-1501-3814-9dd1-f3e1fbe7265e', '0', 'home', 3, 0, '/home', 'mdi:alpha-f-box-outline', 'self', '["SYS_ADMIN","TENANT_ADMIN"]'::json, 'Home Page', '2024-02-26 16:07:20.202', 'home', 'route.home', 'layout.base$view.home');

INSERT INTO public.sys_ui_elements (id, parent_id, element_code, element_type, orders, param1, param2, param3, authority, description, created_at, remark, multilingual, route_path) 
VALUES ('9a11b3e4-9982-a0f0-996c-a9be6e738947', '0', 'data-service', 1, 2, '/data-service', 'mdi:monitor-dashboard', '1', '["SYS_ADMIN"]'::json, 'Data Service', '2024-03-07 17:05:04.101', '', 'route.data-service', 'layout.base');

INSERT INTO public.sys_ui_elements (id, parent_id, element_code, element_type, orders, param1, param2, param3, authority, description, created_at, remark, multilingual, route_path) 
VALUES ('3aaca04b-2a2e-dfca-9fb4-0b2819362783', '2cc0c5ba-f086-91e5-0b8c-ad0546b1f2a9', 'test_kan-ban-test', 3, 1, '/test/kan-ban-test', '', '1', '["SYS_ADMIN","TENANT_ADMIN"]'::json, 'Kanban Test', '2024-05-21 01:17:16.911', '', 'route.test_kan-ban-test', 'view.test_kan-ban-test');

INSERT INTO public.sys_ui_elements (id, parent_id, element_code, element_type, orders, param1, param2, param3, authority, description, created_at, remark, multilingual, route_path) 
VALUES ('2fe87d7c-627e-9ca3-94dd-6d0249853bd4', '990af72f-06ce-5f23-3af6-1694bd479c96', 'management_user', 3, 1, '/management/user', '', '0', '["SYS_ADMIN","TENANT_ADMIN"]'::json, 'Management User', '2024-09-04 10:04:34.658', '', 'default', '');

INSERT INTO public.sys_ui_elements (id, parent_id, element_code, element_type, orders, param1, param2, param3, authority, description, created_at, remark, multilingual, route_path) 
VALUES ('18892c6e-ca04-f2b5-c243-f2c7230b3f33', '990af72f-06ce-5f23-3af6-1694bd479c96', 'manage_user', 3, 1, '/manage/user', '', '0', '["TENANT_ADMIN","SYS_ADMIN"]'::json, 'User Management', '2024-09-04 10:05:06.377', '', 'default', '');

INSERT INTO public.logo (id, system_name, logo_cache, logo_background, logo_loading, home_background, remark) 
VALUES ('a', 'Thingsly', '', '', '', '', NULL);

ALTER TABLE "public"."scene_action_info"
ALTER COLUMN "action_param" TYPE varchar(50) COLLATE "pg_catalog"."default";