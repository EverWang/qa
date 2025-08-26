# Windows 开发模式启动脚本
# 刷刷题项目 - 开发环境快速启动

param(
    [switch]$Backend,
    [switch]$Admin,
    [switch]$Mini,
    [switch]$All,
    [switch]$Stop,
    [switch]$Status,
    [switch]$Help
)

# 颜色定义
$Green = "Green"
$Blue = "Blue"
$Yellow = "Yellow"
$Red = "Red"
$Cyan = "Cyan"
$Magenta = "Magenta"

# 显示帮助信息
function Show-Help {
    Write-Host "🛠️ 刷刷题项目开发模式启动脚本" -ForegroundColor $Green
    Write-Host ""
    Write-Host "用法:" -ForegroundColor $Blue
    Write-Host "  .\dev-windows.ps1                  # 启动所有服务"
    Write-Host "  .\dev-windows.ps1 -All            # 启动所有服务"
    Write-Host "  .\dev-windows.ps1 -Backend        # 仅启动后端服务"
    Write-Host "  .\dev-windows.ps1 -Admin          # 仅启动管理端"
    Write-Host "  .\dev-windows.ps1 -Mini           # 仅启动小程序端"
    Write-Host "  .\dev-windows.ps1 -Stop           # 停止所有服务"
    Write-Host "  .\dev-windows.ps1 -Status         # 查看服务状态"
    Write-Host "  .\dev-windows.ps1 -Help           # 显示帮助信息"
    Write-Host ""
    Write-Host "服务说明:" -ForegroundColor $Blue
    Write-Host "  后端服务: Go + Gin (端口 8080)"
    Write-Host "  管理端: Vue3 + Element Plus (端口 3000)"
    Write-Host "  小程序端: Vue3 + Uni-app (端口 5173)"
    Write-Host ""
    Write-Host "开发特性:" -ForegroundColor $Blue
    Write-Host "  ✅ 热重载支持"
    Write-Host "  ✅ 实时调试"
    Write-Host "  ✅ 源码映射"
    Write-Host "  ✅ 自动重启"
    Write-Host ""
}

# 检查 WSL 状态
function Test-WSL {
    try {
        $result = wsl --status 2>$null
        if ($LASTEXITCODE -eq 0) {
            return $true
        }
    } catch {
        return $false
    }
    return $false
}

# 检查环境变量
function Test-Environment {
    Write-Host "🔍 检查环境配置..." -ForegroundColor $Blue
    
    if (-not (Test-Path ".env")) {
        if (Test-Path ".env.example") {
            Write-Host "📝 复制环境变量配置..." -ForegroundColor $Yellow
            Copy-Item ".env.example" ".env"
            Write-Host "✅ 已创建 .env 文件" -ForegroundColor $Green
        } else {
            Write-Host "⚠️ 未找到环境变量配置文件" -ForegroundColor $Yellow
            return $false
        }
    }
    
    Write-Host "✅ 环境配置检查通过" -ForegroundColor $Green
    return $true
}

# 检查项目依赖
function Test-Dependencies {
    Write-Host "📦 检查项目依赖..." -ForegroundColor $Blue
    
    $missing = @()
    
    # 检查后端依赖
    if (-not (Test-Path "server/go.mod")) {
        $missing += "server/go.mod"
    }
    
    # 检查管理端依赖
    if (-not (Test-Path "admin-web/package.json")) {
        $missing += "admin-web/package.json"
    }
    
    # 检查小程序端依赖
    if (-not (Test-Path "miniprogram/package.json")) {
        $missing += "miniprogram/package.json"
    }
    
    if ($missing.Count -gt 0) {
        Write-Host "❌ 缺少项目文件:" -ForegroundColor $Red
        foreach ($file in $missing) {
            Write-Host "  - $file" -ForegroundColor $Red
        }
        return $false
    }
    
    Write-Host "✅ 项目依赖检查通过" -ForegroundColor $Green
    return $true
}

# 安装依赖
function Install-Dependencies {
    Write-Host "📦 安装项目依赖..." -ForegroundColor $Blue
    
    # 安装管理端依赖
    if (Test-Path "admin-web/package.json") {
        Write-Host "  安装管理端依赖..." -ForegroundColor $Cyan
        $result = wsl -d Ubuntu bash -c "cd /mnt/d/GITVIEW/qa/admin-web && npm install"
        if ($LASTEXITCODE -ne 0) {
            Write-Host "❌ 管理端依赖安装失败" -ForegroundColor $Red
            return $false
        }
    }
    
    # 安装小程序端依赖
    if (Test-Path "miniprogram/package.json") {
        Write-Host "  安装小程序端依赖..." -ForegroundColor $Cyan
        $result = wsl -d Ubuntu bash -c "cd /mnt/d/GITVIEW/qa/miniprogram && npm install"
        if ($LASTEXITCODE -ne 0) {
            Write-Host "❌ 小程序端依赖安装失败" -ForegroundColor $Red
            return $false
        }
    }
    
    # 安装后端依赖
    if (Test-Path "server/go.mod") {
        Write-Host "  安装后端依赖..." -ForegroundColor $Cyan
        $result = wsl -d Ubuntu bash -c "cd /mnt/d/GITVIEW/qa/server && go mod download"
        if ($LASTEXITCODE -ne 0) {
            Write-Host "❌ 后端依赖安装失败" -ForegroundColor $Red
            return $false
        }
    }
    
    Write-Host "✅ 所有依赖安装完成" -ForegroundColor $Green
    return $true
}

