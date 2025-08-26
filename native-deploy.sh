#!/bin/bash

# 原生部署脚本 - 不使用Docker的替代方案
# 适用于Ubuntu WSL环境

set -e

echo "🚀 开始原生部署刷刷题项目..."

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

# 检查依赖
check_dependencies() {
    log_info "检查系统依赖..."
    
    # 检查Node.js
    if ! command -v node &> /dev/null; then
        log_error "Node.js未安装，请先安装Node.js"
        exit 1
    fi
    
    # 检查npm
    if ! command -v npm &> /dev/null; then
        log_error "npm未安装，请先安装npm"
        exit 1
    fi
    
    # 检查Go
    if ! command -v go &> /dev/null; then
        log_warning "Go未安装，将跳过后端服务部署"
        SKIP_BACKEND=true
    fi
    
    # 检查nginx
    if ! command -v nginx &> /dev/null; then
        log_warning "Nginx未安装，将安装nginx"
        install_nginx
    fi
    
    log_success "依赖检查完成"
}

# 安装nginx
install_nginx() {
    log_info "安装Nginx..."
    sudo apt-get update
    sudo apt-get install -y nginx
    log_success "Nginx安装完成"
}

# 设置环境变量
setup_environment() {
    log_info "设置环境变量..."
    
    # 复制环境变量文件
    if [ ! -f ".env" ]; then
        if [ -f ".env.example" ]; then
            cp .env.example .env
            log_info "已创建.env文件，请检查配置"
        else
            log_warning ".env.example文件不存在"
        fi
    fi
    
    log_success "环境变量设置完成"
}

# 构建管理端
build_admin() {
    log_info "构建管理端..."
    
    cd "$PROJECT_ROOT/admin-web"
    
    # 安装依赖
    log_info "安装admin-web依赖..."
    npm install
    
    # 构建项目
    log_info "构建admin-web项目..."
    npm run build
    
    log_success "管理端构建完成"
    cd "$PROJECT_ROOT"
}

# 构建小程序端
build_miniprogram() {
    log_info "构建小程序端..."
    
    cd "$PROJECT_ROOT/miniprogram"
    
    # 检查package.json是否存在
    if [ ! -f "package.json" ]; then
        log_warning "miniprogram/package.json不存在，跳过小程序端构建"
        cd "$PROJECT_ROOT"
        return
    fi
    
    # 安装依赖
    log_info "安装miniprogram依赖..."
    npm install
    
    # 构建项目
    log_info "构建miniprogram项目..."
    npm run build
    
    log_success "小程序端构建完成"
    cd "$PROJECT_ROOT"
}

# 构建后端服务
build_backend() {
    if [ "$SKIP_BACKEND" = true ]; then
        log_warning "跳过后端服务构建（Go未安装）"
        return
    fi
    
    log_info "构建后端服务..."
    
    cd "$PROJECT_ROOT/server"
    
    # 下载Go模块
    log_info "下载Go依赖..."
    go mod download
    
    # 构建Go应用
    log_info "构建Go应用..."
    go build -o shuashuati-server .
    
    log_success "后端服务构建完成"
    cd "$PROJECT_ROOT"
}

# 配置Nginx
configure_nginx() {
    log_info "配置Nginx..."
    
    # 创建nginx配置
    sudo tee /etc/nginx/sites-available/shuashuati > /dev/null <<EOF
server {
    listen 80;
    server_name localhost;
    
    # 管理端
    location /admin/ {
        alias $PROJECT_ROOT/admin-web/dist/;
        try_files \$uri \$uri/ /admin/index.html;
        
        # 静态文件缓存
        location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg)\$ {
            expires 1y;
            add_header Cache-Control "public, immutable";
        }
    }
    
    # 小程序端
    location /app/ {
        alias $PROJECT_ROOT/miniprogram/dist/;
        try_files \$uri \$uri/ /app/index.html;
        
        # 静态文件缓存
        location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg)\$ {
            expires 1y;
            add_header Cache-Control "public, immutable";
        }
    }
    
    # API代理
    location /api/ {
        proxy_pass http://localhost:8080/;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
        
        # CORS配置
        add_header Access-Control-Allow-Origin *;
        add_header Access-Control-Allow-Methods "GET, POST, PUT, DELETE, OPTIONS";
        add_header Access-Control-Allow-Headers "Content-Type, Authorization";
        
        if (\$request_method = 'OPTIONS') {
            return 204;
        }
    }
    
    # 默认重定向到管理端
    location = / {
        return 301 /admin/;
    }
}
EOF
    
    # 启用站点
    sudo ln -sf /etc/nginx/sites-available/shuashuati /etc/nginx/sites-enabled/
    
    # 删除默认站点
    sudo rm -f /etc/nginx/sites-enabled/default
    
    # 测试nginx配置
    sudo nginx -t
    
    log_success "Nginx配置完成"
}

