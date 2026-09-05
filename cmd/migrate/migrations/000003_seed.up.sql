-- 初始化系统配置预设数据
-- ON CONFLICT DO NOTHING: 兼容重复执行（按 id 主键或 name 唯一索引冲突时跳过）
INSERT INTO "__PREFIX__configs" ("id", "name", "group", "title", "tip", "type", "value", "dict", "rule", "input_extend", "allow_del", "weigh") VALUES
	(1,  'config_group',         'basic',          '配置分组',         NULL,                                   'array',    '[{"key":"basic","value":"基础配置"},{"key":"mail","value":"邮件配置"},{"key":"ai","value":"AI 配置"},{"key":"quick_entrance","value":"快捷配置入口"}]', NULL,    'required', NULL,  0, -1),
	(2,  'name',                 'basic',          '站点名称',         NULL,                                   'string',   'AI GO ADMIN',                                                                                                    NULL,    'required', NULL,  0, 99),
	(3,  'record_number',        'basic',          '工信部备案号',     NULL,                                   'string',   '渝ICP备8888888号-1',                                                                                              NULL,    NULL,         NULL,  0, 98),
	(4,  'ps_record_number',     'basic',          '公安备案号',       NULL,                                   'string',   '渝公网安备50010808888888888号',                                                                                   NULL,    NULL,         NULL,  0, 97),
	(5,  'version',              'basic',          '系统版本号',       NULL,                                   'string',   'v1.0.0',                                                                                                          NULL,    'required', NULL,  0, 96),
	(6,  'no_access_ip',         'basic',          '禁止访问 IP',      '禁止访问站点的IP列表，一行一个',        'textarea',   NULL,                                                                                                              NULL,    NULL,         NULL,  0, 95),
	(7,  'smtp_server',          'mail',           'SMTP 服务器',      NULL,                                   'string',   'smtp.com',                                                                                                       NULL,    NULL,         NULL,  0, 99),
	(8,  'smtp_port',            'mail',           'SMTP 端口',        NULL,                                   'string',   '465',                                                                                                            NULL,    NULL,         NULL,  0, 98),
	(9,  'smtp_user',            'mail',           'SMTP 用户',        NULL,                                   'string',   NULL,                                                                                                              NULL,    NULL,         NULL,  0, 97),
	(10, 'smtp_pass',            'mail',           'SMTP 密码',        NULL,                                   'string',   NULL,                                                                                                              NULL,    NULL,         NULL,  0, 96),
	(11, 'smtp_verification',    'mail',           'SMTP 验证方式',    NULL,                                   'select',   'SSL',                                                                                                           '{"SSL":"SSL","TLS":"TLS"}', NULL, NULL,  0, 95),
	(12, 'smtp_sender_mail',     'mail',           'SMTP 发件人邮箱',  NULL,                                   'string',    NULL,                                                                                                             NULL,    'email',    NULL,  0, 94),
	(13, 'quick_entrance',       'quick_entrance', '快捷配置入口',     NULL,                                   'array',     '[{"key":"个人资料","value":"routine/profile"}]',                                                                  NULL,    NULL,         NULL,  0, 0),
	(14, 'ai_api_url',           'ai',             'API URL',        '请填写 OpenAI Responses API 兼容的 URL，如: https://api.deepseek.com/responses', 'string', '',                                                                          NULL,    'required', NULL,  0, 99),
	(15, 'ai_api_key',           'ai',             'API Key',         NULL,                                   'password',  '',                                                                                                               NULL,    'required', '{"showPassword": true}',  0, 98),
	(16, 'ai_model_list',        'ai',             '可用模型列表',     NULL,                                   'array',    '[{"key":"DeepSeek V4 Pro","value":"deepseek-v4-pro"}]',                                                           NULL,    'required', '{"keyTitle":"模型显示名","valueTitle":"模型请求名"}', 0, 97),
	(17, 'ai_default_model',     'ai',             '默认托底模型',     NULL,                                   'hidden',   'deepseek-v4-pro',                                                                                                 NULL,    'required', NULL,  0, 96)
ON CONFLICT DO NOTHING;

