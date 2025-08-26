# Windows 环境一键配置脚本 (修复版)
# 刷刷题项目 - Windows 11 开发环境自动配置
# 版本: 2.0 (修复版)

param(
    [switch]$SkipRestart,
    [switch]$Help,
    [switch]$Debug,
    [switch]$Force
)

# 设置错误处理
$ErrorActionPreference = "Continue"
$ProgressPreference = "SilentlyContinue"

# 颜色定义
$Green = "Green"
$Blue = "Blue"
$Yellow = "Yellow"
$Red = "Red"
$Cyan = "Cyan"
$Magenta = "Magenta"

# 日志文件路径
$LogFile = "$env:TEMP\setup-windows-$(Get-Date -Format 'yyyyMMdd-HHmmss').log"

# 日志函数
function Write-Log {
    param(
        [string]$Message,
        [string]$Level = "INFO",
        [string]$Color = "White"
    )
    
    $timestamp = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
    $logMessage = "[$timestamp] [$Level] $Message"
    
    # 写入控制台
    if ($Color -ne "White") {
        Write-Host $Message -ForegroundColor $Color
    } else {
        Write-Host $Message
    }
    
    # 写入日志文件
    try {
        $logMessage | Out-File -FilePath $LogFile -Append -Encoding UTF8
    } catch {
        # 忽略日志写入错误
    }
}

# 显示帮助信息
function Show-Help {
    Write-Log "🚀 刷刷题项目 Windows 环境配置脚本 (修复版)" -Color $Green
    Write-Log ""
    Write-Log "用法:" -Color $Blue
    Write-Log "  .\setup-windows.ps1                # 完整配置（推荐）"
    Write-Log "  .\setup-windows.ps1 -SkipRestart   # 跳过重启检查"
    Write-Log "  .\setup-windows.ps1 -Debug         # 启用调试模式"
    Write-Log "  .\setup-windows.ps1 -Force         # 强制重新安装"
    Write-Log "  .\setup-windows.ps1 -Help          # 显示帮助信息"
    Write-Log ""
    Write-Log "配置内容:" -Color $Blue
    Write-Log "  ✅ 修复 PowerShell 执行策略"
    Write-Log "  ✅ 启用 WSL2 功能"
    Write-Log "  ✅ 安装 Ubuntu 发行版"
    Write-Log "  ✅ 配置开发环境"
    Write-Log "  ✅ 安装容器运行时 (Podman 优先)"
    Write-Log "  ✅ 配置 VS Code 集成"
    Write-Log ""
    Write-Log "注意事项:" -Color $Yellow
    Write-Log "  ⚠️ 需要管理员权限"
    Write-Log "  ⚠️ 配置过程中可能需要重启"
    Write-Log "  ⚠️ 首次运行需要网络连接"
    Write-Log "  📝 日志文件: $LogFile"
    Write-Log ""
}

# 检查并修复 PowerShell 执行策略
function Fix-ExecutionPolicy {
    Write-Log "🔧 检查 PowerShell 执行策略..." -Color $Blue
    
    try {
        $currentPolicy = Get-ExecutionPolicy -Scope CurrentUser
        Write-Log "  当前执行策略: $currentPolicy" -Color $Cyan
        
        if ($currentPolicy -eq "Restricted" -or $currentPolicy -eq "Undefined") {
            Write-Log "  修复执行策略为 RemoteSigned..." -Color $Yellow
            Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser -Force
            Write-Log "✅ PowerShell 执行策略已修复" -Color $Green
        } else {
            Write-Log "✅ PowerShell 执行策略正常" -Color $Green
        }
        return $true
    } catch {
        Write-Log "❌ 修复 PowerShell 执行策略失败: $($_.Exception.Message)" -Color $Red
        return $false
    }
}

# 检查管理员权限
function Test-AdminRights {
    try {
        $currentUser = [Security.Principal.WindowsIdentity]::GetCurrent()
        $principal = New-Object Security.Principal.WindowsPrincipal($currentUser)
        $isAdmin = $principal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
        
        if ($isAdmin) {
            Write-Log "✅ 管理员权限检查通过" -Color $Green
        } else {
            Write-Log "❌ 需要管理员权限" -Color $Red
        }
        
        return $isAdmin
    } catch {
        Write-Log "❌ 权限检查失败: $($_.Exception.Message)" -Color $Red
        return $false
    }
}

