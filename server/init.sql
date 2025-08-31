SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;

-- 刷刷题数据库初始化脚本 (修正版)
-- 基于 models.go 结构生成

CREATE DATABASE IF NOT EXISTS shuashuati DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE shuashuati;

-- 用户表 (统一用户和管理员)
CREATE TABLE IF NOT EXISTS `users` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `open_id` VARCHAR(100) UNIQUE NOT NULL, -- Changed from openid to open_id
    `username` VARCHAR(50),
    `email` VARCHAR(100),
    `password` VARCHAR(255),
    `nickname` VARCHAR(100) DEFAULT '',
    `avatar` VARCHAR(500) DEFAULT '',
    `role` ENUM('user','admin') DEFAULT 'user',
    `status` VARCHAR(20) DEFAULT 'active',
    `is_verified` BOOLEAN DEFAULT false,
    `is_guest` BOOLEAN DEFAULT false,
    `total_answered` INT DEFAULT 0,
    `total_correct` INT DEFAULT 0,
    `accuracy_rate` DECIMAL(5,2) DEFAULT 0.00,
    `last_active_time` TIMESTAMP NULL,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- 管理员模型 (Added based on models.go)
CREATE TABLE IF NOT EXISTS `admins` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `username` VARCHAR(50) UNIQUE NOT NULL,
    `password_hash` VARCHAR(255) NOT NULL,
    `email` VARCHAR(100),
    `role` VARCHAR(20) DEFAULT 'admin',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='管理员模型';

