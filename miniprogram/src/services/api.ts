// API 基础配置
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1'

// API 响应接口
export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
  success: boolean
}

// 分页响应接口 - 适配后端返回的数据结构
export interface PaginatedResponse<T> {
  data: T[]  // 后端返回的是data字段，不是items
  total: number
  page: number
  size: number  // 后端返回的是size字段，不是pageSize
}

// 前端使用的分页响应接口
export interface FrontendPaginatedResponse<T> {
  items: T[]
  total: number
  page: number
  pageSize: number
  totalPages: number
}

// HTTP 请求方法
class ApiClient {
  private baseURL: string
  private token: string | null = null

  constructor(baseURL: string) {
    this.baseURL = baseURL
    this.loadToken()
  }

  // 加载存储的 token
  private loadToken() {
    this.token = localStorage.getItem('auth_token')
  }

  // 设置 token
  setToken(token: string | null) {
    this.token = token
    if (token) {
      localStorage.setItem('auth_token', token)
    } else {
      localStorage.removeItem('auth_token')
    }
  }

  // 获取请求头
  private getHeaders(): HeadersInit {
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
    }

    if (this.token) {
      headers['Authorization'] = `Bearer ${this.token}`
    }

    return headers
  }

  // 处理响应
  private async handleResponse<T>(response: Response): Promise<ApiResponse<T>> {
    if (!response.ok) {
      // 处理401未授权错误
      if (response.status === 401) {
        this.setToken(null)
        // 动态导入以避免循环依赖
        const { useAuthStore } = await import('../stores/auth')
        const authStore = useAuthStore()
        await authStore.logout()
        
        // 如果当前不在登录页面，则跳转到登录页面
        if (window.location.pathname !== '/auth/login') {
          window.location.href = '/auth/login'
        }
        
        throw new Error('登录已过期，请重新登录')
      }
      
      const errorData = await response.json().catch(() => ({ message: '网络错误' }))
      throw new Error(errorData.message || `HTTP ${response.status}`)
    }

    const data = await response.json()
    
    // 如果后端返回的数据已经是ApiResponse格式，直接返回
    if (data && typeof data === 'object' && 'success' in data && 'code' in data && 'message' in data) {
      return data
    }
    
    // 否则包装成ApiResponse格式
    return {
      code: response.status,
      message: 'success',
      data: data,
      success: true
    }
  }

  // 转换分页响应格式
  private transformPaginatedResponse<T>(backendResponse: ApiResponse<PaginatedResponse<T>>): ApiResponse<FrontendPaginatedResponse<T>> {
    const { data, total, page, size } = backendResponse.data
    return {
      ...backendResponse,
      data: {
        items: data,
        total,
        page,
        pageSize: size,
        totalPages: Math.ceil(total / size)
      }
    }
  }

  // GET 请求
  async get<T>(endpoint: string, params?: Record<string, any>): Promise<ApiResponse<T>> {
    const url = new URL(`${this.baseURL}${endpoint}`)
    
    if (params) {
      Object.keys(params).forEach(key => {
        if (params[key] !== undefined && params[key] !== null) {
          url.searchParams.append(key, String(params[key]))
        }
      })
    }

    const response = await fetch(url.toString(), {
      method: 'GET',
      headers: this.getHeaders(),
    })

    return this.handleResponse<T>(response)
  }

  // 分页GET请求 - 自动转换数据格式
  async getPaginated<T>(endpoint: string, params?: Record<string, any>): Promise<ApiResponse<FrontendPaginatedResponse<T>>> {
    const backendResponse = await this.get<PaginatedResponse<T>>(endpoint, params)
    return this.transformPaginatedResponse(backendResponse)
  }

  // POST 请求
  async post<T>(endpoint: string, data?: any): Promise<ApiResponse<T>> {
    const response = await fetch(`${this.baseURL}${endpoint}`, {
      method: 'POST',
      headers: this.getHeaders(),
      body: data ? JSON.stringify(data) : undefined,
    })

    return this.handleResponse<T>(response)
  }

  // PUT 请求
  async put<T>(endpoint: string, data?: any): Promise<ApiResponse<T>> {
    const response = await fetch(`${this.baseURL}${endpoint}`, {
      method: 'PUT',
      headers: this.getHeaders(),
      body: data ? JSON.stringify(data) : undefined,
    })

    return this.handleResponse<T>(response)
  }

  // DELETE 请求
  async delete<T>(endpoint: string): Promise<ApiResponse<T>> {
    const response = await fetch(`${this.baseURL}${endpoint}`, {
      method: 'DELETE',
      headers: this.getHeaders(),
    })

    return this.handleResponse<T>(response)
  }

  // 上传文件
  async upload<T>(endpoint: string, file: File, additionalData?: Record<string, any>): Promise<ApiResponse<T>> {
    const formData = new FormData()
    formData.append('file', file)
    
    if (additionalData) {
      Object.keys(additionalData).forEach(key => {
        formData.append(key, String(additionalData[key]))
      })
    }

    const headers: HeadersInit = {}
    if (this.token) {
      headers['Authorization'] = `Bearer ${this.token}`
    }

    const response = await fetch(`${this.baseURL}${endpoint}`, {
      method: 'POST',
      headers,
      body: formData,
    })

    return this.handleResponse<T>(response)
  }
}

// 创建 API 客户端实例
export const apiClient = new ApiClient(API_BASE_URL)

// 导出常用方法
export const { get, post, put, delete: del, upload, setToken } = apiClient

// 错误处理工具
export class ApiError extends Error {
  constructor(
    message: string,
    public code?: number,
    public response?: any
  ) {
    super(message)
    this.name = 'ApiError'
  }
}

// 请求拦截器（模拟）
export const requestInterceptor = {
  onRequest: (config: any) => {
    console.log('API Request:', config)
    return config
  },
  onError: (error: any) => {
    console.error('API Request Error:', error)
    throw error
  }
}

// 响应拦截器（模拟）
export const responseInterceptor = {
  onResponse: (response: any) => {
    console.log('API Response:', response)
    return response
  },
  onError: (error: any) => {
    console.error('API Response Error:', error)
    throw error
  }
}

// 导出默认实例
export default apiClient