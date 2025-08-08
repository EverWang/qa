<template>
  <div class="category-list">
    <div class="page-header">
      <h1 class="page-title">分类管理</h1>
      <div class="header-actions">
        <el-button type="primary" @click="handleCreate">
          <el-icon><Plus /></el-icon>
          添加分类
        </el-button>
      </div>
    </div>

    <div class="page-content">
      <!-- 搜索和筛选 -->
      <el-card class="search-card">
        <el-form :model="searchForm" inline>
          <el-form-item label="分类名称">
            <el-input
              v-model="searchForm.name"
              placeholder="请输入分类名称"
              clearable
              @keyup.enter="handleSearch"
              style="width: 200px"
            />
          </el-form-item>
          <el-form-item label="分类状态">
            <el-select
              v-model="searchForm.status"
              placeholder="请选择状态"
              clearable
              style="width: 120px"
            >
              <el-option label="启用" value="active" />
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

      <!-- 分类表格 -->
      <el-card class="table-card">
        <div class="table-header">
          <div class="table-info">
            <span>共 {{ total }} 个分类</span>
          </div>
          <div class="table-actions">
            <el-button 
              type="success" 
              @click="expandAll"
              :disabled="loading"
            >
              <el-icon><ArrowDown /></el-icon>
              展开全部
            </el-button>
            <el-button 
              @click="collapseAll"
              :disabled="loading"
            >
              <el-icon><ArrowUp /></el-icon>
              收起全部
            </el-button>
          </div>
        </div>
        
        <el-table
          ref="tableRef"
          :data="categories"
          v-loading="loading"
          row-key="id"
          :tree-props="{ children: 'children', hasChildren: 'hasChildren' }"
          :default-expand-all="false"
          border
          stripe
        >
          <el-table-column prop="name" label="分类名称" min-width="200">
            <template #default="{ row }">
              <div class="category-name">
                <el-icon v-if="row.level === 1" class="category-icon">
                  <Folder />
                </el-icon>
                <el-icon v-else class="category-icon">
                  <Document />
                </el-icon>
                <span>{{ row.name }}</span>
                <el-tag v-if="row.level === 1" type="primary" size="small">
                  大类
                </el-tag>
                <el-tag v-else type="info" size="small">
                  小类
                </el-tag>
              </div>
            </template>
          </el-table-column>
          
          <el-table-column prop="description" label="描述" min-width="200">
            <template #default="{ row }">
              <span class="description">{{ row.description || '暂无描述' }}</span>
            </template>
          </el-table-column>
          
          <el-table-column prop="question_count" label="题目数量" width="100" align="center">
            <template #default="{ row }">
              <el-tag type="success" size="small">
                {{ row.question_count || 0 }}
              </el-tag>
            </template>
          </el-table-column>
          
          <el-table-column prop="sort_order" label="排序" width="80" align="center">
            <template #default="{ row }">
              <span class="sort-order">{{ row.sort_order }}</span>
            </template>
          </el-table-column>
          
          <el-table-column prop="status" label="状态" width="80" align="center">
            <template #default="{ row }">
              <el-switch
                v-model="row.status"
                :active-value="1"
                :inactive-value="0"
                @change="handleStatusChange(row)"
              />
            </template>
          </el-table-column>
          
          <el-table-column prop="created_at" label="创建时间" width="160">
            <template #default="{ row }">
              <span>{{ formatDate(row.created_at) }}</span>
            </template>
          </el-table-column>
          
          <el-table-column label="操作" width="200" fixed="right">
            <template #default="{ row }">
              <div class="action-buttons">
                <el-button 
                  v-if="row.level === 1"
                  type="primary" 
                  text 
                  @click="handleAddChild(row)"
                >
                  <el-icon><Plus /></el-icon>
                  添加子类
                </el-button>
                <el-button type="primary" text @click="handleEdit(row)">
                  <el-icon><Edit /></el-icon>
                  编辑
                </el-button>
                <el-button 
                  type="danger" 
                  text 
                  @click="handleDelete(row)"
                  :disabled="row.question_count > 0"
                >
                  <el-icon><Delete /></el-icon>
                  删除
                </el-button>
              </div>
            </template>
          </el-table-column>
        </el-table>
        
        <!-- 树形数据不需要分页 -->
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
        <el-form-item label="分类名称" prop="name">
          <el-input
            v-model="form.name"
            placeholder="请输入分类名称"
            maxlength="50"
            show-word-limit
          />
        </el-form-item>
        
        <el-form-item label="上级分类" prop="parent_id">
          <el-select
            v-model="form.parent_id"
            placeholder="请选择上级分类（可选）"
            clearable
            style="width: 100%"
          >
            <el-option label="无（作为大类）" :value="null" />
            <el-option
              v-for="category in parentCategories"
              :key="category.id"
              :label="category.name"
              :value="category.id"
            />
          </el-select>
        </el-form-item>
        
        <el-form-item label="分类描述">
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="3"
            placeholder="请输入分类描述（可选）"
            maxlength="200"
            show-word-limit
          />
        </el-form-item>
        
        <el-form-item label="排序权重" prop="sort_order">
          <el-input-number
            v-model="form.sort_order"
            :min="0"
            :max="9999"
            placeholder="数值越大排序越靠前"
            style="width: 100%"
          />
        </el-form-item>
        
        <el-form-item label="分类状态">
          <el-radio-group v-model="form.status">
            <el-radio :label="1">启用</el-radio>
            <el-radio :label="0">禁用</el-radio>
          </el-radio-group>
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
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, nextTick } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance } from 'element-plus'
import {
  Plus,
  Search,
  Refresh,
  ArrowDown,
  ArrowUp,
  Folder,
  Document,
  Edit,
  Delete
} from '@element-plus/icons-vue'
import { api } from '@/lib/axios'

