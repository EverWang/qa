<template>
  <div class="question-import">
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
        <h1 class="page-title">批量导入题目</h1>
      </div>
    </div>

    <div class="page-content">
      <el-row :gutter="24">
        <el-col :span="16">
          <!-- 导入步骤 -->
          <el-card class="steps-card">
            <el-steps :active="currentStep" align-center>
              <el-step title="下载模板" description="下载Excel导入模板" />
              <el-step title="上传文件" description="上传填写好的Excel文件" />
              <el-step title="数据预览" description="预览导入数据" />
              <el-step title="确认导入" description="确认并导入题目" />
            </el-steps>
          </el-card>

          <!-- 步骤1: 下载模板 -->
          <el-card v-if="currentStep === 0" class="step-card">
            <template #header>
              <span class="card-title">步骤1: 下载导入模板</span>
            </template>
            
            <div class="template-section">
              <div class="template-info">
                <el-icon class="info-icon"><Document /></el-icon>
                <div class="info-content">
                  <h3>Excel导入模板</h3>
                  <p>请下载标准的Excel模板，按照模板格式填写题目数据</p>
                </div>
              </div>
              
              <div class="template-actions">
                <el-button 
                  type="primary" 
                  @click="downloadTemplate"
                  :loading="downloadLoading"
                >
                  <el-icon><Download /></el-icon>
                  下载模板
                </el-button>
                <el-button @click="nextStep">跳过，直接上传</el-button>
              </div>
              
              <div class="template-rules">
                <h4>填写说明：</h4>
                <ul>
                  <li>题目内容：必填，题目的具体内容</li>
                  <li>题目类型：必填，可选值：single(单选)、multiple(多选)、judge(判断)、fill(填空)</li>
                  <li>难度等级：必填，可选值：easy(简单)、medium(中等)、hard(困难)</li>
                  <li>分类ID：必填，题目所属分类的ID</li>
                  <li>选项A-F：选择题必填，其他题型可留空</li>
                  <li>正确答案：必填，单选填字母(如A)，多选用逗号分隔(如A,C)，判断填true/false，填空填具体答案</li>
                  <li>题目解析：可选，题目的解析说明</li>
                </ul>
              </div>
            </div>
          </el-card>

          <!-- 步骤2: 上传文件 -->
          <el-card v-if="currentStep === 1" class="step-card">
            <template #header>
              <span class="card-title">步骤2: 上传Excel文件</span>
            </template>
            
            <div class="upload-section">
              <el-upload
                ref="uploadRef"
                class="upload-dragger"
                drag
                :auto-upload="false"
                :on-change="handleFileChange"
                :before-upload="beforeUpload"
                accept=".xlsx,.xls"
                :limit="1"
                :file-list="fileList"
              >
                <el-icon class="el-icon--upload"><UploadFilled /></el-icon>
                <div class="el-upload__text">
                  将Excel文件拖到此处，或<em>点击上传</em>
                </div>
                <template #tip>
                  <div class="el-upload__tip">
                    支持.xlsx和.xls格式，文件大小不超过10MB
                  </div>
                </template>
              </el-upload>
              
              <div v-if="uploadFile" class="file-info">
                <div class="file-item">
                  <el-icon><Document /></el-icon>
                  <span class="file-name">{{ uploadFile.name }}</span>
                  <span class="file-size">({{ formatFileSize(uploadFile.size) }})</span>
                  <el-button 
                    type="danger" 
                    text 
                    @click="removeFile"
                  >
                    <el-icon><Delete /></el-icon>
                  </el-button>
                </div>
              </div>
              
              <div class="upload-actions">
                <el-button @click="prevStep">上一步</el-button>
                <el-button 
                  type="primary" 
                  @click="parseFile"
                  :disabled="!uploadFile"
                  :loading="parsing"
                >
                  解析文件
                </el-button>
              </div>
            </div>
          </el-card>

          <!-- 步骤3: 数据预览 -->
          <el-card v-if="currentStep === 2" class="step-card">
            <template #header>
              <div class="preview-header">
                <span class="card-title">步骤3: 数据预览</span>
                <div class="preview-stats">
                  <el-tag type="info">总计: {{ previewData.length }}题</el-tag>
                  <el-tag type="success">有效: {{ validCount }}题</el-tag>
                  <el-tag v-if="errorCount > 0" type="danger">错误: {{ errorCount }}题</el-tag>
                </div>
              </div>
            </template>
            
            <div class="preview-section">
              <div v-if="errorCount > 0" class="error-summary">
                <el-alert
                  title="数据验证失败"
                  type="error"
                  :description="`发现 ${errorCount} 条错误数据，请修正后重新上传`"
                  show-icon
                  :closable="false"
                />
              </div>
              
              <el-table 
                :data="previewData" 
                border 
                stripe
                max-height="400"
                class="preview-table"
              >
                <el-table-column type="index" label="行号" width="60" />
                <el-table-column prop="content" label="题目内容" min-width="200">
                  <template #default="{ row }">
                    <div class="content-cell" :class="{ 'error': row.errors?.content }">
                      {{ row.content || '(空)' }}
                      <el-tooltip v-if="row.errors?.content" :content="row.errors.content">
                        <el-icon class="error-icon"><Warning /></el-icon>
                      </el-tooltip>
                    </div>
                  </template>
                </el-table-column>
                <el-table-column prop="type" label="类型" width="80">
                  <template #default="{ row }">
                    <div :class="{ 'error': row.errors?.type }">
                      <el-tag :type="getTypeTagType(row.type)" size="small">
                        {{ getTypeLabel(row.type) }}
                      </el-tag>
                      <el-tooltip v-if="row.errors?.type" :content="row.errors.type">
                        <el-icon class="error-icon"><Warning /></el-icon>
                      </el-tooltip>
                    </div>
                  </template>
                </el-table-column>
                <el-table-column prop="difficulty" label="难度" width="80">
                  <template #default="{ row }">
                    <div :class="{ 'error': row.errors?.difficulty }">
                      <el-tag :type="getDifficultyTagType(row.difficulty)" size="small">
                        {{ getDifficultyLabel(row.difficulty) }}
                      </el-tag>
                      <el-tooltip v-if="row.errors?.difficulty" :content="row.errors.difficulty">
                        <el-icon class="error-icon"><Warning /></el-icon>
                      </el-tooltip>
                    </div>
                  </template>
                </el-table-column>
                <el-table-column prop="categoryId" label="分类" width="100">
                  <template #default="{ row }">
                    <div :class="{ 'error': row.errors?.categoryId }">
                      {{ getCategoryName(row.categoryId) }}
                      <el-tooltip v-if="row.errors?.categoryId" :content="row.errors.categoryId">
                        <el-icon class="error-icon"><Warning /></el-icon>
                      </el-tooltip>
                    </div>
                  </template>
                </el-table-column>
                <el-table-column prop="answer" label="答案" width="100">
                  <template #default="{ row }">
                    <div :class="{ 'error': row.errors?.answer }">
                      {{ formatAnswer(row.answer, row.type) }}
                      <el-tooltip v-if="row.errors?.answer" :content="row.errors.answer">
                        <el-icon class="error-icon"><Warning /></el-icon>
                      </el-tooltip>
                    </div>
                  </template>
                </el-table-column>
                <el-table-column label="状态" width="80">
                  <template #default="{ row }">
                    <el-tag :type="row.valid ? 'success' : 'danger'" size="small">
                      {{ row.valid ? '有效' : '错误' }}
                    </el-tag>
                  </template>
                </el-table-column>
              </el-table>
              
              <div class="preview-actions">
                <el-button @click="prevStep">上一步</el-button>
                <el-button 
                  type="primary" 
                  @click="nextStep"
                  :disabled="errorCount > 0"
                >
                  确认数据
                </el-button>
              </div>
            </div>
          </el-card>

          <!-- 步骤4: 确认导入 -->
          <el-card v-if="currentStep === 3" class="step-card">
            <template #header>
              <span class="card-title">步骤4: 确认导入</span>
            </template>
            
            <div class="confirm-section">
              <div class="import-summary">
                <el-descriptions title="导入摘要" :column="2" border>
                  <el-descriptions-item label="文件名称">
                    {{ uploadFile?.name }}
                  </el-descriptions-item>
                  <el-descriptions-item label="文件大小">
                    {{ formatFileSize(uploadFile?.size || 0) }}
                  </el-descriptions-item>
                  <el-descriptions-item label="题目总数">
                    {{ validCount }}
                  </el-descriptions-item>
                  <el-descriptions-item label="导入时间">
                    {{ new Date().toLocaleString() }}
                  </el-descriptions-item>
                </el-descriptions>
              </div>
              
              <div class="import-options">
                <h4>导入选项：</h4>
                <el-checkbox v-model="importOptions.skipDuplicates">
                  跳过重复题目（基于题目内容判断）
                </el-checkbox>
                <el-checkbox v-model="importOptions.updateExisting">
                  更新已存在的题目
                </el-checkbox>
              </div>
              
              <div class="confirm-actions">
                <el-button @click="prevStep">上一步</el-button>
                <el-button 
                  type="primary" 
                  @click="confirmImport"
                  :loading="importing"
                >
                  确认导入
                </el-button>
              </div>
            </div>
          </el-card>
        </el-col>
        
        <el-col :span="8">
          <el-card class="help-card">
            <template #header>
              <span class="card-title">导入帮助</span>
            </template>
            
            <div class="help-content">
              <div class="help-section">
                <h4>支持的题目类型：</h4>
                <ul>
                  <li><strong>单选题 (single)</strong>：有2-6个选项，只有一个正确答案</li>
                  <li><strong>多选题 (multiple)</strong>：有2-6个选项，可有多个正确答案</li>
                  <li><strong>判断题 (judge)</strong>：只有正确/错误两个选项</li>
                  <li><strong>填空题 (fill)</strong>：需要填写具体答案</li>
                </ul>
              </div>
              
              <div class="help-section">
                <h4>答案格式说明：</h4>
                <ul>
                  <li><strong>单选题</strong>：填写选项字母，如 A</li>
                  <li><strong>多选题</strong>：用逗号分隔，如 A,C,D</li>
                  <li><strong>判断题</strong>：填写 true 或 false</li>
                  <li><strong>填空题</strong>：填写具体答案，多个答案用逗号分隔</li>
                </ul>
              </div>
              
              <div class="help-section">
                <h4>常见问题：</h4>
                <ul>
                  <li>确保Excel文件格式正确</li>
                  <li>题目内容不能为空</li>
                  <li>分类ID必须存在</li>
                  <li>选择题必须有选项</li>
                  <li>答案格式要正确</li>
                </ul>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox, type UploadFile, type UploadFiles } from 'element-plus'
