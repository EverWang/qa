import { apiClient } from './api'
import type { ApiResponse, FrontendPaginatedResponse } from './api'

// 题目难度枚举
export enum QuestionDifficulty {
  EASY = 'easy',
  MEDIUM = 'medium',
  HARD = 'hard'
}

// 题目类型枚举
export enum QuestionType {
  SINGLE_CHOICE = 'single_choice',
  MULTIPLE_CHOICE = 'multiple_choice',
  TRUE_FALSE = 'true_false',
  FILL_BLANK = 'fill_blank',
  ESSAY = 'essay'
}

// 选项接口
export interface QuestionOption {
  id: string
  content: string
  isCorrect: boolean
}

// 题目接口
export interface Question {
  id: number
  title: string
  content: string
  type: QuestionType
  difficulty: QuestionDifficulty
  options: QuestionOption[]
  correctAnswer: string | string[]
  explanation: string
  categoryId: number
  categoryName: string
  category?: Category
  tags: string[]
  createdAt: string
  updatedAt: string
  // 用户相关字段
  isAnswered?: boolean
  userAnswer?: string | string[]
  isCorrect?: boolean
  answeredAt?: string
}

// 分类接口
export interface Category {
  id: number
  name: string
  description: string
  parentId?: number
  level: number
  questionCount: number
  icon?: string
  color?: string
  createdAt: string
  updatedAt: string
}

// 答题记录接口
export interface AnswerRecord {
  id: number
  questionId: number
  userId: number
  userAnswer: number
  isCorrect: boolean
  timeSpent: number // 秒
  answeredAt: string
  question?: Question
}

// 错题本接口
export interface MistakeBook {
  id: number
  questionId: number
  userId: number
  addedAt: string
  reviewCount: number
  lastReviewAt?: string
  isMastered: boolean
  question?: Question
}

// 搜索参数接口
export interface QuestionSearchParams {
  keyword?: string
  categoryId?: number
  difficulty?: QuestionDifficulty
  type?: QuestionType
  tags?: string[]
  page?: number
  pageSize?: number
  sortBy?: 'createdAt' | 'difficulty' | 'random'
  sortOrder?: 'asc' | 'desc'
}

// 答题请求接口
export interface SubmitAnswerRequest {
  questionId: number
  userAnswer: number // 后端期望数字索引
  timeSpent: number
}

// 题目服务
export class QuestionService {
  // 获取题目列表
  static async getQuestions(params?: QuestionSearchParams): Promise<ApiResponse<FrontendPaginatedResponse<Question>>> {
    try {
      console.log('QuestionService.getQuestions 调用参数:', params)
      
      // 转换前端参数名为后端期望的格式
      const backendParams: any = {}
      if (params) {
        if (params.categoryId !== undefined) backendParams.categoryId = params.categoryId
        if (params.difficulty !== undefined) backendParams.difficulty = params.difficulty
        if (params.keyword !== undefined) backendParams.keyword = params.keyword
        if (params.type !== undefined) backendParams.type = params.type
        if (params.page !== undefined) backendParams.page = params.page
        if (params.pageSize !== undefined) backendParams.size = params.pageSize
        if (params.sortBy !== undefined) backendParams.sortBy = params.sortBy
        if (params.sortOrder !== undefined) backendParams.sortOrder = params.sortOrder
      }
      console.log('QuestionService.getQuestions 转换后的参数:', backendParams)
      
      const response = await apiClient.getPaginated<any>('/api/v1/questions', backendParams)
      console.log('QuestionService.getQuestions 原始响应:', response)
      
      // 转换后端数据格式到前端期望格式
      if (response.success && response.data && response.data.items) {
        response.data.items = response.data.items.map((item: any) => {
          console.log('转换题目数据:', item)
          return {
            id: item.id,
            title: item.title,
            content: item.content,
            type: item.type,
            difficulty: item.difficulty,
            options: Array.isArray(item.options) ? item.options.map((opt: any, index: number) => ({
              id: String.fromCharCode(65 + index), // A, B, C, D
              content: opt,
              isCorrect: item.correct_answer === index
            })) : [],
            correctAnswer: item.correct_answer,
            explanation: item.explanation || '',
            categoryId: item.category_id,
            categoryName: item.category?.name || '未分类',
            tags: item.tags || [],
            createdAt: item.createdAt,
            updatedAt: item.updatedAt
          }
        })
        console.log('转换后的题目数据:', response.data.items)
      }
      
      return response
    } catch (error) {
      console.error('QuestionService.getQuestions 错误:', error)
      throw error
    }
  }

