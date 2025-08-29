# 项目开发规则 - 刷刷题容器化部署

## 概述
本文档定义了刷刷题项目在容器化部署过程中需要遵循的技术规则和开发规范。

## 项目架构规则

### 1. 技术栈规范

#### 1.1 后端技术栈
- **语言**：Go 1.23.4+
- **框架**：Gin
- **数据库**：MySQL 8.0+
- **ORM**：GORM
- **容器化**：Podman + podman-compose

#### 1.2 前端技术栈
- **管理端**：Vue3 + TypeScript + Element Plus
- **用户端**：Vue3 + TypeScript + Uni-app
- **构建工具**：Vite
- **样式**：SCSS + Tailwind CSS

#### 1.3 部署环境
- **开发环境**：WSL2 + Ubuntu 22.04
- **容器运行时**：Podman 4.9.3+
- **编排工具**：podman-compose 1.0.6+

### 2. 项目结构规则

#### 2.1 目录结构
```
qa/
├── server/                 # 后端Go服务
│   ├── main.go            # 主程序入口
│   ├── main-linux         # Linux二进制文件
│   ├── Dockerfile         # 后端容器配置
│   └── ...
├── admin-web/             # 管理端前端
│   ├── dist/              # 构建产物
│   ├── nginx.conf         # Nginx配置
│   ├── Dockerfile         # 前端容器配置
│   └── ...
├── miniprogram/           # 用户端前端
│   ├── dist/              # 构建产物
│   ├── nginx.conf         # Nginx配置
│   ├── Dockerfile         # 前端容器配置
│   └── ...
├── docker-compose.yml     # 容器编排配置
├── user_rules.md          # 用户操作规则
└── project_rules.md       # 项目开发规则
```

#### 2.2 命名规范
- **容器名称**：使用`shuashuati_`前缀
- **镜像名称**：使用`qa_`前缀
- **网络名称**：使用`qa_shuashuati_network`
- **数据卷名称**：使用`qa_`前缀

### 3. 容器化规则

#### 3.1 Dockerfile规范

##### 后端Dockerfile规则
```dockerfile
# 使用scratch基础镜像（最小化）
FROM scratch

# 复制Linux二进制文件
COPY main-linux /main

# 暴露端口
EXPOSE 8080

# 启动应用
CMD ["/main"]
```

**关键规则**：
- 必须使用Linux编译的二进制文件
- 使用scratch镜像减小镜像大小
- 端口号通过环境变量配置

##### 前端Dockerfile规则
```dockerfile
# 使用nginx:alpine基础镜像
FROM docker.m.daocloud.io/library/nginx:alpine

# 复制构建产物
COPY dist /usr/share/nginx/html

# 复制nginx配置
COPY nginx.conf /etc/nginx/nginx.conf

# 暴露端口
EXPOSE 80

# 启动nginx
CMD ["nginx", "-g", "daemon off;"]
```

**关键规则**：
- 使用已构建的dist目录
- 使用alpine镜像减小大小
- 配置nginx代理到后端服务

#### 3.2 docker-compose.yml规范

##### 服务定义规则
```yaml
version: '3.8'

services:
  # MySQL数据库
  mysql:
    image: docker.m.daocloud.io/library/mysql:8.0
    container_name: shuashuati_mysql
    environment:
      MYSQL_ROOT_PASSWORD: 123456
      MYSQL_DATABASE: shuashuati
      MYSQL_USER: shuashuati
      MYSQL_PASSWORD: shuashuati123
    ports:
      - "3306:3306"
    networks:
      - shuashuati_network
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost"]
      interval: 30s
      timeout: 10s
      retries: 5

  # 后端服务
  backend:
    build:
      context: ./server
      dockerfile: Dockerfile
    container_name: shuashuati_backend
    environment:
      PORT: 8080
      DB_HOST: mysql
      DB_PORT: 3306
      DB_USER: shuashuati
      DB_PASSWORD: shuashuati123
      DB_NAME: shuashuati
    ports:
      - "8080:8080"
    networks:
      - shuashuati_network
```

**关键规则**：
- 避免使用复杂的depends_on条件
- 使用环境变量传递配置
- 统一使用自定义网络
- 设置健康检查（MySQL）

### 4. 网络配置规则

#### 4.1 容器间通信
- 所有服务必须在同一个自定义网络中
- 使用服务名作为主机名进行通信
- 后端服务名固定为`backend`
- 数据库服务名固定为`mysql`

#### 4.2 端口映射
- MySQL：3306:3306
- 后端API：8080:8080
- 用户端前端：3000:80
- 管理端前端：3001:80

#### 4.3 Nginx代理配置
```nginx
# API代理配置
location /api/ {
    proxy_pass http://backend:8080;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
}
```

### 5. 环境变量规则

