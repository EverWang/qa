<template>
  <div class="min-h-screen bg-gray-50">
    <!-- 搜索头部 -->
    <div class="bg-white shadow-sm border-b sticky top-0 z-10">
      <div class="max-w-4xl mx-auto px-4 py-3">
        <div class="flex items-center space-x-3">
          <button @click="goBack" class="flex items-center text-gray-600 hover:text-gray-800">
            <ArrowLeft class="w-5 h-5" />
          </button>
          <div class="flex-1 relative">
            <input
              ref="searchInput"
              v-model="searchQuery"
              @keyup.enter="performSearch"
              type="text"
              placeholder="搜索题目、关键词..."
              class="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            >
            <Search class="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-gray-400" />
          </div>
          <button 
            @click="performSearch"
            :disabled="!searchQuery.trim()"
            class="bg-blue-600 hover:bg-blue-700 disabled:bg-gray-300 text-white px-4 py-2 rounded-lg transition-colors duration-200"
          >
            搜索
          </button>
        </div>
      </div>
    </div>

    <div class="max-w-4xl mx-auto p-4">
      <!-- 搜索历史 -->
      <div v-if="!hasSearched && searchHistory.length > 0" class="mb-6">
        <div class="flex items-center justify-between mb-3">
          <h3 class="text-lg font-semibold text-gray-800">搜索历史</h3>
          <button @click="clearHistory" class="text-sm text-gray-500 hover:text-gray-700">
            清空
          </button>
        </div>
        <div class="flex flex-wrap gap-2">
          <button
            v-for="(item, index) in searchHistory"
            :key="index"
            @click="searchFromHistory(item)"
            class="bg-gray-100 hover:bg-gray-200 text-gray-700 px-3 py-1 rounded-full text-sm transition-colors duration-200"
          >
            {{ item }}
          </button>
        </div>
      </div>

      <!-- 热门搜索 -->
      <div v-if="!hasSearched" class="mb-6">
        <h3 class="text-lg font-semibold text-gray-800 mb-3">热门搜索</h3>
        <div class="flex flex-wrap gap-2">
          <button
            v-for="(item, index) in hotSearches"
            :key="index"
            @click="searchFromHistory(item)"
            class="bg-blue-50 hover:bg-blue-100 text-blue-700 px-3 py-1 rounded-full text-sm transition-colors duration-200"
          >
            {{ item }}
          </button>
        </div>
      </div>

      <!-- 搜索结果 -->
      <div v-if="hasSearched">
        <!-- 搜索筛选 -->
        <div class="bg-white rounded-lg shadow-sm p-4 mb-4">
          <div class="flex items-center justify-between mb-3">
            <h3 class="font-medium text-gray-800">搜索结果</h3>
            <span class="text-sm text-gray-500">找到 {{ totalResults }} 个结果</span>
          </div>
          <div class="flex items-center space-x-4">
            <select v-model="searchType" @change="performSearch" class="text-sm border border-gray-300 rounded-lg px-3 py-1 focus:ring-2 focus:ring-blue-500 focus:border-transparent">
              <option value="all">全部</option>
              <option value="question">题目</option>
              <option value="category">分类</option>
            </select>
            <select v-model="difficultyFilter" @change="performSearch" class="text-sm border border-gray-300 rounded-lg px-3 py-1 focus:ring-2 focus:ring-blue-500 focus:border-transparent">
              <option value="">全部难度</option>
              <option value="easy">简单</option>
              <option value="medium">中等</option>
              <option value="hard">困难</option>
            </select>
            <select v-model="sortBy" @change="performSearch" class="text-sm border border-gray-300 rounded-lg px-3 py-1 focus:ring-2 focus:ring-blue-500 focus:border-transparent">
              <option value="relevance">相关度</option>
              <option value="latest">最新</option>
              <option value="difficulty">难度</option>
            </select>
          </div>
        </div>

        <!-- 分类结果 -->
        <div v-if="categoryResults.length > 0" class="mb-6">
          <h4 class="text-md font-medium text-gray-800 mb-3">相关分类</h4>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
            <div
              v-for="category in categoryResults"
              :key="category.id"
              @click="goToCategory(category.id)"
              class="bg-white rounded-lg p-3 shadow-sm hover:shadow-md transition-shadow duration-200 cursor-pointer border hover:border-blue-200"
            >
              <div class="flex items-center space-x-3">
                <div class="w-8 h-8 rounded-lg flex items-center justify-center" :style="{ backgroundColor: category.color + '20', color: category.color }">
                  <component :is="category.icon" class="w-4 h-4" />
                </div>
                <div class="flex-1">
                  <h5 class="font-medium text-gray-800">{{ category.name }}</h5>
                  <p class="text-xs text-gray-500">{{ category.questionCount }} 道题目</p>
                </div>
                <ChevronRight class="w-4 h-4 text-gray-400" />
              </div>
            </div>
          </div>
        </div>

        <!-- 题目结果 -->
        <div v-if="questionResults.length > 0">
          <h4 class="text-md font-medium text-gray-800 mb-3">相关题目</h4>
          <div class="space-y-3">
            <div
              v-for="question in questionResults"
              :key="question.id"
              @click="goToQuestion(question.id)"
              class="bg-white rounded-lg p-4 shadow-sm hover:shadow-md transition-shadow duration-200 cursor-pointer border hover:border-blue-200"
            >
              <div class="flex items-start justify-between mb-2">
                <div class="flex items-center space-x-2">
                  <span class="px-2 py-1 text-xs rounded-full" :class="getDifficultyClass(question.difficulty)">
                    {{ getDifficultyText(question.difficulty) }}
                  </span>
                  <span class="px-2 py-1 text-xs rounded-full bg-blue-100 text-blue-800">
                    {{ question.categoryName }}
                  </span>
                  <span v-if="question.isAnswered" class="px-2 py-1 text-xs rounded-full" :class="question.isCorrect ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'">
                    {{ question.isCorrect ? '已答对' : '答错了' }}
                  </span>
                </div>
                <ChevronRight class="w-4 h-4 text-gray-400" />
              </div>
              <h5 class="font-medium text-gray-800 mb-2" v-html="highlightSearchTerm(question.title)"></h5>
              <p class="text-sm text-gray-600 line-clamp-2 mb-2" v-html="highlightSearchTerm(question.content)"></p>
              <div class="flex items-center justify-between text-xs text-gray-500">
                <span>题目 #{{ question.id }}</span>
                <span>{{ formatDate(question.createdAt) }}</span>
              </div>
            </div>
          </div>

          <!-- 加载更多 -->
          <div v-if="hasMore" class="text-center mt-6">
            <button 
              @click="loadMore"
              :disabled="isLoading"
              class="bg-gray-100 hover:bg-gray-200 disabled:bg-gray-50 text-gray-700 py-2 px-6 rounded-lg transition-colors duration-200"
            >
              {{ isLoading ? '加载中...' : '加载更多' }}
            </button>
          </div>
        </div>

        <!-- 无结果 -->
        <div v-if="hasSearched && questionResults.length === 0 && categoryResults.length === 0" class="bg-white rounded-lg shadow-sm p-8 text-center">
          <Search class="w-16 h-16 mx-auto mb-4 text-gray-300" />
          <h3 class="text-lg font-medium text-gray-800 mb-2">未找到相关结果</h3>
          <p class="text-gray-500 mb-4">试试其他关键词或浏览分类</p>
          <div class="flex justify-center space-x-3">
            <button 
              @click="clearSearch"
              class="bg-gray-100 hover:bg-gray-200 text-gray-700 py-2 px-4 rounded-lg transition-colors duration-200"
            >
              重新搜索
            </button>
            <button 
              @click="goToHome"
              class="bg-blue-600 hover:bg-blue-700 text-white py-2 px-4 rounded-lg transition-colors duration-200"
            >
              返回首页
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import {
  ArrowLeft, Search, ChevronRight, Gavel, Stethoscope,
  HardHat, GraduationCap, BookOpen, Calculator, FileText
} from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'

