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
          style="width: 260px; margin-left: 12px;"
        />
      </div>
      <div class="filters" v-if="isAdmin">
        <el-select v-model="selectedRoles" multiple collapse-tags collapse-tags-tooltip placeholder="角色" clearable style="width: 260px;">
          <el-option label="admin" value="admin" />
          <el-option label="default" value="default" />
          <el-option label="restricted" value="restricted" />
        </el-select>
        <el-select v-model="selectedCategories" multiple collapse-tags collapse-tags-tooltip placeholder="账号类型" clearable style="width: 280px; margin-left: 12px;">
          <el-option label="system" value="system" />
          <el-option label="member" value="member" />
          <el-option label="external" value="external" />
        </el-select>
      </div>
    </div>

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

    <div class="pagination">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="pageSizes"
        :total="filteredTotal"
        background
        layout="total, sizes, prev, pager, next, jumper"
      />
    </div>
  </div>
    <!-- 编辑对话框（仅管理员） -->
    <el-dialog v-model="editVisible" title="编辑用户" width="520px">
      <el-form label-width="96px">
        <el-form-item label="用户名">
          <el-input v-model="editForm.username" disabled />
        </el-form-item>
        <el-form-item label="角色">
          <el-select v-model="editForm.role" placeholder="选择角色" style="width: 240px;">
            <el-option label="admin" value="admin" />
            <el-option label="default" value="default" />
            <el-option label="restricted" value="restricted" />
          </el-select>
          <el-button style="margin-left: 12px;" type="primary" :loading="savingRole" @click="onSaveRole">保存角色</el-button>
        </el-form-item>
        <el-form-item label="账号类型">
          <el-select v-model="editForm.category" placeholder="选择类型" style="width: 240px;">
            <el-option label="system" value="system" />
            <el-option label="member" value="member" />
            <el-option label="external" value="external" />
          </el-select>
          <el-button style="margin-left: 12px;" type="primary" :loading="savingCategory" @click="onSaveCategory">保存类型</el-button>
        </el-form-item>
        <el-form-item label="新密码">
          <el-input v-model.trim="editForm.password" type="password" show-password style="width: 240px;" />
          <el-button style="margin-left: 12px;" type="warning" :disabled="!editForm.password" :loading="savingPwd" @click="onChangePwd">修改密码</el-button>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="editVisible = false">关闭</el-button>
        </span>
      </template>
    </el-dialog>
    <!-- 新建用户对话框（仅管理员） -->
    <el-dialog v-model="createVisible" title="新建用户" width="520px">
      <el-form label-width="96px">
        <el-form-item label="用户名">
          <el-input v-model.trim="createForm.username" />
        </el-form-item>
        <el-form-item label="姓">
          <el-input v-model.trim="createForm.surName" />
        </el-form-item>
        <el-form-item label="名">
          <el-input v-model.trim="createForm.givenName" />
        </el-form-item>
        <el-form-item label="邮箱">
          <el-input v-model.trim="createForm.mail" />
        </el-form-item>
        <el-form-item label="角色">
          <el-select v-model="createForm.role" placeholder="选择角色" style="width: 240px;">
            <el-option label="admin" value="admin" />
            <el-option label="default" value="default" />
            <el-option label="restricted" value="restricted" />
          </el-select>
        </el-form-item>
        <el-form-item label="账号类型">
          <el-select v-model="createForm.category" placeholder="选择类型" style="width: 240px;">
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
import { defineProps, defineEmits, computed, ref, watch } from 'vue'
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

// ===== 编辑逻辑（管理员） =====
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

// ===== 新建用户（管理员） =====
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
.filters {
  display: flex;
  align-items: center;
}
.pagination {
  margin-top: 12px;
  display: flex;
  justify-content: flex-end;
}

/* 小屏表格滚动与工具栏换行 */
@media (max-width: 992px) {
  .toolbar { flex-wrap: wrap; gap: 8px; }
}

@media (max-width: 768px) {
  .toolbar { flex-direction: column; align-items: stretch; }
}
</style>

