# 刷刷题项目构建脚本

## Windows PowerShell 脚本

### 1. 全量构建脚本 (build-all.ps1)

```powershell
# 刷刷题项目全量构建脚本
# 用途：构建所有前端项目

Write-Host "=== 刷刷题项目全量构建开始 ===" -ForegroundColor Green

# 检查Node.js环境
if (!(Get-Command "node" -ErrorAction SilentlyContinue)) {
    Write-Host "错误: 未找到Node.js，请先安装Node.js" -ForegroundColor Red
    exit 1
}

# 检查Go环境
if (!(Get-Command "go" -ErrorAction SilentlyContinue)) {
    Write-Host "警告: 未找到Go环境，将跳过后端构建" -ForegroundColor Yellow
    $skipBackend = $true
} else {
    $skipBackend = $false
}

# 获取项目根目录
$projectRoot = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location $projectRoot

Write-Host "项目根目录: $projectRoot" -ForegroundColor Cyan

# 构建管理端
Write-Host "\n=== 构建管理端Web应用 ===" -ForegroundColor Yellow
if (Test-Path "admin-web") {
    Set-Location "admin-web"
    
    Write-Host "安装依赖..." -ForegroundColor Cyan
    npm install
    if ($LASTEXITCODE -ne 0) {
        Write-Host "管理端依赖安装失败" -ForegroundColor Red
        exit 1
    }
    
    Write-Host "构建项目..." -ForegroundColor Cyan
    npm run build
    if ($LASTEXITCODE -ne 0) {
        Write-Host "管理端构建失败" -ForegroundColor Red
        exit 1
    }
    
    Write-Host "管理端构建完成 ✓" -ForegroundColor Green
    Set-Location ..
} else {
    Write-Host "未找到admin-web目录" -ForegroundColor Red
}

# 构建小程序端
Write-Host "\n=== 构建小程序端 ===" -ForegroundColor Yellow
if (Test-Path "miniprogram") {
    Set-Location "miniprogram"
    
    Write-Host "安装依赖..." -ForegroundColor Cyan
    npm install
    if ($LASTEXITCODE -ne 0) {
        Write-Host "小程序端依赖安装失败" -ForegroundColor Red
        exit 1
    }
    
    Write-Host "构建微信小程序..." -ForegroundColor Cyan
    npm run build:mp-weixin
    if ($LASTEXITCODE -ne 0) {
        Write-Host "小程序端构建失败" -ForegroundColor Red
        exit 1
    }
    
    Write-Host "小程序端构建完成 ✓" -ForegroundColor Green
    Write-Host "请使用微信开发者工具打开: miniprogram/dist/build/mp-weixin" -ForegroundColor Cyan
    Set-Location ..
} else {
    Write-Host "未找到miniprogram目录" -ForegroundColor Red
}

# 构建后端
if (-not $skipBackend) {
    Write-Host "\n=== 构建后端服务 ===" -ForegroundColor Yellow
    if (Test-Path "server") {
        Set-Location "server"
        
        Write-Host "下载Go依赖..." -ForegroundColor Cyan
        go mod tidy
        if ($LASTEXITCODE -ne 0) {
            Write-Host "Go依赖下载失败" -ForegroundColor Red
            exit 1
        }
        
        Write-Host "编译后端服务..." -ForegroundColor Cyan
        go build -o qaminiprogram-server.exe main.go
        if ($LASTEXITCODE -ne 0) {
            Write-Host "后端服务编译失败" -ForegroundColor Red
            exit 1
        }
        
        Write-Host "后端服务构建完成 ✓" -ForegroundColor Green
        Write-Host "可执行文件: server/qaminiprogram-server.exe" -ForegroundColor Cyan
        Set-Location ..
    } else {
        Write-Host "未找到server目录" -ForegroundColor Red
    }
}

Write-Host "\n=== 构建完成 ===" -ForegroundColor Green
Write-Host "构建输出:" -ForegroundColor Cyan
Write-Host "  - 管理端: admin-web/dist/" -ForegroundColor White
Write-Host "  - 小程序: miniprogram/dist/build/mp-weixin/" -ForegroundColor White
if (-not $skipBackend) {
    Write-Host "  - 后端: server/qaminiprogram-server.exe" -ForegroundColor White
}

Write-Host "\n下一步操作:" -ForegroundColor Cyan
Write-Host "  1. 部署管理端到Web服务器" -ForegroundColor White
Write-Host "  2. 使用微信开发者工具发布小程序" -ForegroundColor White
Write-Host "  3. 部署后端服务到云服务器" -ForegroundColor White
```

