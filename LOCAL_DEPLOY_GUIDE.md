# 刷刷题项目本地部署指南

本指南将帮助您在 Ubuntu WSL 环境中部署刷刷题项目的完整开发环境。

## 📋 目录

- [系统要求](#系统要求)
- [快速开始](#快速开始)
- [手动部署](#手动部署)
- [服务访问](#服务访问)
- [开发指南](#开发指南)
- [故障排除](#故障排除)
- [常用命令](#常用命令)

## 🔧 系统要求

### 必需软件

- **Ubuntu WSL** (推荐 Ubuntu 20.04 或更高版本)
- **Docker** (20.10 或更高版本)
- **Docker Compose** (2.0 或更高版本)
- **Node.js** (18.x 或更高版本)
- **Go** (1.21 或更高版本)
- **Git**

### 硬件要求

- **内存**: 至少 4GB RAM (推荐 8GB)
- **存储**: 至少 10GB 可用空间
- **CPU**: 双核或更高

## 🚀 快速开始

### 方法一：自动化脚本（推荐）

1. **克隆项目**
   ```bash
   git clone <your-repo-url>
   cd qa
   ```

2. **运行自动化设置脚本**
   ```bash
   chmod +x setup.sh
   ./setup.sh
   ```

   脚本将自动：
   - 检查并安装所需软件
   - 配置环境变量
   - 设置本地域名
   - 构建和启动所有服务

3. **等待部署完成**
   
   脚本运行完成后，您将看到访问地址和常用命令。

### 方法二：Docker Compose 快速启动

如果您已经安装了所需软件：

```bash
# 启动所有服务
docker-compose -f docker-compose.local.yml up -d

# 查看服务状态
docker-compose -f docker-compose.local.yml ps

# 查看日志
docker-compose -f docker-compose.local.yml logs -f
```

## 🛠 手动部署

### 1. 环境准备

#### 安装 Docker

```bash
# 更新包索引
sudo apt-get update

# 安装必要的包
sudo apt-get install -y apt-transport-https ca-certificates curl gnupg lsb-release

# 添加 Docker 官方 GPG 密钥
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg

# 设置稳定版仓库
echo "deb [arch=amd64 signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

# 安装 Docker Engine
sudo apt-get update
sudo apt-get install -y docker-ce docker-ce-cli containerd.io

# 将当前用户添加到 docker 组
sudo usermod -aG docker $USER

# 重新登录或运行
newgrp docker
```

#### 安装 Docker Compose

```bash
# 下载 Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/download/v2.20.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose

# 添加执行权限
sudo chmod +x /usr/local/bin/docker-compose

# 验证安装
docker-compose --version
```

#### 安装 Node.js

```bash
# 安装 NodeSource 仓库
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -

# 安装 Node.js
sudo apt-get install -y nodejs

# 验证安装
node --version
npm --version
```

#### 安装 Go

```bash
# 下载 Go
cd /tmp
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz

# 解压到 /usr/local
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz

# 添加到 PATH
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
echo 'export GOPATH=$HOME/go' >> ~/.bashrc
echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.bashrc

# 重新加载配置
source ~/.bashrc

# 验证安装
go version
```

### 2. 项目配置

#### 配置环境变量

项目已包含配置好的 `.env` 文件，包含 Supabase 连接信息。如需修改：

```bash
# 编辑环境变量
nano .env
```

主要配置项：
- `SUPABASE_URL`: Supabase 项目 URL
- `SUPABASE_ANON_KEY`: Supabase 匿名密钥
- `SUPABASE_SERVICE_ROLE_KEY`: Supabase 服务角色密钥
- `JWT_SECRET`: JWT 签名密钥

#### 配置本地域名（可选）

```bash
# 编辑 hosts 文件
sudo nano /etc/hosts

# 添加以下内容
127.0.0.1 admin.qa.local
127.0.0.1 app.qa.local
127.0.0.1 api.qa.local
```

### 3. 构建和启动服务

#### 构建 Docker 镜像

```bash
# 构建所有服务的镜像
docker-compose -f docker-compose.local.yml build
```

#### 启动服务

```bash
# 启动所有服务
docker-compose -f docker-compose.local.yml up -d

# 查看服务状态
docker-compose -f docker-compose.local.yml ps
```

#### 查看日志

```bash
# 查看所有服务日志
docker-compose -f docker-compose.local.yml logs -f

# 查看特定服务日志
docker-compose -f docker-compose.local.yml logs -f server
docker-compose -f docker-compose.local.yml logs -f admin-web
docker-compose -f docker-compose.local.yml logs -f miniprogram
```

## 🌐 服务访问

部署完成后，您可以通过以下地址访问各个服务：

### 主要访问地址

| 服务 | 地址 | 说明 |
|------|------|------|
| **管理端** | http://localhost/admin/ | 题库管理系统 |
| **小程序端** | http://localhost/app/ | 用户答题界面 |
| **API 接口** | http://localhost/api/ | 后端 API 服务 |

### 直接访问地址

| 服务 | 地址 | 说明 |
|------|------|------|
| **管理端** | http://localhost:3000 | 直接访问管理端 |
| **小程序端** | http://localhost:3001 | 直接访问小程序端 |
| **后端 API** | http://localhost:8080 | 直接访问 API |

### 本地域名访问（如果配置了 hosts）

| 服务 | 地址 | 说明 |
|------|------|------|
| **管理端** | http://admin.qa.local | 管理端域名 |
| **小程序端** | http://app.qa.local | 小程序端域名 |
| **API 接口** | http://api.qa.local | API 域名 |

## 💻 开发指南

### 开发模式启动

如果您需要进行开发，可以使用以下方式：

#### 1. 只启动后端服务

```bash
# 只启动后端和数据库
docker-compose -f docker-compose.local.yml up -d server
```

#### 2. 本地开发前端

```bash
# 管理端开发
cd admin-web
npm install
npm run dev

# 小程序端开发
cd miniprogram
npm install
npm run dev
```

#### 3. 本地开发后端

```bash
# 进入后端目录
cd server

# 安装依赖
go mod download

# 运行开发服务器
go run main.go
```

### 代码热重载

开发模式下，代码修改会自动重新加载：

- **前端**: Vite 提供热重载
- **后端**: 使用 `air` 工具实现热重载

```bash
# 安装 air（Go 热重载工具）
go install github.com/cosmtrek/air@latest

# 在 server 目录下运行
air
```

### 数据库管理

项目使用 Supabase 作为数据库，您可以：

1. **访问 Supabase 控制台**: https://app.supabase.com
2. **查看数据库表**: 在 Supabase 项目中的 Table Editor
3. **执行 SQL**: 在 SQL Editor 中运行查询
4. **查看 API 文档**: 在 API 部分查看自动生成的文档

## 🔧 故障排除

### Docker相关问题

#### 1. Docker连接失败
**错误信息**: `FileNotFoundError: [Errno 2] No such file or directory`

**解决方案**:
```bash
# 方法1: 使用Docker修复脚本（推荐）
chmod +x docker-fix.sh
./docker-fix.sh

# 方法2: 手动启动Docker daemon
sudo dockerd > /dev/null 2>&1 &

# 方法3: 使用原生部署（不使用Docker）
chmod +x native-deploy.sh
./native-deploy.sh
```

#### 2. Docker服务未运行
```bash
# 检查Docker状态
docker version

# 启动Docker服务
sudo systemctl start docker
sudo systemctl enable docker

# WSL环境手动启动
./start-docker.sh
```

#### 3. Docker权限问题
```bash
# 将用户添加到docker组
sudo usermod -aG docker $USER
newgrp docker

# 重新登录后测试
docker run hello-world
```

#### 4. WSL Docker配置问题
```bash
# 创建Docker配置目录
sudo mkdir -p /etc/docker

# 配置Docker daemon for WSL
sudo tee /etc/docker/daemon.json > /dev/null <<EOF
{
    "hosts": ["unix:///var/run/docker.sock"],
    "iptables": false,
    "bridge": "none"
}
EOF

# 重启Docker
sudo systemctl restart docker
```

### 常见问题

#### 1. 端口被占用

```bash
# 查看端口占用
sudo netstat -tulpn | grep :8080
sudo netstat -tulpn | grep :3000
sudo netstat -tulpn | grep :3001

# 杀死占用端口的进程
sudo kill -9 <PID>
```

#### 2. 构建失败
```bash
# 清理Docker缓存
docker system prune -a

# 重新构建
docker-compose -f docker-compose.local.yml build --no-cache
```

#### 3. 网络问题
```bash
# 重置Docker网络
docker network prune

# 重新创建网络
docker-compose -f docker-compose.local.yml down
docker-compose -f docker-compose.local.yml up -d
```

#### 4. 服务启动失败

```bash
# 查看详细错误日志
docker-compose -f docker-compose.local.yml logs <service-name>

# 重新构建镜像
docker-compose -f docker-compose.local.yml build --no-cache <service-name>

# 重启服务
docker-compose -f docker-compose.local.yml restart <service-name>
```

#### 4. 数据库连接问题

检查 `.env` 文件中的 Supabase 配置：

```bash
# 验证 Supabase 连接
curl -H "apikey: YOUR_ANON_KEY" https://YOUR_PROJECT_ID.supabase.co/rest/v1/
```

#### 5. 前端构建失败

```bash
# 清理 node_modules 和重新安装
cd admin-web  # 或 miniprogram
rm -rf node_modules package-lock.json
npm install

# 重新构建
npm run build
```

### 日志查看

```bash
# 查看所有服务日志
docker-compose -f docker-compose.local.yml logs -f

# 查看特定服务日志
docker-compose -f docker-compose.local.yml logs -f server
docker-compose -f docker-compose.local.yml logs -f admin-web
docker-compose -f docker-compose.local.yml logs -f miniprogram
docker-compose -f docker-compose.local.yml logs -f nginx

# 查看最近的日志
docker-compose -f docker-compose.local.yml logs --tail=100 server
```

### 重置环境

如果遇到严重问题，可以重置整个环境：

```bash
# 停止所有服务
docker-compose -f docker-compose.local.yml down

# 删除所有容器和镜像
docker-compose -f docker-compose.local.yml down --rmi all --volumes --remove-orphans

# 清理 Docker 系统
docker system prune -a

# 重新构建和启动
docker-compose -f docker-compose.local.yml up --build -d
```

## 📝 常用命令

### Docker Compose 命令

```bash
# 启动所有服务
docker-compose -f docker-compose.local.yml up -d

# 停止所有服务
docker-compose -f docker-compose.local.yml down

# 重启服务
docker-compose -f docker-compose.local.yml restart

# 重新构建并启动
docker-compose -f docker-compose.local.yml up --build -d

# 查看服务状态
docker-compose -f docker-compose.local.yml ps

# 查看日志
docker-compose -f docker-compose.local.yml logs -f

# 进入容器
docker-compose -f docker-compose.local.yml exec server bash
docker-compose -f docker-compose.local.yml exec admin-web sh
```

### 开发命令

```bash
# 前端开发
cd admin-web && npm run dev
cd miniprogram && npm run dev

# 后端开发
cd server && go run main.go
cd server && air  # 热重载

# 构建前端
cd admin-web && npm run build
cd miniprogram && npm run build

# 构建后端
cd server && go build -o main .
```

### 数据库命令

```bash
# 应用数据库迁移
cd server && go run migrate.go

# 查看 Supabase 表
curl -H "apikey: YOUR_ANON_KEY" https://YOUR_PROJECT_ID.supabase.co/rest/v1/users
```

## 🔒 安全注意事项

1. **环境变量**: 不要将 `.env` 文件提交到版本控制系统
2. **Supabase 密钥**: 服务角色密钥只能在后端使用
3. **JWT 密钥**: 生产环境中使用强密码
4. **CORS 配置**: 根据实际需要配置允许的域名
5. **防火墙**: 确保只开放必要的端口

## 📚 相关文档

- [Docker 官方文档](https://docs.docker.com/)
- [Docker Compose 文档](https://docs.docker.com/compose/)
- [Supabase 文档](https://supabase.com/docs)
- [Vue.js 文档](https://vuejs.org/)
- [Go 官方文档](https://golang.org/doc/)
- [Gin 框架文档](https://gin-gonic.com/docs/)

## 🆘 获取帮助

如果您遇到问题：

1. 查看本文档的故障排除部分
2. 检查服务日志获取错误信息
3. 搜索相关错误信息
4. 联系开发团队获取支持

---

**祝您使用愉快！** 🎉