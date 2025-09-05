<template>
  <div class="profile-page">
    <div class="profile-wrapper">
    <el-card class="profile-card">
      <template #header>
        <div class="card-header">
          <h3>个人资料</h3>
          <el-button type="primary" link @click="goDashboard">返回仪表板</el-button>
        </div>
      </template>

      <div class="info-list">
        <div class="info-item" v-for="item in infoItems" :key="item.label">
          <div class="info-label">{{ item.label }}</div>
          <div class="info-value">{{ item.value }}</div>
        </div>
      </div>
    </el-card>

    <el-card class="password-card" style="margin-top: 16px;">
      <template #header>
        <div class="card-header">
          <h3>修改密码</h3>
        </div>
      </template>

      <el-form :model="pwdForm" :rules="pwdRules" ref="pwdFormRef" label-width="100px">
        <el-form-item prop="password" label="新密码">
          <el-input v-model.trim="pwdForm.password" type="password" show-password />
        </el-form-item>
        <el-form-item prop="confirm" label="确认密码">
          <el-input v-model.trim="pwdForm.confirm" type="password" show-password />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="saving" :disabled="!canSubmit" @click="onChangePassword">保存</el-button>
        </el-form-item>
      </el-form>
    </el-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
const router = useRouter()
import type { FormInstance, FormRules } from 'element-plus'
import { getUserProfile } from '@/utils/auth'
import { changePassword } from '@/api/user'
import { useSuccessTip, useFailedTip } from '@/utils/msgTip'
import type { User } from '@/api/types'

const emptyProfile: User = { username: '', givenName: '', surName: '', mail: '' }
const profile = reactive<User>({ ...emptyProfile, ...(getUserProfile() || {}) })
const infoItems = computed(() => [
  { label: '用户名', value: profile.username || '-' },
  { label: '姓', value: profile.surName || '-' },
  { label: '名', value: profile.givenName || '-' },
  { label: '邮箱', value: profile.mail || '-' },
  { label: '角色', value: (profile as any)?.role || '-' },
  { label: '类型', value: (profile as any)?.category || '-' }
])

// 密码表单
const pwdFormRef = ref<FormInstance>()
const saving = ref(false)
const pwdForm = reactive({ password: '', confirm: '' })
const canSubmit = computed(() => pwdForm.password.length > 0 && pwdForm.password === pwdForm.confirm && !saving.value)

const pwdRules: FormRules<typeof pwdForm> = {
  password: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '至少6位字符', trigger: 'blur' }
  ],
  confirm: [
    { required: true, message: '请再次输入密码', trigger: 'blur' },
    {
      validator: (_r, v, cb) => {
        if (v !== pwdForm.password) cb(new Error('两次输入的密码不一致'))
        else cb()
      },
      trigger: 'blur'
    }
  ]
}

onMounted(() => {
  const up = getUserProfile()
  if (up) Object.assign(profile, up)
})

const onChangePassword = async () => {
  if (!pwdFormRef.value) return
  try {
    await pwdFormRef.value.validate()
  } catch {
    return
  }
  saving.value = true
  try {
    await changePassword('me', { password: pwdForm.password })
    useSuccessTip('密码已更新')
    pwdForm.password = ''
    pwdForm.confirm = ''
  } catch (e: any) {
    const raw = e?.msg || e?.message || e?.data || e?.response?.data
    const text = typeof raw === 'string' ? raw : ''
    if (text.includes('密码强度不够')) {
      useFailedTip('密码强度不够')
    } else {
      useFailedTip(text || '修改密码失败')
    }
  } finally {
    saving.value = false
  }
}

const goDashboard = () => {
  router.push('/dashboard')
}
</script>

<style scoped>
.profile-page {
  padding: 16px;
}
.profile-wrapper {
  width: 50%;
  min-width: 520px;
  max-width: 800px;
  margin: 0 auto;
}
.card-header h3 {
  margin: 0;
}
.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.profile-descriptions :deep(.el-descriptions__label) {
  width: 100px;
}

/* 新的信息列表样式，更现代的卡片风格 */
.info-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}
.info-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 14px;
  background: #fff;
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 10px;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.03);
}
.info-label {
  color: #606266;
}
.info-value {
  color: #303133;
  font-weight: 500;
}
</style>