import { 
  ArrowLeft, 
  Document, 
  Download, 
  UploadFilled, 
  Delete, 
  Warning 
} from '@element-plus/icons-vue'
import { api } from '@/lib/axios'
import * as XLSX from 'xlsx'

// 路由
const router = useRouter()

// 响应式数据
const currentStep = ref(0)
const uploadRef = ref()
const uploadFile = ref<File | null>(null)
const fileList = ref<UploadFiles>([])
const downloadLoading = ref(false)
const parsing = ref(false)
const importing = ref(false)
const categories = ref<Array<{ id: number; name: string }>>([])
const previewData = ref<any[]>([])

// 导入选项
const importOptions = reactive({
  skipDuplicates: true,
  updateExisting: false
})

// 计算属性
const validCount = computed(() => {
  return previewData.value.filter(item => item.valid).length
})

const errorCount = computed(() => {
  return previewData.value.filter(item => !item.valid).length
})

// 方法
const handleBack = () => {
  router.push('/questions')
}

const nextStep = () => {
  if (currentStep.value < 3) {
    currentStep.value++
  }
}

const prevStep = () => {
  if (currentStep.value > 0) {
    currentStep.value--
  }
}

const downloadTemplate = async () => {
  try {
    downloadLoading.value = true
    
    // 创建模板数据
    const templateData = [
      {
        '题目内容': '以下哪个是JavaScript的数据类型？',
        '题目类型': 'single',
        '难度等级': 'easy',
        '分类ID': 1,
        '选项A': 'string',
        '选项B': 'number',
        '选项C': 'boolean',
        '选项D': 'object',
        '选项E': '',
        '选项F': '',
        '正确答案': 'A,B,C,D',
        '题目解析': 'JavaScript有多种基本数据类型'
      },
      {
        '题目内容': 'JavaScript是一种编译型语言',
        '题目类型': 'judge',
        '难度等级': 'medium',
        '分类ID': 1,
        '选项A': '',
        '选项B': '',
        '选项C': '',
        '选项D': '',
        '选项E': '',
        '选项F': '',
        '正确答案': 'false',
        '题目解析': 'JavaScript是解释型语言，不是编译型语言'
      }
    ]
    
    // 创建工作簿
    const wb = XLSX.utils.book_new()
    const ws = XLSX.utils.json_to_sheet(templateData)
    
    // 设置列宽
    ws['!cols'] = [
      { wch: 30 }, // 题目内容
      { wch: 10 }, // 题目类型
      { wch: 10 }, // 难度等级
      { wch: 8 },  // 分类ID
      { wch: 15 }, // 选项A
      { wch: 15 }, // 选项B
      { wch: 15 }, // 选项C
      { wch: 15 }, // 选项D
      { wch: 15 }, // 选项E
      { wch: 15 }, // 选项F
      { wch: 15 }, // 正确答案
      { wch: 30 }  // 题目解析
    ]
    
    XLSX.utils.book_append_sheet(wb, ws, '题目导入模板')
    
    // 下载文件
    XLSX.writeFile(wb, '题目导入模板.xlsx')
    
    ElMessage.success('模板下载成功')
  } catch (error) {
    console.error('下载模板失败:', error)
    ElMessage.error('下载模板失败')
  } finally {
    downloadLoading.value = false
  }
}

