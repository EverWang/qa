<template>
  <div class="system-settings">
    <div class="page-header">
      <h1 class="page-title">系统设置</h1>
    </div>

    <div class="page-content">
      <el-tabs v-model="activeTab" type="border-card">
        <!-- 基础设置 -->
        <el-tab-pane label="基础设置" name="basic">
          <el-card>
            <template #header>
              <div class="card-header">
                <span>系统基础配置</span>
                <el-button type="primary" @click="handleSaveBasic" :loading="saving">
                  保存设置
                </el-button>
              </div>
            </template>
            
            <el-form
              ref="basicFormRef"
              :model="basicForm"
              :rules="basicRules"
              label-width="120px"
              class="settings-form"
            >
              <el-form-item label="系统名称" prop="system_name">
                <el-input
                  v-model="basicForm.system_name"
                  placeholder="请输入系统名称"
                  maxlength="50"
                />
              </el-form-item>
              
              <el-form-item label="系统描述" prop="system_description">
                <el-input
                  v-model="basicForm.system_description"
                  type="textarea"
                  :rows="3"
                  placeholder="请输入系统描述"
                  maxlength="200"
                />
              </el-form-item>
              
              <el-form-item label="系统版本" prop="system_version">
                <el-input
                  v-model="basicForm.system_version"
                  placeholder="请输入系统版本"
                  maxlength="20"
                />
              </el-form-item>
              
              <el-form-item label="联系邮箱" prop="contact_email">
                <el-input
                  v-model="basicForm.contact_email"
                  type="email"
                  placeholder="请输入联系邮箱"
                  maxlength="100"
                />
              </el-form-item>
              
              <el-form-item label="系统状态">
                <el-radio-group v-model="basicForm.system_status">
                  <el-radio label="normal">正常运行</el-radio>
                  <el-radio label="maintenance">维护模式</el-radio>
                </el-radio-group>
              </el-form-item>
              
              <el-form-item label="维护公告" v-if="basicForm.system_status === 'maintenance'">
                <el-input
                  v-model="basicForm.maintenance_notice"
                  type="textarea"
                  :rows="3"
                  placeholder="请输入维护公告内容"
                  maxlength="500"
                />
              </el-form-item>
            </el-form>
          </el-card>
        </el-tab-pane>

        <!-- 答题设置 -->
        <el-tab-pane label="答题设置" name="quiz">
          <el-card>
            <template #header>
              <div class="card-header">
                <span>答题相关配置</span>
                <el-button type="primary" @click="handleSaveQuiz" :loading="saving">
                  保存设置
                </el-button>
              </div>
            </template>
            
            <el-form
              ref="quizFormRef"
              :model="quizForm"
              :rules="quizRules"
              label-width="150px"
              class="settings-form"
            >
              <el-form-item label="每日答题限制" prop="daily_limit">
                <el-input-number
                  v-model="quizForm.daily_limit"
                  :min="0"
                  :max="1000"
                  placeholder="0表示无限制"
                />
                <span class="form-tip">用户每日最多可答题数量，0表示无限制</span>
              </el-form-item>
              
              <el-form-item label="答题时间限制" prop="time_limit">
                <el-input-number
                  v-model="quizForm.time_limit"
                  :min="0"
                  :max="3600"
                  placeholder="0表示无限制"
                />
                <span class="form-tip">每题答题时间限制（秒），0表示无限制</span>
              </el-form-item>
              
              <el-form-item label="积分奖励">
                <el-checkbox v-model="quizForm.enable_points">启用积分系统</el-checkbox>
              </el-form-item>
              
              <div v-if="quizForm.enable_points" class="points-config">
                <el-form-item label="正确答题积分" prop="correct_points">
                  <el-input-number
                    v-model="quizForm.correct_points"
                    :min="0"
                    :max="100"
                  />
                </el-form-item>
                
                <el-form-item label="错误答题积分" prop="wrong_points">
                  <el-input-number
                    v-model="quizForm.wrong_points"
                    :min="-100"
                    :max="0"
                  />
                </el-form-item>
              </div>
              
              <el-form-item label="答题模式">
                <el-checkbox-group v-model="quizForm.quiz_modes">
                  <el-checkbox label="random">随机模式</el-checkbox>
                  <el-checkbox label="category">分类模式</el-checkbox>
                  <el-checkbox label="difficulty">难度模式</el-checkbox>
                  <el-checkbox label="exam">考试模式</el-checkbox>
                </el-checkbox-group>
              </el-form-item>
              
              <el-form-item label="显示解析">
                <el-radio-group v-model="quizForm.show_explanation">
                  <el-radio label="always">总是显示</el-radio>
                  <el-radio label="after_answer">答题后显示</el-radio>
                  <el-radio label="never">从不显示</el-radio>
                </el-radio-group>
              </el-form-item>
            </el-form>
          </el-card>
        </el-tab-pane>

        <!-- 数据统计 -->
        <el-tab-pane label="数据统计" name="statistics">
          <div class="statistics-content">
            <!-- 概览统计 -->
            <el-row :gutter="20" class="stats-overview">
              <el-col :span="6">
                <el-card class="stat-card">
                  <div class="stat-item">
                    <div class="stat-value">{{ statistics.total_users }}</div>
                    <div class="stat-label">总用户数</div>
                  </div>
                </el-card>
              </el-col>
              <el-col :span="6">
                <el-card class="stat-card">
                  <div class="stat-item">
                    <div class="stat-value">{{ statistics.total_questions }}</div>
                    <div class="stat-label">总题目数</div>
                  </div>
                </el-card>
              </el-col>
              <el-col :span="6">
                <el-card class="stat-card">
                  <div class="stat-item">
                    <div class="stat-value">{{ statistics.total_answers }}</div>
                    <div class="stat-label">总答题数</div>
                  </div>
                </el-card>
              </el-col>
              <el-col :span="6">
                <el-card class="stat-card">
                  <div class="stat-item">
                    <div class="stat-value">{{ statistics.active_users }}</div>
                    <div class="stat-label">活跃用户</div>
                  </div>
                </el-card>
              </el-col>
            </el-row>

            <!-- 详细统计 -->
            <el-row :gutter="20" class="detailed-stats">
              <el-col :span="12">
                <el-card>
                  <template #header>
                    <span>用户统计</span>
                  </template>
                  <el-table :data="statistics.user_stats" border>
                    <el-table-column prop="date" label="日期" width="120" />
                    <el-table-column prop="new_users" label="新增用户" />
                    <el-table-column prop="active_users" label="活跃用户" />
                  </el-table>
                </el-card>
              </el-col>
              <el-col :span="12">
                <el-card>
                  <template #header>
                    <span>答题统计</span>
                  </template>
                  <el-table :data="statistics.answer_stats" border>
                    <el-table-column prop="date" label="日期" width="120" />
                    <el-table-column prop="total_answers" label="总答题" />
                    <el-table-column prop="correct_rate" label="正确率" />
                  </el-table>
                </el-card>
              </el-col>
            </el-row>

            <!-- 刷新按钮 -->
            <div class="stats-actions">
              <el-button type="primary" @click="fetchStatistics" :loading="loadingStats">
                <el-icon><Refresh /></el-icon>
                刷新数据
              </el-button>
              <el-button @click="exportStatistics">
                <el-icon><Download /></el-icon>
                导出报表
              </el-button>
            </div>
          </div>
        </el-tab-pane>

        <!-- 操作日志 -->
        <el-tab-pane label="操作日志" name="logs">
          <el-card>
            <template #header>
              <div class="card-header">
                <span>系统操作日志</span>
                <div class="log-filters">
                  <el-select v-model="logFilter.type" placeholder="操作类型" clearable style="width: 120px">
                    <el-option label="登录" value="login" />
                    <el-option label="创建" value="create" />
                    <el-option label="更新" value="update" />
                    <el-option label="删除" value="delete" />
                  </el-select>
                  <el-date-picker
                    v-model="logFilter.date_range"
                    type="daterange"
                    range-separator="至"
                    start-placeholder="开始日期"
                    end-placeholder="结束日期"
                    format="YYYY-MM-DD"
                    value-format="YYYY-MM-DD"
                    style="width: 240px"
                  />
                  <el-button type="primary" @click="fetchLogs">
                    <el-icon><Search /></el-icon>
                    查询
                  </el-button>
                </div>
              </div>
            </template>
            
            <el-table :data="logs" v-loading="loadingLogs" border>
              <el-table-column prop="id" label="ID" width="80" />
              <el-table-column prop="user_name" label="操作用户" width="120" />
              <el-table-column prop="action" label="操作类型" width="100">
                <template #default="{ row }">
                  <el-tag :type="getActionTagType(row.action)" size="small">
                    {{ getActionLabel(row.action) }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="resource" label="操作对象" width="120" />
              <el-table-column prop="description" label="操作描述" min-width="200" />
              <el-table-column prop="ip_address" label="IP地址" width="140" />
              <el-table-column prop="user_agent" label="用户代理" min-width="200" show-overflow-tooltip />
              <el-table-column prop="created_at" label="操作时间" width="160">
                <template #default="{ row }">
                  {{ formatDate(row.created_at) }}
                </template>
              </el-table-column>
            </el-table>
            
            <!-- 分页 -->
            <div class="pagination-wrapper">
              <el-pagination
                v-model:current-page="logPagination.page"
                v-model:page-size="logPagination.size"
                :page-sizes="[10, 20, 50, 100]"
                :total="logTotal"
                layout="total, sizes, prev, pager, next, jumper"
                @size-change="handleLogSizeChange"
                @current-change="handleLogPageChange"
              />
            </div>
          </el-card>
        </el-tab-pane>
      </el-tabs>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, type FormInstance } from 'element-plus'
