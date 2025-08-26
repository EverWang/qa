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