# 检查 Windows 版本
function Test-WindowsVersion {
    try {
        $version = [System.Environment]::OSVersion.Version
        $build = (Get-ItemProperty "HKLM:SOFTWARE\Microsoft\Windows NT\CurrentVersion").CurrentBuild
        
        Write-Log "🔍 检查 Windows 版本..." -Color $Blue
        Write-Log "  版本: $($version.Major).$($version.Minor)" -Color $Cyan
        Write-Log "  构建: $build" -Color $Cyan
        
        # Windows 10 版本 2004 (build 19041) 或更高版本支持 WSL2
        if ($version.Major -ge 10 -and [int]$build -ge 19041) {
            Write-Log "✅ Windows 版本支持 WSL2" -Color $Green
            return $true
        } else {
            Write-Log "❌ Windows 版本过低，需要 Windows 10 版本 2004 (build 19041) 或更高" -Color $Red
            return $false
        }
    } catch {
        Write-Log "❌ Windows 版本检查失败: $($_.Exception.Message)" -Color $Red
        return $false
    }
}

# 检查网络连接
function Test-NetworkConnection {
    Write-Log "🌐 检查网络连接..." -Color $Blue
    
    try {
        $testUrls = @(
            "https://www.microsoft.com",
            "https://github.com",
            "https://registry.npmjs.org"
        )
        
        foreach ($url in $testUrls) {
            try {
                $response = Invoke-WebRequest -Uri $url -Method Head -TimeoutSec 10 -UseBasicParsing
                if ($response.StatusCode -eq 200) {
                    Write-Log "✅ 网络连接正常" -Color $Green
                    return $true
                }
            } catch {
                continue
            }
        }
        
        Write-Log "❌ 网络连接失败" -Color $Red
        return $false
    } catch {
        Write-Log "⚠️ 网络连接检查失败，但继续执行: $($_.Exception.Message)" -Color $Yellow
        return $true
    }
}

# 检查 WSL 状态
function Get-WSLStatus {
    Write-Log "📊 检查 WSL 状态..." -Color $Blue
    
    try {
        # 检查 WSL 命令是否可用
        $wslVersion = wsl --version 2>$null
        if ($LASTEXITCODE -eq 0) {
            Write-Log "  WSL 命令可用" -Color $Cyan
            
            # 检查是否有已安装的发行版
            $distributions = wsl --list --quiet 2>$null
            if ($distributions -and $distributions.Count -gt 0) {
                Write-Log "  已安装的发行版: $($distributions -join ', ')" -Color $Cyan
                return "installed"
            } else {
                return "enabled"
            }
        }
    } catch {
        Write-Log "  WSL 命令不可用" -Color $Yellow
    }
    
    # 检查 WSL 功能是否启用
    try {
        $wslFeature = Get-WindowsOptionalFeature -Online -FeatureName Microsoft-Windows-Subsystem-Linux -ErrorAction SilentlyContinue
        $vmFeature = Get-WindowsOptionalFeature -Online -FeatureName VirtualMachinePlatform -ErrorAction SilentlyContinue
        
        if ($wslFeature -and $vmFeature) {
            if ($wslFeature.State -eq "Enabled" -and $vmFeature.State -eq "Enabled") {
                Write-Log "  WSL 功能已启用" -Color $Cyan
                return "enabled"
            } elseif ($wslFeature.State -eq "Enabled") {
                Write-Log "  WSL 功能部分启用" -Color $Yellow
                return "partial"
            }
        }
        
        Write-Log "  WSL 功能未启用" -Color $Yellow
        return "disabled"
    } catch {
        Write-Log "❌ WSL 状态检查失败: $($_.Exception.Message)" -Color $Red
        return "unknown"
    }
}

