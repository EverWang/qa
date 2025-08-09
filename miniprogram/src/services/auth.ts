import { apiClient } from './api'
import type { ApiResponse } from './api'

// 用户信息接口
export interface User {
  id: number
  username: string
  nickname: string
  avatar?: string
  email?: string
  phone?: string
  isGuest: boolean
  createdAt: string
  updatedAt: string
}

// 登录请求接口
export interface LoginRequest {
  username?: string
  password?: string
  code?: string // 微信授权码
  type: 'password' | 'wechat' | 'guest'
}

// 登录响应接口
export interface LoginResponse {
  user: User
  token: string
  refreshToken: string
  expiresIn: number
}

// 注册请求接口
export interface RegisterRequest {
  username: string
  password: string
  nickname: string
  email?: string
  phone?: string
}

// 用户统计信息接口
// 用户统计接口（与后端 AnswerStatistics 结构匹配）
export interface UserStats {
  totalAnswered: number
  correctAnswered: number
  accuracyRate: number
  totalTimeSpent: number
  averageTime: number
  todayAnswered: number
  weekAnswered: number
  monthAnswered: number
}

// 用户认证服务
export class AuthService {
  // 密码登录
  static async loginWithPassword(username: string, password: string): Promise<ApiResponse<LoginResponse>> {
    return apiClient.post<LoginResponse>('/auth/login', {
      username,
      password,
      type: 'password'
    })
  }

  // 微信登录
  static async loginWithWechat(code: string): Promise<ApiResponse<LoginResponse>> {
    return apiClient.post<LoginResponse>('/auth/login', {
      code,
      type: 'wechat'
    })
  }



  // 注册
  static async register(data: RegisterRequest): Promise<ApiResponse<LoginResponse>> {
    return apiClient.post<LoginResponse>('/auth/register', data)
  }

  // 退出登录
  static async logout(): Promise<ApiResponse<void>> {
    return apiClient.post<void>('/auth/logout')
  }

  // 刷新 token
  static async refreshToken(refreshToken: string): Promise<ApiResponse<{ token: string; expiresIn: number }>> {
    return apiClient.post('/auth/refresh', { refreshToken })
  }

  // 获取当前用户信息
  static async getCurrentUser(): Promise<ApiResponse<User>> {
    return apiClient.get<User>('/auth/me')
  }

  // 更新用户信息
  static async updateProfile(data: Partial<User>): Promise<ApiResponse<User>> {
    return apiClient.put<User>('/auth/profile', data)
  }

  // 修改密码
  static async changePassword(oldPassword: string, newPassword: string): Promise<ApiResponse<void>> {
    return apiClient.put<void>('/auth/password', {
      oldPassword,
      newPassword
    })
  }

  // 绑定微信
  static async bindWechat(code: string): Promise<ApiResponse<void>> {
    return apiClient.post<void>('/auth/bind-wechat', { code })
  }

  // 解绑微信
  static async unbindWechat(): Promise<ApiResponse<void>> {
    return apiClient.delete<void>('/auth/bind-wechat')
  }

  // 获取用户统计信息
  static async getUserStats(): Promise<ApiResponse<UserStats>> {
    return apiClient.get<UserStats>('/answers/statistics')
  }

  // 上传头像
  static async uploadAvatar(file: File): Promise<ApiResponse<{ avatar: string }>> {
    return apiClient.upload<{ avatar: string }>('/auth/avatar', file)
  }

  // 发送验证码
  static async sendVerificationCode(phone: string): Promise<ApiResponse<void>> {
    return apiClient.post<void>('/auth/send-code', { phone })
  }

  // 验证手机号
  static async verifyPhone(phone: string, code: string): Promise<ApiResponse<void>> {
    return apiClient.post<void>('/auth/verify-phone', { phone, code })
  }

  // 重置密码
  static async resetPassword(phone: string, code: string, newPassword: string): Promise<ApiResponse<void>> {
    return apiClient.post<void>('/auth/reset-password', {
      phone,
      code,
      newPassword
    })
  }

  // 检查用户名是否可用
  static async checkUsername(username: string): Promise<ApiResponse<{ available: boolean }>> {
    return apiClient.get<{ available: boolean }>('/auth/check-username', { username })
  }

  // 检查邮箱是否可用
  static async checkEmail(email: string): Promise<ApiResponse<{ available: boolean }>> {
    return apiClient.get<{ available: boolean }>('/auth/check-email', { email })
  }

