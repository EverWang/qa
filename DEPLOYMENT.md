# 刷刷题项目部署指南

## 项目概述

本项目包含三个主要部分：
- **小程序用户端** (`miniprogram/`) - Vue3 + TypeScript + uni-app
- **管理端Web应用** (`admin-web/`) - Vue3 + TypeScript + Element Plus
- **后端服务** (`server/`) - Golang + Gin + MySQL

## 文件大小分析结果

根据分析，项目各目录大小如下：
- `admin-web/`: 255.26 MB (主要是node_modules)
- `miniprogram/`: 174.66 MB (主要是node_modules)
- `node_modules/`: 167.21 MB (根目录依赖)
- `server/`: 15.34 MB (Go后端代码)
- 题库文档: 3.01 MB

**问题分析**: 项目总大小超过612MB，远超Vercel 10MB限制，主要原因是多个node_modules目录。

## 部署策略

### 方案一：分离部署（推荐）

#### 1. 小程序端部署

**特点**: 小程序需要通过微信开发者工具发布，不需要传统Web部署。

**步骤**:
```bash
# 1. 进入小程序目录
cd miniprogram

# 2. 安装依赖
npm install

# 3. 构建项目
npm run build:mp-weixin

# 4. 使用微信开发者工具打开 dist/build/mp-weixin 目录
# 5. 在微信开发者工具中预览和发布
```

#### 2. 管理端Web应用部署

**Vercel部署**:
```bash
# 1. 进入管理端目录
cd admin-web

# 2. 安装Vercel CLI
npm install -g vercel

# 3. 登录Vercel
vercel login

# 4. 部署
vercel --prod
```

**Netlify部署**:
```bash
# 1. 构建项目
cd admin-web
npm run build

# 2. 将 dist/ 目录拖拽到 Netlify 部署页面
# 或使用 Netlify CLI
npm install -g netlify-cli
netlify deploy --prod --dir=dist
```

**传统服务器部署**:
```bash
# 1. 构建项目
cd admin-web
npm run build

# 2. 将 dist/ 目录上传到服务器
# 3. 配置Nginx
```

#### 3. 后端服务部署

**本地运行**:
```bash
# 1. 确保MySQL运行在 localhost:3306
# 2. 进入服务端目录
cd server

# 3. 安装依赖
go mod tidy

# 4. 运行服务
go run main.go
```

**云服务器部署**:
```bash
# 1. 编译
cd server
go build -o qaminiprogram-server main.go

# 2. 上传到服务器并运行
./qaminiprogram-server
```

### 方案二：Docker容器化部署

#### 1. 创建Docker配置

**后端Dockerfile** (`server/Dockerfile`):
```dockerfile
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
EXPOSE 8080
CMD ["./main"]
```

**管理端Dockerfile** (`admin-web/Dockerfile`):
```dockerfile
FROM node:18-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production
COPY . .
RUN npm run build

FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

**Docker Compose** (`docker-compose.yml`):
```yaml
version: '3.8'
services:
  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: 123456
      MYSQL_DATABASE: qaminiprogram
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql

  backend:
    build: ./server
    ports:
      - "8080:8080"
    depends_on:
      - mysql
    environment:
      DB_HOST: mysql
      DB_PORT: 3306
      DB_USER: root
      DB_PASSWORD: 123456
      DB_NAME: qaminiprogram

  admin-web:
    build: ./admin-web
    ports:
      - "80:80"
    depends_on:
      - backend

volumes:
  mysql_data:
```

#### 2. 运行Docker部署
```bash
# 1. 构建并启动所有服务
docker-compose up -d

# 2. 查看服务状态
docker-compose ps

# 3. 查看日志
docker-compose logs -f
```

## 本地开发环境搭建

### 1. 环境要求
- Node.js 18+
- Go 1.21+
- MySQL 8.0+
- 微信开发者工具

### 2. 数据库初始化
```sql
-- 创建数据库
CREATE DATABASE qaminiprogram CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 导入表结构
-- 运行 server/sql/ 目录下的SQL文件
```

### 3. 后端服务启动
```bash
cd server
go mod tidy
go run main.go
```

### 4. 管理端启动
```bash
cd admin-web
npm install
npm run dev
```

### 5. 小程序端启动
```bash
cd miniprogram
npm install
npm run dev:mp-weixin
# 使用微信开发者工具打开 dist/dev/mp-weixin 目录
```

## 生产环境配置

### 1. 环境变量配置

**后端环境变量** (`.env`):
```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=123456
DB_NAME=qaminiprogram
JWT_SECRET=your-jwt-secret
WECHAT_APP_ID=your-wechat-app-id
WECHAT_APP_SECRET=your-wechat-app-secret
```

**前端环境变量** (`.env.production`):
```env
VITE_API_BASE_URL=https://your-api-domain.com
VITE_APP_TITLE=刷刷题管理系统
```

### 2. Nginx配置示例
```nginx
server {
    listen 80;
    server_name your-domain.com;
    
    # 管理端静态文件
    location / {
        root /var/www/admin-web/dist;
        try_files $uri $uri/ /index.html;
    }
    
    # API代理
    location /api/ {
        proxy_pass http://localhost:8080/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
```

## 故障排除

### 1. 部署文件过大问题
- **原因**: node_modules目录过大
- **解决**: 使用优化后的.vercelignore文件，或采用分离部署策略

### 2. 跨域问题
- **原因**: 前后端域名不同
- **解决**: 配置CORS或使用代理

### 3. 数据库连接问题
- **检查**: 数据库服务是否启动
- **检查**: 连接参数是否正确
- **检查**: 防火墙设置

### 4. 微信小程序发布问题
- **检查**: AppID配置是否正确
- **检查**: 域名是否已备案并配置到微信后台
- **检查**: HTTPS证书是否有效

## 监控和维护

### 1. 日志监控
```bash
# 查看后端日志
tail -f server/logs/app.log

# 查看Nginx日志
tail -f /var/log/nginx/access.log
```

### 2. 性能监控
- 使用PM2管理Node.js进程
- 配置数据库慢查询日志
- 设置服务器资源监控

### 3. 备份策略
```bash
# 数据库备份
mysqldump -u root -p qaminiprogram > backup_$(date +%Y%m%d).sql

# 代码备份
tar -czf code_backup_$(date +%Y%m%d).tar.gz /path/to/project
```

## 联系支持

如遇到部署问题，请检查：
1. 环境要求是否满足
2. 配置文件是否正确
3. 网络连接是否正常
4. 日志文件中的错误信息

---

**注意**: 本指南基于当前项目结构编写，实际部署时请根据具体环境调整配置参数。