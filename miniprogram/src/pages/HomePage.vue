<template>
  <div class="min-h-screen bg-gray-50">
    <!-- 顶部导航栏 -->
    <div class="bg-white shadow-sm border-b sticky top-0 z-10">
      <div class="max-w-6xl mx-auto px-4 py-3">
        <div class="flex items-center justify-between">
          <div class="flex items-center space-x-3">
            <div class="w-8 h-8 bg-gradient-to-r from-blue-500 to-indigo-600 rounded-lg flex items-center justify-center">
              <BookOpen class="w-5 h-5 text-white" />
            </div>
            <h1 class="text-xl font-bold text-gray-800">刷刷题</h1>
          </div>
          <div class="flex items-center space-x-3">
            <button @click="goToSearch" class="p-2 text-gray-600 hover:text-gray-800 hover:bg-gray-100 rounded-lg">
              <Search class="w-5 h-5" />
            </button>
            <button @click="goToProfile" class="flex items-center space-x-2 text-gray-600 hover:text-gray-800">
              <User class="w-5 h-5" />
              <span class="text-sm">{{ authStore.user?.nickname || '登录' }}</span>
            </button>
          </div>
        </div>
      </div>
    </div>

    <div class="max-w-6xl mx-auto p-4">
      <!-- 欢迎横幅 -->
      <div class="bg-gradient-to-r from-blue-500 to-indigo-600 rounded-xl p-6 mb-6 text-white">
        <h2 class="text-2xl font-bold mb-2">欢迎来到刷刷题</h2>
        <p class="text-blue-100 mb-4">专业的在线刷题平台，助你轻松通过考试</p>
        <div class="flex items-center space-x-6 text-sm">
          <div class="flex items-center space-x-1">
            <Target class="w-4 h-4" />
            <span>{{ totalQuestions }}+ 道题目</span>
          </div>
          <div class="flex items-center space-x-1">
            <Users class="w-4 h-4" />
            <span>{{ totalUsers }}+ 用户</span>
          </div>
          <div class="flex items-center space-x-1">
            <TrendingUp class="w-4 h-4" />
            <span>{{ passRate }}% 通过率</span>
          </div>
        </div>
      </div>

      <!-- 分类导航 -->
      <div class="mb-6">
        <h3 class="text-lg font-semibold text-gray-800 mb-4">题目分类</h3>
        <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
          <div
            v-for="category in categories"
            :key="category.id"
            @click="goToCategory(category.id)"
            class="bg-white rounded-lg p-4 shadow-sm hover:shadow-md transition-shadow duration-200 cursor-pointer border hover:border-blue-200"
          >
            <div class="flex items-center space-x-3">
              <div class="w-10 h-10 rounded-lg flex items-center justify-center" :style="{ backgroundColor: category.color + '20', color: category.color }">
                <component :is="category.icon" class="w-5 h-5" />
              </div>
              <div class="flex-1">
                <h4 class="font-medium text-gray-800">{{ category.name }}</h4>
                <p class="text-sm text-gray-500">{{ category.questionCount || category.question_count || 0 }} 道题</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 推荐题目 -->
      <div class="mb-6">
        <div class="flex items-center justify-between mb-4">
          <h3 class="text-lg font-semibold text-gray-800">推荐题目</h3>
          <button @click="refreshRecommended" class="text-blue-600 hover:text-blue-700 text-sm flex items-center space-x-1">
            <RotateCcw class="w-4 h-4" />
            <span>换一批</span>
          </button>
        </div>
        <div class="grid gap-4">
          <div
            v-for="question in recommendedQuestions"
            :key="question.id"
            @click="goToQuestion(question.id)"
            class="bg-white rounded-lg p-4 shadow-sm hover:shadow-md transition-shadow duration-200 cursor-pointer border hover:border-blue-200"
          >
            <div class="flex items-start justify-between mb-3">
              <div class="flex items-center space-x-2">
                  <span class="px-2 py-1 bg-blue-100 text-blue-800 text-xs rounded-full">
                    {{ getCategoryName(question.category.id) }}
                  </span>
                  <span class="px-2 py-1 text-xs rounded-full" :class="getDifficultyClass(question.difficulty)">
                    {{ getDifficultyText(question.difficulty) }}
                  </span>
                </div>
              <ChevronRight class="w-5 h-5 text-gray-400" />
            </div>
            <h4 class="font-medium text-gray-800 mb-2 line-clamp-2">{{ question.title }}</h4>
              <p class="text-sm text-gray-600 line-clamp-3">{{ question.content }}</p>
          </div>
        </div>
      </div>

      <!-- 学习统计 -->
      <div v-if="authStore.isLoggedIn && !authStore.isGuest" class="bg-white rounded-lg p-6 shadow-sm">
        <h3 class="text-lg font-semibold text-gray-800 mb-4">学习统计</h3>
        <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
          <div class="text-center">
            <div class="text-2xl font-bold text-blue-600">{{ userStats.totalAnswered }}</div>
            <div class="text-sm text-gray-500">已答题目</div>
          </div>
          <div class="text-center">
            <div class="text-2xl font-bold text-green-600">{{ userStats.correctRate }}%</div>
            <div class="text-sm text-gray-500">正确率</div>
          </div>
          <div class="text-center">
            <div class="text-2xl font-bold text-orange-600">{{ userStats.mistakeCount }}</div>
            <div class="text-sm text-gray-500">错题数量</div>
          </div>
          <div class="text-center">
            <div class="text-2xl font-bold text-purple-600">{{ userStats.studyDays }}</div>
            <div class="text-sm text-gray-500">学习天数</div>
          </div>
        </div>
      </div>
    </div>

    <!-- 底部导航 -->
    <div class="fixed bottom-0 left-0 right-0 bg-white border-t border-gray-200 px-4 py-2">
      <div class="max-w-6xl mx-auto">
        <div class="flex items-center justify-around">
          <button @click="goToHome" class="flex flex-col items-center space-y-1 text-blue-600">
            <Home class="w-5 h-5" />
            <span class="text-xs">首页</span>
          </button>
          <button @click="goToCategory()" class="flex flex-col items-center space-y-1 text-gray-500">
            <Grid3X3 class="w-5 h-5" />
            <span class="text-xs">分类</span>
          </button>
          <button @click="goToMistakes" class="flex flex-col items-center space-y-1 text-gray-500">
            <BookmarkMinus class="w-5 h-5" />
            <span class="text-xs">错题本</span>
          </button>
          <button @click="goToProfile" class="flex flex-col items-center space-y-1 text-gray-500">
            <User class="w-5 h-5" />
            <span class="text-xs">我的</span>
          </button>
        </div>
      </div>
    </div>

    <!-- 底部占位 -->
    <div class="h-16"></div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import {
  BookOpen, Search, User, Target, Users, TrendingUp, ChevronRight,
  RotateCcw, Home, Grid3X3, BookmarkMinus, Gavel, Stethoscope,
  HardHat, GraduationCap
} from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'
import { QuestionService } from '@/services/question'
import { StatisticsService } from '@/services/statistics'
import type { Question, Category } from '@/services/question'
import type { OverviewStatistics } from '@/services/statistics'

