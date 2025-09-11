<template>
  <div id="theme-root">
    <router-view />
    <button class="theme-toggle" @click="toggleDarkMode">{{ isDark ? '🌙' : '☀️' }}</button>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'

const THEME_KEY = 'theme:dark'
const isDark = ref(false)

function applyTheme(dark: boolean) {
  const root = document.documentElement
  if (dark) {
    root.classList.add('dark')
  } else {
    root.classList.remove('dark')
  }
}

function toggleDarkMode() {
  isDark.value = !isDark.value
  localStorage.setItem(THEME_KEY, isDark.value ? '1' : '0')
  applyTheme(isDark.value)
}

onMounted(() => {
  const saved = localStorage.getItem(THEME_KEY)
  isDark.value = saved === '1'
  applyTheme(isDark.value)
})
</script>

<style>
/* 基础尺寸与重置 */
html,
body {
  min-height: 100vh;
  margin: 0;
}

#app {
  min-height: 100vh;
}

.el-container {
  min-height: 100vh;
}

/* 仅在桌面端应用较大的最小宽度约束，移动端放开以适配 */
@media (min-width: 1200px) {
  #app,
  html,
  body {
    min-width: 1200px;
  }
  .el-main {
    min-width: 800px;
  }
  .el-table {
    min-width: 600px;
  }
  .el-form {
    min-width: 400px;
  }
}

/* 主题变量（亮色为默认） */
:root {
  --bg-color: #ffffff;
  --text-color: #111111;
  --muted-text: #666666;
  --card-bg: #ffffff;
  --border-color: #e5e7eb;
  --link-color: #2563eb;
}

/* 暗色模式覆盖（通过 <html> 添加 .dark 类启用） */
html.dark {
  --bg-color: #0b0f14;
  --text-color: #e5e7eb;
  --muted-text: #a3a3a3;
  --card-bg: #121821;
  --border-color: #273445;
  --link-color: #60a5fa;
}

/* 将变量应用到全局元素 */
html, body, #app, #theme-root {
  background: var(--bg-color);
  color: var(--text-color);
}

a { color: var(--link-color); }

/* 简单适配一些常见容器/表格/表单背景与边框 */
.el-card,
.el-main,
.el-form,
.el-table,
.el-container {
  background: var(--card-bg);
  color: var(--text-color);
  border-color: var(--border-color);
}

/* 浮动切换按钮样式 */
.theme-toggle {
  position: fixed;
  right: 16px;
  bottom: 16px;
  z-index: 9999;
  width: 44px;
  height: 44px;
  border-radius: 9999px;
  border: 1px solid var(--border-color);
  background: var(--card-bg);
  color: var(--text-color);
  cursor: pointer;
  box-shadow: 0 2px 8px rgba(0,0,0,0.15);
}
.theme-toggle:hover {
  filter: brightness(1.05);
}
</style>
