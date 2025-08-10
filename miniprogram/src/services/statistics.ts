import { apiClient } from './api'
import type { ApiResponse } from './api'

// 概览统计接口
export interface OverviewStatistics {
  totalUsers: number
  totalQuestions: number
  totalCategories: number
  totalAnswers: number
  todayUsers: number
  todayAnswers: number
  weekUsers: number
  weekAnswers: number
  monthUsers: number
  monthAnswers: number
}

// 用户统计接口
export interface UserStatistics {
  totalAnswered: number
  correctAnswered: number
  accuracyRate: number
  totalTimeSpent: number
  todayAnswered: number
  weekAnswered: number
  monthAnswered: number
}

// 统计服务
export class StatisticsService {
  /**
   * 获取概览统计数据（公开接口，不需要认证）
   */
  static async getOverviewStatistics(): Promise<ApiResponse<OverviewStatistics>> {
    try {
      console.log('StatisticsService.getOverviewStatistics 开始调用')
      const response = await apiClient.get<OverviewStatistics>('/api/v1/statistics/overview')
      console.log('StatisticsService.getOverviewStatistics 响应:', response)
      return response
    } catch (error) {
      console.error('StatisticsService.getOverviewStatistics 错误:', error)
      throw error
    }
  }

  /**
   * 获取用户统计数据（需要认证）
   */
  static async getUserStatistics(): Promise<ApiResponse<UserStatistics>> {
    try {
      console.log('StatisticsService.getUserStatistics 开始调用')
      const response = await apiClient.get<UserStatistics>('/api/v1/answers/statistics')
      console.log('StatisticsService.getUserStatistics 响应:', response)
      return response
    } catch (error) {
      console.error('StatisticsService.getUserStatistics 错误:', error)
      throw error
    }
  }
}

export default StatisticsService