### 2. 开发环境启动脚本 (start-dev.ps1)

```powershell
# 刷刷题项目开发环境启动脚本
# 用途：同时启动所有开发服务

Write-Host "=== 刷刷题开发环境启动 ===" -ForegroundColor Green

# 获取项目根目录
$projectRoot = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location $projectRoot

# 启动后端服务
Write-Host "启动后端服务..." -ForegroundColor Yellow
if (Test-Path "server") {
    Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd '$projectRoot\server'; Write-Host '后端服务启动中...' -ForegroundColor Green; go run main.go"
    Write-Host "后端服务已在新窗口启动 (端口: 8080)" -ForegroundColor Green
} else {
    Write-Host "未找到server目录" -ForegroundColor Red
}

# 等待2秒
Start-Sleep -Seconds 2

# 启动管理端
Write-Host "启动管理端开发服务器..." -ForegroundColor Yellow
if (Test-Path "admin-web") {
    Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd '$projectRoot\admin-web'; Write-Host '管理端开发服务器启动中...' -ForegroundColor Green; npm run dev"
    Write-Host "管理端已在新窗口启动 (端口: 5173)" -ForegroundColor Green
} else {
    Write-Host "未找到admin-web目录" -ForegroundColor Red
}

# 等待2秒
Start-Sleep -Seconds 2

# 启动小程序开发
Write-Host "启动小程序开发模式..." -ForegroundColor Yellow
if (Test-Path "miniprogram") {
    Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd '$projectRoot\miniprogram'; Write-Host '小程序开发模式启动中...' -ForegroundColor Green; npm run dev:mp-weixin"
    Write-Host "小程序开发模式已在新窗口启动" -ForegroundColor Green
    Write-Host "请使用微信开发者工具打开: miniprogram/dist/dev/mp-weixin" -ForegroundColor Cyan
} else {
    Write-Host "未找到miniprogram目录" -ForegroundColor Red
}

Write-Host "\n=== 开发环境启动完成 ===" -ForegroundColor Green
Write-Host "服务地址:" -ForegroundColor Cyan
Write-Host "  - 后端API: http://localhost:8080" -ForegroundColor White
Write-Host "  - 管理端: http://localhost:5173" -ForegroundColor White
Write-Host "  - 小程序: 使用微信开发者工具打开" -ForegroundColor White

Write-Host "\n按任意键退出..." -ForegroundColor Yellow
$null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
```

### 3. 清理脚本 (clean.ps1)

