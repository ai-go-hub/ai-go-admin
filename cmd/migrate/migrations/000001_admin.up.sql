CREATE TABLE IF NOT EXISTS "__PREFIX__admins" (
    "id"              bigserial PRIMARY KEY,
    "username"        varchar(64),
    "nickname"        varchar(64),
    "avatar"          varchar(255),
    "email"           varchar(128),
    "mobile"          varchar(16),
    "login_failure"   bigint NOT NULL DEFAULT 0,
    "last_login_at"   timestamptz,
    "last_login_ip"   varchar(64),
    "password"        varchar(255),
    "bio"             varchar(255),
    "status"          varchar(64),
    "updated_at"      timestamptz,
    "created_at"      timestamptz,
    "deleted_at"      timestamptz
);

CREATE INDEX IF NOT EXISTS "__PREFIX__idx_admins_deleted_at" ON "__PREFIX__admins" ("deleted_at");

COMMENT ON TABLE "__PREFIX__admins" IS '管理员表';
COMMENT ON COLUMN "__PREFIX__admins"."id" IS 'ID';
COMMENT ON COLUMN "__PREFIX__admins"."username" IS '用户名';
COMMENT ON COLUMN "__PREFIX__admins"."nickname" IS '昵称';
COMMENT ON COLUMN "__PREFIX__admins"."avatar" IS '头像';
COMMENT ON COLUMN "__PREFIX__admins"."email" IS '邮箱';
COMMENT ON COLUMN "__PREFIX__admins"."mobile" IS '手机号';
COMMENT ON COLUMN "__PREFIX__admins"."login_failure" IS '连续登录失败次数';
COMMENT ON COLUMN "__PREFIX__admins"."last_login_at" IS '上次登录时间';
COMMENT ON COLUMN "__PREFIX__admins"."last_login_ip" IS '上次登录IP';
COMMENT ON COLUMN "__PREFIX__admins"."password" IS '密码';
COMMENT ON COLUMN "__PREFIX__admins"."bio" IS '个人简介';
COMMENT ON COLUMN "__PREFIX__admins"."status" IS '状态:enable=启用,disable=禁用';
COMMENT ON COLUMN "__PREFIX__admins"."updated_at" IS '更新时间';
COMMENT ON COLUMN "__PREFIX__admins"."created_at" IS '创建时间';
COMMENT ON COLUMN "__PREFIX__admins"."deleted_at" IS '删除时间';

-- ===== admin_rule 菜单和权限规则表 =====
CREATE TABLE IF NOT EXISTS "__PREFIX__admin_rules" (
    "id"          bigserial PRIMARY KEY,
    "pid"         bigint DEFAULT NULL,
    "type"        varchar(50) NOT NULL DEFAULT '',
    "title"       varchar(50) NOT NULL DEFAULT '',
    "name"        varchar(50) NOT NULL DEFAULT '',
    "path"        varchar(255) NOT NULL DEFAULT '',
    "icon"        varchar(50) NOT NULL DEFAULT '',
    "open_type"   varchar(50) NOT NULL DEFAULT '',
    "url"         varchar(255) NOT NULL DEFAULT '',
    "component"   varchar(255) NOT NULL DEFAULT '',
    "keepalive"   smallint NOT NULL DEFAULT 0,
    "extend"      varchar(50) NOT NULL DEFAULT '',
    "remark"      varchar(255) NOT NULL DEFAULT '',
    "weigh"       bigint NOT NULL DEFAULT 0,
    "status"      smallint NOT NULL DEFAULT 1,
    "update_time" timestamptz,
    "create_time" timestamptz
);

CREATE INDEX IF NOT EXISTS "__PREFIX__idx_admin_rules_pid" ON "__PREFIX__admin_rules" ("pid");

COMMENT ON TABLE "__PREFIX__admin_rules" IS '菜单和权限规则表';
COMMENT ON COLUMN "__PREFIX__admin_rules"."id" IS 'ID';
COMMENT ON COLUMN "__PREFIX__admin_rules"."pid" IS '上级规则';
COMMENT ON COLUMN "__PREFIX__admin_rules"."type" IS '规则类型:dir=规则目录,menu=菜单项,node=权限节点';
COMMENT ON COLUMN "__PREFIX__admin_rules"."title" IS '规则标题';
COMMENT ON COLUMN "__PREFIX__admin_rules"."name" IS '规则名称';
COMMENT ON COLUMN "__PREFIX__admin_rules"."path" IS '菜单路由路径';
COMMENT ON COLUMN "__PREFIX__admin_rules"."icon" IS '菜单图标';
COMMENT ON COLUMN "__PREFIX__admin_rules"."open_type" IS '菜单打开方式:tab=选项卡,link=链接,iframe=Iframe';
COMMENT ON COLUMN "__PREFIX__admin_rules"."url" IS '菜单URL';
COMMENT ON COLUMN "__PREFIX__admin_rules"."component" IS '菜单组件路径';
COMMENT ON COLUMN "__PREFIX__admin_rules"."keepalive" IS '缓存:0=关闭,1=开启';
COMMENT ON COLUMN "__PREFIX__admin_rules"."extend" IS '扩展属性:add_route_only=只添加为路由,add_menu_only=只添加为菜单';
COMMENT ON COLUMN "__PREFIX__admin_rules"."remark" IS '备注';
COMMENT ON COLUMN "__PREFIX__admin_rules"."weigh" IS '权重';
COMMENT ON COLUMN "__PREFIX__admin_rules"."status" IS '状态:0=禁用,1=启用';
COMMENT ON COLUMN "__PREFIX__admin_rules"."update_time" IS '更新时间';
COMMENT ON COLUMN "__PREFIX__admin_rules"."create_time" IS '创建时间';
