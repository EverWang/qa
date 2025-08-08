<template>
  <view class="user-records">
    <!-- 顶部导航 -->
    <view class="header">
      <view class="nav-bar">
        <view class="nav-left" @click="goBack">
          <text class="icon-back">‹</text>
        </view>
        <view class="nav-title">做题记录</view>
        <view class="nav-right"></view>
      </view>
    </view>

    <!-- 统计信息 -->
    <view class="stats-section">
      <view class="stats-card">
        <view class="stats-item">
          <text class="stats-number">{{ totalQuestions }}</text>
          <text class="stats-label">总题数</text>
        </view>
        <view class="stats-item">
          <text class="stats-number">{{ correctRate }}%</text>
          <text class="stats-label">正确率</text>
        </view>
        <view class="stats-item">
          <text class="stats-number">{{ studyDays }}</text>
          <text class="stats-label">学习天数</text>
        </view>
      </view>
    </view>

    <!-- 筛选选项 -->
    <view class="filter-section">
      <scroll-view class="filter-scroll" scroll-x>
        <view class="filter-item" 
              :class="{ active: selectedFilter === 'all' }"
              @click="selectFilter('all')">
          全部
        </view>
        <view class="filter-item" 
              :class="{ active: selectedFilter === 'correct' }"
              @click="selectFilter('correct')">
          答对
        </view>
        <view class="filter-item" 
              :class="{ active: selectedFilter === 'wrong' }"
              @click="selectFilter('wrong')">
          答错
        </view>
        <view class="filter-item" 
              :class="{ active: selectedFilter === 'today' }"
              @click="selectFilter('today')">
          今日
        </view>
        <view class="filter-item" 
              :class="{ active: selectedFilter === 'week' }"
              @click="selectFilter('week')">
          本周
        </view>
      </scroll-view>
    </view>

    <!-- 记录列表 -->
    <view class="records-list">
      <view v-if="filteredRecords.length === 0" class="empty-state">
        <text class="empty-icon">📝</text>
        <text class="empty-text">暂无做题记录</text>
        <text class="empty-desc">快去刷题吧！</text>
      </view>
      
      <view v-else>
        <view v-for="record in filteredRecords" :key="record.id" class="record-item">
          <view class="record-header">
            <view class="record-title">{{ record.questionTitle }}</view>
            <view class="record-status" :class="record.isCorrect ? 'correct' : 'wrong'">
              {{ record.isCorrect ? '✓' : '✗' }}
            </view>
          </view>
          
          <view class="record-info">
            <text class="record-category">{{ record.category }}</text>
            <text class="record-difficulty" :class="`difficulty-${record.difficulty}`">
              {{ getDifficultyText(record.difficulty) }}
            </text>
          </view>
          
          <view class="record-meta">
            <text class="record-time">{{ formatTime(record.answeredAt) }}</text>
            <text class="record-duration">用时 {{ record.duration }}s</text>
          </view>
          
          <view class="record-actions">
            <view class="action-btn" @click="viewQuestion(record.questionId)">
              查看题目
            </view>
            <view v-if="!record.isCorrect" class="action-btn primary" @click="retryQuestion(record.questionId)">
              重新练习
            </view>
          </view>
        </view>
      </view>
    </view>

    <!-- 加载更多 -->
    <view v-if="hasMore && filteredRecords.length > 0" class="load-more" @click="loadMore">
      <text class="load-more-text">{{ isLoadingMore ? '加载中...' : '加载更多' }}</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'

interface QuestionRecord {
  id: string
  questionId: string
  questionTitle: string
  category: string
  difficulty: 'easy' | 'medium' | 'hard'
  isCorrect: boolean
  answeredAt: string
  duration: number
}

const router = useRouter()

// 响应式数据
const records = ref<QuestionRecord[]>([])
const selectedFilter = ref('all')
const isLoadingMore = ref(false)
const hasMore = ref(true)

// 统计数据
const totalQuestions = computed(() => records.value.length)
const correctRate = computed(() => {
  if (records.value.length === 0) return 0
  const correctCount = records.value.filter(r => r.isCorrect).length
  return Math.round((correctCount / records.value.length) * 100)
})
const studyDays = ref(15)

// 筛选后的记录
const filteredRecords = computed(() => {
  let filtered = records.value
  
  switch (selectedFilter.value) {
    case 'correct':
      filtered = records.value.filter(r => r.isCorrect)
      break
    case 'wrong':
      filtered = records.value.filter(r => !r.isCorrect)
      break
    case 'today':
      const today = new Date().toDateString()
      filtered = records.value.filter(r => new Date(r.answeredAt).toDateString() === today)
      break
    case 'week':
      const weekAgo = new Date()
      weekAgo.setDate(weekAgo.getDate() - 7)
      filtered = records.value.filter(r => new Date(r.answeredAt) >= weekAgo)
      break
  }
  
  return filtered.sort((a, b) => new Date(b.answeredAt).getTime() - new Date(a.answeredAt).getTime())
})

// 方法
const goBack = () => {
  router.back()
}

const selectFilter = (filter: string) => {
  selectedFilter.value = filter
}

const getDifficultyText = (difficulty: string) => {
  const map = {
    easy: '简单',
    medium: '中等',
    hard: '困难'
  }
  return map[difficulty as keyof typeof map] || difficulty
}

const formatTime = (time: string) => {
  const date = new Date(time)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  
  if (diff < 60000) return '刚刚'
  if (diff < 3600000) return `${Math.floor(diff / 60000)}分钟前`
  if (diff < 86400000) return `${Math.floor(diff / 3600000)}小时前`
  if (diff < 604800000) return `${Math.floor(diff / 86400000)}天前`
  
  return date.toLocaleDateString()
}

