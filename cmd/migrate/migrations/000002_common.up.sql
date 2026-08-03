-- ===== tokens 令牌表 =====
CREATE TABLE IF NOT EXISTS "__PREFIX__tokens" (
    "token"      varchar(64) PRIMARY KEY,
    "type"       varchar(32) NOT NULL,
    "user_id"    bigint NOT NULL,
    "created_at" timestamptz,
    "expired_at" timestamptz NOT NULL
);
CREATE INDEX IF NOT EXISTS "__PREFIX__idx_tokens_user_id" ON "__PREFIX__tokens" ("user_id");
CREATE INDEX IF NOT EXISTS "__PREFIX__idx_tokens_expired_at" ON "__PREFIX__tokens" ("expired_at");
COMMENT ON TABLE "__PREFIX__tokens" IS '令牌表';
COMMENT ON COLUMN "__PREFIX__tokens"."token" IS '令牌';
COMMENT ON COLUMN "__PREFIX__tokens"."type" IS '令牌类型';
COMMENT ON COLUMN "__PREFIX__tokens"."user_id" IS '用户ID';
COMMENT ON COLUMN "__PREFIX__tokens"."created_at" IS '创建时间';
COMMENT ON COLUMN "__PREFIX__tokens"."expired_at" IS '过期时间';

-- ===== captchas 验证码表 =====
CREATE TABLE IF NOT EXISTS "__PREFIX__captchas" (
    "key"        varchar(64) PRIMARY KEY,
    "code"       varchar(255),
    "info"       text,
    "created_at" timestamptz,
    "expired_at" timestamptz NOT NULL
);
CREATE INDEX IF NOT EXISTS "__PREFIX__idx_captchas_expired_at" ON "__PREFIX__captchas" ("expired_at");
COMMENT ON TABLE "__PREFIX__captchas" IS '验证码表';
COMMENT ON COLUMN "__PREFIX__captchas"."key" IS '验证码查询键';
COMMENT ON COLUMN "__PREFIX__captchas"."code" IS '验证码值（加密后）';
COMMENT ON COLUMN "__PREFIX__captchas"."info" IS '验证码详细信息';
COMMENT ON COLUMN "__PREFIX__captchas"."created_at" IS '创建时间';
COMMENT ON COLUMN "__PREFIX__captchas"."expired_at" IS '过期时间';

-- ===== config 系统配置表 =====
CREATE TABLE IF NOT EXISTS "__PREFIX__configs" (
    "id"           bigserial PRIMARY KEY,
    "name"         varchar(64) NOT NULL DEFAULT '',
    "group"        varchar(64) NOT NULL DEFAULT '',
    "title"        varchar(64) NOT NULL DEFAULT '',
    "type"         varchar(64) NOT NULL DEFAULT '',
    "tip"          varchar(128) DEFAULT NULL,
    "value"        text DEFAULT NULL,
    "dict"         text DEFAULT NULL,
    "rule"         varchar(128) DEFAULT NULL,
    "input_extend" varchar(255) DEFAULT NULL,
    "allow_del"    smallint NOT NULL DEFAULT 1,
    "weigh"        bigint NOT NULL DEFAULT 0
);
CREATE UNIQUE INDEX IF NOT EXISTS "__PREFIX__idx_configs_name" ON "__PREFIX__configs" ("name");
COMMENT ON TABLE "__PREFIX__configs" IS '系统配置表';
COMMENT ON COLUMN "__PREFIX__configs"."id" IS 'ID';
COMMENT ON COLUMN "__PREFIX__configs"."name" IS '变量名';
COMMENT ON COLUMN "__PREFIX__configs"."group" IS '分组';
COMMENT ON COLUMN "__PREFIX__configs"."title" IS '变量标题';
COMMENT ON COLUMN "__PREFIX__configs"."tip" IS '变量输入提示';
COMMENT ON COLUMN "__PREFIX__configs"."type" IS '变量输入组件类型';
COMMENT ON COLUMN "__PREFIX__configs"."value" IS '变量值';
COMMENT ON COLUMN "__PREFIX__configs"."dict" IS '字典数据';
COMMENT ON COLUMN "__PREFIX__configs"."rule" IS '验证规则';
COMMENT ON COLUMN "__PREFIX__configs"."input_extend" IS '输入框扩展属性';
COMMENT ON COLUMN "__PREFIX__configs"."allow_del" IS '允许删除:0=否,1=是';
COMMENT ON COLUMN "__PREFIX__configs"."weigh" IS '权重';

