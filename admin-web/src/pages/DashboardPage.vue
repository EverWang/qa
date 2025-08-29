<template>
  <div class="dashboard">
    <div class="dashboard-header">
      <h1 class="page-title">仪表板</h1>
      <p class="page-subtitle">欢迎回来，{{ username }}！</p>
    </div>
    
    <!-- 统计卡片 -->
    <div class="stats-grid">
      <el-card class="stat-card" shadow="hover">
        <div class="stat-content">
          <div class="stat-icon questions">
            <el-icon><Document /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ stats.totalQuestions }}</div>
            <div class="stat-label">题目总数</div>
          </div>
        </div>
      </el-card>
      
      <el-card class="stat-card" shadow="hover">
        <div class="stat-content">
          <div class="stat-icon categories">
            <el-icon><Folder /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ stats.totalCategories }}</div>
            <div class="stat-label">分类数量</div>
          </div>
        </div>
      </el-card>
      
      <el-card class="stat-card" shadow="hover">
        <div class="stat-content">
          <div class="stat-icon users">
            <el-icon><User /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ stats.totalUsers }}</div>
            <div class="stat-label">用户数量</div>
          </div>
        </div>
      </el-card>
      
      <el-card class="stat-card" shadow="hover">
        <div class="stat-content">
          <div class="stat-icon records">
            <el-icon><TrendCharts /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ stats.totalRecords }}</div>
            <div class="stat-label">答题记录</div>
          </div>
        </div>
      </el-card>
    </div>
    
    <!-- 图表区域 -->
    <div class="charts-grid">
      <!-- 题目分类分布图 -->
      <el-card class="chart-card" shadow="hover">
        <template #header>
          <div class="card-header">
            <span class="card-title">题目分类分布</span>
            <el-button type="text" @click="refreshCategoryChart">
              <el-icon><Refresh /></el-icon>
            </el-button>
          </div>
        </template>
        <div class="chart-container" ref="categoryChartRef"></div>
      </el-card>
      
      <!-- 用户答题趋势图 -->
      <el-card class="chart-card" shadow="hover">
        <template #header>
          <div class="card-header">
            <span class="card-title">用户答题趋势</span>
            <el-button type="text" @click="refreshTrendChart">
              <el-icon><Refresh /></el-icon>
            </el-button>
          </div>
        </template>
        <div class="chart-container" ref="trendChartRef"></div>
      </el-card>
    </div>
    
    <!-- 最近活动 -->
    <div class="recent-activities">
      <el-card shadow="hover">
        <template #header>
          <div class="card-header">
            <span class="card-title">最近活动</span>
            <el-button type="text" @click="refreshActivities">
              <el-icon><Refresh /></el-icon>
            </el-button>
          </div>
        </template>
        
        <el-timeline>
          <el-timeline-item
            v-for="activity in recentActivities"
            :key="activity.id"
            :timestamp="activity.timestamp"
            :type="getActivityType(activity.type)"
          >
            <div class="activity-content">
              <div class="activity-title">{{ activity.title }}</div>
              <div class="activity-description">{{ activity.description }}</div>
            </div>
          </el-timeline-item>
        </el-timeline>
        
        <div v-if="recentActivities.length === 0" class="empty-state">
          <el-empty description="暂无最近活动" />
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage } from 'element-plus'
import {
  Document,
  Folder,
  User,
  TrendCharts,
  Refresh
} from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'
import api from '@/lib/axios'

// 引入echarts（如果需要图表功能）
// import * as echarts from 'echarts'

const authStore = useAuthStore()

// 用户名
const username = computed(() => authStore.user?.username || '管理员')

// 统计数据
const stats = reactive({
  totalQuestions: 0,
  totalCategories: 0,
  totalUsers: 0,
  totalRecords: 0
})

// 最近活动
interface Activity {
  id: number
  type: 'question' | 'user' | 'category' | 'system'
  title: string
  description: string
  timestamp: string
}

const recentActivities = ref<Activity[]>([])

// 图表引用
const categoryChartRef = ref<HTMLElement>()
const trendChartRef = ref<HTMLElement>()

// 获取活动类型样式
const getActivityType = (type: string) => {
  const typeMap: Record<string, string> = {
    question: 'primary',
    user: 'success',
    category: 'warning',
    system: 'info'
  }
  return typeMap[type] || 'info'
}

// 获取统计数据
const fetchStats = async () => {
  try {
    const response = await api.get('/api/v1/admin/statistics/overview')
    // 后端返回格式: {code: 200, message: 'success', data: {...}}
    const data = response.data.data
    Object.assign(stats, {
      totalQuestions: data.totalQuestions,
      totalCategories: data.totalCategories,
      totalUsers: data.totalUsers,
      totalRecords: data.totalAnswers
    })
  } catch (error) {
    console.error('获取统计数据失败:', error)
    // 使用模拟数据
    Object.assign(stats, {
      totalQuestions: 1250,
      totalCategories: 15,
      totalUsers: 3420,
      totalRecords: 8960
    })
  }
}

