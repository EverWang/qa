# 使用Node.js作为构建环境
FROM docker.m.daocloud.io/node:18-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制package文件
COPY package*.json ./

# 安装依赖（包括开发依赖，用于构建）
RUN npm ci

# 复制源代码
COPY . .

# 设置环境变量
ENV VITE_API_BASE_URL=http://localhost:8080/api/v1

# 构建应用
RUN npm run build

# 使用nginx作为运行环境
FROM docker.m.daocloud.io/nginx:alpine

# 复制构建产物到nginx目录
COPY --from=builder /app/dist /usr/share/nginx/html

# 复制nginx配置
COPY nginx.conf /etc/nginx/conf.d/default.conf

# 暴露端口
EXPOSE 80

# 启动nginx
CMD ["nginx", "-g", "daemon off;"]