  // 删除账户
  static async deleteAccount(password: string): Promise<ApiResponse<void>> {
    return apiClient.delete<void>('/auth/account')
  }
}

// 模拟数据（开发阶段使用）
export class MockAuthService {
  private static mockUser: User = {
    id: 1,
    username: 'testuser',
    nickname: '测试用户',
    avatar: 'https://trae-api-us.mchost.guru/api/ide/v1/text_to_image?prompt=cute%20cartoon%20avatar%20of%20a%20student%20with%20glasses%20and%20books&image_size=square',
    email: 'test@example.com',
    phone: '13800138000',
    isGuest: false,
    createdAt: '2024-01-01T00:00:00Z',
    updatedAt: '2024-01-01T00:00:00Z'
  }

  private static mockStats: UserStats = {
    totalAnswered: 1250,
    correctAnswered: 980,
    accuracyRate: 78.4,
    totalTimeSpent: 140400, // 39小时转换为秒
    averageTime: 112.32, // 平均每题用时（秒）
    todayAnswered: 15,
    weekAnswered: 85,
    monthAnswered: 320
  }

  static async loginWithPassword(username: string, password: string): Promise<ApiResponse<LoginResponse>> {
    // 模拟网络延迟
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    if (username === 'admin' && password === 'admin') {
      return {
        code: 200,
        message: '登录成功',
        success: true,
        data: {
          user: this.mockUser,
          token: 'mock_token_' + Date.now(),
          refreshToken: 'mock_refresh_token_' + Date.now(),
          expiresIn: 7200
        }
      }
    } else {
      throw new Error('用户名或密码错误')
    }
  }

  static async loginWithWechat(code: string): Promise<ApiResponse<LoginResponse>> {
    await new Promise(resolve => setTimeout(resolve, 1500))
    
    return {
      code: 200,
      message: '微信登录成功',
      success: true,
      data: {
        user: { ...this.mockUser, nickname: '微信用户' },
        token: 'mock_wechat_token_' + Date.now(),
        refreshToken: 'mock_wechat_refresh_token_' + Date.now(),
        expiresIn: 7200
      }
    }
  }

  static async loginAsGuest(): Promise<ApiResponse<LoginResponse>> {
    await new Promise(resolve => setTimeout(resolve, 500))
    
    return {
      code: 200,
      message: '游客登录成功',
      success: true,
      data: {
        user: {
          ...this.mockUser,
          id: 0,
          username: 'guest',
          nickname: '游客用户',
          isGuest: true
        },
        token: 'mock_guest_token_' + Date.now(),
        refreshToken: 'mock_guest_refresh_token_' + Date.now(),
        expiresIn: 3600
      }
    }
  }

  static async getCurrentUser(): Promise<ApiResponse<User>> {
    await new Promise(resolve => setTimeout(resolve, 300))
    
    return {
      code: 200,
      message: '获取用户信息成功',
      success: true,
      data: this.mockUser
    }
  }

  static async getUserStats(): Promise<ApiResponse<UserStats>> {
    await new Promise(resolve => setTimeout(resolve, 500))
    
    return {
      code: 200,
      message: '获取用户统计成功',
      success: true,
      data: this.mockStats
    }
  }

  static async updateProfile(data: Partial<User>): Promise<ApiResponse<User>> {
    await new Promise(resolve => setTimeout(resolve, 500))
    
    // 更新模拟用户数据
    this.mockUser = { ...this.mockUser, ...data }
    
    return {
      code: 200,
      message: '更新用户信息成功',
      success: true,
      data: this.mockUser
    }
  }

  static async refreshToken(refreshToken: string): Promise<ApiResponse<{ token: string; expiresIn: number }>> {
    await new Promise(resolve => setTimeout(resolve, 300))
    
    return {
      code: 200,
      message: '刷新token成功',
      success: true,
      data: {
        token: 'mock_refreshed_token_' + Date.now(),
        expiresIn: 7200
      }
    }
  }

  static async logout(): Promise<ApiResponse<void>> {
    await new Promise(resolve => setTimeout(resolve, 300))
    
    return {
      code: 200,
      message: '退出登录成功',
      success: true,
      data: undefined
    }
  }
}

// 使用真实的认证服务
export const authService = AuthService

export default authService