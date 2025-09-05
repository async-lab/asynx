<script setup lang="ts">
</script>
<template>
    <router-view />
</template>

<style>
/* 全局最小尺寸设置 */
html, body {
  min-width: 1200px;
  min-height: 100vh;
  margin: 0;
}

#app {
  min-width: 1200px;
  min-height: 100vh;
}

/* 确保所有容器都有最小尺寸 */
.el-container {
  min-height: 100vh;
}

.el-main {
  min-width: 800px;
}

.el-card {
  min-height: 200px;
}

.el-table {
  min-width: 600px;
}

.el-form {
  min-width: 400px;
}
</style>