const viewQuestion = (questionId: string) => {
  router.push(`/question/${questionId}`)
}

const retryQuestion = (questionId: string) => {
  router.push(`/question/${questionId}?retry=true`)
}

const loadMore = async () => {
  if (isLoadingMore.value) return
  
  isLoadingMore.value = true
  
  // 模拟加载更多数据
  await new Promise(resolve => setTimeout(resolve, 1000))
  
  // 这里应该调用API加载更多记录
  // const newRecords = await questionService.getUserRecords({ page: currentPage + 1 })
  
  isLoadingMore.value = false
}

const loadRecords = async () => {
  // 模拟加载数据
  await new Promise(resolve => setTimeout(resolve, 500))
  
  // 模拟数据
  const mockRecords: QuestionRecord[] = [
    {
      id: '1',
      questionId: 'q1',
      questionTitle: 'JavaScript中var、let、const的区别是什么？',
      category: 'JavaScript',
      difficulty: 'medium',
      isCorrect: true,
      answeredAt: new Date().toISOString(),
      duration: 45
    },
    {
      id: '2',
      questionId: 'q2',
      questionTitle: 'Vue3的响应式原理是什么？',
      category: 'Vue.js',
      difficulty: 'hard',
      isCorrect: false,
      answeredAt: new Date(Date.now() - 3600000).toISOString(),
      duration: 120
    },
    {
      id: '3',
      questionId: 'q3',
      questionTitle: 'CSS盒模型包括哪些部分？',
      category: 'CSS',
      difficulty: 'easy',
      isCorrect: true,
      answeredAt: new Date(Date.now() - 7200000).toISOString(),
      duration: 30
    }
  ]
  
  records.value = mockRecords
}

onMounted(() => {
  loadRecords()
})
</script>

<style scoped>
.user-records {
  min-height: 100vh;
  background-color: #f5f5f5;
}

.header {
  background: white;
  border-bottom: 1px solid #eee;
}

.nav-bar {
  display: flex;
  align-items: center;
  height: 44px;
  padding: 0 16px;
}

.nav-left {
  width: 44px;
  height: 44px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.icon-back {
  font-size: 24px;
  color: #333;
}

.nav-title {
  flex: 1;
  text-align: center;
  font-size: 18px;
  font-weight: 600;
  color: #333;
}

.nav-right {
  width: 44px;
}

.stats-section {
  padding: 16px;
}

.stats-card {
  background: white;
  border-radius: 12px;
  padding: 20px;
  display: flex;
  justify-content: space-around;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.stats-item {
  text-align: center;
}

.stats-number {
  display: block;
  font-size: 24px;
  font-weight: bold;
  color: #007AFF;
  margin-bottom: 4px;
}

.stats-label {
  font-size: 14px;
  color: #666;
}

.filter-section {
  padding: 0 16px 16px;
}

.filter-scroll {
  white-space: nowrap;
}

.filter-item {
  display: inline-block;
  padding: 8px 16px;
  margin-right: 12px;
  background: white;
  border-radius: 20px;
  font-size: 14px;
  color: #666;
  border: 1px solid #eee;
}

.filter-item.active {
  background: #007AFF;
  color: white;
  border-color: #007AFF;
}

.records-list {
  padding: 0 16px;
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
  background: white;
  border-radius: 12px;
}

.empty-icon {
  display: block;
  font-size: 48px;
  margin-bottom: 16px;
}

.empty-text {
  display: block;
  font-size: 18px;
  color: #333;
  margin-bottom: 8px;
}

.empty-desc {
  font-size: 14px;
  color: #999;
}

.record-item {
  background: white;
  border-radius: 12px;
  padding: 16px;
  margin-bottom: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.record-header {
  display: flex;
  align-items: flex-start;
  margin-bottom: 12px;
}

.record-title {
  flex: 1;
  font-size: 16px;
  color: #333;
  line-height: 1.4;
  margin-right: 12px;
}

.record-status {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: bold;
}

.record-status.correct {
  background: #34C759;
  color: white;
}

.record-status.wrong {
  background: #FF3B30;
  color: white;
}

.record-info {
  display: flex;
  align-items: center;
  margin-bottom: 8px;
}

.record-category {
  font-size: 12px;
  color: #007AFF;
  background: #E3F2FD;
  padding: 2px 8px;
  border-radius: 10px;
  margin-right: 8px;
}

.record-difficulty {
  font-size: 12px;
  padding: 2px 8px;
  border-radius: 10px;
}

.difficulty-easy {
  color: #34C759;
  background: #E8F5E8;
}

.difficulty-medium {
  color: #FF9500;
  background: #FFF3E0;
}

.difficulty-hard {
  color: #FF3B30;
  background: #FFEBEE;
}

.record-meta {
  display: flex;
  justify-content: space-between;
  margin-bottom: 12px;
}

.record-time,
.record-duration {
  font-size: 12px;
  color: #999;
}

.record-actions {
  display: flex;
  gap: 12px;
}

.action-btn {
  flex: 1;
  text-align: center;
  padding: 8px 16px;
  border-radius: 6px;
  font-size: 14px;
  border: 1px solid #ddd;
  color: #666;
}

.action-btn.primary {
  background: #007AFF;
  color: white;
  border-color: #007AFF;
}

.load-more {
  text-align: center;
  padding: 20px;
}

.load-more-text {
  color: #007AFF;
  font-size: 14px;
}
</style>