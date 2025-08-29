# 刷刷题 - 在线刷题平台

一个基于 Vue3 + Go + MySQL 的现代化在线刷题平台，支持题目管理、在线答题、错题本等功能。

## 项目架构

```
刷刷题/
├── admin-web/          # 管理端前端 (Vue3 + Element Plus)
├── miniprogram/        # 用户端前端 (Vue3 + Vite)
├── server/             # 后端服务 (Go + Gin + GORM)
├── docker-compose.yml  # 容器编排配置
├── DEPLOYMENT.md       # 详细部署指南
└── README.md          # 项目说明
```

## 技术栈

### 前端
- **框架**: Vue 3 + TypeScript
- **构建工具**: Vite
- **UI 组件库**: 
  - 管理端: Element Plus
  - 用户端: 自定义组件
- **样式**: Tailwind CSS + SCSS
- **状态管理**: Pinia
- **路由**: Vue Router

### 后端
- **语言**: Go 1.21+
- **框架**: Gin
- **ORM**: GORM
- **数据库**: MySQL 8.0
- **认证**: JWT
- **文档**: 内置 API 文档

### 部署
- **容器化**: Podman/Docker
- **编排**: Podman Compose
- **环境**: WSL2 + Ubuntu 22.04
- **反向代理**: Nginx

## 功能特性

### 用户端功能
- ✅ 用户注册登录
- ✅ 题目分类浏览
- ✅ 在线答题
- ✅ 答题记录
- ✅ 错题本管理
- ✅ 成绩统计
- ✅ 题目搜索
- ✅ 个人中心

### 管理端功能
- ✅ 管理员登录
- ✅ 题目管理 (增删改查)
- ✅ 分类管理
- ✅ 用户管理
- ✅ 数据统计
- ✅ 系统设置
- ✅ 操作日志

## 快速开始

### 环境要求
- Windows 11 + WSL2 + Ubuntu 22.04
- Podman 4.0+ 或 Docker 20.0+
- Git

### 一键部署

```bash
# 1. 克隆项目
git clone <项目地址>
cd shuashuati

# 2. 创建数据目录
mkdir -p ./data/mysql

# 3. 配置环境变量
cp .env.example .env
# 编辑 .env 文件，配置数据库密码等

# 4. 启动服务
podman-compose up -d

# 5. 查看状态
podman-compose ps
```

### 访问地址

- **用户端**: http://localhost:3000
- **管理端**: http://localhost:3001
- **后端API**: http://localhost:8080

### 默认账号

**管理员账号**:
- 用户名: `admin`
- 密码: `123456`

**测试用户**:
- 用户名: `testuser`
- 密码: `123456`

## 开发指南

### 本地开发

```bash
# 后端开发
cd server
go mod tidy
go run main.go

# 管理端前端开发
cd admin-web
npm install
npm run dev

# 用户端前端开发
cd miniprogram
npm install
npm run dev
```

### 代码规范

- **命名规范**: 统一使用驼峰命名
- **数据库**: 使用下划线命名
- **提交规范**: 遵循 Conventional Commits
- **代码格式**: 使用 ESLint + Prettier (前端), gofmt (后端)

### 项目结构

```
server/                 # 后端服务
├── api/               # API 路由定义
├── config/            # 配置文件
├── controllers/       # 控制器
├── models/           # 数据模型
├── middleware/       # 中间件
├── routes/           # 路由配置
├── init.sql          # 数据库初始化脚本
└── main.go           # 入口文件

admin-web/             # 管理端前端
├── src/
│   ├── components/   # 组件
│   ├── pages/        # 页面
│   ├── stores/       # 状态管理
│   ├── lib/          # 工具库
│   └── router/       # 路由配置
└── package.json

miniprogram/           # 用户端前端
├── src/
│   ├── components/   # 组件
│   ├── pages/        # 页面
│   ├── stores/       # 状态管理
│   ├── services/     # API 服务
│   └── router/       # 路由配置
└── package.json
```

## API 文档

### 认证接口
- `POST /api/auth/login` - 用户登录
- `POST /api/auth/register` - 用户注册
- `POST /api/auth/logout` - 用户登出

### 题目接口
- `GET /api/questions` - 获取题目列表
- `GET /api/questions/:id` - 获取题目详情
- `POST /api/questions` - 创建题目
- `PUT /api/questions/:id` - 更新题目
- `DELETE /api/questions/:id` - 删除题目

### 分类接口
- `GET /api/categories` - 获取分类列表
- `POST /api/categories` - 创建分类
- `PUT /api/categories/:id` - 更新分类
- `DELETE /api/categories/:id` - 删除分类

### 答题接口
- `POST /api/answers` - 提交答案
- `GET /api/answers/records` - 获取答题记录
- `GET /api/answers/statistics` - 获取答题统计

## 部署指南

详细的部署指南请参考 [DEPLOYMENT.md](./DEPLOYMENT.md)

### 生产环境部署要点

1. **安全配置**
   - 修改默认密码
   - 配置 HTTPS
   - 设置防火墙规则

2. **性能优化**
   - 调整数据库配置
   - 配置缓存
   - 启用 Gzip 压缩

3. **监控和日志**
   - 配置日志收集
   - 设置监控告警
   - 定期备份数据

## 测试

```bash
# 运行前端测试
cd admin-web && npm run test
cd miniprogram && npm run test

# 运行后端测试
cd server && go test ./...

# 运行 E2E 测试
npm run test:e2e
```

## 贡献指南

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情

## 更新日志

### v1.0.0 (2024-01-XX)
- ✨ 初始版本发布
- ✨ 完整的题目管理功能
- ✨ 用户答题和错题本功能
- ✨ 管理端和用户端界面
- ✨ 容器化部署支持

## 技术支持

- 📧 邮箱: admin@shuashuati.com
- 🐛 问题反馈: [GitHub Issues](https://github.com/your-repo/issues)
- 📖 文档: [项目文档](https://docs.shuashuati.com)

---

**注意**: 本项目仅供学习和研究使用，请勿用于商业用途。
