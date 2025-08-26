# 刷刷题项目 - Windows 开发环境

> 🚀 Windows 11 轻量级开发环境，支持 WSL2 + 容器技术，实现 Windows 开发 Linux 运行

## 📋 项目概述

刷刷题是一个完整的在线题库系统，包含三个核心模块：

- 🎨 **管理端**: Vue3 + Element Plus + TypeScript
- 📱 **小程序端**: Vue3 + Uni-app + TypeScript  
- 🔧 **后端服务**: Golang + Gin + Supabase

## 🎯 快速开始

### 一键启动（推荐）

```powershell
# 智能启动 - 自动检测环境并选择最佳方式
.\quick-start.ps1
```

### 首次使用

```powershell
# 1. 环境配置（需要管理员权限）
.\setup-windows.ps1

# 2. 重启计算机（如提示需要）

# 3. 启动项目
.\quick-start.ps1
```

## 🎮 启动模式

| 模式 | 命令 | 特点 | 适用场景 |
|------|------|------|----------|
| 🤖 **智能模式** | `quick-start.ps1` | 自动选择最佳方式 | 日常开发（推荐） |
| 🛠️ **开发模式** | `dev-windows.ps1` | 热重载、快速启动 | 开发调试 |
| 🐳 **容器模式** | `start-windows.ps1` | 环境隔离、生产一致 | 生产测试 |

## 🌐 访问地址

### 开发模式
- 管理端: http://localhost:3000
- 小程序端: http://localhost:5173
- 后端 API: http://localhost:8080

### 容器模式
- 管理端: http://localhost/admin/
- 小程序端: http://localhost/app/
- 后端 API: http://localhost/api/

## 🔧 常用命令

### 服务管理
```powershell
# 查看服务状态
.\dev-windows.ps1 -Status

# 停止所有服务
.\dev-windows.ps1 -Stop

# 重启服务
.\dev-windows.ps1 -Stop && .\dev-windows.ps1
```

### 分别启动服务
```powershell
# 仅启动后端
.\dev-windows.ps1 -Backend

# 仅启动管理端
.\dev-windows.ps1 -Admin

# 仅启动小程序端
.\dev-windows.ps1 -Mini

# 启动前端（管理端 + 小程序端）
.\dev-windows.ps1 -Admin -Mini
```

### 容器管理
```powershell
# 查看容器状态
wsl -d Ubuntu bash -c "cd /mnt/d/GITVIEW/qa && docker-compose -f docker-compose.local.yml ps"

# 查看容器日志
wsl -d Ubuntu bash -c "cd /mnt/d/GITVIEW/qa && docker-compose -f docker-compose.local.yml logs -f"

# 停止容器
wsl -d Ubuntu bash -c "cd /mnt/d/GITVIEW/qa && docker-compose -f docker-compose.local.yml down"
```

## 📦 环境要求

### 系统要求
- Windows 11 或 Windows 10 版本 2004+
- 至少 8GB 内存（推荐 16GB）
- 至少 20GB 可用磁盘空间

### 自动安装的组件
- ✅ WSL2 (Windows Subsystem for Linux)
- ✅ Ubuntu 发行版
- ✅ Node.js 18+
- ✅ Go 1.21+
- ✅ Docker 或 Podman
- ✅ 开发工具链

## 🛠️ 开发特性

### 开发模式优势
- ⚡ **快速启动**: 原生运行，启动速度快
- 🔥 **热重载**: 代码修改实时生效
- 🐛 **易调试**: 直接访问源码和日志
- 💾 **低资源**: 相比容器模式资源占用更少

### 容器模式优势
- 🔒 **环境隔离**: 完全隔离的运行环境
- 🎯 **生产一致**: 与生产环境完全一致
- 📦 **一键部署**: 容器化部署，简单可靠
- 🔄 **版本管理**: 支持多版本并行运行

## 📁 项目结构

```
qa/
├── 📁 admin-web/          # 管理端（Vue3 + Element Plus）
├── 📁 miniprogram/        # 小程序端（Vue3 + Uni-app）
├── 📁 server/             # 后端服务（Go + Gin）
├── 📁 nginx/              # Nginx 配置
├── 📁 .vscode/            # VS Code 配置
├── 🐳 docker-compose.*.yml # Docker 配置
├── 🚀 quick-start.ps1     # 智能启动脚本
├── 🛠️ dev-windows.ps1      # 开发模式脚本
├── 🐳 start-windows.ps1   # 容器模式脚本
├── ⚙️ setup-windows.ps1    # 环境配置脚本
└── 📖 WINDOWS_DEV_GUIDE.md # 详细开发指南
```

## 🔍 故障排除

### 常见问题

#### 1. WSL 相关问题
```powershell
# 重启 WSL
wsl --shutdown
wsl -d Ubuntu

# 检查 WSL 状态
wsl --status
wsl --list --verbose
```

#### 2. 容器相关问题
```bash
# 在 WSL 中检查容器状态
podman info
docker info

# 重启容器服务
sudo systemctl restart docker
```

#### 3. 端口冲突
```powershell
# 查看端口占用
netstat -ano | findstr :8080
netstat -ano | findstr :3000
netstat -ano | findstr :5173

# 停止占用端口的进程
taskkill /PID <PID> /F
```

#### 4. 权限问题
```bash
# 在 WSL 中修复文件权限
sudo chown -R $USER:$USER /mnt/d/GITVIEW/qa
chmod +x *.sh
```

### 获取帮助

```powershell
# 查看脚本帮助
.\quick-start.ps1 -Help
.\dev-windows.ps1 -Help
.\start-windows.ps1 -Help
```

## 📚 详细文档

- 📖 [Windows 开发环境详细指南](WINDOWS_DEV_GUIDE.md)
- 🐳 [本地部署指南](LOCAL_DEPLOY_GUIDE.md)
- 🚀 [Sealos 云端部署指南](SEALOS_DEPLOY_GUIDE.md)

## 🎨 开发工作流

### 日常开发
1. 启动开发环境: `quick-start.ps1`
2. 修改代码（支持热重载）
3. 测试功能
4. 提交代码

### 生产测试
1. 启动容器环境: `start-windows.ps1`
2. 验证功能
3. 性能测试
4. 部署上线

## 🤝 贡献指南

1. Fork 项目
2. 创建功能分支: `git checkout -b feature/new-feature`
3. 提交更改: `git commit -am 'Add new feature'`
4. 推送分支: `git push origin feature/new-feature`
5. 提交 Pull Request

## 📄 许可证

MIT License - 详见 [LICENSE](LICENSE) 文件

## 🙋‍♂️ 支持

如果遇到问题或需要帮助：

1. 📖 查看 [详细开发指南](WINDOWS_DEV_GUIDE.md)
2. 🔍 搜索已知问题
3. 💬 提交 Issue
4. 📧 联系开发团队

---

> 💡 **提示**: 推荐使用智能启动模式 `quick-start.ps1`，它会自动检测环境并选择最佳的启动方式！