interface CategoryWithIcon extends Category {
  color: string
  icon: any
}

interface UserStats {
  totalAnswered: number
  correctRate: number
  mistakeCount: number
  studyDays: number
}

const router = useRouter()
const authStore = useAuthStore()

const totalQuestions = ref(0)
const totalUsers = ref(0)
const passRate = ref(0)

const categories = ref<CategoryWithIcon[]>([])
const loading = ref(false)

const recommendedQuestions = ref<{
  id: number
  title: string
  content: string
  difficulty: 'easy' | 'medium' | 'hard'
  category: {
    id: number
    name: string
  }
}[]>([])

const userStats = ref<UserStats>({
  totalAnswered: 0,
  correctRate: 0,
  mistakeCount: 0,
  studyDays: 0
})

// 难度样式
const getDifficultyClass = (difficulty: string) => {
  switch (difficulty) {
    case 'easy':
      return 'bg-green-100 text-green-800'
    case 'medium':
      return 'bg-yellow-100 text-yellow-800'
    case 'hard':
      return 'bg-red-100 text-red-800'
    default:
      return 'bg-gray-100 text-gray-800'
  }
}

const getDifficultyText = (difficulty: string) => {
  switch (difficulty) {
    case 'easy':
      return '简单'
    case 'medium':
      return '中等'
    case 'hard':
      return '困难'
    default:
      return '未知'
  }
}