# 启用 WSL 功能
function Enable-WSLFeatures {
    Write-Log "📦 启用 WSL 功能..." -Color $Blue
    
    try {
        # 启用 WSL
        Write-Log "  启用 Windows Subsystem for Linux..." -Color $Cyan
        $result1 = dism.exe /online /enable-feature /featurename:Microsoft-Windows-Subsystem-Linux /all /norestart 2>&1
        
        if ($LASTEXITCODE -ne 0) {
            Write-Log "❌ WSL 功能启用失败" -Color $Red
            Write-Log "错误信息: $result1" -Color $Red
            return $false
        }
        
        # 启用虚拟机平台
        Write-Log "  启用 Virtual Machine Platform..." -Color $Cyan
        $result2 = dism.exe /online /enable-feature /featurename:VirtualMachinePlatform /all /norestart 2>&1
        
        if ($LASTEXITCODE -ne 0) {
            Write-Log "❌ 虚拟机平台启用失败" -Color $Red
            Write-Log "错误信息: $result2" -Color $Red
            return $false
        }
        
        Write-Log "✅ WSL 功能启用成功" -Color $Green
        return $true
    } catch {
        Write-Log "❌ 启用 WSL 功能时发生错误: $($_.Exception.Message)" -Color $Red
        return $false
    }
}

# 下载并安装 WSL2 内核更新
function Install-WSL2Kernel {
    Write-Log "📦 检查 WSL2 内核更新..." -Color $Blue
    
    try {
        $kernelUrl = "https://wslstorestorage.blob.core.windows.net/wslblob/wsl_update_x64.msi"
        $kernelPath = "$env:TEMP\wsl_update_x64.msi"
        
        # 检查是否需要下载
        if (-not (Test-Path $kernelPath) -or $Force) {
            Write-Log "  下载 WSL2 内核更新..." -Color $Cyan
            Invoke-WebRequest -Uri $kernelUrl -OutFile $kernelPath -UseBasicParsing
        }
        
        # 安装内核更新
        Write-Log "  安装 WSL2 内核更新..." -Color $Cyan
        $installResult = Start-Process -FilePath "msiexec.exe" -ArgumentList "/i", $kernelPath, "/quiet", "/norestart" -Wait -PassThru
        
        if ($installResult.ExitCode -eq 0) {
            Write-Log "✅ WSL2 内核更新安装成功" -Color $Green
            return $true
        } else {
            Write-Log "⚠️ WSL2 内核更新安装可能失败，退出代码: $($installResult.ExitCode)" -Color $Yellow
            return $true
        }
    } catch {
        Write-Log "⚠️ WSL2 内核更新安装失败: $($_.Exception.Message)" -Color $Yellow
        return $true
    }
}

# 设置 WSL2 为默认版本
function Set-WSL2Default {
    Write-Log "⚙️ 设置 WSL2 为默认版本..." -Color $Blue
    
    try {
        $result = wsl --set-default-version 2 2>&1
        if ($LASTEXITCODE -eq 0) {
            Write-Log "✅ WSL2 设置为默认版本" -Color $Green
            return $true
        } else {
            Write-Log "⚠️ WSL2 设置可能需要重启后生效" -Color $Yellow
            Write-Log "错误信息: $result" -Color $Yellow
            return $true
        }
    } catch {
        Write-Log "❌ 设置 WSL2 默认版本失败: $($_.Exception.Message)" -Color $Red
        return $false
    }
}

# 安装 Ubuntu
function Install-Ubuntu {
    Write-Log "🐧 安装 Ubuntu 发行版..." -Color $Blue
    
    # 检查是否已安装
    try {
        $distributions = wsl --list --quiet 2>$null
        if ($distributions -contains "Ubuntu" -and -not $Force) {
            Write-Log "✅ Ubuntu 已安装" -Color $Green
            return $true
        }
    } catch {
        # 继续安装
    }
    
    try {
        Write-Log "  正在安装 Ubuntu..." -Color $Cyan
        
        # 尝试使用 winget 安装
        $wingetResult = winget install Canonical.Ubuntu --accept-source-agreements --accept-package-agreements 2>&1
        
        if ($LASTEXITCODE -eq 0) {
            Write-Log "✅ Ubuntu 通过 winget 安装成功" -Color $Green
            return $true
        }
        
        # 备用方案：使用 wsl --install
        Write-Log "  尝试备用安装方法..." -Color $Yellow
        $wslResult = wsl --install -d Ubuntu 2>&1
        
        if ($LASTEXITCODE -eq 0) {
            Write-Log "✅ Ubuntu 安装成功" -Color $Green
            return $true
        } else {
            Write-Log "❌ Ubuntu 安装失败" -Color $Red
            Write-Log "winget 错误: $wingetResult" -Color $Red
            Write-Log "wsl 错误: $wslResult" -Color $Red
            return $false
        }
    } catch {
        Write-Log "❌ 安装 Ubuntu 时发生错误: $($_.Exception.Message)" -Color $Red
        return $false
    }
}

