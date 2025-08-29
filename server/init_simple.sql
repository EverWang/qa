-- 刷刷题数据库初始化脚本（简化版）
-- 创建数据库和表结构，插入初始数据

-- 创建数据库（如果不存在）
CREATE DATABASE IF NOT EXISTS shuashuati DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE shuashuati;

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(36) PRIMARY KEY DEFAULT (UUID()),
    openid VARCHAR(100) UNIQUE COMMENT '微信OpenID',
    username VARCHAR(50) UNIQUE COMMENT '用户名',
    email VARCHAR(100) COMMENT '邮箱',
    password VARCHAR(255) COMMENT '密码哈希',
    nickname VARCHAR(100) NOT NULL COMMENT '用户昵称',
    avatar VARCHAR(500) DEFAULT '' COMMENT '头像URL',
    role VARCHAR(20) DEFAULT 'user' COMMENT '用户角色',
    status VARCHAR(20) DEFAULT 'active' COMMENT '用户状态',
    is_verified BOOLEAN DEFAULT FALSE COMMENT '是否验证',
    is_guest BOOLEAN DEFAULT FALSE COMMENT '是否游客',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_openid (openid),
    INDEX idx_username (username),
    INDEX idx_role (role)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- 管理员表
CREATE TABLE IF NOT EXISTS admins (
    id VARCHAR(36) PRIMARY KEY DEFAULT (UUID()),
    username VARCHAR(50) NOT NULL UNIQUE COMMENT '用户名',
    password_hash VARCHAR(255) NOT NULL COMMENT '密码哈希',
    email VARCHAR(100) DEFAULT '' COMMENT '邮箱',
    role VARCHAR(20) DEFAULT 'admin' COMMENT '角色',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_username (username)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='管理员表';

-- 分类表
CREATE TABLE IF NOT EXISTS categories (
    id VARCHAR(36) PRIMARY KEY DEFAULT (UUID()),
    name VARCHAR(100) NOT NULL COMMENT '分类名称',
    description VARCHAR(255) COMMENT '分类描述',
    parent_id VARCHAR(36) NULL COMMENT '父分类ID',
    level INT DEFAULT 1 COMMENT '分类层级',
    sort INT DEFAULT 0 COMMENT '排序权重',
    status INT DEFAULT 1 COMMENT '状态 1-启用 0-禁用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_parent_id (parent_id),
    INDEX idx_level (level),
    INDEX idx_sort (sort)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='分类表';

-- 题目表
CREATE TABLE IF NOT EXISTS questions (
    id VARCHAR(36) PRIMARY KEY DEFAULT (UUID()),
    title VARCHAR(500) NOT NULL COMMENT '题目标题',
    content TEXT NOT NULL COMMENT '题目内容',
    type VARCHAR(20) DEFAULT 'single' COMMENT '题目类型',
    options JSON NOT NULL COMMENT '选项数组',
    correct_answer INT NOT NULL COMMENT '正确答案索引',
    explanation TEXT COMMENT '答案解析',
    difficulty VARCHAR(20) DEFAULT 'medium' COMMENT '难度等级',
    category_id VARCHAR(36) NOT NULL COMMENT '分类ID',
    creator_id VARCHAR(36) COMMENT '创建者ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_category_id (category_id),
    INDEX idx_difficulty (difficulty),
    INDEX idx_creator_id (creator_id),
    FULLTEXT idx_title_content (title, content)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='题目表';

-- 答题记录表
CREATE TABLE IF NOT EXISTS answer_records (
    id VARCHAR(36) PRIMARY KEY DEFAULT (UUID()),
    user_id VARCHAR(36) NOT NULL COMMENT '用户ID',
    question_id VARCHAR(36) NOT NULL COMMENT '题目ID',
    user_answer INT NOT NULL COMMENT '用户答案',
    is_correct BOOLEAN NOT NULL COMMENT '是否正确',
    time_spent INT DEFAULT 0 COMMENT '答题耗时(秒)',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_question_id (question_id),
    INDEX idx_is_correct (is_correct),
    INDEX idx_user_question (user_id, question_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='答题记录表';

-- 错题本表
CREATE TABLE IF NOT EXISTS mistake_books (
    id VARCHAR(36) PRIMARY KEY DEFAULT (UUID()),
    user_id VARCHAR(36) NOT NULL COMMENT '用户ID',
    question_id VARCHAR(36) NOT NULL COMMENT '题目ID',
    is_mastered BOOLEAN DEFAULT FALSE COMMENT '是否已掌握',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_user_question (user_id, question_id),
    INDEX idx_user_id (user_id),
    INDEX idx_question_id (question_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='错题本表';

-- 系统设置表
CREATE TABLE IF NOT EXISTS system_settings (
    id VARCHAR(36) PRIMARY KEY DEFAULT (UUID()),
    setting_key VARCHAR(100) NOT NULL UNIQUE COMMENT '设置键',
    setting_value TEXT COMMENT '设置值',
    description VARCHAR(255) COMMENT '设置描述',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_key (setting_key)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统设置表';

-- 操作日志表
CREATE TABLE IF NOT EXISTS operation_logs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    operator VARCHAR(50) NOT NULL COMMENT '操作者',
    action VARCHAR(20) NOT NULL COMMENT '操作动作',
    resource VARCHAR(50) NOT NULL COMMENT '操作资源',
    description VARCHAR(255) COMMENT '操作描述',
    ip VARCHAR(45) COMMENT 'IP地址',
    user_agent VARCHAR(500) COMMENT '用户代理',
    request_data TEXT COMMENT '请求数据',
    response_data TEXT COMMENT '响应数据',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_operator (operator),
    INDEX idx_action (action),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='操作日志表';

-- 插入初始分类数据
INSERT IGNORE INTO categories (id, name, parent_id, level, sort) VALUES
('cat-1', '道路交通安全法律、法规', NULL, 1, 1),
('cat-2', '交通信号', NULL, 1, 2),
('cat-3', '安全行车、文明驾驶基础知识', NULL, 1, 3),
('cat-4', '机动车驾驶操作相关基础知识', NULL, 1, 4),
('cat-5', '道路交通安全法律、法规和规章', 'cat-1', 2, 1),
('cat-6', '地方性法规', 'cat-1', 2, 2),
('cat-7', '道路交通信号灯', 'cat-2', 2, 1),
('cat-8', '道路交通标志', 'cat-2', 2, 2),
('cat-9', '道路交通标线', 'cat-2', 2, 3),
('cat-10', '交通警察手势', 'cat-2', 2, 4),
('cat-11', '安全行车常识', 'cat-3', 2, 1),
('cat-12', '文明驾驶常识', 'cat-3', 2, 2),
('cat-13', '机动车构造常识', 'cat-4', 2, 1),
('cat-14', '机动车检验常识', 'cat-4', 2, 2);

-- 插入示例题目数据
INSERT IGNORE INTO questions (id, title, content, options, correct_answer, explanation, difficulty, category_id) VALUES
('q-1', '机动车驾驶人违法驾驶造成重大交通事故构成犯罪的，依法追究什么责任？', 
 '机动车驾驶人违法驾驶造成重大交通事故构成犯罪的，依法追究什么责任？', 
 '["民事责任", "经济责任", "刑事责任", "直接责任"]', 
 2, 
 '根据《道路交通安全法》规定，违法驾驶造成重大交通事故构成犯罪的，依法追究刑事责任。', 
 'easy', 
 'cat-5'),

('q-2', '驾驶机动车在道路上违反道路交通安全法的行为，属于什么行为？', 
 '驾驶机动车在道路上违反道路交通安全法的行为，属于什么行为？', 
 '["违章行为", "违法行为", "过失行为", "违规行为"]', 
 1, 
 '违反道路交通安全法的行为属于违法行为，应当承担相应的法律责任。', 
 'easy', 
 'cat-5'),

('q-3', '这个标志是何含义？', 
 '这个标志是何含义？（假设显示禁止通行标志）', 
 '["禁止通行", "禁止驶入", "禁止车辆临时停放", "禁止车辆长时停放"]', 
 0, 
 '红色圆形标志表示禁止通行，任何车辆和行人都不得通过。', 
 'easy', 
 'cat-8');

-- 数据库初始化完成
SELECT 'Database initialization completed successfully!' as message;