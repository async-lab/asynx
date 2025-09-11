<template>
  <div class="users-page">
    <div class="toolbar">
      <div class="left-actions">
        <el-button v-if="isAdmin" type="primary" @click="onCreate">
          新建用户
        </el-button>
        <el-input
          v-model.trim="searchQuery"
          placeholder="搜索用户名或姓名"
          clearable
          class="search-input"
        />
      </div>
      <div class="filters" v-if="isAdmin">
        <div class="filters-row">
          <el-select v-model="selectedRoles" multiple collapse-tags collapse-tags-tooltip placeholder="角色" clearable class="filter-select">
          <el-option label="admin" value="admin" />
          <el-option label="default" value="default" />
          <el-option label="restricted" value="restricted" />
          </el-select>
          <el-select v-model="selectedCategories" multiple collapse-tags collapse-tags-tooltip placeholder="账号类型" clearable class="filter-select ml-12">
          <el-option label="system" value="system" />
          <el-option label="member" value="member" />
          <el-option label="external" value="external" />
          </el-select>
        </div>
      </div>
    </div>

    <!-- 桌面端表格 -->
    <div class="table-wrapper" v-show="!isMobile">
      <el-table :data="pagedUsers" style="margin-top: 12px;" v-loading="loading">
        <el-table-column prop="username" label="用户名" />
        <el-table-column label="姓名">
          <template #default="scope">
            {{ (scope.row.surName || '') + (scope.row.givenName || '') }}
          </template>
        </el-table-column>
        <el-table-column prop="mail" label="邮箱" />
        <el-table-column prop="role" label="角色" />
        <el-table-column prop="category" label="账号类型" />
        <el-table-column v-if="isAdmin" label="操作" width="180">
          <template #default="scope">
            <el-button v-if="isAdmin" size="small" @click="onEdit(scope.row)">编辑</el-button>
            <el-button v-if="isAdmin" size="small" type="danger" @click="onDelete(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <!-- 移动端卡片列表 -->
    <div class="mobile-cards" v-show="isMobile">
      <el-empty v-if="!loading && pagedUsers.length === 0" description="暂无数据" />
      <el-card v-for="user in pagedUsers" :key="user.username" class="user-card" shadow="hover">
        <div class="card-row main">
          <div class="name">{{ (user.surName || '') + (user.givenName || '') || '未命名' }}</div>
          <div class="username">@{{ user.username }}</div>
        </div>
        <div class="card-row">
          <span class="label">邮箱</span>
          <span class="value break">{{ user.mail || '-' }}</span>
        </div>
        <div class="card-row meta">
          <span class="tag">{{ user.role }}</span>
          <span class="tag">{{ user.category }}</span>
        </div>
        <div class="card-actions" v-if="isAdmin">
          <el-button size="small" @click="onEdit(user)">编辑</el-button>
          <el-button size="small" type="danger" @click="onDelete(user)">删除</el-button>
        </div>
      </el-card>
    </div>

    <div class="pagination">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="pageSizes"
        :total="filteredTotal"
        background
        :small="isMobile"
        :layout="paginationLayout"
      />
    </div>
  </div>
    <!-- 编辑对话框（仅管理员） -->
    <el-dialog v-model="editVisible" title="编辑用户" :width="isMobile ? '96%' : '640px'" >
      <el-form :label-width="isMobile ? '0' : '96px'" :label-position="isMobile ? 'top' : 'right'" class="edit-form">
        <el-form-item label="用户名">
          <el-input v-model="editForm.username" disabled style="width: 100%" />
        </el-form-item>

        <div class="edit-grid single">
          <el-form-item label="角色" class="compact-item">
            <div class="control-with-action">
              <el-select v-model="editForm.role" placeholder="选择角色" class="control">
                <el-option label="admin" value="admin" />
                <el-option label="default" value="default" />
                <el-option label="restricted" value="restricted" />
              </el-select>
              <el-button type="primary" :loading="savingRole" @click="onSaveRole">保存角色</el-button>
            </div>
          </el-form-item>

          <el-form-item label="账号类型" class="compact-item">
            <div class="control-with-action">
              <el-select v-model="editForm.category" placeholder="选择类型" class="control">
                <el-option label="system" value="system" />
                <el-option label="member" value="member" />
                <el-option label="external" value="external" />
              </el-select>
              <el-button type="primary" :loading="savingCategory" @click="onSaveCategory">保存类型</el-button>
            </div>
          </el-form-item>
        </div>

        <el-form-item label="新密码" class="compact-item">
          <div class="control-with-action">
            <el-input v-model.trim="editForm.password" type="password" show-password style="width: 100%" />
            <el-button type="warning" :disabled="!editForm.password" :loading="savingPwd" @click="onChangePwd">修改密码</el-button>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="editVisible = false">关闭</el-button>
        </span>
      </template>
    </el-dialog>
    <!-- 新建用户对话框（仅管理员） -->
    <el-dialog v-model="createVisible" title="新建用户" :width="isMobile ? '96%' : '520px'" >
      <el-form :label-width="isMobile ? '0' : '96px'" :label-position="isMobile ? 'top' : 'right'">
        <el-form-item label="用户名">
          <el-input v-model.trim="createForm.username" style="width: 100%" />
        </el-form-item>
        <el-form-item label="姓">
          <el-input v-model.trim="createForm.surName" style="width: 100%" />
        </el-form-item>
        <el-form-item label="名">
          <el-input v-model.trim="createForm.givenName" style="width: 100%" />
        </el-form-item>
        <el-form-item label="邮箱">
          <el-input v-model.trim="createForm.mail" style="width: 100%" />
        </el-form-item>
        <el-form-item label="角色">
          <el-select v-model="createForm.role" placeholder="选择角色" style="width: 100%;">
            <el-option label="admin" value="admin" />
            <el-option label="default" value="default" />
            <el-option label="restricted" value="restricted" />
          </el-select>
        </el-form-item>
        <el-form-item label="账号类型">
          <el-select v-model="createForm.category" placeholder="选择类型" style="width: 100%;">
            <el-option label="system" value="system" />
            <el-option label="member" value="member" />
            <el-option label="external" value="external" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="createVisible = false">取消</el-button>
          <el-button type="primary" :loading="creating" @click="onSubmitCreate">创建</el-button>
        </span>
      </template>
    </el-dialog>