  // 获取题目详情
  static async getQuestion(id: number): Promise<ApiResponse<Question>> {
    try {
      console.log('QuestionService.getQuestion 调用参数:', id)
      const response = await apiClient.get<any>(`/api/v1/questions/${id}`)
      console.log('QuestionService.getQuestion 原始响应:', response)
      
      // 转换后端数据格式到前端期望格式
      if (response.success && response.data) {
        const item = response.data
        console.log('转换题目详情数据:', item)
        
        const transformedData = {
          id: item.id,
          title: item.title,
          content: item.content,
          type: item.type,
          difficulty: item.difficulty,
          options: Array.isArray(item.options) ? item.options : [],
          correctAnswer: item.correct_answer,
          explanation: item.explanation || '',
          categoryId: item.category_id,
          categoryName: item.category?.name || '未分类',
          tags: item.tags || [],
          createdAt: item.createdAt,
          updatedAt: item.updatedAt,
          category: item.category
        }
        console.log('转换后的题目详情数据:', transformedData)
        
        return {
          ...response,
          data: transformedData
        }
      }
      
      return response
    } catch (error) {
      console.error('QuestionService.getQuestion 错误:', error)
      throw error
    }
  }

  // 搜索题目
  static async searchQuestions(params: QuestionSearchParams): Promise<ApiResponse<FrontendPaginatedResponse<Question>>> {
    return apiClient.getPaginated<Question>('/api/v1/questions/search', params)
  }

  // 获取随机题目
  static async getRandomQuestions(count: number, categoryId?: number, difficulty?: QuestionDifficulty): Promise<ApiResponse<Question[]>> {
    return apiClient.get<Question[]>('/api/v1/questions/random', {
      count,
      categoryId,
      difficulty
    })
  }

  // 提交答案
  static async submitAnswer(data: SubmitAnswerRequest): Promise<ApiResponse<AnswerRecord>> {
    return apiClient.post<AnswerRecord>('/api/v1/answers', data)
  }

  // 获取题目统计
  static async getQuestionStats(questionId: number): Promise<ApiResponse<{
    totalAnswers: number
    correctAnswers: number
    correctRate: number
    averageTime: number
  }>> {
    return apiClient.get(`/api/v1/questions/${questionId}/stats`)
  }

  // 获取分类列表
  static async getCategories(parentId?: number): Promise<ApiResponse<Category[]>> {
    try {
      console.log('QuestionService.getCategories 调用参数:', { parentId })
      const response = await apiClient.get<any>('/api/v1/categories', { parentId })
      console.log('QuestionService.getCategories 原始响应:', response)
      
      // 转换后端数据格式到前端期望格式
      if (response.success && response.data) {
        console.log('response.data类型:', typeof response.data)
        console.log('response.data内容:', response.data)
        console.log('response.data是否为数组:', Array.isArray(response.data))
        
        // 如果response.data是对象且包含data字段，则使用data字段
        let actualData = response.data
        if (typeof response.data === 'object' && !Array.isArray(response.data) && response.data.data) {
          actualData = response.data.data
        }
        
        // 确保actualData是数组
        const dataArray = Array.isArray(actualData) ? actualData : []
        console.log('QuestionService.getCategories 数据数组:', dataArray)
        console.log('数据数组长度:', dataArray.length)
        console.log('第一个数据项:', dataArray[0])
        
        const transformedData = dataArray.map((item: any) => {
          console.log('转换分类数据:', item)
          return {
            id: item.id,
            name: item.name,
            description: item.description || '',
            parentId: item.parentId,
            level: item.level,
            status: item.status, // 添加status字段
            questionCount: item.questionCount || 0,
            createdAt: item.createdAt,
            updatedAt: item.updatedAt
          }
        })
        console.log('转换后的分类数据:', transformedData)
        
        return {
          ...response,
          data: transformedData
        }
      }
      
      return {
        ...response,
        data: []
      }
    } catch (error) {
      console.error('QuestionService.getCategories 错误:', error)
      throw error
    }
  }

