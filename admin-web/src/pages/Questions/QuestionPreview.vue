<template>
  <div class="question-preview">
    <div class="preview-header">
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
    
    <div class="preview-content">
      <!-- 题目内容 -->
      <div class="question-content">
        <h3 class="content-title">题目内容</h3>
        <div class="content-text">
          {{ question.content || '(题目内容为空)' }}
        </div>
      </div>
      
      <!-- 选项（单选题和多选题） -->
      <div v-if="question.options && question.options.length > 0" class="question-options">
        <h3 class="content-title">选项</h3>
        <div class="options-list">
          <div
            v-for="(option, index) in question.options"
            :key="index"
            class="option-item"
            :class="{ 'correct-option': isCorrectOption(option, index) }"
          >
            <div class="option-content">
              <span class="option-label">{{ getOptionLabel(index) }}.</span>
              <span class="option-text">{{ option || '(选项内容为空)' }}</span>
            </div>
            <el-icon v-if="isCorrectOption(option, index)" class="correct-icon">
              <Check />
            </el-icon>
          </div>
        </div>
      </div>
      
      <!-- 正确答案 -->
      <div class="question-answer">
        <h3 class="content-title">正确答案</h3>
        <div class="answer-content">
          <template v-if="question.type === 'judge'">
            <el-tag 
              :type="question.answer === 'true' ? 'success' : 'danger'" 
              size="large"
            >
              {{ question.answer === 'true' ? '正确' : question.answer === 'false' ? '错误' : '(未设置答案)' }}
            </el-tag>
          </template>
          <template v-else-if="question.type === 'fill'">
            <div class="fill-answer">
              {{ getFormattedFillAnswer() }}
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
                <span v-if="question.answer.length === 0" class="empty-answer">
                  (未设置答案)
                </span>
              </template>
              <template v-else>
                <el-tag v-if="question.answer" type="success" size="large">
                  {{ question.answer }}
                </el-tag>
                <span v-else class="empty-answer">(未设置答案)</span>
              </template>
            </div>
          </template>
        </div>
      </div>
      
      <!-- 解析 -->
      <div v-if="question.explanation" class="question-explanation">
        <h3 class="content-title">题目解析</h3>
        <div class="explanation-content">
          {{ question.explanation }}
        </div>
      </div>
      
      <!-- 空状态提示 -->
      <div v-if="!question.content" class="empty-state">
        <el-empty description="请先填写题目内容" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Check } from '@element-plus/icons-vue'

// 题目接口
interface Question {
  id: number
  content: string
  type: 'single' | 'multiple' | 'judge' | 'fill'
  difficulty: 'easy' | 'medium' | 'hard'
  options?: string[]
  answer: string | string[]
  explanation?: string
  category_id: number | string
  category?: {
    id: number | string
    name: string
  }
  created_at: string
  updated_at: string
}

// Props
const props = defineProps<{
  question: Question
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

// 格式化填空题答案
const getFormattedFillAnswer = () => {
  if (!props.question.answer) {
    return '(未设置答案)'
  }
  
  if (Array.isArray(props.question.answer)) {
    return props.question.answer.length > 0 
      ? props.question.answer.join('、') 
      : '(未设置答案)'
  }
  
  return props.question.answer.toString().trim() || '(未设置答案)'
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
</script>

<style scoped>
.question-preview {
  max-height: 70vh;
  overflow-y: auto;
}

.preview-header {
  margin-bottom: 24px;
  padding-bottom: 16px;
  border-bottom: 1px solid #e5e7eb;
}

.question-meta {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.preview-content {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.content-title {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin: 0 0 12px 0;
  padding-bottom: 8px;
  border-bottom: 2px solid #e5e7eb;
}

.question-content .content-text {
  font-size: 16px;
  line-height: 1.6;
  color: #374151;
  padding: 16px;
  background-color: #f9fafb;
  border-radius: 8px;
  border-left: 4px solid #3b82f6;
  min-height: 60px;
  display: flex;
  align-items: center;
}

.options-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.option-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
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

.option-content {
  display: flex;
  align-items: center;
  flex: 1;
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
  min-height: 60px;
  display: flex;
  align-items: center;
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
  align-items: center;
}

.answer-tag {
  font-size: 14px;
  font-weight: 600;
}

.empty-answer {
  color: #9ca3af;
  font-style: italic;
}

.explanation-content {
  font-size: 15px;
  line-height: 1.6;
  color: #374151;
  padding: 16px;
  background-color: #fffbeb;
  border-radius: 8px;
  border-left: 4px solid #f59e0b;
  min-height: 60px;
  display: flex;
  align-items: center;
}

.empty-state {
  padding: 40px 0;
  text-align: center;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .question-meta {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .option-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }
  
  .option-content {
    width: 100%;
  }
  
  .correct-icon {
    align-self: flex-end;
    margin-left: 0;
  }
  
  .choice-answer {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>