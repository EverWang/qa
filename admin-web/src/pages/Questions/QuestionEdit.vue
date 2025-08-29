<template>
  <div class="question-edit">
    <div class="page-header">
      <div class="header-left">
        <el-button 
          type="text" 
          @click="handleBack"
          class="back-button"
        >
          <el-icon><ArrowLeft /></el-icon>
          返回题目列表
        </el-button>
        <h1 class="page-title">编辑题目</h1>
      </div>
      <div class="header-actions">
        <el-button @click="handlePreview">预览题目</el-button>
        <el-button 
          type="primary" 
          @click="handleSave"
          :loading="saving"
        >
          保存题目
        </el-button>
      </div>
    </div>

    <div class="page-content">
      <el-row :gutter="24">
        <el-col :span="16">
          <el-card class="form-card">
            <template #header>
              <span class="card-title">题目信息</span>
            </template>
            
            <el-form
              ref="formRef"
              :model="form"
              :rules="rules"
              label-width="100px"
              label-position="top"
            >
              <!-- 基本信息 -->
              <el-row :gutter="16">
                <el-col :span="8">
                  <el-form-item label="题目类型" prop="type">
                    <el-select 
                      v-model="form.type" 
                      placeholder="请选择题目类型"
                      @change="handleTypeChange"
                      style="width: 100%"
                    >
                      <el-option label="单选题" value="single" />
                      <el-option label="多选题" value="multiple" />
                      <el-option label="判断题" value="judge" />
                      <el-option label="填空题" value="fill" />
                    </el-select>
                  </el-form-item>
                </el-col>
                <el-col :span="8">
                  <el-form-item label="题目分类" prop="categoryId">
                    <el-select 
                      v-model="form.categoryId" 
                      placeholder="请选择分类"
                      style="width: 100%"
                      filterable
                    >
                      <el-option
                        v-for="category in categories"
                        :key="category.id"
                        :label="category.name"
                        :value="category.id"
                      />
                    </el-select>
                  </el-form-item>
                </el-col>
                <el-col :span="8">
                  <el-form-item label="难度等级" prop="difficulty">
                    <el-select 
                      v-model="form.difficulty" 
                      placeholder="请选择难度"
                      style="width: 100%"
                    >
                      <el-option label="简单" value="easy" />
                      <el-option label="中等" value="medium" />
                      <el-option label="困难" value="hard" />
                    </el-select>
                  </el-form-item>
                </el-col>
              </el-row>

              <!-- 题目标题 -->
              <el-form-item label="题目标题" prop="title">
                <el-input
                  v-model="form.title"
                  placeholder="请输入题目标题"
                  maxlength="200"
                  show-word-limit
                />
              </el-form-item>

              <!-- 题目内容 -->
              <el-form-item label="题目内容" prop="content">
                <el-input
                  v-model="form.content"
                  type="textarea"
                  :rows="4"
                  placeholder="请输入题目内容"
                  maxlength="1000"
                  show-word-limit
                />
              </el-form-item>

              <!-- 选项设置（单选题和多选题） -->
              <template v-if="form.type === 'single' || form.type === 'multiple'">
                <el-form-item label="选项设置" prop="options">
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
                      <el-button
                        v-if="form.options.length > 2"
                        type="danger"
                        text
                        @click="removeOption(index)"
                        class="remove-btn"
                      >
                        <el-icon><Delete /></el-icon>
                      </el-button>
                    </div>
                    <el-button
                      v-if="form.options.length < 6"
                      type="primary"
                      text
                      @click="addOption"
                      class="add-option-btn"
                    >
                      <el-icon><Plus /></el-icon>
                      添加选项
                    </el-button>
                  </div>
                </el-form-item>
              </template>

              <!-- 答案设置 -->
              <el-form-item label="正确答案" prop="answer">
                <!-- 单选题答案 -->
                <template v-if="form.type === 'single'">
                  <el-radio-group v-model="form.answer">
                    <el-radio
                      v-for="(option, index) in form.options"
                      :key="index"
                      :label="getOptionLabel(index)"
                      :disabled="!option.trim()"
                    >
                      {{ getOptionLabel(index) }}. {{ option || '(选项内容为空)' }}
                    </el-radio>
                  </el-radio-group>
                </template>
                
                <!-- 多选题答案 -->
                <template v-else-if="form.type === 'multiple'">
                  <el-checkbox-group v-model="form.answer">
                    <el-checkbox
                      v-for="(option, index) in form.options"
                      :key="index"
                      :label="getOptionLabel(index)"
                      :disabled="!option.trim()"
                    >
                      {{ getOptionLabel(index) }}. {{ option || '(选项内容为空)' }}
                    </el-checkbox>
                  </el-checkbox-group>
                </template>
                
                <!-- 判断题答案 -->
                <template v-else-if="form.type === 'judge'">
                  <el-radio-group v-model="form.answer">
                    <el-radio label="true">正确</el-radio>
                    <el-radio label="false">错误</el-radio>
                  </el-radio-group>
                </template>
                
                <!-- 填空题答案 -->
                <template v-else-if="form.type === 'fill'">
                  <div class="fill-answers">
                    <div 
                      v-for="(answer, index) in form.answer" 
                      :key="index"
                      class="fill-answer-item"
                    >
                      <el-input
                        v-model="form.answer[index]"
                        placeholder="请输入答案"
                        maxlength="100"
                      />
                      <el-button
                        v-if="form.answer.length > 1"
                        type="danger"
                        text
                        @click="removeFillAnswer(index)"
                      >
                        <el-icon><Delete /></el-icon>
                      </el-button>
                    </div>
                    <el-button
                      type="primary"
                      text
                      @click="addFillAnswer"
                      class="add-answer-btn"
                    >
                      <el-icon><Plus /></el-icon>
                      添加答案
                    </el-button>
                    <div class="fill-tip">
                      <el-icon><InfoFilled /></el-icon>
                      填空题支持多个正确答案，用户答对其中任意一个即可
                    </div>
                  </div>
                </template>
              </el-form-item>

              <!-- 题目解析 -->
              <el-form-item label="题目解析">
                <el-input
                  v-model="form.explanation"
                  type="textarea"
                  :rows="3"
                  placeholder="请输入题目解析（可选）"
                  maxlength="500"
                  show-word-limit
                />
              </el-form-item>
            </el-form>
          </el-card>
        </el-col>
        
        <el-col :span="8">
          <el-card class="preview-card">
            <template #header>
              <span class="card-title">题目预览</span>
            </template>
            <QuestionPreview :question="previewQuestion" />
          </el-card>
        </el-col>
      </el-row>
    </div>

    <!-- 预览对话框 -->
    <el-dialog
      v-model="previewVisible"
      title="题目预览"
      width="800px"
      :before-close="handleClosePreview"
    >
      <QuestionPreview :question="previewQuestion" />
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, ElMessageBox, type FormInstance } from 'element-plus'
import { ArrowLeft, Delete, Plus, InfoFilled } from '@element-plus/icons-vue'
import QuestionPreview from './QuestionPreview.vue'
import { api } from '@/lib/axios'

