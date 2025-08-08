<template>
  <div class="question-create">
    <div class="page-header">
      <div class="header-left">
        <el-button @click="$router.back()">
          <el-icon><ArrowLeft /></el-icon>
          返回
        </el-button>
        <div class="header-title">
          <h1 class="page-title">添加题目</h1>
          <p class="page-subtitle">创建新的题目，支持单选、多选、判断和填空题</p>
        </div>
      </div>
      <div class="header-right">
        <el-button @click="handlePreview">
          <el-icon><View /></el-icon>
          预览
        </el-button>
        <el-button type="primary" @click="handleSave" :loading="saving">
          <el-icon><Check /></el-icon>
          保存题目
        </el-button>
      </div>
    </div>
    
    <el-card class="form-card" shadow="never">
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="120px"
        size="large"
      >
        <!-- 基本信息 -->
        <div class="form-section">
          <h3 class="section-title">基本信息</h3>
          
          <el-form-item label="题目类型" prop="type">
            <el-radio-group v-model="form.type" @change="handleTypeChange">
              <el-radio value="single">单选题</el-radio>
              <el-radio value="multiple">多选题</el-radio>
              <el-radio value="judge">判断题</el-radio>
              <el-radio value="fill">填空题</el-radio>
            </el-radio-group>
          </el-form-item>
          
          <el-form-item label="题目分类" prop="categoryId">
            <el-select
              v-model="form.categoryId"
              placeholder="请选择题目分类"
              style="width: 300px"
            >
              <el-option
                v-for="category in categories"
                :key="category.id"
                :label="category.name"
                :value="category.id"
              />
            </el-select>
          </el-form-item>
          
          <el-form-item label="难度等级" prop="difficulty">
            <el-radio-group v-model="form.difficulty">
              <el-radio value="easy">简单</el-radio>
              <el-radio value="medium">中等</el-radio>
              <el-radio value="hard">困难</el-radio>
            </el-radio-group>
          </el-form-item>
        </div>
        
        <!-- 题目内容 -->
        <div class="form-section">
          <h3 class="section-title">题目内容</h3>
          
          <el-form-item label="题目描述" prop="content">
            <el-input
              v-model="form.content"
              type="textarea"
              :rows="4"
              placeholder="请输入题目内容"
              maxlength="1000"
              show-word-limit
            />
          </el-form-item>
        </div>
        
        <!-- 选项设置（单选题和多选题） -->
        <div v-if="form.type === 'single' || form.type === 'multiple'" class="form-section">
          <h3 class="section-title">选项设置</h3>
          
          <div class="options-container">
            <div
              v-for="(option, index) in form.options"
              :key="index"
              class="option-item"
            >
              <div class="option-input">
                <span class="option-label">{{ getOptionLabel(index) }}.</span>
                <el-input
                  v-model="form.options[index]"
                  placeholder="请输入选项内容"
                  maxlength="200"
                />
              </div>
              <div class="option-actions">
                <el-button
                  v-if="form.options.length > 2"
                  type="danger"
                  link
                  @click="removeOption(index)"
                >
                  <el-icon><Delete /></el-icon>
                </el-button>
              </div>
            </div>
            
            <el-button
              v-if="form.options.length < 6"
              type="primary"
              link
              @click="addOption"
            >
              <el-icon><Plus /></el-icon>
              添加选项
            </el-button>
          </div>
          
          <el-form-item label="正确答案" prop="answer">
            <template v-if="form.type === 'single'">
              <el-radio-group v-model="form.answer">
                <el-radio
                  v-for="(option, index) in form.options"
                  :key="index"
                  :value="getOptionLabel(index)"
                  :disabled="!option.trim()"
                >
                  {{ getOptionLabel(index) }}. {{ option || '(空选项)' }}
                </el-radio>
              </el-radio-group>
            </template>
            <template v-else>
              <el-checkbox-group v-model="form.answer">
                <el-checkbox
                  v-for="(option, index) in form.options"
                  :key="index"
                  :value="getOptionLabel(index)"
                  :disabled="!option.trim()"
                >
                  {{ getOptionLabel(index) }}. {{ option || '(空选项)' }}
                </el-checkbox>
              </el-checkbox-group>
            </template>
          </el-form-item>
        </div>
        
        <!-- 判断题答案 -->
        <div v-if="form.type === 'judge'" class="form-section">
          <h3 class="section-title">正确答案</h3>
          
          <el-form-item label="判断结果" prop="answer">
            <el-radio-group v-model="form.answer">
              <el-radio value="true">正确</el-radio>
              <el-radio value="false">错误</el-radio>
            </el-radio-group>
          </el-form-item>
        </div>
        
        <!-- 填空题答案 -->
        <div v-if="form.type === 'fill'" class="form-section">
          <h3 class="section-title">正确答案</h3>
          
          <el-form-item label="答案内容" prop="answer">
            <el-input
              v-model="form.answer"
              placeholder="请输入正确答案，多个答案用英文逗号分隔"
              maxlength="500"
            />
            <div class="form-tip">
              提示：如果有多个正确答案，请用英文逗号分隔，如：答案1,答案2,答案3
            </div>
          </el-form-item>
        </div>
        
        <!-- 题目解析 -->
        <div class="form-section">
          <h3 class="section-title">题目解析</h3>
          
          <el-form-item label="解析内容" prop="explanation">
            <el-input
              v-model="form.explanation"
              type="textarea"
              :rows="3"
              placeholder="请输入题目解析（可选）"
              maxlength="500"
              show-word-limit
            />
          </el-form-item>
        </div>
      </el-form>
    </el-card>
    
    <!-- 预览对话框 -->
    <el-dialog
      v-model="previewVisible"
      title="题目预览"
      width="800px"
      destroy-on-close
    >
      <QuestionPreview :question="previewQuestion" />
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElForm } from 'element-plus'
import {
  ArrowLeft,
  View,
  Check,
  Delete,
  Plus
} from '@element-plus/icons-vue'
import { api } from '@/lib/axios'
import QuestionPreview from './QuestionPreview.vue'

