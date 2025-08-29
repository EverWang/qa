# 用户操作规则 - 容器化部署指南

## 概述
本文档总结了在WSL+Ubuntu+Podman环境下进行容器化部署时用户需要遵循的操作规则和最佳实践。

## 环境要求

### 基础环境
- Windows 11 + WSL2
- Ubuntu 22.04 LTS
- Podman 4.9.3+
- podman-compose 1.0.6+

### 开发环境
- Go 1.23.4+ (用于后端编译)
- Node.js 22.9.0+ (用于前端构建)

## 容器化部署规则

### 1. 环境准备规则

#### 1.1 WSL环境配置
- **必须在WSL环境中执行所有容器相关操作**
- 不允许在Windows宿主机上直接运行服务
- 确保WSL中已安装podman和podman-compose

#### 1.2 网络代理处理
- 如果WSL环境中存在代理配置，需要在Dockerfile中显式清空代理环境变量
- 设置以下环境变量为空值：
  ```dockerfile
  ENV HTTP_PROXY=
  ENV HTTPS_PROXY=
  ENV http_proxy=
  ENV https_proxy=
  ```

### 2. 后端服务部署规则

#### 2.1 Go二进制文件编译
- **必须编译Linux版本的二进制文件**
- Windows版本的.exe文件无法在Linux容器中运行
- 编译命令：
  ```bash
  GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main-linux .
  ```

#### 2.2 端口配置规则
- **避免在代码中硬编码端口号**
- 使用环境变量PORT或SERVER_PORT
- 优先级：PORT > SERVER_PORT > 默认值
- 确保docker-compose.yml中的端口映射与代码中的端口一致

#### 2.3 Dockerfile最佳实践
- 使用多阶段构建或直接使用预编译的二进制文件
- 对于简单的Go应用，可以使用scratch基础镜像
- 确保二进制文件具有执行权限

### 3. 前端服务部署规则

#### 3.1 Nginx配置规则
- **确保upstream服务器名称与docker-compose.yml中的服务名一致**
- 后端服务名应为`backend`而不是`qa_server_1`
- API代理配置示例：
  ```nginx
  location /api/ {
      proxy_pass http://backend:8080;
      proxy_set_header Host $host;
      proxy_set_header X-Real-IP $remote_addr;
      proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
      proxy_set_header X-Forwarded-Proto $scheme;
  }
  ```

#### 3.2 构建优化
- 使用已构建的dist目录而不是在容器中重新构建
- 简化Dockerfile，直接复制dist目录到nginx html目录
- 使用nginx:alpine基础镜像以减小镜像大小

### 4. 数据库部署规则

#### 4.1 MySQL配置
- 使用官方MySQL 8.0镜像
- 配置健康检查确保数据库完全启动后再启动依赖服务
- 设置合适的字符集和排序规则（utf8mb4）

### 5. 容器编排规则

#### 5.1 docker-compose.yml配置
- **避免使用复杂的depends_on条件**
- podman-compose对`condition: service_healthy`支持有限
- 使用简单的depends_on列表：
  ```yaml
  depends_on:
    - mysql
  ```

#### 5.2 网络配置
- 使用自定义网络确保容器间通信
- 设置合适的子网范围避免冲突
- 使用服务名作为容器间通信的主机名

### 6. 故障排除规则

#### 6.1 容器启动失败
- 检查容器日志：`podman logs <container_name>`
- 检查端口冲突和网络配置
- 验证镜像是否正确构建

#### 6.2 依赖问题处理
- 如果出现容器依赖错误，清理所有容器：`podman rm -af`
- 清理网络：`podman network prune -f`
- 重新启动服务：`podman-compose up -d`

#### 6.3 镜像更新
- 代码修改后需要重新构建镜像
- 停止并删除旧容器后再启动新容器
- 使用`podman build -t <image_name> .`重新构建

### 7. 安全规则

#### 7.1 权限管理
- 避免在容器中使用root用户运行应用
- 设置合适的文件权限
- 使用非特权端口（>1024）

#### 7.2 敏感信息处理
- 使用环境变量传递敏感配置
- 不要在镜像中硬编码密码和密钥
- 使用.env文件管理环境变量

### 8. 性能优化规则

#### 8.1 镜像优化
- 使用多阶段构建减小镜像大小
- 选择合适的基础镜像（alpine vs ubuntu）
- 清理不必要的文件和依赖

#### 8.2 资源限制
- 为容器设置合适的内存和CPU限制
- 监控容器资源使用情况
- 根据实际需求调整资源配置

## 常见错误及解决方案

### 错误1：二进制文件无法执行
**现象**：`exec container process (missing dynamic library?) '/app/./main': No such file or directory`
**原因**：使用了Windows版本的二进制文件
**解决**：重新编译Linux版本的二进制文件

### 错误2：端口冲突
**现象**：`bind: address already in use`
**原因**：端口被其他进程占用或配置不一致
**解决**：检查端口配置，确保代码和docker-compose.yml中的端口一致

### 错误3：nginx upstream错误
**现象**：`host not found in upstream "qa_server_1"`
**原因**：nginx配置中的服务名与实际服务名不匹配
**解决**：修改nginx.conf中的upstream服务名为正确的服务名

### 错误4：容器依赖错误
**现象**：`container depends on container not found`
**原因**：容器依赖关系配置错误或旧容器残留
**解决**：清理所有容器和网络，重新启动服务

## 检查清单

部署前请确认以下项目：

- [ ] WSL环境已正确配置
- [ ] Go环境已安装（如需编译）
- [ ] 已编译Linux版本的后端二进制文件
- [ ] Dockerfile中的端口配置正确
- [ ] nginx配置中的upstream服务名正确
- [ ] docker-compose.yml中的服务依赖关系简化
- [ ] 环境变量配置正确
- [ ] 网络配置无冲突

## 维护建议

1. **定期更新**：定期更新基础镜像和依赖
2. **监控日志**：定期检查容器日志，及时发现问题
3. **备份数据**：定期备份数据库和重要配置文件
4. **文档更新**：及时更新部署文档和操作手册
5. **测试验证**：每次部署后进行功能测试验证

---

*本文档基于实际部署经验总结，如有问题请及时更新和完善。*