-- 初始化权限规则数据
INSERT INTO "__PREFIX__admin_rules" ("id", "pid", "type", "title",           "name",                     "path",                 "icon",                    "open_type", "url", "component",                                     "keepalive", "extend", "remark", "weigh", "status", "updated_at", "created_at") VALUES
    (1,  NULL, 'menu', '控制台',         'dashboard',                  'dashboard',            'lucide-chart-pie',        'tab',       '',    '/src/views/admin/dashboard.vue',                1,           '',       '',       999,     1,        NOW(),        NOW()),
    (2,  1,    'node', '查看',           'dashboard/read',             '',                     '',                        '',          '',    '',                                              1,           '',       '',       0,       1,        NOW(),        NOW()),
    (3,  NULL, 'dir',  '权限管理',       'auth',                       'auth',                 'lucide-users',            '',          '',    '',                                              1,           '',       '',       998,     1,        NOW(),        NOW()),
    (4,  3,    'menu', '角色组管理',     'auth/group',                 'auth/group',           'lucide-users-round',      'tab',       '',    '/src/views/admin/auth/group/index.vue',         1,           '',       '',       99,      1,        NOW(),        NOW()),
    (5,  4,    'node', '查看',           'auth/group/read',            '',                     '',                        '',          '',    '',                                              1,           '',       '',       0,       1,        NOW(),        NOW()),
    (6,  4,    'node', '添加',           'auth/group/create',          '',                     '',                        '',          '',    '',                                              1,           '',       '',       0,       1,        NOW(),        NOW()),
    (7,  4,    'node', '更新',           'auth/group/update',          '',                     '',                        '',          '',    '',                                              1,           '',       '',       0,       1,        NOW(),        NOW()),
    (8,  4,    'node', '删除',           'auth/group/delete',          '',                     '',                        '',          '',    '',                                              1,           '',       '',       0,       1,        NOW(),        NOW()),
    (9,  3,    'menu', '管理员管理',     'auth/admin',                 'auth/admin',           'lucide-user-key',         'tab',       '',    '/src/views/admin/auth/admin/index.vue',         1,           '',       '',       98,      1,        NOW(),        NOW()),
    (10, 9,    'node', '查看',           'auth/admin/read',            '',                     '',                        '',          '',    '',                                              1,           '',       '',       0,       1,        NOW(),        NOW()),
    (11, 9,    'node', '添加',           'auth/admin/create',          '',                     '',                        '',          '',    '',                                              1,           '',       '',       0,       1,        NOW(),        NOW()),
    (12, 9,    'node', '更新',           'auth/admin/update',          '',                     '',                        '',          '',    '',                                              1,           '',       '',       0,       1,        NOW(),        NOW()),
    (13, 9,    'node', '删除',           'auth/admin/delete',          '',                     '',                        '',          '',    '',                                              1,           '',       '',       0,       1,        NOW(),        NOW()),
    (14, 3,    'menu', '菜单规则管理',   'auth/rule',                  'auth/rule',            'lucide-list-tree',        'tab',       '',    '/src/views/admin/auth/rule/index.vue',          1,           '',       '',       97,      1,        NOW(),        NOW()),
    (15, 14,   'node', '查看',           'auth/rule/read',             '',                     '',                        '',          '',    '',                                              1,           '',       '',       0,       1,        NOW(),        NOW()),
    (16, 14,   'node', '添加',           'auth/rule/create',           '',                     '',                        '',          '',    '',                                              1,           '',       '',       0,       1,        NOW(),        NOW()),
    (17, 14,   'node', '更新',           'auth/rule/update',           '',                     '',                        '',          '',    '',                                              1,           '',       '',       0,       1,        NOW(),        NOW()),
    (18, 14,   'node', '删除',           'auth/rule/delete',           '',                     '',                        '',          '',    '',                                              1,           '',       '',       0,       1,        NOW(),        NOW()),
    (19, 14,   'node', '排序',           'auth/rule/sort',             '',                     '',                        '',          '',    '',                                              1,           '',       '',       0,       1,        NOW(),        NOW()),
    (20, 3,    'menu', '管理员日志管理',  'auth/log',                   'auth/log',            'lucide-scroll-text',      'tab',       '',    '/src/views/admin/auth/log/index.vue',            1,           '',       '',       96,      1,        NOW(),        NOW()),
    (21, 20,   'node', '查看',           'auth/log/read',              '',                     '',                        '',          '',    '',                                              1,           '',       '',       0,       1,        NOW(),        NOW()),
    (22, NULL, 'dir',  '常规管理',       'routine',                    'routine',              'lucide-settings',         '',          '',    '',                                              1,           '',       '',       997,     1,        NOW(),        NOW()),
    (23, 22,   'menu', '系统配置',       'routine/config',             'routine/config',       'lucide-wrench',           'tab',       '',    '/src/views/admin/routine/config/index.vue',     1,           '',       '',       99,      1,        NOW(),        NOW()),
    (24, 23,   'node', '查看',           'routine/config/read',        '',                     '',                        '',          '',    '',                                              1,           '',       '',       0,       1,        NOW(),        NOW()),
    (25, 23,   'node', '更新',           'routine/config/update',      '',                     '',                        '',          '',    '',                                              1,           '',       '',       0,       1,        NOW(),        NOW()),
    (26, 22,   'menu', '附件管理',       'routine/attachment',         'routine/attachment',   'lucide-folder',           'tab',       '',    '/src/views/admin/routine/attachment/index.vue', 1,           '',       '',       98,      1,        NOW(),        NOW()),
    (27, 26,   'node', '查看',           'routine/attachment/read',    '',                     '',                        '',          '',    '',                                              1,           '',       '',       0,       1,        NOW(),        NOW()),
    (28, 26,   'node', '更新',           'routine/attachment/update',  '',                     '',                        '',          '',    '',                                              1,           '',       '',       0,       1,        NOW(),        NOW()),
    (29, 26,   'node', '删除',           'routine/attachment/delete',  '',                     '',                        '',          '',    '',                                              1,           '',       '',       0,       1,        NOW(),        NOW()),
    (30, 22,   'menu', '个人资料',       'routine/profile',            'routine/profile',      'lucide-user',             'tab',       '',    '/src/views/admin/routine/profile.vue',          1,           '',       '',       97,      1,        NOW(),        NOW()),
    (31, 30,   'node', '查看',           'routine/profile/read',       '',                     '',                        '',          '',    '',                                              1,           '',       '',       0,       1,        NOW(),        NOW()),
    (32, 30,   'node', '更新',           'routine/profile/update',     '',                     '',                        '',          '',    '',                                              1,           '',       '',       0,       1,        NOW(),        NOW()),
    (33, NULL, 'menu', 'CRUD代码生成',   'crud/crud',                  'crud/crud',            'lucide-code-xml',         'tab',       '',    '/src/views/admin/crud/index.vue',               1,           '',       '',       996,     1,        NOW(),        NOW()),
    (34, 33,   'node', '查看',           'crud/crud/read',             '',                     '',                        '',          '',    '',                                              1,           '',       '',       0,       1,        NOW(),        NOW()),
    (35, 33,   'node', '创建',           'crud/crud/create',           '',                     '',                        '',          '',    '',                                              1,           '',       '',       0,       1,        NOW(),        NOW()),
    (36, 33,   'node', '删除',           'crud/crud/delete',           '',                     '',                        '',          '',    '',                                              1,           '',       '',       0,       1,        NOW(),        NOW())
