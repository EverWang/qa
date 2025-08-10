@echo off
chcp 65001 >nul
setlocal enabledelayedexpansion

REM 刷刷题系统Docker部署脚本 (Windows版本)
REM 使用方法：
REM deploy.bat start    - 启动所有服务
REM deploy.bat stop     - 停止所有服务
REM deploy.bat restart  - 重启所有服务
REM deploy.bat logs     - 查看日志
REM deploy.bat build    - 重新构建镜像

REM 检查Docker是否安装
:check_docker
docker --version >nul 2>&1
if errorlevel 1 (
    echo [错误] Docker未安装，请先安装Docker Desktop
    pause
    exit /b 1
)

docker-compose --version >nul 2>&1
if errorlevel 1 (
    echo [错误] Docker Compose未安装，请先安装Docker Compose
    pause
    exit /b 1
)

REM 主逻辑
if "%1"=="" goto show_help
if "%1"=="start" goto start_services
if "%1"=="stop" goto stop_services
if "%1"=="restart" goto restart_services
if "%1"=="logs" goto view_logs
if "%1"=="build" goto build_images
if "%1"=="clean" goto clean_up
if "%1"=="backup" goto backup_database
if "%1"=="help" goto show_help
if "%1"=="--help" goto show_help
if "%1"=="-h" goto show_help

echo [错误] 未知命令: %1
goto show_help

REM 启动服务
:start_services
echo [信息] 正在启动刷刷题系统...
docker-compose up -d
if errorlevel 1 (
    echo [错误] 服务启动失败
    pause
    exit /b 1
)

echo [信息] 等待服务启动...
timeout /t 10 /nobreak >nul

echo [信息] 检查服务状态...
docker-compose ps

echo.
echo [成功] 服务启动完成！
echo [信息] 访问地址：
echo   小程序端: http://localhost
echo   管理端: http://localhost/admin 或 http://admin.localhost
echo   API接口: http://localhost/api
echo   数据库: localhost:3306 (用户名: qaminiprogram, 密码: qaminiprogram123)
echo.
pause
exit /b 0

REM 停止服务
:stop_services
echo [信息] 正在停止刷刷题系统...
docker-compose down
echo [成功] 服务已停止
pause
exit /b 0

REM 重启服务
:restart_services
echo [信息] 正在重启刷刷题系统...
docker-compose restart
echo [成功] 服务已重启
pause
exit /b 0

REM 查看日志
:view_logs
echo [信息] 查看服务日志...
docker-compose logs -f
exit /b 0

REM 重新构建镜像
:build_images
echo [信息] 正在重新构建镜像...
docker-compose build --no-cache
echo [成功] 镜像构建完成
pause
exit /b 0

REM 清理资源
:clean_up
echo [信息] 正在清理Docker资源...
docker-compose down -v
docker system prune -f
echo [成功] 清理完成
pause
exit /b 0

REM 备份数据库
:backup_database
echo [信息] 正在备份数据库...
if not exist backups mkdir backups
set backup_file=backups\backup_%date:~0,4%%date:~5,2%%date:~8,2%_%time:~0,2%%time:~3,2%%time:~6,2%.sql
set backup_file=%backup_file: =0%
docker-compose exec mysql mysqldump -u qaminiprogram -pqaminiprogram123 qaminiprogram > "%backup_file%"
echo [成功] 数据库备份完成: %backup_file%
pause
exit /b 0

REM 显示帮助信息
:show_help
echo 刷刷题系统Docker部署脚本 (Windows版本)
echo.
echo 用法: %0 [命令]
echo.
echo 命令:
echo   start     启动所有服务
echo   stop      停止所有服务
echo   restart   重启所有服务
echo   logs      查看服务日志
echo   build     重新构建镜像
echo   clean     清理Docker资源
echo   backup    备份数据库
echo   help      显示帮助信息
echo.
pause
exit /b 0