import {
  Refresh,
  Download,
  Search
} from '@element-plus/icons-vue'
import { api } from '@/lib/axios'

// 接口定义
interface BasicSettings {
  system_name: string
  system_description: string
  system_version: string
  contact_email: string
  system_status: 'normal' | 'maintenance'
  maintenance_notice: string
}

interface QuizSettings {
  daily_limit: number
  time_limit: number
  enable_points: boolean
  correct_points: number
  wrong_points: number
  quiz_modes: string[]
  show_explanation: 'always' | 'after_answer' | 'never'
}

interface Statistics {
  total_users: number
  total_questions: number
  total_answers: number
  active_users: number
  user_stats: Array<{
    date: string
    new_users: number
    active_users: number
  }>
  answer_stats: Array<{
    date: string
    total_answers: number
    correct_rate: string
  }>
}

interface Log {
  id: number
  user_name: string
  action: string
  resource: string
  description: string
  ip_address: string
  user_agent: string
  created_at: string
}

// 响应式数据
const activeTab = ref('basic')
const saving = ref(false)
const loadingStats = ref(false)
const loadingLogs = ref(false)
const basicFormRef = ref<FormInstance>()
const quizFormRef = ref<FormInstance>()

// 基础设置表单
const basicForm = reactive<BasicSettings>({
  system_name: '',
  system_description: '',
  system_version: '',
  contact_email: '',
  system_status: 'normal',
  maintenance_notice: ''
})

