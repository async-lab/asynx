<template>
  <div class="dashboard">
    <!-- 顶部导航栏 -->
    <el-header class="dashboard-header">
      <div class="header-left">
        <h2>AsyncLab 仪表板</h2>
      </div>
      <div class="header-right">
        <el-dropdown @command="handleCommand">
          <span class="user-info">
            <span class="username">{{ fullName }}</span>
            <el-icon><ArrowDown /></el-icon>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="profile">设置</el-dropdown-item>
              <el-dropdown-item divided command="logout">退出登录</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </el-header>

    <!-- 主要内容区域 -->
    <el-container class="dashboard-container">
      <!-- 侧边栏 -->
      <el-aside width="200px" class="dashboard-sidebar">
        <el-menu
          :default-active="activeMenu"
          class="sidebar-menu"
          @select="handleMenuSelect"
        >
          <el-menu-item index="overview">
            <el-icon><User /></el-icon>
            <span>概览</span>
          </el-menu-item>
          <el-menu-item index="users">
            <el-icon><User /></el-icon>
            <span>用户管理</span>
          </el-menu-item>
        </el-menu>
      </el-aside>

      <!-- 主内容区 -->
      <el-main class="dashboard-main">
        <el-card>
          <template #header>
            <div class="card-header">
              <h3>{{ getPageTitle() }}</h3>
            </div>
          </template>
          
          <div class="dashboard-content">
            <div v-if="activeMenu === 'overview'" class="overview-page">
              <h4>系统概览</h4>
              <p>欢迎使用 AsyncLab 仪表板</p>
            </div>
            <UsersPage 
              v-else-if="activeMenu === 'users'"
              :users="users"
              @create="createUser"
              @edit="editUser"
              @delete="deleteUser"
            />
          </div>
        </el-card>
      </el-main>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { removeToken, getUserProfile, clearUserProfile } from '@/utils/auth'
import { useWarningConfirm } from '@/utils/msgTip'
import { 
  ArrowDown, 
  User 
} from '@element-plus/icons-vue'
import UsersPage from '@/components/dashboard/UsersPage.vue'

type MenuKey = 'overview' | 'projects' | 'users'

interface UserItem {
  username: string
  email: string
  role: string
  status: 'active' | 'disabled'
}

const router = useRouter()
const activeMenu = ref<MenuKey>('overview')
const fullName = computed(() => {
  const profile = getUserProfile() as any
  const name = (profile?.surName ?? '') + (profile?.givenName ?? '')
  return name || '用户'
})

// 用户数据
const users = ref<UserItem[]>([
  { username: 'admin', email: 'admin@example.com', role: '管理员', status: 'active' },
  { username: 'user1', email: 'user1@example.com', role: '普通用户', status: 'active' },
  { username: 'user2', email: 'user2@example.com', role: '普通用户', status: 'disabled' }
])

// 获取页面标题
const getPageTitle = () => {
  const titles: Record<MenuKey, string> = {
    overview: '系统概览',
    projects: '项目管理',
    users: '用户管理'
  }
  return titles[activeMenu.value]
}

// 处理菜单选择
const handleMenuSelect = (index: string) => {
  activeMenu.value = index as MenuKey
}

// 处理用户下拉菜单命令
const handleCommand = async (command: string) => {
  switch (command) {
    case 'profile':
      router.push('/profile')
      break
    case 'logout':
      await handleLogout()
      break
  }
}

// 处理退出登录
const handleLogout = async () => {
  try {
    await useWarningConfirm('确定要退出登录吗？')
    removeToken()
    clearUserProfile()
    router.push('/login')
  } catch {
    // 用户取消退出
  }
}

// 用户相关方法
const createUser = () => {
  console.log('创建用户')
}

const editUser = (user: any) => {
  console.log('编辑用户:', user)
}

const deleteUser = (user: any) => {
  console.log('删除用户:', user)
}


</script>

<style scoped>
.dashboard {
  min-height: 100vh;
  min-width: 1200px;
  display: flex;
  flex-direction: column;
}

.dashboard-header {
  min-height: 60px;
  background: #fff;
  border-bottom: 1px solid #e4e7ed;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.header-left h2 {
  margin: 0;
  color: #409eff;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}

.username {
  font-weight: 500;
}

.dashboard-container {
  flex: 1;
  overflow: hidden;
  min-height: calc(100vh - 60px);
}

.dashboard-sidebar {
  background: #fff;
  border-right: 1px solid #e4e7ed;
  min-width: 200px;
  width: 200px;
}

.sidebar-menu {
  border-right: none;
}

.dashboard-main {
  padding: 20px;
  background: #f5f7fa;
  min-width: 800px;
  min-height: calc(100vh - 60px);
}

.dashboard-content {
  min-height: 400px;
  min-width: 600px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  min-height: 50px;
}

.card-header h3 {
  margin: 0;
  color: #303133;
}

.overview-page,
.projects-page {
  min-height: 500px;
  min-width: 600px;
}
</style> 