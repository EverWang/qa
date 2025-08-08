import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '@/lib/axios'

// 用户信息接口
export interface UserInfo {
  id: number
  username: string
  email: string
  role: string
  avatar?: string
  created_at: string
}

// 登录请求接口
export interface LoginRequest {
  username: string
  password: string
}

// 登录响应接口
export interface LoginResponse {
  token: string
  user: UserInfo
}

export const useAuthStore = defineStore('auth', () => {
  // 状态
  const token = ref<string | null>(localStorage.getItem('admin_token'))
  const user = ref<UserInfo | null>(null)
  const loading = ref(false)

  // 计算属性
  const isAuthenticated = computed(() => !!token.value)
  const isAdmin = computed(() => user.value?.role === 'admin')

  // 设置api默认请求头
  if (token.value) {
    api.defaults.headers.common['Authorization'] = `Bearer ${token.value}`
  }

  // 登录方法
  const login = async (credentials: LoginRequest): Promise<boolean> => {
    try {
      loading.value = true
      const response = await api.post('/api/v1/admin/login', credentials)
      
      // 后端返回格式: {code: 200, message: 'success', data: {token, user}}
      const { token: newToken, user: userInfo } = response.data.data
      
      // 保存token和用户信息
      token.value = newToken
      user.value = userInfo
      
      // 保存到localStorage
      localStorage.setItem('admin_token', newToken)
      localStorage.setItem('admin_user', JSON.stringify(userInfo))
      
      // 设置api默认请求头
      api.defaults.headers.common['Authorization'] = `Bearer ${newToken}`
      
      // 确保用户信息完整，如果需要可以再次获取
      if (!userInfo.role) {
        await fetchUserInfo()
      }
      
      return true
    } catch (error) {
      console.error('登录失败:', error)
      return false
    } finally {
      loading.value = false
    }
  }

  // 登出方法
  const logout = async () => {
    try {
      // 调用后端登出接口
      await api.post('/api/v1/admin/logout')
    } catch (error) {
      console.error('登出请求失败:', error)
    } finally {
      // 清除本地数据
      token.value = null
      user.value = null
      localStorage.removeItem('admin_token')
      localStorage.removeItem('admin_user')
      delete api.defaults.headers.common['Authorization']
    }
  }

  // 获取用户信息
  const fetchUserInfo = async (): Promise<boolean> => {
    try {
      const response = await api.get('/api/v1/admin/profile')
      // 后端返回格式: {code: 200, message: 'success', data: userInfo}
      user.value = response.data.data
      localStorage.setItem('admin_user', JSON.stringify(response.data.data))
      return true
    } catch (error) {
      console.error('获取用户信息失败:', error)
      // 如果获取用户信息失败，可能token已过期，执行登出
      await logout()
      return false
    }
  }

  // 初始化用户信息（从localStorage恢复）
  const initializeAuth = () => {
    const savedUser = localStorage.getItem('admin_user')
    if (savedUser && token.value) {
      try {
        user.value = JSON.parse(savedUser)
      } catch (error) {
        console.error('解析用户信息失败:', error)
        logout()
      }
    }
  }

  // 检查token有效性
  const checkTokenValidity = async (): Promise<boolean> => {
    if (!token.value) return false
    
    try {
      await api.get('/api/v1/admin/check-token')
      return true
    } catch (error) {
      console.error('Token验证失败:', error)
      await logout()
      return false
    }
  }

  return {
    // 状态
    token,
    user,
    loading,
    // 计算属性
    isAuthenticated,
    isAdmin,
    // 方法
    login,
    logout,
    fetchUserInfo,
    initializeAuth,
    checkTokenValidity
  }
})