</template>

<script setup lang="ts">
import { defineProps, defineEmits, computed, ref, watch, onMounted, onBeforeUnmount } from 'vue'
import type { User } from '@/api/types'
import { modifyUserRole, modifyUserCategory, deleteUser, changePassword, registerUser } from '@/api/user'
import { useSuccessTip, useFailedTip, useWarningConfirm } from '@/utils/msgTip'

const props = defineProps<{ users: User[]; isAdmin?: boolean; loading?: boolean }>()
const emit = defineEmits(['refresh'])

// 多选筛选
const selectedRoles = ref<string[]>([])
const selectedCategories = ref<string[]>([])

// 分页
const currentPage = ref<number>(1)
const pageSize = ref<number>(20)
const pageSizes = [10, 20, 50, 100]

const filteredUsers = computed(() => {
  let list = props.users || []
  // 关键字搜索（用户名或姓名：姓+名）
  if (searchQuery.value) {
    const kw = searchQuery.value.toLowerCase()
    list = list.filter(u => {
      const username = (u.username || '').toLowerCase()
      const fullName = `${u.surName || ''}${u.givenName || ''}`.toLowerCase()
      return username.includes(kw) || fullName.includes(kw)
    })
  }
  if (props.isAdmin) {
    if (selectedRoles.value.length > 0) {
      list = list.filter(u => selectedRoles.value.includes(u.role))
    }
    if (selectedCategories.value.length > 0) {
      list = list.filter(u => selectedCategories.value.includes(u.category))
    }
  }
  return list
})

const filteredTotal = computed(() => filteredUsers.value.length)

const pagedUsers = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredUsers.value.slice(start, end)
})

// 当过滤条件变化时，重置到第一页
watch([selectedRoles, selectedCategories, () => props.users], () => {
  currentPage.value = 1
})

