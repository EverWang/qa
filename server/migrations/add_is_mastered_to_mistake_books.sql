-- 添加is_mastered字段到mistake_books表
ALTER TABLE mistake_books 
ADD COLUMN is_mastered BOOLEAN DEFAULT FALSE COMMENT '是否已掌握';

-- 添加updated_at字段到mistake_books表
ALTER TABLE mistake_books 
ADD COLUMN updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;

-- 更新现有记录的updated_at字段
UPDATE mistake_books SET updated_at = created_at WHERE updated_at IS NULL;