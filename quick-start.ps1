# Windows 快速启动脚本
# 刷刷题项目 - 一键启动解决方案

param(
    [ValidateSet("container", "dev", "auto")]
    [string]$Mode = "auto",
    [switch]$Setup,
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
    Write-Host "🚀 刷刷题项目快速启动脚本" -ForegroundColor $Green
    Write-Host ""
    Write-Host "用法:" -ForegroundColor $Blue
    Write-Host "  .\quick-start.ps1                  # 自动选择最佳启动方式"
    Write-Host "  .\quick-start.ps1 -Mode container  # 强制使用容器模式"
    Write-Host "  .\quick-start.ps1 -Mode dev        # 强制使用开发模式"
    Write-Host "  .\quick-start.ps1 -Setup           # 首次环境配置"
    Write-Host "  .\quick-start.ps1 -Help            # 显示帮助信息"
    Write-Host ""
    Write-Host "启动模式:" -ForegroundColor $Blue
    Write-Host "  🐳 容器模式: 使用 Docker/Podman 容器运行"
    Write-Host "     - 优点: 环境隔离、生产一致性"
    Write-Host "     - 缺点: 启动较慢、资源占用高"
    Write-Host ""
    Write-Host "  🛠️ 开发模式: 原生运行各个服务"
    Write-Host "     - 优点: 启动快速、热重载、易调试"
    Write-Host "     - 缺点: 需要本地环境配置"
    Write-Host ""
    Write-Host "  🤖 自动模式: 智能选择最佳方式"
    Write-Host "     - 检测环境配置"
    Write-Host "     - 自动选择最适合的启动方式"
    Write-Host ""
}

# 检测系统环境
function Get-SystemInfo {
    $info = @{
        HasWSL = $false
        HasDocker = $false
        HasPodman = $false
        HasNode = $false
        HasGo = $false
        WSLDistributions = @()
        ContainerRuntime = "none"
    }
    
    # 检测 WSL
    try {
        $wslResult = wsl --status 2>$null
        if ($LASTEXITCODE -eq 0) {
            $info.HasWSL = $true
            $distributions = wsl --list --quiet 2>$null
            if ($distributions) {
                $info.WSLDistributions = $distributions -split "`n" | Where-Object { $_ -ne "" }
            }
        }
    } catch {
        # WSL 不可用
    }
    
    # 在 WSL 中检测开发环境
    if ($info.HasWSL -and $info.WSLDistributions -contains "Ubuntu") {
        try {
            # 检测 Docker
            $dockerCheck = wsl -d Ubuntu bash -c "command -v docker" 2>$null
            if ($LASTEXITCODE -eq 0) {
                $info.HasDocker = $true
                $info.ContainerRuntime = "docker"
            }
            
            # 检测 Podman
            $podmanCheck = wsl -d Ubuntu bash -c "command -v podman" 2>$null
            if ($LASTEXITCODE -eq 0) {
                $info.HasPodman = $true
                if ($info.ContainerRuntime -eq "none") {
                    $info.ContainerRuntime = "podman"
                }
            }
            
            # 检测 Node.js
            $nodeCheck = wsl -d Ubuntu bash -c "command -v node" 2>$null
            if ($LASTEXITCODE -eq 0) {
                $info.HasNode = $true
            }
            
            # 检测 Go
            $goCheck = wsl -d Ubuntu bash -c "command -v go" 2>$null
            if ($LASTEXITCODE -eq 0) {
                $info.HasGo = $true
            }
        } catch {
            # 检测失败
        }
    }
    
    return $info
}

# 显示系统信息
function Show-SystemInfo {
    param([hashtable]$Info)
    
    Write-Host "🔍 系统环境检测" -ForegroundColor $Blue
    Write-Host "================" -ForegroundColor $Blue
    Write-Host ""
    
    # WSL 状态
    if ($Info.HasWSL) {
        Write-Host "✅ WSL2 已安装" -ForegroundColor $Green
        if ($Info.WSLDistributions.Count -gt 0) {
            Write-Host "  发行版: $($Info.WSLDistributions -join ', ')" -ForegroundColor $Cyan
        }
    } else {
        Write-Host "❌ WSL2 未安装" -ForegroundColor $Red
    }
    
    # 容器运行时
    if ($Info.HasDocker) {
        Write-Host "✅ Docker 已安装" -ForegroundColor $Green
    } else {
        Write-Host "❌ Docker 未安装" -ForegroundColor $Red
    }
    
    if ($Info.HasPodman) {
        Write-Host "✅ Podman 已安装" -ForegroundColor $Green
    } else {
        Write-Host "❌ Podman 未安装" -ForegroundColor $Red
    }
    
    # 开发环境
    if ($Info.HasNode) {
        Write-Host "✅ Node.js 已安装" -ForegroundColor $Green
    } else {
        Write-Host "❌ Node.js 未安装" -ForegroundColor $Red
    }
    
    if ($Info.HasGo) {
        Write-Host "✅ Go 已安装" -ForegroundColor $Green
    } else {
        Write-Host "❌ Go 未安装" -ForegroundColor $Red
    }
    
    Write-Host ""
}

