-- 添加游客支持的数据库迁移
-- 添加 is_guest 字段
ALTER TABLE users ADD COLUMN is_guest BOOLEAN DEFAULT FALSE;

-- 更新 role 字段的枚举值以包含 guest
ALTER TABLE users MODIFY COLUMN role ENUM('user', 'admin', 'guest') DEFAULT 'user';

-- 为 is_guest 字段添加索引以提高查询性能
CREATE INDEX idx_users_is_guest ON users(is_guest);