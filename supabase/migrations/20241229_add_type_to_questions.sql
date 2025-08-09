-- 添加type字段到questions表
ALTER TABLE questions ADD COLUMN IF NOT EXISTS type VARCHAR(20) DEFAULT 'single';

-- 添加类型约束
ALTER TABLE questions ADD CONSTRAINT questions_type_check 
  CHECK (type IN ('single', 'multiple', 'judge', 'fill'));

-- 更新现有数据的type字段（根据options字段判断）
UPDATE questions 
SET type = CASE 
  WHEN JSON_ARRAY_LENGTH(options) >= 2 THEN 'single'
  WHEN JSON_ARRAY_LENGTH(options) = 0 THEN 'judge'
  ELSE 'single'
END
WHERE type IS NULL OR type = '';

-- 修改options字段为可空
ALTER TABLE questions ALTER COLUMN options DROP NOT NULL;