# 启动后端服务
function Start-Backend {
    Write-Host "🔧 启动后端服务..." -ForegroundColor $Blue
    
    $title = "🔧 后端服务 - Go Server (端口: 8080)"
    $command = "Write-Host '$title' -ForegroundColor Green; Write-Host ''; wsl -d Ubuntu bash -c 'cd /mnt/d/GITVIEW/qa/server && echo \"🚀 启动 Go 服务器...\" && go run main.go'"
    
    Start-Process powershell -ArgumentList @(
        "-NoExit",
        "-Command",
        $command
    )
    
    Write-Host "✅ 后端服务已启动" -ForegroundColor $Green
}

# 启动管理端
function Start-Admin {
    Write-Host "🎨 启动管理端..." -ForegroundColor $Blue
    
    $title = "🎨 管理端 - Vue3 + Element Plus (端口: 3000)"
    $command = "Write-Host '$title' -ForegroundColor Blue; Write-Host ''; wsl -d Ubuntu bash -c 'cd /mnt/d/GITVIEW/qa/admin-web && echo \"🚀 启动管理端开发服务器...\" && npm run dev'"
    
    Start-Process powershell -ArgumentList @(
        "-NoExit",
        "-Command",
        $command
    )
    
    Write-Host "✅ 管理端已启动" -ForegroundColor $Green
}

# 启动小程序端
function Start-MiniProgram {
    Write-Host "📱 启动小程序端..." -ForegroundColor $Blue
    
    $title = "📱 小程序端 - Vue3 + Uni-app (端口: 5173)"
    $command = "Write-Host '$title' -ForegroundColor Magenta; Write-Host ''; wsl -d Ubuntu bash -c 'cd /mnt/d/GITVIEW/qa/miniprogram && echo \"🚀 启动小程序开发服务器...\" && npm run dev'"
    
    Start-Process powershell -ArgumentList @(
        "-NoExit",
        "-Command",
        $command
    )
    
    Write-Host "✅ 小程序端已启动" -ForegroundColor $Green
}

# 停止所有服务
function Stop-AllServices {
    Write-Host "🛑 停止所有开发服务..." -ForegroundColor $Yellow
    
    # 查找并关闭相关进程
    $processes = @(
        "go",
        "node",
        "npm",
        "vite"
    )
    
    foreach ($processName in $processes) {
        try {
            $runningProcesses = Get-Process -Name $processName -ErrorAction SilentlyContinue
            if ($runningProcesses) {
                Write-Host "  停止 $processName 进程..." -ForegroundColor $Cyan
                $runningProcesses | Stop-Process -Force
            }
        } catch {
            # 忽略错误
        }
    }
    
    # 停止 WSL 中的进程
    try {
        Write-Host "  停止 WSL 中的开发服务..." -ForegroundColor $Cyan
        wsl -d Ubuntu bash -c "pkill -f 'go run' 2>/dev/null || true"
        wsl -d Ubuntu bash -c "pkill -f 'npm run dev' 2>/dev/null || true"
        wsl -d Ubuntu bash -c "pkill -f 'vite' 2>/dev/null || true"
    } catch {
        # 忽略错误
    }
    
    Write-Host "✅ 所有服务已停止" -ForegroundColor $Green
}

# 查看服务状态
function Show-ServiceStatus {
    Write-Host "📊 服务状态检查" -ForegroundColor $Blue
    Write-Host "================" -ForegroundColor $Blue
    Write-Host ""
    
    # 检查端口占用
    $ports = @(
        @{Port=8080; Service="后端服务 (Go)"},
        @{Port=3000; Service="管理端 (Vue3)"},
        @{Port=5173; Service="小程序端 (Uni-app)"}
    )
    
    foreach ($portInfo in $ports) {
        try {
            $connection = Test-NetConnection -ComputerName localhost -Port $portInfo.Port -WarningAction SilentlyContinue
            if ($connection.TcpTestSucceeded) {
                Write-Host "✅ $($portInfo.Service) - 端口 $($portInfo.Port) 运行中" -ForegroundColor $Green
            } else {
                Write-Host "❌ $($portInfo.Service) - 端口 $($portInfo.Port) 未运行" -ForegroundColor $Red
            }
        } catch {
            Write-Host "❌ $($portInfo.Service) - 端口 $($portInfo.Port) 检查失败" -ForegroundColor $Red
        }
    }
    
    Write-Host ""
    Write-Host "🌐 访问地址:" -ForegroundColor $Blue
    Write-Host "  后端 API: http://localhost:8080" -ForegroundColor $Cyan
    Write-Host "  管理端: http://localhost:3000" -ForegroundColor $Cyan
    Write-Host "  小程序端: http://localhost:5173" -ForegroundColor $Cyan
    Write-Host ""
}