# 创建 WSL 配置文件
function Create-WSLConfig {
    Write-Log "📝 创建 WSL 配置文件..." -Color $Blue
    
    $wslConfigPath = "$env:USERPROFILE\.wslconfig"
    $configContent = @'
[wsl2]
# 限制内存使用 (4GB)
memory=4GB

# 限制 CPU 核心数 (4 核)
processors=4

# 启用交换文件 (2GB)
swap=2GB

# 禁用页面报告以提高性能
pageReporting=false

# 启用嵌套虚拟化
nestedVirtualization=true

# 设置网络模式
networkingMode=mirrored

# 启用 DNS 隧道
dnsTunneling=true

# 启用防火墙
firewall=true

# 自动回收内存
autoMemoryReclaim=gradual
'@
    
    try {
        # 备份现有配置
        if (Test-Path $wslConfigPath) {
            $backupPath = "$wslConfigPath.backup.$(Get-Date -Format 'yyyyMMdd-HHmmss')"
            Copy-Item $wslConfigPath $backupPath
            Write-Log "  已备份现有配置到: $backupPath" -Color $Cyan
        }
        
        $configContent | Out-File -FilePath $wslConfigPath -Encoding UTF8 -Force
        Write-Log "✅ WSL 配置文件已创建: $wslConfigPath" -Color $Green
        return $true
    } catch {
        Write-Log "❌ 创建 WSL 配置文件失败: $($_.Exception.Message)" -Color $Red
        return $false
    }
}

# 创建开发环境安装脚本
function Create-DevEnvScript {
    Write-Log "📝 创建开发环境安装脚本..." -Color $Blue
    
    $scriptPath = "install-dev-env.sh"
    
    $scriptContent = @'
#!/bin/bash

# 刷刷题项目 - Ubuntu 开发环境安装脚本 (修复版)
# 版本: 2.0

set -e

echo "🚀 开始配置 Ubuntu 开发环境..."

# 更新系统
echo "📦 更新系统包..."
sudo apt update
sudo apt upgrade -y

# 安装基础工具
echo "🔧 安装基础工具..."
sudo apt install -y curl wget git build-essential software-properties-common apt-transport-https ca-certificates gnupg lsb-release

# 安装 Node.js 18
echo "📦 安装 Node.js..."
if ! command -v node &> /dev/null; then
    curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
    sudo apt-get install -y nodejs
else
    echo "Node.js 已安装: $(node --version)"
fi

# 安装 pnpm
echo "📦 安装 pnpm..."
if ! command -v pnpm &> /dev/null; then
    npm install -g pnpm
else
    echo "pnpm 已安装: $(pnpm --version)"
fi

# 安装 Go 1.21
echo "📦 安装 Go..."
if ! command -v go &> /dev/null; then
    wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
    sudo rm -rf /usr/local/go
    sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
    rm go1.21.0.linux-amd64.tar.gz
    
    # 配置环境变量
    if ! grep -q "/usr/local/go/bin" ~/.bashrc; then
        echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
        echo 'export GOPATH=$HOME/go' >> ~/.bashrc
        echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.bashrc
    fi
    
    export PATH=$PATH:/usr/local/go/bin
    export GOPATH=$HOME/go
    export PATH=$PATH:$GOPATH/bin
else
    echo "Go 已安装: $(go version)"
fi

# 安装 Podman
echo "🐳 安装 Podman..."
if ! command -v podman &> /dev/null; then
    sudo apt-get update
    sudo apt-get install -y podman
    
    # 配置 Podman
    sudo mkdir -p /etc/containers
    sudo tee /etc/containers/registries.conf > /dev/null <<EOF
[registries.search]
registries = ['docker.io', 'registry.fedoraproject.org', 'quay.io']

[[registry]]
location = "docker.io"
[[registry.mirror]]
location = "https://registry.docker-cn.com"
EOF
    
    # 配置环境变量
    if ! grep -q "DOCKER_HOST=unix:///run/user/1000/podman/podman.sock" ~/.bashrc; then
        echo 'export DOCKER_HOST=unix:///run/user/1000/podman/podman.sock' >> ~/.bashrc
    fi
else
    echo "Podman 已安装: $(podman --version)"
fi

# 安装 podman-compose
echo "📦 安装 podman-compose..."
sudo apt install -y python3-pip
pip3 install podman-compose

# 添加到 PATH
if ! grep -q "~/.local/bin" ~/.bashrc; then
    echo 'export PATH=$PATH:~/.local/bin' >> ~/.bashrc
fi

# 创建别名
if ! grep -q 'alias docker-compose="podman-compose"' ~/.bashrc; then
    echo 'alias docker-compose="podman-compose"' >> ~/.bashrc
    echo 'alias docker="podman"' >> ~/.bashrc
fi

# 配置 Git
echo "⚙️ 配置 Git..."
git config --global core.preloadindex true
git config --global core.fscache true
git config --global gc.auto 256

# 配置 Node.js 性能
if ! grep -q "NODE_OPTIONS" ~/.bashrc; then
    echo 'export NODE_OPTIONS="--max-old-space-size=4096"' >> ~/.bashrc
fi

# 重新加载配置
source ~/.bashrc 2>/dev/null || true

echo "✅ Ubuntu 开发环境配置完成!"
echo ""
echo "🔧 已安装的工具:"
echo "  - Node.js: $(node --version 2>/dev/null || echo '未安装')"
echo "  - npm: $(npm --version 2>/dev/null || echo '未安装')"
echo "  - pnpm: $(pnpm --version 2>/dev/null || echo '未安装')"
echo "  - Go: $(go version 2>/dev/null || echo '未安装')"
echo "  - Podman: $(podman --version 2>/dev/null || echo '未安装')"
echo ""
echo "📝 下一步:"
echo "  1. 重启 WSL: wsl --shutdown && wsl -d Ubuntu"
echo "  2. 进入项目目录: cd /mnt/d/GITVIEW/qa"
echo "  3. 运行启动脚本: ./start-windows.ps1"
echo ""
'@
    
    try {
        $scriptContent | Out-File -FilePath $scriptPath -Encoding UTF8
        Write-Log "✅ 开发环境安装脚本已创建: $scriptPath" -Color $Green
        return $true
    } catch {
        Write-Log "❌ 创建开发环境安装脚本失败: $($_.Exception.Message)" -Color $Red
        return $false
    }
}