// 根据分类ID获取分类名称
const getCategoryName = (categoryId: number) => {
  const category = categories.value.find(cat => cat.id === categoryId)
  return category ? category.name : '未分类'
}

// 导航方法
const goToHome = () => {
  router.push('/')
}

const goToSearch = () => {
  router.push('/search')
}

const goToProfile = () => {
  router.push('/user/center')
}

const goToCategory = (categoryId?: number) => {
  if (categoryId) {
    router.push(`/category/${categoryId}`)
  } else {
    router.push('/category')
  }
}

const goToQuestion = (questionId: number) => {
  router.push(`/question/${questionId}`)
}

const goToMistakes = () => {
  router.push('/user/mistakes')
}

/**
 * 刷新推荐题目
 */
const refreshRecommended = async () => {
  try {
    loading.value = true
    await loadRecommendedQuestions()
  } catch (error) {
    console.error('刷新推荐题目失败:', error)
  } finally {
    loading.value = false
  }
}

/**
 * 加载推荐题目 - 使用真实API，支持随机获取
 */
const loadRecommendedQuestions = async () => {
  try {
    console.log('开始加载推荐题目...')
    
    // 先获取题目总数
    const totalResponse = await QuestionService.getQuestions({ page: 1, pageSize: 1 })
    const totalQuestions = totalResponse.data?.total || 0
    
    if (totalQuestions === 0) {
      recommendedQuestions.value = []
      return
    }
    
    // 随机选择页码，确保能获取到不同的题目
    const maxPage = Math.ceil(totalQuestions / 3)
    const randomPage = Math.floor(Math.random() * maxPage) + 1
    
    const response = await QuestionService.getQuestions({
      page: randomPage,
      pageSize: 3
    })
    
    console.log('推荐题目API响应:', response)
    
    if (response.success && response.data) {
      recommendedQuestions.value = response.data.items.map(q => ({
        id: q.id,
        title: q.title,
        content: q.content,
        difficulty: q.difficulty as 'easy' | 'medium' | 'hard',
        category: {
          id: q.category_id || q.categoryId,
          name: q.category?.name || q.categoryName || '未分类'
        }
      }))
      console.log('推荐题目加载成功，数量:', recommendedQuestions.value.length)
    } else {
      console.log('推荐题目API响应失败:', response)
    }
  } catch (error) {
    console.error('加载推荐题目失败:', error)
    // 保持空数组，不显示错误给用户
    recommendedQuestions.value = []
  }
}

/**
 * 加载分类数据 - 使用真实API
 */