# 等待服务启动
function Wait-ForServices {
    param(
        [array]$Ports,
        [int]$TimeoutSeconds = 60
    )
    
    Write-Host "⏳ 等待服务启动..." -ForegroundColor $Blue
    
    $startTime = Get-Date
    $allStarted = $false
    
    while (-not $allStarted -and ((Get-Date) - $startTime).TotalSeconds -lt $TimeoutSeconds) {
        $allStarted = $true
        
        foreach ($port in $Ports) {
            try {
                $connection = Test-NetConnection -ComputerName localhost -Port $port -WarningAction SilentlyContinue
                if (-not $connection.TcpTestSucceeded) {
                    $allStarted = $false
                    break
                }
            } catch {
                $allStarted = $false
                break
            }
        }
        
        if (-not $allStarted) {
            Start-Sleep -Seconds 2
            Write-Host "." -NoNewline -ForegroundColor $Yellow
        }
    }
    
    Write-Host ""
    
    if ($allStarted) {
        Write-Host "✅ 所有服务已启动" -ForegroundColor $Green
        return $true
    } else {
        Write-Host "⚠️ 部分服务可能未完全启动" -ForegroundColor $Yellow
        return $false
    }
}

# 打开浏览器
function Open-Browser {
    Write-Host "🌐 打开浏览器..." -ForegroundColor $Blue
    
    try {
        Start-Process "http://localhost:3000"
        Start-Sleep -Seconds 1
        Start-Process "http://localhost:5173"
        Write-Host "✅ 浏览器已打开" -ForegroundColor $Green
    } catch {
        Write-Host "⚠️ 无法自动打开浏览器" -ForegroundColor $Yellow
    }
}

# 主函数
function Main {
    # 显示标题
    Write-Host ""
    Write-Host "🛠️ 刷刷题项目 - 开发模式" -ForegroundColor $Green
    Write-Host "==========================" -ForegroundColor $Green
    Write-Host ""
    
    # 处理帮助参数
    if ($Help) {
        Show-Help
        return
    }
    
    # 处理停止参数
    if ($Stop) {
        Stop-AllServices
        return
    }
    
    # 处理状态参数
    if ($Status) {
        Show-ServiceStatus
        return
    }
    
    # 检查 WSL
    if (-not (Test-WSL)) {
        Write-Host "❌ WSL 未正确配置" -ForegroundColor $Red
        Write-Host "请先运行: .\setup-windows.ps1" -ForegroundColor $Yellow
        return
    }
    
    # 检查环境
    if (-not (Test-Environment)) {
        Write-Host "❌ 环境配置检查失败" -ForegroundColor $Red
        return
    }
    
    # 检查依赖
    if (-not (Test-Dependencies)) {
        Write-Host "❌ 项目依赖检查失败" -ForegroundColor $Red
        return
    }
    
    # 安装依赖
    Write-Host "📦 准备开发环境..." -ForegroundColor $Blue
    if (-not (Install-Dependencies)) {
        Write-Host "❌ 依赖安装失败" -ForegroundColor $Red
        return
    }
    
    # 根据参数启动服务
    $servicesToWait = @()
    
    if ($Backend -or $All -or (-not $Backend -and -not $Admin -and -not $Mini)) {
        Start-Backend
        $servicesToWait += 8080
        Start-Sleep -Seconds 2
    }
    
    if ($Admin -or $All -or (-not $Backend -and -not $Admin -and -not $Mini)) {
        Start-Admin
        $servicesToWait += 3000
        Start-Sleep -Seconds 2
    }
    
    if ($Mini -or $All -or (-not $Backend -and -not $Admin -and -not $Mini)) {
        Start-MiniProgram
        $servicesToWait += 5173
        Start-Sleep -Seconds 2
    }
    
    # 等待服务启动
    if ($servicesToWait.Count -gt 0) {
        Write-Host ""
        Write-Host "📝 所有服务已在独立终端中启动" -ForegroundColor $Blue
        
        # 等待服务就绪
        if (Wait-ForServices -Ports $servicesToWait -TimeoutSeconds 60) {
            # 显示服务状态
            Show-ServiceStatus
            
            # 打开浏览器
            if ($Admin -or $Mini -or $All -or (-not $Backend -and -not $Admin -and -not $Mini)) {
                Open-Browser
            }
        }
    }
    
    Write-Host ""
    Write-Host "📋 常用命令:" -ForegroundColor $Blue
    Write-Host "  查看状态: .\dev-windows.ps1 -Status" -ForegroundColor $Cyan
    Write-Host "  停止服务: .\dev-windows.ps1 -Stop" -ForegroundColor $Cyan
    Write-Host "  重启服务: .\dev-windows.ps1 -Stop && .\dev-windows.ps1" -ForegroundColor $Cyan
    Write-Host ""
    Write-Host "💡 提示: 所有服务支持热重载，修改代码后会自动重新加载" -ForegroundColor $Yellow
    Write-Host ""
}

# 执行主函数
Main