# Windows PowerShell 启动脚本
# 刷刷题项目 - Windows 环境快速启动

param(
    [switch]$Dev,
    [switch]$Container,
    [switch]$Help
)

# 颜色定义
$Green = "Green"
$Blue = "Blue"
$Yellow = "Yellow"
$Red = "Red"
$Cyan = "Cyan"

# 显示帮助信息
function Show-Help {
    Write-Host "🚀 刷刷题项目 Windows 启动脚本" -ForegroundColor $Green
    Write-Host ""
    Write-Host "用法:" -ForegroundColor $Blue
    Write-Host "  .\start-windows.ps1                # 默认容器模式启动"
    Write-Host "  .\start-windows.ps1 -Container     # 容器模式启动"
    Write-Host "  .\start-windows.ps1 -Dev          # 开发模式启动"
    Write-Host "  .\start-windows.ps1 -Help         # 显示帮助信息"
    Write-Host ""
    Write-Host "模式说明:" -ForegroundColor $Blue
    Write-Host "  容器模式: 使用 Docker/Podman 容器运行所有服务"
    Write-Host "  开发模式: 原生运行，支持热重载和调试"
    Write-Host ""
}

# 检查 WSL 状态
function Test-WSL {
    Write-Host "🔍 检查 WSL 状态..." -ForegroundColor $Blue
    
    try {
        $wslStatus = wsl --status 2>$null
        if ($LASTEXITCODE -ne 0) {
            Write-Host "❌ WSL 未正确安装或启动" -ForegroundColor $Red
            Write-Host "请运行以下命令安装 WSL:" -ForegroundColor $Yellow
            Write-Host "  wsl --install -d Ubuntu" -ForegroundColor $Cyan
            return $false
        }
        
        # 检查 Ubuntu 是否安装
        $distributions = wsl --list --quiet
        if ($distributions -notcontains "Ubuntu") {
            Write-Host "❌ Ubuntu 发行版未安装" -ForegroundColor $Red
            Write-Host "请运行以下命令安装 Ubuntu:" -ForegroundColor $Yellow
            Write-Host "  wsl --install -d Ubuntu" -ForegroundColor $Cyan
            return $false
        }
        
        Write-Host "✅ WSL 状态正常" -ForegroundColor $Green
        return $true
    }
    catch {
        Write-Host "❌ WSL 检查失败: $($_.Exception.Message)" -ForegroundColor $Red
        return $false
    }
}

# 检查项目文件
function Test-ProjectFiles {
    Write-Host "📁 检查项目文件..." -ForegroundColor $Blue
    
    $requiredFiles = @(
        "docker-compose.local.yml",
        "server\main.go",
        "admin-web\package.json",
        "miniprogram\package.json"
    )
    
    $missing = @()
    foreach ($file in $requiredFiles) {
        if (-not (Test-Path $file)) {
            $missing += $file
        }
    }
    
    if ($missing.Count -gt 0) {
        Write-Host "❌ 缺少必要文件:" -ForegroundColor $Red
        foreach ($file in $missing) {
            Write-Host "  - $file" -ForegroundColor $Red
        }
        return $false
    }
    
    Write-Host "✅ 项目文件检查通过" -ForegroundColor $Green
    return $true
}

# 容器模式启动
function Start-ContainerMode {
    Write-Host "🐳 启动容器模式..." -ForegroundColor $Green
    
    # 检查并创建启动脚本
    $startScript = "start-containers.sh"
    if (-not (Test-Path $startScript)) {
        Write-Host "📝 创建容器启动脚本..." -ForegroundColor $Blue
        Create-ContainerScript
    }
    
    # 启动容器
    Write-Host "📦 启动所有容器服务..." -ForegroundColor $Blue
    $result = wsl -d Ubuntu bash -c "cd /mnt/d/GITVIEW/qa && chmod +x start-containers.sh && ./start-containers.sh"
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ 容器启动成功!" -ForegroundColor $Green
        
        # 等待服务启动
        Write-Host "⏳ 等待服务启动..." -ForegroundColor $Blue
        Start-Sleep -Seconds 10
        
        # 检查服务状态
        Write-Host "📊 检查服务状态..." -ForegroundColor $Blue
        wsl -d Ubuntu bash -c "cd /mnt/d/GITVIEW/qa && docker-compose -f docker-compose.local.yml ps"
        
        # 打开浏览器
        Open-Browser
    } else {
        Write-Host "❌ 容器启动失败" -ForegroundColor $Red
        Write-Host "请检查 WSL 和容器运行时配置" -ForegroundColor $Yellow
    }
}