const router = useRouter()

// 表单引用
const formRef = ref<InstanceType<typeof ElForm>>()

// 分类数据
interface Category {
  id: number
  name: string
}

const categories = ref<Category[]>([])

// 表单数据
const form = reactive({
  type: 'single' as 'single' | 'multiple' | 'judge' | 'fill',
  categoryId: '',
  difficulty: 'easy' as 'easy' | 'medium' | 'hard',
  content: '',
  options: ['', ''],
  answer: '' as string | string[],
  explanation: ''
})

// 保存状态
const saving = ref(false)

// 预览
const previewVisible = ref(false)
const previewQuestion = computed(() => ({
  id: 0,
  content: form.content,
  type: form.type as 'single' | 'multiple' | 'judge' | 'fill',
  difficulty: form.difficulty as 'easy' | 'medium' | 'hard',
  options: form.type === 'single' || form.type === 'multiple' ? form.options.filter(opt => opt.trim()) : undefined,
  answer: form.answer,
  explanation: form.explanation,
  category_id: form.categoryId,
  category: categories.value.find(c => c.id.toString() === form.categoryId),
  created_at: new Date().toISOString(),
  updated_at: new Date().toISOString()
}))

// 表单验证规则
const rules = {
  type: [
    { required: true, message: '请选择题目类型', trigger: 'change' }
  ],
  categoryId: [
    { required: true, message: '请选择题目分类', trigger: 'change' }
  ],
  difficulty: [
    { required: true, message: '请选择难度等级', trigger: 'change' }
  ],
  content: [
    { required: true, message: '请输入题目内容', trigger: 'blur' },
    { min: 5, message: '题目内容至少5个字符', trigger: 'blur' }
  ],
  answer: [
    { required: true, message: '请设置正确答案', trigger: 'change' }
  ]
}

// 获取分类列表
const fetchCategories = async () => {
  try {
    const response = await api.get('/api/v1/admin/categories')
    categories.value = response.data.data
  } catch (error) {
    console.error('获取分类列表失败:', error)
    // 使用模拟数据
    categories.value = [
      { id: 1, name: '前端开发' },
      { id: 2, name: '后端开发' },
      { id: 3, name: '数据库' },
      { id: 4, name: '算法与数据结构' }
    ]
  }
}

// 题目类型变化处理
const handleTypeChange = (type: string) => {
  // 重置答案
  if (type === 'single' || type === 'judge' || type === 'fill') {
    form.answer = ''
  } else if (type === 'multiple') {
    form.answer = [] as string[]
  }
  
  // 重置选项
  if (type === 'single' || type === 'multiple') {
    if (form.options.length < 2) {
      form.options = ['', '']
    }
  }
}

// 获取选项标签
const getOptionLabel = (index: number) => {
  const labels = ['A', 'B', 'C', 'D', 'E', 'F']
  return labels[index] || (index + 1).toString()
}

// 添加选项
const addOption = () => {
  if (form.options.length < 6) {
    form.options.push('')
  }
}

