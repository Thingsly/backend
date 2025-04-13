-- Create service_plugins table
CREATE TABLE service_plugins (
     id VARCHAR(36) PRIMARY KEY, -- Service ID
     name VARCHAR(255) NOT NULL, -- Service name
     service_identifier VARCHAR(100) NOT NULL, -- Service identifier
     service_type INT NOT NULL CHECK (service_type IN (1, 2)), -- Service type: 1 - Access Protocol, 2 - Access Service
     last_active_time TIMESTAMP, -- Last active time of the service
     version VARCHAR(100), -- Version number
     create_at TIMESTAMP NOT NULL, -- Creation time
     update_at TIMESTAMP NOT NULL, -- Update time
     description VARCHAR(255), -- Description
     service_config JSON, -- Service configuration
     remark VARCHAR(255) -- Remarks
);

-- Add unique constraints
ALTER TABLE service_plugins
    ADD CONSTRAINT unique_service_identifier UNIQUE (service_identifier); -- Ensure unique service identifier

ALTER TABLE service_plugins
    ADD CONSTRAINT unique_name UNIQUE (name); -- Ensure unique service name

-- Alter columns to use timestamptz type
ALTER TABLE "public"."service_plugins"
ALTER COLUMN "create_at" TYPE timestamptz USING "create_at"::timestamptz, -- Change create_at to timestamptz
  ALTER COLUMN "update_at" TYPE timestamptz USING "update_at"::timestamptz; -- Change update_at to timestamptz

-- Create service_access table
CREATE TABLE service_access (
    id VARCHAR(36) PRIMARY KEY, -- Access ID
    name VARCHAR(100) NOT NULL, -- Access name
    service_plugin_id VARCHAR(36) NOT NULL, -- Service ID
    voucher VARCHAR(999) NOT NULL, -- Voucher
    description VARCHAR(255), -- Description
    service_access_config JSON, -- Service access configuration
    remark VARCHAR(255), -- Remarks
    CONSTRAINT fk_service_plugin
        FOREIGN KEY (service_plugin_id)
            REFERENCES service_plugins (id)
            ON DELETE RESTRICT -- Prevent deletion of the referenced service plugin
);

-- Add new columns to service_access
ALTER TABLE "public"."service_access"
    ADD COLUMN "create_at" timestamptz, -- Add create_at column
  ADD COLUMN "update_at" timestamptz, -- Add update_at column
  ADD COLUMN "tenant_id" varchar(36) NOT NULL; -- Add tenant_id column

-- Set not null constraints on create_at and update_at
ALTER TABLE "public"."service_access"
    ALTER COLUMN "create_at" SET NOT NULL, -- Set create_at as NOT NULL
ALTER COLUMN "update_at" SET NOT NULL; -- Set update_at as NOT NULL

-- Add foreign key to devices table
ALTER TABLE public.devices ADD service_access_id varchar(36) NULL; -- Add service_access_id column
ALTER TABLE public.devices ADD CONSTRAINT devices_service_access_fk FOREIGN KEY (service_access_id) REFERENCES public.service_access(id) ON DELETE RESTRICT; -- Add foreign key constraint

-- Add comments to service_plugins table and columns
COMMENT ON TABLE service_plugins IS 'Service Management';
COMMENT ON COLUMN service_plugins.id IS 'Service ID';
COMMENT ON COLUMN service_plugins.name IS 'Service name';
COMMENT ON COLUMN service_plugins.service_identifier IS 'Service identifier';
COMMENT ON COLUMN service_plugins.service_type IS 'Service type: 1 - Access Protocol, 2 - Access Service';
COMMENT ON COLUMN service_plugins.last_active_time IS 'Last active time';
COMMENT ON COLUMN service_plugins.version IS 'Version number';
COMMENT ON COLUMN service_plugins.create_at IS 'Creation time';
COMMENT ON COLUMN service_plugins.update_at IS 'Update time';
COMMENT ON COLUMN service_plugins.description IS 'Description';
COMMENT ON COLUMN service_plugins.service_config IS 'Service configuration';
COMMENT ON COLUMN service_plugins.remark IS 'Remarks';