// 分类接口
interface Category {
  id: number
  name: string
  description?: string
  parent_id?: number
  level: number
  sort_order: number
  status: number // 1-启用 0-禁用
  question_count: number
  children?: Category[]
  hasChildren?: boolean
  created_at: string
  updated_at: string
}

// 响应式数据
const tableRef = ref()
const formRef = ref<FormInstance>()
const loading = ref(false)
const submitting = ref(false)
const dialogVisible = ref(false)
const isEdit = ref(false)
const categories = ref<Category[]>([])
const total = ref(0)

// 搜索表单
const searchForm = reactive({
  name: '',
  status: ''
})

// 移除分页相关代码，因为使用树形数据

// 表单数据
const form = reactive({
  id: 0,
  name: '',
  description: '',
  parent_id: null as number | null,
  sort_order: 0,
  status: 1 // 1-启用 0-禁用
})

// 表单验证规则
const rules = {
  name: [
    { required: true, message: '请输入分类名称', trigger: 'blur' },
    { min: 2, max: 50, message: '分类名称长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  sort_order: [
    { required: true, message: '请输入排序权重', trigger: 'blur' }
  ]
}

// 计算属性
const dialogTitle = computed(() => {
  return isEdit.value ? '编辑分类' : '添加分类'
})

const parentCategories = computed(() => {
  // 只显示大类作为父级选项
  return categories.value.filter(cat => cat.level === 1 && cat.id !== form.id)
})

// 方法
const handleCreate = () => {
  resetForm()
  isEdit.value = false
  dialogVisible.value = true
}

const handleAddChild = (parent: Category) => {
  resetForm()
  form.parent_id = parent.id
  isEdit.value = false
  dialogVisible.value = true
}

const handleEdit = (row: Category) => {
  form.id = row.id
  form.name = row.name
  form.description = row.description || ''
  form.parent_id = row.parent_id || null
  form.sort_order = row.sort_order
  form.status = row.status
  isEdit.value = true
  dialogVisible.value = true
}

const handleDelete = async (row: Category) => {
  if (row.question_count > 0) {
    ElMessage.warning('该分类下还有题目，无法删除')
    return
  }
  
  try {
    await ElMessageBox.confirm(
      `确定要删除分类「${row.name}」吗？`,
      '删除确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await api.delete(`/api/v1/admin/categories/${row.id}`)
    ElMessage.success('删除成功')
    await fetchCategories()
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('删除分类失败:', error)
      ElMessage.error(error.response?.data?.message || '删除失败')
    }
  }
}

const handleStatusChange = async (row: Category) => {
  const originalStatus = row.status
  try {
    await api.put(`/api/v1/admin/categories/${row.id}/status`, {
      status: row.status
    })
    ElMessage.success('状态更新成功')
  } catch (error: any) {
    console.error('更新状态失败:', error)
    ElMessage.error('状态更新失败')
    // 恢复原状态
    row.status = originalStatus === 1 ? 0 : 1
  }
}

const handleSearch = () => {
  fetchCategories()
}

const handleReset = () => {
  searchForm.name = ''
  searchForm.status = ''
  fetchCategories()
}

const expandAll = async () => {
  if (!tableRef.value) return
  
  // 展开所有节点
  const expandNode = (nodes: Category[]) => {
    nodes.forEach(node => {
      if (node.children && node.children.length > 0) {
        tableRef.value.toggleRowExpansion(node, true)
        expandNode(node.children)
      }
    })
  }
  
  await nextTick()
  expandNode(categories.value)
}

const collapseAll = async () => {
  if (!tableRef.value) return
  
  // 收起所有节点
  const collapseNode = (nodes: Category[]) => {
    nodes.forEach(node => {
      tableRef.value.toggleRowExpansion(node, false)
      if (node.children && node.children.length > 0) {
        collapseNode(node.children)
      }
    })
  }
  
  await nextTick()
  collapseNode(categories.value)
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
    
    const categoryData = {
      name: form.name,
      description: form.description || undefined,
      parent_id: form.parent_id || undefined,
      sort_order: form.sort_order,
      status: form.status
    }
    
    if (isEdit.value) {
      await api.put(`/api/v1/admin/categories/${form.id}`, categoryData)
      ElMessage.success('更新成功')
    } else {
      await api.post('/api/v1/admin/categories', categoryData)
      ElMessage.success('创建成功')
    }
    
    dialogVisible.value = false
    await fetchCategories()
  } catch (error: any) {
    console.error('保存分类失败:', error)
    ElMessage.error(error.response?.data?.message || '保存失败')
  } finally {
    submitting.value = false
  }
}

