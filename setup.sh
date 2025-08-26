#!/bin/bash

# 刷刷题项目本地部署脚本
# 适用于Ubuntu WSL环境

set -e

echo "🚀 开始部署刷刷题项目..."

# 项目根目录
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$PROJECT_ROOT"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查是否在 WSL 环境中
check_wsl() {
    if grep -q Microsoft /proc/version; then
        echo -e "${GREEN}✓ 检测到 WSL 环境${NC}"
    else
        echo -e "${YELLOW}⚠ 未检测到 WSL 环境，脚本可能需要调整${NC}"
    fi
}

# 检查Docker
check_docker() {
    log_info "检查Docker安装状态..."
    
    if ! command -v docker &> /dev/null; then
        log_warning "Docker未安装"
        echo ""
        echo "选择安装方式："
        echo "1) 自动修复Docker配置 (推荐)"
        echo "2) 使用原生部署 (不使用Docker)"
        echo "3) 手动安装Docker后重试"
        read -p "请选择 [1-3]: " choice
        
        case $choice in
            1)
                log_info "运行Docker修复脚本..."
                chmod +x docker-fix.sh
                ./docker-fix.sh
                ;;
            2)
                log_info "使用原生部署方式..."
                chmod +x native-deploy.sh
                ./native-deploy.sh
                exit 0
                ;;
            3)
                log_info "请手动安装Docker后重新运行此脚本"
                echo "安装命令: sudo apt-get update && sudo apt-get install docker.io"
                exit 1
                ;;
            *)
                log_error "无效选择"
                exit 1
                ;;
        esac
    fi
    
    # 检查Docker服务状态
    if ! docker version &> /dev/null; then
        log_warning "Docker服务未运行，尝试启动..."
        
        # 尝试启动Docker服务
        if sudo systemctl start docker 2>/dev/null; then
            sudo systemctl enable docker
            log_success "Docker服务启动成功"
        else
            log_warning "系统服务启动失败，尝试手动启动Docker daemon..."
            
            # 检查是否有Docker修复脚本
            if [ -f "start-docker.sh" ]; then
                ./start-docker.sh
            elif [ -f "docker-fix.sh" ]; then
                log_info "运行Docker修复脚本..."
                chmod +x docker-fix.sh
                ./docker-fix.sh
            else
                log_error "Docker启动失败，建议运行修复脚本"
                echo "请运行: chmod +x docker-fix.sh && ./docker-fix.sh"
                exit 1
            fi
        fi
    fi
    
    log_success "Docker检查完成"
}

# 检查并安装 Docker Compose
install_docker_compose() {
    echo -e "${BLUE}📦 检查 Docker Compose 安装状态...${NC}"
    
    if ! command -v docker-compose &> /dev/null; then
        echo -e "${YELLOW}⚠ Docker Compose 未安装，开始安装...${NC}"
        
        # 下载 Docker Compose
        sudo curl -L "https://github.com/docker/compose/releases/download/v2.20.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
        
        # 添加执行权限
        sudo chmod +x /usr/local/bin/docker-compose
        
        echo -e "${GREEN}✓ Docker Compose 安装完成${NC}"
    else
        echo -e "${GREEN}✓ Docker Compose 已安装${NC}"
    fi
}

# 检查并安装 Node.js
install_nodejs() {
    echo -e "${BLUE}📦 检查 Node.js 安装状态...${NC}"
    
    if ! command -v node &> /dev/null; then
        echo -e "${YELLOW}⚠ Node.js 未安装，开始安装...${NC}"
        
        # 安装 NodeSource 仓库
        curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
        
        # 安装 Node.js
        sudo apt-get install -y nodejs
        
        echo -e "${GREEN}✓ Node.js 安装完成${NC}"
    else
        echo -e "${GREEN}✓ Node.js 已安装 ($(node --version))${NC}"
    fi
}

# 检查并安装 Go
install_golang() {
    echo -e "${BLUE}📦 检查 Go 安装状态...${NC}"
    
    if ! command -v go &> /dev/null; then
        echo -e "${YELLOW}⚠ Go 未安装，开始安装...${NC}"
        
        # 下载 Go
        cd /tmp
        wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
        
        # 解压到 /usr/local
        sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
        
        # 添加到 PATH
        echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
        echo 'export GOPATH=$HOME/go' >> ~/.bashrc
        echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.bashrc
        
        # 重新加载 bashrc
        source ~/.bashrc
        
        echo -e "${GREEN}✓ Go 安装完成${NC}"
        echo -e "${YELLOW}⚠ 请重新打开终端或运行 'source ~/.bashrc'${NC}"
    else
        echo -e "${GREEN}✓ Go 已安装 ($(go version))${NC}"
    fi
}

