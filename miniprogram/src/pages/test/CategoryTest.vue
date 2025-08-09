<template>
  <div class="p-4">
    <h1 class="text-2xl font-bold mb-4">分类数据测试</h1>
    
    <div class="mb-4">
      <button @click="testCategories" class="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600">
        测试分类API
      </button>
    </div>
    
    <div v-if="loading" class="text-blue-600">
      加载中...
    </div>
    
    <div v-if="error" class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">
      <strong>错误:</strong> {{ error }}
    </div>
    
    <div v-if="categories.length > 0" class="bg-green-100 border border-green-400 text-green-700 px-4 py-3 rounded mb-4">
      <strong>成功:</strong> 加载了 {{ categories.length }} 个分类
    </div>
    
    <div v-if="rawResponse" class="mb-4">
      <h3 class="text-lg font-semibold mb-2">原始API响应:</h3>
      <pre class="bg-gray-100 p-4 rounded overflow-auto text-sm">{{ JSON.stringify(rawResponse, null, 2) }}</pre>
    </div>
    
    <div v-if="categories.length > 0">
      <h3 class="text-lg font-semibold mb-2">处理后的分类数据:</h3>
      <div class="grid gap-4">
        <div v-for="category in categories" :key="category.id" class="bg-white border rounded p-4">
          <h4 class="font-medium">{{ category.name }}</h4>
          <p class="text-sm text-gray-600">ID: {{ category.id }}</p>
          <p class="text-sm text-gray-600">级别: {{ category.level }}</p>
          <p class="text-sm text-gray-600">题目数量: {{ category.questionCount }}</p>
          <p class="text-sm text-gray-600">描述: {{ category.description || '无' }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { QuestionService } from '@/services/question'
import type { Category } from '@/services/question'

const loading = ref(false)
const error = ref('')
const categories = ref<Category[]>([])
const rawResponse = ref<any>(null)

/**
 * 测试分类API调用
 */
const testCategories = async () => {
  loading.value = true
  error.value = ''
  categories.value = []
  rawResponse.value = null
  
  try {
    console.log('开始测试分类API...')
    const response = await QuestionService.getCategories()
    console.log('分类API响应:', response)
    
    rawResponse.value = response
    
    if (response.success && response.data) {
      categories.value = response.data
      console.log('分类数据加载成功:', categories.value)
    } else {
      error.value = `API调用失败: ${response.message || '未知错误'}`
    }
  } catch (err: any) {
    console.error('分类API调用错误:', err)
    error.value = `网络错误: ${err.message || err.toString()}`
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
/* 样式 */
</style>