  // 获取分类详情
  static async getCategory(id: number): Promise<ApiResponse<Category>> {
    return apiClient.get<Category>(`/api/v1/categories/${id}`)
  }

  // 获取分类下的题目
  static async getCategoryQuestions(categoryId: number, params?: Omit<QuestionSearchParams, 'categoryId'>): Promise<ApiResponse<FrontendPaginatedResponse<Question>>> {
    return apiClient.getPaginated<Question>(`/api/v1/categories/${categoryId}/questions`, params)
  }

  // 获取用户答题记录
  static async getAnswerRecords(params?: {
    page?: number
    pageSize?: number
    categoryId?: number
    isCorrect?: boolean
  }): Promise<ApiResponse<FrontendPaginatedResponse<AnswerRecord>>> {
    return apiClient.getPaginated<AnswerRecord>('/api/v1/answers', params)
  }

  // 获取错题本
  static async getMistakeBook(params?: {
    page?: number
    pageSize?: number
    categoryId?: number
    isMastered?: boolean
  }): Promise<ApiResponse<FrontendPaginatedResponse<MistakeBook>>> {
    return apiClient.getPaginated<MistakeBook>('/api/v1/mistakes', params)
  }

  // 添加到错题本
  static async addToMistakeBook(questionId: number): Promise<ApiResponse<MistakeBook>> {
    return apiClient.post<MistakeBook>('/api/v1/mistakes', { questionId: questionId })
  }

  // 从错题本移除
  static async removeFromMistakeBook(mistakeId: number): Promise<ApiResponse<void>> {
    return apiClient.delete<void>(`/api/v1/mistakes/${mistakeId}`)
  }

  // 标记错题为已掌握
  static async markMistakeAsMastered(questionId: number): Promise<ApiResponse<void>> {
    return apiClient.put<void>(`/api/v1/mistakes/${questionId}/master`)
  }

  // 重置错题状态
  static async resetMistakeStatus(questionId: number): Promise<ApiResponse<void>> {
    return apiClient.put<void>(`/api/v1/mistakes/${questionId}/reset`)
  }

  // 获取分类进度
  static async getCategoryProgress(categoryId: number): Promise<ApiResponse<{
    categoryId: number
    categoryName: string
    totalAnswered: number
    correctAnswered: number
    accuracyRate: number
  }>> {
    return apiClient.get(`/api/v1/categories/${categoryId}/progress`)
  }

  // 清空错题本
  static async clearMistakeBook(): Promise<ApiResponse<{
    message: string
    deletedCount: number
  }>> {
    return apiClient.delete('/api/v1/mistakes/clear')
  }

  // 获取热门标签
  static async getPopularTags(limit?: number): Promise<ApiResponse<string[]>> {
    return apiClient.get<string[]>('/api/v1/questions/tags/popular', { limit })
  }

  // 收藏题目
  static async favoriteQuestion(questionId: number): Promise<ApiResponse<void>> {
    return apiClient.post<void>(`/api/v1/questions/${questionId}/favorite`)
  }

  // 取消收藏题目
  static async unfavoriteQuestion(questionId: number): Promise<ApiResponse<void>> {
    return apiClient.delete<void>(`/api/v1/questions/${questionId}/favorite`)
  }

  // 获取收藏的题目
  static async getFavoriteQuestions(params?: {
    page?: number
    pageSize?: number
  }): Promise<ApiResponse<FrontendPaginatedResponse<Question>>> {
    return apiClient.get<FrontendPaginatedResponse<Question>>('/api/v1/questions/favorites', params)
  }

