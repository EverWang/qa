-- 创建操作日志表
CREATE TABLE IF NOT EXISTS operation_logs (
    id INT AUTO_INCREMENT PRIMARY KEY,
    operator VARCHAR(50) NOT NULL COMMENT '操作人',
    action VARCHAR(20) NOT NULL COMMENT '操作类型',
    resource VARCHAR(50) NOT NULL COMMENT '操作资源',
    description VARCHAR(255) COMMENT '操作描述',
    ip VARCHAR(45) COMMENT 'IP地址',
    user_agent VARCHAR(500) COMMENT '用户代理',
    request_data TEXT COMMENT '请求数据',
    response_data TEXT COMMENT '响应数据',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX idx_operator (operator),
    INDEX idx_action (action),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='操作日志表';