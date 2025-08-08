<template>
  <div class="admin-layout">
    <!-- 侧边栏 -->
    <aside class="sidebar" :class="{ collapsed: sidebarCollapsed }">
      <div class="sidebar-header">
        <div class="logo">
          <el-icon class="logo-icon"><Management /></el-icon>
          <span v-show="!sidebarCollapsed" class="logo-text">刷刷题管理</span>
        </div>
      </div>
      
      <nav class="sidebar-nav">
        <el-menu
          :default-active="activeMenu"
          :collapse="sidebarCollapsed"
          :unique-opened="true"
          router
          class="sidebar-menu"
        >
          <el-menu-item index="/">
            <el-icon><House /></el-icon>
            <template #title>仪表板</template>
          </el-menu-item>
          
          <el-sub-menu index="questions">
            <template #title>
              <el-icon><Document /></el-icon>
              <span>题目管理</span>
            </template>
            <el-menu-item index="/questions">题目列表</el-menu-item>
            <el-menu-item index="/questions/create">添加题目</el-menu-item>
            <el-menu-item index="/questions/import">批量导入</el-menu-item>
          </el-sub-menu>
          
          <el-menu-item index="/categories">
            <el-icon><Folder /></el-icon>
            <template #title>分类管理</template>
          </el-menu-item>
          
          <el-menu-item index="/users">
            <el-icon><User /></el-icon>
            <template #title>用户管理</template>
          </el-menu-item>
          
          <el-sub-menu index="system">
            <template #title>
              <el-icon><Setting /></el-icon>
              <span>系统管理</span>
            </template>
            <el-menu-item index="/settings">系统设置</el-menu-item>
            <el-menu-item index="/operation-logs">操作日志</el-menu-item>
          </el-sub-menu>
        </el-menu>
      </nav>
    </aside>
    
    <!-- 主内容区域 -->
    <div class="main-container">
      <!-- 顶部栏 -->
      <header class="header">
        <div class="header-left">
          <el-button
            type="text"
            class="sidebar-toggle"
            @click="toggleSidebar"
          >
            <el-icon><Expand v-if="sidebarCollapsed" /><Fold v-else /></el-icon>
          </el-button>
          
          <el-breadcrumb class="breadcrumb">
            <el-breadcrumb-item
              v-for="item in breadcrumbs"
              :key="item.path"
              :to="item.path"
            >
              {{ item.title }}
            </el-breadcrumb-item>
          </el-breadcrumb>
        </div>
        
        <div class="header-right">
          <!-- 用户信息下拉菜单 -->
          <el-dropdown @command="handleUserCommand">
            <div class="user-info">
              <el-avatar :size="32" :src="userAvatar">
                <el-icon><User /></el-icon>
              </el-avatar>
              <span class="username">{{ username }}</span>
              <el-icon class="dropdown-icon"><ArrowDown /></el-icon>
            </div>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="profile">
                  <el-icon><User /></el-icon>
                  个人资料
                </el-dropdown-item>
                <el-dropdown-item command="settings">
                  <el-icon><Setting /></el-icon>
                  账户设置
                </el-dropdown-item>
                <el-dropdown-item divided command="logout">
                  <el-icon><SwitchButton /></el-icon>
                  退出登录
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </header>
      
      <!-- 页面内容 -->
      <main class="content">
        <router-view />
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessageBox, ElMessage } from 'element-plus'
import {
  Management,
  House,
  Document,
  Folder,
  User,
  Setting,
  Expand,
  Fold,
  ArrowDown,
  SwitchButton
} from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

// 侧边栏折叠状态
const sidebarCollapsed = ref(false)

// 当前激活的菜单
const activeMenu = computed(() => route.path)

// 用户信息
const username = computed(() => authStore.user?.username || '管理员')
const userAvatar = computed(() => authStore.user?.avatar || '')

// 面包屑导航
const breadcrumbs = computed(() => {
  const matched = route.matched.filter(item => item.meta && item.meta.title)
  const breadcrumbList = matched.map(item => ({
    path: item.path,
    title: item.meta?.title || ''
  }))
  
  // 如果没有匹配到，根据路径生成默认面包屑
  if (breadcrumbList.length === 0) {
    const pathMap: Record<string, string> = {
      '/': '仪表板',
      '/questions': '题目管理',
      '/questions/create': '添加题目',
      '/questions/import': '批量导入',
      '/categories': '分类管理',
      '/users': '用户管理',
      '/settings': '系统设置',
      '/operation-logs': '操作日志'
    }
    
    const title = pathMap[route.path] || '未知页面'
    return [{ path: route.path, title }]
  }
  
  return breadcrumbList
})