-- ===== attachments 附件表 =====
CREATE TABLE IF NOT EXISTS "__PREFIX__attachments" (
    "id"             bigserial PRIMARY KEY,
    "topic"          varchar(64) NOT NULL DEFAULT '',
    "user_id"        bigint NOT NULL DEFAULT 0,
    "user_type"      varchar(64) NOT NULL DEFAULT '',
    "url"            varchar(255) NOT NULL DEFAULT '',
    "name"           varchar(255) NOT NULL DEFAULT '',
    "size"           bigint NOT NULL DEFAULT 0,
    "mimetype"       varchar(64) NOT NULL DEFAULT '',
    "quote"          bigint NOT NULL DEFAULT 0,
    "driver"         varchar(64) NOT NULL DEFAULT '',
    "sha1"           varchar(64) NOT NULL DEFAULT '',
    "created_at"     timestamptz DEFAULT NULL,
    "last_upload_at" timestamptz DEFAULT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS "__PREFIX__idx_attachments_sha1" ON "__PREFIX__attachments" ("sha1");
COMMENT ON TABLE "__PREFIX__attachments" IS '附件表';
COMMENT ON COLUMN "__PREFIX__attachments"."id" IS 'ID';
COMMENT ON COLUMN "__PREFIX__attachments"."topic" IS '主题分类';
COMMENT ON COLUMN "__PREFIX__attachments"."user_id" IS '上传用户ID';
COMMENT ON COLUMN "__PREFIX__attachments"."user_type" IS '上传用户身份类型';
COMMENT ON COLUMN "__PREFIX__attachments"."url" IS '存储路径';
COMMENT ON COLUMN "__PREFIX__attachments"."name" IS '原始名称';
COMMENT ON COLUMN "__PREFIX__attachments"."size" IS '大小';
COMMENT ON COLUMN "__PREFIX__attachments"."mimetype" IS 'MIME类型';
COMMENT ON COLUMN "__PREFIX__attachments"."quote" IS '上传（引用）次数';
COMMENT ON COLUMN "__PREFIX__attachments"."driver" IS '存储驱动';
COMMENT ON COLUMN "__PREFIX__attachments"."sha1" IS 'SHA1编码';
COMMENT ON COLUMN "__PREFIX__attachments"."created_at" IS '创建时间';
COMMENT ON COLUMN "__PREFIX__attachments"."last_upload_at" IS '最后上传时间';

-- ===== areas 省份地区表 =====
CREATE TABLE IF NOT EXISTS "__PREFIX__areas" (
    "id"    bigserial PRIMARY KEY,
    "pid"   bigint NOT NULL DEFAULT 0,
    "name"  varchar(64) NOT NULL DEFAULT '',
    "level" smallint NOT NULL DEFAULT 0,
    "code"  varchar(32) NOT NULL DEFAULT '',
    "zip"   varchar(32) NOT NULL DEFAULT ''
);
CREATE INDEX IF NOT EXISTS "__PREFIX__idx_areas_pid" ON "__PREFIX__areas" ("pid");
CREATE INDEX IF NOT EXISTS "__PREFIX__idx_areas_level" ON "__PREFIX__areas" ("level");
COMMENT ON TABLE "__PREFIX__areas" IS '省份地区数据表';
COMMENT ON COLUMN "__PREFIX__areas"."id" IS 'ID';
COMMENT ON COLUMN "__PREFIX__areas"."pid" IS '上级ID';
COMMENT ON COLUMN "__PREFIX__areas"."name" IS '名称';
COMMENT ON COLUMN "__PREFIX__areas"."level" IS '等级';
COMMENT ON COLUMN "__PREFIX__areas"."code" IS '行政区划代码';
COMMENT ON COLUMN "__PREFIX__areas"."zip" IS '邮编';
