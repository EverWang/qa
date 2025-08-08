import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { setupRouterGuards } from './guards'
import LoginPage from '@/pages/LoginPage.vue'
import AdminLayout from '@/components/Layout/AdminLayout.vue'

// 定义路由配置
const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'login',
    component: LoginPage,
    meta: {
      title: '登录',
      requiresAuth: false
    }
  },
  {
    path: '/',
    component: AdminLayout,
    redirect: '/dashboard',
    meta: {
      requiresAuth: true
    },
    children: [
      {
        path: 'dashboard',
        name: 'dashboard',
        component: () => import('@/pages/DashboardPage.vue'),
        meta: {
          title: '仪表板',
          requiresAuth: true
        }
      },
      {
        path: 'questions',
        name: 'QuestionList',
        component: () => import('@/pages/Questions/QuestionList.vue'),
        meta: {
          title: '题目管理',
          requiresAuth: true
        }
      },
      {
        path: 'questions/create',
        name: 'QuestionCreate',
        component: () => import('@/pages/Questions/QuestionCreate.vue'),
        meta: {
          title: '创建题目',
          requiresAuth: true
        }
      },
      {
        path: 'questions/import',
        name: 'QuestionImport',
        component: () => import('@/pages/Questions/QuestionImport.vue'),
        meta: {
          title: '批量导入',
          requiresAuth: true
        }
      },
      {
        path: 'questions/:id',
        name: 'QuestionDetail',
        component: () => import('@/pages/Questions/QuestionDetail.vue'),
        meta: {
          title: '题目详情',
          requiresAuth: true
        }
      },
      {
        path: 'questions/:id/edit',
        name: 'QuestionEdit',
        component: () => import('@/pages/Questions/QuestionEdit.vue'),
        meta: {
          title: '编辑题目',
          requiresAuth: true
        }
      },
      {
        path: 'categories',
        name: 'categories',
        component: () => import('@/pages/Categories/CategoryList.vue'),
        meta: {
          title: '分类管理',
          requiresAuth: true
        }
      },
      {
        path: 'users',
        name: 'UserManagement',
        component: () => import('@/pages/Users/UserList.vue'),
        meta: { title: '用户管理', requiresAuth: true, roles: ['admin'] }
      },
      {
        path: 'settings',
        name: 'SystemSettings',
        component: () => import('@/pages/System/SystemSettings.vue'),
        meta: { title: '系统设置', requiresAuth: true, roles: ['admin'] }
      },
      {
        path: 'operation-logs',
        name: 'OperationLogs',
        component: () => import('@/pages/System/OperationLogs.vue'),
        meta: { title: '操作日志', requiresAuth: true, roles: ['admin'] }
      }
    ]
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'not-found',
    component: () => import('@/pages/NotFound.vue'),
    meta: {
      title: '页面未找到'
    }
  }
]

// 创建路由实例
const router = createRouter({
  history: createWebHistory(),
  routes,
})

// 设置路由守卫
setupRouterGuards(router)

export default router