// 删除选项
const removeOption = (index: number) => {
  if (form.options.length > 2) {
    form.options.splice(index, 1)
    
    // 如果删除的选项是正确答案，需要重置答案
    const optionLabel = getOptionLabel(index)
    if (form.type === 'single' && form.answer === optionLabel) {
      form.answer = ''
    } else if (form.type === 'multiple' && Array.isArray(form.answer)) {
      const answerIndex = form.answer.indexOf(optionLabel)
      if (answerIndex > -1) {
        form.answer.splice(answerIndex, 1)
      }
    }
  }
}

// 获取正确答案索引（后端期望的格式）
const getCorrectAnswerIndex = () => {
  if (form.type === 'single') {
    // 单选题：返回选中选项的索引
    const labels = ['A', 'B', 'C', 'D', 'E', 'F']
    const index = labels.indexOf(form.answer as string)
    return index >= 0 ? index : 0
  } else if (form.type === 'multiple') {
    // 多选题：返回第一个选中选项的索引（简化处理）
    const answers = form.answer as string[]
    if (answers.length > 0) {
      const labels = ['A', 'B', 'C', 'D', 'E', 'F']
      const index = labels.indexOf(answers[0])
      return index >= 0 ? index : 0
    }
    return 0
  } else if (form.type === 'judge') {
    // 判断题：true=0, false=1
    return form.answer === 'true' ? 0 : 1
  } else {
    // 填空题：返回0
    return 0
  }
}

// 预览题目
const handlePreview = () => {
  if (!form.content.trim()) {
    ElMessage.warning('请先输入题目内容')
    return
  }
  previewVisible.value = true
}

// 保存题目
const handleSave = async () => {
  if (!formRef.value) return
  
  try {
    // 验证表单
    const valid = await formRef.value.validate()
    if (!valid) return
    
    // 验证选择题选项
    if ((form.type === 'single' || form.type === 'multiple') && 
        form.options.filter(opt => opt.trim()).length < 2) {
      ElMessage.error('选择题至少需要2个选项')
      return
    }
    
    saving.value = true
    
    // 准备提交数据
    const submitData = {
      title: form.content.trim(), // 后端期望title字段
      content: form.content.trim(),
      options: (form.type === 'single' || form.type === 'multiple') 
        ? form.options.filter(opt => opt.trim()) 
        : ['true', 'false'], // 判断题和填空题也需要提供options
      correct_answer: getCorrectAnswerIndex(), // 后端期望correct_answer字段
      explanation: form.explanation.trim() || '',
      difficulty: form.difficulty,
      category_id: parseInt(form.categoryId.toString())
    }
    
    await api.post('/api/v1/admin/questions', submitData)
    ElMessage.success('题目创建成功')
    router.push('/questions')
  } catch (error) {
    console.error('保存题目失败:', error)
    ElMessage.error('保存失败，请稍后重试')
  } finally {
    saving.value = false
  }
}

// 组件挂载时获取数据
onMounted(() => {
  fetchCategories()
})
</script>

<style scoped>
.question-create {
  max-width: 1000px;
  margin: 0 auto;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 24px;
}

.header-left {
  display: flex;
  align-items: flex-start;
  gap: 16px;
}

.header-title {
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

.form-card {
  margin-bottom: 20px;
}

.form-card :deep(.el-card__body) {
  padding: 32px;
}

.form-section {
  margin-bottom: 32px;
  padding-bottom: 24px;
  border-bottom: 1px solid #e5e7eb;
}

.form-section:last-child {
  border-bottom: none;
  margin-bottom: 0;
}

.section-title {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
  margin: 0 0 20px 0;
  padding-bottom: 8px;
  border-bottom: 2px solid #e5e7eb;
}

.options-container {
  margin-bottom: 20px;
}

.option-item {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
}

.option-input {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 8px;
}

.option-label {
  font-weight: 600;
  color: #374151;
  min-width: 24px;
}

.option-actions {
  width: 40px;
  display: flex;
  justify-content: center;
}

.form-tip {
  font-size: 12px;
  color: #6b7280;
  margin-top: 4px;
  line-height: 1.4;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    gap: 16px;
  }
  
  .header-left {
    flex-direction: column;
    gap: 12px;
    width: 100%;
  }
  
  .header-right {
    width: 100%;
    justify-content: flex-start;
  }
  
  .form-card :deep(.el-card__body) {
    padding: 20px;
  }
  
  .option-item {
    flex-direction: column;
    align-items: stretch;
    gap: 8px;
  }
  
  .option-actions {
    width: auto;
    justify-content: flex-end;
  }
}
</style>