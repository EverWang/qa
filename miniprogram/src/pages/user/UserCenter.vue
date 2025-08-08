<template>
  <div class="min-h-screen bg-gray-50">
    <!-- 顶部导航 -->
    <div class="bg-white shadow-sm border-b">
      <div class="max-w-4xl mx-auto px-4 py-3 flex items-center justify-between">
        <button @click="goBack" class="flex items-center text-gray-600 hover:text-gray-800">
          <ArrowLeft class="w-5 h-5 mr-1" />
          <span>返回</span>
        </button>
        <h1 class="text-lg font-medium text-gray-800">个人中心</h1>
        <button @click="goToSettings" class="flex items-center text-gray-600 hover:text-gray-800">
          <Settings class="w-5 h-5" />
        </button>
      </div>
    </div>

    <div class="max-w-4xl mx-auto p-4">
      <!-- 用户信息卡片 -->
      <div class="bg-white rounded-lg shadow-sm p-6 mb-6">
        <div class="flex items-center space-x-4">
          <div class="w-16 h-16 rounded-full overflow-hidden bg-gray-200">
            <img 
              :src="authStore.user?.avatar || defaultAvatar" 
              :alt="authStore.user?.nickname"
              class="w-full h-full object-cover"
            />
          </div>
          <div class="flex-1">
            <h2 class="text-xl font-semibold text-gray-800">{{ authStore.user?.nickname || '未登录' }}</h2>
            <p class="text-gray-500 text-sm">{{ authStore.isGuest ? '游客模式' : '正式用户' }}</p>
            <div class="flex items-center space-x-4 mt-2">
              <div class="flex items-center text-sm text-gray-600">
                <Calendar class="w-4 h-4 mr-1" />
                <span>加入 {{ joinDays }} 天</span>
              </div>
              <div class="flex items-center text-sm text-gray-600">
                <Trophy class="w-4 h-4 mr-1" />
                <span>等级 {{ userLevel }}</span>
              </div>
            </div>
          </div>
          <button 
            v-if="authStore.isGuest"
            @click="goToLogin"
            class="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white text-sm rounded-lg transition-colors duration-200"
          >
            立即登录
          </button>
        </div>
      </div>

      <!-- 学习统计 -->
      <div class="bg-white rounded-lg shadow-sm p-6 mb-6">
        <h3 class="text-lg font-semibold text-gray-800 mb-4">学习统计</h3>
        <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
          <div class="text-center p-4 bg-blue-50 rounded-lg">
            <div class="text-2xl font-bold text-blue-600">{{ stats.totalAnswered }}</div>
            <div class="text-sm text-gray-600 mt-1">已答题目</div>
          </div>
          <div class="text-center p-4 bg-green-50 rounded-lg">
            <div class="text-2xl font-bold text-green-600">{{ stats.correctRate }}%</div>
            <div class="text-sm text-gray-600 mt-1">正确率</div>
          </div>
          <div class="text-center p-4 bg-orange-50 rounded-lg">
            <div class="text-2xl font-bold text-orange-600">{{ stats.mistakeCount }}</div>
            <div class="text-sm text-gray-600 mt-1">错题数量</div>
          </div>
          <div class="text-center p-4 bg-purple-50 rounded-lg">
            <div class="text-2xl font-bold text-purple-600">{{ stats.studyDays }}</div>
            <div class="text-sm text-gray-600 mt-1">学习天数</div>
          </div>
        </div>
      </div>

      <!-- 功能菜单 -->
      <div class="bg-white rounded-lg shadow-sm mb-6">
        <div class="p-4 border-b border-gray-100">
          <h3 class="text-lg font-semibold text-gray-800">功能菜单</h3>
        </div>
        <div class="divide-y divide-gray-100">
          <button 
            @click="goToRecords"
            class="w-full flex items-center justify-between p-4 hover:bg-gray-50 transition-colors duration-200"
          >
            <div class="flex items-center space-x-3">
              <div class="w-10 h-10 bg-blue-100 rounded-lg flex items-center justify-center">
                <FileText class="w-5 h-5 text-blue-600" />
              </div>
              <div class="text-left">
                <div class="font-medium text-gray-800">做题记录</div>
                <div class="text-sm text-gray-500">查看历史答题记录</div>
              </div>
            </div>
            <ChevronRight class="w-5 h-5 text-gray-400" />
          </button>
          
          <button 
            @click="goToMistakes"
            class="w-full flex items-center justify-between p-4 hover:bg-gray-50 transition-colors duration-200"
          >
            <div class="flex items-center space-x-3">
              <div class="w-10 h-10 bg-orange-100 rounded-lg flex items-center justify-center">
                <BookmarkMinus class="w-5 h-5 text-orange-600" />
              </div>
              <div class="text-left">
                <div class="font-medium text-gray-800">错题本</div>
                <div class="text-sm text-gray-500">复习错误题目</div>
              </div>
            </div>
            <div class="flex items-center space-x-2">
              <span class="text-sm text-orange-600 font-medium">{{ stats.mistakeCount }}</span>
              <ChevronRight class="w-5 h-5 text-gray-400" />
            </div>
          </button>
          
          <button 
            @click="goToSettings"
            class="w-full flex items-center justify-between p-4 hover:bg-gray-50 transition-colors duration-200"
          >
            <div class="flex items-center space-x-3">
              <div class="w-10 h-10 bg-gray-100 rounded-lg flex items-center justify-center">
                <Settings class="w-5 h-5 text-gray-600" />
              </div>
              <div class="text-left">
                <div class="font-medium text-gray-800">设置</div>
                <div class="text-sm text-gray-500">个人偏好设置</div>
              </div>
            </div>
            <ChevronRight class="w-5 h-5 text-gray-400" />
          </button>
        </div>
      </div>

      <!-- 最近学习 -->
      <div class="bg-white rounded-lg shadow-sm mb-6">
        <div class="p-4 border-b border-gray-100">
          <h3 class="text-lg font-semibold text-gray-800">最近学习</h3>
        </div>
        <div class="p-4">
          <div v-if="recentQuestions.length > 0" class="space-y-3">
            <div 
              v-for="question in recentQuestions"
              :key="question.id"
              @click="goToQuestion(question.id)"
              class="flex items-center justify-between p-3 bg-gray-50 rounded-lg hover:bg-gray-100 cursor-pointer transition-colors duration-200"
            >
              <div class="flex-1">
                <div class="font-medium text-gray-800 mb-1">{{ question.title }}</div>
                <div class="flex items-center space-x-2 text-sm text-gray-500">
                  <span class="px-2 py-1 bg-blue-100 text-blue-800 text-xs rounded-full">
                    {{ question.category.name }}
                  </span>
                  <span>{{ formatDate(question.answeredAt) }}</span>
                </div>
              </div>
              <div class="flex items-center space-x-2">
                <CheckCircle v-if="question.isCorrect" class="w-5 h-5 text-green-600" />
                <XCircle v-else class="w-5 h-5 text-red-600" />
                <ChevronRight class="w-4 h-4 text-gray-400" />
              </div>
            </div>
          </div>
          <div v-else class="text-center py-8 text-gray-500">
            <BookOpen class="w-12 h-12 mx-auto mb-3 text-gray-300" />
            <p>还没有答题记录</p>
            <button 
              @click="goToHome"
              class="mt-3 text-blue-600 hover:text-blue-700 text-sm"
            >
              去刷题 →
            </button>
          </div>
        </div>
      </div>

      <!-- 退出登录 -->
      <div v-if="!authStore.isGuest" class="bg-white rounded-lg shadow-sm">
        <button 
          @click="handleLogout"
          class="w-full p-4 text-red-600 hover:bg-red-50 rounded-lg transition-colors duration-200 flex items-center justify-center space-x-2"
        >
          <LogOut class="w-5 h-5" />
          <span>退出登录</span>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import {
  ArrowLeft, Settings, Calendar, Trophy, FileText, BookmarkMinus,
  ChevronRight, CheckCircle, XCircle, BookOpen, LogOut
} from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'

