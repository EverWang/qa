# 刷刷题后端API服务

## 项目简介

刷刷题是一个微信小程序答题系统的后端服务，基于Golang + Gin + GORM + MySQL开发，提供用户认证、题目管理、答题记录、错题本等功能。

## 技术栈

- **语言**: Go 1.21+
- **框架**: Gin
- **数据库**: MySQL 8.0+
- **ORM**: GORM
- **认证**: JWT
- **配置**: 环境变量

## 项目结构

```
server/
├── config/          # 配置文件
│   └── database.go  # 数据库配置
├── controllers/     # 控制器
│   ├── auth.go      # 认证控制器
│   ├── user.go      # 用户控制器
│   ├── question.go  # 题目控制器
│   ├── category.go  # 分类控制器
│   ├── answer.go    # 答题记录控制器
│   ├── mistake.go   # 错题本控制器
│   ├── statistics.go # 统计控制器
│   └── base.go      # 基础控制器
├── middleware/      # 中间件
│   └── middleware.go # JWT、CORS等中间件
├── models/          # 数据模型
│   └── models.go    # 数据库模型定义
├── routes/          # 路由
│   └── routes.go    # 路由配置
├── sql/             # SQL文件
│   └── init_database.sql # 数据库初始化脚本
├── .env             # 环境变量配置
├── go.mod           # Go模块文件
├── go.sum           # Go依赖锁定文件
├── main.go          # 程序入口
└── README.md        # 项目说明
```

## 快速开始

### 1. 环境准备

- Go 1.21+
- MySQL 8.0+
- Git

### 2. 克隆项目

```bash
git clone <repository-url>
cd qaminiprogram/server
```

### 3. 安装依赖

```bash
go mod download
```

### 4. 配置环境变量

复制 `.env` 文件并修改配置：

```bash
cp .env.example .env
```

修改 `.env` 文件中的配置：

```env
# 数据库配置
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=qaminiprogram

# JWT配置
JWT_SECRET=your_jwt_secret_key
JWT_EXPIRE_HOURS=24

# 微信小程序配置
WX_APPID=your_wechat_appid
WX_SECRET=your_wechat_secret

# 服务器配置
SERVER_PORT=8080
SERVER_MODE=debug
```

### 5. 初始化数据库

```bash
# 创建数据库
mysql -u root -p -e "CREATE DATABASE qaminiprogram DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# 导入初始化脚本
mysql -u root -p qaminiprogram < sql/init_database.sql
```

### 6. 运行服务

```bash
go run main.go
```

服务将在 `http://localhost:8080` 启动。

## API文档

### 基础信息

- **Base URL**: `http://localhost:8080/api/v1`
- **认证方式**: JWT Bearer Token
- **数据格式**: JSON

### 响应格式

#### 成功响应
```json
{
  "code": 200,
  "message": "success",
  "data": {}
}
```

#### 分页响应
```json
{
  "code": 200,
  "message": "success",
  "data": [],
  "total": 100,
  "page": 1,
  "size": 10
}
```

#### 错误响应
```json
{
  "code": 400,
  "message": "错误信息"
}
```

### 认证接口

#### 微信登录
```http
POST /auth/login
Content-Type: application/json

{
  "code": "微信登录code"
}
```

#### 刷新令牌
```http
POST /auth/refresh
Content-Type: application/json

{
  "token": "旧的JWT令牌"
}
```

### 用户接口

#### 获取用户信息
```http
GET /user/profile
Authorization: Bearer <token>
```

#### 更新用户信息
```http
PUT /user/profile
Authorization: Bearer <token>
Content-Type: application/json

{
  "nickname": "新昵称",
  "avatar": "头像URL"
}
```

### 分类接口

#### 获取分类列表
```http
GET /categories?tree=true&parent_id=0
```

#### 获取分类详情
```http
GET /categories/{id}
```

### 题目接口

#### 获取题目列表
```http
GET /questions?page=1&size=10&category_id=1&difficulty=medium&keyword=关键词
```

#### 获取题目详情
```http
GET /questions/{id}
```

#### 获取随机题目
```http
GET /questions/random?count=10&category_id=1&difficulty=medium
```

