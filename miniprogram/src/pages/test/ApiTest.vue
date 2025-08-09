<template>
  <div class="api-test-page">
    <h1>API测试页面</h1>
    
    <div class="test-section">
      <h2>分类API测试</h2>
      <button @click="testCategories" :disabled="loading">测试分类API</button>
      <div v-if="categoriesResult">
        <h3>分类结果:</h3>
        <pre>{{ JSON.stringify(categoriesResult, null, 2) }}</pre>
      </div>
    </div>
    
    <div class="test-section">
      <h2>题目API测试</h2>
      <button @click="testQuestions" :disabled="loading">测试题目API</button>
      <div v-if="questionsResult">
        <h3>题目结果:</h3>
        <pre>{{ JSON.stringify(questionsResult, null, 2) }}</pre>
      </div>
    </div>
    
    <div class="test-section">
      <h2>统计API测试</h2>
      <button @click="testStatistics" :disabled="loading">测试统计API</button>
      <div v-if="statisticsResult">
        <h3>统计结果:</h3>
        <pre>{{ JSON.stringify(statisticsResult, null, 2) }}</pre>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { QuestionService } from '@/services/question'
import { StatisticsService } from '@/services/statistics'

const loading = ref(false)
const categoriesResult = ref(null)
const questionsResult = ref(null)
const statisticsResult = ref(null)

const testCategories = async () => {
  loading.value = true
  try {
    console.log('测试分类API...')
    const response = await QuestionService.getCategories()
    console.log('分类API响应:', response)
    categoriesResult.value = response
  } catch (error) {
    console.error('分类API测试失败:', error)
    categoriesResult.value = { error: error.message }
  } finally {
    loading.value = false
  }
}

const testQuestions = async () => {
  loading.value = true
  try {
    console.log('测试题目API...')
    const response = await QuestionService.getQuestions({ page: 1, pageSize: 5 })
    console.log('题目API响应:', response)
    questionsResult.value = response
  } catch (error) {
    console.error('题目API测试失败:', error)
    questionsResult.value = { error: error.message }
  } finally {
    loading.value = false
  }
}

const testStatistics = async () => {
  loading.value = true
  try {
    console.log('测试统计API...')
    const response = await StatisticsService.getOverviewStatistics()
    console.log('统计API响应:', response)
    statisticsResult.value = response
  } catch (error) {
    console.error('统计API测试失败:', error)
    statisticsResult.value = { error: error.message }
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.api-test-page {
  padding: 20px;
  max-width: 800px;
  margin: 0 auto;
}

.test-section {
  margin-bottom: 30px;
  padding: 20px;
  border: 1px solid #ddd;
  border-radius: 8px;
}

button {
  padding: 10px 20px;
  background: #007bff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

button:disabled {
  background: #ccc;
  cursor: not-allowed;
}

pre {
  background: #f5f5f5;
  padding: 10px;
  border-radius: 4px;
  overflow-x: auto;
  max-height: 300px;
  overflow-y: auto;
}
</style>