# 推荐启动模式
function Get-RecommendedMode {
    param([hashtable]$Info)
    
    # 评分系统
    $containerScore = 0
    $devScore = 0
    
    # WSL 基础分
    if ($Info.HasWSL) {
        $containerScore += 20
        $devScore += 20
    }
    
    # 容器运行时分数
    if ($Info.HasDocker -or $Info.HasPodman) {
        $containerScore += 30
    }
    
    # 开发环境分数
    if ($Info.HasNode) {
        $devScore += 25
    }
    if ($Info.HasGo) {
        $devScore += 25
    }
    
    # 完整性加分
    if ($Info.HasNode -and $Info.HasGo) {
        $devScore += 10  # 开发环境完整
    }
    
    if (($Info.HasDocker -or $Info.HasPodman) -and $Info.HasWSL) {
        $containerScore += 10  # 容器环境完整
    }
    
    # 决定推荐模式
    if ($containerScore -ge 50 -and $devScore -ge 50) {
        # 两种模式都可用，优先推荐开发模式（更快）
        return @{
            Primary = "dev"
            Secondary = "container"
            Reason = "开发环境完整，推荐开发模式（启动更快）"
        }
    } elseif ($devScore -ge 50) {
        return @{
            Primary = "dev"
            Secondary = $null
            Reason = "开发环境完整，推荐开发模式"
        }
    } elseif ($containerScore -ge 50) {
        return @{
            Primary = "container"
            Secondary = $null
            Reason = "容器环境可用，推荐容器模式"
        }
    } else {
        return @{
            Primary = "setup"
            Secondary = $null
            Reason = "环境未配置，需要先运行环境配置"
        }
    }
}

# 环境配置
function Start-EnvironmentSetup {
    Write-Host "⚙️ 开始环境配置..." -ForegroundColor $Blue
    
    if (Test-Path "setup-windows.ps1") {
        Write-Host "🚀 运行环境配置脚本..." -ForegroundColor $Green
        & ".\setup-windows.ps1"
    } else {
        Write-Host "❌ 未找到环境配置脚本" -ForegroundColor $Red
        Write-Host "请确保 setup-windows.ps1 文件存在" -ForegroundColor $Yellow
        return $false
    }
    
    return $true
}

# 容器模式启动
function Start-ContainerMode {
    Write-Host "🐳 启动容器模式..." -ForegroundColor $Blue
    
    if (Test-Path "start-windows.ps1") {
        Write-Host "🚀 运行容器启动脚本..." -ForegroundColor $Green
        & ".\start-windows.ps1" -Container
    } else {
        Write-Host "❌ 未找到容器启动脚本" -ForegroundColor $Red
        return $false
    }
    
    return $true
}

# 开发模式启动
function Start-DevMode {
    Write-Host "🛠️ 启动开发模式..." -ForegroundColor $Blue
    
    if (Test-Path "dev-windows.ps1") {
        Write-Host "🚀 运行开发启动脚本..." -ForegroundColor $Green
        & ".\dev-windows.ps1" -All
    } else {
        Write-Host "❌ 未找到开发启动脚本" -ForegroundColor $Red
        return $false
    }
    
    return $true
}

# 显示启动后信息
function Show-PostStartInfo {
    param([string]$StartMode)
    
    Write-Host ""
    Write-Host "🎉 项目启动完成!" -ForegroundColor $Green
    Write-Host "===============" -ForegroundColor $Green
    Write-Host ""
    
    Write-Host "📊 启动模式: $StartMode" -ForegroundColor $Blue
    Write-Host ""
    
    Write-Host "🌐 访问地址:" -ForegroundColor $Blue
    if ($StartMode -eq "container") {
        Write-Host "  管理端: http://localhost/admin/" -ForegroundColor $Cyan
        Write-Host "  用户端: http://localhost/app/" -ForegroundColor $Cyan
        Write-Host "  API: http://localhost/api/" -ForegroundColor $Cyan
    } else {
        Write-Host "  管理端: http://localhost:3000" -ForegroundColor $Cyan
        Write-Host "  用户端: http://localhost:5173" -ForegroundColor $Cyan
        Write-Host "  API: http://localhost:8080" -ForegroundColor $Cyan
    }
    Write-Host ""
    
    Write-Host "📋 常用命令:" -ForegroundColor $Blue
    if ($StartMode -eq "container") {
        Write-Host "  查看状态: wsl -d Ubuntu bash -c 'cd /mnt/d/GITVIEW/qa && docker-compose -f docker-compose.local.yml ps'" -ForegroundColor $Cyan
        Write-Host "  查看日志: wsl -d Ubuntu bash -c 'cd /mnt/d/GITVIEW/qa && docker-compose -f docker-compose.local.yml logs -f'" -ForegroundColor $Cyan
        Write-Host "  停止服务: wsl -d Ubuntu bash -c 'cd /mnt/d/GITVIEW/qa && docker-compose -f docker-compose.local.yml down'" -ForegroundColor $Cyan
    } else {
        Write-Host "  查看状态: .\dev-windows.ps1 -Status" -ForegroundColor $Cyan
        Write-Host "  停止服务: .\dev-windows.ps1 -Stop" -ForegroundColor $Cyan
        Write-Host "  重启服务: .\dev-windows.ps1 -Stop && .\dev-windows.ps1" -ForegroundColor $Cyan
    }
    Write-Host ""
    
    Write-Host "💡 提示:" -ForegroundColor $Yellow
    if ($StartMode -eq "dev") {
        Write-Host "  - 开发模式支持热重载，修改代码后会自动重新加载" -ForegroundColor $Yellow
        Write-Host "  - 可以在各个终端窗口中查看实时日志" -ForegroundColor $Yellow
    } else {
        Write-Host "  - 容器模式提供生产环境一致性" -ForegroundColor $Yellow
        Write-Host "  - 使用 docker-compose 命令管理服务" -ForegroundColor $Yellow
    }
    Write-Host ""
}