# 创建环境变量文件
setup_env_files() {
    echo -e "${BLUE}📝 设置环境变量文件...${NC}"
    
    # 创建 .env 文件
    if [ ! -f .env ]; then
        cat > .env << EOF
# Supabase 配置
SUPABASE_URL=your_supabase_url
SUPABASE_ANON_KEY=your_supabase_anon_key
SUPABASE_SERVICE_ROLE_KEY=your_supabase_service_role_key

# JWT 密钥
JWT_SECRET=your-jwt-secret-key-change-this-in-production

# 数据库配置（如果使用本地数据库）
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=shuashuati

# 应用配置
APP_ENV=development
APP_PORT=8080
APP_DEBUG=true

# 前端配置
VITE_API_BASE_URL=http://localhost:8080
VITE_SUPABASE_URL=your_supabase_url
VITE_SUPABASE_ANON_KEY=your_supabase_anon_key
EOF
        echo -e "${GREEN}✓ 创建 .env 文件${NC}"
        echo -e "${YELLOW}⚠ 请编辑 .env 文件，填入正确的 Supabase 配置${NC}"
    else
        echo -e "${GREEN}✓ .env 文件已存在${NC}"
    fi
}

# 设置 hosts 文件
setup_hosts() {
    echo -e "${BLUE}🌐 设置本地域名...${NC}"
    
    # 检查是否已添加本地域名
    if ! grep -q "qa.local" /etc/hosts; then
        echo -e "${YELLOW}⚠ 添加本地域名到 /etc/hosts${NC}"
        
        sudo tee -a /etc/hosts > /dev/null << EOF

# 刷刷题项目本地域名
127.0.0.1 admin.qa.local
127.0.0.1 app.qa.local
127.0.0.1 api.qa.local
EOF
        echo -e "${GREEN}✓ 本地域名设置完成${NC}"
    else
        echo -e "${GREEN}✓ 本地域名已设置${NC}"
    fi
}

# 构建和启动服务
start_services() {
    echo -e "${BLUE}🚀 构建和启动服务...${NC}"
    
    # 停止可能正在运行的服务
    docker-compose -f docker-compose.local.yml down
    
    # 构建镜像
    echo -e "${BLUE}🔨 构建 Docker 镜像...${NC}"
    docker-compose -f docker-compose.local.yml build
    
    # 启动服务
    echo -e "${BLUE}▶️ 启动服务...${NC}"
    docker-compose -f docker-compose.local.yml up -d
    
    # 等待服务启动
    echo -e "${BLUE}⏳ 等待服务启动...${NC}"
    sleep 10
    
    # 检查服务状态
    echo -e "${BLUE}📊 检查服务状态...${NC}"
    docker-compose -f docker-compose.local.yml ps
}

# 显示访问信息
show_access_info() {
    echo -e "${GREEN}🎉 部署完成！${NC}"
    echo -e "${BLUE}📱 访问地址：${NC}"
    echo -e "  管理端: ${GREEN}http://localhost/admin/${NC} 或 ${GREEN}http://admin.qa.local${NC}"
    echo -e "  小程序端: ${GREEN}http://localhost/app/${NC} 或 ${GREEN}http://app.qa.local${NC}"
    echo -e "  API 接口: ${GREEN}http://localhost/api/${NC} 或 ${GREEN}http://api.qa.local${NC}"
    echo -e "  直接访问: ${GREEN}http://localhost:3000${NC} (管理端), ${GREEN}http://localhost:3001${NC} (小程序端), ${GREEN}http://localhost:8080${NC} (API)"
    echo ""
    echo -e "${BLUE}🛠 常用命令：${NC}"
    echo -e "  查看日志: ${YELLOW}docker-compose -f docker-compose.local.yml logs -f${NC}"
    echo -e "  停止服务: ${YELLOW}docker-compose -f docker-compose.local.yml down${NC}"
    echo -e "  重启服务: ${YELLOW}docker-compose -f docker-compose.local.yml restart${NC}"
    echo -e "  重新构建: ${YELLOW}docker-compose -f docker-compose.local.yml up --build -d${NC}"
    echo ""
    echo -e "${YELLOW}⚠ 注意：请确保已正确配置 .env 文件中的 Supabase 连接信息${NC}"
}

# 主函数
main() {
    log_info "开始部署刷刷题项目"
    echo "项目目录: $PROJECT_ROOT"
    echo ""
    
    check_wsl
    check_docker
    install_docker_compose
    install_nodejs
    install_golang
    setup_env_files
    setup_hosts
    start_services
    show_access_info
    
    echo ""
    log_success "部署完成！"
}

# 检查是否以root权限运行
if [ "$EUID" -eq 0 ]; then
    log_error "请不要以root权限运行此脚本"
    exit 1
fi

# 运行主函数
main "$@"