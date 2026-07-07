-- 初始化系统配置预设数据
-- ON CONFLICT DO NOTHING: 兼容重复执行（按 id 主键或 name 唯一索引冲突时跳过）
INSERT INTO "__PREFIX__config" ("id", "name", "group", "title", "tip", "type", "value", "content", "rule", "extend", "input_extend", "allow_del", "weigh") VALUES
    (1, 'config_group', 'basic', '配置分组', '', 'array', '[{"key":"basic","value":"基础配置"},{"key":"mail","value":"邮件配置"},{"key":"config_quick_entrance","value":"快捷配置入口"}]', NULL, 'required', '', '', 0, -1),
    (2, 'name', 'basic', '站点名称', '', 'string', 'AI GO ADMIN', NULL, 'required', '', '', 0, 99),
    (3, 'entrance', 'basic', '自定义后台入口', '', 'string', '/admin', NULL, 'required', '', '', 0, 98),
    (4, 'record_number', 'basic', '域名备案号', '', 'string', '渝ICP备8888888号-1', NULL, '', '', '', 0, 97),
    (5, 'version', 'basic', '系统版本号', '', 'string', 'v1.0.0', NULL, 'required', '', '', 0, 96),
    (6, 'timezone', 'basic', '时区', '', 'string', 'Asia/Shanghai', NULL, 'required', '', '', 0, 95),
    (7, 'no_access_ip', 'basic', '禁止访问 IP', '禁止访问站点的ip列表，一行一个', 'textarea', NULL, NULL, '', '', '', 0, 94),
    (8, 'smtp_server', 'mail', 'SMTP 服务器', '', 'string', 'smtp.com', NULL, '', '', '', 0, 99),
    (9, 'smtp_port', 'mail', 'SMTP 端口', '', 'string', '465', NULL, '', '', '', 0, 98),
    (10, 'smtp_user', 'mail', 'SMTP 用户', '', 'string', NULL, NULL, '', '', '', 0, 97),
    (11, 'smtp_pass', 'mail', 'SMTP 密码', '', 'string', NULL, NULL, '', '', '', 0, 96),
    (12, 'smtp_verification', 'mail', 'SMTP 验证方式', '', 'select', 'SSL', '{"SSL":"SSL","TLS":"TLS"}', '', '', '', 0, 95),
    (13, 'smtp_sender_mail', 'mail', 'SMTP 发件人邮箱', '', 'string', NULL, NULL, 'email', '', '', 0, 94),
    (14, 'config_quick_entrance', 'quick_entrance', '快捷配置入口', '', 'array', '[]', NULL, '', '', '', 0, 0)
ON CONFLICT DO NOTHING;

-- 手动插入了显式 id，需将序列推进到当前最大值，避免后续自增插入冲突
SELECT setval('__PREFIX__config_id_seq', (SELECT COALESCE(MAX("id"), 1) FROM "__PREFIX__config"));