# 交互式模式选择
function Select-ModeInteractive {
    param([hashtable]$Recommendation)
    
    Write-Host "🤖 智能推荐" -ForegroundColor $Blue
    Write-Host "============" -ForegroundColor $Blue
    Write-Host ""
    Write-Host "推荐模式: $($Recommendation.Primary)" -ForegroundColor $Green
    Write-Host "推荐理由: $($Recommendation.Reason)" -ForegroundColor $Cyan
    Write-Host ""
    
    Write-Host "请选择启动模式:" -ForegroundColor $Blue
    Write-Host "  1. 🛠️ 开发模式 (推荐) - 快速启动，支持热重载" -ForegroundColor $Green
    Write-Host "  2. 🐳 容器模式 - 环境隔离，生产一致性" -ForegroundColor $Blue
    Write-Host "  3. ⚙️ 环境配置 - 首次使用或重新配置" -ForegroundColor $Yellow
    Write-Host "  4. ❌ 退出" -ForegroundColor $Red
    Write-Host ""
    
    do {
        $choice = Read-Host "请输入选择 (1-4)"
        switch ($choice) {
            "1" { return "dev" }
            "2" { return "container" }
            "3" { return "setup" }
            "4" { return "exit" }
            default {
                Write-Host "❌ 无效选择，请输入 1-4" -ForegroundColor $Red
            }
        }
    } while ($true)
}

# 主函数
function Main {
    # 显示标题
    Write-Host ""
    Write-Host "🚀 刷刷题项目 - 快速启动" -ForegroundColor $Green
    Write-Host "==========================" -ForegroundColor $Green
    Write-Host ""
    
    # 处理帮助参数
    if ($Help) {
        Show-Help
        return
    }
    
    # 处理环境配置参数
    if ($Setup) {
        if (Start-EnvironmentSetup) {
            Write-Host "✅ 环境配置完成，请重新运行启动脚本" -ForegroundColor $Green
        }
        return
    }
    
    # 检测系统环境
    Write-Host "🔍 检测系统环境..." -ForegroundColor $Blue
    $systemInfo = Get-SystemInfo
    Show-SystemInfo -Info $systemInfo
    
    # 获取推荐模式
    $recommendation = Get-RecommendedMode -Info $systemInfo
    
    # 确定启动模式
    $selectedMode = $Mode
    if ($Mode -eq "auto") {
        if ($recommendation.Primary -eq "setup") {
            Write-Host "⚠️ 环境未完全配置" -ForegroundColor $Yellow
            $selectedMode = Select-ModeInteractive -Recommendation $recommendation
        } else {
            # 自动选择推荐模式
            $selectedMode = $recommendation.Primary
            Write-Host "🤖 自动选择模式: $selectedMode" -ForegroundColor $Green
            Write-Host "理由: $($recommendation.Reason)" -ForegroundColor $Cyan
            Write-Host ""
        }
    }
    
    # 处理退出
    if ($selectedMode -eq "exit") {
        Write-Host "👋 再见!" -ForegroundColor $Blue
        return
    }
    
    # 执行启动
    $success = $false
    switch ($selectedMode) {
        "setup" {
            $success = Start-EnvironmentSetup
        }
        "container" {
            $success = Start-ContainerMode
        }
        "dev" {
            $success = Start-DevMode
        }
        default {
            Write-Host "❌ 未知的启动模式: $selectedMode" -ForegroundColor $Red
            return
        }
    }
    
    # 显示启动后信息
    if ($success -and $selectedMode -ne "setup") {
        Show-PostStartInfo -StartMode $selectedMode
    } elseif ($selectedMode -eq "setup") {
        Write-Host "✅ 环境配置完成，请重新运行启动脚本" -ForegroundColor $Green
    } else {
        Write-Host "❌ 启动失败，请检查错误信息" -ForegroundColor $Red
    }
}

# 执行主函数
Main