#### 5.1 后端环境变量
```yaml
environment:
  # 服务器配置
  PORT: 8080                    # 服务端口
  SERVER_PORT: 8080             # 备用端口配置
  SERVER_MODE: release          # 运行模式
  GIN_MODE: release             # Gin模式
  
  # 数据库配置
  DB_HOST: mysql                # 数据库主机
  DB_PORT: 3306                 # 数据库端口
  DB_USER: shuashuati           # 数据库用户
  DB_PASSWORD: shuashuati123    # 数据库密码
  DB_NAME: shuashuati           # 数据库名
  DB_CHARSET: utf8mb4           # 字符集
  
  # JWT配置
  JWT_SECRET: shuashuati_jwt_secret_key_2024
  JWT_EXPIRE_HOURS: 24
  
  # CORS配置
  ALLOWED_ORIGINS: "http://localhost:5173,http://localhost:5174,http://localhost:3000,http://localhost:3001"
```

#### 5.2 环境变量优先级
1. PORT（最高优先级）
2. SERVER_PORT
3. 代码中的默认值（最低优先级）

### 6. 代码规范

#### 6.1 Go代码规范
```go
// 端口配置示例
func main() {
    // 获取端口配置，优先级：PORT > SERVER_PORT > 默认值
    port := os.Getenv("PORT")
    if port == "" {
        port = os.Getenv("SERVER_PORT")
        if port == "" {
            port = "8080"
        }
    }
    
    log.Printf("Server starting on port %s", port)
    if err := r.Run(":" + port); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}
```

**关键规则**：
- 避免硬编码端口号
- 使用环境变量配置
- 添加详细的日志输出
- 实现优雅的错误处理

#### 6.2 前端代码规范
- 使用TypeScript进行类型检查
- 遵循Vue3 Composition API规范
- 使用Tailwind CSS进行样式管理
- API请求统一使用相对路径`/api/`

### 7. 构建规则

#### 7.1 后端构建
```bash
# Linux二进制文件编译
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main-linux .

# 镜像构建
podman build -t qa_backend ./server
```

#### 7.2 前端构建
```bash
# 构建前端资源
npm run build

# 镜像构建
podman build -t qa_admin-web ./admin-web
podman build -t qa_miniprogram ./miniprogram
```

### 8. 部署流程规则

#### 8.1 标准部署流程
1. **环境检查**：确认WSL和Podman环境
2. **代码编译**：编译Linux版本的后端二进制文件
3. **镜像构建**：构建所有服务的Docker镜像
4. **服务启动**：使用podman-compose启动所有服务
5. **健康检查**：验证所有服务正常运行
6. **功能测试**：执行端到端功能测试

#### 8.2 部署命令序列
```bash
# 1. 编译后端
cd server
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main-linux .

# 2. 构建镜像
cd ..
podman build -t qa_backend ./server
podman build -t qa_admin-web ./admin-web
podman build -t qa_miniprogram ./miniprogram

# 3. 启动服务
podman-compose up -d

# 4. 检查状态
podman ps
podman logs shuashuati_backend
```

### 9. 监控和日志规则

#### 9.1 日志管理
- 使用结构化日志格式
- 记录关键操作和错误信息
- 定期清理日志文件

#### 9.2 健康检查
- 后端：`GET /health`
- 前端：HTTP 200响应
- 数据库：mysqladmin ping

#### 9.3 监控指标
- 容器运行状态
- 内存和CPU使用率
- 网络连接状态
- 应用响应时间

### 10. 故障恢复规则

#### 10.1 服务重启
```bash
# 重启单个服务
podman restart shuashuati_backend

# 重启所有服务
podman-compose restart

# 完全重新部署
podman-compose down
podman rm -af
podman-compose up -d
```

#### 10.2 数据备份
- 定期备份MySQL数据
- 备份重要配置文件
- 版本控制所有代码和配置

### 11. 安全规则

#### 11.1 网络安全
- 使用内部网络进行容器间通信
- 只暴露必要的端口
- 配置防火墙规则

#### 11.2 数据安全
- 使用强密码
- 定期更新密钥
- 加密敏感数据传输

### 12. 性能优化规则

#### 12.1 镜像优化
- 使用多阶段构建
- 选择合适的基础镜像
- 清理不必要的文件

#### 12.2 运行时优化
- 设置合适的资源限制
- 使用缓存机制
- 优化数据库查询

## 版本控制规则

### 1. Git规范
- 使用语义化版本号
- 编写清晰的提交信息
- 使用分支管理功能开发

### 2. 发布规范
- 标记重要版本
- 维护更新日志
- 测试后再发布

## 文档维护规则

### 1. 文档更新
- 及时更新技术文档
- 记录重要变更
- 维护操作手册

### 2. 知识沉淀
- 总结最佳实践
- 记录常见问题
- 分享经验教训

---

*本文档是项目开发和部署的技术规范，所有团队成员都应严格遵循。如有疑问或建议，请及时沟通和更新。*