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
