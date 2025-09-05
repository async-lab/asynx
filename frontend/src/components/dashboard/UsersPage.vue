<template>
  <div class="users-page">
    <el-button type="primary" @click="$emit('create')">
      新建用户
    </el-button>
    <el-table :data="users" style="margin-top: 20px;">
      <el-table-column prop="username" label="用户名" />
      <el-table-column prop="email" label="邮箱" />
      <el-table-column prop="role" label="角色" />
      <el-table-column prop="status" label="状态">
        <template #default="scope">
          <el-tag :type="scope.row.status === 'active' ? 'success' : 'danger'">
            {{ scope.row.status === 'active' ? '活跃' : '禁用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作">
        <template #default="scope">
          <el-button size="small" @click="$emit('edit', scope.row)">
            编辑
          </el-button>
          <el-button size="small" type="danger" @click="$emit('delete', scope.row)">
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script setup lang="ts">
import { defineProps } from 'vue'

interface UserItem {
  username: string
  email: string
  role: string
  status: 'active' | 'disabled'
}

defineProps<{ users: UserItem[] }>()
</script>

<style scoped>
.users-page {
  min-height: 500px;
  min-width: 600px;
}
</style>

