<template>
  <div class="min-h-screen bg-gray-50">
    <!-- 顶部导航 -->
    <div class="bg-white shadow-sm border-b">
      <div class="max-w-4xl mx-auto px-4 py-3 flex items-center justify-between">
        <button @click="goBack" class="flex items-center text-gray-600 hover:text-gray-800">
          <ArrowLeft class="w-5 h-5 mr-1" />
          <span>返回</span>
        </button>
        <h1 class="text-lg font-medium text-gray-800">题目详情</h1>
        <button @click="shareQuestion" class="flex items-center text-blue-600 hover:text-blue-700">
          <Share2 class="w-5 h-5 mr-1" />
          <span>分享</span>
        </button>
      </div>
    </div>

    <div class="max-w-4xl mx-auto p-4">
      <!-- 题目信息 -->
      <div class="bg-white rounded-lg shadow-sm p-6 mb-4">
        <div class="flex items-center justify-between mb-4">
          <div class="flex items-center space-x-2">
            <span class="px-2 py-1 bg-blue-100 text-blue-800 text-xs rounded-full">
              {{ question?.category?.name }}
            </span>
            <span class="px-2 py-1 text-xs rounded-full" :class="difficultyClass">
              {{ difficultyText }}
            </span>
          </div>
          <div class="text-sm text-gray-500">
            题目 {{ questionId }}
          </div>
        </div>

        <h2 class="text-xl font-medium text-gray-800 mb-4">{{ question?.title }}</h2>
        <div class="text-gray-700 mb-6 leading-relaxed" v-html="question?.content"></div>

        <!-- 选项 -->
        <div class="space-y-3">
          <div
            v-for="(option, index) in question?.options"
            :key="index"
            @click="selectOption(index)"
            class="p-4 border rounded-lg cursor-pointer transition-all duration-200"
            :class="getOptionClass(index)"
          >
            <div class="flex items-center">
              <div class="w-6 h-6 rounded-full border-2 mr-3 flex items-center justify-center" :class="getOptionIndicatorClass(index)">
                <span class="text-sm font-medium">{{ String.fromCharCode(65 + index) }}</span>
              </div>
              <span class="flex-1">{{ option }}</span>
              <CheckCircle v-if="showAnswer && index === question?.correctAnswer" class="w-5 h-5 text-green-600" />
              <XCircle v-if="showAnswer && selectedOption === index && index !== question?.correctAnswer" class="w-5 h-5 text-red-600" />
            </div>
          </div>
        </div>

        <!-- 答题按钮 -->
        <div class="mt-6 flex space-x-3">
          <button
            v-if="!showAnswer"
            @click="submitAnswer"
            :disabled="selectedOption === null"
            class="flex-1 bg-blue-600 hover:bg-blue-700 disabled:bg-gray-400 text-white font-medium py-3 px-4 rounded-lg transition-colors duration-200"
          >
            提交答案
          </button>
          <button
            v-if="showAnswer"
            @click="nextQuestion"
            class="flex-1 bg-green-600 hover:bg-green-700 text-white font-medium py-3 px-4 rounded-lg transition-colors duration-200"
          >
            下一题
          </button>
          <button
            v-if="showAnswer && !isCorrect"
            @click="addToMistakes"
            :disabled="isAddingToMistakes"
            class="px-4 py-3 bg-orange-100 hover:bg-orange-200 text-orange-700 rounded-lg transition-colors duration-200 flex items-center"
          >
            <BookmarkPlus class="w-4 h-4 mr-1" />
            {{ isAddingToMistakes ? '添加中...' : '加入错题本' }}
          </button>
        </div>
      </div>

      <!-- 答案解析 -->
      <div v-if="showAnswer && question?.explanation" class="bg-white rounded-lg shadow-sm p-6">
        <h3 class="text-lg font-medium text-gray-800 mb-3 flex items-center">
          <Lightbulb class="w-5 h-5 mr-2 text-yellow-500" />
          答案解析
        </h3>
        <div class="text-gray-700 leading-relaxed" v-html="question.explanation"></div>
      </div>
    </div>

    <!-- 答题结果弹窗 -->
    <div v-if="showResultModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg p-6 m-4 max-w-sm w-full">
        <div class="text-center">
          <div class="w-16 h-16 mx-auto mb-4 rounded-full flex items-center justify-center" :class="isCorrect ? 'bg-green-100' : 'bg-red-100'">
            <CheckCircle v-if="isCorrect" class="w-8 h-8 text-green-600" />
            <XCircle v-else class="w-8 h-8 text-red-600" />
          </div>
          <h3 class="text-lg font-medium mb-2" :class="isCorrect ? 'text-green-800' : 'text-red-800'">
            {{ isCorrect ? '回答正确！' : '回答错误' }}
          </h3>
          <p class="text-gray-600 mb-4">
            {{ isCorrect ? '恭喜你答对了这道题！' : `正确答案是 ${String.fromCharCode(65 + (question?.correctAnswer || 0))}` }}
          </p>
          <button
            @click="closeResultModal"
            class="w-full bg-blue-600 hover:bg-blue-700 text-white font-medium py-2 px-4 rounded-lg transition-colors duration-200"
          >
            查看解析
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ArrowLeft, Share2, CheckCircle, XCircle, Lightbulb, BookmarkPlus } from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'

