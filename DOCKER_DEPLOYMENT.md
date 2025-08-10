# 刷刷题系统 Docker 部署指南

本文档详细说明如何使用 Docker 在自建服务器上部署刷刷题系统。

## 系统架构

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Nginx 反向代理  │    │   小程序前端      │    │   管理端前端      │
│   (端口: 80)     │    │   (端口: 3000)   │    │   (端口: 3001)   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
         ┌─────────────────────────────────────────────────┐
         │                Golang 后端服务                   │
         │                (端口: 8080)                     │
         └─────────────────────────────────────────────────┘
                                 │
         ┌─────────────────────────────────────────────────┐
         │                MySQL 数据库                     │
         │                (端口: 3306)                     │
         └─────────────────────────────────────────────────┘
```

## 前置要求

### 系统要求
- 操作系统：Linux (推荐 Ubuntu 20.04+) 或 Windows 10/11
- 内存：至少 4GB RAM
- 存储：至少 10GB 可用空间
- 网络：可访问互联网（用于下载镜像）

### 软件要求
- Docker Engine 20.10+
- Docker Compose 2.0+

## 安装 Docker

### Linux (Ubuntu/Debian)
```bash
# 更新包索引
sudo apt update

# 安装必要的包
sudo apt install apt-transport-https ca-certificates curl gnupg lsb-release

# 添加Docker官方GPG密钥
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg

# 添加Docker仓库
echo "deb [arch=amd64 signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

# 安装Docker
sudo apt update
sudo apt install docker-ce docker-ce-cli containerd.io docker-compose-plugin

# 启动Docker服务
sudo systemctl start docker
sudo systemctl enable docker

# 将当前用户添加到docker组
sudo usermod -aG docker $USER

# 重新登录或执行以下命令
newgrp docker
```

### Windows
1. 下载并安装 [Docker Desktop for Windows](https://www.docker.com/products/docker-desktop)
2. 启动 Docker Desktop
3. 确保 WSL 2 已启用（推荐）

### macOS
1. 下载并安装 [Docker Desktop for Mac](https://www.docker.com/products/docker-desktop)
2. 启动 Docker Desktop

## 部署步骤

### 1. 获取项目代码
```bash
# 克隆项目（如果使用Git）
git clone <your-repository-url>
cd qaminiprogram

# 或者直接上传项目文件到服务器
```

### 2. 配置环境变量（可选）
如需自定义配置，可以修改以下文件：
- `docker-compose.yml` - Docker服务配置
- `server/.env` - 后端环境变量
- `miniprogram/.env` - 小程序端环境变量

### 3. 启动服务

#### Linux/macOS
```bash
# 给脚本执行权限
chmod +x deploy.sh

# 启动所有服务
./deploy.sh start
```

#### Windows
```cmd
# 启动所有服务
deploy.bat start
```

#### 手动启动（所有平台）
```bash
# 构建并启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f
```

### 4. 验证部署

等待所有服务启动完成后（约1-2分钟），访问以下地址验证部署：

- **小程序端**: http://localhost 或 http://your-server-ip
- **管理端**: http://localhost/admin 或 http://your-server-ip/admin
- **API接口**: http://localhost/api 或 http://your-server-ip/api

### 5. 默认账号

系统会自动创建以下默认账号：

**管理员账号**
- 用户名：admin
- 密码：123456

**测试用户账号**
- 用户名：testuser
- 密码：123456

**数据库账号**
- 主机：localhost:3306
- 数据库：qaminiprogram
- 用户名：qaminiprogram
- 密码：qaminiprogram123

## 常用命令

### 服务管理
```bash
# 启动服务
docker-compose up -d

# 停止服务
docker-compose down

# 重启服务
docker-compose restart

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f [service_name]

# 重新构建镜像
docker-compose build --no-cache
```

### 数据库管理
```bash
# 进入MySQL容器
docker-compose exec mysql mysql -u qaminiprogram -pqaminiprogram123 qaminiprogram

# 备份数据库
docker-compose exec mysql mysqldump -u qaminiprogram -pqaminiprogram123 qaminiprogram > backup.sql

# 恢复数据库
docker-compose exec -T mysql mysql -u qaminiprogram -pqaminiprogram123 qaminiprogram < backup.sql
```

### 容器管理
```bash
# 查看所有容器
docker ps -a