#### 根据分类获取题目
```http
GET /questions/category/{category_id}?page=1&size=10
```

### 答题接口

#### 提交答案
```http
POST /answers
Authorization: Bearer <token>
Content-Type: application/json

{
  "question_id": 1,
  "user_answer": 0,
  "time_spent": 30
}
```

#### 获取答题历史
```http
GET /answers/history?page=1&size=10&category_id=1&is_correct=true
Authorization: Bearer <token>
```

#### 获取答题统计
```http
GET /answers/statistics
Authorization: Bearer <token>
```

### 错题本接口

#### 获取错题本
```http
GET /mistakes?page=1&size=10&category_id=1&difficulty=medium
Authorization: Bearer <token>
```

#### 添加到错题本
```http
POST /mistakes
Authorization: Bearer <token>
Content-Type: application/json

{
  "question_id": 1
}
```

#### 从错题本移除
```http
DELETE /mistakes/{id}
Authorization: Bearer <token>
```

#### 清空错题本
```http
DELETE /mistakes/clear
Authorization: Bearer <token>
```

### 管理员接口

所有管理员接口都需要管理员权限，路径前缀为 `/admin`。

#### 管理员登录
```http
POST /admin/login
Content-Type: application/json

{
  "username": "admin",
  "password": "admin123"
}
```

#### 用户管理
```http
# 获取用户列表
GET /admin/users?page=1&size=10

# 获取用户详情
GET /admin/users/{id}

# 更新用户
PUT /admin/users/{id}

# 删除用户
DELETE /admin/users/{id}
```

#### 分类管理
```http
# 创建分类
POST /admin/categories

# 更新分类
PUT /admin/categories/{id}

# 删除分类
DELETE /admin/categories/{id}
```

#### 题目管理
```http
# 创建题目
POST /admin/questions

# 更新题目
PUT /admin/questions/{id}

# 删除题目
DELETE /admin/questions/{id}

# 批量导入题目
POST /admin/questions/import

# 导出题目
GET /admin/questions/export
```

#### 统计接口
```http
# 概览统计
GET /admin/statistics/overview

# 题目统计
GET /admin/statistics/questions?page=1&size=10&sort_by=total_answered&sort_order=DESC

# 用户统计
GET /admin/statistics/users?page=1&size=10&sort_by=total_answered&sort_order=DESC
```

## 数据库设计

### 主要表结构

- **users**: 用户表
- **categories**: 分类表
- **questions**: 题目表
- **answer_records**: 答题记录表
- **mistake_books**: 错题本表
- **admins**: 管理员表

详细的表结构请参考 `sql/init_database.sql` 文件。

## 开发指南

### 添加新的API接口

1. 在 `models/models.go` 中定义数据模型
2. 在 `controllers/` 目录下创建对应的控制器
3. 在 `routes/routes.go` 中添加路由
4. 更新API文档

### 中间件使用

- `middleware.CORS()`: 跨域处理
- `middleware.Logger()`: 日志记录
- `middleware.JWTAuth()`: JWT认证
- `middleware.AdminAuth()`: 管理员权限验证

### 错误处理

使用统一的错误响应格式，通过 `ErrorResponse()` 函数返回错误信息。

## 部署

### Docker部署

```dockerfile
# Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

EXPOSE 8080
CMD ["./main"]
```

### 生产环境配置

1. 设置 `SERVER_MODE=release`
2. 使用强密码和安全的JWT密钥
3. 配置HTTPS
4. 设置适当的CORS策略
5. 配置日志轮转
6. 设置数据库连接池

## 测试

### 单元测试

```bash
go test ./...
```

### API测试

可以使用Postman、curl或其他API测试工具进行测试。项目提供了完整的API文档和示例。

## 常见问题

### 1. 数据库连接失败

检查 `.env` 文件中的数据库配置是否正确，确保MySQL服务正在运行。

### 2. JWT令牌验证失败

确保请求头中包含正确的Authorization字段：`Bearer <token>`

### 3. 微信登录失败

检查微信小程序的AppID和AppSecret配置是否正确。

## 贡献指南

1. Fork项目
2. 创建功能分支
3. 提交更改
4. 推送到分支
5. 创建Pull Request

## 许可证

MIT License