// 答题设置表单
const quizForm = reactive<QuizSettings>({
  daily_limit: 0,
  time_limit: 0,
  enable_points: false,
  correct_points: 1,
  wrong_points: 0,
  quiz_modes: ['random'],
  show_explanation: 'after_answer'
})

// 统计数据
const statistics = reactive<Statistics>({
  total_users: 0,
  total_questions: 0,
  total_answers: 0,
  active_users: 0,
  user_stats: [],
  answer_stats: []
})

// 日志相关
const logs = ref<Log[]>([])
const logTotal = ref(0)
const logFilter = reactive({
  type: '',
  date_range: [] as string[]
})
const logPagination = reactive({
  page: 1,
  size: 20
})

// 表单验证规则
const basicRules = {
  system_name: [
    { required: true, message: '请输入系统名称', trigger: 'blur' }
  ],
  contact_email: [
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ]
}

const quizRules = {
  daily_limit: [
    { type: 'number', min: 0, max: 1000, message: '请输入0-1000之间的数字', trigger: 'blur' }
  ],
  time_limit: [
    { type: 'number', min: 0, max: 3600, message: '请输入0-3600之间的数字', trigger: 'blur' }
  ]
}

// 方法
const handleSaveBasic = async () => {
  if (!basicFormRef.value) return
  
  try {
    await basicFormRef.value.validate()
    saving.value = true
    
    await api.put('/admin/settings/basic', basicForm)
    ElMessage.success('基础设置保存成功')
  } catch (error: any) {
    console.error('保存基础设置失败:', error)
    ElMessage.error(error.response?.data?.message || '保存失败')
  } finally {
    saving.value = false
  }
}

const handleSaveQuiz = async () => {
  if (!quizFormRef.value) return
  
  try {
    await quizFormRef.value.validate()
    saving.value = true
    
    await api.put('/admin/settings/quiz', quizForm)
    ElMessage.success('答题设置保存成功')
  } catch (error: any) {
    console.error('保存答题设置失败:', error)
    ElMessage.error(error.response?.data?.message || '保存失败')
  } finally {
    saving.value = false
  }
}

