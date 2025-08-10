<template>
  <div class="min-h-screen bg-gray-50">
    <!-- 顶部导航 -->
    <div class="bg-white shadow-sm border-b">
      <div class="max-w-4xl mx-auto px-4 py-3 flex items-center justify-between">
        <button @click="goBack" class="flex items-center text-gray-600 hover:text-gray-800">
          <ArrowLeft class="w-5 h-5 mr-1" />
          <span>返回</span>
        </button>
        <h1 class="text-lg font-medium text-gray-800">错题本</h1>
        <button @click="clearAllMistakes" class="text-red-600 hover:text-red-700 text-sm">
          清空
        </button>
      </div>
    </div>

    <div class="max-w-4xl mx-auto p-4">
      <!-- 统计信息 -->
      <div class="bg-white rounded-lg shadow-sm p-4 mb-6">
        <div class="grid grid-cols-3 gap-4 text-center">
          <div>
            <div class="text-2xl font-bold text-red-600">{{ mistakeQuestions.length }}</div>
            <div class="text-sm text-gray-500">错题总数</div>
          </div>
          <div>
            <div class="text-2xl font-bold text-orange-600">{{ reviewedCount }}</div>
            <div class="text-sm text-gray-500">已复习</div>
          </div>
          <div>
            <div class="text-2xl font-bold text-green-600">{{ masteredCount }}</div>
            <div class="text-sm text-gray-500">已掌握</div>
          </div>
        </div>
      </div>

      <!-- 筛选和排序 -->
      <div class="bg-white rounded-lg shadow-sm p-4 mb-4">
        <div class="flex items-center justify-between mb-4">
          <h3 class="font-medium text-gray-800">筛选条件</h3>
          <button @click="resetFilters" class="text-blue-600 hover:text-blue-700 text-sm">
            重置
          </button>
        </div>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <!-- 分类筛选 -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">题目分类</label>
            <select v-model="selectedCategory" class="w-full p-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent">
              <option value="">全部分类</option>
              <option v-for="category in categories" :key="category.id" :value="category.id">
                {{ category.name }}
              </option>
            </select>
          </div>
          
          <!-- 难度筛选 -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">难度等级</label>
            <select v-model="selectedDifficulty" class="w-full p-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent">
              <option value="">全部难度</option>
              <option value="easy">简单</option>
              <option value="medium">中等</option>
              <option value="hard">困难</option>
            </select>
          </div>
          
          <!-- 状态筛选 -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">复习状态</label>
            <select v-model="selectedStatus" class="w-full p-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent">
              <option value="">全部状态</option>
              <option value="not_reviewed">未复习</option>
              <option value="reviewed">已复习</option>
              <option value="mastered">已掌握</option>
            </select>
          </div>
        </div>
      </div>

      <!-- 错题列表 -->
      <div v-if="filteredQuestions.length > 0" class="space-y-4">
        <div 
          v-for="mistake in filteredQuestions"
          :key="mistake.id"
          class="bg-white rounded-lg shadow-sm border hover:shadow-md transition-shadow duration-200"
        >
          <div class="p-4">
            <div class="flex items-start justify-between mb-3">
              <div class="flex items-center space-x-2">
                <span class="px-2 py-1 bg-blue-100 text-blue-800 text-xs rounded-full">
                    {{ getCategoryName(mistake.question.categoryId) }}
                  </span>
                <span class="px-2 py-1 text-xs rounded-full" :class="getDifficultyClass(mistake.question.difficulty)">
                  {{ getDifficultyText(mistake.question.difficulty) }}
                </span>
                <span class="px-2 py-1 text-xs rounded-full" :class="getStatusClass(mistake.status)">
                  {{ getStatusText(mistake.status) }}
                </span>
              </div>
              <div class="flex items-center space-x-2">
                <button 
                  @click="toggleMastered(mistake)"
                  :class="mistake.status === 'mastered' ? 'text-green-600' : 'text-gray-400'"
                  class="hover:text-green-600 transition-colors duration-200"
                >
                  <CheckCircle class="w-5 h-5" />
                </button>
                <button 
                  @click="removeMistake(mistake.id)"
                  class="text-gray-400 hover:text-red-600 transition-colors duration-200"
                >
                  <Trash2 class="w-5 h-5" />
                </button>
              </div>
            </div>
            
            <h4 class="font-medium text-gray-800 mb-2">{{ mistake.question.title }}</h4>
            <p class="text-gray-600 text-sm mb-3 line-clamp-2">{{ mistake.question.content }}</p>
            
            <div class="flex items-center justify-between text-sm text-gray-500 mb-4">
              <span>错误时间：{{ formatDate(mistake.createdAt) }}</span>
              <span v-if="mistake.lastReviewAt">最后复习：{{ formatDate(mistake.lastReviewAt) }}</span>
            </div>
            
            <div class="flex space-x-3">
              <button 
                @click="reviewQuestion(mistake)"
                class="flex-1 bg-blue-600 hover:bg-blue-700 text-white py-2 px-4 rounded-lg transition-colors duration-200 flex items-center justify-center space-x-2"
              >
                <RotateCcw class="w-4 h-4" />
                <span>重新练习</span>
              </button>
              <button 
                @click="viewQuestion(mistake.question.id)"
                class="flex-1 bg-gray-100 hover:bg-gray-200 text-gray-700 py-2 px-4 rounded-lg transition-colors duration-200 flex items-center justify-center space-x-2"
              >
                <Eye class="w-4 h-4" />
                <span>查看详情</span>
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- 空状态 -->
      <div v-else class="bg-white rounded-lg shadow-sm p-8 text-center">
        <BookmarkMinus class="w-16 h-16 mx-auto mb-4 text-gray-300" />
        <h3 class="text-lg font-medium text-gray-800 mb-2">
          {{ mistakeQuestions.length === 0 ? '暂无错题' : '没有符合条件的错题' }}
        </h3>
        <p class="text-gray-500 mb-4">
          {{ mistakeQuestions.length === 0 ? '继续努力，保持正确率！' : '试试调整筛选条件' }}
        </p>
        <button 
          @click="goToHome"
          class="bg-blue-600 hover:bg-blue-700 text-white py-2 px-4 rounded-lg transition-colors duration-200"
        >
          去刷题
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import {
  ArrowLeft, CheckCircle, Trash2, RotateCcw, Eye, BookmarkMinus
} from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'
import { questionService, type Question, type Category, type MistakeBook } from '@/services/question'