```powershell
# 刷刷题项目清理脚本
# 用途：清理构建文件和依赖

Write-Host "=== 刷刷题项目清理 ===" -ForegroundColor Green

# 获取项目根目录
$projectRoot = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location $projectRoot

Write-Host "项目根目录: $projectRoot" -ForegroundColor Cyan

# 清理管理端
Write-Host "\n清理管理端..." -ForegroundColor Yellow
if (Test-Path "admin-web") {
    Set-Location "admin-web"
    
    if (Test-Path "node_modules") {
        Write-Host "删除 node_modules..." -ForegroundColor Cyan
        Remove-Item -Recurse -Force "node_modules"
    }
    
    if (Test-Path "dist") {
        Write-Host "删除 dist..." -ForegroundColor Cyan
        Remove-Item -Recurse -Force "dist"
    }
    
    if (Test-Path "package-lock.json") {
        Write-Host "删除 package-lock.json..." -ForegroundColor Cyan
        Remove-Item -Force "package-lock.json"
    }
    
    Write-Host "管理端清理完成 ✓" -ForegroundColor Green
    Set-Location ..
}

# 清理小程序端
Write-Host "\n清理小程序端..." -ForegroundColor Yellow
if (Test-Path "miniprogram") {
    Set-Location "miniprogram"
    
    if (Test-Path "node_modules") {
        Write-Host "删除 node_modules..." -ForegroundColor Cyan
        Remove-Item -Recurse -Force "node_modules"
    }
    
    if (Test-Path "dist") {
        Write-Host "删除 dist..." -ForegroundColor Cyan
        Remove-Item -Recurse -Force "dist"
    }
    
    if (Test-Path "unpackage") {
        Write-Host "删除 unpackage..." -ForegroundColor Cyan
        Remove-Item -Recurse -Force "unpackage"
    }
    
    if (Test-Path "package-lock.json") {
        Write-Host "删除 package-lock.json..." -ForegroundColor Cyan
        Remove-Item -Force "package-lock.json"
    }
    
    Write-Host "小程序端清理完成 ✓" -ForegroundColor Green
    Set-Location ..
}

# 清理后端
Write-Host "\n清理后端..." -ForegroundColor Yellow
if (Test-Path "server") {
    Set-Location "server"
    
    if (Test-Path "qaminiprogram-server.exe") {
        Write-Host "删除可执行文件..." -ForegroundColor Cyan
        Remove-Item -Force "qaminiprogram-server.exe"
    }
    
    if (Test-Path "logs") {
        Write-Host "删除日志文件..." -ForegroundColor Cyan
        Remove-Item -Recurse -Force "logs"
    }
    
    Write-Host "后端清理完成 ✓" -ForegroundColor Green
    Set-Location ..
}

# 清理根目录
Write-Host "\n清理根目录..." -ForegroundColor Yellow
if (Test-Path "node_modules") {
    Write-Host "删除根目录 node_modules..." -ForegroundColor Cyan
    Remove-Item -Recurse -Force "node_modules"
}

if (Test-Path "package-lock.json") {
    Write-Host "删除根目录 package-lock.json..." -ForegroundColor Cyan
    Remove-Item -Force "package-lock.json"
}

Write-Host "\n=== 清理完成 ===" -ForegroundColor Green
Write-Host "已清理的内容:" -ForegroundColor Cyan
Write-Host "  - 所有 node_modules 目录" -ForegroundColor White
Write-Host "  - 所有构建输出目录" -ForegroundColor White
Write-Host "  - 所有 package-lock.json 文件" -ForegroundColor White
Write-Host "  - 后端可执行文件" -ForegroundColor White
```

### 4. 部署脚本 (deploy.ps1)