const fetchSettings = async () => {
  try {
    const [basicResponse, quizResponse] = await Promise.all([
      api.get('/admin/settings/basic'),
      api.get('/admin/settings/quiz')
    ])
    
    Object.assign(basicForm, basicResponse.data)
    Object.assign(quizForm, quizResponse.data)
  } catch (error) {
    console.error('获取设置失败:', error)
    ElMessage.error('获取设置失败')
  }
}

const fetchStatistics = async () => {
  try {
    loadingStats.value = true
    const response = await api.get('/admin/statistics')
    Object.assign(statistics, response.data)
  } catch (error) {
    console.error('获取统计数据失败:', error)
    ElMessage.error('获取统计数据失败')
  } finally {
    loadingStats.value = false
  }
}

const exportStatistics = async () => {
  try {
    const response = await api.get('/admin/statistics/export', {
      responseType: 'blob'
    })
    
    const blob = new Blob([response.data], {
      type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'
    })
    
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `统计报表_${new Date().toISOString().split('T')[0]}.xlsx`
    link.click()
    window.URL.revokeObjectURL(url)
    
    ElMessage.success('报表导出成功')
  } catch (error) {
    console.error('导出报表失败:', error)
    ElMessage.error('导出报表失败')
  }
}

const fetchLogs = async () => {
  try {
    loadingLogs.value = true
    
    const params = {
      page: logPagination.page,
      size: logPagination.size,
      type: logFilter.type || undefined,
      start_date: logFilter.date_range[0] || undefined,
      end_date: logFilter.date_range[1] || undefined
    }
    
    const response = await api.get('/admin/logs', { params })
    const data = response.data
    
    logs.value = data.data || []
    logTotal.value = data.total || 0
  } catch (error) {
    console.error('获取操作日志失败:', error)
    ElMessage.error('获取操作日志失败')
  } finally {
    loadingLogs.value = false
  }
}

const handleLogPageChange = (page: number) => {
  logPagination.page = page
  fetchLogs()
}

const handleLogSizeChange = (size: number) => {
  logPagination.size = size
  logPagination.page = 1
  fetchLogs()
}

const getActionLabel = (action: string) => {
  const actionMap: Record<string, string> = {
    login: '登录',
    create: '创建',
    update: '更新',
    delete: '删除',
    export: '导出',
    import: '导入'
  }
  return actionMap[action] || action
}

const getActionTagType = (action: string) => {
  const typeMap: Record<string, string> = {
    login: 'success',
    create: 'primary',
    update: 'warning',
    delete: 'danger',
    export: 'info',
    import: 'info'
  }
  return typeMap[action] || 'info'
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleString('zh-CN')
}

// 生命周期
onMounted(() => {
  fetchSettings()
  fetchStatistics()
  fetchLogs()
})
</script>

<style scoped>
.system-settings {
  padding: 24px;
}

.page-header {
  margin-bottom: 24px;
}

.page-title {
  font-size: 24px;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
}

.page-content {
  background: #fff;
  border-radius: 8px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.settings-form {
  max-width: 600px;
}

.form-tip {
  color: #6b7280;
  font-size: 12px;
  margin-left: 8px;
}

.points-config {
  margin-left: 24px;
  padding-left: 16px;
  border-left: 2px solid #e5e7eb;
}

.statistics-content {
  padding: 20px;
}

.stats-overview {
  margin-bottom: 24px;
}

.stat-card {
  text-align: center;
}

.stat-item {
  padding: 16px;
}

.stat-value {
  font-size: 32px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 8px;
}

.stat-label {
  color: #6b7280;
  font-size: 14px;
}

.detailed-stats {
  margin-bottom: 24px;
}

.stats-actions {
  display: flex;
  justify-content: center;
  gap: 12px;
}

.log-filters {
  display: flex;
  gap: 12px;
  align-items: center;
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: 24px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .system-settings {
    padding: 16px;
  }
  
  .card-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }
  
  .log-filters {
    flex-direction: column;
    align-items: flex-start;
    width: 100%;
  }
  
  .stats-overview .el-col {
    margin-bottom: 16px;
  }
  
  .detailed-stats .el-col {
    margin-bottom: 16px;
  }
  
  .stats-actions {
    flex-direction: column;
  }
}
</style>