interface MistakeQuestion extends MistakeBook {
  question: Question
  status: 'not_reviewed' | 'reviewed' | 'mastered'
  reviewCount: number
  isMastered: boolean
  createdAt: string
  lastReviewAt?: string
}

const router = useRouter()
const authStore = useAuthStore()

const mistakeQuestions = ref<MistakeQuestion[]>([])
const categories = ref<Category[]>([])

// 筛选条件
const selectedCategory = ref('')
const selectedDifficulty = ref('')
const selectedStatus = ref('')

// 统计数据
const reviewedCount = computed(() => 
  mistakeQuestions.value.filter(m => m.status === 'reviewed').length
)

const masteredCount = computed(() => 
  mistakeQuestions.value.filter(m => m.status === 'mastered').length
)

// 筛选后的题目
const filteredQuestions = computed(() => {
  return mistakeQuestions.value.filter(mistake => {
    if (selectedCategory.value && mistake.question.categoryId !== parseInt(selectedCategory.value)) {
      return false
    }
    if (selectedDifficulty.value && mistake.question.difficulty !== selectedDifficulty.value) {
      return false
    }
    if (selectedStatus.value && mistake.status !== selectedStatus.value) {
      return false
    }
    return true
  })
})

// 样式方法
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

const getStatusClass = (status: string) => {
  switch (status) {
    case 'not_reviewed':
      return 'bg-red-100 text-red-800'
    case 'reviewed':
      return 'bg-orange-100 text-orange-800'
    case 'mastered':
      return 'bg-green-100 text-green-800'
    default:
      return 'bg-gray-100 text-gray-800'
  }
}

