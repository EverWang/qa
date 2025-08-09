import type { Router, RouteRecordRaw } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useAuthStore } from '@/stores/auth'

// 不需要登录的页面
const whiteList = ['/login']

// 设置路由守卫
export function setupRouterGuards(router: Router) {
  // 全局前置守卫
  router.beforeEach(async (to, from, next) => {
    // 显示加载状态
    const authStore = useAuthStore()
    
    // 如果是白名单页面，直接放行
    if (whiteList.includes(to.path)) {
      // 如果已经登录，跳转到首页
      if (authStore.isAuthenticated) {
        next('/')
        return
      }
      next()
      return
    }
    
    // 检查是否已登录
    if (!authStore.isAuthenticated) {
      ElMessage.warning('请先登录')
      next('/login')
      return
    }
    
    // 如果有token但没有用户信息，尝试获取用户信息
    if (!authStore.user) {
      try {
        const success = await authStore.fetchUserInfo()
        if (!success) {
          next('/login')
          return
        }
      } catch (error) {
        console.error('获取用户信息失败:', error)
        next('/login')
        return
      }
    }
    
    // 检查用户权限
    if (!checkPermission(to.path, authStore.user?.role)) {
      ElMessage.error('权限不足，无法访问该页面')
      next('/')
      return
    }
    
    next()
  })
  
  // 全局后置守卫
  router.afterEach((to) => {
    // 设置页面标题
    const title = getPageTitle(to.path)
    document.title = title
  })
}

// 检查权限
function checkPermission(path: string, userRole?: string): boolean {
  // 如果没有用户角色信息，拒绝访问
  if (!userRole) {
    return false
  }
  
  // 管理员权限映射
  const adminRoutes = [
    '/',
    '/questions',
    '/questions/create',
    '/questions/import',
    '/categories',
    '/users',
    '/settings',
    '/operation-logs'
  ]
  
  // 根据用户角色检查权限
  switch (userRole) {
    case 'admin':
      // 管理员可以访问所有页面
      return adminRoutes.some(route => path.startsWith(route))
    case 'editor':
      // 编辑员只能访问题目和分类管理
      const editorRoutes = ['/', '/questions', '/categories']
      return editorRoutes.some(route => path.startsWith(route))
    default:
      return false
  }
}

// 获取页面标题
function getPageTitle(path: string): string {
  const titleMap: Record<string, string> = {
    '/': '仪表板 - 刷刷题管理后台',
    '/login': '登录 - 刷刷题管理后台',
    '/questions': '题目管理 - 刷刷题管理后台',
    '/questions/create': '添加题目 - 刷刷题管理后台',
    '/questions/import': '批量导入 - 刷刷题管理后台',
    '/categories': '分类管理 - 刷刷题管理后台',
    '/users': '用户管理 - 刷刷题管理后台',
    '/settings': '系统设置 - 刷刷题管理后台',
    '/operation-logs': '操作日志 - 刷刷题管理后台'
  }
  
  return titleMap[path] || '刷刷题管理后台'
}

// 动态添加路由（用于权限路由）
export function addDynamicRoutes(router: Router, userRole: string) {
  // 根据用户角色动态添加路由
  const routes = getRoutesByRole(userRole)
  
  routes.forEach(route => {
    router.addRoute(route)
  })
}

// 根据角色获取路由
function getRoutesByRole(role: string): RouteRecordRaw[] {
  const baseRoutes: RouteRecordRaw[] = [
    {
      path: '/',
      name: 'dashboard',
      component: () => import('@/pages/DashboardPage.vue'),
      meta: { title: '仪表板', requiresAuth: true }
    }
  ]
  
  const adminRoutes: RouteRecordRaw[] = [
    {
      path: '/questions',
      name: 'questions',
      component: () => import('@/pages/Questions/QuestionList.vue'),
      meta: { title: '题目管理', requiresAuth: true }
    },
    {
      path: '/categories',
      name: 'categories',
      component: () => import('@/pages/Categories/CategoryList.vue'),
      meta: { title: '分类管理', requiresAuth: true }
    },
    {
      path: '/users',
      name: 'users',
      component: () => import('@/pages/Users/UserList.vue'),
      meta: { title: '用户管理', requiresAuth: true }
    },
    {
      path: '/system',
      name: 'system',
      component: () => import('@/pages/System/SystemSettings.vue'),
      meta: { title: '系统管理', requiresAuth: true }
    }
  ]
  
  switch (role) {
    case 'admin':
      return [...baseRoutes, ...adminRoutes]
    case 'editor':
      return [...baseRoutes, ...adminRoutes.slice(0, 2)] // 只包含题目和分类管理
    default:
      return baseRoutes
  }
}