const handleFileChange = (file: UploadFile, files: UploadFiles) => {
  uploadFile.value = file.raw || null
  fileList.value = files
}

const beforeUpload = (file: File) => {
  const isExcel = file.type === 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' || 
                  file.type === 'application/vnd.ms-excel'
  const isLt10M = file.size / 1024 / 1024 < 10
  
  if (!isExcel) {
    ElMessage.error('只能上传Excel文件！')
    return false
  }
  if (!isLt10M) {
    ElMessage.error('文件大小不能超过10MB！')
    return false
  }
  return false // 阻止自动上传
}

const removeFile = () => {
  uploadFile.value = null
  fileList.value = []
  uploadRef.value?.clearFiles()
}

const parseFile = async () => {
  if (!uploadFile.value) {
    ElMessage.error('请先选择文件')
    return
  }
  
  try {
    parsing.value = true
    
    const data = await readExcelFile(uploadFile.value)
    const validatedData = validateData(data)
    
    previewData.value = validatedData
    
    if (validatedData.length === 0) {
      ElMessage.error('文件中没有有效数据')
      return
    }
    
    nextStep()
    
  } catch (error: any) {
    console.error('解析文件失败:', error)
    ElMessage.error(error.message || '解析文件失败')
  } finally {
    parsing.value = false
  }
}