-- 分类表
CREATE TABLE IF NOT EXISTS `categories` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `name` VARCHAR(100) NOT NULL,
    `description` VARCHAR(255),
    `parent_id` BIGINT UNSIGNED,
    `level` INT DEFAULT 1,
    `sort` INT DEFAULT 0,
    `status` INT DEFAULT 1 COMMENT '状态 1-启用 0-禁用',
    `question_count` INT DEFAULT 0,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`parent_id`) REFERENCES `categories`(`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='分类表';

-- 题目表
CREATE TABLE IF NOT EXISTS `questions` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `title` VARCHAR(500) NOT NULL,
    `content` TEXT NOT NULL,
    `type` VARCHAR(20) DEFAULT 'single',
    `options` JSON NOT NULL,
    `correct_answer` INT NOT NULL,
    `explanation` TEXT,
    `difficulty` VARCHAR(20) DEFAULT 'medium',
    `category_id` BIGINT UNSIGNED NOT NULL,
    `creator_id` BIGINT UNSIGNED,
    `total_answered` INT DEFAULT 0,
    `total_correct` INT DEFAULT 0,
    `accuracy_rate` DECIMAL(5,2) DEFAULT 0.00,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`category_id`) REFERENCES `categories`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`creator_id`) REFERENCES `users`(`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='题目表';

-- 答题记录表
CREATE TABLE IF NOT EXISTS `answer_records` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `user_id` BIGINT UNSIGNED NOT NULL,
    `question_id` BIGINT UNSIGNED NOT NULL,
    `user_answer` INT NOT NULL,
    `is_correct` BOOLEAN NOT NULL,
    `time_spent` INT DEFAULT 0,
    `answered_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Reverted to original
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Added created_at
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP -- Added updated_at
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='答题记录表';

-- 错题本表
CREATE TABLE IF NOT EXISTS `mistake_books` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `user_id` BIGINT UNSIGNED NOT NULL,
    `question_id` BIGINT UNSIGNED NOT NULL,
    `is_mastered` BOOLEAN DEFAULT false,
    `added_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Added created_at
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP -- Added updated_at
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='错题本表';

-- 操作日志表
CREATE TABLE IF NOT EXISTS `operation_logs` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `operator` VARCHAR(50) NOT NULL,
    `action` VARCHAR(20) NOT NULL,
    `resource` VARCHAR(50) NOT NULL,
    `description` VARCHAR(255),
    `ip` VARCHAR(45),
    `user_agent` VARCHAR(500),
    `request_data` TEXT,
    `response_data` TEXT,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='操作日志表';

-- 系统设置表
CREATE TABLE IF NOT EXISTS `system_settings` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `key` VARCHAR(255) UNIQUE NOT NULL COMMENT '设置键名',
    `value` TEXT COMMENT '设置值(JSON格式)',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` TIMESTAMP NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统设置表';

-- 插入默认管理员账号
-- 密码是 123456 的 bcrypt 哈希值
INSERT INTO `users` (`open_id`, `username`, `password`, `nickname`, `role`, `is_verified`) VALUES
('admin_unique_openid_123', 'admin', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', '系统管理员', 'admin', true);

-- 插入测试用户
INSERT INTO `users` (`open_id`, `username`, `password`, `nickname`, `role`, `is_verified`) VALUES
('testuser_unique_openid_456', 'testuser', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', '测试用户', 'user', true);

-- 插入初始分类数据
INSERT INTO `categories` (`id`, `name`, `parent_id`, `level`, `sort`, `description`) VALUES
(1, '道路交通安全法律、法规', NULL, 1, 1, '关于交通法规的题目'),
(2, '交通信号', NULL, 1, 2, '关于交通信号灯、标志、标线的题目'),
(3, '安全行车、文明驾驶基础知识', NULL, 1, 3, '关于安全驾驶和文明驾驶的题目'),
(4, '机动车驾驶操作相关基础知识', NULL, 1, 4, '关于驾驶操作和车辆基础知识的题目');

-- 插入测试题目数据
INSERT INTO `questions` (`title`, `content`, `type`, `options`, `correct_answer`, `explanation`, `difficulty`, `category_id`, `creator_id`) VALUES
('以下哪项不是交通信号灯的颜色？', '交通信号灯通常由红、黄、绿三种颜色组成。', 'single', '["红色", "蓝色", "黄色", "绿色"]', 1, '交通信号灯没有蓝色。', 'easy', 2, NULL),
('在没有交通信号灯的路口，车辆应该如何通行？', '在没有交通信号灯的路口，车辆应遵守交通标志、标线，并注意避让行人。', 'single', '["加速通过", "鸣笛示意", "减速慢行，注意观察", "随意通行"]', 2, '在没有交通信号灯的路口，应减速慢行，注意观察，确保安全。', 'medium', 1, NULL),
('驾驶机动车在高速公路上行驶，遇到雾、雨、雪、沙尘、冰雹等低能见度气象条件时，应当如何操作？', '在高速公路低能见度条件下行驶，应降低车速，保持安全距离，并开启相关灯光。', 'single', '["开启远光灯", "开启危险报警闪光灯", "开启示廓灯和前后位灯", "开启雾灯、近光灯、示廓灯和前后位灯"]', 3, '在低能见度条件下，应开启雾灯、近光灯、示廓灯和前后位灯，并降低车速。', 'hard', 3, NULL),
('以下哪种行为属于文明驾驶？', '文明驾驶是驾驶员应具备的基本素质。', 'single', '["随意变道", "超速行驶", "礼让行人", "酒后驾车"]', 2, '礼让行人是文明驾驶的重要体现。', 'easy', 3, NULL),
('机动车驾驶人在道路上发生交通事故，造成人身伤亡的，应当立即抢救受伤人员，并迅速报告公安机关交通管理部门。对吗？', '交通事故处理的基本原则。', 'single', '["对", "错"]', 0, '发生交通事故，造成人身伤亡的，应立即抢救受伤人员，并迅速报告公安机关交通管理部门。', 'medium', 1, NULL),
('驾驶机动车通过没有交通信号灯控制也没有交通警察指挥的交叉路口，相对方向行驶的右转弯的机动车让左转弯的机动车先行。对吗？', '交叉路口通行规则。', 'single', '["对", "错"]', 1, '相对方向行驶的右转弯的机动车应让左转弯的机动车先行。', 'hard', 1, NULL),
('以下哪项是驾驶机动车时禁止的行为？', '驾驶机动车时应遵守交通法规。', 'single', '["系安全带", "使用手机", "开启转向灯", "保持安全距离"]', 1, '驾驶机动车时禁止使用手机。', 'easy', 3, NULL),
('在高速公路上行驶，车速超过规定时速50%的，公安机关交通管理部门可以吊销机动车驾驶证。对吗？', '超速行驶的处罚。', 'single', '["对", "错"]', 0, '在高速公路上行驶，车速超过规定时速50%的，可以吊销机动车驾驶证。', 'medium', 1, NULL),
('以下哪种情况可以不避让正在执行紧急任务的警车？', '特种车辆的避让原则。', 'single', '["在非机动车道行驶", "在高速公路行驶", "在同方向只有一条机动车道的道路上行驶", "在执行紧急任务的警车前方行驶"]', 2, '在同方向只有一条机动车道的道路上，可以不避让正在执行紧急任务的警车。', 'hard', 1, NULL),
('驾驶机动车在道路上行驶，遇有前方车辆停车排队等候或者缓慢行驶时，不得借道超车或者占用对面车道、穿插等候的车辆。对吗？', '交通拥堵时的驾驶行为。', 'single', '["对", "错"]', 0, '在交通拥堵时，不得借道超车或者占用对面车道、穿插等候的车辆。', 'medium', 3, NULL);
SELECT 'Database initialization completed successfully!' as message;