-- Add comments to service_access table and columns
COMMENT ON TABLE service_access IS 'Service Access (Tenant-side)';
COMMENT ON COLUMN service_access.id IS 'Access ID';
COMMENT ON COLUMN service_access.name IS 'Name';
COMMENT ON COLUMN service_access.service_plugin_id IS 'Service ID';
COMMENT ON COLUMN service_access.voucher IS 'Voucher';
COMMENT ON COLUMN service_access.description IS 'Description';
COMMENT ON COLUMN service_access.service_access_config IS 'Service configuration';
COMMENT ON COLUMN service_access.create_at IS 'Creation time';
COMMENT ON COLUMN service_access.update_at IS 'Update time';
COMMENT ON COLUMN service_access.tenant_id IS 'Tenant ID';
COMMENT ON COLUMN service_access.remark IS 'Remarks';

-- Add additional comment for service_plugins service_config column
COMMENT ON COLUMN service_plugins.service_config IS 'Service configuration: Access Protocol and Access Service configurations';

-- Alter columns in scene_action_info table
ALTER TABLE "public"."scene_action_info"
ALTER COLUMN "action_param" TYPE varchar(50) COLLATE "pg_catalog"."default"; -- Change action_param to varchar(50)

-- Modify foreign key constraints in telemetry_set_logs table
ALTER TABLE public.telemetry_set_logs DROP CONSTRAINT telemetry_set_logs_device_id_fkey;
ALTER TABLE public.telemetry_set_logs ADD CONSTRAINT telemetry_set_logs_device_id_fkey FOREIGN KEY (device_id) REFERENCES public.devices(id) ON DELETE CASCADE;

-- Modify foreign key constraints in attribute_set_logs table
ALTER TABLE public.attribute_set_logs DROP CONSTRAINT attribute_set_logs_device_id_fkey;
ALTER TABLE public.attribute_set_logs ADD CONSTRAINT attribute_set_logs_device_id_fkey FOREIGN KEY (device_id) REFERENCES public.devices(id) ON DELETE CASCADE;

-- Modify foreign key constraints in command_set_logs table
ALTER TABLE public.command_set_logs DROP CONSTRAINT command_set_logs_device_id_fkey;
ALTER TABLE public.command_set_logs ADD CONSTRAINT command_set_logs_device_id_fkey FOREIGN KEY (device_id) REFERENCES public.devices(id) ON DELETE CASCADE;

-- Alter last_active_time column in service_plugins table
ALTER TABLE public.service_plugins ALTER COLUMN last_active_time TYPE timestamptz USING last_active_time::timestamptz;

-- Insert UI elements into sys_ui_elements table (example)
-- DELETE FROM public.sys_ui_elements WHERE id='367dbdb9-f28b-7a49-b8cd-23a915015093';
-- INSERT INTO public.sys_ui_elements (id, parent_id, element_code, element_type, orders, param1, param2, param3, authority, description, created_at, remark, multilingual, route_path) 
-- VALUES('075d9f19-5618-bb9b-6ccd-f382bfd3292b', '5373a6a2-1861-af35-eb4c-adfd5ca55ecd', 'device_service-access', 3, 1129, '/device/service-access', 'mdi:ab-testing', '0', '["TENANT_ADMIN"]'::json, 'Service Access Management', '2024-07-01 21:52:09.402', '', 'route.device_service_access', '');
-- INSERT INTO public.sys_ui_elements (id, parent_id, element_code, element_type, orders, param1, param2, param3, authority, description, created_at, remark, multilingual, route_path) 
-- VALUES('f960c45c-6d5b-e67a-c4ff-1f0e869c1625', '5373a6a2-1861-af35-eb4c-adfd5ca55ecd', 'device_service-details', 3, 1130, '/device/service-details', 'ph:align-bottom', '1', '["TENANT_ADMIN"]'::json, 'Service Details', '2024-07-01 23:16:56.668', '', 'route.device_service_details', '');
-- INSERT INTO public.sys_ui_elements (id, parent_id, element_code, element_type, orders, param1, param2, param3, authority, description, created_at, remark, multilingual, route_path) 
-- VALUES('29a684f9-c2bb-1a6f-6045-314944bef580', 'a2c53126-029f-7138-4d7a-f45491f396da', 'plug_in', 3, 32, '/apply/plugin', 'mdi:emoticon', '0', '["SYS_ADMIN"]'::json, 'Plugin Management', '2024-06-29 01:04:51.301', '', 'route.apply_in', '');