const readExcelFile = (file: File): Promise<any[]> => {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    
    reader.onload = (e) => {
      try {
        const data = new Uint8Array(e.target?.result as ArrayBuffer)
        const workbook = XLSX.read(data, { type: 'array' })
        const sheetName = workbook.SheetNames[0]
        const worksheet = workbook.Sheets[sheetName]
        const jsonData = XLSX.utils.sheet_to_json(worksheet)
        
        resolve(jsonData)
      } catch (error) {
        reject(new Error('文件格式错误，请检查Excel文件'))
      }
    }
    
    reader.onerror = () => {
      reject(new Error('文件读取失败'))
    }
    
    reader.readAsArrayBuffer(file)
  })
}

const validateData = (data: any[]) => {
  return data.map((row, index) => {
    const errors: Record<string, string> = {}
    let valid = true
    
    // 验证题目内容
    if (!row['题目内容'] || !row['题目内容'].toString().trim()) {
      errors.content = '题目内容不能为空'
      valid = false
    }
    
    // 验证题目类型
    const validTypes = ['single', 'multiple', 'judge', 'fill']
    if (!row['题目类型'] || !validTypes.includes(row['题目类型'])) {
      errors.type = '题目类型必须是: single, multiple, judge, fill 之一'
      valid = false
    }
    
    // 验证难度等级
    const validDifficulties = ['easy', 'medium', 'hard']
    if (!row['难度等级'] || !validDifficulties.includes(row['难度等级'])) {
      errors.difficulty = '难度等级必须是: easy, medium, hard 之一'
      valid = false
    }
    
    // 验证分类ID
    const categoryId = parseInt(row['分类ID'])
    if (!categoryId || !categories.value.find(cat => cat.id === categoryId)) {
      errors.categoryId = '分类ID不存在'
      valid = false
    }
    
    // 验证选项（选择题）
    if (row['题目类型'] === 'single' || row['题目类型'] === 'multiple') {
      const options = [row['选项A'], row['选项B'], row['选项C'], row['选项D'], row['选项E'], row['选项F']]
        .filter(opt => opt && opt.toString().trim())
      
      if (options.length < 2) {
        errors.options = '选择题至少需要2个选项'
        valid = false
      }
    }
    
    // 验证答案
    if (!row['正确答案'] || !row['正确答案'].toString().trim()) {
      errors.answer = '正确答案不能为空'
      valid = false
    }
    
    return {
      rowIndex: index + 2, // Excel行号（从2开始，因为第1行是标题）
      content: row['题目内容'],
      type: row['题目类型'],
      difficulty: row['难度等级'],
      categoryId: categoryId,
      options: [row['选项A'], row['选项B'], row['选项C'], row['选项D'], row['选项E'], row['选项F']]
        .filter(opt => opt && opt.toString().trim()),
      answer: row['正确答案'],
      explanation: row['题目解析'] || '',
      errors,
      valid
    }
  })
}