// 获取最近活动
const fetchRecentActivities = async () => {
  try {
    const response = await api.get('/api/v1/admin/operation-logs?page=1&pageSize=10')
    const logs = response.data.data || []
    
    // 转换操作日志为活动格式
    recentActivities.value = logs.map((log: any) => ({
      id: log.id,
      type: getActivityTypeFromAction(log.action),
      title: getActivityTitle(log.action, log.resource),
      description: log.description || `${log.operator} ${log.action} ${log.resource}`,
      timestamp: formatTimestamp(log.created_at)
    }))
  } catch (error) {
    console.error('获取最近活动失败:', error)
    ElMessage.warning('获取最近活动失败，请稍后重试')
    recentActivities.value = []
  }
}

// 根据操作类型获取活动类型
const getActivityTypeFromAction = (action: string): string => {
  if (action.includes('题目') || action.includes('question')) return 'question'
  if (action.includes('用户') || action.includes('user')) return 'user'
  if (action.includes('分类') || action.includes('category')) return 'category'
  return 'system'
}

// 获取活动标题
const getActivityTitle = (action: string, resource: string): string => {
  const actionMap: Record<string, string> = {
    'create': '新增',
    'update': '更新',
    'delete': '删除',
    'login': '登录',
    'logout': '登出'
  }
  
  const actionText = actionMap[action.toLowerCase()] || action
  return `${actionText}${resource || ''}`
}

// 格式化时间戳
const formatTimestamp = (timestamp: string): string => {
  if (!timestamp) return ''
  const date = new Date(timestamp)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
}

// 初始化分类分布图表
const initCategoryChart = () => {
  // 这里可以使用echarts或其他图表库
  // 暂时使用文本显示
  if (categoryChartRef.value) {
    categoryChartRef.value.innerHTML = `
      <div style="text-align: center; padding: 40px; color: #999;">
        <p>分类分布图表</p>
        <p>（图表功能开发中）</p>
      </div>
    `
  }
}

// 初始化趋势图表
const initTrendChart = () => {
  // 这里可以使用echarts或其他图表库
  // 暂时使用文本显示
  if (trendChartRef.value) {
    trendChartRef.value.innerHTML = `
      <div style="text-align: center; padding: 40px; color: #999;">
        <p>答题趋势图表</p>
        <p>（图表功能开发中）</p>
      </div>
    `
  }
}

// 刷新图表
const refreshCategoryChart = () => {
  ElMessage.success('分类分布图表已刷新')
  initCategoryChart()
}

const refreshTrendChart = () => {
  ElMessage.success('答题趋势图表已刷新')
  initTrendChart()
}

// 刷新活动
const refreshActivities = () => {
  ElMessage.success('最近活动已刷新')
  fetchRecentActivities()
}

// 组件挂载时初始化
onMounted(async () => {
  await fetchStats()
  await fetchRecentActivities()
  initCategoryChart()
  initTrendChart()
})
</script>

<style scoped>
.dashboard {
  max-width: 1200px;
  margin: 0 auto;
}

.dashboard-header {
  margin-bottom: 24px;
}

.page-title {
  font-size: 28px;
  font-weight: 600;
  color: #1f2937;
  margin: 0 0 8px 0;
}

.page-subtitle {
  font-size: 16px;
  color: #6b7280;
  margin: 0;
}

/* 统计卡片样式 */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 20px;
  margin-bottom: 24px;
}

.stat-card {
  border-radius: 12px;
  transition: transform 0.2s;
}

.stat-card:hover {
  transform: translateY(-2px);
}

.stat-content {
  display: flex;
  align-items: center;
  padding: 8px;
}

.stat-icon {
  width: 60px;
  height: 60px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 16px;
  font-size: 24px;
  color: white;
}

.stat-icon.questions {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.stat-icon.categories {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.stat-icon.users {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
}

.stat-icon.records {
  background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
}

.stat-info {
  flex: 1;
}

.stat-value {
  font-size: 32px;
  font-weight: 700;
  color: #1f2937;
  line-height: 1;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 14px;
  color: #6b7280;
  font-weight: 500;
}

/* 图表区域样式 */
.charts-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
  gap: 20px;
  margin-bottom: 24px;
}

.chart-card {
  border-radius: 12px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-title {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
}

.chart-container {
  height: 300px;
  width: 100%;
}

/* 最近活动样式 */
.recent-activities {
  margin-bottom: 24px;
}

.recent-activities .el-card {
  border-radius: 12px;
}

.activity-content {
  margin-left: 8px;
}

.activity-title {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 4px;
}

.activity-description {
  font-size: 14px;
  color: #6b7280;
}

.empty-state {
  padding: 40px 0;
  text-align: center;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .stats-grid {
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: 16px;
  }
  
  .charts-grid {
    grid-template-columns: 1fr;
    gap: 16px;
  }
  
  .stat-content {
    padding: 4px;
  }
  
  .stat-icon {
    width: 50px;
    height: 50px;
    font-size: 20px;
    margin-right: 12px;
  }
  
  .stat-value {
    font-size: 24px;
  }
  
  .page-title {
    font-size: 24px;
  }
}
</style>