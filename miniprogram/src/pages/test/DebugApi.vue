<template>
  <div class="p-4">
    <h1 class="text-2xl font-bold mb-4">API调试页面</h1>
    
    <div class="space-y-4">
      <div class="bg-white p-4 rounded-lg shadow">
        <h2 class="text-lg font-semibold mb-2">分类API测试</h2>
        <button @click="testCategories" class="bg-blue-500 text-white px-4 py-2 rounded mr-2">测试分类API</button>
        <div v-if="categoriesResult" class="mt-2">
          <h3 class="font-medium">结果:</h3>
          <pre class="bg-gray-100 p-2 rounded text-sm overflow-auto">{{ JSON.stringify(categoriesResult, null, 2) }}</pre>
        </div>
      </div>
      
      <div class="bg-white p-4 rounded-lg shadow">
        <h2 class="text-lg font-semibold mb-2">题目API测试</h2>
        <button @click="testQuestions" class="bg-green-500 text-white px-4 py-2 rounded mr-2">测试题目API</button>
        <div v-if="questionsResult" class="mt-2">
          <h3 class="font-medium">结果:</h3>
          <pre class="bg-gray-100 p-2 rounded text-sm overflow-auto">{{ JSON.stringify(questionsResult, null, 2) }}</pre>
        </div>
      </div>
      
      <div class="bg-white p-4 rounded-lg shadow">
        <h2 class="text-lg font-semibold mb-2">统计API测试</h2>
        <button @click="testStatistics" class="bg-purple-500 text-white px-4 py-2 rounded mr-2">测试统计API</button>
        <div v-if="statisticsResult" class="mt-2">
          <h3 class="font-medium">结果:</h3>
          <pre class="bg-gray-100 p-2 rounded text-sm overflow-auto">{{ JSON.stringify(statisticsResult, null, 2) }}</pre>
        </div>
      </div>
      
      <div class="bg-white p-4 rounded-lg shadow">
        <h2 class="text-lg font-semibold mb-2">直接API测试</h2>
        <button @click="testDirectApi" class="bg-red-500 text-white px-4 py-2 rounded mr-2">直接测试API</button>
        <div v-if="directResult" class="mt-2">
          <h3 class="font-medium">结果:</h3>
          <pre class="bg-gray-100 p-2 rounded text-sm overflow-auto">{{ JSON.stringify(directResult, null, 2) }}</pre>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { QuestionService } from '@/services/question'
import { StatisticsService } from '@/services/statistics'
import { apiClient } from '@/services/api'

const categoriesResult = ref(null)
const questionsResult = ref(null)
const statisticsResult = ref(null)
const directResult = ref(null)

const testCategories = async () => {
  try {
    console.log('开始测试分类API...')
    const result = await QuestionService.getCategories()
    console.log('分类API结果:', result)
    categoriesResult.value = result
  } catch (error) {
    console.error('分类API测试失败:', error)
    categoriesResult.value = { error: error.message }
  }
}

const testQuestions = async () => {
  try {
    console.log('开始测试题目API...')
    const result = await QuestionService.getQuestions({ page: 1, pageSize: 3 })
    console.log('题目API结果:', result)
    questionsResult.value = result
  } catch (error) {
    console.error('题目API测试失败:', error)
    questionsResult.value = { error: error.message }
  }
}

const testStatistics = async () => {
  try {
    console.log('开始测试统计API...')
    const result = await StatisticsService.getOverviewStatistics()
    console.log('统计API结果:', result)
    statisticsResult.value = result
  } catch (error) {
    console.error('统计API测试失败:', error)
    statisticsResult.value = { error: error.message }
  }
}

const testDirectApi = async () => {
  try {
    console.log('开始直接测试API...')
    const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'
    const result = await fetch(`${API_BASE_URL}/api/v1/categories`)
    const data = await result.json()
    console.log('直接API结果:', data)
    directResult.value = data
  } catch (error) {
    console.error('直接API测试失败:', error)
    directResult.value = { error: error.message }
  }
}
</script>

<style scoped>
pre {
  max-height: 300px;
}
</style>