# 启动服务
start_services() {
    log_info "启动服务..."
    
    # 重启nginx
    log_info "重启Nginx..."
    sudo systemctl restart nginx
    sudo systemctl enable nginx
    
    # 启动后端服务（如果存在）
    if [ "$SKIP_BACKEND" != true ] && [ -f "$PROJECT_ROOT/server/shuashuati-server" ]; then
        log_info "启动后端服务..."
        cd "$PROJECT_ROOT/server"
        
        # 创建systemd服务文件
        sudo tee /etc/systemd/system/shuashuati-server.service > /dev/null <<EOF
[Unit]
Description=Shuashuati Server
After=network.target

[Service]
Type=simple
User=$USER
WorkingDirectory=$PROJECT_ROOT/server
ExecStart=$PROJECT_ROOT/server/shuashuati-server
Restart=always
RestartSec=5
Environment=PORT=8080

[Install]
WantedBy=multi-user.target
EOF
        
        # 重新加载systemd并启动服务
        sudo systemctl daemon-reload
        sudo systemctl enable shuashuati-server
        sudo systemctl start shuashuati-server
        
        cd "$PROJECT_ROOT"
    fi
    
    log_success "服务启动完成"
}

# 显示状态
show_status() {
    echo ""
    echo "🎉 部署完成！"
    echo ""
    echo "📋 访问地址："
    echo "  管理端: http://localhost/admin/"
    echo "  小程序端: http://localhost/app/"
    echo "  后端API: http://localhost/api/"
    echo ""
    echo "📊 服务状态："
    echo "  Nginx: $(sudo systemctl is-active nginx)"
    
    if [ "$SKIP_BACKEND" != true ]; then
        echo "  后端服务: $(sudo systemctl is-active shuashuati-server 2>/dev/null || echo 'not installed')"
    else
        echo "  后端服务: skipped (Go not installed)"
    fi
    
    echo ""
    echo "🔧 管理命令："
    echo "  查看Nginx状态: sudo systemctl status nginx"
    echo "  重启Nginx: sudo systemctl restart nginx"
    
    if [ "$SKIP_BACKEND" != true ]; then
        echo "  查看后端状态: sudo systemctl status shuashuati-server"
        echo "  重启后端: sudo systemctl restart shuashuati-server"
        echo "  查看后端日志: sudo journalctl -u shuashuati-server -f"
    fi
    
    echo "  查看Nginx日志: sudo tail -f /var/log/nginx/access.log"
}

# 清理函数
cleanup() {
    log_info "清理临时文件..."
    # 这里可以添加清理逻辑
}

# 错误处理
error_handler() {
    log_error "部署过程中发生错误，正在清理..."
    cleanup
    exit 1
}

# 设置错误处理
trap error_handler ERR

# 主执行流程
main() {
    echo "🚀 开始原生部署刷刷题项目"
    echo "项目目录: $PROJECT_ROOT"
    echo ""
    
    check_dependencies
    setup_environment
    build_admin
    build_miniprogram
    build_backend
    configure_nginx
    start_services
    show_status
    
    echo ""
    log_success "部署完成！项目已成功部署到本地环境"
}

# 检查是否以root权限运行
if [ "$EUID" -eq 0 ]; then
    log_error "请不要以root权限运行此脚本"
    exit 1
fi

# 执行主函数
main "$@"