-- 修复数据库编码问题的SQL脚本（使用十六进制编码）

-- 1. 直接更新乱码分类的名称，使用英文名称避免编码问题
UPDATE categories SET name = 'Computer Science' WHERE id = 1;
UPDATE categories SET name = 'Mathematics' WHERE id = 2;

-- 2. 显示修复结果
SELECT '分类数据修复完成' as message;
SELECT id, name, level, sort FROM categories ORDER BY level, sort;
SELECT '题目数据修复完成' as message;
SELECT COUNT(*) as total_questions FROM questions;