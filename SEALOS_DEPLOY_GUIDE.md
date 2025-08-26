# 刷刷题项目 Sealos 部署指南

本指南将帮助您将刷刷题项目部署到 Sealos 云平台，包括使用 DevBox 进行云端开发。

## 📋 目录

- [项目概述](#项目概述)
- [前置要求](#前置要求)
- [DevBox 云端开发环境](#devbox-云端开发环境)
- [部署步骤](#部署步骤)
- [访问地址](#访问地址)
- [监控和维护](#监控和维护)
- [故障排除](#故障排除)

## 🎯 项目概述

刷刷题项目包含三个主要组件：

- **管理端 (admin-web)**: Vue3 + Element Plus 管理界面
- **小程序端 (miniprogram)**: Vue3 + uni-app 用户界面
- **后端服务 (server)**: Golang + Gin API 服务
- **数据库**: Supabase PostgreSQL

## 🔧 前置要求

1. **Sealos 账户**: 注册 [Sealos](https://cloud.sealos.io) 账户
2. **域名配置**: 确保您有权限配置以下域名（或使用 Sealos 提供的默认域名）：
   - `shuashuati-admin.cloud.sealos.io` (管理端)
   - `shuashuati-app.cloud.sealos.io` (小程序端)
   - `shuashuati-server.cloud.sealos.io` (后端API)
3. **Supabase 配置**: 确保 Supabase 项目已配置完成

## 🚀 DevBox 云端开发环境

### 创建 DevBox

1. 登录 Sealos 控制台
2. 点击 "DevBox" 应用
3. 点击 "创建 DevBox"
4. 上传项目的 `devbox.yaml` 配置文件
5. 点击 "创建" 开始部署

### DevBox 功能特性

- **多端口支持**: 
  - 5173: 管理端开发服务器
  - 5174: 小程序端开发服务器
  - 8080: 后端API服务
- **预装环境**: Node.js 18, Go 1.21
- **VS Code 集成**: 支持在线代码编辑
- **SSH 访问**: 支持终端访问
- **自动保存**: 每5分钟自动保存工作进度

### 使用 DevBox 开发

```bash
# 进入 DevBox 后，启动各个服务

# 启动管理端开发服务器
cd /workspace/admin-web
npm run dev

# 启动小程序端开发服务器
cd /workspace/miniprogram
npm run dev

# 启动后端服务（如果需要本地开发）
cd /workspace/server
go run main.go
```

## 📦 部署步骤

### 1. 准备源码

首先需要将源码上传到 Sealos，可以通过以下方式：

#### 方式一：使用 Git（推荐）

```bash
# 在 DevBox 中克隆代码
git clone <your-repo-url> /workspace
```

#### 方式二：创建 ConfigMap

```bash
# 为每个组件创建源码 ConfigMap
kubectl create configmap admin-web-source --from-file=admin-web/
kubectl create configmap miniprogram-source --from-file=miniprogram/
kubectl create configmap server-source --from-file=server/
```

### 2. 部署后端服务

```bash
# 部署后端服务
kubectl apply -f server/sealos-deploy.yaml

# 检查部署状态
kubectl get pods -l app=server
kubectl get svc server-service
kubectl get ingress server-ingress
```

### 3. 部署管理端

```bash
# 部署管理端
kubectl apply -f admin-web/sealos-deploy.yaml

# 检查部署状态
kubectl get pods -l app=admin-web
kubectl get svc admin-web-service
kubectl get ingress admin-web-ingress
```

### 4. 部署小程序端

```bash
# 部署小程序端
kubectl apply -f miniprogram/sealos-deploy.yaml

# 检查部署状态
kubectl get pods -l app=miniprogram
kubectl get svc miniprogram-service
kubectl get ingress miniprogram-ingress
```

### 5. 验证部署

```bash
# 检查所有服务状态
kubectl get all

# 检查 Ingress 状态
kubectl get ingress

# 查看服务日志
kubectl logs -l app=server
kubectl logs -l app=admin-web
kubectl logs -l app=miniprogram
```

## 🌐 访问地址

部署完成后，您可以通过以下地址访问各个服务：

- **管理端**: https://shuashuati-admin.cloud.sealos.io
- **小程序端**: https://shuashuati-app.cloud.sealos.io
- **后端API**: https://shuashuati-server.cloud.sealos.io
- **DevBox**: 通过 Sealos 控制台访问

## 📊 监控和维护

### 查看资源使用情况

```bash
# 查看 Pod 资源使用
kubectl top pods

# 查看节点资源使用
kubectl top nodes

# 查看 HPA 状态
kubectl get hpa
```

### 扩缩容操作

```bash
# 手动扩容
kubectl scale deployment server --replicas=5
kubectl scale deployment admin-web --replicas=3
kubectl scale deployment miniprogram --replicas=3

# 查看扩容状态
kubectl get deployments
```

### 更新部署

```bash
# 更新配置
kubectl apply -f server/sealos-deploy.yaml
kubectl apply -f admin-web/sealos-deploy.yaml
kubectl apply -f miniprogram/sealos-deploy.yaml

# 滚动更新
kubectl rollout restart deployment/server
kubectl rollout restart deployment/admin-web
kubectl rollout restart deployment/miniprogram

# 查看更新状态
kubectl rollout status deployment/server
```

## 🔧 故障排除

### 常见问题

#### 1. Pod 启动失败

```bash
# 查看 Pod 详细信息
kubectl describe pod <pod-name>

# 查看 Pod 日志
kubectl logs <pod-name>

# 查看事件
kubectl get events --sort-by=.metadata.creationTimestamp
```

#### 2. 服务无法访问

```bash
# 检查服务状态
kubectl get svc
kubectl describe svc <service-name>

# 检查 Ingress 状态
kubectl get ingress
kubectl describe ingress <ingress-name>

# 检查 DNS 解析
nslookup shuashuati-server.cloud.sealos.io
```

#### 3. 构建失败

```bash
# 查看 InitContainer 日志
kubectl logs <pod-name> -c build-admin-web
kubectl logs <pod-name> -c build-miniprogram

# 检查 ConfigMap
kubectl get configmap
kubectl describe configmap <configmap-name>
```

#### 4. 数据库连接问题

```bash
# 检查环境变量
kubectl exec -it <pod-name> -- env | grep SUPABASE

# 测试数据库连接
kubectl exec -it <pod-name> -- curl -I $SUPABASE_URL
```

### 日志收集

```bash
# 收集所有服务日志
kubectl logs -l app=server > server.log
kubectl logs -l app=admin-web > admin-web.log
kubectl logs -l app=miniprogram > miniprogram.log

# 实时查看日志
kubectl logs -f -l app=server
```

### 性能优化

1. **资源调优**:
   - 根据实际使用情况调整 CPU 和内存限制
   - 配置合适的 HPA 策略

2. **缓存优化**:
   - 启用 Nginx 静态资源缓存
   - 配置 CDN 加速

3. **数据库优化**:
   - 优化 Supabase 查询
   - 配置连接池

## 🔐 安全配置

### 环境变量管理

```bash
# 更新 Secret
kubectl create secret generic server-secrets \
  --from-literal=SUPABASE_URL=<your-url> \
  --from-literal=SUPABASE_SERVICE_ROLE_KEY=<your-key> \
  --from-literal=JWT_SECRET=<your-secret> \
  --dry-run=client -o yaml | kubectl apply -f -
```

### HTTPS 配置

所有服务都已配置 HTTPS，使用 Let's Encrypt 自动证书。

### CORS 配置

后端服务已配置 CORS，允许来自管理端和小程序端的请求。

## 📞 技术支持

如果在部署过程中遇到问题，可以：

1. 查看 [Sealos 官方文档](https://docs.sealos.io/)
2. 检查项目的 GitHub Issues
3. 联系技术支持团队

---

**注意**: 请确保在生产环境中使用强密码和安全的环境变量配置。定期备份数据和配置文件。