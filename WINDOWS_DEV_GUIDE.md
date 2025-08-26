# Windows 11 轻量级开发环境配置指南

本指南专为刷刷题项目在 Windows 11 环境下的开发而设计，提供轻量级容器解决方案，避免 Docker Desktop 的资源消耗。

## 📋 目录

- [环境概述](#环境概述)
- [方案选择](#方案选择)
- [WSL2 + Podman 方案](#wsl2--podman-方案)
- [WSL2 + Docker 方案](#wsl2--docker-方案)
- [VS Code 开发环境](#vs-code-开发环境)
- [项目配置](#项目配置)
- [常用命令](#常用命令)
- [故障排除](#故障排除)

## 🎯 环境概述

### 目标
- **Windows 开发，Linux 运行**：确保开发环境与生产环境一致
- **轻量级**：避免 Docker Desktop 的资源消耗
- **高效**：快速启动和部署
- **兼容性**：支持完整的容器生态

### 技术栈
- **前端**：Vue3 + TypeScript + Vite
- **后端**：Golang + Gin
- **数据库**：Supabase PostgreSQL
- **容器**：Podman 或 Docker (无 Desktop)
- **开发环境**：WSL2 + VS Code

## 🚀 方案选择

### 方案对比

| 方案 | 资源占用 | 启动速度 | 兼容性 | 推荐度 |
|------|----------|----------|--------|---------|
| WSL2 + Podman | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ | 🥇 推荐 |
| WSL2 + Docker | ⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | 🥈 备选 |
| Docker Desktop | ⭐⭐ | ⭐⭐ | ⭐⭐⭐⭐⭐ | ❌ 不推荐 |

## 🐧 WSL2 + Podman 方案（推荐）

### 1. 安装 WSL2

```powershell
# 以管理员身份运行 PowerShell

# 启用 WSL 功能
dism.exe /online /enable-feature /featurename:Microsoft-Windows-Subsystem-Linux /all /norestart

# 启用虚拟机平台
dism.exe /online /enable-feature /featurename:VirtualMachinePlatform /all /norestart

# 重启计算机
Restart-Computer

# 设置 WSL2 为默认版本
wsl --set-default-version 2

# 安装 Ubuntu
wsl --install -d Ubuntu
```

### 2. 配置 Ubuntu 环境

```bash
# 更新系统
sudo apt update && sudo apt upgrade -y

# 安装基础工具
sudo apt install -y curl wget git build-essential

# 安装 Node.js
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt-get install -y nodejs

# 安装 Go
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

### 3. 安装 Podman

```bash
# 添加 Podman 仓库
sudo apt-get update
sudo apt-get -y install software-properties-common
sudo add-apt-repository -y ppa:projectatomic/ppa
sudo apt-get update

# 安装 Podman
sudo apt-get -y install podman

# 配置 Podman
sudo mkdir -p /etc/containers
sudo tee /etc/containers/registries.conf > /dev/null <<EOF
[registries.search]
registries = ['docker.io', 'registry.fedoraproject.org', 'quay.io', 'registry.access.redhat.com', 'registry.centos.org']
EOF

# 启用 rootless 模式
echo 'export DOCKER_HOST=unix:///run/user/1000/podman/podman.sock' >> ~/.bashrc
source ~/.bashrc
```

### 4. 安装 Podman Compose

```bash
# 安装 pip
sudo apt install -y python3-pip

# 安装 podman-compose
pip3 install podman-compose

# 添加到 PATH
echo 'export PATH=$PATH:~/.local/bin' >> ~/.bashrc
source ~/.bashrc

# 创建 docker-compose 别名
echo 'alias docker-compose="podman-compose"' >> ~/.bashrc
echo 'alias docker="podman"' >> ~/.bashrc
source ~/.bashrc
```

## 🐳 WSL2 + Docker 方案（备选）

### 1. 安装 Docker Engine

```bash
# 卸载旧版本
sudo apt-get remove docker docker-engine docker.io containerd runc

# 安装依赖
sudo apt-get update
sudo apt-get install -y \
    ca-certificates \
    curl \
    gnupg \
    lsb-release

# 添加 Docker GPG 密钥
sudo mkdir -p /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg

# 添加 Docker 仓库
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

# 安装 Docker Engine
sudo apt-get update
sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

# 将用户添加到 docker 组
sudo usermod -aG docker $USER
newgrp docker

# 启动 Docker 服务
sudo systemctl enable docker
sudo systemctl start docker
```

### 2. 配置 Docker for WSL

```bash
# 创建 Docker 配置
sudo mkdir -p /etc/docker
sudo tee /etc/docker/daemon.json > /dev/null <<EOF
{
  "hosts": ["unix:///var/run/docker.sock"],
  "iptables": false
}
EOF

# 重启 Docker
sudo systemctl restart docker
```

## 💻 VS Code 开发环境

### 1. 安装必要扩展

在 Windows 上安装以下 VS Code 扩展：

```json
{
  "recommendations": [
    "ms-vscode-remote.remote-wsl",
    "ms-vscode-remote.remote-containers",
    "ms-vscode.vscode-typescript-next",
    "Vue.volar",
    "golang.go",
    "ms-vscode.vscode-json",
    "bradlc.vscode-tailwindcss",
    "esbenp.prettier-vscode",
    "ms-vscode.vscode-eslint"
  ]
}
```

### 2. WSL 集成配置

在 WSL 中打开项目：

```bash
# 在 WSL 中导航到项目目录
cd /mnt/d/GITVIEW/qa

# 使用 VS Code 打开项目
code .
```

### 3. 容器开发配置

创建 `.vscode/settings.json`：

```json
{
  "remote.containers.defaultExtensions": [
    "Vue.volar",
    "golang.go",
    "ms-vscode.vscode-typescript-next"
  ],
  "go.toolsManagement.checkForUpdates": "local",
  "go.useLanguageServer": true,
  "typescript.preferences.includePackageJsonAutoImports": "on",
  "vue.server.hybridMode": true
}
```

## 🔧 项目配置

### 1. 创建 Windows 专用启动脚本

创建 `start-windows.ps1`：

```powershell
# Windows PowerShell 启动脚本
Write-Host "🚀 启动刷刷题项目 (Windows 环境)" -ForegroundColor Green

# 检查 WSL 状态
$wslStatus = wsl --status
if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ WSL 未正确安装或启动" -ForegroundColor Red
    exit 1
}

# 启动 WSL 并运行项目
Write-Host "📦 启动容器服务..." -ForegroundColor Blue
wsl -d Ubuntu bash -c "cd /mnt/d/GITVIEW/qa && ./start-containers.sh"

# 等待服务启动
Start-Sleep -Seconds 5

# 打开浏览器
Write-Host "🌐 打开浏览器..." -ForegroundColor Blue
Start-Process "http://localhost/admin/"
Start-Process "http://localhost/app/"

Write-Host "✅ 项目启动完成！" -ForegroundColor Green
Write-Host "管理端: http://localhost/admin/" -ForegroundColor Yellow
Write-Host "用户端: http://localhost/app/" -ForegroundColor Yellow
Write-Host "API: http://localhost/api/" -ForegroundColor Yellow
```

### 2. 创建容器启动脚本

创建 `start-containers.sh`：

```bash
#!/bin/bash

# 容器启动脚本
echo "🐳 启动刷刷题项目容器..."

# 检查容器运行时
if command -v podman &> /dev/null; then
    CONTAINER_CMD="podman-compose"
    echo "✅ 使用 Podman"
elif command -v docker &> /dev/null; then
    CONTAINER_CMD="docker-compose"
    echo "✅ 使用 Docker"
else
    echo "❌ 未找到容器运行时"
    exit 1
fi

# 启动服务
echo "📦 启动所有服务..."
$CONTAINER_CMD -f docker-compose.local.yml up -d

# 检查服务状态
echo "📊 检查服务状态..."
$CONTAINER_CMD -f docker-compose.local.yml ps

echo "✅ 容器启动完成！"
```

### 3. 创建开发模式脚本

创建 `dev-windows.ps1`：

```powershell
# Windows 开发模式启动脚本
Write-Host "🛠️ 启动开发模式" -ForegroundColor Green

# 启动后端
Write-Host "🔧 启动后端服务..." -ForegroundColor Blue
Start-Process powershell -ArgumentList "-NoExit", "-Command", "wsl -d Ubuntu bash -c 'cd /mnt/d/GITVIEW/qa/server && go run main.go'"

# 启动管理端
Write-Host "🎨 启动管理端..." -ForegroundColor Blue
Start-Process powershell -ArgumentList "-NoExit", "-Command", "wsl -d Ubuntu bash -c 'cd /mnt/d/GITVIEW/qa/admin-web && npm run dev'"

# 启动小程序端
Write-Host "📱 启动小程序端..." -ForegroundColor Blue
Start-Process powershell -ArgumentList "-NoExit", "-Command", "wsl -d Ubuntu bash -c 'cd /mnt/d/GITVIEW/qa/miniprogram && npm run dev'"

Write-Host "✅ 开发环境启动完成！" -ForegroundColor Green
```

## 📋 常用命令

### 容器管理

```bash
# 启动所有服务
podman-compose -f docker-compose.local.yml up -d
# 或
docker-compose -f docker-compose.local.yml up -d

# 查看服务状态
podman-compose -f docker-compose.local.yml ps

# 查看日志
podman-compose -f docker-compose.local.yml logs -f

# 停止服务
podman-compose -f docker-compose.local.yml down

# 重新构建
podman-compose -f docker-compose.local.yml build --no-cache
```

### 开发命令

```bash
# 进入 WSL
wsl -d Ubuntu

# 在 WSL 中打开 VS Code
code /mnt/d/GITVIEW/qa

# 查看容器状态
podman ps
# 或
docker ps

# 进入容器
podman exec -it <container_name> bash
# 或
docker exec -it <container_name> bash
```

### Windows 快捷命令

```powershell
# 启动项目（容器模式）
.\start-windows.ps1

# 启动项目（开发模式）
.\dev-windows.ps1

# 停止所有服务
wsl -d Ubuntu bash -c "cd /mnt/d/GITVIEW/qa && podman-compose -f docker-compose.local.yml down"
```

## 🔧 故障排除

### WSL 相关问题

#### 1. WSL 启动失败
```powershell
# 重启 WSL
wsl --shutdown
wsl -d Ubuntu

# 检查 WSL 版本
wsl --list --verbose

# 设置默认版本
wsl --set-version Ubuntu 2
```

#### 2. 文件权限问题
```bash
# 在 WSL 中修复权限
sudo chown -R $USER:$USER /mnt/d/GITVIEW/qa
chmod +x *.sh
```

### 容器相关问题

#### 1. Podman 权限问题
```bash
# 重新配置 rootless 模式
podman system reset
podman system migrate

# 检查 Podman 状态
podman info
```

#### 2. 端口冲突
```bash
# 查看端口占用
sudo netstat -tlnp | grep :80
sudo netstat -tlnp | grep :3000

# 停止占用端口的进程
sudo kill -9 <PID>
```

#### 3. 镜像拉取失败
```bash
# 配置镜像源
sudo tee /etc/containers/registries.conf > /dev/null <<EOF
[[registry]]
location = "docker.io"
[[registry.mirror]]
location = "https://registry.docker-cn.com"
EOF

# 重启 Podman
podman system restart
```

### 网络问题

#### 1. WSL 网络连接问题
```powershell
# 重置 WSL 网络
wsl --shutdown
Restart-Service LxssManager
```

#### 2. 容器网络问题
```bash
# 重置容器网络
podman network prune
podman-compose -f docker-compose.local.yml down
podman-compose -f docker-compose.local.yml up -d
```

## 🎯 性能优化

### 1. WSL 性能优化

创建 `.wslconfig` 文件在 `C:\Users\<用户名>\.wslconfig`：

```ini
[wsl2]
# 限制内存使用
memory=4GB
# 限制 CPU 核心数
processors=4
# 启用交换文件
swap=2GB
# 设置交换文件路径
swapFile=C:\temp\wsl-swap.vhdx
# 禁用页面报告
pageReporting=false
# 启用嵌套虚拟化
nestedVirtualization=true
```

### 2. 容器性能优化

```bash
# 配置 Podman 性能
echo 'export PODMAN_USERNS=keep-id' >> ~/.bashrc
echo 'export BUILDAH_ISOLATION=chroot' >> ~/.bashrc
source ~/.bashrc
```

### 3. 开发环境优化

```bash
# 配置 Git 性能
git config --global core.preloadindex true
git config --global core.fscache true
git config --global gc.auto 256

# 配置 Node.js 性能
echo 'export NODE_OPTIONS="--max-old-space-size=4096"' >> ~/.bashrc
source ~/.bashrc
```

## 📊 资源监控

### 1. 系统资源监控

```powershell
# 查看 WSL 资源使用
wsl --status
Get-Process -Name "wsl*" | Select-Object Name, CPU, WorkingSet

# 查看内存使用
Get-WmiObject -Class Win32_OperatingSystem | Select-Object TotalVisibleMemorySize, FreePhysicalMemory
```

### 2. 容器资源监控

```bash
# 查看容器资源使用
podman stats
# 或
docker stats

# 查看系统资源
htop
# 或
top
```

## 🚀 快速开始

### 🎯 一键启动（推荐）

最简单的方式是使用智能启动脚本：

```powershell
# 智能启动 - 自动检测环境并选择最佳方式
.\quick-start.ps1

# 或者指定启动模式
.\quick-start.ps1 -Mode dev        # 开发模式
.\quick-start.ps1 -Mode container  # 容器模式
.\quick-start.ps1 -Setup           # 环境配置
```

### 📋 完整配置流程

#### 1. 首次环境配置

```powershell
# 以管理员身份运行 PowerShell
.\setup-windows.ps1
```

#### 2. 重启计算机（如需要）

#### 3. 完成 Ubuntu 初始化

```bash
# 在 WSL Ubuntu 中运行
bash /mnt/d/GITVIEW/qa/install-dev-env.sh
```

#### 4. 启动项目

```powershell
# 方式一：智能启动（推荐）
.\quick-start.ps1

# 方式二：容器模式
.\start-windows.ps1

# 方式三：开发模式
.\dev-windows.ps1
```

### 🎮 启动模式对比

| 模式 | 启动脚本 | 启动速度 | 资源占用 | 热重载 | 适用场景 |
|------|----------|----------|----------|--------|-----------|
| 🤖 智能模式 | `quick-start.ps1` | ⭐⭐⭐⭐⭐ | 自适应 | ✅ | 日常开发（推荐） |
| 🛠️ 开发模式 | `dev-windows.ps1` | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ | ✅ | 开发调试 |
| 🐳 容器模式 | `start-windows.ps1` | ⭐⭐⭐ | ⭐⭐ | ❌ | 生产测试 |

### 🔧 常用操作

```powershell
# 查看服务状态
.\dev-windows.ps1 -Status

# 停止所有服务
.\dev-windows.ps1 -Stop

# 重启服务
.\dev-windows.ps1 -Stop && .\dev-windows.ps1

# 仅启动后端
.\dev-windows.ps1 -Backend

# 仅启动前端
.\dev-windows.ps1 -Admin -Mini
```

## 📝 总结

这个配置方案提供了：

- ✅ **轻量级**：避免 Docker Desktop 的资源消耗
- ✅ **高性能**：WSL2 提供接近原生 Linux 性能
- ✅ **兼容性**：完全兼容 Linux 生产环境
- ✅ **易用性**：一键启动脚本和 VS Code 集成
- ✅ **灵活性**：支持容器和原生开发模式

选择 **WSL2 + Podman** 方案可以获得最佳的资源效率和开发体验！