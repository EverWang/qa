<template>
  <div class="min-h-screen bg-gray-50">
    <!-- 顶部导航栏 -->
    <div class="bg-white shadow-sm border-b sticky top-0 z-10">
      <div class="max-w-4xl mx-auto px-4 py-3">
        <div class="flex items-center justify-between">
          <div class="flex items-center space-x-3">
            <button @click="goBack" class="flex items-center text-gray-600 hover:text-gray-800">
              <ArrowLeft class="w-5 h-5" />
            </button>
            <h1 class="text-lg font-semibold text-gray-800">题目详情</h1>
          </div>
          <div class="flex items-center space-x-2">
            <button 
              @click="shareQuestion" 
              class="p-2 text-gray-600 hover:text-gray-800 hover:bg-gray-100 rounded-lg"
            >
              <Share2 class="w-5 h-5" />
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 题目内容 -->
    <div class="max-w-4xl mx-auto p-4" v-if="question">
      <!-- 题目信息 -->
      <div class="bg-white rounded-lg shadow-sm p-6 mb-4">
        <div class="flex items-center space-x-2 mb-4">
          <span class="px-3 py-1 bg-blue-100 text-blue-800 text-sm rounded-full">
            {{ question.categoryName }}
          </span>
          <span class="px-3 py-1 text-sm rounded-full" :class="difficultyClass">
            {{ difficultyText }}
          </span>
        </div>
        
        <h2 class="text-xl font-semibold text-gray-800 mb-4">{{ question.title }}</h2>
        <div class="text-gray-700 mb-6" v-html="question.content"></div>
        
        <!-- 选项 -->
        <div class="space-y-3" v-if="question.options && question.options.length > 0">
          <div 
            v-for="(option, index) in question.options" 
            :key="index"
            @click="selectOption(index)"
            class="flex items-center p-4 border rounded-lg cursor-pointer transition-colors"
            :class="getOptionClass(index)"
          >
            <div class="flex items-center justify-center w-8 h-8 rounded-full mr-3" :class="getOptionIndicatorClass(index)">
              <span class="text-sm font-medium">{{ String.fromCharCode(65 + index) }}</span>
            </div>
            <div class="flex-1 text-gray-800">{{ option.content || option }}</div>
          </div>
        </div>
        
        <!-- 提交按钮 -->
        <div class="mt-6 flex justify-center" v-if="!isAnswered">
          <button 
            @click="submitAnswer"
            :disabled="selectedOption === null"
            class="px-8 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:bg-gray-300 disabled:cursor-not-allowed transition-colors"
          >
            提交答案
          </button>
        </div>
      </div>
      
      <!-- 答题结果 -->
      <div v-if="showResult" class="bg-white rounded-lg shadow-sm p-6 mb-4">
        <div class="flex items-center space-x-2 mb-4">
          <CheckCircle v-if="isCorrect" class="w-6 h-6 text-green-600" />
          <XCircle v-else class="w-6 h-6 text-red-600" />
          <span class="text-lg font-semibold" :class="isCorrect ? 'text-green-600' : 'text-red-600'">
            {{ isCorrect ? '回答正确！' : '回答错误' }}
          </span>
        </div>
        
        <div v-if="!isCorrect" class="mb-4">
          <p class="text-gray-600 mb-2">正确答案：</p>
          <p class="text-gray-800 font-medium">{{ getCorrectAnswerText() }}</p>
        </div>
        
        <div v-if="question.explanation" class="mb-4">
          <p class="text-gray-600 mb-2">题目解析：</p>
          <div class="text-gray-800" v-html="question.explanation"></div>
        </div>
        
        <div class="flex space-x-3">
          <button 
            @click="nextQuestion"
            class="flex-1 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
          >
            下一题
          </button>
          <button 
            v-if="!isCorrect"
            @click="addToMistakes"
            class="px-4 py-2 bg-orange-600 text-white rounded-lg hover:bg-orange-700 transition-colors"
          >
            加入错题本
          </button>
        </div>
      </div>
    </div>
    
    <!-- 加载状态 -->
    <div v-else class="flex items-center justify-center min-h-screen">
      <div class="text-center">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
        <p class="text-gray-600">加载中...</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { 
  ArrowLeft, Share2, CheckCircle, XCircle 
} from 'lucide-vue-next'
import { QuestionService, type Question } from '@/services/question'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

// 响应式数据
const question = ref<Question | null>(null)
const selectedOption = ref<number | null>(null)
const isAnswered = ref(false)
const showResult = ref(false)
const isCorrect = ref(false)
const startTime = ref<number>(Date.now())

