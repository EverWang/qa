<template>
  <div class="operation-logs">
    <div class="page-header">
      <h1 class="page-title">操作日志</h1>
    </div>

    <div class="page-content">
      <!-- 搜索筛选 -->
      <el-card class="search-card">
        <el-form :model="searchForm" inline>
          <el-form-item label="操作人">
            <el-input
              v-model="searchForm.operator"
              placeholder="请输入操作人"
              clearable
              style="width: 200px"
            />
          </el-form-item>
          <el-form-item label="操作类型">
            <el-select
              v-model="searchForm.action"
              placeholder="请选择操作类型"
              clearable
              style="width: 150px"
            >
              <el-option label="创建" value="create" />
              <el-option label="更新" value="update" />
              <el-option label="删除" value="delete" />
              <el-option label="登录" value="login" />
              <el-option label="登出" value="logout" />
            </el-select>
          </el-form-item>
          <el-form-item label="时间范围">
            <el-date-picker
              v-model="searchForm.dateRange"
              type="datetimerange"
              range-separator="至"
              start-placeholder="开始时间"
              end-placeholder="结束时间"
              format="YYYY-MM-DD HH:mm:ss"
              value-format="YYYY-MM-DD HH:mm:ss"
            />
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

      <!-- 日志列表 -->
      <el-card class="table-card">
        <el-table
          v-loading="loading"
          :data="logs"
          stripe
          style="width: 100%"
        >
          <el-table-column prop="id" label="ID" width="80" />
          <el-table-column prop="operator" label="操作人" width="120" />
          <el-table-column prop="action" label="操作类型" width="100">
            <template #default="{ row }">
              <el-tag
                :type="getActionTagType(row.action)"
                size="small"
              >
                {{ getActionText(row.action) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="resource" label="操作对象" width="150" />
          <el-table-column prop="description" label="操作描述" min-width="200" />
          <el-table-column prop="ip" label="IP地址" width="120" />
          <el-table-column prop="userAgent" label="用户代理" width="200" show-overflow-tooltip />
          <el-table-column prop="createdAt" label="操作时间" width="180">
            <template #default="{ row }">
              {{ formatDateTime(row.createdAt) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="100" fixed="right">
            <template #default="{ row }">
              <el-button
                type="primary"
                size="small"
                text
                @click="handleViewDetail(row)"
              >
                详情
              </el-button>
            </template>
          </el-table-column>
        </el-table>

        <!-- 分页 -->
        <div class="pagination-wrapper">
          <el-pagination
            v-model:current-page="pagination.page"
            v-model:page-size="pagination.pageSize"
            :total="pagination.total"
            :page-sizes="[10, 20, 50, 100]"
            layout="total, sizes, prev, pager, next, jumper"
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
          />
        </div>
      </el-card>
    </div>

    <!-- 详情对话框 -->
    <el-dialog
      v-model="detailDialogVisible"
      title="操作日志详情"
      width="600px"
    >
      <el-descriptions :column="1" border>
        <el-descriptions-item label="操作人">
          {{ currentLog?.operator }}
        </el-descriptions-item>
        <el-descriptions-item label="操作类型">
          <el-tag :type="getActionTagType(currentLog?.action)" size="small">
            {{ getActionText(currentLog?.action) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="操作对象">
          {{ currentLog?.resource }}
        </el-descriptions-item>
        <el-descriptions-item label="操作描述">
          {{ currentLog?.description }}
        </el-descriptions-item>
        <el-descriptions-item label="IP地址">
          {{ currentLog?.ip }}
        </el-descriptions-item>
        <el-descriptions-item label="用户代理">
          {{ currentLog?.userAgent }}
        </el-descriptions-item>
        <el-descriptions-item label="操作时间">
          {{ formatDateTime(currentLog?.createdAt) }}
        </el-descriptions-item>
        <el-descriptions-item label="请求参数" v-if="currentLog?.requestData">
          <pre>{{ JSON.stringify(JSON.parse(currentLog.requestData), null, 2) }}</pre>
        </el-descriptions-item>
        <el-descriptions-item label="响应数据" v-if="currentLog?.responseData">
          <pre>{{ JSON.stringify(JSON.parse(currentLog.responseData), null, 2) }}</pre>
        </el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Search, Refresh } from '@element-plus/icons-vue'
import { api } from '@/lib/axios'

// 接口类型定义
interface OperationLog {
  id: number
  operator: string
  action: string
  resource: string
  description: string
  ip: string
  userAgent: string
  requestData?: string
  responseData?: string
  createdAt: string
}

// 响应式数据
const loading = ref(false)
const logs = ref<OperationLog[]>([])
const detailDialogVisible = ref(false)
const currentLog = ref<OperationLog | null>(null)

// 搜索表单
const searchForm = reactive({
  operator: '',
  action: '',
  dateRange: [] as string[]
})

// 分页信息
const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

/**
 * 获取操作日志列表
 */
const fetchLogs = async () => {
  try {
    loading.value = true
    const params = {
      page: pagination.page,
      pageSize: pagination.pageSize,
      operator: searchForm.operator || undefined,
      action: searchForm.action || undefined,
      startTime: searchForm.dateRange[0] || undefined,
      endTime: searchForm.dateRange[1] || undefined
    }

    const response = await api.get('/api/v1/admin/operation-logs', { params })
    
    if (response && response.data && response.data.code === 200) {
      logs.value = response.data.data.data || []
      pagination.total = response.data.data.total || 0
    } else {
      logs.value = []
      pagination.total = 0
      ElMessage.error('获取操作日志失败')
    }
  } catch (error) {
    console.error('获取操作日志失败:', error)
    ElMessage.error('获取操作日志失败: ' + (error.response?.data?.message || error.message))
    logs.value = []
    pagination.total = 0
  } finally {
    loading.value = false
  }
}

/**
 * 搜索操作
 */
const handleSearch = () => {
  pagination.page = 1
  fetchLogs()
}

/**
 * 重置搜索
 */
const handleReset = () => {
  searchForm.operator = ''
  searchForm.action = ''
  searchForm.dateRange = []
  pagination.page = 1
  fetchLogs()
}

/**
 * 查看详情
 */
const handleViewDetail = (log: OperationLog) => {
  currentLog.value = log
  detailDialogVisible.value = true
}

/**
 * 分页大小改变
 */
const handleSizeChange = (size: number) => {
  pagination.pageSize = size
  pagination.page = 1
  fetchLogs()
}

/**
 * 当前页改变
 */
const handleCurrentChange = (page: number) => {
  pagination.page = page
  fetchLogs()
}

/**
 * 获取操作类型标签类型
 */
const getActionTagType = (action: string) => {
  const typeMap: Record<string, string> = {
    create: 'success',
    update: 'warning',
    delete: 'danger',
    login: 'info',
    logout: 'info'
  }
  return typeMap[action] || 'info'
}

/**
 * 获取操作类型文本
 */
const getActionText = (action: string) => {
  const textMap: Record<string, string> = {
    create: '创建',
    update: '更新',
    delete: '删除',
    login: '登录',
    logout: '登出'
  }
  return textMap[action] || action
}

/**
 * 格式化日期时间
 */
const formatDateTime = (dateTime: string | undefined) => {
  if (!dateTime) return '-'
  return new Date(dateTime).toLocaleString('zh-CN')
}

// 组件挂载时获取数据
onMounted(() => {
  fetchLogs()
})
</script>

<style scoped lang="scss">
.operation-logs {
  .page-header {
    margin-bottom: 20px;
    
    .page-title {
      font-size: 24px;
      font-weight: 600;
      color: #303133;
      margin: 0;
    }
  }

  .page-content {
    .search-card {
      margin-bottom: 20px;
    }

    .table-card {
      .pagination-wrapper {
        display: flex;
        justify-content: center;
        margin-top: 20px;
      }
    }
  }

  pre {
    background-color: #f5f5f5;
    padding: 10px;
    border-radius: 4px;
    font-size: 12px;
    max-height: 200px;
    overflow-y: auto;
  }
}
</style>