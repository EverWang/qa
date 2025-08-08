<template>
  <div class="min-h-screen bg-gray-50">
    <!-- 顶部导航 -->
    <div class="bg-white shadow-sm border-b">
      <div class="max-w-4xl mx-auto px-4 py-3 flex items-center justify-between">
        <button @click="goBack" class="flex items-center text-gray-600 hover:text-gray-800">
          <ArrowLeft class="w-5 h-5 mr-1" />
          <span>返回</span>
        </button>
        <h1 class="text-lg font-medium text-gray-800">设置</h1>
        <div class="w-16"></div>
      </div>
    </div>

    <div class="max-w-4xl mx-auto p-4">
      <!-- 账户设置 -->
      <div class="bg-white rounded-lg shadow-sm mb-4">
        <div class="p-4 border-b border-gray-100">
          <h2 class="text-lg font-semibold text-gray-800">账户设置</h2>
        </div>
        <div class="divide-y divide-gray-100">
          <!-- 个人信息 -->
          <div class="p-4 flex items-center justify-between">
            <div class="flex items-center space-x-3">
              <User class="w-5 h-5 text-gray-500" />
              <div>
                <div class="font-medium text-gray-800">个人信息</div>
                <div class="text-sm text-gray-500">修改昵称、头像等</div>
              </div>
            </div>
            <ChevronRight class="w-5 h-5 text-gray-400" />
          </div>
          
          <!-- 账户安全 -->
          <div class="p-4 flex items-center justify-between">
            <div class="flex items-center space-x-3">
              <Shield class="w-5 h-5 text-gray-500" />
              <div>
                <div class="font-medium text-gray-800">账户安全</div>
                <div class="text-sm text-gray-500">密码、绑定手机等</div>
              </div>
            </div>
            <ChevronRight class="w-5 h-5 text-gray-400" />
          </div>
          
          <!-- 隐私设置 -->
          <div class="p-4 flex items-center justify-between">
            <div class="flex items-center space-x-3">
              <Lock class="w-5 h-5 text-gray-500" />
              <div>
                <div class="font-medium text-gray-800">隐私设置</div>
                <div class="text-sm text-gray-500">数据使用、隐私保护</div>
              </div>
            </div>
            <ChevronRight class="w-5 h-5 text-gray-400" />
          </div>
        </div>
      </div>

      <!-- 学习设置 -->
      <div class="bg-white rounded-lg shadow-sm mb-4">
        <div class="p-4 border-b border-gray-100">
          <h2 class="text-lg font-semibold text-gray-800">学习设置</h2>
        </div>
        <div class="divide-y divide-gray-100">
          <!-- 每日目标 -->
          <div class="p-4">
            <div class="flex items-center justify-between mb-3">
              <div class="flex items-center space-x-3">
                <Target class="w-5 h-5 text-gray-500" />
                <div>
                  <div class="font-medium text-gray-800">每日答题目标</div>
                  <div class="text-sm text-gray-500">设置每天的答题数量</div>
                </div>
              </div>
              <span class="text-blue-600 font-medium">{{ dailyGoal }} 题</span>
            </div>
            <input
              v-model="dailyGoal"
              type="range"
              min="5"
              max="100"
              step="5"
              class="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer slider"
            >
            <div class="flex justify-between text-xs text-gray-500 mt-1">
              <span>5题</span>
              <span>100题</span>
            </div>
          </div>
          
          <!-- 答题模式 -->
          <div class="p-4">
            <div class="flex items-center space-x-3 mb-3">
              <BookOpen class="w-5 h-5 text-gray-500" />
              <div>
                <div class="font-medium text-gray-800">答题模式</div>
                <div class="text-sm text-gray-500">选择默认的答题方式</div>
              </div>
            </div>
            <div class="space-y-2">
              <label class="flex items-center space-x-3">
                <input v-model="answerMode" type="radio" value="practice" class="text-blue-600" />
                <span class="text-gray-700">练习模式（可查看答案）</span>
              </label>
              <label class="flex items-center space-x-3">
                <input v-model="answerMode" type="radio" value="exam" class="text-blue-600" />
                <span class="text-gray-700">考试模式（完成后查看）</span>
              </label>
            </div>
          </div>
          
          <!-- 自动收藏错题 -->
          <div class="p-4 flex items-center justify-between">
            <div class="flex items-center space-x-3">
              <BookmarkPlus class="w-5 h-5 text-gray-500" />
              <div>
                <div class="font-medium text-gray-800">自动收藏错题</div>
                <div class="text-sm text-gray-500">答错的题目自动加入错题本</div>
              </div>
            </div>
            <label class="relative inline-flex items-center cursor-pointer">
              <input v-model="autoCollectMistakes" type="checkbox" class="sr-only peer" />
              <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"></div>
            </label>
          </div>
        </div>
      </div>

      <!-- 通知设置 -->
      <div class="bg-white rounded-lg shadow-sm mb-4">
        <div class="p-4 border-b border-gray-100">
          <h2 class="text-lg font-semibold text-gray-800">通知设置</h2>
        </div>
        <div class="divide-y divide-gray-100">
          <!-- 学习提醒 -->
          <div class="p-4 flex items-center justify-between">
            <div class="flex items-center space-x-3">
              <Bell class="w-5 h-5 text-gray-500" />
              <div>
                <div class="font-medium text-gray-800">学习提醒</div>
                <div class="text-sm text-gray-500">每日学习时间提醒</div>
              </div>
            </div>
            <label class="relative inline-flex items-center cursor-pointer">
              <input v-model="studyReminder" type="checkbox" class="sr-only peer" />
              <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"></div>
            </label>
          </div>
          
          <!-- 提醒时间 -->
          <div v-if="studyReminder" class="p-4">
            <div class="flex items-center space-x-3 mb-3">
              <Clock class="w-5 h-5 text-gray-500" />
              <div>
                <div class="font-medium text-gray-800">提醒时间</div>
                <div class="text-sm text-gray-500">设置每日提醒的时间</div>
              </div>
            </div>
            <input
              v-model="reminderTime"
              type="time"
              class="border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            >
          </div>
          
          <!-- 成就通知 -->
          <div class="p-4 flex items-center justify-between">
            <div class="flex items-center space-x-3">
              <Award class="w-5 h-5 text-gray-500" />
              <div>
                <div class="font-medium text-gray-800">成就通知</div>
                <div class="text-sm text-gray-500">获得成就时的通知</div>
              </div>
            </div>
            <label class="relative inline-flex items-center cursor-pointer">
              <input v-model="achievementNotification" type="checkbox" class="sr-only peer" />
              <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"></div>
            </label>
          </div>
        </div>
      </div>

      <!-- 其他设置 -->
      <div class="bg-white rounded-lg shadow-sm mb-4">
        <div class="p-4 border-b border-gray-100">
          <h2 class="text-lg font-semibold text-gray-800">其他</h2>
        </div>
        <div class="divide-y divide-gray-100">
          <!-- 主题设置 -->
          <div class="p-4">
            <div class="flex items-center space-x-3 mb-3">
              <Palette class="w-5 h-5 text-gray-500" />
              <div>
                <div class="font-medium text-gray-800">主题设置</div>
                <div class="text-sm text-gray-500">选择应用主题</div>
              </div>
            </div>
            <div class="space-y-2">
              <label class="flex items-center space-x-3">
                <input v-model="theme" type="radio" value="light" class="text-blue-600" />
                <span class="text-gray-700">浅色主题</span>
              </label>
              <label class="flex items-center space-x-3">
                <input v-model="theme" type="radio" value="dark" class="text-blue-600" />
                <span class="text-gray-700">深色主题</span>
              </label>
              <label class="flex items-center space-x-3">
                <input v-model="theme" type="radio" value="auto" class="text-blue-600" />
                <span class="text-gray-700">跟随系统</span>
              </label>
            </div>
          </div>
          
          <!-- 语言设置 -->
          <div class="p-4 flex items-center justify-between">
            <div class="flex items-center space-x-3">
              <Globe class="w-5 h-5 text-gray-500" />
              <div>
                <div class="font-medium text-gray-800">语言</div>
                <div class="text-sm text-gray-500">简体中文</div>
              </div>
            </div>
            <ChevronRight class="w-5 h-5 text-gray-400" />
          </div>
          
          <!-- 清除缓存 -->
          <div @click="clearCache" class="p-4 flex items-center justify-between cursor-pointer hover:bg-gray-50">
            <div class="flex items-center space-x-3">
              <Trash2 class="w-5 h-5 text-gray-500" />
              <div>
                <div class="font-medium text-gray-800">清除缓存</div>
                <div class="text-sm text-gray-500">清除本地缓存数据</div>
              </div>
            </div>
            <span class="text-sm text-gray-500">{{ cacheSize }}</span>
          </div>
        </div>
      </div>

      <!-- 关于 -->
      <div class="bg-white rounded-lg shadow-sm mb-4">
        <div class="p-4 border-b border-gray-100">
          <h2 class="text-lg font-semibold text-gray-800">关于</h2>
        </div>
        <div class="divide-y divide-gray-100">
          <!-- 版本信息 -->
          <div class="p-4 flex items-center justify-between">
            <div class="flex items-center space-x-3">
              <Info class="w-5 h-5 text-gray-500" />
              <div>
                <div class="font-medium text-gray-800">版本信息</div>
                <div class="text-sm text-gray-500">当前版本 v1.0.0</div>
              </div>
            </div>
            <button class="text-blue-600 text-sm hover:text-blue-700">检查更新</button>
          </div>
          
          <!-- 用户协议 -->
          <div class="p-4 flex items-center justify-between">
            <div class="flex items-center space-x-3">
              <FileText class="w-5 h-5 text-gray-500" />
              <div>
                <div class="font-medium text-gray-800">用户协议</div>
                <div class="text-sm text-gray-500">查看用户服务协议</div>
              </div>
            </div>
            <ChevronRight class="w-5 h-5 text-gray-400" />
          </div>
          
          <!-- 隐私政策 -->
          <div class="p-4 flex items-center justify-between">
            <div class="flex items-center space-x-3">
              <Shield class="w-5 h-5 text-gray-500" />
              <div>
                <div class="font-medium text-gray-800">隐私政策</div>
                <div class="text-sm text-gray-500">查看隐私保护政策</div>
              </div>
            </div>
            <ChevronRight class="w-5 h-5 text-gray-400" />
          </div>
          
          <!-- 意见反馈 -->
          <div class="p-4 flex items-center justify-between">
            <div class="flex items-center space-x-3">
              <MessageSquare class="w-5 h-5 text-gray-500" />
              <div>
                <div class="font-medium text-gray-800">意见反馈</div>
                <div class="text-sm text-gray-500">提交问题和建议</div>
              </div>
            </div>
            <ChevronRight class="w-5 h-5 text-gray-400" />
          </div>
        </div>
      </div>

      <!-- 退出登录 -->
      <div v-if="authStore.isLoggedIn && !authStore.isGuest" class="bg-white rounded-lg shadow-sm mb-6">
        <button 
          @click="logout"
          class="w-full p-4 text-red-600 font-medium hover:bg-red-50 rounded-lg transition-colors duration-200"
        >
          退出登录
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import {
  ArrowLeft, User, Shield, Lock, Target, BookOpen, BookmarkPlus,
  Bell, Clock, Award, Palette, Globe, Trash2, Info, FileText,
  MessageSquare, ChevronRight
} from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

