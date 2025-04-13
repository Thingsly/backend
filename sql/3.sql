-- Insert a new function into the sys_function table
INSERT INTO public.sys_function (id, "name", enable_flag, description, remark) 
VALUES('function_3', 'frontend_res', 'disable', 'Frontend RSA Encryption', NULL);

-- Alter columns in the casbin_rule table to varchar(200)
ALTER TABLE "public"."casbin_rule"
ALTER COLUMN "v0" TYPE varchar(200) COLLATE "pg_catalog"."default",
  ALTER COLUMN "v1" TYPE varchar(200) COLLATE "pg_catalog"."default",
  ALTER COLUMN "v2" TYPE varchar(200) COLLATE "pg_catalog"."default",
  ALTER COLUMN "v3" TYPE varchar(200) COLLATE "pg_catalog"."default",
  ALTER COLUMN "v4" TYPE varchar(200) COLLATE "pg_catalog"."default",
  ALTER COLUMN "v5" TYPE varchar(200) COLLATE "pg_catalog"."default";

-- Definition for public.device_model_custom_control

-- Drop table (if necessary)
-- DROP TABLE public.device_model_custom_control;

CREATE TABLE public.device_model_custom_control (
	id varchar(36) NOT NULL, -- ID
	device_template_id varchar(36) NOT NULL, -- Device Template ID
	"name" varchar(255) NOT NULL, -- Name
	control_type varchar NOT NULL, -- 1. Control type 2. Telemetry 3. Attributes
	description varchar(500) NULL, -- Description
	"content" text NULL, -- Command content
	enable_status varchar(10) NOT NULL, -- Enable status 'enable' or 'disable'
	created_at timestamp NOT NULL, -- Creation time
	updated_at timestamp NOT NULL, -- Update time
	remark varchar(255) NULL, -- Remarks
	tenant_id varchar(36) NOT NULL,
	CONSTRAINT device_model_custom_control_pk PRIMARY KEY (id),
	CONSTRAINT device_model_custom_control_device_templates_fk FOREIGN KEY (device_template_id) REFERENCES public.device_templates(id) ON DELETE CASCADE
);

-- Column comments for device_model_custom_control
COMMENT ON COLUMN public.device_model_custom_control.id IS 'ID';
COMMENT ON COLUMN public.device_model_custom_control.device_template_id IS 'Device Template ID';
COMMENT ON COLUMN public.device_model_custom_control."name" IS 'Name';
COMMENT ON COLUMN public.device_model_custom_control.control_type IS '1. Control type 2. Telemetry 3. Attributes';
COMMENT ON COLUMN public.device_model_custom_control.description IS 'Description';
COMMENT ON COLUMN public.device_model_custom_control."content" IS 'Command content';
COMMENT ON COLUMN public.device_model_custom_control.enable_status IS 'Enable status: enable or disable';
COMMENT ON COLUMN public.device_model_custom_control.created_at IS 'Creation time';
COMMENT ON COLUMN public.device_model_custom_control.updated_at IS 'Update time';
COMMENT ON COLUMN public.device_model_custom_control.remark IS 'Remarks';

-- Definition for public.expected_datas

-- Drop table (if necessary)
-- DROP TABLE public.expected_datas;

CREATE TABLE public.expected_datas (
	id varchar(36) NOT NULL, -- Unique command identifier (UUID)
	device_id varchar(36) NOT NULL, -- Target device ID
	send_type varchar(50) NOT NULL, -- Command type (e.g., telemetry, attribute, command)
	payload jsonb NOT NULL, -- Command content (specific command parameters)
	created_at timestamptz(6) NOT NULL, -- Command creation time
	send_time timestamptz(6) NULL, -- Actual send time (if sent)
	status varchar(50) NOT NULL DEFAULT 'pending'::character varying, -- Command status (pending, sent, expired), default pending
	message text NULL, -- Status additional information (e.g., reason for failure)
	expiry_time timestamptz(6) NULL, -- Command expiration time (optional)
	"label" varchar(100) NULL, -- Command label (optional)
	tenant_id varchar(36) NOT NULL, -- Tenant ID (for multi-tenant systems)
	CONSTRAINT expected_datas_pkey PRIMARY KEY (id)
);

-- Column comments for expected_datas
COMMENT ON COLUMN public.expected_datas.id IS 'Unique command identifier (UUID)';
COMMENT ON COLUMN public.expected_datas.device_id IS 'Target device ID';
COMMENT ON COLUMN public.expected_datas.send_type IS 'Command type (e.g., telemetry, attribute, command)';
COMMENT ON COLUMN public.expected_datas.payload IS 'Command content (specific command parameters)';
COMMENT ON COLUMN public.expected_datas.created_at IS 'Command creation time';
COMMENT ON COLUMN public.expected_datas.send_time IS 'Actual send time (if sent)';
COMMENT ON COLUMN public.expected_datas.status IS 'Command status (pending, sent, expired), default pending';
COMMENT ON COLUMN public.expected_datas.message IS 'Status additional information (e.g., reason for failure)';
COMMENT ON COLUMN public.expected_datas.expiry_time IS 'Command expiration time (optional)';
COMMENT ON COLUMN public.expected_datas."label" IS 'Command label (optional)';
COMMENT ON COLUMN public.expected_datas.tenant_id IS 'Tenant ID (for multi-tenant systems)';

-- Foreign key constraint for expected_datas table
ALTER TABLE public.expected_datas ADD CONSTRAINT expected_datas_devices_fk FOREIGN KEY (device_id) REFERENCES public.devices(id) ON DELETE CASCADE ON UPDATE CASCADE;