// 搜索与筛选
const searchQuery = ref<string>('')
watch(searchQuery, () => { currentPage.value = 1 })

const isMobile = ref<boolean>(false)

const updateIsMobile = () => {
  isMobile.value = window.matchMedia('(max-width: 768px)').matches
}

onMounted(() => {
  updateIsMobile()
  // 根据当前设备类型设置分页大小
  pageSize.value = isMobile.value ? 10 : 20
  window.addEventListener('resize', updateIsMobile)
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', updateIsMobile)
})

// 监听移动端切换动态调整每页数量
watch(isMobile, (v) => {
  pageSize.value = v ? 10 : 20
})

const paginationLayout = computed(() => isMobile.value ? 'prev, pager, next' : 'total, sizes, prev, pager, next, jumper')

const editVisible = ref(false)
const savingRole = ref(false)
const savingCategory = ref(false)
const savingPwd = ref(false)
const editForm = ref<{ username: string; role: string; category: string; password: string }>({
  username: '',
  role: '',
  category: '',
  password: ''
})

const onEdit = (row: User) => {
  editForm.value = {
    username: row.username,
    role: row.role,
    category: row.category,
    password: ''
  }
  editVisible.value = true
}

const onSaveRole = async () => {
  if (!props.isAdmin) return
  try {
    savingRole.value = true
    await modifyUserRole(editForm.value.username, { role: editForm.value.role } as any)
    useSuccessTip('角色已更新')
    emit('refresh')
  } catch (e: any) {
    useFailedTip(e?.msg || e?.message || '更新角色失败')
  } finally {
    savingRole.value = false
  }
}

const onSaveCategory = async () => {
  if (!props.isAdmin) return
  try {
    savingCategory.value = true
    await modifyUserCategory(editForm.value.username, { category: editForm.value.category } as any)
    useSuccessTip('账号类型已更新')
    emit('refresh')
  } catch (e: any) {
    useFailedTip(e?.msg || e?.message || '更新账号类型失败')
  } finally {
    savingCategory.value = false
  }
}

const onChangePwd = async () => {
  if (!props.isAdmin) return
  if (!editForm.value.password) return
  try {
    savingPwd.value = true
    // 管理员修改目标用户密码
    await changePassword(editForm.value.username, { password: editForm.value.password })
    useSuccessTip('密码已更新')
    editForm.value.password = ''
  } catch (e: any) {
    useFailedTip(e?.msg || e?.message || '更新密码失败')
  } finally {
    savingPwd.value = false
  }
}

const onDelete = async (row: User) => {
  if (!props.isAdmin) return
  try {
    await useWarningConfirm(`确认删除用户「${row.username}」吗？`)
  } catch {
    return
  }
  try {
    await deleteUser(row.username)
    useSuccessTip('用户已删除')
    emit('refresh')
  } catch (e: any) {
    useFailedTip(e?.msg || e?.message || '删除失败')
  }
}

const createVisible = ref(false)
const creating = ref(false)
const createForm = ref<{ username: string; surName: string; givenName: string; mail: string; role: string; category: string }>({
  username: '',
  surName: '',
  givenName: '',
  mail: '',
  role: 'default',
  category: 'member'
})

const onCreate = () => {
  if (!props.isAdmin) return
  createVisible.value = true
}

const resetCreateForm = () => {
  createForm.value = { username: '', surName: '', givenName: '', mail: '', role: 'default', category: 'member' }
}

const onSubmitCreate = async () => {
  if (!props.isAdmin) return
  if (!createForm.value.username) { useFailedTip('请输入用户名'); return }
  try {
    creating.value = true
    await registerUser(createForm.value as any)
    useSuccessTip('用户已创建')
    createVisible.value = false
    resetCreateForm()
    emit('refresh')
  } catch (e: any) {
    useFailedTip(e?.msg || e?.message || '创建失败')
  } finally {
    creating.value = false
  }
}
</script>

