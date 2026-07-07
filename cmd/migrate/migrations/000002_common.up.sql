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
    "expired_at" timestamptz NOT NULL,
    "created_at" timestamptz
);
CREATE INDEX IF NOT EXISTS "__PREFIX__idx_captchas_expired_at" ON "__PREFIX__captchas" ("expired_at");
COMMENT ON TABLE "__PREFIX__captchas" IS '验证码表';
COMMENT ON COLUMN "__PREFIX__captchas"."key" IS '验证码查询键';
COMMENT ON COLUMN "__PREFIX__captchas"."code" IS '验证码值（加密后）';
COMMENT ON COLUMN "__PREFIX__captchas"."info" IS '验证码详细信息';
COMMENT ON COLUMN "__PREFIX__captchas"."expired_at" IS '过期时间';
COMMENT ON COLUMN "__PREFIX__captchas"."created_at" IS '创建时间';

-- ===== config 系统配置表 =====
CREATE TABLE IF NOT EXISTS "__PREFIX__config" (
    "id"           bigserial PRIMARY KEY,
    "name"         varchar(50) NOT NULL DEFAULT '',
    "group"        varchar(50) NOT NULL DEFAULT '',
    "title"        varchar(50) NOT NULL DEFAULT '',
    "tip"          varchar(100) NOT NULL DEFAULT '',
    "type"         varchar(50) NOT NULL DEFAULT '',
    "value"        text,
    "content"      text,
    "rule"         varchar(100) NOT NULL DEFAULT '',
    "extend"       varchar(255) NOT NULL DEFAULT '',
    "input_extend" varchar(255) NOT NULL DEFAULT '',
    "allow_del"    smallint NOT NULL DEFAULT 0,
    "weigh"        bigint NOT NULL DEFAULT 0
);
CREATE UNIQUE INDEX IF NOT EXISTS "__PREFIX__idx_config_name" ON "__PREFIX__config" ("name");
COMMENT ON TABLE "__PREFIX__config" IS '系统配置表';
COMMENT ON COLUMN "__PREFIX__config"."id" IS 'ID';
COMMENT ON COLUMN "__PREFIX__config"."name" IS '变量名';
COMMENT ON COLUMN "__PREFIX__config"."group" IS '分组';
COMMENT ON COLUMN "__PREFIX__config"."title" IS '变量标题';
COMMENT ON COLUMN "__PREFIX__config"."tip" IS '变量描述';
COMMENT ON COLUMN "__PREFIX__config"."type" IS '变量输入组件类型';
COMMENT ON COLUMN "__PREFIX__config"."value" IS '变量值';
COMMENT ON COLUMN "__PREFIX__config"."content" IS '字典数据';
COMMENT ON COLUMN "__PREFIX__config"."rule" IS '验证规则';
COMMENT ON COLUMN "__PREFIX__config"."extend" IS '扩展属性';
COMMENT ON COLUMN "__PREFIX__config"."input_extend" IS '输入框扩展属性';
COMMENT ON COLUMN "__PREFIX__config"."allow_del" IS '允许删除:0=否,1=是';
COMMENT ON COLUMN "__PREFIX__config"."weigh" IS '权重';
