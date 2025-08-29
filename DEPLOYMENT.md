# 刷刷题项目部署指南

## 环境要求

### 系统环境
- **操作系统**: Windows 11
- **WSL**: Windows Subsystem for Linux 2 (WSL2)
- **Linux 发行版**: Ubuntu 22.04 LTS
- **容器运行时**: Podman 4.0+
- **容器编排**: Podman Compose

### 硬件要求
- **内存**: 最少 4GB，推荐 8GB+
- **存储**: 最少 10GB 可用空间
- **CPU**: 2 核心以上

## 环境准备

### 1. 安装 WSL2 和 Ubuntu 22.04

```powershell
# 在 Windows PowerShell (管理员模式) 中执行
wsl --install -d Ubuntu-22.04

# 重启计算机后，设置 Ubuntu 用户名和密码
```

### 2. 配置 WSL2 环境

```bash
# 进入 WSL2 Ubuntu 环境
wsl -d Ubuntu-22.04

# 更新系统包
sudo apt update && sudo apt upgrade -y

# 安装必要的工具
sudo apt install -y curl wget git vim nano
```

### 3. 安装 Podman

```bash
# 添加 Podman 官方仓库
sudo apt update
sudo apt install -y software-properties-common
sudo add-apt-repository -y ppa:projectatomic/ppa

# 安装 Podman
sudo apt update
sudo apt install -y podman

# 验证安装
podman --version
```

### 4. 安装 Podman Compose

```bash
# 安装 Python3 和 pip
sudo apt install -y python3 python3-pip

# 安装 podman-compose
pip3 install podman-compose

# 添加到 PATH (添加到 ~/.bashrc)
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc

# 验证安装
podman-compose --version
```

### 5. 配置 Podman 无根模式

```bash
# 启用无根模式
echo 'export DOCKER_HOST="unix://$XDG_RUNTIME_DIR/podman/podman.sock"' >> ~/.bashrc
source ~/.bashrc

# 启动 Podman 服务
systemctl --user enable --now podman.socket
```

## 项目部署

### 1. 获取项目代码

```bash
# 克隆项目到 WSL2 环境
cd ~
git clone <项目仓库地址> shuashuati
cd shuashuati
```

### 2. 创建数据目录

```bash
# 创建 MySQL 数据持久化目录
mkdir -p ./data/mysql
sudo chown -R 999:999 ./data/mysql
```

### 3. 配置环境变量

```bash
# 复制环境变量模板
cp .env.example .env

# 编辑环境变量
nano .env
```

环境变量配置示例：
```bash
# 数据库配置
DB_HOST=mysql
DB_PORT=3306
DB_USER=shuashuati
DB_PASSWORD=shuashuati123
DB_NAME=shuashuati
DB_CHARSET=utf8mb4

# JWT 配置
JWT_SECRET=shuashuati_jwt_secret_key_2024
JWT_EXPIRE_HOURS=24

# 服务器配置
SERVER_PORT=8080
SERVER_MODE=release
GIN_MODE=release

# CORS 配置
ALLOWED_ORIGINS=http://localhost:5173,http://localhost:5174,http://localhost:3000,http://localhost:3001
```

### 4. 构建和启动服务

```bash
# 构建所有服务镜像
podman-compose build

# 启动所有服务
podman-compose up -d

# 查看服务状态
podman-compose ps

# 查看服务日志
podman-compose logs -f
```

### 5. 验证部署

```bash
# 检查容器状态
podman ps

# 检查网络连接
curl http://localhost:8080/api/health

# 检查前端服务
curl http://localhost:3000
curl http://localhost:3001
```

## 服务访问

部署成功后，可以通过以下地址访问各个服务：

- **后端 API**: http://localhost:8080
- **用户端前端**: http://localhost:3000
- **管理端前端**: http://localhost:3001
- **MySQL 数据库**: localhost:3306

## 常用管理命令

### 服务管理

```bash
# 启动所有服务
podman-compose up -d

# 停止所有服务
podman-compose down

# 重启特定服务
podman-compose restart backend

# 查看服务状态
podman-compose ps

# 查看服务日志
podman-compose logs backend
podman-compose logs mysql
```

### 数据库管理

```bash
# 连接到 MySQL 数据库
podman exec -it shuashuati_mysql mysql -u shuashuati -p

# 备份数据库
podman exec shuashuati_mysql mysqldump -u shuashuati -pshuashuati123 shuashuati > backup.sql

# 恢复数据库
podman exec -i shuashuati_mysql mysql -u shuashuati -pshuashuati123 shuashuati < backup.sql
```

### 容器管理

```bash
# 查看所有容器
podman ps -a

# 进入容器
podman exec -it shuashuati_backend /bin/sh

# 查看容器资源使用情况
podman stats

# 清理未使用的镜像和容器
podman system prune -a
```

## 故障排除

### 常见问题

1. **容器启动失败**
   ```bash
   # 查看详细错误日志
   podman-compose logs <service_name>
   
   # 检查端口占用
   netstat -tulpn | grep <port>
   ```

2. **数据库连接失败**
   ```bash
   # 检查 MySQL 容器状态
   podman logs shuashuati_mysql
   
   # 验证数据库连接
   podman exec shuashuati_mysql mysql -u root -p123456 -e "SHOW DATABASES;"
   ```

3. **前端访问失败**
   ```bash
   # 检查前端容器日志
   podman logs shuashuati_miniprogram
   podman logs shuashuati_admin
   
   # 检查 Nginx 配置
   podman exec shuashuati_miniprogram cat /etc/nginx/nginx.conf
   ```

4. **后端 API 错误**
   ```bash
   # 查看后端日志
   podman logs shuashuati_backend
   
   # 检查环境变量
   podman exec shuashuati_backend env | grep DB_
   ```

### 性能优化

1. **调整 MySQL 配置**
   ```bash
   # 编辑 docker-compose.yml 中的 MySQL 命令参数
   # 根据服务器内存调整 innodb-buffer-pool-size
   ```

2. **调整容器资源限制**
   ```yaml
   # 在 docker-compose.yml 中添加资源限制
   deploy:
     resources:
       limits:
         memory: 512M
         cpus: '0.5'
   ```

## 安全建议

1. **修改默认密码**
   - 修改 MySQL root 密码
   - 修改应用数据库用户密码
   - 修改 JWT 密钥

2. **网络安全**
   - 配置防火墙规则
   - 使用 HTTPS (生产环境)
   - 限制数据库访问

3. **数据备份**
   - 定期备份数据库
   - 备份应用配置文件
   - 测试恢复流程

## 更新和维护

### 应用更新

```bash
# 拉取最新代码
git pull origin main

# 重新构建镜像
podman-compose build

# 重启服务
podman-compose down
podman-compose up -d
```

### 系统维护

```bash
# 清理系统
podman system prune -a

# 更新系统包
sudo apt update && sudo apt upgrade -y

# 检查磁盘空间
df -h
```

## 技术支持

如遇到问题，请按以下步骤排查：

1. 查看相关服务日志
2. 检查网络连接
3. 验证配置文件
4. 查看系统资源使用情况
5. 参考故障排除章节

---

**注意**: 本部署指南适用于开发和测试环境。生产环境部署需要额外的安全配置和性能优化。