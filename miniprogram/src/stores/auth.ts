import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authService, type User, type UserStats } from '@/services/auth'
import { apiClient } from '@/services/api'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const token = ref<string | null>(null)
  const userStats = ref<UserStats | null>(null)
  const isLoading = ref(false)

  // 计算属性
  const isLoggedIn = computed(() => !!user.value)
  const isGuest = computed(() => user.value?.isGuest || false)

  // 设置用户信息
  const setUser = (userData: User) => {
    user.value = userData
    localStorage.setItem('user_info', JSON.stringify(userData))
  }

  // 设置token
  const setToken = (tokenValue: string) => {
    token.value = tokenValue
    localStorage.setItem('auth_token', tokenValue)
    // 同时设置API客户端的token
    apiClient.setToken(tokenValue)
  }

  // 设置用户统计
  const setUserStats = (stats: UserStats) => {
    userStats.value = stats
  }

  // 密码登录
  const loginWithPassword = async (username: string, password: string) => {
    try {
      isLoading.value = true
      const response = await authService.loginWithPassword(username, password)
      
      if (response.success) {
        setUser(response.data.user)
        setToken(response.data.token)
        
        // 获取用户统计信息
        await fetchUserStats()
        
        return response.data
      } else {
        throw new Error(response.message)
      }
    } catch (error) {
      console.error('密码登录失败:', error)
      throw error
    } finally {
      isLoading.value = false
    }
  }

  // 微信登录
  const loginWithWechat = async (code: string) => {
    try {
      isLoading.value = true
      const response = await authService.loginWithWechat(code)
      
      if (response.success) {
        setUser(response.data.user)
        setToken(response.data.token)
        
        // 获取用户统计信息
        await fetchUserStats()
        
        return response.data
      } else {
        throw new Error(response.message)
      }
    } catch (error) {
      console.error('微信登录失败:', error)
      throw error
    } finally {
      isLoading.value = false
    }
  }

  // 游客登录
  const loginAsGuest = async () => {
    try {
      isLoading.value = true
      const response = await authService.loginAsGuest()
      
      if (response.success) {
        setUser(response.data.user)
        setToken(response.data.token)
        
        // 游客不需要获取统计信息，但需要保存refresh token（如果有）
        if (response.data.refreshToken) {
          localStorage.setItem('refresh_token', response.data.refreshToken)
        }
        
        return response.data
      } else {
        throw new Error(response.message)
      }
    } catch (error) {
      console.error('游客登录失败:', error)
      throw error
    } finally {
      isLoading.value = false
    }
  }

  // 获取当前用户信息
  const fetchCurrentUser = async () => {
    try {
      const response = await authService.getCurrentUser()
      if (response.success) {
        setUser(response.data)
        return response.data
      }
    } catch (error) {
      console.error('获取用户信息失败:', error)
      // 如果是游客用户，不执行登出操作
      if (!isGuest.value) {
        // 如果获取用户信息失败，可能token已过期，执行登出
        logout()
      }
    }
  }

  // 获取用户统计信息
  const fetchUserStats = async () => {
    try {
      if (isGuest.value) return // 游客不获取统计信息
      
      const response = await authService.getUserStats()
      if (response.success) {
        setUserStats(response.data)
        return response.data
      }
    } catch (error) {
      console.error('获取用户统计失败:', error)
    }
  }

  // 更新用户信息
  const updateProfile = async (data: Partial<User>) => {
    try {
      isLoading.value = true
      const response = await authService.updateProfile(data)
      
      if (response.success) {
        setUser(response.data)
        return response.data
      } else {
        throw new Error(response.message)
      }
    } catch (error) {
      console.error('更新用户信息失败:', error)
      throw error
    } finally {
      isLoading.value = false
    }
  }

  // 退出登录
  const logout = async () => {
    try {
      // 调用后端登出接口
      if (token.value && !isGuest.value) {
        await authService.logout()
      }
    } catch (error) {
      console.error('登出请求失败:', error)
    } finally {
      // 无论后端请求是否成功，都清除本地状态
      user.value = null
      token.value = null
      userStats.value = null
      localStorage.removeItem('auth_token')
      localStorage.removeItem('user_info')
      apiClient.setToken(null)
    }
  }

  // 初始化认证状态
  const initAuth = async () => {
    const savedToken = localStorage.getItem('auth_token')
    const savedUser = localStorage.getItem('user_info')
    
    if (savedToken && savedUser) {
      try {
        token.value = savedToken
        user.value = JSON.parse(savedUser)
        apiClient.setToken(savedToken)
        
        // 验证token是否有效
        await fetchCurrentUser()
        
        // 获取最新的用户统计
        if (!isGuest.value) {
          await fetchUserStats()
        }
      } catch (error) {
        console.error('初始化认证状态失败:', error)
        // 如果验证失败，清除无效的认证信息并执行游客登录
        logout()
        await loginAsGuest()
      }
    } else {
      // 如果没有保存的认证信息，自动执行游客登录
      try {
        await loginAsGuest()
      } catch (error) {
        console.error('自动游客登录失败:', error)
      }
    }
  }

  // 刷新token
  const refreshToken = async () => {
    try {
      const refreshTokenValue = localStorage.getItem('refresh_token')
      if (!refreshTokenValue) {
        throw new Error('没有刷新token')
      }
      
      const response = await authService.refreshToken(refreshTokenValue)
      if (response.success) {
        setToken(response.data.token)
        return response.data.token
      } else {
        throw new Error(response.message)
      }
    } catch (error) {
      console.error('刷新token失败:', error)
      logout()
      throw error
    }
  }

  return {
    user,
    token,
    userStats,
    isLoading,
    isLoggedIn,
    isGuest,
    setUser,
    setToken,
    setUserStats,
    loginWithPassword,
    loginWithWechat,
    loginAsGuest,
    fetchCurrentUser,
    fetchUserStats,
    updateProfile,
    logout,
    initAuth,
    refreshToken
  }
})