interface UserStats {
  totalAnswered: number
  correctRate: number
  mistakeCount: number
  studyDays: number
}

interface RecentQuestion {
  id: number
  title: string
  category: {
    id: number
    name: string
  }
  isCorrect: boolean
  answeredAt: string
}

const router = useRouter()
const authStore = useAuthStore()

const defaultAvatar = 'https://trae-api-us.mchost.guru/api/ide/v1/text_to_image?prompt=simple%20user%20avatar%20profile%20picture&image_size=square'

const joinDays = ref(15)
const userLevel = ref(3)

const stats = ref<UserStats>({
  totalAnswered: 0,
  correctRate: 0,
  mistakeCount: 0,
  studyDays: 0
})

const recentQuestions = ref<RecentQuestion[]>([])

// 格式化日期
const formatDate = (dateString: string) => {
  const date = new Date(dateString)
  const now = new Date()
  const diffTime = now.getTime() - date.getTime()
  const diffDays = Math.floor(diffTime / (1000 * 60 * 60 * 24))
  
  if (diffDays === 0) {
    return '今天'
  } else if (diffDays === 1) {
    return '昨天'
  } else if (diffDays < 7) {
    return `${diffDays}天前`
  } else {
    return date.toLocaleDateString()
  }
}

// 导航方法
const goBack = () => {
  router.back()
}