interface Category {
  id: number
  name: string
  description: string
  questionCount: number
  color: string
  icon: any
}

interface Question {
  id: number
  title: string
  content: string
  difficulty: 'easy' | 'medium' | 'hard'
  categoryName: string
  isAnswered: boolean
  isCorrect?: boolean
  createdAt: string
}

const router = useRouter()
const authStore = useAuthStore()
const searchInput = ref<HTMLInputElement>()

const searchQuery = ref('')
const hasSearched = ref(false)
const isLoading = ref(false)
const hasMore = ref(true)
const page = ref(1)
const totalResults = ref(0)

const searchType = ref('all')
const difficultyFilter = ref('')
const sortBy = ref('relevance')

const questionResults = ref<Question[]>([])
const categoryResults = ref<Category[]>([])

// 搜索历史
const searchHistory = ref<string[]>([
  '法律基础', '医学常识', '工程管理', '考试真题'
])

// 热门搜索
const hotSearches = ref<string[]>([
  '法考', '医考', '建造师', '公务员', '教师资格', '会计师', '英语四级', '计算机二级'
])

// 所有分类数据
const allCategories = ref<Category[]>([
  { id: 1, name: '法考题', description: '法律职业资格考试题目', questionCount: 1200, color: '#3B82F6', icon: Gavel },
  { id: 2, name: '医考题', description: '医师资格考试题目', questionCount: 800, color: '#EF4444', icon: Stethoscope },
  { id: 3, name: '工程类', description: '各类工程考试题目', questionCount: 1500, color: '#F59E0B', icon: HardHat },
  { id: 4, name: '其他考试', description: '其他各类考试题目', questionCount: 1500, color: '#8B5CF6', icon: GraduationCap },
  { id: 11, name: '法律基础', description: '法律基础知识', questionCount: 400, color: '#3B82F6', icon: BookOpen },
  { id: 12, name: '法律条款', description: '法律条款解读', questionCount: 500, color: '#3B82F6', icon: FileText },
  { id: 21, name: '临床医学', description: '临床医学知识', questionCount: 400, color: '#EF4444', icon: Stethoscope },
  { id: 22, name: '基础医学', description: '基础医学理论', questionCount: 400, color: '#EF4444', icon: BookOpen },
  { id: 31, name: '道路工程', description: '道路工程设计与施工', questionCount: 750, color: '#F59E0B', icon: HardHat },
  { id: 32, name: '桥隧工程', description: '桥梁隧道工程', questionCount: 750, color: '#F59E0B', icon: Calculator }
])

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

