# 刷刷题系统

一个基于微信小程序的在线刷题系统，包含用户端小程序、管理端Web应用和后端API服务。

## 项目概述

刷刷题系统是一个完整的在线学习平台，支持多种题型的在线练习和管理。系统采用前后端分离架构，提供了丰富的功能和良好的用户体验。

### 主要功能

**用户端（微信小程序）**
- 用户注册登录
- 题目分类浏览
- 在线答题练习
- 错题本管理
- 学习记录统计
- 题目搜索
- 社交分享

**管理端（Web应用）**
- 管理员登录
- 题目批量导入
- 题目分类管理
- 用户管理
- 数据统计分析
- 系统配置

**后端服务**
- RESTful API接口
- 用户认证授权
- 数据库管理
- 文件上传处理
- 微信小程序集成

## 技术栈

### 前端技术
- **小程序端**: Vue 3 + TypeScript + uni-app + uni-ui
- **管理端**: Vue 3 + TypeScript + Element Plus + Vite
- **样式**: SCSS

### 后端技术
- **框架**: Go + Gin
- **数据库**: MySQL 8.0
- **ORM**: GORM
- **认证**: JWT
- **文档**: Swagger

### 部署技术
- **容器化**: Docker + Docker Compose
- **反向代理**: Nginx
- **云部署**: Vercel（可选）

## 快速开始

### Docker部署（推荐）

推荐使用Docker进行快速部署，详细说明请参考 [Docker部署文档](./DOCKER_DEPLOYMENT.md)。

#### 环境要求
- Docker Engine 20.10+
- Docker Compose 2.0+
- 至少4GB内存和10GB存储空间

#### 快速部署

**Linux/macOS:**
```bash
# 给脚本执行权限
chmod +x deploy.sh

# 启动所有服务
./deploy.sh start
```

**Windows:**
```cmd
# 启动所有服务
deploy.bat start
```

#### 手动部署
```bash
# 构建并启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f
```

### 本地开发

#### 环境要求
- Node.js 18+
- Go 1.21+
- MySQL 8.0+

#### 1. 克隆项目
```bash
git clone <repository-url>
cd qaminiprogram
```

#### 2. 启动数据库
```bash
# 使用Docker启动MySQL
docker run -d \
  --name mysql-dev \
  -e MYSQL_ROOT_PASSWORD=123456 \
  -e MYSQL_DATABASE=qaminiprogram \
  -p 3306:3306 \
  mysql:8.0
```

#### 3. 启动后端服务
```bash
cd server

# 安装依赖
go mod tidy

# 配置环境变量
cp .env.example .env
# 编辑 .env 文件，配置数据库连接等信息

# 启动服务
go run main.go
```

#### 4. 启动管理端
```bash
cd admin-web

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

#### 5. 启动小程序端
```bash
cd miniprogram

# 安装依赖
npm install

# 启动开发服务器
npm run dev:mp-weixin
```

## 访问地址

部署完成后，可以通过以下地址访问系统：

- **小程序端**: http://localhost
- **管理端**: http://localhost/admin
- **API文档**: http://localhost/api/docs
- **API接口**: http://localhost/api/v1

## 默认账号

**管理员账号**
- 用户名：admin
- 密码：123456

**测试用户**
- 用户名：testuser
- 密码：123456

## 常见问题

### Q: 如何修改数据库连接配置？
A: 修改 `server/.env` 文件中的数据库连接参数，或在Docker部署时修改 `docker-compose.yml` 中的环境变量。

### Q: 如何添加新的API接口？
A: 在 `server/controllers/` 中添加控制器，在 `server/api/` 中注册路由，并更新Swagger文档。

### Q: 如何自定义小程序样式？
A: 修改 `miniprogram/src/` 目录下的SCSS文件，或在 `miniprogram/static/` 中添加自定义样式。

### Q: 如何配置微信小程序？
A: 在微信开发者工具中配置AppID，并在 `server/.env` 中配置微信小程序的AppSecret。

### Q: 如何备份数据？
A: 使用 `deploy.sh backup` 或 `deploy.bat backup` 命令，或手动导出MySQL数据库。

## 项目结构

```
qaminiprogram/
├── miniprogram/              # 微信小程序端
├── admin-web/               # 管理端Web应用
├── server/                  # 后端服务
├── nginx/                   # Nginx配置
├── docker-compose.yml       # Docker编排文件
├── deploy.sh               # Linux部署脚本
├── deploy.bat              # Windows部署脚本
├── DOCKER_DEPLOYMENT.md    # Docker部署文档
└── README.md               # 项目说明文档
```

详细的项目结构说明请参考 [Docker部署文档](./DOCKER_DEPLOYMENT.md)。