const goToHome = () => {
  router.push('/')
}

const goToLogin = () => {
  router.push('/auth/login')
}

const goToSettings = () => {
  router.push('/user/settings')
}

const goToRecords = () => {
  router.push('/user/records')
}

const goToMistakes = () => {
  router.push('/user/mistakes')
}

const goToQuestion = (questionId: number) => {
  router.push(`/question/${questionId}`)
}

// 退出登录
const handleLogout = () => {
  if (confirm('确定要退出登录吗？')) {
    authStore.logout()
    router.push('/')
  }
}

// 加载用户统计
const loadUserStats = async () => {
  if (authStore.isGuest) {
    stats.value = {
      totalAnswered: 0,
      correctRate: 0,
      mistakeCount: 0,
      studyDays: 0
    }
    return
  }
  
  try {
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 500))
    
    stats.value = {
      totalAnswered: 156,
      correctRate: 78,
      mistakeCount: 34,
      studyDays: 12
    }
  } catch (error) {
    console.error('加载用户统计失败:', error)
  }
}

// 加载最近答题记录
const loadRecentQuestions = async () => {
  if (authStore.isGuest) {
    recentQuestions.value = []
    return
  }
  
  try {
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 300))
    
    recentQuestions.value = [
      {
        id: 1,
        title: '道路工程施工安全规范',
        category: { id: 3, name: '道路工程' },
        isCorrect: true,
        answeredAt: new Date(Date.now() - 1000 * 60 * 60 * 2).toISOString() // 2小时前
      },
      {
        id: 2,
        title: '法律条文理解',
        category: { id: 1, name: '法考题' },
        isCorrect: false,
        answeredAt: new Date(Date.now() - 1000 * 60 * 60 * 24).toISOString() // 1天前
      },
      {
        id: 3,
        title: '临床医学基础',
        category: { id: 2, name: '医考题' },
        isCorrect: true,
        answeredAt: new Date(Date.now() - 1000 * 60 * 60 * 24 * 2).toISOString() // 2天前
      }
    ]
  } catch (error) {
    console.error('加载最近答题记录失败:', error)
  }
}

onMounted(() => {
  loadUserStats()
  loadRecentQuestions()
})
</script>

<style scoped>
/* 自定义样式 */
</style>