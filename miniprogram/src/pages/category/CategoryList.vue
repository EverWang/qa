<template>
  <div class="min-h-screen bg-gray-50">
    <!-- 顶部导航 -->
    <div class="bg-white shadow-sm border-b">
      <div class="max-w-4xl mx-auto px-4 py-3 flex items-center justify-between">
        <button @click="goBack" class="flex items-center text-gray-600 hover:text-gray-800">
          <ArrowLeft class="w-5 h-5 mr-1" />
          <span>返回</span>
        </button>
        <h1 class="text-lg font-medium text-gray-800">
          {{ currentCategory ? currentCategory.name : '题目分类' }}
        </h1>
        <button @click="goToSearch" class="flex items-center text-gray-600 hover:text-gray-800">
          <Search class="w-5 h-5" />
        </button>
      </div>
    </div>

    <div class="max-w-4xl mx-auto p-4">
      <!-- 面包屑导航 -->
      <div v-if="breadcrumbs.length > 0" class="mb-4">
        <nav class="flex items-center space-x-2 text-sm text-gray-600">
          <button @click="goToCategory()" class="hover:text-blue-600">全部分类</button>
          <template v-for="(crumb, index) in breadcrumbs" :key="crumb.id">
            <ChevronRight class="w-4 h-4" />
            <button 
              @click="goToCategory(crumb.id)"
              :class="index === breadcrumbs.length - 1 ? 'text-blue-600 font-medium' : 'hover:text-blue-600'"
            >
              {{ crumb.name }}
            </button>
          </template>
        </nav>
      </div>

      <!-- 分类统计 -->
      <div v-if="currentCategory" class="bg-white rounded-lg shadow-sm p-6 mb-6">
        <div class="flex items-center space-x-4 mb-4">
          <div class="w-12 h-12 rounded-lg flex items-center justify-center" :style="{ backgroundColor: currentCategory.color + '20', color: currentCategory.color }">
            <component :is="currentCategory.icon" class="w-6 h-6" />
          </div>
          <div>
            <h2 class="text-xl font-semibold text-gray-800">{{ currentCategory.name }}</h2>
            <p class="text-gray-600">{{ currentCategory.description }}</p>
          </div>
        </div>
        <div class="grid grid-cols-3 gap-4 text-center">
          <div>
            <div class="text-2xl font-bold text-blue-600">{{ currentCategory.questionCount }}</div>
            <div class="text-sm text-gray-500">题目总数</div>
          </div>
          <div>
            <div class="text-2xl font-bold text-green-600">{{ userProgress.answered }}</div>
            <div class="text-sm text-gray-500">已答题目</div>
          </div>
          <div>
            <div class="text-2xl font-bold text-orange-600">{{ userProgress.correctRate }}%</div>
            <div class="text-sm text-gray-500">正确率</div>
          </div>
        </div>
      </div>

      <!-- 子分类 -->
      <div v-if="subCategories.length > 0" class="mb-6">
        <h3 class="text-lg font-semibold text-gray-800 mb-4">子分类</h3>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div
            v-for="category in subCategories"
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
                <p class="text-sm text-gray-500">{{ category.questionCount }} 道题目</p>
              </div>
              <ChevronRight class="w-5 h-5 text-gray-400" />
            </div>
          </div>
        </div>
      </div>

      <!-- 题目列表 -->
      <div class="mb-6">
        <div class="flex items-center justify-between mb-4">
          <h3 class="text-lg font-semibold text-gray-800">题目列表</h3>
          <div class="flex items-center space-x-2">
            <!-- 分类筛选器 -->
            <select v-model="selectedCategoryFilter" @change="onCategoryFilterChange" class="text-sm border border-gray-300 rounded-lg px-3 py-1 focus:ring-2 focus:ring-blue-500 focus:border-transparent">
              <option value="">全部分类</option>
              <option v-for="category in allCategories" :key="category.id" :value="category.id">
                {{ category.parentId ? '　　' + category.name : category.name }}
              </option>
            </select>
            <select v-model="sortBy" class="text-sm border border-gray-300 rounded-lg px-3 py-1 focus:ring-2 focus:ring-blue-500 focus:border-transparent">
              <option value="default">默认排序</option>
              <option value="difficulty">按难度</option>
              <option value="latest">最新题目</option>
            </select>
            <select v-model="difficultyFilter" class="text-sm border border-gray-300 rounded-lg px-3 py-1 focus:ring-2 focus:ring-blue-500 focus:border-transparent">
              <option value="">全部难度</option>
              <option value="easy">简单</option>
              <option value="medium">中等</option>
              <option value="hard">困难</option>
            </select>
          </div>
        </div>

        <!-- 题目卡片 -->
        <div v-if="filteredQuestions.length > 0" class="space-y-4">
          <div
            v-for="question in filteredQuestions"
            :key="question.id"
            @click="goToQuestion(question.id)"
            class="bg-white rounded-lg p-4 shadow-sm hover:shadow-md transition-shadow duration-200 cursor-pointer border hover:border-blue-200"
          >
            <div class="flex items-start justify-between mb-3">
              <div class="flex items-center space-x-2">
                <span class="px-2 py-1 text-xs rounded-full" :class="getDifficultyClass(question.difficulty)">
                  {{ getDifficultyText(question.difficulty) }}
                </span>
                <span v-if="question.isAnswered" class="px-2 py-1 text-xs rounded-full" :class="question.isCorrect ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'">
                  {{ question.isCorrect ? '已答对' : '答错了' }}
                </span>
              </div>
              <ChevronRight class="w-5 h-5 text-gray-400" />
            </div>
            <h4 class="font-medium text-gray-800 mb-2">{{ question.title }}</h4>
            <p class="text-sm text-gray-600 line-clamp-2 mb-3">{{ question.content }}</p>
            <div class="flex items-center justify-between text-xs text-gray-500">
              <span>题目 #{{ question.id }}</span>
              <span>{{ formatDate(question.createdAt) }}</span>
            </div>
          </div>
        </div>

        <!-- 空状态 -->
        <div v-else class="bg-white rounded-lg shadow-sm p-8 text-center">
          <FileText class="w-16 h-16 mx-auto mb-4 text-gray-300" />
          <h3 class="text-lg font-medium text-gray-800 mb-2">暂无题目</h3>
          <p class="text-gray-500 mb-4">该分类下暂时没有题目</p>
          <button 
            @click="goToHome"
            class="bg-blue-600 hover:bg-blue-700 text-white py-2 px-4 rounded-lg transition-colors duration-200"
          >
            返回首页
          </button>
        </div>
      </div>

      <!-- 加载更多 -->
      <div v-if="hasMore && filteredQuestions.length > 0" class="text-center">
        <button 
          @click="loadMore"
          :disabled="isLoading"
          class="bg-gray-100 hover:bg-gray-200 disabled:bg-gray-50 text-gray-700 py-2 px-6 rounded-lg transition-colors duration-200"
        >
          {{ isLoading ? '加载中...' : '加载更多' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import {
  ArrowLeft, Search, ChevronRight, FileText, Gavel, Stethoscope,
  HardHat, GraduationCap, BookOpen, Calculator
} from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'
import { QuestionService, type Category, type Question } from '@/services/question'

interface CategoryWithIcon extends Category {
  color: string
  icon: any
}

interface UserProgress {
  answered: number
  correctRate: number
}

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const categoryId = computed(() => route.params.id ? parseInt(route.params.id as string) : null)

const currentCategory = ref<CategoryWithIcon | null>(null)
const subCategories = ref<CategoryWithIcon[]>([])
const questions = ref<Question[]>([])
const breadcrumbs = ref<CategoryWithIcon[]>([])
const userProgress = ref<UserProgress>({ answered: 0, correctRate: 0 })

const sortBy = ref('default')
const difficultyFilter = ref('')
const selectedCategoryFilter = ref('')
const isLoading = ref(false)
const hasMore = ref(true)
const page = ref(1)

// 分类数据
const allCategories = ref<CategoryWithIcon[]>([])

// 为分类添加图标和颜色的辅助函数
const addCategoryIcon = (category: Category): CategoryWithIcon => {
  const iconMap: Record<string, any> = {
    '法考题': Gavel,
    '医考题': Stethoscope,
    '工程类': HardHat,
    '其他考试': GraduationCap,
    '法律基础': BookOpen,
    '法律条款': FileText,
    '法律解释': Calculator,
    '临床医学': Stethoscope,
    '基础医学': BookOpen,
    '道路工程': HardHat,
    '桥隧工程': Calculator
  }
  
  const colorMap: Record<string, string> = {
    '法考题': '#3B82F6',
    '医考题': '#EF4444',
    '工程类': '#F59E0B',
    '其他考试': '#8B5CF6',
    '法律基础': '#3B82F6',
    '法律条款': '#3B82F6',
    '法律解释': '#3B82F6',
    '临床医学': '#EF4444',
    '基础医学': '#EF4444',
    '道路工程': '#F59E0B',
    '桥隧工程': '#F59E0B'
  }
  
  return {
    ...category,
    icon: iconMap[category.name] || BookOpen,
    color: colorMap[category.name] || '#6B7280'
  }
}

// 筛选后的题目
const filteredQuestions = computed(() => {
  let filtered = questions.value
  
  if (difficultyFilter.value) {
    filtered = filtered.filter(q => q.difficulty === difficultyFilter.value)
  }
  
  // 排序
  if (sortBy.value === 'difficulty') {
    const difficultyOrder = { 'easy': 1, 'medium': 2, 'hard': 3 }
    filtered = [...filtered].sort((a, b) => difficultyOrder[a.difficulty] - difficultyOrder[b.difficulty])
  } else if (sortBy.value === 'latest') {
    filtered = [...filtered].sort((a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime())
  }
  
  return filtered
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

// 格式化日期
const formatDate = (dateString: string) => {
  const date = new Date(dateString)
  return date.toLocaleDateString()
}

// 分类筛选处理
const onCategoryFilterChange = async () => {
  try {
    isLoading.value = true
    page.value = 1
    
    if (selectedCategoryFilter.value) {
      // 如果选择了特定分类，加载该分类下的题目
      const categoryIdToFilter = parseInt(selectedCategoryFilter.value)
      console.log('CategoryList.onCategoryFilterChange 筛选分类:', categoryIdToFilter)
      
      const questionsResponse = await QuestionService.getQuestions({
        page: 1,
        pageSize: 10,
        categoryId: categoryIdToFilter
      })
      
      if (questionsResponse.success) {
        questions.value = questionsResponse.data.items || []
        hasMore.value = questionsResponse.data.page < questionsResponse.data.totalPages
        console.log('CategoryList.onCategoryFilterChange 筛选后的题目:', questions.value)
      } else {
        console.error('CategoryList.onCategoryFilterChange 筛选题目失败:', questionsResponse.message)
        questions.value = []
        hasMore.value = false
      }
    } else {
      // 如果选择了"全部分类"，重新加载当前分类的所有题目
      await loadCategoryData(categoryId.value || undefined)
    }
  } catch (error) {
    console.error('CategoryList.onCategoryFilterChange 筛选失败:', error)
    questions.value = []
    hasMore.value = false
  } finally {
    isLoading.value = false
  }
}

// 导航方法
const goBack = () => {
  router.back()
}

const goToHome = () => {
  router.push('/')
}

const goToSearch = () => {
  router.push('/search')
}

const goToCategory = (id?: number) => {
  if (id) {
    router.push(`/category/${id}`)
  } else {
    router.push('/category')
  }
}

const goToQuestion = (questionId: number) => {
  router.push(`/question/${questionId}`)
}

// 加载更多
const loadMore = async () => {
  if (isLoading.value || !hasMore.value) return
  
  try {
    isLoading.value = true
    page.value++
    
    let moreQuestions: Question[] = []
    if (categoryId.value) {
      // 加载特定分类的题目
      const response = await QuestionService.getQuestions({
          page: page.value,
          pageSize: 10,
          categoryId: categoryId.value
        })
      if (response.success) {
          moreQuestions = response.data.items
          hasMore.value = response.data.page < response.data.totalPages
        }
    } else {
      // 加载所有题目
      const response = await QuestionService.getQuestions({
      page: page.value + 1,
      pageSize: 10
    })
    if (response.success) {
       moreQuestions = response.data.items
       hasMore.value = response.data.page < response.data.totalPages
     }
    }
    
    questions.value.push(...moreQuestions)
  } catch (error) {
    console.error('加载更多失败:', error)
  } finally {
    isLoading.value = false
  }
}

// 构建面包屑导航
const buildBreadcrumbs = (categoryId: number) => {
  const breadcrumbs: CategoryWithIcon[] = []
  let current = allCategories.value.find(c => c.id === categoryId)
  
  while (current) {
    breadcrumbs.unshift(current)
    if (current.parentId) {
      current = allCategories.value.find(c => c.id === current!.parentId)
    } else {
      break
    }
  }
  
  return breadcrumbs
}

// 加载分类数据
const loadCategoryData = async (id?: number) => {
  try {
    console.log('CategoryList.loadCategoryData 开始加载，分类ID:', id)
    isLoading.value = true
    
    // 首先加载所有分类数据
    if (allCategories.value.length === 0) {
      console.log('CategoryList.loadCategoryData 加载所有分类')
      const categoriesResponse = await QuestionService.getCategories()
      console.log('CategoryList.loadCategoryData 分类响应:', categoriesResponse)
      if (categoriesResponse.success) {
        allCategories.value = (categoriesResponse.data || []).map(addCategoryIcon)
        console.log('CategoryList.loadCategoryData 处理后的分类:', allCategories.value)
      }
    }
    
    if (id) {
      console.log('CategoryList.loadCategoryData 加载特定分类:', id)
      // 加载特定分类 - 从已加载的分类中查找
      const foundCategory = allCategories.value.find(c => c.id === id)
      if (foundCategory) {
        currentCategory.value = foundCategory
        console.log('CategoryList.loadCategoryData 找到当前分类:', currentCategory.value)
      } else {
        console.log('CategoryList.loadCategoryData 未找到分类，尝试从API获取')
        try {
           const categoryResponse = await QuestionService.getCategory(id)
           if (categoryResponse.success) {
             currentCategory.value = addCategoryIcon(categoryResponse.data)
             console.log('CategoryList.loadCategoryData API获取的分类:', currentCategory.value)
           }
         } catch (error) {
           console.error('CategoryList.loadCategoryData 获取分类详情失败:', error)
           currentCategory.value = null
         }
      }
      
      // 获取子分类
      subCategories.value = allCategories.value.filter(c => c.parentId === id)
      console.log('CategoryList.loadCategoryData 子分类:', subCategories.value)
      breadcrumbs.value = buildBreadcrumbs(id)
      console.log('CategoryList.loadCategoryData 面包屑:', breadcrumbs.value)
      
      // 加载用户进度
      if (authStore.isLoggedIn && !authStore.isGuest) {
        try {
          // 暂时使用默认值，后续可以实现具体的进度API
          userProgress.value = {
            answered: 0,
            correctRate: 0
          }
        } catch (error) {
          console.error('获取用户进度失败:', error)
          // 如果获取失败，设置默认值
          userProgress.value = {
            answered: 0,
            correctRate: 0
          }
        }
      }
      
      // 加载该分类下的题目
      console.log('CategoryList.loadCategoryData 加载分类题目，分类ID:', id)
      const questionsResponse = await QuestionService.getQuestions({
        page: 1,
        pageSize: 10,
        categoryId: id
      })
      console.log('CategoryList.loadCategoryData 题目响应:', questionsResponse)
      if (questionsResponse.success) {
        questions.value = questionsResponse.data.items || []
        hasMore.value = questionsResponse.data.page < questionsResponse.data.totalPages
        console.log('CategoryList.loadCategoryData 加载的题目:', questions.value)
        console.log('CategoryList.loadCategoryData 是否有更多:', hasMore.value)
      } else {
        console.error('CategoryList.loadCategoryData 加载题目失败:', questionsResponse.message || questionsResponse.success)
        questions.value = []
        hasMore.value = false
      }
    } else {
      console.log('CategoryList.loadCategoryData 加载根分类')
      // 加载根分类
      currentCategory.value = null
      subCategories.value = allCategories.value.filter(c => !c.parentId)
      breadcrumbs.value = []
      console.log('CategoryList.loadCategoryData 根分类子分类:', subCategories.value)
      
      // 加载所有题目
      console.log('CategoryList.loadCategoryData 加载所有题目')
      const response = await QuestionService.getQuestions({
        page: 1,
        pageSize: 10
      })
      console.log('CategoryList.loadCategoryData 所有题目响应:', response)
      if (response.success) {
        questions.value = response.data.items || []
        hasMore.value = response.data.page < response.data.totalPages
        console.log('CategoryList.loadCategoryData 加载的所有题目:', questions.value)
      } else {
        console.error('CategoryList.loadCategoryData 加载所有题目失败:', response.message || response.success)
        questions.value = []
        hasMore.value = false
      }
    }
    
    // 重置分页
    page.value = 1
    
  } catch (error) {
    console.error('CategoryList.loadCategoryData 加载分类数据失败:', error)
    questions.value = []
    hasMore.value = false
  } finally {
    isLoading.value = false
    console.log('CategoryList.loadCategoryData 加载完成')
  }
}

// 监听路由变化
watch(() => route.params.id, (newId) => {
  const id = newId ? parseInt(newId as string) : undefined
  loadCategoryData(id)
}, { immediate: true })

onMounted(() => {
  loadCategoryData(categoryId.value || undefined)
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