<style scoped>
.users-page {
  min-height: 500px;
}
.toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.left-actions {
  display: flex;
  align-items: center;
}
.search-input { width: 260px; margin-left: 12px; }
.filters {
  display: flex;
  align-items: center;
}
.filters-row { display: flex; align-items: center; }
.filter-select { width: 260px; }
.ml-12 { margin-left: 12px; }
.pagination {
  margin-top: 12px;
  display: flex;
  justify-content: flex-end;
}

/* 表格容器横向滚动，避免列过多溢出 */
.table-wrapper { overflow-x: auto; }
.table-wrapper .el-table { min-width: 720px; }

/* 移动端卡片列表样式 */
.mobile-cards { display: grid; grid-template-columns: 1fr; gap: 12px; margin-top: 12px; }
.user-card { border-radius: 10px; }
.card-row { display: flex; align-items: center; justify-content: space-between; gap: 8px; margin-bottom: 8px; }
.card-row.main { margin-bottom: 4px; }
.name { font-weight: 600; color: #303133; }
.username { color: #909399; font-size: 12px; }
.label { color: #909399; min-width: 48px; }
.value { color: #606266; }
.value.break { word-break: break-all; overflow-wrap: anywhere; }
.card-row.meta { justify-content: flex-start; gap: 8px; }
.tag { background: #f4f4f5; color: #606266; padding: 2px 8px; border-radius: 10px; font-size: 12px; }
.card-actions { display: flex; gap: 8px; margin-top: 8px; }

/* 小屏表格滚动与工具栏换行 */
@media (max-width: 992px) {
  .toolbar { flex-wrap: wrap; gap: 8px; align-items: stretch; }
  .left-actions { width: 100%; }
  /* 搜索框小屏固定宽度，不跟随页面变化 */
  .search-input { width: 240px; margin-left: 0; }
  .filters { width: 100%; justify-content: space-between; }
  .filters-row { width: 100%; gap: 8px; margin-top: 4px; flex-wrap: wrap; }
  .filter-select { width: 100%; }
  .ml-12 { margin-left: 0; }
  .pagination { justify-content: center; }
}

@media (max-width: 768px) {
  .toolbar { flex-direction: column; align-items: stretch; }
  /* 对话框内联按钮在小屏改为换行占满 */
  :deep(.el-dialog__body) .dialog-inline-action {
    margin-left: 0 !important;
    margin-top: 8px;
    width: 100%;
  }
  :deep(.el-dialog) {
    margin: 0 !important;
  }
  :deep(.el-dialog__footer) {
    padding: 10px 12px;
  }
}

/* 编辑对话框响应式网格与“控件+动作”布局 */
.edit-form .edit-grid { display: grid; gap: 12px 0; width: 100%;}
.edit-form .edit-grid.single { grid-template-columns: 1fr; }
.edit-form .compact-item :deep(.el-form-item__content) {
  width: 100%;
}
.control-with-action {
  display: grid;
  grid-template-columns: 1fr max-content;
  gap: 12px;
  align-items: center;
}
/* 桌面端给下拉控件一个舒适的最小宽度，避免被过度压缩 */
.control-with-action .control { min-width: 350px; }


@media (max-width: 992px) {
  .control-with-action { grid-template-columns: 1fr; }
  .control-with-action .control { min-width: 180px; width: 100%; }
}

/* 暗色模式：用户管理页面文字与标签对比度增强 */
html.dark .users-page .name { color: var(--el-text-color-primary); }
html.dark .users-page .username { color: var(--el-text-color-secondary); }
html.dark .users-page .label { color: var(--el-text-color-secondary); }
html.dark .users-page .value { color: var(--el-text-color-regular); }
html.dark .users-page .tag {
  background: rgba(255, 255, 255, 0.06);
  color: var(--el-text-color-regular);
  border: 1px solid var(--border-color);
}

/* 表格内标题与正文在暗色下略提亮（Element Plus 已做大部分适配，这里微调） */
html.dark .users-page .el-table th.el-table__cell {
  color: var(--el-text-color-primary);
}
html.dark .users-page .el-table .cell {
  color: var(--el-text-color-regular);
}
</style>