  // 举报题目
  static async reportQuestion(questionId: number, reason: string, description?: string): Promise<ApiResponse<void>> {
    return apiClient.post<void>(`/api/v1/questions/${questionId}/report`, {
      reason,
      description
    })
  }
}

// 模拟数据服务
export class MockQuestionService {
  private static mockCategories: Category[] = [
    { id: 1, name: '法考题', description: '法律职业资格考试题目', level: 1, questionCount: 1200, icon: 'gavel', color: '#3B82F6', createdAt: '2024-01-01T00:00:00Z', updatedAt: '2024-01-01T00:00:00Z' },
    { id: 2, name: '医考题', description: '医师资格考试题目', level: 1, questionCount: 800, icon: 'stethoscope', color: '#EF4444', createdAt: '2024-01-01T00:00:00Z', updatedAt: '2024-01-01T00:00:00Z' },
    { id: 3, name: '工程类', description: '各类工程考试题目', level: 1, questionCount: 1500, icon: 'hard-hat', color: '#F59E0B', createdAt: '2024-01-01T00:00:00Z', updatedAt: '2024-01-01T00:00:00Z' },
    { id: 11, name: '法律基础', description: '法律基础知识', parentId: 1, level: 2, questionCount: 400, icon: 'book-open', color: '#3B82F6', createdAt: '2024-01-01T00:00:00Z', updatedAt: '2024-01-01T00:00:00Z' },
    { id: 12, name: '法律条款', description: '法律条款解读', parentId: 1, level: 2, questionCount: 500, icon: 'file-text', color: '#3B82F6', createdAt: '2024-01-01T00:00:00Z', updatedAt: '2024-01-01T00:00:00Z' }
  ]

  private static generateMockQuestion(id: number, categoryId: number = 1): Question {
    const difficulties = [QuestionDifficulty.EASY, QuestionDifficulty.MEDIUM, QuestionDifficulty.HARD]
    const types = [QuestionType.SINGLE_CHOICE, QuestionType.MULTIPLE_CHOICE, QuestionType.TRUE_FALSE]
    const category = this.mockCategories.find(c => c.id === categoryId) || this.mockCategories[0]
    
    const difficulty = difficulties[Math.floor(Math.random() * difficulties.length)]
    const type = types[Math.floor(Math.random() * types.length)]
    
    const options: QuestionOption[] = [
      { id: 'A', content: '选项A的内容', isCorrect: true },
      { id: 'B', content: '选项B的内容', isCorrect: false },
      { id: 'C', content: '选项C的内容', isCorrect: false },
      { id: 'D', content: '选项D的内容', isCorrect: false }
    ]
    
    if (type === QuestionType.MULTIPLE_CHOICE) {
      options[2].isCorrect = true // 多选题
    } else if (type === QuestionType.TRUE_FALSE) {
      options.splice(2) // 判断题只有两个选项
      options[0].content = '正确'
      options[1].content = '错误'
    }
    
    return {
      id,
      title: `第${id}题 - ${category.name}相关题目`,
      content: `这是一道关于${category.name}的${difficulty === 'easy' ? '简单' : difficulty === 'medium' ? '中等' : '困难'}题目。请仔细阅读题目内容，选择正确答案。题目内容会根据实际情况进行详细描述，包含相关的背景信息和具体要求。`,
      type,
      difficulty,
      options,
      correctAnswer: type === QuestionType.MULTIPLE_CHOICE ? ['A', 'C'] : 'A',
      explanation: `这道题的正确答案是${type === QuestionType.MULTIPLE_CHOICE ? 'A和C' : 'A'}。解析：根据相关理论和实践经验，选项A${type === QuestionType.MULTIPLE_CHOICE ? '和选项C' : ''}是正确的，因为它们符合题目要求的条件和标准。其他选项存在明显的错误或不完整之处。`,
      categoryId,
      categoryName: category.name,
      tags: ['基础知识', '重点题型', '常考题'],
      createdAt: new Date(Date.now() - Math.random() * 30 * 24 * 60 * 60 * 1000).toISOString(),
      updatedAt: new Date().toISOString(),
      isAnswered: Math.random() > 0.7,
      isCorrect: Math.random() > 0.3
    }
  }