# 进入容器
docker exec -it <container_name> /bin/sh

# 查看容器日志
docker logs <container_name>

# 清理未使用的资源
docker system prune -f
```

## 生产环境配置

### 1. 域名配置

修改 `nginx/conf.d/default.conf`，将 `localhost` 替换为你的域名：

```nginx
server {
    listen 80;
    server_name your-domain.com;  # 替换为你的域名
    # ... 其他配置
}

server {
    listen 80;
    server_name admin.your-domain.com;  # 管理端域名
    # ... 其他配置
}
```

### 2. HTTPS 配置

在 `nginx/conf.d/default.conf` 中添加 SSL 配置：

```nginx
server {
    listen 443 ssl;
    server_name your-domain.com;
    
    ssl_certificate /etc/nginx/ssl/cert.pem;
    ssl_certificate_key /etc/nginx/ssl/key.pem;
    
    # ... 其他配置
}

server {
    listen 80;
    server_name your-domain.com;
    return 301 https://$server_name$request_uri;
}
```

### 3. 环境变量配置

修改 `docker-compose.yml` 中的环境变量：

```yaml
services:
  backend:
    environment:
      JWT_SECRET: "your-production-jwt-secret"  # 生产环境JWT密钥
      SERVER_MODE: "release"  # 生产模式
      ALLOWED_ORIGINS: "https://your-domain.com"  # 允许的域名
```

### 4. 数据持久化

确保数据库数据持久化，`docker-compose.yml` 中已配置：

```yaml
volumes:
  mysql_data:  # 数据库数据卷
```

## 故障排除

### 常见问题

1. **端口冲突**
   - 检查端口是否被占用：`netstat -tlnp | grep :80`
   - 修改 `docker-compose.yml` 中的端口映射

2. **服务启动失败**
   - 查看日志：`docker-compose logs <service_name>`
   - 检查配置文件语法
   - 确保Docker有足够的资源

3. **数据库连接失败**
   - 检查数据库容器是否正常运行
   - 验证数据库连接参数
   - 查看后端服务日志

4. **前端无法访问API**
   - 检查nginx配置
   - 验证API路径配置
   - 查看网络连接

### 日志查看
```bash
# 查看所有服务日志
docker-compose logs

# 查看特定服务日志
docker-compose logs backend
docker-compose logs mysql
docker-compose logs nginx

# 实时查看日志
docker-compose logs -f
```

### 性能监控
```bash
# 查看容器资源使用情况
docker stats

# 查看系统资源
top
df -h
free -h
```

## 备份与恢复

### 数据备份
```bash
# 使用脚本备份
./deploy.sh backup

# 手动备份数据库
mkdir -p backups
docker-compose exec mysql mysqldump -u qaminiprogram -pqaminiprogram123 qaminiprogram > backups/backup_$(date +%Y%m%d_%H%M%S).sql

# 备份上传的文件（如果有）
tar -czf backups/uploads_$(date +%Y%m%d_%H%M%S).tar.gz uploads/
```

### 数据恢复
```bash
# 恢复数据库
docker-compose exec -T mysql mysql -u qaminiprogram -pqaminiprogram123 qaminiprogram < backups/backup_20240101_120000.sql

# 恢复上传的文件
tar -xzf backups/uploads_20240101_120000.tar.gz
```

## 更新部署

### 更新应用代码
```bash
# 停止服务
docker-compose down

# 更新代码
git pull  # 或上传新的代码文件

# 重新构建镜像
docker-compose build --no-cache

# 启动服务
docker-compose up -d
```

### 更新Docker镜像
```bash
# 拉取最新镜像
docker-compose pull

# 重启服务
docker-compose up -d
```

## 安全建议

1. **修改默认密码**
   - 登录后立即修改管理员密码
   - 修改数据库密码

2. **网络安全**
   - 使用防火墙限制访问
   - 配置HTTPS
   - 定期更新系统和Docker

3. **数据安全**
   - 定期备份数据
   - 限制数据库访问权限
   - 监控系统日志

## 支持

如果在部署过程中遇到问题，请：

1. 查看本文档的故障排除部分
2. 检查系统日志和容器日志
3. 确认系统要求是否满足
4. 联系技术支持

---

**注意**: 本部署方案适用于开发、测试和小规模生产环境。对于大规模生产环境，建议使用 Kubernetes 等容器编排平台。