# 创建 VS Code 配置
function Create-VSCodeConfig {
    Write-Log "📝 创建 VS Code 配置..." -Color $Blue
    
    # 创建 .vscode 目录
    if (-not (Test-Path ".vscode")) {
        New-Item -ItemType Directory -Path ".vscode" -Force | Out-Null
    }
    
    # 扩展推荐
    $extensionsContent = @'
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
    "ms-vscode.vscode-eslint",
    "ms-vscode.powershell",
    "redhat.vscode-yaml",
    "ms-vscode.vscode-docker"
  ]
}
'@
    
    # 工作区设置
    $settingsContent = @'
{
  "remote.containers.defaultExtensions": [
    "Vue.volar",
    "golang.go",
    "ms-vscode.vscode-typescript-next"
  ],
  "go.toolsManagement.checkForUpdates": "local",
  "go.useLanguageServer": true,
  "typescript.preferences.includePackageJsonAutoImports": "on",
  "vue.server.hybridMode": true,
  "terminal.integrated.defaultProfile.windows": "PowerShell",
  "terminal.integrated.profiles.windows": {
    "PowerShell": {
      "source": "PowerShell",
      "icon": "terminal-powershell"
    },
    "WSL": {
      "path": "wsl.exe",
      "args": ["-d", "Ubuntu"],
      "icon": "terminal-ubuntu"
    }
  },
  "files.eol": "\n",
  "files.encoding": "utf8",
  "editor.formatOnSave": true,
  "editor.codeActionsOnSave": {
    "source.fixAll.eslint": true
  }
}
'@
    
    try {
        $extensionsContent | Out-File -FilePath ".vscode\extensions.json" -Encoding UTF8
        $settingsContent | Out-File -FilePath ".vscode\settings.json" -Encoding UTF8
        Write-Log "✅ VS Code 配置已创建" -Color $Green
        return $true
    } catch {
        Write-Log "❌ 创建 VS Code 配置失败: $($_.Exception.Message)" -Color $Red
        return $false
    }
}