const confirmImport = async () => {
  try {
    importing.value = true
    
    const validData = previewData.value.filter(item => item.valid)
    
    const response = await api.post('/api/v1/admin/questions/import', {
      questions: validData.map(item => ({
        content: item.content,
        type: item.type,
        difficulty: item.difficulty,
        categoryId: item.categoryId,
        options: item.options.length > 0 ? item.options : undefined,
        answer: item.answer,
        explanation: item.explanation || undefined
      })),
      options: importOptions
    })
    
    const result = response.data
    
    ElMessage.success(`导入成功！共导入 ${result.imported} 题，跳过 ${result.skipped} 题`)
    
    router.push('/questions')
    
  } catch (error: any) {
    console.error('导入失败:', error)
    ElMessage.error(error.response?.data?.message || '导入失败')
  } finally {
    importing.value = false
  }
}

const formatFileSize = (size: number) => {
  if (size < 1024) {
    return size + ' B'
  } else if (size < 1024 * 1024) {
    return (size / 1024).toFixed(1) + ' KB'
  } else {
    return (size / 1024 / 1024).toFixed(1) + ' MB'
  }
}

const getCategoryName = (categoryId: number) => {
  const category = categories.value.find(cat => cat.id === categoryId)
  return category ? category.name : `分类${categoryId}`
}

const formatAnswer = (answer: string, type: string) => {
  if (!answer) return '(空)'
  
  if (type === 'judge') {
    return answer === 'true' ? '正确' : answer === 'false' ? '错误' : answer
  }
  
  return answer.toString()
}

