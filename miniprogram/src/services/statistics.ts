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
    return apiClient.get<OverviewStatistics>('/statistics/overview')
  }

  /**
   * 获取用户个人统计数据（需要认证）
   */
  static async getUserStatistics(): Promise<ApiResponse<UserStatistics>> {
    return apiClient.get<UserStatistics>('/answers/statistics')
  }
}

export default StatisticsService