const getStatusText = (status: string) => {
  switch (status) {
    case 'not_reviewed':
      return '未复习'
    case 'reviewed':
      return '已复习'
    case 'mastered':
      return '已掌握'
    default:
      return '未知'
  }
}

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

    // 根据分类ID获取分类名称
    const getCategoryName = (categoryId: number) => {
      const category = categories.value.find(cat => cat.id === categoryId)
      return category ? category.name : '未知分类'
    }

// 操作方法
const goBack = () => {
  router.back()
}

const goToHome = () => {
  router.push('/')
}

const viewQuestion = (questionId: number) => {
  router.push(`/question/${questionId}`)
}

const reviewQuestion = (mistake: MistakeQuestion) => {
  // 更新复习状态
  mistake.status = 'reviewed'
  mistake.lastReviewAt = new Date().toISOString()
  mistake.reviewCount++
  
  // 跳转到题目页面
  router.push(`/question/${mistake.question.id}`)
}

const toggleMastered = async (mistake: MistakeQuestion) => {
  try {
    if (mistake.status === 'mastered') {
      // 重置掌握状态
      await questionService.resetMistakeStatus(mistake.question.id)
      mistake.status = 'reviewed'
    } else {
      // 标记为已掌握
      await questionService.markMistakeAsMastered(mistake.question.id)
      mistake.status = 'mastered'
      mistake.lastReviewAt = new Date().toISOString()
    }
  } catch (error) {
    console.error('更新掌握状态失败:', error)
    alert('操作失败，请重试')
  }
}

const removeMistake = async (mistakeId: number) => {
  if (!confirm('确定要从错题本中移除这道题吗？')) {
    return
  }
  
  try {
    // 找到对应的错题记录
    const mistake = mistakeQuestions.value.find(m => m.id === mistakeId)
    if (!mistake) {
      alert('错题记录不存在')
      return
    }
    
    // 调用API删除错题记录，传递错题本记录ID而不是题目ID
    await questionService.removeFromMistakeBook(mistakeId)
    
    // 从本地列表中移除
    const index = mistakeQuestions.value.findIndex(m => m.id === mistakeId)
    if (index > -1) {
      mistakeQuestions.value.splice(index, 1)
    }
  } catch (error) {
    console.error('移除错题失败:', error)
    alert('移除失败，请重试')
  }
}

const clearAllMistakes = async () => {
  if (!confirm('确定要清空所有错题吗？此操作不可恢复。')) {
    return
  }
  
  try {
    const response = await questionService.clearMistakeBook()
    if (response.success) {
      mistakeQuestions.value = []
      alert(`错题本已清空，共删除 ${response.data.deleted_count} 道题目`)
    }
  } catch (error) {
    console.error('清空错题本失败:', error)
    alert('清空失败，请重试')
  }
}

const resetFilters = () => {
  selectedCategory.value = ''
  selectedDifficulty.value = ''
  selectedStatus.value = ''
}

// 加载数据
const loadMistakeQuestions = async () => {
  try {
    const response = await questionService.getMistakeBook({
      page: 1,
      pageSize: 100 // 加载所有错题
    })
    
    if (response.success) {
      // 转换数据格式，添加状态字段
      mistakeQuestions.value = response.data.items.map(mistake => {
        // 根据isMastered字段确定状态
        let status: 'not_reviewed' | 'reviewed' | 'mastered' = 'not_reviewed'
        if (mistake.isMastered) {
          status = 'mastered'
        } else if (mistake.lastReviewAt) {
          status = 'reviewed'
        }
        
        return {
          ...mistake,
          status,
          createdAt: mistake.addedAt || new Date().toISOString()
        } as MistakeQuestion
      })
    }
  } catch (error) {
    console.error('加载错题本失败:', error)
    if (error.message && error.message.includes('登录已过期')) {
      alert('登录已过期，请重新登录')
      router.push('/auth/login')
    } else {
      alert('加载错题本失败，请重试')
    }
  }
}

const loadCategories = async () => {
  try {
    const response = await questionService.getCategories()
    if (response.success) {
      categories.value = response.data
    }
  } catch (error) {
    console.error('加载分类失败:', error)
  }
}

// 添加登录跳转方法
const goToLogin = () => {
  router.push('/auth/login')
}

onMounted(() => {
  loadMistakeQuestions()
  loadCategories()
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