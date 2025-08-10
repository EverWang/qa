#!/bin/bash

# 刷刷题系统Docker部署脚本
# 使用方法：
# ./deploy.sh start    - 启动所有服务
# ./deploy.sh stop     - 停止所有服务
# ./deploy.sh restart  - 重启所有服务
# ./deploy.sh logs     - 查看日志
# ./deploy.sh build    - 重新构建镜像

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 打印带颜色的消息
print_message() {
    echo -e "${2}[$(date +'%Y-%m-%d %H:%M:%S')] $1${NC}"
}

# 检查Docker是否安装
check_docker() {
    if ! command -v docker &> /dev/null; then
        print_message "Docker未安装，请先安装Docker" $RED
        exit 1
    fi
    
    if ! command -v docker-compose &> /dev/null; then
        print_message "Docker Compose未安装，请先安装Docker Compose" $RED
        exit 1
    fi
}

# 启动服务
start_services() {
    print_message "正在启动刷刷题系统..." $BLUE
    docker-compose up -d
    
    print_message "等待服务启动..." $YELLOW
    sleep 10
    
    print_message "检查服务状态..." $BLUE
    docker-compose ps
    
    print_message "服务启动完成！" $GREEN
    print_message "访问地址：" $GREEN
    print_message "  小程序端: http://localhost" $GREEN
    print_message "  管理端: http://localhost/admin 或 http://admin.localhost" $GREEN
    print_message "  API接口: http://localhost/api" $GREEN
    print_message "  数据库: localhost:3306 (用户名: qaminiprogram, 密码: qaminiprogram123)" $GREEN
}

# 停止服务
stop_services() {
    print_message "正在停止刷刷题系统..." $BLUE
    docker-compose down
    print_message "服务已停止" $GREEN
}

# 重启服务
restart_services() {
    print_message "正在重启刷刷题系统..." $BLUE
    docker-compose restart
    print_message "服务已重启" $GREEN
}

# 查看日志
view_logs() {
    print_message "查看服务日志..." $BLUE
    docker-compose logs -f
}

# 重新构建镜像
build_images() {
    print_message "正在重新构建镜像..." $BLUE
    docker-compose build --no-cache
    print_message "镜像构建完成" $GREEN
}

# 清理资源
clean_up() {
    print_message "正在清理Docker资源..." $BLUE
    docker-compose down -v
    docker system prune -f
    print_message "清理完成" $GREEN
}

# 备份数据库
backup_database() {
    print_message "正在备份数据库..." $BLUE
    mkdir -p backups
    docker-compose exec mysql mysqldump -u qaminiprogram -pqaminiprogram123 qaminiprogram > "backups/backup_$(date +%Y%m%d_%H%M%S).sql"
    print_message "数据库备份完成" $GREEN
}

# 恢复数据库
restore_database() {
    if [ -z "$2" ]; then
        print_message "请指定备份文件路径" $RED
        print_message "用法: ./deploy.sh restore <backup_file>" $YELLOW
        exit 1
    fi
    
    print_message "正在恢复数据库..." $BLUE
    docker-compose exec -T mysql mysql -u qaminiprogram -pqaminiprogram123 qaminiprogram < "$2"
    print_message "数据库恢复完成" $GREEN
}

# 显示帮助信息
show_help() {
    echo "刷刷题系统Docker部署脚本"
    echo ""
    echo "用法: $0 [命令]"
    echo ""
    echo "命令:"
    echo "  start     启动所有服务"
    echo "  stop      停止所有服务"
    echo "  restart   重启所有服务"
    echo "  logs      查看服务日志"
    echo "  build     重新构建镜像"
    echo "  clean     清理Docker资源"
    echo "  backup    备份数据库"
    echo "  restore   恢复数据库"
    echo "  help      显示帮助信息"
    echo ""
}

# 主函数
main() {
    check_docker
    
    case "$1" in
        start)
            start_services
            ;;
        stop)
            stop_services
            ;;
        restart)
            restart_services
            ;;
        logs)
            view_logs
            ;;
        build)
            build_images
            ;;
        clean)
            clean_up
            ;;
        backup)
            backup_database
            ;;
        restore)
            restore_database "$@"
            ;;
        help|--help|-h)
            show_help
            ;;
        "")
            show_help
            ;;
        *)
            print_message "未知命令: $1" $RED
            show_help
            exit 1
            ;;
    esac
}

# 执行主函数
main "$@"