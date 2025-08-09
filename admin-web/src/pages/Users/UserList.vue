<template>
  <div class="user-list">
    <div class="page-header">
      <h1 class="page-title">用户管理</h1>
      <div class="header-actions">
        <el-button type="primary" @click="handleCreate">
          <el-icon><Plus /></el-icon>
          添加用户
        </el-button>
      </div>
    </div>

    <div class="page-content">
      <!-- 搜索和筛选 -->
      <el-card class="search-card">
        <el-form :model="searchForm" inline>
          <el-form-item label="用户名">
            <el-input
              v-model="searchForm.username"
              placeholder="请输入用户名"
              clearable
              @keyup.enter="handleSearch"
              style="width: 200px"
            />
          </el-form-item>
          <el-form-item label="邮箱">
            <el-input
              v-model="searchForm.email"
              placeholder="请输入邮箱"
              clearable
              @keyup.enter="handleSearch"
              style="width: 200px"
            />
          </el-form-item>
          <el-form-item label="角色">
            <el-select
              v-model="searchForm.role"
              placeholder="请选择角色"
              clearable
              style="width: 120px"
            >
              <el-option label="管理员" value="admin" />
              <el-option label="普通用户" value="user" />
            </el-select>
          </el-form-item>
          <el-form-item label="状态">
            <el-select
              v-model="searchForm.status"
              placeholder="请选择状态"
              clearable
              style="width: 120px"
            >
              <el-option label="正常" value="active" />
              <el-option label="禁用" value="inactive" />
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

      <!-- 用户表格 -->
      <el-card class="table-card">
        <div class="table-header">
          <div class="table-info">
            <span>共 {{ total }} 个用户</span>
          </div>
          <div class="table-actions">
            <el-button 
              type="danger" 
              @click="handleBatchDelete"
              :disabled="selectedUsers.length === 0"
            >
              <el-icon><Delete /></el-icon>
              批量删除
            </el-button>
          </div>
        </div>
        
        <el-table
          :data="users"
          v-loading="loading"
          @selection-change="handleSelectionChange"
          border
          stripe
        >
          <el-table-column type="selection" width="55" />
          
          <el-table-column prop="avatar" label="头像" width="80">
            <template #default="{ row }">
              <el-avatar 
                :src="row.avatar" 
                :size="40"
                class="user-avatar"
              >
                <el-icon><User /></el-icon>
              </el-avatar>
            </template>
          </el-table-column>
          
          <el-table-column prop="username" label="用户名" min-width="120">
            <template #default="{ row }">
              <div class="username-cell">
                <span class="username">{{ row.username }}</span>
                <el-tag v-if="row.is_verified" type="success" size="small">
                  已验证
                </el-tag>
              </div>
            </template>
          </el-table-column>
          
          <el-table-column prop="email" label="邮箱" min-width="180">
            <template #default="{ row }">
              <span>{{ row.email }}</span>
            </template>
          </el-table-column>
          
          <el-table-column prop="role" label="角色" width="100">
            <template #default="{ row }">
              <el-tag :type="getRoleTagType(row.role)" size="small">
                {{ getRoleLabel(row.role) }}
              </el-tag>
            </template>
          </el-table-column>
          
          <el-table-column prop="status" label="状态" width="100">
            <template #default="{ row }">
              <el-switch
                v-model="row.status"
                active-value="active"
                inactive-value="inactive"
                @change="handleStatusChange(row)"
              />
            </template>
          </el-table-column>
          
          <el-table-column prop="last_login_at" label="最后登录" width="160">
            <template #default="{ row }">
              <span>{{ formatDate(row.last_login_at) }}</span>
            </template>
          </el-table-column>
          
          <el-table-column prop="created_at" label="注册时间" width="160">
            <template #default="{ row }">
              <span>{{ formatDate(row.created_at) }}</span>
            </template>
          </el-table-column>
          
          <el-table-column label="操作" width="200" fixed="right">
            <template #default="{ row }">
              <div class="action-buttons">
                <el-button type="primary" text @click="handleView(row)">
                  <el-icon><View /></el-icon>
                  查看
                </el-button>
                <el-button type="primary" text @click="handleEdit(row)">
                  <el-icon><Edit /></el-icon>
                  编辑
                </el-button>
                <el-button 
                  type="danger" 
                  text 
                  @click="handleDelete(row)"
                  :disabled="row.role === 'admin' && row.id === currentUserId"
                >
                  <el-icon><Delete /></el-icon>
                  删除
                </el-button>
              </div>
            </template>
          </el-table-column>
        </el-table>
        
        <!-- 分页 -->
        <div class="pagination-wrapper">
          <el-pagination
            v-model:current-page="pagination.page"
            v-model:page-size="pagination.size"
            :page-sizes="[10, 20, 50, 100]"
            :total="total"
            layout="total, sizes, prev, pager, next, jumper"
            @size-change="handleSizeChange"
            @current-change="handlePageChange"
          />
        </div>
      </el-card>
    </div>

    <!-- 创建/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="600px"
      :before-close="handleDialogClose"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="100px"
      >
        <el-form-item label="用户名" prop="username">
          <el-input
            v-model="form.username"
            placeholder="请输入用户名"
            :disabled="isEdit"
            maxlength="50"
          />
        </el-form-item>
        
        <el-form-item label="邮箱" prop="email">
          <el-input
            v-model="form.email"
            placeholder="请输入邮箱"
            type="email"
            maxlength="100"
          />
        </el-form-item>
        
        <el-form-item v-if="!isEdit" label="密码" prop="password">
          <el-input
            v-model="form.password"
            placeholder="请输入密码"
            type="password"
            show-password
            maxlength="50"
          />
        </el-form-item>
        
        <el-form-item v-if="!isEdit" label="确认密码" prop="confirmPassword">
          <el-input
            v-model="form.confirmPassword"
            placeholder="请再次输入密码"
            type="password"
            show-password
            maxlength="50"
          />
        </el-form-item>
        
        <el-form-item label="角色" prop="role">
          <el-select v-model="form.role" placeholder="请选择角色" style="width: 100%">
            <el-option label="管理员" value="admin" />
            <el-option label="普通用户" value="user" />
          </el-select>
        </el-form-item>
        
        <el-form-item label="状态">
          <el-radio-group v-model="form.status">
            <el-radio label="active">正常</el-radio>
            <el-radio label="inactive">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
        
        <el-form-item label="邮箱验证">
          <el-checkbox v-model="form.is_verified">已验证邮箱</el-checkbox>
        </el-form-item>
      </el-form>
      
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="handleDialogClose">取消</el-button>
          <el-button 
            type="primary" 
            @click="handleSubmit"
            :loading="submitting"
          >
            确定
          </el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 用户详情对话框 -->
    <el-dialog
      v-model="detailVisible"
      title="用户详情"
      width="600px"
    >
      <div v-if="selectedUser" class="user-detail">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="头像">
            <el-avatar :src="selectedUser.avatar" :size="60">
              <el-icon><User /></el-icon>
            </el-avatar>
          </el-descriptions-item>
          <el-descriptions-item label="用户名">
            {{ selectedUser.username }}
          </el-descriptions-item>
          <el-descriptions-item label="邮箱">
            {{ selectedUser.email }}
          </el-descriptions-item>
          <el-descriptions-item label="角色">
            <el-tag :type="getRoleTagType(selectedUser.role)">
              {{ getRoleLabel(selectedUser.role) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="selectedUser.status === 'active' ? 'success' : 'danger'">
              {{ selectedUser.status === 'active' ? '正常' : '禁用' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="邮箱验证">
            <el-tag :type="selectedUser.is_verified ? 'success' : 'warning'">
              {{ selectedUser.is_verified ? '已验证' : '未验证' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="注册时间">
            {{ formatDate(selectedUser.created_at) }}
          </el-descriptions-item>
          <el-descriptions-item label="最后登录">
            {{ formatDate(selectedUser.last_login_at) }}
          </el-descriptions-item>
          <el-descriptions-item label="答题统计" span="2">
            <div class="stats-info">
              <span>总答题: {{ selectedUser.total_answers || 0 }}次</span>
              <span>正确率: {{ ((selectedUser.correct_answers || 0) / (selectedUser.total_answers || 1) * 100).toFixed(1) }}%</span>
            </div>
          </el-descriptions-item>
        </el-descriptions>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance } from 'element-plus'
import {
  Plus,
  Search,
  Refresh,
  Delete,
  User,
  View,
  Edit
} from '@element-plus/icons-vue'
import api from '@/lib/axios'
import { useAuthStore } from '@/stores/auth'

// 用户接口
interface User {
  id: number
  username: string
  email: string
  avatar?: string
  role: 'admin' | 'user'
  status: 'active' | 'inactive'
  is_verified: boolean
  last_login_at?: string
  created_at: string
  updated_at: string
  total_answers?: number
  correct_answers?: number
}

// 响应式数据
const authStore = useAuthStore()
const formRef = ref<FormInstance>()
const loading = ref(false)
const submitting = ref(false)
const dialogVisible = ref(false)
const detailVisible = ref(false)
const isEdit = ref(false)
const users = ref<User[]>([])
const selectedUsers = ref<User[]>([])
const selectedUser = ref<User | null>(null)
const total = ref(0)

// 搜索表单
const searchForm = reactive({
  username: '',
  email: '',
  role: '',
  status: ''
})

// 分页
const pagination = reactive({
  page: 1,
  size: 20
})

// 表单数据
const form = reactive({
  id: 0,
  username: '',
  email: '',
  password: '',
  confirmPassword: '',
  role: 'user' as 'admin' | 'user',
  status: 'active' as 'active' | 'inactive',
  is_verified: false
})

// 表单验证规则
const rules = computed(() => ({
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 50, message: '用户名长度在 3 到 50 个字符', trigger: 'blur' },
    { pattern: /^[a-zA-Z0-9_]+$/, message: '用户名只能包含字母、数字和下划线', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  password: isEdit.value ? [] : [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 50, message: '密码长度在 6 到 50 个字符', trigger: 'blur' }
  ],
  confirmPassword: isEdit.value ? [] : [
    { required: true, message: '请再次输入密码', trigger: 'blur' },
    {
      validator: (rule: any, value: string, callback: Function) => {
        if (value !== form.password) {
          callback(new Error('两次输入的密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ],
  role: [
    { required: true, message: '请选择角色', trigger: 'change' }
  ]
}))

// 计算属性
const dialogTitle = computed(() => {
  return isEdit.value ? '编辑用户' : '添加用户'
})

const currentUserId = computed(() => {
  return authStore.user?.id
})

// 方法
const handleCreate = () => {
  resetForm()
  isEdit.value = false
  dialogVisible.value = true
}

const handleView = (row: User) => {
  selectedUser.value = row
  detailVisible.value = true
}

const handleEdit = (row: User) => {
  form.id = row.id
  form.username = row.username
  form.email = row.email
  form.password = ''
  form.confirmPassword = ''
  form.role = row.role
  form.status = row.status
  form.is_verified = row.is_verified
  isEdit.value = true
  dialogVisible.value = true
}

const handleDelete = async (row: User) => {
  if (row.role === 'admin' && row.id === currentUserId.value) {
    ElMessage.warning('不能删除自己的账号')
    return
  }
  
  try {
    await ElMessageBox.confirm(
      `确定要删除用户「${row.username}」吗？`,
      '删除确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await api.delete(`/api/v1/admin/users/${row.id}`)
    ElMessage.success('删除成功')
    await fetchUsers()
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('删除用户失败:', error)
      ElMessage.error(error.response?.data?.message || '删除失败')
    }
  }
}

const handleBatchDelete = async () => {
  if (selectedUsers.value.length === 0) {
    ElMessage.warning('请选择要删除的用户')
    return
  }
  
  // 检查是否包含当前用户
  const hasCurrentUser = selectedUsers.value.some(user => user.id === currentUserId.value)
  if (hasCurrentUser) {
    ElMessage.warning('不能删除自己的账号')
    return
  }
  
  try {
    await ElMessageBox.confirm(
      `确定要删除选中的 ${selectedUsers.value.length} 个用户吗？`,
      '批量删除确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    const userIds = selectedUsers.value.map(user => user.id)
    await api.delete('/api/v1/admin/users/batch', { data: { ids: userIds } })
    
    ElMessage.success('批量删除成功')
    await fetchUsers()
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('批量删除失败:', error)
      ElMessage.error(error.response?.data?.message || '批量删除失败')
    }
  }
}

const handleStatusChange = async (row: User) => {
  try {
    await api.put(`/api/v1/admin/users/${row.id}/status`, {
      status: row.status
    })
    ElMessage.success('状态更新成功')
  } catch (error: any) {
    console.error('更新状态失败:', error)
    ElMessage.error('状态更新失败')
    // 恢复原状态
    row.status = row.status === 'active' ? 'inactive' : 'active'
  }
}

const handleSelectionChange = (selection: User[]) => {
  selectedUsers.value = selection
}

const handleSearch = () => {
  pagination.page = 1
  fetchUsers()
}

const handleReset = () => {
  searchForm.username = ''
  searchForm.email = ''
  searchForm.role = ''
  searchForm.status = ''
  pagination.page = 1
  fetchUsers()
}

const handlePageChange = (page: number) => {
  pagination.page = page
  fetchUsers()
}

const handleSizeChange = (size: number) => {
  pagination.size = size
  pagination.page = 1
  fetchUsers()
}

const handleDialogClose = () => {
  dialogVisible.value = false
  resetForm()
}

const handleSubmit = async () => {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    submitting.value = true
    
    const userData = {
      username: form.username,
      email: form.email,
      role: form.role,
      status: form.status,
      is_verified: form.is_verified,
      ...(isEdit.value ? {} : { password: form.password })
    }
    
    if (isEdit.value) {
      await api.put(`/api/v1/admin/users/${form.id}`, userData)
      ElMessage.success('更新成功')
    } else {
      await api.post('/api/v1/admin/users', userData)
      ElMessage.success('创建成功')
    }
    
    dialogVisible.value = false
    await fetchUsers()
  } catch (error: any) {
    console.error('保存用户失败:', error)
    ElMessage.error(error.response?.data?.message || '保存失败')
  } finally {
    submitting.value = false
  }
}

const resetForm = () => {
  form.id = 0
  form.username = ''
  form.email = ''
  form.password = ''
  form.confirmPassword = ''
  form.role = 'user'
  form.status = 'active'
  form.is_verified = false
  
  if (formRef.value) {
    formRef.value.clearValidate()
  }
}

const fetchUsers = async () => {
  try {
    loading.value = true
    
    const params = {
      page: pagination.page,
      size: pagination.size,
      username: searchForm.username || undefined,
      email: searchForm.email || undefined,
      role: searchForm.role || undefined,
      status: searchForm.status || undefined
    }
    
    const response = await api.get('/api/v1/admin/users', { params })
    const data = response.data
    
    users.value = data.data || []
    total.value = data.total || 0
    
  } catch (error) {
    console.error('获取用户列表失败:', error)
    ElMessage.error('获取用户列表失败')
  } finally {
    loading.value = false
  }
}

const getRoleLabel = (role: string) => {
  const roleMap: Record<string, string> = {
    admin: '管理员',
    user: '普通用户'
  }
  return roleMap[role] || role
}

const getRoleTagType = (role: string) => {
  const roleMap: Record<string, string> = {
    admin: 'danger',
    user: 'primary'
  }
  return roleMap[role] || 'info'
}

const formatDate = (dateString?: string) => {
  if (!dateString) return '-'
  return new Date(dateString).toLocaleString('zh-CN')
}

// 生命周期
onMounted(() => {
  fetchUsers()
})
</script>

<style scoped>
.user-list {
  padding: 24px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.page-title {
  font-size: 24px;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
}

.header-actions {
  display: flex;
  gap: 12px;
}

.page-content {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.search-card {
  padding: 16px;
}

.table-card {
  flex: 1;
}

.table-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  padding: 0 4px;
}

.table-info {
  color: #6b7280;
  font-size: 14px;
}

.table-actions {
  display: flex;
  gap: 8px;
}

.user-avatar {
  border: 2px solid #e5e7eb;
}

.username-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.username {
  font-weight: 500;
}

.action-buttons {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: 24px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.user-detail {
  padding: 16px 0;
}

.stats-info {
  display: flex;
  gap: 16px;
  color: #6b7280;
  font-size: 14px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .user-list {
    padding: 16px;
  }
  
  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }
  
  .header-actions {
    width: 100%;
    justify-content: flex-end;
  }
  
  .table-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }
  
  .table-actions {
    width: 100%;
    justify-content: flex-end;
  }
  
  .action-buttons {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .pagination-wrapper {
    overflow-x: auto;
  }
  
  .stats-info {
    flex-direction: column;
    gap: 8px;
  }
}
</style>