<template>
  <div class="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 flex items-center justify-center p-4">
    <div class="bg-white rounded-2xl shadow-xl p-8 w-full max-w-md">
      <!-- Logo和标题 -->
      <div class="text-center mb-8">
        <div class="w-20 h-20 bg-gradient-to-r from-blue-500 to-indigo-600 rounded-full mx-auto mb-4 flex items-center justify-center">
          <BookOpen class="w-10 h-10 text-white" />
        </div>
        <h1 class="text-2xl font-bold text-gray-800 mb-2">刷刷题</h1>
        <p class="text-gray-600">专业的在线刷题平台</p>
      </div>

      <!-- 登录表单 -->
      <form @submit.prevent="handleLogin" class="space-y-6">
        <!-- 用户名输入框 -->
        <div>
          <label for="username" class="block text-sm font-medium text-gray-700 mb-2">用户名</label>
          <input
            id="username"
            v-model="username"
            type="text"
            required
            placeholder="请输入用户名"
            class="w-full px-3 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          />
        </div>

        <!-- 密码输入框 -->
        <div>
          <label for="password" class="block text-sm font-medium text-gray-700 mb-2">密码</label>
          <input
            id="password"
            v-model="password"
            type="password"
            required
            placeholder="请输入密码"
            class="w-full px-3 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          />
        </div>

        <!-- 登录按钮 -->
        <button
          type="submit"
          :disabled="isLoading || !username || !password"
          class="w-full bg-blue-500 hover:bg-blue-600 disabled:bg-gray-400 text-white font-medium py-3 px-4 rounded-lg transition-colors duration-200 flex items-center justify-center space-x-2"
        >
          <Loader2 v-if="isLoading" class="w-5 h-5 animate-spin" />
          <User v-else class="w-5 h-5" />
          <span>{{ isLoading ? '登录中...' : '登录' }}</span>
        </button>
      </form>

      <!-- 服务条款 -->
      <div class="mt-8 text-center">
        <p class="text-xs text-gray-500">
          登录即表示同意
          <a href="#" class="text-blue-500 hover:underline">用户协议</a>
          和
          <a href="#" class="text-blue-500 hover:underline">隐私政策</a>
        </p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { BookOpen, User, Loader2 } from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const isLoading = ref(false)
const username = ref('')
const password = ref('')

// 用户名密码登录
const handleLogin = async () => {
  try {
    isLoading.value = true
    
    // 调用store的密码登录方法
    await authStore.loginWithPassword(username.value, password.value)
    
    // 获取重定向地址，如果没有则跳转到首页
    const redirect = route.query.redirect as string
    router.push(redirect || '/')
  } catch (error) {
    console.error('登录失败:', error)
    alert(error instanceof Error ? error.message : '登录失败，请重试')
  } finally {
    isLoading.value = false
  }
}

// 检查是否已登录
onMounted(() => {
  if (authStore.isLoggedIn) {
    const redirect = route.query.redirect as string
    router.push(redirect || '/')
  }
})
</script>

<style scoped>
/* 自定义样式 */
.bg-gradient-to-br {
  background: linear-gradient(to bottom right, var(--tw-gradient-stops));
}
</style>