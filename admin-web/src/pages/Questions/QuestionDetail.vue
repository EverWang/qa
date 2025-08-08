<template>
  <div class="question-detail">
    <div class="detail-header">
      <div class="question-meta">
        <el-tag :type="getTypeTagType(question.type)" size="large">
          {{ getTypeLabel(question.type) }}
        </el-tag>
        <el-tag :type="getDifficultyTagType(question.difficulty)" size="large">
          {{ getDifficultyLabel(question.difficulty) }}
        </el-tag>
        <el-tag v-if="question.category" type="info" size="large">
          {{ question.category.name }}
        </el-tag>
      </div>
    </div>
    
    <div class="detail-content">
      <!-- 题目内容 -->
      <div class="content-section">
        <h3 class="section-title">题目内容</h3>
        <div class="question-content">
          {{ question.content }}
        </div>
      </div>
      
      <!-- 选项（单选题和多选题） -->
      <div v-if="question.options && question.options.length > 0" class="content-section">
        <h3 class="section-title">选项</h3>
        <div class="options-list">
          <div
            v-for="(option, index) in question.options"
            :key="index"
            class="option-item"
            :class="{ 'correct-option': isCorrectOption(option, index) }"
          >
            <span class="option-label">{{ getOptionLabel(index) }}.</span>
            <span class="option-text">{{ option }}</span>
            <el-icon v-if="isCorrectOption(option, index)" class="correct-icon">
              <Check />
            </el-icon>
          </div>
        </div>
      </div>
      
      <!-- 正确答案 -->
      <div class="content-section">
        <h3 class="section-title">正确答案</h3>
        <div class="answer-content">
          <template v-if="question.type === 'judge'">
            <el-tag :type="question.answer === 'true' ? 'success' : 'danger'" size="large">
              {{ question.answer === 'true' ? '正确' : '错误' }}
            </el-tag>
          </template>
          <template v-else-if="question.type === 'fill'">
            <div class="fill-answer">
              {{ Array.isArray(question.answer) ? question.answer.join('、') : question.answer }}
            </div>
          </template>
          <template v-else>
            <div class="choice-answer">
              <template v-if="Array.isArray(question.answer)">
                <el-tag
                  v-for="ans in question.answer"
                  :key="ans"
                  type="success"
                  size="large"
                  class="answer-tag"
                >
                  {{ ans }}
                </el-tag>
              </template>
              <template v-else>
                <el-tag type="success" size="large">
                  {{ question.answer }}
                </el-tag>
              </template>
            </div>
          </template>
        </div>
      </div>
      
      <!-- 解析 -->
      <div v-if="question.explanation" class="content-section">
        <h3 class="section-title">题目解析</h3>
        <div class="explanation-content">
          {{ question.explanation }}
        </div>
      </div>
      
      <!-- 题目信息 -->
      <div class="content-section">
        <h3 class="section-title">题目信息</h3>
        <div class="info-grid">
          <div class="info-item">
            <span class="info-label">题目ID：</span>
            <span class="info-value">{{ question.id }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">创建时间：</span>
            <span class="info-value">{{ formatDate(question.created_at) }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">更新时间：</span>
            <span class="info-value">{{ formatDate(question.updated_at) }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">所属分类：</span>
            <span class="info-value">{{ question.category?.name || '未分类' }}</span>
          </div>
        </div>
      </div>
    </div>
    
    <div class="detail-footer">
      <el-button @click="$emit('close')">
        关闭
      </el-button>
      <el-button type="primary" @click="handleEdit">
        <el-icon><Edit /></el-icon>
        编辑题目
      </el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Check, Edit } from '@element-plus/icons-vue'

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

// Props
const props = defineProps<{
  question: Question
}>()

// Emits
const emit = defineEmits<{
  close: []
}>()

// 判断是否为正确选项
const isCorrectOption = (option: string, index: number) => {
  if (props.question.type === 'single') {
    const answerIndex = getOptionIndex(props.question.answer as string)
    return index === answerIndex
  } else if (props.question.type === 'multiple') {
    const answers = Array.isArray(props.question.answer) 
      ? props.question.answer 
      : [props.question.answer]
    return answers.some(ans => {
      const answerIndex = getOptionIndex(ans)
      return index === answerIndex
    })
  }
  return false
}

// 获取选项索引
const getOptionIndex = (answer: string) => {
  const optionMap: Record<string, number> = {
    'A': 0, 'B': 1, 'C': 2, 'D': 3, 'E': 4, 'F': 5
  }
  return optionMap[answer.toUpperCase()] ?? -1
}

// 获取选项标签
const getOptionLabel = (index: number) => {
  const labels = ['A', 'B', 'C', 'D', 'E', 'F']
  return labels[index] || (index + 1).toString()
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

// 编辑题目
const handleEdit = () => {
  // 跳转到编辑页面
  router.push(`/questions/${props.question.id}/edit`)
  emit('close')
}
</script>

<style scoped>
.question-detail {
  max-height: 70vh;
  overflow-y: auto;
}

.detail-header {
  margin-bottom: 24px;
  padding-bottom: 16px;
  border-bottom: 1px solid #e5e7eb;
}

.question-meta {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.detail-content {
  margin-bottom: 24px;
}

.content-section {
  margin-bottom: 24px;
}

.section-title {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin: 0 0 12px 0;
  padding-bottom: 8px;
  border-bottom: 2px solid #e5e7eb;
}

.question-content {
  font-size: 16px;
  line-height: 1.6;
  color: #374151;
  padding: 16px;
  background-color: #f9fafb;
  border-radius: 8px;
  border-left: 4px solid #3b82f6;
}

.options-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.option-item {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  background-color: #f9fafb;
  border-radius: 8px;
  border: 1px solid #e5e7eb;
  transition: all 0.2s;
}

.option-item.correct-option {
  background-color: #f0f9ff;
  border-color: #3b82f6;
  color: #1e40af;
}

.option-label {
  font-weight: 600;
  margin-right: 12px;
  min-width: 24px;
}

.option-text {
  flex: 1;
  line-height: 1.5;
}

.correct-icon {
  color: #10b981;
  font-size: 18px;
  margin-left: 8px;
}

.answer-content {
  padding: 16px;
  background-color: #f0f9ff;
  border-radius: 8px;
  border-left: 4px solid #10b981;
}

.fill-answer {
  font-size: 16px;
  font-weight: 600;
  color: #059669;
}

.choice-answer {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.answer-tag {
  font-size: 14px;
  font-weight: 600;
}

.explanation-content {
  font-size: 15px;
  line-height: 1.6;
  color: #374151;
  padding: 16px;
  background-color: #fffbeb;
  border-radius: 8px;
  border-left: 4px solid #f59e0b;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 16px;
}

.info-item {
  display: flex;
  align-items: center;
  padding: 12px;
  background-color: #f9fafb;
  border-radius: 6px;
}

.info-label {
  font-weight: 600;
  color: #6b7280;
  margin-right: 8px;
  min-width: 80px;
}

.info-value {
  color: #374151;
  flex: 1;
}

.detail-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding-top: 16px;
  border-top: 1px solid #e5e7eb;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .question-meta {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .info-grid {
    grid-template-columns: 1fr;
  }
  
  .detail-footer {
    flex-direction: column;
  }
  
  .detail-footer .el-button {
    width: 100%;
  }
}
</style>