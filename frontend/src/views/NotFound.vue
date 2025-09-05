<template>
  <div class="not-found">
    <div class="not-found-content">
      <div class="error-code">404</div>
      <h1 class="error-title">页面未找到</h1>
      <p class="error-description">
        抱歉，您访问的页面不存在或已被移除。
      </p>
      <div class="error-actions">
        <el-button type="primary" @click="goHome">
          返回首页
        </el-button>
        <el-button @click="goBack">
          返回上页
        </el-button>
      </div>
      
      <div class="suggestions">
        <h3>您可以尝试：</h3>
        <ul>
          <li>检查URL是否正确</li>
          <li>使用导航菜单浏览其他页面</li>
          <li>联系管理员获取帮助</li>
        </ul>
      </div>
    </div>
    
    <div class="not-found-illustration">
      <el-icon size="120" color="#e4e7ed"><QuestionFilled /></el-icon>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()

// 返回首页
const goHome = () => {
  router.push('/')
}

// 返回上一页
const goBack = () => {
  if (window.history.length > 1) {
    router.go(-1)
  } else {
    router.push('/')
  }
}
</script>

<style scoped>
.not-found {
  min-height: 100vh;
  min-width: 800px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  padding: 20px;
}

.not-found-content {
  text-align: center;
  max-width: 600px;
  min-width: 400px;
  margin-right: 60px;
  min-height: 400px;
}

.error-code {
  font-size: 120px;
  font-weight: bold;
  color: #409eff;
  line-height: 1;
  margin-bottom: 20px;
  text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.1);
}

.error-title {
  font-size: 32px;
  color: #303133;
  margin: 0 0 16px 0;
}

.error-description {
  font-size: 16px;
  color: #606266;
  margin: 0 0 32px 0;
  line-height: 1.6;
}

.error-actions {
  margin-bottom: 40px;
}

.error-actions .el-button {
  margin: 0 8px;
}

.suggestions {
  text-align: left;
  background: rgba(255, 255, 255, 0.8);
  padding: 20px;
  border-radius: 8px;
  border: 1px solid #e4e7ed;
}

.suggestions h3 {
  margin: 0 0 12px 0;
  color: #303133;
  font-size: 16px;
}

.suggestions ul {
  margin: 0;
  padding-left: 20px;
  color: #606266;
}

.suggestions li {
  margin-bottom: 8px;
  line-height: 1.5;
}

.not-found-illustration {
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0.6;
}

@media (max-width: 768px) {
  .not-found {
    flex-direction: column;
    text-align: center;
  }
  
  .not-found-content {
    margin-right: 0;
    margin-bottom: 40px;
  }
  
  .error-code {
    font-size: 80px;
  }
  
  .error-title {
    font-size: 24px;
  }
  
  .suggestions {
    text-align: center;
  }
}
</style> 