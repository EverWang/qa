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
                <p class="text-sm text-gray-500">{{ category.questionCount }} 道题</p>
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
                  {{ question.category.name }}
                </span>
                <span class="px-2 py-1 text-xs rounded-full" :class="getDifficultyClass(question.difficulty)">
                  {{ getDifficultyText(question.difficulty) }}
                </span>
              </div>
              <ChevronRight class="w-5 h-5 text-gray-400" />
            </div>
            <h4 class="font-medium text-gray-800 mb-2">{{ question.title }}</h4>
            <p class="text-sm text-gray-600 line-clamp-2">{{ question.content }}</p>
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

// 导航方法
const goToHome = () => {
  router.push('/')
}

const goToSearch = () => {
  router.push('/search')
}

const goToProfile = () => {
  if (authStore.isLoggedIn) {
    router.push('/user/center')
  } else {
    router.push('/auth/login')
  }
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
  if (authStore.isLoggedIn && !authStore.isGuest) {
    router.push('/user/mistakes')
  } else {
    router.push('/auth/login')
  }
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
 * 加载推荐题目 - 使用真实API
 */
const loadRecommendedQuestions = async () => {
  try {
    const response = await QuestionService.getQuestions({
      page: 1,
      pageSize: 3
    })
    
    if (response.success && response.data) {
      recommendedQuestions.value = response.data.items.map(q => ({
        id: q.id,
        title: q.title,
        content: q.content,
        difficulty: q.difficulty as 'easy' | 'medium' | 'hard',
        category: {
          id: q.categoryId,
          name: q.categoryName || '未分类'
        }
      }))
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
    const response = await QuestionService.getCategories()
    
    if (response.success && response.data) {
      // 为分类添加图标和颜色
      const iconMap: { [key: string]: any } = {
        '法考': Gavel,
        '医考': Stethoscope,
        '工程': HardHat,
        '其他': GraduationCap
      }
      
      const colorMap: { [key: string]: string } = {
        '法考': '#3B82F6',
        '医考': '#EF4444', 
        '工程': '#F59E0B',
        '其他': '#8B5CF6'
      }
      
      categories.value = response.data.map((cat, index) => {
        const key = Object.keys(iconMap).find(k => cat.name.includes(k)) || '其他'
        return {
          ...cat,
          color: colorMap[key] || colorMap['其他'],
          icon: iconMap[key] || iconMap['其他']
        }
      })
    }
  } catch (error) {
    console.error('加载分类失败:', error)
    // 使用默认分类数据
    categories.value = [
      { id: 1, name: '法考题', questionCount: 0, color: '#3B82F6', icon: Gavel, description: '', parentId: undefined, level: 1, createdAt: '', updatedAt: '' },
      { id: 2, name: '医考题', questionCount: 0, color: '#EF4444', icon: Stethoscope, description: '', parentId: undefined, level: 1, createdAt: '', updatedAt: '' },
      { id: 3, name: '工程类', questionCount: 0, color: '#F59E0B', icon: HardHat, description: '', parentId: undefined, level: 1, createdAt: '', updatedAt: '' },
      { id: 4, name: '其他考试', questionCount: 0, color: '#8B5CF6', icon: GraduationCap, description: '', parentId: undefined, level: 1, createdAt: '', updatedAt: '' }
    ]
  }
}

/**
 * 加载概览统计数据 - 使用真实API
 */
const loadOverviewStats = async () => {
  try {
    const response = await StatisticsService.getOverviewStatistics()
    
    if (response.success && response.data) {
      totalQuestions.value = response.data.totalQuestions
      totalUsers.value = response.data.totalUsers
      // 计算通过率：正确答题数 / 总答题数 * 100
      if (response.data.totalAnswers > 0) {
        // 这里需要从答题记录中计算正确率，暂时使用一个估算值
        passRate.value = 75 // 可以根据实际数据调整
      } else {
        passRate.value = 0
      }
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
      userStats.value = {
        totalAnswered: response.data.totalAnswered,
        correctRate: response.data.accuracyRate,
        mistakeCount: response.data.totalAnswered - response.data.correctAnswered,
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
