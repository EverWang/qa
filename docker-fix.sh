#!/bin/bash

# Docker修复脚本 - 适用于Ubuntu WSL环境
# 解决Docker连接问题和配置WSL中的Docker

set -e

echo "🔧 开始修复Ubuntu WSL中的Docker配置..."

# 检查是否在WSL环境中
if ! grep -q microsoft /proc/version; then
    echo "❌ 此脚本仅适用于WSL环境"
    exit 1
fi

# 函数：安装Docker
install_docker() {
    echo "📦 安装Docker..."
    
    # 更新包索引
    sudo apt-get update
    
    # 安装必要的包
    sudo apt-get install -y \
        ca-certificates \
        curl \
        gnupg \
        lsb-release
    
    # 添加Docker官方GPG密钥
    sudo mkdir -p /etc/apt/keyrings
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
    
    # 设置Docker仓库
    echo \
        "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
        $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
    
    # 更新包索引
    sudo apt-get update
    
    # 安装Docker Engine
    sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
    
    echo "✅ Docker安装完成"
}

# 函数：配置Docker for WSL
configure_docker_wsl() {
    echo "⚙️ 配置Docker for WSL..."
    
    # 将用户添加到docker组
    sudo usermod -aG docker $USER
    
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
    
    echo "✅ Docker WSL配置完成"
}

# 函数：启动Docker服务
start_docker() {
    echo "🚀 启动Docker服务..."
    
    # 在WSL中启动Docker daemon
    if ! pgrep dockerd > /dev/null; then
        echo "启动Docker daemon..."
        sudo dockerd > /dev/null 2>&1 &
        sleep 5
    fi
    
    # 检查Docker是否运行
    if docker version > /dev/null 2>&1; then
        echo "✅ Docker服务运行正常"
    else
        echo "❌ Docker服务启动失败"
        return 1
    fi
}

# 函数：安装Docker Compose
install_docker_compose() {
    echo "📦 安装Docker Compose..."
    
    # 下载Docker Compose
    sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    
    # 设置执行权限
    sudo chmod +x /usr/local/bin/docker-compose
    
    # 创建符号链接
    sudo ln -sf /usr/local/bin/docker-compose /usr/bin/docker-compose
    
    echo "✅ Docker Compose安装完成"
}

# 函数：创建Docker启动脚本
create_docker_start_script() {
    echo "📝 创建Docker启动脚本..."
    
    cat > start-docker.sh << 'EOF'
#!/bin/bash
# Docker启动脚本 for WSL

echo "🐳 启动Docker服务..."

# 检查Docker是否已运行
if docker version > /dev/null 2>&1; then
    echo "✅ Docker已在运行"
    exit 0
fi

# 启动Docker daemon
echo "启动Docker daemon..."
sudo dockerd > /dev/null 2>&1 &

# 等待Docker启动
echo "等待Docker启动..."
for i in {1..30}; do
    if docker version > /dev/null 2>&1; then
        echo "✅ Docker启动成功"
        exit 0
    fi
    sleep 1
done

echo "❌ Docker启动失败"
exit 1
EOF
    
    chmod +x start-docker.sh
    echo "✅ Docker启动脚本创建完成: ./start-docker.sh"
}

# 主执行流程
main() {
    echo "🔍 检查Docker安装状态..."
    
    # 检查Docker是否已安装
    if ! command -v docker &> /dev/null; then
        echo "Docker未安装，开始安装..."
        install_docker
    else
        echo "✅ Docker已安装"
    fi
    
    # 配置Docker for WSL
    configure_docker_wsl
    
    # 检查Docker Compose
    if ! command -v docker-compose &> /dev/null; then
        echo "Docker Compose未安装，开始安装..."
        install_docker_compose
    else
        echo "✅ Docker Compose已安装"
    fi
    
    # 创建启动脚本
    create_docker_start_script
    
    # 尝试启动Docker
    start_docker
    
    echo ""
    echo "🎉 Docker修复完成！"
    echo ""
    echo "📋 使用说明："
    echo "1. 重新登录或运行: newgrp docker"
    echo "2. 启动Docker: ./start-docker.sh"
    echo "3. 运行项目: docker-compose -f docker-compose.local.yml up -d"
    echo ""
    echo "💡 提示：每次重启WSL后需要运行 ./start-docker.sh 启动Docker"
}

# 执行主函数
main