const getTypeLabel = (type: string) => {
  const typeMap: Record<string, string> = {
    single: '单选',
    multiple: '多选',
    judge: '判断',
    fill: '填空'
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

const fetchCategories = async () => {
  try {
    const response = await api.get('/api/v1/admin/categories', { params: { tree: 'true' } })
    categories.value = response.data.data || []
  } catch (error) {
    console.error('获取分类列表失败:', error)
    ElMessage.error('获取分类列表失败')
  }
}

// 生命周期
onMounted(() => {
  fetchCategories()
})
</script>

<style scoped>
.question-import {
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

.page-content {
  max-width: 1400px;
}

.steps-card {
  margin-bottom: 24px;
}

.step-card {
  min-height: 400px;
}

.card-title {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
}

.template-section {
  padding: 20px 0;
}

.template-info {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 24px;
  padding: 20px;
  background-color: #f8fafc;
  border-radius: 8px;
  border-left: 4px solid #3b82f6;
}

.info-icon {
  font-size: 32px;
  color: #3b82f6;
}

.info-content h3 {
  margin: 0 0 8px 0;
  font-size: 18px;
  color: #1f2937;
}

.info-content p {
  margin: 0;
  color: #6b7280;
}

.template-actions {
  display: flex;
  gap: 12px;
  margin-bottom: 24px;
}

.template-rules {
  background-color: #fffbeb;
  padding: 16px;
  border-radius: 8px;
  border-left: 4px solid #f59e0b;
}

.template-rules h4 {
  margin: 0 0 12px 0;
  color: #92400e;
}

.template-rules ul {
  margin: 0;
  padding-left: 20px;
  color: #78350f;
}

.template-rules li {
  margin-bottom: 8px;
}

.upload-section {
  padding: 20px 0;
}

.upload-dragger {
  margin-bottom: 24px;
}

.file-info {
  margin-bottom: 24px;
}

.file-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  background-color: #f8fafc;
  border-radius: 8px;
  border: 1px solid #e5e7eb;
}

.file-name {
  flex: 1;
  font-weight: 500;
}

.file-size {
  color: #6b7280;
  font-size: 14px;
}

.upload-actions {
  display: flex;
  gap: 12px;
}

.preview-section {
  padding: 20px 0;
}

.preview-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.preview-stats {
  display: flex;
  gap: 8px;
}

.error-summary {
  margin-bottom: 16px;
}

.preview-table {
  margin-bottom: 24px;
}

.content-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.content-cell.error {
  color: #ef4444;
}

.error-icon {
  color: #ef4444;
  font-size: 16px;
}

.preview-actions {
  display: flex;
  gap: 12px;
}

.confirm-section {
  padding: 20px 0;
}

.import-summary {
  margin-bottom: 24px;
}

.import-options {
  margin-bottom: 24px;
  padding: 16px;
  background-color: #f8fafc;
  border-radius: 8px;
}

.import-options h4 {
  margin: 0 0 12px 0;
  color: #1f2937;
}

.import-options .el-checkbox {
  display: block;
  margin-bottom: 8px;
}

.confirm-actions {
  display: flex;
  gap: 12px;
}

.help-card {
  position: sticky;
  top: 24px;
}

.help-content {
  font-size: 14px;
}

.help-section {
  margin-bottom: 20px;
}

.help-section:last-child {
  margin-bottom: 0;
}

.help-section h4 {
  margin: 0 0 8px 0;
  color: #1f2937;
  font-size: 14px;
}

.help-section ul {
  margin: 0;
  padding-left: 16px;
  color: #6b7280;
}

.help-section li {
  margin-bottom: 4px;
  line-height: 1.5;
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .page-content :deep(.el-col) {
    width: 100%;
  }
  
  .help-card {
    position: static;
    margin-top: 24px;
  }
}

@media (max-width: 768px) {
  .question-import {
    padding: 16px;
  }
  
  .template-info {
    flex-direction: column;
    text-align: center;
  }
  
  .template-actions {
    flex-direction: column;
  }
  
  .upload-actions,
  .preview-actions,
  .confirm-actions {
    flex-direction: column;
  }
  
  .preview-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }
  
  .preview-stats {
    flex-wrap: wrap;
  }
}
</style>