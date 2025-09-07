<script setup lang="ts">
</script>
<template>
    <router-view />
</template>

<style>
/* 基础尺寸与重置 */
html, body {
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
  #app, html, body { min-width: 1200px; }
  .el-main { min-width: 800px; }
  .el-table { min-width: 600px; }
  .el-form { min-width: 400px; }
}
</style>