# 开发模式启动
function Start-DevMode {
    Write-Host "🛠️ 启动开发模式..." -ForegroundColor $Green
    
    # 检查环境变量
    if (-not (Test-Path ".env")) {
        Write-Host "⚠️ 未找到 .env 文件，使用示例配置" -ForegroundColor $Yellow
        if (Test-Path ".env.example") {
            Copy-Item ".env.example" ".env"
            Write-Host "📝 已复制 .env.example 到 .env" -ForegroundColor $Blue
        }
    }
    
    # 启动后端服务
    Write-Host "🔧 启动后端服务..." -ForegroundColor $Blue
    Start-Process powershell -ArgumentList @(
        "-NoExit",
        "-Command",
        "Write-Host '🔧 后端服务 - Go Server' -ForegroundColor Green; wsl -d Ubuntu bash -c 'cd /mnt/d/GITVIEW/qa/server && go run main.go'"
    )
    
    Start-Sleep -Seconds 2
    
    # 启动管理端
    Write-Host "🎨 启动管理端..." -ForegroundColor $Blue
    Start-Process powershell -ArgumentList @(
        "-NoExit",
        "-Command",
        "Write-Host '🎨 管理端 - Admin Web' -ForegroundColor Blue; wsl -d Ubuntu bash -c 'cd /mnt/d/GITVIEW/qa/admin-web && npm run dev'"
    )
    
    Start-Sleep -Seconds 2
    
    # 启动小程序端
    Write-Host "📱 启动小程序端..." -ForegroundColor $Blue
    Start-Process powershell -ArgumentList @(
        "-NoExit",
        "-Command",
        "Write-Host '📱 小程序端 - MiniProgram' -ForegroundColor Magenta; wsl -d Ubuntu bash -c 'cd /mnt/d/GITVIEW/qa/miniprogram && npm run dev'"
    )
    
    Write-Host "✅ 开发环境启动完成!" -ForegroundColor $Green
    Write-Host "📝 所有服务已在独立终端中启动" -ForegroundColor $Blue
    
    # 等待服务启动
    Write-Host "⏳ 等待服务启动..." -ForegroundColor $Blue
    Start-Sleep -Seconds 15
    
    # 打开浏览器
    Open-Browser
}

# 创建容器启动脚本
function Create-ContainerScript {
    $scriptContent = @'
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
    echo "请先安装 Docker 或 Podman"
    exit 1
fi

# 检查 docker-compose 文件
if [ ! -f "docker-compose.local.yml" ]; then
    echo "❌ 未找到 docker-compose.local.yml 文件"
    exit 1
fi

# 停止现有容器
echo "🛑 停止现有容器..."
$CONTAINER_CMD -f docker-compose.local.yml down 2>/dev/null || true

# 启动服务
echo "📦 启动所有服务..."
$CONTAINER_CMD -f docker-compose.local.yml up -d

if [ $? -eq 0 ]; then
    echo "✅ 容器启动成功！"
    
    # 等待服务就绪
    echo "⏳ 等待服务就绪..."
    sleep 5
    
    # 检查服务状态
    echo "📊 检查服务状态..."
    $CONTAINER_CMD -f docker-compose.local.yml ps
    
    echo ""
    echo "🌐 服务访问地址:"
    echo "  管理端: http://localhost/admin/"
    echo "  用户端: http://localhost/app/"
    echo "  API: http://localhost/api/"
else
    echo "❌ 容器启动失败"
    echo "请检查 Docker/Podman 配置和网络连接"
    exit 1
fi
'@
    
    $scriptContent | Out-File -FilePath "start-containers.sh" -Encoding UTF8
    Write-Host "📝 已创建 start-containers.sh" -ForegroundColor $Green
}

# 打开浏览器
function Open-Browser {
    Write-Host "🌐 打开浏览器..." -ForegroundColor $Blue
    
    try {
        Start-Process "http://localhost/admin/"
        Start-Sleep -Seconds 1
        Start-Process "http://localhost/app/"
        
        Write-Host "✅ 浏览器已打开" -ForegroundColor $Green
    }
    catch {
        Write-Host "⚠️ 无法自动打开浏览器" -ForegroundColor $Yellow
    }
    
    Write-Host ""
    Write-Host "🌐 服务访问地址:" -ForegroundColor $Green
    Write-Host "  管理端: http://localhost/admin/" -ForegroundColor $Cyan
    Write-Host "  用户端: http://localhost/app/" -ForegroundColor $Cyan
    Write-Host "  API: http://localhost/api/" -ForegroundColor $Cyan
    Write-Host ""
}

# 主函数
function Main {
    # 显示标题
    Write-Host "" 
    Write-Host "🚀 刷刷题项目 - Windows 启动器" -ForegroundColor $Green
    Write-Host "=================================" -ForegroundColor $Green
    Write-Host ""
    
    # 处理帮助参数
    if ($Help) {
        Show-Help
        return
    }
    
    # 检查 WSL
    if (-not (Test-WSL)) {
        return
    }
    
    # 检查项目文件
    if (-not (Test-ProjectFiles)) {
        return
    }
    
    # 根据参数选择启动模式
    if ($Dev) {
        Start-DevMode
    } elseif ($Container -or (-not $Dev -and -not $Container)) {
        # 默认使用容器模式
        Start-ContainerMode
    }
    
    Write-Host ""
    Write-Host "📋 常用命令:" -ForegroundColor $Blue
    Write-Host "  停止服务: wsl -d Ubuntu bash -c 'cd /mnt/d/GITVIEW/qa && docker-compose -f docker-compose.local.yml down'" -ForegroundColor $Cyan
    Write-Host "  查看日志: wsl -d Ubuntu bash -c 'cd /mnt/d/GITVIEW/qa && docker-compose -f docker-compose.local.yml logs -f'" -ForegroundColor $Cyan
    Write-Host "  重启服务: .\start-windows.ps1" -ForegroundColor $Cyan
    Write-Host ""
}

# 执行主函数
Main