// 学习设置
const dailyGoal = ref(20)
const answerMode = ref('practice')
const autoCollectMistakes = ref(true)

// 通知设置
const studyReminder = ref(true)
const reminderTime = ref('20:00')
const achievementNotification = ref(true)

// 其他设置
const theme = ref('light')
const cacheSize = ref('12.5 MB')

// 导航方法
const goBack = () => {
  router.back()
}

// 清除缓存
const clearCache = async () => {
  try {
    // 模拟清除缓存
    await new Promise(resolve => setTimeout(resolve, 1000))
    cacheSize.value = '0 MB'
    
    // 显示成功提示
    alert('缓存清除成功')
    
    // 重新计算缓存大小
    setTimeout(() => {
      cacheSize.value = '2.1 MB'
    }, 2000)
  } catch (error) {
    console.error('清除缓存失败:', error)
    alert('清除缓存失败，请重试')
  }
}

// 退出登录
const logout = async () => {
  try {
    const confirmed = confirm('确定要退出登录吗？')
    if (!confirmed) return
    
    // 执行退出登录
    authStore.logout()
    
    // 跳转到登录页
    router.push('/login')
  } catch (error) {
    console.error('退出登录失败:', error)
    alert('退出登录失败，请重试')
  }
}

// 保存设置
const saveSettings = () => {
  try {
    const settings = {
      dailyGoal: dailyGoal.value,
      answerMode: answerMode.value,
      autoCollectMistakes: autoCollectMistakes.value,
      studyReminder: studyReminder.value,
      reminderTime: reminderTime.value,
      achievementNotification: achievementNotification.value,
      theme: theme.value
    }
    
    // 保存到本地存储
    localStorage.setItem('userSettings', JSON.stringify(settings))
    console.log('设置已保存:', settings)
  } catch (error) {
    console.error('保存设置失败:', error)
  }
}