interface Question {
  id: number
  title: string
  content: string
  options: string[]
  correctAnswer: number
  explanation: string
  difficulty: 'easy' | 'medium' | 'hard'
  category: {
    id: number
    name: string
  }
}

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const questionId = computed(() => route.params.id as string)
const question = ref<Question | null>(null)
const selectedOption = ref<number | null>(null)
const showAnswer = ref(false)
const showResultModal = ref(false)
const isCorrect = ref(false)
const isAddingToMistakes = ref(false)

// 难度样式
const difficultyClass = computed(() => {
  switch (question.value?.difficulty) {
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
  switch (question.value?.difficulty) {
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

// 选项样式
const getOptionClass = (index: number) => {
  if (!showAnswer.value) {
    return selectedOption.value === index
      ? 'border-blue-500 bg-blue-50'
      : 'border-gray-200 hover:border-gray-300 hover:bg-gray-50'
  }
  
  if (index === question.value?.correctAnswer) {
    return 'border-green-500 bg-green-50'
  }
  
  if (selectedOption.value === index && index !== question.value?.correctAnswer) {
    return 'border-red-500 bg-red-50'
  }
  
  return 'border-gray-200 bg-gray-50'
}

const getOptionIndicatorClass = (index: number) => {
  if (!showAnswer.value) {
    return selectedOption.value === index
      ? 'border-blue-500 bg-blue-500 text-white'
      : 'border-gray-300 text-gray-600'
  }
  
  if (index === question.value?.correctAnswer) {
    return 'border-green-500 bg-green-500 text-white'
  }
  
  if (selectedOption.value === index && index !== question.value?.correctAnswer) {
    return 'border-red-500 bg-red-500 text-white'
  }
  
  return 'border-gray-300 text-gray-600'
}

// 选择选项
const selectOption = (index: number) => {
  if (!showAnswer.value) {
    selectedOption.value = index
  }
}

// 提交答案
const submitAnswer = async () => {
  if (selectedOption.value === null) return
  
  isCorrect.value = selectedOption.value === question.value?.correctAnswer
  showResultModal.value = true
  
  // 记录答题结果
  try {
    // 这里调用API记录答题结果
    console.log('记录答题结果:', {
      questionId: questionId.value,
      userAnswer: selectedOption.value,
      isCorrect: isCorrect.value,
      timeSpent: 0 // 实际项目中需要计算答题时间
    })
  } catch (error) {
    console.error('记录答题结果失败:', error)
  }
}

// 关闭结果弹窗
const closeResultModal = () => {
  showResultModal.value = false
  showAnswer.value = true
}

// 下一题
const nextQuestion = () => {
  const nextId = parseInt(questionId.value) + 1
  router.push(`/question/${nextId}`)
}

// 加入错题本
const addToMistakes = async () => {
  if (!authStore.isLoggedIn || authStore.isGuest) {
    alert('请先登录后再使用错题本功能')
    return
  }
  
  try {
    isAddingToMistakes.value = true
    
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    console.log('添加到错题本:', questionId.value)
    alert('已添加到错题本')
  } catch (error) {
    console.error('添加到错题本失败:', error)
    alert('添加失败，请重试')
  } finally {
    isAddingToMistakes.value = false
  }
}

// 分享题目
const shareQuestion = () => {
  if (navigator.share) {
    navigator.share({
      title: question.value?.title || '刷刷题',
      text: '来看看这道题目',
      url: window.location.href
    })
  } else {
    // 复制链接到剪贴板
    navigator.clipboard.writeText(window.location.href)
    alert('链接已复制到剪贴板')
  }
}

// 返回
const goBack = () => {
  router.back()
}

// 加载题目数据
const loadQuestion = async () => {
  try {
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 500))
    
    // 模拟题目数据
    question.value = {
      id: parseInt(questionId.value),
      title: `第${questionId.value}题`,
      content: '以下关于道路工程施工安全的说法，正确的是？',
      options: [
        '施工现场可以不设置安全警示标志',
        '施工人员必须佩戴安全帽和反光背心',
        '夜间施工不需要特殊照明设备',
        '施工区域可以随意堆放材料'
      ],
      correctAnswer: 1,
      explanation: '根据《公路工程施工安全技术规范》，施工人员在施工现场必须佩戴安全帽和反光背心，这是保障施工安全的基本要求。其他选项都是错误的做法，会带来安全隐患。',
      difficulty: 'medium',
      category: {
        id: 1,
        name: '道路工程'
      }
    }
  } catch (error) {
    console.error('加载题目失败:', error)
    alert('加载题目失败，请重试')
  }
}

onMounted(() => {
  loadQuestion()
})
</script>

<style scoped>
/* 自定义样式 */
</style>