// 高亮搜索关键词
const highlightSearchTerm = (text: string) => {
  if (!searchQuery.value.trim()) return text
  
  const regex = new RegExp(`(${searchQuery.value.trim()})`, 'gi')
  return text.replace(regex, '<mark class="bg-yellow-200">$1</mark>')
}

// 格式化日期
const formatDate = (dateString: string) => {
  const date = new Date(dateString)
  return date.toLocaleDateString()
}

// 导航方法
const goBack = () => {
  router.back()
}

const goToHome = () => {
  router.push('/')
}

const goToCategory = (categoryId: number) => {
  router.push(`/category/${categoryId}`)
}

const goToQuestion = (questionId: number) => {
  router.push(`/question/${questionId}`)
}

// 搜索方法
const performSearch = async () => {
  const query = searchQuery.value.trim()
  if (!query) return
  
  try {
    isLoading.value = true
    hasSearched.value = true
    page.value = 1
    hasMore.value = true
    
    // 添加到搜索历史
    if (!searchHistory.value.includes(query)) {
      searchHistory.value.unshift(query)
      if (searchHistory.value.length > 10) {
        searchHistory.value = searchHistory.value.slice(0, 10)
      }
    }
    
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 500))
    
    // 搜索分类
    if (searchType.value === 'all' || searchType.value === 'category') {
      categoryResults.value = allCategories.value.filter(category => 
        category.name.toLowerCase().includes(query.toLowerCase()) ||
        category.description.toLowerCase().includes(query.toLowerCase())
      )
    } else {
      categoryResults.value = []
    }
    
    // 搜索题目
    if (searchType.value === 'all' || searchType.value === 'question') {
      questionResults.value = generateMockQuestions(query, 10)
    } else {
      questionResults.value = []
    }
    
    totalResults.value = categoryResults.value.length + questionResults.value.length
    
  } catch (error) {
    console.error('搜索失败:', error)
  } finally {
    isLoading.value = false
  }
}

