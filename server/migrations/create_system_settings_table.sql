-- 创建系统设置表
CREATE TABLE IF NOT EXISTS `system_settings` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `key` varchar(255) NOT NULL COMMENT '设置键名',
  `value` text COMMENT '设置值(JSON格式)',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_system_settings_key` (`key`),
  KEY `idx_system_settings_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统设置表';

-- 插入默认设置
INSERT INTO `system_settings` (`key`, `value`, `created_at`, `updated_at`) VALUES
('basic', '{"system_name":"刷刷题","system_description":"专业的在线刷题平台","system_version":"1.0.0","contact_email":"admin@example.com","system_status":"normal","maintenance_notice":""}', NOW(), NOW()),
('quiz', '{"daily_limit":0,"time_limit":0,"enable_points":true,"correct_points":1,"wrong_points":0,"quiz_modes":["random","category"],"show_explanation":"after_answer"}', NOW(), NOW())
ON DUPLICATE KEY UPDATE
`updated_at` =