# 刷刷题用户端前端 Dockerfile
# 使用已构建的dist目录

FROM docker.m.daocloud.io/library/nginx:alpine
# 复制构建产物
COPY dist /usr/share/nginx/html

# 复制 nginx 配置
COPY nginx.conf /etc/nginx/nginx.conf

# 暴露端口
EXPOSE 80

# 启动 nginx
CMD ["nginx", "-g", "daemon off;"]