// 从历史搜索
const searchFromHistory = (query: string) => {
  searchQuery.value = query
  performSearch()
}

// 清空搜索
const clearSearch = () => {
  searchQuery.value = ''
  hasSearched.value = false
  questionResults.value = []
  categoryResults.value = []
  totalResults.value = 0
  nextTick(() => {
    searchInput.value?.focus()
  })
}

// 清空历史
const clearHistory = () => {
  searchHistory.value = []
}

// 加载更多
const loadMore = async () => {
  if (isLoading.value) return
  
  try {
    isLoading.value = true
    page.value++
    
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 500))
    
    // 模拟加载更多题目
    const moreQuestions = generateMockQuestions(searchQuery.value, 10, page.value)
    questionResults.value.push(...moreQuestions)
    
    // 模拟没有更多数据
    if (page.value >= 3) {
      hasMore.value = false
    }
  } catch (error) {
    console.error('加载更多失败:', error)
  } finally {
    isLoading.value = false
  }
}

// 生成模拟题目数据
const generateMockQuestions = (query: string, count: number, pageNum: number = 1): Question[] => {
  const difficulties: ('easy' | 'medium' | 'hard')[] = ['easy', 'medium', 'hard']
  const categories = ['法考题', '医考题', '工程类', '其他考试']
  const questions: Question[] = []
  
  for (let i = 0; i < count; i++) {
    const id = (pageNum - 1) * count + i + 1
    const category = categories[Math.floor(Math.random() * categories.length)]
    const difficulty = difficulties[Math.floor(Math.random() * difficulties.length)]
    
    // 根据搜索词生成相关题目
    let title = `关于${query}的第${id}题`
    let content = `这是一道关于${query}的${getDifficultyText(difficulty)}题目，请仔细阅读题目内容并选择正确答案。题目涉及${category}相关知识点。`
    
    // 应用难度筛选
    if (difficultyFilter.value && difficulty !== difficultyFilter.value) {
      continue
    }
    
    questions.push({
      id,
      title,
      content,
      difficulty,
      categoryName: category,
      isAnswered: Math.random() > 0.7,
      isCorrect: Math.random() > 0.3,
      createdAt: new Date(Date.now() - Math.random() * 30 * 24 * 60 * 60 * 1000).toISOString()
    })
  }
  
  // 排序
  if (sortBy.value === 'latest') {
    questions.sort((a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime())
  } else if (sortBy.value === 'difficulty') {
    const difficultyOrder = { 'easy': 1, 'medium': 2, 'hard': 3 }
    questions.sort((a, b) => difficultyOrder[a.difficulty] - difficultyOrder[b.difficulty])
  }
  
  return questions
}

onMounted(() => {
  // 自动聚焦搜索框
  nextTick(() => {
    searchInput.value?.focus()
  })
})
</script>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

:deep(mark) {
  background-color: #fef08a;
  padding: 0 2px;
  border-radius: 2px;
}
</style>