// 切换侧边栏
const toggleSidebar = () => {
  sidebarCollapsed.value = !sidebarCollapsed.value
  // 保存状态到localStorage
  localStorage.setItem('sidebar_collapsed', sidebarCollapsed.value.toString())
}

// 处理用户下拉菜单命令
const handleUserCommand = async (command: string) => {
  switch (command) {
    case 'profile':
      ElMessage.info('个人资料功能开发中')
      break
    case 'settings':
      ElMessage.info('账户设置功能开发中')
      break
    case 'logout':
      await handleLogout()
      break
  }
}

// 处理退出登录
const handleLogout = async () => {
  try {
    await ElMessageBox.confirm(
      '确定要退出登录吗？',
      '提示',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await authStore.logout()
    ElMessage.success('已退出登录')
    router.push('/login')
  } catch (error) {
    // 用户取消操作
  }
}

// 初始化侧边栏状态
const initSidebarState = () => {
  const saved = localStorage.getItem('sidebar_collapsed')
  if (saved !== null) {
    sidebarCollapsed.value = saved === 'true'
  }
}

// 监听路由变化，更新面包屑
watch(
  () => route.path,
  () => {
    // 可以在这里添加页面切换的逻辑
  },
  { immediate: true }
)

// 初始化
initSidebarState()
</script>

<style scoped>
.admin-layout {
  display: flex;
  height: 100vh;
  background-color: #f5f5f5;
}

/* 侧边栏样式 */
.sidebar {
  width: 240px;
  background-color: #001529;
  transition: width 0.3s;
  overflow: hidden;
}

.sidebar.collapsed {
  width: 64px;
}

.sidebar-header {
  height: 64px;
  display: flex;
  align-items: center;
  padding: 0 16px;
  border-bottom: 1px solid #1f2937;
}

.logo {
  display: flex;
  align-items: center;
  color: white;
}

.logo-icon {
  font-size: 24px;
  margin-right: 8px;
}

.logo-text {
  font-size: 18px;
  font-weight: 600;
  white-space: nowrap;
}

.sidebar-nav {
  height: calc(100vh - 64px);
  overflow-y: auto;
}

.sidebar-menu {
  border: none;
  background-color: transparent;
}

.sidebar-menu :deep(.el-menu-item),
.sidebar-menu :deep(.el-sub-menu__title) {
  color: #bfbfbf;
  border-bottom: none;
}

.sidebar-menu :deep(.el-menu-item:hover),
.sidebar-menu :deep(.el-sub-menu__title:hover) {
  background-color: #1f2937;
  color: white;
}

.sidebar-menu :deep(.el-menu-item.is-active) {
  background-color: #1890ff;
  color: white;
}

/* 主内容区域样式 */
.main-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

/* 顶部栏样式 */
.header {
  height: 64px;
  background-color: white;
  border-bottom: 1px solid #e8e8e8;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 16px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.08);
}

.header-left {
  display: flex;
  align-items: center;
}

.sidebar-toggle {
  margin-right: 16px;
  font-size: 18px;
}

.breadcrumb {
  font-size: 14px;
}

.header-right {
  display: flex;
  align-items: center;
}

.user-info {
  display: flex;
  align-items: center;
  cursor: pointer;
  padding: 8px 12px;
  border-radius: 6px;
  transition: background-color 0.3s;
}

.user-info:hover {
  background-color: #f5f5f5;
}

.username {
  margin: 0 8px;
  font-size: 14px;
  color: #333;
}

.dropdown-icon {
  font-size: 12px;
  color: #999;
}

/* 内容区域样式 */
.content {
  flex: 1;
  padding: 24px;
  overflow-y: auto;
  background-color: #f5f5f5;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .sidebar {
    position: fixed;
    left: 0;
    top: 0;
    z-index: 1000;
    height: 100vh;
  }
  
  .sidebar.collapsed {
    left: -240px;
  }
  
  .main-container {
    margin-left: 0;
  }
  
  .content {
    padding: 16px;
  }
}
</style>