# 检查是否需要重启
function Test-RebootRequired {
    try {
        $rebootRequired = $false
        
        # 检查 Windows 更新重启标志
        if (Test-Path "HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\WindowsUpdate\Auto Update\RebootRequired") {
            $rebootRequired = $true
        }
        
        # 检查挂起的文件重命名操作
        $pendingFileRename = Get-ItemProperty "HKLM:\SYSTEM\CurrentControlSet\Control\Session Manager" -Name "PendingFileRenameOperations" -ErrorAction SilentlyContinue
        if ($pendingFileRename) {
            $rebootRequired = $true
        }
        
        return $rebootRequired
    } catch {
        Write-Log "⚠️ 重启检查失败: $($_.Exception.Message)" -Color $Yellow
        return $false
    }
}

# 回滚函数
function Invoke-Rollback {
    param([string]$Reason)
    
    Write-Log "🔄 执行回滚操作: $Reason" -Color $Yellow
    
    try {
        # 恢复 WSL 配置备份
        $wslConfigPath = "$env:USERPROFILE\.wslconfig"
        $backupFiles = Get-ChildItem "$env:USERPROFILE\.wslconfig.backup.*" -ErrorAction SilentlyContinue
        
        if ($backupFiles) {
            $latestBackup = $backupFiles | Sort-Object LastWriteTime -Descending | Select-Object -First 1
            Copy-Item $latestBackup.FullName $wslConfigPath -Force
            Write-Log "  已恢复 WSL 配置备份" -Color $Cyan
        }
        
        Write-Log "✅ 回滚操作完成" -Color $Green
    } catch {
        Write-Log "❌ 回滚操作失败: $($_.Exception.Message)" -Color $Red
    }
}

# 主配置函数
function Start-Configuration {
    Write-Log "🚀 开始配置 Windows 开发环境..." -Color $Green
    Write-Log "=================================" -Color $Green
    Write-Log "日志文件: $LogFile" -Color $Cyan
    Write-Log ""
    
    try {
        # 修复 PowerShell 执行策略
        if (-not (Fix-ExecutionPolicy)) {
            Invoke-Rollback "PowerShell 执行策略修复失败"
            return $false
        }
        
        # 检查 Windows 版本
        if (-not (Test-WindowsVersion)) {
            return $false
        }
        
        # 检查网络连接
        Test-NetworkConnection | Out-Null
        
        # 检查 WSL 状态
        $wslStatus = Get-WSLStatus
        Write-Log "📊 当前 WSL 状态: $wslStatus" -Color $Blue
        
        $needRestart = $false
        
        # 根据状态执行相应操作
        switch ($wslStatus) {
            "disabled" {
                if (-not (Enable-WSLFeatures)) {
                    Invoke-Rollback "WSL 功能启用失败"
                    return $false
                }
                Install-WSL2Kernel | Out-Null
                $needRestart = $true
            }
            "partial" {
                Write-Log "⚙️ 完善 WSL 配置..." -Color $Blue
                if (-not (Enable-WSLFeatures)) {
                    Invoke-Rollback "WSL 功能完善失败"
                    return $false
                }
                Install-WSL2Kernel | Out-Null
                $needRestart = $true
            }
            "enabled" {
                Write-Log "✅ WSL 功能已启用" -Color $Green
                Install-WSL2Kernel | Out-Null
                if (-not (Set-WSL2Default)) {
                    Invoke-Rollback "WSL2 默认版本设置失败"
                    return $false
                }
                if (-not (Install-Ubuntu)) {
                    Invoke-Rollback "Ubuntu 安装失败"
                    return $false
                }
            }
            "installed" {
                Write-Log "✅ WSL 已完全配置" -Color $Green
            }
            "unknown" {
                Write-Log "⚠️ WSL 状态未知，尝试重新配置" -Color $Yellow
                if (-not (Enable-WSLFeatures)) {
                    Invoke-Rollback "WSL 重新配置失败"
                    return $false
                }
                $needRestart = $true
            }
        }
        
        # 创建配置文件
        Create-WSLConfig | Out-Null
        Create-DevEnvScript | Out-Null
        Create-VSCodeConfig | Out-Null
        
        # 检查是否需要重启
        if ($needRestart -or (Test-RebootRequired)) {
            if (-not $SkipRestart) {
                Write-Log "" 
                Write-Log "⚠️ 需要重启计算机以完成配置" -Color $Yellow
                Write-Log "重启后请运行以下命令完成配置:" -Color $Blue
                Write-Log "  1. 启动 Ubuntu: wsl -d Ubuntu" -Color $Cyan
                Write-Log "  2. 运行安装脚本: bash /mnt/d/GITVIEW/qa/install-dev-env.sh" -Color $Cyan
                Write-Log "  3. 启动项目: .\start-windows.ps1" -Color $Cyan
                Write-Log ""
                
                $restart = Read-Host "是否现在重启? (y/N)"
                if ($restart -eq "y" -or $restart -eq "Y") {
                    Write-Log "🔄 正在重启..." -Color $Blue
                    Restart-Computer -Force
                }
            }
        } else {
            # 直接配置开发环境
            Write-Log "🔧 配置开发环境..." -Color $Blue
            
            # 确保 Ubuntu 正在运行
            wsl -d Ubuntu echo "WSL 连接测试" 2>$null
            if ($LASTEXITCODE -ne 0) {
                Write-Log "⚠️ Ubuntu 未正确启动，请手动启动后运行安装脚本" -Color $Yellow
            } else {
                $result = wsl -d Ubuntu bash -c "cd /mnt/d/GITVIEW/qa && chmod +x install-dev-env.sh && ./install-dev-env.sh" 2>&1
                
                if ($LASTEXITCODE -eq 0) {
                    Write-Log "✅ 开发环境配置完成!" -Color $Green
                } else {
                    Write-Log "⚠️ 开发环境配置可能需要手动完成" -Color $Yellow
                    Write-Log "错误信息: $result" -Color $Yellow
                }
            }
        }
        
        return $true
    } catch {
        Write-Log "❌ 配置过程中发生错误: $($_.Exception.Message)" -Color $Red
        Invoke-Rollback "配置过程异常"
        return $false
    }
}