const resetForm = () => {
  form.id = 0
  form.name = ''
  form.description = ''
  form.parent_id = null
  form.sort_order = 0
  form.status = 1 // 1-启用
  
  if (formRef.value) {
    formRef.value.clearValidate()
  }
}

const fetchCategories = async () => {
  try {
    loading.value = true
    
    const params = {
      tree: 'true', // 获取树形数据以支持展开/收起功能
      search: searchForm.name || undefined,
      status: searchForm.status || undefined
    }
    
    const response = await api.get('/api/v1/admin/categories', { params })
    const data = response.data
    
    // 树形数据结构
    categories.value = data.data || []
    total.value = data.total || categories.value.length
    
  } catch (error) {
    console.error('获取分类列表失败:', error)
    ElMessage.error('获取分类列表失败')
  } finally {
    loading.value = false
  }
}

const formatDate = (dateString: string) => {
  if (!dateString) return '-'
  return new Date(dateString).toLocaleString('zh-CN')
}

// 生命周期
onMounted(() => {
  fetchCategories()
})
</script>

<style scoped>
.category-list {
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

.category-name {
  display: flex;
  align-items: center;
  gap: 8px;
}

.category-icon {
  color: #6b7280;
}

.description {
  color: #6b7280;
  font-size: 14px;
}

.sort-order {
  font-weight: 600;
  color: #374151;
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

/* 响应式设计 */
@media (max-width: 768px) {
  .category-list {
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
}
</style>