const loadCategories = async () => {
  try {
    console.log('开始加载分类数据...')
    const response = await QuestionService.getCategories()
    
    console.log('分类API响应:', response)
    
    if (response.success && response.data) {
      // 为分类添加图标和颜色
      const getIconAndColor = (categoryName: string) => {
        const name = categoryName.toLowerCase()
        
        // 法律相关
        if (name.includes('法') || name.includes('律') || name.includes('司法')) {
          return { icon: Gavel, color: '#3B82F6' }
        }
        // 医学相关
        if (name.includes('医') || name.includes('药') || name.includes('护理') || name.includes('健康')) {
          return { icon: Stethoscope, color: '#EF4444' }
        }
        // 建筑工程相关
        if (name.includes('建筑') || name.includes('工程') || name.includes('施工') || name.includes('道路') || name.includes('桥梁')) {
          return { icon: HardHat, color: '#F59E0B' }
        }
        // 默认
        return { icon: GraduationCap, color: '#8B5CF6' }
      }
      
      // 只显示一级分类（没有父分类的分类）且状态为启用的分类
      const topLevelCategories = response.data.filter(cat => 
        (!cat.parent_id || cat.parent_id === 0) && cat.status === 1
      )
      
      categories.value = topLevelCategories.map((cat) => {
        const { icon, color } = getIconAndColor(cat.name)
        console.log('处理分类:', cat.name, '图标:', icon.name, '颜色:', color)
        return {
          ...cat,
          color,
          icon
        }
      })
      
      console.log('分类加载成功，数量:', categories.value.length)
      console.log('分类详情:', categories.value)
    } else {
      console.log('分类API响应失败:', response)
    }
  } catch (error) {
      console.error('加载分类失败:', error)
      // 不使用默认分类数据，保持空数组让用户知道加载失败
      categories.value = []
    }
}

/**
 * 加载概览统计数据 - 使用真实的统计API
 */
const loadOverviewStats = async () => {
  try {
    console.log('开始加载概览统计数据...')
    const response = await StatisticsService.getOverviewStatistics()
    
    console.log('概览统计API响应:', response)
    
    if (response.success && response.data) {
      const data = response.data
      totalQuestions.value = data.totalQuestions || 0
      totalUsers.value = data.totalUsers || 0
      // 计算通过率：如果有答题数据，计算通过率，否则使用默认值
      if (data.totalAnswers && data.totalAnswers > 0) {
        passRate.value = Math.round((data.correctAnswers || 0) / data.totalAnswers * 100)
      } else {
        passRate.value = 75 // 默认通过率
      }
      console.log('概览统计加载成功:', { totalQuestions: totalQuestions.value, totalUsers: totalUsers.value, passRate: passRate.value })
    } else {
      console.log('概览统计API响应失败:', response)
      // 使用默认值
      totalQuestions.value = 0
      totalUsers.value = 0
      passRate.value = 0
    }
  } catch (error) {
    console.error('加载概览统计失败:', error)
    // 使用默认值
    totalQuestions.value = 0
    totalUsers.value = 0
    passRate.value = 0
  }
}

/**
 * 加载用户统计 - 使用真实API
 */
const loadUserStats = async () => {
  // 只有在用户已登录且不是游客时才加载用户统计
  if (!authStore.isLoggedIn || authStore.isGuest || !authStore.token) {
    // 重置为默认值
    userStats.value = {
      totalAnswered: 0,
      correctRate: 0,
      mistakeCount: 0,
      studyDays: 0
    }
    return
  }
  
  try {
    const response = await StatisticsService.getUserStatistics()
    
    if (response.success && response.data) {
      const data = response.data
      userStats.value = {
        totalAnswered: data.totalAnswered || 0,
        correctRate: Math.round(data.accuracyRate || 0),
        mistakeCount: (data.totalAnswered || 0) - (data.correctAnswered || 0),
        studyDays: Math.ceil((Date.now() - new Date(authStore.user?.createdAt || Date.now()).getTime()) / (1000 * 60 * 60 * 24))
      }
    }
  } catch (error) {
    console.error('加载用户统计失败:', error)
    // 使用默认值
    userStats.value = {
      totalAnswered: 0,
      correctRate: 0,
      mistakeCount: 0,
      studyDays: 0
    }
  }
}

onMounted(async () => {
  loading.value = true
  try {
    // 并行加载数据
    await Promise.all([
      loadOverviewStats(),
      loadCategories(),
      loadRecommendedQuestions(),
      loadUserStats()
    ])
  } catch (error) {
    console.error('页面初始化失败:', error)
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
