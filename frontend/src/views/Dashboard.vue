<template>
  <div class="dashboard">
    <!-- 顶部导航栏 -->
    <el-header class="dashboard-header">
      <div class="header-left">
        <h2>AsyncLab</h2>
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
              <el-dropdown-item divided command="OIDC">OIDC设置</el-dropdown-item>
              <el-dropdown-item divided command="logout">退出登录</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </el-header>

    <!-- 移动端顶部导航（<= 992px 显示） -->
    <div class="mobile-nav">
      <el-menu
        mode="horizontal"
        :default-active="activeMenu"
        @select="handleMenuSelect"
      >
        <el-menu-item index="overview">
          <el-icon><House /></el-icon>
          <span>概览</span>
        </el-menu-item>
        <el-menu-item index="users" v-if="!isRestricted">
          <el-icon><User /></el-icon>
          <span>用户管理</span>
        </el-menu-item>
      </el-menu>
    </div>

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
            <el-icon><House /></el-icon>
            <span>概览</span>
          </el-menu-item>
          <el-menu-item index="users" v-if="!isRestricted">
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
              <HomeHero compact />
            </div>
            <UsersPage 
              v-else-if="activeMenu === 'users'"
              :users="users"
              :is-admin="isAdmin"
              :loading="usersLoading"
              @create="createUser"
              @refresh="loadUsers"
            />
          </div>
        </el-card>
      </el-main>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watchEffect, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { removeToken, getUserProfile, clearUserProfile } from '@/utils/auth'
import { useWarningConfirm } from '@/utils/msgTip'
import { 
  ArrowDown, 
  House, 
  User 
} from '@element-plus/icons-vue'
import UsersPage from '@/components/dashboard/UsersPage.vue'
import HomeHero from '@/components/HomeHero.vue'
import { getUserList } from '@/api/user'
import type { User as ApiUser } from '@/api/types'

type MenuKey = 'overview' | 'projects' | 'users'

const router = useRouter()
const activeMenu = ref<MenuKey>('overview')
const fullName = computed(() => {
  const profile = getUserProfile() as any
  const name = (profile?.surName ?? '') + (profile?.givenName ?? '')
  return name || '用户'
})

// 用户数据
const users = ref<ApiUser[]>([])
const usersLoading = ref<boolean>(false)

const profile = computed(() => (getUserProfile() as any) || {})
const isAdmin = computed(() => profile.value?.role === 'admin')
const isDefault = computed(() => profile.value?.role === 'default')
const isRestricted = computed(() => profile.value?.role === 'restricted')

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
    case 'OIDC':
      window.open('https://keycloak.internal.asynclab.club/realms/asynclab/account  ', '_blank')
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

// 已由 UsersPage 内部处理编辑/删除，通过 refresh 事件回传，这里无需实现

// 加载用户列表并按权限过滤
onMounted(async () => {
  if (activeMenu.value !== 'users') return
})

// 进入用户管理时拉取数据
const loadUsers = async () => {
  usersLoading.value = true
  const list = await getUserList() as any
  let data: ApiUser[] = Array.isArray(list?.data) ? list.data : (Array.isArray(list) ? list : [])
  if (isAdmin.value) {
    users.value = data
  } else if (isDefault.value) {
    // 仅同组（按 category 分组）
    const myCategory = profile.value?.category
    users.value = data.filter((u: any) => u.category === myCategory && u.role === 'default')
  } else if (isRestricted.value) {
    users.value = []
  }
  await nextTick()
  usersLoading.value = false
}

// 切换到用户管理菜单时触发
watchEffect(() => {
  if (activeMenu.value === 'users') {
    if (isRestricted.value) {
      // 无权限访问
      users.value = []
      usersLoading.value = false
    } else {
      loadUsers()
    }
  }
})


</script>

<style scoped>
.dashboard {
  min-height: 100vh;
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
  min-height: calc(100vh - 60px);
}

.dashboard-content {
  min-height: 400px;
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
}

@media (max-width: 1200px) {
  .dashboard-header { padding: 0 12px; }
  .dashboard-main { padding: 12px; }
}

@media (max-width: 992px) {
  .dashboard-container { display: block; }
  .dashboard-sidebar { display: none; }
}

@media (max-width: 768px) {
  .dashboard-header { min-height: 56px; }
  .dashboard-main { padding: 10px; }
}

.mobile-nav {
  display: none;
  background: #fff;
  border-bottom: 1px solid #e4e7ed;
}

@media (max-width: 992px) {
  .mobile-nav { display: block; }
}

/* 概览为空白占位，不添加多余内容 */

.overview-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.brand-logo {
  width: 120px;
  height: 120px;
}
</style> 
<style>
/* 暗色模式覆盖：Dashboard 头部/侧栏/容器/移动导航 */
html.dark .dashboard-header {
  background: var(--card-bg);
  border-bottom-color: var(--border-color);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.35);
}
html.dark .header-left h2 { color: var(--el-color-primary); }

html.dark .dashboard-sidebar {
  background: var(--card-bg);
  border-right-color: var(--border-color);
}

html.dark .dashboard-main {
  background: #0f131a;
}

html.dark .mobile-nav {
  background: var(--card-bg);
  border-bottom-color: var(--border-color);
}

html.dark .card-header h3 { color: var(--el-text-color-primary); }
</style>