// 路由
const router = useRouter()
const route = useRoute()

// 响应式数据
const formRef = ref<FormInstance>()
const loading = ref(false)
const saving = ref(false)
const previewVisible = ref(false)
const categories = ref<Array<{ id: number; name: string }>>([])

// 表单数据
const form = reactive({
  id: 0,
  title: '',
  content: '',
  type: 'single' as 'single' | 'multiple' | 'judge' | 'fill',
  difficulty: 'medium' as 'easy' | 'medium' | 'hard',
  options: ['', ''],
  answer: '' as string | string[],
  explanation: '',
  categoryId: ''
})

// 表单验证规则
const rules = {
  title: [
    { required: true, message: '请输入题目标题', trigger: 'blur' },
    { min: 3, message: '题目标题至少3个字符', trigger: 'blur' }
  ],
  content: [
    { required: true, message: '请输入题目内容', trigger: 'blur' },
    { min: 5, message: '题目内容至少5个字符', trigger: 'blur' }
  ],
  type: [
    { required: true, message: '请选择题目类型', trigger: 'change' }
  ],
  difficulty: [
    { required: true, message: '请选择难度等级', trigger: 'change' }
  ],
  categoryId: [
    { required: true, message: '请选择题目分类', trigger: 'change' }
  ],
  options: [
    {
      validator: (rule: any, value: string[], callback: Function) => {
        if (form.type === 'single' || form.type === 'multiple') {
          const validOptions = value.filter(option => option.trim())
          if (validOptions.length < 2) {
            callback(new Error('至少需要2个有效选项'))
          } else {
            callback()
          }
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ],
  answer: [
    {
      validator: (rule: any, value: any, callback: Function) => {
        if (form.type === 'single') {
          if (!value) {
            callback(new Error('请选择正确答案'))
          } else {
            callback()
          }
        } else if (form.type === 'multiple') {
          if (!Array.isArray(value) || value.length === 0) {
            callback(new Error('请至少选择一个正确答案'))
          } else {
            callback()
          }
        } else if (form.type === 'judge') {
          if (!value) {
            callback(new Error('请选择正确答案'))
          } else {
            callback()
          }
        } else if (form.type === 'fill') {
          if (!Array.isArray(value) || value.length === 0 || !value.some(ans => ans.trim())) {
            callback(new Error('请至少输入一个答案'))
          } else {
            callback()
          }
        }
        callback()
      },
      trigger: 'change'
    }
  ]
}

// 计算属性
const previewQuestion = computed(() => ({
  id: form.id,
  content: form.content,
  type: form.type,
  difficulty: form.difficulty,
  options: form.options,
  answer: form.answer,
  explanation: form.explanation,
  category_id: form.categoryId,
  category: {
    id: form.categoryId,
    name: categories.value.find(c => c.id.toString() === form.categoryId.toString())?.name || ''
  },
  created_at: new Date().toISOString(),
  updated_at: new Date().toISOString()
}))

// 方法
const handleBack = () => {
  router.push('/questions')
}

const handlePreview = () => {
  previewVisible.value = true
}

const handleClosePreview = () => {
  previewVisible.value = false
}

const handleTypeChange = () => {
  // 只在用户手动切换类型时重置答案和选项
  // 如果是从后端加载数据时的类型变化，不应该重置
  if (form.id === 0) {
    // 新建题目时才重置
    if (form.type === 'single') {
      form.answer = ''
      form.options = ['', '']
    } else if (form.type === 'multiple') {
      form.answer = []
      form.options = ['', '']
    } else if (form.type === 'judge') {
      form.answer = ''
      form.options = []
    } else if (form.type === 'fill') {
      form.answer = ['']
      form.options = []
    }
  }
}

const getOptionLabel = (index: number) => {
  const labels = ['A', 'B', 'C', 'D', 'E', 'F']
  return labels[index] || (index + 1).toString()
}

const addOption = () => {
  if (form.options.length < 6) {
    form.options.push('')
  }
}

const removeOption = (index: number) => {
  if (form.options.length > 2) {
    form.options.splice(index, 1)
    // 如果删除的选项是已选答案，需要重置答案
    const optionLabel = getOptionLabel(index)
    if (form.type === 'single' && form.answer === optionLabel) {
      form.answer = ''
    } else if (form.type === 'multiple' && Array.isArray(form.answer)) {
      form.answer = form.answer.filter(ans => ans !== optionLabel)
    }
  }
}

const addFillAnswer = () => {
  if (Array.isArray(form.answer)) {
    form.answer.push('')
  }
}

const removeFillAnswer = (index: number) => {
  if (Array.isArray(form.answer) && form.answer.length > 1) {
    form.answer.splice(index, 1)
  }
}

const handleSave = async () => {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    saving.value = true
    
    console.log('保存前的表单数据:', form) // 调试日志
    
    // 处理答案字段的转换
    let correctAnswer = 0
    if (form.type === 'single') {
      // 单选题：将选项标签转换为索引
      const labels = ['A', 'B', 'C', 'D', 'E', 'F']
      const answerIndex = labels.indexOf(form.answer as string)
      if (answerIndex === -1) {
        ElMessage.error('请选择有效的答案选项')
        return
      }
      correctAnswer = answerIndex
    } else if (form.type === 'multiple') {
      // 多选题：将选项标签数组转换为位掩码
      const labels = ['A', 'B', 'C', 'D', 'E', 'F']
      const answers = form.answer as string[]
      if (!answers || answers.length === 0) {
        ElMessage.error('多选题至少需要选择一个正确答案')
        return
      }
      let mask = 0
      answers.forEach(ans => {
        const index = labels.indexOf(ans)
        if (index >= 0 && index < form.options.length) {
          mask |= (1 << index)
        }
      })
      correctAnswer = mask
      
      // 验证位掩码是否在有效范围内
      const maxValidMask = (1 << form.options.length) - 1
      if (mask > maxValidMask) {
        ElMessage.error('选择的答案超出了选项范围')
        return
      }
    } else if (form.type === 'judge') {
      // 判断题：将true/false转换为索引
      if (form.answer !== 'true' && form.answer !== 'false') {
        ElMessage.error('请选择判断题的正确答案')
        return
      }
      correctAnswer = form.answer === 'true' ? 0 : 1
    } else if (form.type === 'fill') {
      // 填空题：索引设为0，答案存储在explanation中
      const fillAnswers = form.answer as string[]
      if (!fillAnswers || fillAnswers.length === 0 || !fillAnswers.some(ans => ans.trim())) {
        ElMessage.error('填空题至少需要一个有效答案')
        return
      }
      correctAnswer = 0
    }
    
    // 确保correctAnswer是有效的整数
    if (!Number.isInteger(correctAnswer) || correctAnswer < 0) {
      ElMessage.error('答案格式错误')
      return
    }
    
    // 确保所有必需字段都有值
    if (!form.title || !form.title.trim()) {
      ElMessage.error('题目标题不能为空')
      return
    }
    if (!form.content || !form.content.trim()) {
      ElMessage.error('题目内容不能为空')
      return
    }
    
    // 验证分类ID
    if (!form.categoryId || form.categoryId === '') {
      ElMessage.error('请选择题目分类')
      return
    }
    
    const questionData = {
      title: form.title.trim(), // 题目标题
      content: form.content.trim(), // 题目内容
      type: form.type,
      difficulty: form.difficulty,
      options: (form.type === 'single' || form.type === 'multiple') ? form.options.filter(opt => opt.trim()) : [],
      correctAnswer: correctAnswer, // 后端期望驼峰命名
      explanation: form.type === 'fill' ? (form.answer as string[])[0] : form.explanation,
      categoryId: form.categoryId // 后端期望UUID字符串，不需要parseInt
    }
    
    console.log('发送到后端的数据:', questionData) // 调试日志
    console.log('form.title值:', form.title) // 调试title字段
    console.log('form对象:', form) // 调试整个form对象
    
    await api.put(`/api/v1/admin/questions/${form.id}`, questionData)
    
    ElMessage.success('题目更新成功')
    router.push('/questions')
  } catch (error: any) {
    console.error('更新题目失败:', error)
    ElMessage.error(error.response?.data?.message || '更新题目失败')
  } finally {
    saving.value = false
  }
}

const fetchCategories = async () => {
  try {
    const response = await api.get('/api/v1/admin/categories')
    categories.value = response.data.data || []
  } catch (error) {
    console.error('获取分类列表失败:', error)
    ElMessage.error('获取分类列表失败')
  }
}

const fetchQuestion = async (id: string) => {
  try {
    loading.value = true
    const response = await api.get(`/api/v1/admin/questions/${id}`)
    const question = response.data.data
    
    console.log('获取到的题目数据:', question) // 调试日志
    
    // 填充表单数据
    form.id = question.id
    form.title = question.title || ''
    form.content = question.content || ''
    form.type = question.type || 'single'
    form.difficulty = question.difficulty
    form.options = question.options || []
    form.explanation = question.explanation || ''
    form.categoryId = (question.categoryId || question.categoryId).toString()
    
    console.log('设置后的form.type:', form.type) // 调试日志
    console.log('设置后的form.options:', form.options) // 调试日志
    
    // 处理答案字段的映射
    if (form.type === 'single') {
      // 单选题：后端correctAnswer是索引，前端需要转换为选项标签
      const answerIndex = question.correctAnswer || question.correctAnswer || 0
      form.answer = getOptionLabel(answerIndex)
    } else if (form.type === 'multiple') {
      // 多选题：后端correctAnswer是位掩码，前端需要转换为选项标签数组
      const answerMask = question.correctAnswer || question.correctAnswer || 0
      const answerLabels = []
      for (let i = 0; i < form.options.length; i++) {
        if (answerMask & (1 << i)) {
          answerLabels.push(getOptionLabel(i))
        }
      }
      form.answer = answerLabels
    } else if (form.type === 'judge') {
      // 判断题：后端correctAnswer是索引（0=正确，1=错误）
      const answerIndex = question.correctAnswer || question.correctAnswer || 0
      form.answer = answerIndex === 0 ? 'true' : 'false'
    } else if (form.type === 'fill') {
      // 填空题：使用explanation字段存储答案
      form.answer = question.explanation ? [question.explanation] : ['']
    }
    
    // 确保选项数组至少有2个元素（对于选择题）
    if ((form.type === 'single' || form.type === 'multiple') && form.options.length < 2) {
      while (form.options.length < 2) {
        form.options.push('')
      }
    }
    
  } catch (error: any) {
    console.error('获取题目详情失败:', error)
    ElMessage.error('获取题目详情失败')
    router.push('/questions')
  } finally {
    loading.value = false
  }
}

// 生命周期
onMounted(async () => {
  const questionId = route.params.id as string
  if (!questionId) {
    ElMessage.error('题目ID不存在')
    router.push('/questions')
    return
  }
  
  await Promise.all([
    fetchCategories(),
    fetchQuestion(questionId)
  ])
})
</script>

<style scoped>
.question-edit {
  padding: 24px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  padding-bottom: 16px;
  border-bottom: 1px solid #e5e7eb;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.back-button {
  display: flex;
  align-items: center;
  gap: 4px;
  color: #6b7280;
  font-size: 14px;
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
  max-width: 1400px;
}

.form-card,
.preview-card {
  height: fit-content;
}

.card-title {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
}

.options-container {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.option-item {
  display: flex;
  align-items: center;
  gap: 12px;
}

.option-input {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
}

.option-label {
  font-weight: 600;
  min-width: 24px;
  color: #374151;
}

.remove-btn {
  color: #ef4444;
}

.add-option-btn {
  align-self: flex-start;
  margin-top: 8px;
}

.fill-answers {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.fill-answer-item {
  display: flex;
  align-items: center;
  gap: 12px;
}

.add-answer-btn {
  align-self: flex-start;
}

.fill-tip {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #6b7280;
  font-size: 14px;
  margin-top: 8px;
}

.preview-card {
  position: sticky;
  top: 24px;
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .page-content :deep(.el-col) {
    width: 100%;
  }
  
  .preview-card {
    position: static;
    margin-top: 24px;
  }
}

@media (max-width: 768px) {
  .question-edit {
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
  
  .option-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }
  
  .option-input {
    width: 100%;
  }
  
  .fill-answer-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }
}
</style>