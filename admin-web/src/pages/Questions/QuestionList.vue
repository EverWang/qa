<template>
  <div class="question-list">
    <div class="page-header">
      <div class="header-left">
        <h1 class="page-title">题目管理</h1>
        <p class="page-subtitle">管理所有题目，支持搜索、编辑和删除</p>
      </div>
      <div class="header-right">
        <el-button type="primary" @click="$router.push('/questions/create')">
          <el-icon><Plus /></el-icon>
          添加题目
        </el-button>
        <el-button @click="$router.push('/questions/import')">
          <el-icon><Upload /></el-icon>
          批量导入
        </el-button>
      </div>
    </div>
    
    <!-- 搜索和筛选 -->
    <el-card class="search-card" shadow="never">
      <el-form :model="searchForm" inline>
        <el-form-item label="题目内容">
          <el-input
            v-model="searchForm.content"
            placeholder="请输入题目内容关键词"
            clearable
            style="width: 200px"
          />
        </el-form-item>
        <el-form-item label="题目类型">
          <el-select
            v-model="searchForm.type"
            placeholder="请选择题目类型"
            clearable
            style="width: 150px"
          >
            <el-option label="单选题" value="single" />
            <el-option label="多选题" value="multiple" />
            <el-option label="判断题" value="judge" />
            <el-option label="填空题" value="fill" />
          </el-select>
        </el-form-item>
        <el-form-item label="分类">
          <el-select
            v-model="searchForm.categoryId"
            placeholder="请选择分类"
            clearable
            style="width: 150px"
          >
            <el-option
              v-for="category in categories"
              :key="category.id"
              :label="category.name"
              :value="category.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="难度">
          <el-select
            v-model="searchForm.difficulty"
            placeholder="请选择难度"
            clearable
            style="width: 120px"
          >
            <el-option label="简单" value="easy" />
            <el-option label="中等" value="medium" />
            <el-option label="困难" value="hard" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>
            搜索
          </el-button>
          <el-button @click="handleReset">
            <el-icon><Refresh /></el-icon>
            重置
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>
    
    <!-- 题目列表 -->
    <el-card class="table-card" shadow="never">
      <template #header>
        <div class="card-header">
          <span>题目列表（共 {{ total }} 条）</span>
          <div class="header-actions">
            <el-button
              type="danger"
              :disabled="selectedQuestions.length === 0"
              @click="handleBatchDelete"
            >
              <el-icon><Delete /></el-icon>
              批量删除
            </el-button>
          </div>
        </div>
      </template>
      
      <el-table
        v-loading="loading"
        :data="questions"
        @selection-change="handleSelectionChange"
        stripe
        style="width: 100%"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="content" label="题目内容" min-width="300">
          <template #default="{ row }">
            <div class="question-content">
              <p class="content-text">{{ row.content }}</p>
              <div class="content-meta">
                <el-tag :type="getTypeTagType(row.type)" size="small">
                  {{ getTypeLabel(row.type) }}
                </el-tag>
                <el-tag :type="getDifficultyTagType(row.difficulty)" size="small">
                  {{ getDifficultyLabel(row.difficulty) }}
                </el-tag>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="category" label="分类" width="120">
          <template #default="{ row }">
            <span>{{ row.category?.name || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="创建时间" width="180">
          <template #default="{ row }">
            <span>{{ formatDate(row.created_at) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleView(row)">
              <el-icon><View /></el-icon>
              查看
            </el-button>
            <el-button type="primary" link @click="handleEdit(row)">
              <el-icon><Edit /></el-icon>
              编辑
            </el-button>
            <el-button type="danger" link @click="handleDelete(row)">
              <el-icon><Delete /></el-icon>
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.size"
          :total="total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>
    
    <!-- 题目详情对话框 -->
    <el-dialog
      v-model="detailDialogVisible"
      title="题目详情"
      width="800px"
      destroy-on-close
    >
      <QuestionDetail
        v-if="currentQuestion"
        :question="currentQuestion"
        @close="detailDialogVisible = false"
      />
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Plus,
  Upload,
  Search,
  Refresh,
  Delete,
  View,
  Edit
} from '@element-plus/icons-vue'
import api from '@/lib/axios'
import QuestionDetail from './QuestionDetail.vue'

// 路由
const router = useRouter()

// 题目接口
interface Question {
  id: number
  content: string
  type: 'single' | 'multiple' | 'judge' | 'fill'
  difficulty: 'easy' | 'medium' | 'hard'
  options?: string[]
  answer: string | string[]
  explanation?: string
  category_id: number
  category?: {
    id: number
    name: string
  }
  created_at: string
  updated_at: string
}

// 分类接口
interface Category {
  id: number
  name: string
}

// 搜索表单
const searchForm = reactive({
  content: '',
  type: '',
  categoryId: '',
  difficulty: ''
})

// 分页
const pagination = reactive({
  page: 1,
  size: 20
})

// 数据
const questions = ref<Question[]>([])
const categories = ref<Category[]>([])
const selectedQuestions = ref<Question[]>([])
const total = ref(0)
const loading = ref(false)

// 对话框
const detailDialogVisible = ref(false)
const currentQuestion = ref<Question | null>(null)

// 获取题目列表
const fetchQuestions = async () => {
  try {
    loading.value = true
    const params = {
      page: pagination.page,
      size: pagination.size,
      keyword: searchForm.content || undefined,
      type: searchForm.type || undefined, // 添加题目类型筛选参数
      category_id: searchForm.categoryId || undefined,
      difficulty: searchForm.difficulty || undefined
    }
    
    const response = await api.get('/api/v1/admin/questions', { params })
    const data = response.data
    
    // 后端返回的分页数据结构：{ code, message, data: [...], total, page, size }
    questions.value = data.data || []
    total.value = data.total || 0
  } catch (error) {
    console.error('获取题目列表失败:', error)
    ElMessage.error('获取题目列表失败')
    // 使用模拟数据
    questions.value = [
      {
        id: 1,
        content: 'Vue.js 是什么？',
        type: 'single',
        difficulty: 'easy',
        options: ['A. 前端框架', 'B. 后端框架', 'C. 数据库', 'D. 编程语言'],
        answer: 'A',
        explanation: 'Vue.js 是一个用于构建用户界面的渐进式前端框架。',
        categoryId: 1,
        category: { id: 1, name: '前端开发' },
        createdAt: '2024-01-15 10:30:00',
        updatedAt: '2024-01-15 10:30:00'
      },
      {
        id: 2,
        content: 'JavaScript 中的闭包是什么？',
        type: 'single',
        difficulty: 'medium',
        options: ['A. 函数', 'B. 对象', 'C. 变量', 'D. 函数和其词法环境的组合'],
        answer: 'D',
        explanation: '闭包是函数和声明该函数的词法环境的组合。',
        categoryId: 1,
        category: { id: 1, name: '前端开发' },
        createdAt: '2024-01-15 11:00:00',
        updatedAt: '2024-01-15 11:00:00'
      }
    ]
    total.value = 2
  } finally {
    loading.value = false
  }
}

// 获取分类列表
const fetchCategories = async () => {
  try {
    const response = await api.get('/api/v1/admin/categories', { params: { tree: 'true' } })
    const data = response.data
    
    // 如果是树形结构，直接使用；如果是分页结构，取data字段
    categories.value = Array.isArray(data.data) ? data.data : (data.data || [])
  } catch (error) {
    console.error('获取分类列表失败:', error)
    // 使用模拟数据
    categories.value = [
      { id: 1, name: '前端开发' },
      { id: 2, name: '后端开发' },
      { id: 3, name: '数据库' }
    ]
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  fetchQuestions()
}

// 重置
const handleReset = () => {
  Object.assign(searchForm, {
    content: '',
    type: '',
    categoryId: '',
    difficulty: ''
  })
  pagination.page = 1
  fetchQuestions()
}

// 分页变化
const handleSizeChange = (size: number) => {
  pagination.size = size
  fetchQuestions()
}

const handleCurrentChange = (page: number) => {
  pagination.page = page
  fetchQuestions()
}

// 选择变化
const handleSelectionChange = (selection: Question[]) => {
  selectedQuestions.value = selection
}

// 查看题目
const handleView = (question: Question) => {
  currentQuestion.value = question
  detailDialogVisible.value = true
}

// 编辑题目
const handleEdit = (question: Question) => {
  // 跳转到编辑页面
  router.push(`/questions/${question.id}/edit`)
}

// 删除题目
const handleDelete = async (question: Question) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除题目"${question.content}"吗？`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await api.delete(`/api/v1/admin/questions/${question.id}`)
    ElMessage.success('删除成功')
    fetchQuestions()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除题目失败:', error)
      ElMessage.error('删除失败')
    }
  }
}

// 批量删除
const handleBatchDelete = async () => {
  try {
    await ElMessageBox.confirm(
      `确定要删除选中的 ${selectedQuestions.value.length} 道题目吗？`,
      '确认批量删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    const ids = selectedQuestions.value.map(q => q.id)
    await api.delete('/api/v1/admin/questions/batch', { data: { ids } })
    ElMessage.success('批量删除成功')
    fetchQuestions()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('批量删除失败:', error)
      ElMessage.error('批量删除失败')
    }
  }
}

// 工具函数
const getTypeLabel = (type: string) => {
  const typeMap: Record<string, string> = {
    single: '单选题',
    multiple: '多选题',
    judge: '判断题',
    fill: '填空题'
  }
  return typeMap[type] || type
}

const getTypeTagType = (type: string) => {
  const typeMap: Record<string, string> = {
    single: 'primary',
    multiple: 'success',
    judge: 'warning',
    fill: 'info'
  }
  return typeMap[type] || 'info'
}

const getDifficultyLabel = (difficulty: string) => {
  const difficultyMap: Record<string, string> = {
    easy: '简单',
    medium: '中等',
    hard: '困难'
  }
  return difficultyMap[difficulty] || difficulty
}

const getDifficultyTagType = (difficulty: string) => {
  const difficultyMap: Record<string, string> = {
    easy: 'success',
    medium: 'warning',
    hard: 'danger'
  }
  return difficultyMap[difficulty] || 'info'
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleString('zh-CN')
}

// 组件挂载时获取数据
onMounted(() => {
  fetchQuestions()
  fetchCategories()
})
</script>

<style scoped>
.question-list {
  max-width: 1200px;
  margin: 0 auto;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 24px;
}

.header-left {
  flex: 1;
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

.header-right {
  display: flex;
  gap: 12px;
}

.search-card {
  margin-bottom: 20px;
}

.search-card :deep(.el-card__body) {
  padding: 20px;
}

.table-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-actions {
  display: flex;
  gap: 12px;
}

.question-content {
  padding: 8px 0;
}

.content-text {
  margin: 0 0 8px 0;
  line-height: 1.5;
  color: #1f2937;
}

.content-meta {
  display: flex;
  gap: 8px;
}

.pagination-container {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    gap: 16px;
  }
  
  .header-right {
    width: 100%;
    justify-content: flex-start;
  }
  
  .search-card :deep(.el-form) {
    flex-direction: column;
  }
  
  .search-card :deep(.el-form-item) {
    margin-right: 0;
    margin-bottom: 16px;
  }
}
</style>