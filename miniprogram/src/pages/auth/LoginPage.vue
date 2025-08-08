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
      <div class="space-y-6">
        <!-- 微信登录按钮 -->
        <button
          @click="handleWechatLogin"
          :disabled="isLoading"
          class="w-full bg-green-500 hover:bg-green-600 disabled:bg-gray-400 text-white font-medium py-3 px-4 rounded-lg transition-colors duration-200 flex items-center justify-center space-x-2"
        >
          <Smartphone v-if="!isLoading" class="w-5 h-5" />
          <Loader2 v-else class="w-5 h-5 animate-spin" />
          <span>{{ isLoading ? '登录中...' : '微信授权登录' }}</span>
        </button>

        <!-- 游客登录 -->
        <button
          @click="handleGuestLogin"
          :disabled="isLoading"
          class="w-full bg-gray-100 hover:bg-gray-200 disabled:bg-gray-50 text-gray-700 font-medium py-3 px-4 rounded-lg transition-colors duration-200 flex items-center justify-center space-x-2"
        >
          <User class="w-5 h-5" />
          <span>游客体验</span>
        </button>
      </div>

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
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { BookOpen, Smartphone, User, Loader2 } from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const isLoading = ref(false)

// 微信授权登录
const handleWechatLogin = async () => {
  try {
    isLoading.value = true
    
    // 模拟获取微信授权码
    const mockCode = 'mock_wechat_code_' + Date.now()
    
    // 调用store的微信登录方法
    await authStore.loginWithWechat(mockCode)
    
    // 获取重定向地址，如果没有则跳转到首页
    const redirect = route.query.redirect as string
    router.push(redirect || '/')
  } catch (error) {
    console.error('微信登录失败:', error)
    alert(error instanceof Error ? error.message : '微信登录失败，请重试')
  } finally {
    isLoading.value = false
  }
}

// 游客登录
const handleGuestLogin = async () => {
  try {
    isLoading.value = true
    
    // 调用store的游客登录方法
    await authStore.loginAsGuest()
    
    // 获取重定向地址，如果没有则跳转到首页
    const redirect = route.query.redirect as string
    router.push(redirect || '/')
  } catch (error) {
    console.error('游客登录失败:', error)
    alert(error instanceof Error ? error.message : '游客登录失败，请重试')
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