// 计算属性
const difficultyClass = computed(() => {
  if (!question.value) return ''
  switch (question.value.difficulty) {
    case 'easy':
      return 'bg-green-100 text-green-800'
    case 'medium':
      return 'bg-yellow-100 text-yellow-800'
    case 'hard':
      return 'bg-red-100 text-red-800'
    default:
      return 'bg-gray-100 text-gray-800'
  }
})

const difficultyText = computed(() => {
  if (!question.value) return ''
  switch (question.value.difficulty) {
    case 'easy':
      return '简单'
    case 'medium':
      return '中等'
    case 'hard':
      return '困难'
    default:
      return '未知'
  }
})

// 方法
const getOptionClass = (index: number) => {
  if (!showResult.value) {
    return selectedOption.value === index 
      ? 'border-blue-500 bg-blue-50' 
      : 'border-gray-200 hover:border-gray-300'
  }
  
  // 显示结果时的样式
  if (question.value?.options && Array.isArray(question.value.options)) {
    const option = question.value.options[index]
    if (option && typeof option === 'object' && option.isCorrect) {
      return 'border-green-500 bg-green-50'
    }
  }
  
  if (selectedOption.value === index && !isCorrect.value) {
    return 'border-red-500 bg-red-50'
  }
  
  return 'border-gray-200'
}

const getOptionIndicatorClass = (index: number) => {
  if (!showResult.value) {
    return selectedOption.value === index 
      ? 'bg-blue-600 text-white' 
      : 'bg-gray-100 text-gray-600'
  }
  
  // 显示结果时的样式
  if (question.value?.options && Array.isArray(question.value.options)) {
    const option = question.value.options[index]
    if (option && typeof option === 'object' && option.isCorrect) {
      return 'bg-green-600 text-white'
    }
  }
  
  if (selectedOption.value === index && !isCorrect.value) {
    return 'bg-red-600 text-white'
  }
  
  return 'bg-gray-100 text-gray-600'
}

const selectOption = (index: number) => {
  if (!isAnswered.value) {
    selectedOption.value = index
  }
}

const submitAnswer = async () => {
  if (selectedOption.value === null || !question.value) return
  
  try {
    const timeSpent = Math.floor((Date.now() - startTime.value) / 1000)
    // 后端期望的是数字索引，不是字符串
    const userAnswer = selectedOption.value
    
    const response = await QuestionService.submitAnswer({
      questionId: question.value.id,
      userAnswer,
      timeSpent
    })
    
    if (response.success) {
      isAnswered.value = true
      showResult.value = true
      isCorrect.value = response.data?.isCorrect || false
      
      if (isCorrect.value) {
        alert('回答正确！')
      } else {
        alert('回答错误，请查看解析')
      }
    }
  } catch (error) {
    console.error('提交答案失败:', error)
    alert('提交答案失败，请重试')
  }
}

const getCorrectAnswerText = () => {
  if (!question.value?.options) return ''
  
  if (Array.isArray(question.value.options)) {
    const correctIndex = question.value.options.findIndex((opt: any) => 
      typeof opt === 'object' ? opt.isCorrect : false
    )
    if (correctIndex !== -1) {
      const option = question.value.options[correctIndex]
      const content = typeof option === 'object' ? option.content : option
      return `${String.fromCharCode(65 + correctIndex)}. ${content}`
    }
  }
  
  return question.value.correctAnswer?.toString() || ''
}

const nextQuestion = () => {
  // 随机获取下一题
  router.push('/questions/random')
}

const addToMistakes = async () => {
  if (!question.value || !authStore.isLoggedIn) {
    alert('请先登录')
    return
  }
  
  try {
    await QuestionService.addToMistakeBook(question.value.id)
    alert('已添加到错题本')
  } catch (error) {
    console.error('添加到错题本失败:', error)
    alert('添加到错题本失败')
  }
}

const shareQuestion = () => {
  if (!question.value) return
  
  const url = window.location.href
  const title = `刷刷题 - ${question.value.title}`
  
  if (navigator.share) {
    navigator.share({
      title,
      url
    })
  } else {
    // 复制链接到剪贴板
    navigator.clipboard.writeText(url).then(() => {
      alert('链接已复制到剪贴板')
    })
  }
}

const goBack = () => {
  router.back()
}

const loadQuestion = async () => {
  const questionId = Number(route.params.id)
  if (!questionId) {
    alert('题目ID无效')
    router.back()
    return
  }
  
  try {
    const response = await QuestionService.getQuestion(questionId)
    if (response.success && response.data) {
      question.value = response.data
      startTime.value = Date.now()
    } else {
      alert('加载题目失败')
      router.back()
    }
  } catch (error) {
    console.error('加载题目失败:', error)
    alert('加载题目失败')
    router.back()
  }
}

// 生命周期
onMounted(() => {
  loadQuestion()
})
</script>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.line-clamp-3 {
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>