import { createRouter, createWebHistory } from 'vue-router'
import HomePage from '@/pages/HomePage.vue'
import LoginPage from '@/pages/auth/LoginPage.vue'
import QuestionDetail from '@/pages/question/QuestionDetail.vue'
import UserCenter from '@/pages/user/UserCenter.vue'
import UserRecords from '@/pages/user/UserRecords.vue'
import MistakeBook from '@/pages/user/MistakeBook.vue'
import CategoryList from '@/pages/category/CategoryList.vue'
import SearchPage from '@/pages/search/SearchPage.vue'
import SettingsPage from '@/pages/user/SettingsPage.vue'
import ApiTest from '@/pages/test/ApiTest.vue'
import CategoryTest from '@/pages/test/CategoryTest.vue'

// 定义路由配置
const routes = [
  {
    path: '/',
    name: 'home',
    component: HomePage,
    meta: { title: '首页' }
  },
  {
    path: '/auth/login',
    name: 'login',
    component: LoginPage,
    meta: { title: '登录授权' }
  },
  {
    path: '/question/:id',
    name: 'questionDetail',
    component: QuestionDetail,
    meta: { title: '题目详情' }
  },
  {
    path: '/user/center',
    name: 'userCenter',
    component: UserCenter,
    meta: { title: '个人中心', requiresAuth: true }
  },
  {
    path: '/user/records',
    name: 'userRecords',
    component: UserRecords,
    meta: { title: '做题记录', requiresAuth: true }
  },
  {
    path: '/user/mistakes',
    name: 'mistakeBook',
    component: MistakeBook,
    meta: { title: '错题本', requiresAuth: true }
  },
  {
    path: '/category/:id?',
    name: 'categoryList',
    component: CategoryList,
    meta: { title: '分类题目' }
  },
  {
    path: '/search',
    name: 'search',
    component: SearchPage,
    meta: { title: '搜索题目' }
  },
  {
    path: '/user/settings',
    name: 'settings',
    component: SettingsPage,
    meta: { title: '设置', requiresAuth: true }
  },
  {
    path: '/test/api',
    name: 'apiTest',
    component: ApiTest,
    meta: { title: 'API测试' }
  },
  {
    path: '/test/debug',
    name: 'debugApi',
    component: () => import('@/pages/test/DebugApi.vue'),
    meta: { title: 'API调试' }
  },
  {
    path: '/test/category',
    name: 'categoryTest',
    component: CategoryTest,
    meta: { title: '分类测试' }
  }
]

// 创建路由实例
const router = createRouter({
  history: createWebHistory(),
  routes,
})

// 路由守卫
router.beforeEach((to, from, next) => {
  // 设置页面标题
  if (to.meta.title) {
    document.title = `${to.meta.title} - 刷刷题`
  }
  
  // 检查是否需要登录
  if (to.meta.requiresAuth) {
    const token = localStorage.getItem('auth_token')
    const userInfo = localStorage.getItem('user_info')
    
    if (!token || !userInfo) {
      next({ name: 'login', query: { redirect: to.fullPath } })
      return
    }
    
    // 检查用户信息是否有效
    try {
      const user = JSON.parse(userInfo)
      if (user && (user.id > 0 || user.isGuest)) {
        // 正式用户和游客用户都可以访问需要认证的页面
        next()
        return
      }
    } catch (error) {
      console.error('解析用户信息失败:', error)
      next({ name: 'login', query: { redirect: to.fullPath } })
      return
    }
    
    next({ name: 'login', query: { redirect: to.fullPath } })
    return
  }
  
  next()
})

export default router
