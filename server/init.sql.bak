-- 刷刷题数据库初始化脚本
-- 创建数据库和表结构，插入初始数据

-- 创建数据库（如果不存在）
CREATE DATABASE IF NOT EXISTS shuashuati DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE shuashuati;

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    openid VARCHAR(100) NOT NULL UNIQUE COMMENT '微信OpenID',
    nickname VARCHAR(100) DEFAULT '' COMMENT '用户昵称',
    avatar VARCHAR(500) DEFAULT '' COMMENT '头像URL',
    role ENUM('user', 'admin') DEFAULT 'user' COMMENT '用户角色',
    total_answered INT DEFAULT 0 COMMENT '总答题数',
    total_correct INT DEFAULT 0 COMMENT '总正确数',
    accuracy_rate DECIMAL(5,2) DEFAULT 0.00 COMMENT '正确率',
    last_active_time TIMESTAMP NULL COMMENT '最后活跃时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_openid (openid),
    INDEX idx_role (role),
    INDEX idx_last_active (last_active_time)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- 分类表
CREATE TABLE IF NOT EXISTS categories (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL COMMENT '分类名称',
    parent_id BIGINT NULL COMMENT '父分类ID',
    level INT DEFAULT 1 COMMENT '分类层级',
    sort INT DEFAULT 0 COMMENT '排序权重',
    question_count INT DEFAULT 0 COMMENT '题目数量',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_parent_id (parent_id),
    INDEX idx_level (level),
    INDEX idx_sort (sort),
    FOREIGN KEY (parent_id) REFERENCES categories(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='分类表';

-- 题目表
CREATE TABLE IF NOT EXISTS questions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(500) NOT NULL COMMENT '题目标题',
    content TEXT NOT NULL COMMENT '题目内容',
    options JSON NOT NULL COMMENT '选项数组',
    correct_answer INT NOT NULL COMMENT '正确答案索引',
    explanation TEXT COMMENT '答案解析',
    difficulty ENUM('easy', 'medium', 'hard') DEFAULT 'medium' COMMENT '难度等级',
    category_id BIGINT NOT NULL COMMENT '分类ID',
    total_answered INT DEFAULT 0 COMMENT '总答题次数',
    total_correct INT DEFAULT 0 COMMENT '总正确次数',
    accuracy_rate DECIMAL(5,2) DEFAULT 0.00 COMMENT '正确率',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_category_id (category_id),
    INDEX idx_difficulty (difficulty),
    INDEX idx_accuracy_rate (accuracy_rate),
    INDEX idx_total_answered (total_answered),
    FULLTEXT idx_title_content (title, content),
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='题目表';

-- 答题记录表
CREATE TABLE IF NOT EXISTS answer_records (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL COMMENT '用户ID',
    question_id BIGINT NOT NULL COMMENT '题目ID',
    user_answer INT NOT NULL COMMENT '用户答案',
    is_correct BOOLEAN NOT NULL COMMENT '是否正确',
    time_spent INT DEFAULT 0 COMMENT '答题耗时(秒)',
    answered_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '答题时间',
    INDEX idx_user_id (user_id),
    INDEX idx_question_id (question_id),
    INDEX idx_is_correct (is_correct),
    INDEX idx_answered_at (answered_at),
    INDEX idx_user_question (user_id, question_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (question_id) REFERENCES questions(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='答题记录表';

-- 错题本表
CREATE TABLE IF NOT EXISTS mistake_books (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL COMMENT '用户ID',
    question_id BIGINT NOT NULL COMMENT '题目ID',
    added_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
    UNIQUE KEY uk_user_question (user_id, question_id),
    INDEX idx_user_id (user_id),
    INDEX idx_question_id (question_id),
    INDEX idx_added_at (added_at),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (question_id) REFERENCES questions(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='错题本表';

-- 插入初始分类数据
INSERT INTO categories (id, name, parent_id, level, sort) VALUES
(1, '道路交通安全法律、法规', NULL, 1, 1),
(2, '交通信号', NULL, 1, 2),
(3, '安全行车、文明驾驶基础知识', NULL, 1, 3),
(4, '机动车驾驶操作相关基础知识', NULL, 1, 4),
(5, '道路交通安全法律、法规和规章', 1, 2, 1),
(6, '地方性法规', 1, 2, 2),
(7, '道路交通信号灯', 2, 2, 1),
(8, '道路交通标志', 2, 2, 2),
(9, '道路交通标线', 2, 2, 3),
(10, '交通警察手势', 2, 2, 4),
(11, '安全行车常识', 3, 2, 1),
(12, '文明驾驶常识', 3, 2, 2),
(13, '机动车构造常识', 4, 2, 1),
(14, '机动车检验常识', 4, 2, 2);

-- 更新分类的题目数量（初始为0）
UPDATE categories SET question_count = 0;

-- 管理员表
CREATE TABLE IF NOT EXISTS admins (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(50) NOT NULL UNIQUE COMMENT '用户名',
    password_hash VARCHAR(255) NOT NULL COMMENT '密码哈希',
    email VARCHAR(100) DEFAULT '' COMMENT '邮箱',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_username (username)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='管理员表';

-- 插入默认管理员账号（用户名：admin，密码：123456）
-- 密码哈希是 123456 的 bcrypt 哈希值
INSERT INTO admins (username, password_hash, email, created_at) VALUES
('admin', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'admin@example.com', NOW());

-- 插入管理员用户到users表（用于微信登录）
INSERT INTO users (openid, nickname, avatar, role, created_at) VALUES
('admin_openid_123456', '系统管理员', '', 'admin', NOW());

-- 插入示例题目数据
INSERT INTO questions (title, content, options, correct_answer, explanation, difficulty, category_id) VALUES
('机动车驾驶人违法驾驶造成重大交通事故构成犯罪的，依法追究什么责任？', 
 '机动车驾驶人违法驾驶造成重大交通事故构成犯罪的，依法追究什么责任？', 
 '["民事责任", "经济责任", "刑事责任", "直接责任"]', 
 2, 
 '根据《道路交通安全法》规定，违法驾驶造成重大交通事故构成犯罪的，依法追究刑事责任。', 
 'easy', 
 5),

('驾驶机动车在道路上违反道路交通安全法的行为，属于什么行为？', 
 '驾驶机动车在道路上违反道路交通安全法的行为，属于什么行为？', 
 '["违章行为", "违法行为", "过失行为", "违规行为"]', 
 1, 
 '违反道路交通安全法的行为属于违法行为，应当承担相应的法律责任。', 
 'easy', 
 5),

('机动车驾驶证有效期满换证时，需要提交什么？', 
 '机动车驾驶证有效期满换证时，需要提交什么？', 
 '["机动车登记证书", "机动车保险单", "身体条件证明", "机动车技术检验合格证明"]', 
 2, 
 '驾驶证有效期满换证时，需要提交身体条件证明，确保驾驶人身体状况符合驾驶要求。', 
 'medium', 
 5),

('这个标志是何含义？', 
 '这个标志是何含义？（假设显示禁止通行标志）', 
 '["禁止通行", "禁止驶入", "禁止车辆临时停放", "禁止车辆长时停放"]', 
 0, 
 '红色圆形标志表示禁止通行，任何车辆和行人都不得通过。', 
 'easy', 
 8),

('这个标志是何含义？', 
 '这个标志是何含义？（假设显示限速标志）', 
 '["最低限速60公里/小时", "最高限速60公里/小时", "解除限速60公里/小时", "建议速度60公里/小时"]', 
 1, 
 '圆形标志内显示数字表示最高限速，此标志表示最高限速60公里/小时。', 
 'medium', 
 8),

('在冰雪道路上行车时，车辆的稳定性降低，加速过急时车轮极易空转或溜滑。', 
 '在冰雪道路上行车时，车辆的稳定性降低，加速过急时车轮极易空转或溜滑。', 
 '["正确", "错误"]', 
 0, 
 '在冰雪道路上，路面摩擦力减小，车辆稳定性确实会降低，加速过急容易导致车轮空转或溜滑。', 
 'easy', 
 11),

('驾驶机动车在雾天行车开启雾灯和危险报警闪光灯。', 
 '驾驶机动车在雾天行车开启雾灯和危险报警闪光灯。', 
 '["正确", "错误"]', 
 0, 
 '雾天行车应当开启雾灯和危险报警闪光灯，提高车辆的可见性，确保行车安全。', 
 'easy', 
 11),

('机动车仪表板上（如图所示）亮表示什么？', 
 '机动车仪表板上（如图所示）亮表示什么？（假设显示发动机故障灯）', 
 '["发动机故障", "燃油不足", "充电故障", "冷却液不足"]', 
 0, 
 '发动机故障指示灯亮起表示发动机系统出现故障，应及时检修。', 
 'medium', 
 13),

('机动车在高速公路上发生故障时，在来车方向50至100米处设置警告标志。', 
 '机动车在高速公路上发生故障时，在来车方向50至100米处设置警告标志。', 
 '["正确", "错误"]', 
 1, 
 '在高速公路上发生故障时，应在来车方向150米以外设置警告标志，而不是50至100米。', 
 'hard', 
 11),

('申请机动车驾驶证的人，应当符合国务院公安部门规定的驾驶许可条件。', 
 '申请机动车驾驶证的人，应当符合国务院公安部门规定的驾驶许可条件。', 
 '["正确", "错误"]', 
 0, 
 '根据《道路交通安全法》规定，申请驾驶证的人应当符合国务院公安部门规定的驾驶许可条件。', 
 'easy', 
 5);

-- 更新分类的题目数量
UPDATE categories c SET question_count = (
    SELECT COUNT(*) FROM questions q WHERE q.category_id = c.id
);

-- 创建触发器：插入题目时更新分类题目数量
DELIMITER //
CREATE TRIGGER tr_questions_insert_update_category_count
AFTER INSERT ON questions
FOR EACH ROW
BEGIN
    UPDATE categories SET question_count = question_count + 1 WHERE id = NEW.category_id;
END//

-- 创建触发器：删除题目时更新分类题目数量
CREATE TRIGGER tr_questions_delete_update_category_count
AFTER DELETE ON questions
FOR EACH ROW
BEGIN
    UPDATE categories SET question_count = question_count - 1 WHERE id = OLD.category_id;
END//

-- 创建触发器：更新题目分类时更新分类题目数量
CREATE TRIGGER tr_questions_update_category_count
AFTER UPDATE ON questions
FOR EACH ROW
BEGIN
    IF OLD.category_id != NEW.category_id THEN
        UPDATE categories SET question_count = question_count - 1 WHERE id = OLD.category_id;
        UPDATE categories SET question_count = question_count + 1 WHERE id = NEW.category_id;
    END IF;
END//

-- 创建触发器：插入答题记录时更新题目和用户统计
CREATE TRIGGER tr_answer_records_insert_update_stats
AFTER INSERT ON answer_records
FOR EACH ROW
BEGIN
    -- 更新题目统计
    UPDATE questions SET 
        total_answered = total_answered + 1,
        total_correct = total_correct + IF(NEW.is_correct, 1, 0),
        accuracy_rate = ROUND((total_correct + IF(NEW.is_correct, 1, 0)) * 100.0 / (total_answered + 1), 2)
    WHERE id = NEW.question_id;
    
    -- 更新用户统计
    UPDATE users SET 
        total_answered = total_answered + 1,
        total_correct = total_correct + IF(NEW.is_correct, 1, 0),
        accuracy_rate = ROUND((total_correct + IF(NEW.is_correct, 1, 0)) * 100.0 / (total_answered + 1), 2),
        last_active_time = NEW.answered_at
    WHERE id = NEW.user_id;
END//

DELIMITER ;

-- 创建视图：用户答题统计视图
CREATE VIEW v_user_answer_stats AS
SELECT 
    u.id,
    u.openid,
    u.nickname,
    u.avatar,
    u.total_answered,
    u.total_correct,
    u.accuracy_rate,
    u.last_active_time,
    COUNT(DISTINCT ar.question_id) as unique_questions_answered,
    COUNT(mb.id) as mistake_count,
    DATE(u.created_at) as join_date
FROM users u
LEFT JOIN answer_records ar ON u.id = ar.user_id
LEFT JOIN mistake_books mb ON u.id = mb.user_id
WHERE u.role = 'user'
GROUP BY u.id;

-- 创建视图：题目答题统计视图
CREATE VIEW v_question_answer_stats AS
SELECT 
    q.id,
    q.title,
    q.difficulty,
    q.category_id,
    c.name as category_name,
    q.total_answered,
    q.total_correct,
    q.accuracy_rate,
    COUNT(mb.id) as mistake_count,
    q.created_at
FROM questions q
LEFT JOIN categories c ON q.category_id = c.id
LEFT JOIN mistake_books mb ON q.id = mb.question_id
GROUP BY q.id;

-- 创建视图：分类统计视图
CREATE VIEW v_category_stats AS
SELECT 
    c.id,
    c.name,
    c.parent_id,
    c.level,
    c.question_count,
    COALESCE(SUM(q.total_answered), 0) as total_answered,
    COALESCE(SUM(q.total_correct), 0) as total_correct,
    CASE 
        WHEN SUM(q.total_answered) > 0 THEN ROUND(SUM(q.total_correct) * 100.0 / SUM(q.total_answered), 2)
        ELSE 0.00
    END as accuracy_rate
FROM categories c
LEFT JOIN questions q ON c.id = q.category_id
GROUP BY c.id;

-- 插入一些示例答题记录（可选）
-- INSERT INTO answer_records (user_id, question_id, user_answer, is_correct, time_spent) VALUES
-- (1, 1, 2, true, 15),
-- (1, 2, 1, true, 20),
-- (1, 3, 0, false, 25);

-- 创建索引优化查询性能
CREATE INDEX idx_answer_records_user_answered_at ON answer_records(user_id, answered_at);
CREATE INDEX idx_answer_records_question_answered_at ON answer_records(question_id, answered_at);
CREATE INDEX idx_users_role_last_active ON users(role, last_active_time);
CREATE INDEX idx_questions_category_difficulty ON questions(category_id, difficulty);

-- 数据库初始化完成
SELECT 'Database initialization completed successfully!' as message;