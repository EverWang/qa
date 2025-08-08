-- 刷刷题数据库初始化脚本
-- 创建数据库
CREATE DATABASE IF NOT EXISTS qaminiprogram DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE qaminiprogram;

-- 创建用户表
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    openid VARCHAR(100) UNIQUE NOT NULL COMMENT '微信openid',
    nickname VARCHAR(50) NOT NULL COMMENT '用户昵称',
    avatar VARCHAR(255) DEFAULT '' COMMENT '用户头像',
    role ENUM('user', 'admin') DEFAULT 'user' COMMENT '用户角色',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

-- 创建用户表索引
CREATE INDEX idx_users_openid ON users(openid);
CREATE INDEX idx_users_role ON users(role);

-- 创建分类表
CREATE TABLE categories (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL COMMENT '分类名称',
    parent_id INT DEFAULT NULL COMMENT '父分类ID',
    level TINYINT DEFAULT 1 COMMENT '分类层级',
    sort INT DEFAULT 0 COMMENT '排序',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (parent_id) REFERENCES categories(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='题目分类表';

-- 创建分类表索引
CREATE INDEX idx_categories_parent_id ON categories(parent_id);
CREATE INDEX idx_categories_level ON categories(level);

-- 初始化分类数据
INSERT INTO categories (name, parent_id, level, sort) VALUES
('法考题', NULL, 1, 1),
('医考题', NULL, 1, 2),
('工程类', NULL, 1, 3),
('法律基础', 1, 2, 1),
('法律条款', 1, 2, 2),
('法律解释', 1, 2, 3),
('临床医学', 2, 2, 1),
('基础医学', 2, 2, 2),
('道路工程', 3, 2, 1),
('桥隧工程', 3, 2, 2);

-- 创建题目表
CREATE TABLE questions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(200) NOT NULL COMMENT '题目标题',
    content TEXT NOT NULL COMMENT '题目内容',
    options JSON NOT NULL COMMENT '选项数组',
    correct_answer TINYINT NOT NULL COMMENT '正确答案索引',
    explanation TEXT DEFAULT '' COMMENT '答案解析',
    difficulty ENUM('easy', 'medium', 'hard') DEFAULT 'medium' COMMENT '难度等级',
    category_id INT NOT NULL COMMENT '分类ID',
    creator_id INT DEFAULT NULL COMMENT '创建者ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE,
    FOREIGN KEY (creator_id) REFERENCES users(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='题目表';

-- 创建题目表索引
CREATE INDEX idx_questions_category_id ON questions(category_id);
CREATE INDEX idx_questions_difficulty ON questions(difficulty);
CREATE INDEX idx_questions_creator_id ON questions(creator_id);
CREATE FULLTEXT INDEX idx_questions_content ON questions(title, content);

-- 创建答题记录表
CREATE TABLE answer_records (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL COMMENT '用户ID',
    question_id INT NOT NULL COMMENT '题目ID',
    user_answer TINYINT NOT NULL COMMENT '用户答案',
    is_correct BOOLEAN NOT NULL COMMENT '是否正确',
    time_spent INT DEFAULT 0 COMMENT '答题用时(秒)',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (question_id) REFERENCES questions(id) ON DELETE CASCADE,
    UNIQUE KEY uk_user_question (user_id, question_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='答题记录表';

-- 创建答题记录表索引
CREATE INDEX idx_answer_records_user_id ON answer_records(user_id);
CREATE INDEX idx_answer_records_question_id ON answer_records(question_id);
CREATE INDEX idx_answer_records_created_at ON answer_records(created_at);

-- 创建错题本表
CREATE TABLE mistake_books (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL COMMENT '用户ID',
    question_id INT NOT NULL COMMENT '题目ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (question_id) REFERENCES questions(id) ON DELETE CASCADE,
    UNIQUE KEY uk_user_question (user_id, question_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='错题本表';

-- 创建错题本表索引
CREATE INDEX idx_mistake_books_user_id ON mistake_books(user_id);
CREATE INDEX idx_mistake_books_created_at ON mistake_books(created_at);

-- 创建管理员表
CREATE TABLE admins (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL COMMENT '用户名',
    password_hash VARCHAR(255) NOT NULL COMMENT '密码哈希',
    email VARCHAR(100) DEFAULT '' COMMENT '邮箱',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='管理员表';

-- 创建管理员表索引
CREATE INDEX idx_admins_username ON admins(username);

-- 初始化管理员数据（密码：admin123）
INSERT INTO admins (username, password_hash, email) VALUES
('admin', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iKVjzieMwkOBSaEVFVGGvQTdO5wS', 'admin@example.com');

-- 插入示例题目数据
INSERT INTO questions (title, content, options, correct_answer, explanation, difficulty, category_id) VALUES
('道路工程基础知识', '下列关于道路工程的描述，正确的是？', '["道路工程只包括路面工程", "道路工程包括路基、路面、桥梁、隧道等", "道路工程不需要考虑排水系统", "道路工程只需要考虑车辆通行"]', 1, '道路工程是一个综合性工程，包括路基工程、路面工程、桥梁工程、隧道工程、排水工程等多个方面。', 'easy', 9),
('桥梁结构类型', '下列哪种桥梁结构适用于大跨度桥梁？', '["简支梁桥", "连续梁桥", "悬索桥", "拱桥"]', 2, '悬索桥是目前世界上跨度最大的桥梁结构形式，适用于大跨度桥梁建设。', 'medium', 10),
('法律基础概念', '下列关于法律效力的说法，正确的是？', '["法律只对公民有约束力", "法律对所有人都有约束力", "法律只对政府有约束力", "法律效力因人而异"]', 1, '法律面前人人平等，法律对所有人都具有约束力，包括公民、法人和其他组织。', 'easy', 4);