```powershell
# 刷刷题项目部署脚本
# 用途：自动化部署流程

param(
    [string]$Target = "all",  # all, admin, miniprogram, server
    [string]$Environment = "production"  # development, production
)

Write-Host "=== 刷刷题项目部署 ===" -ForegroundColor Green
Write-Host "部署目标: $Target" -ForegroundColor Cyan
Write-Host "部署环境: $Environment" -ForegroundColor Cyan

# 获取项目根目录
$projectRoot = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location $projectRoot

# 检查Vercel CLI
function Test-VercelCLI {
    if (!(Get-Command "vercel" -ErrorAction SilentlyContinue)) {
        Write-Host "未找到Vercel CLI，正在安装..." -ForegroundColor Yellow
        npm install -g vercel
        if ($LASTEXITCODE -ne 0) {
            Write-Host "Vercel CLI安装失败" -ForegroundColor Red
            return $false
        }
    }
    return $true
}

# 部署管理端
function Deploy-AdminWeb {
    Write-Host "\n=== 部署管理端 ===" -ForegroundColor Yellow
    
    if (!(Test-Path "admin-web")) {
        Write-Host "未找到admin-web目录" -ForegroundColor Red
        return
    }
    
    Set-Location "admin-web"
    
    # 安装依赖
    Write-Host "安装依赖..." -ForegroundColor Cyan
    npm install
    if ($LASTEXITCODE -ne 0) {
        Write-Host "依赖安装失败" -ForegroundColor Red
        return
    }
    
    # 构建项目
    Write-Host "构建项目..." -ForegroundColor Cyan
    if ($Environment -eq "production") {
        npm run build
    } else {
        npm run build:dev
    }
    
    if ($LASTEXITCODE -ne 0) {
        Write-Host "项目构建失败" -ForegroundColor Red
        return
    }
    
    # 部署到Vercel
    if (Test-VercelCLI) {
        Write-Host "部署到Vercel..." -ForegroundColor Cyan
        if ($Environment -eq "production") {
            vercel --prod
        } else {
            vercel
        }
        
        if ($LASTEXITCODE -eq 0) {
            Write-Host "管理端部署成功 ✓" -ForegroundColor Green
        } else {
            Write-Host "管理端部署失败" -ForegroundColor Red
        }
    }
    
    Set-Location ..
}

# 构建小程序
function Build-Miniprogram {
    Write-Host "\n=== 构建小程序 ===" -ForegroundColor Yellow
    
    if (!(Test-Path "miniprogram")) {
        Write-Host "未找到miniprogram目录" -ForegroundColor Red
        return
    }
    
    Set-Location "miniprogram"
    
    # 安装依赖
    Write-Host "安装依赖..." -ForegroundColor Cyan
    npm install
    if ($LASTEXITCODE -ne 0) {
        Write-Host "依赖安装失败" -ForegroundColor Red
        return
    }
    
    # 构建小程序
    Write-Host "构建小程序..." -ForegroundColor Cyan
    npm run build:mp-weixin
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "小程序构建成功 ✓" -ForegroundColor Green
        Write-Host "请使用微信开发者工具打开: miniprogram/dist/build/mp-weixin" -ForegroundColor Cyan
        Write-Host "然后在微信开发者工具中上传代码并发布" -ForegroundColor Cyan
    } else {
        Write-Host "小程序构建失败" -ForegroundColor Red
    }
    
    Set-Location ..
}

# 构建后端
function Build-Server {
    Write-Host "\n=== 构建后端 ===" -ForegroundColor Yellow
    
    if (!(Test-Path "server")) {
        Write-Host "未找到server目录" -ForegroundColor Red
        return
    }
    
    if (!(Get-Command "go" -ErrorAction SilentlyContinue)) {
        Write-Host "未找到Go环境" -ForegroundColor Red
        return
    }
    
    Set-Location "server"
    
    # 下载依赖
    Write-Host "下载Go依赖..." -ForegroundColor Cyan
    go mod tidy
    if ($LASTEXITCODE -ne 0) {
        Write-Host "Go依赖下载失败" -ForegroundColor Red
        return
    }
    
    # 编译
    Write-Host "编译后端服务..." -ForegroundColor Cyan
    go build -o qaminiprogram-server.exe main.go
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "后端服务构建成功 ✓" -ForegroundColor Green
        Write-Host "可执行文件: server/qaminiprogram-server.exe" -ForegroundColor Cyan
        Write-Host "请将可执行文件上传到云服务器并运行" -ForegroundColor Cyan
    } else {
        Write-Host "后端服务构建失败" -ForegroundColor Red
    }
    
    Set-Location ..
}

# 执行部署
switch ($Target.ToLower()) {
    "admin" { Deploy-AdminWeb }
    "miniprogram" { Build-Miniprogram }
    "server" { Build-Server }
    "all" {
        Deploy-AdminWeb
        Build-Miniprogram
        Build-Server
    }
    default {
        Write-Host "无效的部署目标: $Target" -ForegroundColor Red
        Write-Host "可用选项: all, admin, miniprogram, server" -ForegroundColor Yellow
    }
}

Write-Host "\n=== 部署完成 ===" -ForegroundColor Green
```

## 使用说明

### 1. 脚本权限设置

在首次使用前，需要设置PowerShell执行策略：

```powershell
# 以管理员身份运行PowerShell，然后执行：
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

### 2. 脚本使用方法

```powershell
# 全量构建
.\build-all.ps1

# 启动开发环境
.\start-dev.ps1

# 清理项目
.\clean.ps1

# 部署所有项目
.\deploy.ps1 -Target all -Environment production

# 只部署管理端
.\deploy.ps1 -Target admin -Environment production

# 只构建小程序
.\deploy.ps1 -Target miniprogram

# 只构建后端
.\deploy.ps1 -Target server
```

### 3. 快速开始

1. **首次设置**：
   ```powershell
   # 清理并重新安装
   .\clean.ps1
   .\build-all.ps1
   ```

2. **日常开发**：
   ```powershell
   # 启动所有开发服务
   .\start-dev.ps1
   ```

3. **生产部署**：
   ```powershell
   # 部署到生产环境
   .\deploy.ps1 -Target all -Environment production
   ```

### 4. 注意事项

- 确保已安装 Node.js 18+ 和 Go 1.21+
- 首次运行需要安装依赖，可能需要较长时间
- 部署前请确保已配置好相关环境变量
- 小程序发布需要微信开发者工具和相应权限
- 后端部署需要配置数据库连接信息

### 5. 故障排除

- **权限错误**：以管理员身份运行PowerShell
- **网络错误**：检查网络连接和防火墙设置
- **依赖安装失败**：清理后重新安装
- **构建失败**：检查代码语法和类型错误