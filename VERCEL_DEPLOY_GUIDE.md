# Vercel 部署指南

## 项目部署状态

✅ **配置文件已修复**
- 已修复 `vercel.json` 配置文件中的 functions 和 builds 属性冲突
- 移除了冲突的 functions 配置
- 优化了路由配置

✅ **构建测试通过**
- admin-web 项目构建成功
- 依赖包已安装完成
- dist 目录正确生成

## 手动部署步骤

### 1. 登录 Vercel
```bash
vercel login
```
选择您的登录方式（GitHub、Google等），在浏览器中完成登录。

### 2. 部署到生产环境
```bash
vercel --prod
```

### 3. 或者使用 Vercel 网页界面部署
1. 访问 [vercel.com](https://vercel.com)
2. 登录您的账户
3. 点击 "New Project"
4. 导入您的 Git 仓库
5. Vercel 会自动检测配置并部署

## 配置说明

当前 `vercel.json` 配置：
```json
{
  "version": 2,
  "name": "shuashuati-admin",
  "builds": [
    {
      "src": "admin-web/package.json",
      "use": "@vercel/static-build",
      "config": {
        "distDir": "admin-web/dist"
      }
    }
  ],
  "routes": [
    {
      "handle": "filesystem"
    },
    {
      "src": "/(.*)",
      "dest": "/admin-web/dist/index.html"
    }
  ],
  "regions": ["hkg1"]
}
```

## 环境变量配置

部署后，请在 Vercel 项目设置中添加以下环境变量：

```
VUE_APP_API_BASE_URL=https://your-api-domain.com
VUE_APP_SUPABASE_URL=https://lvafyknbbxmcatbmzeij.supabase.co
VUE_APP_SUPABASE_ANON_KEY=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

## 后端部署

后端服务已配置为使用 Supabase PostgreSQL，可以部署到：
- Vercel Functions
- Railway
- Render
- 或其他支持 Go 的平台

## 数据库

✅ Supabase 数据库已配置完成
- 数据库表结构已创建
- 后端代码已适配 UUID 主键
- 连接配置已更新

## 部署完成后

1. 访问部署的 URL 测试管理后台
2. 确认 API 接口正常工作
3. 测试小程序端连接

---

**注意**: 如果遇到部署问题，请检查：
1. Vercel 账户权限
2. Git 仓库访问权限
3. 构建日志中的错误信息