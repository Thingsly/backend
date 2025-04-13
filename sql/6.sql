-- ✅2025/3/20 v1.1.6

-- Create a lowercase index for the "name" field
CREATE INDEX idx_lower_name ON public.devices (LOWER(name));

-- Create a lowercase index for the "device_number" field
CREATE INDEX idx_lower_device_number ON public.devices (LOWER(device_number));

-- ✅2025/4/3
CREATE TABLE "public"."message_push_rule_log" (
                                                  "id" varchar(60) COLLATE "pg_catalog"."default" NOT NULL,
                                                  "user_id" varchar(60) COLLATE "pg_catalog"."default" NOT NULL,
                                                  "push_id" varchar(60) COLLATE "pg_catalog"."default" NOT NULL,
                                                  "type" int2 NOT NULL,
                                                  "create_time" timestamp(6) NOT NULL,
                                                  CONSTRAINT "message_pkey" PRIMARY KEY ("id")
)
;

ALTER TABLE "public"."message_push_rule_log"
    OWNER TO "postgres";

COMMENT ON COLUMN "public"."message_push_rule_log"."type" IS '1 Active Expiry, 2 Passive Expiry, 3 Scheduled Task, 4 Automatic Cleanup';

COMMENT ON COLUMN "public"."message_push_rule_log"."create_time" IS 'Effective Time';

COMMENT ON TABLE "public"."message_push_rule_log" IS 'Expiry Rule Log';

CREATE TABLE "public"."message_push_manage" (
                                                "id" varchar(60) COLLATE "pg_catalog"."default" NOT NULL,
                                                "user_id" varchar(60) COLLATE "pg_catalog"."default" NOT NULL,
                                                "push_id" varchar(60) COLLATE "pg_catalog"."default" NOT NULL,
                                                "device_type" varchar(60) COLLATE "pg_catalog"."default" NOT NULL,
                                                "status" int2 NOT NULL DEFAULT 1,
                                                "create_time" timestamp(6) NOT NULL,
                                                "update_time" timestamp(6),
                                                "delete_time" timestamp(6),
                                                "last_push_time" timestamp(6),
                                                "err_count" int4 DEFAULT 0,
                                                "inactive_time" timestamp(6),
                                                CONSTRAINT "message_push_manage_pkey" PRIMARY KEY ("id")
)
;

ALTER TABLE "public"."message_push_manage"
    OWNER TO "postgres";

CREATE UNIQUE INDEX "index_user_push" ON "public"."message_push_manage" USING btree (
    "user_id" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
    "push_id" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
    );

COMMENT ON COLUMN "public"."message_push_manage"."user_id" IS 'User ID';

COMMENT ON COLUMN "public"."message_push_manage"."push_id" IS 'Push ID';

COMMENT ON COLUMN "public"."message_push_manage"."device_type" IS 'Device Type';

COMMENT ON COLUMN "public"."message_push_manage"."status" IS 'Type 1: Normal, 2: Canceled';

COMMENT ON COLUMN "public"."message_push_manage"."create_time" IS 'Creation Time';

COMMENT ON COLUMN "public"."message_push_manage"."update_time" IS 'Update Time';

COMMENT ON COLUMN "public"."message_push_manage"."delete_time" IS 'Delete Time';

COMMENT ON COLUMN "public"."message_push_manage"."last_push_time" IS 'Last Push Time';

COMMENT ON COLUMN "public"."message_push_manage"."err_count" IS 'Contact Push Error Count';

COMMENT ON COLUMN "public"."message_push_manage"."inactive_time" IS 'Inactive Time Mark';

COMMENT ON TABLE "public"."message_push_manage" IS 'Message Push Notification';

CREATE TABLE "public"."message_push_log" (
                                             "id" varchar(60) COLLATE "pg_catalog"."default" NOT NULL,
                                             "user_id" varchar(60) COLLATE "pg_catalog"."default" NOT NULL,
                                             "message_type" int8 NOT NULL,
                                             "content" json NOT NULL,
                                             "status" int2 NOT NULL,
                                             "err_message" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
                                             "create_time" timestamp(6) NOT NULL,
                                             CONSTRAINT "message_push_log_pkey" PRIMARY KEY ("id")
)
;

ALTER TABLE "public"."message_push_log"
    OWNER TO "postgres";

COMMENT ON COLUMN "public"."message_push_log"."user_id" IS 'User ID';

COMMENT ON COLUMN "public"."message_push_log"."message_type" IS 'Message Type 1: Alarm Message';

COMMENT ON COLUMN "public"."message_push_log"."content" IS 'Message Body Content';

COMMENT ON COLUMN "public"."message_push_log"."status" IS '1: Push Successful, 2: Push Failed';

COMMENT ON COLUMN "public"."message_push_log"."err_message" IS 'Error Message';

COMMENT ON COLUMN "public"."message_push_log"."create_time" IS 'Send Time';

COMMENT ON TABLE "public"."message_push_log" IS 'Message Push Log';

CREATE TABLE "public"."message_push_config" (
                                                "id" varchar(60) COLLATE "pg_catalog"."default" NOT NULL,
                                                "url" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
                                                "config_type" int2 NOT NULL DEFAULT 1,
                                                "create_time" timestamp(6) NOT NULL,
                                                "update_time" timestamp(6),
                                                CONSTRAINT "message_push_config_pkey" PRIMARY KEY ("id")
)
;

ALTER TABLE "public"."message_push_config"
    OWNER TO "postgres";

COMMENT ON COLUMN "public"."message_push_config"."url" IS 'Push URL';

COMMENT ON COLUMN "public"."message_push_config"."config_type" IS 'Configuration Type 1: Push URL';

COMMENT ON COLUMN "public"."message_push_config"."create_time" IS 'Creation Time';

COMMENT ON COLUMN "public"."message_push_config"."update_time" IS 'Update Time';

COMMENT ON TABLE "public"."message_push_config" IS 'Message Push Configuration';


-- ✅2025/4/10
-- Create a view to get the latest alarm record for each device
CREATE OR REPLACE VIEW public.latest_device_alarms AS
WITH unnested_devices AS (
  SELECT 
    ah.id,
    ah.alarm_config_id,
    ah.group_id,
    ah.scene_automation_id,
    ah.name,
    ah.description,
    ah.content,
    ah.alarm_status,
    ah.tenant_id,
    ah.remark,
    ah.create_at,
    jsonb_array_elements_text(ah.alarm_device_list) AS device_id
  FROM 
    public.alarm_history ah
),
ranked_alarms AS (
  SELECT 
    *,
    ROW_NUMBER() OVER (PARTITION BY device_id ORDER BY create_at DESC) as rn
  FROM 
    unnested_devices
)
SELECT 
  id,
  alarm_config_id,
  group_id,
  scene_automation_id,
  name,
  description,
  content,
  alarm_status,
  tenant_id,
  remark,
  create_at,
  device_id
FROM 
  ranked_alarms
WHERE 
  rn = 1;