  static async getQuestions(params?: QuestionSearchParams): Promise<ApiResponse<FrontendPaginatedResponse<Question>>> {
    await new Promise(resolve => setTimeout(resolve, 500))
    
    const page = params?.page || 1
    const pageSize = params?.pageSize || 10
    const total = 1000
    
    const questions = Array.from({ length: pageSize }, (_, i) => 
      this.generateMockQuestion((page - 1) * pageSize + i + 1, params?.categoryId)
    )
    
    return {
      code: 200,
      message: '获取题目列表成功',
      success: true,
      data: {
        items: questions,
        total,
        page,
        pageSize,
        totalPages: Math.ceil(total / pageSize)
      }
    }
  }

  static async getQuestion(id: number): Promise<ApiResponse<Question>> {
    await new Promise(resolve => setTimeout(resolve, 300))
    
    return {
      code: 200,
      message: '获取题目详情成功',
      success: true,
      data: this.generateMockQuestion(id)
    }
  }

  static async getCategories(parentId?: number): Promise<ApiResponse<Category[]>> {
    await new Promise(resolve => setTimeout(resolve, 200))
    
    const categories = parentId 
      ? this.mockCategories.filter(c => c.parentId === parentId)
      : this.mockCategories.filter(c => !c.parentId)
    
    return {
      code: 200,
      message: '获取分类列表成功',
      success: true,
      data: categories
    }
  }

  static async submitAnswer(data: SubmitAnswerRequest): Promise<ApiResponse<AnswerRecord>> {
    await new Promise(resolve => setTimeout(resolve, 300))
    
    const question = this.generateMockQuestion(data.questionId)
    const isCorrect = JSON.stringify(data.userAnswer) === JSON.stringify(question.correctAnswer)
    
    return {
      code: 200,
      message: '提交答案成功',
      success: true,
      data: {
        id: Date.now(),
        questionId: data.questionId,
        userId: 1,
        userAnswer: data.userAnswer,
        isCorrect,
        timeSpent: data.timeSpent,
        answeredAt: new Date().toISOString(),
        question
      }
    }
  }

  static async getMistakeBook(params?: any): Promise<ApiResponse<FrontendPaginatedResponse<MistakeBook>>> {
    await new Promise(resolve => setTimeout(resolve, 400))
    
    const page = params?.page || 1
    const pageSize = params?.pageSize || 10
    const total = 50
    
    const mistakes = Array.from({ length: Math.min(pageSize, total - (page - 1) * pageSize) }, (_, i) => ({
      id: (page - 1) * pageSize + i + 1,
      questionId: (page - 1) * pageSize + i + 1,
      userId: 1,
      addedAt: new Date(Date.now() - Math.random() * 7 * 24 * 60 * 60 * 1000).toISOString(),
      reviewCount: Math.floor(Math.random() * 5),
      lastReviewAt: Math.random() > 0.5 ? new Date(Date.now() - Math.random() * 24 * 60 * 60 * 1000).toISOString() : undefined,
      isMastered: Math.random() > 0.7,
      question: this.generateMockQuestion((page - 1) * pageSize + i + 1)
    }))
    
    return {
      code: 200,
      message: '获取错题本成功',
      success: true,
      data: {
        items: mistakes,
        total,
        page,
        pageSize,
        totalPages: Math.ceil(total / pageSize)
      }
    }
  }

  static async addToMistakeBook(questionId: number): Promise<ApiResponse<MistakeBook>> {
    await new Promise(resolve => setTimeout(resolve, 300))
    
    return {
      code: 200,
      message: '添加到错题本成功',
      success: true,
      data: {
        id: Date.now(),
        questionId,
        userId: 1,
        addedAt: new Date().toISOString(),
        reviewCount: 0,
        isMastered: false,
        question: this.generateMockQuestion(questionId)
      }
    }
  }
}

// 使用真实的题目服务
export const questionService = QuestionService

export default questionService