ON CONFLICT DO NOTHING;

-- 初始化管理员分组
INSERT INTO "__PREFIX__admin_groups" ("id", "pid", "name", "rules", "status", "updated_at", "created_at") VALUES
    (1, NULL, '超级管理组', '*', 1, NOW(), NOW())
ON CONFLICT DO NOTHING;

-- 初始化管理员
INSERT INTO "__PREFIX__admins" ("id", "username", "nickname", "avatar", "email", "mobile", "login_failure", "last_login_at", "last_login_ip", "password", "bio", "status", "updated_at", "created_at", "deleted_at") VALUES
    (1, 'admin', 'Admin', '/static/images/avatar.png', 'admin@ai-go-hub.com', '18888888888', 0, NOW(), '::1', '$2a$10$T7nw94qiMClADHg0KezIzukmNvvIxHG.V5uObKOGwPTVq4W2ee7B.', NULL, 'enable', NOW(), NOW(), NULL)
ON CONFLICT DO NOTHING;

-- 初始化管理员与分组映射
INSERT INTO "__PREFIX__admin_group_access" ("uid", "group_id") VALUES
    (1, 1)
ON CONFLICT DO NOTHING;

-- setval
SELECT setval('__PREFIX__admins_id_seq', (SELECT COALESCE(MAX("id"), 1) FROM "__PREFIX__admins"));
SELECT setval('__PREFIX__configs_id_seq', (SELECT COALESCE(MAX("id"), 1) FROM "__PREFIX__configs"));
SELECT setval('__PREFIX__admin_rules_id_seq', (SELECT COALESCE(MAX("id"), 1) FROM "__PREFIX__admin_rules"));
SELECT setval('__PREFIX__admin_groups_id_seq', (SELECT COALESCE(MAX("id"), 1) FROM "__PREFIX__admin_groups"));