// 加载设置
const loadSettings = () => {
  try {
    const savedSettings = localStorage.getItem('userSettings')
    if (savedSettings) {
      const settings = JSON.parse(savedSettings)
      dailyGoal.value = settings.dailyGoal || 20
      answerMode.value = settings.answerMode || 'practice'
      autoCollectMistakes.value = settings.autoCollectMistakes !== false
      studyReminder.value = settings.studyReminder !== false
      reminderTime.value = settings.reminderTime || '20:00'
      achievementNotification.value = settings.achievementNotification !== false
      theme.value = settings.theme || 'light'
    }
  } catch (error) {
    console.error('加载设置失败:', error)
  }
}

// 监听设置变化并自动保存
const watchSettings = () => {
  // 这里可以使用 watch 来监听设置变化并自动保存
  // 为了简化，我们在每次变化时都保存
  const saveOnChange = () => {
    setTimeout(saveSettings, 100) // 延迟保存避免频繁操作
  }
  
  // 监听各种设置变化
  // 实际项目中可以使用 Vue 的 watch API
}

onMounted(() => {
  loadSettings()
  watchSettings()
})
</script>

<style scoped>
/* 自定义滑块样式 */
.slider::-webkit-slider-thumb {
  appearance: none;
  height: 20px;
  width: 20px;
  border-radius: 50%;
  background: #3B82F6;
  cursor: pointer;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.slider::-moz-range-thumb {
  height: 20px;
  width: 20px;
  border-radius: 50%;
  background: #3B82F6;
  cursor: pointer;
  border: none;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}
</style>