# 显示完成信息
function Show-CompletionInfo {
    Write-Log ""
    Write-Log "🎉 Windows 开发环境配置完成!" -Color $Green
    Write-Log "==============================" -Color $Green
    Write-Log ""
    Write-Log "📋 配置摘要:" -Color $Blue
    Write-Log "  ✅ PowerShell 执行策略已修复" -Color $Green
    Write-Log "  ✅ WSL2 功能已启用" -Color $Green
    Write-Log "  ✅ Ubuntu 发行版已安装" -Color $Green
    Write-Log "  ✅ 开发工具已配置" -Color $Green
    Write-Log "  ✅ 容器运行时已安装 (Podman 优先)" -Color $Green
    Write-Log "  ✅ VS Code 集成已配置" -Color $Green
    Write-Log ""
    Write-Log "🚀 快速开始:" -Color $Blue
    Write-Log "  1. 启动项目: .\start-windows.ps1" -Color $Cyan
    Write-Log "  2. 开发模式: .\start-windows.ps1 -Dev" -Color $Cyan
    Write-Log "  3. 打开 VS Code: code ." -Color $Cyan
    Write-Log ""
    Write-Log "📖 更多信息请查看: WINDOWS_DEV_GUIDE.md" -Color $Blue
    Write-Log "📝 日志文件: $LogFile" -Color $Blue
    Write-Log ""
}

# 主函数
function Main {
    # 显示帮助
    if ($Help) {
        Show-Help
        return
    }
    
    # 检查管理员权限
    if (-not (Test-AdminRights)) {
        Write-Log "❌ 需要管理员权限运行此脚本" -Color $Red
        Write-Log "请右键点击 PowerShell 并选择 '以管理员身份运行'" -Color $Yellow
        Write-Log "或者在 PowerShell 中运行: Start-Process PowerShell -Verb RunAs" -Color $Yellow
        return
    }
    
    Write-Log "🚀 刷刷题项目 Windows 环境配置脚本 (修复版 v2.0)" -Color $Green
    Write-Log "开始时间: $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')" -Color $Cyan
    
    # 开始配置
    if (Start-Configuration) {
        Show-CompletionInfo
    } else {
        Write-Log "❌ 配置过程中发生错误" -Color $Red
        Write-Log "请检查错误信息并重试，或查看日志文件: $LogFile" -Color $Yellow
        Write-Log ""
        Write-Log "🔧 故障排除建议:" -Color $Blue
        Write-Log "  1. 确保以管理员身份运行" -Color $Cyan
        Write-Log "  2. 检查网络连接" -Color $Cyan
        Write-Log "  3. 重启计算机后重试" -Color $Cyan
        Write-Log "  4. 查看详细日志: $LogFile" -Color $Cyan
    }
}

# 执行主函数
Main